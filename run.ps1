#Requires -Version 5.1
<#
.SYNOPSIS
    Weather API - Windows Development Startup Script
.DESCRIPTION
    Aggressively checks dependencies, sets up environment, and starts both backend and frontend services.
    This script ensures everything is properly configured before starting.
#>

[CmdletBinding()]
param()

$ErrorActionPreference = "Stop"
$ProgressPreference = "Continue"

# Colors for output
$Colors = @{
    Info = "Cyan"
    Success = "Green"
    Warning = "Yellow"
    Error = "Red"
    Emphasis = "Magenta"
}

function Write-Status {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor $Colors.Info
}

function Write-Success {
    param([string]$Message)
    Write-Host "[SUCCESS] $Message" -ForegroundColor $Colors.Success
}

function Write-Warning {
    param([string]$Message)
    Write-Host "[WARNING] $Message" -ForegroundColor $Colors.Warning
}

function Write-Error {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor $Colors.Error
}

function Write-Emphasis {
    param([string]$Message)
    Write-Host $Message -ForegroundColor $Colors.Emphasis
}

function Test-Command {
    param([string]$Command)
    $null = Get-Command $Command -ErrorAction SilentlyContinue
    return $?
}

function Test-Port {
    param([int]$Port)
    $listener = $null
    try {
        $listener = New-Object System.Net.Sockets.TcpListener([System.Net.IPAddress]::Loopback, $Port)
        $listener.Start()
        $listener.Stop()
        return $true
    }
    catch {
        return $false
    }
    finally {
        if ($listener -ne $null) {
            $listener.Stop()
        }
    }
}

function Wait-ForService {
    param(
        [string]$Url,
        [int]$TimeoutSeconds = 30
    )
    $startTime = Get-Date
    while ((Get-Date) - $startTime).TotalSeconds -lt $TimeoutSeconds) {
        try {
            $response = Invoke-WebRequest -Uri $Url -Method GET -TimeoutSec 2 -UseBasicParsing -ErrorAction SilentlyContinue
            if ($response.StatusCode -eq 200) {
                return $true
            }
        }
        catch {
            Start-Sleep -Milliseconds 500
        }
    }
    return $false
}

function Initialize-Environment {
    Write-Emphasis "🚀 Weather API - Windows Development Environment"
    Write-Emphasis "================================================"
    Write-Host ""
    
    # Aggressive dependency checking
    Write-Status "AGGRESSIVELY checking dependencies..."
    
    # Check Go
    if (-not (Test-Command "go")) {
        Write-Error "Go is NOT installed or NOT in PATH!"
        Write-Status "Install Go from: https://go.dev/dl/"
        Write-Status "Make sure to restart your terminal after installation!"
        exit 1
    }
    
    $goVersion = (go version) -replace "go version go", "" -replace " .+", ""
    Write-Success "Go found: v$goVersion"
    
    # Check Bun
    if (-not (Test-Command "bun")) {
        Write-Error "Bun is NOT installed or NOT in PATH!"
        Write-Status "Install Bun with: powershell -c "irm bun.sh/install.ps1|iex""
        Write-Status "Make sure to restart your terminal after installation!"
        exit 1
    }
    
    $bunVersion = (bun --version)
    Write-Success "Bun found: v$bunVersion"
    
    # Check ports
    Write-Status "Checking port availability..."
    if (-not (Test-Port 3000)) {
        Write-Error "Port 3000 is ALREADY IN USE!"
        Write-Status "Kill the process using port 3000 or change the port in main.go"
        exit 1
    }
    
    if (-not (Test-Port 5173)) {
        Write-Error "Port 5173 is ALREADY IN USE!"
        Write-Status "Kill the process using port 5173"
        exit 1
    }
    
    Write-Success "All ports available"
}

