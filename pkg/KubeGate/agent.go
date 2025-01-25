package KubeGate

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/loaynaser3/KubeGate/pkg/config"
	"github.com/loaynaser3/KubeGate/pkg/logging"
	"github.com/loaynaser3/KubeGate/pkg/queue"
	"github.com/loaynaser3/KubeGate/pkg/utils"
	"github.com/sirupsen/logrus"
)

// StartAgent initializes the RabbitMQ consumer and starts processing messages
func StartAgent() error {
	cfg, err := config.LoadAgentConfig()
	if err != nil {
		return fmt.Errorf("failed to load agent config: %v", err)
	}

	// Initialize the messaging backend
	messageQueue, err := queue.NewMessageQueue(cfg.Backend, cfg.RabbitMQURL)
	if err != nil {
		return fmt.Errorf("failed to initialize messaging backend: %v", err)
	}
	defer func() {
		if err := messageQueue.Close(); err != nil {
			logging.Logger.WithError(err).Error("Failed to close message queue")
		}
	}()

	// Connect to backend
	if err := messageQueue.Connect(); err != nil {
		return fmt.Errorf("failed to connect to messaging backend: %v", err)
	}

	logging.Logger.WithFields(logrus.Fields{
		"backend":       cfg.Backend,
		"command_queue": cfg.CommandQueue,
	}).Info("Agent started and listening on command queue")

	// Start consuming messages from the command queue
	return messageQueue.ReceiveMessages(cfg.CommandQueue, func(msg queue.Message) error {
		return handleCommand(msg, messageQueue)
	})
}

// handleCommand processes a single message and sends a response back to the client
func handleCommand(msg queue.Message, mq queue.MessageQueue) error {
	logging.Logger.WithFields(logrus.Fields{
		"command":     msg.Body,
		"reply_queue": msg.ReplyTo,
		"correlation": msg.CorrelationID,
	}).Info("Received command from client")

	// Decode Base64 arguments and prepare the command
	decodedArgs, err := utils.ReplaceBase64WithFile(strings.Split(msg.Body, " "), utils.DecodeBase64StringToFile)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"command": msg.Body,
		}).Error("Failed to decode Base64 arguments")
		return err
	}

	// Execute the command
	logging.Logger.WithField("command", decodedArgs).Info("Executing kubectl command")
	result, err := executeKubectlCommand(strings.Join(decodedArgs, " "))
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"command": decodedArgs,
			"error":   err.Error(),
		}).Error("Failed to execute kubectl command")
		result = fmt.Sprintf("Error: %v", err)
	}

	// Send response using the messaging backend
	if err := mq.PublishResponse(msg.ReplyTo, msg.CorrelationID, result); err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"correlation": msg.CorrelationID,
			"reply_queue": msg.ReplyTo,
			"error":       err.Error(),
		}).Error("Failed to send response to client")
		return err
	}

	logging.Logger.WithFields(logrus.Fields{
		"correlation": msg.CorrelationID,
		"reply_queue": msg.ReplyTo,
	}).Info("Response sent to client")
	return nil
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
			"error":   err,
		}).Error("Failed to execute command")
		return "failed to execute command: %s, error: %v", fmt.Errorf("failed to execute command: %s, error: %v", string(output), err)
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
