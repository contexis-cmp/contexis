# HR Chatbot Example

This example shows a minimal HR chatbot that uses:
- Hugging Face model `openai/gpt-oss-20b`
- `sqlite` memory provider for company policies

Steps:
1) Generate agent
```bash
ctx generate agent HRBot --tools database --memory episodic
```
2) Add policies to `examples/hr-chatbot/policies.txt` (one policy per line).
3) Ingest to memory
```bash
ctx memory ingest --provider sqlite --component HRBot --model bge-small-en --input examples/hr-chatbot/policies.txt
```
4) Start server
```bash
export HF_TOKEN=...; export HF_MODEL_ID=openai/gpt-oss-20b
ctx serve --addr :8000
```
5) Query via API
```bash
curl -s localhost:8000/api/v1/chat \
  -H 'Content-Type: application/json' \
  -d '{"tenant_id":"","context":"HRBot","component":"HRBot","query":"What is our parental leave policy?","top_k":5,"data":{},"prompt_file":"agent_response.md"}'
```
