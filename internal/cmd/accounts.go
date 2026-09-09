package cmd

import (
	"fmt"
	"net/http"
	"io"
	"strings"

	"github.com/spf13/cobra"
)

var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Manage test accounts",
}

var accountsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create and fund a new test account",
	Run: func(cmd *cobra.Command, args []string) {
		label, _ := cmd.Flags().GetString("label")

		body := fmt.Sprintf(`{"label":"%s"}`, label)
		resp, err := http.Post(
			fmt.Sprintf("%s/api/v1/accounts", coreURL),
			"application/json",
			strings.NewReader(body),
		)
		if err != nil {
			fatalf(ExitCoreUnreachable, "cannot reach stellaryard-core at %s. Is it running?", coreURL)
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusCreated {
			fatalf(ExitAppError, "account creation failed: %s", string(respBody))
		}

		if formatStr == "json" {
			fmt.Println(string(respBody))
		} else {
			fmt.Println("Account created successfully")
		}
	},
}

var accountsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all managed accounts",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := http.Get(fmt.Sprintf("%s/api/v1/accounts", coreURL))
		if err != nil {
			fatalf(ExitCoreUnreachable, "cannot reach stellaryard-core at %s. Is it running?", coreURL)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		if formatStr == "json" {
			fmt.Println(string(body))
			return
		}

		fmt.Println("No accounts yet. Use 'stellaryard accounts create' to add one.")
	},
}

func init() {
	accountsCreateCmd.Flags().String("label", "", "Account label")
	accountsCmd.AddCommand(accountsCreateCmd)
	accountsCmd.AddCommand(accountsListCmd)
	rootCmd.AddCommand(accountsCmd)
}
