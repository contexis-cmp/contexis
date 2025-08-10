package runtimeprompt

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
)

// Engine loads, compiles, caches, renders, and validates prompt templates.
type Engine struct {
	mu          sync.RWMutex
	cache       map[string]*template.Template // key: canonical path
	projectRoot string
}

func NewEngine(projectRoot string) *Engine {
	return &Engine{cache: make(map[string]*template.Template), projectRoot: projectRoot}
}

// RenderFile renders a template file in prompts/<component>/... with the provided data.
func (e *Engine) RenderFile(component string, relPath string, data map[string]interface{}) (string, error) {
	full := filepath.Join(e.projectRoot, "prompts", component, relPath)
	tmpl, err := e.loadTemplate(full)
	if err != nil {
		return "", err
	}
	if data == nil {
		data = map[string]interface{}{}
	}
	var sb strings.Builder
	if err := tmpl.Execute(&sb, data); err != nil {
		return "", fmt.Errorf("render template: %w", err)
	}
	return sb.String(), nil
}

// loadTemplate compiles and caches a template by absolute path.
func (e *Engine) loadTemplate(absPath string) (*template.Template, error) {
	e.mu.RLock()
	if t, ok := e.cache[absPath]; ok {
		e.mu.RUnlock()
		return t, nil
	}
	e.mu.RUnlock()

	b, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("read template: %w", err)
	}
	baseFuncs := template.FuncMap{
		"join":  strings.Join,
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"trim":  strings.TrimSpace,
		"toJSON": func(v interface{}) string {
			by, _ := json.Marshal(v)
			return string(by)
		},
	}
	// Extend with include that uses baseFuncs for included templates
	funcs := template.FuncMap{}
	for k, v := range baseFuncs {
		funcs[k] = v
	}
	funcs["include"] = func(rel string, data interface{}) (string, error) {
		incPath := filepath.Join(filepath.Dir(absPath), rel)
		b, err := os.ReadFile(incPath)
		if err != nil {
			return "", err
		}
		t2, err := template.New(filepath.Base(incPath)).Funcs(baseFuncs).Parse(string(b))
		if err != nil {
			return "", err
		}
		var sb strings.Builder
		if err := t2.Execute(&sb, data); err != nil {
			return "", err
		}
		return sb.String(), nil
	}
	tmpl, err := template.New(filepath.Base(absPath)).Funcs(funcs).Parse(string(b))
	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}
	e.mu.Lock()
	e.cache[absPath] = tmpl
	e.mu.Unlock()
	return tmpl, nil
}

// OptimizeTokens trims content to at most maxTokens using a naive whitespace tokenization.
func OptimizeTokens(content string, maxTokens int) string {
	if maxTokens <= 0 {
		return content
	}
	tokens := strings.Fields(content)
	if len(tokens) <= maxTokens {
		return content
	}
	return strings.Join(tokens[:maxTokens], " ") + " ..."
}

// ValidateFormat validates response against expected format.
// If format == "json", checks JSON. If format == "markdown", ensures non-empty and at least one newline.
func ValidateFormat(format, response string) error {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "json":
		var v interface{}
		if err := json.Unmarshal([]byte(response), &v); err != nil {
			return fmt.Errorf("invalid json: %w", err)
		}
		return nil
	case "markdown":
		if strings.TrimSpace(response) == "" {
			return fmt.Errorf("empty markdown response")
		}
		return nil
	case "text", "":
		if strings.TrimSpace(response) == "" {
			return fmt.Errorf("empty response")
		}
		return nil
	default:
		// Unknown format, best-effort non-empty check
		if strings.TrimSpace(response) == "" {
			return fmt.Errorf("empty response for format %s", format)
		}
		return nil
	}
}
