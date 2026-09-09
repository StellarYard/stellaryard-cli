package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Exit codes
const (
	ExitSuccess      = 0
	ExitArgError     = 1
	ExitCoreUnreachable = 2
	ExitAppError     = 3
)

var (
	coreURL   string
	verbose   bool
	formatStr string
)

var rootCmd = &cobra.Command{
	Use:   "stellaryard",
	Short: "A scriptable terminal client for StellarYard",
	Long:  "StellarYard CLI manages local Stellar development environments through stellaryard-core's API.",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&coreURL, "core-url", "http://localhost:8080", "Core API URL")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable verbose output")
	rootCmd.PersistentFlags().StringVar(&formatStr, "format", "table", "Output format: table, json")
}

func fatalf(code int, format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
	os.Exit(code)
}
