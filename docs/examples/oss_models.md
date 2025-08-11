# Using GPT-OSS with Contexis

This guide shows how to use an open-source model (GPT-OSS placeholder) with the Contexis runtime.

## Model Configuration

Create `.contextrc` in your project root:

```yaml
models:
  primary: "gpt-oss-7b-instruct"
  embeddings: "bge-small-en"

runtime:
  temperature: 0.1
  max_tokens: 800
```

> Note: The runtime here focuses on context/prompt/memory orchestration. Model invocation is provider-agnostic and can be plugged in using your preferred OSS serving stack (e.g., llama.cpp, vLLM, TGI).

## Prompt Rendering

Render a prompt using the runtime prompt engine (model-agnostic):

```bash
ctx prompt render --component=CustomerDocs --template=search_response.md --data='{"UserQuery":"returns"}'
```

## Memory Search

Use the file-backed vector store for local development:

```bash
printf "Returns are accepted within 30 days." > docs.txt
ctx memory ingest --provider=sqlite --component=CustomerDocs --input=docs.txt
ctx memory search --provider=sqlite --component=CustomerDocs --query="return policy" --top-k=3
```

## Simple HTTP Chat

Start the server and POST to `/api/v1/chat`:

```bash
ctx serve --addr :8000
curl -X POST http://localhost:8000/api/v1/chat -H 'Content-Type: application/json' \
  -d '{
    "tenant_id":"",
    "context":"SupportBot",
    "component":"CustomerDocs",
    "query":"return policy",
    "top_k":3,
    "data": {"user_input": "What is the return policy?"}
  }'
```

## Notes

- Replace `gpt-oss-7b-instruct` with your local model identifier.
- Use your OSS runtime’s REST/gRPC client in your application layer to complete the generation call after rendering prompts and gathering memory.
- Contexis ensures the context/memory/prompt orchestration is consistent and reproducible across model providers.


