# GoMarkDoc Documentation Guide

This guide explains how to use `gomarkdoc` to generate comprehensive documentation from Go code comments in the Contexis CMP Framework.

## Overview

`gomarkdoc` is a tool that generates Markdown documentation from Go code comments. It parses Go source files and extracts documentation from:
- Package comments
- Function comments
- Struct field comments
- Type comments
- Constant and variable comments

## Installation

### Prerequisites

- Go 1.19 or later
- Git (for version information)

### Install GoMarkDoc

```bash
# Install gomarkdoc
go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest

# Verify installation
gomarkdoc --version
```

## Basic Usage

### Generate Documentation for a Package

```bash
# Generate documentation for the main CLI package
gomarkdoc ./src/cli

# Generate documentation for the core prompt package
gomarkdoc ./src/core/prompt

# Generate documentation for the configuration package
gomarkdoc ./src/cli/config
```

### Generate Documentation for Multiple Packages

```bash
# Generate documentation for all packages in the project
gomarkdoc ./src/...

# Generate documentation for specific packages
gomarkdoc ./src/cli ./src/core ./src/runtime
```

### Output Options

```bash
# Output to stdout (default)
gomarkdoc ./src/cli

# Output to a file
gomarkdoc ./src/cli -o docs/api/cli.md

# Output to multiple files (one per package)
gomarkdoc ./src/... -o docs/api/

# Include version information
gomarkdoc ./src/cli --include-vendor

# Include private functions and types
gomarkdoc ./src/cli --private
```

## Documentation Standards

### Package Documentation

Every package should have comprehensive documentation that explains:
- Purpose and functionality
- Key features
- Example usage
- Dependencies

```go
// Package main provides the Contexis CMP Framework CLI application.
//
// The Contexis CLI is a Rails-inspired command-line interface for building
// reproducible AI applications using the Context-Memory-Prompt (CMP) architecture.
// It provides commands for project initialization, component generation, testing,
// and deployment of AI applications.
//
// Key Features:
//   - Local-first development with out-of-the-box local models
//   - Component generation (RAG, agents, workflows)
//   - Memory management and vector search
//   - Drift detection and testing
//   - Production migration tools
//
// Example Usage:
//
//	# Initialize a new project
//	ctx init my-ai-app
//
//	# Generate a RAG component
//	ctx generate rag CustomerDocs
//
//	# Run tests with drift detection
//	ctx test --drift-detection
//
//	# Start development server
//	ctx serve --addr :8000
package main
```

### Function Documentation

Functions should be documented with:
- Purpose and behavior
- Parameters and their types
- Return values and their types
- Example usage
- Error conditions

```go
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
    // Implementation...
}
```

### Struct Documentation

Structs should be documented with:
- Purpose and usage
- Field descriptions
- Example usage
- Related types

```go
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
```

### Type Documentation

Types should be documented with:
- Purpose and usage
- Valid values or constraints
- Example usage

```go
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
```

## Advanced Usage

### Custom Templates

Create custom templates for different documentation styles:

```bash
# Use a custom template
gomarkdoc ./src/cli --template ./docs/templates/api.md.tmpl

# Use a different output format
gomarkdoc ./src/cli --format html
```

### Configuration File

Create a `.gomarkdoc.yml` configuration file:

```yaml
# .gomarkdoc.yml
output: docs/api/
format: markdown
include-vendor: false
private: false
template: ./docs/templates/api.md.tmpl
packages:
  - ./src/cli
  - ./src/core
  - ./src/runtime
```

### Integration with CI/CD

Add documentation generation to your CI/CD pipeline:

```yaml
# .github/workflows/docs.yml
name: Generate Documentation

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  docs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v3
        with:
          go-version: '1.19'
      
      - name: Install gomarkdoc
        run: go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest
      
      - name: Generate documentation
        run: gomarkdoc ./src/... -o docs/api/
      
      - name: Commit documentation
        run: |
          git config --local user.email "action@github.com"
          git config --local user.name "GitHub Action"
          git add docs/api/
          git commit -m "Update API documentation" || exit 0
          git push
```

## Best Practices

### 1. Consistent Documentation Style

- Use clear, concise language
- Include examples for complex functions
- Document all exported functions and types
- Use consistent formatting and structure

### 2. Comprehensive Examples

```go
// Example usage in package documentation
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
```

### 3. Parameter and Return Documentation

```go
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
    // Implementation...
}
```

### 4. Error Documentation

```go
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
```

### 5. Type Constraints and Valid Values

```go
// Type defines the variable type: "string", "context", "memory", "user"
Type        string `json:"type" yaml:"type"`

// Privacy specifies the privacy level: "user_isolated", "shared", "public"
Privacy    string `json:"privacy" yaml:"privacy"`

// Temperature controls response randomness (0.0 to 1.0)
Temperature float64 `json:"temperature" yaml:"temperature"`
```

## Generated Documentation Structure

The generated documentation follows this structure:

```markdown
# Package Name

Package description and overview.

## Index

- [Types](#types)
- [Functions](#functions)
- [Constants](#constants)
- [Variables](#variables)

## Types

### [TypeName](link-to-type)

Type description.

#### Fields

- `FieldName` `FieldType` - Field description

#### Methods

- `MethodName(params) returns` - Method description

## Functions

### [FunctionName](link-to-function)

Function description.

**Parameters:**
- `param` `type` - Parameter description

**Returns:**
- `type` - Return value description

## Examples

Code examples showing usage.

## Constants

- `ConstantName` = `value` - Constant description

## Variables

- `VariableName` `type` - Variable description
```

## Maintenance

### Regular Updates

- Generate documentation after each release
- Review and update examples regularly
- Keep documentation in sync with code changes

### Quality Checks

```bash
# Check for undocumented exported functions
gomarkdoc ./src/... --check

# Validate documentation completeness
gomarkdoc ./src/... --validate
```

### Documentation Review

- Review generated documentation for accuracy
- Ensure examples are up-to-date
- Check for broken links or references
- Verify formatting and readability

## Troubleshooting

### Common Issues

1. **Missing Documentation**
   - Ensure all exported functions have comments
   - Check that package comments are present
   - Verify struct field comments are complete

2. **Formatting Issues**
   - Use consistent indentation in examples
   - Escape special characters properly
   - Follow Go comment conventions

3. **Generation Errors**
   - Check for syntax errors in Go code
   - Verify file permissions
   - Ensure all dependencies are available

### Debug Commands

```bash
# Verbose output for debugging
gomarkdoc ./src/cli --verbose

# Check specific package
gomarkdoc ./src/cli/main.go

# Generate with debug information
gomarkdoc ./src/cli --debug
```

## Integration with Other Tools

### Hugo Integration

```bash
# Generate documentation for Hugo
gomarkdoc ./src/... -o content/api/ --format hugo
```

### Docusaurus Integration

```bash
# Generate documentation for Docusaurus
gomarkdoc ./src/... -o docs/api/ --format docusaurus
```

### GitHub Pages

```bash
# Generate documentation for GitHub Pages
gomarkdoc ./src/... -o docs/ --format github-pages
```

## Resources

- [GoMarkDoc Documentation](https://github.com/princjef/gomarkdoc)
- [Go Documentation Guidelines](https://golang.org/doc/comment)
- [Effective Go - Documentation](https://golang.org/doc/effective_go.html#commentary)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments#comment-sentences)
