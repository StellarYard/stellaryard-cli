package cmd

import (
	"fmt"
	"net/http"
	"io"
	"encoding/json"

	"github.com/spf13/cobra"
)

var containersCmd = &cobra.Command{
	Use:   "containers",
	Short: "Manage StellarYard containers",
}

var containersStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a container",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			fatalf(ExitArgError, "container name is required (--name)")
		}

		resp, err := http.Post(
			fmt.Sprintf("%s/api/v1/containers/%s/start", coreURL, name),
			"application/json",
			nil,
		)
		if err != nil {
			fatalf(ExitCoreUnreachable, "cannot reach stellaryard-core at %s. Is it running?", coreURL)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			fatalf(ExitAppError, "container start failed: %s", string(body))
		}

		fmt.Printf("Container %s started\n", name)
	},
}

var containersStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a container",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			fatalf(ExitArgError, "container name is required (--name)")
		}

		resp, err := http.Post(
			fmt.Sprintf("%s/api/v1/containers/%s/stop", coreURL, name),
			"application/json",
			nil,
		)
		if err != nil {
			fatalf(ExitCoreUnreachable, "cannot reach stellaryard-core at %s. Is it running?", coreURL)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			fatalf(ExitAppError, "container stop failed: %s", string(body))
		}

		fmt.Printf("Container %s stopped\n", name)
	},
}

var containersStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show container status",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := http.Get(fmt.Sprintf("%s/api/v1/containers", coreURL))
		if err != nil {
			fatalf(ExitCoreUnreachable, "cannot reach stellaryard-core at %s. Is it running?", coreURL)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		if formatStr == "json" {
			fmt.Println(string(body))
			return
		}

		// Parse and format as table
		var containers []struct {
			Name  string `json:"name"`
			State string `json:"state"`
		}
		if err := json.Unmarshal(body, &containers); err != nil {
			fatalf(ExitAppError, "failed to parse response: %v", err)
		}

		if len(containers) == 0 {
			fmt.Println("No containers found")
			return
		}

		fmt.Printf("%-20s %-15s\n", "NAME", "STATE")
		fmt.Printf("%-20s %-15s\n", "----", "-----")
		for _, c := range containers {
			fmt.Printf("%-20s %-15s\n", c.Name, c.State)
		}
	},
}

func init() {
	containersStartCmd.Flags().String("name", "", "Container name (horizon, soroban-rpc)")
	containersStopCmd.Flags().String("name", "", "Container name (horizon, soroban-rpc)")

	containersCmd.AddCommand(containersStartCmd)
	containersCmd.AddCommand(containersStopCmd)
	containersCmd.AddCommand(containersStatusCmd)
	rootCmd.AddCommand(containersCmd)
}
