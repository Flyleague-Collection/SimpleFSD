// Package config
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/half-nothing/simple-fsd/src/interfaces/global"
	"github.com/half-nothing/simple-fsd/src/interfaces/logger"
)

type HttpServerStore struct {
	StoreType       int                        `json:"store_type"`        // 文件存储类型, 0: 本地存储, 1: 阿里云OSS存储, 2: 腾讯云对象存储
	Region          string                     `json:"region"`            // 云存储地域
	Bucket          string                     `json:"bucket"`            // 云存储桶名
	AccessId        string                     `json:"access_id"`         // 访问id
	AccessKey       string                     `json:"access_key"`        // 访问秘钥
	CdnDomain       string                     `json:"cdn_domain"`        // 自定义加速域名
	UseInternalUrl  bool                       `json:"use_internal_url"`  // 上传使用内部域名
	LocalStorePath  string                     `json:"local_store_path"`  // 本地存储路径
	RemoteStorePath string                     `json:"remote_store_path"` // 远程存储路径
	FileLimit       *HttpServerStoreFileLimits `json:"file_limit"`
}

func defaultHttpServerStore() *HttpServerStore {
	return &HttpServerStore{
		StoreType:       0,
		Region:          "",
		Bucket:          "",
		AccessId:        "",
		AccessKey:       "",
		CdnDomain:       "",
		UseInternalUrl:  false,
		LocalStorePath:  "uploads",
		RemoteStorePath: "",
		FileLimit:       defaultHttpServerStoreFileLimits(),
	}
}

func (config *HttpServerStore) checkValid(logger logger.LoggerInterface) *ValidResult {
	if result := config.FileLimit.checkValid(logger); result.IsFail() {
		return result
	}
	if config.LocalStorePath == "" {
		return ValidFail(errors.New("invalid json field http_server.store.local_store_path, path cannot be empty"))
	}
	if err := os.MkdirAll(filepath.Clean(config.LocalStorePath), global.DefaultDirectoryPermission); err != nil {
		return ValidFailWith(fmt.Errorf("error while creating local store path(%s)", config.LocalStorePath), err)
	}
	if result := config.FileLimit.CreateDir(logger, config.LocalStorePath, config.RemoteStorePath); result.IsFail() {
		return result
	}
	switch config.StoreType {
	case 0:
		if result := config.FileLimit.CheckLocalStore(logger, true); result.IsFail() {
			return result
		}
		// 本地存储
		// 不用任何额外操作, 仅占位使用
	case 1, 2:
		if result := config.FileLimit.CheckLocalStore(logger, false); result.IsFail() {
			return result
		}
		// 阿里云OSS存储或者腾讯云对象存储
		if config.Region == "" {
			return ValidFail(errors.New("invalid json field http_server.store.region, region cannot be empty"))
		}
		if config.Bucket == "" {
			return ValidFail(errors.New("invalid json field http_server.store.bucket, bucket cannot be empty"))
		}
		if config.AccessId == "" {
			return ValidFail(errors.New("invalid json field http_server.store.access_id, access_id cannot be empty"))
		}
		if config.AccessKey == "" {
			return ValidFail(errors.New("invalid json field http_server.store.access_key, access_key cannot be empty"))
		}
	default:
		return ValidFail(fmt.Errorf("invalid json field http_server.store_type %d, only support 0, 1, 2", config.StoreType))
	}
	return ValidPass()
}
