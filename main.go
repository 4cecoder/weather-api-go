// @title Weather API
// @version 1.0
// @description A weather service that provides forecasted weather based on latitude and longitude coordinates using the National Weather Service API.

// @host localhost:3000
// @BasePath /

// @contact.name API Support
// @contact.email support@weather-api.example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/mattn/go-sqlite3"
	"github.com/redis/go-redis/v9"

	swagger "github.com/arsmn/fiber-swagger/v2"
)

var db *sql.DB
var rdb *redis.Client
var ctx = context.Background()

type WeatherResponse struct {
	Forecast     string  `json:"forecast" example:"Partly Cloudy"`
	Temperature  string  `json:"temperature" example:"moderate"`
	TemperatureC float64 `json:"temperature_c" example:"22.5"`
}

func initRedis() *redis.Client {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: "",
		DB:       0,
	})

	// Test connection
	if _, err := client.Ping(ctx).Result(); err != nil {
		log.Printf("Redis connection failed: %v, falling back to SQLite only", err)
		return nil
	}

	log.Println("Connected to Redis")
	return client
}

type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid latitude parameter"`
	Details string `json:"details,omitempty" example:"Latitude must be between -90 and 90"`
}

type HealthResponse struct {
	Status    string `json:"status" example:"healthy"`
	Timestamp string `json:"timestamp" example:"2024-01-15T10:30:00Z"`
}

type WeatherCache struct {
	Latitude  float64
	Longitude float64
	Forecast  string
	TempC     float64
	TempF     float64
	Timestamp time.Time
}

