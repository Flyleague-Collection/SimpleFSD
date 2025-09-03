package database

import (
	"fmt"
	c "github.com/half-nothing/simple-fsd/internal/config"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"github.com/half-nothing/simple-fsd/internal/utils"
)

type FlightPlanOperation struct {
	config *c.Config
}

func NewFlightPlanOperation(config *c.Config) *FlightPlanOperation {
	return &FlightPlanOperation{config: config}
}

func (flightPlanOperation *FlightPlanOperation) NewFlightPlan(user *User, callsign string, flightPlanData []string) (flightPlan *FlightPlan, err error) {
	if len(flightPlanData) < 17 {
		return nil, ErrFlightPlanDataTooShort
	}
	flightPlan.Cid = user.Cid
	flightPlan.Callsign = callsign
	flightPlanOperation.updateFlightPlanData(flightPlan, flightPlanData)
	return flightPlan, nil
}

func (flightPlanOperation *FlightPlanOperation) UpdateFlightPlan(flightPlan *FlightPlan, callsign string, flightPlanData []string, isAtc bool) (err error) {
	if flightPlan.Locked && !isAtc {
		return ErrFlightPlanLocked
	}
	if len(flightPlanData) < 17 {
		return ErrFlightPlanDataTooShort
	}
	flightPlan.Callsign = callsign
	flightPlanOperation.updateFlightPlanData(flightPlan, flightPlanData)
	return nil
}

func (flightPlanOperation *FlightPlanOperation) updateFlightPlanData(flightPlan *FlightPlan, flightPlanData []string) {
	flightPlan.FlightType = flightPlanData[2]
	flightPlan.AircraftType = flightPlanData[3]
	flightPlan.Tas = utils.StrToInt(flightPlanData[4], 100)
	flightPlan.DepartureAirport = flightPlanData[5]
	flightPlan.DepartureTime = utils.StrToInt(flightPlanData[6], 0)
	flightPlan.AtcDepartureTime = utils.StrToInt(flightPlanData[7], 0)
	flightPlan.CruiseAltitude = flightPlanData[8]
	flightPlan.ArrivalAirport = flightPlanData[9]
	flightPlan.RouteTimeHour = flightPlanData[10]
	flightPlan.RouteTimeMinute = flightPlanData[11]
	flightPlan.FuelTimeHour = flightPlanData[12]
	flightPlan.FuelTimeMinute = flightPlanData[13]
	flightPlan.AlternateAirport = flightPlanData[14]
	flightPlan.Remarks = flightPlanData[15]
	flightPlan.Route = flightPlanData[16]
}

func (flightPlanOperation *FlightPlanOperation) ToString(flightPlan *FlightPlan, receiver string) string {
	return fmt.Sprintf("$FP%s:%s:%s:%s:%d:%s:%d:%d:%s:%s:%s:%s:%s:%s:%s:%s:%s\r\n",
		flightPlan.Callsign, receiver, flightPlan.FlightType, flightPlan.AircraftType, flightPlan.Tas,
		flightPlan.DepartureAirport, flightPlan.DepartureTime, flightPlan.AtcDepartureTime, flightPlan.CruiseAltitude,
		flightPlan.ArrivalAirport, flightPlan.RouteTimeHour, flightPlan.RouteTimeMinute, flightPlan.FuelTimeHour,
		flightPlan.FuelTimeMinute, flightPlan.AlternateAirport, flightPlan.Remarks, flightPlan.Route)
}
