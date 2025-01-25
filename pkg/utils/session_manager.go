package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/loaynaser3/KubeGate/pkg/logging"
	"github.com/loaynaser3/KubeGate/pkg/queue"
)

// Session holds information about the reply queue
type Session struct {
	QueueName string    `json:"queue_name"`
	Timestamp time.Time `json:"timestamp"`
}

// SessionManager handles session-related logic for reply queues
type SessionManager struct {
	FilePath        string
	SessionDuration time.Duration
	MessageQueue    queue.MessageQueue
}

// NewSessionManager creates a new instance of SessionManager
func NewSessionManager(filePath string, duration time.Duration, mq queue.MessageQueue) *SessionManager {
	// Expand the file path for home directory
	filePath = strings.Replace(filePath, "~", os.Getenv("HOME"), 1)
	logging.Logger.WithField("file_path", filePath).Info("SessionManager initialized")
	return &SessionManager{
		FilePath:        filePath,
		SessionDuration: duration,
		MessageQueue:    mq,
	}
}

// GetOrCreateReplyQueue manages the lifecycle of the reply queue
func (sm *SessionManager) GetOrCreateReplyQueue() (string, error) {
	logging.Logger.WithField("file_path", sm.FilePath).Info("Attempting to get or create reply queue")

	// Check if the session file exists
	if _, err := os.Stat(sm.FilePath); err == nil {
		logging.Logger.Info("Session file found. Attempting to load it.")

		// Load the session file
		data, err := os.ReadFile(sm.FilePath)
		if err == nil {
			var session Session
			if err := json.Unmarshal(data, &session); err == nil {
				logging.Logger.WithFields(map[string]interface{}{
					"queue_name": session.QueueName,
					"timestamp":  session.Timestamp,
				}).Info("Session file loaded successfully")

				// Check if the session is still valid
				if time.Since(session.Timestamp) < sm.SessionDuration {
					logging.Logger.WithField("queue_name", session.QueueName).Info("Session is valid. Reusing the queue.")
					return session.QueueName, nil // Reuse the existing queue
				}

				logging.Logger.WithField("queue_name", session.QueueName).Info("Session expired. Deleting the old queue.")
				_ = sm.MessageQueue.DeleteQueue(session.QueueName) // Clean up expired queue
			} else {
				logging.Logger.WithError(err).Error("Failed to parse session file")
			}
		} else {
			logging.Logger.WithError(err).Error("Failed to read session file")
		}
	} else {
		logging.Logger.WithError(err).Info("Session file not found. Proceeding to create a new queue.")
	}

	// Create a new reply queue
	newQueue := "reply-queue-" + uuid.New().String()
	logging.Logger.WithField("queue_name", newQueue).Info("Creating new reply queue")
	if err := sm.MessageQueue.CreateQueue(newQueue); err != nil {
		logging.Logger.WithError(err).Error("Failed to create reply queue")
		return "", fmt.Errorf("failed to create reply queue: %w", err)
	}

	// Save the new session to the session file
	session := Session{
		QueueName: newQueue,
		Timestamp: time.Now().UTC(),
	}
	fileData, err := json.Marshal(session)
	if err != nil {
		logging.Logger.WithError(err).Error("Failed to marshal session data")
		return "", err
	}

	if err := os.WriteFile(sm.FilePath, fileData, 0644); err != nil {
		logging.Logger.WithError(err).Error("Failed to write session file")
		return "", err
	}

	logging.Logger.WithField("queue_name", newQueue).Info("Reply queue created and session saved successfully")
	return newQueue, nil
}

// CleanupStaleSession removes stale session files and queues
func (sm *SessionManager) CleanupStaleSession() error {
	logging.Logger.WithField("file_path", sm.FilePath).Info("Attempting to clean up stale session")

	// Check if the session file exists
	if _, err := os.Stat(sm.FilePath); err == nil {
		logging.Logger.Info("Session file found. Attempting to clean up.")

		// Load the session file
		data, err := os.ReadFile(sm.FilePath)
		if err == nil {
			var session Session
			if err := json.Unmarshal(data, &session); err == nil {
				logging.Logger.WithField("queue_name", session.QueueName).Info("Deleting queue from stale session")
				_ = sm.MessageQueue.DeleteQueue(session.QueueName)
			} else {
				logging.Logger.WithError(err).Error("Failed to parse session file during cleanup")
			}
		} else {
			logging.Logger.WithError(err).Error("Failed to read session file during cleanup")
		}

		// Remove the session file
		if err := os.Remove(sm.FilePath); err != nil {
			logging.Logger.WithError(err).Error("Failed to remove session file during cleanup")
			return err
		}
		logging.Logger.Info("Stale session cleanup completed successfully")
	} else {
		logging.Logger.WithError(err).Info("No session file found. No cleanup needed")
	}

	return nil
}
