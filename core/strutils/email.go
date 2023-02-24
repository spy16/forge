package strutils

import (
	"crypto/md5"
	"fmt"
	"net/mail"
	"strings"
)

// GravatarURL returns a valid Gravatar URL for the given email id. Size parameter
// is passed to the URL to make thumbnail of given size.
func GravatarURL(email string, size int) string {
	if size <= 0 || size >= 2048 {
		size = 128
	}
	email = strings.TrimSpace(email)
	if !IsValidEmail(email) {
		return fmt.Sprintf("https://www.gravatar.com/avatar?s=%d", size)
	}
	hash := md5.Sum([]byte(email))
	return fmt.Sprintf("https://www.gravatar.com/avatar/%x?s=%d", hash, size)
}

// IsValidEmail returns true if the given email is valid.
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
