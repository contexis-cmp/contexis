# CLI Guide

Install and verify:
```bash
make install
ctx version
```

Generators:
```bash
# Agent
ctx generate agent SupportBot --tools web_search,database --memory episodic

# RAG system
ctx generate rag CustomerDocs --db sqlite --embeddings bge-small-en

# Workflow
ctx generate workflow ContentPipeline --steps research,write,review
```

Context operations:
```bash
# Validate a context
ctx context validate SupportBot

# Clear runtime context cache
ctx context reload
```

Prompt operations:
```bash
# Render a prompt template
ctx prompt render --component SupportBot --template agent_response.md --data '{"user":"Alice"}'
```

Memory operations:
```bash
# Ingest policies (one document per line) into sqlite vector store
ctx memory ingest --provider sqlite --component CustomerDocs --model bge-small-en --input policies.txt

# Search memory
ctx memory search --provider sqlite --component CustomerDocs --query "return policy" --top-k 5

# Optimize (optional)
ctx memory optimize --provider sqlite --component CustomerDocs --version <version-id>
```

Run HTTP server:
```bash
# Start server on :8000
ctx serve --addr :8000
```

Deploy:
```bash
# Build container
ctx build --image contexis-cmp/contexis --tag latest

# Deploy with Docker
ctx deploy --target docker --image contexis-cmp/contexis:latest --ports 8000:8000 --detach
```
