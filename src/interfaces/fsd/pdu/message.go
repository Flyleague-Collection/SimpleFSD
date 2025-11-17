// Package pdu
package pdu

import (
	"github.com/half-nothing/simple-fsd/src/interfaces/fsd"
)

type Message struct {
	*Base
	Message string
}

func (c *Message) Parse(data []string, raw []byte) (Interface, *fsd.CommandResult) {
	// ZSHA_CTR ZSSS_APP 111
	// [   0  ] [   1  ] [2]
	if r := c.CheckLength(fsd.ClientCommandMessage, len(data)); r != nil {
		return nil, r
	}
	command := &Message{
		Base:    NewBase(fsd.ClientCommandMessage, data[0], data[1]),
		Message: data[2],
	}
	command.raw = raw
	return command, fsd.CommandResultSuccess()
}
