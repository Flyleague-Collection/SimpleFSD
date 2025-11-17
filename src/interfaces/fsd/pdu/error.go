// Package pdu
package pdu

import (
	"github.com/half-nothing/simple-fsd/src/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/src/interfaces/global"
)

type Error struct {
	*Base
	ErrorType fsd.ClientError
	Param     string
	Message   string
}

func NewError(to string, errorType fsd.ClientError, param string, message string) *Error {
	return &Error{
		Base:      NewBase(fsd.ClientCommandError, global.FSDServerName, to),
		ErrorType: errorType,
		Param:     param,
		Message:   message,
	}
}

func (c *Error) Build() []byte {
	return MakeProtocolDataUnitPacket(
		c.GetType(),
		c.From,
		c.To,
		c.ErrorType.Data,
		c.Param,
		c.Message,
	)
}
