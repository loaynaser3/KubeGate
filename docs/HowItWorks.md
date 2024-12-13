How It Works Now

    Client Workflow:
        The client:
            Creates a unique replyQueue and declares it.
            Sends the command to the kubegate-commands queue with:
                The replyQueue name in the ReplyTo field.
                A unique CorrelationId.
            Waits on the replyQueue for a response with the matching CorrelationId.

    Agent Workflow:
        The agent:
            Listens for commands on the kubegate-commands queue.
            Processes the command (e.g., by executing kubectl).
            Sends the response back to the replyQueue specified in the ReplyTo field, including the original CorrelationId.

    Response Handling:
        The client receives the response from its replyQueue and verifies the CorrelationId before processing it.