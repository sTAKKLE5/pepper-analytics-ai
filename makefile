# Binary name
BINARY_NAME=server

.PHONY: all build run generate clean prod dev tidy deps

# Default target
all: generate build run-dev

# Tidy and verify dependencies
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

# Install dependencies
deps: tidy
	@echo "Verifying dependencies..."
	go mod verify

# Generate templ files
generate:
	@echo "Generating templ files..."
	templ generate

# Build the application
build: deps
	@echo "Building..."
	go build -o ${BINARY_NAME} cmd/server/main.go

# Run the application in development
run-dev:
	@echo "Running in development mode..."
	GO_ENV=development ./${BINARY_NAME}

# Run the application in production
run-prod:
	@echo "Running in production mode..."
	GO_ENV=production ./${BINARY_NAME}

# Development target (default)
dev: generate build run-dev

# Production target
prod: generate build run-prod

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f ${BINARY_NAME}
	go clean

# Help target
help:
	@echo "Available targets:"
	@echo "  make dev              - Run in development mode (default)"
	@echo "  make prod             - Run in production mode"
	@echo "  make tidy             - Tidy go modules"
	@echo "  make deps             - Verify and tidy dependencies"
	@echo "  make clean            - Clean build artifacts"
	@echo "  make help             - Show this help message"