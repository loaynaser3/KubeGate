package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kubegate",
	Short: "KubeGate is a proxy tool for kubectl to access private Kubernetes clusters.",
	Long: `KubeGate is a CLI tool that provides secure proxy access to private Kubernetes clusters
using a message queue-based architecture. It supports single-command execution and an interactive shell.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to KubeGate! Use --help to see available commands.")
	},
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
