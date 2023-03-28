package cenc

import (
	"fmt"
	"io"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"

	"go.uber.org/zap/zapcore"
)

// formatRegexp 匹配案例：
// %{color:red}%{level:debug}  匹配结果：共找到两处匹配：%{color:red}和%{level:debug}；
// %{color}%{messagee}  匹配结果：共找到一处匹配：%{color}；
// %{id:123}xxxx%{module::p2p}  匹配结果：共找到两处匹配：%{id:123}和%{module::p2p}。
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
	case "level":
		return newLevelFormatter(format), nil
	case "message":
		return newMessageFormatter(format), nil
	case "shortfunc":
		return newShortFuncFormatter(format), nil
	case "time":
		return newTimeFormatter(format), nil
	case "id":
		return newSequenceFormatter(format), nil
	case "module":
		return newModuleFormatter(format), nil
	default:
		return nil, fmt.Errorf("unknown verb: %s", verb)
	}
}

// => MultiFormatter

type MultiFormatter struct {
	mutex sync.Mutex
	formatters []Formatter
}

func NewMultiFormatter(formatters ...Formatter) *MultiFormatter {
	return &MultiFormatter{
		formatters: formatters,
	}
}

func (mf *MultiFormatter) Format(w io.Writer, entry zapcore.Entry, fields []zapcore.Field) {
	mf.mutex.Lock()
	for _, formatter := range mf.formatters {
		formatter.Format(w, entry, fields)
	}
	mf.mutex.Unlock()
}

func (mf *MultiFormatter) SetFormatters(formatters []Formatter) {
	mf.mutex.Lock()
	mf.formatters = formatters
	mf.mutex.Unlock()
}

// => ColorFormatter

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
		fmt.Fprint(w, cf.LevelColor(entry.Level).Bold())
	case cf.Reset:
		fmt.Fprint(w, ResetColor())
	default:
		fmt.Fprint(w, cf.LevelColor(entry.Level).Normal())
	}
}

// => LevelFormatter

type LevelFormatter struct {
	FormatVerb string
}

func newLevelFormatter(fv string) LevelFormatter {
	return LevelFormatter{FormatVerb: "%" + stringOrDefault(fv, "s")}
}

func (lf LevelFormatter) Format(w io.Writer, entry zapcore.Entry, fields []zapcore.Field) {
	fmt.Fprintf(w, lf.FormatVerb, entry.Level.CapitalString())
}

// => MessageFormatter

type MessageFormatter struct {
	FormatVerb string
}

func newMessageFormatter(fv string) MessageFormatter {
	return MessageFormatter{FormatVerb: "%" + stringOrDefault(fv, "s")}
}

func (mf MessageFormatter) Format(w io.Writer, entry zapcore.Entry, fields []zapcore.Field) {
	fmt.Fprintf(w, mf.FormatVerb, strings.TrimRight(entry.Message, "\n"))
}

// => ShortFuncFormatter

type ShortFuncFormatter struct {
	FormatVerb string
}

func newShortFuncFormatter(fv string) ShortFuncFormatter {
	return ShortFuncFormatter{FormatVerb: "%" + stringOrDefault(fv, "s")}
}

func (sf ShortFuncFormatter) Format(w io.Writer, entry zapcore.Entry, fields []zapcore.Field) {
	f := runtime.FuncForPC(entry.Caller.PC)
	if f == nil {
		fmt.Fprintf(w, sf.FormatVerb, "(unknown)")
		return
	}
	fname := f.Name()
	funcIdx := strings.LastIndex(fname, ".")
	fmt.Fprintf(w, sf.FormatVerb, fname[funcIdx+1:])
}

// => TimeFormatter

type TimeFormatter struct {
	Layout string
}

func newTimeFormatter(layout string) TimeFormatter {
	return TimeFormatter{Layout: stringOrDefault(layout, "2006-01-02T15:04:05.999Z07:00")}
}

func (tf TimeFormatter) Format(w io.Writer, entry zapcore.Entry, fields []zapcore.Field) {
	fmt.Fprint(w, entry.Time.Format(tf.Layout))
}

// => SequenceFormatter

// 全局变量，供所有 SequenceFormatter 实例使用。
var sequence uint64

func SetSequence(s uint64) {
	atomic.StoreUint64(&sequence, s)
}

type SequenceFormatter struct {
	FormatVerb string
}

func newSequenceFormatter(fv string) SequenceFormatter {
	return SequenceFormatter{FormatVerb: "%" + stringOrDefault(fv, "d")}
}

func (sf SequenceFormatter) Format(w io.Writer, entry zapcore.Entry, fields []zapcore.Field) {
	fmt.Fprintf(w, sf.FormatVerb, atomic.AddUint64(&sequence, 1))
}

// => ModuleFormatter

type ModuleFormatter struct {
	FormatVerb string
}

func newModuleFormatter(fv string) ModuleFormatter {
	return ModuleFormatter{FormatVerb: "%" + stringOrDefault(fv, "s")}
}

func (mf ModuleFormatter) Format(w io.Writer, entry zapcore.Entry, fields []zapcore.Field) {
	fmt.Fprintf(w, mf.FormatVerb, entry.LoggerName)
}

// => StringFormatter

type StringFormatter struct {
	Value string
}

// Format StringFormatter 直接就是将字符串 Value 写入到 io.Writer 中。
func (s StringFormatter) Format(w io.Writer, entry zapcore.Entry, fields []zapcore.Field) {
	fmt.Fprintf(w, "%s", s.Value)
}

func stringOrDefault(str, dlt string) string {
	if str != "" {
		return str
	}
	return dlt
}
