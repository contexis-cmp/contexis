package prompt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"
	"time"
)

// Prompt represents pure templates hydrated at runtime in the CMP framework
type Prompt struct {
	Name        string `json:"name" yaml:"name"`
	Version     string `json:"version" yaml:"version"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	Template  string       `json:"template" yaml:"template"`
	Variables []Variable   `json:"variables,omitempty" yaml:"variables,omitempty"`
	Config    PromptConfig `json:"config" yaml:"config"`

	CreatedAt time.Time         `json:"created_at" yaml:"created_at"`
	UpdatedAt time.Time         `json:"updated_at" yaml:"updated_at"`
	Metadata  map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// Variable defines a template variable
type Variable struct {
	Name        string `json:"name" yaml:"name"`
	Type        string `json:"type" yaml:"type"` // "string", "context", "memory", "user"
	Required    bool   `json:"required" yaml:"required"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Default     string `json:"default,omitempty" yaml:"default,omitempty"`
}

// PromptConfig defines prompt behavior
type PromptConfig struct {
	MaxTokens   int     `json:"max_tokens" yaml:"max_tokens"`
	Temperature float64 `json:"temperature" yaml:"temperature"`
	Format      string  `json:"format,omitempty" yaml:"format,omitempty"` // "json", "markdown", "text"
	Language    string  `json:"language,omitempty" yaml:"language,omitempty"`
}

// New creates a new Prompt with default values
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

// Render hydrates the template with provided data
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

// AddVariable adds a new template variable
func (p *Prompt) AddVariable(name, varType string, required bool, description string) {
	variable := Variable{
		Name:        name,
		Type:        varType,
		Required:    required,
		Description: description,
	}
	p.Variables = append(p.Variables, variable)
}

// Validate ensures the prompt is properly configured
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

// ToJSON converts the prompt to JSON
func (p *Prompt) ToJSON() ([]byte, error) {
	return json.MarshalIndent(p, "", "  ")
}

// FromJSON creates a prompt from JSON
func FromJSON(data []byte) (*Prompt, error) {
	var prompt Prompt
	if err := json.Unmarshal(data, &prompt); err != nil {
		return nil, err
	}
	return &prompt, nil
}

// GetSHA returns a content-based hash for versioning
func (p *Prompt) GetSHA() (string, error) {
	data, err := p.ToJSON()
	if err != nil {
		return "", err
	}
	// TODO: Implement proper SHA256 hashing
	return fmt.Sprintf("sha256:%x", data), nil
}

// Example templates for common use cases
const (
	// RAG response template
	RAGResponseTemplate = `Based on the following context, answer the user's question:

Context:
{{range .context}}
- {{.content}}
{{end}}

Question: {{.question}}

Answer:`

	// Agent conversation template
	AgentConversationTemplate = `You are {{.role.persona}}.

Your capabilities include: {{range .role.capabilities}}{{.}}, {{end}}

Your limitations: {{range .role.limitations}}{{.}}, {{end}}

Current conversation context:
{{.conversation_history}}

User: {{.user_input}}

Assistant:`

	// Workflow step template
	WorkflowStepTemplate = `Execute the following step in the workflow:

Step: {{.step_name}}
Description: {{.step_description}}

Input data: {{.input_data}}

Instructions: {{.instructions}}

Output format: {{.output_format}}

Result:`
)
