//go:build !windows

package voice_server

import (
	"net"

	"github.com/half-nothing/simple-fsd/src/interfaces/voice"
	"golang.org/x/net/ipv4"
)

func (s *VoiceServer) broadcastToTargets(targets []*net.UDPAddr, rawData []byte, client *voice.ClientInfo) {
	packetConn := ipv4.NewPacketConn(s.udpConn)
	messages := make([]ipv4.Message, len(targets))

	for i, addr := range targets {
		messages[i] = ipv4.Message{
			Buffers: [][]byte{rawData},
			Addr:    addr,
		}
	}

	n, err := packetConn.WriteBatch(messages, 0)
	if err != nil {
		client.Logger.ErrorF("Failed to batch send voice data: %v", err)
	} else if n < len(messages) {
		client.Logger.WarnF("Partial batch send: %d/%d", n, len(messages))
	}
}
