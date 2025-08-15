# Contexis Overview

Contexis is a framework for building reproducible AI applications using the **Context-Memory-Prompt (CMP)** architecture. It provides a **local-first development experience** that mirrors production but uses smaller, local components.

## Core Architecture

### CMP Components

- **Context**: Declarative agent configuration in `.ctx` YAML (persona, tools, guardrails, memory, testing)
- **Memory**: Versioned knowledge bases (vector store) and episodic logs, tenant-aware
- **Prompt**: Pure templates rendered at runtime with data and context

### Local-First Development

Contexis provides out-of-the-box local models for development:

- **Phi-3.5-Mini** (~2GB) for text generation
- **Sentence Transformers** (~90MB) for embeddings
- **Chroma(SQLite)** for local vector database

**No external API keys required** - start developing immediately!

## Key Properties

- **Versioned and validated contexts** (`src/core/schema/context_schema.json`)
- **Local model integration** with automatic download and management
- **Memory providers**: `sqlite` vector store (file-backed JSONL) and `episodic` conversation logs
- **Prompt engine** with include functions and helper funcs (`src/runtime/prompt/engine.go`)
- **HTTP runtime server** for chat (`ctx serve`), with local model inference
- **Security (optional)**: API key auth, rate limiting, audit trail
- **Production migration**: Easy switch to OpenAI, Anthropic, or other providers

## Project Layout

```
project/
├── contexts/<Component>/*.ctx
├── memory/<Component>/
├── prompts/<Component>/*.md
├── tools/<Component>/*.py
├── tests/<Component>/
├── config/environments/
│   ├── development.yaml  # Local models by default
│   └── production.yaml   # External providers
└── context.lock.json
```

## Quick Start

```bash
# Install
make install

# Create project (local-first by default)
ctx init my-ai-app
cd my-ai-app

# Set up environment
cp .env.example .env
pip install -r requirements.txt

# Optional: Pre-download local models
ctx models warmup

# Generate RAG system
ctx generate rag CustomerDocs --db=sqlite --embeddings=sentence-transformers

# Add knowledge
echo 'Your company policies here...' > memory/CustomerDocs/documents/policies.txt
ctx memory ingest --provider=sqlite --component=CustomerDocs --input=memory/CustomerDocs/documents/policies.txt

# Test with local models
ctx run CustomerDocs "What are your policies?"
```

## Development Workflow

### 1. Local Development
- **No API keys needed** - everything runs locally
- **Fast iteration** - no external API calls
- **Cost effective** - no usage costs
- **Privacy** - all data stays local

### 2. Production Migration
```yaml
# config/environments/production.yaml
providers:
  openai:
    api_key: ${OPENAI_API_KEY}
    model: gpt-4o-mini
```

The same code works seamlessly across environments!

## Features

### Local Models
- **Automatic download** on first use
- **Model warmup** command for pre-downloading
- **Progress tracking** during downloads
- **Error handling** for network issues

### Enhanced Developer Experience
- **Rails-like logging** with colored output
- **Guided workflows** with step-by-step instructions
- **Comprehensive testing** with drift detection
- **Hot reload** for development server

### Enterprise Ready
- **Multi-tenancy** with tenant isolation
- **Security controls** (optional auth, RBAC, rate limiting)
- **Monitoring** with Prometheus metrics
- **Kubernetes** deployment with Helm charts

## Model Providers

### Local (Default)
- **Text Generation**: Phi-3.5-Mini via local Python provider
- **Embeddings**: Sentence Transformers (all-MiniLM-L6-v2)
- **Vector DB**: Chroma with SQLite backend

### External
- **OpenAI**: GPT-4, GPT-3.5, text-embedding-3-small
- **Anthropic**: Claude-3 models
- **Hugging Face**: Any model via Inference API

## Performance

| Environment | Startup Time | Response Time | Memory Usage | Cost |
|-------------|-------------|---------------|--------------|------|
| Local Development | ~30s (first run) | 2-5s | ~2GB | Free |
| Production | Instant | 1-3s | Minimal | Per token |

## Getting Started

Start with the [Getting Started Guide](docs/guides/getting-started.md) for a complete walkthrough of building your first AI application with local models.
