// Package command
package command

import (
	"fmt"
	"github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
	"maps"
)

type CommandBlock struct {
	processor   fsd.CommandProcessor
	requirement *fsd.CommandRequirement
}

type CommandHandler struct {
	commands         map[fsd.ClientCommand]*CommandBlock
	possibleCommands [][]byte
}

func NewCommandHandler() *CommandHandler {
	return &CommandHandler{
		commands:         make(map[fsd.ClientCommand]*CommandBlock),
		possibleCommands: make([][]byte, 0),
	}
}

func (handler *CommandHandler) GetPossibleCommands() [][]byte {
	return handler.possibleCommands
}

func (handler *CommandHandler) GeneratePossibleCommands() {
	for command := range maps.Keys(handler.commands) {
		handler.possibleCommands = append(handler.possibleCommands, []byte(command))
	}
}

func (handler *CommandHandler) Register(command fsd.ClientCommand, processor fsd.CommandProcessor, requirement *fsd.CommandRequirement) {
	handler.commands[command] = &CommandBlock{processor: processor, requirement: requirement}
}

func (handler *CommandHandler) Call(command fsd.ClientCommand, session fsd.SessionInterface, data []string, rawLine []byte) *fsd.Result {
	block, ok := handler.commands[command]
	if !ok || block.processor == nil {
		return nil
	}
	length := len(data)
	if block.requirement != nil && length < block.requirement.RequireLength {
		return fsd.ResultError(fsd.Syntax, block.requirement.Fatal, session.Callsign(), fmt.Errorf("datapack length too short, require %d but got %d", block.requirement.RequireLength, length))
	}
	return block.processor(session, data, rawLine)
}
