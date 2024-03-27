package parser

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

func backendRefsToKongStateBackends(
	logger logrus.FieldLogger,
	route client.Object,
	backendRefs []gatewayv1.BackendRef,
	allowed map[gatewayv1.Namespace][]gatewayv1beta1.ReferenceGrantTo,
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
			logger.Errorf(
				"%s requested backendRef to %s %s/%s, but no ReferenceGrant permits it, skipping...",
				objName, kind, namespace, backendRef.Name)
		}
	}

	return backends
}
