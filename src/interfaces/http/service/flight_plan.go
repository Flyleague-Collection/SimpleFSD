// Package service
package service

import (
	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
)

var (
	ErrFlightPlanNotFound       = NewApiStatus("FLIGHT_PLAN_NOT_FOUND", "飞行计划不存在", NotFound)
	ErrFlightPlanLocked         = NewApiStatus("FLIGHT_PLAN_LOCKED", "飞行计划已锁定", Conflict)
	ErrFlightPlanUnlocked       = NewApiStatus("FLIGHT_PLAN_UNLOCKED", "飞行计划未锁定", Conflict)
	SuccessSubmitFlightPlan     = NewApiStatus("SUBMIT_FLIGHT_PLAN", "成功提交计划", Ok)
	SuccessGetFlightPlan        = NewApiStatus("GET_FLIGHT_PLAN", "成功获取计划", Ok)
	SuccessGetFlightPlans       = NewApiStatus("GET_FLIGHT_PLANS", "成功获取计划", Ok)
	SuccessDeleteSelfFlightPlan = NewApiStatus("DELETE_SELF_FLIGHT_PLAN", "成功删除自己的飞行计划", Ok)
	SuccessDeleteFlightPlan     = NewApiStatus("DELETE_FLIGHT_PLAN", "成功删除飞行计划", Ok)
	SuccessLockFlightPlan       = NewApiStatus("LOCK_FLIGHT_PLAN", "成功修改计划锁定状态", Ok)
)

type FlightPlanServiceInterface interface {
	SubmitFlightPlan(req *RequestSubmitFlightPlan) *ApiResponse[ResponseSubmitFlightPlan]
	GetFlightPlan(req *RequestGetFlightPlan) *ApiResponse[ResponseGetFlightPlan]
	GetFlightPlans(req *RequestGetFlightPlans) *ApiResponse[ResponseGetFlightPlans]
	DeleteSelfFlightPlan(req *RequestDeleteSelfFlightPlan) *ApiResponse[ResponseDeleteSelfFlightPlan]
	DeleteFlightPlan(req *RequestDeleteFlightPlan) *ApiResponse[ResponseDeleteFlightPlan]
	LockFlightPlan(req *RequestLockFlightPlan) *ApiResponse[ResponseLockFlightPlan]
}

type RequestSubmitFlightPlan struct {
	JwtHeader
	*entity.FlightPlan
}

type ResponseSubmitFlightPlan bool

type RequestGetFlightPlan struct {
	JwtHeader
}

type ResponseGetFlightPlan struct {
	*entity.FlightPlan
}

type RequestGetFlightPlans struct {
	JwtHeader
	Page     int `query:"page_number"`
	PageSize int `query:"page_size"`
}

type ResponseGetFlightPlans struct {
	Items    []*entity.FlightPlan `json:"items"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"page_size"`
	Total    int64                `json:"total"`
}

type RequestDeleteSelfFlightPlan struct {
	JwtHeader
	EchoContentHeader
}

type ResponseDeleteSelfFlightPlan bool

type RequestDeleteFlightPlan struct {
	JwtHeader
	EchoContentHeader
	TargetCid int `param:"cid"`
}

type ResponseDeleteFlightPlan bool

type RequestLockFlightPlan struct {
	JwtHeader
	EchoContentHeader
	TargetCid int `param:"cid"`
	Lock      bool
}

type ResponseLockFlightPlan bool
