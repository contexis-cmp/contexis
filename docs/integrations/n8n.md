# n8n Integration Guide

This guide shows how to integrate Contexis with n8n for workflow automation and AI-powered business processes.

## Overview

n8n is a powerful workflow automation platform that can integrate with Contexis through:
- **HTTP API calls** to Contexis endpoints
- **Webhook triggers** for real-time events
- **Data transformation** between systems
- **Error handling** and retry logic

## Prerequisites

- n8n instance running (self-hosted or cloud)
- Contexis server running with API access
- API key for Contexis (if authentication enabled)

## Basic Integration

### 1. HTTP Request Node

The primary integration method is using n8n's HTTP Request node to call Contexis APIs.

#### Chat API Integration

```javascript
// n8n HTTP Request node configuration
{
  "method": "POST",
  "url": "http://localhost:8080/api/v1/chat",
  "headers": {
    "Content-Type": "application/json",
    "Authorization": "Bearer YOUR_API_KEY"
  },
  "body": {
    "context": "CustomerDocs",
    "query": "{{ $json.customer_question }}",
    "tenant_id": "{{ $json.tenant_id }}",
    "top_k": 5
  }
}
```

#### Complete Workflow Example

```javascript
// n8n workflow: Customer Support Automation
{
  "nodes": [
    {
      "name": "Customer Inquiry",
      "type": "n8n-nodes-base.webhook",
      "parameters": {
        "httpMethod": "POST",
        "path": "customer-inquiry",
        "responseMode": "responseNode"
      }
    },
    {
      "name": "Contexis Chat",
      "type": "n8n-nodes-base.httpRequest",
      "parameters": {
        "method": "POST",
        "url": "http://localhost:8080/api/v1/chat",
        "headers": {
          "Content-Type": "application/json",
          "Authorization": "Bearer {{ $env.CONTEXIS_API_KEY }}"
        },
        "body": {
          "context": "CustomerDocs",
          "query": "{{ $json.question }}",
          "tenant_id": "{{ $json.customer_id }}"
        }
      }
    },
    {
      "name": "Process Response",
      "type": "n8n-nodes-base.set",
      "parameters": {
        "values": {
          "customer_id": "{{ $json.customer_id }}",
          "ai_response": "{{ $('Contexis Chat').item.json.rendered }}",
          "timestamp": "{{ $now }}"
        }
      }
    },
    {
      "name": "Send Response",
      "type": "n8n-nodes-base.httpRequest",
      "parameters": {
        "method": "POST",
        "url": "{{ $json.webhook_url }}",
        "body": {
          "response": "{{ $json.ai_response }}",
          "customer_id": "{{ $json.customer_id }}"
        }
      }
    }
  ]
}
```

### 2. Webhook Integration

Set up Contexis webhooks to trigger n8n workflows on specific events.

#### Contexis Webhook Configuration

```bash
# Create webhook in Contexis
curl -X POST http://localhost:8080/api/v1/webhooks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -d '{
    "url": "https://your-n8n-instance.com/webhook/contexis-events",
    "events": ["memory.updated", "context.created", "test.completed"],
    "secret": "your-webhook-secret"
  }'
```

#### n8n Webhook Receiver

```javascript
// n8n webhook node configuration
{
  "name": "Contexis Webhook",
  "type": "n8n-nodes-base.webhook",
  "parameters": {
    "httpMethod": "POST",
    "path": "contexis-events",
    "responseMode": "responseNode",
    "options": {
      "responseHeaders": {
        "X-Webhook-Secret": "{{ $env.WEBHOOK_SECRET }}"
      }
    }
  }
}
```

## Advanced Workflows

### 1. Document Processing Pipeline

