// Package pdu
package pdu

import (
	"github.com/half-nothing/simple-fsd/src/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/src/utils"
)

type FastPositionFast struct {
	*Base
	Lat               float64
	Lon               float64
	AltitudeTrue      float64
	AltitudeAgl       float64
	Pitch             float64
	Heading           float64
	Bank              float64
	OnGround          bool
	VelocityLongitude float64
	VelocityAltitude  float64
	VelocityLatitude  float64
	VelocityPitch     float64
	VelocityHeading   float64
	VelocityBank      float64
	NoseGearAngle     float64
}

func (c *FastPositionFast) Parse(data []string, raw []byte) (Interface, *fsd.CommandResult) {
	// BAW421 35.338120 25.179485 109.37 0.41 4282382704 0.1836 -0.0125 0.2964 -0.0006 0.0003 -0.0001 -1.64
	// [ 0  ] [    1  ] [    2  ] [  3 ] [ 4] [    5   ] [  6 ] [   7 ] [  8 ] [   9 ] [ 10 ] [  11 ] [ 12]
	if r := c.CheckLength(fsd.ClientCommandFastPositionFast, len(data)); r != nil {
		return nil, r
	}
	command := &FastPositionFast{
		Base:              NewBase(fsd.ClientCommandFastPositionFast, data[0], ""),
		Lat:               utils.StrToFloat(data[1], 0),
		Lon:               utils.StrToFloat(data[2], 0),
		AltitudeTrue:      utils.StrToFloat(data[3], 0),
		AltitudeAgl:       utils.StrToFloat(data[4], 0),
		VelocityLongitude: utils.StrToFloat(data[6], 0),
		VelocityAltitude:  utils.StrToFloat(data[7], 0),
		VelocityLatitude:  utils.StrToFloat(data[8], 0),
		VelocityPitch:     utils.StrToFloat(data[9], 0),
		VelocityHeading:   utils.StrToFloat(data[10], 0),
		VelocityBank:      utils.StrToFloat(data[11], 0),
		NoseGearAngle:     utils.StrToFloat(data[12], 0),
	}
	command.Pitch, command.Heading, command.Bank, command.OnGround = utils.UnpackPBH(uint32(utils.StrToInt(data[5], 0)))
	command.raw = raw
	return command, fsd.CommandResultSuccess()
}
