package KubeGate

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/loaynaser3/KubeGate/pkg/queue"
	"github.com/streadway/amqp"
)

// StartAgent initializes the RabbitMQ consumer and starts processing messages
func StartAgent() error {
	rabbitURL := "amqp://guest:guest@localhost:5672/"
	conn, err := queue.Connect(rabbitURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	log.Println("Agent is listening for commands...")
	return queue.ConsumeMessages(conn, "kubegate-commands", handleCommand)
}

func handleCommand(ch *amqp.Channel, msg amqp.Delivery) {
	log.Printf("Received command: %s\n", string(msg.Body))

	// Execute the command
	command := string(msg.Body)
	result, err := executeKubectlCommand(command)
	if err != nil {
		result = fmt.Sprintf("Error executing command: %v", err)
	}

	// Send response to reply queue
	err = ch.Publish(
		"",          // exchange
		msg.ReplyTo, // reply queue
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			Body:          []byte(result),
			CorrelationId: msg.CorrelationId, // Match request correlation ID
		},
	)
	if err != nil {
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
	println(output)
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %s, error: %v", string(output), err)
	}

	return string(output), nil
}

// ExecuteCommand simulates executing a kubectl command.
func ExecuteCommand(command string) (string, error) {
	// Simulate processing the command
	// In production, use `client-go` to interact with the Kubernetes API
	fmt.Printf("Executing: %s\n", command)
	return fmt.Sprintf("Simulated result for command: %s", command), nil
}
