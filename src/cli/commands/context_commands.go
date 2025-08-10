package commands

import (
	"fmt"
	"os"
	"path/filepath"

	runtimecontext "github.com/contexis-cmp/contexis/src/runtime/context"
	"github.com/spf13/cobra"
)

// GetContextCommand returns the parent context command with subcommands.
func GetContextCommand(projectRoot string) *cobra.Command {
	ctxCmd := &cobra.Command{
		Use:   "context",
		Short: "Context operations (validate, reload)",
	}

	ctxCmd.AddCommand(newContextValidateCmd(projectRoot))
	ctxCmd.AddCommand(newContextReloadCmd(projectRoot))
	return ctxCmd
}

func newContextValidateCmd(projectRoot string) *cobra.Command {
	var tenantID string
	cmd := &cobra.Command{
		Use:   "validate [contextName]",
		Short: "Validate a context (.ctx) file for a tenant",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			contextName := args[0]
			if projectRoot == "" {
				cwd, _ := os.Getwd()
				projectRoot = cwd
			}
			svc := runtimecontext.NewContextService(projectRoot)
			ctx, err := svc.ResolveContext(tenantID, contextName)
			if err != nil {
				return err
			}
			sha, _ := ctx.GetSHA()
			fmt.Fprintf(cmd.OutOrStdout(), "Context '%s' is valid (SHA %s)\n", contextName, sha)
			return nil
		},
	}
	cmd.Flags().StringVar(&tenantID, "tenant", "", "Tenant ID for tenant-specific context resolution")
	return cmd
}

func newContextReloadCmd(projectRoot string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reload",
		Short: "Reload contexts by clearing runtime cache",
		RunE: func(cmd *cobra.Command, args []string) error {
			if projectRoot == "" {
				cwd, _ := os.Getwd()
				projectRoot = cwd
			}
			svc := runtimecontext.NewContextService(projectRoot)
			if err := svc.ReloadContext(filepath.Join(projectRoot, "contexts")); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Context cache cleared")
			return nil
		},
	}
	return cmd
}
