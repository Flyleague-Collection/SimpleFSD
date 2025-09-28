// Package service
package service

import (
	"github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
)

var (
	ErrSendMessage              = NewApiStatus("FAIL_SEND_MESSAGE", "发送消息失败", ServerInternalError)
	ErrClientNotFound           = NewApiStatus("CLIENT_NOT_FOUND", "指定客户端不存在", NotFound)
	SuccessSendMessage          = NewApiStatus("SEND_MESSAGE", "发送成功", Ok)
	SuccessKillClient           = NewApiStatus("KILL_CLIENT", "成功踢出客户端", Ok)
	SuccessGetClientPath        = NewApiStatus("GET_CLIENT_PATH", "获取客户端飞行路径", Ok)
	SuccessSendBroadcastMessage = NewApiStatus("SEND_BROADCAST_MESSAGE", "获取客户端飞行路径", Ok)
)

type ClientServiceInterface interface {
	GetOnlineClients() *fsd.OnlineClients
	SendMessageToClient(req *RequestSendMessageToClient) *ApiResponse[ResponseSendMessageToClient]
	KillClient(req *RequestKillClient) *ApiResponse[ResponseKillClient]
	GetClientFlightPath(req *RequestClientPath) *ApiResponse[ResponseClientPath]
	SendBroadcastMessage(req *RequestSendBroadcastMessage) *ApiResponse[ResponseSendBroadcastMessage]
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

type RequestSendBroadcastMessage struct {
	JwtHeader
	EchoContentHeader
	Target  string `json:"target"`
	Message string `json:"message"`
}

type ResponseSendBroadcastMessage bool
