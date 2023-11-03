package parser

import (
	"fmt"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

func backendRefsToKongStateBackends(
	logger logr.Logger,
	route client.Object,
	backendRefs []gatewayapi.BackendRef,
	allowed map[gatewayapi.Namespace][]gatewayapi.ReferenceGrantTo,
) kongstate.ServiceBackends {
	backends := kongstate.ServiceBackends{}

	for _, backendRef := range backendRefs {
		if util.IsBackendRefGroupKindSupported(
			backendRef.Group,
			backendRef.Kind,
		) && newRefChecker(backendRef).IsRefAllowedByGrant(allowed) {
			port := int32(-1)
			if backendRef.Port != nil {
				port = int32(*backendRef.Port)
			}
			backend := kongstate.ServiceBackend{
				Name: string(backendRef.Name),
				PortDef: kongstate.PortDef{
					Mode:   kongstate.PortModeByNumber,
					Number: port,
				},
				Weight: backendRef.Weight,
			}
			if backendRef.Namespace != nil {
				backend.Namespace = string(*backendRef.Namespace)
			}
			backends = append(backends, backend)
		} else {
			// we log impermissible refs rather than failing the entire rule. while we cannot actually route to
			// these, we do not want a single impermissible ref to take the entire rule offline. in the case of edits,
			// failing the entire rule could potentially delete routes that were previously online and in use, and
			// that remain viable (because they still have some permissible backendRefs)
			var (
				namespace = route.GetNamespace()
				kind      string
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
			logger.Error(nil, "Object requested backendRef to target, but no ReferenceGrant permits it, skipping...",
				"object_name", objName,
				"target_kind", kind,
				"target_namespace", namespace,
				"target_name", backendRef.Name,
			)
		}
	}

	return backends
}
