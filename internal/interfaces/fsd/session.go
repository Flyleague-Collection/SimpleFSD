// Package fsd
package fsd

import (
	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"net"
)

type SessionInterface interface {
	SendError(result *Result)
	HandleConnection()
	Callsign() string
	SetCallsign(callsign string)
	User() *operation.User
	SetUser(user *operation.User)
	ConnId() string
	Conn() net.Conn
	SetDisconnected(disconnect bool)
}
