package commands

import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
)

// GetRootCommand returns the root command
func GetRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "ctx",
		Short: "Contexis CMP Framework CLI",
		Long:  `A comprehensive CLI for the Context-Memory-Prompt (CMP) framework.`,
	}

    // Add subcommands
    rootCmd.AddCommand(GetGenerateCommand())
    rootCmd.AddCommand(GetPluginCommand(getProjectRoot()))
	rootCmd.AddCommand(GetVersionCommand())

	return rootCmd
}

// getProjectRoot returns current working directory as project root
func getProjectRoot() string {
    cwd, _ := os.Getwd()
    return cwd
}

// GetGenerateCommand returns the generate command
func GetGenerateCommand() *cobra.Command {
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate CMP components",
		Long:  `Generate RAG systems, agents, and workflows using CMP templates.`,
	}

	// Add subcommands
	generateCmd.AddCommand(GetAgentCommand())
	generateCmd.AddCommand(GetRAGCommand())
	generateCmd.AddCommand(GetWorkflowCommand())

	return generateCmd
}

// GetAgentCommand returns the agent command
func GetAgentCommand() *cobra.Command {
    agentCmd := &cobra.Command{
        Use:   "agent",
		Short: "Generate a conversational agent",
		Long:  `Generate a conversational agent with specified tools and memory type.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			tools, _ := cmd.Flags().GetString("tools")
			memory, _ := cmd.Flags().GetString("memory")

			// Validate inputs
			if name == "" {
				return fmt.Errorf("agent name is required")
			}

			// Generate agent
			return GenerateAgent(cmd.Context(), name, tools, memory)
		},
	}

	// Add flags
	agentCmd.Flags().StringP("tools", "t", "", "Comma-separated list of tools (web_search,database,api,file_system,email)")
	agentCmd.Flags().StringP("memory", "m", "episodic", "Memory type (episodic,none)")

	return agentCmd
}

// GetRAGCommand returns the RAG command
func GetRAGCommand() *cobra.Command {
	ragCmd := &cobra.Command{
        Use:   "rag",
		Short: "Generate a RAG system",
		Long:  `Generate a Retrieval-Augmented Generation (RAG) system.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = args[0] // name parameter
			// TODO: Implement RAG generation
			return fmt.Errorf("RAG generation not yet implemented")
		},
	}

	return ragCmd
}

// GetWorkflowCommand returns the workflow command
func GetWorkflowCommand() *cobra.Command {
	workflowCmd := &cobra.Command{
        Use:   "workflow",
		Short: "Generate a workflow",
		Long:  `Generate a multi-step AI workflow.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			steps, _ := cmd.Flags().GetString("steps")

			// Validate inputs
			if name == "" {
				return fmt.Errorf("workflow name is required")
			}

			// Generate workflow
			return GenerateWorkflow(cmd.Context(), name, steps)
		},
	}

	// Add flags
	workflowCmd.Flags().StringP("steps", "s", "", "Comma-separated list of steps (research,write,review,extract,transform,load,analyze,generate,validate,deploy)")

	return workflowCmd
}

// GetVersionCommand returns the version command
func GetVersionCommand() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Contexis CMP Framework v0.1.0")
		},
	}

	return versionCmd
}
