package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// getOpenAPISpec returns the OpenAPI specification
func getOpenAPISpec() map[string]interface{} {
	return map[string]interface{}{
		"openapi": "3.0.0",
		"info": map[string]interface{}{
			"title":       "Weather API",
			"version":     "1.0.0",
			"description": "A modern weather service providing forecast data with dual-layer caching",
			"contact": map[string]interface{}{
				"name":  "API Support",
				"email": "support@weather-api.example.com",
			},
			"license": map[string]interface{}{
				"name": "MIT",
				"url":  "https://opensource.org/licenses/MIT",
			},
		},
		"servers": []map[string]interface{}{
			{
				"url":         "http://localhost:3000/api",
				"description": "Local development server",
			},
		},
		"paths": map[string]interface{}{
			"/weather": map[string]interface{}{
				"get": map[string]interface{}{
					"summary":     "Get weather forecast",
					"description": "Returns current weather forecast for given coordinates",
					"tags":        []string{"Weather"},
					"parameters": []map[string]interface{}{
						{
							"name":        "lat",
							"in":          "query",
							"required":    true,
							"schema":      map[string]interface{}{"type": "number"},
							"description": "Latitude (-90 to 90)",
							"example":     40.7128,
						},
						{
							"name":        "lon",
							"in":          "query",
							"required":    true,
							"schema":      map[string]interface{}{"type": "number"},
							"description": "Longitude (-180 to 180)",
							"example":     -74.0060,
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Weather data retrieved successfully",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"type": "object",
										"properties": map[string]interface{}{
											"forecast": map[string]interface{}{
												"type":        "string",
												"example":     "Partly Cloudy",
												"description": "Short weather forecast",
											},
											"temperature": map[string]interface{}{
												"type":        "string",
												"enum":        []string{"hot", "cold", "moderate"},
												"description": "Temperature classification",
											},
											"temperature_c": map[string]interface{}{
												"type":        "number",
												"example":     22.5,
												"description": "Temperature in Celsius",
											},
											"temperature_f": map[string]interface{}{
												"type":        "number",
												"example":     72.5,
												"description": "Temperature in Fahrenheit",
											},
										},
									},
								},
							},
						},
						"400": map[string]interface{}{
							"description": "Invalid parameters",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"type": "object",
										"properties": map[string]interface{}{
											"error":   map[string]interface{}{"type": "string"},
											"details": map[string]interface{}{"type": "string"},
										},
									},
								},
							},
						},
					},
				},
			},
			"/health": map[string]interface{}{
				"get": map[string]interface{}{
					"summary":     "Health check",
					"description": "Check API health status",
					"tags":        []string{"System"},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Service is healthy",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"type": "object",
										"properties": map[string]interface{}{
											"status":    map[string]interface{}{"type": "string", "example": "healthy"},
											"timestamp": map[string]interface{}{"type": "string"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"tags": []map[string]interface{}{
			{"name": "Weather", "description": "Weather forecast operations"},
			{"name": "System", "description": "System health and status"},
		},
	}
}

// GetAPIDocsHTML returns the futuristic API documentation HTML
func GetAPIDocsHTML() string {
	apiSpec, _ := json.Marshal(getOpenAPISpec())

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Weather API - Interactive Documentation</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet">
    <script src="https://unpkg.com/@stoplight/elements@8.0.0/web-components.min.js"></script>
    <link rel="stylesheet" href="https://unpkg.com/@stoplight/elements@8.0.0/styles.min.css">
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
            background: linear-gradient(135deg, #0f0f23 0%%, #1a1a2e 50%%, #16213e 100%%);
            min-height: 100vh;
            color: #e4e4e7;
        }
        
        .header {
            background: rgba(15, 15, 35, 0.8);
            backdrop-filter: blur(20px);
            border-bottom: 1px solid rgba(255, 255, 255, 0.1);
            padding: 1.5rem 2rem;
            position: sticky;
            top: 0;
            z-index: 100;
        }
        
        .header-content {
            max-width: 1400px;
            margin: 0 auto;
            display: flex;
            align-items: center;
            justify-content: space-between;
        }
        
        .logo {
            display: flex;
            align-items: center;
            gap: 0.75rem;
        }
        
        .logo-icon {
            width: 40px;
            height: 40px;
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            border-radius: 12px;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 1.25rem;
        }
        
        .logo-text {
            font-size: 1.25rem;
            font-weight: 600;
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
        }
        
        .badge {
            background: rgba(102, 126, 234, 0.2);
            color: #667eea;
            padding: 0.25rem 0.75rem;
            border-radius: 100px;
            font-size: 0.75rem;
            font-weight: 600;
            letter-spacing: 0.05em;
        }
        
        .header-actions {
            display: flex;
            gap: 1rem;
        }
        
        .btn {
            display: inline-flex;
            align-items: center;
            gap: 0.5rem;
            padding: 0.625rem 1.25rem;
            border-radius: 8px;
            font-size: 0.875rem;
            font-weight: 500;
            text-decoration: none;
            transition: all 0.2s ease;
        }
        
        .btn-primary {
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white;
            border: none;
        }
        
        .btn-primary:hover {
            transform: translateY(-2px);
            box-shadow: 0 8px 25px rgba(102, 126, 234, 0.4);
        }
        
        .btn-secondary {
            background: rgba(255, 255, 255, 0.05);
            color: #e4e4e7;
            border: 1px solid rgba(255, 255, 255, 0.1);
        }
        
        .btn-secondary:hover {
            background: rgba(255, 255, 255, 0.1);
            border-color: rgba(255, 255, 255, 0.2);
        }
        
        .main-content {
            max-width: 1400px;
            margin: 0 auto;
            padding: 2rem;
            height: calc(100vh - 80px);
        }
        
        .elements-container {
            background: rgba(255, 255, 255, 0.03);
            border: 1px solid rgba(255, 255, 255, 0.1);
            border-radius: 16px;
            height: 100%%;
            overflow: hidden;
            box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
        }
        
        .glow {
            position: fixed;
            width: 600px;
            height: 600px;
            background: radial-gradient(circle, rgba(102, 126, 234, 0.15) 0%%, transparent 70%%);
            top: -300px;
            right: -300px;
            pointer-events: none;
            z-index: 0;
        }
        
        .glow-2 {
            position: fixed;
            width: 400px;
            height: 400px;
            background: radial-gradient(circle, rgba(118, 75, 162, 0.1) 0%%, transparent 70%%);
            bottom: -200px;
            left: -200px;
            pointer-events: none;
            z-index: 0;
        }
    </style>
</head>
<body>
    <div class="glow"></div>
    <div class="glow-2"></div>
    
    <header class="header">
        <div class="header-content">
            <div class="logo">
                <div class="logo-icon">🌤️</div>
                <span class="logo-text">Weather API</span>
                <span class="badge">v1.0.0</span>
            </div>
            <div class="header-actions">
                <a href="/" class="btn btn-secondary">Back to App</a>
                <a href="https://github.com/4cecoder/weather-api-go" target="_blank" class="btn btn-primary">View on GitHub</a>
            </div>
        </div>
    </header>
    
    <main class="main-content">
        <div class="elements-container">
            <elements-api
                apiDescriptionDocument='%s'
                router="hash"
                layout="sidebar"
                hideSchemas="false"
                logo="false"
            />
        </div>
    </main>
</body>
</html>`, string(apiSpec))
}

// ServeAPIDocs serves the futuristic API documentation page for Fiber
func ServeAPIDocs(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(GetAPIDocsHTML())
}
