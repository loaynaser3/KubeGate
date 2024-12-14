package queue

import (
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

const DefaultTimeout = 30 // Timeout for waiting for responses in seconds

// Connect establishes a connection to RabbitMQ
func Connect(rabbitMQURL string) (*amqp.Connection, error) {
	return amqp.Dial(rabbitMQURL)
}

// SendWithReply sends a command to RabbitMQ with a reply-to queue
func SendWithReply(conn *amqp.Connection, queueName string, message string, replyQueue string) (string, error) {
	ch, err := conn.Channel()
	if err != nil {
		return "", err
	}
	defer ch.Close()

	// Generate a unique correlation ID
	correlationID := uuid.New().String()

	// Publish the message
	err = ch.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			Body:          []byte(message),
			CorrelationId: correlationID,
			ReplyTo:       replyQueue,
		},
	)
	if err != nil {
		return "", err
	}

	return correlationID, nil
}

// PublishResponse sends a response message to the specified reply queue
func PublishResponse(conn *amqp.Connection, replyQueue, correlationID, message string) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return ch.Publish(
		"",         // exchange
		replyQueue, // reply queue
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			Body:          []byte(message),
			CorrelationId: correlationID,
		},
	)
}
