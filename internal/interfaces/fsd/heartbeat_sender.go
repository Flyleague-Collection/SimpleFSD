package fsd

import (
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"time"
)

type Heartbeat func() error

type HeartbeatSender struct {
	logger   log.LoggerInterface
	interval time.Duration
	ticker   *time.Ticker
	stopChan chan struct{}
	sendFunc Heartbeat
}

func NewHeartbeatSender(
	logger log.LoggerInterface,
	interval time.Duration,
	sendFunc Heartbeat,
) *HeartbeatSender {
	return &HeartbeatSender{
		logger:   logger,
		interval: interval,
		stopChan: make(chan struct{}),
		sendFunc: sendFunc,
	}
}

func (h *HeartbeatSender) Start() {
	h.ticker = time.NewTicker(h.interval)

	go func() {
		for {
			select {
			case <-h.ticker.C:
				err := h.sendFunc()
				if err != nil {
					h.logger.ErrorF("Error sending heartbeat: %v\n", err)
				}
			case <-h.stopChan:
				return
			}
		}
	}()
}

func (h *HeartbeatSender) Stop() {
	if h.ticker != nil {
		h.ticker.Stop()
	}
	close(h.stopChan)
}
