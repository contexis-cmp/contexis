package commands

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/contexis-cmp/contexis/src/cli/logger"
	"go.uber.org/zap"
)

// RAGConfig holds configuration for RAG system generation
type RAGConfig struct {
	Name        string
	DBType      string
	Embeddings  string
	Description string
	Version     string
}

// generateRAG creates a complete RAG system
func generateRAG(ctx context.Context, name, dbType, embeddings string) error {
	log := logger.WithContext(ctx)

	// Set defaults if not provided
	if dbType == "" {
		dbType = "sqlite"
	}
	if embeddings == "" {
		embeddings = "sentence-transformers"
	}

	// Validate configuration
	if err := validateRAGConfig(dbType, embeddings); err != nil {
		log.Error("RAG configuration validation failed", zap.Error(err))
		return fmt.Errorf("invalid RAG configuration: %w", err)
	}

	config := RAGConfig{
		Name:        name,
		DBType:      dbType,
		Embeddings:  embeddings,
		Description: fmt.Sprintf("RAG system for %s", name),
		Version:     "1.0.0",
	}

	log.Info("generating RAG system",
		zap.String("name", name),
		zap.String("db_type", dbType),
		zap.String("embeddings", embeddings))

	// Create RAG-specific directory structure
	ragDirs := []string{
		fmt.Sprintf("contexts/%s", name),
		fmt.Sprintf("memory/%s", name),
		fmt.Sprintf("memory/%s/documents", name),
		fmt.Sprintf("memory/%s/embeddings", name),
		fmt.Sprintf("prompts/%s", name),
		fmt.Sprintf("tools/%s", name),
		fmt.Sprintf("tests/%s", name),
	}

	for _, dir := range ragDirs {
		if err := os.MkdirAll(dir, 0750); err != nil {
			log.Error("failed to create directory", zap.String("directory", dir), zap.Error(err))
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Generate RAG components
	if err := generateRAGContext(ctx, config); err != nil {
		log.Error("failed to generate RAG context", zap.Error(err))
		return fmt.Errorf("failed to generate RAG context: %w", err)
	}

	if err := generateRAGMemory(ctx, config); err != nil {
		log.Error("failed to generate RAG memory", zap.Error(err))
		return fmt.Errorf("failed to generate RAG memory: %w", err)
	}

	if err := generateRAGPrompts(ctx, config); err != nil {
		log.Error("failed to generate RAG prompts", zap.Error(err))
		return fmt.Errorf("failed to generate RAG prompts: %w", err)
	}

	if err := generateRAGTools(ctx, config); err != nil {
		log.Error("failed to generate RAG tools", zap.Error(err))
		return fmt.Errorf("failed to generate RAG tools: %w", err)
	}

	if err := generateRAGTests(ctx, config); err != nil {
		log.Error("failed to generate RAG tests", zap.Error(err))
		return fmt.Errorf("failed to generate RAG tests: %w", err)
	}

	if err := generateRAGConfig(ctx, config); err != nil {
		log.Error("failed to generate RAG config", zap.Error(err))
		return fmt.Errorf("failed to generate RAG config: %w", err)
	}

	log.Info("RAG system generated successfully",
		zap.String("name", name),
		zap.String("path", fmt.Sprintf("contexts/%s", name)))

	fmt.Printf("✅ Successfully generated RAG system: %s\n", name)
	fmt.Printf("📁 RAG components created in: contexts/%s/\n", name)
	fmt.Printf("\n🚀 Next steps:\n")
	fmt.Printf("  # Add documents to your knowledge base\n")
	fmt.Printf("  cp your-documents/* memory/%s/documents/\n", name)
	fmt.Printf("  \n")
	fmt.Printf("  # Test your RAG system\n")
	fmt.Printf("  ctx test %s\n", name)
	fmt.Printf("  \n")
	fmt.Printf("  # Run a query\n")
	fmt.Printf("  ctx run %s \"What is your question?\"\n", name)

	return nil
}

// validateRAGConfig validates RAG configuration parameters
func validateRAGConfig(dbType, embeddings string) error {
	validDBs := []string{"sqlite", "postgres", "chroma"}
	validEmbeddings := []string{"sentence-transformers", "openai", "cohere", "bge-small-en"}

	// Validate database type
	isValidDB := false
	for _, valid := range validDBs {
		if dbType == valid {
			isValidDB = true
			break
		}
	}
	if !isValidDB {
		return fmt.Errorf("invalid database type '%s'. Valid types: %s", dbType, strings.Join(validDBs, ", "))
	}

	// Validate embeddings model
	isValidEmbeddings := false
	for _, valid := range validEmbeddings {
		if embeddings == valid {
			isValidEmbeddings = true
			break
		}
	}
	if !isValidEmbeddings {
		return fmt.Errorf("invalid embeddings model '%s'. Valid models: %s", embeddings, strings.Join(validEmbeddings, ", "))
	}

	return nil
}
