// Package service
package service

import "github.com/half-nothing/simple-fsd/src/interfaces/operation"

var (
	SuccessGetAuditLog          = NewApiStatus("GET_AUDIT_LOG", "成功获取审计日志", Ok)
	SuccessLogUnlawfulOverreach = NewApiStatus("LOG_UNLAWFUL_OVERREACH", "成功记录非法访问", Ok)
)

type AuditServiceInterface interface {
	GetAuditLogPage(req *RequestGetAuditLog) *ApiResponse[ResponseGetAuditLog]
	LogUnlawfulOverreach(req *RequestLogUnlawfulOverreach) *ApiResponse[ResponseLogUnlawfulOverreach]
}

type RequestGetAuditLog struct {
	JwtHeader
	Page     int `query:"page_number"`
	PageSize int `query:"page_size"`
}

type ResponseGetAuditLog struct {
	Items    []*operation.AuditLog `json:"items"`
	Page     int                   `json:"page"`
	PageSize int                   `json:"page_size"`
	Total    int64                 `json:"total"`
}

type RequestLogUnlawfulOverreach struct {
	JwtHeader
	EchoContentHeader
	AccessPath string `json:"access_path"`
}

type ResponseLogUnlawfulOverreach bool
