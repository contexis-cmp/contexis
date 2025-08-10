package runtimecontext

import (
    "errors"
    "fmt"
    "io/fs"
    "os"
    "path/filepath"
    "strings"
    "sync"

    corectx "github.com/contexis-cmp/contexis/src/core/context"
    coreval "github.com/contexis-cmp/contexis/src/core/schema"
    "gopkg.in/yaml.v3"
)

// ContextService is responsible for resolving, validating, and caching contexts.
type ContextService struct {
    mu          sync.RWMutex
    cache       map[string]*corectx.Context // key: tenantID|contextName
    projectRoot string
}

// NewContextService creates a new ContextService instance.
func NewContextService(projectRoot string) *ContextService {
    return &ContextService{
        cache:       make(map[string]*corectx.Context),
        projectRoot: projectRoot,
    }
}

// ResolveContext loads, validates, and returns a Context for a tenant and context name.
// Resolution order:
// 1) contexts/tenants/<tenantID>/<contextName>.ctx
// 2) contexts/<contextName>/*.ctx (first .ctx file found in that directory)
func (s *ContextService) ResolveContext(tenantID, contextName string) (*corectx.Context, error) {
    if contextName == "" {
        return nil, fmt.Errorf("context name is required")
    }

    key := fmt.Sprintf("%s|%s", tenantID, contextName)

    s.mu.RLock()
    if ctx, ok := s.cache[key]; ok {
        s.mu.RUnlock()
        return ctx, nil
    }
    s.mu.RUnlock()

    // Resolve path
    candidatePaths := s.candidatePaths(tenantID, contextName)
    var loaded *corectx.Context
    var loadErr error
    for _, p := range candidatePaths {
        data, err := os.ReadFile(p)
        if err != nil {
            if errors.Is(err, fs.ErrNotExist) {
                continue
            }
            loadErr = fmt.Errorf("read context '%s': %w", p, err)
            continue
        }

        // Validate against schema (lightweight)
        if err := coreval.ValidateContextYAML(data); err != nil {
            loadErr = fmt.Errorf("schema validation failed for '%s': %w", p, err)
            continue
        }

        // Process extends/include and decode into model
        mergedMap, err := s.loadAndMergeYAML(p, 0)
        if err != nil {
            loadErr = fmt.Errorf("merge failed for '%s': %w", p, err)
            continue
        }

        jsonBytes, err := yamlToJSON(mergedMap)
        if err != nil {
            loadErr = fmt.Errorf("yaml->json failed for '%s': %w", p, err)
            continue
        }

        ctxModel, err := corectx.FromJSON(jsonBytes)
        if err != nil {
            loadErr = fmt.Errorf("parse context model for '%s': %w", p, err)
            continue
        }

        if err := ctxModel.Validate(); err != nil {
            loadErr = fmt.Errorf("context validation for '%s': %w", p, err)
            continue
        }

        loaded = ctxModel
        break
    }

    if loaded == nil {
        if loadErr != nil {
            return nil, loadErr
        }
        return nil, fmt.Errorf("context '%s' not found", contextName)
    }

    s.mu.Lock()
    s.cache[key] = loaded
    s.mu.Unlock()

    return loaded, nil
}

// ReloadContext clears the cache so subsequent calls re-read from disk.
// If a path is supplied, this is currently a no-op beyond clearing the cache.
func (s *ContextService) ReloadContext(_ string) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.cache = make(map[string]*corectx.Context)
    return nil
}

