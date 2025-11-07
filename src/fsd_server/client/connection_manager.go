// Package client
package client

import (
	"slices"
	"sync"

	. "github.com/half-nothing/simple-fsd/src/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/src/interfaces/log"
)

type ConnectionManager struct {
	logger      log.LoggerInterface
	connections map[int][]ClientInterface
	lock        sync.RWMutex
}

func NewConnectionManager(
	logger log.LoggerInterface,
) *ConnectionManager {
	return &ConnectionManager{
		logger:      log.NewLoggerAdapter(logger, "ConnectionManager"),
		connections: make(map[int][]ClientInterface),
		lock:        sync.RWMutex{},
	}
}

func (cm *ConnectionManager) AddConnection(client ClientInterface) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	cid := client.User().Cid
	cm.logger.DebugF("New connection: %d(%s)", cid, client.Callsign())
	if val, ok := cm.connections[cid]; ok {
		cm.connections[cid] = append(val, client)
	} else {
		cm.connections[cid] = []ClientInterface{client}
	}
}

func (cm *ConnectionManager) RemoveConnection(client ClientInterface) error {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	cid := client.User().Cid
	cm.logger.DebugF("Removing connection: %d(%s)", cid, client.Callsign())
	val, ok := cm.connections[cid]
	if !ok {
		return ErrCidNotFound
	}
	index := slices.Index(val, client)
	if index == -1 {
		return ErrConnectionNotFound
	}
	val = slices.Delete(val, index, index+1)
	if len(val) == 0 {
		delete(cm.connections, cid)
	} else {
		cm.connections[cid] = val
	}
	return nil
}

func (cm *ConnectionManager) GetConnections(cid int) ([]ClientInterface, error) {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	val, ok := cm.connections[cid]
	if !ok {
		return nil, ErrCidNotFound
	}
	return val, nil
}
