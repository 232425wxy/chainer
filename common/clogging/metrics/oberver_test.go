package metrics

import (
	"testing"

	commonmetrics "github.com/232425wxy/chainer/common/metrics"
	"github.com/232425wxy/chainer/common/metrics/metricsfakes"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestNewObserver(t *testing.T) {
	provider := &metricsfakes.Provider{}
	checkedCounter := &metricsfakes.Counter{}
	writtenCounter := &metricsfakes.Counter{}

	provider.NewCounterStub = func(opts commonmetrics.CounterOpts) commonmetrics.Counter {
		switch opts.Name {
		case "entries_checked":
			require.Equal(t, CheckedCountOpts, opts)
			return checkedCounter
		case "entries_written":
			require.Equal(t, WriteCountOpts, opts)
			return writtenCounter
		default:
			return nil
		}
	}

	expectedObserver := &Observer{
		CheckedCounter: checkedCounter,
		WrittenCounter: writtenCounter,
	}

	m := NewObserver(provider)
	require.Equal(t, expectedObserver, m)
	require.Equal(t, 2, provider.NewCounterCallCount())
}

func TestCheck(t *testing.T) {
	counter := &metricsfakes.Counter{}
	counter.SetWithReturns(counter)

	m := Observer{CheckedCounter: counter}
	entry := zapcore.Entry{Level: zapcore.DebugLevel}
	checkedEntry := &zapcore.CheckedEntry{}
	m.Check(entry, checkedEntry)

	require.Equal(t, 1, counter.WithCallCount())
	require.Equal(t, []string{"level", "debug"}, counter.WithArgsForCall(0))

	require.Equal(t, 1, counter.AddCallCount())
	require.Equal(t, 1.0, counter.AddArgsForCall(0))
}

func TestWrite(t *testing.T) {
	counter := &metricsfakes.Counter{}
	counter.SetWithReturns(counter)

	m := Observer{WrittenCounter: counter}
	entry := zapcore.Entry{Level: zapcore.DebugLevel}
	m.WriteEntry(entry, nil)

	require.Equal(t, 1, counter.AddCallCount())
	require.Equal(t, 1.0, counter.AddArgsForCall(0))

	require.Equal(t, 1, counter.WithCallCount())
	require.Equal(t, []string{"level", "debug"}, counter.WithArgsForCall(0))
}