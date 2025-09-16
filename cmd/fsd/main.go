package main

import (
	"flag"
	"fmt"
	"github.com/half-nothing/simple-fsd/internal/base"
	"github.com/half-nothing/simple-fsd/internal/database"
	"github.com/half-nothing/simple-fsd/internal/fsd_server"
	"github.com/half-nothing/simple-fsd/internal/fsd_server/packet"
	"github.com/half-nothing/simple-fsd/internal/http_server"
	"github.com/half-nothing/simple-fsd/internal/interfaces"
	"github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/internal/interfaces/global"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"github.com/half-nothing/simple-fsd/internal/interfaces/queue"
	"github.com/half-nothing/simple-fsd/internal/message"
	"github.com/half-nothing/simple-fsd/internal/utils"
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

func checkIntEnv(envKey string, target *int, defaultValue int) {
	value := os.Getenv(envKey)
	if value != "" {
		*target = utils.StrToInt(value, defaultValue)
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
	checkBoolEnv(global.EnvSkipEmailVerification, global.SkipEmailVerification)
	checkBoolEnv(global.EnvUpdateConfig, global.UpdateConfig)
	checkBoolEnv(global.EnvNoLogs, global.NoLogs)
	checkIntEnv(global.EnvMessageQueueChannelSize, global.MessageQueueChannelSize, 128)
	checkStringEnv(global.EnvDownloadPrefix, global.DownloadPrefix)

	defer recoverFromError()

	mainLogger := base.NewLogger()
	mainLogger.Init(global.MainLogPath, global.MainLogName, *global.DebugMode, *global.NoLogs)

	fsdLogger := base.NewLogger()
	fsdLogger.Init(global.FsdLogPath, global.FsdLogName, *global.DebugMode, *global.NoLogs)

	httpLogger := base.NewLogger()
	httpLogger.Init(global.HttpLogPath, global.HttpLogName, *global.DebugMode, *global.NoLogs)

	grpcLogger := base.NewLogger()
	grpcLogger.Init(global.GrpcLogPath, global.GrpcLogName, *global.DebugMode, *global.NoLogs)

	logger := log.NewLoggers(mainLogger, fsdLogger, httpLogger, grpcLogger)

	mainLogger.Info("Application initializing...")

	mainLogger.Info("Reading configuration...")
	configManager := base.NewManager(mainLogger)
	config := configManager.Config()

	mainLogger.Info("Creating cleaner...")
	cleaner := base.NewCleaner(mainLogger)
	cleaner.Init()
	defer cleaner.Clean()

	cleaner.Add(fsdLogger.ShutdownCallback())
	cleaner.Add(httpLogger.ShutdownCallback())
	cleaner.Add(grpcLogger.ShutdownCallback())

	if err := fsd.SyncRatingConfig(config); err != nil {
		mainLogger.FatalF("Error occurred while handle rating base, details: %v", err)
		return
	}

	mainLogger.Info("Connecting to database...")
	shutdownCallback, databaseOperation, err := database.ConnectDatabase(mainLogger, config, *global.DebugMode)
	if err != nil {
		mainLogger.FatalF("Error occurred while initializing operation, details: %v", err)
		return
	}

	cleaner.Add(shutdownCallback)

	mainLogger.InfoF("Initialize message queue with channel size %d", *global.MessageQueueChannelSize)
	messageQueue := message.NewAsyncMessageQueue(mainLogger, *global.MessageQueueChannelSize)

	cleaner.Add(messageQueue.ShutdownCallback())

	clientManager := packet.NewClientManager(fsdLogger, config)

	messageQueue.Subscribe(queue.KickClientFromServer, clientManager.HandleKickClientFromServerMessage)
	messageQueue.Subscribe(queue.SendMessageToClient, clientManager.HandleSendMessageToClientMessage)

	mainLogger.Info("Creating application content...")
	applicationContent := interfaces.NewApplicationContent(logger, cleaner, configManager, clientManager, messageQueue, databaseOperation)

	mainLogger.Info("Application initialized. Starting application...")

	if config.Server.HttpServer.Enabled {
		go http_server.StartHttpServer(applicationContent)
	}

	//if config.Server.GRPCServer.Enabled {
	//	go grpc_server.StartGRPCServer(applicationContent)
	//}

	fsd_server.StartFSDServer(applicationContent)
}
