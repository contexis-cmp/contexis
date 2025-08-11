package commands

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "io/fs"
    "strconv"
    "os/exec"
    "os"
    "path/filepath"
    "regexp"
    "sort"
    "strings"

    "gopkg.in/yaml.v3"
)

// Drift test spec structures map to tests/**/rag_drift_test.yaml
type driftTestSpec struct {
    TestCases            []driftTestCase     `yaml:"test_cases"`
    DriftThresholds      driftThresholds     `yaml:"drift_thresholds"`
    BusinessRules        []businessRule      `yaml:"business_rules"`
    PerformanceBenchmarks map[string]any     `yaml:"performance_benchmarks"`
}

type driftTestCase struct {
    Name              string   `yaml:"name"`
    Input             string   `yaml:"input"`
    ExpectedSimilarity float64  `yaml:"expected_similarity"`
    RequiredKeywords  []string `yaml:"required_keywords"`
    ForbiddenPhrases  []string `yaml:"forbidden_phrases"`
    RequiredResponse  string   `yaml:"required_response"`
    ExpectedFormat    string   `yaml:"expected_format"`
    RequiredSections  []string `yaml:"required_sections"`
    ForbiddenSections []string `yaml:"forbidden_sections"`
}

type driftThresholds struct {
    SimilarityThreshold   float64 `yaml:"similarity_threshold"`
    ResponseTimeThreshold int     `yaml:"response_time_threshold"`
    TokenCountThreshold   int     `yaml:"token_count_threshold"`
}

type businessRule struct {
    Name        string `yaml:"name"`
    Description string `yaml:"description"`
    Validation  string `yaml:"validation"`
}

// Public report structure written as JSON
type DriftRunReport struct {
    Component string             `json:"component"`
    SpecPath  string             `json:"spec_path"`
    Passed    int                `json:"passed"`
    Failed    int                `json:"failed"`
    Total     int                `json:"total"`
    Results   []DriftTestResult  `json:"results"`
}

type DriftTestResult struct {
    Name       string   `json:"name"`
    Status     string   `json:"status"` // PASSED | FAILED | ERROR
    Similarity float64  `json:"similarity,omitempty"`
    Threshold  float64  `json:"threshold,omitempty"`
    Reasons    []string `json:"reasons,omitempty"`
}

// RunDriftDetection discovers and executes all rag_drift_test.yaml specs
type DriftOptions struct {
    OutDir          string
    UpdateBaseline  bool
    UseSemantic     bool
    ComponentFilter string
    WriteJUnit      bool
}

func RunDriftDetection(ctx context.Context, projectRoot string, opts DriftOptions) error {
    if projectRoot == "" {
        var err error
        projectRoot, err = os.Getwd()
        if err != nil {
            return err
        }
    }

    specs, err := findDriftSpecs(projectRoot)
    if err != nil {
        return err
    }
    if len(specs) == 0 {
        return errors.New("no drift specs found (looking for tests/**/rag_drift_test.yaml)")
    }

    // Ensure output directory exists
    outDir := opts.OutDir
    if outDir == "" {
        outDir = filepath.Join(projectRoot, "tests", "reports")
    }
    if err := os.MkdirAll(outDir, 0o755); err != nil {
        return fmt.Errorf("create report dir: %w", err)
    }

    var (
        overallErr error
        index       = make([]DriftRunReport, 0, len(specs))
    )
    for _, specPath := range specs {
        comp := componentFromSpecPath(specPath)
        if opts.ComponentFilter != "" && !strings.EqualFold(opts.ComponentFilter, comp) {
            continue
        }
        rep, err := runSingleSpec(ctx, projectRoot, specPath, opts)
        if err != nil {
            overallErr = err
        }
        // Write JSON report per component
        if comp == "" {
            comp = "unknown"
        }
        outFile := filepath.Join(outDir, fmt.Sprintf("drift_%s.json", sanitizeFileName(comp)))
        if by, mErr := json.MarshalIndent(rep, "", "  "); mErr == nil {
            _ = os.WriteFile(outFile, by, 0o644)
        }
        index = append(index, rep)
        // Also print concise summary line
        fmt.Printf("Drift: %-16s passed=%d failed=%d total=%d\n", comp, rep.Passed, rep.Failed, rep.Total)
    }

    // Write index summary
    if by, mErr := json.MarshalIndent(index, "", "  "); mErr == nil {
        _ = os.WriteFile(filepath.Join(outDir, "drift_index.json"), by, 0o644)
    }
    if opts.WriteJUnit {
        _ = writeJUnit(filepath.Join(outDir, "junit-drift.xml"), index)
    }

    return overallErr
}

