// Package content
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

// ApplicationContent 应用程序上下文结构体，包含所有核心组件的接口
type ApplicationContent struct {
	configManager     config.ManagerInterface        // 配置管理器
	cleaner           cleaner.Interface              // 清理器
	clientManager     fsd.ClientManagerInterface     // 客户端管理器
	connectionManager fsd.ConnectionManagerInterface // 连接管理器
	logger            logger.HandlerInterface        // 日志处理器
	messageQueue      queue.MessageQueueInterface    // 消息队列
	metarManager      metar.ManagerInterface         // METAR气象数据管理器
	operations        repository.DatabaseInterface   // 数据库操作接口
}

// ConfigManager 获取配置管理器实例
func (app *ApplicationContent) ConfigManager() config.ManagerInterface {
	return app.configManager
}

// Cleaner 获取清理器实例
func (app *ApplicationContent) Cleaner() cleaner.Interface { return app.cleaner }

// ClientManager 获取客户端管理器实例
func (app *ApplicationContent) ClientManager() fsd.ClientManagerInterface { return app.clientManager }

// ConnectionManager 获取连接管理器实例
func (app *ApplicationContent) ConnectionManager() fsd.ConnectionManagerInterface {
	return app.connectionManager
}

// Logger 获取日志处理器实例
func (app *ApplicationContent) Logger() logger.HandlerInterface { return app.logger }

// MessageQueue 获取消息队列实例
func (app *ApplicationContent) MessageQueue() queue.MessageQueueInterface { return app.messageQueue }

// MetarManager 获取METAR气象数据管理器实例
func (app *ApplicationContent) MetarManager() metar.ManagerInterface { return app.metarManager }

// Operations 获取数据库操作接口实例
func (app *ApplicationContent) Operations() repository.DatabaseInterface { return app.operations }
