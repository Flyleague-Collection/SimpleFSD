// Package service
package service

import (
	"github.com/half-nothing/simple-fsd/src/interfaces/operation"
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
	ErrActivityLocked         = NewApiStatus("ACTIVITY_LOCKED", "活动报名信息已锁定", Conflict)
	ErrActivityIdMismatch     = NewApiStatus("ACTIVITY_ID_MISMATCH", "活动ID不正确", Conflict)
	ErrActivityNotFound       = NewApiStatus("ACTIVITY_NOT_FOUND", "活动不存在", NotFound)
	ErrFacilityNotFound       = NewApiStatus("FACILITY_NOT_FOUND", "管制席位不存在", NotFound)
	ErrParseTime              = NewApiStatus("TIME_FORMAT_ERROR", "格式错误", BadRequest)
	ErrRatingTooLow           = NewApiStatus("RATING_TOO_LOW", "管制权限不够", PermissionDenied)
	ErrFacilityAlreadyExist   = NewApiStatus("FACILITY_ALREADY_EXIST", "你不能同时报名两个以上的席位", Conflict)
	ErrFacilityAlreadySigned  = NewApiStatus("FACILITY_ALREADY_SIGNED", "已有其他管制员报名", Conflict)
	ErrFacilityUnSigned       = NewApiStatus("FACILITY_UNSIGNED", "该席位尚未有人报名", Conflict)
	ErrFacilityNotYourSign    = NewApiStatus("FACILITY_NOT_YOUR_SIGN", "这不是你报名的席位", Conflict)
	ErrAlreadySigned          = NewApiStatus("ALREADY_SIGNED", "你已经报名该活动了", Conflict)
	ErrCallsignUsed           = NewApiStatus("CALLSIGN_USED", "呼号已被占用", Conflict)
	ErrNoSigned               = NewApiStatus("NO_SIGNED", "你还没有报名该活动", Conflict)
	SuccessGetActivities      = NewApiStatus("GET_ACTIVITIES", "成功获取活动", Ok)
	SuccessGetActivitiesPage  = NewApiStatus("GET_ACTIVITIES_PAGE", "成功获取活动分页", Ok)
	SuccessGetActivityInfo    = NewApiStatus("GET_ACTIVITY_INFO", "成功获取活动信息", Ok)
	SuccessAddActivity        = NewApiStatus("ADD_ACTIVITY", "成功添加活动", Ok)
	SuccessDeleteActivity     = NewApiStatus("DELETE_ACTIVITY", "成功删除活动", Ok)
	SuccessSignFacility       = NewApiStatus("SIGNED_FACILITY", "报名成功", Ok)
	SuccessUnsignFacility     = NewApiStatus("UNSIGNED_FACILITY", "成功取消报名", Ok)
	SuccessSignedActivity     = NewApiStatus("SIGNED_ACTIVITY", "报名成功", Ok)
	SuccessUnsignedActivity   = NewApiStatus("UNSIGNED_ACTIVITY", "取消报名成功", Ok)
	SuccessEditActivity       = NewApiStatus("EDIT_ACTIVITY", "修改活动成功", Ok)
	SuccessEditActivityStatus = NewApiStatus("EDIT_ACTIVITY_STATUS", "成功修改活动状态", Ok)
	SuccessEditPilotsStatus   = NewApiStatus("EDIT_PILOTS_STATUS", "成功修改活动机组状态", Ok)
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
