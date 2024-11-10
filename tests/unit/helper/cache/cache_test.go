package cache_test

import (
	"github.com/mandarine-io/Backend/internal/helper/cache"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CacheUtil_CreateCacheKey(t *testing.T) {
	t.Run(
		"single part", func(t *testing.T) {
			key := cache.CreateCacheKey("part1")
			assert.Equal(t, "part1", key)
		},
	)

	t.Run(
		"multiple parts", func(t *testing.T) {
			key := cache.CreateCacheKey("part1", "part2", "part3")
			expectedKey := "part1.part2.part3"
			assert.Equal(t, expectedKey, key)
		},
	)

	t.Run(
		"empty parts", func(t *testing.T) {
			key := cache.CreateCacheKey()
			assert.Equal(t, "", key)
		},
	)

	t.Run(
		"parts with empty strings", func(t *testing.T) {
			key := cache.CreateCacheKey("part1", "", "part3")
			expectedKey := "part1..part3"
			assert.Equal(t, expectedKey, key)
		},
	)

	t.Run(
		"parts with dots", func(t *testing.T) {
			key := cache.CreateCacheKey("part1.", ".part2", "part3.")
			expectedKey := "part1...part2.part3."
			assert.Equal(t, expectedKey, key)
		},
	)
}
