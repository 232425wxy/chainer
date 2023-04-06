package clogging_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/232425wxy/chainer/common/clogging"
	"github.com/232425wxy/chainer/common/clogging/cenc"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestChainerLoggerEncoding(t *testing.T) {
	formatters, err := cenc.ParseFormat("%{color}[%{module}] %{shortfunc} -> %{level}%{color:reset} %{message}")
	require.NoError(t, err)
	enc := cenc.NewFormatEncoder(formatters...)
	buf := &bytes.Buffer{}
	core := zapcore.NewCore(enc, zapcore.AddSync(buf), zap.NewAtomicLevel())
	zl := clogging.NewZapLogger(core).Named("test").With(zap.String("age", "18"))
	cl := clogging.NewChainerLogger(zl)

	buf.Reset()
	cl.Info("string value", 0, 1, 1.23, struct{}{})
	require.Equal(t, "\x1b[34m[test] TestChainerLoggerEncoding -> INFO\x1b[0m string value 0 1 1.23 {} age=18\n", buf.String())

	buf.Reset()
	cl.Infof("string %s, %d, %.3f, %v", "strval", 0, 1.23, struct{}{})
	require.Equal(t, "\x1b[34m[test] TestChainerLoggerEncoding -> INFO\x1b[0m string strval, 0, 1.230, {} age=18\n", buf.String())

	buf.Reset()
	cl.Infow("this is a message", "int", 0, "float", 1.23, "struct", struct{}{})
	require.Equal(t, "\x1b[34m[test] TestChainerLoggerEncoding -> INFO\x1b[0m this is a message age=18 int=0 float=1.23 struct={}\n", buf.String())

}

