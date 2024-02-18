package types

import (
	"reflect"
	"testing"
)

func TestNewCommand(t *testing.T) {
	tests := []struct {
		name          string
		message       string
		expectedType  CommandType
		expectedArgs  []string
		expectedError bool
	}{
		{
			name:          "Valid addItem command",
			message:       "addItem('key', 'val')",
			expectedType:  AddItem,
			expectedArgs:  []string{"key", "val"},
			expectedError: false,
		},
		{
			name:          "Valid deleteItem command",
			message:       "deleteItem('key')",
			expectedType:  DeleteItem,
			expectedArgs:  []string{"key"},
			expectedError: false,
		},
		{
			name:          "Valid getItem command",
			message:       "getItem('key')",
			expectedType:  GetItem,
			expectedArgs:  []string{"key"},
			expectedError: false,
		},
		{
			name:          "Valid getAllItems command",
			message:       "getAllItems()",
			expectedType:  GetAllItems,
			expectedArgs:  []string{},
			expectedError: false,
		},
		{
			name:          "Invalid message",
			message:       "Invalid message",
			expectedType:  Undefined,
			expectedArgs:  nil,
			expectedError: true,
		},
		{
			name:          "Invalid addItem command",
			message:       "addItem('key')",
			expectedType:  Undefined,
			expectedArgs:  nil,
			expectedError: true,
		},
		{
			name:          "Invalid getItem command",
			message:       "getItem()",
			expectedType:  Undefined,
			expectedArgs:  nil,
			expectedError: true,
		},
		{
			name:          "Invalid getAllItems command",
			message:       "getAllItems('key')",
			expectedType:  Undefined,
			expectedArgs:  nil,
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			command, err := ParseCommand(test.message)

			if test.expectedError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !test.expectedError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if command.Type != test.expectedType {
				t.Errorf("Expected type %s but got %s", test.expectedType, command.Type)
			}

			if !reflect.DeepEqual(command.args, test.expectedArgs) {
				t.Errorf("Expected args %v but got %v", test.expectedArgs, command.args)
			}
		})
	}
}

func TestCommand_isValid(t *testing.T) {
	tests := []struct {
		name           string
		command        Command
		expectedResult bool
	}{
		{
			name:           "Valid addItem command",
			command:        Command{Type: "addItem", args: []string{"key", "val"}},
			expectedResult: true,
		},
		{
			name:           "Valid deleteItem command",
			command:        Command{Type: "deleteItem", args: []string{"key"}},
			expectedResult: true,
		},
		{
			name:           "Valid getItem command",
			command:        Command{Type: "getItem", args: []string{"key"}},
			expectedResult: true,
		},
		{
			name:           "Valid getAllItems command",
			command:        Command{Type: "getAllItems", args: []string{}},
			expectedResult: true,
		},
		{
			name:           "Invalid addItem command",
			command:        Command{Type: "addItem", args: []string{"key"}},
			expectedResult: false,
		},
		{
			name:           "Invalid getItem command",
			command:        Command{Type: "getItem", args: []string{}},
			expectedResult: false,
		},
		{
			name:           "Invalid command type",
			command:        Command{Type: "InvalidType", args: []string{}},
			expectedResult: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.command.isValid()
			if result != test.expectedResult {
				t.Errorf("Expected %v but got %v", test.expectedResult, result)
			}
		})
	}
}
