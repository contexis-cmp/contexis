# Contexis CMP Framework Makefile

.PHONY: help build test clean install install-local dev docs

# Default target
help:
	@echo "Contexis CMP Framework - Available Commands:"
	@echo ""
	@echo "Development:"
	@echo "  make build     - Build the Go CLI and Python packages"
	@echo "  make test      - Run all tests (Go and Python)"
	@echo "  make clean     - Clean build artifacts"
	@echo "  make install   - Install the CLI tool (system-wide if possible, local if not)"
	@echo "  make install-local - Install the CLI tool to user's local directory"
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
	@echo " Go CLI built successfully"
	python -m pip install -e .
	@echo " Python packages installed"

# Run tests
test:
	@echo "Running all tests..."
	go test ./tests/unit/... -v
	go test ./tests/integration/... -v
	@if [ -n "$$($(shell which go) list ./tests/e2e/... 2>/dev/null)" ]; then \
		go test ./tests/e2e/... -v; \
	else \
		echo "Skipping E2E tests: no packages"; \
	fi
	@echo " All tests passed"

# Run unit tests only
test-unit:
	@echo "Running unit tests..."
	go test ./tests/unit/... -v -coverprofile=tests/coverage/unit.out
	@echo " Unit tests passed"

# Run integration tests only
test-integration:
	@echo "Running integration tests..."
	go test ./tests/integration/... -v -coverprofile=tests/coverage/integration.out
	@echo " Integration tests passed"

# Run e2e tests only
test-e2e:
	@echo "Running end-to-end tests..."
	go test ./tests/e2e/... -v -coverprofile=tests/coverage/e2e.out
	@echo " E2E tests passed"

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	mkdir -p tests/coverage
	go test ./tests/... -v -coverprofile=tests/coverage/all.out -covermode=atomic
	go tool cover -html=tests/coverage/all.out -o tests/coverage/coverage.html
	@echo " Coverage report generated at tests/coverage/coverage.html"

# Run performance tests
test-performance:
	@echo "Running performance tests..."
	go test ./tests/... -v -tags=performance
	@echo " Performance tests passed"

# Run security tests
test-security:
	@echo "Running security tests..."
	go test ./tests/... -v -tags=security
	@echo " Security tests passed"

# Run specific test category
test-category:
	@echo "Running $(CATEGORY) tests..."
	go test ./tests/... -v -run $(CATEGORY)
	@echo " $(CATEGORY) tests passed"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf build/
	rm -rf dist/
	rm -rf *.egg-info/
	find . -type f -name "*.pyc" -delete
	find . -type d -name "__pycache__" -delete
	@echo " Cleaned successfully"

# Install the CLI tool (system-wide if possible, local if not)
install: build
	@echo "Installing CLI tool..."
	@if [ -w /usr/local/bin ]; then \
		cp bin/ctx /usr/local/bin/; \
		echo " CLI installed to /usr/local/bin/ctx"; \
	else \
		echo "Installing to user's local bin directory..."; \
		mkdir -p $(HOME)/.local/bin; \
		cp bin/ctx $(HOME)/.local/bin/; \
		echo " CLI installed to $(HOME)/.local/bin/ctx"; \
		echo ""; \
		echo "To use the CLI, add the following to your shell profile:"; \
		echo "export PATH=\"$(HOME)/.local/bin:\$$PATH\""; \
		echo ""; \
		echo "For bash/zsh, add to ~/.bashrc or ~/.zshrc:"; \
		echo "export PATH=\"$(HOME)/.local/bin:\$$PATH\""; \
		echo ""; \
		echo "For fish, add to ~/.config/fish/config.fish:"; \
		echo "set -gx PATH $(HOME)/.local/bin \$PATH"; \
	fi

# Install the CLI tool to user's local directory only
install-local: build
	@echo "Installing CLI tool to user's local directory..."
	mkdir -p $(HOME)/.local/bin
	cp bin/ctx $(HOME)/.local/bin/
	@echo " CLI installed to $(HOME)/.local/bin/ctx"
	@echo ""
	@echo "To use the CLI, add the following to your shell profile:"
	@echo "export PATH=\"$(HOME)/.local/bin:\$$PATH\""
	@echo ""
	@echo "For bash/zsh, add to ~/.bashrc or ~/.zshrc:"
	@echo "export PATH=\"$(HOME)/.local/bin:\$$PATH\""
	@echo ""
	@echo "For fish, add to ~/.config/fish/config.fish:"
	@echo "set -gx PATH $(HOME)/.local/bin \$PATH"

# Start development server
dev:
	@echo "Starting development server..."
	python -m contexis.dev.server

# Generate documentation
docs:
	@echo "Generating documentation..."
	mkdir -p docs/build
	sphinx-build -b html docs/ docs/build/html
	@echo " Documentation generated"

# Build examples
examples:
	@echo "Building example applications..."
	cd examples/rag && ctx build
	cd examples/agent && ctx build
	cd examples/workflow && ctx build
	@echo " Examples built"

# Run linting
lint:
	@echo "Running linting..."
	golangci-lint run
	black --check src/
	isort --check-only src/
	flake8 src/
	@echo " Linting passed"

# Format code
format:
	@echo "Formatting code..."
	gofmt -w src/
	black src/
	isort src/
	@echo " Code formatted"

# Security checks
security:
	@echo "Running security checks..."
	gosec ./...
	bandit -r src/
	@echo " Security checks passed"

# Build Docker image
docker:
	@echo "Building Docker image..."
	docker build -t contexis-cmp/contexis:latest .
	@echo " Docker image built"

# Create release
release:
	@echo "Creating release..."
	git tag -a v$(shell cat VERSION) -m "Release v$(shell cat VERSION)"
	git push origin v$(shell cat VERSION)
	@echo " Release v$(shell cat VERSION) created"

# Development setup
setup:
	@echo "Setting up development environment..."
	python -m pip install -r requirements-dev.txt
	pre-commit install
	@echo " Development environment ready"

# Quick validation
validate:
	@echo "Validating framework..."
	go vet ./...
	python -c "import contexis; print(' Python package valid')"
	@echo " Framework validation passed" 