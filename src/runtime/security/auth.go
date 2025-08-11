package security

import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "errors"
    "net/http"
    "os"
    "strings"
    "sync"
)

// Principal represents an authenticated caller
type Principal struct {
    KeyID    string
    TenantID string
    Scopes   []string
}

// APIKey represents an API key entry in the keystore
type APIKey struct {
    KeyID     string   `json:"key_id"`
    Hash      string   `json:"hash"`       // hex-encoded SHA-256 of the token
    TenantID  string   `json:"tenant_id"`  // associated tenant
    Scopes    []string `json:"scopes"`     // permissions like "chat:invoke", "context:read"
    RateLimit int      `json:"rate_limit"` // requests per minute (optional)
}

// APIKeyStore provides lookup for API keys
type APIKeyStore struct {
    mu   sync.RWMutex
    byID map[string]APIKey
    byH  map[string]APIKey
}

// NewAPIKeyStoreFromEnv loads keys from env. Supports two formats:
// 1) CMP_API_KEYS as JSON array of APIKey
// 2) CMP_API_TOKENS as comma-separated list of "token@tenant:scope1|scope2" (for local dev)
func NewAPIKeyStoreFromEnv() *APIKeyStore {
    store := &APIKeyStore{byID: make(map[string]APIKey), byH: make(map[string]APIKey)}

    if raw := os.Getenv("CMP_API_KEYS"); strings.TrimSpace(raw) != "" {
        var keys []APIKey
        if err := json.Unmarshal([]byte(raw), &keys); err == nil {
            for _, k := range keys {
                store.add(k)
            }
        }
    }

    if raw := os.Getenv("CMP_API_TOKENS"); strings.TrimSpace(raw) != "" {
        // token@tenant:scope1|scope2
        entries := strings.Split(raw, ",")
        for i, e := range entries {
            e = strings.TrimSpace(e)
            if e == "" {
                continue
            }
            var token, tenant, scopesStr string
            parts := strings.SplitN(e, "@", 2)
            if len(parts) == 2 {
                token = parts[0]
                rest := parts[1]
                parts2 := strings.SplitN(rest, ":", 2)
                tenant = parts2[0]
                if len(parts2) == 2 {
                    scopesStr = parts2[1]
                }
            } else {
                // token only
                token = e
            }
            hash := sha256Sum(token)
            key := APIKey{
                KeyID:     "env-" + itoa(i+1),
                Hash:      hash,
                TenantID:  tenant,
                Scopes:    splitScopes(scopesStr),
                RateLimit: 0,
            }
            store.add(key)
        }
    }

    return store
}

func (s *APIKeyStore) add(k APIKey) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.byID[k.KeyID] = k
    if k.Hash != "" {
        s.byH[strings.ToLower(k.Hash)] = k
    }
}

func sha256Sum(s string) string {
    h := sha256.Sum256([]byte(s))
    return strings.ToLower(hex.EncodeToString(h[:]))
}

func splitScopes(s string) []string {
    if s == "" {
        return nil
    }
    parts := strings.Split(s, "|")
    out := make([]string, 0, len(parts))
    for _, p := range parts {
        p = strings.TrimSpace(p)
        if p != "" {
            out = append(out, p)
        }
    }
    return out
}

// Authenticate extracts the Bearer token and verifies it against the store
func (s *APIKeyStore) Authenticate(r *http.Request) (*Principal, error) {
    auth := r.Header.Get("Authorization")
    if auth == "" {
        return nil, errors.New("missing Authorization header")
    }
    if !strings.HasPrefix(strings.ToLower(auth), "bearer ") {
        return nil, errors.New("unsupported auth scheme")
    }
    token := strings.TrimSpace(auth[len("Bearer "):])
    if token == "" {
        return nil, errors.New("empty bearer token")
    }
    h := sha256Sum(token)
    s.mu.RLock()
    key, ok := s.byH[h]
    s.mu.RUnlock()
    if !ok {
        return nil, errors.New("invalid token")
    }
    return &Principal{KeyID: key.KeyID, TenantID: key.TenantID, Scopes: key.Scopes}, nil
}

type principalKey struct{}

// WithPrincipal stores the principal in context
func WithPrincipal(ctx context.Context, p *Principal) context.Context {
    return context.WithValue(ctx, principalKey{}, p)
}

// FromPrincipal retrieves the principal from context
func FromPrincipal(ctx context.Context) (*Principal, bool) {
    p, ok := ctx.Value(principalKey{}).(*Principal)
    return p, ok
}

// HasScope returns true if the principal has the given scope
func HasScope(p *Principal, scope string) bool {
    for _, s := range p.Scopes {
        if s == scope || s == "*" {
            return true
        }
    }
    return false
}

// itoa is a tiny helper to avoid importing strconv in this file
func itoa(i int) string {
    // minimal conversion for small ints
    digits := "0123456789"
    if i == 0 {
        return "0"
    }
    neg := false
    if i < 0 {
        neg = true
        i = -i
    }
    var b []byte
    for i > 0 {
        d := i % 10
        b = append([]byte{digits[d]}, b...)
        i /= 10
    }
    if neg {
        b = append([]byte{'-'}, b...)
    }
    return string(b)
}


