package http_test

import (
	"github.com/mandarine-io/Backend/internal/helper/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HttpUtil_IsPublicOrigin(t *testing.T) {
	t.Run(
		"origin with transport and localhost", func(t *testing.T) {
			origin := "http://localhost"
			result := http.IsPublicOrigin(origin)
			assert.False(t, result)
		},
	)

	t.Run(
		"origin with https and localhost", func(t *testing.T) {
			origin := "https://localhost"
			result := http.IsPublicOrigin(origin)
			assert.False(t, result)
		},
	)

	t.Run(
		"origin with transport and localhost with port", func(t *testing.T) {
			origin := "http://localhost:8080"
			result := http.IsPublicOrigin(origin)
			assert.False(t, result)
		},
	)

	t.Run(
		"origin with public persistence", func(t *testing.T) {
			origin := "https://example.com"
			result := http.IsPublicOrigin(origin)
			assert.True(t, result)
		},
	)

	t.Run(
		"origin with transport and public persistence", func(t *testing.T) {
			origin := "http://example.com"
			result := http.IsPublicOrigin(origin)
			assert.True(t, result)
		},
	)

	t.Run(
		"origin without protocol", func(t *testing.T) {
			origin := "example.com"
			result := http.IsPublicOrigin(origin)
			assert.True(t, result)
		},
	)

	t.Run(
		"origin with trailing slash", func(t *testing.T) {
			origin := "https://example.com/"
			result := http.IsPublicOrigin(origin)
			assert.True(t, result)
		},
	)

	t.Run(
		"localhost with trailing slash", func(t *testing.T) {
			origin := "http://localhost/"
			result := http.IsPublicOrigin(origin)
			assert.False(t, result)
		},
	)
}
