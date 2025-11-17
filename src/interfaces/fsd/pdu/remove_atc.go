// Package pdu
package pdu

import (
	"github.com/half-nothing/simple-fsd/src/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/src/interfaces/global"
)

type RemoveAtc struct {
	*Base
}

func NewRemoveATC() *RemoveAtc {
	return &RemoveAtc{
		Base: NewBase(fsd.ClientCommandRemoveAtc, "", global.FSDServerName),
	}
}

func (c *RemoveAtc) Build() []byte {
	if c.raw == nil {
		return MakeProtocolDataUnitPacket(
			c.GetType(),
			c.From,
			c.To,
		)
	}
	return c.raw
}

func (c *RemoveAtc) Parse(data []string, raw []byte) (Interface, *fsd.CommandResult) {
	// ZGGG_CTR SERVER
	// [   0  ] [  1 ]
	if r := c.CheckLength(fsd.ClientCommandRemoveAtc, len(data)); r != nil {
		return nil, r
	}
	command := &RemoveAtc{
		Base: NewBase(fsd.ClientCommandRemoveAtc, data[0], global.FSDServerName),
	}
	command.raw = raw
	return command, fsd.CommandResultSuccess()
}
