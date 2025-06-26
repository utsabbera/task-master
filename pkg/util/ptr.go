package util

// Ptr is a utility function that returns a pointer to the given value.
func Ptr[T any](v T) *T {
	return &v
}

// Val returns the value pointed to by the given pointer.
func Val[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}

	return *p
}