func findDriftSpecs(root string) ([]string, error) {
    var matches []string
    testsDir := filepath.Join(root, "tests")
    err := filepath.WalkDir(testsDir, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if !d.IsDir() && strings.EqualFold(filepath.Base(path), "rag_drift_test.yaml") {
            matches = append(matches, path)
        }
        return nil
    })
    if err != nil && !os.IsNotExist(err) {
        return nil, err
    }
    sort.Strings(matches)
    return matches, nil
}

func runSingleSpec(ctx context.Context, projectRoot, specPath string, opts DriftOptions) (DriftRunReport, error) {
    by, err := os.ReadFile(specPath)
    if err != nil {
        return DriftRunReport{}, fmt.Errorf("read spec: %w", err)
    }
    var spec driftTestSpec
    if err := yaml.Unmarshal(by, &spec); err != nil {
        return DriftRunReport{}, fmt.Errorf("parse yaml: %w", err)
    }
    component := componentFromSpecPath(specPath)

    // Load documents for naive similarity
    docs := loadComponentDocuments(projectRoot, component)

    report := DriftRunReport{Component: component, SpecPath: specPath}
    // Load baseline (if any)
    baseSim := loadBaseline(projectRoot, component)

    for _, tc := range spec.TestCases {
        var r DriftTestResult
        if opts.UseSemantic {
            r = evaluateTestCaseSemantic(projectRoot, component, tc, spec)
        } else {
            r = evaluateTestCase(tc, spec, docs)
        }

        // Compare against baseline if present and not updating
        if !opts.UpdateBaseline {
            if prev, ok := baseSim[tc.Name]; ok {
                delta := prev - r.Similarity
                // default alert threshold 0.15 unless overridden via env
                alert := 0.15
                if v := os.Getenv("DRIFT_ALERT_THRESHOLD"); v != "" {
                    if f, perr := parseFloat(v); perr == nil { alert = f }
                }
                if delta > alert {
                    r.Status = "FAILED"
                    r.Reasons = append(r.Reasons, fmt.Sprintf("drift delta %.3f exceeds alert threshold %.3f (baseline %.3f -> current %.3f)", delta, alert, prev, r.Similarity))
                }
            }
        }
        report.Results = append(report.Results, r)
        if r.Status == "PASSED" {
            report.Passed++
        } else {
            report.Failed++
        }
        report.Total++
    }

    // Update baseline if requested
    if opts.UpdateBaseline {
        sims := make(map[string]float64, len(report.Results))
        for _, r := range report.Results {
            sims[r.Name] = r.Similarity
        }
        _ = saveBaseline(projectRoot, component, sims)
    }
    return report, nil
}

func componentFromSpecPath(p string) string {
    // tests/<Component>/rag_drift_test.yaml
    dir := filepath.Dir(p)
    return filepath.Base(dir)
}

func loadComponentDocuments(root, component string) []string {
    // Read markdown and text files under memory/<Component>/documents
    var docs []string
    base := filepath.Join(root, "memory", component, "documents")
    _ = filepath.WalkDir(base, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return nil
        }
        if d.IsDir() {
            return nil
        }
        lower := strings.ToLower(path)
        if strings.HasSuffix(lower, ".md") || strings.HasSuffix(lower, ".txt") {
            if by, err := os.ReadFile(path); err == nil {
                docs = append(docs, string(by))
            }
        }
        return nil
    })
    return docs
}

