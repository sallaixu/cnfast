//go:build windows
// +build windows

package colors

import (
	"os"
	"syscall"
	"unsafe"
)

// isTTY 判断 stdout 是否为终端
func isTTY() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

// enableVT 为 Windows 控制台启用虚拟终端序列支持
// 使 ANSI 转义码在 cmd.exe / PowerShell 中正常渲染
func enableVT() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	procGetStdHandle := kernel32.NewProc("GetStdHandle")
	procSetConsoleMode := kernel32.NewProc("SetConsoleMode")
	procGetConsoleMode := kernel32.NewProc("GetConsoleMode")

	handle, _, _ := procGetStdHandle.Call(uintptr(0xfffffff5)) // STD_OUTPUT_HANDLE = -11
	const enableVirtualTerminalProcessing = 0x0004

	var mode uint32
	if ret, _, _ := procGetConsoleMode.Call(handle, uintptr(unsafe.Pointer(&mode))); ret == 0 {
		return
	}
	procSetConsoleMode.Call(handle, uintptr(mode|enableVirtualTerminalProcessing))
}
