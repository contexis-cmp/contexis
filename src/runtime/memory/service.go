package runtimememory

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "path/filepath"
    "strings"
)

// Config for memory service instantiation.
type Config struct {
    Provider        string            // sqlite, episodic
    RootDir         string            // project root for file-backed stores
    ComponentName   string            // e.g., CustomerDocs, SupportBot
    EmbeddingModel  string            // model identifier for vector stores
    Settings        map[string]string // provider-specific settings
    TenantID        string            // tenant isolation
}

// NewStore creates a MemoryStore based on config.
func NewStore(cfg Config) (MemoryStore, error) {
    // Merge component memory_config.yaml if present
    _ = LoadComponentMemoryConfig(&cfg)
    switch strings.ToLower(cfg.Provider) {
    case "sqlite":
        return newSQLiteVectorStore(cfg)
    case "episodic":
        return newEpisodicStore(cfg)
    default:
        return nil, fmt.Errorf("unsupported memory provider: %s", cfg.Provider)
    }
}

// DerivePath returns a tenant-aware path under memory/.
func DerivePath(root, component, tenantID, subpath string) string {
    base := filepath.Join(root, "memory", component)
    if tenantID != "" {
        base = filepath.Join(root, "memory", component, fmt.Sprintf("tenant_%s", sanitize(tenantID)))
    }
    return filepath.Join(base, subpath)
}

func sanitize(s string) string {
    s = strings.ReplaceAll(s, "..", "")
    s = strings.ReplaceAll(s, "/", "_")
    s = strings.ReplaceAll(s, "\\", "_")
    return s
}

// contentSHA computes a deterministic identifier for ingested content.
func contentSHA(chunks []string, model string) string {
    h := sha256.New()
    h.Write([]byte(model))
    for _, c := range chunks {
        h.Write([]byte{0})
        h.Write([]byte(c))
    }
    return hex.EncodeToString(h.Sum(nil))
}


