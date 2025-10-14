// Package config
package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
)

type GRPCServerConfig struct {
	Enabled       bool          `json:"enabled"`
	Host          string        `json:"host"`
	Port          uint          `json:"port"`
	Address       string        `json:"-"`
	CacheTime     string        `json:"whazzup_cache_time"`
	CacheDuration time.Duration `json:"-"`
}

func defaultGRPCServerConfig() *GRPCServerConfig {
	return &GRPCServerConfig{
		Enabled:   false,
		Host:      "0.0.0.0",
		Port:      6811,
		CacheTime: "15s",
	}
}

func (config *GRPCServerConfig) checkValid(_ log.LoggerInterface) *ValidResult {
	if config.Enabled {
		if result := checkPort(config.Port); result.IsFail() {
			return result
		}
		config.Address = fmt.Sprintf("%s:%d", config.Host, config.Port)

		if duration, err := time.ParseDuration(config.CacheTime); err != nil {
			return ValidFailWith(errors.New("invalid json field grpc_server.cache_time"), err)
		} else {
			config.CacheDuration = duration
		}
	}
	return ValidPass()
}
