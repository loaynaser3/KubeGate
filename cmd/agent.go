package cmd

import (
	"fmt"
	"log"

	"github.com/loaynaser3/KubeGate/pkg/KubeGate"
	"github.com/spf13/cobra"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Start the KubeGate agent",
	Long: `The agent listens for messages from RabbitMQ, executes the requested commands
on the Kubernetes cluster, and sends the results back to the client.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting KubeGate agent...")
		err := KubeGate.StartAgent()
		if err != nil {
			log.Fatalf("Agent encountered an error: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(agentCmd)
}
