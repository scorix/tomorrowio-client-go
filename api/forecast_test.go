package api_test

import (
	"testing"

	"github.com/scorix/tomorrowio-client-go/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetWeatherForecastURL(t *testing.T) {
	url := api.GetWeatherForecastURL("https://api.tomorrow.io", 43.15616989135742, -75.8449935913086, "key1")
	assert.Equal(t, "https://api.tomorrow.io/v4/weather/forecast?location=43.156170,-75.844994&apikey=key1", url)
}

func TestGetWeatherForecast(t *testing.T) {
	mockServer := NewMockServer(t)

	forecast, err := api.GetWeatherForecast(t.Context(), nil, mockServer.URL(), "key1", 43.15616989135742, -75.8449935913086)
	require.NoError(t, err)
	require.NotNil(t, forecast)

	assert.Equal(t, 43.15616989135742, forecast.Location.Lat)
	assert.Equal(t, -75.8449935913086, forecast.Location.Lon)
	assert.Equal(t, "New York, United States", forecast.Location.Name)
	assert.Equal(t, "administrative", forecast.Location.Type)
	assert.Equal(t, 60, len(forecast.Timelines.Minutely))
	assert.Equal(t, 120, len(forecast.Timelines.Hourly))
	assert.Equal(t, 7, len(forecast.Timelines.Daily))
}
