// Package pdu
package pdu

import (
	"fmt"
	"strconv"

	"github.com/half-nothing/simple-fsd/src/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/src/interfaces/fsd/rating"
	"github.com/half-nothing/simple-fsd/src/interfaces/global"
	"github.com/half-nothing/simple-fsd/src/utils"
)

type AddPilot struct {
	*Base
	Cid      string
	Password string
	Rating   *rating.Rating
	Protocol int
	SimType  int
	RealName string
}

func NewAddPilot() *AddPilot {
	return &AddPilot{
		Base: NewBase(fsd.ClientCommandAddPilot, "", global.FSDServerName),
	}
}

func (c *AddPilot) Build() []byte {
	return MakeProtocolDataUnitPacket(
		c.GetType(),
		c.From,
		c.To,
		c.Cid,
		"",
		strconv.Itoa(c.Rating.Value),
		strconv.Itoa(c.Protocol),
		strconv.Itoa(c.SimType),
		c.RealName,
	)
}

func (c *AddPilot) Parse(data []string, raw []byte) (Interface, *fsd.CommandResult) {
	// CES2352 SERVER 2352 123456  1   9  16  Half_nothing ZGHA
	// [  0  ] [  1 ] [ 2] [  3 ] [4] [5] [6] [       7       ]

	// BAW421 SERVER 2352 token  1  101  6  Half_nothing ZGHA
	// [  0 ] [  1 ] [ 2] [ 3 ] [4] [5] [6] [       7       ]
	if r := c.CheckLength(fsd.ClientCommandAddPilot, len(data)); r != nil {
		return nil, r
	}
	command := &AddPilot{
		Base: NewBase(fsd.ClientCommandAddPilot, data[0], global.FSDServerName),
	}
	command.Cid = data[2]
	command.Password = data[3]
	reqRating := utils.StrToInt(data[4], -1)
	if reqRating == -1 || !rating.Ratings.IsValidEnum(reqRating) {
		return nil, fsd.CommandSyntaxError(false, "Rating", fmt.Errorf("rating(%s) not vaild", data[4]))
	}
	command.Protocol = utils.StrToInt(data[5], -1)
	if command.Protocol == -1 {
		return nil, fsd.CommandResultError(fsd.ClientErrorInvalidProtocolVision, false, "Protocol", fmt.Errorf("unsupport Protocol %s", data[5]))
	}
	command.SimType = utils.StrToInt(data[6], -1)
	if command.SimType == -1 {
		return nil, fsd.CommandSyntaxError(false, "SimType", fmt.Errorf("unsupport SimType %s", data[6]))
	}
	command.RealName = data[7]
	command.raw = raw
	return command, fsd.CommandResultSuccess()
}
