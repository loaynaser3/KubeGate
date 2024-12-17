package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type AgentConfig struct {
	RabbitMQURL  string `yaml:"rabbitmq-url"`
	CommandQueue string `yaml:"command-queue"`
	Backend      string `yaml:"backend"`
}

var agentConfigFile = filepath.Join(os.Getenv("HOME"), ".kubegate", "agent-config.yaml")

// LoadAgentConfig loads the agent configuration from a YAML file
func LoadAgentConfig() (*AgentConfig, error) {
	// Check for environment variables first
	rabbitURL := os.Getenv("KUBEGATE_RABBITMQ_URL")
	commandQueue := os.Getenv("KUBEGATE_RABBITMQ_QUEUE")

	// If both environment variables are set, use them
	if rabbitURL != "" && commandQueue != "" {
		return &AgentConfig{
			RabbitMQURL:  rabbitURL,
			CommandQueue: commandQueue,
		}, nil
	}

	// Otherwise, fall back to the YAML file
	file, err := os.ReadFile(agentConfigFile)
	if err != nil {
		return nil, err
	}

	var cfg AgentConfig
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
