package commands

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/contexis-cmp/contexis/src/cli/logger"
	"go.uber.org/zap"
)

// WorkflowConfig holds configuration for workflow generation
type WorkflowConfig struct {
	Name             string
	Steps            []string
	Description      string
	Version          string
	Author           string
	CreatedDate      string
	LastModified     string
	MaxConcurrency   int
	RetryAttempts    int
	RetryDelay       int
	Timeout          int
	StatePersistence bool
	ErrorHandling    string
	Logging          string
	Monitoring       bool
	ResourceLimits   ResourceLimits
	StepsConfig      []StepConfig
}

// StepConfig holds configuration for individual workflow steps
type StepConfig struct {
	Name          string
	Description   string
	Type          string
	Input         map[string]interface{}
	Output        map[string]interface{}
	Dependencies  []string
	Timeout       int
	RetryAttempts int
	Parallel      bool
	Condition     string
	ErrorHandling string
	Resources     ResourceLimits
}

// ResourceLimits holds resource constraints for workflows and steps
type ResourceLimits struct {
	CPU     string
	Memory  string
	Storage string
	Network string
}

// WorkflowStep represents a step in the workflow
type WorkflowStep struct {
	Name         string
	Type         string
	Description  string
	Dependencies []string
}

// GenerateWorkflow creates a multi-step AI processing pipeline
func GenerateWorkflow(ctx context.Context, name, steps string) error {
	log := logger.WithContext(ctx)

	// Parse steps string
	stepList := []string{}
	if steps != "" {
		stepList = strings.Split(steps, ",")
		for i, step := range stepList {
			stepList[i] = strings.TrimSpace(step)
		}
	}

	// Validate configuration
	if err := ValidateWorkflowConfig(stepList); err != nil {
		log.Error("workflow configuration validation failed", zap.Error(err))
		return fmt.Errorf("invalid workflow configuration: %w", err)
	}

	config := WorkflowConfig{
		Name:             name,
		Steps:            stepList,
		Description:      fmt.Sprintf("Multi-step AI processing pipeline for %s", name),
		Version:          "1.0.0",
		Author:           "CMP Framework",
		CreatedDate:      time.Now().Format("2006-01-02"),
		LastModified:     time.Now().Format("2006-01-02"),
		MaxConcurrency:   3,
		RetryAttempts:    3,
		RetryDelay:       5,
		Timeout:          300,
		StatePersistence: true,
		ErrorHandling:    "continue_on_error",
		Logging:          "structured",
		Monitoring:       true,
		ResourceLimits: ResourceLimits{
			CPU:     "2",
			Memory:  "4Gi",
			Storage: "10Gi",
			Network: "100Mbps",
		},
		StepsConfig: generateStepConfigs(stepList),
	}

	log.Info("generating workflow",
		zap.String("name", name),
		zap.Strings("steps", stepList))

	// Create workflow-specific directory structure
	workflowDirs := []string{
		fmt.Sprintf("workflows/%s", name),
		fmt.Sprintf("contexts/%s", name),
		fmt.Sprintf("prompts/%s", name),
		fmt.Sprintf("prompts/%s/step_templates", name),
		fmt.Sprintf("memory/%s", name),
		fmt.Sprintf("memory/%s/workflow_state", name),
		fmt.Sprintf("tests/%s", name),
	}

	for _, dir := range workflowDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Error("failed to create workflow directory", zap.String("dir", dir), zap.Error(err))
			return fmt.Errorf("failed to create workflow directory %s: %w", dir, err)
		}
	}

	// Generate workflow components
	if err := generateWorkflowDefinition(ctx, config); err != nil {
		log.Error("failed to generate workflow definition", zap.Error(err))
		return fmt.Errorf("failed to generate workflow definition: %w", err)
	}

	if err := generateWorkflowContext(ctx, config); err != nil {
		log.Error("failed to generate workflow context", zap.Error(err))
		return fmt.Errorf("failed to generate workflow context: %w", err)
	}

	if err := generateWorkflowPrompts(ctx, config); err != nil {
		log.Error("failed to generate workflow prompts", zap.Error(err))
		return fmt.Errorf("failed to generate workflow prompts: %w", err)
	}

	if err := generateWorkflowMemory(ctx, config); err != nil {
		log.Error("failed to generate workflow memory", zap.Error(err))
		return fmt.Errorf("failed to generate workflow memory: %w", err)
	}

	if err := generateWorkflowTests(ctx, config); err != nil {
		log.Error("failed to generate workflow tests", zap.Error(err))
		return fmt.Errorf("failed to generate workflow tests: %w", err)
	}

	log.Info("workflow generation completed successfully",
		zap.String("name", name),
		zap.Strings("steps", stepList))

	return nil
}

