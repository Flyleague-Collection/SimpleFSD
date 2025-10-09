package packet

import (
	"net"
	"sync/atomic"

	. "github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
)

type Session struct {
	conn          net.Conn
	connId        string
	callsign      string
	facilityIdent Facility
	user          *operation.User
	close         atomic.Bool
	client        ClientInterface
}

func NewSession(conn net.Conn) *Session {
	return &Session{
		conn:     conn,
		connId:   conn.RemoteAddr().String(),
		callsign: "unknown",
		client:   nil,
		user:     nil,
		close:    atomic.Bool{},
	}
}

func (session *Session) Callsign() string { return session.callsign }

func (session *Session) SetCallsign(callsign string) { session.callsign = callsign }

func (session *Session) User() *operation.User { return session.user }

func (session *Session) SetUser(user *operation.User) { session.user = user }

func (session *Session) ConnId() string { return session.connId }

func (session *Session) Conn() net.Conn { return session.conn }

func (session *Session) SetDisconnected(disconnect bool) { session.close.Store(disconnect) }

func (session *Session) Client() ClientInterface { return session.client }

func (session *Session) SetClient(client ClientInterface) { session.client = client }

func (session *Session) FacilityIdent() Facility { return session.facilityIdent }

func (session *Session) SetFacilityIdent(facility Facility) { session.facilityIdent = facility }
