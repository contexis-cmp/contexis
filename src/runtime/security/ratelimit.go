package security

import (
    "net"
    "net/http"
    "sync"
    "time"

    "golang.org/x/time/rate"
)

// LimiterKey identifies a limiter bucket
type LimiterKey struct {
    APIKeyID string
    TenantID string
    IP       string
}

// RateLimiter provides token-bucket rate limiting across multiple dimensions
type RateLimiter struct {
    mu       sync.Mutex
    buckets  map[LimiterKey]*rate.Limiter
    defaultR rate.Limit
    defaultB int
}

func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
    return &RateLimiter{
        buckets:  make(map[LimiterKey]*rate.Limiter),
        defaultR: r,
        defaultB: b,
    }
}

func (rl *RateLimiter) Allow(k LimiterKey, perMinuteOverride int) bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    lim, ok := rl.buckets[k]
    if !ok {
        r := rl.defaultR
        b := rl.defaultB
        if perMinuteOverride > 0 {
            r = rate.Limit(float64(perMinuteOverride) / 60.0)
            b = perMinuteOverride / 10
            if b < 1 {
                b = 1
            }
        }
        lim = rate.NewLimiter(r, b)
        rl.buckets[k] = lim
    }
    return lim.Allow()
}

// ExtractIP returns the remote IP from request
func ExtractIP(r *http.Request) string {
    if xf := r.Header.Get("X-Forwarded-For"); xf != "" {
        // take first IP
        for _, part := range splitAndTrim(xf, ',') {
            if net.ParseIP(part) != nil {
                return part
            }
        }
    }
    host, _, err := net.SplitHostPort(r.RemoteAddr)
    if err == nil {
        return host
    }
    return r.RemoteAddr
}

func splitAndTrim(s string, sep rune) []string {
    out := []string{}
    cur := []rune{}
    for _, r := range s {
        if r == sep {
            if len(cur) > 0 {
                out = append(out, stringTrimSpace(string(cur)))
                cur = cur[:0]
            }
        } else {
            cur = append(cur, r)
        }
    }
    if len(cur) > 0 {
        out = append(out, stringTrimSpace(string(cur)))
    }
    return out
}

func stringTrimSpace(s string) string {
    // avoiding importing strings here; we already have a dependency there in auth
    i := 0
    j := len(s)
    for i < j && (s[i] == ' ' || s[i] == '\t' || s[i] == '\n' || s[i] == '\r') {
        i++
    }
    for j > i && (s[j-1] == ' ' || s[j-1] == '\t' || s[j-1] == '\n' || s[j-1] == '\r') {
        j--
    }
    return s[i:j]
}

// Retry-After calculation helper (1 second default)
func RetryAfter() string { return time.Now().Add(1 * time.Second).UTC().Format(http.TimeFormat) }


