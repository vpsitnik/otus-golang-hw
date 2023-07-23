package logger

import (
	"strings"

	"go.uber.org/zap"
)

type Logger struct {
	instance *zap.Logger
}

const (
	debugLevel   = "debug"
	infoLevel    = "info"
	warningLevel = "warning"
	errorLevel   = "error"
)

func New(level string) *Logger {
	logConfig := zap.NewProductionConfig()

	switch strings.ToLower(level) {
	case debugLevel:
		logConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case infoLevel:
		logConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case warningLevel:
		logConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case errorLevel:
		logConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	}

	logger, err := logConfig.Build()
	if err != nil {
		return nil
	}

	return &Logger{instance: logger}
}

func (l Logger) Debug(msg string) {
	l.instance.Debug(msg)
}

func (l Logger) Info(msg string) {
	l.instance.Info(msg)
}

func (l Logger) Warning(msg string) {
	l.instance.Warn(msg)
}

func (l Logger) Error(msg string) {
	l.instance.Error(msg)
}
