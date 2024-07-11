package translator

import (
	"errors"
	"fmt"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

// backendRefsToKongStateBackends takes a list of BackendRefs and returns a list of ServiceBackends.
// The backendRefs are checked for the following conditions. If any of these conditions are met, the BackendRef is
// not included in the returned list:
// - If a BackendRef is not permitted by the provided ReferenceGrantTo set,
// - If a BackendRef is not found,
// - If a BackendRef Group & Kind pair is not supported (currently only Service is supported),
// - If a BackendRef is missing a port.
// The provided client is used to retrieve the Backend referenced by the BackendRef
// to check if it exists.
func backendRefsToKongStateBackends(
	logger logr.Logger,
	storer store.Storer,
	route client.Object,
	backendRefs []gatewayapi.BackendRef,
	allowed map[gatewayapi.Namespace][]gatewayapi.ReferenceGrantTo,
) kongstate.ServiceBackends {
	backends := kongstate.ServiceBackends{}

	for _, backendRef := range backendRefs {
		logger := loggerForBackendRef(logger, route, backendRef)

		nn := client.ObjectKey{
			Name:      string(backendRef.Name),
			Namespace: route.GetNamespace(),
		}
		if backendRef.Namespace != nil {
			nn.Namespace = string(*backendRef.Namespace)
		}

		if backendRef.Kind == nil {
			// This should never happen as the default value defined by Gateway API is 'Service'. Checking just in case.
			logger.Error(nil, "Object requested backendRef to target, but no Kind was specified, skipping...")
			continue
		}

		var err error
		switch *backendRef.Kind {
		case "Service":
			_, err = storer.GetService(nn.Namespace, nn.Name)
		default:
			err = fmt.Errorf("unsupported kind %q, only 'Service' is supported", *backendRef.Kind)
		}
		if err != nil {
			if errors.As(err, &store.NotFoundError{}) {
				logger.Error(err, "Object requested backendRef to target, but it does not exist, skipping...")
			} else {
				logger.Error(err, "Object requested backendRef to target, but an error occurred, skipping...")
			}
			continue
		}

		if !util.IsBackendRefGroupKindSupported(backendRef.Group, backendRef.Kind) ||
			!gatewayapi.NewRefCheckerForRoute(logger, route, backendRef).IsRefAllowedByGrant(allowed) {
			// we log impermissible refs rather than failing the entire rule. while we cannot actually route to
			// these, we do not want a single impermissible ref to take the entire rule offline. in the case of edits,
			// failing the entire rule could potentially delete routes that were previously online and in use, and
			// that remain viable (because they still have some permissible backendRefs)
			logger.Error(nil, "Object requested backendRef to target, but no ReferenceGrant permits it, skipping...")
			continue
		}

		port := int32(-1)
		if backendRef.Port != nil {
			port = int32(*backendRef.Port)
		}
		backend, err := kongstate.NewServiceBackendForService(
			nn,
			kongstate.PortDef{
				Mode:   kongstate.PortModeByNumber,
				Number: port,
			},
		)
		if err != nil {
			logger.Error(err, "failed to create ServiceBackend for backendRef")
			continue
		}
		if backendRef.Weight != nil {
			backend.SetWeight(*backendRef.Weight)
		}
		backends = append(backends, backend)
	}

	return backends
}

func loggerForBackendRef(logger logr.Logger, route client.Object, backendRef gatewayapi.BackendRef) logr.Logger {
	var (
		namespace = route.GetNamespace()
		kind      = "unknown"
	)
	if backendRef.Namespace != nil {
		namespace = string(*backendRef.Namespace)
	}
	if backendRef.Kind != nil {
		kind = string(*backendRef.Kind)
	}

	objName := fmt.Sprintf("%s %s/%s",
		route.GetObjectKind().GroupVersionKind().String(),
		route.GetNamespace(),
		route.GetName())
	return logger.WithValues(
		"object_name", objName,
		"target_kind", kind,
		"target_namespace", namespace,
		"target_name", backendRef.Name,
	)
}
