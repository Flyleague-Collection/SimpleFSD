// Package voice_server
package voice_server

import (
	"context"
	"time"
)

type ShutdownCallback struct {
	server *VoiceServer
}

func NewShutdownCallback(server *VoiceServer) *ShutdownCallback {
	return &ShutdownCallback{
		server: server,
	}
}

func (callback *ShutdownCallback) Invoke(ctx context.Context) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		callback.server.Stop()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-timeoutCtx.Done():
		return timeoutCtx.Err()
	}
}
