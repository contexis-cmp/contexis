package helpers

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/contexis-cmp/contexis/src/cli/commands"
	"github.com/contexis-cmp/contexis/src/cli/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// TestContext provides a test context with logging
func TestContext(t *testing.T) context.Context {
	ctx := context.Background()

	// Initialize test logger
	if err := logger.InitLogger("debug", "console"); err != nil {
		t.Fatalf("Failed to initialize test logger: %v", err)
	}

	return ctx
}

// TestLogger returns a test logger instance
func TestLogger(t *testing.T) *zap.Logger {
	return logger.WithContext(TestContext(t))
}

// CreateTempDir creates a temporary directory for testing
func CreateTempDir(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "cmp_test_*")
	require.NoError(t, err, "Failed to create temp directory")

	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})

	return tempDir
}

// ChangeToTempDir changes to a temporary directory and returns cleanup function
func ChangeToTempDir(t *testing.T) (string, func()) {
	tempDir := CreateTempDir(t)

	originalDir, err := os.Getwd()
	require.NoError(t, err, "Failed to get current directory")

	err = os.Chdir(tempDir)
	require.NoError(t, err, "Failed to change to temp directory")

	cleanup := func() {
		os.Chdir(originalDir)
	}

	return tempDir, cleanup
}

// AssertFileExists checks if a file exists
func AssertFileExists(t *testing.T, path string) {
	_, err := os.Stat(path)
	assert.NoError(t, err, "File should exist: %s", path)
}

// AssertFileNotExists checks if a file does not exist
func AssertFileNotExists(t *testing.T, path string) {
	_, err := os.Stat(path)
	assert.True(t, os.IsNotExist(err), "File should not exist: %s", path)
}

// AssertDirExists checks if a directory exists
func AssertDirExists(t *testing.T, path string) {
	info, err := os.Stat(path)
	assert.NoError(t, err, "Directory should exist: %s", path)
	assert.True(t, info.IsDir(), "Path should be a directory: %s", path)
}

// AssertFileContent checks if a file contains expected content
func AssertFileContent(t *testing.T, path, expectedContent string) {
	content, err := os.ReadFile(path)
	require.NoError(t, err, "Failed to read file: %s", path)
	assert.Contains(t, string(content), expectedContent, "File should contain expected content")
}

// CreateTestFile creates a test file with given content
func CreateTestFile(t *testing.T, path, content string) {
	err := os.MkdirAll(filepath.Dir(path), 0755)
	require.NoError(t, err, "Failed to create directory for test file")

	err = os.WriteFile(path, []byte(content), 0644)
	require.NoError(t, err, "Failed to create test file")
}

// AssertYAMLValid checks if YAML content is valid
func AssertYAMLValid(t *testing.T, content string) {
	// Basic YAML validation - check for common syntax errors
	assert.NotContains(t, content, "{{", "YAML should not contain unprocessed template variables")
	assert.NotContains(t, content, "}}", "YAML should not contain unprocessed template variables")
}

// AssertTemplateProcessed checks if template was processed correctly
func AssertTemplateProcessed(t *testing.T, content string) {
	// Check that template variables were replaced
	assert.NotContains(t, content, "{{", "Template should be processed")
	assert.NotContains(t, content, "}}", "Template should be processed")
}

// TestCase represents a test case with input and expected output
type TestCase struct {
	Name        string
	Input       interface{}
	Expected    interface{}
	ExpectError bool
	Description string
}

// RunTestCases runs a slice of test cases
func RunTestCases[T any, U any](t *testing.T, testCases []TestCase, testFunc func(T) (U, error)) {
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			input, ok := tc.Input.(T)
			require.True(t, ok, "Invalid input type for test case: %s", tc.Name)

			result, err := testFunc(input)

			if tc.ExpectError {
				assert.Error(t, err, tc.Description)
			} else {
				assert.NoError(t, err, tc.Description)
				if tc.Expected != nil {
					assert.Equal(t, tc.Expected, result, tc.Description)
				}
			}
		})
	}
}

