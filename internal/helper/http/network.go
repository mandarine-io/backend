package http

import "strings"

func IsPublicOrigin(origin string) bool {
	origin = strings.TrimRight(origin, "/")
	origin = strings.TrimPrefix(origin, "https://")
	origin = strings.TrimPrefix(origin, "http://")

	return !strings.HasPrefix(origin, "localhost")
}
