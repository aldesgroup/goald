package core

import (
	"bytes"
	"math/rand"
	"strings"
	"unicode"
)

// ------------------------------------------------------------------------------------------------
// Strings
// ------------------------------------------------------------------------------------------------

// pascalToSeparated converts a PascalCase string to seperated case like kebab-case or snake_case
func pascalToSeparated(s string, sep rune) string {
	if s == "" {
		return s
	}

	var buffer bytes.Buffer

	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 && i+1 < len(s) && (unicode.IsLower(rune(s[i+1])) || unicode.IsLower(rune(s[i-1]))) {
				buffer.WriteRune(sep)
			}
			buffer.WriteRune(unicode.ToLower(r))
		} else {
			buffer.WriteRune(r)
		}
	}

	return buffer.String()
}

// PascalToKebab converts a PascalCase string to kebab-case.
func PascalToKebab(s string) string {
	return pascalToSeparated(s, '-')
}

// PascalToSnake converts a PascalCase string to snake_case.
func PascalToSnake(s string) string {
	return pascalToSeparated(s, '_')
}

// separatedToPascal converts a separated case string to PascalCase
func separatedToPascal(s string, sep string) string {
	words := []string{}

	for _, word := range strings.Split(s, sep) {
		if len(word) > 0 {
			// Capitalize the first letter of each word manually.
			words = append(words, strings.ToUpper(word[:1])+word[1:])
		}
	}

	return strings.Join(words, "")
}

// KebabToPascal converts a kebab-case string to PascalCase
func KebabToPascal(s string) string {
	return separatedToPascal(s, "-")
}

// PascalToCamel converts a PascalCase string to camelCase.
func PascalToCamel(s string) string {
	runes := []rune(s)
	size := len(runes)

	if size == 0 || unicode.IsLower(runes[0]) {
		return s
	}

	result := []rune{unicode.ToLower(runes[0])}

	var i int
	for i = 1; i < size; i++ {
		if unicode.IsLower(runes[i]) || i < size-1 && unicode.IsUpper(runes[i]) && unicode.IsLower(runes[i+1]) {
			break
		}

		result = append(result, unicode.ToLower(runes[i]))
	}

	return string(result) + string(runes[i:size])
}

// ToPascal converts string to PascalCase
func ToPascal(s string) string {
	if s == "" {
		return s
	}

	return string(unicode.ToUpper(rune(s[0]))) + s[1:]
}

const charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// RandomString generates a random string of the given length
func RandomString(length int) string {
	// Create a slice to store the characters of the random string
	randomString := make([]byte, length)

	// Populate the slice with random characters
	for i := 0; i < length; i++ {
		randomString[i] = charSet[rand.Intn(len(charSet))]
	}

	// Convert the slice to a string and return
	return string(randomString)
}
