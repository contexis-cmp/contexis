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
	// Set defaults if not provided
	if dbType == "" {
		dbType = "sqlite"
	}
	if embeddings == "" {
		embeddings = "sentence-transformers"
	}

	// Validate configuration
	if err := validateRAGConfig(dbType, embeddings); err != nil {
		logger.LogErrorColored(ctx, "RAG configuration validation failed", err)
		return fmt.Errorf("invalid RAG configuration: %w", err)
	}

	config := RAGConfig{
		Name:        name,
		DBType:      dbType,
		Embeddings:  embeddings,
		Description: fmt.Sprintf("RAG system for %s", name),
		Version:     "1.0.0",
	}

	logger.LogInfo(ctx, "Generating RAG system",
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

	logger.LogInfo(ctx, "Creating RAG directory structure")
	for _, dir := range ragDirs {
		if err := os.MkdirAll(dir, 0750); err != nil {
			logger.LogErrorColored(ctx, "failed to create directory", err, zap.String("directory", dir))
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		logger.LogDebugWithContext(ctx, "Created directory", zap.String("path", dir))
	}

	// Generate RAG components
	logger.LogInfo(ctx, "Generating RAG context")
	if err := generateRAGContext(ctx, config); err != nil {
		logger.LogErrorColored(ctx, "failed to generate RAG context", err)
		return fmt.Errorf("failed to generate RAG context: %w", err)
	}

	logger.LogInfo(ctx, "Generating RAG memory configuration")
	if err := generateRAGMemory(ctx, config); err != nil {
		logger.LogErrorColored(ctx, "failed to generate RAG memory", err)
		return fmt.Errorf("failed to generate RAG memory: %w", err)
	}

	logger.LogInfo(ctx, "Generating RAG prompts")
	if err := generateRAGPrompts(ctx, config); err != nil {
		logger.LogErrorColored(ctx, "failed to generate RAG prompts", err)
		return fmt.Errorf("failed to generate RAG prompts: %w", err)
	}

	logger.LogInfo(ctx, "Generating RAG tools")
	if err := generateRAGTools(ctx, config); err != nil {
		logger.LogErrorColored(ctx, "failed to generate RAG tools", err)
		return fmt.Errorf("failed to generate RAG tools: %w", err)
	}

	logger.LogInfo(ctx, "Generating RAG tests")
	if err := generateRAGTests(ctx, config); err != nil {
		logger.LogErrorColored(ctx, "failed to generate RAG tests", err)
		return fmt.Errorf("failed to generate RAG tests: %w", err)
	}

	logger.LogInfo(ctx, "Generating RAG configuration")
	if err := generateRAGConfig(ctx, config); err != nil {
		logger.LogErrorColored(ctx, "failed to generate RAG config", err)
		return fmt.Errorf("failed to generate RAG config: %w", err)
	}

	// Show generated structure and development flow
	showRAGStructure(name, config)
	showRAGDevelopmentFlow(name, config)

	return nil
}

// showRAGStructure displays the generated RAG system structure
func showRAGStructure(name string, config RAGConfig) {
	fmt.Printf("\n")
	logger.LogSuccess(context.Background(), "RAG system generated successfully",
		zap.String("name", name),
		zap.String("db_type", config.DBType),
		zap.String("embeddings", config.Embeddings))

	fmt.Printf("\n📁 Generated RAG Structure:\n")
	fmt.Printf("  contexts/%s/\n", name)
	fmt.Printf("  ├── 📄 %s.ctx\n", name)
	fmt.Printf("  ├── 📁 memory/%s/\n", name)
	fmt.Printf("  │   ├── 📄 memory_config.yaml\n")
	fmt.Printf("  │   ├── 📁 documents/\n")
	fmt.Printf("  │   └── 📁 embeddings/\n")
	fmt.Printf("  ├── 📁 prompts/%s/\n", name)
	fmt.Printf("  │   ├── 📄 search_response.md\n")
	fmt.Printf("  │   └── 📄 no_results.md\n")
	fmt.Printf("  ├── 📁 tools/%s/\n", name)
	fmt.Printf("  │   ├── 📄 semantic_search.py\n")
	fmt.Printf("  │   └── 📄 requirements.txt\n")
	fmt.Printf("  ├── 📁 tests/%s/\n", name)
	fmt.Printf("  │   ├── 📄 test_rag.py\n")
	fmt.Printf("  │   └── 📄 rag_drift_test.yaml\n")
	fmt.Printf("  └── 📄 config.yaml\n")
}

// showRAGDevelopmentFlow displays the RAG development workflow
func showRAGDevelopmentFlow(name string, config RAGConfig) {
	fmt.Printf("\n🚀 RAG Development Flow:\n")

	fmt.Printf("\n1️⃣  Add documents to your knowledge base:\n")
	fmt.Printf("   cp your-documents/* memory/%s/documents/\n", name)
	fmt.Printf("   # Or create sample documents:\n")
	fmt.Printf("   echo 'Your company policies here...' > memory/%s/documents/policies.txt\n", name)

	fmt.Printf("\n2️⃣  Ingest documents into the vector store:\n")
	fmt.Printf("   ctx memory ingest --provider=%s --component=%s --input=memory/%s/documents/policies.txt\n",
		config.DBType, name, name)

	fmt.Printf("\n3️⃣  Test your RAG system:\n")
	fmt.Printf("   ctx test %s\n", name)
	fmt.Printf("   ctx test --drift-detection --component=%s\n", name)

	fmt.Printf("\n4️⃣  Run queries against your knowledge base:\n")
	fmt.Printf("   ctx run %s \"What are your company policies?\"\n", name)
	fmt.Printf("   ctx run %s \"How do I request a refund?\"\n", name)

	fmt.Printf("\n5️⃣  Start development server:\n")
	fmt.Printf("   ctx serve --addr :8000\n")

	fmt.Printf("\n6️⃣  Monitor and improve:\n")
	fmt.Printf("   # Check drift detection results\n")
	fmt.Printf("   cat tests/reports/drift_%s.json\n", name)
	fmt.Printf("   \n")
	fmt.Printf("   # Add more documents as needed\n")
	fmt.Printf("   ctx memory ingest --provider=%s --component=%s --input=new_documents.txt\n",
		config.DBType, name)

	fmt.Printf("\n📚 Configuration Details:\n")
	fmt.Printf("   • Database: %s\n", config.DBType)
	fmt.Printf("   • Embeddings: %s\n", config.Embeddings)
	fmt.Printf("   • Context: contexts/%s/%s.ctx\n", name, name)
	fmt.Printf("   • Memory: memory/%s/memory_config.yaml\n", name)

	fmt.Printf("\n🎉 Your RAG system is ready! Start adding documents and testing queries.\n")
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
