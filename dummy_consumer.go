package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main_1() {
	rabbitURL := "amqp://guest:guest@localhost:5672/"
	replyQueue := "reply-queue"

	// Connect to RabbitMQ
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare the reply queue
	_, err = ch.QueueDeclare(
		replyQueue, // name
		false,      // durable
		true,       // delete when unused
		true,       // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	// Start consuming messages from the reply queue
	msgs, err := ch.Consume(
		replyQueue, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Printf("Waiting for messages on queue: %s", replyQueue)
	for msg := range msgs {
		fmt.Printf("Received response: %s\n", msg.Body)
	}
}
