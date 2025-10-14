// Package config
package config

import (
	"errors"

	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"golang.org/x/crypto/bcrypt"
)

type GeneralConfig struct {
	SimulatorServer bool `json:"simulator_server"`
	BcryptCost      int  `json:"bcrypt_cost"`
}

func defaultOtherConfig() *GeneralConfig {
	return &GeneralConfig{
		SimulatorServer: false,
		BcryptCost:      12,
	}
}

func (config *GeneralConfig) checkValid(_ log.LoggerInterface) *ValidResult {
	if config.BcryptCost < bcrypt.MinCost || config.BcryptCost > bcrypt.MaxCost {
		return ValidFail(errors.New("bcrypt_cost out of range, must between 4 and 31"))
	}
	return ValidPass()
}
