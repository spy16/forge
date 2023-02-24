package strutils

import "unicode"

// Ptr returns the given string as a string-pointer.
func Ptr(s string) *string { return &s }

func OneOf(item string, arr []string) bool {
	for _, s := range arr {
		if item == s {
			return true
		}
	}
	return false
}

// SnakeCase converts the given string to snake_case version.
func SnakeCase(input string) string {
	var output string
	for i, ch := range input {
		if unicode.IsUpper(ch) {
			if i > 0 && input[i-1] != ' ' && !unicode.IsUpper(rune(input[i-1])) {
				output += "_"
			}
			output += string(unicode.ToLower(ch))
		} else {
			output += string(ch)
		}
	}
	return output
}
