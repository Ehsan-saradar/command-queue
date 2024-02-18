package server

import (
	"bytes"
	"context"
	"os"
	"testing"

	"command-queue/internal/types"

	"github.com/stretchr/testify/assert"

	"command-queue/internal/util/logger"
	"command-queue/internal/util/queue"
)

func TestServer_Start(t *testing.T) {
	// Create a context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize memQueue
	memQ := queue.NewMemQueue(5)

	// Create a server with memQueue and a logger
	s := NewServer(ctx, memQ, logger.NewConsoleLogger())

	// Define input commands
	inputCommands := []string{"addItem('key1,'value1')", "deleteItem('key2')", "getAllItems()"}

	// Create a buffer with input commands
	inputBuffer := bytes.NewBufferString("")
	for _, cmd := range inputCommands {
		inputBuffer.WriteString(cmd + "\n")
	}

	// Start the server in a separate goroutine
	go func() {
		err := s.Start()
		assert.Nilf(t, err, "Start returned an error: %v", err)
	}()

	// Send input commands to the memQueue
	for _, cmd := range inputCommands {
		err := memQ.SendMessage(cmd)
		assert.Nilf(t, err, "SendMessage returned an error: %v", err)
	}
	err := s.Stop()
	assert.Nilf(t, err, "Stop returned an error: %v", err)

	// Simulate cancellation of the context to stop the server
	cancel()
}

func TestProcessCommand(t *testing.T) {
	server := NewServer(context.TODO(), nil, logger.NewConsoleLogger())

	// delete test files if they exist
	os.Remove("key2")
	os.Remove("allItems")

	tests := []struct {
		name    string
		message string
	}{
		{name: "AddItem1", message: types.NewAddCommand("key1", "value1").String()},
		{name: "AddItem2", message: types.NewAddCommand("key2", "value2").String()},
		{name: "DeleteItem1", message: types.NewDeleteCommand("key1").String()},
		{name: "GetItem1", message: types.NewGetCommand("key2").String()},
		{name: "GetAllItems1", message: types.NewGetAllCommand().String()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server.processCommand(tt.message)
		})
	}
	keys, values := server.orderedMap.GetAll()
	assert.Equal(t, []string{"key2"}, keys)
	assert.Equal(t, []interface{}{"value2"}, values)

	bt, err := os.ReadFile("key2")
	assert.Nilf(t, err, "Error reading file: %v", err)
	assert.Equal(t, "key2 : value2\n", string(bt))
	os.Remove("key2")

	bt, err = os.ReadFile("allItems")
	assert.Nilf(t, err, "Error reading file: %v", err)
	assert.Equal(t, "key2 : value2\n", string(bt))
	os.Remove("allItems")
}
