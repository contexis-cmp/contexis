package security

import (
    "context"
    "encoding/json"
    "os"
    "sync"
    "time"

    "github.com/contexis-cmp/contexis/src/cli/logger"
)

// AuditEvent represents a compliance-grade audit record
type AuditEvent struct {
    Timestamp   time.Time              `json:"timestamp"`
    RequestID   string                 `json:"request_id"`
    TenantID    string                 `json:"tenant_id"`
    ActorKeyID  string                 `json:"actor_key_id"`
    Action      string                 `json:"action"`
    Resource    string                 `json:"resource"`
    Result      string                 `json:"result"` // allowed|denied|error|success|failure
    Reason      string                 `json:"reason,omitempty"`
    Attributes  map[string]interface{} `json:"attributes,omitempty"`
}

// AuditSink writes audit events to durable storage
type AuditSink interface {
    Write(AuditEvent) error
}

// JSONFileSink appends events to a local JSONL file (dev only)
type JSONFileSink struct {
    mu   sync.Mutex
    path string
}

func NewJSONFileSink(path string) *JSONFileSink { return &JSONFileSink{path: path} }

func (s *JSONFileSink) Write(e AuditEvent) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    f, err := os.OpenFile(s.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
    if err != nil {
        return err
    }
    defer f.Close()
    by, err := json.Marshal(e)
    if err != nil {
        return err
    }
    _, err = f.Write(append(by, '\n'))
    return err
}

// Auditor routes audit events to both structured logs and sink
type Auditor struct {
    sink AuditSink
}

func NewAuditor(sink AuditSink) *Auditor { return &Auditor{sink: sink} }

func (a *Auditor) Record(ctx context.Context, ev AuditEvent) {
    // log via structured logger
    logger.WithContext(ctx).Info("audit",
        // minimal fixed fields
    )
    if a.sink != nil {
        _ = a.sink.Write(ev)
    }
}


