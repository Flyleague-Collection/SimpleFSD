// Package interfaces
package interfaces

import (
	"github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"github.com/half-nothing/simple-fsd/internal/interfaces/queue"
)

type ApplicationContent struct {
	configManager     ConfigManagerInterface
	cleaner           CleanerInterface
	clientManager     fsd.ClientManagerInterface
	connectionManager fsd.ConnectionManagerInterface
	logger            *log.Loggers
	messageQueue      queue.MessageQueueInterface
	metarManager      MetarManagerInterface
	operations        *operation.DatabaseOperations
}

func NewApplicationContent(
	logger *log.Loggers,
	cleaner CleanerInterface,
	configManager ConfigManagerInterface,
	clientManager fsd.ClientManagerInterface,
	connectionManager fsd.ConnectionManagerInterface,
	messageQueue queue.MessageQueueInterface,
	metarManager MetarManagerInterface,
	db *operation.DatabaseOperations,
) *ApplicationContent {
	return &ApplicationContent{
		configManager:     configManager,
		cleaner:           cleaner,
		clientManager:     clientManager,
		connectionManager: connectionManager,
		logger:            logger,
		messageQueue:      messageQueue,
		metarManager:      metarManager,
		operations:        db,
	}
}

func (app *ApplicationContent) ConfigManager() ConfigManagerInterface {
	return app.configManager
}

func (app *ApplicationContent) Cleaner() CleanerInterface { return app.cleaner }

func (app *ApplicationContent) ClientManager() fsd.ClientManagerInterface { return app.clientManager }

func (app *ApplicationContent) ConnectionManager() fsd.ConnectionManagerInterface {
	return app.connectionManager
}

func (app *ApplicationContent) Logger() *log.Loggers { return app.logger }

func (app *ApplicationContent) MessageQueue() queue.MessageQueueInterface { return app.messageQueue }

func (app *ApplicationContent) MetarManager() MetarManagerInterface { return app.metarManager }

func (app *ApplicationContent) Operations() *operation.DatabaseOperations { return app.operations }
