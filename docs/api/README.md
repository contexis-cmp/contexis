# CMP API Reference

This document provides the complete API reference for the Contexis CMP Framework.

## Overview

The CMP Framework provides APIs for:
- **Context Management**: Creating and managing AI agent contexts
- **Memory Operations**: Storing and querying knowledge bases
- **Prompt Templating**: Generating and managing response templates
- **Tool Integration**: Integrating external tools and services

## Authentication

Optionally enable authentication by setting `CMP_AUTH_ENABLED=true`. Use Bearer tokens in the Authorization header when enabled:

```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
     https://api.contexis.dev/v1/contexts
```

## Base URL

- **Development**: `http://localhost:8080`
- **Production**: `https://api.contexis.dev`

## API Versioning

The API is versioned using URL paths: `/v1/`, `/v2/`, etc.

## Error Handling

All errors follow a consistent format:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid project name",
    "details": {
      "field": "name",
      "reason": "Contains invalid characters"
    }
  }
}
```

## Rate Limiting

If authentication is enabled, per-API key/tenant/IP token bucket limits apply. Rate limit headers include:

Rate limit headers are included in all responses:
- `X-RateLimit-Limit`: Request limit per hour
- `X-RateLimit-Remaining`: Remaining requests
- `X-RateLimit-Reset`: Time until limit resets

## Context Management API

### List Contexts

**GET** `/v1/contexts`

List all contexts for the authenticated user.

**Query Parameters:**
- `limit` (integer): Maximum number of contexts to return (default: 20)
- `offset` (integer): Number of contexts to skip (default: 0)
- `search` (string): Search contexts by name or description

**Response:**
```json
{
  "contexts": [
    {
      "id": "ctx_123",
      "name": "Customer Support Agent",
      "version": "1.0.0",
      "description": "Handles customer inquiries",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "total": 100,
    "limit": 20,
    "offset": 0,
    "has_more": true
  }
}
```

### Get Context

**GET** `/v1/contexts/{context_id}`

Get a specific context by ID.

**Response:**
```json
{
  "id": "ctx_123",
  "name": "Customer Support Agent",
  "version": "1.0.0",
  "description": "Handles customer inquiries",
  "role": {
    "persona": "Professional, helpful customer service representative",
    "capabilities": ["answer_questions", "escalate_issues"],
    "limitations": ["no_refunds_over_policy"]
  },
  "tools": [
    {
      "name": "knowledge_search",
      "uri": "mcp://search.knowledge_base",
      "description": "Search company knowledge"
    }
  ],
  "guardrails": {
    "tone": "professional",
    "format": "json",
    "max_tokens": 500,
    "temperature": 0.1
  },
  "memory": {
    "episodic": true,
    "max_history": 10,
    "privacy": "user_isolated"
  },
  "testing": {
    "drift_threshold": 0.85,
    "business_rules": ["must_include_policy_references"]
  },
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Create Context

**POST** `/v1/contexts`

Create a new context.

**Request Body:**
```json
{
  "name": "Customer Support Agent",
  "version": "1.0.0",
  "description": "Handles customer inquiries",
  "role": {
    "persona": "Professional, helpful customer service representative",
    "capabilities": ["answer_questions", "escalate_issues"],
    "limitations": ["no_refunds_over_policy"]
  },
  "tools": [
    {
      "name": "knowledge_search",
      "uri": "mcp://search.knowledge_base",
      "description": "Search company knowledge"
    }
  ],
  "guardrails": {
    "tone": "professional",
    "format": "json",
    "max_tokens": 500,
    "temperature": 0.1
  },
  "memory": {
    "episodic": true,
    "max_history": 10,
    "privacy": "user_isolated"
  },
  "testing": {
    "drift_threshold": 0.85,
    "business_rules": ["must_include_policy_references"]
  }
}
```

**Response:**
```json
{
  "id": "ctx_123",
  "name": "Customer Support Agent",
  "version": "1.0.0",
  "description": "Handles customer inquiries",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Update Context

**PUT** `/v1/contexts/{context_id}`

Update an existing context.

**Request Body:** Same as Create Context

**Response:** Updated context object

### Delete Context

**DELETE** `/v1/contexts/{context_id}`

Delete a context.

**Response:**
```json
{
  "message": "Context deleted successfully"
}
```

## Memory Management API

### List Memory Stores

**GET** `/v1/memory`

List all memory stores for the authenticated user.

**Query Parameters:**
- `limit` (integer): Maximum number of stores to return (default: 20)
- `offset` (integer): Number of stores to skip (default: 0)

**Response:**
```json
{
  "memory_stores": [
    {
      "id": "mem_123",
      "name": "Customer Knowledge",
      "document_count": 150,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "total": 5,
    "limit": 20,
    "offset": 0,
    "has_more": false
  }
}
```

### Get Memory Store

**GET** `/v1/memory/{memory_id}`

Get a specific memory store by ID.

**Response:**
```json
{
  "id": "mem_123",
  "name": "Customer Knowledge",
  "document_count": 150,
  "embeddings_model": "text-embedding-3-small",
  "vector_db": "chroma",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Create Memory Store

**POST** `/v1/memory`

Create a new memory store.

**Request Body:**
```json
{
  "name": "Customer Knowledge",
  "embeddings_model": "text-embedding-3-small",
  "vector_db": "chroma"
}
```

**Response:**
```json
{
  "id": "mem_123",
  "name": "Customer Knowledge",
  "document_count": 0,
  "embeddings_model": "text-embedding-3-small",
  "vector_db": "chroma",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Add Documents

**POST** `/v1/memory/{memory_id}/documents`

Add documents to a memory store.

**Request Body:**
```json
{
  "documents": [
    {
      "content": "Our return policy allows returns within 30 days...",
      "metadata": {
        "source": "return_policy.pdf",
        "category": "policies",
        "tags": ["returns", "policy"]
      }
    }
  ]
}
```

**Response:**
```json
{
  "added_count": 1,
  "document_ids": ["doc_123"],
  "processing_time": "2.5s"
}
```

### Search Documents

**POST** `/v1/memory/{memory_id}/search`

Search documents in a memory store.

**Request Body:**
```json
{
  "query": "What is your return policy?",
  "top_k": 5,
  "threshold": 0.8
}
```

**Response:**
```json
{
  "results": [
    {
      "id": "doc_123",
      "content": "Our return policy allows returns within 30 days...",
      "score": 0.95,
      "metadata": {
        "source": "return_policy.pdf",
        "category": "policies"
      }
    }
  ],
  "total_results": 1,
  "search_time": "0.15s"
}
```

## Prompt Management API

### List Prompts

**GET** `/v1/prompts`

List all prompts for the authenticated user.

**Response:**
```json
{
  "prompts": [
    {
      "id": "prompt_123",
      "name": "Support Response",
      "version": "1.0.0",
      "description": "Customer support response template",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### Get Prompt

**GET** `/v1/prompts/{prompt_id}`

Get a specific prompt by ID.

**Response:**
```json
{
  "id": "prompt_123",
  "name": "Support Response",
  "version": "1.0.0",
  "description": "Customer support response template",
  "template": "# Response Template\n\nBased on the user query: {{ user_query }}\n\n## Response\n\n{{ response_text }}",
  "variables": ["user_query", "response_text"],
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Create Prompt

**POST** `/v1/prompts`

Create a new prompt.

**Request Body:**
```json
{
  "name": "Support Response",
  "version": "1.0.0",
  "description": "Customer support response template",
  "template": "# Response Template\n\nBased on the user query: {{ user_query }}\n\n## Response\n\n{{ response_text }}"
}
```

**Response:**
```json
{
  "id": "prompt_123",
  "name": "Support Response",
  "version": "1.0.0",
  "description": "Customer support response template",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Render Prompt

**POST** `/v1/prompts/{prompt_id}/render`

Render a prompt template with variables.

**Request Body:**
```json
{
  "variables": {
    "user_query": "What is your return policy?",
    "response_text": "Our return policy allows returns within 30 days of purchase."
  }
}
```

**Response:**
```json
{
  "rendered": "# Response Template\n\nBased on the user query: What is your return policy?\n\n## Response\n\nOur return policy allows returns within 30 days of purchase."
}
```

## Tool Management API

### List Tools

**GET** `/v1/tools`

List all available tools.

**Response:**
```json
{
  "tools": [
    {
      "id": "tool_123",
      "name": "knowledge_search",
      "uri": "mcp://search.knowledge_base",
      "description": "Search company knowledge base",
      "version": "1.0.0",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### Get Tool

**GET** `/v1/tools/{tool_id}`

Get a specific tool by ID.

**Response:**
```json
{
  "id": "tool_123",
  "name": "knowledge_search",
  "uri": "mcp://search.knowledge_base",
  "description": "Search company knowledge base",
  "version": "1.0.0",
  "parameters": {
    "query": {
      "type": "string",
      "required": true,
      "description": "Search query"
    },
    "top_k": {
      "type": "integer",
      "required": false,
      "default": 5,
      "description": "Number of results to return"
    }
  },
  "created_at": "2024-01-01T00:00:00Z"
}
```

### Call Tool

**POST** `/v1/tools/{tool_id}/call`

Call a tool with parameters.

**Request Body:**
```json
{
  "parameters": {
    "query": "return policy",
    "top_k": 3
  }
}
```

**Response:**
```json
{
  "result": [
    {
      "content": "Our return policy allows returns within 30 days...",
      "score": 0.95,
      "source": "return_policy.pdf"
    }
  ],
  "execution_time": "0.25s"
}
```

## Testing API

### Run Tests

**POST** `/v1/tests/run`

Run tests for a context or project.

**Request Body:**
```json
{
  "context_id": "ctx_123",
  "test_types": ["drift", "correctness", "performance"],
  "options": {
    "drift_threshold": 0.85,
    "max_latency": "2s"
  }
}
```

**Response:**
```json
{
  "test_id": "test_123",
  "status": "running",
  "created_at": "2024-01-01T00:00:00Z"
}
```

### Get Test Results

**GET** `/v1/tests/{test_id}`

Get test results.

**Response:**
```json
{
  "id": "test_123",
  "status": "completed",
  "results": {
    "drift": {
      "passed": true,
      "score": 0.92,
      "threshold": 0.85
    },
    "correctness": {
      "passed": true,
      "rules_tested": 5,
      "rules_passed": 5
    },
    "performance": {
      "passed": true,
      "avg_latency": "1.2s",
      "max_latency": "1.8s"
    }
  },
  "created_at": "2024-01-01T00:00:00Z",
  "completed_at": "2024-01-01T00:00:05Z"
}
```

## Webhooks

### Create Webhook

**POST** `/v1/webhooks`

Create a webhook for notifications.

**Request Body:**
```json
{
  "url": "https://your-app.com/webhook",
  "events": ["context.created", "memory.updated", "test.completed"],
  "secret": "your-webhook-secret"
}
```

**Response:**
```json
{
  "id": "webhook_123",
  "url": "https://your-app.com/webhook",
  "events": ["context.created", "memory.updated", "test.completed"],
  "created_at": "2024-01-01T00:00:00Z"
}
```

## SDKs and Libraries

### Go SDK

```go
package main

import (
    "context"
    "log"
    
    "github.com/contexis-cmp/contexis-go-sdk"
)

func main() {
    client := cmp.NewClient("your-api-key")
    
    ctx, err := client.GetContext(context.Background(), "ctx_123")
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Context: %s", ctx.Name)
}
```

### Python SDK

```python
import contexis

client = contexis.Client("your-api-key")

context = client.get_context("ctx_123")
print(f"Context: {context.name}")

# Search memory
results = client.search_memory("mem_123", "return policy")
for result in results:
    print(f"Score: {result.score}, Content: {result.content}")
```

## Error Codes

| Code | Description |
|------|-------------|
| `VALIDATION_ERROR` | Request validation failed |
| `NOT_FOUND` | Resource not found |
| `UNAUTHORIZED` | Authentication required |
| `FORBIDDEN` | Access denied |
| `RATE_LIMITED` | Rate limit exceeded |
| `INTERNAL_ERROR` | Internal server error |

## Support

For API support:
- **Documentation**: [docs.contexis.dev](https://docs.contexis.dev)
- **Email**: api@contexis.dev
- **Discord**: [Contexis Community](https://discord.gg/contexis)
