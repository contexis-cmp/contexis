// Package config provides configuration management for the Contexis CLI.
//
// The config package handles loading, validation, and management of configuration
// files for CMP contexts, environments, and application settings. It provides
// structured configuration with security validation and sensitive data handling.
//
// Key Features:
//   - YAML/JSON configuration loading and validation
//   - Security validation for project names and file paths
//   - Sensitive data redaction for logging
//   - Environment-specific configuration management
//   - Schema validation for configuration files
//
// Example Usage:
//
//	// Load a context configuration
//	ctxConfig, err := config.LoadContext("path/to/context.yaml")
//
//	// Load environment configuration
//	envConfig, err := config.LoadEnvironment("path/to/environment.yaml")
//
//	// Validate project name
//	err = config.ValidateProjectName("my-project")
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// ContextConfig represents a CMP context configuration.
// It defines the structure and behavior of an AI agent or component
// within the CMP framework, including its role, tools, and constraints.
type ContextConfig struct {
	// Name is the unique identifier for the context
	Name        string `json:"name" yaml:"name"`
	
	// Version follows semantic versioning (e.g., "1.0.0")
	Version     string `json:"version" yaml:"version"`
	
	// Description provides human-readable information about the context
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// Role defines the agent's persona, capabilities, and limitations
	Role       RoleConfig    `json:"role" yaml:"role"`
	
	// Tools defines the available tools and their configurations
	Tools      []ToolConfig  `json:"tools,omitempty" yaml:"tools,omitempty"`
	
	// Guardrails defines safety constraints and behavior limits
	Guardrails GuardrailConfig `json:"guardrails,omitempty" yaml:"guardrails,omitempty"`
	
	// Memory defines memory settings and configuration
	Memory     MemoryConfig  `json:"memory,omitempty" yaml:"memory,omitempty"`
	
	// Testing defines testing parameters and validation rules
	Testing    TestingConfig `json:"testing,omitempty" yaml:"testing,omitempty"`

	// CreatedAt is the timestamp when the context was created
	CreatedAt time.Time         `json:"created_at" yaml:"created_at"`
	
	// UpdatedAt is the timestamp when the context was last modified
	UpdatedAt time.Time         `json:"updated_at" yaml:"updated_at"`
	
	// Metadata contains additional key-value pairs for extensibility
	Metadata  map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// RoleConfig defines the agent's role and capabilities.
// It specifies the persona, capabilities, and limitations that guide
// the agent's behavior and responses.
type RoleConfig struct {
	// Persona defines the agent's character and personality
	Persona      string   `json:"persona" yaml:"persona"`
	
	// Capabilities lists the agent's abilities and functions
	Capabilities []string `json:"capabilities" yaml:"capabilities"`
	
	// Limitations lists the agent's restrictions and boundaries
	Limitations  []string `json:"limitations" yaml:"limitations"`
}

// ToolConfig defines an MCP (Model Context Protocol) tool.
// It specifies the configuration for external tools that the agent
// can use to extend its capabilities.
type ToolConfig struct {
	// Name is the unique identifier for the tool
	Name        string `json:"name" yaml:"name"`
	
	// URI is the MCP URI for the tool
	URI         string `json:"uri" yaml:"uri"`
	
	// Description provides human-readable information about the tool
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

// GuardrailConfig defines safety constraints for the agent.
// It provides mechanisms to ensure the agent operates within
// defined boundaries and safety parameters.
type GuardrailConfig struct {
	// Tone specifies the desired communication style
	Tone        string  `json:"tone,omitempty" yaml:"tone,omitempty"`
	
	// Format specifies the output format requirements
	Format      string  `json:"format,omitempty" yaml:"format,omitempty"`
	
	// MaxTokens limits the maximum response length
	MaxTokens   int     `json:"max_tokens,omitempty" yaml:"max_tokens,omitempty"`
	
	// Temperature controls response randomness (0.0 to 1.0)
	Temperature float64 `json:"temperature,omitempty" yaml:"temperature,omitempty"`
}

// MemoryConfig defines memory settings for the agent.
// It controls how the agent stores and retrieves information
// across conversations and interactions.
type MemoryConfig struct {
	// Episodic enables episodic memory for conversation history
	Episodic   bool   `json:"episodic" yaml:"episodic"`
	
	// MaxHistory limits the number of conversation turns to remember
	MaxHistory int    `json:"max_history" yaml:"max_history"`
	
	// Privacy specifies the privacy level: "user_isolated", "shared", "public"
	Privacy    string `json:"privacy" yaml:"privacy"`
}

// TestingConfig defines testing parameters for the context.
// It specifies how the context should be tested and validated
// to ensure proper functionality and behavior.
type TestingConfig struct {
	// DriftThreshold defines the similarity threshold for drift detection
	DriftThreshold float64  `json:"drift_threshold" yaml:"drift_threshold"`
	
	// BusinessRules lists business logic rules for validation
	BusinessRules  []string `json:"business_rules" yaml:"business_rules"`
}

// EnvironmentConfig represents environment-specific settings.
// It defines configuration that varies between development,
// staging, and production environments.
type EnvironmentConfig struct {
	// Environment specifies the environment name
	Environment string `json:"environment" yaml:"environment"`
	
	// Database defines database connection settings
	Database DatabaseConfig `json:"database" yaml:"database"`
	
	// Providers defines AI model provider configurations
	Providers map[string]ProviderConfig `json:"providers" yaml:"providers"`
	
	// Embeddings defines embedding model configurations
	Embeddings EmbeddingConfig `json:"embeddings" yaml:"embeddings"`
	
	// VectorDB defines vector database configurations
	VectorDB VectorDBConfig `json:"vector_db" yaml:"vector_db"`
	
	// Logging defines logging configuration
	Logging LoggingConfig `json:"logging" yaml:"logging"`
	
	// Features defines feature flags and toggles
	Features FeatureConfig `json:"features" yaml:"features"`
}

// DatabaseConfig defines database settings.
// It specifies connection parameters and configuration for
// various database types used by the application.
type DatabaseConfig struct {
	// Type specifies the database type (e.g., "sqlite", "postgres")
	Type     string            `json:"type" yaml:"type"`
	
	// URL is the database connection string
	URL      string            `json:"url" yaml:"url"`
	
	// Options contains additional database-specific options
	Options  map[string]string `json:"options,omitempty" yaml:"options,omitempty"`
	
	// MaxConnections limits the number of database connections
	MaxConnections int `json:"max_connections,omitempty" yaml:"max_connections,omitempty"`
}

// ProviderConfig defines AI model provider settings.
// It specifies configuration for external AI model providers
// like OpenAI, Anthropic, or Hugging Face.
type ProviderConfig struct {
	// APIKey is the provider's API key (sensitive)
	APIKey string `json:"api_key" yaml:"api_key"`
	
	// Model specifies the model to use
	Model string `json:"model" yaml:"model"`
	
	// BaseURL is the provider's API base URL
	BaseURL string `json:"base_url,omitempty" yaml:"base_url,omitempty"`
}

// EmbeddingConfig defines embedding model settings.
// It specifies configuration for text embedding models used
// for semantic search and vector operations.
type EmbeddingConfig struct {
	// Provider specifies the embedding provider
	Provider   string `json:"provider" yaml:"provider"`
	
	// Model specifies the embedding model to use
	Model      string `json:"model" yaml:"model"`
	
	// Dimensions specifies the embedding vector dimensions
	Dimensions int    `json:"dimensions" yaml:"dimensions"`
}

// VectorDBConfig defines vector database settings.
// It specifies configuration for vector databases used to
// store and search embeddings.
type VectorDBConfig struct {
	// Provider specifies the vector database provider
	Provider string `json:"provider" yaml:"provider"`
	
	// URL is the vector database connection string
	URL      string `json:"url" yaml:"url"`
	
	// IndexName specifies the index to use
	IndexName string `json:"index_name" yaml:"index_name"`
	
	// Options contains additional provider-specific options
	Options  map[string]string `json:"options,omitempty" yaml:"options,omitempty"`
}

// LoggingConfig defines logging settings.
// It specifies how the application should log events,
// errors, and debugging information.
type LoggingConfig struct {
	// Level specifies the logging level (debug, info, warn, error)
	Level string `json:"level" yaml:"level"`
	
	// Format specifies the log format (json, text)
	Format string `json:"format" yaml:"format"`
	
	// Output specifies the log output destination
	Output string `json:"output" yaml:"output"`
}

// FeatureConfig defines feature flags and toggles.
// It provides a way to enable or disable features
// without code changes.
type FeatureConfig struct {
	// Enabled features list
	Enabled []string `json:"enabled" yaml:"enabled"`
	
	// Disabled features list
	Disabled []string `json:"disabled" yaml:"disabled"`
	
	// Options contains feature-specific configuration
	Options map[string]interface{} `json:"options,omitempty" yaml:"options,omitempty"`
}

// LoadContext loads and validates a context configuration file.
// It reads a YAML file and unmarshals it into a ContextConfig structure,
// performing validation to ensure the configuration is valid.
//
// Parameters:
//   - path: File path to the context configuration file
//
// Returns:
//   - *ContextConfig: Loaded and validated context configuration
//   - error: Any error that occurred during loading or validation
func LoadContext(path string) (*ContextConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read context file: %w", err)
	}

	var config ContextConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse context file: %w", err)
	}

	// Validate required fields
	if config.Name == "" {
		return nil, fmt.Errorf("context name is required")
	}
	if config.Version == "" {
		return nil, fmt.Errorf("context version is required")
	}

	return &config, nil
}

