package tommorrowio

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/scorix/tomorrowio-client-go/types"
)

type Client interface {
	GetWeatherForecast(ctx context.Context, lat, lon float64) (*types.WeatherForecast, error)
}

type client struct {
	apiKeyPicker  APIKeyPicker
	forecastCache *lru.Cache
}

func NewClient(apiKeys []string, cacheSize int, maxRequestsPerDay int, rpmLimit int) (*client, error) {
	forecastCache, err := lru.New(cacheSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create cache: %w", err)
	}
	apiKeyPicker := NewAPIKeyPicker(apiKeys, maxRequestsPerDay, rpmLimit, time.Now())

	return &client{apiKeyPicker: apiKeyPicker, forecastCache: forecastCache}, nil
}

func (c *client) GetWeatherForecast(ctx context.Context, lat, lon float64) (*types.WeatherForecast, error) {
	key := fmt.Sprintf("%f,%f", lat, lon)
	fn := func() (*types.WeatherForecast, error) {
		return c.getWeatherForecast(ctx, lat, lon)
	}

	return withCache(c.forecastCache, key, fn)
}

func (c *client) getWeatherForecast(ctx context.Context, lat, lon float64) (*types.WeatherForecast, error) {
	apiKey, err := c.apiKeyPicker.GetAPIKey(ctx, time.Now)
	if err != nil {
		return nil, fmt.Errorf("failed to get API key: %w", err)
	}
	url := GetWeatherForecastURL(lat, lon, apiKey)

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
