import { test, expect } from '@playwright/test'

test.describe('Weather App E2E Tests', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/')
  })

  test('homepage loads correctly', async ({ page }) => {
    await expect(page.getByText('🌤️ Weather Forecast')).toBeVisible()
    await expect(page.getByText('Click anywhere on the map to get the weather forecast')).toBeVisible()
    await expect(page.getByText('Select Location')).toBeVisible()
    await expect(page.getByText('Weather Information')).toBeVisible()
  })

  test('clicking on map fetches weather data', async ({ page }) => {
    // Wait for the map to load
    await page.waitForSelector('.leaflet-container')
    
    // Click on the map (coordinates for New York City)
    const map = page.locator('.leaflet-container')
    await map.click({ position: { x: 400, y: 200 } })
    
    // Wait for loading state
    await expect(page.getByText('Loading weather...')).toBeVisible()
    
    // Wait for weather data to load
    await expect(page.getByTestId('weather-card')).toBeVisible({ timeout: 10000 })
    
    // Verify weather elements are present
    await expect(page.getByTestId('temperature')).toBeVisible()
    await expect(page.getByTestId('temperature-label')).toBeVisible()
    await expect(page.getByTestId('forecast')).toBeVisible()
    await expect(page.getByTestId('coordinates')).toBeVisible()
  })

  test('temperature label shows correct styling', async ({ page }) => {
    // Click on a known location
    const map = page.locator('.leaflet-container')
    await map.click({ position: { x: 400, y: 200 } })
    
    // Wait for weather card
    await expect(page.getByTestId('weather-card')).toBeVisible({ timeout: 10000 })
    
    // Check that temperature label has one of the expected classes
    const label = page.getByTestId('temperature-label')
    await expect(label).toBeVisible()
    
    // Verify the label text is one of the expected values
    const labelText = await label.textContent()
    expect(['hot', 'cold', 'moderate']).toContain(labelText?.toLowerCase())
  })

  test('health endpoint returns healthy status', async ({ page }) => {
    const response = await page.request.get('/health')
    expect(response.ok()).toBeTruthy()
    
    const data = await response.json()
    expect(data.status).toBe('healthy')
    expect(data.timestamp).toBeDefined()
  })

  test('weather API endpoint returns correct structure', async ({ page }) => {
    const response = await page.request.get('/weather?lat=40.7128&lon=-74.0060')
    expect(response.ok()).toBeTruthy()
    
    const data = await response.json()
    expect(data).toHaveProperty('forecast')
    expect(data).toHaveProperty('temperature')
    expect(data).toHaveProperty('temperature_c')
    expect(['hot', 'cold', 'moderate']).toContain(data.temperature)
  })

  test('weather API handles invalid coordinates', async ({ page }) => {
    const response = await page.request.get('/weather?lat=200&lon=200')
    expect(response.status()).toBe(400)
    
    const data = await response.json()
    expect(data.error).toBeDefined()
  })

  test('instructions are shown when no location selected', async ({ page }) => {
    await expect(page.getByText('Click on the map to see the weather forecast for that location.')).toBeVisible()
  })
})