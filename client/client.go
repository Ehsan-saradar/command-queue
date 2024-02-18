package client

import (
	"bufio"
	"context"
	"fmt"
	"io"

	"command-queue/internal/types"
	"command-queue/internal/util/queue"
)

type Client struct {
	inputSource io.Reader
	queue       queue.Queue
}

func NewClient(inputSource io.Reader, queue queue.Queue) *Client {
	return &Client{
		inputSource: inputSource,
		queue:       queue,
	}
}

func (c *Client) Start(ctx context.Context) error {
	// Read commands from the input source and send them to the server.
	scanner := bufio.NewScanner(c.inputSource)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			// Context cancelled. Stop processing.
			return nil
		default:
			str := scanner.Text()

			_, err := types.ParseCommand(str)
			if err != nil {
				return fmt.Errorf("error creating command %v", err)
			}
			// check if command is valid or not

			err = c.queue.SendMessage(str)
			if err != nil {
				return fmt.Errorf("error sending command %v", err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %v", err)
	}

	return nil
}

// Stop stops the client, preventing it from sending further commands to the server.
func (c *Client) Stop() error {
	// Perform any cleanup or shutdown logic here.
	return nil
}
