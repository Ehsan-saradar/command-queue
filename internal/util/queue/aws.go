package queue

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// SQSQueue implements the Queue interface for AWS SQS.
type SQSQueue struct {
	svc              *sqs.SQS
	queueURL         string
	stopReceiving    chan struct{} // Channel to signal stop receiving messages
	receivingStopped chan struct{} // Channel to signal that receiving has stopped
}

// NewSQSQueue creates a new instance of SQSQueue.
func NewSQSQueue(region, queueURL string) (*SQSQueue, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	},
	)
	if err != nil {
		return nil, err
	}

	svc := sqs.New(sess)

	return &SQSQueue{
		svc:      svc,
		queueURL: queueURL,
	}, nil
}

// SendMessage sends a message to the AWS SQS queue.
func (q *SQSQueue) SendMessage(message string) error {
	_, err := q.svc.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(message),
		QueueUrl:    aws.String(q.queueURL),
	})
	return err
}

func (c *SQSQueue) ReceiveMessage() (<-chan string, error) {
	messageChannel := make(chan string)
	go func() {
		for {
			select {
			case <-c.stopReceiving:
				close(messageChannel)
				c.receivingStopped <- struct{}{}
				return
			default:
				receiveParams := &sqs.ReceiveMessageInput{
					QueueUrl:            aws.String(c.queueURL),
					MaxNumberOfMessages: aws.Int64(1),
					WaitTimeSeconds:     aws.Int64(20), // Adjust as needed
				}

				receiveResp, err := c.svc.ReceiveMessage(receiveParams)
				if err != nil {
					// Handle error
					continue
				}

				for _, msg := range receiveResp.Messages {
					messageChannel <- *msg.Body
				}
			}
		}
	}()

	return messageChannel, nil
}

func (c *SQSQueue) Close() error {
	close(c.stopReceiving)
	// make sure receiving has stopped before closing the connection
	<-c.receivingStopped
	return nil
}
