import { useState } from 'react'
import { MapContainer, TileLayer, Marker, useMapEvents } from 'react-leaflet'
import { QueryClient, QueryClientProvider, useQuery } from '@tanstack/react-query'
import { 
  Cloud, 
  MapPin, 
  Thermometer, 
  Navigation, 
  Droplets, 
  Wind,
  Sun,
  Snowflake,
  CloudRain,
  CloudLightning,
  Loader2,
  AlertCircle,
  MousePointer2,
  ToggleLeft,
  ToggleRight
} from 'lucide-react'
import 'leaflet/dist/leaflet.css'

const queryClient = new QueryClient()

interface WeatherData {
  forecast: string
  temperature: string
  temperature_c: number
  temperature_f: number
}

// Get weather icon based on forecast and temperature
function getWeatherIcon(forecast: string, temperature: string) {
  const forecastLower = forecast.toLowerCase()
  
  if (forecastLower.includes('rain') || forecastLower.includes('shower')) {
    return <CloudRain className="weather-icon rain" size={64} strokeWidth={1.5} />
  }
  if (forecastLower.includes('snow') || forecastLower.includes('blizzard')) {
    return <Snowflake className="weather-icon cold" size={64} strokeWidth={1.5} />
  }
  if (forecastLower.includes('storm') || forecastLower.includes('thunder')) {
    return <CloudLightning className="weather-icon storm" size={64} strokeWidth={1.5} />
  }
  if (forecastLower.includes('cloud')) {
    return <Cloud className="weather-icon" size={64} strokeWidth={1.5} />
  }
  if (forecastLower.includes('sun') || forecastLower.includes('clear')) {
    return <Sun className="weather-icon hot" size={64} strokeWidth={1.5} />
  }
  
  // Default based on temperature
  if (temperature === 'hot') {
    return <Sun className="weather-icon hot" size={64} strokeWidth={1.5} />
  }
  if (temperature === 'cold') {
    return <Snowflake className="weather-icon cold" size={64} strokeWidth={1.5} />
  }
  
  return <Cloud className="weather-icon" size={64} strokeWidth={1.5} />
}

function LocationMarker({ onLocationSelect }: { onLocationSelect: (lat: number, lng: number) => void }) {
  useMapEvents({
    click(e) {
      onLocationSelect(e.latlng.lat, e.latlng.lng)
    },
  })
  return null
}

function WeatherDisplay({ lat, lng, unit, onUnitChange }: { 
  lat: number; 
  lng: number;
  unit: 'C' | 'F';
  onUnitChange: () => void;
}) {
  const { data, isLoading, error } = useQuery<WeatherData>({
    queryKey: ['weather', lat, lng],
    queryFn: async () => {
      const response = await fetch(`/api/weather?lat=${lat}&lon=${lng}`)
      if (!response.ok) throw new Error('Failed to fetch weather')
      return response.json()
    },
  })

  if (isLoading) {
    return (
      <div className="weather-state loading-state">
        <Loader2 className="spinner" size={32} strokeWidth={1.5} />
        <p>Retrieving forecast...</p>
      </div>
    )
  }

  if (error) {
    return (
      <div className="weather-state error-state">
        <AlertCircle size={32} strokeWidth={1.5} />
        <p>Unable to retrieve weather data</p>
        <span className="error-message">{error.message}</span>
      </div>
    )
  }

  if (!data) return null

  const tempValue = unit === 'C' ? data.temperature_c : data.temperature_f
  const unitLabel = unit === 'C' ? '°C' : '°F'

  return (
    <div className="weather-card" data-testid="weather-card">
      <div className="weather-visual">
        {getWeatherIcon(data.forecast, data.temperature)}
      </div>
      
      <div className="temperature-display">
        <span className="temperature-value" data-testid="temperature">
          {tempValue.toFixed(1)}
        </span>
        <span className="temperature-unit">{unitLabel}</span>
      </div>
      
      <button 
        className="unit-toggle" 
        onClick={onUnitChange}
        aria-label={`Switch to ${unit === 'C' ? 'Fahrenheit' : 'Celsius'}`}
      >
        {unit === 'C' ? (
          <>
            <ToggleLeft size={20} strokeWidth={2} />
            <span>°C</span>
          </>
        ) : (
          <>
            <ToggleRight size={20} strokeWidth={2} />
            <span>°F</span>
          </>
        )}
      </button>
      
      <div className={`temperature-badge ${data.temperature}`} data-testid="temperature-label">
        <Thermometer size={14} strokeWidth={2} />
        <span>{data.temperature}</span>
      </div>
      
      <div className="forecast-text" data-testid="forecast">
        {data.forecast}
      </div>
      
      <div className="coordinates-display" data-testid="coordinates">
        <Navigation size={14} strokeWidth={2} />
        <span>{lat.toFixed(4)}, {lng.toFixed(4)}</span>
      </div>
    </div>
  )
}

function AppContent() {
  const [selectedLocation, setSelectedLocation] = useState<{ lat: number; lng: number } | null>(null)
  const [unit, setUnit] = useState<'C' | 'F'>('C')

  const toggleUnit = () => {
    setUnit(prev => prev === 'C' ? 'F' : 'C')
  }

  return (
    <div className="app">
      <header className="header">
        <div className="header-icon">
          <Cloud size={28} strokeWidth={1.5} />
        </div>
        <h1>Weather Forecast</h1>
        <p className="header-subtitle">Interactive meteorological visualization</p>
      </header>

      <main className="main-content">
        <section className="panel map-panel">
          <div className="panel-header">
            <MapPin size={18} strokeWidth={1.5} />
            <h2>Select Location</h2>
          </div>
          <div className="map-wrapper">
            <MapContainer
              center={[39.8283, -98.5795]}
              zoom={4}
              scrollWheelZoom={true}
            >
              <TileLayer
                attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
                url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
              />
              <LocationMarker onLocationSelect={(lat, lng) => setSelectedLocation({ lat, lng })} />
              {selectedLocation && (
                <Marker position={[selectedLocation.lat, selectedLocation.lng]} />
              )}
            </MapContainer>
          </div>
          <div className="map-hint">
            <MousePointer2 size={14} strokeWidth={2} />
            <span>Click anywhere on the map to retrieve weather data</span>
          </div>
        </section>

        <section className="panel weather-panel">
          <div className="panel-header">
            <Thermometer size={18} strokeWidth={1.5} />
            <h2>Weather Information</h2>
          </div>
          {selectedLocation ? (
            <WeatherDisplay 
              lat={selectedLocation.lat} 
              lng={selectedLocation.lng} 
              unit={unit}
              onUnitChange={toggleUnit}
            />
          ) : (
            <div className="empty-state">
              <div className="empty-icon">
                <Droplets size={48} strokeWidth={1} />
              </div>
              <p className="empty-title">No location selected</p>
              <p className="empty-description">
                Interact with the map to view current weather conditions for any location.
              </p>
            </div>
          )}
        </section>
      </main>
      
      <footer className="footer">
        <Wind size={14} strokeWidth={2} />
        <span>Precision meteorological data</span>
      </footer>
    </div>
  )
}

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <AppContent />
    </QueryClientProvider>
  )
}

export default App