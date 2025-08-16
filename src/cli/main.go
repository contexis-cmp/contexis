// Package main provides the Contexis CMP Framework CLI application.
//
// The Contexis CLI is a Rails-inspired command-line interface for building
// reproducible AI applications using the Context-Memory-Prompt (CMP) architecture.
// It provides commands for project initialization, component generation, testing,
// and deployment of AI applications.
//
// Key Features:
//   - Local-first development with out-of-the-box local models
//   - Component generation (RAG, agents, workflows)
//   - Memory management and vector search
//   - Drift detection and testing
//   - Production migration tools
//
// Example Usage:
//
//	# Initialize a new project
//	ctx init my-ai-app
//
//	# Generate a RAG component
//	ctx generate rag CustomerDocs
//
//	# Run tests with drift detection
//	ctx test --drift-detection
//
//	# Start development server
//	ctx serve --addr :8000
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

// rootCmd is the main command for the Contexis CLI.
// It serves as the entry point for all subcommands and provides
// the overall application description and version information.
var rootCmd = &cobra.Command{
	Use:   "ctx",
	Short: "Contexis CMP Framework CLI",
	Long: `Contexis is a Rails-inspired framework for building reproducible AI applications.
	
The Context-Memory-Prompt (CMP) architecture treats AI components as version-controlled,
first-class citizens, bringing architectural discipline to AI application engineering.`,
	Version: "0.1.14",
}

// init initializes the CLI by adding all subcommands to the root command.
// This function is called automatically when the package is imported.
func init() {
	// Add subcommands
	rootCmd.AddCommand(commands.InitCmd)
	rootCmd.AddCommand(commands.GenerateCmd)
	
	// Plugin commands (use current working directory as project root)
	if cwd, err := os.Getwd(); err == nil {
		rootCmd.AddCommand(commands.GetPluginCommand(cwd))
	} else {
		rootCmd.AddCommand(commands.GetPluginCommand(""))
	}
	
	// Context command with runtime ops
	rootCmd.AddCommand(commands.GetContextCommand(""))
	
	// Memory command
	rootCmd.AddCommand(commands.GetMemoryCommand())
	
	// Prompt command
	rootCmd.AddCommand(commands.GetPromptCommand())
	
	// Lock command
	rootCmd.AddCommand(commands.GetLockCommand())
	rootCmd.AddCommand(commands.GetPromptLintCommand())
	rootCmd.AddCommand(testCmd)
	
	// Build/Deploy commands
	rootCmd.AddCommand(commands.GetBuildCommand())
	rootCmd.AddCommand(commands.GetDeployCommand())
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(commands.GetServeCommand())
	rootCmd.AddCommand(commands.GetRunCommand())
	rootCmd.AddCommand(commands.GetWorkerCommand())
	rootCmd.AddCommand(commands.GetHFCommand())
	rootCmd.AddCommand(commands.GetModelsCommand())
	rootCmd.AddCommand(commands.GetMigrateCommand())
}

// main is the entry point for the Contexis CLI application.
// It initializes the logger, creates a request context, and executes
// the root command with proper error handling.
func main() {
	// Initialize colored logger
	if err := logger.InitColoredLogger("info"); err != nil {
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

// generateRequestID generates a simple request ID for tracing.
// It uses the current process ID to create a unique identifier
// for request tracking and debugging purposes.
//
// Returns:
//   - string: A request ID in the format "req_<pid>"
func generateRequestID() string {
	return fmt.Sprintf("req_%d", os.Getpid())
}

// testCmd provides comprehensive testing functionality for CMP components.
// It supports both drift detection tests and traditional Go test suites
// with various configuration options for different testing scenarios.
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run CMP tests",
	Long:  `Execute drift detection, correctness tests, and other CMP-specific validations.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running CMP tests...")
		
		// Parse drift detection flags
		driftOnly, _ := cmd.Flags().GetBool("drift-detection")
		outDir, _ := cmd.Flags().GetString("out")
		updateBaseline, _ := cmd.Flags().GetBool("update-baseline")
		useSemantic, _ := cmd.Flags().GetBool("semantic")
		component, _ := cmd.Flags().GetString("component")
		writeJUnit, _ := cmd.Flags().GetBool("junit")

		if driftOnly {
			// Execute drift detection tests
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

// initTestFlags initializes the command-line flags for the test command.
// This function is called automatically when the package is imported.
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

// versionCmd displays version information for the Contexis CLI.
// It shows the current version number and can be extended to include
// additional build information like commit hash and build date.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Contexis CMP Framework v0.1.14")
	},
}
