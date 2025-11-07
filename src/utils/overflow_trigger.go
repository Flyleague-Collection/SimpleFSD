// Package utils
package utils

import (
	"sync"
)

type OverflowTrigger struct {
	mu          sync.Mutex
	count       int
	targetValue int
	callback    func()
}

func NewOverflowTrigger(targetValue int, callback func()) *OverflowTrigger {
	return &OverflowTrigger{
		mu:          sync.Mutex{},
		count:       0,
		targetValue: targetValue,
		callback:    callback,
	}
}

func (trigger *OverflowTrigger) Tick() {
	if trigger.targetValue <= 0 {
		return
	}
	if trigger.targetValue == 1 {
		trigger.callback()
		return
	}
	trigger.mu.Lock()
	defer trigger.mu.Unlock()
	trigger.count++
	if trigger.count >= trigger.targetValue {
		trigger.callback()
		trigger.count = 0
	}
}

func (trigger *OverflowTrigger) Reset() {
	trigger.mu.Lock()
	defer trigger.mu.Unlock()
	trigger.count = 0
}
