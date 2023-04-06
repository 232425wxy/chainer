package clogging

import (
	"io"

	"google.golang.org/grpc/grpclog"
)

var Global *Logging

func init() {
	logging, err := New(Config{})
	if err != nil {
		panic(err)
	}
	Global = logging
	grpcLogger := Global.ZapLogger("grpc")
	grpclog.SetLoggerV2(NewGRPCLogger(grpcLogger))
}

func Init(config Config) {
	if err := Global.Apply(config); err != nil {
		panic(err)
	}
}

func Reset() {
	Global.Apply(Config{})
}

// LoggerName 返回给定日志记录器名对应的日志记录级别。
func LoggerLevel(loggerName string) string {
	return Global.Level(loggerName).String()
}

func MustGetLogger(loggerName string) *ChainerLogger {
	return Global.Logger(loggerName)
}

func ActivateSpec(spec string) {
	if err := Global.ActivateSpec(spec); err != nil {
		panic(err)
	}
}

func DefaultLevel() string {
	return defaultLevel.String()
}

func SetWriter(w io.Writer) io.Writer {
	return Global.SetWriter(w)
}

func SetObserver(observer Observer) Observer {
	return Global.SetObserver(observer)
}