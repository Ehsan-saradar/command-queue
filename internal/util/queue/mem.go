package queue

// memQueue implements the Queue interface for in-memory queue and is used for testing.
type memQueue struct {
	messages chan Message
}

func NewMemQueue(bufferLength int) Queue {
	return &memQueue{
		messages: make(chan Message, bufferLength),
	}
}

func (q *memQueue) SendMessage(message string) error {
	q.messages <- Message{Body: message}
	return nil
}

func (q *memQueue) ReceiveMessage() (<-chan Message, error) {
	return q.messages, nil
}

func (q *memQueue) Close() error {
	close(q.messages)
	return nil
}
