package commands

import (
	"fmt"
	"os"
	"path/filepath"

	runtimeserver "github.com/contexis-cmp/contexis/src/runtime/server"
	"github.com/spf13/cobra"
)

// GetServeCommand returns the `serve` command.
//
// The `serve` command runs the Contexis HTTP server with sensible local-first
// defaults. It auto-detects a project virtualenv (`.venv/bin/python`) to use
// for the local Python provider, sets `CMP_LOCAL_MODELS=true` when unset, and
// exports `CMP_PROJECT_ROOT` to the current working directory for resolving
// contexts, prompts and memory paths. A warning is printed if `contexts/`
// is not found at the project root.
func GetServeCommand() *cobra.Command {
	var addr string
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run a simple HTTP server for chat",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Auto-detect project virtualenv python if not explicitly set
			if os.Getenv("CMP_PYTHON_BIN") == "" {
				if wd, err := os.Getwd(); err == nil {
					cand := filepath.Join(wd, ".venv", "bin", "python")
					if _, err := os.Stat(cand); err == nil {
						_ = os.Setenv("CMP_PYTHON_BIN", cand)
					}
				}
			}
			// Default to local-first provider if unset
			if os.Getenv("CMP_LOCAL_MODELS") == "" {
				_ = os.Setenv("CMP_LOCAL_MODELS", "true")
			}
			// Ensure project root is set for provider/template resolution
			if os.Getenv("CMP_PROJECT_ROOT") == "" {
				if wd, err := os.Getwd(); err == nil {
					_ = os.Setenv("CMP_PROJECT_ROOT", wd)
				}
			}
			// Warn if contexts directory is missing
			root := os.Getenv("CMP_PROJECT_ROOT")
			if _, err := os.Stat(filepath.Join(root, "contexts")); err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), "warning: 'contexts/' not found in project root; ensure you are in the project directory or set CMP_PROJECT_ROOT")
			}
			return runtimeserver.Serve(addr)
		},
	}
	cmd.Flags().StringVar(&addr, "addr", ":8000", "Listen address")
	return cmd
}
