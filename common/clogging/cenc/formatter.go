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
		start, end := m[0], m[1] // start 匹配到的是 `%` 的下标位置，end 匹配到的是 `}` 后一个字符下标的位置。
		verbStart, verbEnd := m[2], m[3] // verbStart 匹配到的是字母 `l` 的下标位置，verbEnd 匹配到的是冒号 `:` 的下标位置。
		formatStart, formatEnd := m[4], m[5] // formatStart 匹配到的是点号 `.` 的下标位置，formatEnd 匹配到的是 `}` 下标的位置。

		if start > cursor {
			// 代表 `%{` 前面还有内容
			formatters = append(formatters, StringFormatter{Value: spec[cursor:start]})
		}
	}
}

type StringFormatter struct {
	Value string
}

// Format StringFormatter 直接就是将字符串 Value 写入到 io.Writer 中。
func (s StringFormatter) Format(w io.Writer, entry zapcore.Entry, fields []zapcore.Field) {
	fmt.Fprintf(w, "%s", s.Value)
}