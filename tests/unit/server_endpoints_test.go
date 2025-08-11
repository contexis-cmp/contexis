package unit

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    runtimeserver "github.com/contexis-cmp/contexis/src/runtime/server"
)

func TestHealthz(t *testing.T) {
    h := runtimeserver.NewHandler(t.TempDir())
    req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
    w := httptest.NewRecorder()
    h.ServeHTTP(w, req)
    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d", w.Code)
    }
}

func TestReadyz(t *testing.T) {
    h := runtimeserver.NewHandler(t.TempDir())
    req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
    w := httptest.NewRecorder()
    h.ServeHTTP(w, req)
    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d", w.Code)
    }
}

func TestVersion(t *testing.T) {
    h := runtimeserver.NewHandler(t.TempDir())
    req := httptest.NewRequest(http.MethodGet, "/version", nil)
    w := httptest.NewRecorder()
    h.ServeHTTP(w, req)
    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d", w.Code)
    }
    var got map[string]string
    if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
        t.Fatalf("invalid json: %v", err)
    }
    if got["version"] == "" {
        t.Fatalf("missing version field")
    }
}


