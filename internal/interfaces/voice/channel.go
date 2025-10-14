// Package voice
package voice

import (
	"sync"
	"time"
)

type ChannelFrequency int

type Channel struct {
	Frequency    ChannelFrequency
	ClientsMutex sync.RWMutex
	Clients      map[int]*Transmitter
	CreatedAt    time.Time
}
