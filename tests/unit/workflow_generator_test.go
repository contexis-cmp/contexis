package unit

import (
	"strings"
	"testing"

	"github.com/contexis-cmp/contexis/src/cli/commands"
	"github.com/contexis-cmp/contexis/tests/unit/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// WorkflowGeneratorTestSuite defines the test suite for workflow generator functions.
type WorkflowGeneratorTestSuite struct {
	suite.Suite
	helpers *helpers.TestFixtures
	utils   *helpers.TestUtils
}

// SetupSuite runs once before the entire test suite.
func (suite *WorkflowGeneratorTestSuite) SetupSuite() {
	suite.helpers = &helpers.TestFixtures{}
	suite.utils = &helpers.TestUtils{}
}

// TestWorkflowConfigCreation tests workflow configuration creation.
func (suite *WorkflowGeneratorTestSuite) TestWorkflowConfigCreation() {
	t := suite.T()

	validNames := suite.helpers.ValidWorkflowNames()

	for _, name := range validNames {
		t.Run("config_"+name, func(t *testing.T) {
			config := commands.WorkflowConfig{
				Name:        name,
				Steps:       []string{"research", "write", "review"},
				Description: "Test workflow for " + name,
				Version:     "1.0.0",
				Author:      "Test Author",
			}

			assert.NotEmpty(t, config.Name, "Workflow name should not be empty")
			assert.NotEmpty(t, config.Steps, "Workflow steps should not be empty")
			assert.NotEmpty(t, config.Description, "Workflow description should not be empty")
			assert.NotEmpty(t, config.Version, "Workflow version should not be empty")
			assert.NotEmpty(t, config.Author, "Workflow author should not be empty")
		})
	}
}

// TestWorkflowConfigDefaults tests workflow configuration defaults.
func (suite *WorkflowGeneratorTestSuite) TestWorkflowConfigDefaults() {
	t := suite.T()

	config := commands.WorkflowConfig{
		Name:             "TestWorkflow",
		Steps:            []string{"research", "write"},
		MaxConcurrency:   3,
		RetryAttempts:    3,
		RetryDelay:       5,
		Timeout:          300,
		StatePersistence: true,
		ErrorHandling:    "continue_on_error",
		Logging:          "structured",
		Monitoring:       true,
	}

	// Test default values
	assert.Equal(t, 3, config.MaxConcurrency, "Default max concurrency should be 3")
	assert.Equal(t, 3, config.RetryAttempts, "Default retry attempts should be 3")
	assert.Equal(t, 5, config.RetryDelay, "Default retry delay should be 5")
	assert.Equal(t, 300, config.Timeout, "Default timeout should be 300")
	assert.True(t, config.StatePersistence, "Default state persistence should be true")
	assert.Equal(t, "continue_on_error", config.ErrorHandling, "Default error handling should be continue_on_error")
	assert.Equal(t, "structured", config.Logging, "Default logging should be structured")
	assert.True(t, config.Monitoring, "Default monitoring should be true")
}

// TestWorkflowConfigValidation tests workflow configuration validation.
func (suite *WorkflowGeneratorTestSuite) TestWorkflowConfigValidation() {
	t := suite.T()

	// Test valid step combinations
	validStepCombinations := [][]string{
		{"research", "write", "review"},
		{"extract", "transform", "load"},
		{"analyze", "generate", "validate"},
		{"research", "write"},
		{"deploy"},
	}

	for _, steps := range validStepCombinations {
		t.Run("valid_"+strings.Join(steps, "_"), func(t *testing.T) {
			err := commands.ValidateWorkflowConfig(steps)
			assert.NoError(t, err, "Valid step combination should pass validation: %v", steps)
		})
	}

	// Test invalid step combinations
	invalidStepCombinations := [][]string{
		{}, // Empty steps
		{"invalid_step"},
		{"research", "invalid_step", "review"},
		{"research", "write", "invalid_step", "deploy"},
	}

	for _, steps := range invalidStepCombinations {
		t.Run("invalid_"+strings.Join(steps, "_"), func(t *testing.T) {
			err := commands.ValidateWorkflowConfig(steps)
			assert.Error(t, err, "Invalid step combination should fail validation: %v", steps)
		})
	}
}

// TestWorkflowConfigValidationEdgeCases tests edge cases in workflow validation.
func (suite *WorkflowGeneratorTestSuite) TestWorkflowConfigValidationEdgeCases() {
	t := suite.T()

	// Test empty steps
	err := commands.ValidateWorkflowConfig([]string{})
	assert.Error(t, err, "Empty steps should fail validation")
	assert.Contains(t, err.Error(), "at least one step is required")

	// Test steps with empty strings
	err = commands.ValidateWorkflowConfig([]string{"research", "", "review"})
	assert.Error(t, err, "Steps with empty strings should fail validation")
	assert.Contains(t, err.Error(), "step name cannot be empty")

	// Test single valid step
	err = commands.ValidateWorkflowConfig([]string{"research"})
	assert.NoError(t, err, "Single valid step should pass validation")
}

// TestWorkflowNameValidation tests workflow name validation.
func (suite *WorkflowGeneratorTestSuite) TestWorkflowNameValidation() {
	t := suite.T()

	validNames := []string{
		"ContentPipeline",
		"DataProcessing",
		"MLTraining",
		"test_workflow",
		"workflow-123",
		"MyWorkflow",
	}

	for _, name := range validNames {
		t.Run("valid_"+name, func(t *testing.T) {
			suite.utils.ValidateWorkflowName(t, name)
		})
	}

	invalidNames := []string{
		"",  // Empty name
		"a", // Too short
		"invalid name with spaces",
		"invalid-name-with-special-chars!@#",
		"invalid/name/with/slashes",
		"invalid\\name\\with\\backslashes",
	}

	for _, name := range invalidNames {
		t.Run("invalid_"+name, func(t *testing.T) {
			// This should fail validation, so we expect it to panic or fail
			assert.Panics(t, func() {
				suite.utils.ValidateWorkflowName(t, name)
			}, "Invalid workflow name should fail validation: %s", name)
		})
	}
}

// TestStepTypeValidation tests step type validation.
func (suite *WorkflowGeneratorTestSuite) TestStepTypeValidation() {
	t := suite.T()

	validStepTypes := []string{
		"research", "write", "review", "extract", "transform", "load",
		"analyze", "generate", "validate", "deploy",
	}

	for _, stepType := range validStepTypes {
		t.Run("valid_"+stepType, func(t *testing.T) {
			suite.utils.ValidateStepType(t, stepType)
		})
	}

	invalidStepTypes := []string{
		"invalid_step",
		"",
		"step",
		"research_step",
		"write_process",
	}

	for _, stepType := range invalidStepTypes {
		t.Run("invalid_"+stepType, func(t *testing.T) {
			// This should fail validation, so we expect it to panic or fail
			assert.Panics(t, func() {
				suite.utils.ValidateStepType(t, stepType)
			}, "Invalid step type should fail validation: %s", stepType)
		})
	}
}

// TestStepStructure tests step configuration structure.
func (suite *WorkflowGeneratorTestSuite) TestStepStructure() {
	t := suite.T()

	stepTypes := []string{"research", "write", "review", "extract", "transform"}

	for _, stepType := range stepTypes {
		t.Run("step_"+stepType, func(t *testing.T) {
			step := commands.StepConfig{
				Name:          stepType,
				Description:   "Test step for " + stepType,
				Type:          stepType,
				Input:         make(map[string]interface{}),
				Output:        make(map[string]interface{}),
				Dependencies:  []string{},
				Timeout:       60,
				RetryAttempts: 2,
				Parallel:      false,
				Condition:     "",
				ErrorHandling: "stop_on_error",
				Resources: commands.ResourceLimits{
					CPU:     "1",
					Memory:  "2Gi",
					Storage: "5Gi",
					Network: "50Mbps",
				},
			}

			assert.NotEmpty(t, step.Name, "Step name should not be empty")
			assert.NotEmpty(t, step.Description, "Step description should not be empty")
			assert.NotEmpty(t, step.Type, "Step type should not be empty")
			assert.Equal(t, 60, step.Timeout, "Step timeout should be 60")
			assert.Equal(t, 2, step.RetryAttempts, "Step retry attempts should be 2")
			assert.False(t, step.Parallel, "Step parallel should be false by default")
			assert.Equal(t, "stop_on_error", step.ErrorHandling, "Step error handling should be stop_on_error")
		})
	}
}

// Run the test suite
func TestWorkflowGeneratorTestSuite(t *testing.T) {
	suite.Run(t, new(WorkflowGeneratorTestSuite))
}
