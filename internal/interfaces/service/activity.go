// Package service
package service

import (
	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
)

type ActivityModel struct {
	Id         uint   `json:"id"`
	Publisher  int    `json:"publisher"`
	Title      string `json:"title"`
	ImageUrl   string `json:"image_url"`
	ActiveTime string `json:"active_time"`
	Departure  string `json:"departure"`
	Arrival    string `json:"arrival"`
	Route      string `json:"route"`
	Distance   int    `json:"distance"`
	Status     int    `json:"status"`
	NOTAMS     string `json:"notams"`
}

var (
	ErrActivityLocked     = NewApiStatus("ACTIVITY_LOCKED", "活动报名信息已锁定", Conflict)
	ErrActivityIdMismatch = NewApiStatus("ACTIVITY_ID_MISMATCH", "活动ID不正确", Conflict)
)

type ActivityServiceInterface interface {
	GetActivities(req *RequestGetActivities) *ApiResponse[ResponseGetActivities]
	GetActivitiesPage(req *RequestGetActivitiesPage) *ApiResponse[ResponseGetActivitiesPage]
	GetActivityInfo(req *RequestActivityInfo) *ApiResponse[ResponseActivityInfo]
	AddActivity(req *RequestAddActivity) *ApiResponse[ResponseAddActivity]
	DeleteActivity(req *RequestDeleteActivity) *ApiResponse[ResponseDeleteActivity]
	ControllerJoin(req *RequestControllerJoin) *ApiResponse[ResponseControllerJoin]
	ControllerLeave(req *RequestControllerLeave) *ApiResponse[ResponseControllerLeave]
	PilotJoin(req *RequestPilotJoin) *ApiResponse[ResponsePilotJoin]
	PilotLeave(req *RequestPilotLeave) *ApiResponse[ResponsePilotLeave]
	EditActivity(req *RequestEditActivity) *ApiResponse[ResponseEditActivity]
	EditPilotStatus(req *RequestEditPilotStatus) *ApiResponse[ResponseEditPilotStatus]
	EditActivityStatus(req *RequestEditActivityStatus) *ApiResponse[ResponseEditActivityStatus]
}

type RequestGetActivities struct {
	Time string `query:"time"`
}

type ResponseGetActivities []*operation.Activity

type RequestGetActivitiesPage struct {
	JwtHeader
	Page     int `query:"page_number"`
	PageSize int `query:"page_size"`
}

type ResponseGetActivitiesPage struct {
	Items    []*operation.Activity `json:"items"`
	Page     int                   `json:"page"`
	PageSize int                   `json:"page_size"`
	Total    int64                 `json:"total"`
}

type RequestActivityInfo struct {
	ActivityId uint `param:"activity_id"`
}

type ResponseActivityInfo operation.Activity

type RequestAddActivity struct {
	JwtHeader
	EchoContentHeader
	*operation.Activity
}

type ResponseAddActivity bool

type RequestDeleteActivity struct {
	JwtHeader
	EchoContentHeader
	ActivityId uint `param:"activity_id"`
}

type ResponseDeleteActivity bool

type RequestControllerJoin struct {
	JwtHeader
	Rating     int
	ActivityId uint `param:"activity_id"`
	FacilityId uint `param:"facility_id"`
}

type ResponseControllerJoin bool

type RequestControllerLeave struct {
	JwtHeader
	ActivityId uint `param:"activity_id"`
	FacilityId uint `param:"facility_id"`
}

type ResponseControllerLeave bool

type RequestPilotJoin struct {
	JwtHeader
	ActivityId   uint   `param:"activity_id"`
	Callsign     string `json:"callsign"`
	AircraftType string `json:"aircraft_type"`
}

type ResponsePilotJoin bool

type RequestPilotLeave struct {
	JwtHeader
	ActivityId uint `param:"activity_id"`
}

type ResponsePilotLeave bool

type RequestEditActivity struct {
	JwtHeader
	EchoContentHeader
	*operation.Activity
}

type ResponseEditActivity bool

type RequestEditActivityStatus struct {
	JwtHeader
	ActivityId uint `param:"activity_id"`
	Status     int  `json:"status"`
}

type ResponseEditActivityStatus bool

type RequestEditPilotStatus struct {
	JwtHeader
	ActivityId uint `param:"activity_id"`
	UserId     uint `param:"user_id"`
	Status     int  `json:"status"`
}

type ResponseEditPilotStatus bool
