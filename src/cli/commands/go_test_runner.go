package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type TestRunOptions struct {
	OutDir     string
	RunUnit    bool
	RunInt     bool
	RunE2E     bool
	UseAll     bool
	Category   string
	Coverage   bool
	WriteJUnit bool
}

type testConfig struct {
	TestSuites map[string]struct {
		Enabled           bool   `yaml:"enabled"`
		Timeout           string `yaml:"timeout"`
		Parallel          bool   `yaml:"parallel"`
		CoverageThreshold int    `yaml:"coverage_threshold"`
	} `yaml:"test_suites"`
	TestCategories []struct {
		Name        string   `yaml:"name"`
		Description string   `yaml:"description"`
		Suites      []string `yaml:"suites"`
	} `yaml:"test_categories"`
}

type GoSuiteResult struct {
	Suite       string  `json:"suite"`
	Passed      bool    `json:"passed"`
	CoveragePct float64 `json:"coverage_pct,omitempty"`
	Threshold   int     `json:"threshold,omitempty"`
	OutputPath  string  `json:"output_path"`
	ProfilePath string  `json:"profile_path,omitempty"`
	Error       string  `json:"error,omitempty"`
}

type GoTestsReport struct {
	Results []GoSuiteResult `json:"results"`
}

func RunGoTests(ctx context.Context, projectRoot string, opts TestRunOptions) error {
	if projectRoot == "" {
		var err error
		projectRoot, err = os.Getwd()
		if err != nil {
			return err
		}
	}
	cfg, err := loadTestConfig(projectRoot)
	if err != nil {
		return fmt.Errorf("failed to load test config: %w", err)
	}

	outDir := opts.OutDir
	if outDir == "" {
		outDir = filepath.Join(projectRoot, "tests", "reports")
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	covDir := filepath.Join(projectRoot, "tests", "coverage")
	if opts.Coverage {
		_ = os.MkdirAll(covDir, 0o755)
	}

	// Determine suites to run
	suites := determineSuites(cfg, opts)
	if len(suites) == 0 {
		return errors.New("no test suites selected")
	}

	var report GoTestsReport
	overallPassed := true

	fmt.Println("Running CMP tests...")
	fmt.Println()

	for _, s := range suites {
		res := runGoSuite(projectRoot, s, outDir, covDir, cfg, opts)
		report.Results = append(report.Results, res)

		// Print detailed test results
		printTestResult(s, res, opts)

		if !res.Passed {
			overallPassed = false
		}
	}

	// Print summary
	fmt.Println()
	printTestSummary(report, overallPassed)

	// Write aggregated JSON
	by, _ := json.MarshalIndent(report, "", "  ")
	_ = os.WriteFile(filepath.Join(outDir, "go_tests.json"), by, 0o644)

	if opts.WriteJUnit {
		_ = writeGoJUnit(filepath.Join(outDir, "junit-go.xml"), report)
	}

	if !overallPassed {
		return errors.New("one or more Go test suites failed or did not meet coverage thresholds")
	}
	return nil
}

func determineSuites(cfg testConfig, opts TestRunOptions) []string {
	if opts.Category != "" {
		for _, c := range cfg.TestCategories {
			if strings.EqualFold(c.Name, opts.Category) {
				return c.Suites
			}
		}
	}
	var suites []string
	if opts.UseAll || (!opts.RunUnit && !opts.RunInt && !opts.RunE2E) {
		// respect enabled flags
		if cfg.TestSuites["unit"].Enabled {
			suites = append(suites, "unit")
		}
		if cfg.TestSuites["integration"].Enabled {
			suites = append(suites, "integration")
		}
		if cfg.TestSuites["e2e"].Enabled {
			suites = append(suites, "e2e")
		}

		return suites
	}
	if opts.RunUnit {
		suites = append(suites, "unit")
	}
	if opts.RunInt {
		suites = append(suites, "integration")
	}
	if opts.RunE2E {
		suites = append(suites, "e2e")
	}
	return suites
}

func runGoSuite(root, suite, outDir, covDir string, cfg testConfig, opts TestRunOptions) GoSuiteResult {
	res := GoSuiteResult{Suite: suite}
	pattern := "./tests/" + suite + "/..."
	args := []string{"test", pattern, "-v"}
	var profile string
	if opts.Coverage {
		profile = filepath.Join(covDir, suite+".out")
		args = append(args, "-coverprofile="+profile)
		res.ProfilePath = profile
	}
	cmd := exec.Command("go", args...)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	outPath := filepath.Join(outDir, "go_"+suite+".txt")
	_ = os.WriteFile(outPath, out, 0o644)
	res.OutputPath = outPath

	// Check if the error is due to no packages to test
	if err != nil && strings.Contains(string(out), "no packages to test") {
		// This is not a real failure - just no tests to run
		res.Passed = true
		res.Error = "no test packages found (suite may be empty)"
	} else if err != nil {
		res.Passed = false
		res.Error = err.Error()
	} else {
		res.Passed = true
	}

	// Coverage enforcement
	if opts.Coverage && profile != "" {
		pct := coveragePercent(root, profile)
		res.CoveragePct = pct
		res.Threshold = cfg.TestSuites[suite].CoverageThreshold
		if res.Threshold > 0 && pct < float64(res.Threshold) {
			res.Passed = false
			if res.Error == "" {
				res.Error = fmt.Sprintf("coverage %.1f%% below threshold %d%%", pct, res.Threshold)
			} else {
				res.Error += "; " + fmt.Sprintf("coverage %.1f%% < %d%%", pct, res.Threshold)
			}
		}
	}
	return res
}

// printTestResult prints detailed information about a test suite result
func printTestResult(suite string, res GoSuiteResult, opts TestRunOptions) {
	// Print suite header
	fmt.Printf("📋 %s Tests:\n", strings.Title(suite))

	if res.Passed {
		fmt.Printf("   ✅ Status: PASSED\n")
	} else {
		fmt.Printf("   ❌ Status: FAILED\n")
	}

	// Print error details if failed
	if !res.Passed && res.Error != "" {
		fmt.Printf("   💥 Error: %s\n", res.Error)
	}

	// Print coverage information
	if opts.Coverage && res.CoveragePct > 0 {
		coverageStatus := "✅"
		if res.Threshold > 0 && res.CoveragePct < float64(res.Threshold) {
			coverageStatus = "❌"
		}
		fmt.Printf("   📊 Coverage: %s %.1f%% (threshold: %d%%)\n", coverageStatus, res.CoveragePct, res.Threshold)
	}

	// Print output file location
	if res.OutputPath != "" {
		fmt.Printf("   📄 Details: %s\n", res.OutputPath)
	}

	// If failed, show a snippet of the error output
	if !res.Passed && res.OutputPath != "" {
		if content, err := os.ReadFile(res.OutputPath); err == nil {
			lines := strings.Split(string(content), "\n")
			fmt.Printf("   🔍 Error Details:\n")
			// Show last 5 lines of output for context
			start := len(lines) - 5
			if start < 0 {
				start = 0
			}
			for i := start; i < len(lines); i++ {
				if strings.TrimSpace(lines[i]) != "" {
					fmt.Printf("      %s\n", lines[i])
				}
			}
		}
	}

	fmt.Println()
}

// printTestSummary prints a summary of all test results
func printTestSummary(report GoTestsReport, overallPassed bool) {
	total := len(report.Results)
	passed := 0
	failed := 0

	for _, res := range report.Results {
		if res.Passed {
			passed++
		} else {
			failed++
		}
	}

	fmt.Printf("📊 Test Summary:\n")
	fmt.Printf("   Total Suites: %d\n", total)
	fmt.Printf("   Passed: %d\n", passed)
	fmt.Printf("   Failed: %d\n", failed)

	if overallPassed {
		fmt.Printf("   🎉 Overall Status: ALL TESTS PASSED\n")
	} else {
		fmt.Printf("   💥 Overall Status: SOME TESTS FAILED\n")
		fmt.Printf("\n💡 Tip: Check the detailed output files in tests/reports/ for more information\n")
	}
}

func coveragePercent(root, profile string) float64 {
	// go tool cover -func=profile | tail -n 1 | parse total
	cmd := exec.Command("go", "tool", "cover", "-func="+profile)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0
	}
	lines := strings.Split(string(out), "\n")
	totalRe := regexp.MustCompile(`^total:\s*\(statements\)\s*([0-9.]+)%`)
	for _, ln := range lines {
		if m := totalRe.FindStringSubmatch(strings.TrimSpace(ln)); len(m) == 2 {
			return parseFloatFallback(m[1])
		}
	}
	return 0
}

