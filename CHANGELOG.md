# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### **Phase 3: Production Migration Complete**

**Major Feature: Seamless Environment Migration**

- **Migration Commands**: Added `ctx migrate` command for seamless transitions between environments
- **Production Templates**: Generated production configurations for OpenAI, Anthropic, and Hugging Face
- **Deployment Files**: Created Docker and Kubernetes deployment templates
- **Validation Tools**: Added configuration validation and connectivity testing
- **Documentation**: Comprehensive migration guides and local development documentation

**Migration Features:**
- **Local to Production**: `ctx migrate local-to-production --provider=openai|anthropic|huggingface`
- **Production to Local**: `ctx migrate production-to-local` for rollback
- **Configuration Validation**: `ctx migrate validate` for testing setup
- **Dry Run Mode**: `--dry-run` flag to preview changes without applying them

**Generated Files:**
- Production configuration templates (`config/environments/production.yaml`)
- Environment variable templates (`.env.production`)
- Docker deployment files (`Dockerfile`, `docker-compose.yml`)
- Migration guides (`MIGRATION.md`, `LOCAL_DEVELOPMENT.md`)

**Developer Experience:**
- **One-Command Migration**: Single command to switch from local to production
- **Provider Flexibility**: Support for multiple production providers
- **Rollback Capability**: Easy return to local development
- **Validation**: Built-in configuration and connectivity testing

---

### **Phase 2: Local Model Integration Complete**

**Major Feature: Local-First Development Experience**

- **Local Model Integration**: Successfully integrated Phi-3.5-Mini and Sentence Transformers for out-of-the-box local development
- **Robust Path Resolution**: Fixed local provider path resolution to work from both repo root and generated project directories
- **Model Compatibility**: Resolved transformers library compatibility issues with DynamicCache and flash-attention
- **Template System**: Fixed Go template syntax in RAG prompts to use proper dot notation
- **Complete Happy Path**: End-to-end local development workflow now working seamlessly

**Technical Achievements:**
- **Local AI Provider**: Go-side interface for Python local model execution with proper subprocess management
- **Local Embeddings**: Integrated Sentence Transformers for local embeddings with Chroma(SQLite)
- **Model Warmup**: Added `ctx models warmup` command for pre-downloading local models
- **Enhanced Logging**: Improved error handling and progress tracking for local model operations
- **Production Migration**: Clear path from local development to production providers

**Key Fixes:**
- Fixed `DynamicCache` compatibility issues with transformers 4.55.0
- Resolved template variable syntax (`{{.user_query}}` vs `{{user_query}}`)
- Fixed path resolution for local provider when running from sub-project directories
- Improved timeout handling for initial model downloads
- Enhanced error messages and debugging information

**Developer Experience:**
- **No External Dependencies**: Start developing immediately without API keys
- **Automatic Model Management**: Models download on first use with progress tracking
- **Seamless Migration**: Easy switch to production providers when ready
- **Comprehensive Testing**: All core functionality tested and verified

**Breaking Changes:**
- Removed `--dev-only` flag from `ctx init` - local models are now the default
- Updated template syntax to use Go template dot notation

---

## [0.1.14] - 2024-01-15

### Added
- **Dev-Only Package Implementation**: Added `--dev-only` flag to `ctx init` for local-first development
- **Local Model Support**: Integrated Phi-3.5-Mini and Sentence Transformers for local development
- **Model Warmup Command**: Added `ctx models warmup` for pre-downloading local models
- **Enhanced Project Generation**: Improved `ctx init` output with detailed file structure and development flow
- **Local Provider Integration**: Go-side interface for Python local model execution

### Changed
- **Default Development**: Made local models the default for new projects
- **Configuration Files**: Updated generated `development.yaml`, `requirements.txt`, and `.env.example` for local-first setup
- **Documentation**: Updated README and guides to reflect local-first development approach

### Fixed
- **Template Generation**: Fixed RAG template generation to create `agent_response.md` properly
- **Path Resolution**: Improved local provider path resolution for different execution contexts
- **Model Compatibility**: Resolved transformers library compatibility issues

---

## [0.1.13] - 2024-01-14

### Added
- **Enhanced Test Output**: Improved `ctx test` command with detailed per-component results and colored logging
- **Drift Detection Logging**: Enhanced drift detection output with detailed test results and similarity scores
- **Test Configuration**: Added `test_config.yaml` generation in `ctx init` with comprehensive test setup

### Changed
- **Test Runner**: Enhanced Go test runner to handle "no packages to test" gracefully
- **Drift Detection**: Improved drift detection output with detailed per-component results and summary
- **Project Structure**: Added test subdirectories (`unit`, `integration`, `e2e`, `coverage`, `reports`, `fixtures`) in `ctx init`

### Fixed
- **Test Command**: Fixed `ctx test` command to handle empty test suites properly
- **Test Logging**: Enhanced test output with detailed error messages and status indicators

---

## [0.1.12] - 2024-01-13

### Added
- **RAG Template Optimization**: Made `agent_response.md` a symlink to `search_response.md` to avoid redundancy
- **Enhanced RAG Generation**: Improved RAG generator logging to show all generated files

### Changed
- **Template Structure**: Optimized RAG template generation to use symlinks for shared content
- **Logging**: Enhanced RAG generator output to explicitly show generated files

