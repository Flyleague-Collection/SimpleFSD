// Package pdu
package pdu

import (
	"fmt"

	"github.com/half-nothing/simple-fsd/src/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/src/utils"
)

type AtcEditPlan struct {
	*Base
	Callsign  string
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

func (c *AtcEditPlan) Parse(data []string, raw []byte) (Interface, *fsd.CommandResult) {
	// ZYSH_CTR SERVER CPA421  I  H/A320/L 474 ZYTL 1115  0  FL371 ZYHB  11   8    22   6   ZYCC
	// [   0  ] [  1 ] [  2 ] [3] [   4  ] [5] [ 6] [ 7] [8] [ 9 ] [10] [11] [12] [13] [14] [15]
	// /V/ SEL/AHFL CHI19D/28 VENOS A588 NULRA W206 MAGBI W656 ISLUK W629 LARUN
	// [     16   ] [                             17                          ]
	if r := c.CheckLength(fsd.ClientCommandAtcEditPlan, len(data)); r != nil {
		return nil, r
	}
	command := &AtcEditPlan{
		Base:      NewBase(fsd.ClientCommandAtcEditPlan, data[0], data[1]),
		Callsign:  data[2],
		Type:      data[3],
		Aircraft:  data[4],
		TAS:       utils.StrToInt(data[5], -1),
		DEP:       data[6],
		EOBT:      data[7],
		ACTEOBT:   data[8],
		Altitude:  data[9],
		ARR:       data[10],
		RouteHour: data[11],
		RouteMin:  data[12],
		AirHour:   data[13],
		AirMin:    data[14],
		ALTE:      data[15],
		Remarks:   data[16],
		Route:     data[17],
	}
	if command.TAS == -1 {
		return nil, fsd.CommandSyntaxError(false, "TAS", fmt.Errorf("TAS(%s) not vaild", data[4]))
	}
	command.raw = raw
	return command, fsd.CommandResultSuccess()
}
