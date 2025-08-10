package runtimecontext

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, base string, rel string, content string) string {
	t.Helper()
	path := filepath.Join(base, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdirs failed: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}
	return path
}

func TestResolveContext_Global(t *testing.T) {
	root := t.TempDir()
	// Minimal valid context
	yaml := `name: "Foo"
version: "1.0.0"
role:
  persona: "Test Persona"
`
	writeFile(t, root, filepath.Join("contexts", "Foo", "Foo.ctx"), yaml)

	svc := NewContextService(root)
	ctx, err := svc.ResolveContext("", "Foo")
	if err != nil {
		t.Fatalf("ResolveContext failed: %v", err)
	}
	if ctx == nil || ctx.Name != "Foo" {
		t.Fatalf("unexpected context: %#v", ctx)
	}
}

func TestResolveContext_TenantOverrides(t *testing.T) {
	root := t.TempDir()
	// Global
	global := `name: "Foo"
version: "1.0.0"
role:
  persona: "Global Persona"
`
	writeFile(t, root, filepath.Join("contexts", "Foo", "Foo.ctx"), global)
	// Tenant override
	tenant := `name: "Foo"
version: "1.0.0"
role:
  persona: "Tenant Persona"
`
	writeFile(t, root, filepath.Join("contexts", "tenants", "acme", "Foo.ctx"), tenant)

	svc := NewContextService(root)
	ctx, err := svc.ResolveContext("acme", "Foo")
	if err != nil {
		t.Fatalf("ResolveContext failed: %v", err)
	}
	if got := ctx.Role.Persona; got != "Tenant Persona" {
		t.Fatalf("expected tenant persona, got %q", got)
	}
}

func TestResolveContext_ExtendsIncludeMerge(t *testing.T) {
	root := t.TempDir()
	// Base file with capabilities
	base := `name: "Base"
version: "1.0.0"
role:
  persona: "Base Persona"
  capabilities: ["a", "b"]
tools:
  - name: base_tool
    uri: mcp://base
`
	writeFile(t, root, filepath.Join("contexts", "Foo", "base.yaml"), base)
	// Fragment file to include tools
	frag := `tools:
  - name: include_tool
    uri: mcp://inc
`
	writeFile(t, root, filepath.Join("contexts", "Foo", "frag.yaml"), frag)
	// Main Foo ctx extends base and includes frag, overrides persona and adds capability
	main := `extends: base.yaml
include: [frag.yaml]
name: "Foo"
version: "1.0.0"
role:
  persona: "Main Persona"
  capabilities: ["b", "c"]
`
	writeFile(t, root, filepath.Join("contexts", "Foo", "Foo.ctx"), main)

	svc := NewContextService(root)
	ctx, err := svc.ResolveContext("", "Foo")
	if err != nil {
		t.Fatalf("ResolveContext failed: %v", err)
	}
	if ctx.Role.Persona != "Main Persona" {
		t.Fatalf("expected merged persona, got %q", ctx.Role.Persona)
	}
	// capabilities union: [a, b, c] (order not guaranteed; check contents)
	caps := map[string]bool{}
	for _, c := range ctx.Role.Capabilities {
		caps[c] = true
	}
	if !(caps["a"] && caps["b"] && caps["c"]) {
		t.Fatalf("expected unioned capabilities, got %v", ctx.Role.Capabilities)
	}
	// tools union includes include_tool
	foundInclude := false
	for _, t2 := range ctx.Tools {
		if t2.Name == "include_tool" {
			foundInclude = true
			break
		}
	}
	if !foundInclude {
		t.Fatalf("expected include_tool in merged tools, got %v", ctx.Tools)
	}
}

func TestReloadContext_ClearsCache(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join("contexts", "Foo", "Foo.ctx")
	contentV1 := `name: "Foo"
version: "1.0.0"
role:
  persona: "P1"
`
	writeFile(t, root, path, contentV1)

	svc := NewContextService(root)
	ctx1, err := svc.ResolveContext("", "Foo")
	if err != nil {
		t.Fatalf("ResolveContext failed: %v", err)
	}
	if ctx1.Role.Persona != "P1" {
		t.Fatalf("unexpected persona: %q", ctx1.Role.Persona)
	}

	// Update file
	contentV2 := `name: "Foo"
version: "1.0.1"
role:
  persona: "P2"
`
	writeFile(t, root, path, contentV2)

	// Without reload, cache would return old; after reload, should pick new
	if err := svc.ReloadContext(""); err != nil {
		t.Fatalf("ReloadContext failed: %v", err)
	}
	ctx2, err := svc.ResolveContext("", "Foo")
	if err != nil {
		t.Fatalf("ResolveContext failed: %v", err)
	}
	if ctx2.Role.Persona != "P2" {
		t.Fatalf("expected updated persona after reload, got %q", ctx2.Role.Persona)
	}
}

func TestResolveContext_InvalidYAML(t *testing.T) {
	root := t.TempDir()
	invalid := `name: 123
version: "1.0.0"
role:
  persona: "ok"
`
	writeFile(t, root, filepath.Join("contexts", "Foo", "Foo.ctx"), invalid)

	svc := NewContextService(root)
	if _, err := svc.ResolveContext("", "Foo"); err == nil {
		t.Fatalf("expected validation error for invalid YAML, got nil")
	}
}
