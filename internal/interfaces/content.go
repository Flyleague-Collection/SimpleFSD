// Package interfaces
package interfaces

import (
	"github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"github.com/half-nothing/simple-fsd/internal/interfaces/queue"
)

type ApplicationContent struct {
	configManager ConfigManagerInterface
	cleaner       CleanerInterface
	clientManager fsd.ClientManagerInterface
	logger        *log.Loggers
	messageQueue  queue.MessageQueueInterface
	operations    *operation.DatabaseOperations
}

func NewApplicationContent(
	logger *log.Loggers,
	cleaner CleanerInterface,
	configManager ConfigManagerInterface,
	clientManager fsd.ClientManagerInterface,
	messageQueue queue.MessageQueueInterface,
	db *operation.DatabaseOperations,
) *ApplicationContent {
	return &ApplicationContent{
		configManager: configManager,
		cleaner:       cleaner,
		clientManager: clientManager,
		logger:        logger,
		messageQueue:  messageQueue,
		operations:    db,
	}
}

func (app *ApplicationContent) ConfigManager() ConfigManagerInterface {
	return app.configManager
}

func (app *ApplicationContent) Cleaner() CleanerInterface { return app.cleaner }

func (app *ApplicationContent) ClientManager() fsd.ClientManagerInterface { return app.clientManager }

func (app *ApplicationContent) Logger() *log.Loggers { return app.logger }

func (app *ApplicationContent) MessageQueue() queue.MessageQueueInterface { return app.messageQueue }

func (app *ApplicationContent) Operations() *operation.DatabaseOperations { return app.operations }