func TestChainerLogger(t *testing.T) {
	var enabler zap.LevelEnablerFunc = func(l zapcore.Level) bool { return true }

	var tests = []struct {
		desc    string
		f       func(fl *clogging.ChainerLogger)
		level   zapcore.Level
		message string
		fields  []zapcore.Field
		panics  bool
	}{
		{
			desc:    "DPanic",
			f:       func(fl *clogging.ChainerLogger) { fl.DPanic("arg1", "arg2") },
			level:   zapcore.DPanicLevel,
			message: "arg1 arg2",
			fields:  []zapcore.Field{},
		},
		{
			desc:    "DPanicf",
			f:       func(fl *clogging.ChainerLogger) { fl.DPanicf("panic: %s, %d", "reason", 99) },
			level:   zapcore.DPanicLevel,
			message: "panic: reason, 99",
			fields:  []zapcore.Field{},
		},
		{
			desc:    "DPanicw",
			f:       func(fl *clogging.ChainerLogger) { fl.DPanicw("I'm in a panic", "reason", "something", "code", 99) },
			level:   zapcore.DPanicLevel,
			message: "I'm in a panic",
			fields:  []zapcore.Field{zap.String("reason", "something"), zap.Int("code", 99)},
		},
		{
			desc:    "Debug",
			f:       func(fl *clogging.ChainerLogger) { fl.Debug("arg1", "arg2") },
			level:   zapcore.DebugLevel,
			message: "arg1 arg2",
			fields:  []zapcore.Field{},
		},
		{
			desc:    "Debugf",
			f:       func(fl *clogging.ChainerLogger) { fl.Debugf("debug: %s, %d", "goo", 99) },
			level:   zapcore.DebugLevel,
			message: "debug: goo, 99",
			fields:  []zapcore.Field{},
		},
		{
			desc:    "Debugw",
			f:       func(fl *clogging.ChainerLogger) { fl.Debugw("debug data", "key", "value") },
			level:   zapcore.DebugLevel,
			message: "debug data",
			fields:  []zapcore.Field{zap.String("key", "value")},
		},
		{
			desc:    "Error",
			f:       func(fl *clogging.ChainerLogger) { fl.Error("oh noes", errors.New("bananas")) },
			level:   zapcore.ErrorLevel,
			message: "oh noes bananas",
			fields:  []zapcore.Field{},
		},
		{
			desc:    "Errorf",
			f:       func(fl *clogging.ChainerLogger) { fl.Errorf("error: %s", errors.New("bananas")) },
			level:   zapcore.ErrorLevel,
			message: "error: bananas",
			fields:  []zapcore.Field{},
		},
		{
			desc:    "Errorw",
			f:       func(fl *clogging.ChainerLogger) { fl.Errorw("something failed", "err", errors.New("bananas")) },
			level:   zapcore.ErrorLevel,
			message: "something failed",
			fields:  []zapcore.Field{zap.NamedError("err", errors.New("bananas"))},
		},
		{
			desc:    "Info",
			f:       func(fl *clogging.ChainerLogger) { fl.Info("fyi", "things are great") },
			level:   zapcore.InfoLevel,
			message: "fyi things are great",
			fields:  []zapcore.Field{},
		},
		{
			desc:    "Infof",
			f:       func(fl *clogging.ChainerLogger) { fl.Infof("fyi: %s", "things are great") },
			level:   zapcore.InfoLevel,
			message: "fyi: things are great",
			fields:  []zapcore.Field{},
		},
		{
			desc:    "Infow",
			f:       func(fl *clogging.ChainerLogger) { fl.Infow("fyi", "fish", "are smelly", "fruit", "is sweet") },
			level:   zapcore.InfoLevel,
			message: "fyi",
			fields:  []zapcore.Field{zap.String("fish", "are smelly"), zap.String("fruit", "is sweet")},
		},
		{
			desc:    "Panic",
			f:       func(fl *clogging.ChainerLogger) { fl.Panic("oh noes", errors.New("platypus")) },
			level:   zapcore.PanicLevel,
			message: "oh noes platypus",
			fields:  []zapcore.Field{},
			panics:  true,
		},
		{
			desc:    "Panicf",
			f:       func(fl *clogging.ChainerLogger) { fl.Panicf("error: %s", errors.New("platypus")) },
			level:   zapcore.PanicLevel,
			message: "error: platypus",
			fields:  []zapcore.Field{},
			panics:  true,
		},
		{
			desc:    "Panicw",
			f:       func(fl *clogging.ChainerLogger) { fl.Panicw("something failed", "err", errors.New("platypus")) },
			level:   zapcore.PanicLevel,
			message: "something failed",
			fields:  []zapcore.Field{zap.NamedError("err", errors.New("platypus"))},
			panics:  true,
		},
		{
			desc:    "Warn",
			f:       func(fl *clogging.ChainerLogger) { fl.Warn("oh noes", errors.New("monkeys")) },
			level:   zapcore.WarnLevel,
			message: "oh noes monkeys",
			fields:  []zapcore.Field{},
		},
		{
			desc:    "Warnf",
			f:       func(fl *clogging.ChainerLogger) { fl.Warnf("error: %s", errors.New("monkeys")) },
			level:   zapcore.WarnLevel,
			message: "error: monkeys",
			fields:  []zapcore.Field{},
		},
		{
			desc:    "Warnw",
			f:       func(fl *clogging.ChainerLogger) { fl.Warnw("something is weird", "err", errors.New("monkeys")) },
			level:   zapcore.WarnLevel,
			message: "something is weird",
			fields:  []zapcore.Field{zap.NamedError("err", errors.New("monkeys"))},
		},
		{
			desc:    "With",
			f:       func(fl *clogging.ChainerLogger) { fl.With("key", "value").Debug("cool messages", "and stuff") },
			level:   zapcore.DebugLevel,
			message: "cool messages and stuff",
			fields:  []zapcore.Field{zap.String("key", "value")},
		},
		{
			desc: "WithOptions",
			f: func(fl *clogging.ChainerLogger) {
				fl.WithOptions(zap.Fields(zap.String("optionkey", "optionvalue"))).Debug("cool messages", "and stuff")
			},
			level:   zapcore.DebugLevel,
			message: "cool messages and stuff",
			fields:  []zapcore.Field{zap.String("optionkey", "optionvalue")},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			core, logs := observer.New(enabler)
			fl := clogging.NewChainerLogger(zap.New(core)).Named("lname")

			if tc.panics {
				require.Panics(t, func() { tc.f(fl) })
			} else {
				tc.f(fl)
			}

			err := fl.Sync()
			require.NoError(t, err)

			entries := logs.All()
			require.Len(t, entries, 1)
			entry := entries[0]

			require.Equal(t, tc.level, entry.Level)
			require.Equal(t, tc.message, entry.Message)
			require.Equal(t, tc.fields, entry.Context)
			require.Equal(t, "lname", entry.LoggerName)
		})
	}
}
