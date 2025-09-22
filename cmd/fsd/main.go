package main

import (
	"flag"
	"fmt"
	"github.com/half-nothing/simple-fsd/internal/base"
	"github.com/half-nothing/simple-fsd/internal/database"
	"github.com/half-nothing/simple-fsd/internal/fsd_server"
	"github.com/half-nothing/simple-fsd/internal/fsd_server/packet"
	"github.com/half-nothing/simple-fsd/internal/interfaces"
	"github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/internal/interfaces/global"
	"os"
	"strconv"
)

func recoverFromError() {
	if r := recover(); r != nil {
		fmt.Printf("It looks like there are some serious errors, the details are as follows: %v", r)
	}
}

func checkStringEnv(envKey string, target *string) {
	value := os.Getenv(envKey)
	if value != "" {
		*target = value
	}
}

func checkIntEnv(envKey string, target *int64, defaultValue int64) {
	value := os.Getenv(envKey)
	if value != "" {
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			*target = defaultValue
		} else {
			*target = val
		}
	}
}

func checkBoolEnv(envKey string, target *bool) {
	value := os.Getenv(envKey)
	if val, err := strconv.ParseBool(value); err == nil && val {
		*target = true
	}
}

func main() {
	flag.Parse()

	checkBoolEnv(global.EnvDebugMode, global.DebugMode)
	checkStringEnv(global.EnvConfigFilePath, global.ConfigFilePath)
	checkIntEnv(global.EnvFlushInterval, global.FlushInterval, 5)
	checkBoolEnv(global.EnvNoLogs, global.NoLogs)

	defer recoverFromError()

	logger := base.NewLogger()
	logger.Init(global.MainLogPath, global.MainLogName, *global.DebugMode, *global.NoLogs)

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

	if err := fsd.SyncFacilityConfig(config); err != nil {
		logger.FatalF("Error occurred while handle facility addition, details: %v", err)
		return
	}

	fsd.SyncRangeLimit(config.RangeLimit)

	shutdownCallback, databaseOperation, err := database.ConnectDatabase(logger, config, *global.DebugMode)
	if err != nil {
		logger.FatalF("Error occurred while initializing operation, details: %v", err)
		return
	}

	cleaner.Add(shutdownCallback)

	clientManager := packet.NewClientManager(logger, config)

	applicationContent := interfaces.NewApplicationContent(logger, cleaner, configManager, clientManager, databaseOperation)

	fsd_server.StartFSDServer(applicationContent)
}
