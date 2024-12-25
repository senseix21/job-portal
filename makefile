# Define the Go binary name
BINARY_NAME=job-service

# Define the Go source files
GO_FILES=$(shell find . -name '*.go')

# Default Go options
GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_TEST=$(GO_CMD) test
GO_RUN=$(GO_CMD) run

# Directories
BIN_DIR=bin
BUILD_DIR=build
VENDOR_DIR=vendor

# Targets

# Build the application
build:
	@echo "Building the Go application..."
	$(GO_BUILD) -o $(BIN_DIR)/$(BINARY_NAME)

# Run the application
run: build
	@echo "Running the Go application..."
	$(BIN_DIR)/$(BINARY_NAME)

# Clean the generated files
clean:
	@echo "Cleaning up..."
	rm -rf $(BIN_DIR) $(BUILD_DIR)

# Test the Go application
test:
	@echo "Running tests..."
	$(GO_TEST) -v ./...

# Format Go code using `gofmt`
fmt:
	@echo "Formatting Go code..."
	$(GO_CMD) fmt $(GO_FILES)

# Lint Go code using `golint`
lint:
	@echo "Linting Go code..."
	@golint $(GO_FILES)

# Install dependencies
install-deps:
	@echo "Installing dependencies..."
	$(GO_CMD) mod tidy

# Run tests and build before committing
precommit: test fmt lint

# Run the application in Docker (if using Docker)
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME) .

docker-run: docker-build
	@echo "Running Docker container..."
	docker run -p 8080:8080 $(BINARY_NAME)

# Help command to show available targets
help:
	@echo "Makefile commands:"
	@echo "  build         - Build the Go application"
	@echo "  run           - Run the Go application"
	@echo "  clean         - Clean the build artifacts"
	@echo "  test          - Run the tests"
	@echo "  fmt           - Format the Go code"
	@echo "  lint          - Lint the Go code"
	@echo "  install-deps  - Install dependencies"
	@echo "  precommit     - Run tests, format, and lint before commit"
	@echo "  docker-build  - Build the Docker image"
	@echo "  docker-run    - Run the application inside Docker"
	@echo "  help          - Show this help message"
