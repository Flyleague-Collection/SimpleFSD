// Package interfaces
package interfaces

import (
	"github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
)

type ApplicationContent struct {
	configManager ConfigManagerInterface
	cleaner       CleanerInterface
	clientManager fsd.ClientManagerInterface
	logger        log.LoggerInterface
	operations    *operation.DatabaseOperations
}

func NewApplicationContent(
	logger log.LoggerInterface,
	cleaner CleanerInterface,
	configManager ConfigManagerInterface,
	clientManager fsd.ClientManagerInterface,
	db *operation.DatabaseOperations,
) *ApplicationContent {
	return &ApplicationContent{
		configManager: configManager,
		cleaner:       cleaner,
		clientManager: clientManager,
		logger:        logger,
		operations:    db,
	}
}

func (app *ApplicationContent) ConfigManager() ConfigManagerInterface {
	return app.configManager
}

func (app *ApplicationContent) Cleaner() CleanerInterface { return app.cleaner }

func (app *ApplicationContent) ClientManager() fsd.ClientManagerInterface { return app.clientManager }

func (app *ApplicationContent) Logger() log.LoggerInterface { return app.logger }

func (app *ApplicationContent) Operations() *operation.DatabaseOperations { return app.operations }
