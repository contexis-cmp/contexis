package commands

import (
	"encoding/json"
	"fmt"
	"os"

	runtimeprompt "github.com/contexis-cmp/contexis/src/runtime/prompt"
	"github.com/spf13/cobra"
)

func GetPromptCommand() *cobra.Command {
	pc := &cobra.Command{Use: "prompt", Short: "Prompt operations (render, validate)"}
	pc.AddCommand(newPromptRenderCmd())
	pc.AddCommand(newPromptValidateCmd())
	return pc
}

func newPromptRenderCmd() *cobra.Command {
	var component string
	var templatePath string
	var dataJSON string
	cmd := &cobra.Command{
		Use:   "render",
		Short: "Render a prompt template with JSON data",
		RunE: func(cmd *cobra.Command, args []string) error {
			if component == "" || templatePath == "" {
				return fmt.Errorf("--component and --template are required")
			}
			eng := runtimeprompt.NewEngine(mustGetwd())
			data := map[string]interface{}{}
			if dataJSON != "" {
				if err := json.Unmarshal([]byte(dataJSON), &data); err != nil {
					return fmt.Errorf("invalid --data json: %w", err)
				}
			}
			out, err := eng.RenderFile(component, templatePath, data)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), out)
			return nil
		},
	}
	cmd.Flags().StringVar(&component, "component", "", "Component name (e.g., CustomerDocs, SupportBot)")
	cmd.Flags().StringVar(&templatePath, "template", "", "Relative template path under prompts/<component>/, e.g. search_response.md")
	cmd.Flags().StringVar(&dataJSON, "data", "", "JSON object with template data")
	return cmd
}

func newPromptValidateCmd() *cobra.Command {
	var format string
	var inputPath string
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate a rendered response for a given format",
		RunE: func(cmd *cobra.Command, args []string) error {
			var content []byte
			var err error
			if inputPath == "" {
				content, err = os.ReadFile("-")
				if err != nil {
					return fmt.Errorf("reading stdin not supported; provide --input path")
				}
			} else {
				content, err = os.ReadFile(inputPath)
				if err != nil {
					return err
				}
			}
			if err := runtimeprompt.ValidateFormat(format, string(content)); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "valid")
			return nil
		},
	}
	cmd.Flags().StringVar(&format, "format", "markdown", "Expected format: json|markdown|text")
	cmd.Flags().StringVar(&inputPath, "input", "", "Path to file with rendered content")
	return cmd
}
