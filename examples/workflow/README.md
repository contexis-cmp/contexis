# Workflow Pipeline Example

This example demonstrates how to build automated workflow pipelines using the Contexis CMP framework.

## Overview

This workflow pipeline is designed to process data through multiple stages with AI-powered decision making. It can:

- Extract data from various sources
- Transform and enrich data with AI
- Load processed data to destinations
- Handle errors and retries
- Monitor pipeline performance

## Quick Start

```bash
# Initialize the project
ctx init workflow-pipeline
cd workflow-pipeline

# Generate the workflow
ctx generate workflow DataProcessor --steps=extract,transform,load

# Configure data sources
ctx workflow add-source --name=api --type=http --config=./config/sources/api.yaml
ctx workflow add-source --name=database --type=postgresql --config=./config/sources/database.yaml

# Configure destinations
ctx workflow add-destination --name=warehouse --type=bigquery --config=./config/destinations/warehouse.yaml

# Test the workflow
ctx test

# Run the pipeline
ctx run DataProcessor "Process data"
```

## Generated Structure

```
workflow-pipeline/
├── contexts/
│   └── data_processor.ctx          # Workflow agent definition
├── memory/
│   ├── processed_data/             # Processed data storage
│   └── workflow_state/             # Workflow state management
├── prompts/
│   ├── data_extraction.md          # Data extraction prompts
│   ├── data_transformation.md      # Data transformation prompts
│   └── data_validation.md          # Data validation prompts
├── tools/
│   ├── extractors/                 # Data extraction tools
│   │   ├── api_extractor.py       # API data extraction
│   │   └── db_extractor.py        # Database extraction
│   ├── transformers/               # Data transformation tools
│   │   ├── ai_enricher.py         # AI-powered enrichment
│   │   └── data_cleaner.py        # Data cleaning
│   └── loaders/                    # Data loading tools
│       ├── warehouse_loader.py     # Data warehouse loading
│       └── api_loader.py          # API data loading
├── tests/
│   ├── workflow_integration.py     # Workflow integration tests
│   └── data_quality.py            # Data quality tests
└── context.lock.json               # Version locks
```

## Key Features

### Data Processing
- **Multi-source Extraction**: Extract from APIs, databases, files
- **AI-powered Transformation**: Enrich data with AI insights
- **Quality Validation**: Ensure data quality and consistency
- **Error Handling**: Robust error handling and retries

### Workflow Management
- **State Tracking**: Track workflow execution state
- **Progress Monitoring**: Real-time progress updates
- **Dependency Management**: Handle step dependencies
- **Parallel Processing**: Parallel execution where possible

### Integration
- **API Integration**: Connect to external APIs
- **Database Support**: Support for multiple databases
- **File Processing**: Handle various file formats
- **Real-time Updates**: Live workflow updates

## Configuration

### Workflow Context

Edit `contexts/data_processor.ctx`:

```yaml
name: "Data Processor"
version: "1.0.0"
description: "AI-powered data processing workflow"

role:
  persona: "Efficient data processing agent"
  capabilities: ["data_extraction", "data_transformation", "data_loading", "quality_validation"]
  limitations: ["no_sensitive_data", "no_pii_processing"]

tools:
  - name: "api_extractor"
    uri: "mcp://extract.api"
    description: "Extract data from APIs"
  - name: "db_extractor"
    uri: "mcp://extract.database"
    description: "Extract data from databases"
  - name: "ai_enricher"
    uri: "mcp://transform.ai"
    description: "Enrich data with AI insights"
  - name: "data_validator"
    uri: "mcp://validate.data"
    description: "Validate data quality"
  - name: "warehouse_loader"
    uri: "mcp://load.warehouse"
    description: "Load data to warehouse"

guardrails:
  tone: "technical"
  format: "json"
  max_tokens: 2000
  temperature: 0.1
  
memory:
  episodic: true
  max_history: 100
  privacy: "data_isolated"

testing:
  drift_threshold: 0.85
  business_rules:
    - "ensure_data_quality"
    - "maintain_data_lineage"
    - "handle_errors_gracefully"
```

### Workflow Steps

#### Data Extraction

Edit `prompts/data_extraction.md`:

```markdown
# Data Extraction Template

**Source**: {{ source_name }}
**Configuration**: {{ source_config }}
**Extraction Parameters**: {{ extraction_params }}

## Extraction Guidelines
- Extract data efficiently
- Handle rate limits
- Validate extracted data
- Log extraction progress

## Extracted Data

{{ extracted_data }}

**Records Processed**: {{ record_count }}
**Extraction Time**: {{ extraction_time }}
**Success Rate**: {{ success_rate }}
```

#### Data Transformation

Edit `prompts/data_transformation.md`:

```markdown
# Data Transformation Template

**Input Data**: {{ input_data }}
**Transformation Rules**: {{ transformation_rules }}
**AI Enrichment**: {{ ai_enrichment }}

## Transformation Guidelines
- Apply transformation rules
- Enrich with AI insights
- Maintain data quality
- Handle transformation errors

## Transformed Data

{{ transformed_data }}

**Transformations Applied**: {{ transformations_applied }}
**AI Insights Added**: {{ ai_insights }}
**Quality Score**: {{ quality_score }}
```

## Usage Examples

### Basic Workflow

```bash
# Run complete workflow
ctx run DataProcessor "Process data"

# Run specific step
ctx run DataProcessor "Extract data" --data '{"step":"extract"}'

# Run with custom parameters
ctx run DataProcessor "Process data" --data '{"batch_size": 1000}'
```

