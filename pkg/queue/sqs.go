package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// SQS implementation of the MessageQueue interface
type SQS struct {
	Client   *sqs.Client   // AWS SQS client for interacting with the service
	QueueURL string        // URL of the SQS queue
	Config   aws.Config    // AWS configuration for the client
	Timeout  time.Duration // Timeout for receiving messages
}

// Connect initializes the SQS client and validates the connection
func (s *SQS) Connect() error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %v", err)
	}
	s.Config = cfg
	s.Client = sqs.NewFromConfig(cfg)
	return nil
}

// SendMessage sends a message to the SQS queue
func (s *SQS) SendMessage(queueName, message, correlationID, replyTo string) error {
	input := &sqs.SendMessageInput{
		QueueUrl:    &s.QueueURL,
		MessageBody: aws.String(message),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"CorrelationID": {DataType: aws.String("String"), StringValue: aws.String(correlationID)},
			"ReplyTo":       {DataType: aws.String("String"), StringValue: aws.String(replyTo)},
		},
	}
	_, err := s.Client.SendMessage(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	return nil
}

// ReceiveMessages listens for messages from the SQS queue
func (s *SQS) ReceiveMessages(queueName string, handler func(Message) error) error {
	for {
		input := &sqs.ReceiveMessageInput{
			QueueUrl:            &s.QueueURL,
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     int32(s.Timeout.Seconds()),
			MessageAttributeNames: []string{
				"All",
			},
		}
		output, err := s.Client.ReceiveMessage(context.TODO(), input)
		if err != nil {
			return fmt.Errorf("failed to receive messages: %v", err)
		}

		for _, msg := range output.Messages {
			message := Message{
				Body:          aws.ToString(msg.Body),
				CorrelationID: aws.ToString(msg.MessageAttributes["CorrelationID"].StringValue),
				ReplyTo:       aws.ToString(msg.MessageAttributes["ReplyTo"].StringValue),
			}
			if err := handler(message); err != nil {
				return fmt.Errorf("failed to handle message: %v", err)
			}

			// Delete the message after successful processing
			_, err := s.Client.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
				QueueUrl:      &s.QueueURL,
				ReceiptHandle: msg.ReceiptHandle,
			})
			if err != nil {
				return fmt.Errorf("failed to delete message: %v", err)
			}
		}
	}
}

// PublishResponse sends a response message to the reply queue
func (s *SQS) PublishResponse(replyTo, correlationID, response string) error {
	input := &sqs.SendMessageInput{
		QueueUrl:    &replyTo,
		MessageBody: aws.String(response),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"CorrelationID": {DataType: aws.String("String"), StringValue: aws.String(correlationID)},
		},
	}
	_, err := s.Client.SendMessage(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}
	return nil
}

// Close cleans up resources (no-op for SQS)
func (s *SQS) Close() error {
	return nil
}
