// Package config
package config

import "github.com/half-nothing/simple-fsd/src/interfaces/logger"

type SSLConfig struct {
	Enable          bool   `json:"enable"`
	EnableHSTS      bool   `json:"enable_hsts"`
	ForceSSL        bool   `json:"force_ssl"`
	HstsExpiredTime int    `json:"hsts_expired_time"`
	IncludeDomain   bool   `json:"include_domain"`
	CertFile        string `json:"cert_file"`
	KeyFile         string `json:"key_file"`
}

func defaultSSLConfig() *SSLConfig {
	return &SSLConfig{
		Enable:          false,
		EnableHSTS:      false,
		ForceSSL:        false,
		HstsExpiredTime: 5184000,
		IncludeDomain:   false,
		CertFile:        "",
		KeyFile:         "",
	}
}

func (config *SSLConfig) checkValid(logger logger.LoggerInterface) *ValidResult {
	if config.Enable {
		if config.CertFile == "" || config.KeyFile == "" {
			logger.WarnF("HTTPS server requires both cert and key files. Cert: %s, Key: %s. Falling back to HTTP", config.CertFile, config.KeyFile)
			config.Enable = false
		}
	}
	if !config.Enable && config.EnableHSTS {
		logger.Warn("You can not enable HSTS when ssl is not enable!")
		config.EnableHSTS = false
		config.HstsExpiredTime = 0
		config.IncludeDomain = false
	}
	if !config.Enable && config.ForceSSL {
		logger.Warn("You can not force ssl when ssl is not enable!")
		config.EnableHSTS = false
		config.ForceSSL = false
	}
	return ValidPass()
}
