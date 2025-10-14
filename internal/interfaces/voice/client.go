// Package voice
package voice

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"sync/atomic"

	"github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
)

type Transmitter struct {
	Id         int
	ClientInfo *ClientInfo
	Frequency  ChannelFrequency
	UDPAddr    *net.UDPAddr
}

type ClientInfo struct {
	Cid               int
	Callsign          string
	Client            fsd.ClientInterface
	Logger            log.LoggerInterface
	TCPConn           net.Conn
	Decoder           *json.Decoder
	Encoder           *json.Encoder
	Disconnected      atomic.Bool
	ActiveTransmitter int
	TransmitterMutex  sync.Mutex
	Transmitters      []*Transmitter
}

func NewClientInfo(
	logger log.LoggerInterface,
	cid int,
	callsign string,
	conn net.Conn,
	client fsd.ClientInterface,
) *ClientInfo {
	return &ClientInfo{
		Cid:              cid,
		Callsign:         callsign,
		Client:           client,
		TCPConn:          conn,
		Logger:           logger,
		Decoder:          json.NewDecoder(conn),
		Encoder:          json.NewEncoder(conn),
		Disconnected:     atomic.Bool{},
		TransmitterMutex: sync.Mutex{},
		Transmitters:     make([]*Transmitter, 0),
	}
}

func (client *ClientInfo) MessageReceive(message []byte) {
	_ = client.SendMessage(Message, string(message))
}

func (client *ClientInfo) ConnectionDisconnect() {
	_ = client.SendError("fsd connection disconnected")
	_ = client.TCPConn.Close()
}

func (client *ClientInfo) SendError(msg string) error {
	return client.SendMessage(Error, msg)
}

func (client *ClientInfo) SendMessage(messageType MessageType, msg string) error {
	message := &ControlMessage{
		Type: messageType,
		Data: msg,
	}
	return client.SendControlMessage(message)
}

func (client *ClientInfo) SendControlMessage(msg *ControlMessage) error {
	msg.Cid = client.Cid
	msg.Callsign = client.Callsign
	if err := client.Encoder.Encode(msg); err != nil {
		return fmt.Errorf("failed to write control message: %v", err)
	}
	return nil
}
