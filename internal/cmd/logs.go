package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs <container-name>",
	Short: "Stream container logs",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		follow, _ := cmd.Flags().GetBool("follow")
		// TODO: connect to core WS endpoint /containers/{name}/logs
		// TODO: handle Ctrl+C cleanly (close WS connection, not just os.Exit)
		fmt.Printf("streaming logs for %s (follow=%v): not implemented\n", name, follow)
		return nil
	},
}

func init() {
	logsCmd.Flags().Bool("follow", false, "follow log output (like tail -f)")

	rootCmd.AddCommand(logsCmd)
}