func evaluateTestCase(tc driftTestCase, spec driftTestSpec, documents []string) DriftTestResult {
    res := DriftTestResult{Name: tc.Name}
    // Compute naive similarity as best match among docs
    best := 0.0
    for _, doc := range documents {
        sim := jaccardSimilarity(tokenize(tc.Input), tokenize(doc))
        if sim > best {
            best = sim
        }
    }
    res.Similarity = best
    // Determine threshold
    thr := tc.ExpectedSimilarity
    if thr == 0 {
        thr = spec.DriftThresholds.SimilarityThreshold
        if thr == 0 {
            thr = 0.7
        }
    }
    res.Threshold = thr

    // Evaluate business rule style checks against a pseudo-response (top doc content)
    var reasons []string
    if best < thr {
        reasons = append(reasons, fmt.Sprintf("similarity %.3f < threshold %.3f", best, thr))
    }
    // For keyword checks, use concatenation of all documents (cheap heuristic)
    joined := strings.ToLower(strings.Join(documents, "\n"))
    for _, kw := range tc.RequiredKeywords {
        if !containsWord(joined, strings.ToLower(kw)) {
            reasons = append(reasons, fmt.Sprintf("missing required keyword: %q", kw))
        }
    }
    for _, fp := range tc.ForbiddenPhrases {
        if strings.Contains(joined, strings.ToLower(fp)) {
            reasons = append(reasons, fmt.Sprintf("contains forbidden phrase: %q", fp))
        }
    }
    // expected_format markdown: basic check for a heading or markdown constructs
    if f := strings.ToLower(tc.ExpectedFormat); f == "markdown" {
        if !strings.Contains(joined, "#") && !strings.Contains(joined, "**") && !strings.Contains(joined, "-") {
            reasons = append(reasons, "expected markdown-like content not found")
        }
    }

    if len(reasons) == 0 {
        res.Status = "PASSED"
    } else {
        res.Status = "FAILED"
        res.Reasons = reasons
    }
    return res
}

// evaluateTestCaseSemantic attempts to use tools/<Component>/semantic_search.py to compute similarity
func evaluateTestCaseSemantic(projectRoot, component string, tc driftTestCase, spec driftTestSpec) DriftTestResult {
    r := DriftTestResult{Name: tc.Name}
    // Build a small python snippet to import the tool and run a search
    className := fmt.Sprintf("%sSemanticSearch", component)
    toolDir := filepath.Join(projectRoot, "tools", component)
    py := strings.Join([]string{
        "import sys, os, json",
        fmt.Sprintf("sys.path.append(%q)", toolDir),
        "from semantic_search import " + className + " as S",
        "s = S()",
        fmt.Sprintf("res = s.search(%q, top_k=1, threshold=0.0)", escapePyString(tc.Input)),
        "print(json.dumps(res[0]['similarity'] if res else 0.0))",
    }, ";")

    out, err := runPython(py)
    if err != nil {
        // fallback to naive if python fails
        return evaluateTestCase(tc, spec, loadComponentDocuments(projectRoot, component))
    }
    sim := parseJSONFloat(strings.TrimSpace(out))
    r.Similarity = sim
    thr := tc.ExpectedSimilarity
    if thr == 0 { thr = spec.DriftThresholds.SimilarityThreshold; if thr == 0 { thr = 0.7 } }
    r.Threshold = thr
    if sim >= thr {
        r.Status = "PASSED"
    } else {
        r.Status = "FAILED"
        r.Reasons = []string{fmt.Sprintf("similarity %.3f < threshold %.3f", sim, thr)}
    }
    return r
}

// --- Baseline helpers ---

func baselinePath(root, component string) string {
    return filepath.Join(root, "tests", component, "baselines", "drift_baseline.json")
}

