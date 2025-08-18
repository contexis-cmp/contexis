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

func init() {}

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

	// Use colored logger (already initialized in main.go)

	// Log operation start with colored logging
	operationName := "project_init"
	done := logger.LogOperationColored(ctx, operationName,
		zap.String("project_name", projectName))

	// Validate project name
	if err := validateProjectName(projectName); err != nil {
		logger.LogErrorColored(ctx, "project name validation failed", err)
		return fmt.Errorf("invalid project name: %w", err)
	}

	// Create project directory with secure permissions
	projectPath := filepath.Join(".", projectName)
	if err := os.MkdirAll(projectPath, 0750); err != nil {
		logger.LogErrorColored(ctx, "failed to create project directory", err)
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Initialize project configuration
	config := ProjectConfig{
		Name:        projectName,
		Description: fmt.Sprintf("A CMP application: %s", projectName),
		Version:     "0.1.14",
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

	logger.LogInfo(ctx, "Creating project directory structure")
	for _, dir := range dirs {
		dirPath := filepath.Join(projectPath, dir)
		if err := os.MkdirAll(dirPath, 0750); err != nil {
			logger.LogErrorColored(ctx, "failed to create directory", err, zap.String("directory", dir))
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		logger.LogDebugWithContext(ctx, "Created directory", zap.String("path", dir))
	}

	// Create configuration files
	logger.LogInfo(ctx, "Creating configuration files")
	if err := createConfigFiles(projectPath, config); err != nil {
		logger.LogErrorColored(ctx, "failed to create config files", err)
		return fmt.Errorf("failed to create config files: %w", err)
	}

	// Create basic templates
	logger.LogInfo(ctx, "Creating basic templates")
	if err := createBasicTemplates(projectPath, config); err != nil {
		logger.LogErrorColored(ctx, "failed to create templates", err)
		return fmt.Errorf("failed to create templates: %w", err)
	}

	// Create test files
	logger.LogInfo(ctx, "Creating test files")
	if err := createTestFiles(projectPath, config); err != nil {
		logger.LogErrorColored(ctx, "failed to create test files", err)
		return fmt.Errorf("failed to create test files: %w", err)
	}

	// Create documentation
	logger.LogInfo(ctx, "Creating documentation")
	if err := createDocumentation(projectPath, config); err != nil {
		logger.LogErrorColored(ctx, "failed to create documentation", err)
		return fmt.Errorf("failed to create documentation: %w", err)
	}

	// Create scripts
	logger.LogInfo(ctx, "Creating scripts and requirements")
	if err := createScripts(projectPath, config); err != nil {
		logger.LogErrorColored(ctx, "failed to create scripts", err)
		return fmt.Errorf("failed to create scripts: %w", err)
	}

	// Complete the operation
	done()

	// Show project structure
	showProjectStructure(projectPath, projectName)

	// Show development flow
	showDevelopmentFlow(projectName)

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

	// Create development.yaml (local-first defaults with production hints)
	devConfig := `environment: development

# Database Configuration
database:
  provider: sqlite
  path: ./data/development.db
  pool_size: 10

# AI Provider Configuration (Local by default)
# To switch to production providers (e.g., OpenAI), replace this section accordingly.
providers:
  local:
    # Local text generation model
    model: microsoft/Phi-3-mini-4k-instruct   # For production, use: openai.gpt-4o-mini or similar
    temperature: 0.1
    max_tokens: 1000
    device: auto            # auto, cpu, cuda
    load_in_8bit: false     # set true to save memory
    load_in_4bit: false     # set true to save memory

# Embeddings Configuration (Local Sentence Transformers)
embeddings:
  provider: sentence-transformers   # For production, set provider: openai
  model: all-MiniLM-L6-v2
  dimensions: 384
  device: auto

# Vector Database Configuration (Local Chroma with SQLite)
vector_db:
  provider: chroma
  path: ./data/embeddings
  collection_name: development_knowledge
  persist_directory: ./data/chroma
  anonymized_telemetry: false

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
  enable_telemetry: false
  local_models: true
  offline_mode: true

# Model Download Configuration
model_cache:
  directory: ./data/models
  auto_download: true
  verify_checksums: true

# ---
# Production Hints:
# To use OpenAI in production, configure like:
# providers:
#   openai:
#     api_key: ${OPENAI_API_KEY}
#     model: gpt-4o-mini
# embeddings:
#   provider: openai
#   model: text-embedding-3-small
`

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
	defaultPrompt := `# Default Response

Provide a helpful and concise answer using available knowledge.
`

	promptPath := filepath.Join(projectPath, "prompts", "default_response.md")
	return os.WriteFile(promptPath, []byte(defaultPrompt), 0644)
}

func createTestFiles(projectPath string, config ProjectConfig) error {
	// Create test subdirectories
	testDirs := []string{"unit", "integration", "e2e", "coverage", "reports", "fixtures"}
	for _, dir := range testDirs {
		dirPath := filepath.Join(projectPath, "tests", dir)
		if err := os.MkdirAll(dirPath, 0750); err != nil {
			return err
		}
	}

	// Create test configuration
	testConfig := `# CMP Framework Test Configuration

# Test Suite Configuration
test_suites:
  unit:
    enabled: true
    timeout: 30s
    parallel: true
    coverage_threshold: 80
    
  integration:
    enabled: true
    timeout: 120s
    parallel: false
    coverage_threshold: 70
    
  e2e:
    enabled: true
    timeout: 300s
    parallel: false
    coverage_threshold: 60
    
  performance:
    enabled: false
    timeout: 600s
    parallel: false
    
  security:
    enabled: true
    timeout: 60s
    parallel: false

# Test Categories
test_categories:
  - name: "agent_generator"
    description: "Tests for agent generation functionality"
    suites: ["unit", "integration"]
    
  - name: "cli_commands"
    description: "Tests for CLI command functionality"
    suites: ["unit", "integration"]
    
  - name: "rag_generator"
    description: "Tests for RAG generation functionality"
    suites: ["unit", "integration"]
    
  - name: "workflow_generator"
    description: "Tests for workflow generation functionality"
    suites: ["unit", "integration"]
    
  - name: "core_components"
    description: "Tests for core CMP components"
    suites: ["unit", "integration"]
    
  - name: "security"
    description: "Security and privacy tests"
    suites: ["security"]
    
  - name: "performance"
    description: "Performance and load tests"
    suites: ["performance", "e2e"]

# Test Data Configuration
test_data:
  fixtures_dir: "tests/unit/fixtures"
  temp_dir: "tests/temp"
  cleanup_after: true
  
# Coverage Configuration
coverage:
  enabled: true
  output_format: "html"
  output_dir: "tests/coverage"
  exclude_patterns:
    - "*/vendor/*"
    - "*/testdata/*"
    - "*/mocks/*"
    - "*/examples/*"

# Reporting Configuration
reporting:
  enabled: true
  output_format: ["text", "json", "html"]
  output_dir: "tests/reports"
  include_failures: true
  include_skipped: true

# Environment Configuration
environments:
  test:
    variables:
      CMP_ENV: "test"
      CMP_LOG_LEVEL: "debug"
      CMP_TEMP_DIR: "tests/temp"
      
  integration:
    variables:
      CMP_ENV: "integration"
      CMP_LOG_LEVEL: "info"
      CMP_TEMP_DIR: "tests/temp"

# Test Timeouts
timeouts:
  short: 5s
  medium: 30s
  long: 120s
  very_long: 300s

# Test Retries
retries:
  enabled: true
  max_attempts: 3
  backoff_multiplier: 2.0
  initial_delay: 1s

# Parallel Execution
parallel:
  enabled: true
  max_workers: 4
  test_timeout: 30s`

	testConfigPath := filepath.Join(projectPath, "tests", "test_config.yaml")
	if err := os.WriteFile(testConfigPath, []byte(testConfig), 0644); err != nil {
		return err
	}

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
	// Create requirements.txt (local-first defaults; production providers can be added as needed)
	requirements := `# Local-first development stack (no external APIs required)

# Local AI Models
transformers>=4.35.0
torch>=2.0.0
accelerate>=0.20.0
safetensors>=0.4.0
tokenizers>=0.15.0

# Vector Database (SQLite backend)
chromadb>=0.4.0
sentence-transformers>=2.2.0

# Local Embeddings
numpy>=1.24.0
scikit-learn>=1.3.0

# Core utilities
pydantic>=2.4.0
pyyaml>=6.0
click>=8.0.0
rich>=13.0.0
structlog>=23.0.0

# Development tools
black>=23.0.0
isort>=5.12.0
flake8>=6.0.0
pytest>=7.0.0
pytest-asyncio>=0.21.0

# To switch to production providers, add e.g.:
# openai>=1.0.0
# anthropic>=0.7.0
`

	requirementsPath := filepath.Join(projectPath, "requirements.txt")
	if err := os.WriteFile(requirementsPath, []byte(requirements), 0644); err != nil {
		return err
	}

	// Create .env.example with comprehensive variables (commented where appropriate)
	envExample := `# Contexis Environment (.env)
# Local-first defaults (no external API keys needed). Uncomment to override.

# --- General ---
# CMP_ENV=development
# CMP_PROJECT_ROOT=./
# CMP_LOG_LEVEL=debug
# CMP_LOG_FORMAT=json

# --- Local-first Provider ---
# CMP_LOCAL_MODELS=true
# CMP_OFFLINE_MODE=true
# CMP_PYTHON_BIN=.venv/bin/python
# CMP_LOCAL_TIMEOUT_SECONDS=600
# CMP_LOCAL_MODEL_ID=microsoft/Phi-3-mini-4k-instruct
# CMP_MODEL_CACHE_DIR=./data/models

# --- Server ---
# CMP_AUTH_ENABLED=false
# CMP_PI_ENFORCEMENT=false
# CMP_REQUIRE_CITATION=false
# CMP_TENANT_ID=

# --- Memory / Vector DB ---
# CMP_DB_PROVIDER=sqlite
# CMP_DB_PATH=./data/development/development.db
# CMP_VECTOR_DB_PROVIDER=chroma
# CMP_VECTOR_DB_PATH=./data/embeddings
# CMP_CHROMA_PERSIST_DIR=./data/chroma

# --- Security / Policies ---
# CMP_OOB_REQUIRED_ACTIONS=delete_user,wire_transfer
# CMP_PII_MODE=redact   # redact|block|allow
# CMP_EPISODIC_KEY=changeme
# CMP_API_KEYS=       # comma-separated list of apiKeyId:secret
# CMP_API_TOKENS=     # comma-separated list of tokenId:secret

# --- Hugging Face (optional) ---
# HF_TOKEN=
# HF_MODEL_ID=
# HF_ENDPOINT=

# --- Production: OpenAI / Anthropic (optional) ---
# OPENAI_API_KEY=
# ANTHROPIC_API_KEY=

# --- Integrations (optional) ---
# PINECONE_API_KEY=
# PINECONE_ENVIRONMENT=
# PINECONE_INDEX=
`

	envExamplePath := filepath.Join(projectPath, ".env.example")
	if err := os.WriteFile(envExamplePath, []byte(envExample), 0644); err != nil {
		return err
	}

	return nil
}

func getCurrentTimestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// showProjectStructure displays the created project structure
func showProjectStructure(projectPath, projectName string) {
	fmt.Printf("\n")
	logger.LogSuccess(context.Background(), "Project created successfully",
		zap.String("project_name", projectName),
		zap.String("project_path", projectPath))

	fmt.Printf("\n📁 Project Structure:\n")
	fmt.Printf("  %s/\n", projectName)
	fmt.Printf("  ├── 📄 README.md\n")
	fmt.Printf("  ├── 📄 requirements.txt\n")
	fmt.Printf("  ├── 📄 .env.example\n")
	fmt.Printf("  ├── 📄 context.lock.json\n")
	fmt.Printf("  ├── 📁 contexts/\n")
	fmt.Printf("  │   └── 📄 default_agent.ctx\n")
	fmt.Printf("  ├── 📁 memory/\n")
	fmt.Printf("  │   ├── 📁 documents/\n")
	fmt.Printf("  │   └── 📁 embeddings/\n")
	fmt.Printf("  ├── 📁 prompts/\n")
	fmt.Printf("  │   └── 📄 default_response.md\n")
	fmt.Printf("  ├── 📁 tools/\n")
	fmt.Printf("  ├── 📁 tests/\n")
	fmt.Printf("  │   ├── 📄 test_drift.py\n")
	fmt.Printf("  │   └── 📄 test_correctness.py\n")
	fmt.Printf("  ├── 📁 config/\n")
	fmt.Printf("  │   └── 📁 environments/\n")
	fmt.Printf("  │       ├── 📄 development.yaml\n")
	fmt.Printf("  │       └── 📄 production.yaml\n")
	fmt.Printf("  ├── 📁 docs/\n")
	fmt.Printf("  └── 📁 scripts/\n")
}

// showDevelopmentFlow displays the development workflow
func showDevelopmentFlow(projectName string) {
	fmt.Printf("\n🚀 Development Flow (Local-first):\n")
	fmt.Printf("\n1️⃣  Navigate to your project:\n")
	fmt.Printf("   cd %s\n", projectName)

	fmt.Printf("\n2️⃣  Set up your environment:\n")
	fmt.Printf("   cp .env.example .env\n")
	fmt.Printf("   # Local defaults enabled. No API keys needed.\n")

	fmt.Printf("\n3️⃣  Install dependencies:\n")
	fmt.Printf("   pip install -r requirements.txt\n")
	fmt.Printf("   # This will download local models as needed (Phi-3.5-Mini, Sentence Transformers).\n")

	fmt.Printf("\n4️⃣  Create your first RAG system:\n")
	fmt.Printf("   ctx generate rag MyFirstRAG --db=sqlite --embeddings=sentence-transformers\n")

	fmt.Printf("\n5️⃣  Add knowledge to your system:\n")
	fmt.Printf("   echo 'Your company policies here...' > memory/documents/policies.txt\n")
	fmt.Printf("   ctx memory ingest --provider=sqlite --component=MyFirstRAG --input=memory/documents/policies.txt\n")

	fmt.Printf("\n6️⃣  Test your system:\n")
	fmt.Printf("   ctx test\n")
	fmt.Printf("   ctx run MyFirstRAG \"What are your policies?\"\n")

	fmt.Printf("\n7️⃣  Start development server:\n")
	fmt.Printf("   ctx serve --addr :8000\n")

	fmt.Printf("\n💡 Switching to production providers later is easy: update config/environments/development.yaml to set provider 'openai' and export OPENAI_API_KEY.\n")

	fmt.Printf("\n📚 Next Steps:\n")
	fmt.Printf("   • Customize contexts/ in contexts/\n")
	fmt.Printf("   • Add prompts in prompts/\n")
	fmt.Printf("   • Create tools in tools/\n")
	fmt.Printf("   • Configure AI providers in config/environments/\n")
	fmt.Printf("   • Run tests with ctx test\n")
	fmt.Printf("   • Deploy with ctx deploy\n")

	fmt.Printf("\n🎉 Happy building! Check out docs/ for more information.\n")
}
