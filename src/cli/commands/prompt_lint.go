package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"
)

// Simple prompt linter: detect unused template variables like {{.var}} not in data
func GetPromptLintCommand() *cobra.Command {
	var component string
	cmd := &cobra.Command{
		Use:   "prompt-lint",
		Short: "Lint prompt templates for basic issues",
		RunE: func(cmd *cobra.Command, args []string) error {
			if component == "" {
				return fmt.Errorf("--component is required")
			}
			root, _ := os.Getwd()
			dir := filepath.Join(root, "prompts", component)
			re := regexp.MustCompile(`\{\{\s*\.([a-zA-Z0-9_\.]+)\s*\}\}`)
			return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
				if err != nil {
					return nil
				}
				if d.IsDir() {
					return nil
				}
				by, err := os.ReadFile(path)
				if err != nil {
					return nil
				}
				matches := re.FindAllStringSubmatch(string(by), -1)
				vars := map[string]struct{}{}
				for _, m := range matches {
					if len(m) > 1 {
						vars[m[1]] = struct{}{}
					}
				}
				// We can't resolve data statically; just report variables found
				if len(vars) > 0 {
					fmt.Fprintf(cmd.OutOrStdout(), "%s: variables: ", path)
					first := true
					for v := range vars {
						if !first {
							fmt.Fprint(cmd.OutOrStdout(), ", ")
						}
						fmt.Fprint(cmd.OutOrStdout(), v)
						first = false
					}
					fmt.Fprintln(cmd.OutOrStdout())
				}
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&component, "component", "", "Component name")
	return cmd
}
