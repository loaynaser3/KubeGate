package queue

import (
	"fmt"
	"log"
	"time"

	"github.com/loaynaser3/KubeGate/pkg/logging"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// ConsumeMessages starts consuming messages from a queue - not used currently
func ConsumeMessages(conn *amqp.Connection, queueName string, handler func(ch *amqp.Channel, msg amqp.Delivery)) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Declare the queue
	_, err = ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	// Start consuming messages
	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Printf("Error consuming messages: %v", err)
		return err
	}

	// Process messages
	for msg := range msgs {
		go handler(ch, msg) // Handle each message concurrently
	}

	return nil
}

// WaitForResponse listens for a response on a reply queue with a specific correlation ID
func WaitForResponse(conn *amqp.Connection, replyQueue string, correlationID string, timeout int) (string, error) {
	ch, err := conn.Channel()
	if err != nil {
		return "", err
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
		return "", fmt.Errorf("failed to declare reply queue: %v", err)
	}

	// Consume messages from the reply queue
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
		return "", fmt.Errorf("failed to consume messages: %v", err)
	}

	// Wait for the response
	timeoutChan := time.After(time.Duration(timeout) * time.Second)
	for {
		select {
		case msg := <-msgs:
			if msg.CorrelationId == correlationID {
				return string(msg.Body), nil
			}
		case <-timeoutChan:
			return "", fmt.Errorf("timeout waiting for response")
		}
	}
}

// ConsumeMessagesWithCustomHandler consumes messages and processes them using a custom handler
func ConsumeMessagesWithCustomHandler(conn *amqp.Connection, queueName string, handler func(msg amqp.Delivery)) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Declare the queue
	_, err = ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	// Start consuming messages
	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	logging.Logger.WithField("queue", queueName).Info("Started listening to queue")

	for msg := range msgs {
		logging.Logger.WithFields(logrus.Fields{
			"message_body": string(msg.Body),
			"delivery_tag": msg.DeliveryTag,
		}).Debug("Message consumed")

		go handler(msg) // Process each message concurrently
	}

	return nil
}
