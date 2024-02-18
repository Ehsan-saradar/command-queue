package client

import (
	"bufio"
	"command-queue/internal/types"
	"command-queue/internal/util/queue"
	"context"
	"fmt"
	"io"
)

type Client struct {
	ctx         context.Context
	inputSource io.Reader
	queue       queue.Queue
}

func NewClient(ctx context.Context, inputSource io.Reader, queue queue.Queue) *Client {
	return &Client{
		ctx:         ctx,
		inputSource: inputSource,
		queue:       queue,
	}
}

func (c *Client) Start() error {
	// Use the context's Done channel to check for cancellation signals.
	done := c.ctx.Done()

	// Read commands from the input source and send them to the server.
	scanner := bufio.NewScanner(c.inputSource)
	for scanner.Scan() {
		select {
		case <-done:
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
