package geo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ToMorse(t *testing.T) {
	testCases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "173.177.164.160",
			in:   "173.177.164.160",
			want: ".---- --... ...-- .-.-.- .---- --... --... .-.-.- .---- -.... ....- .-.-.- .---- -.... -----",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := ToMorse(tc.in)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tc.want, got, "Failed encoding to Morse code. Should be equal")
		})
	}
}
