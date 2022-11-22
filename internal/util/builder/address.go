package builder

func addressOf[T any](v T) *T {
	return &v
}