func initDB() error {
	var err error
	db, err = sql.Open("sqlite3", "./weather_cache.db")
	if err != nil {
		return err
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
	return err
}

func getTemperatureCharacterization(tempC float64) string {
	if tempC >= 30.0 {
		return "hot"
	} else if tempC <= 10.0 {
		return "cold"
	}
	return "moderate"
}

func getWeatherFromNWS(lat float64, lon float64) (*WeatherCache, error) {
	url := fmt.Sprintf("https://api.weather.gov/points/%f,%f", lat, lon)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("NWS API returned status: %d", resp.StatusCode)
	}

	var nwsResponse struct {
		Properties struct {
			Forecast string `json:"forecast"`
		} `json:"properties"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&nwsResponse); err != nil {
		return nil, err
	}

	if nwsResponse.Properties.Forecast == "" {
		return nil, fmt.Errorf("no forecast URL found")
	}

	forecastResp, err := http.Get(nwsResponse.Properties.Forecast)
	if err != nil {
		return nil, err
	}
	defer forecastResp.Body.Close()

	if forecastResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("NWS forecast API returned status: %d", forecastResp.StatusCode)
	}

	var forecastResponse struct {
		Properties struct {
			Periods []struct {
				ShortForecast   string  `json:"shortForecast"`
				Temperature     float64 `json:"temperature"`
				TemperatureUnit string  `json:"temperatureUnit"`
			} `json:"periods"`
		} `json:"properties"`
	}

	if err := json.NewDecoder(forecastResp.Body).Decode(&forecastResponse); err != nil {
		return nil, err
	}

	if len(forecastResponse.Properties.Periods) == 0 {
		return nil, fmt.Errorf("no forecast periods found")
	}

	today := forecastResponse.Properties.Periods[0]

	var tempC float64
	if today.TemperatureUnit == "F" {
		tempC = (today.Temperature - 32) * 5 / 9
	} else {
		tempC = today.Temperature
	}

	return &WeatherCache{
		Latitude:  lat,
		Longitude: lon,
		Forecast:  today.ShortForecast,
		TempC:     tempC,
		TempF:     today.Temperature,
		Timestamp: time.Now(),
	}, nil
}

func getCachedWeather(lat float64, lon float64) (*WeatherCache, error) {
	// Try Redis first
	if rdb != nil {
		key := fmt.Sprintf("weather:%.6f:%.6f", lat, lon)
		data, err := rdb.Get(ctx, key).Result()
		if err == nil {
			var cache WeatherCache
			if err := json.Unmarshal([]byte(data), &cache); err == nil {
				log.Println("Cache hit from Redis")
				return &cache, nil
			}
		}
	}

	// Fallback to SQLite
	var cache WeatherCache
	err := db.QueryRow("SELECT forecast, temp_c, temp_f, timestamp FROM weather_cache WHERE latitude = ? AND longitude = ? ORDER BY timestamp DESC LIMIT 1",
		lat, lon).Scan(&cache.Forecast, &cache.TempC, &cache.TempF, &cache.Timestamp)
	if err != nil {
		return nil, err
	}
	cache.Latitude = lat
	cache.Longitude = lon

	// Populate Redis with SQLite data
	if rdb != nil {
		cacheWeather(&cache)
	}

	return &cache, nil
}

func cacheWeather(weather *WeatherCache) error {
	// Cache in Redis
	if rdb != nil {
		key := fmt.Sprintf("weather:%.6f:%.6f", weather.Latitude, weather.Longitude)
		data, err := json.Marshal(weather)
		if err == nil {
			rdb.Set(ctx, key, data, 1*time.Hour)
		}
	}

	// Also cache in SQLite for persistence
	_, err := db.Exec("INSERT INTO weather_cache (latitude, longitude, forecast, temp_c, temp_f) VALUES (?, ?, ?, ?, ?)",
		weather.Latitude, weather.Longitude, weather.Forecast, weather.TempC, weather.TempF)
	return err
}

// @Summary Get weather forecast
// @Description Returns the short forecast and temperature characterization for the specified latitude and longitude
// @Tags weather
// @Accept json
// @Produce json
// @Param lat query number true "Latitude coordinate (-90 to 90)" example(40.7128)
// @Param lon query number true "Longitude coordinate (-180 to 180)" example(-74.0060)
// @Success 200 {object} WeatherResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /weather [get]
func getWeatherHandler(c *fiber.Ctx) error {
	latStr := c.Query("lat")
	if latStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Missing latitude parameter",
			Details: "Latitude is required (e.g., lat=40.7128)",
		})
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid latitude parameter",
			Details: "Latitude must be a valid float number",
		})
	}

	lonStr := c.Query("lon")
	if lonStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Missing longitude parameter",
			Details: "Longitude is required (e.g., lon=-74.0060)",
		})
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid longitude parameter",
			Details: "Longitude must be a valid float number",
		})
	}

	if lat < -90 || lat > 90 || lon < -180 || lon > 180 {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid coordinates",
			Details: "Latitude must be between -90 and 90, Longitude between -180 and 180",
		})
	}

	cachedWeather, err := getCachedWeather(lat, lon)
	if err == nil {
		if time.Since(cachedWeather.Timestamp) < time.Hour {
			return c.JSON(WeatherResponse{
				Forecast:     cachedWeather.Forecast,
				Temperature:  getTemperatureCharacterization(cachedWeather.TempC),
				TemperatureC: cachedWeather.TempC,
			})
		}
	}

	weather, err := getWeatherFromNWS(lat, lon)
	if err != nil {
		if cachedWeather != nil {
			return c.JSON(WeatherResponse{
				Forecast:     cachedWeather.Forecast,
				Temperature:  getTemperatureCharacterization(cachedWeather.TempC),
				TemperatureC: cachedWeather.TempC,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to get weather data",
			Details: err.Error(),
		})
	}

	if err := cacheWeather(weather); err != nil {
		log.Printf("Failed to cache weather data: %v", err)
	}

	return c.JSON(WeatherResponse{
		Forecast:     weather.Forecast,
		Temperature:  getTemperatureCharacterization(weather.TempC),
		TemperatureC: weather.TempC,
	})
}

// @Summary Health check
// @Description Check if the weather service is running
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func healthHandler(c *fiber.Ctx) error {
	return c.JSON(HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

func main() {
	app := fiber.New()

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	if err := initDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize Redis
	rdb = initRedis()
	if rdb != nil {
		defer rdb.Close()
	}

	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:          "/swagger/doc.json",
		DeepLinking:  true,
		DocExpansion: "list",
	}))

	app.Get("/weather", getWeatherHandler)
	app.Get("/health", healthHandler)

	log.Println("Starting weather service on port 3000...")
	log.Println("API Documentation available at: http://localhost:3000/swagger/index.html")

	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
