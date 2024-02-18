package logger

import (
	"fmt"
)

// ConsoleLogger implements the Logger interface for logging to the console.
type ConsoleLogger struct{}

// NewConsoleLogger creates a new instance of ConsoleLogger.
func NewConsoleLogger() Logger {
	return &ConsoleLogger{}
}

// Logf logs a formatted message to the console.
func (c *ConsoleLogger) Logf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
