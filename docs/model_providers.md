# Model Providers

Contexis supports multiple model providers for text generation and embeddings, with **local models as the default** for development.

## Local Models (Default)

Contexis provides out-of-the-box local models for development:

### Text Generation
- **Model**: Phi-3.5-Mini (Microsoft)
- **Size**: ~2GB RAM
- **Performance**: CPU-based, 2-5 seconds per query
- **Quality**: Good for development and testing

### Embeddings
- **Model**: Sentence Transformers (all-MiniLM-L6-v2)
- **Size**: ~90MB RAM
- **Dimensions**: 384
- **Performance**: Fast local embeddings

### Vector Database
- **Provider**: Chroma with SQLite backend
- **Storage**: Local file-based
- **Performance**: Suitable for development and small production workloads

## Configuration

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

### Environment Variables
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

## External Providers

### OpenAI
```bash
export OPENAI_API_KEY=your_openai_api_key
```

```yaml
# config/environments/production.yaml
providers:
  openai:
    api_key: ${OPENAI_API_KEY}
    model: gpt-4o-mini
    temperature: 0.1
    max_tokens: 1000

embeddings:
  provider: openai
  model: text-embedding-3-small
```

### Hugging Face Inference API
```bash
export HF_TOKEN=your_hf_token
export HF_MODEL_ID=meta-llama/Meta-Llama-3.1-8B-Instruct
```

```yaml
# config/environments/production.yaml
providers:
  huggingface:
    token: ${HF_TOKEN}
    model_id: ${HF_MODEL_ID}
    temperature: 0.1
    max_tokens: 1000
```

### Anthropic
```bash
export ANTHROPIC_API_KEY=your_anthropic_api_key
```

```yaml
# config/environments/production.yaml
providers:
  anthropic:
    api_key: ${ANTHROPIC_API_KEY}
    model: claude-3-sonnet-20240229
    temperature: 0.1
    max_tokens: 1000
```

## Model Warmup

Pre-download local models for faster startup:
```bash
ctx models warmup
```

This command:
- Downloads Phi-3.5-Mini text generation model
- Downloads Sentence Transformers embedding model
- Initializes models for first use
- Shows progress and status

## Performance Comparison

| Provider | Startup Time | Response Time | Memory Usage | Cost |
|----------|-------------|---------------|--------------|------|
| Local | ~30s (first run) | 2-5s | ~2GB | Free |
| OpenAI | Instant | 1-3s | Minimal | Per token |
| Hugging Face | Instant | 2-4s | Minimal | Per token |
| Anthropic | Instant | 1-3s | Minimal | Per token |

## Migration Guide

### From Local to Production

1. **Update Configuration**
```yaml
# config/environments/production.yaml
providers:
  openai:
    api_key: ${OPENAI_API_KEY}
    model: gpt-4o-mini
```

2. **Set Environment Variables**
```bash
export OPENAI_API_KEY=your_openai_api_key
```

3. **Deploy**
```bash
ctx serve --addr :8000
```

The same code works seamlessly across environments!

### From Production to Local

1. **Update Configuration**
```yaml
# config/environments/development.yaml
providers:
  local:
    model: microsoft/DialoGPT-medium
```

2. **Enable Local Models**
```bash
export CMP_LOCAL_MODELS=true
```

3. **Warm Up Models**
```bash
ctx models warmup
```

## Troubleshooting

### Local Model Issues
```bash
# Check model status
ctx models warmup

# Verify Python environment
python3 -c "import transformers; print('OK')"

# Check disk space
df -h ./data/models

# Check memory usage
htop
```

### External Provider Issues
```bash
# Test OpenAI
curl -H "Authorization: Bearer $OPENAI_API_KEY" \
  https://api.openai.com/v1/models

# Test Hugging Face
curl -H "Authorization: Bearer $HF_TOKEN" \
  https://huggingface.co/api/models
```

### Performance Issues
```bash
# Monitor resource usage
htop

# Check model loading
tail -f logs/contexis.log

# Profile response times
time ctx run CustomerDocs "test query"
```
