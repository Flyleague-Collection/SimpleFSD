// Package global
package global

import (
	"flag"
	"time"
)

var (
	DebugMode               = flag.Bool("debug", false, "Enable debug mode")
	ConfigFilePath          = flag.String("config", "./config.json", "Path to configuration file")
	SkipEmailVerification   = flag.Bool("skip_email_verification", false, "Skip email verification")
	UpdateConfig            = flag.Bool("update_config", false, "Auto update configuration")
	NoLogs                  = flag.Bool("no_logs", false, "Disable logging to file")
	MessageQueueChannelSize = flag.Int("message_queue_channel_size", 128, "Message Queue channel size")
	DownloadPrefix          = flag.String("download_prefix", "https://raw.githubusercontent.com/Flyleague-Collection/SimpleFSD/refs/heads/main", "auto download prefix")
	MetarCacheCleanInterval = flag.Duration("metar_cache_clean_interval", 30*time.Minute, "metar cache cleanup interval")
	MetarQueryThread        = flag.Int("metar_query_thread", 32, "metar query thread")
)

const (
	AppVersion    = "0.8.0"
	ConfigVersion = "0.8.0"

	EnvDebugMode               = "DEBUG_MODE"
	EnvConfigFilePath          = "CONFIG_FILE_PATH"
	EnvSkipEmailVerification   = "SKIP_EMAIL_VERIFICATION"
	EnvUpdateConfig            = "UPDATE_CONFIG"
	EnvNoLogs                  = "NO_LOGS"
	EnvMessageQueueChannelSize = "MESSAGE_QUEUE_CHANNEL_SIZE"
	EnvDownloadPrefix          = "DOWNLOAD_PREFIX"
	EnvMetarCacheCleanInterval = "METAR_CACHE_CLEAN_INTERVAL"
	EnvMetarQueryThread        = "METAR_QUERY_THREAD"

	LogFilePath = "logs"
	MainLogName = "main"
	MainLogPath = LogFilePath + "/" + MainLogName + ".log"
	FsdLogName  = "fsd"
	FsdLogPath  = LogFilePath + "/" + FsdLogName + ".log"
	HttpLogName = "http"
	HttpLogPath = LogFilePath + "/" + HttpLogName + ".log"
	GrpcLogName = "grpc"
	GrpcLogPath = LogFilePath + "/" + GrpcLogName + ".log"

	AirportDataFilePath              = "/data/airport.json"
	EmailVerifyTemplateFilePath      = "/template/email_verify.template"
	ATCRatingChangeTemplateFilePath  = "/template/atc_rating_change.template"
	PermissionChangeTemplateFilePath = "/template/permission_change.template"
	KickedFromServerTemplateFilePath = "/template/kicked_from_server.template"
	PasswordChangeTemplateFilePath   = "/template/password_change.template"

	DefaultFilePermissions     = 0644
	DefaultDirectoryPermission = 0755

	FSDServerName      = "SERVER"
	FSDDisconnectDelay = 100 * time.Millisecond
)
