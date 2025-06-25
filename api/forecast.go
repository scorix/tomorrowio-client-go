package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	lru "github.com/hashicorp/golang-lru"
	"github.com/scorix/tomorrowio-client-go/cache"
	"github.com/scorix/tomorrowio-client-go/types"
)

const WEATHER_FORECAST_V4_PATH = "/v4/weather/forecast"

func GetWeatherForecastURL(baseURL string, lat, lon float64, apiKey string) string {
	u, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}

	u.Path = WEATHER_FORECAST_V4_PATH
	v := u.Query()
	v.Set("location", fmt.Sprintf("%f,%f", lat, lon))
	v.Set("apikey", apiKey)
	u.RawQuery = v.Encode()

	return u.String()
}

func GetWeatherForecast(ctx context.Context, lruCache *lru.Cache, baseURL string, apiKey string, lat, lon float64) (*types.WeatherForecast, error) {
	key := fmt.Sprintf("%f,%f", lat, lon)
	fn := func() (*types.WeatherForecast, error) {
		return getWeatherForecast(ctx, baseURL, apiKey, lat, lon)
	}

	return cache.WithLRU(lruCache, key, fn)
}

func getWeatherForecast(ctx context.Context, baseURL string, apiKey string, lat, lon float64) (*types.WeatherForecast, error) {
	url := GetWeatherForecastURL(baseURL, lat, lon, apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var weather types.WeatherForecast
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &weather, nil
}
