//go:build !windows
// +build !windows

package colors

import (
	"os"
)

// isTTY 判断 stdout 是否为终端
func isTTY() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

// enableVT 在非 Windows 平台无需额外操作
func enableVT() {}
