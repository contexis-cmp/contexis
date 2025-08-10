package integration

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/contexis/cmp/src/cli/commands"
	"github.com/contexis/cmp/tests/unit/helpers"
)

// AgentGeneratorIntegrationTestSuite provides integration tests for agent generation
type AgentGeneratorIntegrationTestSuite struct {
	suite.Suite
	helpers *helpers.TestFixtures
	utils   *helpers.TestUtils
	tempDir string
}

// SetupSuite sets up the test suite
func (suite *AgentGeneratorIntegrationTestSuite) SetupSuite() {
	suite.helpers = &helpers.TestFixtures{}
	suite.utils = &helpers.TestUtils{}
}

// SetupTest sets up each test
func (suite *AgentGeneratorIntegrationTestSuite) SetupTest() {
	suite.tempDir = helpers.CreateTempDir(suite.T())
	
	// Create test project structure
	suite.utils.CreateTestProjectStructure(suite.T(), suite.tempDir)
	suite.utils.CreateTestTemplates(suite.T(), suite.tempDir)
	
	// Change to temp directory
	originalDir, err := os.Getwd()
	require.NoError(suite.T(), err, "Failed to get current directory")
	
	err = os.Chdir(suite.tempDir)
	require.NoError(suite.T(), err, "Failed to change to temp directory")
	
	suite.T().Cleanup(func() {
		os.Chdir(originalDir)
	})
}

// TestAgentGenerationCompleteFlow tests the complete agent generation flow
func (suite *AgentGeneratorIntegrationTestSuite) TestAgentGenerationCompleteFlow() {
	t := suite.T()
	
	testCases := []struct {
		name    string
		agent   string
		tools   string
		memory  string
		expectError bool
	}{
		{
			name:    "simple_agent",
			agent:   "SimpleBot",
			tools:   "web_search",
			memory:  "episodic",
			expectError: false,
		},
		{
			name:    "multi_tool_agent",
			agent:   "MultiToolBot",
			tools:   "web_search,database,api",
			memory:  "episodic",
			expectError: false,
		},
		{
			name:    "no_tools_agent",
			agent:   "NoToolBot",
			tools:   "",
			memory:  "none",
			expectError: false,
		},
		{
			name:    "invalid_tools",
			agent:   "InvalidBot",
			tools:   "invalid_tool",
			memory:  "episodic",
			expectError: true,
		},
		{
			name:    "invalid_memory",
			agent:   "InvalidMemoryBot",
			tools:   "web_search",
			memory:  "invalid_memory",
			expectError: true,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := helpers.TestContext(t)
			
			// Generate agent
			err := commands.GenerateAgent(ctx, tc.agent, tc.tools, tc.memory)
			
			if tc.expectError {
				assert.Error(t, err, "Expected error for invalid configuration")
				return
			}
			
			assert.NoError(t, err, "Agent generation should succeed")
			
			// Verify directory structure
			expectedDirs := suite.helpers.ExpectedDirectoryStructure(tc.agent)
			suite.utils.AssertDirectoryStructure(t, ".", expectedDirs)
			
			// Verify files
			expectedFiles := suite.helpers.ExpectedFiles(tc.agent)
			suite.utils.AssertFileStructure(t, ".", expectedFiles)
			
			// Verify specific file contents
			suite.verifyAgentContextFile(t, tc.agent)
			suite.verifyAgentBehaviorFile(t, tc.agent)
			suite.verifyMemoryConfigFile(t, tc.agent)
			suite.verifyRequirementsFile(t, tc.agent)
			
			// Verify tool files if tools were specified
			if tc.tools != "" {
				suite.verifyToolFiles(t, tc.agent, tc.tools)
			}
		})
	}
}

