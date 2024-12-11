```
KubeGate/
├── cmd/
│   ├── root.go             # Root command, initializes Kobra and subcommands
│   ├── single/
│   │   └── single.go       # Single command execution logic
│   ├── shell/
│   │   └── shell.go        # Interactive shell logic
│   ├── agent/
│       └── agent.go        # Command to start the agent
├── pkg/
│   ├── k8s/
│   │   └── k8s.go          # Kubernetes interaction logic
│   ├── queue/
│   │   ├── producer.go     # Logic for publishing messages to RabbitMQ
│   │   └── consumer.go     # Logic for consuming messages from RabbitMQ
│   ├── shell/
│   │   └── interactive.go  # Logic for handling interactive commands
│   ├── config/
│       └── config.go       # Configuration loader (e.g., for RabbitMQ, Kubernetes contexts)
├── internal/
│   ├── logger/
│   │   └── logger.go       # Logging setup
│   ├── utils/
│       └── utils.go        # Helper functions used across the project
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