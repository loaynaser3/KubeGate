package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

// Initialize sets up the logger with default settings
func Initialize() {
	// Set log level from environment variable or default to INFO
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	Logger.SetLevel(level)

	// Set log output format
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}
