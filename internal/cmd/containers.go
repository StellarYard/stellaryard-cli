package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var containersCmd = &cobra.Command{
	Use:   "containers",
	Short: "Manage Horizon and Soroban RPC containers",
}

var containersStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a container",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		// TODO: call core API POST /containers/{name}/start
		fmt.Printf("starting container: %s\n", name)
		return nil
	},
}

var containersStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a container",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		// TODO: call core API POST /containers/{name}/stop
		fmt.Printf("stopping container: %s\n", name)
		return nil
	},
}

var containersStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "List container statuses",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: call core API GET /containers
		// Format output based on --format flag
		fmt.Println("container status: not implemented")
		return nil
	},
}

func init() {
	containersStartCmd.Flags().String("name", "horizon", "container name (horizon or soroban-rpc)")
	containersStopCmd.Flags().String("name", "horizon", "container name (horizon or soroban-rpc)")

	containersCmd.AddCommand(containersStartCmd)
	containersCmd.AddCommand(containersStopCmd)
	containersCmd.AddCommand(containersStatusCmd)

	rootCmd.AddCommand(containersCmd)
}