function Initialize-Backend {
    Write-Host ""
    Write-Status "Setting up BACKEND..."
    
    # Clean and rebuild
    if (Test-Path "weather-api.exe") {
        Write-Status "Removing old binary..."
        Remove-Item "weather-api.exe" -Force
    }
    
    # Download dependencies
    Write-Status "Downloading Go modules..."
    $output = go mod download 2>&1
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Failed to download Go modules!"
        Write-Error $output
        exit 1
    }
    
    Write-Status "Tidying Go modules..."
    $output = go mod tidy 2>&1
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Failed to tidy Go modules!"
        Write-Error $output
        exit 1
    }
    
    Write-Success "Backend setup complete"
}

function Initialize-Frontend {
    Write-Host ""
    Write-Status "Setting up FRONTEND..."
    
    Push-Location "frontend"
    
    try {
        # Clean node_modules if it exists but is corrupted
        if (Test-Path "node_modules") {
            Write-Status "Checking existing node_modules..."
            $packageCount = (Get-ChildItem "node_modules" -Directory | Measure-Object).Count
            if ($packageCount -lt 10) {
                Write-Warning "node_modules appears corrupted or incomplete!"
                Write-Status "Removing and reinstalling..."
                Remove-Item "node_modules" -Recurse -Force
            }
        }
        
        # Install dependencies
        if (-not (Test-Path "node_modules")) {
            Write-Status "Installing frontend dependencies with Bun..."
            $output = bun install 2>&1
            if ($LASTEXITCODE -ne 0) {
                Write-Error "Failed to install frontend dependencies!"
                Write-Error $output
                exit 1
            }
        }
        else {
            Write-Status "Frontend dependencies already installed"
        }
        
        Write-Success "Frontend setup complete"
    }
    finally {
        Pop-Location
    }
}

function Start-BackendService {
    Write-Host ""
    Write-Status "Starting BACKEND server..."
    
    # Build first
    Write-Status "Building backend binary..."
    $buildOutput = go build -o weather-api.exe . 2>&1
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Failed to build backend!"
        Write-Error $buildOutput
        exit 1
    }
    
    # Start backend
    Write-Status "Launching backend process..."
    $backendJob = Start-Job -ScriptBlock {
        Set-Location $using:PWD
        .\weather-api.exe 2>&1
    }
    
    # Wait and verify
    Write-Status "Waiting for backend to initialize..."
    Start-Sleep -Seconds 3
    
    if (Wait-ForService -Url "http://localhost:3000/health" -TimeoutSeconds 10) {
        Write-Success "Backend is RUNNING on http://localhost:3000"
        return $backendJob
    }
    else {
        Write-Error "Backend failed to start!"
        $jobOutput = Receive-Job $backendJob
        Write-Error "Backend output:"
        Write-Error $jobOutput
        Stop-Job $backendJob -ErrorAction SilentlyContinue
        Remove-Job $backendJob -ErrorAction SilentlyContinue
        exit 1
    }
}

function Start-FrontendService {
    Write-Host ""
    Write-Status "Starting FRONTEND dev server..."
    
    Push-Location "frontend"
    
    try {
        # Start frontend
        $frontendJob = Start-Job -ScriptBlock {
            Set-Location $using:PWD
            bun run dev -- --host 2>&1
        }
        
        # Wait and verify
        Write-Status "Waiting for frontend to initialize..."
        Start-Sleep -Seconds 5
        
        if (Wait-ForService -Url "http://localhost:5173" -TimeoutSeconds 15) {
            Write-Success "Frontend is RUNNING on http://localhost:5173"
            return $frontendJob
        }
        else {
            Write-Warning "Frontend may still be starting..."
            Write-Success "Frontend process launched (check http://localhost:5173)"
            return $frontendJob
        }
    }
    finally {
        Pop-Location
    }
}

