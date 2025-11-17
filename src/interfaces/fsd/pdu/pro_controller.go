// Package pdu
package pdu

import (
	"github.com/half-nothing/simple-fsd/src/interfaces/fsd"
)

type ProController struct {
	*Base
	Type    string
	SubType string
	Target  string
}

func (c *ProController) Parse(data []string, raw []byte) (Interface, *fsd.CommandResult) {
	// ZSHA_CTR ZSSS_APP CCP  HC CES2352
	// [   0  ] [   1  ] [2] [3] [  4  ]
	if r := c.CheckLength(fsd.ClientCommandProController, len(data)); r != nil {
		return nil, r
	}
	command := &ProController{
		Base: NewBase(fsd.ClientCommandProController, data[0], data[1]),
	}
	command.Type = data[2]
	command.SubType = data[3]
	command.Target = data[4]
	command.raw = raw
	return command, fsd.CommandResultSuccess()
}
