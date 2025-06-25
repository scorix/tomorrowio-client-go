package types

/* https://docs.tomorrow.io/reference/weather-forecast */

type WeatherForecast struct {
	Timelines WeatherTimeline `json:"timelines"`
	Location  WeatherLocation `json:"location"`
}

type WeatherTimeline struct {
	Minutely []WeatherMinutely `json:"minutely"`
	Hourly   []WeatherHourly   `json:"hourly"`
	Daily    []WeatherDaily    `json:"daily"`
}

type WeatherMinutely struct {
	Time   string                `json:"time"`
	Values WeatherMinutelyValues `json:"values"`
}

type WeatherMinutelyValues struct {
	Temperature            float64 `json:"temperature"`
	PrecipitationIntensity float64 `json:"precipitationIntensity"`
	PrecipitationType      int     `json:"precipitationType"`
	CloudCover             float64 `json:"cloudCover"`
	Humidity               float64 `json:"humidity"`
	WindSpeed              float64 `json:"windSpeed"`
	WindDirection          float64 `json:"windDirection"`
	WindGust               float64 `json:"windGust"`
	Visibility             float64 `json:"visibility"`
}

type WeatherHourly struct {
	Time   string              `json:"time"`
	Values WeatherHourlyValues `json:"values"`
}

type WeatherHourlyValues struct {
	Temperature              float64 `json:"temperature"`
	TemperatureApparent      float64 `json:"temperatureApparent"`
	PrecipitationIntensity   float64 `json:"precipitationIntensity"`
	PrecipitationType        int     `json:"precipitationType"`
	PrecipitationProbability float64 `json:"precipitationProbability"`
	CloudCover               float64 `json:"cloudCover"`
	CloudBase                float64 `json:"cloudBase"`
	CloudCeiling             float64 `json:"cloudCeiling"`
	Humidity                 float64 `json:"humidity"`
	Pressure                 float64 `json:"pressure"`
	WindSpeed                float64 `json:"windSpeed"`
	WindDirection            float64 `json:"windDirection"`
	WindGust                 float64 `json:"windGust"`
	Visibility               float64 `json:"visibility"`
	UVIndex                  float64 `json:"uvIndex"`
}

type WeatherDaily struct {
	Time   string             `json:"time"`
	Values WeatherDailyValues `json:"values"`
}

type WeatherDailyValues struct {
	TemperatureMax           float64 `json:"temperatureMax"`
	TemperatureMin           float64 `json:"temperatureMin"`
	TemperatureApparentMax   float64 `json:"temperatureApparentMax"`
	TemperatureApparentMin   float64 `json:"temperatureApparentMin"`
	PrecipitationIntensity   float64 `json:"precipitationIntensity"`
	PrecipitationType        int     `json:"precipitationType"`
	PrecipitationProbability float64 `json:"precipitationProbability"`
	CloudCoverAvg            float64 `json:"cloudCoverAvg"`
	HumidityAvg              float64 `json:"humidityAvg"`
	PressureSeaLevelAvg      float64 `json:"pressureSeaLevelAvg"`
	WindSpeedAvg             float64 `json:"windSpeedAvg"`
	WindDirectionAvg         float64 `json:"windDirectionAvg"`
	WindGustMax              float64 `json:"windGustMax"`
	VisibilityAvg            float64 `json:"visibilityAvg"`
	UVIndexMax               float64 `json:"uvIndexMax"`
	SunriseTime              string  `json:"sunriseTime"`
	SunsetTime               string  `json:"sunsetTime"`
}

type WeatherLocation struct {
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
	Name string  `json:"name"`
	Type string  `json:"type"`
}