### Monitoring and Debugging

```bash
# Monitor workflow progress
ctx monitor workflow --name=DataProcessor

# View workflow logs
ctx logs --workflow=DataProcessor

# Debug workflow issues
ctx run DataProcessor "Debug transform step" --data '{"step":"transform","debug":true}' --debug
```

## Testing

### Workflow Testing

```bash
# Test complete workflow
ctx test --workflow=DataProcessor

# Test specific steps
ctx test --workflow=DataProcessor --step=extract
ctx test --workflow=DataProcessor --step=transform
ctx test --workflow=DataProcessor --step=load

# Test data quality
ctx test --workflow=DataProcessor --type=quality
```

### Automated Testing

The generated test suite includes:

```python
# tests/workflow_integration.py
def test_data_extraction():
    """Test data extraction step"""
    workflow = DataProcessorWorkflow()
    
    # Test API extraction
    api_data = workflow.extract_data("api", {"endpoint": "/users"})
    assert len(api_data) > 0
    
    # Test database extraction
    db_data = workflow.extract_data("database", {"table": "users"})
    assert len(db_data) > 0

def test_data_transformation():
    """Test data transformation step"""
    workflow = DataProcessorWorkflow()
    
    # Test AI enrichment
    enriched_data = workflow.transform_data(test_data)
    assert "ai_insights" in enriched_data
    
    # Test data quality validation
    quality_score = workflow.validate_quality(enriched_data)
    assert quality_score > 0.8

def test_workflow_integration():
    """Test complete workflow integration"""
    workflow = DataProcessorWorkflow()
    
    # Run complete workflow
    result = workflow.run()
    
    # Verify results
    assert result.success
    assert result.records_processed > 0
    assert result.quality_score > 0.8
```

## Deployment

### Local Development

```bash
# Start development environment
ctx dev --port=8080

# Test workflow locally
ctx run DataProcessor "Process data" --data '{"environment":"development"}'
```

### Production Deployment

```bash
# Build for production
ctx build --environment=production

# Deploy to container
ctx deploy --target=docker --image=workflow-pipeline:latest

# Deploy to Kubernetes
ctx deploy --target=kubernetes --namespace=workflows
```

## Monitoring

### Workflow Metrics

```bash
# View workflow metrics
ctx metrics --workflow=DataProcessor

# Monitor performance
ctx metrics --type=performance --workflow=DataProcessor

# Track data quality
ctx metrics --type=quality --workflow=DataProcessor
```

### Alerts

```bash
# Set up alerts
ctx alerts create --workflow=DataProcessor --condition="quality_score < 0.8"

# View active alerts
ctx alerts list --workflow=DataProcessor
```

## Customization

### Adding New Sources

1. **Create Source Configuration**

```yaml
# config/sources/custom_api.yaml
name: "custom_api"
type: "http"
endpoint: "https://api.example.com/data"
headers:
  Authorization: "Bearer ${API_KEY}"
parameters:
  limit: 1000
  offset: 0
```

2. **Add Source to Workflow**

```bash
ctx workflow add-source --name=custom_api --type=http --config=./config/sources/custom_api.yaml
```

3. **Test New Source**

```bash
ctx test --source=custom_api
```

### Custom Transformations

Create custom transformation tools:

```python
# tools/transformers/custom_transformer.py
class CustomTransformer:
    def __init__(self, config):
        self.config = config
    
    async def transform(self, data):
        # Custom transformation logic
        transformed_data = self.apply_transformations(data)
        return transformed_data
```

### Error Handling

Configure error handling strategies:

```yaml
# config/error_handling.yaml
retry_strategy:
  max_retries: 3
  backoff_multiplier: 2
  initial_delay: 1s

error_actions:
  - condition: "extraction_failed"
    action: "notify_admin"
  - condition: "quality_score < 0.8"
    action: "stop_pipeline"
```

## Troubleshooting

### Common Issues

1. **Data Extraction Failures**
   - Check API endpoints and authentication
   - Verify database connections
   - Review rate limits and quotas

2. **Transformation Errors**
   - Validate transformation rules
   - Check data formats
   - Review AI enrichment configuration

3. **Load Failures**
   - Verify destination configurations
   - Check permissions and quotas
   - Review data format requirements

### Debug Mode

```bash
# Enable debug logging
ctx run workflow --name=DataProcessor --debug

# View detailed logs
ctx logs --level=debug --workflow=DataProcessor
```

## Performance Optimization

### Parallel Processing

```bash
# Enable parallel processing
ctx run DataProcessor "Process data" --data '{"parallel":true,"max_workers":4}'

# Configure batch processing
ctx run DataProcessor "Process data" --data '{"batch_size":1000}'
```

### Caching

```bash
# Enable caching
ctx workflow cache --enable --ttl=1h

# Clear cache
ctx workflow cache --clear
```

## Next Steps

- Add more data sources
- Implement advanced transformations
- Add real-time processing
- Implement data lineage tracking
- Add workflow orchestration

## Resources

- [Workflow Design Guide](https://docs.contexis.dev/guides/workflow-design)
- [Data Processing Patterns](https://docs.contexis.dev/guides/data-processing)
- [Error Handling Strategies](https://docs.contexis.dev/guides/error-handling)
- [Performance Optimization](https://docs.contexis.dev/guides/performance)
