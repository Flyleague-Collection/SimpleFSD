// Package config
package config

import "github.com/half-nothing/simple-fsd/internal/interfaces/log"

type ServerConfig struct {
	General     *GeneralConfig     `json:"general"`
	FSDServer   *FSDServerConfig   `json:"fsd_server"`
	HttpServer  *HttpServerConfig  `json:"http_server"`
	VoiceServer *VoiceServerConfig `json:"voice_server"`
	GRPCServer  *GRPCServerConfig  `json:"grpc_server"`
}

func defaultServerConfig() *ServerConfig {
	return &ServerConfig{
		General:     defaultOtherConfig(),
		FSDServer:   defaultFSDServerConfig(),
		HttpServer:  defaultHttpServerConfig(),
		VoiceServer: defaultVoiceServerConfig(),
		GRPCServer:  defaultGRPCServerConfig(),
	}
}

func (config *ServerConfig) checkValid(logger log.LoggerInterface) *ValidResult {
	if result := config.General.checkValid(logger); result.IsFail() {
		return result
	}
	if result := config.FSDServer.checkValid(logger); result.IsFail() {
		return result
	}
	if result := config.HttpServer.checkValid(logger); result.IsFail() {
		return result
	}
	if result := config.VoiceServer.checkValid(logger); result.IsFail() {
		return result
	}
	if result := config.GRPCServer.checkValid(logger); result.IsFail() {
		return result
	}
	return ValidPass()
}
