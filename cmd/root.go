package cmd

import (
	"fmt"
	"os"

	"github.com/loaynaser3/KubeGate/pkg/logging"
	"github.com/sirupsen/logrus"
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
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		logging.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to start")
		os.Exit(1)
	}
}
