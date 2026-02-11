# Weather API Service

<p align="center">
  <a href="https://github.com/4cecoder">
    <img src="https://img.shields.io/badge/GitHub-4cecoder-667eea?style=for-the-badge&logo=github" alt="GitHub">
  </a>
  <a href="https://github.com/4cecoder/weather-api-go">
    <img src="https://img.shields.io/badge/Repo-weather--api--go-764ba2?style=for-the-badge" alt="Repository">
  </a>
</p>

A production-ready weather service with a modern layered architecture, featuring an interactive web interface, futuristic API documentation, and dual-layer caching. Built with Go, React + TanStack, and OpenStreetMap.

![Weather API Screenshot](https://via.placeholder.com/800x400/0f0f23/667eea?text=Weather+API+Interface)

## ✨ Features

### Backend (Go)
- **🏗️ Layered Architecture**: Clean separation with handlers → services → repository → models
- **🔄 Dual-Layer Caching**: Redis (fast in-memory) + SQLite (persistent storage)
- **🌡️ Temperature Conversion**: API returns both Celsius and Fahrenheit
- **🗺️ NWS Integration**: Uses National Weather Service API for accurate forecasts
- **📊 Health Monitoring**: Built-in health check endpoint
- **🔒 Error Handling**: Comprehensive validation and error responses

### Frontend (React + TanStack)
- **🎨 Swiss Luxury Design**: Premium, minimalist aesthetic inspired by high-end spas
- **🗺️ Interactive Map**: OpenStreetMap via Leaflet with click-to-weather functionality
- **🔄 Temperature Toggle**: Switch between Celsius and Fahrenheit on the fly
- **📱 Fully Responsive**: Adapts gracefully from desktop to mobile
- **🎯 Lucide Icons**: No emojis - only high-quality Lucide React icons
- **⚡ Bun**: Fast package management and builds

### API Documentation
- **🚀 Futuristic UI**: Modern sci-fi inspired design at `/docs`
- **🔧 Stoplight Elements**: Interactive OpenAPI documentation (not Swagger)
- **🎨 Dark Theme**: Gradient backgrounds with glass-morphism effects
- **📝 Auto-Generated**: OpenAPI spec generated automatically from Go code

## 🚀 Quick Start

### Option 1: One-Command Startup (Recommended)

```bash
# Clone the repository
git clone https://github.com/4cecoder/weather-api-go.git
cd weather-api-go

# Run everything (backend + frontend)
./run.sh          # macOS/Linux
# OR
.\run.ps1         # Windows
```

Then open:
- **Main App**: http://localhost:3000
- **API Docs**: http://localhost:3000/docs

### Option 2: Docker

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Option 3: Manual Development

```bash
# Backend
go build -o weather-api .
./weather-api

# Frontend (in another terminal)
cd frontend
bun install
bun run dev
```

## 🌐 Available URLs

| Service | URL | Description |
|---------|-----|-------------|
| Main App | http://localhost:3000 | Interactive weather map |
| Frontend Dev | http://localhost:5173 | React dev server (if running) |
| **API Docs** | **http://localhost:3000/docs** | **Futuristic API documentation** |
| Health Check | http://localhost:3000/api/health | Service health status |
| Weather API | http://localhost:3000/api/weather?lat=40.7128&lon=-74.0060 | Get weather data |

## 📡 API Endpoints

### GET /api/weather
Returns current weather forecast for coordinates with both Celsius and Fahrenheit.

**Parameters:**
- `lat` (required): Latitude (-90 to 90)
- `lon` (required): Longitude (-180 to 180)

**Example Request:**
```bash
curl "http://localhost:3000/api/weather?lat=40.7128&lon=-74.0060"
```

**Example Response:**
```json
{
  "forecast": "Partly Cloudy",
  "temperature": "moderate",
  "temperature_c": 22.5,
  "temperature_f": 72.5
}
```

### GET /api/health
Health check endpoint.

**Example Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### GET /docs
**Futuristic interactive API documentation** - Stoplight Elements with:
- Auto-generated from OpenAPI spec
- Try-it-out functionality
- Dark gradient theme with glow effects
- Links to GitHub repo

## 🏗️ Architecture

### Layered Backend Structure
```
internal/
├── handlers/      # HTTP handlers (Fiber)
│   ├── weather.go # Weather endpoint handlers
│   └── docs.go    # API documentation
├── services/      # Business logic
│   ├── weather.go # Weather service with temp conversion
│   └── nws_client.go # NWS API client
├── repository/    # Data access layer
│   └── weather.go # Redis + SQLite caching
└── models/        # Data structures
    └── weather.go # Request/response types
```

### Frontend Architecture
```
frontend/
├── src/
│   ├── App.tsx       # Main app with TanStack Query
│   ├── App.test.tsx  # Unit tests
│   └── index.css     # Swiss luxury spa styling
├── e2e/              # Playwright E2E tests
└── package.json      # Bun dependencies
```

### Caching Strategy
1. **Redis** (Primary): Sub-millisecond response times
2. **SQLite** (Fallback): Persistent storage for durability

**Cache TTL**: 1 hour

### Temperature Classification
- **Hot**: ≥ 30°C (86°F) - shown in coral
- **Cold**: ≤ 10°C (50°F) - shown in blue
- **Moderate**: 10°C - 30°C - shown in green

## 🧪 Testing

### Backend Tests
```bash
# Run all Go tests
make backend-test

# With coverage
make test-coverage
```

### Frontend Tests
```bash
# Run unit tests
make frontend-test

# Run E2E tests
make e2e-test
```

### Full Test Suite
```bash
# Run everything (backend + frontend + build)
make ci
```

## 🛠️ Development

### Build Pipeline (Organized Stages)
```bash
# Backend
make backend-test      # Stage 1: Run tests
make backend-build     # Stage 2: Build binary
make backend-run       # Stage 3: Run locally

# Frontend
make frontend-deps     # Stage 1: Install deps
make frontend-test     # Stage 2: Run tests
make frontend-build    # Stage 3: Build production

# Full CI
make ci               # Run complete pipeline
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | 3000 |
| `REDIS_URL` | Redis connection URL | localhost:6379 |
| `DATABASE_URL` | SQLite database path | ./weather_cache.db |

## 📁 Project Structure

```
weather-api-go/
├── cmd/weather-api/           # Application entry point
├── internal/
│   ├── handlers/              # HTTP handlers
│   ├── services/              # Business logic
│   ├── repository/            # Data access
│   └── models/                # Data structures
├── frontend/                  # React + TanStack frontend
│   ├── src/
│   ├── e2e/                   # Playwright tests
│   └── package.json
├── .github/workflows/         # CI/CD pipeline
├── dist/frontend/            # Built frontend files
├── Dockerfile               # Multi-stage Docker build
├── docker-compose.yml        # Service orchestration
├── run.sh                    # Unix startup script
├── run.ps1                   # Windows startup script
├── Makefile                  # Organized build pipeline
└── README.md                 # This file
```

## 🎨 Design

Premium, minimalist aesthetic with Swiss luxury influences. Uses Lucide React icons exclusively (no emojis), neutral color palette, and responsive design from desktop to mobile.

## 📝 Development Notes

### Trade-offs & Future Improvements

1. **Caching**: Currently uses 1-hour TTL. For production:
   - Consider stale-while-revalidate pattern
   - Implement cache warming strategies
   - Different TTLs for varying freshness needs

2. **Database**: SQLite for simplicity. For production:
   - PostgreSQL or MySQL for better concurrency
   - Connection pooling
   - Database migrations

3. **Testing**: Unit + E2E tests present. Consider adding:
   - Load tests
   - Chaos engineering tests
   - Contract tests

4. **Monitoring**: Add for production:
   - Structured logging (e.g., Zap)
   - Metrics collection (Prometheus)
   - Distributed tracing

## 📄 License

MIT License - See LICENSE file for details

## 🤝 Contributing

Contributions welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Follow the existing code style
4. Add tests for new features
5. Submit a pull request

## 🙏 Acknowledgments

- [National Weather Service](https://api.weather.gov/) for weather data
- [OpenStreetMap](https://www.openstreetmap.org/) for map tiles
- [Stoplight](https://stoplight.io/) for Elements documentation
- [TanStack](https://tanstack.com/) for Query and modern React patterns

---

<p align="center">
  <sub>Built with ❤️ by <a href="https://github.com/4cecoder">4cecoder</a></sub>
</p>