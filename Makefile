.PHONY: help build test coverage clean lint fmt deps run env hooks mocks mock docs

# Binary name
BINARY_NAME=tasks

# Directories
CMD_DIR=./cmd/api
BIN_DIR=./out

# Default target
all: build

# Build the application
build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BINARY_NAME) $(CMD_DIR)

# Run tests with standard output
test:
	gotestsum --format pkgname -- ./...

# Generate test coverage report
coverage:
	mkdir -p ./out
	gotestsum --format pkgname -- -coverprofile=./out/coverage.out ./...
	go tool cover -html=./out/coverage.out -o ./out/coverage.html

# Run linter
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Clean build artifacts
clean:
	rm -rf $(BIN_DIR)
	go clean

# Install dependencies
deps:
	go mod tidy
	go mod verify

# Run the application
run: build
	$(BIN_DIR)/$(BINARY_NAME)

# Update environment
env:
	@echo "Updating direnv environment..."
	@direnv reload

# Install and setup git hooks with Lefthook
hooks:
	@echo "Installing git hooks with Lefthook..."
	lefthook install
	
# Generate mocks
mocks:
	@echo "Generating mocks..."
	go generate ./...

# Generate mock
mock:
	@echo "Generating mock for ${GOPACKAGE}/${GOFILE}..."

# Generate Swagger docs
docs:
	swag init -g cmd/api/main.go -o docs/swagger

# Help target
help:
	@echo "Available targets:"
	@echo "  build - Build the application"
	@echo "  test  - Run tests using gotestsum"
	@echo "  lint  - Run linter"
	@echo "  fmt   - Format code"
	@echo "  clean - Remove build artifacts"
	@echo "  deps  - Install dependencies"
	@echo "  run   - Build and run the application"
	@echo "  env   - Reload direnv environment"
	@echo "  hooks - Install git hooks with Lefthook"
	@echo "  mocks - Generate mocks"
	@echo "  docs  - Generate Swagger docs"
	@echo "  coverage - Generate test coverage report"
