package queue

type Message struct {
	Body      string
	TimeStamp int64
}

// Queue interface defines methods for interacting with a message queue.
type Queue interface {
	// SendMessage sends a message to the queue.
	SendMessage(message string) error

	// ReceiveMessage receives a channel of messages from the queue.
	ReceiveMessage() (<-chan Message, error)
	Close() error
}
