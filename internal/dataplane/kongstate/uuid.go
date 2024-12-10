package kongstate

// StaticUUIDGenerator is a UUIDGenerator that always returns the same UUID. It is used for testing.
type StaticUUIDGenerator struct {
	UUID string
}

func (s StaticUUIDGenerator) NewString() string {
	return s.UUID
}
