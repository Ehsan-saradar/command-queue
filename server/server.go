package server

import (
	"context"
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	"command-queue/internal/types"
	"command-queue/internal/util/logger"
	"command-queue/internal/util/orderedmap"
	"command-queue/internal/util/queue"
)

// Server implements the Server interface.
type Server struct {
	queue      queue.Queue
	orderedMap orderedmap.OrderedMap
	fileMutex  sync.Mutex
	log        logger.Logger
	semaphore  chan interface{}
	cnt        atomic.Uint64
}

// NewServer creates a new instance of Server.
func NewServer(q queue.Queue, log logger.Logger, maxWorkers int) *Server {
	return &Server{
		queue:      q,
		orderedMap: orderedmap.NewOrderedMap(),
		fileMutex:  sync.Mutex{},
		log:        log,
		semaphore:  make(chan interface{}, maxWorkers),
		cnt:        atomic.Uint64{},
	}
}

// Start starts the server, allowing it to read messages from the queue and process commands.
func (s *Server) Start(ctx context.Context) error {
	// Start reading messages from the queue in a separate goroutine.
	messages, err := s.queue.ReceiveMessage()
	if err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		case message, ok := <-messages:
			if !ok {
				return nil
			}
			// Acquire a semaphore slot
			s.semaphore <- struct{}{}

			// Launch a new goroutine to process the message concurrently
			go func(msg queue.Message) {
				defer func() {
					// Release the semaphore slot when the goroutine completes
					<-s.semaphore
				}()
				s.processCommand(msg)
			}(message)
		}
	}
}

// Stop stops the server, preventing it from reading messages and processing commands.
func (s *Server) Stop() error {
	// Server has stopped successfully.
	return nil
}

func (s *Server) processCommand(message queue.Message) {
	command, err := types.ParseCommand(message.Body)
	if err != nil {
		s.log.Logf("Error parsing command: %v\n", err)
		return
	}
	switch command.Type {
	case types.AddItem:
		s.orderedMap.Set(command.Key(), command.Value(), message.TimeStamp)
	case types.DeleteItem:
		s.orderedMap.DeleteItem(command.Key())
	case types.GetItem:
		val, ok := s.orderedMap.Get(command.Key())
		if ok {
			s.writeToFile(command.Key(), fmt.Sprintf("%s : %s\n", command.Key(), val))
		}
	case types.GetAllItems:
		keys, values := s.orderedMap.GetAll()
		result := ""
		for i, key := range keys {
			result += fmt.Sprintf("%s : %s\n", key, values[i])
		}
		s.writeToFile("allItems", result)
	}
}

func (s *Server) writeToFile(filename, content string) {
	// read the counter value and increment it
	index := s.cnt.Add(1)
	filename = fmt.Sprintf("%s_%d", filename, index)
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		s.log.Logf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		s.log.Logf("Error writing to file: %v\n", err)
	}
}
