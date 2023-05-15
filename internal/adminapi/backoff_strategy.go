package adminapi

// UpdateBackoffStrategy keeps state of an update backoff strategy.
type UpdateBackoffStrategy interface {
	// CanUpdate tells whether we're allowed to make an update attempt for a given config hash.
	// In case it returns false, the second return value is a human-readable explanation of why the update cannot
	// be performed at this point in time.
	CanUpdate([]byte) (bool, string)

	// RegisterUpdateSuccess resets the backoff strategy, effectively making it allow next update straight away.
	RegisterUpdateSuccess()

	// RegisterUpdateFailure registers an update failure along with its failure reason passed as a generic error, and
	// a config hash that we failed to push.
	RegisterUpdateFailure(failureReason error, configHash []byte)
}
