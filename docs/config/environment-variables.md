# Environment Variables Reference

This page lists supported environment variables, their purpose, defaults, and valid values. Variables are optional unless noted. Local-first defaults require no API keys.

## General
- CMP_ENV: Runtime environment. Default: development. Values: development|test|integration|production.
- CMP_PROJECT_ROOT: Project root path. Default: current working directory.
- CMP_LOG_LEVEL: Logging level. Default: info (dev may set debug). Values: debug|info|warn|error.
- CMP_LOG_FORMAT: Log format. Default: json. Values: json|console.

## Local-first provider (Python subprocess)
- CMP_LOCAL_MODELS: Enable local model provider. Default: true for dev flow. Values: true|false.
- CMP_OFFLINE_MODE: Avoid outbound calls. Default: false. Values: true|false.
- CMP_PYTHON_BIN: Python interpreter path for local provider. Default: auto-detected .venv/bin/python or python3.
- CMP_LOCAL_TIMEOUT_SECONDS: Inference subprocess timeout. Default: 600.
- CMP_LOCAL_MODEL_ID: Hugging Face model id (e.g., microsoft/Phi-3-mini-4k-instruct). Tiny models recommended for smoke tests.
- CMP_MODEL_CACHE_DIR: HF model cache directory. Default: ./data/models.
- CMP_PYTHON_SCRIPT: Override path to local_provider.py. Default: auto-discovered.

## Server toggles (runtime/security)
- CMP_AUTH_ENABLED: Enable API key auth and RBAC. Default: false. Values: true|false.
- CMP_PI_ENFORCEMENT: Enable prompt injection detection/sanitization. Default: false. Values: true|false.
- CMP_REQUIRE_CITATION: Require cited sources when memory is used. Default: false. Values: true|false.
- CMP_TENANT_ID: Default tenant id for CLI requests (sent via X-Tenant-ID).

## Memory / Vector database
- CMP_DB_PROVIDER: Structured DB provider. Default: sqlite. Values: sqlite|postgres.
- CMP_DB_PATH: SQLite database path (when sqlite).
- CMP_VECTOR_DB_PROVIDER: Vector store provider. Default: chroma. Values: chroma|pinecone.
- CMP_VECTOR_DB_PATH: Local vector data path (for chroma).
- CMP_CHROMA_PERSIST_DIR: Chroma persistence directory.

## Security / Policies
- CMP_OOB_REQUIRED_ACTIONS: Comma-separated actions requiring out-of-band confirmation (e.g., delete_user,wire_transfer).
- CMP_PII_MODE: PII handling mode. Default: allow. Values: block|redact|allow.
- CMP_EPISODIC_KEY: Encryption key for episodic memory store (if enabled).
- CMP_API_KEYS: Comma-separated apiKeyId:secret pairs for API-key auth.
- CMP_API_TOKENS: Comma-separated tokenId:secret pairs for bearer tokens.

## Hugging Face provider
- HF_TOKEN: HF API token (for remote inference API).
- HF_MODEL_ID: Model id for HF Inference API.
- HF_ENDPOINT: Optional custom HF Inference endpoint.

## OpenAI / Anthropic (production)
- OPENAI_API_KEY: API key for OpenAI providers.
- ANTHROPIC_API_KEY: API key for Anthropic providers.

## Integrations
- PINECONE_API_KEY: Pinecone API key.
- PINECONE_ENVIRONMENT: Pinecone environment/region.
- PINECONE_INDEX: Default Pinecone index name.

## Defaults and precedence
- CLI auto-detects `.venv/bin/python`, sets `CMP_LOCAL_MODELS=true`, and `CMP_PROJECT_ROOT` for `serve`/`run` if unset.
- `config/environments/*.yaml` defines provider defaults; env vars override at runtime where applicable.

See also:
- docs/security.md (security features and policies)
- docs/cli.md (command flags and usage)
