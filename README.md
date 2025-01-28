# KubeGate

## Description
KubeGate is a Kubernetes proxy tool that simplifies secure and efficient command execution on private clusters. It bridges the gap between remote clients and cluster resources using RabbitMQ for seamless communication, ensuring multi-tenant compatibility and robust role-based permissions.

With KubeGate, you can execute Kubernetes commands remotely while maintaining control over agent access and command handling. Its dynamic configuration and context management make it a versatile tool for DevOps workflows.

## Features
- **Queue Segmentation**: Ensures isolated command handling for multiple agents and contexts.
- **Dynamic RabbitMQ Configuration**:
  - **Agent**: Supports RabbitMQ details via environment variables or YAML configuration.
  - **Client**: Supports multiple contexts, similar to `kubectl`, with configurable RabbitMQ queues.
- **Command Execution**: Executes `kubectl` commands on behalf of the client using agents running inside the target clusters.
- **Response Management**: Uses unique correlation IDs for each client request, ensuring responses are routed back to the correct client.

## KubeGate Diagram

![KubeGate Diagram](docs/KubeGate-Diagram.jpeg)

The diagram illustrates how the client interacts with the RabbitMQ messaging system, communicates commands to the agent, and receives responses securely.

## Build Package and Deployment

### Prerequisites
- GoLang Version 1.23 or higher.
- Docker installed.
- Access to a Kubernetes cluster.
- RabbitMQ instance (e.g., set up locally or via a cloud provider like CloudAMQP). *you can get a free dev instance from [CloudAMQ](https://www.cloudamqp.com/) free.*
 

### Build Client Binary Locally
1. Run the build command:
   ```bash 
   go build -o kubegate .
  *Alternatively, download the compiled package from the [releases page](https://github.com/loaynaser3/KubeGate/releases).*

2. Verify KubeGate is running:
```bash
  ./kubegate --help
  ```
### Build Docker image
 ```bash
  docker build -t kubegate:latest .
  ```
### Deploy Agent on K8S
```bash
cd helm/kubegate-agent
# Edit values.yaml env section with relevant RabbitMQ details
# Easiest way to get a RabbitMQ dev instance is via CloudAMQP for free (suitable for personal/dev use)
# Change the image repository or keep the default if accessible
helm install kubegate . -f values.yaml
```

## Supported Commands
- Run Kubernetes commands:
  ```bash
  kubeGate run [command] [args...]
  ```
  or 
  ```bash
  kubeGate [command] [args...]
  ```
- Alias:

  ```bash 
  alias kubectl="path/to/kubegate"
  ```

### Examples
1. **Basic Command Execution**:
   ```bash
   kubeGate run get pods -n default
   ```
   Executes `kubectl get pods -n default` in the current context.

2. **Switching Contexts:**
   ```bash
   kubeGate config set-context dev --rabbitmq-url amqp://user:pass@dev-rabbitmq:5672/ --command-queue kubegate-commands-dev
   kubeGate config use-context dev
   ```

3. **List Available Contexts**:
   ```bash
   kubeGate config get-contexts
   ```

## TODO
### Current
- **Interactive Shell**: Add support for running multiple commands in a single session.
- **Enhanced Security**: Explore TLS for RabbitMQ connections and stricter command validation.
- **Prometheus Metrics**: Export metrics such as command execution count and latency.
- **Error Logging**: Enhance error visibility and structured logging.
- **Advanced Configuration**: Add support for namespace-based configurations.


---

KubeGate simplifies Kubernetes operations in private clusters by bridging the gap between clients and agents. Contributions and suggestions are always welcome!


## Planned Folder Structure:

```
KubeGate/
├── cmd/
│   ├── agent.go            # Command to start the agent
│   ├── config.go           # Add or update a context       
│   ├── root.go             # Root command, initializes Kobra and subcommands
│   ├── run.go              # Single command execution logic
│
├── pkg/
│   ├── k8s/
│   │   └── k8s.go          # Kubernetes interaction logic
│   ├── queue/
│   │   ├── producer.go     # Logic for publishing messages to RabbitMQ
│   │   └── consumer.go     # Logic for consuming messages from RabbitMQ
│   ├── KubeGate/
│   │   ├── agent.go        # Logic for Agent
│   │   └── run.go          # Logic for run command
│   ├── shell/
│   │   └── interactive.go  # Logic for handling interactive commands
│   ├── config/
│       └── agent.go       # Configuration loader for agent (for RabbitMQ)
│       └── config.go       # Configuration loader for client (e.g., for RabbitMQ, Kubernetes contexts)
├── internal/
│   ├── logger/
│   ├── utils/
│
├── scripts/
│   ├── setup-rabbitmq.sh   # Script to set up RabbitMQ locally or in a cluster
│   └── deploy-agent.yaml   # Kubernetes manifest for deploying the agent
├── docs/
│   ├── README.md           # Project documentation
│   ├── CONTRIBUTING.md     # Contribution guidelines
│   └── CLI_USAGE.md        # Documentation for CLI usage
├── go.mod                  # Go module file
├── go.sum                  # Dependency lock file
└── main.go                 # Entry point for the application
```

Description of Components

    cmd/:
        Contains Kobra commands.
        Each subcommand (e.g., single, shell, agent) is handled in its own file for modularity.

    pkg/:
        Reusable packages for core functionality.
        k8s/: Handles interactions with Kubernetes via client-go.
        queue/: Implements RabbitMQ producer and consumer logic.
        shell/: Manages logic for the interactive shell mode.
        config/: Centralized configuration management.

    internal/:
        Contains internal utilities like logging and helper functions.
        These are not intended to be imported by external packages.

    scripts/:
        Scripts for setting up RabbitMQ and deploying the agent in Kubernetes.

    docs/:
        Contains project documentation, CLI usage guides, and contribution instructions.

    main.go:
        Initializes the application by setting up Kobra and registering commands.

Project Name: KubeGate

    A Kubernetes proxy tool to access private clusters.
    Stored in the GitHub repository: loaynaser3/KubeGate.

Key Features:

    Single-command execution.
    Interactive shell mode.
    Context management similar to kubectl config (e.g., storing RabbitMQ URLs, namespaces, and credentials).

Architecture:

    Client:
        Sends commands to RabbitMQ.
        Receives responses from RabbitMQ (via reply queues and correlation IDs).
    Agent:
        Processes commands received from RabbitMQ.
        Publishes results back to RabbitMQ.
