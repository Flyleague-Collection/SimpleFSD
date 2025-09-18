// Package store
// 存放 service.StoreServiceInterface 的实现
// 本文件存放腾讯云COS存储实现
package store

import (
	"context"
	"fmt"
	"github.com/half-nothing/simple-fsd/internal/interfaces/config"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"github.com/half-nothing/simple-fsd/internal/interfaces/queue"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/service"
	"github.com/tencentyun/cos-go-sdk-v5"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

type TencentCosStoreService struct {
	logger            log.LoggerInterface
	localStore        StoreServiceInterface
	config            *config.HttpServerStore
	endpoint          *url.URL
	client            *cos.Client
	messageQueue      queue.MessageQueueInterface
	auditLogOperation operation.AuditLogOperationInterface
}

func NewTencentCosStoreService(
	logger log.LoggerInterface,
	config *config.HttpServerStore,
	localStore StoreServiceInterface,
	messageQueue queue.MessageQueueInterface,
	auditLogOperation operation.AuditLogOperationInterface,
) *TencentCosStoreService {
	service := &TencentCosStoreService{
		logger:            logger,
		localStore:        localStore,
		config:            config,
		messageQueue:      messageQueue,
		auditLogOperation: auditLogOperation,
	}
	bucketUrl, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", config.Bucket, strings.ToLower(config.Region)))
	serviceUrl, _ := url.Parse(fmt.Sprintf("https://cos.%s.myqcloud.com", strings.ToLower(config.Region)))
	baseUrl := &cos.BaseURL{BucketURL: bucketUrl, ServiceURL: serviceUrl}
	service.client = cos.NewClient(baseUrl, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.AccessId,
			SecretKey: config.AccessKey,
		},
	})
	if config.CdnDomain != "" {
		service.endpoint, _ = url.Parse(config.CdnDomain)
	} else {
		service.endpoint = service.client.BaseURL.BucketURL
	}
	return service
}

func (store *TencentCosStoreService) SaveImageFile(file *multipart.FileHeader) (*StoreInfo, *ApiStatus) {
	storeInfo, res := store.localStore.SaveImageFile(file)
	if res != nil {
		return nil, res
	}

	storeInfo.RemotePath = strings.Replace(filepath.Join(store.config.RemoteStorePath, storeInfo.FileName), "\\", "/", -1)

	reader, err := file.Open()
	if err != nil {
		store.logger.ErrorF("TencentCosStoreService.SaveImageFile open form file errors: %v", err)
		return nil, ErrFileUploadFail
	}

	_, err = store.client.Object.Put(context.Background(), storeInfo.RemotePath, reader, nil)
	if err != nil {
		store.logger.ErrorF("TencentCosStoreService.SaveImageFile upload image to remote storage error: %v", err)
		return nil, ErrFileUploadFail
	}
	return storeInfo, nil
}

func (store *TencentCosStoreService) DeleteImageFile(file string) (*StoreInfo, error) {
	storeInfo, err := store.localStore.DeleteImageFile(file)
	if err != nil {
		return nil, err
	}

	storeInfo.RemotePath = strings.Replace(filepath.Join(store.config.RemoteStorePath, storeInfo.FileName), "\\", "/", -1)

	_, err = store.client.Object.Delete(context.Background(), storeInfo.RemotePath)
	if err != nil {
		store.logger.ErrorF("TencentCosStoreService.DeleteImageFile delete image from remote storage errors: %v", err)
		return nil, err
	}
	return storeInfo, nil
}

func (store *TencentCosStoreService) SaveUploadImages(req *RequestUploadFile) *ApiResponse[ResponseUploadFile] {
	storeInfo, res := store.SaveImageFile(req.File)
	if res != nil {
		return NewApiResponse[ResponseUploadFile](res, nil)
	}

	accessUrl, err := url.JoinPath(store.endpoint.String(), storeInfo.RemotePath)
	if err != nil {
		return NewApiResponse[ResponseUploadFile](ErrFilePathFail, nil)
	}

	store.messageQueue.Publish(&queue.Message{
		Type: queue.AuditLog,
		Data: store.auditLogOperation.NewAuditLog(
			operation.FileUpload,
			req.Cid,
			storeInfo.RemotePath,
			req.Ip,
			req.UserAgent,
			nil,
		),
	})

	return NewApiResponse(SuccessUploadFile, &ResponseUploadFile{
		FileSize:   req.File.Size,
		AccessPath: accessUrl,
	})
}
