# Define variables
BINARY_NAME = netlinx-language-server
BUILD_DIR = build
VERSION = 0.1.0
LDFLAGS = -ldflags "-X main.Version=$(VERSION)"

.PHONY: build clean test test-cover test-report run fmt lint install

# Default target
all: fmt lint test build

# Build the application
build:
	@echo "Building $(BINARY_NAME) v$(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/netlinx-language-server

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)

# Run tests
test:
	@echo "Running tests..."
	@go test -v -count=1 ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -v -cover ./...

test-report:
	@echo "Generating test report..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

fmt:
	@echo "Formatting code..."
	@go fmt ./...

lint:
	@echo "Linting code..."
	@go vet ./..

install: build
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/

# Run the binary
run: build
	@echo "Running $(BINARY_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME)
