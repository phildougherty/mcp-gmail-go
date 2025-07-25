# Gmail MCP Server Makefile

.PHONY: build run test clean auth docker-build docker-run help

# Variables
BINARY_NAME=gmail-mcp-server
MAIN_FILE=main.go
DOCKER_IMAGE=gmail-mcp-server
DOCKER_TAG=latest

# Default target
help: ## Show this help message
	@echo "Gmail MCP Server - Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Build targets
build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) $(MAIN_FILE)
	@echo "Build complete: $(BINARY_NAME)"

build-linux: ## Build for Linux
	@echo "Building $(BINARY_NAME) for Linux..."
	@GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux $(MAIN_FILE)
	@echo "Linux build complete: $(BINARY_NAME)-linux"

build-darwin: ## Build for macOS
	@echo "Building $(BINARY_NAME) for macOS..."
	@GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-darwin $(MAIN_FILE)
	@echo "macOS build complete: $(BINARY_NAME)-darwin"

build-windows: ## Build for Windows
	@echo "Building $(BINARY_NAME) for Windows..."
	@GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME)-windows.exe $(MAIN_FILE)
	@echo "Windows build complete: $(BINARY_NAME)-windows.exe"

build-all: build-linux build-darwin build-windows ## Build for all platforms

# Run targets
run: build ## Build and run the server
	@echo "Starting Gmail MCP Server..."
	@./$(BINARY_NAME)

run-debug: build ## Build and run with debug logging
	@echo "Starting Gmail MCP Server with debug logging..."
	@./$(BINARY_NAME) -debug

run-dev: ## Run without building (development)
	@echo "Running in development mode..."
	@go run $(MAIN_FILE) -debug

auth: build ## Run OAuth authentication flow
	@echo "Starting OAuth authentication..."
	@./$(BINARY_NAME) -auth

auth-dev: ## Run OAuth authentication in development
	@echo "Running OAuth authentication in development mode..."
	@go run $(MAIN_FILE) -auth

# Test targets
test: ## Run all tests
	@echo "Running tests..."
	@go test ./...

test-verbose: ## Run tests with verbose output
	@echo "Running tests with verbose output..."
	@go test -v ./...

test-cover: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Docker targets
docker-build: ## Build Docker image
	@echo "Building Docker image $(DOCKER_IMAGE):$(DOCKER_TAG)..."
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)"

docker-run: docker-build ## Build and run Docker container
	@echo "Running Docker container..."
	@docker run -p 8080:8080 -v $(PWD)/gcp-oauth.keys.json:/root/.gmail-mcp/gcp-oauth.keys.json:ro $(DOCKER_IMAGE):$(DOCKER_TAG)

docker-compose-up: ## Start with docker-compose
	@echo "Starting services with docker-compose..."
	@docker-compose up -d

docker-compose-down: ## Stop docker-compose services
	@echo "Stopping docker-compose services..."
	@docker-compose down

docker-compose-logs: ## View docker-compose logs
	@docker-compose logs -f

# Development targets
deps: ## Download and tidy dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies updated"

fmt: ## Format Go code
	@echo "Formatting code..."
	@go fmt ./...

lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	@golangci-lint run

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

# Utility targets
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME) $(BINARY_NAME)-* coverage.out coverage.html
	@docker rmi $(DOCKER_IMAGE):$(DOCKER_TAG) 2>/dev/null || true
	@echo "Clean complete"

install: build ## Install binary to GOPATH/bin
	@echo "Installing $(BINARY_NAME) to $(GOPATH)/bin..."
	@cp $(BINARY_NAME) $(GOPATH)/bin/
	@echo "Installation complete"

setup-oauth: ## Instructions for setting up OAuth credentials
	@echo "=== Gmail MCP Server OAuth Setup ==="
	@echo ""
	@echo "1. Go to Google Cloud Console: https://console.cloud.google.com/"
	@echo "2. Create a new project or select an existing one"
	@echo "3. Enable the Gmail API"
	@echo "4. Create credentials:"
	@echo "   - Go to 'Credentials' in the sidebar"
	@echo "   - Click 'Create Credentials' > 'OAuth client ID'"
	@echo "   - Choose 'Desktop application'"
	@echo "   - Download the JSON file"
	@echo "5. Rename the file to 'gcp-oauth.keys.json'"
	@echo "6. Place it in the current directory"
	@echo "7. Run 'make auth' to authenticate"
	@echo ""

check-oauth: ## Check if OAuth credentials are configured
	@if [ -f "gcp-oauth.keys.json" ]; then \
		echo "✓ OAuth credentials found: gcp-oauth.keys.json"; \
	else \
		echo "✗ OAuth credentials not found. Run 'make setup-oauth' for instructions."; \
		exit 1; \
	fi

health-check: ## Check if server is running
	@echo "Checking server health..."
	@curl -f http://localhost:8080/health || (echo "Server is not responding" && exit 1)
	@echo "Server is healthy!"

# Combined targets
setup: deps check-oauth ## Setup development environment
	@echo "Development environment setup complete!"

start: check-oauth run ## Check OAuth and start server

deploy: build-linux docker-build ## Build for production deployment