// ValidateWorkflowConfig validates workflow configuration parameters
func ValidateWorkflowConfig(steps []string) error {
    // Validate steps
    if len(steps) == 0 {
        return fmt.Errorf("at least one step is required")
    }

	// Validate individual steps
	validStepTypes := []string{"research", "write", "review", "extract", "transform", "load", "analyze", "generate", "validate", "deploy"}

	for _, step := range steps {
        if step == "" {
            return fmt.Errorf("step name cannot be empty")
        }

		// Check if step type is valid
		valid := false
		for _, validType := range validStepTypes {
			if step == validType {
				valid = true
				break
			}
		}

        if !valid {
            return fmt.Errorf("invalid step type '%s'. Valid types: %s", step, strings.Join(validStepTypes, ", "))
        }
	}

	return nil
}

// generateStepConfigs generates step configurations based on step names
func generateStepConfigs(steps []string) []StepConfig {
	var stepConfigs []StepConfig

	for i, step := range steps {
		config := StepConfig{
			Name:          step,
			Description:   fmt.Sprintf("Step %d: %s", i+1, step),
			Type:          step,
			Input:         make(map[string]interface{}),
			Output:        make(map[string]interface{}),
			Dependencies:  []string{},
			Timeout:       60,
			RetryAttempts: 2,
			Parallel:      false,
			Condition:     "",
			ErrorHandling: "stop_on_error",
			Resources: ResourceLimits{
				CPU:     "1",
				Memory:  "2Gi",
				Storage: "5Gi",
				Network: "50Mbps",
			},
		}

		// Add dependencies for sequential steps
		if i > 0 {
			config.Dependencies = []string{steps[i-1]}
		}

		stepConfigs = append(stepConfigs, config)
	}

	return stepConfigs
}

// generateWorkflowDefinition creates the main workflow definition file
func generateWorkflowDefinition(ctx context.Context, config WorkflowConfig) error {
	log := logger.WithContext(ctx)

	workflowPath := fmt.Sprintf("workflows/%s/%s.yaml", config.Name, config.Name)

	// Read template
	templatePath := "templates/workflow/workflow_definition.yaml"
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		log.Error("failed to read workflow template", zap.String("path", templatePath), zap.Error(err))
		return fmt.Errorf("failed to read workflow template: %w", err)
	}

	// Parse template
	tmpl, err := template.New("workflow_definition").Parse(string(templateContent))
	if err != nil {
		log.Error("failed to parse workflow template", zap.Error(err))
		return fmt.Errorf("failed to parse workflow template: %w", err)
	}

	// Create file
	file, err := os.Create(workflowPath)
	if err != nil {
		log.Error("failed to create workflow file", zap.String("path", workflowPath), zap.Error(err))
		return fmt.Errorf("failed to create workflow file: %w", err)
	}
	defer file.Close()

	// Execute template
	if err := tmpl.Execute(file, config); err != nil {
		log.Error("failed to execute workflow template", zap.Error(err))
		return fmt.Errorf("failed to execute workflow template: %w", err)
	}

	log.Info("workflow definition generated", zap.String("path", workflowPath))
	return nil
}

// generateWorkflowContext creates the workflow coordinator context
func generateWorkflowContext(ctx context.Context, config WorkflowConfig) error {
	log := logger.WithContext(ctx)

	contextPath := fmt.Sprintf("contexts/%s/workflow_coordinator.ctx", config.Name)

	// Read template
	templatePath := "templates/workflow/workflow_coordinator.ctx"
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		log.Error("failed to read workflow context template", zap.String("path", templatePath), zap.Error(err))
		return fmt.Errorf("failed to read workflow context template: %w", err)
	}

	// Parse template
	tmpl, err := template.New("workflow_coordinator").Parse(string(templateContent))
	if err != nil {
		log.Error("failed to parse workflow context template", zap.Error(err))
		return fmt.Errorf("failed to parse workflow context template: %w", err)
	}

	// Create file
	file, err := os.Create(contextPath)
	if err != nil {
		log.Error("failed to create workflow context file", zap.String("path", contextPath), zap.Error(err))
		return fmt.Errorf("failed to create workflow context file: %w", err)
	}
	defer file.Close()

	// Execute template
	if err := tmpl.Execute(file, config); err != nil {
		log.Error("failed to execute workflow context template", zap.Error(err))
		return fmt.Errorf("failed to execute workflow context template: %w", err)
	}

	log.Info("workflow context generated", zap.String("path", contextPath))
	return nil
}

