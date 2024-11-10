package security_test

import (
	"github.com/mandarine-io/Backend/internal/helper/security"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SecurityUtil_HashPassword(t *testing.T) {
	t.Run(
		"hash password successfully", func(t *testing.T) {
			password := "strongpassword123"
			hash, err := security.HashPassword(password)
			assert.NoError(t, err)
			assert.NotEmpty(t, hash)
		},
	)

	t.Run(
		"hash password and compare", func(t *testing.T) {
			password := "strongpassword123"
			hash, err := security.HashPassword(password)
			assert.NoError(t, err)

			match := security.CheckPasswordHash(password, hash)
			assert.True(t, match)
		},
	)

	t.Run(
		"compare with incorrect password", func(t *testing.T) {
			password := "strongpassword123"
			hash, err := security.HashPassword(password)
			assert.NoError(t, err)

			wrongPassword := "wrongpassword"
			match := security.CheckPasswordHash(wrongPassword, hash)
			assert.False(t, match)
		},
	)
}