// JoinStrings joins strings with a separator
func JoinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

// ContainsWhitespace checks if a string contains whitespace
func ContainsWhitespace(s string) bool {
	for _, char := range s {
		if char == ' ' || char == '\t' || char == '\n' || char == '\r' {
			return true
		}
	}
	return false
}

// TestFixtures provides common test data
type TestFixtures struct{}

// ValidAgentNames returns valid agent names for testing
func (tf *TestFixtures) ValidAgentNames() []string {
	return []string{
		"SupportBot",
		"EmailBot",
		"FileBot",
		"test_agent",
		"agent-123",
		"MyAgent",
	}
}

// ValidWorkflowNames returns a list of valid workflow names for testing
func (tf *TestFixtures) ValidWorkflowNames() []string {
	return []string{
		"ContentPipeline",
		"DataProcessing",
		"MLTraining",
		"test_workflow",
		"workflow-123",
		"MyWorkflow",
	}
}

// InvalidAgentNames returns invalid agent names for testing
func (tf *TestFixtures) InvalidAgentNames() []string {
	return []string{
		"",
		"a", // too short
		"invalid name with spaces",
		"invalid-name-with-special-chars!@#",
		"invalid/name/with/slashes",
		"invalid\\name\\with\\backslashes",
	}
}

// ValidTools returns valid tool combinations for testing
func (tf *TestFixtures) ValidTools() [][]string {
	return [][]string{
		{"web_search"},
		{"database"},
		{"api"},
		{"file_system"},
		{"email"},
		{"web_search", "database"},
		{"api", "email"},
		{"file_system", "api", "database"},
		{}, // empty tools
	}
}

// InvalidTools returns invalid tool combinations for testing
func (tf *TestFixtures) InvalidTools() [][]string {
	return [][]string{
		{"invalid_tool"},
		{"web_search", "invalid_tool"},
		{"database", "invalid_tool", "api"},
	}
}

// ValidMemoryTypes returns valid memory types for testing
func (tf *TestFixtures) ValidMemoryTypes() []string {
	return []string{
		"episodic",
		"none",
	}
}

// InvalidMemoryTypes returns invalid memory types for testing
func (tf *TestFixtures) InvalidMemoryTypes() []string {
	return []string{
		"invalid_memory",
		"",
		"memory",
		"episodic_memory",
	}
}

// ValidAgentConfigs returns valid agent configurations for testing
func (tf *TestFixtures) ValidAgentConfigs() []commands.AgentConfig {
	return []commands.AgentConfig{
		{
			Name:           "TestAgent1",
			Tools:          []string{"web_search", "database"},
			Memory:         "episodic",
			Description:    "Test agent 1",
			Version:        "1.0.0",
			Persona:        "Professional assistant",
			Capabilities:   []string{"conversation", "tool_usage"},
			Limitations:    []string{"no_personal_data"},
			BusinessRules:  []string{"always_helpful"},
			BaselineDate:   time.Now().Format("2006-01-02"),
			AdminEmail:     "test@example.com",
			Tone:           "professional",
			Format:         "json",
			MaxTokens:      500,
			Temperature:    0.1,
			MemoryType:     "episodic",
			MaxHistory:     10,
			Privacy:        "user_isolated",
			DriftThreshold: 0.85,
		},
		{
			Name:           "TestAgent2",
			Tools:          []string{},
			Memory:         "none",
			Description:    "Test agent 2",
			Version:        "1.0.0",
			Persona:        "Simple assistant",
			Capabilities:   []string{"conversation"},
			Limitations:    []string{"no_personal_data"},
			BusinessRules:  []string{"always_helpful"},
			BaselineDate:   time.Now().Format("2006-01-02"),
			AdminEmail:     "test@example.com",
			Tone:           "friendly",
			Format:         "text",
			MaxTokens:      300,
			Temperature:    0.2,
			MemoryType:     "none",
			MaxHistory:     0,
			Privacy:        "user_isolated",
			DriftThreshold: 0.8,
		},
	}
}

