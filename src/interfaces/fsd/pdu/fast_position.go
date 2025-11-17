// Package pdu
package pdu

import (
	"github.com/half-nothing/simple-fsd/src/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/src/interfaces/global"
)

type FastPosition struct {
	*Base
	Enable bool
}

func NewFastPosition(to string, enable bool) *FastPosition {
	return &FastPosition{
		Base:   NewBase(fsd.ClientCommandSendFastPosition, global.FSDServerName, to),
		Enable: enable,
	}
}

func (c *FastPosition) Build() []byte {
	// SERVER BAW421  1
	// [  0 ] [  1 ] [2]
	if c.Enable {
		return MakeProtocolDataUnitPacket(c.GetType(), c.From, c.To, "1")
	}
	return MakeProtocolDataUnitPacket(c.GetType(), c.From, c.To, "0")
}
