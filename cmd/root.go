package cmd

import (
	"fmt"
	"os"

	"github.com/loaynaser3/KubeGate/pkg/logging"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kubegate",
	Short: "KubeGate: A Kubernetes proxy tool",
	Long: `KubeGate is a proxy tool that facilitates secure access to private Kubernetes clusters
by leveraging RabbitMQ and Kubernetes-native interactions. It supports single-command execution
and an interactive shell.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to KubeGate! Use 'kubegate --help' to explore commands.")
	},
}

// Execute initializes the root command
func Execute() {
	logging.Logger.Info("KubeGate CLI started")
	if err := rootCmd.Execute(); err != nil {
		logging.Logger.WithField("error", err.Error()).Error("Command execution failed")

		os.Exit(1)
	}
}
