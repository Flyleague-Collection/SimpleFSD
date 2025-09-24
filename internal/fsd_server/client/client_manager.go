package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/half-nothing/simple-fsd/internal/interfaces/config"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"github.com/half-nothing/simple-fsd/internal/interfaces/queue"
	"github.com/half-nothing/simple-fsd/internal/utils"
	"sync"
	"sync/atomic"
	"time"
)

type ClientManager struct {
	logger          log.LoggerInterface
	clients         map[string]ClientInterface
	lock            sync.RWMutex
	shuttingDown    atomic.Bool
	config          *config.Config
	clientSlicePool sync.Pool
	whazzupContent  *utils.CachedValue[OnlineClients]
}

func NewClientManager(logger log.LoggerInterface, config *config.Config) *ClientManager {
	clientManager := &ClientManager{
		logger:       logger,
		clients:      make(map[string]ClientInterface),
		shuttingDown: atomic.Bool{},
		config:       config,
		clientSlicePool: sync.Pool{
			New: func() interface{} {
				return make([]ClientInterface, 0, 128)
			},
		},
	}
	clientManager.whazzupContent = utils.NewCachedValue[OnlineClients](config.Server.FSDServer.CacheDuration, func() *OnlineClients { return clientManager.getWhazzupContent() })
	return clientManager
}

func (cm *ClientManager) sendRawMessageTo(from int, to string, message string) error {
	client, exists := cm.GetClient(to)
	if !exists {
		return ErrCallsignNotFound
	}

	packet := MakePacket(Message, fmt.Sprintf("(%04d)", from), to, message)

	client.SendLine(packet)
	return nil
}

func (cm *ClientManager) HandleSendMessageToClientMessage(message *queue.Message) error {
	if val, ok := message.Data.(*SendRawMessageData); ok {
		return cm.sendRawMessageTo(val.From, val.To, val.Message)
	}
	return queue.ErrMessageDataType
}

func (cm *ClientManager) broadcastMessage(message *BroadcastMessageData) error {
	if cm.shuttingDown.Load() {
		return errors.New("client manager shutting down")
	}

	if len(message.Message) == 0 {
		return errors.New("message is empty")
	}

	clients := cm.GetClientSnapshot()
	defer cm.putSlice(clients)

	if len(clients) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, cm.config.Server.FSDServer.MaxBroadcastWorkers)

	var filter ClientFilter

	switch message.Target {
	case AllSup:
		filter = BroadcastToSupClient
	case AllATC:
		message.Target = AllClient
		filter = BroadcastToATCClient
	case AllPilot:
		message.Target = AllClient
		filter = BroadcastToAllPilotClient
	case AllClient:
		filter = BroadcastToAllClient
	}

	packet := MakePacket(Message, fmt.Sprintf("(%04d)", message.From), message.Target.String(), message.Message)

	for _, client := range clients {
		if client.Disconnected() {
			continue
		}

		if filter != nil && !filter(client) {
			continue
		}

		wg.Add(1)
		sem <- struct{}{}
		go func(cl ClientInterface) {
			defer func() {
				<-sem
				wg.Done()
			}()

			cm.logger.DebugF("[Broadcast] -> [%s] %s", cl.Callsign(), string(bytes.TrimRight(packet, "\r\n")))
			err := cl.SendLineWithoutLog(packet)
			if err != nil && errors.Is(err, ErrClientSocketWrite) {
				cl.MarkedDisconnect(false)
			}
		}(client)
	}

	wg.Wait()
	return nil
}

func (cm *ClientManager) HandleBroadcastMessage(message *queue.Message) error {
	if val, ok := message.Data.(*BroadcastMessageData); ok {
		return cm.broadcastMessage(val)
	}
	return queue.ErrMessageDataType
}

func (cm *ClientManager) KickClientFromServer(callsign string, reason string) (ClientInterface, error) {
	client, exists := cm.GetClient(callsign)
	if !exists {
		return nil, ErrCallsignNotFound
	}
	client.SendError(ResultError(Custom, true, callsign, fmt.Errorf("you were kicked from the server, reason is %s", reason)))
	return client, nil
}

func (cm *ClientManager) HandleKickClientFromServerMessage(message *queue.Message) error {
	if val, ok := message.Data.(*KickClientData); ok {
		_, err := cm.KickClientFromServer(val.Callsign, val.Reason)
		return err
	}
	return queue.ErrMessageDataType
}

func (cm *ClientManager) GetWhazzupContent() *OnlineClients {
	return cm.whazzupContent.GetValue()
}

