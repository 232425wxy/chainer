package clogging_test

import (
	"os"
	"testing"

	"github.com/232425wxy/chainer/common/clogging"
	"github.com/232425wxy/chainer/common/clogging/cenc"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestFormatEncoding(t *testing.T) {
	format := "%{color}[%{module}] %{shortfunc} => %{level}: %{message}%{color:reset}"
	formatters, err := cenc.ParseFormat(format)
	require.NoError(t, err)
	encoder := cenc.NewFormatEncoder(formatters...)
	buf := os.Stdout // 输出到显示屏。
	// buf := &bytes.Buffer{}
	core := zapcore.NewCore(encoder, zapcore.AddSync(buf), zap.NewAtomicLevel())
	zapLogger := clogging.NewZapLogger(core).Named("module_name").With(zap.String("key1", "value1"))
	logger := clogging.NewChainerLogger(zapLogger)

	logger.Info("hello, everyone!")
	logger.Debug("debug")
}

func TestLevel(t *testing.T) {
	format := "%{color}[%{module}] %{shortfunc} => %{level}: %{message}%{color:reset}"
	spec := "debug"

	logging, _ := clogging.New(clogging.Config{
		Format: format,
		LogSpec: spec,
		Writer: os.Stdout,
	})
	
	zl := logging.ZapLogger("module_name")
	cl := clogging.NewChainerLogger(zl)

	cl.Debug("你好呀！")
}