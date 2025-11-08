// Package service
package service

import (
	"time"

	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
)

var (
	ErrRecordNotFound                 = NewApiStatus("RECORD_NOT_FOUND", "管制员履历不存在", NotFound)
	ErrSameRating                     = NewApiStatus("SAME_RATING", "用户已是该权限", BadRequest)
	SuccessGetControllers             = NewApiStatus("GET_CONTROLLER_PAGE", "获取管制员信息分页成功", Ok)
	SuccessGetCurrentControllerRecord = NewApiStatus("GET_CURRENT_CONTROLLER_RECORD", "获取管制员履历成功", Ok)
	SuccessGetControllerRatings       = NewApiStatus("GET_CONTROLLER_RATINGS", "成功获取权限公示", Ok)
	SuccessUpdateControllerRating     = NewApiStatus("UPDATE_CONTROLLER_RATING", "编辑用户管制权限成功", Ok)
	SuccessAddControllerRecord        = NewApiStatus("ADD_CONTROLLER_RECORD", "添加管制员履历成功", Ok)
	SuccessDeleteControllerRecord     = NewApiStatus("DELETE_CONTROLLER_RECORD", "删除管制员履历成功", Ok)
)

type ControllerServiceInterface interface {
	GetControllerList(req *RequestControllerList) *ApiResponse[ResponseControllerList]
	GetCurrentControllerRecord(req *RequestGetCurrentControllerRecord) *ApiResponse[ResponseGetCurrentControllerRecord]
	GetControllerRecord(req *RequestGetControllerRecord) *ApiResponse[ResponseGetControllerRecord]
	GetControllerRatings(req *RequestControllerRatingList) *ApiResponse[ResponseControllerRatingList]
	UpdateControllerRating(req *RequestUpdateControllerRating) *ApiResponse[ResponseUpdateControllerRating]
	AddControllerRecord(req *RequestAddControllerRecord) *ApiResponse[ResponseAddControllerRecord]
	DeleteControllerRecord(req *RequestDeleteControllerRecord) *ApiResponse[ResponseDeleteControllerRecord]
}

type RequestControllerList struct {
	JwtHeader
	Page     int `query:"page_number"`
	PageSize int `query:"page_size"`
}

type ResponseControllerList struct {
	Items    []*entity.User `json:"items"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
	Total    int64          `json:"total"`
}

type RequestGetCurrentControllerRecord struct {
	JwtHeader
	Page     int `query:"page_number"`
	PageSize int `query:"page_size"`
}

type ResponseGetCurrentControllerRecord struct {
	Items    []*entity.ControllerRecord `json:"items"`
	Page     int                        `json:"page"`
	PageSize int                        `json:"page_size"`
	Total    int64                      `json:"total"`
}

type RequestGetControllerRecord struct {
	JwtHeader
	TargetUid uint `param:"uid"`
	Page      int  `query:"page_number"`
	PageSize  int  `query:"page_size"`
}

type ResponseGetControllerRecord struct {
	Items    []*entity.ControllerRecord `json:"items"`
	Page     int                        `json:"page"`
	PageSize int                        `json:"page_size"`
	Total    int64                      `json:"total"`
}

type RequestUpdateControllerRating struct {
	JwtHeader
	EchoContentHeader
	TargetUid    uint      `param:"uid"`
	Guest        bool      `json:"guest"`
	Rating       int       `json:"rating"`
	UnderMonitor bool      `json:"under_monitor"`
	UnderSolo    bool      `json:"under_solo"`
	Tier2        bool      `json:"tier2"`
	SoloUntil    time.Time `json:"solo_until"`
}

type ResponseUpdateControllerRating bool

type RequestAddControllerRecord struct {
	JwtHeader
	EchoContentHeader
	TargetUid uint   `param:"uid"`
	Type      int    `json:"type"`
	Content   string `json:"content"`
}

type ResponseAddControllerRecord bool

type RequestDeleteControllerRecord struct {
	JwtHeader
	EchoContentHeader
	TargetUid    uint `param:"uid"`
	TargetRecord uint `param:"rid"`
}

type ResponseDeleteControllerRecord bool

type RequestControllerRatingList struct {
	Page     int `query:"page_number"`
	PageSize int `query:"page_size"`
}

type ControllerRating struct {
	Cid          int       `json:"cid"`
	Rating       int       `json:"rating"`
	AvatarUrl    string    `json:"avatar_url"`
	UnderMonitor bool      `json:"under_monitor"`
	UnderSolo    bool      `json:"under_solo"`
	SoloUntil    time.Time `json:"solo_until"`
	Tier2        bool      `json:"tier2"`
	IsGuest      bool      `json:"is_guest"`
}

type ResponseControllerRatingList struct {
	Items    []*ControllerRating `json:"items"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"page_size"`
	Total    int64               `json:"total"`
}
