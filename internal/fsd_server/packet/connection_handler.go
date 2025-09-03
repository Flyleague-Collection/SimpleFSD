package packet

import (
	"bufio"
	"fmt"
	c "github.com/half-nothing/simple-fsd/internal/config"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"net"
	"sync/atomic"
	"time"
)

var (
	splitSign    = []byte("\r\n")
	splitSignLen = len(splitSign)
)

type ConnectionHandler struct {
	conn                net.Conn
	connId              string
	callsign            string
	client              ClientInterface
	clientManager       ClientManagerInterface
	user                *operation.User
	disconnected        atomic.Bool
	config              *c.Config
	userOperation       operation.UserOperationInterface
	flightPlanOperation operation.FlightPlanOperationInterface
}

func NewConnectionHandler(
	conn net.Conn,
	config *c.Config,
	cm ClientManagerInterface,
	userOperation operation.UserOperationInterface,
	flightPlanOperation operation.FlightPlanOperationInterface,
) *ConnectionHandler {
	return &ConnectionHandler{
		conn:                conn,
		connId:              conn.RemoteAddr().String(),
		callsign:            "unknown",
		client:              nil,
		clientManager:       cm,
		user:                nil,
		disconnected:        atomic.Bool{},
		config:              config,
		userOperation:       userOperation,
		flightPlanOperation: flightPlanOperation,
	}
}

func (ch *ConnectionHandler) SendError(result *Result) {
	if result.Success {
		return
	}
	if ch.client != nil {
		ch.client.SendError(result)
		return
	}
	packet := makePacket(Error, "SERVER", ch.callsign, fmt.Sprintf("%03d", result.Errno.Index()), result.Env, result.Errno.String())
	c.DebugF("[%s](%s) <- %s", ch.connId, ch.callsign, packet[:len(packet)-splitSignLen])
	_, _ = ch.conn.Write(packet)
	if result.Fatal {
		ch.disconnected.Store(true)
		time.AfterFunc(500*time.Millisecond, func() {
			_ = ch.conn.Close()
		})
	}
}

func (ch *ConnectionHandler) handleLine(line []byte) {
	if ch.disconnected.Load() {
		return
	}
	command, data := parserCommandLine(line)
	result := ch.handleCommand(command, data, line)
	if result == nil {
		c.WarnF("[%s](%s) handleCommand return a nil result", ch.connId, ch.callsign)
		return
	}
	if !result.Success {
		c.ErrorF("[%s](%s) handleCommand fail, %s, %s", ch.connId, ch.callsign, result.Errno.String(), result.Err.Error())
		ch.SendError(result)
	}
}

func (ch *ConnectionHandler) HandleConnection() {
	defer func() {
		c.DebugF("[%s](%s) x Connection closed", ch.connId, ch.callsign)
		if err := ch.conn.Close(); err != nil && !isNetClosedError(err) {
			c.WarnF("[%s](%s) Error occurred while closing connection, details: %v", ch.connId, ch.callsign, err)
		}
	}()
	scanner := bufio.NewScanner(ch.conn)
	scanner.Split(createSplitFunc(splitSign))
	for scanner.Scan() {
		line := scanner.Bytes()
		c.DebugF("[%s](%s) -> %s", ch.connId, ch.callsign, line)
		ch.handleLine(line)
		if ch.disconnected.Load() {
			break
		}
	}

	if ch.client != nil {
		if ch.client.IsAtc() {
			ch.clientManager.BroadcastMessage(makePacket(RemoveAtc, ch.client.Callsign(), "SERVER"), ch.client, BroadcastToClientInRange)
		} else {
			ch.clientManager.BroadcastMessage(makePacket(RemovePilot, ch.client.Callsign(), "SERVER"), ch.client, BroadcastToClientInRange)
		}
		ch.client.MarkedDisconnect(false)
	}
}

func (ch *ConnectionHandler) Callsign() string { return ch.callsign }

func (ch *ConnectionHandler) SetCallsign(callsign string) { ch.callsign = callsign }

func (ch *ConnectionHandler) User() *operation.User { return ch.user }

func (ch *ConnectionHandler) SetUser(user *operation.User) { ch.user = user }

func (ch *ConnectionHandler) ConnId() string { return ch.connId }

func (ch *ConnectionHandler) Conn() net.Conn { return ch.conn }

func (ch *ConnectionHandler) SetDisconnected(disconnect bool) { ch.disconnected.Store(disconnect) }
