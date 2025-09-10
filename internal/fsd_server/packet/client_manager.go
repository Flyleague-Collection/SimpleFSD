package packet

import (
	"context"
	"fmt"
	"github.com/half-nothing/simple-fsd/internal/interfaces"
	"github.com/half-nothing/simple-fsd/internal/interfaces/config"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/internal/interfaces/global"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type ClientManager struct {
	clients            map[string]ClientInterface
	lock               sync.RWMutex
	shuttingDown       atomic.Bool
	config             *config.Config
	heartbeatSender    *HeartbeatSender
	clientSlicePool    sync.Pool
	applicationContent *interfaces.ApplicationContent
}

var (
	clientManager *ClientManager
	once          sync.Once
)

func NewClientManager(applicationContent *interfaces.ApplicationContent) *ClientManager {
	once.Do(func() {
		c := applicationContent.ConfigManager().Config()
		if clientManager == nil {
			clientManager = &ClientManager{
				clients:            make(map[string]ClientInterface),
				shuttingDown:       atomic.Bool{},
				config:             c,
				applicationContent: applicationContent,
				clientSlicePool: sync.Pool{
					New: func() interface{} {
						return make([]ClientInterface, 0, 128)
					},
				},
			}
			clientManager.heartbeatSender = NewHeartbeatSender(applicationContent.Logger(), c.Server.FSDServer.HeartbeatDuration, clientManager.SendHeartBeat)
		}
	})
	return clientManager
}

func (cm *ClientManager) PutSlice(clients []ClientInterface) {
	cm.clientSlicePool.Put(clients)
}

func (cm *ClientManager) Shutdown(ctx context.Context) error {
	if !cm.shuttingDown.CompareAndSwap(false, true) {
		return fmt.Errorf("shutting down already in progress")
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cm.heartbeatSender.Stop()

	clients := cm.GetClientSnapshot()
	defer cm.PutSlice(clients)

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

func (cm *ClientManager) SendHeartBeat() error {
	if cm.shuttingDown.Load() {
		return nil
	}
	randomInt := rand.Int()
	packet := makePacket(WindDelta, global.FSDServerName, string(AllClient), strconv.Itoa(randomInt%11-5), strconv.Itoa(randomInt%21-10))
	cm.BroadcastMessage(packet, nil, BroadcastToAll)
	return nil
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

func (cm *ClientManager) SendRawMessageTo(from int, to string, message string) error {
	client, exists := cm.GetClient(to)
	if !exists {
		return ErrCallsignNotFound
	}

	bytes := makePacket(Message, fmt.Sprintf("(%04d)", from), to, message)

	client.SendLine(bytes)
	return nil
}

func (cm *ClientManager) BroadcastMessage(message []byte, fromClient ClientInterface, filter BroadcastFilter) {
	if cm.shuttingDown.Load() || len(message) == 0 {
		return
	}

	clients := cm.GetClientSnapshot()
	defer cm.PutSlice(clients) // 重置并放回池中

	if len(clients) == 0 {
		return
	}

	// 准备完整消息（包含分割符）
	fullMsg := make([]byte, len(message), len(message)+len(splitSign))
	copy(fullMsg, message)
	fullMsg = append(fullMsg, splitSign...)

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

			cm.applicationContent.Logger().DebugF("[Broadcast] -> [%s] %s", cl.Callsign(), message)
			cl.SendLineWithoutLog(fullMsg)
		}(client)
	}

	wg.Wait()
}
