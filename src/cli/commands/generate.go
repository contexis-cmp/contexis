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

	// Initialize logger
	if err := logger.InitLogger("info", "json"); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	log := logger.WithContext(ctx).With(
		zap.String("generator_type", generatorType),
		zap.String("name", name),
		zap.String("operation", "generate"),
	)

	// Log operation start
	defer logger.LogOperation(ctx, "component_generation",
		zap.String("generator_type", generatorType),
		zap.String("name", name))()

	// Validate generator type
	validTypes := []string{"rag", "agent", "workflow"}
	isValid := false
	for _, validType := range validTypes {
		if generatorType == validType {
			isValid = true
			break
		}
	}

	if !isValid {
		log.Error("invalid generator type", zap.String("type", generatorType))
		return fmt.Errorf("invalid generator type '%s'. Valid types: %s", generatorType, strings.Join(validTypes, ", "))
	}

	// Validate component name
	if err := config.ValidateProjectName(name); err != nil {
		log.Error("component name validation failed", zap.Error(err))
		return fmt.Errorf("invalid component name: %w", err)
	}

	// Get flags
	dbType, _ := cmd.Flags().GetString("db")
	embeddings, _ := cmd.Flags().GetString("embeddings")
	tools, _ := cmd.Flags().GetString("tools")
	memory, _ := cmd.Flags().GetString("memory")
	steps, _ := cmd.Flags().GetString("steps")

	// Generate based on type
	switch generatorType {
	case "rag":
		return generateRAG(ctx, name, dbType, embeddings)
	case "agent":
		return GenerateAgent(ctx, name, tools, memory)
    case "workflow":
		return GenerateWorkflow(ctx, name, steps)
    case "plugin":
        return GeneratePlugin(ctx, name)
	default:
		return fmt.Errorf("generator type '%s' not implemented yet", generatorType)
	}
}

func init() {
	// Add flags for different generator types
	GenerateCmd.Flags().String("db", "sqlite", "Database type for RAG (sqlite, postgres, chroma)")
	GenerateCmd.Flags().String("embeddings", "sentence-transformers", "Embedding model (sentence-transformers, openai, cohere)")
	GenerateCmd.Flags().String("tools", "", "Comma-separated list of tools for agent")
	GenerateCmd.Flags().String("memory", "episodic", "Memory type for agent (episodic, none)")
	GenerateCmd.Flags().String("steps", "", "Comma-separated list of workflow steps")
}
