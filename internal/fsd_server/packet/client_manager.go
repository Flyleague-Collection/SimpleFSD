package packet

import (
	"context"
	"encoding/json"
	"fmt"
	c "github.com/half-nothing/simple-fsd/internal/config"
	"github.com/half-nothing/simple-fsd/internal/interfaces"
	"github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/internal/utils"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type ClientManager struct {
	clients            map[string]fsd.ClientInterface
	lock               sync.RWMutex
	shuttingDown       atomic.Bool
	config             *c.Config
	heartbeatTimer     *utils.IntervalActuator
	whazzupTimer       *utils.IntervalActuator
	clientSlicePool    sync.Pool
	applicationContent *interfaces.ApplicationContent
}

var (
	clientManager *ClientManager
	once          sync.Once
)

func NewClientManager(applicationContent *interfaces.ApplicationContent) *ClientManager {
	once.Do(func() {
		if clientManager == nil {
			clientManager = &ClientManager{
				clients:            make(map[string]fsd.ClientInterface),
				shuttingDown:       atomic.Bool{},
				config:             applicationContent.Config(),
				applicationContent: applicationContent,
				clientSlicePool: sync.Pool{
					New: func() interface{} {
						return make([]fsd.ClientInterface, 0, 128)
					},
				},
			}
			clientManager.heartbeatTimer = utils.NewIntervalActuator(applicationContent.Config().HeartbeatDuration, clientManager.SendHeartBeat)
			clientManager.whazzupTimer = utils.NewIntervalActuator(applicationContent.Config().WhazzupDuration, clientManager.GenerateWhazzupFile)
			clientManager.heartbeatTimer.Start()
			clientManager.whazzupTimer.Start()
		}
	})
	return clientManager
}

func (cm *ClientManager) GenerateWhazzupFile() error {
	data := &fsd.OnlineClients{
		General: &fsd.OnlineGeneral{
			Version:          3,
			ConnectedClients: 0,
			OnlinePilot:      0,
			OnlineController: 0,
		},
		Pilots:      make([]*fsd.OnlinePilot, 0),
		Controllers: make([]*fsd.OnlineController, 0),
	}

	clientCopy := cm.GetClientSnapshot()
	defer cm.PutSlice(clientCopy)

	for _, client := range clientCopy {
		if client == nil || client.Disconnected() {
			continue
		}
		data.General.ConnectedClients++
		if client.IsAtc() {
			data.General.OnlineController++
			controller := &fsd.OnlineController{
				Cid:       client.User().Cid,
				Callsign:  client.Callsign(),
				RealName:  client.RealName(),
				Latitude:  client.Position()[0].Latitude,
				Longitude: client.Position()[0].Longitude,
				Rating:    client.Rating().Index(),
				Facility:  client.Facility().Index(),
				Frequency: client.Frequency() + 100000,
				AtcInfo:   client.AtisInfo(),
				Range:     int(client.VisualRange()),
				LogonTime: client.LogonTime(),
			}
			data.Controllers = append(data.Controllers, controller)
		} else {
			data.General.OnlinePilot++
			pilot := &fsd.OnlinePilot{
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
				LogonTime:   client.LogonTime(),
			}
			data.Pilots = append(data.Pilots, pilot)
		}
	}

	data.General.GenerateTime = time.Now().Format(time.DateTime)

	file, err := os.OpenFile(cm.config.WhazzupFile, os.O_WRONLY|os.O_CREATE, 0655)
	defer file.Close()

	if err != nil {
		return err
	}

	if data, err := json.Marshal(data); err != nil {
		return err
	} else if _, err := file.Write(data); err != nil {
		return err
	}

	return nil
}

func (cm *ClientManager) PutSlice(clients []fsd.ClientInterface) {
	cm.clientSlicePool.Put(clients)
}

func (cm *ClientManager) Shutdown(ctx context.Context) error {
	if !cm.shuttingDown.CompareAndSwap(false, true) {
		return fmt.Errorf("shutting down already in progress")
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cm.heartbeatTimer.Stop()
	cm.whazzupTimer.Stop()

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

func (cm *ClientManager) GetClientSnapshot() []fsd.ClientInterface {
	cm.lock.RLock()
	defer cm.lock.RUnlock()

	// 从池中获取切片
	clients := cm.clientSlicePool.Get().([]fsd.ClientInterface)
	clients = clients[:0]

	// 填充客户端
	for _, client := range cm.clients {
		clients = append(clients, client)
	}
	return clients
}

// 并发断开所有客户端连接
func (cm *ClientManager) disconnectClients(clients []fsd.ClientInterface) {
	if len(clients) == 0 {
		return
	}

	sem := make(chan struct{}, cm.config.MaxBroadcastWorkers)
	var wg sync.WaitGroup

	for _, client := range clients {
		wg.Add(1)
		sem <- struct{}{}

		go func(c fsd.ClientInterface) {
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
	packet := makePacket(fsd.WindDelta, "SERVER", string(fsd.AllClient), strconv.Itoa(randomInt%11-5), strconv.Itoa(randomInt%21-10))
	cm.BroadcastMessage(packet, nil, fsd.BroadcastToAll)
	return nil
}

func (cm *ClientManager) AddClient(client fsd.ClientInterface) error {
	if cm.shuttingDown.Load() {
		return fmt.Errorf("fsd_server shutting down")
	}
	cm.lock.Lock()
	defer cm.lock.Unlock()

	if _, exists := cm.clients[client.Callsign()]; exists {
		return fmt.Errorf("client already registered: %s", client.Callsign())
	}
	cm.clients[client.Callsign()] = client
	return nil
}

func (cm *ClientManager) GetClient(callsign string) (fsd.ClientInterface, bool) {
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
		return fmt.Errorf("fsd_server is shutting down")
	}

	client, exists := cm.GetClient(callsign)
	if !exists {
		return fsd.ErrCallsignNotFound
	}

	client.SendLine(message)
	return nil
}

func (cm *ClientManager) SendRawMessageTo(from int, to string, message string) error {
	client, exists := cm.GetClient(to)
	if !exists {
		return fsd.ErrCallsignNotFound
	}

	bytes := makePacket(fsd.Message, fmt.Sprintf("(%04d)", from), to, message)

	client.SendLine(bytes)
	return nil
}

func (cm *ClientManager) BroadcastMessage(message []byte, fromClient fsd.ClientInterface, filter fsd.BroadcastFilter) {
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
	sem := make(chan struct{}, cm.config.MaxBroadcastWorkers)

	for _, client := range clients {
		if client == fromClient || client.Disconnected() {
			continue
		}

		if filter != nil && !filter(client, fromClient) {
			continue
		}

		wg.Add(1)
		sem <- struct{}{}
		go func(cl fsd.ClientInterface) {
			defer func() {
				<-sem
				wg.Done()
			}()

			c.DebugF("[Broadcast] -> [%s] %s", cl.Callsign(), message)
			cl.SendLineWithoutLog(fullMsg)
		}(client)
	}

	wg.Wait()
}
