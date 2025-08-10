package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

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

func Serve(addr string) error {
	if addr == "" {
		addr = ":8000"
	}
	root, _ := os.Getwd()
	ctxSvc := runtimecontext.NewContextService(root)
	eng := runtimeprompt.NewEngine(root)

	http.HandleFunc("/api/v1/chat", func(w http.ResponseWriter, r *http.Request) {
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
		// Memory search (optional)
		var results []runtimememory.SearchResult
		if req.Component != "" && req.Query != "" {
			store, err := runtimememory.NewStore(runtimememory.Config{Provider: "sqlite", RootDir: root, ComponentName: req.Component, TenantID: req.TenantID})
			if err == nil {
				defer store.Close()
				results, _ = store.Search(r.Context(), req.Query, req.TopK)
			}
		}
		// Render prompt (example path: agent_response.md or search_response.md)
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
	log.Printf("serving on %s", addr)
	return http.ListenAndServe(addr, nil)
}
