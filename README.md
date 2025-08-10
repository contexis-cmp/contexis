# Contexis: Context-Memory-Prompt (CMP) Framework

A Rails-inspired framework for building reproducible AI applications with architectural discipline.

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![Python Version](https://img.shields.io/badge/Python-3.8+-green.svg)](https://python.org)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)](https://github.com/contexis-cmp/contexis/actions)

## Overview

Contexis introduces the **Context-Memory-Prompt (CMP)** architecture that treats AI components as version-controlled, first-class citizens. Like Rails brought MVC to web development, CMP brings architectural discipline to AI application engineering.

The framework provides a dual-language approach:
- **Go CLI** for fast, reliable command-line operations
- **Python Core** for AI/ML functionality and integrations

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Python 3.8+
- Git

### Installation

#### Option 1: Automatic Installation (Recommended)
```bash
# Clone the repository
git clone https://github.com/contexis-cmp/contexis.git
cd contexis

# Run the installation script (automatically handles PATH setup)
./scripts/install.sh

# Restart your terminal or source your shell config
source ~/.bashrc  # or ~/.zshrc

# Verify installation
ctx version
```

#### Option 2: Manual Local Installation
```bash
# Clone the repository
git clone https://github.com/contexis-cmp/contexis.git
cd contexis

# Install to user's local directory (no sudo required)
make install-local

# Add to your PATH (choose your shell)
# For bash/zsh, add to ~/.bashrc or ~/.zshrc:
export PATH="$HOME/.local/bin:$PATH"

# For fish, add to ~/.config/fish/config.fish:
set -gx PATH $HOME/.local/bin $PATH

# Reload your shell or source the profile
source ~/.bashrc  # or ~/.zshrc

# Verify installation
ctx version
```

#### Option 2: System Installation (requires sudo)
```bash
# Clone the repository
git clone https://github.com/contexis-cmp/contexis.git
cd contexis

# Install system-wide (requires sudo)
sudo make install

# Verify installation
ctx version
```

### Your First AI Application

```bash
# Initialize a new CMP project
ctx init my-ai-app
cd my-ai-app

# Generate a RAG system
ctx generate rag CustomerDocs --db=sqlite --embeddings=openai

# Ingest documents into memory (one per line)
printf "Returns are accepted within 30 days.\nShipping takes 3-5 business days." > docs.txt
ctx memory ingest --provider=sqlite --component=CustomerDocs --input=docs.txt

# Search the knowledge base
ctx memory search --provider=sqlite --component=CustomerDocs --query="return policy" --top-k=3

# Render a prompt template
ctx prompt render --component=CustomerDocs --template=search_response.md --data='{"UserQuery":"return policy"}'

# Serve a simple API (optional)
ctx serve --addr :8000
# curl -X POST http://localhost:8000/api/v1/chat -H 'Content-Type: application/json' \
#   -d '{"context":"SupportBot","component":"SupportBot","query":"return policy","top_k":3,"data":{"user_input":"Hi"}}'
```

## 🏗️ Core Architecture

### CMP Components

- **Context:** Declarative instructions, agent roles, tool definitions
- **Memory:** Versioned knowledge stores, vector databases, logs  
- **Prompt:** Pure templates hydrated at runtime

### Project Structure

```
my-ai-app/
├── contexts/              # Agent roles and behaviors
│   └── support_agent.ctx
├── memory/               # Knowledge base and embeddings
│   ├── documents/
│   └── embeddings/
├── prompts/              # Response templates
│   └── support_response.md
├── tools/               # Custom integrations
│   └── semantic_search.py
├── tests/               # Drift detection and validation
│   ├── drift_detection.py
│   └── correctness.py
├── config/              # Environment configuration
│   └── environments/
└── context.lock.json    # Version locks
```

## ✨ Key Features

### 🎯 Generator-Driven Development
- Scaffold complete AI applications with best practices
- Pre-built templates for RAG, agents, and workflows
- Automatic test suite generation

### 🔒 Version Control & Reproducibility
- Lock all AI components for reproducible deployments
- Context versioning and drift detection
- Deterministic AI behavior across environments

### 🧪 Testing & Quality Assurance
- **Drift Detection:** Automated testing for AI behavior consistency
- **Correctness Tests:** Business logic validation
- **Performance Monitoring:** Response time and quality metrics

### 🏢 Enterprise Ready
- **Multi-Tenancy:** Built-in context isolation
- **Provider Agnostic:** Switch between AI models without code changes
- **Scalable Architecture:** From prototypes to production systems

### 🔌 Runtime Engine (Phase 2)
- **Context Service:** Tenant-aware resolution, inheritance/merge, schema validation
- **Memory Service:** Pluggable providers (file-backed vector store, episodic logs)
- **Prompt Engine:** Template loading, includes, format validation, simple token trimming
- **Guardrails:** Capability validation, format/max_tokens enforcement
- **HTTP Server:** Minimal `/api/v1/chat` endpoint for local experimentation

## 🛠️ Available Commands

```bash
# Project Management
ctx init <project-name>                   # Create new CMP project
ctx generate rag|agent|workflow <name>    # Generate components

# Context Operations
ctx context validate <name> [--tenant=<id>]   # Validate a .ctx file
ctx context reload                            # Clear runtime context cache

# Memory Operations
ctx memory ingest --provider=sqlite|episodic --component=<Comp> [--tenant=<id>] --input=<file>
ctx memory search --provider=sqlite|episodic --component=<Comp> [--tenant=<id>] --query=<q> [--top-k=5]
ctx memory optimize --provider=sqlite|episodic --component=<Comp> [--tenant=<id>]

# Prompt Operations
ctx prompt render --component=<Comp> --template=<path> --data='{"k":"v"}'
ctx prompt validate --format=json|markdown|text --input=<file>
ctx prompt-lint --component=<Comp>

# Reproducibility
ctx lock generate                         # Write context.lock.json

# Server
ctx serve --addr :8000                    # Start HTTP server
```

## 📚 Examples

### RAG System
```bash
ctx generate rag CustomerDocs --db=sqlite --embeddings=openai
```

### Conversational Agent
```bash
ctx generate agent SupportBot --memory=conversation --tools=api
```

### Workflow Pipeline
```bash
ctx generate workflow DataProcessor --steps=extract,transform,load
```

## 🔧 Configuration

### Environment Setup

Edit `config/environments/development.yaml`:

```yaml
ai:
  provider: openai
  model: gpt-4
  api_key: ${OPENAI_API_KEY}

memory:
  vector_db: chromadb
  embeddings: openai
  chunk_size: 1000

testing:
  drift_threshold: 0.85
  correctness_rules: ./tests/rules.yaml
```

### Custom Contexts

Define agent behaviors in `contexts/`:

```yaml
# contexts/support_agent.ctx
role: Customer Support Agent
capabilities:
  - answer_product_questions
  - process_returns
  - escalate_issues
tools:
  - semantic_search
  - knowledge_base_lookup
constraints:
  - always_be_polite
  - never_share_internal_info
```

## 🧪 Testing

### Drift Detection
Monitors AI behavior consistency over time:

```bash
ctx test --drift --threshold=0.85
```

### Correctness Validation
Ensures business logic compliance:

```bash
ctx test --correctness --rules=./tests/business_rules.yaml
```

### Runtime Package Tests
Run targeted tests for runtime packages:

```bash
go test ./src/runtime/context ./src/runtime/memory ./src/runtime/prompt ./src/runtime/guardrails -v
```

## 🚀 Deployment

### Docker Deployment
```bash
ctx build --environment=production
ctx deploy --target=docker --image=my-ai-app:latest
```

### Kubernetes Deployment
```bash
ctx deploy --target=kubernetes --namespace=ai-apps
```

### Cloud Platforms
```bash
ctx deploy --target=aws --region=us-west-2
ctx deploy --target=gcp --project=my-project
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Setup

```bash
# Clone and setup
git clone https://github.com/contexis-cmp/contexis.git
cd contexis
make setup

# Run tests
make test

# Format code
make format

# Build
make build
```

## 📖 Documentation

- [Technical RFC](docs/technical_rfc.md) - Detailed architecture specification
- [API Reference](docs/api/) - Complete API documentation
- [Examples](examples/) - Working examples and tutorials
- [Guides](docs/guides/) - Step-by-step guides

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

- Inspired by Rails' architectural patterns
- Built on the shoulders of the AI/ML community
- Powered by open source technologies

---

**Contexis** - Bringing architectural discipline to AI applications 🚀
