# Contexis CMP Framework Makefile

.PHONY: help build test clean install dev docs

# Default target
help:
	@echo "Contexis CMP Framework - Available Commands:"
	@echo ""
	@echo "Development:"
	@echo "  make build     - Build the Go CLI and Python packages"
	@echo "  make test      - Run all tests (Go and Python)"
	@echo "  make clean     - Clean build artifacts"
	@echo "  make install   - Install the CLI tool"
	@echo "  make dev       - Start development server"
	@echo ""
	@echo "Documentation:"
	@echo "  make docs      - Generate documentation"
	@echo "  make examples  - Build example applications"
	@echo ""
	@echo "Quality:"
	@echo "  make lint      - Run linting (Go and Python)"
	@echo "  make format    - Format code"
	@echo "  make security  - Run security checks"
	@echo ""
	@echo "Deployment:"
	@echo "  make docker    - Build Docker image"
	@echo "  make release   - Create release artifacts"

# Build the framework
build:
	@echo "Building Contexis CMP Framework..."
	go build -o bin/ctx src/cli/main.go
	@echo "✓ Go CLI built successfully"
	python -m pip install -e .
	@echo "✓ Python packages installed"

# Run tests
test:
	@echo "Running tests..."
	go test ./...
	python -m pytest tests/ -v
	@echo "✓ All tests passed"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf build/
	rm -rf dist/
	rm -rf *.egg-info/
	find . -type f -name "*.pyc" -delete
	find . -type d -name "__pycache__" -delete
	@echo "✓ Cleaned successfully"

# Install the CLI tool
install: build
	@echo "Installing CLI tool..."
	cp bin/ctx /usr/local/bin/
	@echo "✓ CLI installed to /usr/local/bin/ctx"

# Start development server
dev:
	@echo "Starting development server..."
	python -m contexis.dev.server

# Generate documentation
docs:
	@echo "Generating documentation..."
	mkdir -p docs/build
	sphinx-build -b html docs/ docs/build/html
	@echo "✓ Documentation generated"

# Build examples
examples:
	@echo "Building example applications..."
	cd examples/rag && ctx build
	cd examples/agent && ctx build
	cd examples/workflow && ctx build
	@echo "✓ Examples built"

# Run linting
lint:
	@echo "Running linting..."
	golangci-lint run
	black --check src/
	isort --check-only src/
	flake8 src/
	@echo "✓ Linting passed"

# Format code
format:
	@echo "Formatting code..."
	gofmt -w src/
	black src/
	isort src/
	@echo "✓ Code formatted"

# Security checks
security:
	@echo "Running security checks..."
	gosec ./...
	bandit -r src/
	@echo "✓ Security checks passed"

# Build Docker image
docker:
	@echo "Building Docker image..."
	docker build -t contexis/cmp:latest .
	@echo "✓ Docker image built"

# Create release
release:
	@echo "Creating release..."
	git tag -a v$(shell cat VERSION) -m "Release v$(shell cat VERSION)"
	git push origin v$(shell cat VERSION)
	@echo "✓ Release v$(shell cat VERSION) created"

# Development setup
setup:
	@echo "Setting up development environment..."
	python -m pip install -r requirements-dev.txt
	pre-commit install
	@echo "✓ Development environment ready"

# Quick validation
validate:
	@echo "Validating framework..."
	go vet ./...
	python -c "import contexis; print('✓ Python package valid')"
	@echo "✓ Framework validation passed" 