package store

// NotFoundError error is returned when a lookup results in no resource.
// This type is meant to be used for error handling using `errors.As()`.
type NotFoundError struct {
	Message string
}

func (e NotFoundError) Error() string {
	if e.Message == "" {
		return "not found"
	}
	return e.Message
}
