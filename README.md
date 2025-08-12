# Contexis: Context-Memory-Prompt (CMP) Framework

A Rails-inspired framework for building reproducible AI applications with architectural discipline.

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![Python Version](https://img.shields.io/badge/Python-3.8+-green.svg)](https://python.org)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Build Status](https://github.com/contexis-cmp/contexis/actions/workflows/build-and-helm.yml/badge.svg)](https://github.com/contexis-cmp/contexis/actions/workflows/build-and-helm.yml)
[![Vulnerability Scan](https://github.com/contexis-cmp/contexis/actions/workflows/security.yml/badge.svg)](https://github.com/contexis-cmp/contexis/actions/workflows/security.yml)

## Overview

Contexis introduces the **Context-Memory-Prompt (CMP)** architecture that treats AI components as version-controlled, first-class citizens. Like Rails brought MVC to web development, CMP brings architectural discipline to AI application engineering.

The framework provides a dual-language approach:
- **Go CLI** for fast, reliable command-line operations
- **Python Core** for AI/ML functionality and integrations

##  Quick Start

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

## ️ Core Architecture

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

##  Key Features

###  Generator-Driven Development
- Scaffold complete AI applications with best practices
- Pre-built templates for RAG, agents, and workflows
- Automatic test suite generation

###  Version Control & Reproducibility
- Lock all AI components for reproducible deployments
- Context versioning and drift detection
- Deterministic AI behavior across environments

###  Testing & Quality Assurance
- **Drift Detection:** Automated testing for AI behavior consistency
- **Correctness Tests:** Business logic validation
- **Performance Monitoring:** Response time and quality metrics

###  Enterprise Ready
- **Multi-Tenancy:** Built-in context isolation
- **Security & Compliance (opt-in):** API key auth, RBAC, rate limiting, audit logs, encryption-at-rest
  - Toggles:
    - `CMP_AUTH_ENABLED=true` – enable API key auth/RBAC/rate limiting
    - `CMP_PI_ENFORCEMENT=true` – enable prompt-injection classification and blocking
    - `CMP_REQUIRE_CITATION=true` – enforce citations on memory-backed responses
    - `CMP_PII_MODE=off|redact|block` – control PII handling in rendered outputs
- **Provider Agnostic:** Switch between AI models without code changes
- **Scalable Architecture:** From prototypes to production systems

###  Runtime Engine
- **Context Service:** Tenant-aware resolution, inheritance/merge, schema validation
- **Memory Service:** Pluggable providers (file-backed vector store, episodic logs)
- **Prompt Engine:** Template loading, includes, format validation, simple token trimming
- **Guardrails:** Capability validation, format/max_tokens enforcement
- **HTTP Server:** `/api/v1/chat` with health `/healthz`, readiness `/readyz`, version `/version`, metrics `/metrics` (supports optional auth/RBAC/rate limiting)
- **Worker:** `ctx worker` exposes `/healthz` and `/metrics` for background processing

## ️ Available Commands

```bash
# Project Management
ctx init <project-name>                   # Create new CMP project
ctx generate rag|agent|workflow|plugin <name>    # Generate components (incl. plugin scaffolding)

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

# Plugins
ctx plugin list                           # List installed plugins
ctx plugin info <name>                    # Show plugin details
ctx plugin install <path|zip_url|git_url> # Install plugin from local path, ZIP URL, or Git URL (supports #ref)
ctx plugin remove <name>                  # Uninstall plugin

## Security Controls

- Input PI Guard (optional): when `CMP_PI_ENFORCEMENT=true`, user inputs are classified for risky phrases (e.g., “ignore previous instructions”). High-risk requests return 403 and are audited.
- Source-constrained mode (optional): when `CMP_REQUIRE_CITATION=true` and a memory-backed component is used, the request requires retrieved results; if none, returns 424. Rendered outputs must contain citations (`Source:`) or return 422.
- PII handling (optional): `CMP_PII_MODE=off|redact|block`. When `block`, responses containing PII return 422; when `redact`, PII is masked before returning.
- OOB Action Gating (optional): declare sensitive actions via `CMP_OOB_REQUIRED_ACTIONS`, require `X-OOB-Confirmed: true` header (or `data.oob_confirmed=true`) for approval; otherwise 403.

### Response examples

```bash
# 403 Forbidden (OOB required)
curl -X POST http://localhost:8000/api/v1/chat \
  -H 'Content-Type: application/json' \
  -H 'X-OOB-Confirmed: false' \
  -d '{"context":"SupportBot","component":"SupportBot","query":"delete account","data":{"action":"account_action"}}'

# 422 Unprocessable Entity (missing citations or PII when blocked)
curl -X POST http://localhost:8000/api/v1/chat \
  -H 'Content-Type: application/json' \
  -d '{"context":"SupportBot","component":"SupportBot","query":"q"}'

# 424 Failed Dependency (no sources found)
curl -X POST http://localhost:8000/api/v1/chat \
  -H 'Content-Type: application/json' \
  -d '{"context":"SupportBot","component":"SupportBot","query":"nonexistent"}'
```
- Metrics: Prometheus exposes `cmp_security_prompt_injection_detections_total`, `cmp_security_policy_violations_total`, `cmp_security_blocked_responses_total`.
```

##  Examples

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

##  Configuration

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

##  Testing

### Run Go test suites
```bash
# All suites (unit, integration, e2e)
ctx test --all --coverage --junit --out tests/reports

# Specific suites
ctx test --unit --coverage
ctx test --integration
ctx test --e2e

# By category defined in tests/test_config.yaml
ctx test --category=agent_generator --coverage
```

### Drift detection
```bash
# Run drift tests for all components found under tests/**/rag_drift_test.yaml
ctx test --drift-detection --out tests/reports

# Limit to a component; write JUnit and update baselines
ctx test --drift-detection --component CustomerDocs --semantic --junit --update-baseline --out tests/reports
```

### Reports and artifacts
- Drift: `tests/reports/drift_<Component>.json`, `tests/reports/drift_index.json`, optional `junit-drift.xml`
- Go tests: `tests/reports/go_<suite>.txt`, `tests/reports/go_tests.json`, optional `junit-go.xml`
- Coverage: profiles under `tests/coverage/*.out` (when `--coverage` is used)

##  Deployment

### Docker Deployment
```bash
ctx build --environment=production
ctx build --image contexis-cmp/contexis --tag dev
ctx deploy --target=docker --image=contexis-cmp/contexis:dev
```

### Kubernetes Deployment (Helm)
```bash
helm upgrade --install contexis charts/contexis-app \
  --set image.repository=contexis-cmp/contexis \
  --set image.tag=dev \
  --set env.CMP_ENV=production
```

### Rollouts and External Secrets
- Enable Argo Rollouts: `--set rollouts.enabled=true` (requires CRDs)
- Enable ExternalSecret: `--set externalSecrets.enabled=true` and configure your SecretStore

### Observability
- Prometheus metrics at `/metrics` on app and worker
- Health endpoints: `/healthz`, `/readyz`

### Hugging Face Inference (optional)
- Set `HF_TOKEN` and `HF_MODEL_ID` in the environment (and optionally `HF_ENDPOINT`)
- When configured, the server will call the HF Inference API after rendering prompts

CLI quick test:
```bash
export HF_TOKEN=...  # your token
export HF_MODEL_ID=meta-llama/Meta-Llama-3.1-8B-Instruct
ctx hf test-model "Hello from Contexis!"
```

##  Contributing

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

##  Documentation

- [Technical RFC](docs/technical_rfc.md) - Detailed architecture specification
- [API Reference](docs/api/) - Complete API documentation
- [Examples](examples/) - Working examples and tutorials
- [Guides](docs/guides/) - Step-by-step guides
  - [Hugging Face Integration](docs/guides/hugging-face.md)

##  License

MIT License - see [LICENSE](LICENSE) for details.

##  Acknowledgments

- Inspired by Rails' architectural patterns
- Built on the shoulders of the AI/ML community
- Powered by open source technologies

---

**Contexis** - Bringing architectural discipline to AI applications 
