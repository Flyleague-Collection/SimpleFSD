// Package message
package message

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/half-nothing/simple-fsd/src/interfaces/global"
	"github.com/half-nothing/simple-fsd/src/interfaces/log"
	"github.com/half-nothing/simple-fsd/src/interfaces/queue"
	"golang.org/x/sync/errgroup"
)

type ShutdownCallback struct {
	asyncMessageQueue queue.MessageQueueInterface
}

func NewShutdownCallback(asyncMessageQueue queue.MessageQueueInterface) *ShutdownCallback {
	return &ShutdownCallback{
		asyncMessageQueue: asyncMessageQueue,
	}
}

func (shutdownCallback *ShutdownCallback) Invoke(ctx context.Context) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	done := make(chan struct{})
	go func() {
		shutdownCallback.asyncMessageQueue.Stop()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-timeoutCtx.Done():
		return timeoutCtx.Err()
	}
}

type AsyncMessageQueue struct {
	logger           log.LoggerInterface
	messageCh        chan *queue.Message
	subscribeMap     map[queue.MessageType][]queue.Subscriber
	shutdownCallback global.Callable
	wg               sync.WaitGroup
	lock             sync.RWMutex
}

func NewAsyncMessageQueue(
	logger log.LoggerInterface,
	cacheSize int,
) *AsyncMessageQueue {
	asyncMessageQueue := &AsyncMessageQueue{
		logger:       logger,
		messageCh:    make(chan *queue.Message, cacheSize),
		subscribeMap: make(map[queue.MessageType][]queue.Subscriber),
		wg:           sync.WaitGroup{},
		lock:         sync.RWMutex{},
	}

	asyncMessageQueue.shutdownCallback = NewShutdownCallback(asyncMessageQueue)

	asyncMessageQueue.Start()

	return asyncMessageQueue
}

func (asyncMessageQueue *AsyncMessageQueue) handleMessage(message *queue.Message) error {
	asyncMessageQueue.lock.RLock()
	subscribers, exist := asyncMessageQueue.subscribeMap[message.Type]
	asyncMessageQueue.lock.RUnlock()
	if !exist {
		asyncMessageQueue.logger.WarnF("No subscribers for message type %s", message.Type.String())
		return fmt.Errorf("no subscribers for message type %s", message.Type.String())
	}
	var eg errgroup.Group
	for _, subscriber := range subscribers {
		eg.Go(func() error { return subscriber(message) })
	}
	if err := eg.Wait(); err != nil {
		asyncMessageQueue.logger.ErrorF("Error in handling message type %s: %s", message.Type.String(), err.Error())
		return err
	}
	return nil
}

func (asyncMessageQueue *AsyncMessageQueue) Start() {
	go func() {
		for message := range asyncMessageQueue.messageCh {
			asyncMessageQueue.wg.Add(1)
			go func() {
				defer asyncMessageQueue.wg.Done()
				asyncMessageQueue.handleMessage(message)
			}()
		}
	}()
}

func (asyncMessageQueue *AsyncMessageQueue) ShutdownCallback() global.Callable {
	return asyncMessageQueue.shutdownCallback
}

func (asyncMessageQueue *AsyncMessageQueue) Stop() {
	close(asyncMessageQueue.messageCh)
	asyncMessageQueue.wg.Wait()
}

func (asyncMessageQueue *AsyncMessageQueue) Publish(message *queue.Message) {
	asyncMessageQueue.messageCh <- message
}

func (asyncMessageQueue *AsyncMessageQueue) SyncPublish(message *queue.Message) error {
	return asyncMessageQueue.handleMessage(message)
}

func (asyncMessageQueue *AsyncMessageQueue) Subscribe(messageType queue.MessageType, handler queue.Subscriber) {
	asyncMessageQueue.lock.Lock()
	defer asyncMessageQueue.lock.Unlock()
	if subscribers, exist := asyncMessageQueue.subscribeMap[messageType]; !exist {
		asyncMessageQueue.subscribeMap[messageType] = []queue.Subscriber{handler}
	} else {
		asyncMessageQueue.subscribeMap[messageType] = append(subscribers, handler)
	}
}
