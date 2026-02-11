#!/bin/bash

# Weather API - Development Startup Script
# This script starts both the backend and frontend in development mode

set -e

echo "🚀 Starting Weather API Development Environment"
echo "================================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if required tools are installed
check_dependencies() {
    print_status "Checking dependencies..."
    
    # Check Go
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install Go 1.24 or later."
        exit 1
    fi
    
    # Check Bun
    if ! command -v bun &> /dev/null; then
        print_error "Bun is not installed. Please install Bun."
        print_status "Visit: https://bun.sh/docs/installation"
        exit 1
    fi
    
    print_success "All dependencies found"
}

# Setup backend
setup_backend() {
    print_status "Setting up backend..."
    
    # Download Go dependencies
    go mod download
    
    print_success "Backend setup complete"
}

# Setup frontend
setup_frontend() {
    print_status "Setting up frontend..."
    
    cd frontend
    
    # Install Bun dependencies
    if [ ! -d "node_modules" ]; then
        print_status "Installing frontend dependencies with Bun..."
        bun install
    else
        print_status "Frontend dependencies already installed"
    fi
    
    cd ..
    print_success "Frontend setup complete"
}

# Start backend
start_backend() {
    print_status "Starting backend server..."
    
    # Build and run backend in background
    go build -o weather-api . && ./weather-api &
    BACKEND_PID=$!
    
    # Wait for backend to start
    sleep 2
    
    # Check if backend is running
    if kill -0 $BACKEND_PID 2>/dev/null; then
        print_success "Backend started on http://localhost:3000 (PID: $BACKEND_PID)"
        echo $BACKEND_PID > .backend.pid
    else
        print_error "Failed to start backend"
        exit 1
    fi
}

# Start frontend dev server
start_frontend() {
    print_status "Starting frontend dev server..."
    
    cd frontend
    
    # Start frontend dev server in background
    bun run dev -- --host &
    FRONTEND_PID=$!
    
    # Wait for frontend to start
    sleep 3
    
    # Check if frontend is running
    if kill -0 $FRONTEND_PID 2>/dev/null; then
        print_success "Frontend started on http://localhost:5173 (PID: $FRONTEND_PID)"
        echo $FRONTEND_PID > ../.frontend.pid
    else
        print_error "Failed to start frontend"
        exit 1
    fi
    
    cd ..
}

# Create stop script
create_stop_script() {
    cat > stop.sh << 'EOF'
#!/bin/bash

echo "🛑 Stopping Weather API..."

# Stop backend
if [ -f .backend.pid ]; then
    BACKEND_PID=$(cat .backend.pid)
    if kill -0 $BACKEND_PID 2>/dev/null; then
        kill $BACKEND_PID
        echo "✅ Backend stopped (PID: $BACKEND_PID)"
    fi
    rm .backend.pid
fi

# Stop frontend
if [ -f .frontend.pid ]; then
    FRONTEND_PID=$(cat .frontend.pid)
    if kill -0 $FRONTEND_PID 2>/dev/null; then
        kill $FRONTEND_PID
        echo "✅ Frontend stopped (PID: $FRONTEND_PID)"
    fi
    rm .frontend.pid
fi

echo "✅ All services stopped"
EOF
    chmod +x stop.sh
    print_success "Created stop.sh script"
}

# Main execution
main() {
    echo ""
    
    # Check dependencies
    check_dependencies
    
    echo ""
    
    # Setup
    setup_backend
    setup_frontend
    
    echo ""
    
    # Start services
    start_backend
    start_frontend
    
    echo ""
    echo "================================================"
    print_success "Weather API is running!"
    echo ""
    echo "🌐 Available URLs:"
    echo "   • Backend API:    http://localhost:3000"
    echo "   • Frontend Dev:   http://localhost:5173"
    echo "   • Health Check:   http://localhost:3000/health"
    echo "   • API Weather:    http://localhost:3000/weather?lat=40.7128&lon=-74.0060"
    echo ""
    echo "📝 To stop the services, run: ./stop.sh"
    echo ""
    
    # Create stop script
    create_stop_script
    
    # Keep script running to maintain background processes
    print_status "Press Ctrl+C to stop all services"
    
    # Trap Ctrl+C
    trap 'echo ""; print_status "Stopping services..."; ./stop.sh; exit 0' INT
    
    # Wait indefinitely
    wait
}

# Run main function
main