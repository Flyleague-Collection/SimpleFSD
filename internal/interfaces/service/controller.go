// Package service
package service

import (
	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"time"
)

type ControllerServiceInterface interface {
	GetControllerList(req *RequestControllerList) *ApiResponse[ResponseControllerList]
	GetCurrentControllerRecord(req *RequestGetCurrentControllerRecord) *ApiResponse[ResponseGetCurrentControllerRecord]
	GetControllerRecord(req *RequestGetControllerRecord) *ApiResponse[ResponseGetControllerRecord]
	UpdateControllerRating(req *RequestUpdateControllerRating) *ApiResponse[ResponseUpdateControllerRating]
	UpdateControllerUnderMonitor(req *RequestUpdateControllerUnderMonitor) *ApiResponse[ResponseUpdateControllerUnderMonitor]
	UpdateControllerUnderSolo(req *RequestUpdateControllerUnderSolo) *ApiResponse[ResponseUpdateControllerUnderSolo]
	UpdateControllerGuest(req *RequestUpdateControllerGuest) *ApiResponse[ResponseUpdateControllerGuest]
	AddControllerRecord(req *RequestAddControllerRecord) *ApiResponse[ResponseAddControllerRecord]
	DeleteControllerRecord(req *RequestDeleteControllerRecord) *ApiResponse[ResponseDeleteControllerRecord]
}

type RequestControllerList struct {
	JwtHeader
	Page     int `query:"page_number"`
	PageSize int `query:"page_size"`
}

type ResponseControllerList struct {
	Items    []*operation.User `json:"items"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
	Total    int64             `json:"total"`
}

type RequestGetCurrentControllerRecord struct {
	JwtHeader
	Page     int `query:"page_number"`
	PageSize int `query:"page_size"`
}

type ResponseGetCurrentControllerRecord struct {
	Items    []*operation.ControllerRecord `json:"items"`
	Page     int                           `json:"page"`
	PageSize int                           `json:"page_size"`
	Total    int64                         `json:"total"`
}

type RequestGetControllerRecord struct {
	JwtHeader
	TargetCid int `param:"cid"`
	Page      int `query:"page_number"`
	PageSize  int `query:"page_size"`
}

type ResponseGetControllerRecord struct {
	Items    []*operation.ControllerRecord `json:"items"`
	Page     int                           `json:"page"`
	PageSize int                           `json:"page_size"`
	Total    int64                         `json:"total"`
}

type RequestUpdateControllerRating struct {
	JwtHeader
	EchoContentHeader
	TargetUid uint `param:"uid"`
	Rating    int  `json:"rating"`
}

type ResponseUpdateControllerRating bool

type RequestUpdateControllerUnderMonitor struct {
	JwtHeader
	EchoContentHeader
	TargetUid    uint `param:"uid"`
	UnderMonitor bool
}

type ResponseUpdateControllerUnderMonitor bool

type RequestUpdateControllerUnderSolo struct {
	JwtHeader
	EchoContentHeader
	TargetUid uint `param:"uid"`
	Solo      bool
	EndTime   time.Time `json:"end_time"`
}

type ResponseUpdateControllerUnderSolo bool

type RequestUpdateControllerGuest struct {
	JwtHeader
	EchoContentHeader
	TargetUid uint `param:"uid"`
	Guest     bool
	Rating    int `json:"rating"`
}

type ResponseUpdateControllerGuest bool

type RequestAddControllerRecord struct {
	JwtHeader
	EchoContentHeader
	TargetCid int    `param:"cid"`
	Type      int    `json:"type"`
	Content   string `json:"content"`
}

type ResponseAddControllerRecord bool

type RequestDeleteControllerRecord struct {
	JwtHeader
	EchoContentHeader
	TargetRecord uint `param:"rid"`
}

type ResponseDeleteControllerRecord bool
