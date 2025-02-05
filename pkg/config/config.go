package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Context struct {
	Name         string `yaml:"name"`
	RabbitMQURL  string `yaml:"rabbitmq-url"`
	CommandQueue string `yaml:"command-queue"`
	ReplyQueue   string `yaml:"reply-queue"`
	Backend      string `yaml:"backend"`
}

type Config struct {
	CurrentContext string    `yaml:"current-context"`
	Contexts       []Context `yaml:"contexts"`
}

var configFile = filepath.Join(os.Getenv("HOME"), ".kubegate", "config.yaml")

// LoadConfig loads the configuration from the YAML file and then overrides
// any values with environment variables if they are set.
func LoadConfig() (*Config, error) {
	var cfg Config

	file, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			// If the file doesn't exist, start with an empty configuration.
			cfg = Config{}
		} else {
			return nil, err
		}
	} else {
		if err = yaml.Unmarshal(file, &cfg); err != nil {
			return nil, err
		}
	}

	// override with environment variables if available.
	overrideConfigWithEnv(&cfg)
	return &cfg, nil
}

// overrideConfigWithEnv overrides configuration values with environment variables.
func overrideConfigWithEnv(cfg *Config) {
	// Override the global current context if the environment variable is set.
	if envCurrent := os.Getenv("CURRENT_CONTEXT"); envCurrent != "" {
		cfg.CurrentContext = envCurrent
	}

	for i := range cfg.Contexts {
		if cfg.Contexts[i].Name == cfg.CurrentContext {
			if envRabbitMQ := os.Getenv("RABBITMQ_URL"); envRabbitMQ != "" {
				cfg.Contexts[i].RabbitMQURL = envRabbitMQ
			}
			if envCommandQueue := os.Getenv("COMMAND_QUEUE"); envCommandQueue != "" {
				cfg.Contexts[i].CommandQueue = envCommandQueue
			}
			if envReplyQueue := os.Getenv("REPLY_QUEUE"); envReplyQueue != "" {
				cfg.Contexts[i].ReplyQueue = envReplyQueue
			}
			if envBackend := os.Getenv("BACKEND"); envBackend != "" {
				cfg.Contexts[i].Backend = envBackend
			}
			break
		}
	}
}

// SaveConfig saves the configuration to the YAML file.
func SaveConfig(cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(configFile), 0755)
	if err != nil {
		return err
	}

	return os.WriteFile(configFile, data, 0644)
}

// GetContext retrieves a context by name.
func GetContext(cfg *Config, name string) (*Context, error) {
	for i := range cfg.Contexts {
		if cfg.Contexts[i].Name == name {
			return &cfg.Contexts[i], nil
		}
	}
	return nil, errors.New("context not found")
}

// SetContext adds or updates a context in the configuration.
func SetContext(cfg *Config, ctx Context) {
	for i, c := range cfg.Contexts {
		if c.Name == ctx.Name {
			cfg.Contexts[i] = ctx // Update existing context.
			return
		}
	}
	cfg.Contexts = append(cfg.Contexts, ctx) // Add new context.
}

// UseContext switches the current context in the configuration
func UseContext(cfg *Config, contextName string) error {
	for _, ctx := range cfg.Contexts {
		if ctx.Name == contextName {
			cfg.CurrentContext = contextName
			return SaveConfig(cfg)
		}
	}
	return fmt.Errorf("context not found: %s", contextName)
}

// DeleteContext removes a context by name
func DeleteContext(cfg *Config, contextName string) error {
	index := -1
	for i, ctx := range cfg.Contexts {
		if ctx.Name == contextName {
			index = i
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("context not found: %s", contextName)
	}

	// Remove the context from the list
	cfg.Contexts = append(cfg.Contexts[:index], cfg.Contexts[index+1:]...)

	// Reset current context if deleted
	if cfg.CurrentContext == contextName {
		cfg.CurrentContext = ""
	}

	return SaveConfig(cfg)
}
