.PHONY: run build test clean help docker-build docker-run docker-dev lint format

help:
	@echo "🚀 Tuity - Twitter-like MVP"
	@echo "======================="
	@echo "Available commands:"
	@echo "  make run            - Start the server (port 8080)"
	@echo "  make build          - Build the application"
	@echo "  make test           - Run all tests"
	@echo "  make test-coverage  - Run tests with coverage"
	@echo "  make format         - Format code"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-run     - Run with Docker Compose"
	@echo "  make docker-dev     - Run development environment"
	@echo "  make docker-stop    - Stop Docker containers"
	@echo "  make help           - Show this help"

# Run the application
run:
	@echo "🚀 Starting Tuity MVP..."
	go run cmd/server/main.go

# Build the application
build:
	@echo "🔨 Building Tuity..."
	go build -o bin/tuity cmd/server/main.go
	@echo "✅ Build complete: bin/tuity"

# Run tests
test:
	@echo "🧪 Running tests..."
	go test ./... -v

# Test with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	go test ./... -v -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report: coverage.html"


# Format code
format:
	@echo "🎨 Formatting code..."
	go fmt ./...
	@if command -v goimports > /dev/null; then \
		goimports -w .; \
	else \
		echo "📝 goimports not found. Install with: go install golang.org/x/tools/cmd/goimports@latest"; \
	fi

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	docker-compose down --volumes --remove-orphans || true
	@echo "✅ Clean complete"

# Docker commands
docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t tuity:latest .

docker-run:
	@echo "🐳 Running with Docker Compose..."
	docker-compose up -d

docker-dev:
	@echo "🐳 Running development environment..."
	docker-compose --profile dev up -d tuity-dev

docker-stop:
	@echo "🐳 Stopping Docker containers..."
	docker-compose down

# Development mode (with hot reload if you have air installed)
dev:
	@if command -v air > /dev/null; then \
		echo "🔥 Starting development server with hot reload..."; \
		air; \
	else \
		echo "📝 Hot reload not available. Install with: go install github.com/cosmtrek/air@latest"; \
		make run; \
	fi

# Install development dependencies
deps:
	@echo "📦 Installing development dependencies..."
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest

# Quick setup for new developers
setup: deps
	@echo "🚀 Setting up Tuity development environment..."
	@echo "✅ Development environment ready!"
	@echo "Run 'make dev' to start with hot reload"
	@echo "Run 'make test' to run tests" 