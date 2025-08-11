package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/contexis-cmp/contexis/src/cli/config"
	"github.com/contexis-cmp/contexis/src/cli/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var InitCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a new CMP project",
	Long: `Create a new Contexis project with the standard directory structure and configuration files.

This command will create:
- Project directory structure
- Configuration files
- Basic templates
- Test setup
- Documentation files`,
	Args: cobra.ExactArgs(1),
	RunE: runInit,
}

type ProjectConfig struct {
	Name        string
	Description string
	Version     string
	Author      string
	Email       string
}

func runInit(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	projectName := args[0]

	// Initialize logger
	if err := logger.InitLogger("info", "json"); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	log := logger.WithContext(ctx).With(
		zap.String("project_name", projectName),
		zap.String("operation", "init"),
	)

	// Log operation start
	defer logger.LogOperation(ctx, "project_init",
		zap.String("project_name", projectName))()

	// Validate project name
	if err := validateProjectName(projectName); err != nil {
		log.Error("project name validation failed", zap.Error(err))
		return fmt.Errorf("invalid project name: %w", err)
	}

	// Create project directory with secure permissions
	projectPath := filepath.Join(".", projectName)
	if err := os.MkdirAll(projectPath, 0750); err != nil {
		log.Error("failed to create project directory", zap.Error(err))
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Initialize project configuration
	config := ProjectConfig{
		Name:        projectName,
		Description: fmt.Sprintf("A CMP application: %s", projectName),
		Version:     "0.1.0",
		Author:      "CMP Developer",
		Email:       "developer@example.com",
	}

	// Create directory structure with secure permissions
	dirs := []string{
		"contexts",
		"memory/documents",
		"memory/embeddings",
		"prompts",
		"tools",
		"tests",
		"config/environments",
		"config/providers",
		"docs",
		"scripts",
	}

	for _, dir := range dirs {
		dirPath := filepath.Join(projectPath, dir)
		if err := os.MkdirAll(dirPath, 0750); err != nil {
			log.Error("failed to create directory", zap.String("directory", dir), zap.Error(err))
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create configuration files
	if err := createConfigFiles(projectPath, config); err != nil {
		log.Error("failed to create config files", zap.Error(err))
		return fmt.Errorf("failed to create config files: %w", err)
	}

	// Create basic templates
	if err := createBasicTemplates(projectPath, config); err != nil {
		log.Error("failed to create templates", zap.Error(err))
		return fmt.Errorf("failed to create templates: %w", err)
	}

	// Create test files
	if err := createTestFiles(projectPath, config); err != nil {
		log.Error("failed to create test files", zap.Error(err))
		return fmt.Errorf("failed to create test files: %w", err)
	}

	// Create documentation
	if err := createDocumentation(projectPath, config); err != nil {
		log.Error("failed to create documentation", zap.Error(err))
		return fmt.Errorf("failed to create documentation: %w", err)
	}

	// Create scripts
	if err := createScripts(projectPath, config); err != nil {
		log.Error("failed to create scripts", zap.Error(err))
		return fmt.Errorf("failed to create scripts: %w", err)
	}

	log.Info("project created successfully",
		zap.String("project_path", projectPath),
		zap.String("project_name", projectName))

	fmt.Printf(" Successfully created CMP project: %s\n", projectName)
	fmt.Printf(" Project structure created at: %s\n", projectPath)
	fmt.Printf("\n Next steps:\n")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Printf("  ctx generate rag MyFirstRAG --db=sqlite --embeddings=openai\n")
	fmt.Printf("  ctx test\n")

	return nil
}

func validateProjectName(name string) error {
	// Use the centralized validation function
	if err := config.ValidateProjectName(name); err != nil {
		return err
	}

	// Check if directory already exists
	if _, err := os.Stat(name); err == nil {
		return fmt.Errorf("directory '%s' already exists", name)
	}

	return nil
}

func createConfigFiles(projectPath string, config ProjectConfig) error {
	// Create context.lock.json
	lockContent := `{
  "version": "1.0.0",
  "project": "` + config.Name + `",
  "created": "` + getCurrentTimestamp() + `",
  "contexts": {},
  "memory": {},
  "prompts": {},
  "tools": {}
}`

	lockPath := filepath.Join(projectPath, "context.lock.json")
	if err := os.WriteFile(lockPath, []byte(lockContent), 0644); err != nil {
		return err
	}

	// Create development.yaml
	devConfig := `environment: development

# Database Configuration
database:
  provider: sqlite
  path: ./data/development.db
  pool_size: 10

# AI Provider Configuration
providers:
  openai:
    api_key: ${OPENAI_API_KEY}
    model: gpt-4o-mini
    temperature: 0.1
    max_tokens: 1000
  
  anthropic:
    api_key: ${ANTHROPIC_API_KEY}
    model: claude-3-sonnet-20240229
    temperature: 0.1
    max_tokens: 1000

# Embeddings Configuration
embeddings:
  provider: openai
  model: text-embedding-3-small
  dimensions: 1536

# Vector Database Configuration
vector_db:
  provider: chroma
  path: ./data/embeddings
  collection_name: development_knowledge

# Testing Configuration
testing:
  drift_threshold: 0.85
  similarity_threshold: 0.8
  max_test_duration: 300s

# Logging Configuration
logging:
  level: debug
  format: json
  output: stdout

# Development Features
features:
  hot_reload: true
  debug_mode: true
  mock_providers: false
  enable_telemetry: false`

	devConfigPath := filepath.Join(projectPath, "config", "environments", "development.yaml")
	if err := os.WriteFile(devConfigPath, []byte(devConfig), 0644); err != nil {
		return err
	}

	// Create production.yaml
	prodConfig := `environment: production

# Database Configuration
database:
  provider: postgresql
  host: ${DB_HOST}
  port: ${DB_PORT}
  name: ${DB_NAME}
  user: ${DB_USER}
  password: ${DB_PASSWORD}
  pool_size: 20

# AI Provider Configuration
providers:
  openai:
    api_key: ${OPENAI_API_KEY}
    model: gpt-4
    temperature: 0.1
    max_tokens: 1000

# Embeddings Configuration
embeddings:
  provider: openai
  model: text-embedding-3-small
  dimensions: 1536

# Vector Database Configuration
vector_db:
  provider: pinecone
  api_key: ${PINECONE_API_KEY}
  environment: ${PINECONE_ENVIRONMENT}
  index_name: ${PINECONE_INDEX}

# Testing Configuration
testing:
  drift_threshold: 0.9
  similarity_threshold: 0.85
  max_test_duration: 60s

# Logging Configuration
logging:
  level: info
  format: json
  output: file
  file_path: ./logs/app.log

# Production Features
features:
  hot_reload: false
  debug_mode: false
  mock_providers: false
  enable_telemetry: true`

	prodConfigPath := filepath.Join(projectPath, "config", "environments", "production.yaml")
	return os.WriteFile(prodConfigPath, []byte(prodConfig), 0644)
}

func createBasicTemplates(projectPath string, config ProjectConfig) error {
	// Create default context
	defaultContext := `name: "Default Agent"
version: "1.0.0"
description: "Default agent for ` + config.Name + `"

role:
  persona: "Helpful AI assistant"
  capabilities: ["answer_questions", "process_requests"]
  limitations: ["no_harmful_content", "no_personal_data"]

tools:
  - name: "knowledge_search"
    uri: "mcp://search.knowledge_base"
    description: "Search project knowledge base"

guardrails:
  tone: "helpful"
  format: "text"
  max_tokens: 500
  temperature: 0.1
  
memory:
  episodic: true
  max_history: 5
  privacy: "user_isolated"

testing:
  drift_threshold: 0.85
  business_rules:
    - "always_be_helpful"
    - "no_harmful_content"`

	contextPath := filepath.Join(projectPath, "contexts", "default_agent.ctx")
	if err := os.WriteFile(contextPath, []byte(defaultContext), 0644); err != nil {
		return err
	}

	// Create default prompt template
	defaultPrompt := `# Response Template

Based on the user query: {{ user_query }}

## Context Information
{{#if conversation_history}}
Previous conversation: {{ conversation_history }}
{{/if}}

## Knowledge Base Results
{{#each knowledge_results}}
- **Source**: {{ source }}
- **Content**: {{ content }}
- **Relevance**: {{ relevance_score }}
{{/each}}

## Response Guidelines
- **Tone**: Helpful and informative
- **Format**: Clear, structured response
- **Max Tokens**: 500
- **Include**: Relevant information, next steps if needed

## Response

{{ response_text }}

{{#if confidence_score}}
**Confidence**: {{ confidence_score }}
{{/if}}`

	promptPath := filepath.Join(projectPath, "prompts", "default_response.md")
	return os.WriteFile(promptPath, []byte(defaultPrompt), 0644)
}

func createTestFiles(projectPath string, config ProjectConfig) error {
	// Create drift detection test
	driftTest := `import pytest
from contexis.testing.drift import DriftDetector

def test_response_drift():
    """Test for response drift detection"""
    detector = DriftDetector(
        threshold=0.85,
        test_queries=[
            "What is this project about?",
            "How can I get help?",
            "What are the main features?"
        ]
    )
    
    results = detector.run_tests()
    assert results.all_passed, "Drift detected in responses"

def test_context_consistency():
    """Test that context remains consistent"""
    # Add your context consistency tests here
    pass`

	driftPath := filepath.Join(projectPath, "tests", "test_drift.py")
	if err := os.WriteFile(driftPath, []byte(driftTest), 0644); err != nil {
		return err
	}

	// Create correctness test
	correctnessTest := `import pytest
from contexis.testing.correctness import CorrectnessTester

def test_business_logic():
    """Test business logic compliance"""
    tester = CorrectnessTester(
        rules_file="./tests/business_rules.yaml"
    )
    
    test_cases = [
        {
            "input": "What is the return policy?",
            "expected": "should_mention_30_days",
            "forbidden": ["refund_immediately"]
        }
    ]
    
    results = tester.run_tests(test_cases)
    assert results.all_passed, "Business logic violations detected"

def test_response_format():
    """Test response format compliance"""
    # Add your format validation tests here
    pass`

	correctnessPath := filepath.Join(projectPath, "tests", "test_correctness.py")
	return os.WriteFile(correctnessPath, []byte(correctnessTest), 0644)
}

func createDocumentation(projectPath string, config ProjectConfig) error {
	// Create README
	readmeContent := "# " + config.Name + `

A CMP (Context-Memory-Prompt) application built with Contexis.

## Quick Start

` + "```" + `bash
# Install dependencies
pip install -r requirements.txt

# Set up environment
cp .env.example .env
# Edit .env with your API keys

# Run tests
ctx test

# Start development
ctx dev
` + "```" + `

## Project Structure

- ` + "`contexts/`" + ` - Agent definitions and behaviors
- ` + "`memory/`" + ` - Knowledge base and embeddings
- ` + "`prompts/`" + ` - Response templates
- ` + "`tools/`" + ` - Custom integrations
- ` + "`tests/`" + ` - Drift detection and validation
- ` + "`config/`" + ` - Environment configuration

## Configuration

Edit ` + "`config/environments/development.yaml`" + ` to customize:
- AI provider settings
- Vector database configuration
- Testing thresholds
- Logging preferences

## Testing

` + "```" + `bash
# Run all tests
ctx test

# Run drift detection only
ctx test --drift

# Run correctness tests
ctx test --correctness
` + "```" + `

## Deployment

` + "```" + `bash
# Build for production
ctx build --environment=production

# Deploy
ctx deploy --target=docker
` + "```" + `
`

	readmePath := filepath.Join(projectPath, "README.md")
	return os.WriteFile(readmePath, []byte(readmeContent), 0644)
}

func createScripts(projectPath string, config ProjectConfig) error {
	// Create requirements.txt
	requirements := `# Core dependencies
contexis>=0.1.0
pydantic>=2.0.0
pyyaml>=6.0
click>=8.0.0
rich>=13.0.0

# AI providers
openai>=1.0.0
anthropic>=0.7.0

# Vector databases
chromadb>=0.4.0
sentence-transformers>=2.2.0

# Testing
pytest>=7.0.0
pytest-asyncio>=0.21.0

# Development
black>=23.0.0
isort>=5.12.0
flake8>=6.0.0`

	requirementsPath := filepath.Join(projectPath, "requirements.txt")
	if err := os.WriteFile(requirementsPath, []byte(requirements), 0644); err != nil {
		return err
	}

	// Create .env.example
	envExample := `# AI Provider Keys
OPENAI_API_KEY=your_openai_api_key_here
ANTHROPIC_API_KEY=your_anthropic_api_key_here

# Database Configuration (for production)
DB_HOST=localhost
DB_PORT=5432
DB_NAME=cmp_app
DB_USER=cmp_user
DB_PASSWORD=your_password_here

# Vector Database (for production)
PINECONE_API_KEY=your_pinecone_api_key_here
PINECONE_ENVIRONMENT=us-west1-gcp
PINECONE_INDEX=your_index_name

# Application Settings
LOG_LEVEL=debug
ENVIRONMENT=development`

	envExamplePath := filepath.Join(projectPath, ".env.example")
	return os.WriteFile(envExamplePath, []byte(envExample), 0644)
}

func getCurrentTimestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}
