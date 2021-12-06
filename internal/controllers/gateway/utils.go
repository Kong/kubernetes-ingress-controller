package gateway

import (
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// -----------------------------------------------------------------------------
// Logging Utilities
// -----------------------------------------------------------------------------

// debug is an alias for the longer log.V(debugLevel).Info for convenience
func debug(log logr.Logger, obj client.Object, msg string, keysAndValues ...interface{}) {
	// temporarily upgrading all debug to info while developing
	// https://github.com/Kong/kubernetes-ingress-controller/issues/1988
	debugLevel := util.InfoLevel
	keysAndValues = append([]interface{}{
		"namespace", obj.GetNamespace(),
		"name", obj.GetName(),
	}, keysAndValues...)
	log.V(debugLevel).Info(msg, keysAndValues...)
}
