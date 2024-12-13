Client (Run Command)

Send Command:
    ExecuteRun creates a unique replyQueue and sends the command with ReplyTo and CorrelationId.
```go
replyQueue := "reply-queue-" + uuid.New().String()
correlationID, err := queue.SendWithReply(conn, "kubegate-commands", kubeCommand, replyQueue)
```

Wait for Response:

    The client listens to its replyQueue for a message with the matching CorrelationId:
```go
case msg := <-msgs:
if msg.CorrelationId == correlationID {
    return string(msg.Body), nil
}

```

Agent

    Process Command:
        The agent reads the command from kubegate-commands and processes it using kubectl.

    Send Response:
        The agent sends the response to the replyQueue specified in ReplyTo with the same CorrelationId:

```go
err := ch.Publish(
    "",           // exchange
    msg.ReplyTo,  // reply queue
    false,        // mandatory
    false,        // immediate
    amqp.Publishing{
        ContentType:   "text/plain",
        Body:          []byte(result),
        CorrelationId: msg.CorrelationId,
    },
)

```        