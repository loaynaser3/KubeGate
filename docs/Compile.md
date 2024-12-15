# Compiling and Running KubeGate

This guide explains how to compile the KubeGate code into an executable binary and run both the agent and client.

---

## Step 1: Compile the Code

1. **Navigate to the project directory**:
   ```bash
   cd kubegate
   ```

2. **Build the binary**:
   ```bash
   go build -o kubeGate main.go
   ```
   - This will generate an executable file named `kubeGate` in the current directory.

3. **Verify the build**:
   - Check if the `kubeGate` binary is present:
     ```bash
     ls kubeGate
     ```

---

## Step 2: Run the Agent

1. **Set the required environment variables**:
   ```bash
   export KUBEGATE_RABBITMQ_URL=amqp://guest:guest@localhost:5672/
   export KUBEGATE_RABBITMQ_QUEUE=kubegate-commands-dev
   ```

2. **Start the agent**:
   ```bash
   ./kubeGate agent
   ```

3. **Verify the agent is running**:
   - Check the logs to confirm it is consuming messages from the RabbitMQ queue.

---

## Step 3: Run the Client

1. Use the `run` command to execute a Kubernetes command:
   ```bash
   ./kubeGate run get pods -n default
   ```

   - Replace `get pods -n default` with your desired `kubectl` command.

2. **Verify the output**:
   - The response from the agent should appear in the terminal.

---

## Troubleshooting

- **Binary Not Found**:
  - Ensure `go build` completed successfully.
  - Verify Go is installed and accessible via the `PATH` environment variable.

- **RabbitMQ Connection Issues**:
  - Ensure RabbitMQ is running and the URL (`KUBEGATE_RABBITMQ_URL`) is correct.
  - Check for firewall or network restrictions.

- **Agent Not Responding**:
  - Confirm the agent is running and processing commands from the queue.

---

This process ensures you have a fully operational KubeGate binary for both the agent and client components. 🚀

