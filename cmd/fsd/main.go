package main

import (
	"flag"
	"fmt"
	"github.com/half-nothing/simple-fsd/internal/base"
	"github.com/half-nothing/simple-fsd/internal/database"
	"github.com/half-nothing/simple-fsd/internal/fsd_server"
	"github.com/half-nothing/simple-fsd/internal/interfaces"
	"github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/internal/interfaces/global"
)

func recoverFromError() {
	if r := recover(); r != nil {
		fmt.Printf("It looks like there are some serious errors, the details are as follows: %v", r)
	}
}

func main() {
	flag.Parse()

	defer recoverFromError()

	logger := base.NewLogger()
	logger.Init(*global.DebugMode)

	logger.Info("Application initializing...")

	cleaner := base.NewCleaner(logger)
	cleaner.Init()
	defer cleaner.Clean()

	configManager := base.NewManager(logger)
	config := configManager.Config()

	if err := fsd.SyncRatingConfig(config); err != nil {
		logger.FatalF("Error occurred while handle rating base, details: %v", err)
		return
	}

	shutdownCallback, databaseOperation, err := database.ConnectDatabase(logger, config, *global.DebugMode)
	if err != nil {
		logger.FatalF("Error occurred while initializing operation, details: %v", err)
		return
	}

	cleaner.Add(shutdownCallback)

	applicationContent := interfaces.NewApplicationContent(configManager, cleaner, logger, databaseOperation)

	fsd_server.StartFSDServer(applicationContent)
}