func (cm *ClientManager) getWhazzupContent() *OnlineClients {
	data := &OnlineClients{
		General: OnlineGeneral{
			Version:          3,
			ConnectedClients: 0,
			OnlinePilot:      0,
			OnlineController: 0,
		},
		Pilots:      make([]*OnlinePilot, 0),
		Controllers: make([]*OnlineController, 0),
	}

	clientCopy := cm.GetClientSnapshot()
	defer cm.putSlice(clientCopy)

	for _, client := range clientCopy {
		if client == nil || client.Disconnected() {
			continue
		}
		data.General.ConnectedClients++
		if client.IsAtc() {
			data.General.OnlineController++
			controller := &OnlineController{
				Cid:         client.User().Cid,
				Callsign:    client.Callsign(),
				RealName:    client.RealName(),
				Latitude:    client.Position()[0].Latitude,
				Longitude:   client.Position()[0].Longitude,
				Rating:      client.Rating().Index(),
				Facility:    client.Facility().Index(),
				Frequency:   client.Frequency() + 100000,
				Range:       int(client.VisualRange()),
				OfflineTime: client.LogoffTime(),
				IsBreak:     client.IsBreak(),
				AtcInfo:     client.AtisInfo(),
				LogonTime:   client.History().StartTime.Format(time.DateTime),
			}
			data.Controllers = append(data.Controllers, controller)
		} else {
			data.General.OnlinePilot++
			pilot := &OnlinePilot{
				Cid:         client.User().Cid,
				Callsign:    client.Callsign(),
				RealName:    client.RealName(),
				Latitude:    client.Position()[0].Latitude,
				Longitude:   client.Position()[0].Longitude,
				Transponder: client.Transponder(),
				Heading:     client.Heading(),
				Altitude:    client.Altitude(),
				GroundSpeed: client.GroundSpeed(),
				FlightPlan:  client.FlightPlan(),
				LogonTime:   client.History().StartTime.Format(time.DateTime),
			}
			data.Pilots = append(data.Pilots, pilot)
		}
	}

	data.General.GenerateTime = time.Now().Format(time.DateTime)

	return data
}

func (cm *ClientManager) putSlice(clients []ClientInterface) {
	cm.clientSlicePool.Put(clients)
}

func (cm *ClientManager) Shutdown(ctx context.Context) error {
	if !cm.shuttingDown.CompareAndSwap(false, true) {
		return fmt.Errorf("shutting down already in progress")
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	clients := cm.GetClientSnapshot()
	defer cm.putSlice(clients)

	done := make(chan struct{})
	go func() {
		defer close(done)
		cm.disconnectClients(clients)
	}()

	select {
	case <-done:
		return nil
	case <-timeoutCtx.Done():
		return timeoutCtx.Err()
	}
}

func (cm *ClientManager) GetClientSnapshot() []ClientInterface {
	cm.lock.RLock()
	defer cm.lock.RUnlock()

	// 从池中获取切片
	clients := cm.clientSlicePool.Get().([]ClientInterface)
	clients = clients[:0]

	// 填充客户端
	for _, client := range cm.clients {
		clients = append(clients, client)
	}
	return clients
}

// 并发断开所有客户端连接
func (cm *ClientManager) disconnectClients(clients []ClientInterface) {
	if len(clients) == 0 {
		return
	}

	sem := make(chan struct{}, cm.config.Server.FSDServer.MaxBroadcastWorkers)
	var wg sync.WaitGroup

	for _, client := range clients {
		wg.Add(1)
		sem <- struct{}{}

		go func(c ClientInterface) {
			defer func() {
				<-sem
				wg.Done()
			}()

			c.MarkedDisconnect(true)
		}(client)
	}

	wg.Wait()
}

func (cm *ClientManager) AddClient(client ClientInterface) error {
	if cm.shuttingDown.Load() {
		return fmt.Errorf("server shutting down")
	}
	cm.lock.Lock()
	defer cm.lock.Unlock()

	if _, exists := cm.clients[client.Callsign()]; exists {
		return fmt.Errorf("client already registered: %s", client.Callsign())
	}
	cm.clients[client.Callsign()] = client
	return nil
}

func (cm *ClientManager) GetClient(callsign string) (ClientInterface, bool) {
	if cm.shuttingDown.Load() {
		return nil, false
	}

	cm.lock.RLock()
	defer cm.lock.RUnlock()

	client, exists := cm.clients[callsign]
	return client, exists
}

func (cm *ClientManager) DeleteClient(callsign string) bool {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	if _, exists := cm.clients[callsign]; !exists {
		return false
	}

	delete(cm.clients, callsign)
	return true
}

func (cm *ClientManager) SendMessageTo(callsign string, message []byte) error {
	if cm.shuttingDown.Load() {
		return fmt.Errorf("server is shutting down")
	}

	client, exists := cm.GetClient(callsign)
	if !exists {
		return ErrCallsignNotFound
	}

	client.SendLine(message)
	return nil
}

func (cm *ClientManager) BroadcastMessage(message []byte, fromClient ClientInterface, filter BroadcastFilter) {
	if cm.shuttingDown.Load() || len(message) == 0 {
		return
	}

	clients := cm.GetClientSnapshot()
	defer cm.putSlice(clients) // 重置并放回池中

	if len(clients) == 0 {
		return
	}

	fullMsg := make([]byte, len(message), len(message)+len(SplitSign))
	copy(fullMsg, message)
	fullMsg = append(fullMsg, SplitSign...)

	// 并发广播
	var wg sync.WaitGroup
	sem := make(chan struct{}, cm.config.Server.FSDServer.MaxBroadcastWorkers)

	for _, client := range clients {
		if client == fromClient || client.Disconnected() {
			continue
		}

		if filter != nil && !filter(client, fromClient) {
			continue
		}

		wg.Add(1)
		sem <- struct{}{}
		go func(cl ClientInterface) {
			defer func() {
				<-sem
				wg.Done()
			}()

			cm.logger.DebugF("[Broadcast] -> [%s] %s", cl.Callsign(), message)
			err := cl.SendLineWithoutLog(fullMsg)
			if err != nil && errors.Is(err, ErrClientSocketWrite) {
				cl.MarkedDisconnect(false)
			}
		}(client)
	}

	wg.Wait()
}
