package config

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Context struct {
	Name         string `yaml:"name"`
	RabbitMQURL  string `yaml:"rabbitmq-url"`
	CommandQueue string `yaml:"command-queue"`
	ReplyQueue   string `yaml:"reply-queue"`
}

type Config struct {
	CurrentContext string    `yaml:"current-context"`
	Contexts       []Context `yaml:"contexts"`
}

var configFile = filepath.Join(os.Getenv("HOME"), ".kubegate", "config.yaml")

// LoadConfig loads the configuration from the YAML file
func LoadConfig() (*Config, error) {
	file, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			// Return a default config if the file doesn't exist
			return &Config{}, nil
		}
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

// SaveConfig saves the configuration to the YAML file
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

// GetContext retrieves a context by name
func GetContext(cfg *Config, name string) (*Context, error) {
	for _, ctx := range cfg.Contexts {
		if ctx.Name == name {
			return &ctx, nil
		}
	}
	return nil, errors.New("context not found")
}

// SetContext adds or updates a context in the configuration
func SetContext(cfg *Config, ctx Context) {
	for i, c := range cfg.Contexts {
		if c.Name == ctx.Name {
			cfg.Contexts[i] = ctx // Update existing context
			return
		}
	}
	cfg.Contexts = append(cfg.Contexts, ctx) // Add new context
}
