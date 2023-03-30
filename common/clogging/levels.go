package clogging

import (
	"fmt"
	"math"

	"go.uber.org/zap/zapcore"
)

const (
	// DisabledLevel 代表一个禁用的日志级别，日志级别被设置成 DisabledLevel，则不会打印日志。
	DisabledLevel = zapcore.Level(math.MinInt8)

	// PayloadLevel 表示会输出非常详细日志信息的日志级别，最低的日志级别。
	PayloadLevel = zapcore.Level(zapcore.DebugLevel - 1) 
)

func IsValidLevel(level string) bool {
	_, err := nameToLevel(level)
	return err == nil
}

func NameToLevel(level string) zapcore.Level {
	l, err := nameToLevel(level)
	if err != nil {
		return zapcore.InfoLevel
	}
	return l
}

func nameToLevel(level string) (zapcore.Level, error) {
	switch level {
	case "PAYLOAD", "payload":
		return PayloadLevel, nil
	case "DEBUG", "debug":
		return zapcore.DebugLevel, nil
	case "INFO", "info":
		return zapcore.InfoLevel, nil
	case "WARNING", "WARN", "warning", "warn":
		return zapcore.WarnLevel, nil
	case "ERROR", "error":
		return zapcore.ErrorLevel, nil
	case "DPANIC", "dpanic":
		return zapcore.DPanicLevel, nil
	case "PANIC", "panic":
		return zapcore.PanicLevel, nil
	case "FATAL", "fatal":
		return zapcore.FatalLevel, nil

	default:
		return DisabledLevel, fmt.Errorf("invalid log level: %s", level)
	}
}
