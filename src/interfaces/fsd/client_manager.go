// Package fsd
package fsd

import (
	"context"
	"errors"

	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
	"github.com/half-nothing/simple-fsd/src/interfaces/queue"
)

var (
	ErrCallsignNotFound = errors.New("callsign not found")
	ErrCidMissMatch     = errors.New("cid miss match")
)

type ClientManagerInterface interface {
	GetWhazzupContent() *OnlineClients
	Shutdown(ctx context.Context) error
	GetClientSnapshot() []ClientInterface
	AddClient(client ClientInterface) error
	GetClient(callsign string) (ClientInterface, bool)
	DeleteClient(callsign string) bool
	HandleKickClientFromServerMessage(message *queue.Message) error
	HandleSendMessageToClientMessage(message *queue.Message) error
	HandleBroadcastMessage(message *queue.Message) error
	KickClientFromServer(callsign string, reason string) (ClientInterface, error)
	SendMessageTo(callsign string, message []byte) error
	BroadcastMessage(message []byte, fromClient ClientInterface, filter BroadcastFilter)
}

type BroadcastMessageData struct {
	From    string
	Target  BroadcastTarget
	Message string
}

type LockChange struct {
	TargetCallsign string
	TargetCid      int
	Locked         bool
}

type FlushFlightPlan struct {
	TargetCallsign string
	TargetCid      int
	FlightPlan     *entity.FlightPlan
}

type SendRawMessageData struct {
	From    string
	To      string
	Message string
}

type KickClientData struct {
	Callsign string
	Reason   string
}

type OnlineGeneral struct {
	Version          int    `json:"version"`
	GenerateTime     string `json:"generate_time"`
	ConnectedClients int    `json:"connected_clients"`
	OnlinePilot      int    `json:"online_pilot"`
	OnlineController int    `json:"online_controller"`
}

type OnlinePilot struct {
	Cid         int                `json:"cid"`
	Callsign    string             `json:"callsign"`
	RealName    string             `json:"real_name"`
	Latitude    float64            `json:"latitude"`
	Longitude   float64            `json:"longitude"`
	Transponder string             `json:"transponder"`
	Heading     int                `json:"heading"`
	Altitude    int                `json:"altitude"`
	GroundSpeed int                `json:"ground_speed"`
	FlightPlan  *entity.FlightPlan `json:"flight_plan"`
	LogonTime   string             `json:"logon_time"`
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
