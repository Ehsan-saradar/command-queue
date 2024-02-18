package server

import (
	"context"
	"fmt"
	"os"
	"sync"

	"command-queue/internal/types"
	"command-queue/internal/util/logger"
	"command-queue/internal/util/orderedmap"
	"command-queue/internal/util/queue"
)

// Server implements the Server interface.
type Server struct {
	ctx        context.Context
	queue      queue.Queue
	orderedMap orderedmap.OrderedMap
	fileMutex  sync.Mutex
	log        logger.Logger
}

// NewServer creates a new instance of Server.
func NewServer(ctx context.Context, q queue.Queue, log logger.Logger) *Server {
	return &Server{
		ctx:        ctx,
		queue:      q,
		orderedMap: orderedmap.NewOrderedMap(),
		fileMutex:  sync.Mutex{},
		log:        log,
	}
}

// Start starts the server, allowing it to read messages from the queue and process commands.
func (s *Server) Start() error {
	// Start reading messages from the queue in a separate goroutine.
	messages, err := s.queue.ReceiveMessage()
	if err != nil {
		return err
	}
	for {
		select {
		case <-s.ctx.Done():
			return nil
		case message, ok := <-messages:
			if !ok {
				return nil
			}
			s.processCommand(message)
		}
	}
}

// Stop stops the server, preventing it from reading messages and processing commands.
func (s *Server) Stop() error {
	s.ctx.Done()
	// Server has stopped successfully.
	return nil
}

func (s *Server) processCommand(message string) {
	command, err := types.ParseCommand(message)
	if err != nil {
		s.log.Logf("Error parsing command: %v\n", err)
		return
	}
	switch command.Type {
	case types.AddItem:
		s.orderedMap.Set(command.Key(), command.Value())
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
	s.fileMutex.Lock()
	defer s.fileMutex.Unlock()

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
