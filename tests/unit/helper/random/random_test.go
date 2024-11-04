package random_test

import (
	"github.com/mandarine-io/Backend/internal/api/helper/random"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SecurityUtil_GenerateRandomNumber(t *testing.T) {
	t.Run(
		"generate OTP with valid length", func(t *testing.T) {
			length := 6
			otp, err := random.GenerateRandomNumber(length)
			assert.NoError(t, err)
			assert.Equal(t, length, len(otp))
		},
	)

	t.Run(
		"generate OTP with zero length", func(t *testing.T) {
			length := 0
			otp, err := random.GenerateRandomNumber(length)
			assert.NoError(t, err)
			assert.Equal(t, length, len(otp))
		},
	)

	t.Run(
		"generate OTP with negative length", func(t *testing.T) {
			length := -1
			otp, err := random.GenerateRandomNumber(length)
			assert.EqualError(t, err, "OTP length is negative")
			assert.Equal(t, "", otp)
		},
	)
}
