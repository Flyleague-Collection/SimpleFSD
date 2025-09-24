// Package service
package service

import (
	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"time"
)

type ControllerApplicationServiceInterface interface {
	GetSelfApplication(req *RequestGetSelfApplication) *ApiResponse[ResponseGetSelfApplication]
	GetApplications(req *RequestGetApplications) *ApiResponse[ResponseGetApplications]
	SubmitControllerApplication(req *RequestSubmitControllerApplication) *ApiResponse[ResponseSubmitControllerApplication]
	CancelSelfApplication(req *RequestCancelSelfApplication) *ApiResponse[ResponseCancelSelfApplication]
	UpdateApplicationStatus(req *RequestUpdateApplicationStatus) *ApiResponse[ResponseUpdateApplicationStatus]
}

type RequestGetSelfApplication struct {
	JwtHeader
}

type ResponseGetSelfApplication *operation.ControllerApplication

type RequestGetApplications struct {
	JwtHeader
	PageArguments
}

type ResponseGetApplications *PageResponse[*operation.ControllerApplication]

type RequestSubmitControllerApplication struct {
	JwtHeader
	EchoContentHeader
	*operation.ControllerApplication
}

type ResponseSubmitControllerApplication bool

type RequestCancelSelfApplication struct {
	JwtHeader
	EchoContentHeader
}

type ResponseCancelSelfApplication bool

type RequestUpdateApplicationStatus struct {
	JwtHeader
	EchoContentHeader
	ApplicationId uint        `param:"aid"`
	Status        int         `json:"status"`
	AvailableTime []time.Time `json:"times"`
	Message       string      `json:"message"`
}

type ResponseUpdateApplicationStatus bool
