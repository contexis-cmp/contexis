package commands

import (
    runtimeworker "github.com/contexis-cmp/contexis/src/runtime/worker"
    "github.com/spf13/cobra"
)

func GetWorkerCommand() *cobra.Command {
    var addr string
    cmd := &cobra.Command{
        Use:   "worker",
        Short: "Run background worker endpoints",
        RunE: func(cmd *cobra.Command, args []string) error {
            return runtimeworker.Serve(addr)
        },
    }
    cmd.Flags().StringVar(&addr, "addr", ":9000", "Listen address for worker HTTP endpoints")
    return cmd
}