// candidatePaths computes possible file paths for a given tenant and context.
func (s *ContextService) candidatePaths(tenantID, contextName string) []string {
    var paths []string
    // Tenant-specific direct file path
    if tenantID != "" {
        tenantPath := filepath.Join(s.projectRoot, "contexts", "tenants", sanitizePath(tenantID), fmt.Sprintf("%s.ctx", contextName))
        paths = append(paths, tenantPath)
    }

    // Global directory: find any .ctx file under contexts/<contextName>/
    globalDir := filepath.Join(s.projectRoot, "contexts", contextName)
    // Common filenames to try first
    preferred := []string{
        filepath.Join(globalDir, strings.ToLower(contextName)+".ctx"),
        filepath.Join(globalDir, "rag_agent.ctx"),
        filepath.Join(globalDir, "workflow_coordinator.ctx"),
    }
    paths = append(paths, preferred...)

    // Fallback: first .ctx file in directory
    _ = filepath.WalkDir(globalDir, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return nil
        }
        if d.IsDir() {
            return nil
        }
        if strings.HasSuffix(strings.ToLower(d.Name()), ".ctx") {
            paths = append(paths, path)
            // don't stop; collect all to increase chances
        }
        return nil
    })
    return uniqueStrings(paths)
}

// loadAndMergeYAML loads a YAML file, processes extends/include recursively, and returns a merged map.
func (s *ContextService) loadAndMergeYAML(path string, depth int) (map[string]interface{}, error) {
    if depth > 5 {
        return nil, fmt.Errorf("maximum extends/include depth exceeded at %s", path)
    }
    raw, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    var current map[string]interface{}
    if err := yaml.Unmarshal(raw, &current); err != nil {
        return nil, err
    }

    // Read meta: extends and include
    extendsVal, _ := current["extends"].(string)
    includeVal, _ := current["include"].([]interface{})

    var base map[string]interface{}
    if extendsVal != "" {
        basePath := s.resolveRelative(path, extendsVal)
        var err error
        base, err = s.loadAndMergeYAML(basePath, depth+1)
        if err != nil {
            return nil, fmt.Errorf("failed to load base '%s': %w", basePath, err)
        }
    }

    // Apply includes
    merged := make(map[string]interface{})
    if base != nil {
        merged = deepCopyMap(base)
    }
    if len(includeVal) > 0 {
        for _, inc := range includeVal {
            incStr, _ := inc.(string)
            if incStr == "" {
                continue
            }
            incPath := s.resolveRelative(path, incStr)
            frag, err := s.loadAndMergeYAML(incPath, depth+1)
            if err != nil {
                return nil, fmt.Errorf("failed to load include '%s': %w", incPath, err)
            }
            merged = DeepMerge(merged, frag)
        }
    }

    // Merge current over the accumulated base/includes
    merged = DeepMerge(merged, current)
    // Remove meta keys from final
    delete(merged, "extends")
    delete(merged, "include")
    return merged, nil
}

func (s *ContextService) resolveRelative(baseFile string, relative string) string {
    // Support absolute-like aliases in future; for now treat as relative to file dir
    if filepath.IsAbs(relative) {
        return relative
    }
    baseDir := filepath.Dir(baseFile)
    return filepath.Clean(filepath.Join(baseDir, relative))
}

func sanitizePath(p string) string {
    // Prevent directory traversal; allow simple ids
    s := strings.ReplaceAll(p, "..", "")
    s = strings.ReplaceAll(s, string(filepath.Separator), "_")
    return s
}

func yamlToJSON(m map[string]interface{}) ([]byte, error) {
    // Marshal via yaml then convert to json by re-unmarshal; simplest portable approach
    // We rely on core/context.FromJSON to parse into struct
    by, err := yaml.Marshal(m)
    if err != nil {
        return nil, err
    }
    // Use yaml.Node to JSON: easiest is to unmarshal into interface{} then marshal to JSON via stdlib
    var tmp interface{}
    if err := yaml.Unmarshal(by, &tmp); err != nil {
        return nil, err
    }
    // Convert YAML numbers to JSON-compatible by walking; for now, rely on stdlib json marshaler in FromJSON
    // Return the YAML bytes; FromJSON expects JSON though, so we should marshal to JSON here
    // To avoid importing another lib, we encode via json.Marshal from the decoded map
    // However, core/context.FromJSON expects JSON bytes; do the conversion explicitly here
    // We'll use encoding/json indirectly by calling core/context.FromJSON only after json marshal
    // Implemented in place in this function for clarity
    // Re-unmarshal into map[string]interface{} already done; now marshal to JSON
    return marshalJSON(tmp)
}

