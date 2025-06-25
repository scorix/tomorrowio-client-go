package tommorrowio

import "fmt"

const (
	API_WEATHER_FORECAST_URL = "https://api.tomorrow.io/v4/weather/forecast"
)

func GetWeatherForecastURL(lat, lon float64, apiKey string) string {
	return fmt.Sprintf("%s?location=%f,%f&apikey=%s", API_WEATHER_FORECAST_URL, lat, lon, apiKey)
}
