// Package config
package config

type ManagerInterface interface {
	Config() *Config
	SaveConfig() error
}
