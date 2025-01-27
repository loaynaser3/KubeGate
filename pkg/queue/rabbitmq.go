package queue

import (
	"fmt"
	"strings"

	"github.com/loaynaser3/KubeGate/pkg/logging"
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
	logging.Logger.Info("Connected to RabbitMQ successfully")
	return nil
}

// DeclareQueue ensures a queue exists with consistent attributes
func (r *RabbitMQ) DeclareQueue(queueName string) error {
	_, err := r.Channel.QueueDeclare(
		queueName,
		true,  // Durable
		false, // Auto-delete
		false, // Exclusive
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue %s: %v", queueName, err)
	}

	logging.Logger.WithField("queue", queueName).Info("Queue declared successfully")
	return nil
}

// CreateQueue ensures a queue exists with consistent attributes
func (r *RabbitMQ) CreateQueue(queueName string) error {
	_, err := r.Channel.QueueDeclare(
		queueName, // Queue name
		true,      // Durable
		false,     // Auto-delete
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue %s: %v", queueName, err)
	}

	logging.Logger.WithField("queue", queueName).Info("Queue created successfully")
	return nil
}

// SendMessage publishes a message to a specified queue
func (r *RabbitMQ) SendMessage(queueName, message, correlationID, replyTo string) error {
	// Ensure the queue exists
	if err := r.DeclareQueue(queueName); err != nil {
		return fmt.Errorf("failed to ensure queue exists: %v", err)
	}

	// Publish the message
	err := r.Channel.Publish(
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

	logging.Logger.WithFields(logrus.Fields{
		"queue":        queueName,
		"correlation":  correlationID,
		"message_size": len(message),
	}).Info("Message sent to queue")
	return nil
}

// ReceiveMessages consumes messages from a specified queue and invokes the handler
func (r *RabbitMQ) ReceiveMessages(queueName string, handler func(Message) error) error {
	// Ensure the queue exists
	if err := r.DeclareQueue(queueName); err != nil {
		return fmt.Errorf("failed to ensure queue exists: %v", err)
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

	logging.Logger.WithField("queue", queueName).Info("Started listening to queue")

	// Process messages using the provided handler
	for msg := range msgs {
		message := Message{
			Body:          string(msg.Body),
			CorrelationID: msg.CorrelationId,
			ReplyTo:       msg.ReplyTo,
		}
		if err := handler(message); err != nil {
			logging.Logger.WithFields(logrus.Fields{
				"error":       err.Error(),
				"correlation": msg.CorrelationId,
			}).Error("Failed to process message")
		}
	}
	return nil
}

// PublishResponse sends a response message to the reply queue
func (r *RabbitMQ) PublishResponse(replyTo, correlationID, response string) error {
	if replyTo == "" {
		return fmt.Errorf("replyTo queue name is empty")
	}

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
		return fmt.Errorf("failed to publish response to queue %s: %v", replyTo, err)
	}

	// Log the response with additional metadata
	truncatedResponse := response
	if len(response) > 100 {
		truncatedResponse = response[:100] + "..."
	}

	logging.Logger.WithFields(logrus.Fields{
		"reply_queue":  replyTo,
		"correlation":  correlationID,
		"response_len": len(response),
		"response":     truncatedResponse,
	}).Info("Response published to reply queue")

	return nil
}

// DeleteQueue deletes a RabbitMQ queue
func (r *RabbitMQ) DeleteQueue(queueName string) error {
	ch, err := r.Conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()

	_, err = ch.QueueDelete(
		queueName, // Queue name
		false,     // IfUnused
		false,     // IfEmpty
		false,     // NoWait
	)
	if err != nil {
		return fmt.Errorf("failed to delete queue %s: %w", queueName, err)
	}

	logging.Logger.WithField("queue", queueName).Info("Queue deleted successfully")
	return nil
}

// Close closes the RabbitMQ connection and channel
func (r *RabbitMQ) Close() error {
	var closeErrors []string

	// Close the channel
	if r.Channel != nil {
		if err := r.Channel.Close(); err != nil && !isChannelAlreadyClosedError(err) {
			closeErrors = append(closeErrors, fmt.Sprintf("Failed to close channel: %v", err))
		}
	}

	// Close the connection
	if r.Conn != nil {
		if err := r.Conn.Close(); err != nil && !isConnectionAlreadyClosedError(err) {
			closeErrors = append(closeErrors, fmt.Sprintf("Failed to close connection: %v", err))
		}
	}

	// Log and return any errors encountered during closure
	if len(closeErrors) > 0 {
		errMsg := strings.Join(closeErrors, "; ")
		logging.Logger.WithError(fmt.Errorf(errMsg)).Error("Failed to close RabbitMQ resources")
		return fmt.Errorf(errMsg)
	}

	logging.Logger.Info("RabbitMQ connection closed successfully")
	return nil
}

func isChannelAlreadyClosedError(err error) bool {
	return strings.Contains(err.Error(), "channel/connection is not open")
}

func isConnectionAlreadyClosedError(err error) bool {
	return strings.Contains(err.Error(), "channel/connection is not open")
}
