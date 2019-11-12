package geo

import (
	"github.com/alwindoss/morse"
	"strings"
)

// ToMorse
func ToMorse(s string) (string, error) {
	h := morse.NewHacker()
	m, err := h.Encode(strings.NewReader(s))
	if err != nil {
		return "", err
	}

	return string(m), nil
}
