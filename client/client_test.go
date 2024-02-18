package client

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"command-queue/internal/util/queue"
)

func TestClient_Start(t *testing.T) {
	// Initialize memQueue
	memQ := queue.NewMemQueue(10)

	// Define input commands
	inputCommands := []string{"addItem('key1,'value1')", "deleteItem('key2')", "getAllItems()"}

	// Create a buffer with input commands
	inputBuffer := bytes.NewBufferString("")
	for _, cmd := range inputCommands {
		inputBuffer.WriteString(cmd + "\n")
	}

	// Create a context
	ctx := context.Background()

	// Create a client with memQueue and input buffer
	c := NewClient(ctx, inputBuffer, memQ)

	// Start the client
	err := c.Start()
	assert.Nilf(t, err, "Start returned an error: %v", err)
	err = c.Stop()
	assert.Nilf(t, err, "Stop returned an error: %v", err)

	// Check if messages were sent to the queue
	receivedMessages, _ := memQ.ReceiveMessage()
	for _, expectedMsg := range inputCommands {
		receivedMsg := <-receivedMessages
		if receivedMsg != expectedMsg {
			t.Errorf("Expected message %s, but got %s", expectedMsg, receivedMsg)
		}
	}

	// Check if there are no more messages in the queue
	select {
	case receivedMsg := <-receivedMessages:
		t.Errorf("Unexpected message in the queue: %s", receivedMsg)
	default:
		// No message in the channel, as expected
	}
}
