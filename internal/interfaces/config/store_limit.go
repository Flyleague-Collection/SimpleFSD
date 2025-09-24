// Package config
package config

import (
	"errors"
	"github.com/half-nothing/simple-fsd/internal/interfaces/global"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"os"
	"path/filepath"
)

type HttpServerStoreFileLimit struct {
	MaxFileSize    int64    `json:"max_file_size"`
	AllowedFileExt []string `json:"allowed_file_ext"`
	StorePrefix    string   `json:"store_prefix"`
	StoreInServer  bool     `json:"store_in_server"`
	RootPath       string   `json:"-"`
}

func (config *HttpServerStoreFileLimit) checkValid(_ log.LoggerInterface) *ValidResult {
	if config.MaxFileSize < 0 {
		return ValidFail(errors.New("invalid json field http_server.store.max_file_size, cannot be negative"))
	}
	return ValidPass()
}

type HttpServerStoreFileLimits struct {
	ImageLimit *HttpServerStoreFileLimit `json:"image_limit"`
	FileLimit  *HttpServerStoreFileLimit `json:"file_limit"`
}

func defaultHttpServerStoreFileLimits() *HttpServerStoreFileLimits {
	return &HttpServerStoreFileLimits{
		ImageLimit: &HttpServerStoreFileLimit{
			MaxFileSize:    5 * 1024 * 1024,
			AllowedFileExt: []string{".jpg", ".png", ".bmp", ".jpeg"},
			StorePrefix:    "images",
			StoreInServer:  true,
		},
		FileLimit: &HttpServerStoreFileLimit{
			MaxFileSize:    10 * 1024 * 1024,
			AllowedFileExt: []string{".md", ".txt", ".pdf", ".doc", ".docx"},
			StorePrefix:    "files",
			StoreInServer:  true,
		},
	}
}

func (config *HttpServerStoreFileLimits) checkValid(logger log.LoggerInterface) *ValidResult {
	if result := config.ImageLimit.checkValid(logger); result.IsFail() {
		return result
	}
	return ValidPass()
}

func (config *HttpServerStoreFileLimits) CheckLocalStore(_ log.LoggerInterface, localStore bool) *ValidResult {
	if !localStore {
		return ValidPass()
	}
	if !config.ImageLimit.StoreInServer {
		return ValidFail(errors.New("when you use local store, store_in_server must be true"))
	}
	if !config.FileLimit.StoreInServer {
		return ValidFail(errors.New("when you use local store, store_in_server must be true"))
	}
	return ValidPass()
}

func (config *HttpServerStoreFileLimits) CreateDir(_ log.LoggerInterface, root string) *ValidResult {
	config.ImageLimit.RootPath = root
	if config.ImageLimit.StoreInServer {
		imagePath := filepath.Join(root, config.ImageLimit.StorePrefix)
		if err := os.MkdirAll(imagePath, global.DefaultDirectoryPermission); err != nil {
			return ValidFailWith(errors.New("error creating the image directory"), err)
		}
	}
	config.FileLimit.RootPath = root
	if config.FileLimit.StoreInServer {
		filePath := filepath.Join(root, config.FileLimit.StorePrefix)
		if err := os.MkdirAll(filePath, global.DefaultDirectoryPermission); err != nil {
			return ValidFailWith(errors.New("error creating the image directory"), err)
		}
	}
	return ValidPass()
}
