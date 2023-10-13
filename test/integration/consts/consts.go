package consts

import "time"

const (
	// waitTick is the default timeout tick interval for checking on ingress resources.
	WaitTick = 250 * time.Millisecond

	// ingressWait is the default amount of time to wait for any particular ingress resource to be provisioned.
	IngressWait = time.Minute * 3
)
