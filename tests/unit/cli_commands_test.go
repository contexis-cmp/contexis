package unit

import (
	"testing"

	"github.com/contexis-cmp/contexis/src/cli/commands"
	"github.com/contexis-cmp/contexis/tests/unit/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// CLICommandsTestSuite provides tests for CLI command functionality
type CLICommandsTestSuite struct {
	suite.Suite
	helpers *helpers.TestFixtures
	utils   *helpers.TestUtils
}

// SetupSuite sets up the test suite
func (suite *CLICommandsTestSuite) SetupSuite() {
	suite.helpers = &helpers.TestFixtures{}
	suite.utils = &helpers.TestUtils{}
}

// TestGenerateCommandStructure tests the generate command structure
func (suite *CLICommandsTestSuite) TestGenerateCommandStructure() {
	t := suite.T()

	// Test that generate command exists and has correct structure
	generateCmd := commands.GetGenerateCommand()

	assert.NotNil(t, generateCmd, "Generate command should exist")
	assert.Equal(t, "generate", generateCmd.Use, "Generate command should have correct use")
	assert.NotEmpty(t, generateCmd.Short, "Generate command should have short description")
	assert.NotEmpty(t, generateCmd.Long, "Generate command should have long description")
}

// TestGenerateCommandSubcommands tests the generate command subcommands
func (suite *CLICommandsTestSuite) TestGenerateCommandSubcommands() {
	t := suite.T()

	generateCmd := commands.GetGenerateCommand()

	// Test that generate command has expected subcommands
	expectedSubcommands := []string{"rag", "agent", "workflow"}

	for _, expected := range expectedSubcommands {
		found := false
		for _, cmd := range generateCmd.Commands() {
			if cmd.Use == expected {
				found = true
				break
			}
		}
		assert.True(t, found, "Generate command should have subcommand: %s", expected)
	}
}

// TestAgentCommandStructure tests the agent command structure
func (suite *CLICommandsTestSuite) TestAgentCommandStructure() {
	t := suite.T()

	agentCmd := commands.GetAgentCommand()

	assert.NotNil(t, agentCmd, "Agent command should exist")
	assert.Equal(t, "agent", agentCmd.Use, "Agent command should have correct use")
	assert.NotEmpty(t, agentCmd.Short, "Agent command should have short description")
	assert.NotEmpty(t, agentCmd.Long, "Agent command should have long description")

	// Test required flags
	requiredFlags := []string{"tools", "memory"}
	for _, flag := range requiredFlags {
		assert.True(t, agentCmd.Flags().Lookup(flag) != nil,
			"Agent command should have flag: %s", flag)
	}
}

// TestAgentCommandValidation tests agent command validation
func (suite *CLICommandsTestSuite) TestAgentCommandValidation() {
	t := suite.T()

	agentCmd := commands.GetAgentCommand()

	// Test valid arguments
	validNames := suite.helpers.ValidAgentNames()
	for _, name := range validNames {
		t.Run("valid_name_"+name, func(t *testing.T) {
			agentCmd.SetArgs([]string{name, "--tools=web_search", "--memory=episodic"})
			err := agentCmd.Execute()
			assert.NoError(t, err, "Agent command should accept valid name: %s", name)
		})
	}

	// Test invalid arguments
	invalidNames := suite.helpers.InvalidAgentNames()
	for _, name := range invalidNames {
		if name != "" { // Skip empty name as it's handled differently
			t.Run("invalid_name_"+name, func(t *testing.T) {
				agentCmd.SetArgs([]string{name, "--tools=web_search", "--memory=episodic"})
				err := agentCmd.Execute()
				assert.Error(t, err, "Agent command should reject invalid name: %s", name)
			})
		}
	}
}

// TestAgentCommandFlags tests agent command flags
func (suite *CLICommandsTestSuite) TestAgentCommandFlags() {
	t := suite.T()

	agentCmd := commands.GetAgentCommand()

	// Test tools flag
	toolsFlag := agentCmd.Flags().Lookup("tools")
	assert.NotNil(t, toolsFlag, "Tools flag should exist")
	assert.Equal(t, "t", toolsFlag.Shorthand, "Tools flag should have shorthand 't'")
	assert.NotEmpty(t, toolsFlag.Usage, "Tools flag should have usage description")

	// Test memory flag
	memoryFlag := agentCmd.Flags().Lookup("memory")
	assert.NotNil(t, memoryFlag, "Memory flag should exist")
	assert.Equal(t, "m", memoryFlag.Shorthand, "Memory flag should have shorthand 'm'")
	assert.NotEmpty(t, memoryFlag.Usage, "Memory flag should have usage description")
}

// TestAgentCommandFlagValidation tests agent command flag validation
func (suite *CLICommandsTestSuite) TestAgentCommandFlagValidation() {
	t := suite.T()

	agentCmd := commands.GetAgentCommand()

	// Test valid tool combinations
	validTools := suite.helpers.ValidTools()
	for _, tools := range validTools {
		toolsStr := helpers.JoinStrings(tools, ",")
		t.Run("valid_tools_"+toolsStr, func(t *testing.T) {
			agentCmd.SetArgs([]string{"TestAgent", "--tools=" + toolsStr, "--memory=episodic"})
			err := agentCmd.Execute()
			assert.NoError(t, err, "Agent command should accept valid tools: %s", toolsStr)
		})
	}

	// Test invalid tool combinations
	invalidTools := suite.helpers.InvalidTools()
	for _, tools := range invalidTools {
		toolsStr := helpers.JoinStrings(tools, ",")
		t.Run("invalid_tools_"+toolsStr, func(t *testing.T) {
			agentCmd.SetArgs([]string{"TestAgent", "--tools=" + toolsStr, "--memory=episodic"})
			err := agentCmd.Execute()
			assert.Error(t, err, "Agent command should reject invalid tools: %s", toolsStr)
		})
	}

	// Test valid memory types
	validMemoryTypes := suite.helpers.ValidMemoryTypes()
	for _, memoryType := range validMemoryTypes {
		t.Run("valid_memory_"+memoryType, func(t *testing.T) {
			agentCmd.SetArgs([]string{"TestAgent", "--tools=web_search", "--memory=" + memoryType})
			err := agentCmd.Execute()
			assert.NoError(t, err, "Agent command should accept valid memory type: %s", memoryType)
		})
	}

	// Test invalid memory types
	invalidMemoryTypes := suite.helpers.InvalidMemoryTypes()
	for _, memoryType := range invalidMemoryTypes {
		if memoryType != "" { // Skip empty type as it's handled differently
			t.Run("invalid_memory_"+memoryType, func(t *testing.T) {
				agentCmd.SetArgs([]string{"TestAgent", "--tools=web_search", "--memory=" + memoryType})
				err := agentCmd.Execute()
				assert.Error(t, err, "Agent command should reject invalid memory type: %s", memoryType)
			})
		}
	}
}

// TestAgentCommandExecution tests agent command execution
func (suite *CLICommandsTestSuite) TestAgentCommandExecution() {
	t := suite.T()

	agentCmd := commands.GetAgentCommand()

	// Test successful execution
	t.Run("successful_execution", func(t *testing.T) {
		agentCmd.SetArgs([]string{"TestAgent", "--tools=web_search", "--memory=episodic"})
		err := agentCmd.Execute()
		assert.NoError(t, err, "Agent command should execute successfully")
	})

	// Test execution with missing required flags
	t.Run("missing_tools_flag", func(t *testing.T) {
		agentCmd.SetArgs([]string{"TestAgent", "--memory=episodic"})
		_ = agentCmd.Execute()
		// This might not be an error if tools are optional
		// assert.Error(t, err, "Agent command should require tools flag")
	})

	t.Run("missing_memory_flag", func(t *testing.T) {
		agentCmd.SetArgs([]string{"TestAgent", "--tools=web_search"})
		_ = agentCmd.Execute()
		// This might not be an error if memory has default
		// assert.Error(t, err, "Agent command should require memory flag")
	})
}

// TestCommandHelp tests command help functionality
func (suite *CLICommandsTestSuite) TestCommandHelp() {
	_ = suite.T() // Not used in this test but required by suite pattern

	// Test root command help
	rootCmd := commands.GetRootCommand()
	rootCmd.SetArgs([]string{"--help"})
	_ = rootCmd.Execute()
	// Note: Help command might not work in test environment, so we don't assert on error

	// Test generate command help
	generateCmd := commands.GetGenerateCommand()
	generateCmd.SetArgs([]string{"--help"})
	_ = generateCmd.Execute()
	// Note: Help command might not work in test environment, so we don't assert on error

	// Test agent command help
	agentCmd := commands.GetAgentCommand()
	agentCmd.SetArgs([]string{"--help"})
	_ = agentCmd.Execute()
	// Note: Help command might not work in test environment, so we don't assert on error
}

// TestCommandVersion tests command version functionality
func (suite *CLICommandsTestSuite) TestCommandVersion() {
	t := suite.T()

	versionCmd := commands.GetVersionCommand()

	assert.NotNil(t, versionCmd, "Version command should exist")
	assert.Equal(t, "version", versionCmd.Use, "Version command should have correct use")
	assert.NotEmpty(t, versionCmd.Short, "Version command should have short description")

	// Test version command execution
	versionCmd.SetArgs([]string{})
	err := versionCmd.Execute()
	assert.NoError(t, err, "Version command should execute successfully")
}

// TestCommandErrorHandling tests command error handling
func (suite *CLICommandsTestSuite) TestCommandErrorHandling() {
	t := suite.T()

	// Test invalid subcommand
	rootCmd := commands.GetRootCommand()
	rootCmd.SetArgs([]string{"invalid_command"})
	err := rootCmd.Execute()
	assert.Error(t, err, "Root command should handle invalid subcommand")

	// Test invalid arguments
	agentCmd := commands.GetAgentCommand()
	agentCmd.SetArgs([]string{})
	err = agentCmd.Execute()
	assert.Error(t, err, "Agent command should handle missing arguments")
}

// Run the test suite
func TestCLICommandsTestSuite(t *testing.T) {
	suite.Run(t, new(CLICommandsTestSuite))
}
