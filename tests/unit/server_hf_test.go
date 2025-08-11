package unit

import (
    "bytes"
    "context"
    "encoding/json"
    "os"
    "path/filepath"
    "net/http"
    "net/http/httptest"
    "testing"

    runtimemodel "github.com/contexis-cmp/contexis/src/runtime/model"
    runtimeserver "github.com/contexis-cmp/contexis/src/runtime/server"
)

type fakeProvider struct{ out string; err error }

func (f fakeProvider) Generate(_ context.Context, _ string, _ runtimemodel.Params) (string, error) {
    return f.out, f.err
}

func scaffoldTempRoot(t *testing.T) string {
    t.Helper()
    root := t.TempDir()
    // minimal context file: contexts/SupportBot/support_bot.ctx
    ctxDir := filepath.Join(root, "contexts", "SupportBot")
    if err := os.MkdirAll(ctxDir, 0o755); err != nil { t.Fatal(err) }
    ctxYAML := []byte("name: SupportBot\nversion: '1.0.0'\nrole:\n  persona: 'helper'\n")
    if err := os.WriteFile(filepath.Join(ctxDir, "support_bot.ctx"), ctxYAML, 0o644); err != nil { t.Fatal(err) }
    // minimal prompt template: prompts/SupportBot/agent_response.md
    prDir := filepath.Join(root, "prompts", "SupportBot")
    if err := os.MkdirAll(prDir, 0o755); err != nil { t.Fatal(err) }
    if err := os.WriteFile(filepath.Join(prDir, "agent_response.md"), []byte("TEMPLATE"), 0o644); err != nil { t.Fatal(err) }
    return root
}

func TestChatWithoutProvider_RendersTemplate(t *testing.T) {
    root := scaffoldTempRoot(t)
    h := runtimeserver.NewHandlerWithProvider(root, nil)
    reqBody := runtimeserver.ChatRequest{TenantID: "", Context: "SupportBot", Component: "SupportBot", Query: ""}
    by, _ := json.Marshal(reqBody)
    req := httptest.NewRequest(http.MethodPost, "/api/v1/chat", bytes.NewReader(by))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    h.ServeHTTP(w, req)
    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
    }
    var got runtimeserver.ChatResponse
    if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil { t.Fatal(err) }
    if got.Rendered != "TEMPLATE" {
        t.Fatalf("expected TEMPLATE, got %q", got.Rendered)
    }
}

func TestChatWithProvider_UsesInferenceOutput(t *testing.T) {
    root := scaffoldTempRoot(t)
    prov := fakeProvider{out: "HF OUT"}
    h := runtimeserver.NewHandlerWithProvider(root, prov)
    reqBody := runtimeserver.ChatRequest{TenantID: "", Context: "SupportBot", Component: "SupportBot", Query: ""}
    by, _ := json.Marshal(reqBody)
    req := httptest.NewRequest(http.MethodPost, "/api/v1/chat", bytes.NewReader(by))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    h.ServeHTTP(w, req)
    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
    }
    var got runtimeserver.ChatResponse
    if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil { t.Fatal(err) }
    if got.Rendered != "HF OUT" {
        t.Fatalf("expected HF OUT, got %q", got.Rendered)
    }
}


