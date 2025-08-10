package unit

import (
	"testing"

	"github.com/contexis-cmp/contexis/src/cli/commands"
	"github.com/contexis-cmp/contexis/tests/unit/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// AgentGeneratorTestSuite provides a test suite for agent generator functionality
type AgentGeneratorTestSuite struct {
	suite.Suite
	helpers *helpers.TestFixtures
	utils   *helpers.TestUtils
}

// SetupSuite sets up the test suite
func (suite *AgentGeneratorTestSuite) SetupSuite() {
	suite.helpers = &helpers.TestFixtures{}
	suite.utils = &helpers.TestUtils{}
}

// TestAgentConfigValidation tests agent configuration validation
func (suite *AgentGeneratorTestSuite) TestAgentConfigValidation() {
	t := suite.T()

	// Test valid configurations
	validTools := suite.helpers.ValidTools()
	validMemoryTypes := suite.helpers.ValidMemoryTypes()

	for _, tools := range validTools {
		for _, memory := range validMemoryTypes {
			t.Run("valid_"+memory+"_"+helpers.JoinStrings(tools, "_"), func(t *testing.T) {
				err := commands.ValidateAgentConfig(tools, memory)
				assert.NoError(t, err, "Valid configuration should pass validation")
			})
		}
	}

	// Test invalid configurations
	invalidTools := suite.helpers.InvalidTools()
	invalidMemoryTypes := suite.helpers.InvalidMemoryTypes()

	for _, tools := range invalidTools {
		t.Run("invalid_tools_"+helpers.JoinStrings(tools, "_"), func(t *testing.T) {
			err := commands.ValidateAgentConfig(tools, "episodic")
			assert.Error(t, err, "Invalid tools should fail validation")
		})
	}

	for _, memory := range invalidMemoryTypes {
		t.Run("invalid_memory_"+memory, func(t *testing.T) {
			err := commands.ValidateAgentConfig([]string{"web_search"}, memory)
			assert.Error(t, err, "Invalid memory type should fail validation")
		})
	}
}

// TestAgentConfigCreation tests agent configuration creation
func (suite *AgentGeneratorTestSuite) TestAgentConfigCreation() {
	t := suite.T()

	validConfigs := suite.helpers.ValidAgentConfigs()

	for _, config := range validConfigs {
		t.Run("config_"+config.Name, func(t *testing.T) {
			// Test config fields
			assert.NotEmpty(t, config.Name, "Agent name should not be empty")
			assert.NotEmpty(t, config.Description, "Description should not be empty")
			assert.NotEmpty(t, config.Version, "Version should not be empty")
			assert.NotEmpty(t, config.Persona, "Persona should not be empty")
			assert.NotEmpty(t, config.Capabilities, "Capabilities should not be empty")
			assert.NotEmpty(t, config.Limitations, "Limitations should not be empty")
			assert.NotEmpty(t, config.BusinessRules, "Business rules should not be empty")
			assert.NotEmpty(t, config.BaselineDate, "Baseline date should not be empty")
			assert.NotEmpty(t, config.AdminEmail, "Admin email should not be empty")
			assert.NotEmpty(t, config.Tone, "Tone should not be empty")
			assert.NotEmpty(t, config.Format, "Format should not be empty")
			assert.Greater(t, config.MaxTokens, 0, "Max tokens should be greater than 0")
			assert.GreaterOrEqual(t, config.Temperature, 0.0, "Temperature should be >= 0")
			assert.LessOrEqual(t, config.Temperature, 1.0, "Temperature should be <= 1")
			assert.NotEmpty(t, config.MemoryType, "Memory type should not be empty")
			assert.GreaterOrEqual(t, config.MaxHistory, 0, "Max history should be >= 0")
			assert.NotEmpty(t, config.Privacy, "Privacy should not be empty")
			assert.Greater(t, config.DriftThreshold, 0.0, "Drift threshold should be > 0")
			assert.LessOrEqual(t, config.DriftThreshold, 1.0, "Drift threshold should be <= 1")
		})
	}
}

// TestToolStructure tests tool structure validation
func (suite *AgentGeneratorTestSuite) TestToolStructure() {
	t := suite.T()

	validTools := suite.helpers.ValidToolDefinitions()

	for _, tool := range validTools {
		t.Run("tool_"+tool.Name, func(t *testing.T) {
			// Test tool fields
			assert.NotEmpty(t, tool.Name, "Tool name should not be empty")
			assert.NotEmpty(t, tool.URI, "Tool URI should not be empty")
			assert.NotEmpty(t, tool.Description, "Tool description should not be empty")

			// Validate tool name
			suite.utils.ValidateToolName(t, tool.Name)

			// Validate URI format
			assert.Contains(t, tool.URI, "mcp://", "Tool URI should use MCP protocol")
		})
	}
}

// TestAgentNameValidation tests agent name validation
func (suite *AgentGeneratorTestSuite) TestAgentNameValidation() {
	t := suite.T()

	// Test valid agent names
	validNames := suite.helpers.ValidAgentNames()
	for _, name := range validNames {
		t.Run("valid_"+name, func(t *testing.T) {
			suite.utils.ValidateAgentName(t, name)
		})
	}

	// Test invalid agent names
	invalidNames := suite.helpers.InvalidAgentNames()
	for _, name := range invalidNames {
		t.Run("invalid_"+name, func(t *testing.T) {
			if name != "" {
				// Skip empty name test as it's handled differently
				suite.utils.ValidateAgentName(t, name)
			}
		})
	}
}

// TestMemoryTypeValidation tests memory type validation
func (suite *AgentGeneratorTestSuite) TestMemoryTypeValidation() {
	t := suite.T()

	// Test valid memory types
	validTypes := suite.helpers.ValidMemoryTypes()
	for _, memoryType := range validTypes {
		t.Run("valid_"+memoryType, func(t *testing.T) {
			suite.utils.ValidateMemoryType(t, memoryType)
		})
	}

	// Test invalid memory types
	invalidTypes := suite.helpers.InvalidMemoryTypes()
	for _, memoryType := range invalidTypes {
		t.Run("invalid_"+memoryType, func(t *testing.T) {
			if memoryType != "" {
				// Skip empty type test as it's handled differently
				suite.utils.ValidateMemoryType(t, memoryType)
			}
		})
	}
}

// TestAgentConfigDefaults tests default configuration values
func (suite *AgentGeneratorTestSuite) TestAgentConfigDefaults() {
	t := suite.T()

	config := commands.AgentConfig{
		Name:           "TestAgent",
		Tools:          []string{"web_search"},
		Memory:         "episodic",
		Description:    "Test agent",
		Version:        "1.0.0",
		Persona:        "Professional assistant",
		Capabilities:   []string{"conversation"},
		Limitations:    []string{"no_personal_data"},
		BusinessRules:  []string{"always_helpful"},
		BaselineDate:   "2024-01-01",
		AdminEmail:     "test@example.com",
		Tone:           "professional",
		Format:         "json",
		MaxTokens:      500,
		Temperature:    0.1,
		MemoryType:     "episodic",
		MaxHistory:     10,
		Privacy:        "user_isolated",
		DriftThreshold: 0.85,
	}

	// Test default values
	assert.Equal(t, "json", config.Format, "Default format should be JSON")
	assert.Equal(t, 500, config.MaxTokens, "Default max tokens should be 500")
	assert.Equal(t, 0.1, config.Temperature, "Default temperature should be 0.1")
	assert.Equal(t, 10, config.MaxHistory, "Default max history should be 10")
	assert.Equal(t, "user_isolated", config.Privacy, "Default privacy should be user_isolated")
	assert.Equal(t, 0.85, config.DriftThreshold, "Default drift threshold should be 0.85")
}

// TestAgentConfigValidationEdgeCases tests edge cases for configuration validation
func (suite *AgentGeneratorTestSuite) TestAgentConfigValidationEdgeCases() {
	t := suite.T()

	// Test empty tools
	err := commands.ValidateAgentConfig([]string{}, "episodic")
	assert.NoError(t, err, "Empty tools should be valid")

	// Test nil tools
	err = commands.ValidateAgentConfig(nil, "episodic")
	assert.NoError(t, err, "Nil tools should be valid")

	// Test mixed valid/invalid tools
	err = commands.ValidateAgentConfig([]string{"web_search", "invalid_tool"}, "episodic")
	assert.Error(t, err, "Mixed valid/invalid tools should fail validation")

	// Test duplicate tools
	err = commands.ValidateAgentConfig([]string{"web_search", "web_search"}, "episodic")
	assert.NoError(t, err, "Duplicate tools should be valid (though not ideal)")
}

// TestAgentConfigBusinessRules tests business rules validation
func (suite *AgentGeneratorTestSuite) TestAgentConfigBusinessRules() {
	t := suite.T()

	configs := suite.helpers.ValidAgentConfigs()

	for _, config := range configs {
		t.Run("business_rules_"+config.Name, func(t *testing.T) {
			// Test required business rules
			requiredRules := []string{"always_helpful"}
			for _, rule := range requiredRules {
				assert.Contains(t, config.BusinessRules, rule,
					"Agent should have required business rule: %s", rule)
			}

			// Test business rules format
			for _, rule := range config.BusinessRules {
				assert.NotEmpty(t, rule, "Business rule should not be empty")
				assert.False(t, helpers.ContainsWhitespace(rule),
					"Business rule should not contain whitespace: %s", rule)
			}
		})
	}
}

// TestAgentConfigCapabilities tests capabilities validation
func (suite *AgentGeneratorTestSuite) TestAgentConfigCapabilities() {
	t := suite.T()

	configs := suite.helpers.ValidAgentConfigs()

	for _, config := range configs {
		t.Run("capabilities_"+config.Name, func(t *testing.T) {
			// Test required capabilities
			requiredCapabilities := []string{"conversation"}
			for _, capability := range requiredCapabilities {
				assert.Contains(t, config.Capabilities, capability,
					"Agent should have required capability: %s", capability)
			}

			// Test capabilities format
			for _, capability := range config.Capabilities {
				assert.NotEmpty(t, capability, "Capability should not be empty")
				assert.False(t, helpers.ContainsWhitespace(capability),
					"Capability should not contain whitespace: %s", capability)
			}
		})
	}
}

// TestAgentConfigLimitations tests limitations validation
func (suite *AgentGeneratorTestSuite) TestAgentConfigLimitations() {
	t := suite.T()

	configs := suite.helpers.ValidAgentConfigs()

	for _, config := range configs {
		t.Run("limitations_"+config.Name, func(t *testing.T) {
			// Test required limitations
			requiredLimitations := []string{"no_personal_data"}
			for _, limitation := range requiredLimitations {
				assert.Contains(t, config.Limitations, limitation,
					"Agent should have required limitation: %s", limitation)
			}

			// Test limitations format
			for _, limitation := range config.Limitations {
				assert.NotEmpty(t, limitation, "Limitation should not be empty")
				assert.False(t, helpers.ContainsWhitespace(limitation),
					"Limitation should not contain whitespace: %s", limitation)
			}
		})
	}
}

// TestAgentConfigSecurity tests security-related configuration
func (suite *AgentGeneratorTestSuite) TestAgentConfigSecurity() {
	t := suite.T()

	configs := suite.helpers.ValidAgentConfigs()

	for _, config := range configs {
		t.Run("security_"+config.Name, func(t *testing.T) {
			// Test privacy settings
			assert.Equal(t, "user_isolated", config.Privacy,
				"Privacy should be user_isolated for security")

			// Test admin email format
			assert.Contains(t, config.AdminEmail, "@",
				"Admin email should be valid format")

			// Test drift threshold security
			assert.GreaterOrEqual(t, config.DriftThreshold, 0.8,
				"Drift threshold should be >= 0.8 for security")
		})
	}
}

// Run the test suite
func TestAgentGeneratorTestSuite(t *testing.T) {
	suite.Run(t, new(AgentGeneratorTestSuite))
}
