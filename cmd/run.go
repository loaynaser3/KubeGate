package cmd

import (
	"github.com/loaynaser3/KubeGate/pkg/KubeGate"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [command] [args...]",
	Short: "Execute a Kubernetes command",
	Long: `The run command sends a Kubernetes command to the KubeGate agent
via RabbitMQ for execution and retrieves the response.`,
	Args: cobra.MinimumNArgs(1), // Require at least one argument
	Run: func(cmd *cobra.Command, args []string) {
		// The first argument is the Kubernetes command (e.g., "get pods")
		kubeCommand := args[0]

		// Any additional arguments are passed to the command (e.g., "-n default")
		commandArgs := args[1:]

		// Delegate to the `run` function in the `pkg/KubeGate` package
		KubeGate.ExecuteRun(kubeCommand, commandArgs)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
