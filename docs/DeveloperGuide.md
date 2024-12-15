# Developer Guide for KubeGate

Welcome to the KubeGate development guide! This document provides detailed instructions for developers to contribute to KubeGate and understand its codebase.

---

## Prerequisites

Before working on KubeGate, ensure you have the following installed:

- **Go**: Version 1.20 or later.
- **RabbitMQ**: A running RabbitMQ instance for testing.
- **kubectl**: Installed and configured for Kubernetes testing.
- **Docker**: For running the agent in a containerized environment (optional).
- **Git**: For version control.

---

## Project Structure

- **`cmd/`**: Contains the CLI commands for KubeGate.
  - `run.go`: Handles the `run` command logic.
  - `config.go`: Manages configuration-related commands.
- **`pkg/KubeGate/`**: Core business logic.
  - `run.go`: Implements the `ExecuteRun` function.
  - `agent.go`: Handles agent-side message processing.
- **`pkg/queue/`**: RabbitMQ utilities.
  - `consumer.go`: Handles message consumption.
  - `producer.go`: Handles message publishing.
- **`pkg/config/`**: Configuration management.
  - `config.go`: Handles client-side configuration.
  - `agent.go`: Manages agent-specific configuration.

---

## Getting Started

### 1. Clone the Repository
```bash
git clone https://github.com/your-org/kubegate.git
cd kubegate
```

### 2. Set Up Environment Variables
For the agent:
```bash
export KUBEGATE_RABBITMQ_URL=amqp://guest:guest@localhost:5672/
export KUBEGATE_RABBITMQ_QUEUE=kubegate-commands-dev
```

### 3. Run the Agent
```bash
go run main.go agent
```

### 4. Run the Client
```bash
go run main.go run get pods -n default
```

---

## Development Workflow

### Adding a New Feature

1. **Define the Feature**:
   - For client-side logic, create or update files under `cmd/` or `pkg/KubeGate/`.
   - For agent-side logic, modify files in `pkg/KubeGate/` or `pkg/queue/`.

2. **Implement the Feature**:
   - Use the existing structure for consistency.
   - For RabbitMQ interactions, add utilities to `consumer.go` or `producer.go` if needed.

3. **Write Tests**:
   - Add unit tests for new functions.
   - Use the `testing` package in Go.

4. **Run Tests**:
   ```bash
   go test ./...
   ```

5. **Commit Changes**:
   ```bash
   git add .
   git commit -m "Add [feature]"
   ```

6. **Push and Open a Pull Request**:
   ```bash
   git push origin feature-branch
   ```

### Coding Conventions
- Follow Go best practices.
- Use `log` for logging.
- Keep functions small and focused.

---

## Testing

### Unit Tests
Run all unit tests:
```bash
go test ./...
```

### Manual Testing
1. Start RabbitMQ.
2. Run the agent.
3. Use the client to send commands.
4. Verify responses.

### Example Manual Test
```bash
kubeGate run get pods -n kube-system
```
Expected Output:
```bash
Command Response:
NAME               READY   STATUS    RESTARTS   AGE
coredns            2/2     Running   0          25m
```

---

## Contributing

1. Fork the repository.
2. Create a feature branch:
   ```bash
   git checkout -b feature/new-feature
   ```
3. Make your changes and commit them.
4. Open a pull request.

---

## TODO for Developers
- Improve test coverage.
- Enhance config managment.
- Implement an interactive shell mode.
- Enhance logging with structured logs.

---

Thank you for contributing to KubeGate! 🚀