```javascript
// n8n workflow: Document Processing with Contexis
{
  "nodes": [
    {
      "name": "Document Upload",
      "type": "n8n-nodes-base.webhook",
      "parameters": {
        "httpMethod": "POST",
        "path": "document-upload"
      }
    },
    {
      "name": "Extract Text",
      "type": "n8n-nodes-base.httpRequest",
      "parameters": {
        "method": "POST",
        "url": "https://api.example.com/extract-text",
        "body": {
          "file": "{{ $json.file_url }}"
        }
      }
    },
    {
      "name": "Ingest to Contexis",
      "type": "n8n-nodes-base.httpRequest",
      "parameters": {
        "method": "POST",
        "url": "http://localhost:8080/api/v1/memory/ingest",
        "headers": {
          "Authorization": "Bearer {{ $env.CONTEXIS_API_KEY }}"
        },
        "body": {
          "provider": "sqlite",
          "component": "CustomerDocs",
          "input": "{{ $('Extract Text').item.json.text }}"
        }
      }
    },
    {
      "name": "Generate Summary",
      "type": "n8n-nodes-base.httpRequest",
      "parameters": {
        "method": "POST",
        "url": "http://localhost:8080/api/v1/chat",
        "headers": {
          "Authorization": "Bearer {{ $env.CONTEXIS_API_KEY }}"
        },
        "body": {
          "context": "CustomerDocs",
          "query": "Summarize the key points from this document",
          "data": {
            "document_content": "{{ $('Extract Text').item.json.text }}"
          }
        }
      }
    },
    {
      "name": "Send Notification",
      "type": "n8n-nodes-base.httpRequest",
      "parameters": {
        "method": "POST",
        "url": "{{ $env.SLACK_WEBHOOK }}",
        "body": {
          "text": "Document processed and summarized: {{ $('Generate Summary').item.json.rendered }}"
        }
      }
    }
  ]
}
```

### 2. Multi-Step AI Workflow

```javascript
// n8n workflow: Multi-Step AI Processing
{
  "nodes": [
    {
      "name": "Start Process",
      "type": "n8n-nodes-base.webhook",
      "parameters": {
        "httpMethod": "POST",
        "path": "ai-workflow"
      }
    },
    {
      "name": "Research Phase",
      "type": "n8n-nodes-base.httpRequest",
      "parameters": {
        "method": "POST",
        "url": "http://localhost:8080/api/v1/chat",
        "headers": {
          "Authorization": "Bearer {{ $env.CONTEXIS_API_KEY }}"
        },
        "body": {
          "context": "ResearchBot",
          "query": "Research {{ $json.topic }}",
          "data": {
            "research_depth": "{{ $json.depth }}"
          }
        }
      }
    },
    {
      "name": "Analysis Phase",
      "type": "n8n-nodes-base.httpRequest",
      "parameters": {
        "method": "POST",
        "url": "http://localhost:8080/api/v1/chat",
        "headers": {
          "Authorization": "Bearer {{ $env.CONTEXIS_API_KEY }}"
        },
        "body": {
          "context": "AnalysisBot",
          "query": "Analyze the research findings",
          "data": {
            "research_data": "{{ $('Research Phase').item.json.rendered }}"
          }
        }
      }
    },
    {
      "name": "Report Generation",
      "type": "n8n-nodes-base.httpRequest",
      "parameters": {
        "method": "POST",
        "url": "http://localhost:8080/api/v1/chat",
        "headers": {
          "Authorization": "Bearer {{ $env.CONTEXIS_API_KEY }}"
        },
        "body": {
          "context": "ReportBot",
          "query": "Generate a comprehensive report",
          "data": {
            "analysis": "{{ $('Analysis Phase').item.json.rendered }}",
            "format": "{{ $json.report_format }}"
          }
        }
      }
    }
  ]
}
```

## Error Handling and Retry Logic

### 1. Retry Configuration

