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
}