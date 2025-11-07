// Package config
package config

import (
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
)

type NavigraphConfig struct {
	Enabled bool   `json:"enable"`
	Token   string `json:"token"`
}

func defaultNavigraphConfig() *NavigraphConfig {
	return &NavigraphConfig{
		Enabled: false,
		Token:   "",
	}
}

func (config *NavigraphConfig) checkValid(_ log.LoggerInterface) *ValidResult {
	return ValidPass()
}
