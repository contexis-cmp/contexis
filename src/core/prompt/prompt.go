// Package prompt provides the core prompt management functionality for the CMP framework.
//
// The prompt package implements pure templates that are hydrated at runtime with
// context, memory, and user data. It provides a structured approach to prompt
// engineering with versioning, validation, and template management capabilities.
//
// Key Features:
//   - Template-based prompt generation
//   - Variable substitution and validation
//   - Content-based versioning with SHA-256 hashing
//   - JSON/YAML serialization support
//   - Built-in template examples for common use cases
//
// Example Usage:
//
//	// Create a new prompt
//	prompt := prompt.New("customer_response", "1.0.0", "Hello {{.name}}, how can I help you?")
//
//	// Render with data
//	result, err := prompt.Render(map[string]interface{}{
//		"name": "John",
//	})
//
//	// Validate prompt
//	err = prompt.Validate()
package prompt

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"text/template"
	"time"
)

// Prompt represents pure templates hydrated at runtime in the CMP framework.
// It provides a structured way to manage prompt templates with versioning,
// validation, and rendering capabilities.
type Prompt struct {
	// Name is the unique identifier for the prompt
	Name        string `json:"name" yaml:"name"`
	
	// Version follows semantic versioning (e.g., "1.0.0")
	Version     string `json:"version" yaml:"version"`
	
	// Description provides human-readable information about the prompt
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// Template contains the raw template text with placeholders
	Template  string       `json:"template" yaml:"template"`
	
	// Variables defines the expected template variables and their types
	Variables []Variable   `json:"variables,omitempty" yaml:"variables,omitempty"`
	
	// Config contains prompt behavior configuration
	Config    PromptConfig `json:"config" yaml:"config"`

	// CreatedAt is the timestamp when the prompt was created
	CreatedAt time.Time         `json:"created_at" yaml:"created_at"`
	
	// UpdatedAt is the timestamp when the prompt was last modified
	UpdatedAt time.Time         `json:"updated_at" yaml:"updated_at"`
	
	// Metadata contains additional key-value pairs for extensibility
	Metadata  map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// Variable defines a template variable with type information and validation rules.
// It provides structure for template variables to ensure proper usage and validation.
type Variable struct {
	// Name is the variable identifier used in templates
	Name        string `json:"name" yaml:"name"`
	
	// Type defines the variable type: "string", "context", "memory", "user"
	Type        string `json:"type" yaml:"type"`
	
	// Required indicates if the variable must be provided during rendering
	Required    bool   `json:"required" yaml:"required"`
	
	// Description provides human-readable information about the variable
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	
	// Default provides a fallback value if the variable is not provided
	Default     string `json:"default,omitempty" yaml:"default,omitempty"`
}

// PromptConfig defines prompt behavior and rendering configuration.
// It controls how templates are processed and what output format is expected.
type PromptConfig struct {
	// MaxTokens limits the maximum number of tokens in the rendered output
	MaxTokens   int     `json:"max_tokens" yaml:"max_tokens"`
	
	// Temperature controls randomness in AI model responses (0.0 to 1.0)
	Temperature float64 `json:"temperature" yaml:"temperature"`
	
	// Format specifies the output format: "json", "markdown", "text"
	Format      string  `json:"format,omitempty" yaml:"format,omitempty"`
	
	// Language specifies the language for the prompt content
	Language    string  `json:"language,omitempty" yaml:"language,omitempty"`
}

// New creates a new Prompt with default values.
// It initializes a prompt with sensible defaults and timestamps.
//
// Parameters:
//   - name: Unique identifier for the prompt
//   - version: Semantic version string
//   - template: Raw template text with placeholders
//
// Returns:
//   - *Prompt: A new prompt instance with default configuration
func New(name, version, template string) *Prompt {
	return &Prompt{
		Name:      name,
		Version:   version,
		Template:  template,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Metadata:  make(map[string]string),
		Config: PromptConfig{
			MaxTokens:   1000,
			Temperature: 0.1,
			Format:      "text",
			Language:    "en",
		},
	}
}

// Render hydrates the template with provided data.
// It processes the template using Go's text/template engine and returns
// the rendered string with all placeholders replaced.
//
// Parameters:
//   - data: Map of variable names to values for template substitution
//
// Returns:
//   - string: The rendered template with all placeholders replaced
//   - error: Any error that occurred during template processing
func (p *Prompt) Render(data map[string]interface{}) (string, error) {
	tmpl, err := template.New(p.Name).Parse(p.Template)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}

	return buf.String(), nil
}

