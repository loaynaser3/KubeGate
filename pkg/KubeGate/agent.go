package KubeGate

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/loaynaser3/KubeGate/pkg/config"
	"github.com/loaynaser3/KubeGate/pkg/queue"
	"github.com/streadway/amqp"
)

// StartAgent initializes the RabbitMQ consumer and starts processing messages
func StartAgent() error {
	// Load agent configuration
	cfg, err := config.LoadAgentConfig()
	if err != nil {
		log.Fatalf("Failed to load agent config: %v", err)
	}

	// Connect to RabbitMQ
	conn, err := queue.Connect(cfg.RabbitMQURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Start consuming messages from the dedicated command queue
	log.Printf("Agent is listening on queue: %s", cfg.CommandQueue)
	return queue.ConsumeMessagesWithCustomHandler(conn, cfg.CommandQueue, func(msg amqp.Delivery) {
		handleCommand(conn, msg)
	})
}

// handleCommand processes a single message and sends a response back to the client
func handleCommand(conn *amqp.Connection, msg amqp.Delivery) {
	log.Printf("Received command: %s\n", string(msg.Body))

	// Execute the command
	command := string(msg.Body)
	result, err := executeKubectlCommand(command)
	if err != nil {
		result = fmt.Sprintf("Error executing command: %v", err)
	}

	// Send response to the reply queue
	if err := queue.PublishResponse(conn, msg.ReplyTo, msg.CorrelationId, result); err != nil {
		log.Printf("Failed to send response: %v", err)
	} else {
		log.Printf("Sent response for correlation ID: %s", msg.CorrelationId)
	}
}

// executeKubectlCommand executes a kubectl command in the agent pod
func executeKubectlCommand(command string) (string, error) {
	// Prepare the kubectl command
	kubectlArgs := strings.Fields(command)
	cmd := exec.Command("kubectl", kubectlArgs...)

	// Capture the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %s, error: %v", string(output), err)
	}

	return string(output), nil
}

// ExecuteCommand SIMULATE executing a kubectl command.
func ExecuteCommand(command string) (string, error) {
	// Simulate processing the command
	// In production, use `client-go` to interact with the Kubernetes API
	fmt.Printf("Executing: %s\n", command)
	return fmt.Sprintf("Simulated result for command: %s", command), nil
}
