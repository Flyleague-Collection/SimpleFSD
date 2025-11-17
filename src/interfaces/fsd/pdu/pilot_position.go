// Package pdu
package pdu

import (
	"fmt"

	"github.com/half-nothing/simple-fsd/src/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/src/interfaces/fsd/rating"
	"github.com/half-nothing/simple-fsd/src/interfaces/global"
	"github.com/half-nothing/simple-fsd/src/utils"
)

type PilotPosition struct {
	*Base
	SquawkMode       string
	SquawkingModeC   bool
	Identing         bool
	SquawkCode       string
	Rating           *rating.Rating
	Position         *global.Position
	TrueAltitude     int
	PressureAltitude int
	GroundSpeed      int
	Pitch            float64
	Bank             float64
	Heading          float64
	OnGround         bool
}

func (c *PilotPosition) Parse(data []string, raw []byte) (Interface, *fsd.CommandResult) {
	//  S  CPA421 7000  1  38.96244 121.53479 87   0  4290770974 278
	// [0] [  1 ] [ 2] [3] [   4  ] [   5   ] [6] [7] [    8   ] [9]
	if r := c.CheckLength(fsd.ClientCommandPilotPosition, len(data)); r != nil {
		return nil, r
	}
	command := &PilotPosition{
		Base: NewBase(fsd.ClientCommandPilotPosition, data[1], global.FSDServerName),
	}
	command.SquawkMode = data[0]
	command.SquawkingModeC = data[0] == "N" || data[0] == "Y"
	command.Identing = data[0] == "Y"
	command.SquawkCode = data[2]
	reqRating := utils.StrToInt(data[3], -1)
	if reqRating == -1 || !rating.Ratings.IsValidEnum(reqRating) {
		return nil, fsd.CommandSyntaxError(false, "Rating", fmt.Errorf("rating(%s) not vaild", data[3]))
	}
	command.Rating = (*rating.Rating)(rating.Ratings.GetEnum(reqRating))
	command.Position = &global.Position{
		Latitude:  utils.StrToFloat(data[4], 0),
		Longitude: utils.StrToFloat(data[5], 0),
	}
	if !command.Position.Valid() {
		return nil, fsd.CommandSyntaxError(false, "Position", fmt.Errorf("position(%s,%s) not vaild", data[4], data[5]))
	}
	command.TrueAltitude = utils.StrToInt(data[6], -1)
	if command.TrueAltitude == -1 {
		return nil, fsd.CommandSyntaxError(false, "TrueAltitude", fmt.Errorf("trueAltitude(%s) not vaild", data[6]))
	}
	diff := utils.StrToInt(data[7], -1)
	if diff == -1 {
		return nil, fsd.CommandSyntaxError(false, "AltitudeDiff", fmt.Errorf("altitudeDiff(%s) not vaild", data[7]))
	}
	command.PressureAltitude = command.TrueAltitude + diff
	pbh := utils.StrToInt(data[8], -1)
	if pbh == -1 {
		return nil, fsd.CommandSyntaxError(false, "PitchBankHeading", fmt.Errorf("pitchBankHeading(%s) not vaild", data[8]))
	}
	command.Pitch, command.Bank, command.Heading, command.OnGround = utils.UnpackPBH(uint32(pbh))
	command.GroundSpeed = utils.StrToInt(data[9], -1)
	if command.GroundSpeed == -1 {
		return nil, fsd.CommandSyntaxError(false, "GroundSpeed", fmt.Errorf("groundSpeed(%s) not vaild", data[9]))
	}
	command.raw = raw
	return command, fsd.CommandResultSuccess()
}