func parseFloatFallback(s string) float64 {
	s = strings.TrimSpace(s)
	var f float64
	if _, err := fmt.Sscanf(s, "%f", &f); err == nil {
		return f
	}
	return 0
}

func loadTestConfig(root string) (testConfig, error) {
	path := filepath.Join(root, "tests", "test_config.yaml")
	by, err := os.ReadFile(path)
	if err != nil {
		return testConfig{}, err
	}
	var cfg testConfig
	if err := yaml.Unmarshal(by, &cfg); err != nil {
		return testConfig{}, err
	}
	return cfg, nil
}

// Minimal JUnit writer for Go suites
func writeGoJUnit(path string, rep GoTestsReport) error {
	var total, failures int
	total = len(rep.Results)
	for _, r := range rep.Results {
		if !r.Passed {
			failures++
		}
	}
	b := &strings.Builder{}
	fmt.Fprintf(b, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	fmt.Fprintf(b, "<testsuite name=\"go\" tests=\"%d\" failures=\"%d\">\n", total, failures)
	for _, r := range rep.Results {
		fmt.Fprintf(b, "  <testcase classname=\"%s\" name=\"suite\">\n", xmlEscape(r.Suite))
		if !r.Passed {
			msg := r.Error
			if msg == "" {
				msg = "failed"
			}
			fmt.Fprintf(b, "    <failure message=\"%s\"/>\n", xmlEscape(msg))
		}
		fmt.Fprintf(b, "  </testcase>\n")
	}
	fmt.Fprintf(b, "</testsuite>\n")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}
