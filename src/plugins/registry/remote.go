package registry

import (
    "archive/zip"
    "bytes"
    "errors"
    "fmt"
    "io"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
)

// InstallRemote fetches a plugin from a URL or Git repo and installs it.
// Supported sources:
// - HTTPS zip archive URL (must contain plugin.json at root or in a top-level folder)
// - Git repo URL (requires git installed); optional ref via #ref
func (r *Registry) InstallRemote(src string) (Plugin, error) {
    if strings.HasSuffix(strings.ToLower(src), ".zip") || strings.Contains(src, "http") {
        tmp, err := os.MkdirTemp("", "cmp-plugin-zip-*")
        if err != nil { return Plugin{}, err }
        defer os.RemoveAll(tmp)
        if err := downloadZip(src, filepath.Join(tmp, "plugin.zip")); err != nil {
            return Plugin{}, err
        }
        unzipDir := filepath.Join(tmp, "unzipped")
        if err := unzip(filepath.Join(tmp, "plugin.zip"), unzipDir); err != nil {
            return Plugin{}, err
        }
        // Detect plugin root (either unzipDir or first subdir containing plugin.json)
        root := findPluginRoot(unzipDir)
        if root == "" {
            return Plugin{}, errors.New("plugin.json not found in archive")
        }
        return r.Install(root)
    }
    // Git: support URL#ref
    if strings.HasPrefix(src, "git@") || strings.HasPrefix(src, "https://") || strings.HasPrefix(src, "ssh://") {
        url := src
        ref := ""
        if i := strings.Index(src, "#"); i >= 0 {
            url = src[:i]
            ref = src[i+1:]
        }
        tmp, err := os.MkdirTemp("", "cmp-plugin-git-*")
        if err != nil { return Plugin{}, err }
        defer os.RemoveAll(tmp)
        if err := run("git", "clone", "--depth", "1", url, tmp); err != nil {
            return Plugin{}, fmt.Errorf("git clone: %w", err)
        }
        if ref != "" {
            if err := runIn(tmp, "git", "fetch", "origin", ref, "--depth", "1"); err != nil { return Plugin{}, err }
            if err := runIn(tmp, "git", "checkout", ref); err != nil { return Plugin{}, err }
        }
        root := findPluginRoot(tmp)
        if root == "" { return Plugin{}, errors.New("plugin.json not found in repo") }
        return r.Install(root)
    }
    return Plugin{}, fmt.Errorf("unsupported remote source: %s", src)
}

func downloadZip(url, dst string) error {
    resp, err := http.Get(url)
    if err != nil { return err }
    defer resp.Body.Close()
    if resp.StatusCode >= 400 { return fmt.Errorf("download failed: %s", resp.Status) }
    out, err := os.Create(dst)
    if err != nil { return err }
    defer out.Close()
    _, err = io.Copy(out, resp.Body)
    return err
}

func unzip(src, dst string) error {
    r, err := zip.OpenReader(src)
    if err != nil { return err }
    defer r.Close()
    if err := os.MkdirAll(dst, 0o755); err != nil { return err }
    for _, f := range r.File {
        fp := filepath.Join(dst, f.Name)
        if f.FileInfo().IsDir() {
            if err := os.MkdirAll(fp, 0o755); err != nil { return err }
            continue
        }
        if err := os.MkdirAll(filepath.Dir(fp), 0o755); err != nil { return err }
        rc, err := f.Open()
        if err != nil { return err }
        var buf bytes.Buffer
        if _, err := io.Copy(&buf, rc); err != nil { rc.Close(); return err }
        rc.Close()
        if err := os.WriteFile(fp, buf.Bytes(), 0o644); err != nil { return err }
    }
    return nil
}

func findPluginRoot(root string) string {
    // root itself
    if _, err := os.Stat(filepath.Join(root, "plugin.json")); err == nil {
        return root
    }
    entries, err := os.ReadDir(root)
    if err != nil { return "" }
    for _, e := range entries {
        if !e.IsDir() { continue }
        p := filepath.Join(root, e.Name())
        if _, err := os.Stat(filepath.Join(p, "plugin.json")); err == nil {
            return p
        }
    }
    return ""
}

func run(cmd string, args ...string) error {
    c := exec.Command(cmd, args...)
    c.Stdout = nil
    c.Stderr = nil
    return c.Run()
}

func runIn(dir string, cmd string, args ...string) error {
    c := exec.Command(cmd, args...)
    c.Dir = dir
    c.Stdout = nil
    c.Stderr = nil
    return c.Run()
}


