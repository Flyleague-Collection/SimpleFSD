// Package log
package log

import (
	"github.com/half-nothing/simple-fsd/internal/interfaces/global"
	"log/slog"
)

type LoggerInterface interface {
	Init(logPath, logName string, debug, noLogs bool)
	ShutdownCallback() global.Callable
	LogHandler() *slog.Logger
	Debug(msg string)
	DebugF(msg string, v ...interface{})
	Info(msg string)
	InfoF(msg string, v ...interface{})
	Warn(msg string)
	WarnF(msg string, v ...interface{})
	Error(msg string)
	ErrorF(msg string, v ...interface{})
	Fatal(msg string)
	FatalF(msg string, v ...interface{})
}
