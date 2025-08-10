package memory

import (
	"encoding/json"
	"fmt"
	"time"
)

// Memory represents versioned knowledge stores in the CMP framework
type Memory struct {
	Name        string `json:"name" yaml:"name"`
	Version     string `json:"version" yaml:"version"`
	Type        string `json:"type" yaml:"type"` // "vector", "episodic", "semantic"
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	Content []MemoryItem `json:"content" yaml:"content"`
	Schema  MemorySchema `json:"schema,omitempty" yaml:"schema,omitempty"`
	Config  MemoryConfig `json:"config" yaml:"config"`

	CreatedAt time.Time         `json:"created_at" yaml:"created_at"`
	UpdatedAt time.Time         `json:"updated_at" yaml:"updated_at"`
	Metadata  map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// MemoryItem represents a single piece of knowledge
type MemoryItem struct {
	ID        string                 `json:"id" yaml:"id"`
	Content   string                 `json:"content" yaml:"content"`
	Type      string                 `json:"type" yaml:"type"` // "text", "document", "conversation"
	Embedding []float64              `json:"embedding,omitempty" yaml:"embedding,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	CreatedAt time.Time              `json:"created_at" yaml:"created_at"`
	UpdatedAt time.Time              `json:"updated_at" yaml:"updated_at"`
}

// MemorySchema defines the structure of memory items
type MemorySchema struct {
	Fields   []SchemaField `json:"fields" yaml:"fields"`
	Required []string      `json:"required,omitempty" yaml:"required,omitempty"`
	Indexes  []string      `json:"indexes,omitempty" yaml:"indexes,omitempty"`
}

// SchemaField defines a field in the memory schema
type SchemaField struct {
	Name        string `json:"name" yaml:"name"`
	Type        string `json:"type" yaml:"type"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Required    bool   `json:"required" yaml:"required"`
}

// MemoryConfig defines memory behavior and storage
type MemoryConfig struct {
	Provider   string            `json:"provider" yaml:"provider"`     // "sqlite", "postgres", "chroma", "pinecone"
	Embeddings string            `json:"embeddings" yaml:"embeddings"` // "openai", "sentence-transformers"
	ChunkSize  int               `json:"chunk_size" yaml:"chunk_size"`
	Overlap    int               `json:"overlap" yaml:"overlap"`
	MaxTokens  int               `json:"max_tokens" yaml:"max_tokens"`
	Settings   map[string]string `json:"settings,omitempty" yaml:"settings,omitempty"`
}

// New creates a new Memory with default values
func New(name, version, memoryType string) *Memory {
	return &Memory{
		Name:      name,
		Version:   version,
		Type:      memoryType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Metadata:  make(map[string]string),
		Content:   make([]MemoryItem, 0),
		Config: MemoryConfig{
			ChunkSize: 1000,
			Overlap:   200,
			MaxTokens: 8000,
		},
	}
}

// AddItem adds a new memory item
func (m *Memory) AddItem(content, itemType string, metadata map[string]interface{}) {
	item := MemoryItem{
		ID:        fmt.Sprintf("%s-%d", m.Name, len(m.Content)+1),
		Content:   content,
		Type:      itemType,
		Metadata:  metadata,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m.Content = append(m.Content, item)
	m.UpdatedAt = time.Now()
}

// Search performs semantic search on memory content
func (m *Memory) Search(query string, limit int) ([]MemoryItem, error) {
	// TODO: Implement semantic search with embeddings
	// This is a placeholder implementation
	results := make([]MemoryItem, 0)
	for _, item := range m.Content {
		// Simple text matching for now
		if len(results) < limit {
			results = append(results, item)
		}
	}
	return results, nil
}

// Validate ensures the memory is properly configured
func (m *Memory) Validate() error {
	if m.Name == "" {
		return fmt.Errorf("memory name is required")
	}
	if m.Version == "" {
		return fmt.Errorf("memory version is required")
	}
	if m.Type == "" {
		return fmt.Errorf("memory type is required")
	}
	return nil
}

// ToJSON converts the memory to JSON
func (m *Memory) ToJSON() ([]byte, error) {
	return json.MarshalIndent(m, "", "  ")
}

// FromJSON creates a memory from JSON
func FromJSON(data []byte) (*Memory, error) {
	var mem Memory
	if err := json.Unmarshal(data, &mem); err != nil {
		return nil, err
	}
	return &mem, nil
}

// GetSHA returns a content-based hash for versioning
func (m *Memory) GetSHA() (string, error) {
	data, err := m.ToJSON()
	if err != nil {
		return "", err
	}
	// TODO: Implement proper SHA256 hashing
	return fmt.Sprintf("sha256:%x", data), nil
}
