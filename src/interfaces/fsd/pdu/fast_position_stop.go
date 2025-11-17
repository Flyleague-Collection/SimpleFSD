// Package pdu
package pdu

import (
	"github.com/half-nothing/simple-fsd/src/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/src/utils"
)

type FastPositionStop struct {
	*Base
	Lat           float64
	Lon           float64
	AltitudeTrue  float64
	AltitudeAgl   float64
	Pitch         float64
	Heading       float64
	Bank          float64
	OnGround      bool
	NoseGearAngle float64
}

func (c *FastPositionStop) Parse(data []string, raw []byte) (Interface, *fsd.CommandResult) {
	// BAW421 35.338116 25.179483 109.41 0.41 4282382704 0.00
	// [ 0  ] [    1  ] [    2  ] [  3 ] [ 4] [    5   ] [ 6]
	if r := c.CheckLength(fsd.ClientCommandFastPositionStop, len(data)); r != nil {
		return nil, r
	}
	command := &FastPositionStop{
		Base:          NewBase(fsd.ClientCommandFastPositionStop, data[0], ""),
		Lat:           utils.StrToFloat(data[1], 0),
		Lon:           utils.StrToFloat(data[2], 0),
		AltitudeTrue:  utils.StrToFloat(data[3], 0),
		AltitudeAgl:   utils.StrToFloat(data[4], 0),
		NoseGearAngle: utils.StrToFloat(data[6], 0),
	}
	command.Pitch, command.Heading, command.Bank, command.OnGround = utils.UnpackPBH(uint32(utils.StrToInt(data[5], 0)))
	command.raw = raw
	return command, fsd.CommandResultSuccess()
}
