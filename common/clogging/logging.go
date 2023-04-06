package clogging

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/232425wxy/chainer/common/clogging/cenc"
	zaplogfmt "github.com/sykesm/zap-logfmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	defaultLevel  = zapcore.InfoLevel
	defaultFormat = "%{color}%{time:2006-01-02 15:04:05.000 MST} [%{module}] %{shortfunc} -> %{level:.4s} %{id:03x}%{color:reset} %{message}"
)

type Config struct {
	Format  string
	LogSpec string
	Writer  io.Writer
}

type Logging struct {
	*LoggerLevels
	mutex          sync.RWMutex
	encoding       Encoding
	encoderConfig  zapcore.EncoderConfig
	multiFormatter *cenc.MultiFormatter
	writer         zapcore.WriteSyncer
	observer       Observer
}

func New(c Config) (*Logging, error) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.NameKey = "name"

	l := &Logging{
		LoggerLevels:   &LoggerLevels{defaultLevel: defaultLevel},
		encoderConfig:  encoderConfig,
		multiFormatter: cenc.NewMultiFormatter(),
	}

	if err := l.Apply(c); err != nil {
		return nil, err
	}
	return l, nil
}

func (l *Logging) Apply(c Config) error {
	err := l.SetFormat(c.Format)
	if err != nil {
		return err
	}
	if c.LogSpec == "" {
		c.LogSpec = os.Getenv("CHAINER_LOGGING_SPEC")
	}
	// 也有可能系统环境没有设置 "CHAINER_LOGGING_SPEC"。
	if c.LogSpec == "" {
		c.LogSpec = defaultLevel.String()
	}
	if err = l.LoggerLevels.ActivateSpec(c.LogSpec); err != nil {
		return err
	}

	if c.Writer == nil {
		c.Writer = os.Stderr
	}
	l.SetWriter(c.Writer)
	return nil

}

func (l *Logging) SetFormat(format string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if format == "" {
		format = defaultFormat
	}

	if format == "json" {
		l.encoding = JSON
		return nil
	}

	if format == "logfmt" {
		l.encoding = LOGFMT
		return nil
	}

	formatters, err := cenc.ParseFormat(format) // 可能是默认的格式："%{color}%{time:2006-01-02 15:04:05.000 MST} [%{module}] %{shortfunc} -> %{level:.4s} %{id:03x}%{color:reset} %{message}"。
	if err != nil {
		return err
	}
	l.multiFormatter.SetFormatters(formatters)
	l.encoding = CONSOLE

	return nil
}

// SetWriter控制格式化的日志记录被写入哪个写入器。
// 除了*os.File之外，写程序需要安全地被多个go例程同时使用。
func (l *Logging) SetWriter(w io.Writer) io.Writer {
	var ws zapcore.WriteSyncer

	switch t := w.(type) {
	case *os.File:
		ws = zapcore.Lock(t) // 将 os.File 包裹在一个 mutex 里，以使其能支持并发操作。
	case zapcore.WriteSyncer:
		ws = t
	default:
		ws = zapcore.AddSync(w) // 将 io.Writer 转化为 zapcore.WriteSyncer。
	}

	l.mutex.Lock()
	old := l.writer
	l.writer = ws
	l.mutex.Unlock()
	return old
}

// SetObserver 用于提供一个日志观察者，当日志级别被检查或写入时，它将被调用。
func (l *Logging) SetObserver(observer Observer) Observer {
	l.mutex.Lock()
	old := l.observer
	l.observer = observer
	l.mutex.Unlock()
	return old
}

func (l *Logging) Write(bz []byte) (int, error) {
	l.mutex.RLock()
	w := l.writer
	l.mutex.RUnlock()
	return w.Write(bz)
}

func (l *Logging) Sync() error {
	l.mutex.RLock()
	w := l.writer
	l.mutex.RUnlock()
	return w.Sync()
}

func (l *Logging) Encoding() Encoding {
	l.mutex.RLock()
	e := l.encoding
	l.mutex.RUnlock()
	return e
}

func (l *Logging) ZapLogger(name string) *zap.Logger {
	if !isValidLoggerName(name) {
		panic(fmt.Sprintf("invalid logger name: %s", name))
	}

	l.mutex.RLock()
	core := &Core{
		LevelEnabler: l.LoggerLevels,
		Levels:       l.LoggerLevels,
		Encoders:     map[Encoding]zapcore.Encoder{
			JSON: zapcore.NewJSONEncoder(l.encoderConfig),
			CONSOLE: cenc.NewFormatEncoder(l.multiFormatter),
			LOGFMT: zaplogfmt.NewEncoder(l.encoderConfig),
		},
		Selector:     l,
		Output:       l,
		Observer:     l,
	}
	l.mutex.RUnlock()

	return NewZapLogger(core).Named(name)
}

func (l *Logging) Check(e zapcore.Entry, ce *zapcore.CheckedEntry) {
	l.mutex.RLock()
	observer := l.observer
	l.mutex.RUnlock()
	if observer != nil {
		observer.Check(e, ce)
	}
}

func (l *Logging) WriteEntry(e zapcore.Entry, fields []zapcore.Field) {
	l.mutex.RLock()
	observer := l.observer
	l.mutex.RUnlock()
	if observer != nil {
		observer.WriteEntry(e, fields)
	}
}

func (l *Logging) Logger(name string) *ChainerLogger {
	zl := l.ZapLogger(name)
	return NewChainerLogger(zl)
}
