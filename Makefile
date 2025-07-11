# qomoboro Makefile
# Canonical Hours Task Manager

# Variables
BINARY_NAME = qomoboro
VERSION = 2.0.0
BUILD_DIR = build
INSTALL_DIR = /usr/local/bin
DATA_DIR = ~/.local/share/qomoboro

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod
GOFMT = $(GOCMD) fmt
GOVET = $(GOCMD) vet

# Build flags
LDFLAGS = -ldflags "-X main.version=$(VERSION)"
BUILD_FLAGS = -v

# Default target
.PHONY: all
all: clean format test build

# Build the binary
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BINARY_NAME) .

# Build for release (optimized)
.PHONY: build-release
build-release:
	@echo "Building $(BINARY_NAME) for release..."
	$(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BINARY_NAME) -ldflags "-s -w" .

# Build for multiple platforms
.PHONY: build-all
build-all: clean
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .

# Run the application
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

# Run with development data directory
.PHONY: run-dev
run-dev: build
	@echo "Running $(BINARY_NAME) with development data..."
	XDG_DATA_HOME=./dev-data ./$(BINARY_NAME)

# Quick development cycle
.PHONY: quick
quick: format build run

# Test the application
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Test with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Run tests with race detection
.PHONY: test-race
test-race:
	@echo "Running tests with race detection..."
	$(GOTEST) -v -race ./...

# Format code
.PHONY: format
format:
	@echo "Formatting code..."
	$(GOFMT) ./...

# Lint code
.PHONY: lint
lint:
	@echo "Linting code..."
	$(GOVET) ./...

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	$(GOMOD) tidy
	$(GOMOD) download

# Update dependencies
.PHONY: deps-update
deps-update:
	@echo "Updating dependencies..."
	$(GOMOD) tidy
	$(GOGET) -u ./...

# Install binary to system
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	sudo cp $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	sudo chmod +x $(INSTALL_DIR)/$(BINARY_NAME)

# Uninstall binary from system
.PHONY: uninstall
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	sudo rm -f $(INSTALL_DIR)/$(BINARY_NAME)

# Development setup
.PHONY: setup-dev
setup-dev:
	@echo "Setting up development environment..."
	$(GOMOD) tidy
	@mkdir -p dev-data
	@echo "Development environment ready!"

# Create sample data for development
.PHONY: sample-data
sample-data:
	@echo "Creating sample data..."
	@mkdir -p dev-data/qomoboro
	@echo "Sample data created in dev-data/"

# Backup user data
.PHONY: backup
backup:
	@echo "Creating backup..."
	./$(BINARY_NAME) --backup || (echo "Building first..." && $(MAKE) build && ./$(BINARY_NAME) --backup)

# Show data directory
.PHONY: data-dir
data-dir:
	@echo "Data directory:"
	./$(BINARY_NAME) --data-dir || (echo "Building first..." && $(MAKE) build && ./$(BINARY_NAME) --data-dir)

# Complete CI pipeline
.PHONY: ci
ci: clean format lint test build
	@echo "CI pipeline completed successfully!"

# Development workflow
.PHONY: dev
dev: format lint test build run-dev

# Show help
.PHONY: help
help:
	@echo "qomoboro Makefile - Canonical Hours Task Manager"
	@echo ""
	@echo "Available targets:"
	@echo "  all           - Clean, format, test, and build"
	@echo "  build         - Build the binary"
	@echo "  build-release - Build optimized release binary"
	@echo "  build-all     - Build for multiple platforms"
	@echo "  run           - Build and run the application"
	@echo "  run-dev       - Run with development data directory"
	@echo "  quick         - Quick development cycle (format, build, run)"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  test-race     - Run tests with race detection"
	@echo "  format        - Format code"
	@echo "  lint          - Lint code"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Install dependencies"
	@echo "  deps-update   - Update dependencies"
	@echo "  install       - Install binary to system"
	@echo "  uninstall     - Uninstall binary from system"
	@echo "  setup-dev     - Setup development environment"
	@echo "  sample-data   - Create sample data for development"
	@echo "  backup        - Create backup of user data"
	@echo "  data-dir      - Show data directory location"
	@echo "  ci            - Complete CI pipeline"
	@echo "  dev           - Development workflow"
	@echo "  help          - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make quick    - Quick development cycle"
	@echo "  make ci       - Run complete CI pipeline"
	@echo "  make dev      - Development workflow with live data"
	@echo "  make install  - Install to system"

# Check if binary exists
.PHONY: check-binary
check-binary:
	@if [ ! -f $(BINARY_NAME) ]; then \
		echo "Binary not found, building..."; \
		$(MAKE) build; \
	fi
