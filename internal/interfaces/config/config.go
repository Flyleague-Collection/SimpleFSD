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
	MetarSource   MetarSources    `json:"metar_source"`
	Database      *DatabaseConfig `json:"database"`
	Rating        map[string]int  `json:"rating"`
	Facility      map[string]int  `json:"facility"`
}

func DefaultConfig() *Config {
	return &Config{
		ConfigVersion: ConfVersion.String(),
		Server:        defaultServerConfig(),
		MetarSource:   defaultMetarSources(),
		Database:      defaultDatabaseConfig(),
		Rating:        make(map[string]int),
		Facility:      make(map[string]int),
	}
}

var ErrVersionUnmatch = errors.New("version mismatch")

func (c *Config) CheckValid(logger log.LoggerInterface) *ValidResult {
	if version, err := newVersion(c.ConfigVersion); err != nil {
		return ValidFailWith(errors.New("version string parse fail"), err)
	} else if result := ConfVersion.checkVersion(version); result != AllMatch {
		return ValidFailWith(fmt.Errorf("config version mismatch, expected %s, got %s", ConfVersion.String(), version.String()), ErrVersionUnmatch)
	}
	if result := c.Database.checkValid(logger); result.IsFail() {
		return result
	}
	if result := c.Server.checkValid(logger); result.IsFail() {
		return result
	}
	if result := c.MetarSource.checkValid(logger); result.IsFail() {
		return result
	}
	return ValidPass()
}
