package logger

import (
	"log"
)

// ConsoleLogger implements the Logger interface for logging to the console.
type ConsoleLogger struct{}

// NewConsoleLogger creates a new instance of ConsoleLogger.
func NewConsoleLogger() Logger {
	return &ConsoleLogger{}
}

// Logf logs a formatted message to the console.
func (c *ConsoleLogger) Printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}
