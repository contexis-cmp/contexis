package worker

import (
    "net/http"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    jobsProcessed = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "cmp_worker_jobs_processed_total",
        Help: "Total number of jobs processed by the worker.",
    })
)

func init() {
    prometheus.MustRegister(jobsProcessed)
}

// Serve starts a minimal HTTP endpoint for worker health and metrics
func Serve(addr string) error {
    if addr == "" {
        addr = ":9000"
    }
    mux := http.NewServeMux()
    mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte("ok"))
    })
    mux.Handle("/metrics", promhttp.Handler())
    srv := &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 5 * time.Second}
    return srv.ListenAndServe()
}


