// Package pdu
package pdu

import (
	"github.com/half-nothing/simple-fsd/src/interfaces/fsd"
)

type RequestHandoff struct {
	*Base
	Target string
}

func (c *RequestHandoff) Parse(data []string, raw []byte) (Interface, *fsd.CommandResult) {
	// ZSSS_APP ZSHA_CTR CES2352
	// [  0   ] [   1  ] [  2  ]
	if r := c.CheckLength(fsd.ClientCommandRequestHandoff, len(data)); r != nil {
		return nil, r
	}
	command := &RequestHandoff{
		Base:   NewBase(fsd.ClientCommandRequestHandoff, data[0], data[1]),
		Target: data[2],
	}
	command.raw = raw
	return command, fsd.CommandResultSuccess()
}
