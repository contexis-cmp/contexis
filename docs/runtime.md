# Runtime and API

Start the server or run queries:
```bash
# Run a query directly (uses local models by default)
ctx run CustomerDocs "What is your return policy?"

# Start server manually for continuous use
ctx serve --addr :8000
```

## Local Development

Contexis provides a **local-first development experience** with out-of-the-box local models:

- **Phi-3.5-Mini** (~2GB) for text generation
- **Sentence Transformers** (~90MB) for embeddings  
- **Chroma(SQLite)** for local vector database

No external API keys required - everything runs locally!

## Health Endpoints

- GET `/healthz` → 200 ok
- GET `/readyz` → 200 ready
- GET `/version` → framework version
- GET `/metrics` → Prometheus metrics

## Chat API

- POST `/api/v1/chat`
- Request body:
```json
{
  "tenant_id": "",
  "context": "CustomerDocs",
  "component": "CustomerDocs",
  "query": "What is your return policy?",
  "top_k": 5,
  "data": {},
  "prompt_file": "search_response.md"
}
```
- Response:
```json
{ "rendered": "...model output or rendered prompt..." }
```

## Model Providers

### Local Models (Default)
When no external providers are configured, the server uses local models:
- **Text Generation**: Phi-3.5-Mini via local Python provider
- **Embeddings**: Sentence Transformers (all-MiniLM-L6-v2)
- **Vector Database**: Chroma with SQLite backend

### External Providers
Configure external providers in `config/environments/`:

```yaml
# For OpenAI
providers:
  openai:
    api_key: ${OPENAI_API_KEY}
    model: gpt-4o-mini

# For Hugging Face
providers:
  huggingface:
    token: ${HF_TOKEN}
    model_id: ${HF_MODEL_ID}
```

## Environment Variables

### Local Development
```bash
# Enable local models (default)
CMP_LOCAL_MODELS=true

# Python binary path (optional)
CMP_PYTHON_BIN=python3

# Model cache directory (optional)
CMP_MODEL_CACHE_DIR=./data/models

# Local timeout (optional, default: 600s)
CMP_LOCAL_TIMEOUT_SECONDS=300
```

### Production
```bash
# OpenAI
OPENAI_API_KEY=your_openai_api_key

# Hugging Face
HF_TOKEN=your_hf_token
HF_MODEL_ID=meta-llama/Meta-Llama-3.1-8B-Instruct

# Vector Database
PINECONE_API_KEY=your_pinecone_key
```

## Model Warmup

Pre-download local models for faster startup:
```bash
ctx models warmup
```

This downloads and initializes:
- Phi-3.5-Mini text generation model
- Sentence Transformers embedding model

## Security (Optional)

Enable security features:
```bash
export CMP_AUTH_ENABLED=true
export CMP_API_TOKENS=devtoken@tenantA:chat:execute|context:read
export CMP_PI_ENFORCEMENT=true
export CMP_REQUIRE_CITATION=true
```

## Performance

### Local Models
- **First Run**: Models download automatically (~3GB total)
- **Subsequent Runs**: Models load from cache (~30s startup)
- **Memory Usage**: ~2GB RAM for text generation, ~90MB for embeddings
- **Response Time**: 2-5 seconds per query (CPU-based)

### Production Models
- **Startup**: Instant (no model loading)
- **Response Time**: 1-3 seconds per query (API-based)
- **Cost**: Per-token pricing (varies by provider)

## Troubleshooting

### Local Model Issues
```bash
# Check model status
ctx models warmup

# Verify Python environment
python3 -c "import transformers; print('OK')"

# Check disk space
df -h ./data/models
```

### Server Issues
```bash
# Check server health
curl http://localhost:8000/healthz

# View server logs
ctx serve --addr :8000 --debug
```

### Performance Issues
```bash
# Monitor resource usage
htop

# Check model loading
tail -f logs/contexis.log
```
