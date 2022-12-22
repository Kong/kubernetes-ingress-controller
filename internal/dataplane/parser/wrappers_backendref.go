package parser

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/types"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

type backendRefWrapper[T types.BackendRefT] struct {
	backendRef T
}

func newBackendRefWrapper[T types.BackendRefT](b T) backendRefWrapper[T] {
	return backendRefWrapper[T]{
		backendRef: b,
	}
}

func (brw backendRefWrapper[T]) Group() *gatewayv1beta1.Group {
	switch br := (interface{})(brw.backendRef).(type) {
	case gatewayv1beta1.BackendRef:
		return br.Group
	}
	return nil
}

func (brw backendRefWrapper[T]) Kind() *gatewayv1beta1.Kind {
	switch br := (interface{})(brw.backendRef).(type) {
	case gatewayv1beta1.BackendRef:
		return br.Kind
	}
	return nil
}

func (brw backendRefWrapper[T]) Name() string {
	switch br := (interface{})(brw.backendRef).(type) {
	case gatewayv1beta1.BackendRef:
		return string(br.Name)
	}
	return ""
}

func (brw backendRefWrapper[T]) Namespace() *string {
	switch br := (interface{})(brw.backendRef).(type) {
	case gatewayv1beta1.BackendRef:
		return (*string)(br.Namespace)
	}
	return nil
}

func (brw backendRefWrapper[T]) Port() int32 {
	switch br := (interface{})(brw.backendRef).(type) {
	case gatewayv1beta1.BackendRef:
		if br.Port == nil {
			return -1
		}
		return int32(*br.Port)
	}
	return -1
}

func (brw backendRefWrapper[T]) Weight() *int32 {
	switch br := (interface{})(brw.backendRef).(type) {
	case gatewayv1beta1.BackendRef:
		return br.Weight
	}
	return nil
}

func backendRefsToKongStateBackends[T types.BackendRefT](
	logger logrus.FieldLogger,
	route client.Object,
	backendRefs []T,
	allowed map[gatewayv1beta1.Namespace][]gatewayv1alpha2.ReferenceGrantTo,
) kongstate.ServiceBackends {
	backends := kongstate.ServiceBackends{}

	for _, backendRef := range backendRefs {
		brw := newBackendRefWrapper(backendRef)

		if util.IsBackendRefGroupKindSupported(
			brw.Group(),
			brw.Kind(),
		) && newRefChecker(backendRef).IsRefAllowedByGrant(allowed) {
			backend := kongstate.ServiceBackend{
				Name: brw.Name(),
				PortDef: kongstate.PortDef{
					Mode:   kongstate.PortModeByNumber,
					Number: brw.Port(),
				},
				Weight: brw.Weight(),
			}
			if brw.Namespace() != nil {
				backend.Namespace = *brw.Namespace()
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
			if brw.Namespace() != nil {
				namespace = *brw.Namespace()
			}
			if brw.Kind() != nil {
				kind = string(*brw.Kind())
			}

			objName := fmt.Sprintf("%s %s/%s",
				route.GetObjectKind().GroupVersionKind().String(),
				route.GetNamespace(),
				route.GetName())
			logger.Errorf(
				"%s requested backendRef to %s %s/%s, but no ReferenceGrant permits it, skipping...",
				objName, kind, namespace, brw.Name())
		}
	}

	return backends
}
