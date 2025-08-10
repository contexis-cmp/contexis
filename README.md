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

```bash
# Clone the repository
git clone https://github.com/contexis-cmp/contexis.git
cd contexis

# Install the framework
make install

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

# Add your knowledge base
ctx memory add --file=./docs/company_policies.md

# Test the system
ctx test

# Run a query
ctx run query "What is your return policy?"
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

### 🔌 Rich Integrations
- **Vector Databases:** ChromaDB, Pinecone, Weaviate
- **AI Providers:** OpenAI, Anthropic, Cohere, Local models
- **Databases:** PostgreSQL, SQLite, Redis
- **Deployment:** Docker, Kubernetes, Cloud platforms

## 🛠️ Available Commands

```bash
# Project Management
ctx init <project-name>           # Create new CMP project
ctx generate <type> <name>        # Generate components
ctx build [--environment=prod]    # Build for deployment

# Memory Operations
ctx memory add --file=<path>      # Add documents to memory
ctx memory update --force         # Update embeddings
ctx memory search <query>         # Search knowledge base

# Testing & Validation
ctx test                          # Run all tests
ctx test --drift                 # Drift detection only
ctx test --correctness           # Business logic tests

# Deployment
ctx deploy --target=docker       # Deploy to Docker
ctx deploy --target=kubernetes   # Deploy to K8s
ctx logs --level=debug          # View application logs

# Development
ctx dev                          # Start development server
ctx validate                     # Validate framework
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

### Performance Testing
Validates response times and quality:

```bash
ctx test --performance --max-latency=2s
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
