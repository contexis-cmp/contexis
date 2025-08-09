package commands

import (
	"context"
	"fmt"
	"os"
	"text/template"

	"go.uber.org/zap"
	"github.com/contexis/cmp/src/cli/logger"
)

// generateRAGContext creates the RAG agent context file
func generateRAGContext(ctx context.Context, config RAGConfig) error {
	log := logger.WithContext(ctx)

	contextPath := fmt.Sprintf("contexts/%s/rag_agent.ctx", config.Name)
	
	contextTemplate := `name: "{{.Name}} RAG Agent"
version: "{{.Version}}"
description: "{{.Description}}"

role:
  persona: "Knowledge assistant with access to document repository"
  capabilities: 
    - "semantic_search"
    - "document_retrieval" 
    - "context_aware_responses"
  limitations:
    - "only_answers_based_on_documents"
    - "no_speculation_or_opinion"
    - "requires_relevant_context"

tools:
  - name: "semantic_search"
    uri: "mcp://search.semantic"
    description: "Search documents using semantic similarity"
  - name: "document_retrieval"
    uri: "mcp://documents.get"
    description: "Retrieve specific documents by ID"

guardrails:
  tone: "helpful_and_informative"
  format: "markdown"
  max_tokens: 1000
  temperature: 0.1
  similarity_threshold: 0.7

memory:
  vector_store: "{{.DBType}}"
  embedding_model: "{{.Embeddings}}"
  max_results: 5
  chunk_size: 700
  chunk_overlap: 120

testing:
  drift_threshold: 0.85
  business_rules:
    - "must_cite_source_documents"
    - "no_information_outside_documents"
    - "consistent_response_format"
`

	tmpl, err := template.New("rag_context").Parse(contextTemplate)
	if err != nil {
		log.Error("failed to parse context template", zap.Error(err))
		return fmt.Errorf("failed to parse context template: %w", err)
	}

	file, err := os.Create(contextPath)
	if err != nil {
		log.Error("failed to create context file", zap.String("path", contextPath), zap.Error(err))
		return fmt.Errorf("failed to create context file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, config); err != nil {
		log.Error("failed to execute context template", zap.Error(err))
		return fmt.Errorf("failed to execute context template: %w", err)
	}

	log.Info("RAG context generated", zap.String("path", contextPath))
	return nil
}
