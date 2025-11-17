// Package pdu
package pdu

import (
	"github.com/half-nothing/simple-fsd/src/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/src/interfaces/global"
)

type RemovePilot struct {
	*Base
}

func (c *RemovePilot) Build() []byte {
	if c.raw == nil {
		return MakeProtocolDataUnitPacket(
			c.GetType(),
			c.From,
			c.To,
		)
	}
	return c.raw
}

func (c *RemovePilot) Parse(data []string, raw []byte) (Interface, *fsd.CommandResult) {
	// CES2352 SERVER
	// [  0  ] [  1 ]
	if r := c.CheckLength(fsd.ClientCommandRemovePilot, len(data)); r != nil {
		return nil, r
	}
	command := &RemovePilot{
		Base: NewBase(fsd.ClientCommandRemovePilot, data[0], global.FSDServerName),
	}
	command.raw = raw
	return command, fsd.CommandResultSuccess()
}
