package clogging_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/232425wxy/chainer/common/clogging"
	"github.com/232425wxy/chainer/common/clogging/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestNew(t *testing.T) {
	logging, err := clogging.New(clogging.Config{})
	require.NoError(t, err)
	require.Equal(t, zapcore.InfoLevel, logging.DefaultLevel())

	_, err = clogging.New(clogging.Config{LogSpec: "::=broken=::"})
	require.EqualError(t, err, "invalid logging specification '::=broken=::': bad segment '=broken='")
}

func TestNewWithEnvironment(t *testing.T) {
	oldSpec, set := os.LookupEnv("CHAINER_LOGGING_SPEC")
	if set {
		defer os.Setenv("CHAINER_LOGGING_SPEC", oldSpec)
	}

	os.Setenv("CHAINER_LOGGING_SPEC", "fatal")
	logging, err := clogging.New(clogging.Config{})
	require.NoError(t, err)
	require.Equal(t, zapcore.FatalLevel, logging.DefaultLevel())

	os.Unsetenv("CHAINER_LOGGING_SPEC")
	logging, err = clogging.New(clogging.Config{})
	require.NoError(t, err)
	require.Equal(t, zapcore.InfoLevel, logging.DefaultLevel())
}

func TestLoggingSetWriter(t *testing.T) {
	ws := &mock.WriteSyncer{}
	buf := &bytes.Buffer{}
	logging, err := clogging.New(clogging.Config{Writer: buf})
	require.NoError(t, err)
	old := logging.SetWriter(ws)
	logging.SetWriter(buf)
	original := logging.SetWriter(ws)
	require.Exactly(t, old, original) // Exactly asserts that two objects are equal in value and type.

	_, err = logging.Write([]byte("hello"))
	require.NoError(t, err)
	require.Equal(t, 1, ws.WriteCallCount())
	require.Equal(t, []byte("hello"), ws.WriteArgsForCall(0))

	err = logging.Sync()
	require.NoError(t, err)

	ws.SetSyncReturns(errors.New("phone"))
	err = logging.Sync()
	require.EqualError(t, err, "phone")
}

func TestNamedLogger(t *testing.T) {
	defer clogging.Reset()
	buf := &bytes.Buffer{}
	clogging.Global.SetWriter(buf)

	t.Run("logger and named (child) logger with different levels", func(t *testing.T) {
		defer buf.Reset()
		logger := clogging.MustGetLogger("chameleon")
		logger2 := logger.Named("hash")
		clogging.ActivateSpec("chameleon=debug:chameleon.hash=error")

		logger.Info("from chameleon: info")
		logger.Debug("from chameleon: debug")
		logger2.Info("from hash: info")
		require.Contains(t, buf.String(), "from chameleon: info")
		require.Contains(t, buf.String(), "from chameleon: debug")
		require.NotContains(t, buf.String(), "from hash: info")
	})

	t.Run("named logger where parent logger isn't enabled", func(t *testing.T) {
		logger := clogging.MustGetLogger("chameleon")
		logger2 := logger.Named("dragon")
		clogging.ActivateSpec("chameleon=fatal:chameleon.dragon=error")
		logger.Error("chameleon: error")
		logger2.Error("dragon: error")
		require.NotContains(t, buf.String(), "chameleon: error")
		require.Contains(t, buf.String(), "dragon: error")
	})
}

func TestInvalidLoggerName(t *testing.T) {
	names := []string{"test*", ".test", "test.", ".", ""}
	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			msg := fmt.Sprintf("invalid logger name: %s", name)
			require.PanicsWithValue(t, msg, func ()  {
				clogging.MustGetLogger(name)
			})
		})
	}
}

func TestCheck(t *testing.T) {
	l := &clogging.Logging{}
	observer := &mock.Observer{}
	e := zapcore.Entry{}
	l.SetObserver(observer)
	l.Check(e, nil)
	require.Equal(t, 1, observer.CheckCallCount())
	e, ce := observer.CheckArgsForCall(0)
	require.Equal(t, e, zapcore.Entry{})
	require.Nil(t, ce)

	l.WriteEntry(e, nil)
	require.Equal(t, 1, observer.WriteEntryCallCount())
	e, f := observer.WriteEntryArgsForCall(0)
	require.Equal(t, e, zapcore.Entry{})
	require.Nil(t, f)

	//	remove observer
	l.SetObserver(nil)
	l.Check(zapcore.Entry{}, nil)
	require.Equal(t, 1, observer.CheckCallCount())
}

func TestLoggerCoreCheck(t *testing.T) {
	logging, err := clogging.New(clogging.Config{})
	require.NoError(t, err)

	logger := logging.ZapLogger("foo")

	err = logging.ActivateSpec("info")
	require.NoError(t, err)
	require.False(t, logger.Core().Enabled(zapcore.DebugLevel), "debug should not be enabled at info level")

	err = logging.ActivateSpec("debug")
	require.NoError(t, err)
	require.True(t, logger.Core().Enabled(zapcore.DebugLevel), "debug should now be enabled at debug level")
}