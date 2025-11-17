// Package pdu
package pdu

import (
	"fmt"

	"github.com/half-nothing/simple-fsd/src/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/src/interfaces/global"
	"github.com/half-nothing/simple-fsd/src/utils"
)

type AtcSubPosition struct {
	*Base
	VisIndex int
	Position *global.Position
}

func (c *AtcSubPosition) Parse(data []string, raw []byte) (Interface, *fsd.CommandResult) {
	// ZSHA_CTR  0  36.67349 120.45621
	// [   0  ] [1] [   2  ] [   3   ]
	if r := c.CheckLength(fsd.ClientCommandAtcSubVisPoint, len(data)); r != nil {
		return nil, r
	}
	command := &AtcSubPosition{
		Base: NewBase(fsd.ClientCommandAtcSubVisPoint, data[0], ""),
	}
	command.VisIndex = utils.StrToInt(data[1], -1)
	if command.VisIndex == -1 {
		return nil, fsd.CommandSyntaxError(false, "visIndex", fmt.Errorf("visIndex(%s) not vaild", data[1]))
	}
	command.Position = &global.Position{
		Latitude:  utils.StrToFloat(data[2], 0),
		Longitude: utils.StrToFloat(data[3], 0),
	}
	if !command.Position.Valid() {
		return nil, fsd.CommandSyntaxError(false, "position", fmt.Errorf("position(%s,%s) not vaild", data[2], data[3]))
	}
	command.raw = raw
	return command, fsd.CommandResultSuccess()
}
