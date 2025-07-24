package core

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path"
	"strings"
	"time"
)

// ------------------------------------------------------------------------------------------------
// Files
// ------------------------------------------------------------------------------------------------

// FileExists checks if a file exists at the specified path.
func FileExists(pathParts ...string) bool {
	fullpath := path.Join(pathParts...)
	info, err := os.Stat(fullpath)
	return !os.IsNotExist(err) && !info.IsDir()
}

// DirExists checks if a file exists at the specified path.
func DirExists(pathParts ...string) bool {
	fullpath := path.Join(pathParts...)
	info, err := os.Stat(fullpath)
	return !os.IsNotExist(err) && info.IsDir()
}

// Returns the given file's modification time, or panics
func EnsureModTime(filePath string) time.Time {
	info, err := os.Stat(filePath)
	PanicMsgIfErr(err, "Could not get modification time for file '%s'", filePath)
	return info.ModTime()
}

// EnsureDir makes sure the directory with the given path elements exists
func EnsureDir(pathElem ...string) string {
	dirname := path.Join(pathElem...)

	PanicMsgIfErr(os.MkdirAll(dirname, 0o777), "Could not create directory '%s'", dirname)

	return dirname
}

// WriteToFile writes the given content to the file with the given path
func WriteToFile(content string, filepaths ...string) {
	// creating the file
	fileName := path.Join(filepaths...)

	// creating the missing directory if needed
	if dir := path.Dir(fileName); path.Base(fileName) != fileName {
		EnsureDir(dir)
	}

	// creating the file
	file, errCreate := os.Create(fileName)
	PanicMsgIfErr(errCreate, "Could not create file %s", fileName)

	// ensuring we've got no leak
	defer func() {
		if errClose := file.Close(); errClose != nil {
			slog.Error(fmt.Sprintf("Could not properly close file %s; cause: %s", fileName, errClose))
			os.Exit(1)
		}
	}()

	// writing to file
	if _, errWrite := file.WriteString(content); errWrite != nil {
		PanicMsgIfErr(errWrite, "Could not write file '%s'", fileName)
	}
}

// EnsureNoDir removes the directory with the given path elements
func EnsureNoDir(pathElem ...string) string {
	dirname := path.Join(pathElem...)
	PanicMsgIfErr(os.RemoveAll(dirname), "Could not remove folder '%s'", dirname)

	return dirname
}

// WriteBytesToFile writes the given bytes to the file with the given path
func WriteBytesToFile(filename string, bytes []byte) {
	if filename != path.Base(filename) {
		EnsureDir(path.Dir(filename))
	}

	PanicMsgIfErr(os.WriteFile(filename, bytes, 0o644), "Could not write to file '%s'", filename)
}

// WriteStringToFile writes the given string to the file with the given path
func WriteStringToFile(filename string, content string, params ...any) {
	WriteBytesToFile(filename, []byte(fmt.Sprintf(content, params...)))
}

// WriteJsonObjToFile writes the given JSON object to the file with the given path
func WriteJsonObjToFile(filename string, obj any) {
	jsonBytes, errMarshal := json.MarshalIndent(obj, "", "\t")
	PanicMsgIfErr(errMarshal, "Could not JSON-marshal to file '%s'", filename)
	WriteBytesToFile(filename, jsonBytes)
}

// ReadFile reads the file with the given path and returns the bytes
func ReadFile(filename string, failIfNotExist bool) []byte {
	if !FileExists(filename) {
		PanicMsgIf(failIfNotExist, "File '%s' cannot be found!", filename)
		return nil
	}

	fileBytes, errRead := os.ReadFile(filename)
	PanicMsgIfErr(errRead, "Could not read file '%s'", filename)
	return fileBytes
}

// ReadFileToJson reads the file with the given path and unmarshals the JSON object
func ReadFileToJson[T any, Y *T](filename string, obj Y, failIfNotExist bool) Y {
	if fileBytes := ReadFile(filename, failIfNotExist); fileBytes != nil {
		PanicMsgIfErr(json.Unmarshal(fileBytes, obj), "Could not JSON-unmarshal file '%s'", filename)
	}
	return obj
}

// ReplaceInFile replaces the given replacements in the file with the given path
func ReplaceInFile(filename string, replacements map[string]string) {
	fileContent := string(ReadFile(filename, true))
	for replace, by := range replacements {
		fileContent = strings.ReplaceAll(fileContent, replace, by)
	}
	WriteStringToFile(filename, "%s", fileContent)
}
