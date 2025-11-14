// Package repository
package repository

import (
	"errors"

	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
)

var (
	ErrFlightPlanNotFound     = errors.New("flight plan not found")
	ErrSimulatorServer        = errors.New("simulator server not support flight plan store")
	ErrFlightPlanDataTooShort = errors.New("flight plan data is too short")
	ErrFlightPlanExists       = errors.New("flight plan already exists")
	ErrFlightPlanLocked       = errors.New("flight plan locked")
)

type FlightPlanInterface interface {
	Base[*entity.FlightPlan]
	GetFlightPlanByCid(cid int) (flightPlan *entity.FlightPlan, err error)
	UpsertFlightPlan(user *entity.User, callsign string, flightPlanData []string) (flightPlan *entity.FlightPlan, err error)
	UpdateFlightPlanData(flightPlan *entity.FlightPlan, flightPlanData []string)
	UpdateFlightPlan(flightPlan *entity.FlightPlan, flightPlanData []string, atcEdit bool) (err error)
	SaveFlightPlan(flightPlan *entity.FlightPlan) (err error)
	GetFlightPlans(page, pageSize int) (flightPlans []*entity.FlightPlan, total int64, err error)
	LockFlightPlan(flightPlan *entity.FlightPlan) (err error)
	UnlockFlightPlan(flightPlan *entity.FlightPlan) (err error)
	DeleteSelfFlightPlan(flightPlan *entity.FlightPlan) (err error)
	DeleteFlightPlan(flightPlan *entity.FlightPlan) (err error)
	UpdateCruiseAltitude(flightPlan *entity.FlightPlan, cruiseAltitude string) (err error)
}
