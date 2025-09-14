// Package global
package global

import (
	"flag"
	"time"
)

var (
	DebugMode      = flag.Bool("debug", false, "Enable debug mode")
	ConfigFilePath = flag.String("config", "./config.json", "Path to configuration file")
	FlushInterval  = flag.Int64("flush_interval", 5, "Flush interval")
)

const (
	AppVersion    = "0.6.0"
	ConfigVersion = "0.6.0"

	DefaultFilePermissions     = 0644
	DefaultDirectoryPermission = 0755

	FSDServerName      = "SERVER"
	FSDDisconnectDelay = 100 * time.Millisecond
)
