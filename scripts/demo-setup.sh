#!/bin/bash

# Ark Virtual Channels Mining Pool Demo Setup Script
# This script sets up and runs the mining pool demo

set -e

echo "🚀 Setting up Ark Virtual Channels Mining Pool Demo..."

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

# Check if Go is installed
check_go() {
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install Go 1.21 or later."
        exit 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    print_success "Go version $GO_VERSION found"
}

# Check if we're in the right directory
check_directory() {
    if [ ! -f "go.mod" ]; then
        print_error "go.mod not found. Please run this script from the mining-pool-demo directory."
        exit 1
    fi
    print_success "Running from correct directory"
}

# Install dependencies
install_deps() {
    print_status "Installing dependencies..."
    go mod tidy
    go mod download
    print_success "Dependencies installed"
}

# Build the application
build_app() {
    print_status "Building mining pool demo..."
    go build -o bin/pool-operator cmd/pool-operator/main.go
    print_success "Application built successfully"
}

# Create demo data
create_demo_data() {
    print_status "Creating demo miners..."
    
    # Wait for server to start
    sleep 3
    
    # Add demo miners
    curl -X POST http://localhost:8080/api/v1/miners \
        -H "Content-Type: application/json" \
        -d '{"miner_name":"Alice","address":"bc1qalice123456789","hash_rate":100}' \
        -s > /dev/null
    
    curl -X POST http://localhost:8080/api/v1/miners \
        -H "Content-Type: application/json" \
        -d '{"miner_name":"Bob","address":"bc1qbob123456789","hash_rate":75}' \
        -s > /dev/null
    
    curl -X POST http://localhost:8080/api/v1/miners \
        -H "Content-Type: application/json" \
        -d '{"miner_name":"Charlie","address":"bc1qcharlie123456789","hash_rate":50}' \
        -s > /dev/null
    
    print_success "Demo miners created"
}

# Start the application
start_app() {
    print_status "Starting mining pool demo..."
    
    # Run in background
    ./bin/pool-operator --block-interval 15s &
    SERVER_PID=$!
    
    # Save PID for cleanup
    echo $SERVER_PID > .server.pid
    
    print_success "Server started with PID $SERVER_PID"
    print_status "Dashboard available at: http://localhost:8080"
}

# Show demo instructions
show_instructions() {
    echo ""
    echo "🎯 Demo Instructions:"
    echo "===================="
    echo ""
    echo "1. Open your browser and go to: http://localhost:8080"
    echo "2. You'll see the pool operator dashboard with 3 demo miners"
    echo "3. Click 'Start All Miners' to begin mining simulation"
    echo "4. Watch as block rewards are processed every 15 seconds"
    echo "5. Monitor Virtual Channels being created and updated"
    echo "6. Open individual miner dashboards to see detailed stats"
    echo ""
    echo "📊 Demo Features:"
    echo "- Real-time mining simulation"
    echo "- Virtual Channel creation and management"
    echo "- Automatic block reward distribution"
    echo "- WebSocket-powered live updates"
    echo "- Individual miner dashboards"
    echo ""
    echo "🛑 To stop the demo:"
    echo "Press Ctrl+C or run: ./scripts/demo-stop.sh"
    echo ""
}

# Cleanup function
cleanup() {
    if [ -f .server.pid ]; then
        SERVER_PID=$(cat .server.pid)
        if kill -0 $SERVER_PID 2>/dev/null; then
            print_status "Stopping server (PID: $SERVER_PID)..."
            kill $SERVER_PID
            rm .server.pid
            print_success "Server stopped"
        fi
    fi
}

# Trap Ctrl+C and cleanup
trap cleanup EXIT INT TERM

# Main execution
main() {
    print_status "Starting Ark Virtual Channels Mining Pool Demo setup..."
    
    check_go
    check_directory
    install_deps
    build_app
    start_app
    create_demo_data
    show_instructions
    
    print_success "Demo setup complete! 🎉"
    print_status "The demo will continue running. Press Ctrl+C to stop."
    
    # Keep script running
    wait
}

# Run main function
main "$@" 