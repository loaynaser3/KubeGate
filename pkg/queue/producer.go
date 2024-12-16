package queue

import (
	"github.com/google/uuid"
	"github.com/loaynaser3/KubeGate/pkg/logging"
	"github.com/sirupsen/logrus"
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
	logging.Logger.WithFields(logrus.Fields{
		"queue":       replyQueue,
		"correlation": correlationID,
		"message":     message,
	}).Info("Published message to queue")
	return correlationID, nil
}

// PublishResponse sends a response message to the specified reply queue
func PublishResponse(conn *amqp.Connection, replyQueue, correlationID, message string) error {
	ch, err := conn.Channel()
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"queue": replyQueue,
			"error": err.Error(),
		}).Error("Failed to publish message")
		return err
	}
	defer ch.Close()

	logging.Logger.WithFields(logrus.Fields{
		"queue":       replyQueue,
		"correlation": correlationID,
		"message":     message,
	}).Info("Published message to queue")
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
