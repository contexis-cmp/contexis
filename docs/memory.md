# Memory

Providers:
- sqlite: file-backed vector store stored under `memory/<Component>/vector_store.jsonl`
- episodic: append-only conversation log under `memory/<Component>/episodic/`

Configuration merge:
- Optional `memory/<Component>/memory_config.yaml` is merged at runtime (embedding dims, provider hints).

Commands:
```bash
# Ingest documents (one per line)
ctx memory ingest --provider sqlite --component HRBot --model bge-small-en --input policies.txt

# Search
ctx memory search --provider sqlite --component HRBot --query "parental leave" --top-k 5

# Optimize (no-op for small stores but safe to run)
ctx memory optimize --provider sqlite --component HRBot
```

Tenancy:
- Use `--tenant TENANT_ID` on memory commands. Data is written under `memory/<Component>/tenant_<TENANT_ID>/`.
