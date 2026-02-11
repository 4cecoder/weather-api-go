package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"weather-api-go/internal/models"
	"weather-api-go/internal/services"
)

// WeatherHandler handles weather-related HTTP requests
type WeatherHandler struct {
	service *services.WeatherService
}

// NewWeatherHandler creates a new weather handler
func NewWeatherHandler(service *services.WeatherService) *WeatherHandler {
	return &WeatherHandler{service: service}
}

// GetWeather handles GET /weather requests
// @Summary Get weather forecast
// @Description Returns the short forecast and temperature characterization for the specified latitude and longitude
// @Tags weather
// @Accept json
// @Produce json
// @Param lat query number true "Latitude coordinate (-90 to 90)" example(40.7128)
// @Param lon query number true "Longitude coordinate (-180 to 180)" example(-74.0060)
// @Success 200 {object} models.WeatherResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /weather [get]
func (h *WeatherHandler) GetWeather(c *fiber.Ctx) error {
	latStr := c.Query("lat")
	if latStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Missing latitude parameter",
			Details: "Latitude is required (e.g., lat=40.7128)",
		})
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid latitude parameter",
			Details: "Latitude must be a valid float number",
		})
	}

	lonStr := c.Query("lon")
	if lonStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Missing longitude parameter",
			Details: "Longitude is required (e.g., lon=-74.0060)",
		})
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid longitude parameter",
			Details: "Longitude must be a valid float number",
		})
	}

	if lat < -90 || lat > 90 || lon < -180 || lon > 180 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid coordinates",
			Details: "Latitude must be between -90 and 90, Longitude between -180 and 180",
		})
	}

	weather, err := h.service.GetWeather(lat, lon)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Failed to get weather data",
			Details: err.Error(),
		})
	}

	return c.JSON(weather)
}

// GetHealth handles GET /health requests
// @Summary Health check
// @Description Check if the weather service is running
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} models.HealthResponse
// @Router /health [get]
func (h *WeatherHandler) GetHealth(c *fiber.Ctx) error {
	return c.JSON(models.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
	})
}
