.PHONY: help build test clean lint fmt deps run env hooks

# Binary name
BINARY_NAME=task-master

# Directories
CMD_DIR=./cmd/cli
BIN_DIR=./bin

# Default target
all: build

# Build the application
build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BINARY_NAME) $(CMD_DIR)

# Run tests with standard output
test:
	gotestsum -- ./...

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

# Run hooks manually
hooks:
	@echo "Running git hooks manually..."
	lefthook run pre-commit
	
# Generate mocks
generate-mocks:
	@echo "Generating mocks..."
	go generate ./pkg/...

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
