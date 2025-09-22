// Package global
package global

import (
	"flag"
	"time"
)

var (
	DebugMode      = flag.Bool("debug", false, "Enable debug mode")
	ConfigFilePath = flag.String("config", "./config.json", "Path to configuration file")
	NoLogs         = flag.Bool("no_logs", false, "Disable logging to file")
	FlushInterval  = flag.Int64("flush_interval", 5, "Flush interval")
)

const (
	AppVersion    = "0.7.0"
	ConfigVersion = "0.7.0"

	EnvDebugMode      = "DEBUG_MODE"
	EnvConfigFilePath = "CONFIG_FILE_PATH"
	EnvFlushInterval  = "FLUSH_INTERVAL"
	EnvNoLogs         = "NO_LOGS"

	LogFilePath = "logs"
	MainLogName = "main"
	MainLogPath = LogFilePath + "/" + MainLogName + ".log"

	DefaultFilePermissions     = 0644
	DefaultDirectoryPermission = 0755

	FSDServerName      = "SERVER"
	FSDDisconnectDelay = 100 * time.Millisecond
)
