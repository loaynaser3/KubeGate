package cmd

import (
	"fmt"

	"github.com/loaynaser3/KubeGate/pkg/config"
	"github.com/spf13/cobra"
)

var (
	contextName  string
	rabbitMQURL  string
	commandQueue string
	replyQueue   string
	backend      string
)

// Root command for config
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage KubeGate contexts",
}

// Set context command with flags
var setContextCmd = &cobra.Command{
	Use:   "set-context",
	Short: "Add or update a context",
	Run: func(cmd *cobra.Command, args []string) {
		if contextName == "" || rabbitMQURL == "" || commandQueue == "" || replyQueue == "" || backend == "" {
			fmt.Println("All flags must be provided: --name, --rabbitmq-url, --command-queue, --reply-queue, --backend")
			return
		}

		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Println("Failed to load config:", err)
			return
		}

		ctx := config.Context{
			Name:         contextName,
			RabbitMQURL:  rabbitMQURL,
			CommandQueue: commandQueue,
			ReplyQueue:   replyQueue,
			Backend:      backend,
		}

		// Set the context
		config.SetContext(cfg, ctx)
		cfg.CurrentContext = contextName // Automatically switch to the new context
		err = config.SaveConfig(cfg)
		if err != nil {
			fmt.Println("Failed to save config:", err)
			return
		}

		fmt.Println("Context set and switched to:", contextName)
	},
}

// Get available contexts
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
			fmt.Printf("- %s (BrokerURL: %s, ReplyQueue: %s, CommandQueue: %s, Backend: %s)\n", ctx.Name, ctx.RabbitMQURL, ctx.ReplyQueue, ctx.CommandQueue, ctx.Backend)
		}
		fmt.Printf("Current context: %s\n", cfg.CurrentContext)
	},
}

// Use context command to switch contexts
var useContextCmd = &cobra.Command{
	Use:   "use-context [context-name]",
	Short: "Switch to a specific context",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Println("Failed to load config:", err)
			return
		}

		if err := config.UseContext(cfg, args[0]); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Switched to context:", args[0])
	},
}

var deleteContextCmd = &cobra.Command{
	Use:   "delete-context [context-name]",
	Short: "Delete a specified context",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Println("Failed to load config:", err)
			return
		}

		if err := config.DeleteContext(cfg, args[0]); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Context deleted successfully:", args[0])
	},
}

func init() {
	// Add subcommands to config
	configCmd.AddCommand(setContextCmd)
	configCmd.AddCommand(getContextsCmd)
	configCmd.AddCommand(useContextCmd)
	configCmd.AddCommand(deleteContextCmd)

	// Add flags to set-context
	setContextCmd.Flags().StringVarP(&contextName, "name", "n", "", "Name of the context")
	setContextCmd.Flags().StringVarP(&rabbitMQURL, "queueUrl", "u", "", "Queue URL")
	setContextCmd.Flags().StringVarP(&commandQueue, "commandQueue", "c", "", "Command queue name")
	setContextCmd.Flags().StringVarP(&replyQueue, "replyQueue", "r", "", "Reply queue name")
	setContextCmd.Flags().StringVarP(&backend, "backend", "b", "rabbitmq", "Backend type (rabbitmq/sqs/pubsub)")

	// Attach config command to root
	rootCmd.AddCommand(configCmd)
}
