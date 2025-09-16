// Package service
package service

import (
	"github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
)

type ClientServiceInterface interface {
	GetOnlineClient() *fsd.OnlineClients
	SendMessageToClient(req *RequestSendMessageToClient) *ApiResponse[ResponseSendMessageToClient]
	KillClient(req *RequestKillClient) *ApiResponse[ResponseKillClient]
	GetClientPath(req *RequestClientPath) *ApiResponse[ResponseClientPath]
}

type RequestSendMessageToClient struct {
	JwtHeader
	EchoContentHeader
	Cid     int
	SendTo  string `param:"callsign"`
	Message string `json:"message"`
}

type ResponseSendMessageToClient bool

type RequestKillClient struct {
	JwtHeader
	EchoContentHeader
	Cid            int
	TargetCallsign string `param:"callsign"`
	Reason         string `json:"reason"`
}

type ResponseKillClient bool

type RequestClientPath struct {
	Callsign string `param:"callsign"`
}

type ResponseClientPath []*fsd.PilotPath
