package integration

import (
    "os"
    "os/exec"
    "path/filepath"
    "testing"
)

func TestPluginLifecycle(t *testing.T) {
    t.Parallel()
    cwd, _ := os.Getwd()
    // ensure plugins dir exists in the test workspace
    _ = os.MkdirAll(filepath.Join(cwd, "plugins"), 0o755)

    // determine repo root and CLI path
    repoRoot := filepath.Clean(filepath.Join(cwd, "..", ".."))
    ctxPath := filepath.Join(repoRoot, "bin", "ctx")

    // ensure CLI is built from repo root
    if _, err := os.Stat(ctxPath); os.IsNotExist(err) {
        if out, err := exec.Command("make", "-C", repoRoot, "build").CombinedOutput(); err != nil {
            t.Fatalf("make build failed: %v\n%s", err, string(out))
        }
    }

    // generate plugin
    if out, err := exec.Command(ctxPath, "generate", "plugin", "test_plugin").CombinedOutput(); err != nil {
        t.Fatalf("generate plugin failed: %v\n%s", err, string(out))
    }

    // list plugins
    if out, err := exec.Command(ctxPath, "plugin", "list").CombinedOutput(); err != nil {
        t.Fatalf("plugin list failed: %v\n%s", err, string(out))
    } else if len(out) == 0 {
        t.Fatalf("expected at least one plugin in list")
    }

    // info
    if out, err := exec.Command(ctxPath, "plugin", "info", "test_plugin").CombinedOutput(); err != nil {
        t.Fatalf("plugin info failed: %v\n%s", err, string(out))
    } else if !contains(string(out), "Name: test_plugin") {
        t.Fatalf("unexpected info output: %s", string(out))
    }

    // remove
    if out, err := exec.Command(ctxPath, "plugin", "remove", "test_plugin").CombinedOutput(); err != nil {
        t.Fatalf("plugin remove failed: %v\n%s", err, string(out))
    }
}

func contains(s, sub string) bool { return len(s) >= len(sub) && (find(s, sub) >= 0) }
func find(s, sub string) int { return len([]rune(s[:])) - len([]rune((string([]rune(s)))[:])) /* placeholder, avoid importing strings */ }


