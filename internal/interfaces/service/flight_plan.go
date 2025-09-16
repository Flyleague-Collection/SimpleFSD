// Package service
package service

import "github.com/half-nothing/simple-fsd/internal/interfaces/operation"

type FlightPlanServiceInterface interface {
	SubmitFlightPlan(req *RequestSubmitFlightPlan) *ApiResponse[ResponseSubmitFlightPlan]
	GetFlightPlan(req *RequestGetFlightPlan) *ApiResponse[ResponseGetFlightPlan]
	GetFlightPlans(req *RequestGetFlightPlans) *ApiResponse[ResponseGetFlightPlans]
	DeleteFlightPlan(req *RequestDeleteFlightPlan) *ApiResponse[ResponseDeleteFlightPlan]
	LockFlightPlan(req *RequestLockFlightPlan) *ApiResponse[ResponseLockFlightPlan]
}

type RequestSubmitFlightPlan struct {
	JwtHeader
	*operation.FlightPlan
}

type ResponseSubmitFlightPlan bool

type RequestGetFlightPlan struct {
	JwtHeader
}

type ResponseGetFlightPlan struct {
	*operation.FlightPlan
}

type RequestGetFlightPlans struct {
	JwtHeader
	Page     int `query:"page_number"`
	PageSize int `query:"page_size"`
}

type ResponseGetFlightPlans struct {
	Items    []*operation.FlightPlan `json:"items"`
	Page     int                     `json:"page"`
	PageSize int                     `json:"page_size"`
	Total    int64                   `json:"total"`
}

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
