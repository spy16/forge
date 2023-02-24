package strutils_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/spy16/forge/core/strutils"
)

func TestIsValidEmail(t *testing.T) {
	t.Parallel()

	table := []struct {
		Email string
		Want  bool
	}{
		{"bob@bobmail.com", true},
		{"bob@bobmail", true},
		{"bobmail.com", false},
		{"", false},
		{"foo", false},
	}

	for i, tt := range table {
		t.Run(fmt.Sprintf("Case#%d", i), func(t *testing.T) {
			got := strutils.IsValidEmail(tt.Email)
			assert.Equal(t, tt.Want, got)
		})
	}
}

func TestGravatarURL(t *testing.T) {
	t.Parallel()

	u := strutils.GravatarURL("bob@bobmail.com", 128)
	assert.Equal(t, u, "https://www.gravatar.com/avatar/8d21641bdea72d47e3344c9c0528c208?s=128")

	u = strutils.GravatarURL("bob@bobmail.com", 64)
	assert.Equal(t, u, "https://www.gravatar.com/avatar/8d21641bdea72d47e3344c9c0528c208?s=64")
}
