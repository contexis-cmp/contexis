package commands

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "github.com/contexis-cmp/contexis/src/plugins/registry"
    "github.com/spf13/cobra"
)

// GetPluginCommand returns the parent plugin command with subcommands.
func GetPluginCommand(projectRoot string) *cobra.Command {
    pluginCmd := &cobra.Command{Use: "plugin", Short: "Manage CMP plugins"}

    // list
    pluginCmd.AddCommand(&cobra.Command{
        Use:   "list",
        Short: "List installed plugins",
        RunE: func(cmd *cobra.Command, args []string) error {
            reg := registry.NewRegistry(projectRoot)
            list, err := reg.List()
            if err != nil {
                return err
            }
            if len(list) == 0 {
                fmt.Fprintln(cmd.OutOrStdout(), "No plugins installed")
                return nil
            }
            for _, p := range list {
                fmt.Fprintf(cmd.OutOrStdout(), "%s %s - %s\n", p.Manifest.Name, p.Manifest.Version, p.Manifest.Description)
            }
            return nil
        },
    })

    // info
    pluginCmd.AddCommand(&cobra.Command{
        Use:   "info [name]",
        Short: "Show plugin details",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            reg := registry.NewRegistry(projectRoot)
            p, err := reg.Info(args[0])
            if err != nil {
                return err
            }
            fmt.Fprintf(cmd.OutOrStdout(), "Name: %s\nVersion: %s\nDescription: %s\nCapabilities: %v\n",
                p.Manifest.Name, p.Manifest.Version, p.Manifest.Description, p.Manifest.Capabilities)
            return nil
        },
    })

    // install (local dir path)
    pluginCmd.AddCommand(&cobra.Command{
        Use:   "install [path_or_url]",
        Short: "Install a plugin from local path, zip URL, or Git URL",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            src := args[0]
            if !filepath.IsAbs(src) {
                cwd, _ := os.Getwd()
                src = filepath.Join(cwd, src)
            }
            reg := registry.NewRegistry(projectRoot)
            var (
                p   registry.Plugin
                err error
            )
            if strings.HasPrefix(args[0], "http") || strings.HasPrefix(args[0], "git@") || strings.Contains(args[0], ".zip") {
                p, err = reg.InstallRemote(args[0])
            } else {
                p, err = reg.Install(src)
            }
            if err != nil {
                return err
            }
            fmt.Fprintf(cmd.OutOrStdout(), "Installed plugin: %s %s\n", p.Manifest.Name, p.Manifest.Version)
            return nil
        },
    })

    // remove
    pluginCmd.AddCommand(&cobra.Command{
        Use:   "remove [name]",
        Short: "Remove an installed plugin",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            reg := registry.NewRegistry(projectRoot)
            if err := reg.Remove(args[0]); err != nil {
                return err
            }
            fmt.Fprintf(cmd.OutOrStdout(), "Removed plugin: %s\n", args[0])
            return nil
        },
    })

    return pluginCmd
}


