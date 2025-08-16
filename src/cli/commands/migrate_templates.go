package commands

import (
	"fmt"
	"strings"
)

// Configuration templates
func getProductionConfig(provider string) string {
	switch provider {
	case "openai":
		return `environment: production

# OpenAI Provider Configuration
providers:
  openai:
    api_key: ${OPENAI_API_KEY}
    model: gpt-4o-mini
    temperature: 0.1
    max_tokens: 1000

# Embeddings Configuration
embeddings:
  provider: openai
  model: text-embedding-3-small
  dimensions: 1536

# Vector Database Configuration
vector_db:
  provider: pinecone
  api_key: ${PINECONE_API_KEY}
  environment: ${PINECONE_ENVIRONMENT}
  index_name: ${PINECONE_INDEX_NAME}

# Testing Configuration
testing:
  drift_threshold: 0.85
  similarity_threshold: 0.8
  max_test_duration: 300s

# Security Configuration
security:
  auth_enabled: true
  rate_limiting: true
  audit_logging: true
`
	case "anthropic":
		return `environment: production

# Anthropic Provider Configuration
providers:
  anthropic:
    api_key: ${ANTHROPIC_API_KEY}
    model: claude-3-sonnet-20240229
    temperature: 0.1
    max_tokens: 1000

# Embeddings Configuration
embeddings:
  provider: openai
  model: text-embedding-3-small
  dimensions: 1536

# Vector Database Configuration
vector_db:
  provider: pinecone
  api_key: ${PINECONE_API_KEY}
  environment: ${PINECONE_ENVIRONMENT}
  index_name: ${PINECONE_INDEX_NAME}

# Testing Configuration
testing:
  drift_threshold: 0.85
  similarity_threshold: 0.8
  max_test_duration: 300s

# Security Configuration
security:
  auth_enabled: true
  rate_limiting: true
  audit_logging: true
`
	case "huggingface":
		return `environment: production

# Hugging Face Provider Configuration
providers:
  huggingface:
    token: ${HF_TOKEN}
    model_id: ${HF_MODEL_ID}
    temperature: 0.1
    max_tokens: 1000

# Embeddings Configuration
embeddings:
  provider: sentence-transformers
  model: all-MiniLM-L6-v2
  dimensions: 384

# Vector Database Configuration
vector_db:
  provider: chroma
  path: ./data/embeddings
  collection_name: production_knowledge

# Testing Configuration
testing:
  drift_threshold: 0.85
  similarity_threshold: 0.8
  max_test_duration: 300s

# Security Configuration
security:
  auth_enabled: true
  rate_limiting: true
  audit_logging: true
`
	default:
		return ""
	}
}

func getProductionEnvTemplate(provider string) string {
	base := `# Production Environment Variables
# Copy this file to .env and fill in your actual values

# Application Settings
ENVIRONMENT=production
LOG_LEVEL=info

# Security Settings
CMP_AUTH_ENABLED=true
CMP_API_TOKENS=your_production_token@tenant:chat:execute|context:read
CMP_PI_ENFORCEMENT=true
CMP_REQUIRE_CITATION=true

`

	switch provider {
	case "openai":
		return base + `# OpenAI Configuration
OPENAI_API_KEY=your_openai_api_key_here

# Vector Database (Pinecone)
PINECONE_API_KEY=your_pinecone_api_key_here
PINECONE_ENVIRONMENT=your_pinecone_environment
PINECONE_INDEX_NAME=your_pinecone_index_name
`
	case "anthropic":
		return base + `# Anthropic Configuration
ANTHROPIC_API_KEY=your_anthropic_api_key_here

# Vector Database (Pinecone)
PINECONE_API_KEY=your_pinecone_api_key_here
PINECONE_ENVIRONMENT=your_pinecone_environment
PINECONE_INDEX_NAME=your_pinecone_index_name
`
	case "huggingface":
		return base + `# Hugging Face Configuration
HF_TOKEN=your_huggingface_token_here
HF_MODEL_ID=meta-llama/Meta-Llama-3.1-8B-Instruct

# Optional: Custom HF endpoint
# HF_ENDPOINT=https://api-inference.huggingface.co/models
`
	default:
		return base
	}
}

func getLocalConfig() string {
	return `environment: development

# Local Provider Configuration (Default)
providers:
  local:
    model: microsoft/DialoGPT-medium
    temperature: 0.1
    max_tokens: 1000

# Local Embeddings Configuration
embeddings:
  provider: sentence-transformers
  model: all-MiniLM-L6-v2
  dimensions: 384

# Local Vector Database Configuration
vector_db:
  provider: chroma
  path: ./data/embeddings
  collection_name: development_knowledge

# Testing Configuration
testing:
  drift_threshold: 0.85
  similarity_threshold: 0.8
  max_test_duration: 300s

# Development Settings
development:
  hot_reload: true
  debug_mode: true
`
}

