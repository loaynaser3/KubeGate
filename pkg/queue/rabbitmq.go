package queue

import (
	"fmt"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// RabbitMQ implements the MessageQueue interface for RabbitMQ
type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	URL     string
}

// Connect establishes a connection to RabbitMQ and opens a channel
func (r *RabbitMQ) Connect() error {
	conn, err := amqp.Dial(r.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	r.Conn = conn

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}
	r.Channel = ch
	logrus.Info("Connected to RabbitMQ successfully")
	return nil
}

// SendMessage publishes a message to a specified queue
func (r *RabbitMQ) SendMessage(queueName, message, correlationID, replyTo string) error {
	// Ensure the queue exists
	_, err := r.Channel.QueueDeclare(
		queueName, // Queue name
		false,     // Durable
		false,     // Delete when unused
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %v", err)
	}

	// Publish the message
	err = r.Channel.Publish(
		"",        // Exchange
		queueName, // Routing key
		false,     // Mandatory
		false,     // Immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			Body:          []byte(message),
			CorrelationId: correlationID,
			ReplyTo:       replyTo,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	logrus.WithFields(logrus.Fields{
		"queue":       queueName,
		"correlation": correlationID,
	}).Info("Message sent to queue")
	return nil
}

// ReceiveMessages consumes messages from a specified queue and invokes the handler
func (r *RabbitMQ) ReceiveMessages(queueName string, handler func(Message) error) error {
	// Ensure the queue exists
	_, err := r.Channel.QueueDeclare(
		queueName, // Queue name
		false,     // Durable
		false,     // Delete when unused
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %v", err)
	}

	// Start consuming messages
	msgs, err := r.Channel.Consume(
		queueName, // Queue name
		"",        // Consumer tag
		true,      // Auto-ack
		false,     // Exclusive
		false,     // No-local
		false,     // No-wait
		nil,       // Args
	)
	if err != nil {
		return fmt.Errorf("failed to consume messages: %v", err)
	}

	logrus.WithField("queue", queueName).Info("Started listening to queue")

	// Process messages using the provided handler
	for msg := range msgs {
		message := Message{
			Body:          string(msg.Body),
			CorrelationID: msg.CorrelationId,
			ReplyTo:       msg.ReplyTo,
		}
		if err := handler(message); err != nil {
			logrus.WithFields(logrus.Fields{
				"error":       err.Error(),
				"correlation": msg.CorrelationId,
			}).Error("Failed to process message")
		}
	}
	return nil
}

// PublishResponse sends a response message to the reply queue
func (r *RabbitMQ) PublishResponse(replyTo, correlationID, response string) error {
	err := r.Channel.Publish(
		"",      // Exchange
		replyTo, // Reply queue
		false,   // Mandatory
		false,   // Immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			Body:          []byte(response),
			CorrelationId: correlationID,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish response: %v", err)
	}

	logrus.WithFields(logrus.Fields{
		"reply_queue": replyTo,
		"correlation": correlationID,
		"response":    response,
	}).Info("Response published to reply queue")
	return nil
}

// Close closes the RabbitMQ connection and channel
func (r *RabbitMQ) Close() error {
	if r.Channel != nil {
		if err := r.Channel.Close(); err != nil {
			log.Printf("Failed to close channel: %v", err)
		}
	}
	if r.Conn != nil {
		if err := r.Conn.Close(); err != nil {
			log.Printf("Failed to close connection: %v", err)
		}
	}
	logrus.Info("RabbitMQ connection closed")
	return nil
}
