package goald

import (
	"os"
)

// ------------------------------------------------------------------------------------------------
// Files
// ------------------------------------------------------------------------------------------------

// FileExists checks if a file exists at the specified path.
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// DirExists checks if a file exists at the specified path.
func DirExists(dirPath string) bool {
	info, err := os.Stat(dirPath)
	return !os.IsNotExist(err) && info.IsDir()
}
