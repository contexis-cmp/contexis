package registry

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Manifest describes a plugin package.
//
// The manifest is stored at `plugins/<name>/plugin.json` and controls
// discoverability and capability declaration.
// Example:
//
//	{
//	  "name": "re_ranker",
//	  "version": "0.1.14",
//	  "description": "Reranking component for RAG",
//	  "capabilities": ["memory:rerank"],
//	  "compatibility": {"cmp_min": "0.1.14"}
//	}
type Manifest struct {
	Name          string            `json:"name"`
	Version       string            `json:"version"`
	Description   string            `json:"description"`
	Capabilities  []string          `json:"capabilities"`
	Compatibility map[string]string `json:"compatibility"`
	Signature     string            `json:"signature,omitempty"` // optional detached signature (hex/base64)
}

// Plugin represents a discovered plugin on disk.
type Plugin struct {
	Path     string
	Manifest Manifest
}

// Registry loads and manages plugins rooted at a project directory.
type Registry struct {
	root string
}

// NewRegistry creates a new Registry reading from `<root>/plugins`.
func NewRegistry(root string) *Registry { return &Registry{root: root} }

func (r *Registry) pluginsDir() string { return filepath.Join(r.root, "plugins") }

// List scans the plugins directory and returns parsed plugin manifests.
func (r *Registry) List() ([]Plugin, error) {
	dir := r.pluginsDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []Plugin{}, nil
		}
		return nil, err
	}
	var out []Plugin
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		p := filepath.Join(dir, e.Name())
		mf := filepath.Join(p, "plugin.json")
		by, err := os.ReadFile(mf)
		if err != nil {
			continue
		}
		var m Manifest
		if err := json.Unmarshal(by, &m); err != nil {
			continue
		}
		if m.Name == "" || m.Version == "" {
			continue
		}
		out = append(out, Plugin{Path: p, Manifest: m})
	}
	sort.Slice(out, func(i, j int) bool { return strings.Compare(out[i].Manifest.Name, out[j].Manifest.Name) < 0 })
	return out, nil
}

// Info returns a single plugin manifest by name.
func (r *Registry) Info(name string) (*Plugin, error) {
	list, err := r.List()
	if err != nil {
		return nil, err
	}
	for _, p := range list {
		if p.Manifest.Name == name {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("plugin '%s' not found", name)
}

// Install copies a plugin directory from src into `plugins/<name>`.
func (r *Registry) Install(src string) (Plugin, error) {
	// read manifest from src
	mf := filepath.Join(src, "plugin.json")
	by, err := os.ReadFile(mf)
	if err != nil {
		return Plugin{}, fmt.Errorf("read manifest: %w", err)
	}
	var m Manifest
	if err := json.Unmarshal(by, &m); err != nil {
		return Plugin{}, fmt.Errorf("parse manifest: %w", err)
	}
	if m.Name == "" {
		return Plugin{}, fmt.Errorf("manifest name is required")
	}
	// Capability checks (best-effort)
	for _, c := range m.Capabilities {
		if !IsKnownCapability(c) {
			// allow unknown but warn by returning an error; callers may decide policy later
			// choose permissive: do not error, just continue
			_ = c
		}
	}
	dst := filepath.Join(r.pluginsDir(), m.Name)
	if err := copyDir(src, dst); err != nil {
		return Plugin{}, err
	}
	return Plugin{Path: dst, Manifest: m}, nil
}

// Remove deletes `plugins/<name>` recursively.
func (r *Registry) Remove(name string) error {
	path := filepath.Join(r.pluginsDir(), name)
	if err := os.RemoveAll(path); err != nil {
		return err
	}
	return nil
}

func copyDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		if rel == "." {
			return nil
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		by, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, by, 0o644)
	})
}
