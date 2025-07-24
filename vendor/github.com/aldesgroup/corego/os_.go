//go:build !windows && !linux
// +build !windows,!linux

package core

import (
	"runtime"
)

func IsWindows() bool {
	PanicMsg("Platform '%s' is currently not supported right now!", runtime.GOOS)
	return false
}

func IsLinux() bool {
	PanicMsg("Platform '%s' is currently not supported right now!", runtime.GOOS)
	return false
}

func CopyCmd() string {
	PanicMsg("Platform '%s' is currently not supported right now!", runtime.GOOS)
	return ""
}

func RemoveCmd() string {
	PanicMsg("Platform '%s' is currently not supported right now!", runtime.GOOS)
	return ""
}

func MoveCmd() string {
	PanicMsg("Platform '%s' is currently not supported right now!", runtime.GOOS)
	return ""
}
