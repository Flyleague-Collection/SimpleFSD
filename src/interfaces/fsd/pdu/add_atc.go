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

type AddAtc struct {
	*Base
	RealName string
	Cid      string
	Password string
	Rating   *rating.Rating
	Protocol int
	Position *global.Position
}

func (c *AddAtc) Build() []byte {
	return MakeProtocolDataUnitPacket(
		c.GetType(),
		c.From,
		c.To,
		c.RealName,
		c.Cid,
		"",
		strconv.Itoa(c.Rating.Value),
	)
}

func (c *AddAtc) Parse(data []string, raw []byte) (Interface, *fsd.CommandResult) {
	// 2352_OBS SERVER 2352 2352 token  1  100
	// [   0  ] [  1 ] [ 2] [ 3] [ 4 ] [5] [6]
	if r := c.CheckLength(fsd.ClientCommandAddAtc, len(data)); r != nil {
		return nil, r
	}
	command := &AddAtc{
		Base: NewBase(fsd.ClientCommandAddAtc, data[0], global.FSDServerName),
	}
	command.RealName = data[2]
	command.Cid = data[3]
	command.Password = data[4]
	reqRating := utils.StrToInt(data[5], -1)
	if reqRating == -1 || !rating.Ratings.IsValidEnum(reqRating) {
		return nil, fsd.CommandSyntaxError(false, "rating", fmt.Errorf("rating(%s) not vaild", data[5]))
	}
	command.Rating = (*rating.Rating)(rating.Ratings.GetEnum(reqRating))
	command.Protocol = utils.StrToInt(data[6], -1)
	if command.Protocol == -1 {
		return nil, fsd.CommandSyntaxError(false, "protocol", fmt.Errorf("unsupport Protocol %s", data[6]))
	}
	if *global.Vatsim {
		return command, fsd.CommandResultSuccess()
	}
	// 2352_OBS SERVER 2352 2352 123456  1   9   1   0  29.86379 119.49287 100
	// [   0  ] [  1 ] [ 2] [ 3] [  4 ] [5] [6] [7] [8] [   9  ] [   10  ] [11]
	command.Position = &global.Position{
		Latitude:  utils.StrToFloat(data[9], 0),
		Longitude: utils.StrToFloat(data[10], 0),
	}
	if !command.Position.Valid() {
		return nil, fsd.CommandSyntaxError(false, "position", fmt.Errorf("position(%s,%s) not vaild", data[9], data[10]))
	}
	command.raw = raw
	return command, fsd.CommandResultSuccess()
}
