package queue

import (
	"context"
	"fmt"

	"github.com/streadway/amqp"
)

// RabbitMQQueue implements the Queue interface for RabbitMQ.
type RabbitMQQueue struct {
	ctx          context.Context
	connection   *amqp.Connection
	channel      *amqp.Channel
	queueName    string
	routingName  string
	bufferLength int
}

// NewRabbitMQQueue creates a new instance of RabbitMQQueue.
func NewRabbitMQQueue(ctx context.Context, url string, queueName string, bufferLength int) (Queue, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, err
	}
	return &RabbitMQQueue{
		ctx:          ctx,
		connection:   conn,
		channel:      ch,
		routingName:  q.Name,
		queueName:    queueName,
		bufferLength: bufferLength,
	}, nil
}

// ReceiveMessage receives a channel of messages from the RabbitMQ queue.
func (q *RabbitMQQueue) ReceiveMessage() (<-chan string, error) {
	msgs, err := q.channel.Consume(
		q.queueName, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		return nil, err
	}

	msgChan := make(chan string, q.bufferLength)
	go func() {
		defer close(msgChan)

		for {
			select {
			case <-q.ctx.Done():
				return // Exit goroutine if context is canceled
			case msg, ok := <-msgs:
				if !ok {
					return // Exit goroutine if message channel is closed
				}
				msgChan <- string(msg.Body)
			}
		}
	}()

	return msgChan, nil
}

// DeleteMessage is not implemented for RabbitMQQueue.
func (q *RabbitMQQueue) DeleteMessage(message string) error {
	// Not implemented for RabbitMQQueue.
	return nil
}

// SendMessage sends a message to the RabbitMQ queue.
func (q *RabbitMQQueue) SendMessage(message string) error {
	return q.channel.Publish(
		"",            // exchange
		q.routingName, // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

func (q *RabbitMQQueue) Close() error {
	fmt.Println("Try to close channel")
	err := q.channel.Close()
	if err != nil {
		return err
	}
	return q.connection.Close()
}
