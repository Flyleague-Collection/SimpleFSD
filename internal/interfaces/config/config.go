// Package config
package config

import (
	"errors"
	"fmt"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
)

type Config struct {
	ConfigVersion string          `json:"config_version"`
	Server        *ServerConfig   `json:"server"`
	Database      *DatabaseConfig `json:"database"`
	Rating        map[string]int  `json:"rating"`
}

func DefaultConfig() *Config {
	return &Config{
		ConfigVersion: ConfVersion.String(),
		Server:        defaultServerConfig(),
		Database:      defaultDatabaseConfig(),
		Rating:        make(map[string]int),
	}
}

func (c *Config) CheckValid(logger log.LoggerInterface) *ValidResult {
	if version, err := newVersion(c.ConfigVersion); err != nil {
		return ValidFailWith(errors.New("version string parse fail"), err)
	} else if result := ConfVersion.checkVersion(version); result != AllMatch {
		return ValidFail(fmt.Errorf("config version mismatch, expected %s, got %s", ConfVersion.String(), version.String()))
	}
	if result := c.Database.checkValid(logger); result.IsFail() {
		return result
	}
	if result := c.Server.checkValid(logger); result.IsFail() {
		return result
	}
	return ValidPass()
}
