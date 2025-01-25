package logging

import (
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
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

	// Set log output format (TextFormatter for now)
	useJSON := os.Getenv("LOG_FORMAT") == "json"
	if useJSON {
		Logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	// Configure log output
	logFile := os.Getenv("LOG_FILE")
	if logFile != "" {
		// Use log rotation with lumberjack
		Logger.SetOutput(&lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    10,   // Max megabytes before log is rotated
			MaxBackups: 3,    // Max number of old log files to keep
			MaxAge:     28,   // Max number of days to retain old logs
			Compress:   true, // Compress old logs
		})
	} else {
		// Default to stdout if no log file is provided
		Logger.SetOutput(os.Stdout)
	}

	Logger.Info("Logger initialized successfully")
}
