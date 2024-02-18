package main

import (
	"command-queue/internal/util/logger"
	"command-queue/server"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"command-queue/internal/util/queue"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Queue buffer length
const defaultBufferLength = 1000

func main() {
	// Parse command-line arguments
	queueType := flag.String("queue", "", "Type of queue (rabbitmq or aws)")
	region := flag.String("region", "", "AWS region (required for aws)")
	queueURL := flag.String("queueURL", "", "aws queue URL")
	connectionString := flag.String("conn", "", "RabbitMQ connection string")
	queuName := flag.String("queueName", "", "Queue name")
	flag.Parse()

	// Check if required arguments are provided
	if *queueType == "" {
		fmt.Println("Please provide a queue type (rabbitmq or aws)")
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling for cancellation
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		fmt.Println("\nReceived SIGTERM or SIGINT. Cancelling all operations...")
		cancel()
	}()

	// Initialize queue based on the provided type
	var q queue.Queue
	var err error
	switch *queueType {
	case "rabbitmq":
		if *connectionString == "" {
			fmt.Println("Please provide RabbitMQ connection string")
			os.Exit(1)
		}
		if *queuName == "" {
			fmt.Println("Please provide queue name")
			os.Exit(1)
		}
		q, err = queue.NewRabbitMQQueue(ctx, *connectionString, *queueType, defaultBufferLength)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "aws":
		if *region == "" || *queueURL == "" {
			fmt.Println("Please provide AWS region and SQS queue URL")
			os.Exit(1)
		}

		awsq, err := session.NewSession(&aws.Config{Region: aws.String(*region)})
		if err != nil {
			fmt.Printf("Error creating AWS session: %v\n", err)
			os.Exit(1)
		}

		_ = awsq
		//svc := sqs.New(sess)
		//q = queue.NewSQSQueue(svc, *queueURL)
	default:
		fmt.Println("Invalid queue type. Supported types: rabbitmq, aws")
		os.Exit(1)
	}
	defer q.Close()

	// Initialize server
	s := server.NewServer(ctx, q, logger.NewConsoleLogger())

	// Run the server
	if err := s.Start(); err != nil {
		fmt.Printf("Error running server: %v\n", err)
		os.Exit(1)
	}
}
