//go:build linux
// +build linux

package core

func IsWindows() bool {
	return false
}

func IsLinux() bool {
	return true
}

func CopyCmd() string {
	return "cp -r"
}

func RemoveCmd() string {
	return "rm -fr"
}

func MoveCmd() string {
	return "mv"
}
