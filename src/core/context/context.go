package context

import (
	"encoding/json"
	"fmt"
	"time"
)

// Context represents the declarative instructions and agent roles
// in the CMP framework
type Context struct {
	Name        string            `json:"name" yaml:"name"`
	Version     string            `json:"version" yaml:"version"`
	Description string            `json:"description,omitempty" yaml:"description,omitempty"`
	
	Role        Role              `json:"role" yaml:"role"`
	Tools       []Tool            `json:"tools,omitempty" yaml:"tools,omitempty"`
	Guardrails  Guardrails        `json:"guardrails,omitempty" yaml:"guardrails,omitempty"`
	Memory      MemoryConfig      `json:"memory,omitempty" yaml:"memory,omitempty"`
	Testing     TestingConfig     `json:"testing,omitempty" yaml:"testing,omitempty"`
	
	CreatedAt   time.Time         `json:"created_at" yaml:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" yaml:"updated_at"`
	Metadata    map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// Role defines the agent's persona and capabilities
type Role struct {
	Persona     string   `json:"persona" yaml:"persona"`
	Capabilities []string `json:"capabilities" yaml:"capabilities"`
	Limitations []string `json:"limitations" yaml:"limitations"`
}

// Tool represents an external function or integration
type Tool struct {
	Name        string `json:"name" yaml:"name"`
	URI         string `json:"uri" yaml:"uri"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

// Guardrails define behavioral constraints
type Guardrails struct {
	Tone       string `json:"tone,omitempty" yaml:"tone,omitempty"`
	Format     string `json:"format,omitempty" yaml:"format,omitempty"`
	MaxTokens  int    `json:"max_tokens,omitempty" yaml:"max_tokens,omitempty"`
	Temperature float64 `json:"temperature,omitempty" yaml:"temperature,omitempty"`
}

// MemoryConfig defines memory behavior
type MemoryConfig struct {
	Episodic   bool   `json:"episodic" yaml:"episodic"`
	MaxHistory int    `json:"max_history" yaml:"max_history"`
	Privacy    string `json:"privacy" yaml:"privacy"`
}

// TestingConfig defines testing parameters
type TestingConfig struct {
	DriftThreshold float64  `json:"drift_threshold" yaml:"drift_threshold"`
	BusinessRules  []string `json:"business_rules" yaml:"business_rules"`
}

// New creates a new Context with default values
func New(name, version string) *Context {
	return &Context{
		Name:      name,
		Version:   version,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Metadata:  make(map[string]string),
		Testing: TestingConfig{
			DriftThreshold: 0.85,
		},
	}
}

// Validate ensures the context is properly configured
func (c *Context) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("context name is required")
	}
	if c.Version == "" {
		return fmt.Errorf("context version is required")
	}
	if c.Role.Persona == "" {
		return fmt.Errorf("role persona is required")
	}
	return nil
}

// ToJSON converts the context to JSON
func (c *Context) ToJSON() ([]byte, error) {
	return json.MarshalIndent(c, "", "  ")
}

// FromJSON creates a context from JSON
func FromJSON(data []byte) (*Context, error) {
	var ctx Context
	if err := json.Unmarshal(data, &ctx); err != nil {
		return nil, err
	}
	return &ctx, nil
}

// GetSHA returns a content-based hash for versioning
func (c *Context) GetSHA() (string, error) {
	data, err := c.ToJSON()
	if err != nil {
		return "", err
	}
	// TODO: Implement proper SHA256 hashing
	return fmt.Sprintf("sha256:%x", data), nil
} 