package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

// export LOG_LEVEL=debug
// export LOG_FILE=/var/log/kubegate.log
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

	// Set log output file
	logFile := os.Getenv("LOG_FILE")
	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			Logger.Fatalf("Failed to open log file: %v", err)
		}
		Logger.SetOutput(file)
	} else if logFile == "" {
		logFile = "/tmp/kubegate.log"
		file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			Logger.Fatalf("Failed to open log file: %v", err)
		}
		Logger.SetOutput(file)
		// Logger.SetOutput(os.Stdout) // Default to stdout if no log file is provided
	}
}
