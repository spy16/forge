package strutils

import (
	"math/rand"
)

const (
	CharsetNums  = "0123456789"
	CharsetLower = "abcdefghijklmnopqrstuvwxyz"
	CharsetUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// RandStr returns a random string of length 'n'. Characters are picked from
// the charset if provided or alphanumeric characters are used.
func RandStr(n int, charset ...string) string {
	chars := CharsetLower + CharsetUpper + CharsetNums
	if len(charset) >= 1 {
		chars = ""
		for _, s := range charset {
			chars += s
		}
	}

	s := make([]byte, n)
	for i := range s {
		s[i] = chars[rand.Intn(len(chars))]
	}
	return string(s)
}
