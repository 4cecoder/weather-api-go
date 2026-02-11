# Weather API Service

A production-ready Go service that provides weather forecasts based on latitude and longitude coordinates using the National Weather Service API. The service features dual-layer caching (Redis + SQLite), OpenAPI documentation, and Docker support.

## Features

- **Weather Forecast**: Returns today's short forecast and temperature characterization
- **Dual Caching**: Redis (fast, in-memory) + SQLite (persistent, fallback)
- **Temperature Classification**: Hot (≥30°C), Cold (≤10°C), or Moderate
- **API Documentation**: Auto-generated OpenAPI/Swagger UI
- **Docker Support**: Multi-stage Docker build with docker-compose
- **Health Monitoring**: Health check endpoint
- **Error Handling**: Comprehensive error responses

## Quick Start

### Using Docker (Recommended)

```bash
# Clone the repository
git clone <repo-url>
cd weather-api-go

# Start services with docker-compose
docker-compose up -d

# Access the API
curl "http://localhost:3000/weather?lat=40.7128&lon=-74.0060"

# View API documentation
open http://localhost:3000/swagger/index.html
```

### Local Development

```bash
# Install dependencies
go mod download

# Run the application
go run main.go

# The server will start on http://localhost:3000
```

## API Endpoints

### GET /weather
Returns the weather forecast for given coordinates.

**Parameters:**
- `lat` (required): Latitude (-90 to 90)
- `lon` (required): Longitude (-180 to 180)

**Example Request:**
```bash
curl "http://localhost:3000/weather?lat=40.7128&lon=-74.0060"
```

**Example Response:**
```json
{
  "forecast": "Partly Cloudy",
  "temperature": "moderate",
  "temperature_c": 22.5
}
```

### GET /health
Health check endpoint.

**Example Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### GET /swagger/*
Interactive API documentation (Swagger UI).

## Architecture

### Caching Strategy

1. **Redis** (Primary Cache): In-memory caching for sub-millisecond response times
2. **SQLite** (Fallback Cache): Persistent storage for data durability

Cache TTL: 1 hour

### Temperature Classification

- **Hot**: ≥ 30°C (86°F)
- **Cold**: ≤ 10°C (50°F)
- **Moderate**: Between 10°C and 30°C

## Project Structure

```
weather-api-go/
├── main.go              # Application entry point
├── main_test.go         # Unit tests
├── Dockerfile           # Docker build configuration
├── docker-compose.yml   # Docker Compose orchestration
├── go.mod               # Go module dependencies
├── go.sum               # Go module checksums
├── .env.example         # Environment variables template
├── README.md            # This file
└── weather_cache.db     # SQLite database (auto-created)
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | 3000 |
| `REDIS_URL` | Redis connection URL | localhost:6379 |
| `DATABASE_URL` | SQLite database path | ./weather_cache.db |

## Testing

```bash
# Run all tests
go test -v

# Run tests with coverage
go test -v -cover

# Run benchmarks
go test -bench=.
```

## Docker Commands

```bash
# Build and start
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Rebuild after changes
docker-compose up -d --build
```

## Development Notes

### Shortcuts & Trade-offs

1. **Caching Strategy**: Uses a simple 1-hour TTL. In production, consider:
   - Stale-while-revalidate pattern
   - Cache warming strategies
   - Different TTLs for different data freshness requirements

2. **Error Handling**: Falls back to cached data on NWS API failures. For production:
   - Consider circuit breaker patterns
   - Add structured logging and monitoring
   - Implement rate limiting

3. **Database**: SQLite is used for simplicity. For production:
   - Consider PostgreSQL or MySQL for better concurrency
   - Add connection pooling
   - Implement database migrations

4. **Testing**: Currently has unit tests. Consider adding:
   - Integration tests with NWS API mocking
   - Load tests
   - End-to-end tests

## API Documentation

Once the server is running, access the interactive documentation at:
```
http://localhost:3000/swagger/index.html
```

The documentation includes:
- All available endpoints
- Request/response schemas
- Example requests
- Interactive "Try it out" feature

## Data Source

This service uses the [National Weather Service API](https://api.weather.gov/)

## License

MIT License - See LICENSE file for details