package main

import (
	"context"
	"fmt"
	"os"

	"github.com/contexis-cmp/contexis/src/cli/commands"
	"github.com/contexis-cmp/contexis/src/cli/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "ctx",
	Short: "Contexis CMP Framework CLI",
	Long: `Contexis is a Rails-inspired framework for building reproducible AI applications.
	
The Context-Memory-Prompt (CMP) architecture treats AI components as version-controlled,
first-class citizens, bringing architectural discipline to AI application engineering.`,
	Version: "0.1.0",
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(commands.InitCmd)
	rootCmd.AddCommand(commands.GenerateCmd)
	// Context command with runtime ops
	rootCmd.AddCommand(commands.GetContextCommand(""))
	// Memory command
	rootCmd.AddCommand(commands.GetMemoryCommand())
	// Prompt command
	rootCmd.AddCommand(commands.GetPromptCommand())
	// Lock command
	rootCmd.AddCommand(commands.GetLockCommand())
	rootCmd.AddCommand(commands.GetPromptLintCommand())
	rootCmd.AddCommand(generateCmd)
    rootCmd.AddCommand(testCmd)
    // Build/Deploy commands
    rootCmd.AddCommand(commands.GetBuildCommand())
    rootCmd.AddCommand(commands.GetDeployCommand())
	rootCmd.AddCommand(versionCmd)
    rootCmd.AddCommand(commands.GetServeCommand())
    rootCmd.AddCommand(commands.GetWorkerCommand())
    rootCmd.AddCommand(commands.GetHFCommand())
}

func main() {
	// Initialize logger
	if err := logger.InitLogger("info", "json"); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	log := logger.GetLogger()

	// Create context with request ID
	ctx := context.WithValue(context.Background(), "request_id", generateRequestID())

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		log.Error("command execution failed", zap.Error(err))
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// generateRequestID generates a simple request ID for tracing
func generateRequestID() string {
	return fmt.Sprintf("req_%d", os.Getpid())
}

var generateCmd = &cobra.Command{
	Use:   "generate [type] [name]",
	Short: "Generate a new CMP component",
	Long:  `Generate RAG systems, agents, workflows, or other CMP components using templates.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Generating %s component: %s\n", args[0], args[1])
		// TODO: Implement component generation
	},
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run CMP tests",
	Long:  `Execute drift detection, correctness tests, and other CMP-specific validations.`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Running CMP tests...")
        driftOnly, _ := cmd.Flags().GetBool("drift-detection")
        outDir, _ := cmd.Flags().GetString("out")
        updateBaseline, _ := cmd.Flags().GetBool("update-baseline")
        useSemantic, _ := cmd.Flags().GetBool("semantic")
        component, _ := cmd.Flags().GetString("component")
        writeJUnit, _ := cmd.Flags().GetBool("junit")

        if driftOnly {
            opts := commands.DriftOptions{
                OutDir:          outDir,
                UpdateBaseline:  updateBaseline,
                UseSemantic:     useSemantic,
                ComponentFilter: component,
                WriteJUnit:      writeJUnit,
            }
            if err := commands.RunDriftDetection(cmd.Context(), "", opts); err != nil {
                fmt.Fprintf(os.Stderr, "Drift detection failed: %v\n", err)
                os.Exit(1)
            }
            return
        }

        // Default: run Go unit/integration/e2e tests with optional coverage and category mapping
        runAll, _ := cmd.Flags().GetBool("all")
        unit, _ := cmd.Flags().GetBool("unit")
        integ, _ := cmd.Flags().GetBool("integration")
        e2e, _ := cmd.Flags().GetBool("e2e")
        category, _ := cmd.Flags().GetString("category")
        coverage, _ := cmd.Flags().GetBool("coverage")
        writeJUnit, _ = cmd.Flags().GetBool("junit")

        gOpts := commands.TestRunOptions{
            OutDir:     outDir,
            RunUnit:    unit,
            RunInt:     integ,
            RunE2E:     e2e,
            UseAll:     runAll,
            Category:   category,
            Coverage:   coverage,
            WriteJUnit: writeJUnit,
        }
        if err := commands.RunGoTests(cmd.Context(), "", gOpts); err != nil {
            fmt.Fprintf(os.Stderr, "Go tests failed: %v\n", err)
            os.Exit(1)
        }
    },
}

func init() {
    // Flags for test command
    testCmd.Flags().Bool("drift-detection", false, "Run drift detection tests only")
    testCmd.Flags().String("out", "", "Output directory for test reports")
    testCmd.Flags().Bool("update-baseline", false, "Update drift baselines with current results")
    testCmd.Flags().Bool("semantic", false, "Use semantic similarity via tools/<Component>/semantic_search.py if available")
    testCmd.Flags().String("component", "", "Limit to a single component (e.g., CustomerDocs)")
    testCmd.Flags().Bool("junit", false, "Write JUnit XML report for CI integration")
    // Go test selection
    testCmd.Flags().Bool("all", false, "Run all Go test suites (unit, integration, e2e)")
    testCmd.Flags().Bool("unit", false, "Run Go unit tests only")
    testCmd.Flags().Bool("integration", false, "Run Go integration tests only")
    testCmd.Flags().Bool("e2e", false, "Run Go end-to-end tests only")
    testCmd.Flags().String("category", "", "Run tests for a configured category from tests/test_config.yaml")
    testCmd.Flags().Bool("coverage", false, "Collect coverage and enforce thresholds from tests/test_config.yaml")
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Contexis CMP Framework v0.1.0")
	},
}
