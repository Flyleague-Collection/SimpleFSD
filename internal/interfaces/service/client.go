// Package service
package service

import (
	"github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
)

type ClientServiceInterface interface {
	GetOnlineClient() *OnlineClients
	SendMessageToClient(req *RequestSendMessageToClient) *ApiResponse[ResponseSendMessageToClient]
	KillClient(req *RequestKillClient) *ApiResponse[ResponseKillClient]
	GetClientPath(req *RequestClientPath) *ApiResponse[ResponseClientPath]
}

type OnlineGeneral struct {
	Version          int    `json:"version"`
	GenerateTime     string `json:"generate_time"`
	ConnectedClients int    `json:"connected_clients"`
	OnlinePilot      int    `json:"online_pilot"`
	OnlineController int    `json:"online_controller"`
}

type OnlinePilot struct {
	Cid         int                   `json:"cid"`
	Callsign    string                `json:"callsign"`
	RealName    string                `json:"real_name"`
	Latitude    float64               `json:"latitude"`
	Longitude   float64               `json:"longitude"`
	Transponder string                `json:"transponder"`
	Heading     int                   `json:"heading"`
	Altitude    int                   `json:"altitude"`
	GroundSpeed int                   `json:"ground_speed"`
	FlightPlan  *operation.FlightPlan `json:"flight_plan"`
	LogonTime   string                `json:"logon_time"`
}

type OnlineController struct {
	Cid         int      `json:"cid"`
	Callsign    string   `json:"callsign"`
	RealName    string   `json:"real_name"`
	Latitude    float64  `json:"latitude"`
	Longitude   float64  `json:"longitude"`
	Rating      int      `json:"rating"`
	Facility    int      `json:"facility"`
	Frequency   int      `json:"frequency"`
	Range       int      `json:"range"`
	OfflineTime string   `json:"offline_time"`
	IsBreak     bool     `json:"is_break"`
	AtcInfo     []string `json:"atc_info"`
	LogonTime   string   `json:"logon_time"`
}

type OnlineClients struct {
	General     OnlineGeneral       `json:"general"`
	Pilots      []*OnlinePilot      `json:"pilots"`
	Controllers []*OnlineController `json:"controllers"`
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
	Callsign string `query:"callsign"`
}

type ResponseClientPath []*fsd.PilotPath