// TestAgentGenerationWithDifferentTools tests agent generation with various tool combinations
func (suite *AgentGeneratorIntegrationTestSuite) TestAgentGenerationWithDifferentTools() {
	t := suite.T()
	
	toolCombinations := [][]string{
		{"web_search"},
		{"database"},
		{"api"},
		{"file_system"},
		{"email"},
		{"web_search", "database"},
		{"api", "email"},
		{"file_system", "api", "database"},
	}
	
	for _, tools := range toolCombinations {
		t.Run("tools_"+helpers.JoinStrings(tools, "_"), func(t *testing.T) {
			agentName := "ToolTestBot_" + helpers.JoinStrings(tools, "_")
			toolsStr := helpers.JoinStrings(tools, ",")
			
			ctx := helpers.TestContext(t)
			
			// Generate agent
			err := commands.GenerateAgent(ctx, agentName, toolsStr, "episodic")
			assert.NoError(t, err, "Agent generation should succeed")
			
			// Verify tool files
			suite.verifyToolFiles(t, agentName, toolsStr)
			
			// Verify tool count matches
			toolDir := filepath.Join("tools", agentName)
			fileCount := suite.utils.CountFilesInDirectory(t, toolDir)
			expectedCount := len(tools) + 1 // +1 for requirements.txt
			assert.Equal(t, expectedCount, fileCount, "Tool directory should have correct number of files")
		})
	}
}

// TestAgentGenerationWithDifferentMemoryTypes tests agent generation with different memory types
func (suite *AgentGeneratorIntegrationTestSuite) TestAgentGenerationWithDifferentMemoryTypes() {
	t := suite.T()
	
	memoryTypes := []string{"episodic", "none"}
	
	for _, memoryType := range memoryTypes {
		t.Run("memory_"+memoryType, func(t *testing.T) {
			agentName := "MemoryTestBot_" + memoryType
			
			ctx := helpers.TestContext(t)
			
			// Generate agent
			err := commands.GenerateAgent(ctx, agentName, "web_search", memoryType)
			assert.NoError(t, err, "Agent generation should succeed")
			
			// Verify memory configuration
			memoryConfigPath := filepath.Join("memory", agentName, "memory_config.yaml")
			content := suite.utils.GetFileContent(t, memoryConfigPath)
			
			assert.Contains(t, content, memoryType, "Memory config should contain correct memory type")
			
			// Verify memory directory structure
			if memoryType == "episodic" {
				episodicDir := filepath.Join("memory", agentName, "episodic")
				helpers.AssertDirExists(t, episodicDir)
			}
		})
	}
}

// TestAgentGenerationTemplateProcessing tests that templates are processed correctly
func (suite *AgentGeneratorIntegrationTestSuite) TestAgentGenerationTemplateProcessing() {
	t := suite.T()
	
	agentName := "TemplateTestBot"
	
	ctx := helpers.TestContext(t)
	
	// Generate agent
	err := commands.GenerateAgent(ctx, agentName, "web_search", "episodic")
	assert.NoError(t, err, "Agent generation should succeed")
	
	// Verify template processing
	contextPath := filepath.Join("contexts", agentName, agentName+".ctx")
	content := suite.utils.GetFileContent(t, contextPath)
	
	// Check that template variables were replaced
	assert.Contains(t, content, agentName, "Context should contain agent name")
	assert.NotContains(t, content, "{{ .Name }}", "Template should be processed")
	assert.NotContains(t, content, "{{ .Version }}", "Template should be processed")
	assert.NotContains(t, content, "{{ .Description }}", "Template should be processed")
	
	// Verify YAML structure
	suite.utils.ValidateYAMLContent(t, content)
}

// TestAgentGenerationFilePermissions tests that generated files have correct permissions
func (suite *AgentGeneratorIntegrationTestSuite) TestAgentGenerationFilePermissions() {
	t := suite.T()
	
	agentName := "PermissionTestBot"
	
	ctx := helpers.TestContext(t)
	
	// Generate agent
	err := commands.GenerateAgent(ctx, agentName, "web_search", "episodic")
	assert.NoError(t, err, "Agent generation should succeed")
	
	// Check file permissions
	files := suite.helpers.ExpectedFiles(agentName)
	for _, file := range files {
		info, err := os.Stat(file)
		require.NoError(t, err, "Failed to stat file: %s", file)
		
		// Check that files are readable and writable by owner
		mode := info.Mode()
		assert.True(t, mode.IsRegular(), "File should be a regular file: %s", file)
		assert.True(t, mode&0400 != 0, "File should be readable by owner: %s", file)
		assert.True(t, mode&0200 != 0, "File should be writable by owner: %s", file)
	}
}

