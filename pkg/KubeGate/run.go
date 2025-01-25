package KubeGate

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/loaynaser3/KubeGate/pkg/config"
	"github.com/loaynaser3/KubeGate/pkg/logging"
	"github.com/loaynaser3/KubeGate/pkg/queue"
	"github.com/loaynaser3/KubeGate/pkg/utils"
)

// ExecuteRun handles the "run" command
func ExecuteRun(kubeCommand string, args []string) {
	cfg, err := config.LoadConfig()
	if err != nil {
		logging.Logger.Fatalf("Failed to load config: %v", err)
	}

	context, err := config.GetContext(cfg, cfg.CurrentContext)
	if err != nil {
		logging.Logger.Fatalf("Failed to get current context: %v", err)
	}

	// Initialize message queue
	messageQueue, err := queue.NewMessageQueue(context.Backend, context.RabbitMQURL)
	if err != nil {
		logging.Logger.Fatalf("Failed to initialize messaging backend: %v", err)
	}
	defer func() {
		if err := messageQueue.Close(); err != nil {
			logging.Logger.Printf("Warning: Failed to close message queue: %v", err)
		}
	}()

	if err := messageQueue.Connect(); err != nil {
		logging.Logger.Fatalf("Failed to connect to messaging backend: %v", err)
	}

	// Initialize session manager
	sessionManager := utils.NewSessionManager("~/.kubegate/temp.json", time.Hour, messageQueue)

	// Get or create the reply queue (includes cleanup logging.Loggeric)
	replyQueue, err := sessionManager.GetOrCreateReplyQueue()
	if err != nil {
		logging.Logger.Fatalf("Failed to manage reply queue: %v", err)
	}

	correlationID := uuid.New().String()

	// Manually parse `-f` or `--file` flag
	encodedArgs, err := utils.ReplaceFileWithBase64(args, utils.EncodeFileToBase64String)
	if err != nil {
		logging.Logger.Fatalf("Failed to encode message: %v", err)
	}

	// Send command to agent
	fullCommand := strings.Join(append([]string{kubeCommand}, encodedArgs...), " ")
	err = messageQueue.SendMessage(context.CommandQueue, fullCommand, correlationID, replyQueue)
	if err != nil {
		logging.Logger.Fatalf("Failed to send command: %v", err)
	}

	// fmt.Println("Command sent. Waiting for response...")

	// Wait for response with timeout
	responseChan := make(chan string, 1)
	go func() {
		err := messageQueue.ReceiveMessages(replyQueue, func(msg queue.Message) error {
			if msg.CorrelationID == correlationID {
				responseChan <- msg.Body
			} else {
				fmt.Printf("Skipped message: CorrelationID=%s (expected %s)\n", msg.CorrelationID, correlationID)
			}
			return nil
		})
		if err != nil {
			logging.Logger.Printf("Error receiving messages from queue %s: %v", replyQueue, err)
		}
	}()

	select {
	case response := <-responseChan:
		fmt.Printf("Response received: %s\n", response)
	case <-time.After(60 * time.Second): // Add timeout to avoid indefinite waiting
		logging.Logger.Fatalf("Timeout waiting for response.")
	}
}
