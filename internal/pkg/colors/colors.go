// Package colors 提供 ANSI 颜色输出支持
// 支持 NO_COLOR 环境变量与非 TTY 场景下的自动降级
package colors

import (
	"os"
)

// ANSI 转义码
const (
	reset         = "\033[0m"
	bold          = "\033[1m"
	faint         = "\033[2m"
	brightRed     = "\033[91m"
	brightGreen   = "\033[92m"
	brightYellow  = "\033[93m"
	brightBlue    = "\033[94m"
	brightMagenta = "\033[95m"
	brightCyan    = "\033[96m"
	brightWhite   = "\033[97m"
)

// enabled 是否启用颜色输出
var enabled = false

// terminal 标准输出是否为终端
var terminal = false

func init() {
	terminal = isTTY()
	enabled = terminal && os.Getenv("NO_COLOR") == ""
	enableVT()
}

// IsTerminal 返回标准输出是否为终端（用于决定是否使用 \\r 刷新进度等交互行为）
func IsTerminal() bool { return terminal }

// wrap 包裹 ANSI 转义码，禁用时原样返回
func wrap(code, s string) string {
	if !enabled {
		return s
	}
	return code + s + reset
}

// Bold 加粗
func Bold(s string) string { return wrap(bold, s) }

// Faint 弱化（浅色/暗淡显示）
func Faint(s string) string { return wrap(faint, s) }

// Red 红色
func Red(s string) string { return wrap(brightRed, s) }

// Green 绿色
func Green(s string) string { return wrap(brightGreen, s) }

// Yellow 黄色
func Yellow(s string) string { return wrap(brightYellow, s) }

// Blue 蓝色
func Blue(s string) string { return wrap(brightBlue, s) }

// Magenta 品红
func Magenta(s string) string { return wrap(brightMagenta, s) }

// Cyan 青色
func Cyan(s string) string { return wrap(brightCyan, s) }

// White 白色（加粗亮白）
func White(s string) string { return wrap(brightWhite, s) }

// Success 成功提示（绿色 label）
func Success(s string) string { return wrap(brightGreen, "✔") + " " + s }

// Error 错误提示（红色 label）
func Error(s string) string { return wrap(brightRed, "✘") + " " + s }

// Warn 警告提示（黄色 label）
func Warn(s string) string { return wrap(brightYellow, "⚠") + " " + s }

// Info 信息提示（青色 label）
func Info(s string) string { return wrap(brightCyan, "→") + " " + s }

// Enabled 返回颜色是否启用
func Enabled() bool { return enabled }

// Label 返回带颜色的短标签
func Label(s string) string { return wrap(brightCyan, s) }
