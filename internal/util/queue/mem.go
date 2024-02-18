package queue

// memQueue implements the Queue interface for in-memory queue and is used for testing.
type memQueue struct {
	messages chan string
}

func NewMemQueue(bufferLength int) Queue {
	return &memQueue{
		messages: make(chan string, bufferLength),
	}
}

func (q *memQueue) SendMessage(message string) error {
	q.messages <- message
	return nil
}

func (q *memQueue) ReceiveMessage() (<-chan string, error) {
	return q.messages, nil
}

func (q *memQueue) Close() error {
	close(q.messages)
	return nil
}
