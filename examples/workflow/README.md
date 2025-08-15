# Workflow Pipeline Example

This example demonstrates how to build a data processing workflow using Contexis with **local models by default**.

## Overview

A workflow pipeline provides:
- **Multi-Step Processing**: Chain multiple operations together
- **Data Transformation**: Process and transform data between steps
- **Error Handling**: Robust error handling and recovery
- **Monitoring**: Track progress and performance of each step

## Quick Start

### 1. Initialize Project

```bash
# Create a new project
ctx init workflow-example
cd workflow-example

# Set up environment
cp .env.example .env
pip install -r requirements.txt
```

### 2. Generate Workflow

```bash
# Generate data processing workflow with local models
ctx generate workflow ContentPipeline --steps research,write,review
```

This creates:
- `contexts/ContentPipeline/` - Workflow configuration
- `memory/ContentPipeline/` - Data storage
- `prompts/ContentPipeline/` - Step templates
- `tools/ContentPipeline/` - Processing tools
- `tests/ContentPipeline/` - Test suite

### 3. Test Your Workflow

```bash
# Test with local models (no API keys needed!)
ctx run ContentPipeline "Create a blog post about AI trends"
ctx run ContentPipeline "Generate a product description for wireless headphones"
ctx run ContentPipeline "Write a technical tutorial about machine learning"
```

## Configuration

### Local Development (Default)

The workflow uses local models by default:

```yaml
# config/environments/development.yaml
providers:
  local:
    model: microsoft/DialoGPT-medium
    temperature: 0.1
    max_tokens: 1000

workflow:
  steps:
    - research
    - write
    - review
  error_handling: retry
  max_retries: 3
```

### Production Migration

When ready for production, switch to external providers:

```yaml
# config/environments/production.yaml
providers:
  openai:
    api_key: ${OPENAI_API_KEY}
    model: gpt-4o-mini

workflow:
  steps:
    - research
    - write
    - review
  error_handling: retry
  max_retries: 3
```

## Advanced Usage

### 1. Multi-Step Processing

```bash
# Run complete workflow
ctx run ContentPipeline "Create a blog post about AI trends"

# The workflow automatically:
# 1. Researches the topic
# 2. Writes the content
# 3. Reviews and improves the content
```

### 2. Step-by-Step Execution

```bash
# Run individual steps
ctx run ContentPipeline "Research AI trends" --step=research
ctx run ContentPipeline "Write about AI trends" --step=write
ctx run ContentPipeline "Review AI trends content" --step=review
```

### 3. Test Workflow

```bash
# Test workflow execution
ctx test --drift-detection --component=ContentPipeline

# Test specific workflow scenarios
ctx test --correctness --rules=./tests/workflow_rules.yaml
```

### 4. Start Development Server

```bash
# Start server for continuous development
ctx serve --addr :8000

# Test via HTTP API
curl -X POST http://localhost:8000/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "context": "ContentPipeline",
    "component": "ContentPipeline",
    "query": "Create a blog post about AI trends",
    "data": {"workflow_id": "workflow123"}
  }'
```

## Customization

### 1. Modify Workflow Steps

Edit `contexts/ContentPipeline/ContentPipeline.ctx`:

```yaml
name: "Content Creation Pipeline"
version: "1.0.0"
description: "Multi-step content creation workflow"

role:
  persona: "Professional content creator and editor"
  capabilities: ["research", "writing", "editing", "optimization"]
  limitations: ["no_plagiarism", "fact_checking_required"]

tools:
  - name: "research_tool"
    uri: "mcp://research.gather"
    description: "Gather information and research topics"
  - name: "writing_tool"
    uri: "mcp://content.write"
    description: "Create and write content"
  - name: "review_tool"
    uri: "mcp://content.review"
    description: "Review and improve content"

guardrails:
  tone: "professional"
  format: "markdown"
  max_tokens: 2000
  temperature: 0.2

workflow:
  steps:
    - research
    - write
    - review
  error_handling: retry
  max_retries: 3
```

### 2. Customize Step Templates

Edit `prompts/ContentPipeline/research.md`:

```markdown
# Research Step Template

**Topic:** {{ .topic }}

**Research Guidelines:**
- Gather relevant information from reliable sources
- Identify key points and insights
- Organize information logically
- Note any important statistics or data

**Research Output:**
{{ .research_output }}

**Key Findings:**
{{#each key_findings}}
- {{ . }}
{{/each}}
```

Edit `prompts/ContentPipeline/write.md`:

