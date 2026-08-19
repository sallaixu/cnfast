package util

import "regexp"

// ansiRe 匹配 ANSI 转义序列（用于计算可见宽度）
var ansiRe = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// cjkRe 匹配中日韩等全角字符
var cjkRe = regexp.MustCompile(`[\x{3000}-\x{303F}\x{3040}-\x{30FF}\x{3400}-\x{4DBF}\x{4E00}-\x{9FFF}\x{F900}-\x{FAFF}\x{FF00}-\x{FF60}]`)

// VisualLen 计算去除 ANSI 转义码后的显示宽度（CJK 字符按 2 计）
func VisualLen(s string) int {
	clean := ansiRe.ReplaceAllString(s, "")
	w := 0
	for _, r := range clean {
		if cjkRe.MatchString(string(r)) {
			w += 2
		} else {
			w++
		}
	}
	return w
}

// Pad 将字符串补齐到指定显示宽度（右侧补空格），保持颜色不影响对齐
func Pad(s string, width int) string {
	n := VisualLen(s)
	if n >= width {
		return s
	}
	for i := n; i < width; i++ {
		s += " "
	}
	return s
}

// Truncate 按显示宽度截断字符串，超出部分以省略号结尾
func Truncate(s string, maxWidth int) string {
	if VisualLen(s) <= maxWidth {
		return s
	}

	var out []rune
	w := 0
	for _, r := range []rune(s) {
		rw := 1
		if cjkRe.MatchString(string(r)) {
			rw = 2
		}
		if w+rw > maxWidth-1 {
			break
		}
		out = append(out, r)
		w += rw
	}
	return string(out) + "…"
}
