// Package pdu
package pdu

import (
	"github.com/half-nothing/simple-fsd/src/interfaces/fsd"
)

type KillClient struct {
	*Base
	Reason string
}

func (c *KillClient) Parse(data []string, raw []byte) (Interface, *fsd.CommandResult) {
	// ZSHA_CTR CPA421 test
	// [   0  ] [  1 ] [ 2]
	if r := c.CheckLength(fsd.ClientCommandKillClient, len(data)); r != nil {
		return nil, r
	}
	command := &KillClient{
		Base: NewBase(fsd.ClientCommandKillClient, data[0], data[1]),
	}
	if len(data) == 3 {
		command.Reason = data[2]
	}
	command.raw = raw
	return command, fsd.CommandResultSuccess()
}
