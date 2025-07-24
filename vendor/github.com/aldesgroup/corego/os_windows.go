//go:build windows
// +build windows

package core

func IsWindows() bool {
	return true
}

func IsLinux() bool {
	return false
}

func CopyCmd() string {
	return "powershell -Command Copy-Item -Recurse"
}

func RemoveCmd() string {
	return "powershell -Command Remove-Item -Recurse -Force"
}

func MoveCmd() string {
	return "powershell -Command Move-Item"
}
