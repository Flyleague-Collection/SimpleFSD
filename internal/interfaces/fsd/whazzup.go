// Package fsd
package fsd

import "github.com/half-nothing/simple-fsd/internal/interfaces/operation"

type OnlineGeneral struct {
	Version          int    `json:"version"`
	GenerateTime     string `json:"generate_time"`
	ConnectedClients int    `json:"connected_clients"`
	OnlinePilot      int    `json:"online_pilot"`
	OnlineController int    `json:"online_controller"`
}

type OnlinePilot struct {
	Cid         string                `json:"cid"`
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
	Cid       string   `json:"cid"`
	Callsign  string   `json:"callsign"`
	RealName  string   `json:"real_name"`
	Latitude  float64  `json:"latitude"`
	Longitude float64  `json:"longitude"`
	Rating    int      `json:"rating"`
	Facility  int      `json:"facility"`
	Frequency int      `json:"frequency"`
	Range     int      `json:"range"`
	AtcInfo   []string `json:"atc_info"`
	LogonTime string   `json:"logon_time"`
}

type OnlineClients struct {
	General     *OnlineGeneral      `json:"general"`
	Pilots      []*OnlinePilot      `json:"pilots"`
	Controllers []*OnlineController `json:"controllers"`
}
