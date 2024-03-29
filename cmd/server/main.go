package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"command-queue/internal/util/logger"
	"command-queue/server"

	"command-queue/internal/util/queue"
)

// Queue buffer length
const (
	defaultBufferLength = 1000
	defaultMaxWorkers   = 10
)

func main() {
	// Parse command-line arguments
	queueType := flag.String("queue", "", "Type of queue (rabbitmq or aws)")
	region := flag.String("region", "", "AWS region (required for aws)")
	queueURL := flag.String("queueURL", "", "aws queue URL")
	connectionString := flag.String("conn", "", "RabbitMQ connection string")
	queuName := flag.String("queueName", "", "Queue name")
	maxWorkers := flag.Int("maxWorkers", defaultMaxWorkers, "Maximum number of workers")
	flag.Parse()

	// Check if required arguments are provided
	if *queueType == "" {
		fmt.Println("Please provide a queue type (rabbitmq or aws)")
		os.Exit(1)
	}
	if *maxWorkers < 1 {
		fmt.Println("maxWorkers must be a positive integer")
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
		q, err = queue.NewRabbitMQQueue(ctx, *connectionString, *queuName, defaultBufferLength)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "aws":
		if *region == "" || *queueURL == "" {
			fmt.Println("Please provide AWS region and SQS queue URL")
			os.Exit(1)
		}

		q, err = queue.NewSQSQueue(*region, *queueURL, defaultBufferLength)
		if err != nil {
			fmt.Printf("Error creating SQS queue: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Println("Invalid queue type. Supported types: rabbitmq, aws")
		os.Exit(1)
	}
	defer q.Close()

	// Initialize server
	s := server.NewServer(q, logger.NewConsoleLogger(), *maxWorkers)

	// Run the server
	if err := s.Start(ctx); err != nil {
		fmt.Printf("Error running server: %v\n", err)
		os.Exit(1)
	}
}
