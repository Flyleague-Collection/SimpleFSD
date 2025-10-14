//go:build windows

package voice_server

import (
	"net"
	"sync"

	"github.com/half-nothing/simple-fsd/internal/interfaces/voice"
)

func (s *VoiceServer) broadcastToTargets(targets []*net.UDPAddr, rawData []byte, client *voice.ClientInfo) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, s.config.BroadcastLimit)

	for _, addr := range targets {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(targetAddr *net.UDPAddr) {
			defer func() {
				<-semaphore
				wg.Done()
			}()

			_, err := s.udpConn.WriteToUDP(rawData, targetAddr)
			if err != nil {
				client.Logger.DebugF("Failed to send to %s: %v", targetAddr, err)
			}
		}(addr)
	}

	wg.Wait()
}
