package gateway

import (
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// -----------------------------------------------------------------------------
// Logging Utilities
// -----------------------------------------------------------------------------

// debug is an alias for the longer log.V(util.DebugLevel).Info for convenience
func debug(log logr.Logger, obj client.Object, msg string, keysAndValues ...interface{}) {
	keysAndValues = append([]interface{}{
		"namespace", obj.GetNamespace(),
		"name", obj.GetName(),
	}, keysAndValues...)
	log.V(util.DebugLevel).Info(msg, keysAndValues...)
}

// info is an alias for the longer log.V(util.InfoLevel).Info for convenience
func info(log logr.Logger, obj client.Object, msg string, keysAndValues ...interface{}) {
	keysAndValues = append([]interface{}{
		"namespace", obj.GetNamespace(),
		"name", obj.GetName(),
	}, keysAndValues...)
	log.V(util.InfoLevel).Info(msg, keysAndValues...)
}
