.PHONY: build test clean lint fmt deps run update-env install-hooks hooks help

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

# Run tests
test:
	go test ./...

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
update-env:
	@echo "Updating direnv environment..."
	@direnv reload

# Install and setup git hooks with Lefthook
install-hooks:
	@echo "Installing git hooks with Lefthook..."
	lefthook install

# Run hooks manually
hooks:
	@echo "Running git hooks manually..."
	lefthook run pre-commit

# Help target
help:
	@echo "Available targets:"
	@echo "  build      - Build the application"
	@echo "  test       - Run tests"
	@echo "  lint       - Run linter"
	@echo "  fmt        - Format code"
	@echo "  clean      - Remove build artifacts"
	@echo "  deps       - Install dependencies"
	@echo "  run        - Build and run the application"
	@echo "  update-env - Reload direnv environment"
	@echo "  install-hooks - Install git hooks with Lefthook"
	@echo "  hooks      - Run git hooks manually"
