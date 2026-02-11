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
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/mattn/go-sqlite3"
	"github.com/redis/go-redis/v9"

	"weather-api-go/internal/handlers"
	"weather-api-go/internal/repository"
	"weather-api-go/internal/services"
)

var ctx = context.Background()

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

	if _, err := client.Ping(ctx).Result(); err != nil {
		log.Printf("Redis connection failed: %v, falling back to SQLite only", err)
		return nil
	}

	log.Println("Connected to Redis")
	return client
}

func main() {
	app := fiber.New()

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	// Initialize database
	db, err := repository.InitDB("./weather_cache.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize Redis
	rdb := initRedis()
	if rdb != nil {
		defer rdb.Close()
	}

	// Initialize layered architecture
	weatherRepo := repository.NewWeatherRepository(db, rdb)
	nwsClient := services.NewNWSAPIClient()
	weatherService := services.NewWeatherService(weatherRepo, nwsClient)
	weatherHandler := handlers.NewWeatherHandler(weatherService)

	// Routes
	app.Get("/weather", weatherHandler.GetWeather)
	app.Get("/health", weatherHandler.GetHealth)

	// Serve frontend static files
	app.Static("/", "./dist/frontend")

	// SPA fallback
	app.Get("/*", func(c *fiber.Ctx) error {
		if len(c.Path()) >= 4 && c.Path()[:4] == "/api" {
			return c.Next()
		}
		return c.SendFile("./dist/frontend/index.html")
	})

	log.Println("Starting weather service on port 3000...")
	log.Println("Frontend available at: http://localhost:3000")

	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
