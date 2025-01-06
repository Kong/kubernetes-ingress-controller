package integration

import (
	"time"

	"github.com/kong/kubernetes-ingress-controller/v3/test/integration/consts"
)

// -----------------------------------------------------------------------------
// Testing Timeouts
// -----------------------------------------------------------------------------

const (
	// waitTick is the default timeout tick interval for checking on ingress resources.
	waitTick = consts.WaitTick

	// ingressWait is the default amount of time to wait for any particular ingress resource to be provisioned.
	ingressWait = consts.IngressWait

	// httpcTimeout is the default client timeout for HTTP clients used in tests.
	httpcTimeout = time.Second * 3

	// statusWait is a const duration used in test assertions like .Eventually to
	// wait for object statuses to fulfill a provided predicate.
	statusWait = consts.StatusWait
)
