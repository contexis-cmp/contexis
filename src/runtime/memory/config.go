package runtimememory

import (
    "fmt"
    "os"
    "path/filepath"

    "gopkg.in/yaml.v3"
)

// LoadComponentMemoryConfig loads memory_config.yaml for a component if present and merges into cfg.Settings.
func LoadComponentMemoryConfig(cfg *Config) error {
    path := filepath.Join(cfg.RootDir, "memory", cfg.ComponentName, "memory_config.yaml")
    by, err := os.ReadFile(path)
    if err != nil {
        // optional
        return nil
    }
    var m map[string]interface{}
    if err := yaml.Unmarshal(by, &m); err != nil {
        return fmt.Errorf("parse memory_config.yaml: %w", err)
    }
    if cfg.Settings == nil { cfg.Settings = map[string]string{} }
    // best-effort flatten of a few known keys
    if vs, ok := m["vector_store"].(map[string]interface{}); ok {
        if t, ok := vs["type"].(string); ok { cfg.Settings["vector_store_type"] = t }
        if p, ok := vs["path"].(string); ok { cfg.Settings["vector_store_path"] = p }
    }
    if em, ok := m["embedding_model"].(map[string]interface{}); ok {
        if n, ok := em["name"].(string); ok { cfg.EmbeddingModel = n }
        if d, ok := em["dimensions"].(int); ok { cfg.Settings["embedding_dim"] = fmt.Sprintf("%d", d) }
    }
    if ep, ok := m["episodic"].(map[string]interface{}); ok {
        if en, ok := ep["enabled"].(bool); ok && en { cfg.Provider = "episodic" }
        if _, ok := ep["encryption"].(bool); ok { cfg.Settings["episodic_encryption"] = "true" }
    }
    return nil
}


