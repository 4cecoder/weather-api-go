package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetTemperatureCharacterization(t *testing.T) {
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
			result := getTemperatureCharacterization(tt.tempC)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWeatherResponse(t *testing.T) {
	response := WeatherResponse{
		Forecast:     "Sunny",
		Temperature:  "hot",
		TemperatureC: 32.5,
	}

	data, err := json.Marshal(response)
	assert.NoError(t, err)

	var decoded WeatherResponse
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, response.Forecast, decoded.Forecast)
	assert.Equal(t, response.Temperature, decoded.Temperature)
	assert.Equal(t, response.TemperatureC, decoded.TemperatureC)
}

func TestWeatherCache(t *testing.T) {
	cache := WeatherCache{
		Latitude:  40.7128,
		Longitude: -74.0060,
		Forecast:  "Partly Cloudy",
		TempC:     22.5,
		TempF:     72.5,
		Timestamp: time.Now(),
	}

	assert.Equal(t, 40.7128, cache.Latitude)
	assert.Equal(t, -74.0060, cache.Longitude)
	assert.Equal(t, "Partly Cloudy", cache.Forecast)
	assert.Equal(t, 22.5, cache.TempC)
	assert.Equal(t, 72.5, cache.TempF)
}

func TestHealthHandler(t *testing.T) {
	app := fiber.New()
	app.Get("/health", healthHandler)

	t.Run("returns healthy status", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/health", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var result HealthResponse
		err = json.Unmarshal(body, &result)
		assert.NoError(t, err)
		assert.Equal(t, "healthy", result.Status)
		assert.NotEmpty(t, result.Timestamp)
	})
}

func TestGetWeatherHandler_Validation(t *testing.T) {
	// Setup test database
	testDBPath := "./test_weather_cache.db"
	os.Remove(testDBPath)
	defer os.Remove(testDBPath)

	var err error
	db, err = initDBForTest(testDBPath)
	if err != nil {
		t.Fatalf("Failed to init test DB: %v", err)
	}
	defer db.Close()

	app := fiber.New()
	app.Get("/weather", getWeatherHandler)

	tests := []struct {
		name       string
		url        string
		wantStatus int
		wantError  string
	}{
		{
			name:       "missing latitude",
			url:        "/weather?lon=-74.0060",
			wantStatus: 400,
			wantError:  "Missing latitude parameter",
		},
		{
			name:       "missing longitude",
			url:        "/weather?lat=40.7128",
			wantStatus: 400,
			wantError:  "Missing longitude parameter",
		},
		{
			name:       "invalid latitude",
			url:        "/weather?lat=invalid&lon=-74.0060",
			wantStatus: 400,
			wantError:  "Invalid latitude parameter",
		},
		{
			name:       "invalid longitude",
			url:        "/weather?lat=40.7128&lon=invalid",
			wantStatus: 400,
			wantError:  "Invalid longitude parameter",
		},
		{
			name:       "latitude out of range",
			url:        "/weather?lat=100&lon=-74.0060",
			wantStatus: 400,
			wantError:  "Invalid coordinates",
		},
		{
			name:       "longitude out of range",
			url:        "/weather?lat=40.7128&lon=200",
			wantStatus: 400,
			wantError:  "Invalid coordinates",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			body, _ := io.ReadAll(resp.Body)
			var result ErrorResponse
			err = json.Unmarshal(body, &result)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantError, result.Error)
		})
	}
}

func initDBForTest(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS weather_cache (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			latitude REAL NOT NULL,
			longitude REAL NOT NULL,
			forecast TEXT,
			temp_c REAL,
			temp_f REAL,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return db, err
}

func BenchmarkGetTemperatureCharacterization(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getTemperatureCharacterization(25.0)
	}
}
