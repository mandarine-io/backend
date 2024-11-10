package ref

func SafeDeref[T any](p *T) T {
	if p == nil {
		var v T
		return v
	}
	return *p
}

func SafeRef[T any](v T) *T {
	return &v
}
