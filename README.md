# Contexis CMP Framework

A framework for building reproducible AI applications using the **Context-Memory-Prompt (CMP)** architecture. Contexis treats AI components as version-controlled, first-class citizens, bringing architectural discipline to AI application engineering.

## 🚀 Quick Start

### 1. Install Contexis

```bash
# Clone and install
git clone https://github.com/contexis-cmp/contexis.git
cd contexis
make install

# Verify installation
ctx version
```

### 2. Create Your First AI Application

```bash
# Initialize a new project (local-first by default)
ctx init my-support-bot
cd my-support-bot

# Set up environment
cp .env.example .env
pip install -r requirements.txt

# Create your first RAG system
ctx generate rag MyFirstRAG --db=sqlite --embeddings=sentence-transformers

# Add knowledge
echo 'Our return policy allows returns within 30 days with receipt.' > memory/MyFirstRAG/documents/policies.txt
ctx memory ingest --provider=sqlite --component=MyFirstRAG --input=memory/MyFirstRAG/documents/policies.txt

# Test your system
ctx run MyFirstRAG "What is your return policy?"
```

**That's it!** Your AI application is now running with local models (Phi-3.5-Mini + Chroma + Sentence Transformers) - no external API calls needed.

## 🎯 Key Features

### **Local-First Development**
- **Out-of-the-box local models**: Phi-3.5-Mini for text generation, Sentence Transformers for embeddings
- **Local vector database**: Chroma with SQLite backend
- **No external dependencies**: Start developing immediately without API keys
- **Production-ready**: Easy migration to OpenAI, Anthropic, or other providers

### **CMP Architecture**
- **Context**: Declarative agent configuration (persona, tools, guardrails)
- **Memory**: Versioned knowledge bases with semantic search
- **Prompt**: Pure templates rendered at runtime

### **Developer Experience**
- **Enhanced Logging**: Rails-like colored output with detailed progress tracking
- **Guided Workflows**: Step-by-step instructions for each command
- **Comprehensive Testing**: Drift detection, correctness testing, and performance monitoring
- **Hot Reload**: Development server with automatic reloading

### **Enterprise Ready**
- **Multi-tenancy**: Tenant isolation and RBAC
- **Security Controls**: API key auth, rate limiting, audit logging
- **Monitoring**: Prometheus metrics and OpenTelemetry tracing
- **Kubernetes**: Helm charts and deployment manifests

## 📚 Documentation

- **[Getting Started Guide](docs/guides/getting-started.md)** - Complete walkthrough
- **[CLI Reference](docs/cli.md)** - All available commands
- **[API Reference](docs/api/README.md)** - HTTP API documentation
- **[Deployment Guide](docs/deployment.md)** - Production deployment
- **[Security Guide](docs/security.md)** - Security features and configuration

## 🔧 Available Commands

```bash
# Project Management
ctx init <project-name>          # Create new project
ctx generate <type> <name>       # Generate components (rag, agent, workflow)

# Development
ctx run <context> <query>        # Run queries with local models
ctx serve --addr :8000          # Start development server
ctx models warmup               # Pre-download local models

# Memory Operations
ctx memory ingest --provider=sqlite --component=<name> --input=<file>
ctx memory search --provider=sqlite --component=<name> --query="<query>"

# Testing
ctx test                        # Run all tests
ctx test --drift-detection      # Monitor AI behavior consistency

# Context Management
ctx context validate <name>     # Validate context files
ctx context reload              # Reload context cache
```

## 🏗️ Architecture

### Local Development Stack
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Phi-3.5-Mini  │    │  Sentence       │    │   Chroma        │
│   (Text Gen)    │    │  Transformers   │    │   (Vector DB)   │
│   ~2GB RAM      │    │  (Embeddings)   │    │   SQLite Backend│
│                 │    │  ~90MB RAM      │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Production Migration
```yaml
# config/environments/production.yaml
providers:
  openai:
    api_key: ${OPENAI_API_KEY}
    model: gpt-4o-mini
embeddings:
  provider: openai
  model: text-embedding-3-small
vector_db:
  provider: pinecone
  api_key: ${PINECONE_API_KEY}
```

## 🧪 Testing

```bash
# Run comprehensive tests
ctx test --all --coverage

# Monitor AI behavior drift
ctx test --drift-detection --component=MyFirstRAG

# Test specific scenarios
ctx test --correctness --rules=./tests/business_rules.yaml
```

## 🚀 Deployment

### Local Development
```bash
ctx serve --addr :8000
```

### Docker
```bash
ctx build --image contexis-cmp/contexis --tag latest
ctx deploy --target docker --image contexis-cmp/contexis:latest
```

### Kubernetes
```bash
# Deploy with Helm
helm install contexis ./charts/contexis-app \
  --set environment=production \
  --set replicas=3
```

## 🔒 Security

Enable security features:
```bash
export CMP_AUTH_ENABLED=true
export CMP_API_TOKENS=devtoken@tenantA:chat:execute|context:read
export CMP_PI_ENFORCEMENT=true
export CMP_REQUIRE_CITATION=true
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

## 🆘 Support

- **Documentation**: [docs.contexis.dev](https://docs.contexis.dev)
- **Issues**: [GitHub Issues](https://github.com/contexis-cmp/contexis/issues)
- **Discussions**: [GitHub Discussions](https://github.com/contexis-cmp/contexis/discussions)
- **Community**: [Discord](https://discord.gg/contexis)

---

**Contexis** - Bringing architectural discipline to AI applications 
