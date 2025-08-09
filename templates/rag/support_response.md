# Support Response Template

Based on the customer inquiry: {{ user_query }}

## Context Information
- Customer ID: {{ customer_id }}
- Order History: {{ order_history }}
- Previous Interactions: {{ conversation_history }}

## Knowledge Base Results
{{#each knowledge_results}}
- **Source**: {{ source }}
- **Content**: {{ content }}
- **Relevance**: {{ relevance_score }}
{{/each}}

## Response Guidelines
- **Tone**: Professional and helpful
- **Format**: Structured JSON response
- **Max Tokens**: 500
- **Include**: Policy references, next steps, escalation if needed

## Response Template

```json
{
  "response": "{{ response_text }}",
  "confidence": {{ confidence_score }},
  "policy_references": ["{{ policy_refs }}"],
  "next_actions": ["{{ next_actions }}"],
  "escalation_needed": {{ escalation_required }},
  "context_sha": "{{ context_sha }}"
}
```

## Example Response

For a return policy question:
```json
{
  "response": "Our return policy allows returns within 30 days of purchase with original receipt. I can help you process your return right away.",
  "confidence": 0.95,
  "policy_references": ["return_policy_2024"],
  "next_actions": ["verify_purchase_date", "check_receipt"],
  "escalation_needed": false,
  "context_sha": "sha256:abc123..."
}
``` 