package cmd

import (
	"fmt"
	"io"
	"net/http"
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
		req, err := newRequest("POST", fmt.Sprintf("%s/api/v1/accounts", coreURL), strings.NewReader(body))
		if err != nil {
			fatalf(ExitAppError, "failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
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
		req, err := newRequest("GET", fmt.Sprintf("%s/api/v1/accounts", coreURL), nil)
		if err != nil {
			fatalf(ExitAppError, "failed to create request: %v", err)
		}

		resp, err := http.DefaultClient.Do(req)
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