// ValidToolDefinitions returns valid tool definitions for testing
func (tf *TestFixtures) ValidToolDefinitions() []commands.Tool {
	return []commands.Tool{
		{
			Name:        "web_search",
			URI:         "mcp://web.search",
			Description: "Search the web for current information",
		},
		{
			Name:        "database",
			URI:         "mcp://database.query",
			Description: "Query database for user and order information",
		},
		{
			Name:        "api",
			URI:         "mcp://api.call",
			Description: "Make API calls to external services",
		},
		{
			Name:        "file_system",
			URI:         "mcp://file.read",
			Description: "Read and write files",
		},
		{
			Name:        "email",
			URI:         "mcp://email.send",
			Description: "Send and read emails",
		},
	}
}

// ExpectedDirectoryStructure returns expected directory structure for an agent
func (tf *TestFixtures) ExpectedDirectoryStructure(agentName string) []string {
	return []string{
		"contexts/" + agentName,
		"memory/" + agentName,
		"memory/" + agentName + "/episodic",
		"memory/" + agentName + "/user_preferences",
		"prompts/" + agentName,
		"tools/" + agentName,
		"tests/" + agentName,
	}
}

// ExpectedFiles returns expected files for an agent
func (tf *TestFixtures) ExpectedFiles(agentName string) []string {
	return []string{
		"contexts/" + agentName + "/" + agentName + ".ctx",
		"prompts/" + agentName + "/agent_response.md",
		"tests/" + agentName + "/agent_behavior.yaml",
		"memory/" + agentName + "/memory_config.yaml",
		"tools/" + agentName + "/requirements.txt",
	}
}

// TestUtils provides utility functions for testing
type TestUtils struct{}

// CreateTestProjectStructure creates a test project structure
func (tu *TestUtils) CreateTestProjectStructure(t *testing.T, basePath string) {
	dirs := []string{
		"contexts",
		"memory",
		"prompts",
		"tools",
		"tests",
		"templates",
		"templates/agent",
		"templates/rag",
		"templates/workflow",
	}

	for _, dir := range dirs {
		path := filepath.Join(basePath, dir)
		err := os.MkdirAll(path, 0755)
		require.NoError(t, err, "Failed to create directory: %s", path)
	}
}

// CreateTestTemplates creates test template files
func (tu *TestUtils) CreateTestTemplates(t *testing.T, basePath string) {
	templates := map[string]string{
		"templates/agent/support_bot.ctx": `name: "{{ .Name }}"
version: "{{ .Version }}"
description: "{{ .Description }}"
role:
  persona: "{{ .Persona }}"
  capabilities: {{ .Capabilities }}
  limitations: {{ .Limitations }}`,

		"templates/agent/agent_response.md": `# Agent Response Template
## Conversation Context
- User ID: [USER_ID]
- Session ID: [SESSION_ID]
## Response Guidelines
- **Tone**: Professional and helpful
- **Format**: JSON
- **Max Tokens**: 500`,

		"templates/agent/agent_behavior.yaml": `# Agent Behavior Test Configuration
agent_name: "{{ .Name }}"
test_version: "{{ .Version }}"
description: "{{ .Description }}"`,

		"templates/agent/requirements.txt": `# CMP Agent Tools Requirements
requests>=2.31.0
urllib3>=2.0.0`,

		"templates/agent/web_search.py": `#!/usr/bin/env python3
"""
Web Search Tool for CMP Agents
"""
import requests
import logging`,

		"templates/agent/database.py": `#!/usr/bin/env python3
"""
Database Tool for CMP Agents
"""
import sqlite3
import logging`,
	}

	for path, content := range templates {
		fullPath := filepath.Join(basePath, path)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		require.NoError(t, err, "Failed to create directory for template: %s", fullPath)

		err = os.WriteFile(fullPath, []byte(content), 0644)
		require.NoError(t, err, "Failed to create template file: %s", fullPath)
	}
}

