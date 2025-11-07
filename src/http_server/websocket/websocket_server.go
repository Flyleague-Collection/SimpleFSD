// Package websocket
// File websocket_server.go
package websocket

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/half-nothing/simple-fsd/src/interfaces/config"
	"github.com/half-nothing/simple-fsd/src/interfaces/global"
	"github.com/half-nothing/simple-fsd/src/interfaces/http/service"
	"github.com/half-nothing/simple-fsd/src/interfaces/log"
	"github.com/half-nothing/simple-fsd/src/interfaces/queue"
	"github.com/labstack/echo/v4"
)

type WebSocketServer struct {
	logger       log.LoggerInterface
	clientMap    map[string]*FakeClient
	messageChan  chan *SendMessage
	messageQueue queue.MessageQueueInterface
	ctx          context.Context
	cancel       context.CancelFunc
	config       *config.HttpServerConfig
	lock         sync.RWMutex
	upgrader     *websocket.Upgrader
}

func NewWebSocketServer(
	logger log.LoggerInterface,
	messageQueue queue.MessageQueueInterface,
	config *config.HttpServerConfig,
) *WebSocketServer {
	server := &WebSocketServer{
		logger:       logger,
		clientMap:    make(map[string]*FakeClient),
		messageChan:  make(chan *SendMessage, 128),
		messageQueue: messageQueue,
		config:       config,
		lock:         sync.RWMutex{},
		upgrader:     &websocket.Upgrader{},
	}
	server.ctx, server.cancel = context.WithCancel(context.Background())
	messageQueue.Subscribe(queue.FsdMessageReceived, server.MessageReceiveHandler)
	return server
}

func (server *WebSocketServer) Close(ctx context.Context) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	server.cancel()
	done := make(chan struct{})
	go func() {
		server.logger.Info("Closing WebSocket Server")
		wg := sync.WaitGroup{}
		for _, client := range server.clientMap {
			wg.Add(1)
			go func(client *FakeClient) {
				defer wg.Done()
				client.Disconnect()
			}(client)
		}
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-timeoutCtx.Done():
		server.cancel()
		return timeoutCtx.Err()
	}
}

func (server *WebSocketServer) ConnectToFsd(c echo.Context) error {
	ws, err := server.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	server.newConnection(ws)
	return nil
}

func (server *WebSocketServer) RemoveClient(callsign string) {
	server.lock.Lock()
	defer server.lock.Unlock()
	delete(server.clientMap, callsign)
}

func (server *WebSocketServer) defaultKeyFunc(token *jwt.Token) (interface{}, error) {
	if token.Method.Alg() != global.SigningMethod {
		return nil, errors.New("illegal signature methods")
	}
	return []byte(server.config.JWT.Secret), nil
}

func (server *WebSocketServer) sendError(conn *websocket.Conn, closeCode int, text string) {
	msg := websocket.FormatCloseMessage(closeCode, text)
	_ = conn.WriteControl(websocket.CloseMessage, msg, time.Now().Add(5*time.Second))
	_ = conn.Close()
}

func (server *WebSocketServer) newConnection(conn *websocket.Conn) {
	_ = conn.SetReadDeadline(time.Now().Add(time.Minute))

	conn.SetCloseHandler(func(code int, text string) error {
		server.logger.InfoF("connection closed by client: %d %s", code, text)
		_ = conn.Close()
		return nil
	})

	_, msg, err := conn.ReadMessage()
	if err != nil {
		server.logger.ErrorF("error while reading message: %v", err)
		server.sendError(conn, websocket.CloseInternalServerErr, "Token format error")
		return
	}
	_ = conn.SetReadDeadline(time.Time{})

	claims, err := jwt.ParseWithClaims(string(msg), &service.Claims{}, server.defaultKeyFunc)
	if err != nil {
		server.logger.ErrorF("error while parsing claims: %v", err)
		server.sendError(conn, TokenFormatError, "Token format error")
		return
	}

	token, ok := claims.Claims.(*service.Claims)
	if !ok {
		server.logger.ErrorF("error while parsing claims: %v", claims)
		server.sendError(conn, TokenFormatError, "Token format error")
		return
	}

	if token.ExpiresAt.Before(time.Now()) {
		server.logger.ErrorF("error while parsing claims: %v", token)
		server.sendError(conn, TokenExpired, "Token expired")
		return
	}

	callsign := server.config.FormatCallsign(token.Cid)
	client := NewFakeClient(
		server.logger,
		callsign,
		conn,
		server.messageQueue,
		server.ctx,
		*global.WebsocketTimeout,
		*global.WebsocketHeartbeatInterval,
		*global.WebsocketMessageChannelSize,
		func(callsign string) {
			server.RemoveClient(callsign)
		})
	server.lock.Lock()
	c, ok := server.clientMap[callsign]
	if ok {
		c.Disconnect()
	}
	server.clientMap[callsign] = client
	server.lock.Unlock()
	client.SendMessage(&ReceiveMessage{From: "SERVER", Data: fmt.Sprintf("Your callsign: %s", callsign)})
}

func (server *WebSocketServer) SendMessageTo(message []byte) error {
	select {
	case <-server.ctx.Done():
		return server.ctx.Err()
	default:
		rawData := strings.Split(string(message), ":")
		if len(rawData) != 3 {
			return errors.New("message format error")
		}

		from := strings.TrimLeft(rawData[0], "#TM")
		to := rawData[1]
		msg := rawData[2]

		server.lock.RLock()
		client, ok := server.clientMap[to]
		server.lock.RUnlock()
		if !ok {
			return errors.New("client not exist")
		}

		client.SendMessage(&ReceiveMessage{From: from, Data: msg})
		return nil
	}
}

func (server *WebSocketServer) MessageReceiveHandler(message *queue.Message) error {
	select {
	case <-server.ctx.Done():
		return server.ctx.Err()
	default:
		if val, ok := message.Data.([]byte); ok {
			return server.SendMessageTo(val)
		}
		return queue.ErrMessageDataType
	}
}
