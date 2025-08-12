package registry

// capabilityRegistry holds known capabilities for validation and discovery.
// These strings are intentionally simple "area:action" identifiers.
var capabilityRegistry = map[string]string{
    // Memory / RAG
    "memory:rerank":         "Provide reranking for search results",
    "memory:vector_backend": "Implement a vector store backend",
    "memory:ingest":         "Custom document ingestion pipeline",

    // Tools
    "tool:mcp":          "Expose an MCP tool integration",
    "tool:api":          "Provide an API-based tool",
    "tool:database":     "Database query tool",

    // Prompting
    "prompt:template":   "Prompt template extensions",

    // Runtime
    "runtime:provider":  "Model provider integration",
}

// IsKnownCapability returns true if cap is in the registry.
func IsKnownCapability(cap string) bool {
    _, ok := capabilityRegistry[cap]
    return ok
}

// ListCapabilities returns the known capability keys.
func ListCapabilities() []string {
    out := make([]string, 0, len(capabilityRegistry))
    for k := range capabilityRegistry {
        out = append(out, k)
    }
    return out
}


