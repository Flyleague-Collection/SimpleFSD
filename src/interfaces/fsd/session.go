// Package fsd
package fsd

import (
	"net"

	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
)

type SessionInterface interface {
	Callsign() string
	SetCallsign(callsign string)
	User() *entity.User
	SetUser(user *entity.User)
	ConnId() string
	Conn() net.Conn
	SetDisconnected(disconnect bool)
	Client() ClientInterface
	SetClient(client ClientInterface)
	FacilityIdent() Facility
	SetFacilityIdent(facility Facility)
}
