package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Manage test accounts",
}

var accountsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create and fund a new test account",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		label, _ := cmd.Flags().GetString("label")
		// TODO: call core API POST /accounts
		fmt.Printf("creating account with label: %s\n", label)
		return nil
	},
}

var accountsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List managed accounts",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: call core API GET /accounts
		// Format output based on --format flag (table or json)
		fmt.Println("account list: not implemented")
		return nil
	},
}

var accountsShowCmd = &cobra.Command{
	Use:   "show <publicKey>",
	Short: "Show account details and balance",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		publicKey := args[0]
		// TODO: call core API GET /accounts/{publicKey}
		fmt.Printf("account details for %s: not implemented\n", publicKey)
		return nil
	},
}

func init() {
	accountsCreateCmd.Flags().String("label", "", "human-readable account label")

	accountsCmd.AddCommand(accountsCreateCmd)
	accountsCmd.AddCommand(accountsListCmd)
	accountsCmd.AddCommand(accountsShowCmd)

	rootCmd.AddCommand(accountsCmd)
}
