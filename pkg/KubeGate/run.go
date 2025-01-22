package KubeGate

import (
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/loaynaser3/KubeGate/pkg/config"
	"github.com/loaynaser3/KubeGate/pkg/queue"
	"github.com/loaynaser3/KubeGate/pkg/utils"
)

var (
	filePath    string // To store the value of the `-f` flag
	encodedFile string
)

// ExecuteRun handles the "run" command
func ExecuteRun(kubeCommand string, args []string) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	context, err := config.GetContext(cfg, cfg.CurrentContext)
	if err != nil {
		log.Fatalf("Failed to get current context: %v", err)
	}

	// Initialize message queue
	messageQueue, err := queue.NewMessageQueue(context.Backend, context.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to initialize messaging backend: %v", err)
	}
	defer messageQueue.Close()

	if err := messageQueue.Connect(); err != nil {
		log.Fatalf("Failed to connect to messaging backend: %v", err)
	}

	replyQueue := "reply-queue-" + uuid.New().String()
	correlationID := uuid.New().String()

	// Manually parse `-f` or `--file` flag
	encodedArgs, err := utils.ReplaceFileWithBase64(args, utils.EncodeFileToBase64String)
	if err != nil {
		log.Fatalf("Failed to encode message: %v", err)
	}
	// Send command to agent
	fullCommand := strings.Join(append([]string{kubeCommand}, encodedArgs...), " ")
	err = messageQueue.SendMessage(context.CommandQueue, fullCommand, correlationID, replyQueue)
	if err != nil {
		log.Fatalf("Failed to send command: %v", err)
	}

	fmt.Println("Command sent. Waiting for response...")

	// Wait for response
	responseChan := make(chan string, 1)
	go func() {
		messageQueue.ReceiveMessages(replyQueue, func(msg queue.Message) error {
			if msg.CorrelationID == correlationID {
				responseChan <- msg.Body
			}
			return nil
		})
	}()

	response := <-responseChan
	fmt.Printf("Response: %s\n", response)
}