```markdown
# Writing Step Template

**Topic:** {{ .topic }}

**Research Summary:**
{{ .research_summary }}

**Writing Guidelines:**
- Create engaging, informative content
- Use clear, professional language
- Include relevant examples and data
- Structure content with headings and subheadings

**Content Output:**
{{ .content_output }}
```

Edit `prompts/ContentPipeline/review.md`:

```markdown
# Review Step Template

**Original Content:**
{{ .original_content }}

**Review Guidelines:**
- Check for clarity and readability
- Ensure accuracy and completeness
- Improve flow and structure
- Optimize for target audience

**Review Comments:**
{{ .review_comments }}

**Improved Content:**
{{ .improved_content }}
```

### 3. Add Custom Tools

Create `tools/ContentPipeline/research_tool.py`:

```python
class ResearchTool:
    def __init__(self, config):
        self.config = config
    
    async def research_topic(self, topic):
        # Implement research logic
        # Search web, gather information
        return {
            "topic": topic,
            "findings": ["finding1", "finding2"],
            "sources": ["source1", "source2"]
        }
```

Create `tools/ContentPipeline/writing_tool.py`:

```python
class WritingTool:
    def __init__(self, config):
        self.config = config
    
    async def write_content(self, topic, research_data):
        # Implement writing logic
        # Generate content based on research
        return {
            "content": "Generated content...",
            "word_count": 500,
            "structure": ["intro", "body", "conclusion"]
        }
```

### 4. Configure Workflow

Edit `memory/ContentPipeline/workflow_config.yaml`:

```yaml
workflow:
  steps:
    - research
    - write
    - review
  error_handling: retry
  max_retries: 3
  timeout: 300s
  parallel_execution: false
```

## Testing

### 1. Run All Tests

```bash
# Run comprehensive tests
ctx test --all --coverage
```

### 2. Test Workflow Steps

```bash
# Test individual steps
ctx test --drift-detection --component=ContentPipeline

# Test workflow integration
ctx test --correctness --rules=./tests/workflow_rules.yaml
```

### 3. Performance Testing

```bash
# Test workflow execution time
time ctx run ContentPipeline "Create a blog post about AI trends"

# Monitor resource usage
htop
```

## Workflow Examples

### Content Creation Workflow

```bash
# Input: "Create a blog post about AI trends"
ctx run ContentPipeline "Create a blog post about AI trends"

# Step 1: Research
# - Gather information about AI trends
# - Identify key topics and insights
# - Collect relevant statistics

# Step 2: Write
# - Create engaging introduction
# - Develop main content sections
# - Include examples and data

# Step 3: Review
# - Check for clarity and flow
# - Ensure accuracy and completeness
# - Optimize for readability

# Output: Complete, polished blog post
```

### Data Processing Workflow

```bash
# Input: "Process customer feedback data"
ctx run ContentPipeline "Process customer feedback data"

# Step 1: Research
# - Analyze feedback patterns
# - Identify common themes
# - Gather context information

# Step 2: Write
# - Create analysis report
# - Generate insights and recommendations
# - Structure findings clearly

# Step 3: Review
# - Verify data accuracy
# - Improve report clarity
# - Add actionable recommendations

# Output: Comprehensive analysis report
```

## Troubleshooting

### Common Issues

1. **Workflow Step Failures**
   ```bash
   # Check step configuration
   cat contexts/ContentPipeline/ContentPipeline.ctx
   
   # Test individual steps
   ctx run ContentPipeline "test" --step=research
   ```

2. **Data Flow Issues**
   ```bash
   # Check workflow configuration
   cat memory/ContentPipeline/workflow_config.yaml
   
   # Test data passing between steps
   ctx run ContentPipeline "test" --debug
   ```

3. **Performance Issues**
   ```bash
   # Check step execution times
   time ctx run ContentPipeline "test"
   
   # Monitor resource usage
   htop
   ```

### Performance Optimization

1. **Parallel Execution**
   ```yaml
   # Enable parallel step execution
   workflow:
     parallel_execution: true
   ```

2. **Caching**
   ```yaml
   # Enable step result caching
   workflow:
     caching: true
     cache_ttl: 3600s
   ```

3. **Model Warmup**
   ```bash
   # Pre-download models for faster execution
   ctx models warmup
   ```

## Next Steps

- **Add more workflow steps** for complex processing
- **Implement error recovery** for robust execution
- **Add monitoring** for workflow performance
- **Create custom tools** for specific processing needs
- **Set up scheduling** for automated workflows
- **Migrate to production** providers when ready

## Resources

- [Getting Started Guide](../docs/guides/getting-started.md)
- [CLI Reference](../docs/cli.md)
- [Model Providers](../docs/model_providers.md)
- [Memory Management](../docs/memory.md)
