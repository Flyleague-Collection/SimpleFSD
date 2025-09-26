// Package service
package service

import "github.com/half-nothing/simple-fsd/internal/interfaces/operation"

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
