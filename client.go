package tommorrowio

import (
	"context"
	"fmt"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/scorix/tomorrowio-client-go/api"
	"github.com/scorix/tomorrowio-client-go/types"
)

const (
	BASE_URL = "https://api.tomorrow.io"
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
	apiKey, err := c.apiKeyPicker.GetAPIKey(ctx, time.Now)
	if err != nil {
		return nil, fmt.Errorf("failed to get API key: %w", err)
	}

	return api.GetWeatherForecast(ctx, c.forecastCache, BASE_URL, apiKey, lat, lon)
}
