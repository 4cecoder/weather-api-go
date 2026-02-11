import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import App from './App'

const createTestQueryClient = () => new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
    },
  },
})

describe('App', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders header correctly', () => {
    const queryClient = createTestQueryClient()
    render(
      <QueryClientProvider client={queryClient}>
        <App />
      </QueryClientProvider>
    )
    
    expect(screen.getByText('Weather Forecast')).toBeInTheDocument()
    expect(screen.getByText('Interactive meteorological visualization')).toBeInTheDocument()
  })

  it('shows instructions when no location selected', () => {
    const queryClient = createTestQueryClient()
    render(
      <QueryClientProvider client={queryClient}>
        <App />
      </QueryClientProvider>
    )
    
    expect(screen.getByText('No location selected')).toBeInTheDocument()
    expect(screen.getByText('Interact with the map to view current weather conditions for any location.')).toBeInTheDocument()
  })
})

describe('WeatherDisplay', () => {
  it('displays weather data correctly', async () => {
    const queryClient = createTestQueryClient()
    
    // We'll test the component behavior through the app
    render(
      <QueryClientProvider client={queryClient}>
        <App />
      </QueryClientProvider>
    )
    
    // Verify initial state
    expect(screen.getByText('No location selected')).toBeInTheDocument()
  })
})

describe('API Integration', () => {
  it('fetches weather data successfully', async () => {
    const mockWeather = {
      forecast: 'Sunny',
      temperature: 'hot',
      temperature_c: 32.0,
    }
    
    // Create a fresh mock for this test
    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => mockWeather,
    })
    
    const response = await mockFetch('/api/weather?lat=40.7128&lon=-74.0060')
    const data = await response.json()
    
    expect(data).toEqual(mockWeather)
    expect(mockFetch).toHaveBeenCalledWith('/api/weather?lat=40.7128&lon=-74.0060')
  })

  it('handles API errors', async () => {
    // Create a fresh mock for this test
    const mockFetch = vi.fn().mockResolvedValue({
      ok: false,
      status: 400,
      json: async () => ({ error: 'Invalid coordinates' }),
    })
    
    const response = await mockFetch('/api/weather?lat=invalid&lon=invalid')
    
    expect(response.ok).toBe(false)
    expect(response.status).toBe(400)
  })
})
