# Getting Started with Contexis

Welcome to Contexis! This guide will help you get up and running with your first AI application using the CMP (Context-Memory-Prompt) framework.

## Prerequisites

Before you begin, make sure you have the following installed:

- **Go 1.21+**: [Download from golang.org](https://golang.org/dl/)
- **Python 3.10+**: [Download from python.org](https://python.org/downloads/)
- **Git**: [Download from git-scm.com](https://git-scm.com/downloads)

## Installation

### 1. Install Contexis CLI

```bash
# Clone the repository
git clone https://github.com/contexis-cmp/contexis.git
cd contexis

# Install the CLI tool
make install

# Verify installation
ctx version
```

You should see output like:
```
Contexis CMP Framework v0.1.14
```

### 2. Set Up Your Environment

```bash
# Create a new project directory
mkdir my-first-ai-app
cd my-first-ai-app

# Initialize a new CMP project (local-first by default)
ctx init my-support-bot

# Navigate into your project
cd my-support-bot
```

## Your First AI Application

Let's build a customer support chatbot that can answer questions about your company policies. **No external API keys required** - everything runs locally!

### 1. Project Structure

Your project now has the following structure:

```
my-support-bot/
├── contexts/              # Agent definitions
│   └── default_agent.ctx
├── memory/               # Knowledge base
│   ├── documents/
│   └── embeddings/
├── prompts/              # Response templates
│   └── default_response.md
├── tools/               # Custom integrations
├── tests/               # Test suite
├── config/              # Configuration
│   └── environments/
└── context.lock.json    # Version locks
```

### 2. Configure Your Environment

The project comes with local-first defaults. Edit `config/environments/development.yaml` to see the configuration:

```yaml
environment: development

# Local AI Provider Configuration (default)
providers:
  local:
    model: microsoft/DialoGPT-medium  # Local model for text generation
    temperature: 0.1
    max_tokens: 1000

# Local Embeddings Configuration (default)
embeddings:
  provider: sentence-transformers
  model: all-MiniLM-L6-v2  # Local embeddings model
  dimensions: 384

# Local Vector Database Configuration (default)
vector_db:
  provider: chroma
  path: ./data/embeddings
  collection_name: development_knowledge

# Testing Configuration
testing:
  drift_threshold: 0.85
  similarity_threshold: 0.8
  max_test_duration: 300s
```

### 3. Set Up Your Environment

```bash
# Copy the example environment file
cp .env.example .env

# Install local model dependencies
pip install -r requirements.txt

# Optional: Pre-download local models (recommended for first run)
ctx models warmup
```

### 4. Create Your First Context

Let's create a customer support agent context. Edit `contexts/support_agent.ctx`:

```yaml
name: "Customer Support Agent"
version: "1.0.0"
description: "Handles customer inquiries with company knowledge"

role:
  persona: "Professional, helpful customer service representative"
  capabilities: ["answer_questions", "escalate_issues", "process_returns"]
  limitations: ["no_refunds_over_policy", "no_personal_data_sharing"]

tools:
  - name: "knowledge_search"
    uri: "mcp://search.knowledge_base"
    description: "Search company knowledge for customer questions"

guardrails:
  tone: "professional"
  format: "text"
  max_tokens: 500
  temperature: 0.1
  
memory:
  episodic: true
  max_history: 10
  privacy: "user_isolated"

testing:
  drift_threshold: 0.85
  business_rules:
    - "must_include_policy_references"
    - "no_pricing_speculation"
    - "always_verify_order_before_processing"
```

### 5. Generate a RAG System

Create a RAG system to handle knowledge-based queries:

```bash
# Generate a RAG system with local embeddings
ctx generate rag CustomerDocs --db=sqlite --embeddings=sentence-transformers
```

### 6. Add Knowledge to Memory

Create a document with your company policies. Create `memory/CustomerDocs/documents/company_policies.md`:

```markdown
# Company Policies

## Return Policy
Our return policy allows returns within 30 days of purchase with original receipt. Items must be in original condition and packaging.

## Shipping Policy
We offer free shipping on orders over $50. Standard shipping takes 3-5 business days.

## Customer Service Hours
Our customer service team is available Monday-Friday, 9 AM - 6 PM EST.

## Privacy Policy
We respect your privacy and never share personal information with third parties.
```

Now add this document to your memory:

```bash
# Add the document to memory
ctx memory ingest --provider=sqlite --component=CustomerDocs --input=memory/CustomerDocs/documents/company_policies.md
```

### 7. Create a Response Template

Edit `prompts/CustomerDocs/search_response.md`:

```markdown
# Support Response Template

Based on the customer inquiry: {{ .user_query }}

## Context Information
{{#if conversation_history}}
Previous conversation: {{ conversation_history }}
{{/if}}

## Knowledge Base Results
{{#each knowledge_results}}
- **Source**: {{ source }}
- **Content**: {{ content }}
- **Relevance**: {{ relevance_score }}
{{/each}}

## Response Guidelines
- **Tone**: Professional and helpful
- **Format**: Clear, structured response
- **Max Tokens**: 500
- **Include**: Policy references, next steps, escalation if needed

## Response

{{ .response }}

{{#if confidence_score}}
**Confidence**: {{ confidence_score }}
{{/if}}
```

### 8. Test Your System

Now let's test your support bot with local models:

```bash
# Test a query (uses local Phi-3.5-Mini model)
ctx run CustomerDocs "What is your return policy?"
```

You should see a response like:
```
Our return policy allows returns within 30 days of purchase with original receipt. Items must be in original condition and packaging.

This policy is designed to ensure customer satisfaction while maintaining the quality of our products.

**Confidence**: 0.95
```

### 9. Run Tests

```bash
# Run comprehensive tests
ctx test --all --coverage

# Monitor AI behavior drift
ctx test --drift-detection --component=CustomerDocs

# Test specific scenarios
ctx test --correctness --rules=./tests/business_rules.yaml
```

### 10. Start Development Server

```bash
# Start the development server
ctx serve --addr :8000
```

The server exposes:
- Health: `/healthz`, `/readyz`
- Version: `/version`
- Metrics: `/metrics` (Prometheus)

## Next Steps

### 1. Customize Your Agent

- **Add more capabilities** to your context
- **Create additional tools** for specific functions
- **Fine-tune the personality** and tone

### 2. Expand Your Knowledge Base

- **Add more documents** to your memory
- **Organize content** with metadata and tags
- **Regular updates** to keep information current

### 3. Improve Responses

- **Refine prompt templates** for better responses
- **Add conditional logic** for different scenarios
- **Include more context** in responses

### 4. Add Monitoring

- **Set up drift detection** to monitor AI behavior
- **Configure alerts** for performance issues
- **Track user satisfaction** metrics

## Switching to Production

When you're ready to deploy to production, simply update your configuration:

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

The same code works seamlessly across environments!

## Troubleshooting

### Common Issues

1. **Model Download Issues**
   - Run `ctx models warmup` to pre-download models
   - Check your internet connection for initial downloads
   - Ensure sufficient disk space (~3GB for local models)

2. **Memory Issues**
   - Ensure documents are properly formatted
   - Check that embeddings are being generated

3. **Response Quality**
   - Review and refine your prompt templates
   - Add more context and examples

4. **Performance Issues**
   - Check your vector database configuration
   - Monitor local model performance

### Getting Help

- **Documentation**: [docs.contexis.dev](https://docs.contexis.dev)
- **Issues**: [GitHub Issues](https://github.com/contexis-cmp/contexis/issues)
- **Discussions**: [GitHub Discussions](https://github.com/contexis-cmp/contexis/discussions)
- **Community**: [Discord](https://discord.gg/contexis)

## Advanced Features

### 1. Multi-Tenancy

Set up tenant isolation for enterprise applications:

```yaml
memory:
  privacy: "user_isolated"
  tenant_id: "${TENANT_ID}"
```

### 2. Custom Tools

Create custom tools for specific integrations:

```python
# tools/custom_api.py
class CustomAPITool:
    def __init__(self, api_key):
        self.api_key = api_key
    
    async def call(self, parameters):
        # Custom API integration
        return result
```

### 3. Advanced Testing

Set up comprehensive testing:

```bash
# Run drift detection
ctx test --drift --threshold=0.85

# Run correctness tests
ctx test --correctness --rules=./tests/business_rules.yaml

# Run performance tests
ctx test --performance --max-latency=2s
```

## Conclusion

Congratulations! You've successfully created your first AI application with Contexis. You now have:

-  A working customer support chatbot with local models
-  Knowledge base with company policies
-  Professional response templates
-  Testing and monitoring setup
-  **No external dependencies** - everything runs locally!

Continue exploring the framework to unlock its full potential for building reliable, secure, and scalable AI applications.

---

**Ready for more?** Check out our [advanced guides](https://docs.contexis.dev/guides) and [examples](https://docs.contexis.dev/examples) to learn about more advanced features and use cases.
