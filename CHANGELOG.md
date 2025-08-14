# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added
- **New `ctx run` command**: Added direct query execution without manually starting the server
  - Automatically starts server, sends query, and returns response
  - Supports all HTTP API features (memory search, prompt rendering, model inference)
  - Includes debug mode and timeout configuration
  - Examples: `ctx run SupportBot "What is your return policy?"`
  - Flags: `--addr`, `--component`, `--data`, `--debug`, `--timeout`, `--top-k`
- **Colored Logging System**: Rails-like colored console output for better readability
  - **Color-coded levels**: Green INFO, Yellow WARN, Red ERROR, Cyan DEBUG
  - **Status indicators**: ✓ Success, ⚠ Warning, ✗ Error, ▶ Start, ℹ Info, 🔍 Debug
  - **Structured output**: Colored timestamps, field names, and values
  - **Smart detection**: Automatically disables colors in CI environments or when redirected to files
  - **Industry standards**: Follows Rails backtrace color conventions for familiarity

### Fixed
- **CLI Command Duplication**: Removed duplicate `generate` command from help output
  - Removed unused local `generateCmd` variable that was conflicting with proper `commands.GenerateCmd`
  - Cleaned up main.go file structure
- **Documentation Updates**: Fixed all references to non-existent `ctx run` command
  - Updated `docs/guides/getting-started.md` with correct command syntax
  - Updated `docs/cli.md` to include run command examples
  - Updated `docs/runtime.md` with proper workflow instructions
  - Updated all example READMEs in `examples/` directory:
    - `examples/workflow/README.md`: Fixed workflow execution commands
    - `examples/rag/README.md`: Fixed RAG query commands
    - `examples/agent/README.md`: Fixed conversation commands
  - Updated `contexts/CustomerDocs/README.md` (already had correct format)

### Changed
- **Improved Developer Experience**: Users can now execute queries directly without manual server management
- **Better CLI Workflow**: Simplified workflow from `ctx serve` + HTTP requests to single `ctx run` command
- **Consistent Documentation**: All examples now use the correct, working command syntax

## [0.1.14] - 2024-08-15

### Added
- Initial release of Contexis CMP Framework
- Context-Memory-Prompt (CMP) architecture implementation
- CLI tool with component generation, testing, and deployment
- HTTP server with chat API
- Memory management with vector stores
- Prompt templating system
- Security features (authentication, rate limiting, audit logging)
- Hugging Face integration
- Plugin system
- Comprehensive testing framework

### Features
- Agent generation with customizable tools and memory
- RAG system generation with vector databases
- Workflow generation for data processing pipelines
- Multi-tenant support
- Drift detection and monitoring
- Kubernetes deployment support
- Docker containerization
- Prometheus metrics and OpenTelemetry tracing
