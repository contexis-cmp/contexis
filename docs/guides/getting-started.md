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
Contexis CMP Framework v0.1.0
```

### 2. Set Up Your Environment

```bash
# Create a new project directory
mkdir my-first-ai-app
cd my-first-ai-app

# Initialize a new CMP project
ctx init my-support-bot

# Navigate into your project
cd my-support-bot
```

## Your First AI Application

Let's build a customer support chatbot that can answer questions about your company policies.

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

Edit `config/environments/development.yaml` to set up your AI provider:

```yaml
environment: development

# AI Provider Configuration
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
  provider: chroma
  path: ./data/embeddings
  collection_name: development_knowledge

# Testing Configuration
testing:
  drift_threshold: 0.85
  similarity_threshold: 0.8
  max_test_duration: 300s
```

### 3. Set Up Your API Keys

Create a `.env` file in your project root:

```bash
# Copy the example environment file
cp .env.example .env

# Edit the file with your API keys
nano .env
```

Add your API keys:
```bash
# AI Provider Keys
OPENAI_API_KEY=your_openai_api_key_here

# Application Settings
LOG_LEVEL=debug
ENVIRONMENT=development
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

### 5. Add Knowledge to Memory

Create a document with your company policies. Create `memory/documents/company_policies.md`:

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
ctx memory add --file=memory/documents/company_policies.md
```

### 6. Create a Response Template

Edit `prompts/support_response.md`:

```markdown
# Support Response Template

Based on the customer inquiry: {{ user_query }}

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

{{ response_text }}

{{#if confidence_score}}
**Confidence**: {{ confidence_score }}
{{/if}}
```

### 7. Run and Monitor

Now let's test your support bot:

```bash
# Run Go test suites with coverage and JUnit output
ctx test --all --coverage --junit --out tests/reports

# Run drift detection for your knowledge base
ctx test --drift-detection --component CustomerDocs --semantic --out tests/reports

### 8. Health and Metrics

The server exposes:
- Health: `/healthz`, `/readyz`
- Version: `/version`
- Metrics: `/metrics` (Prometheus)

# Test a query (example)
ctx run query "What is your return policy?"
```

You should see a response like:
```
Our return policy allows returns within 30 days of purchase with original receipt. Items must be in original condition and packaging.

This policy is designed to ensure customer satisfaction while maintaining the quality of our products.

**Confidence**: 0.95
```

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

## Troubleshooting

### Common Issues

1. **API Key Errors**
   - Make sure your API keys are correctly set in `.env`
   - Verify the keys are valid and have sufficient credits

2. **Memory Issues**
   - Ensure documents are properly formatted
   - Check that embeddings are being generated

3. **Response Quality**
   - Review and refine your prompt templates
   - Add more context and examples

4. **Performance Issues**
   - Check your vector database configuration
   - Monitor API rate limits

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

- ✅ A working customer support chatbot
- ✅ Knowledge base with company policies
- ✅ Professional response templates
- ✅ Testing and monitoring setup

Continue exploring the framework to unlock its full potential for building reliable, secure, and scalable AI applications.

---

**Ready for more?** Check out our [advanced guides](https://docs.contexis.dev/guides) and [examples](https://docs.contexis.dev/examples) to learn about more advanced features and use cases.
