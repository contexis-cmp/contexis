# Model Providers

Hugging Face Inference API is supported.

Environment variables:
- `HF_TOKEN`: your Hugging Face access token
- `HF_MODEL_ID`: model identifier, e.g. `openai/gpt-oss-20b`
- `HF_ENDPOINT` (optional): base endpoint, defaults to `https://api-inference.huggingface.co/models`

Server usage:
- When `HF_TOKEN` and `HF_MODEL_ID` are present, `ctx serve` will call HF for completions.

Example:
```bash
export HF_TOKEN=hf_********************************
export HF_MODEL_ID=openai/gpt-oss-20b
ctx serve --addr :8000
```
