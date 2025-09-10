// Package global
package global

import (
	"flag"
	"time"
)

var (
	DebugMode             = flag.Bool("debug", false, "Enable debug mode")
	ConfigFilePath        = flag.String("config", "./config.json", "Path to configuration file")
	SkipEmailVerification = flag.Bool("skip_email_verification", false, "Skip email verification")
)

const (
	AppVersion    = "0.6.0"
	ConfigVersion = "0.6.0"

	DefaultFilePermissions     = 0644
	DefaultDirectoryPermission = 0755

	AirportDataFileUrl              = "https://raw.githubusercontent.com/Flyleague-Collection/SimpleFSD/refs/heads/main/data/airport.json"
	EmailVerifyTemplateFileUrl      = "https://raw.githubusercontent.com/Flyleague-Collection/SimpleFSD/refs/heads/main/template/email_verify.template"
	ATCRatingChangeTemplateFileUrl  = "https://raw.githubusercontent.com/Flyleague-Collection/SimpleFSD/refs/heads/main/template/atc_rating_change.template"
	PermissionChangeTemplateFileUrl = "https://raw.githubusercontent.com/Flyleague-Collection/SimpleFSD/refs/heads/main/template/permission_change.template"
	KickedFromServerTemplateFileUrl = "https://raw.githubusercontent.com/Flyleague-Collection/SimpleFSD/refs/heads/main/template/kicked_from_server.template"

	FSDServerName      = "SERVER"
	FSDDisconnectDelay = 100 * time.Millisecond
)
