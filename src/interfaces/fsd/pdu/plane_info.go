// Package pdu
package pdu

import (
	"fmt"
	"strings"

	"github.com/half-nothing/simple-fsd/src/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/src/utils"
)

type PlaneInfo struct {
	*Base
	SubType   string
	PlaneData map[string]string
}

func (c *PlaneInfo) Parse(data []string, raw []byte) (Interface, *fsd.CommandResult) {
	// ZSHA_CTR CES2352 PIR
	// [   0  ] [  1  ] [2]
	if r := c.CheckLength(fsd.ClientCommandPlaneInfo, len(data)); r != nil {
		return nil, r
	}
	command := &PlaneInfo{
		Base:      NewBase(fsd.ClientCommandPlaneInfo, data[0], data[1]),
		SubType:   data[2],
		PlaneData: make(map[string]string),
	}
	command.raw = raw
	if len(data) == 3 {
		return command, fsd.CommandResultSuccess()
	}
	// ZSHA_CTR CES2352 PI  GEN EQUIPMENT=B738 AIRLINE=CDG
	// [   0  ] [  1  ] [2] [3] [      4     ] [     5   ]
	if len(data) < 5 {
		return nil, fsd.CommandSyntaxError(false, "length", fmt.Errorf("command length %d less than 5", len(data)))
	}
	if data[2] != "PI" || data[3] != "GEN" {
		return command, fsd.CommandResultSuccess()
	}
	utils.ForEach(data[4:], func(index int, element string) {
		kv := strings.Split(element, "=")
		if len(kv) != 2 {
			return
		}
		command.PlaneData[kv[0]] = kv[1]
	})
	return command, fsd.CommandResultSuccess()
}
