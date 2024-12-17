package queue

import (
	"fmt"
)

func NewMessageQueue(backend string, url string) (MessageQueue, error) {
	switch backend {
	case "rabbitmq":
		return &RabbitMQ{URL: url}, nil
	case "sqs":
		// Placeholder for AWS SQS implementation
		return nil, fmt.Errorf("SQS not implemented yet")
	case "pubsub":
		// Placeholder for GCP Pub/Sub implementation
		return nil, fmt.Errorf("Pub/Sub not implemented yet")
	case "azure":
		// Placeholder for Azure Service Bus implementation
		return nil, fmt.Errorf("Azure Service Bus not implemented yet")
	default:
		return nil, fmt.Errorf("unknown backend: %s", backend)
	}
}
