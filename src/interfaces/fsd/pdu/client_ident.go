// Package pdu
package pdu

import "github.com/half-nothing/simple-fsd/src/interfaces/fsd"

type ClientIdent struct {
	*Base
	Id        string
	Name      string
	Major     string
	Minor     string
	Cid       string
	MachineId string
}

func (c *ClientIdent) Parse(data []string, raw []byte) (Interface, *fsd.CommandResult) {
	// ZSHA_CTR SERVER 69d7 EuroScope 3.2  3   2  2352 370598540
	// [   0  ] [  1 ] [ 2] [     3     ] [4] [5] [ 6] [   7   ]
	if r := c.CheckLength(fsd.ClientCommandClientIdent, len(data)); r != nil {
		return nil, r
	}
	command := &ClientIdent{
		Base:      NewBase(fsd.ClientCommandClientIdent, data[0], data[1]),
		Id:        data[2],
		Name:      data[3],
		Major:     data[4],
		Minor:     data[5],
		Cid:       data[6],
		MachineId: data[7],
	}
	command.raw = raw
	return command, fsd.CommandResultSuccess()
}