// AddVariable adds a new template variable to the prompt.
// It provides a convenient way to define expected variables with their
// types and requirements for validation purposes.
//
// Parameters:
//   - name: Variable identifier
//   - varType: Variable type ("string", "context", "memory", "user")
//   - required: Whether the variable is mandatory
//   - description: Human-readable description of the variable
func (p *Prompt) AddVariable(name, varType string, required bool, description string) {
	variable := Variable{
		Name:        name,
		Type:        varType,
		Required:    required,
		Description: description,
	}
	p.Variables = append(p.Variables, variable)
}

// Validate ensures the prompt is properly configured.
// It checks that all required fields are present and valid according
// to the prompt schema and business rules.
//
// Returns:
//   - error: Validation error if the prompt is invalid, nil otherwise
func (p *Prompt) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("prompt name is required")
	}
	if p.Version == "" {
		return fmt.Errorf("prompt version is required")
	}
	if p.Template == "" {
		return fmt.Errorf("prompt template is required")
	}
	return nil
}

// ToJSON converts the prompt to JSON format.
// It serializes the prompt structure to JSON for storage or transmission.
//
// Returns:
//   - []byte: JSON representation of the prompt
//   - error: Any error that occurred during serialization
func (p *Prompt) ToJSON() ([]byte, error) {
	return json.MarshalIndent(p, "", "  ")
}

// FromJSON creates a prompt from JSON data.
// It deserializes JSON data into a prompt structure.
//
// Parameters:
//   - data: JSON bytes representing a prompt
//
// Returns:
//   - *Prompt: Deserialized prompt instance
//   - error: Any error that occurred during deserialization
func FromJSON(data []byte) (*Prompt, error) {
	var prompt Prompt
	if err := json.Unmarshal(data, &prompt); err != nil {
		return nil, err
	}
	return &prompt, nil
}

// GetSHA returns a content-based hash for versioning.
// It generates a SHA-256 hash of the prompt's JSON representation
// to enable content-based versioning and change detection.
//
// Returns:
//   - string: SHA-256 hash prefixed with "sha256:"
//   - error: Any error that occurred during hashing
func (p *Prompt) GetSHA() (string, error) {
	data, err := p.ToJSON()
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return fmt.Sprintf("sha256:%s", hex.EncodeToString(hash[:])), nil
}

// Example templates for common use cases
const (
	// RAG response template provides a structured format for RAG system responses
	RAGResponseTemplate = `# {{.component_name}} Response

**User Query:** {{.user_query}}

**Retrieved Context:**
{{range .knowledge_results}}
### Document {{.index}}: {{.title}}
{{.content}}

---
{{end}}

**Response:**
Based on the retrieved documents, here is the answer to your question:

{{.response}}

**Sources:**
{{range .knowledge_results}}
- {{.title}} (relevance: {{.relevance_score}})
{{end}}`

	// Agent conversation template provides a format for conversational AI responses
	AgentConversationTemplate = `# {{.agent_name}} Conversation

**Context:** {{.conversation_context}}

**User Message:** {{.user_message}}

**Agent Response:**
{{.agent_response}}

**Next Actions:**
{{range .suggested_actions}}
- {{.action}}: {{.description}}
{{end}}`

	// Workflow step template provides a format for workflow step execution
	WorkflowStepTemplate = `# {{.workflow_name}} - Step {{.step_number}}: {{.step_name}}

**Input:** {{.step_input}}

**Processing:**
{{.step_processing}}

**Output:** {{.step_output}}

**Status:** {{.step_status}}
{{if .step_error}}
**Error:** {{.step_error}}
{{end}}`
)
