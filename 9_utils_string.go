package goald

import (
	"bytes"
	"strings"
	"unicode"
)

// ------------------------------------------------------------------------------------------------
// Strings
// ------------------------------------------------------------------------------------------------

// PascalToSnake converts a PascalCase string to snake_case.
func PascalToSnake(s string) string {
	if s == "" {
		return s
	}

	var buffer bytes.Buffer

	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 && i+1 < len(s) && unicode.IsLower(rune(s[i+1])) {
				buffer.WriteRune('_')
			}
			buffer.WriteRune(unicode.ToLower(r))
		} else {
			buffer.WriteRune(r)
		}
	}

	return buffer.String()
}

// SnakeToPascal converts a snake_case string to PascalCase
func SnakeToPascal(s string) string {
	words := []string{}

	for _, word := range strings.Split(s, "_") {
		if len(word) > 0 {
			// Capitalize the first letter of each word manually.
			words = append(words, strings.ToUpper(word[:1])+word[1:])
		}
	}

	return strings.Join(words, "")
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