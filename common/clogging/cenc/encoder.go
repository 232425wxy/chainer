package cenc

import (
	"io"
	"time"

	zaplogfmt "github.com/sykesm/zap-logfmt"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type Formatter interface {
	Format(w io.Writer, entry zapcore.Entry, fields []zapcore.Field)
}

type FormatEncoder struct {
	zapcore.Encoder
	formatters []Formatter
	pool       buffer.Pool
}

func NewFormatEncoder(formatters ...Formatter) *FormatEncoder {
	return &FormatEncoder{
		Encoder: zaplogfmt.NewEncoder(zapcore.EncoderConfig{
			MessageKey:     "",
			LevelKey:       "",
			TimeKey:        "",
			NameKey:        "",
			CallerKey:      "",
			StacktraceKey:  "",
			LineEnding:     "\n",
			EncodeDuration: zapcore.StringDurationEncoder, // 使用其内置的String方法序列化了一个time.Duration。
			EncodeTime: func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
				pae.AppendString(t.Format("2006-01-02T15:04:05.999Z07:00"))
			},
		}),
		formatters: formatters,
		pool: buffer.NewPool(),
	}
}

func (f *FormatEncoder) Clone() zapcore.Encoder {
	return &FormatEncoder{
		Encoder: f.Encoder.Clone(),
		formatters: f.formatters,
		pool: f.pool,
	}
}

func (f *FormatEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	line := f.pool.Get()
	for _, formatter := range f.formatters {
		formatter.Format(line, entry, fields)
	}

	encodedFields, err := f.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return nil, err
	}
	if line.Len() > 0 && encodedFields.Len() != 1 {
		line.AppendString(" ")
	}
	line.AppendString(encodedFields.String())
	encodedFields.Free()

	return line, nil
}