// AssertDirectoryStructure checks if the expected directory structure exists
func (tu *TestUtils) AssertDirectoryStructure(t *testing.T, basePath string, expectedDirs []string) {
	for _, dir := range expectedDirs {
		path := filepath.Join(basePath, dir)
		AssertDirExists(t, path)
	}
}

// AssertFileStructure checks if the expected files exist
func (tu *TestUtils) AssertFileStructure(t *testing.T, basePath string, expectedFiles []string) {
	for _, file := range expectedFiles {
		path := filepath.Join(basePath, file)
		AssertFileExists(t, path)
	}
}

// AssertTemplateContent checks if template content is valid
func (tu *TestUtils) AssertTemplateContent(t *testing.T, content string, expectedVariables []string) {
	// Check that template variables are present
	for _, variable := range expectedVariables {
		assert.Contains(t, content, "{{ ."+variable+" }}",
			"Template should contain variable: %s", variable)
	}

	// Check that template is properly formatted
	assert.NotContains(t, content, "{{{", "Template should not contain malformed variables")
	assert.NotContains(t, content, "}}}", "Template should not contain malformed variables")
}

// AssertGeneratedContent checks if generated content is valid
func (tu *TestUtils) AssertGeneratedContent(t *testing.T, content string, expectedContent []string) {
	for _, expected := range expectedContent {
		assert.Contains(t, content, expected,
			"Generated content should contain: %s", expected)
	}

	// Check that template variables are processed
	assert.NotContains(t, content, "{{ .", "Generated content should not contain unprocessed template variables")
	assert.NotContains(t, content, "}}", "Generated content should not contain unprocessed template variables")
}

// CleanupTestFiles removes test files and directories
func (tu *TestUtils) CleanupTestFiles(t *testing.T, paths []string) {
	for _, path := range paths {
		if err := os.RemoveAll(path); err != nil {
			t.Logf("Warning: Failed to cleanup test file: %s", path)
		}
	}
}

// ValidateYAMLContent validates YAML content structure
func (tu *TestUtils) ValidateYAMLContent(t *testing.T, content string) {
	// Basic YAML validation
	assert.Contains(t, content, ":", "YAML should contain key-value pairs")
	assert.NotContains(t, content, "{{", "YAML should not contain template variables")
	assert.NotContains(t, content, "}}", "YAML should not contain template variables")
}

// ValidateJSONContent validates JSON content structure
func (tu *TestUtils) ValidateJSONContent(t *testing.T, content string) {
	// Basic JSON validation
	assert.Contains(t, content, "{", "JSON should contain opening brace")
	assert.Contains(t, content, "}", "JSON should contain closing brace")
	assert.NotContains(t, content, "{{", "JSON should not contain template variables")
	assert.NotContains(t, content, "}}", "JSON should not contain template variables")
}

// ValidatePythonContent validates Python content structure
func (tu *TestUtils) ValidatePythonContent(t *testing.T, content string) {
	// Basic Python validation
	assert.Contains(t, content, "#!/usr/bin/env python3", "Python file should have shebang")
	assert.Contains(t, content, "import", "Python file should contain imports")
	assert.Contains(t, content, "def", "Python file should contain function definitions")
}

// CountFilesInDirectory counts files in a directory
func (tu *TestUtils) CountFilesInDirectory(t *testing.T, dirPath string) int {
	entries, err := os.ReadDir(dirPath)
	require.NoError(t, err, "Failed to read directory: %s", dirPath)

	count := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			count++
		}
	}

	return count
}

// CountDirectoriesInDirectory counts directories in a directory
func (tu *TestUtils) CountDirectoriesInDirectory(t *testing.T, dirPath string) int {
	entries, err := os.ReadDir(dirPath)
	require.NoError(t, err, "Failed to read directory: %s", dirPath)

	count := 0
	for _, entry := range entries {
		if entry.IsDir() {
			count++
		}
	}

	return count
}

