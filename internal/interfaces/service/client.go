// Package service
package service

import (
	"github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
)

type ClientServiceInterface interface {
	GetOnlineClients() *fsd.OnlineClients
	SendMessageToClient(req *RequestSendMessageToClient) *ApiResponse[ResponseSendMessageToClient]
	KillClient(req *RequestKillClient) *ApiResponse[ResponseKillClient]
	GetClientFlightPath(req *RequestClientPath) *ApiResponse[ResponseClientPath]
}

type RequestSendMessageToClient struct {
	JwtHeader
	EchoContentHeader
	SendTo  string `param:"callsign"`
	Message string `json:"message"`
}

type ResponseSendMessageToClient bool

type RequestKillClient struct {
	JwtHeader
	EchoContentHeader
	TargetCallsign string `param:"callsign"`
	Reason         string `json:"reason"`
}

type ResponseKillClient bool

type RequestClientPath struct {
	Callsign string `param:"callsign"`
}

type ResponseClientPath []*fsd.PilotPath
