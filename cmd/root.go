package cmd

import (
	"fmt"
	"os"

	"github.com/loaynaser3/KubeGate/pkg/KubeGate"
	"github.com/loaynaser3/KubeGate/pkg/logging"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kubegate",
	Short: "KubeGate: A Kubernetes proxy tool",
	Long: `KubeGate is a proxy tool that facilitates secure access to private Kubernetes clusters
by leveraging RabbitMQ and Kubernetes-native interactions. It supports single-command execution
and an interactive shell.`,
	Args: cobra.ArbitraryArgs, // Allow arbitrary arguments
	Run: func(cmd *cobra.Command, args []string) {
		// Handle the `-h` or `--help` flag
		if containsHelpFlag(args) {
			cmd.Help() // Show the default help message
			return
		}

		if len(args) == 0 {
			// No subcommand or arguments provided
			fmt.Println("Welcome to KubeGate! Use 'kubegate --help' to explore commands.")
			return
		}

		// Treat the remaining args as a Kubernetes command
		kubeCommand := args[0]
		commandArgs := args[1:]

		logging.Logger.WithFields(map[string]interface{}{
			"command": kubeCommand,
			"args":    commandArgs,
		}).Info("Running command as kubegate run")
		KubeGate.ExecuteRun(kubeCommand, commandArgs)
	},
}

// Execute initializes the root command
func Execute() {
	logging.Logger.Info("KubeGate CLI started")
	// Disable Cobra's default flag parsing
	rootCmd.DisableFlagParsing = true

	if err := rootCmd.Execute(); err != nil {
		logging.Logger.WithField("error", err.Error()).Error("Command execution failed")
		os.Exit(1)
	}
}

// containsHelpFlag checks if the arguments include `-h` or `--help`
func containsHelpFlag(args []string) bool {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			return true
		}
	}
	return false
}
