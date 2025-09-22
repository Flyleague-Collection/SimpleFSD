// Package fsd
package fsd

import (
	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"net"
	"time"
)

type SessionInterface interface {
	SendError(result *Result)
	HandleConnection(timeout time.Duration)
	Callsign() string
	SetCallsign(callsign string)
	User() *operation.User
	SetUser(user *operation.User)
	ConnId() string
	Conn() net.Conn
	SetDisconnected(disconnect bool)
}
