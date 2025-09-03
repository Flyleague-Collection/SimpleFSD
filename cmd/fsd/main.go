package main

import (
	c "github.com/half-nothing/simple-fsd/internal/config"
	"github.com/half-nothing/simple-fsd/internal/database"
	"github.com/half-nothing/simple-fsd/internal/fsd_server"
	"github.com/half-nothing/simple-fsd/internal/interfaces"
	"github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
)

func main() {
	config := c.GetConfig()
	if err := fsd.SyncRatingConfig(config); err != nil {
		c.FatalF("Error occurred while handle rating config, details: %v", err)
		return
	}
	loggerCallback := c.Init(config)
	c.Info("Application initializing...")
	cleaner := c.GetCleaner()
	cleaner.Init(loggerCallback)
	defer cleaner.Clean()
	databaseOperation, err := database.ConnectDatabase(config)
	if err != nil {
		c.FatalF("Error occurred while initializing operation, details: %v", err)
		return
	}
	applicationContent := interfaces.NewApplicationContent(config, databaseOperation)
	fsd_server.StartFSDServer(applicationContent)
}
