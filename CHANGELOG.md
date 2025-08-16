# Changelog

All notable changes to Contexis CMP Framework will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2025-01-16

### 🚀 Added

#### **Local-First Development Experience**
- **Local Models Integration**: Out-of-the-box support for Phi-3.5-Mini (text generation) and Sentence Transformers (embeddings)
- **Local Vector Database**: Chroma with SQLite backend for development without external dependencies
- **Model Management**: `ctx models warmup` command for pre-downloading local models
- **Zero External Dependencies**: Start developing immediately without API keys or cloud services

#### **Production Migration System**
- **Migration Commands**: 
  - `ctx migrate local-to-production --provider=openai|anthropic|huggingface` for seamless cloud migration
  - `ctx migrate production-to-local` for rollback to local development
  - `ctx migrate validate` for configuration testing
- **Production Templates**: Auto-generated configurations for OpenAI, Anthropic, and Hugging Face
- **Deployment Files**: Generated Docker and Kubernetes deployment templates
- **Environment Validation**: Built-in connectivity and configuration testing

#### **Enhanced Developer Experience**
- **Rails-like Colored Logging**: Beautiful console output with status indicators and progress tracking
- **Enhanced CLI Commands**:
  - New `ctx run` command for direct query execution without server setup
  - Improved `ctx test` with detailed per-component results and colored output
  - Enhanced `ctx generate` with better templates and guided workflows
- **Comprehensive Testing**: Enhanced drift detection, correctness testing, and performance monitoring
- **Guided Workflows**: Step-by-step instructions for every command

### 🔧 Changed

#### **Default Development Workflow**
- **Local-First by Default**: New projects now start with local models instead of requiring external API keys
- **Enhanced Project Generation**: Detailed file structure output and development guidance
- **Updated Documentation**: All examples and guides updated to reflect local-first approach

#### **Template System Improvements**
- **Go Template Syntax**: Fixed template variable syntax to use proper dot notation (`{{.user_query}}`)
- **RAG Template Generation**: Improved to create proper `agent_response.md` files
- **Enhanced Symlink Handling**: Better template generation with relative paths

#### **Runtime Improvements**
- **Local Provider Integration**: Go-side interface for Python local model execution
- **Enhanced Path Resolution**: Fixed local provider path resolution across execution contexts
- **Model Compatibility**: Resolved transformers library compatibility issues

### 🐛 Fixed

#### **Template & Generation Fixes**
- Fixed Go template syntax in RAG prompts (proper dot notation)
- Resolved RAG template generation to create both `search_response.md` and `agent_response.md`
- Fixed path resolution for local provider execution from different directories
- Enhanced symlink creation to use relative paths for portability

#### **Model & Runtime Fixes**
- Resolved `DynamicCache` compatibility issues with transformers 4.55.0
- Fixed timeout handling for initial model downloads
- Improved error messages and debugging information for local models
- Enhanced subprocess management for Python local model execution

#### **CLI & Testing Fixes**
- Fixed duplicate `generate` command in `ctx --help` output
- Enhanced test runner to handle empty test suites gracefully
- Improved drift detection output with detailed results and summaries
- Fixed `ctx test` command logging and error handling

### 💥 Breaking Changes

- **Removed `--dev-only` flag**: Local models are now the default for new projects
- **Updated template syntax**: Templates now use Go template dot notation consistently
- **Default development flow**: New projects start with local models by default

### 📚 Documentation

- **Updated README**: Reflects new local-first approach with complete quick start guide
- **Migration Guides**: Comprehensive documentation for environment transitions (`MIGRATION.md`, `LOCAL_DEVELOPMENT.md`)
- **Enhanced Examples**: All examples updated to use correct `ctx run` command syntax
- **Local Development Guides**: Complete documentation for local-first workflow

---

## [0.1.14] - 2024-01-15

### 🎉 Initial Release

This is the first public release of Contexis, a Rails-inspired framework for building reproducible AI applications using the Context–Memory–Prompt (CMP) architecture.

#### **Core Framework**
- **CMP Architecture**: Implemented Context-Memory-Prompt pattern as first-class, versioned artifacts
- **CLI Tool (`ctx`)**: Complete command-line interface with:
  - `ctx init` - Project initialization
  - `ctx generate` - Component generation (agent, RAG, workflow)
  - `ctx test` - Testing with drift detection (unit/integration/e2e)
  - `ctx serve` - Development server
  - `ctx build` - Build and deployment
  - Context/memory/prompt management utilities

#### **Runtime & Server**
- **HTTP Runtime Server**: Production-ready server with:
  - Health endpoints (`/healthz`, `/readyz`)
  - Version and metrics (`/version`, `/metrics` with Prometheus support)
  - Chat API (`/api/v1/chat`)
  - Optional Hugging Face inference integration

#### **Domain Model**
- **Versioned Components**: 
  - Context engine with schema validation
  - Memory management with vector stores
  - Prompt templating system
  - Deterministic SHA locking for reproducibility
- **Guardrails**: Format validation, token limits, and behavior constraints

#### **Security & Enterprise Features**
- **Authentication & Authorization**:
  - API key authentication
  - Role-Based Access Control (RBAC)
  - Per-key/tenant/IP rate limiting
- **Safety Features**:
  - Prompt injection detection
  - PII redaction and blocking
  - Source-constrained answering with required citations

#### **Deployment & Operations**
- **Kubernetes-Ready**: 
  - Helm chart (`charts/contexis-app`) with hardened defaults
  - Ingress/TLS configuration
  - Optional ExternalSecret integration
  - Argo Rollouts support
- **Monitoring**: Prometheus metrics and structured logging

#### **Plugin System (Alpha)**
- **Plugin Framework**: 
  - `ctx generate plugin` for scaffolding
  - Local and remote plugin installation
  - Capability registry system
  - Example plugin templates

#### **Testing & Quality**
- **Comprehensive Testing**:
  - Drift detection for AI behavior consistency
  - Correctness validation
  - Consolidated test reports
  - Business rule validation

### 🚀 Quick Start

#### **Installation**
```bash
# Download prebuilt binary and add to PATH
ctx version
# Expected: Contexis CMP Framework v0.1.14

# Or build from source (Go 1.21+)
git clone https://github.com/contexis-cmp/contexis.git
cd contexis
make install
```

#### **Create Your First AI App**
```bash
mkdir my-ai-app && cd my-ai-app
ctx init my-support-bot
cd my-support-bot
ctx generate agent SupportBot --tools=web_search,database --memory=episodic
ctx test --all
ctx serve
```

#### **Kubernetes Deployment**
```bash
helm install contexis charts/contexis-app \
  --set image.repository=contexis-cmp/contexis \
  --set image.tag=0.1.14
```

### 🔧 Configuration

#### **Security (Optional)**
```bash
export CMP_AUTH_ENABLED=true
export CMP_PI_ENFORCEMENT=true
export CMP_REQUIRE_CITATION=true
export CMP_API_TOKENS="devtoken@tenantA:chat:execute|context:read"
```

#### **Hugging Face Integration (Optional)**
```bash
export HF_TOKEN=your_token
export HF_MODEL_ID=microsoft/DialoGPT-medium
```

### 📋 Requirements

- **Go**: 1.21+ (for building from source)
- **Git**: Required for remote plugin installations
- **Optional**: Hugging Face account for external model inference

### 🔗 Resources

- **Documentation**: https://docs.contexis.dev
- **GitHub Issues**: https://github.com/contexis-cmp/contexis/issues
- **Discussions**: https://github.com/contexis-cmp/contexis/discussions

---

*Thank you for trying Contexis. We welcome feedback and contributions to shape the roadmap.*