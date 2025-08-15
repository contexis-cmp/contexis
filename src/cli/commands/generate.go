package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/contexis-cmp/contexis/src/cli/config"
	"github.com/contexis-cmp/contexis/src/cli/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var GenerateCmd = &cobra.Command{
	Use:   "generate [type] [name]",
	Short: "Generate a new CMP component",
	Long: `Generate RAG systems, agents, workflows, or other CMP components using templates.

Available generators:
  rag       - Knowledge-based retrieval systems
  agent     - Conversational agents with tools
  workflow  - Multi-step AI processing pipelines
  plugin    - Scaffolds a plugin template

Examples:
  ctx generate rag CustomerDocs --db=sqlite --embeddings=openai
  ctx generate agent SupportBot --tools=web_search,database --memory=episodic
  ctx generate workflow ContentPipeline --steps=research,write,review`,
	Args: cobra.ExactArgs(2),
	RunE: runGenerate,
}

func runGenerate(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	generatorType := strings.ToLower(args[0])
	name := args[1]

	// Use colored logger (already initialized in main.go)

	// Log operation start with colored logging
	done := logger.LogOperationColored(ctx, "component_generation",
		zap.String("generator_type", generatorType),
		zap.String("name", name))

	// Validate generator type
	validTypes := []string{"rag", "agent", "workflow", "plugin"}
	isValid := false
	for _, validType := range validTypes {
		if generatorType == validType {
			isValid = true
			break
		}
	}

	if !isValid {
		logger.LogErrorColored(ctx, "invalid generator type", fmt.Errorf("invalid type '%s'", generatorType), zap.String("type", generatorType))
		return fmt.Errorf("invalid generator type '%s'. Valid types: %s", generatorType, strings.Join(validTypes, ", "))
	}

	// Validate component name
	if err := config.ValidateProjectName(name); err != nil {
		logger.LogErrorColored(ctx, "component name validation failed", err)
		return fmt.Errorf("invalid component name: %w", err)
	}

	// Get flags
	dbType, _ := cmd.Flags().GetString("db")
	embeddings, _ := cmd.Flags().GetString("embeddings")
	tools, _ := cmd.Flags().GetString("tools")
	memory, _ := cmd.Flags().GetString("memory")
	steps, _ := cmd.Flags().GetString("steps")

	// Generate based on type
	var result error
	switch generatorType {
	case "rag":
		result = generateRAG(ctx, name, dbType, embeddings)
	case "agent":
		result = GenerateAgent(ctx, name, tools, memory)
	case "workflow":
		result = GenerateWorkflow(ctx, name, steps)
	case "plugin":
		result = GeneratePlugin(ctx, name)
	default:
		result = fmt.Errorf("generator type '%s' not implemented yet", generatorType)
	}

	// Complete the operation
	done()

	return result
}

func init() {
	// Add flags for different generator types
	GenerateCmd.Flags().String("db", "sqlite", "Database type for RAG (sqlite, postgres, chroma)")
	GenerateCmd.Flags().String("embeddings", "sentence-transformers", "Embedding model (sentence-transformers, openai, cohere)")
	GenerateCmd.Flags().String("tools", "", "Comma-separated list of tools for agent")
	GenerateCmd.Flags().String("memory", "episodic", "Memory type for agent (episodic, none)")
	GenerateCmd.Flags().String("steps", "", "Comma-separated list of workflow steps")
}
