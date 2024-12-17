package queue

type Message struct {
	Body          string
	CorrelationID string
	ReplyTo       string
}

type MessageQueue interface {
	Connect() error                                                      // Establish connection to backend
	SendMessage(queueName, message, correlationID, replyTo string) error // Sends a command to the agent
	ReceiveMessages(queueName string, handler func(Message) error) error // Processes messages from a queue (includes a custom handler)
	PublishResponse(replyTo, correlationID, response string) error       // Sends responses back to the reply queue
	Close() error                                                        // Close connection
}
