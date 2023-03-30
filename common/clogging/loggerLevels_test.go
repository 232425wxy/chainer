package clogging

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestLoggerLevelsActivateSpec(t *testing.T) {
	var tests = []struct {
		spec                 string
		expectedLevels       map[string]zapcore.Level
		expectedDefaultLevel zapcore.Level
	}{
		{
			spec:                 "DEBUG",
			expectedLevels:       map[string]zapcore.Level{},
			expectedDefaultLevel: zapcore.DebugLevel,
		},
		{
			spec:                 "DEBUG:",
			expectedLevels:       map[string]zapcore.Level{},
			expectedDefaultLevel: zapcore.InfoLevel,
		},
		{
			spec: "logger=info:DEBUG",
			expectedLevels: map[string]zapcore.Level{
				"logger":     zapcore.InfoLevel,
				"logger.a":   zapcore.InfoLevel,
				"logger.a.b": zapcore.InfoLevel,
			},
			expectedDefaultLevel: zapcore.DebugLevel,
		},
		{
			spec: "xx.yy.zz=panic:debug:aa.bb=error",
			expectedLevels: map[string]zapcore.Level{
				"xx":             zapcore.DebugLevel, // 没有为 xx 日志记录器设置过日志级别，那么就将默认的日志级别作为它的日志级别。
				"xx.yy.zz.__":    zapcore.PanicLevel,
				"aaa.bb":         zapcore.DebugLevel,
				"aa.bb.cc.dd.ee": zapcore.ErrorLevel,
			},
			expectedDefaultLevel: zapcore.DebugLevel,
		},
		{
			spec: "xx.yy.zz,XX.YY=panic:debug:aa.bb=error",
			expectedLevels: map[string]zapcore.Level{
				"xx":             zapcore.DebugLevel, // 没有为 xx 日志记录器设置过日志级别，那么就将默认的日志级别作为它的日志级别。
				"xx.yy.zz.__":    zapcore.PanicLevel,
				"aaa.bb":         zapcore.DebugLevel,
				"aa.bb.cc.dd.ee": zapcore.ErrorLevel,
				"XX.YY.ZZ":       zapcore.PanicLevel,
				"XX.YY.":         zapcore.PanicLevel,
			},
			expectedDefaultLevel: zapcore.DebugLevel,
		},
		{
			spec: "info:warn",
			expectedLevels: map[string]zapcore.Level{
				"xx":             zapcore.WarnLevel, // 没有为 xx 日志记录器设置过日志级别，那么就将默认的日志级别作为它的日志级别。
				"xx.yy.zz.__":    zapcore.WarnLevel,
				"aaa.bb":         zapcore.WarnLevel,
				"aa.bb.cc.dd.ee": zapcore.WarnLevel,
				"XX.YY.ZZ":       zapcore.WarnLevel,
				"XX.YY.":         zapcore.WarnLevel,
			},
			expectedDefaultLevel: zapcore.WarnLevel,
		},
	}

	for _, test := range tests {
		t.Run(test.spec, func(t *testing.T) {
			ll := &LoggerLevels{}

			err := ll.ActivateSpec(test.spec)
			require.NoError(t, err)
			require.Equal(t, test.expectedDefaultLevel, ll.DefaultLevel())
			for name, lvl := range test.expectedLevels {
				require.Equal(t, lvl, ll.Level(name))
			}
		})
	}
}

func TestLoggerLevelsActivateSpecErrors(t *testing.T) {
	var tests = []struct {
		spec string
		err  error
	}{
		{spec: "=INFO:DEBUG", err: errors.New("invalid logging specification '=INFO:DEBUG': no logger specified in segment '=INFO'")},
		{spec: "=INFO=:DEBUG", err: errors.New("invalid logging specification '=INFO=:DEBUG': bad segment '=INFO='")},
		{spec: "cat", err: errors.New("invalid logging specification 'cat': bad segment 'cat'")},
		{spec: "a$=info", err: errors.New("invalid logging specification 'a$=info': bad logger name 'a$'")},
	}

	for _, test := range tests {
		t.Run(test.spec, func(t *testing.T) {
			ll := &LoggerLevels{}
			err := ll.ActivateSpec("fatal:a=warn")
			require.NoError(t, err)

			err = ll.ActivateSpec(test.spec)
			require.EqualError(t, err, test.err.Error())

			require.Equal(t, zapcore.FatalLevel, ll.DefaultLevel(), "default should not change")
			require.Equal(t, zapcore.WarnLevel, ll.Level("a.b"))
		})
	}
}

func TestSpec(t *testing.T) {
	var tests = []struct{
		input string
		output string
	}{
		{input: "", output: "info"},
		{input: "a.b.c=debug:warn:x.y=panic", output: "a.b.c=debug:x.y=panic:warn"},
		{input: "a.b.c=debug:x.y=panic", output: "a.b.c=debug:x.y=panic:info"},
	}

	for _, test := range tests {
		ll := &LoggerLevels{}
		err := ll.ActivateSpec(test.input)
		require.NoError(t, err)
		require.Equal(t, test.output, ll.Spec())
	}
}

func TestEnables(t *testing.T) {
	var tests = []struct{
		spec string
		enabledAt zapcore.Level
	}{
		{spec: "payload", enabledAt: PayloadLevel},
		{spec: "debug", enabledAt: zapcore.DebugLevel},
		{spec: "info", enabledAt: zapcore.InfoLevel},
		{spec: "warn", enabledAt: zapcore.WarnLevel},
		{spec: "error", enabledAt: zapcore.ErrorLevel},
		{spec: "dpanic", enabledAt: zapcore.DPanicLevel},
		{spec: "panic", enabledAt: zapcore.PanicLevel},
		{spec: "fatal", enabledAt: zapcore.FatalLevel},
		{spec: "fatal", enabledAt: zapcore.FatalLevel},
		{spec: "a=debug:b=error", enabledAt: zapcore.DebugLevel},
	}

	for _, test := range tests {
		t.Run(test.spec, func(t *testing.T) {
			ll := &LoggerLevels{}
			err := ll.ActivateSpec(test.spec)
			require.NoError(t, err)

			for i := PayloadLevel; i <= zapcore.FatalLevel; i++ {
				if test.enabledAt <= i {
					require.Truef(t, ll.Enabled(i), "expected level %s and spec %s to be enabled", zapcore.Level(i), test.spec)
				} else {
					require.False(t, ll.Enabled(i), "expected level %s and spec %s to be disabled", zapcore.Level(i), test.spec)
				}
			}
		})
	}
}