package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var ledgerCmd = &cobra.Command{
	Use:   "ledger",
	Short: "Inspect ledger state and transactions",
}

var ledgerSnapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Show current ledger state summary",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: call core API GET /ledger/snapshot
		fmt.Println("ledger snapshot: not implemented")
		return nil
	},
}

var ledgerTxCmd = &cobra.Command{
	Use:   "tx",
	Short: "Transaction commands",
}

var ledgerTxListCmd = &cobra.Command{
	Use:   "list",
	Short: "List recent transactions",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		// TODO: call core API GET /ledger/transactions?limit=N
		// NOTE: blocked on core Phase 3 XDR-decoding decision
		fmt.Printf("transaction list (limit %d): not implemented\n", limit)
		return nil
	},
}

func init() {
	ledgerTxListCmd.Flags().Int("limit", 20, "maximum number of transactions to return")

	ledgerTxCmd.AddCommand(ledgerTxListCmd)
	ledgerCmd.AddCommand(ledgerSnapshotCmd)
	ledgerCmd.AddCommand(ledgerTxCmd)

	rootCmd.AddCommand(ledgerCmd)
}
