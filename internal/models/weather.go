package models

import "time"

// WeatherResponse represents the API response for weather data
type WeatherResponse struct {
	Forecast     string  `json:"forecast" example:"Partly Cloudy"`
	Temperature  string  `json:"temperature" example:"moderate"`
	TemperatureC float64 `json:"temperature_c" example:"22.5"`
	TemperatureF float64 `json:"temperature_f" example:"72.5"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid latitude parameter"`
	Details string `json:"details,omitempty" example:"Latitude must be between -90 and 90"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string `json:"status" example:"healthy"`
	Timestamp string `json:"timestamp" example:"2024-01-15T10:30:00Z"`
}

// WeatherCache represents cached weather data
type WeatherCache struct {
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Forecast  string    `json:"forecast"`
	TempC     float64   `json:"temp_c"`
	TempF     float64   `json:"temp_f"`
	Timestamp time.Time `json:"timestamp"`
}

// Coordinates represents geographic coordinates
type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// NWSPointsResponse represents the NWS API points endpoint response
type NWSPointsResponse struct {
	Properties struct {
		Forecast string `json:"forecast"`
	} `json:"properties"`
}

// NWSForecastResponse represents the NWS API forecast endpoint response
type NWSForecastResponse struct {
	Properties struct {
		Periods []struct {
			ShortForecast   string  `json:"shortForecast"`
			Temperature     float64 `json:"temperature"`
			TemperatureUnit string  `json:"temperatureUnit"`
		} `json:"periods"`
	} `json:"properties"`
}
