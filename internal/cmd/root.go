package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Exit codes — load-bearing for CI automation.
// Do not add new codes or repurpose existing ones.
const (
	ExitSuccess          = 0
	ExitCommandError     = 1
	ExitCoreUnreachable  = 2
	ExitCoreApplication  = 3
)

var (
	coreURL   string
	formatStr string
)

var rootCmd = &cobra.Command{
	Use:   "stellaryard",
	Short: "StellarYard — local Stellar/Soroban development environment CLI",
	Long:  "Scriptable terminal client for StellarYard. Manages containers, accounts, contracts, and ledger data via stellaryard-core's API.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// TODO: validate coreURL is reachable
		// Classify errors as exit 2 (unreachable) or 3 (application error)
		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&coreURL, "core-url", "http://localhost:8080", "stellaryard-core API URL")
	rootCmd.PersistentFlags().StringVar(&formatStr, "format", "table", "output format: table or json")
}

// Execute runs the root command.
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return err
	}
	return nil
}
