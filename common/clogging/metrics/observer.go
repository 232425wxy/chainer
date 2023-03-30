package metrics

import (
	"github.com/232425wxy/chainer/common/metrics"
	"go.uber.org/zap/zapcore"
)

var (
	CheckedCountOpts = metrics.CounterOpts{
		Namespace:    "logging",
		Name:         "entries_checked",
		Help:         "Number of log entries checked against the active logging level",
		LabelNames:   []string{"level"},
		LabelHelp:    map[string]string{},
		StatsdFormat: "%{#fqname}.%{level}",
	}

	WriteCountOpts = metrics.CounterOpts{
		Namespace:    "logging",
		Name:         "entries_written",
		Help:         "Number of log entries that are written",
		LabelNames:   []string{"level"},
		LabelHelp:    map[string]string{},
		StatsdFormat: "%{#fqname}.%{level}",
	}
)

type Observer struct {
	CheckedCounter metrics.Counter
	WrittenCounter metrics.Counter
}

func NewObserver(provider metrics.Provider) *Observer {
	return &Observer{
		CheckedCounter: provider.NewCounter(CheckedCountOpts),
		WrittenCounter: provider.NewCounter(WriteCountOpts),
	}
}

// Check 传入的两个参数，只用到了第一个参数。
func (o *Observer) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) {
	o.CheckedCounter.With("level", entry.Level.String()).Add(1)
}

func (o *Observer) WriteEntry(entry zapcore.Entry, fields []zapcore.Field) {
	o.WrittenCounter.With("level", entry.Level.String()).Add(1)
}