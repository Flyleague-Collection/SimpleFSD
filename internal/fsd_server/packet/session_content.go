// Package packet
package packet

import (
	"bufio"
	"errors"
	"fmt"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/internal/interfaces/global"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"time"
)

type SessionContent struct {
	logger           log.LoggerInterface
	commandHandler   CommandHandlerInterface
	clientManager    ClientManagerInterface
	heartbeatTimeout time.Duration
	possibleCommands [][]byte
}

func NewSessionContent(
	logger log.LoggerInterface,
	commandHandler CommandHandlerInterface,
	clientManager ClientManagerInterface,
	heartbeatTimeout time.Duration,
) *SessionContent {
	content := &SessionContent{
		logger:           log.NewLoggerAdapter(logger, "SessionManager"),
		commandHandler:   commandHandler,
		clientManager:    clientManager,
		heartbeatTimeout: heartbeatTimeout,
	}
	content.possibleCommands = commandHandler.GetPossibleCommands()
	return content
}

func (content *SessionContent) SendError(session *Session, result *Result) {
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

	packet := MakePacket(Error, global.FSDServerName, session.callsign, fmt.Sprintf("%03d", result.Errno.Index()), result.Env, errString)
	content.logger.DebugF("[%s](%s) <- %s", session.connId, session.callsign, packet[:len(packet)-SplitSignLen])
	_, _ = session.conn.Write(packet)
	if result.Fatal {
		session.close.Store(true)
	}
}

func (content *SessionContent) handleCommand(session *Session, commandType ClientCommand, data []string, rawLine []byte) *Result {
	res := content.commandHandler.Call(commandType, session, data, rawLine)
	if res == nil {
		return ResultError(Syntax, false, string(commandType), errors.New("parse command failed"))
	}
	return res
}

func (content *SessionContent) handleLine(session *Session, line []byte) {
	if session.close.Load() {
		return
	}
	command, data := parserCommandLine(line, content.possibleCommands)
	result := content.handleCommand(session, command, data, line)
	if result == nil {
		content.logger.WarnF("[%s](%s) handleCommand return a nil result", session.connId, session.callsign)
		return
	}
	if !result.Success {
		content.logger.ErrorF("[%s](%s) handleCommand fail, %s, %s", session.connId, session.callsign, result.Errno.String(), result.Err.Error())
		content.SendError(session, result)
	}
}

func (content *SessionContent) HandleConnection(session *Session) {
	defer func() {
		time.AfterFunc(global.FSDDisconnectDelay, func() {
			content.logger.DebugF("[%s](%s) x Connection closed", session.connId, session.callsign)
			if err := session.conn.Close(); err != nil && !isNetClosedError(err) {
				content.logger.WarnF("[%s](%s) Error occurred while closing connection, details: %v", session.connId, session.callsign, err)
			}
		})
	}()
	scanner := bufio.NewScanner(session.conn)
	scanner.Split(createSplitFunc(SplitSign))
	_ = session.conn.SetDeadline(time.Now().Add(content.heartbeatTimeout))
	for scanner.Scan() {
		_ = session.conn.SetDeadline(time.Now().Add(content.heartbeatTimeout))
		if scanner.Err() != nil {
			content.logger.ErrorF("error while scanning, %v", scanner.Err())
			break
		}
		line := scanner.Bytes()
		content.logger.DebugF("[%s](%s) -> %s", session.connId, session.callsign, line)
		if session.client == nil {
			content.handleLine(session, line)
		} else {
			go content.handleLine(session, line)
		}
		if session.close.Load() {
			break
		}
	}

	if session.client != nil {
		if session.client.IsAtc() {
			content.clientManager.BroadcastMessage(MakePacketWithoutSign(RemoveAtc, session.client.Callsign(), global.FSDServerName), session.client, BroadcastToClientInRange)
		} else {
			content.clientManager.BroadcastMessage(MakePacketWithoutSign(RemovePilot, session.client.Callsign(), global.FSDServerName), session.client, BroadcastToClientInRange)
		}
		session.client.MarkedDisconnect(false)
	}
}
