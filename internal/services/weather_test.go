package services

import (
	"testing"
)

func TestGetTemperatureCharacterization(t *testing.T) {
	service := &WeatherService{}

	tests := []struct {
		name     string
		tempC    float64
		expected string
	}{
		{"Hot temperature", 35.0, "hot"},
		{"Hot boundary", 30.0, "hot"},
		{"Cold temperature", 5.0, "cold"},
		{"Cold boundary", 10.0, "cold"},
		{"Moderate low", 15.0, "moderate"},
		{"Moderate high", 25.0, "moderate"},
		{"Moderate middle", 20.0, "moderate"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.GetTemperatureCharacterization(tt.tempC)
			if result != tt.expected {
				t.Errorf("GetTemperatureCharacterization(%f) = %s; want %s", tt.tempC, result, tt.expected)
			}
		})
	}
}
