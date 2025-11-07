// Package fsd
package fsd

type CommandProcessor func(session SessionInterface, data []string, rawLine []byte) *Result

type CommandHandlerInterface interface {
	GetPossibleCommands() [][]byte
	GeneratePossibleCommands()
	Register(command ClientCommand, processor CommandProcessor, requirement *CommandRequirement)
	Call(command ClientCommand, session SessionInterface, data []string, rawLine []byte) *Result
}
