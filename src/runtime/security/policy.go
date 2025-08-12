package security

import (
    "os"
    "strings"
)
// Policy defines action gating requirements for high-risk operations
// Example: map action -> requires_out_of_band_confirmation
type Policy struct {
    RequireOOB map[string]bool
    NoUnsupportedClaims bool // if true, response must cite sources when memory results are used
    PIIMode string // off|redact|block
}

func DefaultPolicy() Policy {
    return Policy{
        RequireOOB: map[string]bool{
            "data_change": true,
            "account_action": true,
        },
        NoUnsupportedClaims: true,
        PIIMode: "off",
    }
}

// RequiresOutOfBand returns true if an action requires external confirmation
func (p Policy) RequiresOutOfBand(action string) bool { return p.RequireOOB[action] }

// MergeEnv allows overriding policy using environment variables
// CMP_OOB_REQUIRED_ACTIONS: comma-separated list of actions requiring OOB (e.g., "data_change,account_action,payment")
// CMP_PII_MODE: one of off|redact|block
func (p Policy) MergeEnv() Policy {
    out := p
    if v := os.Getenv("CMP_OOB_REQUIRED_ACTIONS"); v != "" {
        out.RequireOOB = make(map[string]bool)
        for _, a := range strings.Split(v, ",") {
            if a != "" { out.RequireOOB[a] = true }
        }
    }
    if v := os.Getenv("CMP_PII_MODE"); v != "" {
        switch strings.ToLower(v) {
        case "off", "redact", "block":
            out.PIIMode = strings.ToLower(v)
        }
    }
    return out
}


