#!/bin/bash

# Ark Virtual Channels Mining Pool Demo Stop Script

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Stop the demo server
stop_server() {
    if [ -f .server.pid ]; then
        SERVER_PID=$(cat .server.pid)
        if kill -0 $SERVER_PID 2>/dev/null; then
            print_status "Stopping mining pool demo server (PID: $SERVER_PID)..."
            kill $SERVER_PID
            rm .server.pid
            print_success "Server stopped successfully"
        else
            print_error "Server process not found"
            rm .server.pid
        fi
    else
        print_error "No server PID file found"
    fi
}

# Clean up any remaining processes
cleanup() {
    print_status "Cleaning up..."
    
    # Kill any remaining pool-operator processes
    pkill -f "pool-operator" 2>/dev/null || true
    
    # Remove build artifacts
    rm -rf bin/
    
    print_success "Cleanup complete"
}

# Main execution
main() {
    print_status "Stopping Ark Virtual Channels Mining Pool Demo..."
    
    stop_server
    cleanup
    
    print_success "Demo stopped successfully! 👋"
}

# Run main function
main "$@" 