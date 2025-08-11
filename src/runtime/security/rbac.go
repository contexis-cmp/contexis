package security

import "strings"

// Resource represents a protected resource
type Resource struct {
    Type   string // e.g., context, memory, prompt, tool, metrics
    Name   string // component or identifier
    Tenant string // tenant scope
}

// Action is an operation on a resource
type Action string

const (
    ActionRead    Action = "read"
    ActionWrite   Action = "write"
    ActionExecute Action = "execute"
    ActionAdmin   Action = "admin"
)

// CheckPermission validates if principal has permission on resource
func CheckPermission(p *Principal, r Resource, a Action) bool {
    if p == nil {
        return false
    }
    // Tenant isolation: must match if key is tenant-bound
    if p.TenantID != "" && r.Tenant != "" && !strings.EqualFold(p.TenantID, r.Tenant) {
        return false
    }
    // Simple scope-based RBAC: scopes like "context:read", "memory:write", "chat:execute", "admin:*"
    needed := r.Type + ":" + string(a)
    if HasScope(p, needed) || HasScope(p, "admin:*") {
        return true
    }
    return false
}


