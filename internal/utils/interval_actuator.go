package utils

import (
	"fmt"
	"time"
)

type IntervalActuator struct {
	interval time.Duration
	ticker   *time.Ticker
	stopChan chan struct{}
	callback func() error
}

func NewIntervalActuator(interval time.Duration, callback func() error) *IntervalActuator {
	return &IntervalActuator{
		interval: interval,
		stopChan: make(chan struct{}),
		callback: callback,
	}
}

func (h *IntervalActuator) Start() {
	h.ticker = time.NewTicker(h.interval)

	go func() {
		defer fmt.Println("Actuator stopped")

		for {
			select {
			case <-h.ticker.C:
				err := h.callback()
				if err != nil {
					fmt.Printf("Error actuator function: %v\n", err)
				}
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
