package server

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    runtimecontext "github.com/contexis-cmp/contexis/src/runtime/context"
    runtimememory "github.com/contexis-cmp/contexis/src/runtime/memory"
    runtimeprompt "github.com/contexis-cmp/contexis/src/runtime/prompt"
)

type ChatRequest struct {
	TenantID  string                 `json:"tenant_id"`
	Context   string                 `json:"context"`
	Component string                 `json:"component"`
	Query     string                 `json:"query"`
	TopK      int                    `json:"top_k"`
	Data      map[string]interface{} `json:"data"`
}

type ChatResponse struct {
	Rendered string `json:"rendered"`
}

// NewHandler constructs an http.Handler with health, readiness, version, and chat endpoints.
func NewHandler(root string) http.Handler {
    ctxSvc := runtimecontext.NewContextService(root)
    eng := runtimeprompt.NewEngine(root)

    mux := http.NewServeMux()

    mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte("ok"))
    })

    mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
        if ctxSvc == nil || eng == nil {
            http.Error(w, "not ready", http.StatusServiceUnavailable)
            return
        }
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte("ready"))
    })

    mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
        _ = json.NewEncoder(w).Encode(map[string]string{
            "version": "0.1.0",
        })
    })

    mux.HandleFunc("/api/v1/chat", func(w http.ResponseWriter, r *http.Request) {
        var req ChatRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        ctxModel, err := ctxSvc.ResolveContext(req.TenantID, req.Context)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        var results []runtimememory.SearchResult
        if req.Component != "" && req.Query != "" {
            store, err := runtimememory.NewStore(runtimememory.Config{Provider: "sqlite", RootDir: root, ComponentName: req.Component, TenantID: req.TenantID})
            if err == nil {
                defer store.Close()
                results, _ = store.Search(r.Context(), req.Query, req.TopK)
            }
        }
        data := map[string]interface{}{
            "context": ctxModel,
            "results": results,
        }
        for k, v := range req.Data {
            data[k] = v
        }
        rendered, err := eng.RenderFile(req.Component, "agent_response.md", data)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        _ = json.NewEncoder(w).Encode(ChatResponse{Rendered: rendered})
    })

    return mux
}

func Serve(addr string) error {
	if addr == "" {
		addr = ":8000"
	}
	root, _ := os.Getwd()
    handler := NewHandler(root)
    srv := &http.Server{Addr: addr, Handler: handler}
    // graceful shutdown
    go func() {
        log.Printf("serving on %s", addr)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Printf("server error: %v", err)
        }
    }()

    // wait for termination signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    return srv.Shutdown(ctx)
}
