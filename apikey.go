package tommorrowio

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// ErrNoAPIKeyAvailable is returned when no API key is available due to rate limits
var ErrNoAPIKeyAvailable = errors.New("no API key available: all keys have exceeded rate limits")

type APIKeyPicker interface {
	GetAPIKey(ctx context.Context, nowFn func() time.Time) (string, error)
}

type apiKeyPicker struct {
	apiKeys        []string
	requestsPerDay int
	rpmLimit       int

	// Daily request tracking
	dailyRequests map[string]int
	dailyResetAt  map[string]time.Time

	// Rate limiters for RPM
	rateLimiters map[string]*rate.Limiter

	mu sync.RWMutex
}

// NewAPIKeyPicker creates a new APIKeyPicker with the given API keys and limits.
// The initialTime parameter is used to set the initial time for rate limiting and daily resets.
func NewAPIKeyPicker(apiKeys []string, maxRequestsPerDay, rpmLimit int, initialTime time.Time) APIKeyPicker {
	dailyResetAt := make(map[string]time.Time, len(apiKeys))
	rateLimiters := make(map[string]*rate.Limiter, len(apiKeys))
	dailyRequests := make(map[string]int, len(apiKeys))

	for _, key := range apiKeys {
		dailyResetAt[key] = initialTime.Add(24 * time.Hour)
		rateLimiters[key] = rate.NewLimiter(rate.Every(time.Minute/time.Duration(rpmLimit)), rpmLimit) // Burst equal to RPM limit
		dailyRequests[key] = 0
	}

	return &apiKeyPicker{
		apiKeys:        apiKeys,
		requestsPerDay: maxRequestsPerDay,
		rpmLimit:       rpmLimit,
		dailyRequests:  dailyRequests,
		dailyResetAt:   dailyResetAt,
		rateLimiters:   rateLimiters,
	}
}

func (p *apiKeyPicker) GetAPIKey(ctx context.Context, nowFn func() time.Time) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := nowFn()

	// Create a shuffled copy of API keys to randomize selection
	availableKeys := make([]string, len(p.apiKeys))
	copy(availableKeys, p.apiKeys)
	rand.Shuffle(len(availableKeys), func(i, j int) {
		availableKeys[i], availableKeys[j] = availableKeys[j], availableKeys[i]
	})

	// Try each key in random order
	for _, key := range availableKeys {
		// Reset daily counter if needed
		if now.After(p.dailyResetAt[key]) {
			p.dailyRequests[key] = 0
			// Set next reset time to be 24 hours from now
			p.dailyResetAt[key] = now.Add(24 * time.Hour)
		}

		// Check daily limit
		if p.dailyRequests[key] >= p.requestsPerDay {
			continue
		}

		// Check RPM limit using the provided time
		if !p.rateLimiters[key].AllowN(now, 1) {
			continue
		}

		// Key is available, increment counter and return
		p.dailyRequests[key]++

		return key, nil
	}

	return "", ErrNoAPIKeyAvailable
}
