package KubeGate

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/loaynaser3/KubeGate/pkg/config"
	"github.com/loaynaser3/KubeGate/pkg/logging"
	"github.com/loaynaser3/KubeGate/pkg/queue"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// StartAgent initializes the RabbitMQ consumer and starts processing messages
func StartAgent() error {
	// Load agent configuration
	cfg, err := config.LoadAgentConfig()
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to load agent config")
	}

	// Connect to RabbitMQ
	conn, err := queue.Connect(cfg.RabbitMQURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Start consuming messages from the dedicated command queue
	log.Printf("Agent is listening on queue: %s", cfg.CommandQueue)
	logging.Logger.WithFields(logrus.Fields{
		"rabbitmq_url":  cfg.RabbitMQURL,
		"command_queue": cfg.CommandQueue,
	}).Info("Agent started and connected to RabbitMQ")
	return queue.ConsumeMessagesWithCustomHandler(conn, cfg.CommandQueue, func(msg amqp.Delivery) {
		handleCommand(conn, msg)
	})
}

// handleCommand processes a single message and sends a response back to the client
func handleCommand(conn *amqp.Connection, msg amqp.Delivery) {
	logging.Logger.WithFields(logrus.Fields{
		"command":     string(msg.Body),
		"reply_queue": msg.ReplyTo,
		"correlation": msg.CorrelationId,
	}).Info("Received command from client")

	// Execute the command
	command := string(msg.Body)
	logging.Logger.WithField("command", command).Info("Executing kubectl command")
	result, err := executeKubectlCommand(command)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"command": command,
			"error":   err.Error(),
		}).Error("Failed to execute kubectl command")
	}

	// Send response to the reply queue
	if err := queue.PublishResponse(conn, msg.ReplyTo, msg.CorrelationId, result); err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"correlation": msg.CorrelationId,
			"reply_queue": msg.ReplyTo,
			"Error":       err,
		}).Error("Failed to send response")

	} else {
		logging.Logger.WithFields(logrus.Fields{
			"correlation": msg.CorrelationId,
			"reply_queue": msg.ReplyTo,
		}).Info("Response sent to client")

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
		logging.Logger.WithFields(logrus.Fields{
			"command": string(output),
			"Error":   err,
		}).Error("failed to execute command")
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
