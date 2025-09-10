// Package operation
package operation

import (
	"errors"
)

var (
	ErrFlightPlanDataTooShort = errors.New("flight plan data is too short")
	ErrFlightPlanLocked       = errors.New("flight plan locked")
)

// FlightPlanOperationInterface 飞行计划操作接口定义
type FlightPlanOperationInterface interface {
	// NewFlightPlan 创建或更新飞行计划, 当err为nil时返回值flightPlan有效
	NewFlightPlan(user *User, callsign string, flightPlanData []string) (flightPlan *FlightPlan, err error)
	UpdateFlightPlan(flightPlan *FlightPlan, callsign string, flightPlanData []string, isAtc bool) (err error)
	// ToString 将飞行计划转换为ES和Swift可识别的形式
	ToString(flightPlan *FlightPlan, receiver string) (str string)
}
