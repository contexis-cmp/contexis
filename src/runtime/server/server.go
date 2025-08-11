package server

import (
    "context"
    "encoding/json"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/contexis-cmp/contexis/src/cli/logger"
    runtimecontext "github.com/contexis-cmp/contexis/src/runtime/context"
    runtimememory "github.com/contexis-cmp/contexis/src/runtime/memory"
    runtimeprompt "github.com/contexis-cmp/contexis/src/runtime/prompt"
    runtimesecurity "github.com/contexis-cmp/contexis/src/runtime/security"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    runtimemodel "github.com/contexis-cmp/contexis/src/runtime/model"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/trace"
    "go.uber.org/zap"
)

type ChatRequest struct {
	TenantID  string                 `json:"tenant_id"`
	Context   string                 `json:"context"`
	Component string                 `json:"component"`
	Query     string                 `json:"query"`
	TopK      int                    `json:"top_k"`
	Data      map[string]interface{} `json:"data"`
    PromptFile string                `json:"prompt_file"`
}

type ChatResponse struct {
	Rendered string `json:"rendered"`
}

// Prometheus metrics
var (
    httpRequestsInFlight = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "cmp_http_requests_in_flight",
        Help: "Number of HTTP requests currently being served.",
    })
    httpRequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
        Name:    "cmp_request_duration_seconds",
        Help:    "Duration of HTTP requests.",
        Buckets: prometheus.DefBuckets,
    }, []string{"method", "path", "code"})
    promptRenderDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
        Name:    "cmp_prompt_render_duration_seconds",
        Help:    "Duration of prompt rendering by component.",
        Buckets: prometheus.DefBuckets,
    }, []string{"component"})
    memorySearchDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
        Name:    "cmp_memory_search_duration_seconds",
        Help:    "Duration of memory search by component.",
        Buckets: prometheus.DefBuckets,
    }, []string{"component"})
    driftScoreGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Name: "cmp_drift_score",
        Help: "Latest drift score by component (optional).",
    }, []string{"component"})
    hfInferenceLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
        Name:    "cmp_hf_inference_latency_seconds",
        Help:    "Latency of Hugging Face inference calls.",
        Buckets: prometheus.DefBuckets,
    }, []string{"model"})
    hfInferenceErrors = prometheus.NewCounterVec(prometheus.CounterOpts{
        Name: "cmp_hf_inference_errors_total",
        Help: "Total errors from Hugging Face inference calls.",
    }, []string{"code"})
)

func init() {
    prometheus.MustRegister(httpRequestsInFlight)
    prometheus.MustRegister(httpRequestDuration)
    prometheus.MustRegister(promptRenderDuration)
    prometheus.MustRegister(memorySearchDuration)
    prometheus.MustRegister(driftScoreGauge)
    prometheus.MustRegister(hfInferenceLatency)
    prometheus.MustRegister(hfInferenceErrors)
}

type statusWriter struct {
    http.ResponseWriter
    status int
}

func (w *statusWriter) WriteHeader(code int) {
    w.status = code
    w.ResponseWriter.WriteHeader(code)
}

// generateRequestID returns a simple timestamp-based ID; in prod consider UUIDs
func generateRequestID() string {
    return time.Now().UTC().Format("20060102T150405.000000000Z07:00")
}

// NewHandler constructs an http.Handler with health, readiness, version, and chat endpoints.
func NewHandler(root string) http.Handler {
    prov, _ := runtimemodel.FromEnv()
    return NewHandlerWithProvider(root, prov)
}

