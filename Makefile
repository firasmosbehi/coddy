 # Coddy Makefile

.PHONY: build build-cli build-server test lint fmt clean run docker

# Build variables
BINARY_NAME=coddy
SERVER_BINARY=coddy-server
BUILD_DIR=./build

# Go variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-s -w"

# Default target
all: build

# Build both binaries
build: build-cli build-server

# Build CLI binary
build-cli:
	@echo "Building Coddy CLI..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/coddy

# Build server binary
build-server:
	@echo "Building Coddy Server..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(SERVER_BINARY) ./cmd/server

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -race ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Lint code
lint:
	@echo "Linting..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, using go vet..."; \
		$(GOCMD) vet ./...; \
	fi

# Format code
fmt:
	@echo "Formatting..."
	$(GOCMD) fmt ./...

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	$(GOMOD) tidy
	$(GOMOD) download

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

# Run the CLI
run: build-cli
	@echo "Running Coddy CLI..."
	@$(BUILD_DIR)/$(BINARY_NAME)

# Run the server
run-server: build-server
	@echo "Running Coddy Server..."
	@$(BUILD_DIR)/$(SERVER_BINARY)

# Build Docker sandbox image
docker:
	@echo "Building Docker sandbox image..."
	docker build -t coddy-sandbox:latest -f docker/Dockerfile .

# Install binaries to GOPATH/bin
install: build
	@echo "Installing binaries..."
	$(GOCMD) install ./cmd/coddy
	$(GOCMD) install ./cmd/server

# Development mode - watch and rebuild
watch:
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "Installing air..."; \
		$(GOGET) -u github.com/cosmtrek/air; \
		air; \
	fi

# Help
help:
	@echo "Available targets:"
	@echo "  build          - Build CLI and server binaries"
	@echo "  build-cli      - Build CLI binary only"
	@echo "  build-server   - Build server binary only"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  lint           - Run linter"
	@echo "  fmt            - Format code"
	@echo "  tidy           - Tidy and download dependencies"
	@echo "  clean          - Clean build artifacts"
	@echo "  run            - Build and run CLI"
	@echo "  run-server     - Build and run server"
	@echo "  docker         - Build Docker sandbox image"
	@echo "  install        - Install binaries to GOPATH/bin"
	@echo "  watch          - Watch and rebuild (requires air)"
