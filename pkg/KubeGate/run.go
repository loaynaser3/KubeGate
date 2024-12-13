package KubeGate

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/loaynaser3/KubeGate/pkg/queue"
)

// ExecuteRun handles the "run" command
func ExecuteRun(kubeCommand string) {
	rabbitURL := "amqp://guest:guest@localhost:5672/"
	conn, err := queue.Connect(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Use a unique reply queue for this client
	replyQueue := "reply-queue-" + uuid.New().String()

	// Declare the reply queue
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to create channel: %v", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		replyQueue, // name
		false,      // durable
		true,       // delete when unused
		true,       // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare reply queue: %v", err)
	}

	// Send the command
	correlationID, err := queue.SendWithReply(conn, "kubegate-commands", kubeCommand, replyQueue)
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
