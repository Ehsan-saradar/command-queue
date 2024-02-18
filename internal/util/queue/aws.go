package queue

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// SQSQueue implements the Queue interface for AWS SQS.
type SQSQueue struct {
	ctx          aws.Context
	svc          *sqs.SQS
	queueURL     string
	bufferLength int64
}

// NewSQSQueue creates a new instance of SQSQueue.
func NewSQSQueue(region, queueURL string, bufferLength int64) (*SQSQueue, error) {
	config := aws.NewConfig()
	config.Credentials = credentials.NewEnvCredentials()
	config.Region = aws.String(region)

	sess, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	return &SQSQueue{
		ctx:          aws.BackgroundContext(),
		svc:          sqs.New(sess),
		queueURL:     queueURL,
		bufferLength: bufferLength,
	}, nil
}

// SendMessage sends a message to the AWS SQS queue.
func (s *SQSQueue) SendMessage(message string) error {
	_, err := s.svc.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(message),
		QueueUrl:    aws.String(s.queueURL),
	})
	return err
}

func (s *SQSQueue) ReceiveMessage() (<-chan string, error) {
	messageChannel := make(chan string, s.bufferLength)
	go func() {
		for {
			receiveParams := &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(s.queueURL),
				MaxNumberOfMessages: aws.Int64(10),
			}
			receiveResp, err := s.svc.ReceiveMessageWithContext(s.ctx, receiveParams)
			if err != nil {
				close(messageChannel)
				return
			}

			for _, msg := range receiveResp.Messages {
				messageChannel <- *msg.Body
			}
		}
	}()
	return messageChannel, nil
}

func (s *SQSQueue) Close() error {
	// Cancel the context to signal stop receiving messages
	s.ctx.Done()
	return nil
}
