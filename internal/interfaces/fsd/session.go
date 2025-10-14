// Package fsd
package fsd

import (
	"net"

	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
)

type SessionInterface interface {
	Callsign() string
	SetCallsign(callsign string)
	User() *operation.User
	SetUser(user *operation.User)
	ConnId() string
	Conn() net.Conn
	SetDisconnected(disconnect bool)
	Client() ClientInterface
	SetClient(client ClientInterface)
	FacilityIdent() Facility
	SetFacilityIdent(facility Facility)
}
