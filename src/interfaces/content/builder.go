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
	content *ApplicationContent
}

// NewApplicationContentBuilder 创建一个新的ApplicationContentBuilder实例
func NewApplicationContentBuilder() *ApplicationContentBuilder {
	return &ApplicationContentBuilder{
		content: &ApplicationContent{},
	}
}

// SetConfigManager 设置配置管理器
func (builder *ApplicationContentBuilder) SetConfigManager(configManager config.ManagerInterface) *ApplicationContentBuilder {
	builder.content.configManager = configManager
	return builder
}

// SetCleaner 设置清理器
func (builder *ApplicationContentBuilder) SetCleaner(cleaner cleaner.Interface) *ApplicationContentBuilder {
	builder.content.cleaner = cleaner
	return builder
}

// SetClientManager 设置客户端管理器
func (builder *ApplicationContentBuilder) SetClientManager(clientManager fsd.ClientManagerInterface) *ApplicationContentBuilder {
	builder.content.clientManager = clientManager
	return builder
}

// SetConnectionManager 设置连接管理器
func (builder *ApplicationContentBuilder) SetConnectionManager(connectionManager fsd.ConnectionManagerInterface) *ApplicationContentBuilder {
	builder.content.connectionManager = connectionManager
	return builder
}

// SetLogger 设置日志处理器
func (builder *ApplicationContentBuilder) SetLogger(logger logger.HandlerInterface) *ApplicationContentBuilder {
	builder.content.logger = logger
	return builder
}

// SetMessageQueue 设置消息队列
func (builder *ApplicationContentBuilder) SetMessageQueue(messageQueue queue.MessageQueueInterface) *ApplicationContentBuilder {
	builder.content.messageQueue = messageQueue
	return builder
}

// SetMetarManager 设置METAR管理器
func (builder *ApplicationContentBuilder) SetMetarManager(metarManager metar.ManagerInterface) *ApplicationContentBuilder {
	builder.content.metarManager = metarManager
	return builder
}

// SetOperations 设置数据库操作接口
func (builder *ApplicationContentBuilder) SetOperations(operations repository.DatabaseInterface) *ApplicationContentBuilder {
	builder.content.operations = operations
	return builder
}

// Build 构建并返回ApplicationContent实例
func (builder *ApplicationContentBuilder) Build() *ApplicationContent {
	return builder.content
}
