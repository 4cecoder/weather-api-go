package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"weather-api-go/internal/models"
)

var ctx = context.Background()

// WeatherRepository handles weather data persistence
type WeatherRepository struct {
	db  *sql.DB
	rdb *redis.Client
}

// NewWeatherRepository creates a new weather repository
func NewWeatherRepository(db *sql.DB, rdb *redis.Client) *WeatherRepository {
	return &WeatherRepository{
		db:  db,
		rdb: rdb,
	}
}

// GetFromCache retrieves weather data from cache (Redis first, then SQLite)
func (r *WeatherRepository) GetFromCache(lat, lon float64) (*models.WeatherCache, error) {
	// Try Redis first
	if r.rdb != nil {
		key := fmt.Sprintf("weather:%.6f:%.6f", lat, lon)
		data, err := r.rdb.Get(ctx, key).Result()
		if err == nil {
			var cache models.WeatherCache
			if err := json.Unmarshal([]byte(data), &cache); err == nil {
				return &cache, nil
			}
		}
	}

	// Fallback to SQLite
	var cache models.WeatherCache
	err := r.db.QueryRow(
		"SELECT forecast, temp_c, temp_f, timestamp FROM weather_cache WHERE latitude = ? AND longitude = ? ORDER BY timestamp DESC LIMIT 1",
		lat, lon,
	).Scan(&cache.Forecast, &cache.TempC, &cache.TempF, &cache.Timestamp)

	if err != nil {
		return nil, err
	}

	cache.Latitude = lat
	cache.Longitude = lon
	return &cache, nil
}

// SaveToCache saves weather data to cache (Redis and SQLite)
func (r *WeatherRepository) SaveToCache(weather *models.WeatherCache) error {
	// Cache in Redis
	if r.rdb != nil {
		key := fmt.Sprintf("weather:%.6f:%.6f", weather.Latitude, weather.Longitude)
		data, err := json.Marshal(weather)
		if err == nil {
			r.rdb.Set(ctx, key, data, 1*time.Hour)
		}
	}

	// Also cache in SQLite for persistence
	_, err := r.db.Exec(
		"INSERT INTO weather_cache (latitude, longitude, forecast, temp_c, temp_f) VALUES (?, ?, ?, ?, ?)",
		weather.Latitude, weather.Longitude, weather.Forecast, weather.TempC, weather.TempF,
	)
	return err
}

// IsCacheFresh checks if cached data is still fresh (within 1 hour)
func (r *WeatherRepository) IsCacheFresh(cache *models.WeatherCache) bool {
	return time.Since(cache.Timestamp) < time.Hour
}

// InitDB initializes the database schema
func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
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