```javascript
// n8n HTTP Request with retry logic
{
  "name": "Contexis Chat with Retry",
  "type": "n8n-nodes-base.httpRequest",
  "parameters": {
    "method": "POST",
    "url": "http://localhost:8080/api/v1/chat",
    "headers": {
      "Authorization": "Bearer {{ $env.CONTEXIS_API_KEY }}"
    },
    "body": {
      "context": "CustomerDocs",
      "query": "{{ $json.question }}"
    },
    "options": {
      "timeout": 30000,
      "retry": {
        "enabled": true,
        "maxAttempts": 3,
        "waitTime": 5000
      }
    }
  }
}
```

### 2. Error Handling Workflow

```javascript
// n8n workflow with error handling
{
  "nodes": [
    {
      "name": "Try Contexis",
      "type": "n8n-nodes-base.httpRequest",
      "parameters": {
        "method": "POST",
        "url": "http://localhost:8080/api/v1/chat",
        "body": {
          "context": "CustomerDocs",
          "query": "{{ $json.question }}"
        }
      }
    },
    {
      "name": "Handle Error",
      "type": "n8n-nodes-base.set",
      "parameters": {
        "values": {
          "error": true,
          "fallback_response": "I apologize, but I'm unable to process your request at the moment. Please try again later.",
          "original_error": "{{ $json.error }}"
        }
      },
      "continueOnFail": true
    },
    {
      "name": "Send Response",
      "type": "n8n-nodes-base.set",
      "parameters": {
        "values": {
          "response": "{{ $('Try Contexis').item.json.rendered || $('Handle Error').item.json.fallback_response }}",
          "success": "{{ $('Try Contexis').item.json.rendered ? true : false }}"
        }
      }
    }
  ]
}
```

## Data Transformation

### 1. Input Transformation

```javascript
// n8n Set node for input transformation
{
  "name": "Transform Input",
  "type": "n8n-nodes-base.set",
  "parameters": {
    "values": {
      "context": "{{ $json.customer_type === 'enterprise' ? 'EnterpriseDocs' : 'CustomerDocs' }}",
      "query": "{{ $json.message }}",
      "tenant_id": "{{ $json.customer_id }}",
      "top_k": "{{ $json.search_depth || 5 }}"
    }
  }
}
```

### 2. Response Processing

```javascript
// n8n Set node for response processing
{
  "name": "Process Response",
  "type": "n8n-nodes-base.set",
  "parameters": {
    "values": {
      "customer_response": {
        "message": "{{ $('Contexis Chat').item.json.rendered }}",
        "confidence": "{{ $('Contexis Chat').item.json.confidence || 0.8 }}",
        "sources": "{{ $('Contexis Chat').item.json.sources || [] }}"
      },
      "metadata": {
        "processing_time": "{{ $now.diff($('Start').item.json.timestamp).toSeconds() }}",
        "model_used": "{{ $('Contexis Chat').item.json.model || 'local' }}"
      }
    }
  }
}
```

## Environment Configuration

### 1. n8n Environment Variables

```bash
# n8n environment variables
CONTEXIS_API_KEY=your_api_key_here
CONTEXIS_BASE_URL=http://localhost:8080
WEBHOOK_SECRET=your_webhook_secret
SLACK_WEBHOOK=https://hooks.slack.com/services/your/webhook/url
```

### 2. Contexis Configuration

```yaml
# config/environments/development.yaml
api:
  enabled: true
  port: 8080
  cors:
    enabled: true
    origins: ["https://your-n8n-instance.com"]
  
webhooks:
  enabled: true
  max_retries: 3
  timeout: 30s
```

## Monitoring and Logging

### 1. Workflow Monitoring