func marshalJSON(v interface{}) ([]byte, error) {
    // Local wrapper to avoid importing encoding/json here; defer to core package which already imports json
    // However, we cannot access non-exported json from here. We'll just import encoding/json directly.
    return coreval.JSONMarshal(v)
}

func uniqueStrings(in []string) []string {
    seen := make(map[string]struct{}, len(in))
    out := make([]string, 0, len(in))
    for _, s := range in {
        if s == "" {
            continue
        }
        if _, ok := seen[s]; ok {
            continue
        }
        seen[s] = struct{}{}
        out = append(out, s)
    }
    return out
}

// DeepMerge merges src into dst recursively. Lists are unioned by value; maps are deep-merged; scalars override.
func DeepMerge(dst, src map[string]interface{}) map[string]interface{} {
    if dst == nil && src == nil {
        return map[string]interface{}{}
    }
    if dst == nil {
        return deepCopyMap(src)
    }
    if src == nil {
        return deepCopyMap(dst)
    }
    out := deepCopyMap(dst)
    for k, v := range src {
        if v == nil {
            out[k] = nil
            continue
        }
        if existing, ok := out[k]; ok {
            switch vTyped := v.(type) {
            case map[string]interface{}:
                if exMap, ok := existing.(map[string]interface{}); ok {
                    out[k] = DeepMerge(exMap, vTyped)
                } else {
                    out[k] = deepCopyMap(vTyped)
                }
            case []interface{}:
                // Union lists by string value fallback
                out[k] = unionLists(existing, vTyped)
            default:
                out[k] = v
            }
        } else {
            // New key
            switch vTyped := v.(type) {
            case map[string]interface{}:
                out[k] = deepCopyMap(vTyped)
            case []interface{}:
                out[k] = copyList(vTyped)
            default:
                out[k] = v
            }
        }
    }
    return out
}

func deepCopyMap(m map[string]interface{}) map[string]interface{} {
    if m == nil {
        return nil
    }
    out := make(map[string]interface{}, len(m))
    for k, v := range m {
        switch vv := v.(type) {
        case map[string]interface{}:
            out[k] = deepCopyMap(vv)
        case []interface{}:
            out[k] = copyList(vv)
        default:
            out[k] = vv
        }
    }
    return out
}

func copyList(in []interface{}) []interface{} {
    if in == nil {
        return nil
    }
    out := make([]interface{}, len(in))
    for i, v := range in {
        switch vv := v.(type) {
        case map[string]interface{}:
            out[i] = deepCopyMap(vv)
        case []interface{}:
            out[i] = copyList(vv)
        default:
            out[i] = vv
        }
    }
    return out
}

func unionLists(existing interface{}, incoming []interface{}) []interface{} {
    var base []interface{}
    switch ex := existing.(type) {
    case []interface{}:
        base = copyList(ex)
    default:
        base = []interface{}{}
    }
    seen := make(map[string]struct{}, len(base))
    out := make([]interface{}, 0, len(base)+len(incoming))
    for _, v := range base {
        out = append(out, v)
        seen[valueKey(v)] = struct{}{}
    }
    for _, v := range incoming {
        key := valueKey(v)
        if _, ok := seen[key]; ok {
            continue
        }
        out = append(out, v)
        seen[key] = struct{}{}
    }
    return out
}

func valueKey(v interface{}) string {
    switch t := v.(type) {
    case string:
        return t
    case map[string]interface{}:
        // try to pick stable fields for tools
        if n, ok := t["name"].(string); ok {
            return n
        }
        if u, ok := t["uri"].(string); ok {
            return u
        }
        return fmt.Sprintf("%v", t)
    default:
        return fmt.Sprintf("%v", t)
    }
}


