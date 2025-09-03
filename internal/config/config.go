package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/half-nothing/simple-fsd/internal/utils"
	"os"
	"time"
)

func checkPort(port uint) (bool, error) {
	if port <= 0 {
		return false, errors.New("port must be greater than zero")
	}
	if port > 65535 {
		return false, errors.New("port must be less than 65535")
	}
	if port < 1024 {
		WarnF("The %d port may have a special usage, use it with caution")
	}
	return true, nil
}

type Config struct {
	DebugMode            bool           `json:"debug_mode"` // 是否启用调试模式
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
	Motd                 []string       `json:"motd"`
	Rating               map[string]int `json:"rating"`
}

func defaultConfig() *Config {
	return &Config{
		DebugMode:           false,
		ConfigVersion:       confVersion.String(),
		FSDName:             "SimpleFSD",
		Host:                "localhost",
		Port:                6809,
		HeartbeatInterval:   "60s",
		SessionCleanTime:    "40s",
		MaxWorkers:          128,
		MaxBroadcastWorkers: 128,
		CertFile:            "cert.txt",
		EncryptionType:      int(SHA256),
		Rating:              make(map[string]int),
	}
}

type EncryptionType int

const (
	NoEncryption EncryptionType = iota
	MD5
	SHA256
	BCRYPT
)

var (
	config = utils.NewCachedValue[Config](0, func() *Config {
		if config, err := readConfig(); err != nil {
			FatalF("Error occurred while reading config %v", err)
			panic(err)
		} else {
			return config
		}
	})
)

// readConfig 从配置文件读取配置
func readConfig() (*Config, error) {
	config := defaultConfig()

	// 读取配置文件
	if bytes, err := os.ReadFile("config.json"); err != nil {
		// 如果配置文件不存在，创建默认配置
		if err := config.SaveConfig(); err != nil {
			return nil, err
		}
		return nil, errors.New("the configuration file does not exist and has been created. Please try again after editing the configuration file")
	} else if err := json.Unmarshal(bytes, config); err != nil {
		// 解析JSON配置
		return nil, fmt.Errorf("the configuration file does not contain valid JSON, %v", err)
	} else if pass, err := config.CheckValid(); !pass {
		return nil, err
	}
	return config, nil
}

func (c *Config) CheckValid() (bool, error) {
	if version, err := newVersion(c.ConfigVersion); err != nil {
		return false, err
	} else if result := confVersion.checkVersion(version); result != AllMatch {
		return false, fmt.Errorf("config version mismatch, expected %s, got %s", confVersion.String(), version.String())
	}

	if pass, err := checkPort(c.Port); !pass {
		return pass, err
	}

	c.Address = fmt.Sprintf("%s:%d", c.Host, c.Port)

	if duration, err := time.ParseDuration(c.SessionCleanTime); err != nil {
		return false, fmt.Errorf("invalid json field session_clean_time, duration parse error, %v", err)
	} else {
		c.SessionCleanDuration = duration
	}

	if duration, err := time.ParseDuration(c.HeartbeatInterval); err != nil {
		return false, fmt.Errorf("invalid json field heartbead_interval, duration parse error, %v", err)
	} else {
		c.HeartbeatDuration = duration
	}

	if c.EncryptionType < int(NoEncryption) || c.EncryptionType > int(BCRYPT) {
		return false, fmt.Errorf("invalid encryption type, encryption type must be between %d and %d", NoEncryption, BCRYPT)
	}

	return true, nil
}

func (c *Config) SaveConfig() error {
	if writer, err := os.OpenFile("config.json", os.O_WRONLY|os.O_CREATE, 0655); err != nil {
		return err
	} else if data, err := json.MarshalIndent(c, "", "\t"); err != nil {
		return err
	} else if _, err = writer.Write(data); err != nil {
		return err
	} else if err := writer.Close(); err != nil {
		return err
	}
	return nil
}

func GetConfig() *Config {
	return config.GetValue()
}
