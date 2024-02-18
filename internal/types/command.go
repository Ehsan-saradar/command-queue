package types

import (
	"errors"
	"fmt"
	"strings"
)

type CommandType string

const (
	Undefined   CommandType = ""
	AddItem     CommandType = "addItem"
	DeleteItem  CommandType = "deleteItem"
	GetItem     CommandType = "getItem"
	GetAllItems CommandType = "getAllItems"
)

type Command struct {
	args []string
	Type CommandType
}

func ParseCommand(message string) (Command, error) {
	var command Command
	command.Type = command.GetCommand(message)
	command.args = command.GetArgs(message)
	if !command.isValid() {
		return Command{}, errors.New(fmt.Sprintf("Invalid message: %s\n", message))
	}
	return command, nil
}

func NewAddCommand(key, value string) Command {
	return Command{
		Type: AddItem,
		args: []string{key, value},
	}
}
func NewDeleteCommand(key string) Command {
	return Command{
		Type: DeleteItem,
		args: []string{key},
	}
}
func NewGetCommand(key string) Command {
	return Command{
		Type: GetItem,
		args: []string{key},
	}
}
func NewGetItemCommand(key string) Command {
	return Command{
		Type: GetItem,
		args: []string{key},
	}
}
func NewGetAllCommand() Command {
	return Command{
		Type: GetAllItems,
		args: []string{},
	}
}

func (c Command) GetCommand(message string) CommandType {
	message = strings.TrimSpace(strings.Split(message, "(")[0])
	switch message {
	case "addItem":
		return AddItem
	case "deleteItem":
		return DeleteItem
	case "getItem":
		return GetItem
	case "getAllItems":
		return GetAllItems
	default:
		return Undefined
	}
}

func (c Command) GetArgs(message string) []string {
	parts := strings.Split(message, "(")
	if len(parts) < 2 {
		return nil
	}
	parts = strings.Split(parts[1], ")")
	parts = strings.Split(parts[0], ",")
	for i := 0; i < len(parts); i++ {
		parts[i] = strings.Trim(strings.TrimSpace(parts[i]), "'")
		if len(parts[i]) == 0 {
			parts = append(parts[:i], parts[i+1:]...)
			i--
		}
	}
	return parts
}

func (c Command) isValid() bool {
	switch c.Type {
	case AddItem:
		if len(c.args) == 2 {
			return true
		}
	case DeleteItem, GetItem:
		if len(c.args) == 1 {
			return true
		}
	case GetAllItems:
		if len(c.args) == 0 {
			return true
		}
	}
	return false
}

func (c Command) Key() string {
	return c.args[0]
}
func (c Command) Value() string {
	if len(c.args) < 2 {
		return ""
	}
	return c.args[1]
}

func (c Command) String() string {
	switch len(c.args) {
	case 0:
		return fmt.Sprintf("%s()", c.Type)
	default:
		return fmt.Sprintf("%s('%s')", c.Type, strings.Join(c.args, "', '"))
	}
}