// NewHandlerWithProvider allows injecting a model provider (for tests).
func NewHandlerWithProvider(root string, provider runtimemodel.Provider) http.Handler {
    ctxSvc := runtimecontext.NewContextService(root)
    eng := runtimeprompt.NewEngine(root)
    // Security components (enabled via CMP_AUTH_ENABLED=true)
    authEnabled := os.Getenv("CMP_AUTH_ENABLED") == "true"
    keyStore := runtimesecurity.NewAPIKeyStoreFromEnv()
    rateLimiter := runtimesecurity.NewRateLimiter(10.0/1.0, 5)
    auditor := runtimesecurity.NewAuditor(runtimesecurity.NewJSONFileSink("audit.log"))

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

    // Expose Prometheus metrics
    mux.Handle("/metrics", promhttp.Handler())

    mux.HandleFunc("/api/v1/chat", func(w http.ResponseWriter, r *http.Request) {
        var req ChatRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        // Optional authentication & RBAC
        var principal *runtimesecurity.Principal
        if authEnabled {
            p, err := keyStore.Authenticate(r)
            if err != nil {
                w.Header().Set("WWW-Authenticate", "Bearer")
                http.Error(w, "unauthorized", http.StatusUnauthorized)
                auditor.Record(r.Context(), runtimesecurity.AuditEvent{
                    Timestamp:  time.Now(),
                    RequestID:  r.Context().Value("request_id").(string),
                    TenantID:   req.TenantID,
                    ActorKeyID: "",
                    Action:     "chat:invoke",
                    Resource:   "chat",
                    Result:     "denied",
                    Reason:     "auth_failed",
                })
                return
            }
            principal = p
            // Rate limiting (per key/tenant/ip)
            ip := runtimesecurity.ExtractIP(r)
            if !rateLimiter.Allow(runtimesecurity.LimiterKey{APIKeyID: p.KeyID, TenantID: p.TenantID, IP: ip}, 0) {
                w.Header().Set("Retry-After", runtimesecurity.RetryAfter())
                http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
                auditor.Record(r.Context(), runtimesecurity.AuditEvent{
                    Timestamp:  time.Now(),
                    RequestID:  r.Context().Value("request_id").(string),
                    TenantID:   req.TenantID,
                    ActorKeyID: p.KeyID,
                    Action:     "chat:invoke",
                    Resource:   "chat",
                    Result:     "denied",
                    Reason:     "rate_limited",
                })
                return
            }
            // RBAC: require chat:execute
            res := runtimesecurity.Resource{Type: "chat", Name: "chat", Tenant: req.TenantID}
            if !runtimesecurity.CheckPermission(principal, res, runtimesecurity.ActionExecute) {
                http.Error(w, "forbidden", http.StatusForbidden)
                auditor.Record(r.Context(), runtimesecurity.AuditEvent{
                    Timestamp:  time.Now(),
                    RequestID:  r.Context().Value("request_id").(string),
                    TenantID:   req.TenantID,
                    ActorKeyID: principal.KeyID,
                    Action:     "chat:invoke",
                    Resource:   "chat",
                    Result:     "denied",
                    Reason:     "rbac",
                })
                return
            }
            // Bind principal to context
            r = r.WithContext(runtimesecurity.WithPrincipal(r.Context(), principal))
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
                msStart := time.Now()
                results, _ = store.Search(r.Context(), req.Query, req.TopK)
                memorySearchDuration.WithLabelValues(req.Component).Observe(time.Since(msStart).Seconds())
            }
        }
        data := map[string]interface{}{
            "context": ctxModel,
            "results": results,
        }
        for k, v := range req.Data {
            data[k] = v
        }
        // Safe prompt selection with allowlist
        promptFile := "agent_response.md"
        if req.PromptFile != "" {
            switch req.PromptFile {
            case "agent_response.md", "search_response.md":
                promptFile = req.PromptFile
            default:
                http.Error(w, "unsupported prompt file", http.StatusBadRequest)
                return
            }
        }
        prStart := time.Now()
        rendered, err := eng.RenderFile(req.Component, promptFile, data)
        promptRenderDuration.WithLabelValues(req.Component).Observe(time.Since(prStart).Seconds())
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        // If a provider is configured, perform inference with rendered prompt
        if provider != nil {
            // Tracing span for inference
            tracer := otel.Tracer("contexis/runtime/inference")
            ctx := r.Context()
            ctx, span := tracer.Start(ctx, "huggingface.generate")
            span.SetAttributes(
                attribute.String("provider", "huggingface"),
                attribute.String("model_id", os.Getenv("HF_MODEL_ID")),
            )
            infStart := time.Now()
            out, infErr := provider.Generate(ctx, rendered, runtimemodel.Params{MaxNewTokens: 256})
            hfInferenceLatency.WithLabelValues(os.Getenv("HF_MODEL_ID")).Observe(time.Since(infStart).Seconds())
            if infErr != nil {
                hfInferenceErrors.WithLabelValues("bad_gateway").Inc()
                span.RecordError(infErr)
                span.End()
                http.Error(w, infErr.Error(), http.StatusBadGateway)
                return
            }
            span.End()
            _ = json.NewEncoder(w).Encode(ChatResponse{Rendered: out})
            return
        }
        _ = json.NewEncoder(w).Encode(ChatResponse{Rendered: rendered})
    })

    // Wrap with metrics + tracing + logging context middleware
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
        httpRequestsInFlight.Inc()
        start := time.Now()

        // Correlation and tenant context
        reqID := generateRequestID()
        tenantID := r.Header.Get("X-Tenant-ID")
        ctx := r.Context()
        ctx = context.WithValue(ctx, "request_id", reqID)
        if tenantID != "" {
            ctx = context.WithValue(ctx, "tenant_id", tenantID)
        }

        // Tracing
        tracer := otel.Tracer("contexis/runtime/server")
        var span trace.Span
        ctx, span = tracer.Start(ctx, r.Method+" "+r.URL.Path)
        span.SetAttributes(
            attribute.String("http.method", r.Method),
            attribute.String("http.target", r.URL.Path),
        )
        defer span.End()

        // Serve request with augmented context
        mux.ServeHTTP(sw, r.WithContext(ctx))

        duration := time.Since(start).Seconds()
        httpRequestsInFlight.Dec()
        httpRequestDuration.WithLabelValues(r.Method, r.URL.Path, http.StatusText(sw.status)).Observe(duration)
        logger.WithContext(ctx).Info("request completed",
            zap.String("method", r.Method),
            zap.String("path", r.URL.Path),
            zap.Int("status", sw.status),
            zap.Float64("duration_seconds", duration),
        )
    })
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
        logger.GetLogger().Info("serving", zap.String("addr", addr))
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.GetLogger().Error("server error", zap.Error(err))
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
