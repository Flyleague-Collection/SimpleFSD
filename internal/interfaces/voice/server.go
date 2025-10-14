// Package voice
package voice

type VoicePacket struct {
	Cid         int
	Transmitter int
	Frequency   int
	Callsign    string
	Data        []byte
}

type VoiceServerInterface interface {
	Start() error
	Stop()
}
