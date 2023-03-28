package cenc

import (
	"fmt"
	"io"
	"regexp"

	"go.uber.org/zap/zapcore"
)

// formatRegexp 匹配案例：
// %{color:red}%{level:debug}  匹配结果：共找到两处匹配：%{color:red}和%{level:debug}；
// %{color}%{messagee}  匹配结果：共找到一处匹配：%{color}；
// %{id:123}xxxx%{module::p2p}  匹配结果：共找到两处匹配：%{id:123}和%{module:p2p}。
var formatRegexp = regexp.MustCompile(`%{(color|id|level|message|module|shortfunc|time)(?::(.*?))?}`)

func ParseFormat(spec string) ([]Formatter, error) {
	cursor := 0
	formatters := []Formatter{}

	matches := formatRegexp.FindAllStringSubmatchIndex(spec, -1)
	for _, m := range matches {
		// 假设 spec 是 `%{level:.4s}`。
		start, end := m[0], m[1]             // start 匹配到的是 `%` 的下标位置，end 匹配到的是 `}` 后一个字符下标的位置。
		verbStart, verbEnd := m[2], m[3]     // verbStart 匹配到的是字母 `l` 的下标位置，verbEnd 匹配到的是冒号 `:` 的下标位置。
		formatStart, formatEnd := m[4], m[5] // formatStart 匹配到的是点号 `.` 的下标位置，formatEnd 匹配到的是 `}` 下标的位置。

		if start > cursor {
			// 代表 `%{` 前面还有内容，前面的内容就都当作普通字符串来处理，处理的方式就是直接将其原封不动的写入到 io.Writer 里。
			formatters = append(formatters, StringFormatter{Value: spec[cursor:start]})
		}

		var format string
		if formatStart >= 0 {
			// 匹配不到的情况下，formatStart 等于 -1
			format = spec[formatStart:formatEnd]
		}

		formatter, err := NewFormatter(spec[verbStart:verbEnd], format)
		if err != nil {
			return nil, err
		}
		formatters = append(formatters, formatter)
		cursor = end
	}

	if cursor != len(spec) {
		// 最后一个 `}` 后面还有内容，后面的内容就都当作普通字符串来处理，处理的方式就是直接将其原封不动的写入到 io.Writer 里。
		formatters = append(formatters, StringFormatter{Value: spec[cursor:]})
	}

	return formatters, nil
}

func NewFormatter(verb, format string) (Formatter, error) {
	switch verb {
	case "color":
		return newColorFormatter(format)
	default:
		return nil, fmt.Errorf("unknown verb: %s", verb)
	}
}

type ColorFormatter struct {
	Bold  bool
	Reset bool
}

func newColorFormatter(format string) (ColorFormatter, error) {
	switch format {
	case "bold":
		return ColorFormatter{Bold: true}, nil
	case "reset":
		return ColorFormatter{Reset: true}, nil
	case "": // 空的情况下，既不加粗，也不重置
		return ColorFormatter{}, nil 
	default:
		return ColorFormatter{}, fmt.Errorf("invalid color option: %s", format)
	}
}

func (cf ColorFormatter) LevelColor(l zapcore.Level) Color {
	switch l {
	case zapcore.DebugLevel:
		return ColorCyan
	case zapcore.InfoLevel:
		return ColorBlue
	case zapcore.WarnLevel:
		return ColorYellow
	case zapcore.ErrorLevel:
		return ColorRed
	case zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		return ColorMagenta
	default:
		return ColorNone
	}
}

func (cf ColorFormatter) Format(w io.Writer, entry zapcore.Entry, fields []zapcore.Field) {
	switch {
	case cf.Bold:
		fmt.Fprintf(w, cf.LevelColor(entry.Level).Bold())
	case cf.Reset:
		fmt.Fprintf(w, ResetColor())
	default:
		fmt.Fprintf(w, cf.LevelColor(entry.Level).Normal())
	}
}

type StringFormatter struct {
	Value string
}

// Format StringFormatter 直接就是将字符串 Value 写入到 io.Writer 中。
func (s StringFormatter) Format(w io.Writer, entry zapcore.Entry, fields []zapcore.Field) {
	fmt.Fprintf(w, "%s", s.Value)
}
