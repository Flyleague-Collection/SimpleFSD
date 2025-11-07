package utils

import (
	"time"
)

type IntervalActuator struct {
	interval time.Duration
	ticker   *time.Ticker
	stopChan chan struct{}
	callback func()
}

func NewIntervalActuator(interval time.Duration, callback func()) *IntervalActuator {
	return &IntervalActuator{
		interval: interval,
		stopChan: make(chan struct{}),
		callback: callback,
	}
}

func (h *IntervalActuator) Start() {
	h.ticker = time.NewTicker(h.interval)

	go func() {
		for {
			select {
			case <-h.ticker.C:
				h.callback()
			case <-h.stopChan:
				return
			}
		}
	}()
}

func (h *IntervalActuator) Stop() {
	if h.ticker != nil {
		h.ticker.Stop()
	}
	close(h.stopChan)
}