func loadBaseline(root, component string) map[string]float64 {
    path := baselinePath(root, component)
    by, err := os.ReadFile(path)
    if err != nil {
        return map[string]float64{}
    }
    var m map[string]float64
    if err := json.Unmarshal(by, &m); err != nil {
        return map[string]float64{}
    }
    return m
}

func saveBaseline(root, component string, sims map[string]float64) error {
    path := baselinePath(root, component)
    if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
        return err
    }
    by, _ := json.MarshalIndent(sims, "", "  ")
    return os.WriteFile(path, by, 0o644)
}

// --- Utilities ---

func parseFloat(s string) (float64, error) { return strconv.ParseFloat(strings.TrimSpace(s), 64) }

func escapePyString(s string) string { return strings.ReplaceAll(s, "'", "\\'") }

func runPython(code string) (string, error) {
    cmd := exec.Command("python3", "-c", code)
    out, err := cmd.CombinedOutput()
    if err != nil {
        return string(out), err
    }
    return string(out), nil
}

func parseJSONFloat(s string) float64 {
    var f float64
    if err := json.Unmarshal([]byte(s), &f); err == nil {
        return f
    }
    // fallback parse
    if v, err := strconv.ParseFloat(s, 64); err == nil {
        return v
    }
    return 0
}

// Optional JUnit writer for CI integrations
func writeJUnit(path string, reports []DriftRunReport) error {
    var total, failures int
    for _, r := range reports { total += r.Total; failures += r.Failed }
    b := &strings.Builder{}
    fmt.Fprintf(b, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
    fmt.Fprintf(b, "<testsuite name=\"drift\" tests=\"%d\" failures=\"%d\">\n", total, failures)
    for _, r := range reports {
        for _, t := range r.Results {
            fmt.Fprintf(b, "  <testcase classname=\"%s\" name=\"%s\">\n", xmlEscape(r.Component), xmlEscape(t.Name))
            if t.Status != "PASSED" {
                msg := ""
                if len(t.Reasons) > 0 { msg = strings.Join(t.Reasons, "; ") }
                fmt.Fprintf(b, "    <failure message=\"%s\"/>\n", xmlEscape(msg))
            }
            fmt.Fprintf(b, "  </testcase>\n")
        }
    }
    fmt.Fprintf(b, "</testsuite>\n")
    if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil { return err }
    return os.WriteFile(path, []byte(b.String()), 0o644)
}

func xmlEscape(s string) string {
    r := strings.NewReplacer(
        "&", "&amp;",
        "<", "&lt;",
        ">", "&gt;",
        "\"", "&quot;",
        "'", "&apos;",
    )
    return r.Replace(s)
}



func tokenize(s string) []string {
    s = strings.ToLower(s)
    // keep words and numbers
    rx := regexp.MustCompile(`[a-z0-9]+`)
    return rx.FindAllString(s, -1)
}

func jaccardSimilarity(a, b []string) float64 {
    if len(a) == 0 || len(b) == 0 {
        return 0
    }
    setA := make(map[string]struct{}, len(a))
    setB := make(map[string]struct{}, len(b))
    for _, t := range a { setA[t] = struct{}{} }
    for _, t := range b { setB[t] = struct{}{} }
    var inter, uni int
    for t := range setA {
        uni++
        if _, ok := setB[t]; ok { inter++ }
    }
    for t := range setB {
        if _, ok := setA[t]; !ok { uni++ }
    }
    if uni == 0 { return 0 }
    return float64(inter) / float64(uni)
}

func containsWord(haystack, needle string) bool {
    if needle == "" { return true }
    // word boundary-ish check
    pattern := fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(needle))
    rx := regexp.MustCompile(pattern)
    return rx.FindStringIndex(haystack) != nil
}

func sanitizeFileName(s string) string {
    s = strings.TrimSpace(s)
    s = strings.ReplaceAll(s, " ", "_")
    rx := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
    return rx.ReplaceAllString(s, "-")
}


