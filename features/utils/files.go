package utils

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"time"
)

// ------------------------------------------------------------------------------------------------
// Files
// ------------------------------------------------------------------------------------------------

// FileExists checks if a file exists at the specified path.
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// Returns the given file's modification time, or panics
func EnsureModTime(filePath string) time.Time {
	info, err := os.Stat(filePath)
	PanicErrf(err, "Could not get modification time for file '%s'", filePath)
	return info.ModTime()
}

// DirExists checks if a file exists at the specified path.
func DirExists(dirPath string) bool {
	info, err := os.Stat(dirPath)
	return !os.IsNotExist(err) && info.IsDir()
}

// EnsureDir makes sure the directory with the given path elements exists
func EnsureDir(pathElem ...string) string {
	dirname := path.Join(pathElem...)

	PanicErrf(os.MkdirAll(dirname, 0o777), "Could not create directory '%s'", dirname)

	return dirname
}

// WriteToFile
func WriteToFile(content string, filepaths ...string) {
	// creating the file
	fileName := path.Join(filepaths...)

	// creating the missing directory if needed
	if dir := path.Dir(fileName); path.Base(fileName) != fileName {
		EnsureDir(dir)
	}

	// creating the file
	file, errCreate := os.Create(fileName)
	if errCreate != nil {
		Panicf("Could not create file %s; cause: %s", fileName, errCreate)
	}

	// ensuring we've got no leak
	defer func() {
		if errClose := file.Close(); errClose != nil {
			slog.Error(fmt.Sprintf("Could not properly close file %s; cause: %s", fileName, errClose))
			os.Exit(1)
		}
	}()

	// writing to file
	if _, errWrite := file.WriteString(content); errWrite != nil {
		PanicErrf(errWrite, "Could not write file '%s'", fileName)
	}
}
