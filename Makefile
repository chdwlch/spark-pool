.PHONY: build run clean test help

# Default target
all: build

# Build the pool operator
build:
	@echo "Building mining pool demo..."
	go build -o bin/pool-operator cmd/pool-operator/main.go

# Run the pool operator
run: build
	@echo "Starting mining pool demo..."
	./bin/pool-operator

# Run in development mode
dev:
	@echo "Starting mining pool demo in development mode..."
	go run cmd/pool-operator/main.go --block-interval 10s

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run

# Generate documentation
docs:
	@echo "Generating documentation..."
	godoc -http=:6060

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t mining-pool-demo .

# Docker run
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 mining-pool-demo

# Show help
help:
	@echo "Available targets:"
	@echo "  build        - Build the pool operator"
	@echo "  run          - Build and run the pool operator"
	@echo "  dev          - Run in development mode (10s block interval)"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Run tests"
	@echo "  deps         - Install dependencies"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code"
	@echo "  docs         - Generate documentation"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  help         - Show this help message" 