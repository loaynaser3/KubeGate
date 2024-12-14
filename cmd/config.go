package cmd

import (
	"fmt"

	"github.com/loaynaser3/KubeGate/pkg/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage KubeGate contexts",
}

var setContextCmd = &cobra.Command{
	Use:   "set-context",
	Short: "Add or update a context",
	Args:  cobra.ExactArgs(3), // name, rabbitmq-url, reply-queue
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Println("Failed to load config:", err)
			return
		}

		ctx := config.Context{
			Name:         args[0],
			RabbitMQURL:  args[1],
			CommandQueue: args[2],
			ReplyQueue:   args[3],
		}

		config.SetContext(cfg, ctx)
		err = config.SaveConfig(cfg)
		if err != nil {
			fmt.Println("Failed to save config:", err)
			return
		}

		fmt.Println("Context set successfully.")
	},
}

var getContextsCmd = &cobra.Command{
	Use:   "get-contexts",
	Short: "List all contexts",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Println("Failed to load config:", err)
			return
		}

		fmt.Println("Available contexts:")
		for _, ctx := range cfg.Contexts {
			fmt.Printf("- %s (RabbitMQ: %s, ReplyQueue: %s)\n", ctx.Name, ctx.RabbitMQURL, ctx.ReplyQueue)
		}
		fmt.Printf("Current context: %s\n", cfg.CurrentContext)
	},
}

func init() {
	configCmd.AddCommand(setContextCmd)
	configCmd.AddCommand(getContextsCmd)
	rootCmd.AddCommand(configCmd)
}
