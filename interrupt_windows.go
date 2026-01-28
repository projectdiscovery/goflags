//go:build windows

package goflags

import "syscall"

var (
	kernel32                     = syscall.NewLazyDLL("kernel32.dll")
	procGenerateConsoleCtrlEvent = kernel32.NewProc("GenerateConsoleCtrlEvent")
)

func sendInterrupt() {
	// CTRL_BREAK_EVENT = 1, sends to all processes in the console group
	// Using 0 as process group ID sends to all processes attached to the console
	procGenerateConsoleCtrlEvent.Call(1, 0)
}
