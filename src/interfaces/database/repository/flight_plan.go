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

// FlightPlanInterface 飞行计划操作接口定义
type FlightPlanInterface interface {
	// GetFlightPlanByCid 通过用户cid获取飞行计划, 当err为nil时返回值flightPlan有效
	GetFlightPlanByCid(cid int) (flightPlan *entity.FlightPlan, err error)
	// UpsertFlightPlan 创建或更新飞行计划, 当err为nil时返回值flightPlan有效
	UpsertFlightPlan(user *entity.User, callsign string, flightPlanData []string) (flightPlan *entity.FlightPlan, err error)
	// UpdateFlightPlanData 更新飞行计划(不提交数据库)
	UpdateFlightPlanData(flightPlan *entity.FlightPlan, flightPlanData []string)
	// UpdateFlightPlan 更新飞行计划(提交数据库), 当err为nil时更新成功
	UpdateFlightPlan(flightPlan *entity.FlightPlan, flightPlanData []string, atcEdit bool) (err error)
	SaveFlightPlan(flightPlan *entity.FlightPlan) (err error)
	GetFlightPlans(page, pageSize int) (flightPlans []*entity.FlightPlan, total int64, err error)
	LockFlightPlan(flightPlan *entity.FlightPlan) (err error)
	UnlockFlightPlan(flightPlan *entity.FlightPlan) (err error)
	DeleteSelfFlightPlan(flightPlan *entity.FlightPlan) (err error)
	DeleteFlightPlan(flightPlan *entity.FlightPlan) (err error)
	// UpdateCruiseAltitude 更新巡航高度, 当err为nil时更新成功
	UpdateCruiseAltitude(flightPlan *entity.FlightPlan, cruiseAltitude string) (err error)
	// ToString 将飞行计划转换为ES和Swift可识别的形式
	ToString(flightPlan *entity.FlightPlan) (str string)
}
