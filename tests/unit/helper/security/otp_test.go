package security_test

import (
	"mandarine/internal/api/helper/security"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SecurityUtil_GenerateOTP(t *testing.T) {
	t.Run(
		"generate OTP with valid length", func(t *testing.T) {
			length := 6
			otp, err := security.GenerateOTP(length)
			assert.NoError(t, err)
			assert.Equal(t, length, len(otp))
		},
	)

	t.Run(
		"generate OTP with zero length", func(t *testing.T) {
			length := 0
			otp, err := security.GenerateOTP(length)
			assert.NoError(t, err)
			assert.Equal(t, length, len(otp))
		},
	)

	t.Run(
		"generate OTP with negative length", func(t *testing.T) {
			length := -1
			otp, err := security.GenerateOTP(length)
			assert.EqualError(t, err, "OTP length is negative")
			assert.Equal(t, "", otp)
		},
	)
}
