// Package service
package service

import (
	"fmt"
	"github.com/half-nothing/simple-fsd/internal/interfaces/config"
	"mime/multipart"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

var (
	ErrFilePathFail       = NewApiStatus("FILE_PATH_FAIL", "文件上传失败", ServerInternalError)
	ErrFileSaveFail       = NewApiStatus("FILE_SAVE_FAIL", "文件保存失败", ServerInternalError)
	ErrFileUploadFail     = NewApiStatus("FILE_UPLOAD_FAIL", "文件上传失败", ServerInternalError)
	ErrFileOverSize       = NewApiStatus("FILE_OVER_SIZE", "文件过大", BadRequest)
	ErrFileExtUnsupported = NewApiStatus("FILE_EXT_UNSUPPORTED", "不支持的文件类型", BadRequest)
	ErrFileNameIllegal    = NewApiStatus("FILE_NAME_ILLEGAL", "文件名不合法", BadRequest)
	SuccessUploadFile     = NewApiStatus("UPLOAD_FILE", "文件上传成功", Ok)
)

type FileType int

const (
	IMAGES FileType = iota
	FILES
	UNKNOWN
)

// StoreInfo 文件存储信息
type StoreInfo struct {
	FileType      FileType                         // 文件类型 [FileType]
	FileLimit     *config.HttpServerStoreFileLimit // 该类型文件限制 [config.HttpServerStoreFileLimit]
	RootPath      string                           // 存储根目录
	FilePath      string                           // 文件存储路径
	RemotePath    string                           // 远程文件存储路径
	FileName      string                           // 文件名
	FileExt       string                           // 文件扩展名
	FileContent   *multipart.FileHeader            // 文件内容 [multipart.FileHeader]
	StoreInServer bool                             // 是否保存在本地
}

func NewStoreInfo(fileType FileType, fileLimit *config.HttpServerStoreFileLimit, file *multipart.FileHeader) *StoreInfo {
	return &StoreInfo{
		FileType:      fileType,
		FileLimit:     fileLimit,
		RootPath:      fileLimit.RootPath,
		FilePath:      "",
		FileName:      "",
		RemotePath:    "",
		FileExt:       "",
		FileContent:   file,
		StoreInServer: fileLimit.StoreInServer,
	}
}

func (fileType FileType) GenerateStoreInfo(fileLimit *config.HttpServerStoreFileLimit, file *multipart.FileHeader) (*StoreInfo, *ApiStatus) {
	if strings.Contains(file.Filename, string(filepath.Separator)) {
		return nil, ErrFileNameIllegal
	}

	ext := filepath.Ext(file.Filename)

	if !slices.Contains(fileLimit.AllowedFileExt, ext) {
		return nil, ErrFileExtUnsupported
	}

	if file.Size > fileLimit.MaxFileSize {
		return nil, ErrFileOverSize
	}

	storeInfo := NewStoreInfo(fileType, fileLimit, file)

	storeInfo.FileExt = ext
	storeInfo.FileName = filepath.Join(fileLimit.StorePrefix, fmt.Sprintf("%d%s", time.Now().UnixNano(), ext))
	storeInfo.FilePath = filepath.Join(fileLimit.RootPath, storeInfo.FileName)
	storeInfo.RemotePath = strings.Replace(storeInfo.FileName, "\\", "/", -1)

	return storeInfo, nil
}

type StoreServiceInterface interface {
	SaveImageFile(file *multipart.FileHeader) (*StoreInfo, *ApiStatus)
	DeleteImageFile(file string) (*StoreInfo, error)
	SaveUploadImages(req *RequestUploadFile) *ApiResponse[ResponseUploadFile]
}

type RequestUploadFile struct {
	JwtHeader
	EchoContentHeader
	File *multipart.FileHeader
}

type ResponseUploadFile struct {
	FileSize   int64  `json:"file_size"`
	AccessPath string `json:"access_path"`
}
