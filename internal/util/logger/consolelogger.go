package logger

import (
	"fmt"
	"sync"
)

// ConsoleLogger implements the Logger interface for logging to the console.
type ConsoleLogger struct {
	mu sync.Mutex // Mutex for synchronization
}

// NewConsoleLogger creates a new instance of ConsoleLogger.
func NewConsoleLogger() Logger {
	return &ConsoleLogger{}
}

// Logf logs a formatted message to the console.
func (c *ConsoleLogger) Logf(format string, args ...interface{}) {
	c.mu.Lock()         // Acquire mutex lock to ensure exclusive access to the console
	defer c.mu.Unlock() // Release mutex lock when Logf exits, even in case of panics or early returns
	fmt.Printf(format, args...)
}
