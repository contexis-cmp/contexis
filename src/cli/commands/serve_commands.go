package commands

import (
	runtimeserver "github.com/contexis-cmp/contexis/src/runtime/server"
	"github.com/spf13/cobra"
)

func GetServeCommand() *cobra.Command {
	var addr string
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run a simple HTTP server for chat",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runtimeserver.Serve(addr)
		},
	}
	cmd.Flags().StringVar(&addr, "addr", ":8000", "Listen address")
	return cmd
}
