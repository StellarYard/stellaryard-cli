package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var contractsCmd = &cobra.Command{
	Use:   "contracts",
	Short: "Deploy and invoke Soroban contracts",
}

var contractsDeployCmd = &cobra.Command{
	Use:   "deploy <wasm-path>",
	Short: "Deploy a WASM contract",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		wasmPath := args[0]
		// TODO: validate WASM file exists locally before sending to core
		// TODO: call core API POST /contracts/deploy
		fmt.Printf("deploying contract from %s: not implemented\n", wasmPath)
		return nil
	},
}

var contractsInvokeCmd = &cobra.Command{
	Use:   "invoke <contractId> <method> [args...]",
	Short: "Invoke a contract method",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		contractID := args[0]
		method := args[1]
		contractArgs := args[2:]
		// TODO: call core API POST /contracts/{contractId}/invoke
		fmt.Printf("invoking %s.%s with %v: not implemented\n", contractID, method, contractArgs)
		return nil
	},
}

func init() {
	contractsCmd.AddCommand(contractsDeployCmd)
	contractsCmd.AddCommand(contractsInvokeCmd)

	rootCmd.AddCommand(contractsCmd)
}
