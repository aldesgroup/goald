package utils

import (
	"log"
	"os"
	"path"
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

// WriteToFile
func WriteToFile(content string, filepaths ...string) {
	// creating the file
	fileName := path.Join(filepaths...)

	// creating the missing directory if needed
	if dir := path.Dir(fileName); path.Base(fileName) != fileName {
		PanicErrf(os.MkdirAll(dir, 0o777), "Could not create directory '%s'", dir)
	}

	// creating the file
	file, errCreate := os.Create(fileName)
	if errCreate != nil {
		Panicf("Could not create file %s; cause: %s", fileName, errCreate)
	}

	// ensuring we've got no leak
	defer func() {
		if errClose := file.Close(); errClose != nil {
			log.Fatalf("Could not properly close file %s; cause: %s", fileName, errClose)
		}
	}()

	// writing to file
	if _, errWrite := file.WriteString(content); errWrite != nil {
		PanicErrf(errWrite, "Could not write file '%s'", fileName)
	}
}
