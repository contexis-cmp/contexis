package commands

import (
	"context"
	"fmt"
	"os"
	"text/template"

	"github.com/contexis-cmp/contexis/src/cli/logger"
	"go.uber.org/zap"
)

// generateRAGMemory creates memory store configuration and structure
func generateRAGMemory(ctx context.Context, config RAGConfig) error {
	log := logger.WithContext(ctx)

	// Create memory configuration
	memoryConfigPath := fmt.Sprintf("memory/%s/memory_config.yaml", config.Name)

	memoryConfigTemplate := `# Memory configuration for {{.Name}} RAG system
vector_store:
  type: "{{.DBType}}"
  path: "./memory/{{.Name}}/vector_store.db"
  
embedding_model:
  name: "{{.Embeddings}}"
  dimensions: 384
  max_length: 512

document_processing:
  chunk_size: 700
  chunk_overlap: 120
  supported_formats: ["txt", "md", "pdf", "docx"]
  
indexing:
  batch_size: 100
  parallel_workers: 4
  similarity_metric: "cosine"
`

	tmpl, err := template.New("memory_config").Parse(memoryConfigTemplate)
	if err != nil {
		log.Error("failed to parse memory config template", zap.Error(err))
		return fmt.Errorf("failed to parse memory config template: %w", err)
	}

	file, err := os.Create(memoryConfigPath)
	if err != nil {
		log.Error("failed to create memory config file", zap.String("path", memoryConfigPath), zap.Error(err))
		return fmt.Errorf("failed to create memory config file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, config); err != nil {
		log.Error("failed to execute memory config template", zap.Error(err))
		return fmt.Errorf("failed to execute memory config template: %w", err)
	}

	// Create sample document
	sampleDocPath := fmt.Sprintf("memory/%s/documents/sample.md", config.Name)
	sampleDoc := `# Sample Document

This is a sample document for your {{.Name}} RAG system.

## Key Information

- This document demonstrates the document format
- Documents should be in Markdown format for best results
- Include relevant keywords and concepts for better search

## Usage

1. Replace this sample with your actual documents
2. Ensure documents are well-structured and contain relevant information
3. The RAG system will automatically index and search these documents

## Example Query

Try asking: "What is the key information in this document?"
`

	tmpl, err = template.New("sample_doc").Parse(sampleDoc)
	if err != nil {
		log.Error("failed to parse sample doc template", zap.Error(err))
		return fmt.Errorf("failed to parse sample doc template: %w", err)
	}

	file, err = os.Create(sampleDocPath)
	if err != nil {
		log.Error("failed to create sample document", zap.String("path", sampleDocPath), zap.Error(err))
		return fmt.Errorf("failed to create sample document: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, config); err != nil {
		log.Error("failed to execute sample doc template", zap.Error(err))
		return fmt.Errorf("failed to execute sample doc template: %w", err)
	}

	log.Info("RAG memory configuration generated", zap.String("path", memoryConfigPath))
	return nil
}
