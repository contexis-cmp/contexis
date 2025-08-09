package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
	"github.com/go-playground/validator/v10"
)

// ContextConfig represents a CMP context configuration
type ContextConfig struct {
	Name        string            `yaml:"name" validate:"required"`
	Version     string            `yaml:"version" validate:"required,semver"`
	Description string            `yaml:"description"`
	Role        RoleConfig        `yaml:"role" validate:"required"`
	Tools       []ToolConfig      `yaml:"tools" validate:"dive"`
	Guardrails  GuardrailConfig   `yaml:"guardrails"`
	Memory      MemoryConfig      `yaml:"memory"`
	Testing     TestingConfig     `yaml:"testing"`
}

// RoleConfig defines the agent's role and capabilities
type RoleConfig struct {
	Persona     string   `yaml:"persona" validate:"required"`
	Capabilities []string `yaml:"capabilities" validate:"required,min=1"`
	Limitations []string `yaml:"limitations"`
}

// ToolConfig defines an MCP tool
type ToolConfig struct {
	Name        string `yaml:"name" validate:"required"`
	URI         string `yaml:"uri" validate:"required,startswith=mcp://"`
	Description string `yaml:"description"`
}

// GuardrailConfig defines safety constraints
type GuardrailConfig struct {
	Tone        string  `yaml:"tone" validate:"required"`
	Format      string  `yaml:"format" validate:"required"`
	MaxTokens   int     `yaml:"max_tokens" validate:"required,min=1,max=8000"`
	Temperature float64 `yaml:"temperature" validate:"min=0,max=2"`
}

// MemoryConfig defines memory settings
type MemoryConfig struct {
	Episodic   bool   `yaml:"episodic"`
	MaxHistory int    `yaml:"max_history" validate:"min=0,max=100"`
	Privacy    string `yaml:"privacy" validate:"oneof=user_isolated shared public"`
}

// TestingConfig defines testing parameters
type TestingConfig struct {
	DriftThreshold float64  `yaml:"drift_threshold" validate:"min=0,max=1"`
	BusinessRules  []string `yaml:"business_rules"`
}

// EnvironmentConfig represents environment-specific settings
type EnvironmentConfig struct {
	Environment string                 `yaml:"environment" validate:"required"`
	Database    DatabaseConfig         `yaml:"database" validate:"required"`
	Providers   map[string]ProviderConfig `yaml:"providers"`
	Embeddings  EmbeddingConfig        `yaml:"embeddings" validate:"required"`
	VectorDB    VectorDBConfig         `yaml:"vector_db" validate:"required"`
	Testing     TestingConfig          `yaml:"testing"`
	Logging     LoggingConfig          `yaml:"logging"`
	Features    FeatureConfig          `yaml:"features"`
}

// DatabaseConfig defines database settings
type DatabaseConfig struct {
	Provider string `yaml:"provider" validate:"required,oneof=sqlite postgresql mysql"`
	Path     string `yaml:"path,omitempty"`
	Host     string `yaml:"host,omitempty"`
	Port     int    `yaml:"port,omitempty"`
	Name     string `yaml:"name,omitempty"`
	User     string `yaml:"user,omitempty"`
	Password string `yaml:"password,omitempty"`
	PoolSize int    `yaml:"pool_size" validate:"min=1,max=100"`
}

// ProviderConfig defines AI provider settings
type ProviderConfig struct {
	APIKey      string  `yaml:"api_key" validate:"required"`
	Model       string  `yaml:"model" validate:"required"`
	Temperature float64 `yaml:"temperature" validate:"min=0,max=2"`
	MaxTokens   int     `yaml:"max_tokens" validate:"min=1,max=8000"`
}

// EmbeddingConfig defines embedding settings
type EmbeddingConfig struct {
	Provider   string `yaml:"provider" validate:"required"`
	Model      string `yaml:"model" validate:"required"`
	Dimensions int    `yaml:"dimensions" validate:"min=1"`
}

// VectorDBConfig defines vector database settings
type VectorDBConfig struct {
	Provider       string `yaml:"provider" validate:"required,oneof=chroma pinecone weaviate"`
	Path           string `yaml:"path,omitempty"`
	CollectionName string `yaml:"collection_name,omitempty"`
	APIKey         string `yaml:"api_key,omitempty"`
	Environment    string `yaml:"environment,omitempty"`
	IndexName      string `yaml:"index_name,omitempty"`
}

// LoggingConfig defines logging settings
type LoggingConfig struct {
	Level    string `yaml:"level" validate:"oneof=debug info warn error"`
	Format   string `yaml:"format" validate:"oneof=json text"`
	Output   string `yaml:"output" validate:"oneof=stdout file"`
	FilePath string `yaml:"file_path,omitempty"`
}

// FeatureConfig defines feature flags
type FeatureConfig struct {
	HotReload      bool `yaml:"hot_reload"`
	DebugMode      bool `yaml:"debug_mode"`
	MockProviders  bool `yaml:"mock_providers"`
	EnableTelemetry bool `yaml:"enable_telemetry"`
}

// LoadContext loads and validates a context configuration file
func LoadContext(path string) (*ContextConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read context file: %w", err)
	}
	
	var config ContextConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parse context: %w", err)
	}
	
	validate := validator.New()
	if err := validate.Struct(&config); err != nil {
		return nil, fmt.Errorf("validate context: %w", err)
	}
	
	return &config, nil
}

// LoadEnvironment loads and validates an environment configuration file
func LoadEnvironment(path string) (*EnvironmentConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read environment file: %w", err)
	}
	
	var config EnvironmentConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parse environment: %w", err)
	}
	
	validate := validator.New()
	if err := validate.Struct(&config); err != nil {
		return nil, fmt.Errorf("validate environment: %w", err)
	}
	
	return &config, nil
}

// ValidateProjectName validates project name according to security rules
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

// SafeJoin safely joins paths to prevent directory traversal attacks
func SafeJoin(base, path string) (string, error) {
	full := filepath.Join(base, filepath.Clean(path))
	
	if !strings.HasPrefix(full, base) {
		return "", fmt.Errorf("path escapes base directory")
	}
	
	return full, nil
}

// RedactSensitiveData removes sensitive information from logs
func RedactSensitiveData(data map[string]interface{}) map[string]interface{} {
	sensitiveKeys := []string{"api_key", "password", "token", "secret"}
	
	redacted := make(map[string]interface{})
	for k, v := range data {
		isSensitive := false
		for _, sensitive := range sensitiveKeys {
			if strings.Contains(strings.ToLower(k), sensitive) {
				isSensitive = true
				break
			}
		}
		
		if isSensitive {
			if str, ok := v.(string); ok {
				if len(str) > 8 {
					redacted[k] = str[:4] + "..." + str[len(str)-4:]
				} else {
					redacted[k] = "***"
				}
			} else {
				redacted[k] = "***"
			}
		} else {
			redacted[k] = v
		}
	}
	
	return redacted
}
