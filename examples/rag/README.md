# RAG Application Example

This example demonstrates how to build a Retrieval-Augmented Generation (RAG) system using the Contexis CMP framework.

## Quick Start

```bash
# Initialize the project
ctx init customer-support-rag
cd customer-support-rag

# Generate the RAG system
ctx generate rag CustomerDocs --db=sqlite --embeddings=openai

# Ingest your knowledge base (one document per line)
printf "Returns are accepted within 30 days.\nShipping takes 3-5 business days." > docs.txt
ctx memory ingest --provider=sqlite --component=CustomerDocs --input=docs.txt

# Test the system (optional)
ctx test

# Search the knowledge base
ctx memory search --provider=sqlite --component=CustomerDocs --query="What is your return policy?" --top-k=3

# Render a response with a prompt template (optional)
ctx prompt render --component=CustomerDocs --template=search_response.md --data='{"UserQuery":"What is your return policy?"}'

# Serve a simple API (optional)
ctx serve --addr :8000
# curl -X POST http://localhost:8000/api/v1/chat -H 'Content-Type: application/json' \
#   -d '{"context":"CustomerDocs","component":"CustomerDocs","query":"return policy","top_k":3,"data":{"user_query":"What is your return policy?"}}'
```

## Generated Structure

```
customer-support-rag/
├── contexts/
│   └── support_agent.ctx          # Agent role and behavior
├── memory/
│   ├── documents/                  # Knowledge base files
│   └── embeddings/                 # Vector embeddings
├── prompts/
│   └── support_response.md         # Response template
├── tools/
│   └── semantic_search.py         # Search implementation
├── tests/
│   ├── drift_detection.py         # Similarity monitoring
│   └── correctness.py             # Business logic validation
└── context.lock.json              # Version locks
```

## Key Features

- **Semantic Search**: Find relevant documents using embeddings
- **Context-Aware Responses**: Combine retrieved knowledge with LLM generation
- **Version Control**: Track all components for reproducibility
- **Drift Detection**: Monitor for unexpected behavior changes
- **Multi-Provider Support**: Switch between OpenAI, Anthropic, etc.

## Configuration

Edit `config/environments/development.yaml` to customize:

- AI provider settings
- Vector database configuration
- Testing thresholds
- Logging preferences

## Testing

The generated test suite includes:

- **Drift Detection**: Monitors response similarity over time
- **Correctness Tests**: Validates business logic compliance
- **Performance Tests**: Ensures response time requirements
- **Integration Tests**: End-to-end workflow validation

## Deployment

```bash
# Build for production
ctx build --environment=production

# Deploy to container
ctx deploy --target=docker

# Deploy to Kubernetes
ctx deploy --target=kubernetes
```

## Customization

### Adding New Knowledge

```bash
# Add documents to memory (append to your docs file and re-ingest)
printf "Another policy line" >> docs.txt
ctx memory ingest --provider=sqlite --component=CustomerDocs --input=docs.txt
```

### Modifying Context

Edit `contexts/support_agent.ctx` to change:

- Agent persona and capabilities
- Tool integrations
- Response guardrails
- Testing parameters

### Custom Prompts

Create new templates in `prompts/` and reference them in your context files.

## Troubleshooting

### Common Issues

1. **Low Search Relevance**: Adjust embedding model or chunk size
2. **Response Drift**: Check context boundaries and guardrails
3. **Performance Issues**: Optimize vector database or reduce context size
4. **Provider Errors**: Verify API keys and rate limits

### Debug Mode

```bash
# Enable debug logging
ctx run CustomerDocs "test question" --debug

# View detailed logs
ctx logs --level=debug
```

## Next Steps

- Add more sophisticated tools (database lookups, API calls)
- Implement conversation memory for multi-turn interactions
- Add authentication and multi-tenancy
- Create custom prompt templates for your domain 