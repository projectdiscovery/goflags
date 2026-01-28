//go:build !windows

package goflags

import "os"

func sendInterrupt() {
	if p, err := os.FindProcess(os.Getpid()); err == nil {
		_ = p.Signal(os.Interrupt)
	}
}
