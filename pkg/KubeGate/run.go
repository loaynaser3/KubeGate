package KubeGate

import (
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/loaynaser3/KubeGate/pkg/config"
	"github.com/loaynaser3/KubeGate/pkg/logging"
	"github.com/loaynaser3/KubeGate/pkg/queue"
	"github.com/sirupsen/logrus"
)

// ExecuteRun handles the "run" command
func ExecuteRun(kubeCommand string, args []string) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		logging.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to load config")

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
	logging.Logger.WithFields(logrus.Fields{
		"command":       fullCommand,
		"command_queue": commandQueue,
		"reply_queue":   replyQueue,
	}).Info("Sending command to agent")

	correlationID, err := queue.SendWithReply(conn, commandQueue, fullCommand, replyQueue)
	if err != nil {
		log.Fatalf("Failed to send command: %v", err)
	}

	// Wait for response
	fmt.Println("Command sent. Waiting for response...")
	response, err := queue.WaitForResponse(conn, replyQueue, correlationID, queue.DefaultTimeout)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"correlation": correlationID,
			"error":       err.Error(),
		}).Error("Failed to get response from agent")
	}

	logging.Logger.WithFields(logrus.Fields{
		"correlation": correlationID,
		"response":    response,
	}).Info("Received response from agent")
}
