// Package log
package log

import (
	"log/slog"

	"github.com/half-nothing/simple-fsd/internal/interfaces/global"
)

type Loggers struct {
	mainLogger  LoggerInterface
	fsdLogger   LoggerInterface
	httpLogger  LoggerInterface
	grpcLogger  LoggerInterface
	voiceLogger LoggerInterface
}

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

func NewLoggers(
	mainLogger LoggerInterface,
	fsdLogger LoggerInterface,
	httpLogger LoggerInterface,
	grpcLogger LoggerInterface,
	voiceLogger LoggerInterface,
) *Loggers {
	return &Loggers{
		mainLogger:  mainLogger,
		fsdLogger:   fsdLogger,
		httpLogger:  httpLogger,
		grpcLogger:  grpcLogger,
		voiceLogger: voiceLogger,
	}
}

func (logger *Loggers) MainLogger() LoggerInterface { return logger.mainLogger }

func (logger *Loggers) FsdLogger() LoggerInterface { return logger.fsdLogger }

func (logger *Loggers) HttpLogger() LoggerInterface { return logger.httpLogger }

func (logger *Loggers) GrpcLogger() LoggerInterface { return logger.grpcLogger }

func (logger *Loggers) VoiceLogger() LoggerInterface { return logger.voiceLogger }
