package clogging

import "go.uber.org/zap/zapcore"

type Encoding int8

const (
	CONSOLE = iota
	JSON
	LOGFMT
)

// EncodingSelector 用于决定日志记录被编码成何种格式。
type EncodingSelector interface {
	Encoding() Encoding
}

type Observer interface {
	Check(entry zapcore.Entry, ce *zapcore.CheckedEntry)
	WriteEntry(entry zapcore.Entry, fields []zapcore.Field)
}

type Core struct {
	zapcore.LevelEnabler // LevelEnabler 决定在记录消息时是否启用一个给定的日志级别。
	Levels *LoggerLevels
	Encoders map[Encoding]zapcore.Encoder
	Selector EncodingSelector
	Output zapcore.WriteSyncer
	Observer Observer
}

func (c *Core) With(fields []zapcore.Field) zapcore.Core {
	clones := map[Encoding]zapcore.Encoder{}
	for name, enc := range c.Encoders {
		clone := enc.Clone()
		addFields(clone, fields)
		clones[name] = clone
	}
	return &Core{
		LevelEnabler: c.LevelEnabler,
		Levels:       c.Levels,
		Encoders:     clones,
		Selector:     c.Selector,
		Output:       c.Output,
		Observer:     c.Observer,
	}
}

func (c *Core) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Observer != nil {
		c.Observer.Check(entry, ce)
	}

	if c.Enabled(entry.Level) && c.Levels.Level(entry.LoggerName).Enabled(entry.Level) {
		return ce.AddCore(entry, c)
	}
	return ce
}

func (c *Core) Write(e zapcore.Entry, fields []zapcore.Field) error {
	encoding := c.Selector.Encoding()
	enc := c.Encoders[encoding]

	buf, err := enc.EncodeEntry(e, fields)
	if err != nil {
		return err
	}
	_, err = c.Output.Write(buf.Bytes())
	buf.Free()
	if err != nil {
		return err
	}

	if e.Level >= zapcore.PanicLevel {
		c.Sync()
	}

	if c.Observer != nil {
		c.Observer.WriteEntry(e, fields)
	}

	return nil
}

func (c *Core) Sync() error {
	return c.Output.Sync()
}

func addFields(enc zapcore.ObjectEncoder, fields []zapcore.Field) {
	for i := range fields {
		fields[i].AddTo(enc)
	}
}

