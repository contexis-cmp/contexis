# CLI Guide

Install and verify:
```bash
make install
ctx version
```

## Project Management

```bash
# Initialize a new project (local-first by default)
ctx init my-ai-app

# Generate components
ctx generate rag CustomerDocs --db=sqlite --embeddings=sentence-transformers
ctx generate agent SupportBot --tools web_search,database --memory episodic
ctx generate workflow ContentPipeline --steps research,write,review
```

## Development Commands

```bash
# Run queries directly (uses local models by default)
ctx run CustomerDocs "What is your return policy?"

# Start development server
ctx serve --addr :8000

# Pre-download local models (recommended for first run)
ctx models warmup
```

## Context Operations

```bash
# Validate a context
ctx context validate SupportBot

# Clear runtime context cache
ctx context reload
```

## Prompt Operations

```bash
# Render a prompt template
ctx prompt render --component SupportBot --template agent_response.md --data '{"user":"Alice"}'
```

## Memory Operations

```bash
# Ingest documents into local vector store
ctx memory ingest --provider sqlite --component CustomerDocs --input policies.txt

# Search memory
ctx memory search --provider sqlite --component CustomerDocs --query "return policy" --top-k 5

# Optimize (optional)
ctx memory optimize --provider sqlite --component CustomerDocs --version <version-id>
```

## Testing

```bash
# Run all tests
ctx test

# Run drift detection
ctx test --drift-detection --component CustomerDocs

# Run with coverage
ctx test --all --coverage
```

## Deployment

```bash
# Build container
ctx build --image contexis-cmp/contexis --tag latest

# Deploy with Docker
ctx deploy --target docker --image contexis-cmp/contexis:latest --ports 8000:8000 --detach
```

## Local Development Workflow

### 1. Initialize Project
```bash
ctx init my-support-bot
cd my-support-bot
cp .env.example .env
pip install -r requirements.txt
```

### 2. Set Up Local Models
```bash
# Pre-download models (optional but recommended)
ctx models warmup
```

### 3. Generate Components
```bash
# Create a RAG system
ctx generate rag MyFirstRAG --db=sqlite --embeddings=sentence-transformers

# Add knowledge
echo 'Your company policies here...' > memory/MyFirstRAG/documents/policies.txt
ctx memory ingest --provider=sqlite --component=MyFirstRAG --input=memory/MyFirstRAG/documents/policies.txt
```

### 4. Test and Run
```bash
# Test your system
ctx run MyFirstRAG "What are your policies?"

# Start server for continuous development
ctx serve --addr :8000
```

## Environment Configuration

### Local Development (Default)
```yaml
# config/environments/development.yaml
providers:
  local:
    model: microsoft/DialoGPT-medium
    temperature: 0.1
    max_tokens: 1000

embeddings:
  provider: sentence-transformers
  model: all-MiniLM-L6-v2

vector_db:
  provider: chroma
  path: ./data/embeddings
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

## Command Reference

### Global Flags
- `--debug`: Enable debug output
- `--log-level`: Set log level (debug, info, warn, error)
- `--config`: Specify config file path

### Run Command
```bash
ctx run <context> <query> [flags]
```
- `--addr`: Server address (default: :8000)
- `--component`: Component name (defaults to context name)
- `--data`: Additional JSON data
- `--timeout`: Request timeout in seconds (default: 30)
- `--top-k`: Number of memory results to retrieve (default: 5)

### Generate Command
```bash
ctx generate <type> <name> [flags]
```
Types: `rag`, `agent`, `workflow`

### Memory Commands
```bash
ctx memory ingest --provider <provider> --component <name> --input <file>
ctx memory search --provider <provider> --component <name> --query <query>
ctx memory optimize --provider <provider> --component <name>
```

### Test Command
```bash
ctx test [flags]
```
- `--all`: Run all test suites
- `--drift-detection`: Run drift detection tests
- `--coverage`: Generate coverage reports
- `--component`: Test specific component
- `--out`: Output directory for reports
