package cmd

import (
	"fmt"
	"github.com/loaynaser3/KubeGate/pkg/KubeGate"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute a Kubernetes command",
	Long: `The run command sends a Kubernetes command to the KubeGate agent
via RabbitMQ for execution and retrieves the response.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Error: You must provide a Kubernetes command to execute.")
			return
		}

		// Combine args into a single command string
		kubeCommand := args[0]
		if len(args) > 1 {
			kubeCommand = fmt.Sprintf("%s %s", args[0], args[1:])
		}

		// Delegate to the `run` function in the `pkg/KubeGate` package
		KubeGate.ExecuteRun(kubeCommand)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
