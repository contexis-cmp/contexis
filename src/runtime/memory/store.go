package runtimememory

import "context"

// SearchResult represents a single search match from a memory store.
type SearchResult struct {
    ID       string                 `json:"id"`
    Content  string                 `json:"content"`
    Score    float64                `json:"score"`
    Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// MemoryStore defines the interface for memory providers.
type MemoryStore interface {
    // IngestDocuments ingests raw text documents and returns a version identifier (e.g., memory SHA).
    IngestDocuments(ctx context.Context, documents []string) (string, error)

    // Search performs a semantic search and returns the top results.
    Search(ctx context.Context, query string, topK int) ([]SearchResult, error)

    // Optimize performs any maintenance/compaction tasks.
    Optimize(ctx context.Context, memoryVersion string) error

    // Close releases resources.
    Close() error
}


