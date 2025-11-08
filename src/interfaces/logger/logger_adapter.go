// Package logger
package logger

import (
	"fmt"
	"log/slog"

	"github.com/half-nothing/simple-fsd/src/interfaces/global"
)

type LoggerDecorator struct {
	logger           Interface
	loggerPrefixName string
}

func NewLoggerAdapter(
	logger Interface,
	loggerPrefixName string,
) Interface {
	return &LoggerDecorator{
		logger:           logger,
		loggerPrefixName: loggerPrefixName,
	}
}

func (loggerDecorator *LoggerDecorator) Init(logPath, logName string, debug, noLogs bool) {
	loggerDecorator.logger.Init(logPath, logName, debug, noLogs)
}

func (loggerDecorator *LoggerDecorator) ShutdownCallback() global.Callable {
	return loggerDecorator.logger.ShutdownCallback()
}

func (loggerDecorator *LoggerDecorator) LogHandler() *slog.Logger {
	return loggerDecorator.logger.LogHandler()
}

func (loggerDecorator *LoggerDecorator) Debug(msg string) {
	loggerDecorator.logger.Debug(fmt.Sprintf("%s | %s", loggerDecorator.loggerPrefixName, msg))
}

func (loggerDecorator *LoggerDecorator) DebugF(msg string, v ...interface{}) {
	loggerDecorator.Debug(fmt.Sprintf(msg, v...))
}

func (loggerDecorator *LoggerDecorator) Info(msg string) {
	loggerDecorator.logger.Info(fmt.Sprintf("%s | %s", loggerDecorator.loggerPrefixName, msg))
}

func (loggerDecorator *LoggerDecorator) InfoF(msg string, v ...interface{}) {
	loggerDecorator.Info(fmt.Sprintf(msg, v...))
}

func (loggerDecorator *LoggerDecorator) Warn(msg string) {
	loggerDecorator.logger.Warn(fmt.Sprintf("%s | %s", loggerDecorator.loggerPrefixName, msg))
}

func (loggerDecorator *LoggerDecorator) WarnF(msg string, v ...interface{}) {
	loggerDecorator.Warn(fmt.Sprintf(msg, v...))
}

func (loggerDecorator *LoggerDecorator) Error(msg string) {
	loggerDecorator.logger.Error(fmt.Sprintf("%s | %s", loggerDecorator.loggerPrefixName, msg))
}

func (loggerDecorator *LoggerDecorator) ErrorF(msg string, v ...interface{}) {
	loggerDecorator.Error(fmt.Sprintf(msg, v...))
}

func (loggerDecorator *LoggerDecorator) Fatal(msg string) {
	loggerDecorator.logger.Fatal(fmt.Sprintf("%s | %s", loggerDecorator.loggerPrefixName, msg))
}

func (loggerDecorator *LoggerDecorator) FatalF(msg string, v ...interface{}) {
	loggerDecorator.Debug(fmt.Sprintf(msg, v...))
}
