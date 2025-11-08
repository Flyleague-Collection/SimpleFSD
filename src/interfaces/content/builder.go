// Package content 提供应用内容构建相关功能
package content

import (
	"github.com/half-nothing/simple-fsd/src/interfaces/cleaner"
	"github.com/half-nothing/simple-fsd/src/interfaces/config"
	"github.com/half-nothing/simple-fsd/src/interfaces/database/repository"
	"github.com/half-nothing/simple-fsd/src/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/src/interfaces/logger"
	"github.com/half-nothing/simple-fsd/src/interfaces/metar"
	"github.com/half-nothing/simple-fsd/src/interfaces/queue"
)

// ApplicationContentBuilder 应用内容构建器结构体，用于构建ApplicationContent实例
type ApplicationContentBuilder struct {
	configManager     config.ManagerInterface
	cleaner           cleaner.Interface
	clientManager     fsd.ClientManagerInterface
	connectionManager fsd.ConnectionManagerInterface
	logger            logger.HandlerInterface
	messageQueue      queue.MessageQueueInterface
	metarManager      metar.ManagerInterface
	operations        repository.DatabaseInterface
}

// NewApplicationContentBuilder 创建一个新的ApplicationContentBuilder实例
func NewApplicationContentBuilder() *ApplicationContentBuilder {
	return &ApplicationContentBuilder{}
}

// SetConfigManager 设置配置管理器
func (builder *ApplicationContentBuilder) SetConfigManager(configManager config.ManagerInterface) *ApplicationContentBuilder {
	builder.configManager = configManager
	return builder
}

// SetCleaner 设置清理器
func (builder *ApplicationContentBuilder) SetCleaner(cleaner cleaner.Interface) *ApplicationContentBuilder {
	builder.cleaner = cleaner
	return builder
}

// SetClientManager 设置客户端管理器
func (builder *ApplicationContentBuilder) SetClientManager(clientManager fsd.ClientManagerInterface) *ApplicationContentBuilder {
	builder.clientManager = clientManager
	return builder
}

// SetConnectionManager 设置连接管理器
func (builder *ApplicationContentBuilder) SetConnectionManager(connectionManager fsd.ConnectionManagerInterface) *ApplicationContentBuilder {
	builder.connectionManager = connectionManager
	return builder
}

// SetLogger 设置日志处理器
func (builder *ApplicationContentBuilder) SetLogger(logger logger.HandlerInterface) *ApplicationContentBuilder {
	builder.logger = logger
	return builder
}

// SetMessageQueue 设置消息队列
func (builder *ApplicationContentBuilder) SetMessageQueue(messageQueue queue.MessageQueueInterface) *ApplicationContentBuilder {
	builder.messageQueue = messageQueue
	return builder
}

// SetMetarManager 设置METAR管理器
func (builder *ApplicationContentBuilder) SetMetarManager(metarManager metar.ManagerInterface) *ApplicationContentBuilder {
	builder.metarManager = metarManager
	return builder
}

// SetOperations 设置数据库操作接口
func (builder *ApplicationContentBuilder) SetOperations(operations repository.DatabaseInterface) *ApplicationContentBuilder {
	builder.operations = operations
	return builder
}

// Build 构建并返回ApplicationContent实例
func (builder *ApplicationContentBuilder) Build() *ApplicationContent {
	return &ApplicationContent{
		configManager:     builder.configManager,
		cleaner:           builder.cleaner,
		clientManager:     builder.clientManager,
		connectionManager: builder.connectionManager,
		logger:            builder.logger,
		messageQueue:      builder.messageQueue,
		metarManager:      builder.metarManager,
		operations:        builder.operations,
	}
}
