package queue

import (
	"testing"
)

func TestMemQueue(t *testing.T) {
	// Initialize memQueue with buffer length of 5
	q := NewMemQueue(5)

	// Test SendMessage method
	message := "test message"
	if err := q.SendMessage(message); err != nil {
		t.Errorf("SendMessage failed: %v", err)
	}

	// Test ReceiveMessage method
	receivedMessages, err := q.ReceiveMessage()
	if err != nil {
		t.Errorf("ReceiveMessage failed: %v", err)
	}

	// Verify received message
	select {
	case receivedMessage := <-receivedMessages:
		if receivedMessage != message {
			t.Errorf("Received message doesn't match sent message: got %s, want %s", receivedMessage, message)
		}
	default:
		t.Error("No message received")
	}

	// Test Close method
	if err := q.Close(); err != nil {
		t.Errorf("Close failed: %v", err)
	}

	// Ensure the queue is closed
	select {
	case _, ok := <-q.(*memQueue).messages:
		if ok {
			t.Error("Queue is not closed")
		}
	default:
	}
}
