package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

// GeneratePlugin scaffolds a plugin from the template
func GeneratePlugin(_ context.Context, name string) error {
	if name == "" {
		return fmt.Errorf("plugin name is required")
	}
	cwd, _ := os.Getwd()
	dst := filepath.Join(cwd, "plugins", name)
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}
	// Write minimal manifest
	manifest := fmt.Sprintf(`{
  "name": "%s",
  "version": "0.2.0",
  "description": "Custom CMP plugin",
  "capabilities": [],
  "compatibility": {"cmp_min": "0.2.0"}
}
`, name)
	if err := os.WriteFile(filepath.Join(dst, "plugin.json"), []byte(manifest), 0o644); err != nil {
		return err
	}
	readme := []byte("# " + name + "\n\nThis is a CMP plugin.\n")
	if err := os.WriteFile(filepath.Join(dst, "README.md"), readme, 0o644); err != nil {
		return err
	}
	fmt.Printf("Plugin scaffold created at %s\n", dst)
	return nil
}
