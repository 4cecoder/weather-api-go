package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"weather-api-go/internal/models"
)

// NWSAPIClient handles communication with National Weather Service API
type NWSAPIClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewNWSAPIClient creates a new NWS API client
func NewNWSAPIClient() *NWSAPIClient {
	return &NWSAPIClient{
		baseURL: "https://api.weather.gov",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetForecast fetches weather forecast for given coordinates
func (c *NWSAPIClient) GetForecast(lat, lon float64) (*models.WeatherCache, error) {
	// Step 1: Get forecast URL from points endpoint
	pointsURL := fmt.Sprintf("%s/points/%f,%f", c.baseURL, lat, lon)

	pointsResp, err := c.httpClient.Get(pointsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch points data: %w", err)
	}
	defer pointsResp.Body.Close()

	if pointsResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("NWS points API returned status: %d", pointsResp.StatusCode)
	}

	var pointsData models.NWSPointsResponse
	if err := json.NewDecoder(pointsResp.Body).Decode(&pointsData); err != nil {
		return nil, fmt.Errorf("failed to decode points response: %w", err)
	}

	if pointsData.Properties.Forecast == "" {
		return nil, fmt.Errorf("no forecast URL found in points response")
	}

	// Step 2: Get actual forecast data
	forecastResp, err := c.httpClient.Get(pointsData.Properties.Forecast)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch forecast data: %w", err)
	}
	defer forecastResp.Body.Close()

	if forecastResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("NWS forecast API returned status: %d", forecastResp.StatusCode)
	}

	var forecastData models.NWSForecastResponse
	if err := json.NewDecoder(forecastResp.Body).Decode(&forecastData); err != nil {
		return nil, fmt.Errorf("failed to decode forecast response: %w", err)
	}

	if len(forecastData.Properties.Periods) == 0 {
		return nil, fmt.Errorf("no forecast periods found")
	}

	// Parse first period (today's forecast)
	today := forecastData.Properties.Periods[0]

	// Convert temperature to Celsius if needed
	var tempC float64
	if today.TemperatureUnit == "F" {
		tempC = (today.Temperature - 32) * 5 / 9
	} else {
		tempC = today.Temperature
	}

	return &models.WeatherCache{
		Latitude:  lat,
		Longitude: lon,
		Forecast:  today.ShortForecast,
		TempC:     tempC,
		TempF:     today.Temperature,
		Timestamp: time.Now(),
	}, nil
}