### Fixed
- **Template Generation**: Fixed RAG template generation to create both `search_response.md` and `agent_response.md`
- **Symlink Paths**: Fixed symlink creation to use relative paths for portability

---

## [0.1.11] - 2024-01-12

### Added
- **RAG Template Generation**: Added `agent_response.md` generation in RAG prompts
- **Template Variable Support**: Added support for template variables in RAG response templates

### Fixed
- **Template Parsing**: Fixed template variable parsing issues in RAG response generation
- **Static Templates**: Implemented static template generation to avoid parsing errors

---

## [0.1.10] - 2024-01-11

### Added
- **Enhanced Logging System**: Implemented Rails-like colored logging with status indicators
- **Colored Logger**: Added `src/cli/logger/colored_logger.go` with ANSI color support
- **Enhanced Commands**: Updated all CLI commands to use new colored logging system

### Changed
- **Logging Format**: Switched from JSON to colored console logging by default
- **Command Output**: Enhanced all command outputs with detailed progress tracking and status indicators
- **Project Generation**: Improved `ctx init` output with detailed file structure and development flow
- **Component Generation**: Enhanced `ctx generate` commands with detailed output and development guidance

### Fixed
- **Duplicate Generate Command**: Fixed duplicate `generate` command in `ctx --help` output
- **Logging Integration**: Successfully integrated colored logging across all CLI commands

---

## [0.1.9] - 2024-01-10

### Added
- **Run Command**: Added `ctx run` command for direct query execution without server setup
- **Enhanced Documentation**: Updated all example READMEs to use correct `ctx run` command
- **CLI Help**: Updated help documentation to clarify correct workflow

### Changed
- **Command Structure**: Integrated `ctx run` command into main CLI structure
- **Documentation**: Updated all documentation to reflect new `ctx run` command usage

### Fixed
- **Missing Command**: Added the missing `ctx run` command for better user experience
- **Documentation**: Fixed all example READMEs to use correct commands

---

## [0.1.8] - 2024-01-09

### Added
- **Plugin System**: Added plugin registry and management capabilities
- **Plugin Commands**: Added `ctx plugin` commands for plugin management
- **Plugin Scaffolding**: Added plugin generation capabilities

### Changed
- **Plugin Architecture**: Enhanced plugin system with better discovery and lifecycle management
- **Plugin Documentation**: Updated documentation to include plugin system

---

## [0.1.7] - 2024-01-08

### Added
- **Security Controls**: Added comprehensive security features including API key auth, RBAC, and rate limiting
- **Audit Logging**: Added audit trail for security events
- **Prompt Injection Guard**: Added protection against prompt injection attacks
- **PII Handling**: Added PII detection and handling capabilities

### Changed
- **Security Architecture**: Enhanced security architecture with optional controls
- **Documentation**: Updated security documentation and examples

---

## [0.1.6] - 2024-01-07

### Added
- **Hugging Face Integration**: Added support for Hugging Face Inference API
- **Model Provider Abstraction**: Added model provider interface for different AI providers
- **HF Commands**: Added `ctx hf` commands for Hugging Face model testing

### Changed
- **Model Architecture**: Enhanced model architecture to support multiple providers
- **Documentation**: Added Hugging Face integration documentation

---

## [0.1.5] - 2024-01-06

### Added
- **Memory System**: Added comprehensive memory management with vector stores
- **Memory Commands**: Added `ctx memory` commands for memory operations
- **Vector Database Support**: Added support for ChromaDB and other vector databases

### Changed
- **Memory Architecture**: Enhanced memory system with better provider support
- **Documentation**: Added memory system documentation

---

## [0.1.4] - 2024-01-05

### Added
- **Prompt Engine**: Added comprehensive prompt templating system
- **Prompt Commands**: Added `ctx prompt` commands for prompt management
- **Template Validation**: Added template validation and linting

### Changed
- **Prompt Architecture**: Enhanced prompt system with better templating capabilities
- **Documentation**: Added prompt system documentation

---

## [0.1.3] - 2024-01-04

### Added
- **Context System**: Added comprehensive context management system
- **Context Commands**: Added `ctx context` commands for context operations
- **Schema Validation**: Added context schema validation

### Changed
- **Context Architecture**: Enhanced context system with better validation
- **Documentation**: Added context system documentation

---

## [0.1.2] - 2024-01-03

### Added
- **Generator System**: Added component generators for RAG, agents, and workflows
- **Generator Commands**: Added `ctx generate` commands for component generation
- **Project Scaffolding**: Added project scaffolding capabilities

### Changed
- **Generator Architecture**: Enhanced generator system with better templates
- **Documentation**: Added generator system documentation

---

## [0.1.1] - 2024-01-02

### Added
- **CLI Framework**: Added comprehensive CLI framework using Cobra
- **Basic Commands**: Added basic CLI commands for project management
- **Project Initialization**: Added `ctx init` command for project setup

### Changed
- **CLI Architecture**: Enhanced CLI architecture with better command structure
- **Documentation**: Added CLI documentation

---

## [0.1.0] - 2024-01-01

### Added
- **Initial Release**: First release of Contexis CMP Framework
- **Core Architecture**: Implemented CMP (Context-Memory-Prompt) architecture
- **Basic Framework**: Added basic framework structure and components
- **Documentation**: Added initial documentation and examples

### Changed
- **Project Structure**: Established initial project structure
- **Architecture**: Defined core CMP architecture principles

---
