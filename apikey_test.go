package tomorrowio_test

import (
	"context"
	"testing"
	"time"

	"github.com/scorix/tomorrowio-client-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// timeNowFn returns a function that returns a fixed time plus an optional duration
func timeNowFn(base time.Time, offset ...time.Duration) func() time.Time {
	if len(offset) > 0 {
		return func() time.Time { return base.Add(offset[0]) }
	}
	return func() time.Time { return base }
}

func TestAPIKeyPicker(t *testing.T) {
	baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	t.Run("daily limit", func(t *testing.T) {
		picker := tomorrowio.NewAPIKeyPicker([]string{"key1"}, 2, 60, baseTime) // 2 requests per day, 60 rpm
		ctx := context.Background()

		// First request should succeed
		key1, err := picker.GetAPIKey(ctx, timeNowFn(baseTime))
		require.NoError(t, err)
		assert.Equal(t, "key1", key1)

		// Second request should succeed
		key2, err := picker.GetAPIKey(ctx, timeNowFn(baseTime))
		require.NoError(t, err)
		assert.Equal(t, "key1", key2)

		// Third request should fail (daily limit exceeded)
		_, err = picker.GetAPIKey(ctx, timeNowFn(baseTime))
		assert.ErrorIs(t, err, tomorrowio.ErrNoAPIKeyAvailable)
	})

	t.Run("rpm limit", func(t *testing.T) {
		picker := tomorrowio.NewAPIKeyPicker([]string{"key1"}, 1000, 2, baseTime) // 1000 requests per day, 2 rpm
		ctx := context.Background()

		// First request should succeed
		key1, err := picker.GetAPIKey(ctx, timeNowFn(baseTime))
		require.NoError(t, err)
		assert.Equal(t, "key1", key1)

		// Second request should succeed
		key2, err := picker.GetAPIKey(ctx, timeNowFn(baseTime, time.Second))
		require.NoError(t, err)
		assert.Equal(t, "key1", key2)

		// Third request should fail (rpm limit exceeded)
		_, err = picker.GetAPIKey(ctx, timeNowFn(baseTime, 2*time.Second))
		assert.ErrorIs(t, err, tomorrowio.ErrNoAPIKeyAvailable)

		// Request should succeed again after rate limit reset
		key3, err := picker.GetAPIKey(ctx, timeNowFn(baseTime, 31*time.Second))
		require.NoError(t, err)
		assert.Equal(t, "key1", key3)
	})

	t.Run("multiple keys", func(t *testing.T) {
		picker := tomorrowio.NewAPIKeyPicker([]string{"key1", "key2"}, 1, 60, baseTime) // 1 request per day per key, 60 rpm
		ctx := context.Background()

		// First request should succeed with either key
		key1, err := picker.GetAPIKey(ctx, timeNowFn(baseTime))
		require.NoError(t, err)
		assert.Contains(t, []string{"key1", "key2"}, key1)

		// Second request should succeed with the other key
		key2, err := picker.GetAPIKey(ctx, timeNowFn(baseTime))
		require.NoError(t, err)
		assert.Contains(t, []string{"key1", "key2"}, key2)
		assert.NotEqual(t, key1, key2, "should use different keys")

		// Third request should fail (all keys have reached daily limit)
		_, err = picker.GetAPIKey(ctx, timeNowFn(baseTime))
		assert.ErrorIs(t, err, tomorrowio.ErrNoAPIKeyAvailable)
	})

	t.Run("daily reset", func(t *testing.T) {
		picker := tomorrowio.NewAPIKeyPicker([]string{"key1"}, 1, 60, baseTime) // 1 request per day, 60 rpm
		ctx := context.Background()

		// First request should succeed
		key1, err := picker.GetAPIKey(ctx, timeNowFn(baseTime))
		require.NoError(t, err)
		assert.Equal(t, "key1", key1)

		// Second request should fail (daily limit exceeded)
		_, err = picker.GetAPIKey(ctx, timeNowFn(baseTime))
		assert.ErrorIs(t, err, tomorrowio.ErrNoAPIKeyAvailable)

		// Move time forward by 25 hours to ensure we're in the next day
		nextDay := baseTime.Add(25 * time.Hour)

		// Request should succeed again after daily reset
		key2, err := picker.GetAPIKey(ctx, timeNowFn(nextDay))
		require.NoError(t, err)
		assert.Equal(t, "key1", key2)
	})
}
