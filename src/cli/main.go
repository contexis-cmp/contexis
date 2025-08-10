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
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(versionCmd)
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
		// TODO: Implement test execution
	},
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy CMP application",
	Long:  `Deploy the current CMP application to the configured environment.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Deploying CMP application...")
		// TODO: Implement deployment
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Contexis CMP Framework v0.1.0")
	},
}
