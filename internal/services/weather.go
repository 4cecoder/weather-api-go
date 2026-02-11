package services

import (
	"weather-api-go/internal/models"
	"weather-api-go/internal/repository"
)

// WeatherService handles weather-related business logic
type WeatherService struct {
	repo      *repository.WeatherRepository
	nwsClient *NWSAPIClient
}

// NewWeatherService creates a new weather service
func NewWeatherService(repo *repository.WeatherRepository, nwsClient *NWSAPIClient) *WeatherService {
	return &WeatherService{
		repo:      repo,
		nwsClient: nwsClient,
	}
}

// GetTemperatureCharacterization categorizes temperature as hot, cold, or moderate
func (s *WeatherService) GetTemperatureCharacterization(tempC float64) string {
	if tempC >= 30.0 {
		return "hot"
	} else if tempC <= 10.0 {
		return "cold"
	}
	return "moderate"
}

// GetWeather retrieves weather data with caching
func (s *WeatherService) GetWeather(lat, lon float64) (*models.WeatherResponse, error) {
	// Try to get from cache
	cachedWeather, err := s.repo.GetFromCache(lat, lon)
	if err == nil && s.repo.IsCacheFresh(cachedWeather) {
		return &models.WeatherResponse{
			Forecast:     cachedWeather.Forecast,
			Temperature:  s.GetTemperatureCharacterization(cachedWeather.TempC),
			TemperatureC: cachedWeather.TempC,
			TemperatureF: cachedWeather.TempF,
		}, nil
	}

	// Fetch fresh data from NWS
	weather, err := s.nwsClient.GetForecast(lat, lon)
	if err != nil {
		// Return stale cache if available
		if cachedWeather != nil {
			return &models.WeatherResponse{
				Forecast:     cachedWeather.Forecast,
				Temperature:  s.GetTemperatureCharacterization(cachedWeather.TempC),
				TemperatureC: cachedWeather.TempC,
				TemperatureF: cachedWeather.TempF,
			}, nil
		}
		return nil, err
	}

	// Save to cache (ignore errors, don't fail the request)
	s.repo.SaveToCache(weather)

	return &models.WeatherResponse{
		Forecast:     weather.Forecast,
		Temperature:  s.GetTemperatureCharacterization(weather.TempC),
		TemperatureC: weather.TempC,
		TemperatureF: weather.TempF,
	}, nil
}