// TestAgentGenerationErrorHandling tests error handling during agent generation
func (suite *AgentGeneratorIntegrationTestSuite) TestAgentGenerationErrorHandling() {
	t := suite.T()
	
	ctx := helpers.TestContext(t)
	
	// Test with invalid agent name
	err := commands.GenerateAgent(ctx, "", "web_search", "episodic")
	assert.Error(t, err, "Empty agent name should cause error")
	
	// Test with invalid tool
	err = commands.GenerateAgent(ctx, "TestBot", "invalid_tool", "episodic")
	assert.Error(t, err, "Invalid tool should cause error")
	
	// Test with invalid memory type
	err = commands.GenerateAgent(ctx, "TestBot", "web_search", "invalid_memory")
	assert.Error(t, err, "Invalid memory type should cause error")
}

// Helper methods for verification

func (suite *AgentGeneratorIntegrationTestSuite) verifyAgentContextFile(t *testing.T, agentName string) {
	contextPath := filepath.Join("contexts", agentName, agentName+".ctx")
	content := suite.utils.GetFileContent(t, contextPath)
	
	// Verify YAML structure
	suite.utils.ValidateYAMLContent(t, content)
	
	// Verify required fields
	expectedContent := []string{
		"name:",
		"version:",
		"description:",
		"role:",
		"persona:",
		"capabilities:",
		"limitations:",
	}
	suite.utils.AssertFileContains(t, contextPath, expectedContent)
}

func (suite *AgentGeneratorIntegrationTestSuite) verifyAgentBehaviorFile(t *testing.T, agentName string) {
	behaviorPath := filepath.Join("tests", agentName, "agent_behavior.yaml")
	content := suite.utils.GetFileContent(t, behaviorPath)
	
	// Verify YAML structure
	suite.utils.ValidateYAMLContent(t, content)
	
	// Verify required fields
	expectedContent := []string{
		"agent_name:",
		"test_version:",
		"description:",
		"test_cases:",
	}
	suite.utils.AssertFileContains(t, behaviorPath, expectedContent)
}

func (suite *AgentGeneratorIntegrationTestSuite) verifyMemoryConfigFile(t *testing.T, agentName string) {
	memoryPath := filepath.Join("memory", agentName, "memory_config.yaml")
	content := suite.utils.GetFileContent(t, memoryPath)
	
	// Verify YAML structure
	suite.utils.ValidateYAMLContent(t, content)
	
	// Verify required fields
	expectedContent := []string{
		"memory_type:",
		"max_history:",
		"privacy:",
	}
	suite.utils.AssertFileContains(t, memoryPath, expectedContent)
}

func (suite *AgentGeneratorIntegrationTestSuite) verifyRequirementsFile(t *testing.T, agentName string) {
	requirementsPath := filepath.Join("tools", agentName, "requirements.txt")
	content := suite.utils.GetFileContent(t, requirementsPath)
	
	// Verify Python requirements format
	assert.Contains(t, content, "requests", "Requirements should contain requests")
	assert.Contains(t, content, "urllib3", "Requirements should contain urllib3")
}

func (suite *AgentGeneratorIntegrationTestSuite) verifyToolFiles(t *testing.T, agentName, toolsStr string) {
	tools := helpers.SplitString(toolsStr, ",")
	
	for _, tool := range tools {
		tool = helpers.TrimSpace(tool)
		if tool == "" {
			continue
		}
		
		toolPath := filepath.Join("tools", agentName, tool+".py")
		content := suite.utils.GetFileContent(t, toolPath)
		
		// Verify Python file structure
		suite.utils.ValidatePythonContent(t, content)
		
		// Verify tool-specific content
		assert.Contains(t, content, tool, "Tool file should contain tool name")
		assert.Contains(t, content, "import", "Tool file should contain imports")
		assert.Contains(t, content, "def", "Tool file should contain function definitions")
	}
}

// Run the test suite
func TestAgentGeneratorIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(AgentGeneratorIntegrationTestSuite))
}
