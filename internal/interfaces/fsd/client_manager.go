// Package fsd
package fsd

import (
	"context"
	"errors"
)

var (
	ErrCallsignNotFound = errors.New("callsign not found")
)

type ClientManagerInterface interface {
	PutSlice(clients []ClientInterface)
	Shutdown(ctx context.Context) error
	GetClientSnapshot() []ClientInterface
	AddClient(client ClientInterface) error
	GetClient(callsign string) (ClientInterface, bool)
	DeleteClient(callsign string) bool
	SendMessageTo(callsign string, message []byte) error
	SendRawMessageTo(from int, to string, message string) error
	BroadcastMessage(message []byte, fromClient ClientInterface, filter BroadcastFilter)
	NewClient(callsign string, rating Rating, protocol int, realName string, socket SessionInterface, isAtc bool) ClientInterface
}
