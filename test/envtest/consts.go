package envtest

import "time"

const (
	// waitTime is a time to wait for a condition to be met.
	waitTime = 3 * time.Second
	// tickTime is a time to wait between condition's checks.
	tickTime = 100 * time.Millisecond
)
