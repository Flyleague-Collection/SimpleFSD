// Package voice
package voice

type MessageType string

var (
	Switch       MessageType = "channel"
	Ping         MessageType = "ping"
	Pong         MessageType = "pong"
	Error        MessageType = "error"
	TextReceive  MessageType = "text_receive"
	VoiceReceive MessageType = "voice_receive"
	Message      MessageType = "message"
	Disconnect   MessageType = "disconnect"
)

type ControlMessage struct {
	Type        MessageType `json:"type"`
	Cid         int         `json:"cid"`
	Callsign    string      `json:"callsign"`
	Transmitter int         `json:"transmitter"`
	Data        string      `json:"data"`
}
