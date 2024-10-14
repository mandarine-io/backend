package cache

import "strings"

func CreateCacheKey(parts ...string) string {
	return strings.Join(parts, ".")
}
