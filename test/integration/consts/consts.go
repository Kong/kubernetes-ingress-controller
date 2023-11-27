package consts

import "time"

const (
	// WaitTick is the default timeout tick interval for checking on resources.
	WaitTick = 250 * time.Millisecond

	// IngressWait is the default amount of time to wait for any particular ingress resource to be provisioned.
	IngressWait = 3 * time.Minute

	// StatusWait is the default amount of time to wait for object statuses to fulfill a provided predicate.
	StatusWait = 3 * time.Minute
)