```javascript
// n8n workflow with monitoring
{
  "nodes": [
    {
      "name": "Log Start",
      "type": "n8n-nodes-base.httpRequest",
      "parameters": {
        "method": "POST",
        "url": "{{ $env.MONITORING_WEBHOOK }}",
        "body": {
          "event": "workflow_started",
          "workflow_id": "{{ $workflow.id }}",
          "timestamp": "{{ $now }}"
        }
      }
    },
    {
      "name": "Contexis Call",
      "type": "n8n-nodes-base.httpRequest",
      "parameters": {
        "method": "POST",
        "url": "http://localhost:8080/api/v1/chat",
        "body": {
          "context": "CustomerDocs",
          "query": "{{ $json.question }}"
        }
      }
    },
    {
      "name": "Log Success",
      "type": "n8n-nodes-base.httpRequest",
      "parameters": {
        "method": "POST",
        "url": "{{ $env.MONITORING_WEBHOOK }}",
        "body": {
          "event": "workflow_completed",
          "workflow_id": "{{ $workflow.id }}",
          "success": true,
          "timestamp": "{{ $now }}"
        }
      }
    }
  ]
}
```

### 2. Performance Tracking

```javascript
// n8n workflow with performance tracking
{
  "name": "Track Performance",
  "type": "n8n-nodes-base.set",
  "parameters": {
    "values": {
      "start_time": "{{ $now }}",
      "request_id": "{{ $json.request_id }}"
    }
  }
},
{
  "name": "Contexis Request",
  "type": "n8n-nodes-base.httpRequest",
  "parameters": {
    "method": "POST",
    "url": "http://localhost:8080/api/v1/chat",
    "body": {
      "context": "CustomerDocs",
      "query": "{{ $json.question }}"
    }
  }
},
{
  "name": "Calculate Performance",
  "type": "n8n-nodes-base.set",
  "parameters": {
    "values": {
      "processing_time": "{{ $now.diff($('Track Performance').item.json.start_time).toMilliseconds() }}",
      "response_length": "{{ $('Contexis Request').item.json.rendered.length }}",
      "success": true
    }
  }
}
```

## Best Practices

### 1. API Key Management

- Store API keys in n8n environment variables
- Use different keys for different environments
- Rotate keys regularly
- Monitor API key usage

### 2. Error Handling

- Implement retry logic for transient failures
- Provide fallback responses for critical workflows
- Log errors for debugging
- Set up alerts for repeated failures

### 3. Performance Optimization

- Use connection pooling for HTTP requests
- Implement caching for repeated queries
- Monitor response times
- Set appropriate timeouts

### 4. Security

- Validate webhook signatures
- Use HTTPS for all communications
- Implement rate limiting
- Monitor for suspicious activity

## Troubleshooting

### Common Issues

1. **Connection Timeouts**
   - Check Contexis server status
   - Verify network connectivity
   - Increase timeout values

2. **Authentication Errors**
   - Verify API key is correct
   - Check API key permissions
   - Ensure proper header format

3. **Webhook Failures**
   - Verify webhook URL is accessible
   - Check webhook secret configuration
   - Monitor webhook delivery logs

### Debug Workflow

```javascript
// n8n debug workflow
{
  "nodes": [
    {
      "name": "Debug Info",
      "type": "n8n-nodes-base.set",
      "parameters": {
        "values": {
          "timestamp": "{{ $now }}",
          "workflow_id": "{{ $workflow.id }}",
          "node_id": "{{ $node.id }}",
          "input_data": "{{ JSON.stringify($json) }}"
        }
      }
    },
    {
      "name": "Test Contexis",
      "type": "n8n-nodes-base.httpRequest",
      "parameters": {
        "method": "GET",
        "url": "http://localhost:8080/healthz"
      }
    },
    {
      "name": "Log Debug",
      "type": "n8n-nodes-base.httpRequest",
      "parameters": {
        "method": "POST",
        "url": "{{ $env.DEBUG_WEBHOOK }}",
        "body": {
          "debug_info": "{{ $('Debug Info').item.json }}",
          "health_check": "{{ $('Test Contexis').item.json }}"
        }
      }
    }
  ]
}
```

## Support

For n8n integration support:
- **Documentation**: [docs.contexis.dev/integrations/n8n](https://docs.contexis.dev/integrations/n8n)
- **Community**: [Discord Integration Channel](https://discord.gg/contexis)
- **Issues**: [GitHub Issues](https://github.com/contexis-cmp/contexis/issues)