func getLocalEnvTemplate() string {
	return `# Contexis Development Environment
# Local-first defaults (no external API keys needed)

# Enable local models (default)
CMP_LOCAL_MODELS=true

# Python binary path (optional)
CMP_PYTHON_BIN=python3

# Model cache directory (optional)
CMP_MODEL_CACHE_DIR=./data/models

# Local timeout (optional, default: 600s)
CMP_LOCAL_TIMEOUT_SECONDS=300

# Application Settings
LOG_LEVEL=debug
ENVIRONMENT=development

# ---
# Production Hints (uncomment and configure to switch providers)
# OPENAI_API_KEY=your_openai_api_key_here
# ANTHROPIC_API_KEY=your_anthropic_api_key_here
# HF_TOKEN=your_huggingface_token_here
# HF_MODEL_ID=meta-llama/Meta-Llama-3.1-8B-Instruct
`
}

func getDockerfile(provider string) string {
	return `# Contexis Production Dockerfile
FROM python:3.10-slim

# Set working directory
WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y \
    git \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements and install Python dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY . .

# Create data directory
RUN mkdir -p /app/data

# Expose port
EXPOSE 8000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8000/healthz || exit 1

# Run the application
CMD ["ctx", "serve", "--addr", ":8000"]
`
}

func getDockerCompose(provider string) string {
	return `# Contexis Production Docker Compose
version: '3.8'

services:
  contexis:
    build: .
    ports:
      - "8000:8000"
    environment:
      - ENVIRONMENT=production
      - LOG_LEVEL=info
    env_file:
      - .env.production
    volumes:
      - ./data:/app/data
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8000/healthz"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  # Optional: Add Redis for caching
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    restart: unless-stopped

volumes:
  redis_data:
`
}

func getMigrationGuide(provider string) string {
	return fmt.Sprintf(`# Migration Guide: Local to %s Production

This guide helps you migrate your Contexis project from local development to %s production.

## Prerequisites

1. **%s Account**: Ensure you have an active %s account
2. **API Keys**: Obtain your %s API keys
3. **Vector Database**: Set up a production vector database (Pinecone recommended)

## Migration Steps

### 1. Update Environment Variables

Copy the generated .env.production file and fill in your actual values:

`+"```"+`bash
cp .env.production .env
# Edit .env with your actual API keys
`+"```"+`

### 2. Test Configuration

Validate your configuration:

`+"```"+`bash
ctx migrate validate
`+"```"+`

### 3. Test Production Setup

Test your production configuration:

`+"```"+`bash
# Test with production provider
ctx run YourComponent "Test query"
`+"```"+`

### 4. Deploy to Production

#### Option A: Docker Deployment

`+"```"+`bash
# Build and run with Docker
docker-compose up -d
`+"```"+`

#### Option B: Kubernetes Deployment

`+"```"+`bash
# Deploy to Kubernetes
helm install contexis ./charts/contexis-app --set environment=production --set replicas=3
`+"```"+`

## Configuration Details

### %s Configuration

Your production configuration uses %s for:
- **Text Generation**: %s models
- **Embeddings**: %s embeddings
- **Vector Database**: Pinecone for scalable storage

### Performance Considerations

- **Response Time**: 1-3 seconds per query
- **Cost**: Per-token pricing (varies by model)
- **Scalability**: Horizontal scaling supported
- **Monitoring**: Prometheus metrics available

## Monitoring and Maintenance

### Health Checks

`+"```"+`bash
# Check application health
curl http://localhost:8000/healthz

# View metrics
curl http://localhost:8000/metrics
`+"```"+`

### Logs

`+"```"+`bash
# View application logs
docker-compose logs -f contexis

# View error logs
docker-compose logs -f contexis | grep ERROR
`+"```"+`

## Troubleshooting

### Common Issues

1. **API Key Errors**
   - Verify your %s API key is correct
   - Check API key permissions and quotas

2. **Connection Issues**
   - Verify network connectivity
   - Check firewall settings

3. **Performance Issues**
   - Monitor resource usage
   - Consider model optimization

### Support

- **Documentation**: [docs.contexis.dev](https://docs.contexis.dev)
- **Issues**: [GitHub Issues](https://github.com/contexis-cmp/contexis/issues)
- **Community**: [Discord](https://discord.gg/contexis)

## Rollback

If you need to rollback to local development:

`+"```"+`bash
ctx migrate production-to-local
`+"```"+`

This will restore your local development configuration.
`,
		strings.Title(provider), strings.Title(provider),
		strings.Title(provider), strings.Title(provider), strings.Title(provider),
		strings.Title(provider), strings.Title(provider),
		getProviderModel(provider), getProviderEmbeddings(provider),
		strings.Title(provider))
}

