package clogging

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zapgrpc"
)

type ChainerLogger struct {
	sl *zap.SugaredLogger
}

func NewChainerLogger(l *zap.Logger, options ...zap.Option) *ChainerLogger {
	return &ChainerLogger{
		sl: l.WithOptions(append(options, zap.AddCallerSkip(1))...).Sugar(),
	}
}

func NewGRPCLogger(l *zap.Logger) *zapgrpc.Logger {
	l = l.WithOptions(zap.AddCaller(), zap.AddCallerSkip(3))
	return zapgrpc.NewLogger(l, zapgrpc.WithDebug())
}

func NewZapLogger(core zapcore.Core, options ...zap.Option) *zap.Logger {
	return zap.New(core, append([]zap.Option{zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)}, options...)...)
}

func (cl *ChainerLogger) With(args ...interface{}) *ChainerLogger {
	return &ChainerLogger{sl: cl.sl.With(args...)}
}

func (cl *ChainerLogger) WithOptions(opts ...zap.Option) *ChainerLogger {
	return &ChainerLogger{sl: cl.sl.Desugar().WithOptions(opts...).Sugar()}
}

func (cl *ChainerLogger) IsEnabledFor(level zapcore.Level) bool {
	// SugaredLogger 的 core 是 zap.Logger。
	return cl.sl.Desugar().Core().Enabled(level)
}

func (cl *ChainerLogger) Zap() *zap.Logger {
	return cl.sl.Desugar()
}

func (cl *ChainerLogger) Sync() error {
	return cl.sl.Sync()
}

// Named 增加一个日志记录器的名字。
func (cl *ChainerLogger) Named(name string) *ChainerLogger {
	return &ChainerLogger{sl: cl.sl.Named(name)}
}

func (cl *ChainerLogger) Debug(args ...interface{}) {
	cl.sl.Debugf(formatArgs(args...))
}

func (cl *ChainerLogger) Debugf(template string, args ...interface{}) {
	cl.sl.Debugf(template, args...)
}

func (cl *ChainerLogger) Debugw(msg string, kvs ...interface{}) {
	cl.sl.Debugw(msg, kvs...)
}

func (cl *ChainerLogger) Info(args ...interface{}) {
	cl.sl.Infof(formatArgs(args...))
}

func (cl *ChainerLogger) Infof(template string, args ...interface{}) {
	cl.sl.Infof(template, args...)
}

func (cl *ChainerLogger) Infow(msg string, kvs ...interface{}) {
	cl.sl.Infow(msg, kvs...)
}

func (cl *ChainerLogger) Warn(args ...interface{}) {
	cl.sl.Warnf(formatArgs(args...))
}

func (cl *ChainerLogger) Warnf(template string, args ...interface{}) {
	cl.sl.Warnf(template, args...)
}

func (cl *ChainerLogger) Warnw(msg string, kvs ...interface{}) {
	cl.sl.Warnw(msg, kvs...)
}

func (cl *ChainerLogger) Error(args ...interface{}) {
	cl.sl.Errorf(formatArgs(args...))
}

func (cl *ChainerLogger) Errorf(template string, args ...interface{}) {
	cl.sl.Errorf(template, args...)
}

func (cl *ChainerLogger) Errorw(msg string, kvs ...interface{}) {
	cl.sl.Errorw(msg, kvs...)
}

func (cl *ChainerLogger) DPanic(args ...interface{}) {
	cl.sl.DPanicf(formatArgs(args...))
}

func (cl *ChainerLogger) DPanicf(template string, args ...interface{}) {
	cl.sl.DPanicf(template, args...)
}

func (cl *ChainerLogger) DPanicw(msg string, kvs ...interface{}) {
	cl.sl.DPanicw(msg, kvs...)
}

func (cl *ChainerLogger) Panic(args ...interface{}) {
	cl.sl.Panicf(formatArgs(args...))
}

func (cl *ChainerLogger) Panicf(template string, args ...interface{}) {
	cl.sl.Panicf(template, args...)
}

func (cl *ChainerLogger) Panicw(msg string, kvs ...interface{}) {
	cl.sl.Panicw(msg, kvs...)
}

func (cl *ChainerLogger) Fatal(args ...interface{}) {
	cl.sl.Fatalf(formatArgs(args...))
}

func (cl *ChainerLogger) Fatalf(template string, args ...interface{}) {
	cl.sl.Fatalf(template, args...)
}

func (cl *ChainerLogger) Fatalw(msg string, kvs ...interface{}) {
	cl.sl.Fatalw(msg, kvs...)
}

func formatArgs(args ...interface{}) string {
	return strings.TrimSuffix(fmt.Sprintln(args...), "\n")
}
