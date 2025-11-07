// Package websocket
// File fake_client.go
package websocket

import (
	"context"
	"encoding/json"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/half-nothing/simple-fsd/src/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/src/interfaces/log"
	"github.com/half-nothing/simple-fsd/src/interfaces/queue"
)

type FakeClient struct {
	logger          log.LoggerInterface
	conn            *websocket.Conn
	receiveChannel  chan *ReceiveMessage
	ctx             context.Context
	cancelFunc      context.CancelFunc
	messageQueue    queue.MessageQueueInterface
	timeoutDuration time.Duration
	heartbeatTicker *time.Ticker
	closeChan       chan bool
	callsign        string
	wg              sync.WaitGroup
	afterDisconnect func(callsign string)
	mu              sync.RWMutex
	disconnected    atomic.Bool
}

func NewFakeClient(
	logger log.LoggerInterface,
	callsign string,
	conn *websocket.Conn,
	messageQueue queue.MessageQueueInterface,
	ctx context.Context,
	timeoutDuration time.Duration,
	heartbeatDuration time.Duration,
	channelSize int,
	afterDisconnect func(callsign string),
) *FakeClient {
	client := &FakeClient{
		logger:          log.NewLoggerAdapter(logger, callsign),
		callsign:        callsign,
		conn:            conn,
		messageQueue:    messageQueue,
		timeoutDuration: timeoutDuration,
		heartbeatTicker: time.NewTicker(heartbeatDuration),
		receiveChannel:  make(chan *ReceiveMessage, channelSize),
		wg:              sync.WaitGroup{},
		closeChan:       make(chan bool),
		afterDisconnect: afterDisconnect,
		mu:              sync.RWMutex{},
		disconnected:    atomic.Bool{},
	}
	client.initializeConnection()
	client.ctx, client.cancelFunc = context.WithCancel(ctx)
	client.wg.Add(2)
	go client.receiveHandler()
	go client.sendHandler()
	return client
}

func (client *FakeClient) initializeConnection() {
	client.conn.SetPongHandler(func(appData string) error {
		return client.resetReadDeadline()
	})

	client.conn.SetCloseHandler(func(code int, text string) error {
		client.logger.InfoF("connection closed by client: %d %s", code, text)
		client.Disconnect()
		return nil
	})

	_ = client.resetReadDeadline()
}

func (client *FakeClient) safeDisconnect() {
	client.disconnected.Store(true)

	client.heartbeatTicker.Stop()
	client.cancelFunc()

	msg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Server disconnected, see you next time")
	_ = client.conn.WriteControl(websocket.CloseMessage, msg, time.Now().Add(3*time.Second))
	_ = client.conn.Close()

	client.wg.Wait()

	if client.afterDisconnect != nil {
		client.afterDisconnect(client.callsign)
	}
}

func (client *FakeClient) Disconnect() {
	if client.disconnected.Load() {
		return
	}
	client.safeDisconnect()
}

func (client *FakeClient) SendMessage(msg *ReceiveMessage) {
	if client.disconnected.Load() {
		return
	}

	select {
	case <-client.ctx.Done():
	case client.receiveChannel <- msg:
	default:
		client.logger.Warn("send channel full, dropping message")
	}
}

func (client *FakeClient) resetReadDeadline() error {
	return client.conn.SetReadDeadline(time.Now().Add(client.timeoutDuration))
}

func (client *FakeClient) receiveHandler() {
	defer client.wg.Done()
	for {
		select {
		case <-client.ctx.Done():
			return
		default:
			_, msg, err := client.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					client.logger.ErrorF("unexpected close error: %v", err)
				}
				client.Disconnect()
				return
			}

			if err := client.resetReadDeadline(); err != nil {
				client.logger.ErrorF("error setting read deadline: %v", err)
				client.Disconnect()
				return
			}

			data := &SendMessage{}
			if err := json.Unmarshal(msg, data); err != nil {
				client.logger.ErrorF("error while unmarshalling message: %v", err)
				continue
			}

			if fsd.IsValidBroadcastTarget(data.Target) {
				client.messageQueue.Publish(&queue.Message{
					Type: queue.BroadcastMessage,
					Data: &fsd.BroadcastMessageData{
						From:    client.callsign,
						Target:  fsd.BroadcastTarget(data.Target),
						Message: data.Data,
					},
				})
			} else {
				client.messageQueue.Publish(&queue.Message{
					Type: queue.SendMessageToClient,
					Data: &fsd.SendRawMessageData{
						From:    client.callsign,
						To:      data.Target,
						Message: data.Data,
					},
				})
			}
		}
	}
}

func (client *FakeClient) sendHandler() {
	defer client.wg.Done()
	for {
		select {
		case <-client.ctx.Done():
			return
		case msg := <-client.receiveChannel:
			data, err := json.Marshal(msg)
			if err != nil {
				client.logger.ErrorF("error while marshalling message: %v", err)
				continue
			}
			if err := client.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				client.logger.ErrorF("error while writing message: %v", err)
				client.Disconnect()
				return
			}
		case <-client.heartbeatTicker.C:
			if err := client.conn.WriteMessage(websocket.PingMessage, []byte("Ping")); err != nil {
				client.logger.ErrorF("error while sending ping: %v", err)
				client.Disconnect()
				return
			}
		}
	}
}
