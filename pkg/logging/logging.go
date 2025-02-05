package logging

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sort"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger instance
var Logger = logrus.New()

// CustomFormatter ensures fields appear before caller info
type CustomFormatter struct {
	logrus.TextFormatter
}

// Format controls the order of fields in logs
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Extract caller info first
	fileInfo := ""
	if entry.HasCaller() {
		fileInfo = fmt.Sprintf(" %s:%d", entry.Caller.File, entry.Caller.Line)
	}

	// Format message and fields before file info
	fields := make([]string, 0, len(entry.Data))
	for key, value := range entry.Data {
		fields = append(fields, fmt.Sprintf("%s=%v", key, value))
	}

	// Sort fields to maintain consistent order
	sort.Strings(fields)
	logMsg := fmt.Sprintf("%s [%s] %s %s%s\n", entry.Time.Format("2006-01-02 15:04:05"), entry.Level, entry.Message, fmt.Sprint(fields), fileInfo)

	return []byte(logMsg), nil
}

// Initialize sets up the logger with flexible output options
func Initialize() {
	// Set log level
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	Logger.SetLevel(level)

	// Configure log format (JSON or Text)
	useJSON := os.Getenv("LOG_FORMAT") == "json"
	if useJSON {
		Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			CallerPrettyfier: func(frame *runtime.Frame) (string, string) {
				return "", fmt.Sprintf("%s:%d", frame.File, frame.Line)
			},
		})
	} else {
		Logger.SetFormatter(&CustomFormatter{
			TextFormatter: logrus.TextFormatter{
				FullTimestamp: true,
			},
		})
	}

	// Determine log output behavior
	logOutput := os.Getenv("LOG_OUTPUT") // "stdout", "none", or default (empty)
	logFile := os.Getenv("LOG_FILE")     // Log file path

	var outputWriters []io.Writer

	if logOutput == "stdout" {
		outputWriters = append(outputWriters, os.Stdout)
	}
	if logFile != "" {
		outputWriters = append(outputWriters, &lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    10,   // MB before rotating
			MaxBackups: 3,    // Number of old log files to keep
			MaxAge:     28,   // Days to retain old logs
			Compress:   true, // Compress old logs
		})
	}

	// Set output
	if len(outputWriters) > 0 {
		Logger.SetOutput(io.MultiWriter(outputWriters...))
		Logger.SetReportCaller(true)
		Logger.Info("Logger initialized successfully")
	} else {
		Logger.SetOutput(io.Discard) // Disable logging by default
	}
}
