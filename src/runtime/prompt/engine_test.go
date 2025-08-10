package runtimeprompt

import (
    "os"
    "path/filepath"
    "testing"
)

func TestRenderFile(t *testing.T) {
    root := t.TempDir()
    // create a sample template under prompts/Comp
    dir := filepath.Join(root, "prompts", "Comp")
    if err := os.MkdirAll(dir, 0o755); err != nil { t.Fatal(err) }
    tpl := "Hello {{.name}}!"
    if err := os.WriteFile(filepath.Join(dir, "greet.md"), []byte(tpl), 0o644); err != nil { t.Fatal(err) }

    eng := NewEngine(root)
    out, err := eng.RenderFile("Comp", "greet.md", map[string]interface{}{"name": "World"})
    if err != nil { t.Fatalf("RenderFile: %v", err) }
    if out != "Hello World!" {
        t.Fatalf("unexpected render output: %q", out)
    }
}

func TestOptimizeTokens(t *testing.T) {
    content := "one two three four five"
    out := OptimizeTokens(content, 3)
    if out == content { t.Fatalf("expected trimming") }
}

func TestValidateFormat(t *testing.T) {
    if err := ValidateFormat("json", `{"a":1}`); err != nil {
        t.Fatalf("json validate: %v", err)
    }
    if err := ValidateFormat("markdown", "# Title\nBody"); err != nil {
        t.Fatalf("md validate: %v", err)
    }
}