// GetFileContent reads and returns file content
func (tu *TestUtils) GetFileContent(t *testing.T, filePath string) string {
	content, err := os.ReadFile(filePath)
	require.NoError(t, err, "Failed to read file: %s", filePath)
	return string(content)
}

// AssertFileContains checks if file contains expected content
func (tu *TestUtils) AssertFileContains(t *testing.T, filePath string, expectedContent []string) {
	content := tu.GetFileContent(t, filePath)

	for _, expected := range expectedContent {
		assert.Contains(t, content, expected,
			"File %s should contain: %s", filePath, expected)
	}
}

// AssertFileNotContains checks if file does not contain unexpected content
func (tu *TestUtils) AssertFileNotContains(t *testing.T, filePath string, unexpectedContent []string) {
	content := tu.GetFileContent(t, filePath)

	for _, unexpected := range unexpectedContent {
		assert.NotContains(t, content, unexpected,
			"File %s should not contain: %s", filePath, unexpected)
	}
}

// ValidateAgentName validates agent name format
func (tu *TestUtils) ValidateAgentName(t *testing.T, name string) {
	// Agent name validation rules
	assert.NotEmpty(t, name, "Agent name should not be empty")
	assert.True(t, len(name) >= 2, "Agent name should be at least 2 characters")
	assert.False(t, strings.Contains(name, " "), "Agent name should not contain spaces")
	assert.False(t, strings.Contains(name, "/"), "Agent name should not contain slashes")
	assert.False(t, strings.Contains(name, "\\"), "Agent name should not contain backslashes")
	assert.False(t, strings.Contains(name, "!"), "Agent name should not contain special characters")
}

// ValidateToolName validates tool name format
func (tu *TestUtils) ValidateToolName(t *testing.T, name string) {
	// Tool name validation rules
	validTools := []string{"web_search", "database", "api", "file_system", "email"}

	assert.NotEmpty(t, name, "Tool name should not be empty")
	assert.Contains(t, validTools, name, "Tool name should be valid: %s", name)
}

// ValidateMemoryType validates memory type format
func (tu *TestUtils) ValidateMemoryType(t *testing.T, memoryType string) {
	// Memory type validation rules
	validTypes := []string{"episodic", "none"}

	assert.NotEmpty(t, memoryType, "Memory type should not be empty")
	assert.Contains(t, validTypes, memoryType, "Memory type should be valid: %s", memoryType)
}

// SplitString splits a string by separator
func SplitString(s, sep string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, sep)
}

// TrimSpace removes leading and trailing whitespace
func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

// ValidateWorkflowName validates workflow name format
func (tu *TestUtils) ValidateWorkflowName(t *testing.T, name string) {
	// Workflow name validation rules
	assert.NotEmpty(t, name, "Workflow name should not be empty")
	assert.True(t, len(name) >= 2, "Workflow name should be at least 2 characters")
	assert.False(t, strings.Contains(name, " "), "Workflow name should not contain spaces")
	assert.False(t, strings.Contains(name, "/"), "Workflow name should not contain slashes")
	assert.False(t, strings.Contains(name, "\\"), "Workflow name should not contain backslashes")
	// Check for special characters (excluding hyphens and underscores)
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '-' || char == '_') {
			t.Errorf("Workflow name should not contain special characters")
			return
		}
	}
}

// ValidateStepType validates step type format
func (tu *TestUtils) ValidateStepType(t *testing.T, stepType string) {
	// Step type validation rules
	validTypes := []string{"research", "write", "review", "extract", "transform", "load", "analyze", "generate", "validate", "deploy"}

	assert.NotEmpty(t, stepType, "Step type should not be empty")
	assert.Contains(t, validTypes, stepType, "Step type should be valid: %s", stepType)
}
