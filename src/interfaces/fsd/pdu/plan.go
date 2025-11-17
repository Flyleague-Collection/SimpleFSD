// Package pdu
package pdu

import (
	"fmt"
	"strconv"

	"github.com/half-nothing/simple-fsd/src/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/src/utils"
)

type Plan struct {
	*Base
	Type      string
	Aircraft  string
	TAS       int
	DEP       string
	EOBT      string
	ACTEOBT   string
	Altitude  string
	ARR       string
	RouteHour string
	RouteMin  string
	AirHour   string
	AirMin    string
	ALTE      string
	Remarks   string
	Route     string
}

func NewPlan(
	from string,
	to string,
	type_ string,
	aircraft string,
	tas int,
	dep string,
	eobt string,
	acteobt string,
	altitude string,
	arr string,
	routeHour string,
	routeMin string,
	airHour string,
	airMin string,
	alte string,
	remarks string,
	route string,
) *Plan {
	return &Plan{
		Base:      NewBase(fsd.ClientCommandPlan, from, to),
		Type:      type_,
		Aircraft:  aircraft,
		TAS:       tas,
		DEP:       dep,
		EOBT:      eobt,
		ACTEOBT:   acteobt,
		Altitude:  altitude,
		ARR:       arr,
		RouteHour: routeHour,
		RouteMin:  routeMin,
		AirHour:   airHour,
		AirMin:    airMin,
		ALTE:      alte,
		Remarks:   remarks,
		Route:     route,
	}
}

func (c *Plan) Build() []byte {
	return MakeProtocolDataUnitPacket(
		c.GetType(),
		c.From,
		fsd.BroadcastTargetAllATC.Value,
		c.Type,
		c.Aircraft,
		strconv.Itoa(c.TAS),
		c.DEP,
		c.EOBT,
		c.ACTEOBT,
		c.Altitude,
		c.ARR,
		c.RouteHour,
		c.RouteMin,
		c.AirHour,
		c.AirMin,
		c.ALTE,
		c.Remarks,
		c.Route,
	)
}

func (c *Plan) Parse(data []string, raw []byte) (Interface, *fsd.CommandResult) {
	// CPA421 SERVER  I  H/A320/L 474 ZYTL 1115  0  FL371 ZYHB  1    18   2    26  ZYCC
	// [  0 ] [  1 ] [2] [  3   ] [4] [ 5] [ 6] [7] [ 8 ] [9 ] [10] [11] [12] [13] [14]
	// /V/ SEL/AHFL VENOS A588 NULRA W206 MAGBI W656 ISLUK W629 LARUN
	// [    15    ] [                      16                       ]
	if r := c.CheckLength(fsd.ClientCommandPlan, len(data)); r != nil {
		return nil, r
	}
	command := &Plan{
		Base:      NewBase(fsd.ClientCommandPlan, data[0], data[1]),
		Type:      data[2],
		Aircraft:  data[3],
		TAS:       utils.StrToInt(data[4], -1),
		DEP:       data[5],
		EOBT:      data[6],
		ACTEOBT:   data[7],
		Altitude:  data[8],
		ARR:       data[9],
		RouteHour: data[10],
		RouteMin:  data[11],
		AirHour:   data[12],
		AirMin:    data[13],
		ALTE:      data[14],
		Remarks:   data[15],
		Route:     data[16],
	}
	if command.TAS == -1 {
		return nil, fsd.CommandSyntaxError(false, "TAS", fmt.Errorf("TAS(%s) not vaild", data[4]))
	}
	command.raw = raw
	return command, fsd.CommandResultSuccess()
}