// LoadEnvironment loads and validates an environment configuration file.
// It reads a YAML file and unmarshals it into an EnvironmentConfig structure,
// performing validation to ensure the configuration is valid.
//
// Parameters:
//   - path: File path to the environment configuration file
//
// Returns:
//   - *EnvironmentConfig: Loaded and validated environment configuration
//   - error: Any error that occurred during loading or validation
func LoadEnvironment(path string) (*EnvironmentConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read environment file: %w", err)
	}

	var config EnvironmentConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse environment file: %w", err)
	}

	// Validate required fields
	if config.Environment == "" {
		return nil, fmt.Errorf("environment name is required")
	}

	return &config, nil
}

// ValidateProjectName validates project name according to security rules.
// It ensures project names are safe and don't contain dangerous characters
// or patterns that could lead to security vulnerabilities.
//
// Parameters:
//   - name: Project name to validate
//
// Returns:
//   - error: Validation error if the name is invalid, nil otherwise
func ValidateProjectName(name string) error {
	if len(name) == 0 || len(name) > 50 {
		return fmt.Errorf("project name must be 1-50 characters")
	}
	
	// Only allow alphanumeric, hyphens, and underscores
	if !regexp.MustCompile(`^[a-zA-Z0-9-_]+$`).MatchString(name) {
		return fmt.Errorf("project name contains invalid characters (only a-z, A-Z, 0-9, -, _ allowed)")
	}
	
	// Prevent common dangerous names
	dangerousNames := []string{"config", "system", "admin", "root", "etc", "var", "tmp"}
	for _, dangerous := range dangerousNames {
		if strings.EqualFold(name, dangerous) {
			return fmt.Errorf("project name '%s' is not allowed", name)
		}
	}
	
	return nil
}