func getLocalDevelopmentGuide() string {
	return `# Local Development Guide

This guide helps you set up and use Contexis for local development.

## Quick Start

### 1. Initialize Project

` + "```" + `bash
ctx init my-project
cd my-project
` + "```" + `

### 2. Set Up Environment

` + "```" + `bash
# Copy environment template
cp .env.example .env

# Install dependencies
pip install -r requirements.txt
` + "```" + `

### 3. Download Local Models

` + "```" + `bash
# Pre-download models (recommended)
ctx models warmup
` + "```" + `

### 4. Create Your First Component

` + "```" + `bash
# Generate a RAG system
ctx generate rag MyRAG --db=sqlite --embeddings=sentence-transformers

# Add knowledge
echo 'Your company policies here...' > memory/MyRAG/documents/policies.txt
ctx memory ingest --provider=sqlite --component=MyRAG --input=memory/MyRAG/documents/policies.txt
` + "```" + `

### 5. Test Your System

` + "```" + `bash
# Test with local models
ctx run MyRAG "What are your policies?"
` + "```" + `

## Local Development Features

### Local Models

- **Phi-3.5-Mini**: ~2GB RAM for text generation
- **Sentence Transformers**: ~90MB RAM for embeddings
- **Chroma(SQLite)**: Local vector database

### Benefits

- **No API Keys**: Start developing immediately
- **Offline Capable**: Works without internet connection
- **Cost Effective**: No usage costs
- **Fast Iteration**: No external API calls

### Performance

- **First Run**: Models download automatically (~3GB total)
- **Subsequent Runs**: Models load from cache (~30s startup)
- **Response Time**: 2-5 seconds per query (CPU-based)

## Development Workflow

### 1. Component Development

` + "```" + `bash
# Generate components
ctx generate rag CustomerDocs --db=sqlite --embeddings=sentence-transformers
ctx generate agent SupportBot --tools web_search,database --memory episodic
ctx generate workflow ContentPipeline --steps research,write,review
` + "```" + `

### 2. Testing

` + "```" + `bash
# Run all tests
ctx test --all --coverage

# Test drift detection
ctx test --drift-detection --component=MyRAG

# Test specific scenarios
ctx test --correctness --rules=./tests/business_rules.yaml
` + "```" + `

### 3. Development Server

` + "```" + `bash
# Start development server
ctx serve --addr :8000

# Test via HTTP API
curl -X POST http://localhost:8000/api/v1/chat -H "Content-Type: application/json" -d '{"context": "MyRAG", "component": "MyRAG", "query": "What are your policies?", "top_k": 5}'
` + "```" + `

## Configuration

### Local Configuration

Your local configuration uses:
- **Text Generation**: Phi-3.5-Mini via local Python provider
- **Embeddings**: Sentence Transformers (all-MiniLM-L6-v2)
- **Vector Database**: Chroma with SQLite backend

### Environment Variables

` + "```" + `bash
# Enable local models (default)
CMP_LOCAL_MODELS=true

# Python binary path (optional)
CMP_PYTHON_BIN=python3

# Model cache directory (optional)
CMP_MODEL_CACHE_DIR=./data/models

# Local timeout (optional, default: 600s)
CMP_LOCAL_TIMEOUT_SECONDS=300
` + "```" + `

## Troubleshooting

### Common Issues

1. **Model Download Issues**
   ` + "```" + `bash
   # Pre-download models
   ctx models warmup
   
   # Check disk space
   df -h ./data/models
   ` + "```" + `

2. **Memory Issues**
   ` + "```" + `bash
   # Verify embeddings generation
   ctx memory search --provider=sqlite --component=MyRAG --query="test" --top-k=1
   ` + "```" + `

3. **Performance Issues**
   ` + "```" + `bash
   # Monitor resource usage
   htop
   
   # Check model loading
   tail -f logs/contexis.log
   ` + "```" + `

### Performance Optimization

1. **Model Warmup**
   ` + "```" + `bash
   # Pre-download models for faster startup
   ctx models warmup
   ` + "```" + `

2. **Vector Store Optimization**
   ` + "```" + `bash
   # Optimize vector store
   ctx memory optimize --provider=sqlite --component=MyRAG
   ` + "```" + `

## Next Steps

- **Add more documents** to expand knowledge base
- **Customize prompts** for better responses
- **Add business rules** for compliance
- **Set up monitoring** for development
- **Migrate to production** when ready

## Migration to Production

When you're ready to deploy to production:

` + "```" + `bash
# Migrate to production provider
ctx migrate local-to-production --provider=openai
` + "```" + `

This will:
1. Update your configuration for production
2. Generate environment templates
3. Create deployment files
4. Update documentation

## Resources

- [Getting Started Guide](docs/guides/getting-started.md)
- [CLI Reference](docs/cli.md)
- [Model Providers](docs/model_providers.md)
- [Memory Management](docs/memory.md)
`
}

func getProviderModel(provider string) string {
	switch provider {
	case "openai":
		return "GPT-4o-mini"
	case "anthropic":
		return "Claude-3-Sonnet"
	case "huggingface":
		return "Meta-Llama-3.1-8B-Instruct"
	default:
		return "Unknown"
	}
}

func getProviderEmbeddings(provider string) string {
	switch provider {
	case "openai":
		return "text-embedding-3-small"
	case "anthropic":
		return "text-embedding-3-small"
	case "huggingface":
		return "all-MiniLM-L6-v2"
	default:
		return "Unknown"
	}
}
