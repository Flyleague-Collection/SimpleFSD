// Package store
// 存放 service.StoreServiceInterface 的实现
// 本文件存放阿里云OSS存储实现
package store

import (
	"context"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/half-nothing/simple-fsd/internal/interfaces/config"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"github.com/half-nothing/simple-fsd/internal/interfaces/queue"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/service"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"strings"
)

type ALiYunOssStoreService struct {
	logger            log.LoggerInterface
	localStore        StoreServiceInterface
	config            *config.HttpServerStore
	endpoint          *url.URL
	client            *oss.Client
	messageQueue      queue.MessageQueueInterface
	auditLogOperation operation.AuditLogOperationInterface
}

func NewALiYunOssStoreService(
	logger log.LoggerInterface,
	config *config.HttpServerStore,
	localStore StoreServiceInterface,
	messageQueue queue.MessageQueueInterface,
	auditLogOperation operation.AuditLogOperationInterface,
) *ALiYunOssStoreService {
	service := &ALiYunOssStoreService{
		logger:            log.NewLoggerAdapter(logger, "ALiYunOssStoreService"),
		localStore:        localStore,
		config:            config,
		messageQueue:      messageQueue,
		auditLogOperation: auditLogOperation,
	}
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(config.AccessId, config.AccessKey)).
		WithRegion(config.Region).
		WithUseInternalEndpoint(config.UseInternalUrl)
	service.client = oss.NewClient(cfg)
	if config.CdnDomain != "" {
		service.endpoint, _ = url.Parse(config.CdnDomain)
	} else {
		service.endpoint, _ = url.Parse(strings.Replace(*cfg.Endpoint, "-internal", "", 1))
	}
	return service
}

func (store *ALiYunOssStoreService) SaveImageFile(file *multipart.FileHeader) (*StoreInfo, *ApiStatus) {
	storeInfo, res := store.localStore.SaveImageFile(file)
	if res != nil {
		return nil, res
	}

	storeInfo.RemotePath = strings.Replace(filepath.Join(store.config.RemoteStorePath, storeInfo.FileName), "\\", "/", -1)

	reader, err := file.Open()
	if err != nil {
		store.logger.ErrorF("ALiYunOssStoreService.SaveImageFile open form file error: %v", err)
		return nil, ErrFileUploadFail
	}

	putRequest := &oss.PutObjectRequest{
		Bucket:       oss.Ptr(store.config.Bucket),
		Key:          oss.Ptr(storeInfo.RemotePath),
		StorageClass: oss.StorageClassStandard,
		Body:         reader,
	}

	_, err = store.client.PutObject(context.TODO(), putRequest)
	if err != nil {
		store.logger.ErrorF("ALiYunOssStoreService.SaveImageFile upload image to remote storage error: %v", err)
		return nil, ErrFileUploadFail
	}
	return storeInfo, nil
}

func (store *ALiYunOssStoreService) DeleteImageFile(file string) (*StoreInfo, error) {
	storeInfo, err := store.localStore.DeleteImageFile(file)
	if err != nil {
		return nil, err
	}

	storeInfo.RemotePath = strings.Replace(filepath.Join(store.config.RemoteStorePath, storeInfo.FileName), "\\", "/", -1)

	delRequest := &oss.DeleteObjectRequest{
		Bucket: oss.Ptr(store.config.Bucket),
		Key:    oss.Ptr(storeInfo.RemotePath),
	}

	_, err = store.client.DeleteObject(context.TODO(), delRequest)
	if err != nil {
		store.logger.ErrorF("ALiYunOssStoreService.DeleteImageFile delete image from remote storage error: %v", err)
		return nil, err
	}
	return storeInfo, nil
}

func (store *ALiYunOssStoreService) SaveUploadImages(req *RequestUploadFile) *ApiResponse[ResponseUploadFile] {
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
