# Runtime and API

Start the server or run queries:
```bash
# Run a query directly (starts server temporarily)
ctx run SupportBot "What is your return policy?"

# Start server manually for continuous use
ctx serve --addr :8000
```

Health endpoints:
- GET `/healthz` → 200 ok
- GET `/readyz` → 200 ready
- GET `/version` → framework version
- GET `/metrics` → Prometheus metrics

Chat API:
- POST `/api/v1/chat`
- Request body:
```json
{
  "tenant_id": "",
  "context": "HRBot",
  "component": "HRBot",
  "query": "What is our parental leave policy?",
  "top_k": 5,
  "data": {},
  "prompt_file": "agent_response.md"
}
```
- Response:
```json
{ "rendered": "...model output or rendered prompt..." }
```

Model provider:
- If `HF_TOKEN` and `HF_MODEL_ID` are set, the server calls Hugging Face Inference API with the rendered prompt.
- Otherwise, the raw rendered prompt is returned.

Security (optional):
- Enable with `CMP_AUTH_ENABLED=true` to require Bearer auth and enforce rate limits and RBAC.
