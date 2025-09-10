// Package interfaces
package interfaces

import (
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
)

type ApplicationContent struct {
	configManager ConfigManagerInterface
	cleaner       CleanerInterface
	logger        log.LoggerInterface
	operations    *operation.DatabaseOperations
}

func NewApplicationContent(
	configManager ConfigManagerInterface,
	cleaner CleanerInterface,
	logger log.LoggerInterface,
	db *operation.DatabaseOperations,
) *ApplicationContent {
	return &ApplicationContent{
		configManager: configManager,
		cleaner:       cleaner,
		logger:        logger,
		operations:    db,
	}
}

func (app *ApplicationContent) ConfigManager() ConfigManagerInterface {
	return app.configManager
}

func (app *ApplicationContent) Cleaner() CleanerInterface { return app.cleaner }

func (app *ApplicationContent) Logger() log.LoggerInterface { return app.logger }

func (app *ApplicationContent) Operations() *operation.DatabaseOperations { return app.operations }
