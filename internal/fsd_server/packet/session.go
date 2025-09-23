package packet

import (
	"bufio"
	"fmt"
	"github.com/half-nothing/simple-fsd/internal/interfaces"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/internal/interfaces/global"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"net"
	"sync/atomic"
	"time"
)

var (
	splitSign    = []byte("\r\n")
	splitSignLen = len(splitSign)
)

type Session struct {
	logger              log.LoggerInterface
	conn                net.Conn
	connId              string
	callsign            string
	facilityIdent       Facility
	user                *operation.User
	disconnected        atomic.Bool
	application         *interfaces.ApplicationContent
	client              ClientInterface
	clientManager       ClientManagerInterface
	metarManager        interfaces.MetarManagerInterface
	refuseOutRange      bool
	isSimulatorServer   bool
	userOperation       operation.UserOperationInterface
	flightPlanOperation operation.FlightPlanOperationInterface
}

func NewSession(
	application *interfaces.ApplicationContent,
	conn net.Conn,
) *Session {
	config := application.ConfigManager().Config()
	return &Session{
		logger:              application.Logger().FsdLogger(),
		application:         application,
		conn:                conn,
		connId:              conn.RemoteAddr().String(),
		callsign:            "unknown",
		client:              nil,
		clientManager:       application.ClientManager(),
		user:                nil,
		disconnected:        atomic.Bool{},
		metarManager:        application.MetarManager(),
		refuseOutRange:      config.Server.FSDServer.RangeLimit.RefuseOutRange,
		isSimulatorServer:   config.Server.General.SimulatorServer,
		userOperation:       application.Operations().UserOperation(),
		flightPlanOperation: application.Operations().FlightPlanOperation(),
	}
}

func (session *Session) SendError(result *Result) {
	if result.Success {
		return
	}
	if session.client != nil {
		session.client.SendError(result)
		return
	}

	var errString string
	if result.Errno == Custom {
		errString = result.Err.Error()
	} else {
		errString = result.Errno.String()
	}

	packet := makePacket(Error, global.FSDServerName, session.callsign, fmt.Sprintf("%03d", result.Errno.Index()), result.Env, errString)
	session.logger.DebugF("[%s](%s) <- %s", session.connId, session.callsign, packet[:len(packet)-splitSignLen])
	_, _ = session.conn.Write(packet)
	if result.Fatal {
		session.disconnected.Store(true)
	}
}

func (session *Session) handleLine(line []byte) {
	if session.disconnected.Load() {
		return
	}
	command, data := parserCommandLine(line)
	result := session.handleCommand(command, data, line)
	if result == nil {
		session.logger.WarnF("[%s](%s) handleCommand return a nil result", session.connId, session.callsign)
		return
	}
	if !result.Success {
		session.logger.ErrorF("[%s](%s) handleCommand fail, %s, %s", session.connId, session.callsign, result.Errno.String(), result.Err.Error())
		session.SendError(result)
	}
}

func (session *Session) HandleConnection(timeout time.Duration) {
	defer func() {
		time.AfterFunc(global.FSDDisconnectDelay, func() {
			session.logger.DebugF("[%s](%s) x Connection closed", session.connId, session.callsign)
			if err := session.conn.Close(); err != nil && !isNetClosedError(err) {
				session.logger.WarnF("[%s](%s) Error occurred while closing connection, details: %v", session.connId, session.callsign, err)
			}
		})
	}()
	scanner := bufio.NewScanner(session.conn)
	scanner.Split(createSplitFunc(splitSign))
	_ = session.conn.SetDeadline(time.Now().Add(timeout))
	for scanner.Scan() {
		_ = session.conn.SetDeadline(time.Now().Add(timeout))
		if scanner.Err() != nil {
			session.logger.ErrorF("error while scanning, %v", scanner.Err())
			break
		}
		line := scanner.Bytes()
		session.logger.DebugF("[%s](%s) -> %s", session.connId, session.callsign, line)
		if session.client == nil {
			session.handleLine(line)
		} else {
			go session.handleLine(line)
		}
		if session.disconnected.Load() {
			break
		}
	}

	if session.client != nil {
		if session.client.IsAtc() {
			session.clientManager.BroadcastMessage(makePacket(RemoveAtc, session.client.Callsign(), global.FSDServerName), session.client, BroadcastToClientInRange)
		} else {
			session.clientManager.BroadcastMessage(makePacket(RemovePilot, session.client.Callsign(), global.FSDServerName), session.client, BroadcastToClientInRange)
		}
		session.client.MarkedDisconnect(false)
	}
}

func (session *Session) Callsign() string { return session.callsign }

func (session *Session) SetCallsign(callsign string) { session.callsign = callsign }

func (session *Session) User() *operation.User { return session.user }

func (session *Session) SetUser(user *operation.User) { session.user = user }

func (session *Session) ConnId() string { return session.connId }

func (session *Session) Conn() net.Conn { return session.conn }

func (session *Session) SetDisconnected(disconnect bool) { session.disconnected.Store(disconnect) }
