// Package log
package log

import (
	"github.com/half-nothing/simple-fsd/internal/interfaces/global"
	"log/slog"
)

type Loggers struct {
	mainLogger LoggerInterface
	fsdLogger  LoggerInterface
	httpLogger LoggerInterface
	grpcLogger LoggerInterface
}

type LoggerInterface interface {
	Init(logPath, logName string, debug, noLogs bool)
	ShutdownCallback() global.Callable
	LogHandler() *slog.Logger
	Debug(msg string, v ...interface{})
	DebugF(msg string, v ...interface{})
	Info(msg string, v ...interface{})
	InfoF(msg string, v ...interface{})
	Warn(msg string, v ...interface{})
	WarnF(msg string, v ...interface{})
	Error(msg string, v ...interface{})
	ErrorF(msg string, v ...interface{})
	Fatal(msg string, v ...interface{})
	FatalF(msg string, v ...interface{})
}

func NewLoggers(
	mainLogger LoggerInterface,
	fsdLogger LoggerInterface,
	httpLogger LoggerInterface,
	grpcLogger LoggerInterface,
) *Loggers {
	return &Loggers{
		mainLogger: mainLogger,
		fsdLogger:  fsdLogger,
		httpLogger: httpLogger,
		grpcLogger: grpcLogger,
	}
}

func (logger *Loggers) MainLogger() LoggerInterface { return logger.mainLogger }

func (logger *Loggers) FsdLogger() LoggerInterface { return logger.fsdLogger }

func (logger *Loggers) HttpLogger() LoggerInterface { return logger.httpLogger }

func (logger *Loggers) GrpcLogger() LoggerInterface { return logger.grpcLogger }