function New-StopScript {
    $stopScript = @'
#Requires -Version 5.1
param()

$ErrorActionPreference = "SilentlyContinue"

Write-Host "🛑 Stopping Weather API services..." -ForegroundColor Cyan
Write-Host ""

# Kill backend
$backendProcesses = Get-Process | Where-Object { $_.ProcessName -eq "weather-api" -or $_.CommandLine -like "*weather-api*" }
if ($backendProcesses) {
    $backendProcesses | Stop-Process -Force
    Write-Host "✅ Backend stopped" -ForegroundColor Green
}

# Kill frontend (Node/Bun processes in frontend directory)
$frontendProcesses = Get-Process | Where-Object { 
    ($_.ProcessName -eq "bun" -or $_.ProcessName -eq "node") -and 
    $_.CommandLine -like "*frontend*"
}
if ($frontendProcesses) {
    $frontendProcesses | Stop-Process -Force
    Write-Host "✅ Frontend stopped" -ForegroundColor Green
}

# Kill any processes on our ports
$port3000 = Get-NetTCPConnection -LocalPort 3000 -ErrorAction SilentlyContinue | Select-Object -ExpandProperty OwningProcess
if ($port3000) {
    Stop-Process -Id $port3000 -Force
    Write-Host "✅ Port 3000 freed" -ForegroundColor Green
}

$port5173 = Get-NetTCPConnection -LocalPort 5173 -ErrorAction SilentlyContinue | Select-Object -ExpandProperty OwningProcess
if ($port5173) {
    Stop-Process -Id $port5173 -Force
    Write-Host "✅ Port 5173 freed" -ForegroundColor Green
}

Write-Host ""
Write-Host "✅ All services stopped" -ForegroundColor Green

# Clean up job files
Remove-Job * -Force -ErrorAction SilentlyContinue
'@
    
    $stopScript | Out-File -FilePath "stop.ps1" -Encoding UTF8
    Write-Success "Created stop.ps1 script"
}

# Main execution
try {
    Initialize-Environment
    Initialize-Backend
    Initialize-Frontend
    
    $backendJob = Start-BackendService
    $frontendJob = Start-FrontendService
    
    Write-Host ""
    Write-Emphasis "================================================"
    Write-Success "Weather API is RUNNING!"
    Write-Emphasis "================================================"
    Write-Host ""
    Write-Host "🌐 Available URLs:" -ForegroundColor White
    Write-Host "   • Backend API:    http://localhost:3000" -ForegroundColor Cyan
    Write-Host "   • Frontend Dev:   http://localhost:5173" -ForegroundColor Cyan
    Write-Host "   • API Docs:       http://localhost:3000/docs (Interactive Documentation)" -ForegroundColor Cyan
    Write-Host "   • Health Check:   http://localhost:3000/api/health" -ForegroundColor Cyan
    Write-Host "   • Weather API:    http://localhost:3000/api/weather?lat=40.7128&lon=-74.0060" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "📝 To stop all services:" -ForegroundColor White
    Write-Host "   • Run: .\stop.ps1" -ForegroundColor Yellow
    Write-Host "   • Or press Ctrl+C in this window" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "💡 Press Ctrl+C to stop all services gracefully" -ForegroundColor Magenta
    Write-Host ""
    
    # Create stop script
    New-StopScript
    
    # Monitor jobs
    while ($true) {
        Start-Sleep -Seconds 1
        
        # Check if jobs are still running
        $backendRunning = ($backendJob.State -eq "Running")
        $frontendRunning = ($frontendJob.State -eq "Running")
        
        if (-not $backendRunning) {
            Write-Error "BACKEND has stopped unexpectedly!"
            $output = Receive-Job $backendJob
            Write-Error "Last output:"
            Write-Error $output
            break
        }
        
        if (-not $frontendRunning) {
            Write-Warning "FRONTEND has stopped!"
            $output = Receive-Job $frontendJob
            Write-Warning "Last output:"
            Write-Warning $output
            break
        }
    }
}
catch {
    Write-Error "An error occurred: $_"
    Write-Error $_.ScriptStackTrace
    exit 1
}
finally {
    Write-Host ""
    Write-Status "Cleaning up..."
    
    # Stop all jobs
    Get-Job | Stop-Job -ErrorAction SilentlyContinue
    Get-Job | Remove-Job -ErrorAction SilentlyContinue
    
    Write-Host "✅ Cleanup complete" -ForegroundColor Green
}