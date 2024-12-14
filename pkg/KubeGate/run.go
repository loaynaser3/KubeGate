package KubeGate

import (
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/loaynaser3/KubeGate/pkg/config"
	"github.com/loaynaser3/KubeGate/pkg/queue"
)

// ExecuteRun handles the "run" command
func ExecuteRun(kubeCommand string, args []string) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Get the current context
	context, err := config.GetContext(cfg, cfg.CurrentContext)
	if err != nil {
		log.Fatalf("Failed to get current context: %v", err)
	}

	// Use the command queue for the current context
	commandQueue := context.CommandQueue
	replyQueue := "reply-queue-" + uuid.New().String()

	// Connect to RabbitMQ
	conn, err := queue.Connect(context.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Send the command to the agent
	fullCommand := strings.Join(append([]string{kubeCommand}, args...), " ")
	correlationID, err := queue.SendWithReply(conn, commandQueue, fullCommand, replyQueue)
	if err != nil {
		log.Fatalf("Failed to send command: %v", err)
	}

	// Wait for response
	fmt.Println("Command sent. Waiting for response...")
	response, err := queue.WaitForResponse(conn, replyQueue, correlationID, queue.DefaultTimeout)
	if err != nil {
		log.Fatalf("Failed to get response: %v", err)
	}

	fmt.Printf("Command Response: %s\n", response)
}
