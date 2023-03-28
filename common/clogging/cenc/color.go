package cenc

import "fmt"

type Color uint8

const ColorNone Color = 0

const (
	ColorBlack Color = iota + 30
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
)

// Normal 在终端打印出对应的颜色，一般来讲，打印不同颜色的内容到终端上，其方法如下所示：
// fmt.Printf("\x1b[%dmhello world 31: 红 \x1b[0m\n", 31)
// 可以看到，格式化字符串里，不仅包含 Normal 方法的 `\x1b[%dm`，还包含后半段 `\x1b[0m`，这代表打印完红色之后，
// 就将往终端输出内容的颜色重置为默认值，Normal 方法则缺少了重置这一步骤，这意味着，我们只要调用一次 Normal 方
// 法，在不调用重置方法 ResetColor 的前提下，可以一直打印红色的内容。
func (c Color) Normal() string {
	return fmt.Sprintf("\x1b[%dm", c)
}

func (c Color) Bold() string {
	if c == ColorNone {
		return c.Normal()
	}
	return fmt.Sprintf("\x1b[%d;1m", c)
}

func ResetColor() string { return ColorNone.Normal() }
