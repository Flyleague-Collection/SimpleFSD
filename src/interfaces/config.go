// Package interfaces
package interfaces

import (
	. "github.com/half-nothing/simple-fsd/src/interfaces/config"
)

type ConfigManagerInterface interface {
	Config() *Config
	SaveConfig() error
}