// SafeJoin safely joins paths to prevent directory traversal attacks.
// It ensures that the resulting path doesn't escape the base directory
// and provides protection against path traversal vulnerabilities.
//
// Parameters:
//   - base: Base directory path
//   - path: Path to join with the base
//
// Returns:
//   - string: Safely joined path
//   - error: Error if the path would escape the base directory
func SafeJoin(base, path string) (string, error) {
	full := filepath.Join(base, filepath.Clean(path))
	if !strings.HasPrefix(full, base) {
		return "", fmt.Errorf("path escapes base directory")
	}
	return full, nil
}

// RedactSensitiveData removes sensitive information from logs.
// It identifies and replaces sensitive data like API keys, passwords,
// and tokens with placeholder values to prevent accidental exposure.
//
// Parameters:
//   - data: String containing potentially sensitive data
//
// Returns:
//   - string: String with sensitive data redacted
func RedactSensitiveData(data string) string {
	// Redact API keys (common patterns)
	apiKeyPattern := regexp.MustCompile(`(api[_-]?key|token|password|secret)[\s]*[:=][\s]*["']?[a-zA-Z0-9_-]+["']?`)
	data = apiKeyPattern.ReplaceAllString(data, "$1: [REDACTED]")
	
	// Redact email addresses
	emailPattern := regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`)
	data = emailPattern.ReplaceAllString(data, "[REDACTED_EMAIL]")
	
	// Redact URLs with credentials
	urlPattern := regexp.MustCompile(`https?://[^:]+:[^@]+@[^\s]+`)
	data = urlPattern.ReplaceAllString(data, "[REDACTED_URL]")
	
	return data
}