// generateWorkflowPrompts creates step-specific prompt templates
func generateWorkflowPrompts(ctx context.Context, config WorkflowConfig) error {
	log := logger.WithContext(ctx)

	for _, step := range config.StepsConfig {
		promptPath := fmt.Sprintf("prompts/%s/step_templates/%s.md", config.Name, step.Name)

		// Read template
		templatePath := fmt.Sprintf("templates/workflow/step_%s.md", step.Type)
		templateContent, err := os.ReadFile(templatePath)
		if err != nil {
			log.Error("failed to read step template", zap.String("path", templatePath), zap.Error(err))
			return fmt.Errorf("failed to read step template: %w", err)
		}

		// Parse template
		tmpl, err := template.New(fmt.Sprintf("step_%s", step.Name)).Parse(string(templateContent))
		if err != nil {
			log.Error("failed to parse step template", zap.Error(err))
			return fmt.Errorf("failed to parse step template: %w", err)
		}

		// Create file
		file, err := os.Create(promptPath)
		if err != nil {
			log.Error("failed to create step prompt file", zap.String("path", promptPath), zap.Error(err))
			return fmt.Errorf("failed to create step prompt file: %w", err)
		}
		defer file.Close()

		// Execute template
		if err := tmpl.Execute(file, step); err != nil {
			log.Error("failed to execute step template", zap.Error(err))
			return fmt.Errorf("failed to execute step template: %w", err)
		}

		log.Info("step prompt generated", zap.String("path", promptPath))
	}

	return nil
}

// generateWorkflowMemory creates workflow state management configuration
func generateWorkflowMemory(ctx context.Context, config WorkflowConfig) error {
	log := logger.WithContext(ctx)

	memoryPath := fmt.Sprintf("memory/%s/workflow_state.yaml", config.Name)

	// Read template
	templatePath := "templates/workflow/workflow_state.yaml"
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		log.Error("failed to read workflow state template", zap.String("path", templatePath), zap.Error(err))
		return fmt.Errorf("failed to read workflow state template: %w", err)
	}

	// Parse template
	tmpl, err := template.New("workflow_state").Parse(string(templateContent))
	if err != nil {
		log.Error("failed to parse workflow state template", zap.Error(err))
		return fmt.Errorf("failed to parse workflow state template: %w", err)
	}

	// Create file
	file, err := os.Create(memoryPath)
	if err != nil {
		log.Error("failed to create workflow state file", zap.String("path", memoryPath), zap.Error(err))
		return fmt.Errorf("failed to create workflow state file: %w", err)
	}
	defer file.Close()

	// Execute template
	if err := tmpl.Execute(file, config); err != nil {
		log.Error("failed to execute workflow state template", zap.Error(err))
		return fmt.Errorf("failed to execute workflow state template: %w", err)
	}

	log.Info("workflow state configuration generated", zap.String("path", memoryPath))
	return nil
}

// generateWorkflowTests creates workflow integration tests
func generateWorkflowTests(ctx context.Context, config WorkflowConfig) error {
	log := logger.WithContext(ctx)

	testPath := fmt.Sprintf("tests/%s/workflow_integration.py", config.Name)

	// Read template
	templatePath := "templates/workflow/workflow_integration.py"
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		log.Error("failed to read workflow test template", zap.String("path", templatePath), zap.Error(err))
		return fmt.Errorf("failed to read workflow test template: %w", err)
	}

	// Parse template
	tmpl, err := template.New("workflow_integration").Parse(string(templateContent))
	if err != nil {
		log.Error("failed to parse workflow test template", zap.Error(err))
		return fmt.Errorf("failed to parse workflow test template: %w", err)
	}

	// Create file
	file, err := os.Create(testPath)
	if err != nil {
		log.Error("failed to create workflow test file", zap.String("path", testPath), zap.Error(err))
		return fmt.Errorf("failed to create workflow test file: %w", err)
	}
	defer file.Close()

	// Execute template
	if err := tmpl.Execute(file, config); err != nil {
		log.Error("failed to execute workflow test template", zap.Error(err))
		return fmt.Errorf("failed to execute workflow test template: %w", err)
	}

	log.Info("workflow integration test generated", zap.String("path", testPath))
	return nil
}
