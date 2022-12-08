package address

// Of returns an address of provided argument.
func Of[T any](v T) *T {
	return &v
}
