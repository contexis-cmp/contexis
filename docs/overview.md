# Contexis Overview

Contexis is a framework for building reproducible AI applications using the CMP architecture:

- Context: Declarative agent configuration in `.ctx` YAML (persona, tools, guardrails, memory, testing)
- Memory: Versioned knowledge bases (vector store) and episodic logs, tenant-aware
- Prompt: Pure templates rendered at runtime with data and context

Key properties:
- Versioned and validated contexts (`src/core/schema/context_schema.json`)
- Memory providers: `sqlite` vector store (file-backed JSONL) and `episodic` conversation logs
- Prompt engine with include functions and helper funcs (`src/runtime/prompt/engine.go`)
- HTTP runtime server for chat (`ctx serve`), optional Hugging Face inference
- Security (optional): API key auth, rate limiting, audit trail

Typical project layout:
```
project/
├── contexts/<Component>/*.ctx
├── memory/<Component>/
├── prompts/<Component>/*.md
├── tools/<Component>/*.py
└── tests/<Component>/
```

Start with `QUICKSTART.md` for a guided build of an HR chatbot.
