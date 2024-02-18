package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"command-queue/client"
	"command-queue/internal/util/queue"
)

const (
	defaultBufferLength = 1000
)

func main() {
	// Parse command-line arguments
	queueType := flag.String("queue", "", "Type of queue (rabbitmq or aws)")
	region := flag.String("region", "", "AWS region (required for aws)")
	queueURL := flag.String("queueURL", "", "aws queue URL")
	connectionString := flag.String("conn", "", "RabbitMQ connection string")
	queuName := flag.String("queueName", "", "Queue name")
	filePath := flag.String("file", "", "Input file path")
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

		q, err = queue.NewSQSQueue(*region, *queueURL)
		if err != nil {
			fmt.Printf("Error creating SQS queue: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Println("Invalid queue type. Supported types: rabbitmq, aws")
		os.Exit(1)
	}
	defer q.Close()

	// Initialize client with input source
	var inputSource *os.File
	if *filePath != "" {
		var err error
		inputSource, err = os.Open(*filePath)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			os.Exit(1)
		}
		defer inputSource.Close()
	} else {
		inputSource = os.Stdin
	}

	c := client.NewClient(inputSource, q)

	// Run the client
	if err := c.Start(ctx); err != nil {
		fmt.Printf("Error running client: %v\n", err)
		os.Exit(1)
	}
}
