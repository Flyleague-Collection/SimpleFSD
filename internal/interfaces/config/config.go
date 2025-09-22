// Package config
package config

import (
	"errors"
	"fmt"
	"github.com/half-nothing/simple-fsd/internal/interfaces/global"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"os"
	"time"
)

type EncryptionType int

const (
	NoEncryption EncryptionType = iota
	MD5
	SHA256
	BCRYPT
)

type Config struct {
	ConfigVersion        string         `json:"config_version"`
	FSDName              string         `json:"fsd_name"` // 应用名称
	Host                 string         `json:"host"`
	Port                 uint           `json:"port"`
	Address              string         `json:"-"`
	HeartbeatInterval    string         `json:"heartbeat_interval"`
	HeartbeatDuration    time.Duration  `json:"-"`
	SessionCleanTime     string         `json:"session_clean_time"` // 会话保留时间
	SessionCleanDuration time.Duration  `json:"-"`
	MaxWorkers           int            `json:"max_workers"`           // 并发线程数
	MaxBroadcastWorkers  int            `json:"max_broadcast_workers"` // 广播并发线程数
	CertFile             string         `json:"cert_file"`
	EncryptionType       int            `json:"encryption_type"`
	WhazzupFile          string         `json:"whazzup_file"`
	WhazzupInterval      string         `json:"whazzup_interval"`
	WhazzupDuration      time.Duration  `json:"-"`
	RangeLimit           *FsdRangeLimit `json:"range_limit"`
	FirstMotdLine        string         `json:"first_motd_line"`
	Motd                 []string       `json:"motd"`
	Rating               map[string]int `json:"rating"`
	Facility             map[string]int `json:"facility"`
}

func DefaultConfig() *Config {
	return &Config{
		ConfigVersion:       ConfVersion.String(),
		FSDName:             "SimpleFSD",
		Host:                "localhost",
		Port:                6809,
		HeartbeatInterval:   "40s",
		SessionCleanTime:    "40s",
		MaxWorkers:          128,
		MaxBroadcastWorkers: 128,
		CertFile:            "cert.txt",
		EncryptionType:      int(SHA256),
		WhazzupFile:         "whazzup.json",
		WhazzupInterval:     "15s",
		RangeLimit:          defaultFsdRangeLimitConfig(),
		FirstMotdLine:       "Welcome to use %[1]s v%[2]s",
		Motd:                make([]string, 0),
		Rating:              make(map[string]int),
		Facility:            make(map[string]int),
	}
}

func (config *Config) CheckValid(logger log.LoggerInterface) *ValidResult {
	if version, err := newVersion(config.ConfigVersion); err != nil {
		return ValidFailWith(errors.New("invalid json field config_version"), err)
	} else if result := ConfVersion.checkVersion(version); result != AllMatch {
		return ValidFail(fmt.Errorf("config version mismatch, expected %s, got %s", ConfVersion.String(), version.String()))
	}

	if result := config.RangeLimit.checkValid(logger); result.IsFail() {
		return result
	}

	if result := checkPort(config.Port); result.IsFail() {
		return result
	}

	config.FirstMotdLine = fmt.Sprintf(config.FirstMotdLine, config.FSDName, global.AppVersion)
	data := make([]string, 0, 1+len(config.Motd))
	data = append(data, config.FirstMotdLine)
	data = append(data, config.Motd...)
	config.Motd = data

	config.Address = fmt.Sprintf("%s:%d", config.Host, config.Port)

	if duration, err := time.ParseDuration(config.SessionCleanTime); err != nil {
		return ValidFailWith(errors.New("invalid json field session_clean_time, duration parse error"), err)
	} else {
		config.SessionCleanDuration = duration
	}

	if duration, err := time.ParseDuration(config.HeartbeatInterval); err != nil {
		return ValidFailWith(errors.New("invalid json field heartbead_interval, duration parse error"), err)
	} else if duration <= 25*time.Second {
		return ValidFail(fmt.Errorf("heartbead_interval must larger than 25s, got %.0fs", duration.Seconds()))
	} else {
		config.HeartbeatDuration = duration
	}

	if duration, err := time.ParseDuration(config.WhazzupInterval); err != nil {
		return ValidFailWith(errors.New("invalid json field whazzup_interval, duration parse error"), err)
	} else {
		config.WhazzupDuration = duration
	}

	if _, err := os.Stat(config.WhazzupFile); err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(config.WhazzupFile)
			if err != nil {
				return ValidFailWith(errors.New("error creating whazzup file, %v"), err)
			}
			_ = file.Close()
		} else {
			return ValidFail(fmt.Errorf("can not check whazzup file, %v, %v", config.WhazzupFile, err))
		}
	}

	if config.EncryptionType < int(NoEncryption) || config.EncryptionType > int(BCRYPT) {
		return ValidFail(fmt.Errorf("invalid encryption type, encryption type must be between %d and %d", NoEncryption, BCRYPT))
	}

	return ValidPass()
}
