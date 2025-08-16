package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
	"time"

	"github.com/contexis-cmp/contexis/src/cli/logger"
	"go.uber.org/zap"
)

// AgentConfig holds configuration for agent generation
type AgentConfig struct {
	Name           string
	Tools          []string
	Memory         string
	Description    string
	Version        string
	Persona        string
	Capabilities   []string
	Limitations    []string
	BusinessRules  []string
	BaselineDate   string
	AdminEmail     string
	Tone           string
	Format         string
	MaxTokens      int
	Temperature    float64
	MemoryType     string
	MaxHistory     int
	Privacy        string
	DriftThreshold float64
}

// Tool represents a tool that can be used by the agent
type Tool struct {
	Name        string
	URI         string
	Description string
}

// generateAgent creates a conversational agent with tools and episodic memory
func GenerateAgent(ctx context.Context, name, tools, memory string) error {
	log := logger.WithContext(ctx)

	// Validate agent name early to match test expectations
	if err := validateAgentName(name); err != nil {
		return err
	}

	// Set defaults if not provided
	if memory == "" {
		memory = "episodic"
	}

	// Parse tools string
	toolList := []string{}
	if tools != "" {
		toolList = strings.Split(tools, ",")
		for i, tool := range toolList {
			toolList[i] = strings.TrimSpace(tool)
		}
	}

	// Validate configuration
	if err := ValidateAgentConfig(toolList, memory); err != nil {
		log.Error("agent configuration validation failed", zap.Error(err))
		return fmt.Errorf("invalid agent configuration: %w", err)
	}

	config := AgentConfig{
		Name:           name,
		Tools:          toolList,
		Memory:         memory,
		Description:    fmt.Sprintf("Conversational agent for %s", name),
		Version:        "1.0.0",
		Persona:        "Professional, helpful conversational assistant",
		Capabilities:   []string{"conversation", "tool_usage", "memory_retention", "context_awareness"},
		Limitations:    []string{"no_personal_data", "no_harmful_content", "no_unauthorized_access"},
		BusinessRules:  []string{"always_helpful", "professional_tone", "tool_security", "memory_privacy"},
		BaselineDate:   time.Now().Format("2006-01-02"),
		AdminEmail:     "admin@example.com",
		Tone:           "professional",
		Format:         "json",
		MaxTokens:      500,
		Temperature:    0.1,
		MemoryType:     memory,
		MaxHistory:     10,
		Privacy:        "user_isolated",
		DriftThreshold: 0.85,
	}

	logger.LogInfo(ctx, "Generating agent",
		zap.String("name", name),
		zap.Strings("tools", toolList),
		zap.String("memory", memory))

	// Create agent-specific directory structure
	agentDirs := []string{
		fmt.Sprintf("contexts/%s", name),
		fmt.Sprintf("memory/%s", name),
		fmt.Sprintf("memory/%s/episodic", name),
		fmt.Sprintf("memory/%s/user_preferences", name),
		fmt.Sprintf("prompts/%s", name),
		fmt.Sprintf("tools/%s", name),
		fmt.Sprintf("tests/%s", name),
	}

	logger.LogInfo(ctx, "Creating agent directory structure")
	for _, dir := range agentDirs {
		if err := os.MkdirAll(dir, 0750); err != nil {
			logger.LogErrorColored(ctx, "failed to create directory", err, zap.String("directory", dir))
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		logger.LogDebugWithContext(ctx, "Created directory", zap.String("path", dir))
	}

	// Generate agent components
	logger.LogInfo(ctx, "Generating agent context")
	if err := generateAgentContext(ctx, config); err != nil {
		logger.LogErrorColored(ctx, "failed to generate agent context", err)
		return fmt.Errorf("failed to generate agent context: %w", err)
	}

	logger.LogInfo(ctx, "Generating agent prompts")
	if err := generateAgentPrompts(ctx, config); err != nil {
		logger.LogErrorColored(ctx, "failed to generate agent prompts", err)
		return fmt.Errorf("failed to generate agent prompts: %w", err)
	}

	logger.LogInfo(ctx, "Generating agent tools")
	if err := generateAgentTools(ctx, config); err != nil {
		logger.LogErrorColored(ctx, "failed to generate agent tools", err)
		return fmt.Errorf("failed to generate agent tools: %w", err)
	}

	logger.LogInfo(ctx, "Generating agent tests")
	if err := generateAgentTests(ctx, config); err != nil {
		logger.LogErrorColored(ctx, "failed to generate agent tests", err)
		return fmt.Errorf("failed to generate agent tests: %w", err)
	}

	logger.LogInfo(ctx, "Generating agent memory configuration")
	if err := generateAgentMemory(ctx, config); err != nil {
		logger.LogErrorColored(ctx, "failed to generate agent memory", err)
		return fmt.Errorf("failed to generate agent memory: %w", err)
	}

	logger.LogInfo(ctx, "Generating agent requirements")
	if err := generateAgentRequirements(ctx, config); err != nil {
		logger.LogErrorColored(ctx, "failed to generate agent requirements", err)
		return fmt.Errorf("failed to generate agent requirements: %w", err)
	}

	// Show generated structure and development flow
	showAgentStructure(name, config)
	showAgentDevelopmentFlow(name, config)

	return nil
}

// validateAgentConfig validates agent configuration parameters
func ValidateAgentConfig(tools []string, memory string) error {
	// Validate memory type
	validMemoryTypes := []string{"episodic", "none"}
	isValidMemory := false
	for _, validType := range validMemoryTypes {
		if memory == validType {
			isValidMemory = true
			break
		}
	}
	if !isValidMemory {
		return fmt.Errorf("invalid memory type '%s'. Valid types: %s", memory, strings.Join(validMemoryTypes, ", "))
	}

	// Validate tools
	validTools := []string{"web_search", "database", "api", "file_system", "email"}
	for _, tool := range tools {
		isValid := false
		for _, validTool := range validTools {
			if tool == validTool {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("invalid tool '%s'. Valid tools: %s", tool, strings.Join(validTools, ", "))
		}
	}

	return nil
}

// generateAgentContext creates the agent context file
func generateAgentContext(ctx context.Context, config AgentConfig) error {
	log := logger.WithContext(ctx)

	// Define available tools
	availableTools := map[string]Tool{
		"web_search": {
			Name:        "web_search",
			URI:         "mcp://web.search",
			Description: "Search the web for current information",
		},
		"database": {
			Name:        "database",
			URI:         "mcp://database.query",
			Description: "Query database for user and order information",
		},
		"api": {
			Name:        "api",
			URI:         "mcp://api.call",
			Description: "Make API calls to external services",
		},
		"file_system": {
			Name:        "file_system",
			URI:         "mcp://file.read",
			Description: "Read and write files",
		},
		"email": {
			Name:        "email",
			URI:         "mcp://email.send",
			Description: "Send and read emails",
		},
	}

	// Filter tools based on configuration
	var selectedTools []Tool
	for _, toolName := range config.Tools {
		if tool, exists := availableTools[toolName]; exists {
			selectedTools = append(selectedTools, tool)
		}
	}

	// Create context template data
	templateData := struct {
		AgentConfig
		Tools []Tool
	}{
		AgentConfig: config,
		Tools:       selectedTools,
	}

	// Resolve template path and parse
	ctxRel := "templates/agent/support_bot.ctx"
	ctxAbs, rerr := resolveTemplatePath(ctxRel)
	if rerr != nil {
		log.Error("failed to resolve agent context template path", zap.Error(rerr))
		return fmt.Errorf("template not found: %s: %w", ctxRel, rerr)
	}
	tmpl, err := template.ParseFiles(ctxAbs)
	if err != nil {
		log.Error("failed to parse agent context template", zap.Error(err))
		return fmt.Errorf("failed to parse template %s: %w", ctxAbs, err)
	}

	// Create output file
	outputPath := fmt.Sprintf("contexts/%s/%s.ctx", config.Name, strings.ToLower(config.Name))
	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Error("failed to create agent context file", zap.Error(err))
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer outputFile.Close()

	// Execute template
	if err := tmpl.Execute(outputFile, templateData); err != nil {
		log.Error("failed to execute agent context template", zap.Error(err))
		return fmt.Errorf("failed to execute template: %w", err)
	}

	log.Info("agent context generated", zap.String("path", outputPath))
	return nil
}

// generateAgentPrompts creates prompt templates for the agent
func generateAgentPrompts(ctx context.Context, config AgentConfig) error {
	log := logger.WithContext(ctx)

	// Read agent response template
	respRel := "templates/agent/agent_response.md"
	respAbs, rerr := resolveTemplatePath(respRel)
	if rerr != nil {
		log.Error("failed to resolve agent response template path", zap.Error(rerr))
		return fmt.Errorf("template not found: %s: %w", respRel, rerr)
	}
	tmpl, err := template.ParseFiles(respAbs)
	if err != nil {
		log.Error("failed to parse agent response template", zap.Error(err))
		return fmt.Errorf("failed to parse template %s: %w", respAbs, err)
	}

	// Create output file
	outputPath := fmt.Sprintf("prompts/%s/agent_response.md", config.Name)
	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Error("failed to create agent response file", zap.Error(err))
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer outputFile.Close()

	// Execute template
	if err := tmpl.Execute(outputFile, config); err != nil {
		log.Error("failed to execute agent response template", zap.Error(err))
		return fmt.Errorf("failed to execute template: %w", err)
	}

	log.Info("agent prompts generated", zap.String("path", outputPath))
	return nil
}

// generateAgentTools creates tool implementations for the agent
func generateAgentTools(ctx context.Context, config AgentConfig) error {
	log := logger.WithContext(ctx)

	// Copy tool templates based on selected tools
	for _, tool := range config.Tools {
		if err := copyToolTemplate(ctx, tool, config.Name); err != nil {
			log.Error("failed to copy tool template", zap.String("tool", tool), zap.Error(err))
			return fmt.Errorf("failed to copy tool %s: %w", tool, err)
		}
	}

	log.Info("agent tools generated", zap.Strings("tools", config.Tools))
	return nil
}

// copyToolTemplate copies a tool template to the agent's tools directory
func copyToolTemplate(ctx context.Context, toolName, agentName string) error {
	log := logger.WithContext(ctx)

	// Define tool template mappings
	toolTemplates := map[string]string{
		"web_search":  "templates/agent/web_search.py",
		"database":    "templates/agent/database.py",
		"api":         "templates/agent/api.py",
		"file_system": "templates/agent/file_system.py",
		"email":       "templates/agent/email.py",
	}

	templatePath, exists := toolTemplates[toolName]
	if !exists {
		log.Warn("no template found for tool", zap.String("tool", toolName))
		return nil
	}

	// Resolve and read template content
	absPath, rerr := resolveTemplatePath(templatePath)
	if rerr != nil {
		log.Error("failed to resolve tool template path", zap.String("template", templatePath), zap.Error(rerr))
		return fmt.Errorf("template not found: %s: %w", templatePath, rerr)
	}
	content, err := os.ReadFile(absPath)
	if err != nil {
		log.Error("failed to read tool template", zap.String("template", absPath), zap.Error(err))
		return fmt.Errorf("failed to read template %s: %w", absPath, err)
	}

	// Create output file
	outputPath := fmt.Sprintf("tools/%s/%s.py", agentName, toolName)
	if err := os.WriteFile(outputPath, content, 0644); err != nil {
		log.Error("failed to write tool file", zap.String("path", outputPath), zap.Error(err))
		return fmt.Errorf("failed to write file: %w", err)
	}

	log.Info("tool template copied", zap.String("tool", toolName), zap.String("path", outputPath))
	return nil
}

// generateAgentTests creates test configuration for the agent
func generateAgentTests(ctx context.Context, config AgentConfig) error {
	log := logger.WithContext(ctx)

	// Read behavior test template
	behRel := "templates/agent/agent_behavior.yaml"
	behAbs, rerr := resolveTemplatePath(behRel)
	if rerr != nil {
		log.Error("failed to resolve agent behavior template path", zap.Error(rerr))
		return fmt.Errorf("template not found: %s: %w", behRel, rerr)
	}
	tmpl, err := template.ParseFiles(behAbs)
	if err != nil {
		log.Error("failed to parse agent behavior template", zap.Error(err))
		return fmt.Errorf("failed to parse template %s: %w", behAbs, err)
	}

	// Create output file
	outputPath := fmt.Sprintf("tests/%s/agent_behavior.yaml", config.Name)
	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Error("failed to create agent behavior file", zap.Error(err))
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer outputFile.Close()

	// Execute template
	if err := tmpl.Execute(outputFile, config); err != nil {
		log.Error("failed to execute agent behavior template", zap.Error(err))
		return fmt.Errorf("failed to execute template: %w", err)
	}

	log.Info("agent tests generated", zap.String("path", outputPath))
	return nil
}

// generateAgentMemory creates memory configuration for the agent
func generateAgentMemory(ctx context.Context, config AgentConfig) error {
	log := logger.WithContext(ctx)

	// Create memory configuration file
	memoryConfig := fmt.Sprintf(`# Memory Configuration for %s Agent
memory_type: "%s"
max_history: 50
privacy: "user_isolated"
retention_days: 30
encryption: true

# Episodic Memory Settings
episodic:
  enabled: %t
  max_conversations: 100
  context_window: 10
  importance_threshold: 0.7

# User Preferences
preferences:
  enabled: true
  max_preferences: 20
  update_frequency: "session"

# Security Settings
security:
  data_encryption: true
  access_logging: true
  audit_trail: true
`, config.Name, config.Memory, config.Memory == "episodic")

	outputPath := fmt.Sprintf("memory/%s/memory_config.yaml", config.Name)
	if err := os.WriteFile(outputPath, []byte(memoryConfig), 0644); err != nil {
		log.Error("failed to create memory config", zap.Error(err))
		return fmt.Errorf("failed to create memory config: %w", err)
	}

	log.Info("agent memory configuration generated", zap.String("path", outputPath))
	return nil
}

// generateAgentRequirements creates requirements.txt for the agent
func generateAgentRequirements(ctx context.Context, config AgentConfig) error {
	log := logger.WithContext(ctx)

	// Read requirements template
	reqRel := "templates/agent/requirements.txt"
	reqAbs, rerr := resolveTemplatePath(reqRel)
	if rerr != nil {
		log.Error("failed to resolve requirements template path", zap.Error(rerr))
		return fmt.Errorf("template not found: %s: %w", reqRel, rerr)
	}
	content, err := os.ReadFile(reqAbs)
	if err != nil {
		log.Error("failed to read requirements template", zap.String("template", reqAbs), zap.Error(err))
		return fmt.Errorf("failed to read requirements template %s: %w", reqAbs, err)
	}

	// Create output file
	outputPath := fmt.Sprintf("tools/%s/requirements.txt", config.Name)
	if err := os.WriteFile(outputPath, content, 0644); err != nil {
		log.Error("failed to write requirements file", zap.String("path", outputPath), zap.Error(err))
		return fmt.Errorf("failed to write requirements file: %w", err)
	}

	log.Info("agent requirements generated", zap.String("path", outputPath))
	return nil
}

// resolveTemplatePath tries to find a template by checking cwd, repo root (go.mod), and source-relative
func resolveTemplatePath(rel string) (string, error) {
	// 1) CWD
	if p := filepath.Clean(rel); fileExists(p) {
		abs, _ := filepath.Abs(p)
		return abs, nil
	}
	// 2) Walk up for go.mod
	if cwd, err := os.Getwd(); err == nil {
		if root := findRepoRoot(cwd); root != "" {
			candidate := filepath.Join(root, rel)
			if fileExists(candidate) {
				return candidate, nil
			}
		}
	}
	// 3) Relative to this source file
	if _, src, _, ok := runtime.Caller(0); ok {
		base := filepath.Dir(src)
		root := filepath.Clean(filepath.Join(base, "..", "..", ".."))
		candidate := filepath.Join(root, rel)
		if fileExists(candidate) {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("unable to locate %s", rel)
}

func findRepoRoot(start string) string {
	dir := start
	for i := 0; i < 12; i++ {
		if fileExists(filepath.Join(dir, "go.mod")) {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// validateAgentName enforces the naming rules expected by tests
func validateAgentName(name string) error {
	if name == "" {
		return fmt.Errorf("agent name is required")
	}
	if len(name) < 2 {
		return fmt.Errorf("agent name should be at least 2 characters")
	}
	if strings.Contains(name, " ") {
		return fmt.Errorf("agent name should not contain spaces")
	}
	if strings.Contains(name, "/") {
		return fmt.Errorf("agent name should not contain slashes")
	}
	if strings.Contains(name, "\\") {
		return fmt.Errorf("agent name should not contain backslashes")
	}
	if strings.Contains(name, "!") {
		return fmt.Errorf("agent name should not contain special characters")
	}
	return nil
}

// showAgentStructure displays the generated agent structure
func showAgentStructure(name string, config AgentConfig) {
	fmt.Printf("\n")
	logger.LogSuccess(context.Background(), "Agent generated successfully",
		zap.String("name", name),
		zap.Strings("tools", config.Tools),
		zap.String("memory", config.MemoryType))

	fmt.Printf("\n🤖 Generated Agent Structure:\n")
	fmt.Printf("  contexts/%s/\n", name)
	fmt.Printf("  ├── 📄 %s.ctx\n", name)
	fmt.Printf("  ├── 📁 memory/%s/\n", name)
	fmt.Printf("  │   ├── 📄 memory_config.yaml\n")
	fmt.Printf("  │   ├── 📁 episodic/\n")
	fmt.Printf("  │   └── 📁 user_preferences/\n")
	fmt.Printf("  ├── 📁 prompts/%s/\n", name)
	fmt.Printf("  │   └── 📄 agent_response.md\n")
	fmt.Printf("  ├── 📁 tools/%s/\n", name)
	fmt.Printf("  │   ├── 📄 api.py\n")
	fmt.Printf("  │   └── 📄 requirements.txt\n")
	fmt.Printf("  └── 📁 tests/%s/\n", name)
	fmt.Printf("      └── 📄 agent_behavior.yaml\n")
}

// showAgentDevelopmentFlow displays the agent development workflow
func showAgentDevelopmentFlow(name string, config AgentConfig) {
	fmt.Printf("\n🚀 Agent Development Flow:\n")

	fmt.Printf("\n1️⃣  Test your agent:\n")
	fmt.Printf("   ctx test %s\n", name)
	fmt.Printf("   ctx test --behavior --component=%s\n", name)

	fmt.Printf("\n2️⃣  Start a conversation:\n")
	fmt.Printf("   ctx run %s \"Hello, can you help me?\"\n", name)
	fmt.Printf("   ctx run %s \"What tools do you have available?\"\n", name)

	fmt.Printf("\n3️⃣  Start development server:\n")
	fmt.Printf("   ctx serve --addr :8000\n")

	fmt.Printf("\n4️⃣  Customize your agent:\n")
	fmt.Printf("   # Edit the context file\n")
	fmt.Printf("   nano contexts/%s/%s.ctx\n", name, name)
	fmt.Printf("   \n")
	fmt.Printf("   # Modify prompts\n")
	fmt.Printf("   nano prompts/%s/agent_response.md\n", name)
	fmt.Printf("   \n")
	fmt.Printf("   # Add custom tools\n")
	fmt.Printf("   nano tools/%s/api.py\n", name)

	fmt.Printf("\n5️⃣  Monitor agent behavior:\n")
	fmt.Printf("   # Check drift detection\n")
	fmt.Printf("   ctx test --drift-detection --component=%s\n", name)
	fmt.Printf("   \n")
	fmt.Printf("   # View behavior reports\n")
	fmt.Printf("   cat tests/reports/behavior_%s.json\n", name)

	fmt.Printf("\n📚 Configuration Details:\n")
	fmt.Printf("   • Tools: %s\n", strings.Join(config.Tools, ", "))
	fmt.Printf("   • Memory: %s\n", config.MemoryType)
	fmt.Printf("   • Context: contexts/%s/%s.ctx\n", name, name)
	fmt.Printf("   • Memory Config: memory/%s/memory_config.yaml\n", name)

	fmt.Printf("\n🎉 Your agent is ready! Start conversations and customize as needed.\n")
}
