package queue

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// SQSQueue implements the Queue interface for AWS SQS.
type SQSQueue struct {
	svc      *sqs.SQS
	queueURL string
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

// ReceiveMessage receives a message from the AWS SQS queue.
func (q *SQSQueue) ReceiveMessage() (string, error) {
	result, err := q.svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(q.queueURL),
		MaxNumberOfMessages: aws.Int64(1),
		WaitTimeSeconds:     aws.Int64(1),
	})
	if err != nil {
		return "", err
	}

	if len(result.Messages) == 0 {
		return "", nil
	}

	return *result.Messages[0].Body, nil
}

// DeleteMessage deletes a message from the AWS SQS queue.
func (q *SQSQueue) DeleteMessage(receiptHandle string) error {
	_, err := q.svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(q.queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	})
	return err
}
