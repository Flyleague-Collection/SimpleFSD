// Package message
package message

import (
	"math/rand/v2"
	"sync/atomic"
	"testing"
	"time"

	"github.com/half-nothing/simple-fsd/src/base"
	"github.com/half-nothing/simple-fsd/src/interfaces/queue"
)

func ExampleNewAsyncMessageQueue() {
	logger := base.NewLogger()
	logger.Init("", "", true, true)
	messageQueue := NewAsyncMessageQueue(logger, 128)
	messageQueue.Subscribe(queue.SendVerifyEmail, func(message *queue.Message) error {
		return nil
	})
	messageNumber := 4096
	for i := 0; i < messageNumber; i++ {
		messageQueue.Publish(&queue.Message{
			Type: queue.SendMessageToClient,
			Data: uint64(i),
		})
	}
}

func TestMessageQueue(t *testing.T) {
	logger := base.NewLogger()
	logger.Init("", "", true, true)
	timeReceiveVerifyEmail := atomic.Int32{}
	timeSendVerifyEmail := atomic.Int32{}
	timeReceiveMessageToClient := atomic.Int32{}
	timeSendMessageToClient := atomic.Int32{}
	messageQueue := NewAsyncMessageQueue(logger, 128)
	messageQueue.Subscribe(queue.SendVerifyEmail, func(message *queue.Message) error {
		timeReceiveVerifyEmail.Add(1)
		return nil
	})
	messageQueue.Subscribe(queue.SendMessageToClient, func(message *queue.Message) error {
		timeReceiveMessageToClient.Add(1)
		return nil
	})
	logger.Info("Message publish start")
	messageNumber := 4096
	startTime := time.Now()
	for i := 0; i < messageNumber; i++ {
		if rand.IntN(100) < 50 {
			timeSendVerifyEmail.Add(1)
			messageQueue.Publish(&queue.Message{
				Type: queue.SendVerifyEmail,
				Data: uint64(i),
			})
		} else {
			timeSendMessageToClient.Add(1)
			messageQueue.Publish(&queue.Message{
				Type: queue.SendMessageToClient,
				Data: uint64(i),
			})
		}
	}
	endTime := time.Now()
	logger.InfoF("Message publish end, publish %d messages, cost %dns(around %dms)", messageNumber, endTime.Sub(startTime), endTime.Sub(startTime).Milliseconds())
	logger.InfoF("SendVerifyEmail %d/%d; SendMessageToClient %d/%d", timeSendVerifyEmail.Load(), timeReceiveVerifyEmail.Load(), timeSendMessageToClient.Load(), timeReceiveMessageToClient.Load())
	messageQueue.Stop()
	logger.InfoF("SendVerifyEmail %d/%d; SendMessageToClient %d/%d", timeSendVerifyEmail.Load(), timeReceiveVerifyEmail.Load(), timeSendMessageToClient.Load(), timeReceiveMessageToClient.Load())
	if timeSendVerifyEmail.Load() != timeReceiveVerifyEmail.Load() || timeSendMessageToClient.Load() != timeReceiveMessageToClient.Load() {
		t.Fatal("Unmatch send and receive message")
	}
}
