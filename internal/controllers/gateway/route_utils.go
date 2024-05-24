package gateway

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/go-logr/logr"
	"github.com/samber/lo"
	"github.com/samber/mo"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

// -----------------------------------------------------------------------------
// Route Utilities
// -----------------------------------------------------------------------------

const (
	ConditionTypeProgrammed                                            = "Programmed"
	ConditionReasonProgrammedUnknown   gatewayapi.RouteConditionReason = "Unknown"
	ConditionReasonConfiguredInGateway gatewayapi.RouteConditionReason = "ConfiguredInGateway"
	ConditionReasonTranslationError    gatewayapi.RouteConditionReason = "TranslationError"
)

var (
	ErrNoMatchingListenerHostname = fmt.Errorf("no matching hostnames in listener")
	ErrNoSupportedGateway         = fmt.Errorf("no supported gateway found for route")
)

// supportedGatewayWithCondition is a struct that wraps a gateway and some further info
// such as the condition Status condition Accepted of the gateway and the listenerName.
type supportedGatewayWithCondition struct {
	gateway      *gatewayapi.Gateway
	condition    metav1.Condition
	listenerName string
}

func (g supportedGatewayWithCondition) GetName() string {
	return g.gateway.GetName()
}

func (g supportedGatewayWithCondition) GetNamespace() string {
	return g.gateway.GetNamespace()
}

func (g supportedGatewayWithCondition) GetSectionName() mo.Option[string] {
	if g.listenerName != "" {
		return mo.Some(g.listenerName)
	}
	return mo.None[string]()
}

// parentRefsForRoute provides a list of the parentRefs given a Gateway APIs route object
// (e.g. HTTPRoute, TCPRoute, e.t.c.) which refer to the Gateway resource(s) which manage it.
func parentRefsForRoute[T gatewayapi.RouteT](route T) ([]gatewayapi.ParentReference, error) {
	// Note: Ideally we wouldn't have to do this but it's hard to juggle around types
	// and support ParentReference and gatewayapi.ParentReference
	// at the same time so we just copy v1alpha2 refs to a new v1beta1 slice.
	convertV1Alpha2ToV1Beta1ParentReference := func(
		refsAlpha []gatewayapi.ParentReference,
	) []gatewayapi.ParentReference {
		ret := make([]gatewayapi.ParentReference, len(refsAlpha))
		for i, v := range refsAlpha {
			ret[i] = gatewayapi.ParentReference{
				Group:       v.Group,
				Kind:        v.Kind,
				Namespace:   v.Namespace,
				Name:        v.Name,
				SectionName: v.SectionName,
				Port:        v.Port,
			}
		}
		return ret
	}

	var refs []gatewayapi.ParentReference
	switch r := (interface{})(route).(type) {
	case *gatewayapi.HTTPRoute:
		refs = r.Spec.ParentRefs
	case *gatewayapi.UDPRoute:
		refs = r.Spec.ParentRefs
	case *gatewayapi.TCPRoute:
		refs = r.Spec.ParentRefs
	case *gatewayapi.TLSRoute:
		refs = r.Spec.ParentRefs
	case *gatewayapi.GRPCRoute:
		refs = r.Spec.ParentRefs
	default:
		return nil, fmt.Errorf("can't determine parent Gateway for unsupported route type %s", reflect.TypeOf(route))
	}
	for _, ref := range refs {
		if string(*ref.Group) != gatewayv1.GroupName || string(*ref.Kind) != "Gateway" {
			return nil, fmt.Errorf("unsupported parent kind %s/%s", string(*ref.Group), string(*ref.Kind))
		}
	}

	switch r := (interface{})(route).(type) {
	case *gatewayapi.HTTPRoute:
		return r.Spec.ParentRefs, nil
	case *gatewayapi.UDPRoute:
		return convertV1Alpha2ToV1Beta1ParentReference(r.Spec.ParentRefs), nil
	case *gatewayapi.TCPRoute:
		return convertV1Alpha2ToV1Beta1ParentReference(r.Spec.ParentRefs), nil
	case *gatewayapi.TLSRoute:
		return convertV1Alpha2ToV1Beta1ParentReference(r.Spec.ParentRefs), nil
	case *gatewayapi.GRPCRoute:
		return convertV1Alpha2ToV1Beta1ParentReference(r.Spec.ParentRefs), nil
	default:
		return nil, fmt.Errorf("can't determine parent Gateway for unsupported route type %s", reflect.TypeOf(route))
	}
}

// getSupportedGatewayForRoute will retrieve the Gateway and GatewayClass object for any
// Gateway APIs route object (e.g. HTTPRoute, TCPRoute, e.t.c.) from the provided cached
// client if they match this controller. If there are no gateways present for this route
// OR the present gateways are references to missing objects, this will return a unsupportedGW error.
//
// There is a parameter `specifiedGW` here, which is used to specific the gateway.
func getSupportedGatewayForRoute[T gatewayapi.RouteT](ctx context.Context, logger logr.Logger, mgrc client.Client, route T, specifiedGW controllers.OptionalNamespacedName) ([]supportedGatewayWithCondition, error) {
	// gather the parentrefs for this route object
	parentRefs, err := parentRefsForRoute(route)
	if err != nil {
		return nil, err
	}

	// search each parentRef to see if this controller is one of the supported ones
	gateways := make([]supportedGatewayWithCondition, 0)
	for _, parentRef := range parentRefs {
		// gather the namespace/name for the gateway
		namespace := route.GetNamespace()
		if parentRef.Namespace != nil {
			// TODO: need namespace restrictions implementation done before
			// merging this, need to filter out objects with a disallowed NS.
			// https://github.com/Kong/kubernetes-ingress-controller/issues/2080
			namespace = string(*parentRef.Namespace)
		}
		name := string(parentRef.Name)

		// If the flag `--gateway-to-reconcile` is set, KIC will only reconcile the specified gateway.
		// https://github.com/Kong/kubernetes-ingress-controller/issues/5322
		if gatewayToReconcile, ok := specifiedGW.Get(); ok {
			parentNamespace := route.GetNamespace()
			if parentRef.Namespace != nil {
				parentNamespace = string(*parentRef.Namespace)
			}
			if !(parentNamespace == gatewayToReconcile.Namespace && string(parentRef.Name) == gatewayToReconcile.Name) {
				continue
			}
		}

		// pull the Gateway object from the cached client
		gateway := gatewayapi.Gateway{}
		if err := mgrc.Get(ctx, client.ObjectKey{
			Namespace: namespace,
			Name:      name,
		}, &gateway); err != nil {
			if apierrors.IsNotFound(err) {
				// if a configured gateway is not found it's still possible
				// that there's another gateway, so keep searching through the list.
				continue
			}
			return nil, fmt.Errorf("failed to retrieve gateway for route: %w", err)
		}
		gwLogger := logger.WithValues("parentRef.gateway", fmt.Sprintf("%s/%s", gateway.Namespace, gateway.Name))

		// pull the GatewayClass for the Gateway object from the cached client
		gatewayClass := gatewayapi.GatewayClass{}
		if err := mgrc.Get(ctx, client.ObjectKey{
			Name: string(gateway.Spec.GatewayClassName),
		}, &gatewayClass); err != nil {
			if apierrors.IsNotFound(err) {
				// if a configured gatewayClass is not found it's still possible
				// that there's another properly configured gateway in the parentRefs,
				// so keep searching through the list.
				continue
			}
			return nil, fmt.Errorf("failed to retrieve gatewayclass for gateway: %w", err)
		}

		// If the GatewayClass does not match this controller then skip it
		if gatewayClass.Spec.ControllerName != GetControllerName() {
			continue
		}

		// Otherwise we're all set and this controller should reconcile this route.

		var (
			// Set to true if there exists a listener which wasn't filtered by:
			// - AlowedRoutes
			// - listener name matching
			// - listener status checks
			// - listener and route type checks
			matched = false
			// Set to true if ParentRef specified a hostname and it matches route's hostnames.
			matchingHostname *metav1.ConditionStatus
			// Set to true if ParentRef specifies a Port and a listener matches that Port.
			portMatched = false

			allowedByAllowedRoutes  = false
			allowedBySupportedKinds = false
			allowedByListenerName   = false
			listenerReady           = false
		)

		for _, listener := range gateway.Spec.Listeners {
			listenerLogger := gwLogger.WithValues("listener", string(listener.Name))
			// Check if the route matches listener's AllowedRoutes.
			if ok, err := routeMatchesListenerAllowedRoutes(ctx, mgrc, route, listener, gateway.Namespace, parentRef.Namespace); err != nil {
				return nil, fmt.Errorf("failed matching listener %s to a route %s for gateway %s: %w",
					listener.Name, route.GetName(), gateway.Name, err,
				)
			} else if !ok {
				listenerLogger.V(util.DebugLevel).Info("Route does not match listener's allowed routes")
				continue
			}
			allowedByAllowedRoutes = true

			// Check the listeners statuses:
			// - Check if a listener status exists with a matching type (via SupportedKinds).
			// - Check if it matches the requested listener by name (if specified).
			if err := existsMatchingListenerInStatus(route, listener, gateway.Status.Listeners); err != nil {
				listenerLogger.V(util.DebugLevel).Info("Listener does not support this route", "reason", err.Error())
				continue
			} else { //nolint:revive
				allowedBySupportedKinds = true
			}

			if err := listenerProgrammedInStatus(listener.Name, gateway.Status.Listeners); err != nil {
				listenerLogger.V(util.DebugLevel).Info("Listener is not ready", "reason", err.Error())
				continue
			} else { //nolint:revive
				listenerReady = true
			}

			// Check if listener name matches.
			if parentRef.SectionName != nil {
				if *parentRef.SectionName != "" && *parentRef.SectionName != listener.Name {
					listenerLogger.V(util.DebugLevel).Info(
						"Listener name does not match parentRef.SectionName",
						"parentRef_sectionName", parentRef.SectionName,
					)
					continue
				}
				allowedByListenerName = true
			}

			// Perform the port matching as described in GEP-957.
			if parentRef.Port != nil {
				if *parentRef.Port != listener.Port {
					// This ParentRef has a port specified and it's different
					// than current listener's port.
					listenerLogger.V(util.DebugLevel).Info(
						"Listener port does not match parentRef.Port",
						"listener_port", listener.Port, "parentRef_port", parentRef.Port,
					)
					continue
				}
				portMatched = true
			}

			// Check if listener protocol matches
			if !routeTypeMatchesListenerType(route, listener) {
				listenerLogger.V(util.DebugLevel).Info(
					"Route's type does not match listener's type",
					"route_name", route.GetName(),
				)
				continue
			}

			if routeHostnamesIntersectsWithListenerHostname(route, listener) {
				condTrue := metav1.ConditionTrue
				matchingHostname = &condTrue
			} else {
				condFalse := metav1.ConditionFalse
				matchingHostname = &condFalse
				listenerLogger.V(util.DebugLevel).Info("Route's hostname does not match listener's hostname")
				continue
			}

			matched = true
		}

		if matched {
			var listenerName string
			if parentRef.SectionName != nil && *parentRef.SectionName != "" {
				listenerName = string(*parentRef.SectionName)
			}

			gateways = append(gateways, supportedGatewayWithCondition{
				gateway:      &gateway,
				listenerName: listenerName,
				condition: metav1.Condition{
					Type:               string(gatewayapi.RouteConditionAccepted),
					Status:             metav1.ConditionTrue,
					Reason:             string(gatewayapi.RouteReasonAccepted),
					ObservedGeneration: route.GetGeneration(),
				},
			})
		} else {
			// We failed to match a listener with this route

			// This will also catch a case of not matching listener/section name.
			reason := gatewayapi.RouteReasonNoMatchingParent
			switch {
			case matchingHostname != nil && *matchingHostname == metav1.ConditionFalse:
				// If there is no matchingHostname, the gateway Status Condition Accepted
				// must be set to False with reason NoMatchingListenerHostname
				reason = gatewayapi.RouteReasonNoMatchingListenerHostname
			case parentRef.SectionName != nil && !allowedByListenerName:
				// If ParentRef specified listener names but none of the listeners matches the name,
				// the gateway Status Condition Accepted must be set to False with reason RouteReasonNoMatchingParent.
				reason = gatewayapi.RouteReasonNoMatchingParent
			case !listenerReady:
				reason = gatewayapi.RouteReasonNotAllowedByListeners
			case parentRef.Port != nil && !portMatched:
				// If ParentRef specified a Port but none of the listeners matched, the gateway Status
				// Condition Accepted must be set to False with reason NoMatchingListenerPort
				reason = gatewayapi.RouteReasonNoMatchingParent
			case !allowedByAllowedRoutes || !allowedBySupportedKinds:
				reason = gatewayapi.RouteReasonNotAllowedByListeners
			}

			var listenerName string
			if parentRef.SectionName != nil && *parentRef.SectionName != "" {
				listenerName = string(*parentRef.SectionName)
			}

			gateways = append(gateways, supportedGatewayWithCondition{
				gateway:      &gateway,
				listenerName: listenerName,
				condition: metav1.Condition{
					Type:               string(gatewayapi.RouteConditionAccepted),
					Status:             metav1.ConditionFalse,
					Reason:             string(reason),
					ObservedGeneration: route.GetGeneration(),
				},
			})
		}
	}

	if len(gateways) == 0 {
		return nil, ErrNoSupportedGateway
	}

	return gateways, nil
}

func routeHostnamesIntersectsWithListenerHostname[T gatewayapi.RouteT](route T, listener gatewayapi.Listener) bool {
	switch r := (any)(route).(type) {
	case *gatewayapi.HTTPRoute:
		return listenerHostnameIntersectWithRouteHostnames(listener, r.Spec.Hostnames)
	case *gatewayapi.TCPRoute:
		return true
	case *gatewayapi.UDPRoute:
		return true
	case *gatewayapi.TLSRoute:
		return listenerHostnameIntersectWithRouteHostnames(listener, r.Spec.Hostnames)
	case *gatewayapi.GRPCRoute:
		return listenerHostnameIntersectWithRouteHostnames(listener, r.Spec.Hostnames)
	default:
		return false
	}
}

func routeTypeMatchesListenerType[T gatewayapi.RouteT](route T, listener gatewayapi.Listener) bool {
	switch (any)(route).(type) {
	case *gatewayapi.HTTPRoute:
		// HTTPRoutes support Terminate only
		// Note: this is a guess we are doing as the upstream documentation is unclear at the moment.
		// see https://github.com/kubernetes-sigs/gateway-api/issues/1474
		if !(listener.Protocol == gatewayapi.HTTPProtocolType || listener.Protocol == gatewayapi.HTTPSProtocolType) {
			return false
		}
		if listener.TLS != nil && *listener.TLS.Mode != gatewayapi.TLSModeTerminate {
			return false
		}
	case *gatewayapi.TCPRoute:
		if listener.Protocol != gatewayapi.TCPProtocolType {
			return false
		}
		// TCPRoutes support Terminate only
		// Note: this is a guess we are doing as the upstream documentation is unclear at the moment.
		// see https://github.com/kubernetes-sigs/gateway-api/issues/1474
		if listener.TLS != nil && *listener.TLS.Mode != gatewayapi.TLSModeTerminate {
			return false
		}
	case *gatewayapi.UDPRoute:
		if listener.Protocol != gatewayapi.UDPProtocolType {
			return false
		}
		// TLS should not be set in UDP listeners
		if listener.TLS != nil {
			return false
		}
	case *gatewayapi.TLSRoute:
		if listener.Protocol != gatewayapi.TLSProtocolType {
			return false
		}
		// TLSRoutes currently support Passthrough only
		// Note: this is a guess we are doing as the upstream documentation is unclear at the moment.
		// see https://github.com/kubernetes-sigs/gateway-api/issues/1474
		if listener.TLS != nil && *listener.TLS.Mode != gatewayapi.TLSModePassthrough {
			return false
		}
	case *gatewayapi.GRPCRoute:
		if listener.Protocol != gatewayapi.HTTPSProtocolType && listener.Protocol != gatewayapi.HTTPProtocolType {
			return false
		}
	default:
		return false
	}
	return true
}

// routeMatchesListenerAllowedRoutes checks if the provided route matches the
// criteria defined in listener's AllowedRoutes field.
func routeMatchesListenerAllowedRoutes[T gatewayapi.RouteT](
	ctx context.Context,
	mgrc client.Client,
	route T,
	listener gatewayapi.Listener,
	gatewayNamespace string,
	parentRefNamespace *gatewayapi.Namespace,
) (bool, error) {
	if listener.AllowedRoutes == nil {
		return true, nil
	}

	if len(listener.AllowedRoutes.Kinds) > 0 {
		// Find if the route has a type that's within the listener's supported gatewayapi.
		_, ok := lo.Find(listener.AllowedRoutes.Kinds, func(rgk gatewayapi.RouteGroupKind) bool {
			gvk := route.GetObjectKind().GroupVersionKind()
			return (rgk.Group != nil && string(*rgk.Group) == gvk.Group) && string(rgk.Kind) == gvk.Kind
		})
		if !ok {
			return false, nil
		}
	}

	if listener.AllowedRoutes.Namespaces == nil || listener.AllowedRoutes.Namespaces.From == nil {
		return true, nil
	}

	switch *listener.AllowedRoutes.Namespaces.From {
	case gatewayapi.NamespacesFromAll:
		return true, nil

	case gatewayapi.NamespacesFromSame:
		// If parentRef didn't specify the namespace then we check if
		// the gateway is from the same namespace as the route
		if parentRefNamespace == nil {
			return gatewayNamespace == route.GetNamespace(), nil
		}
		// Otherwise compare routes namespace with parentRef's one.
		return route.GetNamespace() == string(*parentRefNamespace), nil

	case gatewayapi.NamespacesFromSelector:
		namespace := corev1.Namespace{}
		if err := mgrc.Get(ctx, client.ObjectKey{Name: route.GetNamespace()}, &namespace); err != nil {
			return false, fmt.Errorf("failed to get namespace %s: %w", route.GetNamespace(), err)
		}

		s, err := metav1.LabelSelectorAsSelector(listener.AllowedRoutes.Namespaces.Selector)
		if err != nil {
			return false, fmt.Errorf(
				"failed to convert AllowedRoutes LabelSelector %s to Selector for listener %s: %w",
				listener.AllowedRoutes.Namespaces.Selector, listener.Name, err,
			)
		}

		ok := s.Matches(labels.Set(namespace.Labels))
		return ok, nil

	default:
		return false, fmt.Errorf(
			"unknown listener.AllowedRoutes.Namespaces.From value: %s for listener %s",
			*listener.AllowedRoutes.Namespaces.From, listener.Name,
		)
	}
}

var (
	errUnsupportedRouteKind          = errors.New("unsupported route kind")
	errUnmatchedListenerName         = errors.New("unmatched listener name")
	errListenerNoProgrammedCondition = errors.New("no Programmed condition found for listener")
	errListenerNotProgrammedYet      = errors.New("listener not programmed yet")
)

// existsMatchingReadyListenerInStatus checks if:
// - If a listener status exists with a matching type (via SupportedKinds).
// - If it matches the requested listener by name (if specified).
// - And finally check if the provided listener is marked as Ready.
func existsMatchingListenerInStatus[T gatewayapi.RouteT](route T, listener gatewayapi.Listener, lss []gatewayapi.ListenerStatus) error {
	listenerFound := false

	// Find listener's status...
	_, ok := lo.Find(lss, func(ls gatewayapi.ListenerStatus) bool {
		if ls.Name != listener.Name {
			return false
		}
		listenerFound = true

		// Find if the route has a type that's within the supported types, listed
		// in listener's status.
		_, ok := lo.Find(ls.SupportedKinds, func(rgk gatewayapi.RouteGroupKind) bool {
			// The artificially filled in GVK is needed for testing mostly and for
			// situations when the object is not coming from the api server.
			// Related upstream issue: https://github.com/kubernetes/kubernetes/issues/3030
			var gvk schema.GroupVersionKind
			switch any(route).(type) {
			case *gatewayapi.HTTPRoute:
				gvk = schema.GroupVersionKind{
					Group:   gatewayv1.GroupVersion.Group,
					Version: gatewayv1.GroupVersion.Version,
					Kind:    "HTTPRoute",
				}
			default:
				gvk = route.GetObjectKind().GroupVersionKind()
			}
			return (rgk.Group != nil && string(*rgk.Group) == gvk.Group) && string(rgk.Kind) == gvk.Kind
		})
		return ok
	})

	if !ok && !listenerFound {
		return errUnmatchedListenerName // Matching Listener's not found.
	}

	if !ok && listenerFound {
		return errUnsupportedRouteKind // Listener(s) found but none with matching supported kinds.
	}

	return nil
}

func listenerProgrammedInStatus(listenerName gatewayapi.SectionName, lss []gatewayapi.ListenerStatus) error {
	listenerStatus, ok := lo.Find(lss, func(ls gatewayapi.ListenerStatus) bool {
		return ls.Name == listenerName
	})
	if !ok {
		return errUnmatchedListenerName // Matching Listener's not found.
	}

	programmedStatus, ok := lo.Find(listenerStatus.Conditions, func(condition metav1.Condition) bool {
		return condition.Type == string(gatewayapi.ListenerConditionProgrammed)
	})
	if !ok {
		return errListenerNoProgrammedCondition // "Programmed" condition not found in conditions of listener's conditions.
	}

	if programmedStatus.Status != metav1.ConditionTrue {
		return errListenerNotProgrammedYet // "Programmed" condition is not true.
	}

	return nil
}

func listenerHostnameIntersectWithRouteHostnames[H gatewayapi.HostnameT](listener gatewayapi.Listener, hostnames []H) bool {
	if len(hostnames) == 0 {
		return true
	}

	// if the listener has no hostname, all hostnames automatically intersect
	if listener.Hostname == nil || *listener.Hostname == "" {
		return true
	}

	// iterate over all the hostnames and check that at least one intersect with the listener hostname
	for _, hostname := range hostnames {
		if util.HostnamesIntersect(*listener.Hostname, hostname) {
			return true
		}
	}

	return false
}

// isListenerHostnameEffective returns true if the listener can specify an effective
// hostname to match hostnames in requests.
// It basically checks if the listener is using any these protocols: HTTP, HTTPS, or TLS.
func isListenerHostnameEffective(listener gatewayapi.Listener) bool {
	return listener.Protocol == gatewayapi.HTTPProtocolType ||
		listener.Protocol == gatewayapi.HTTPSProtocolType ||
		listener.Protocol == gatewayapi.TLSProtocolType
}

// filterHostnames accepts a HTTPRoute and returns a version of the same object with only a subset of the
// hostnames, the ones matching with the listeners' hostname.
// it returns an error if the intersection of hostname match in httproute and listeners is empty.
func filterHostnames(gateways []supportedGatewayWithCondition, httpRoute *gatewayapi.HTTPRoute) (*gatewayapi.HTTPRoute, error) {
	filteredHostnames := make([]gatewayapi.Hostname, 0)
	// if no hostnames are specified in the route spec, we use the UNION of all hostnames in supported gateways.
	// if any of supported listener has not specified hostname, the hostnames of HTTPRoute remains empty
	// to match **ANY** hostname.
	if len(httpRoute.Spec.Hostnames) == 0 {
		var matchAnyHost bool
		filteredHostnames, matchAnyHost = getUnionOfGatewayHostnames(gateways)
		if matchAnyHost {
			return httpRoute, nil
		}
	} else {
		for _, hostname := range httpRoute.Spec.Hostnames {
			if hostnameMatching := getMinimumHostnameIntersection(gateways, hostname); hostnameMatching != "" {
				filteredHostnames = append(filteredHostnames, hostnameMatching)
			}
		}
		if len(filteredHostnames) == 0 {
			return httpRoute, ErrNoMatchingListenerHostname
		}
	}

	httpRoute.Spec.Hostnames = filteredHostnames
	return httpRoute, nil
}

// getUnionOfGatewayHostnames returns UNION of hostnames specified in supported gateways.
// the second return value is true if any hostname could be matched in at least one listener
// in supported gateways and listeners, so the `HTTPRoute` could match any hostname.
func getUnionOfGatewayHostnames(gateways []supportedGatewayWithCondition) ([]gatewayapi.Hostname, bool) {
	hostnames := make([]gatewayapi.Hostname, 0)
	for _, gateway := range gateways {
		if gateway.listenerName != "" {
			if listener := extractListenerSpecFromGateway(
				gateway.gateway,
				gatewayapi.SectionName(gateway.listenerName),
			); listener != nil {
				// return true if the listener has not specified hostname to match any hostname.
				if listener.Hostname == nil {
					return nil, true
				}
				hostnames = append(hostnames, *listener.Hostname)
			}
		} else {
			for _, listener := range gateway.gateway.Spec.Listeners {
				// here we consider ALL listeners that are able to configure a hostname if no listener attached.
				// may be changed if there is a conclusion on the upstream discussion about it:
				// https://github.com/kubernetes-sigs/gateway-api/discussions/1563
				if isListenerHostnameEffective(listener) {
					if listener.Hostname == nil {
						return nil, true
					}
					hostnames = append(hostnames, *listener.Hostname)
				}
			}
		}
	}
	return hostnames, false
}

// getMinimumHostnameIntersection returns the minimum intersecting hostname, in the sense that:
//
// - if the listener hostname is empty, return the httpRoute hostname
// - if the listener hostname acts as a wildcard for the httpRoute hostname, return the httpRoute hostname
// - if the httpRoute hostname acts as a wildcard for the listener hostname, return the listener hostname
// - if the httpRoute hostname is the same of the listener hostname, return it
// - if none of the above is true, return an empty string.
func getMinimumHostnameIntersection(gateways []supportedGatewayWithCondition, hostname gatewayapi.Hostname) gatewayapi.Hostname {
	for _, gateway := range gateways {
		for _, listener := range gateway.gateway.Spec.Listeners {
			// if the listenerName is specified and matches the name of the gateway listener proceed
			if (gatewayapi.SectionName)(gateway.listenerName) == "" ||
				(gatewayapi.SectionName)(gateway.listenerName) == (listener.Name) {
				if listener.Hostname == nil || *listener.Hostname == "" {
					return hostname
				}
				if util.HostnamesMatch(string(*listener.Hostname), string(hostname)) {
					return hostname
				}
				if util.HostnamesMatch(string(hostname), string(*listener.Hostname)) {
					return (*listener.Hostname)
				}
			}
		}
	}
	return ""
}

func isRouteAccepted(gateways []supportedGatewayWithCondition) bool {
	for _, gateway := range gateways {
		if gateway.condition.Type == string(gatewayapi.RouteConditionAccepted) && gateway.condition.Status == metav1.ConditionTrue {
			return true
		}
	}
	return false
}

// isHTTPReferenceGranted checks that the backendRef referenced by the HTTPRoute is granted by a ReferenceGrant.
func isHTTPReferenceGranted(grantSpec gatewayapi.ReferenceGrantSpec, backendRef gatewayapi.HTTPBackendRef, fromNamespace string) bool {
	var backendRefGroup gatewayapi.Group
	var backendRefKind gatewayapi.Kind

	if backendRef.Group != nil {
		backendRefGroup = *backendRef.Group
	}
	if backendRef.Kind != nil {
		backendRefKind = *backendRef.Kind
	}
	for _, from := range grantSpec.From {
		if from.Group != gatewayv1.GroupName || from.Kind != "HTTPRoute" || fromNamespace != string(from.Namespace) {
			continue
		}

		for _, to := range grantSpec.To {
			if backendRefGroup == to.Group &&
				backendRefKind == to.Kind &&
				(to.Name == nil || *to.Name == backendRef.Name) {
				return true
			}
		}
	}
	return false
}

// sameCondition returns true if the conditions in parameter has the same type, status, reason and message.
func sameCondition(a, b metav1.Condition) bool {
	return a.Type == b.Type &&
		a.Status == b.Status &&
		a.Reason == b.Reason &&
		a.Message == b.Message &&
		a.ObservedGeneration == b.ObservedGeneration
}

func setRouteParentStatusCondition(parentStatus *gatewayapi.RouteParentStatus, newCondition metav1.Condition) bool {
	var conditionFound, changed bool
	for i, condition := range parentStatus.Conditions {
		if condition.Type == newCondition.Type {
			conditionFound = true
			if !sameCondition(condition, newCondition) {
				parentStatus.Conditions[i] = newCondition
				changed = true
			}
		}
	}

	if !conditionFound {
		parentStatus.Conditions = append(parentStatus.Conditions, newCondition)
		changed = true
	}
	return changed
}

func parentStatusHasProgrammedCondition(parentStatus *gatewayapi.RouteParentStatus) bool {
	for _, condition := range parentStatus.Conditions {
		if condition.Type == ConditionTypeProgrammed {
			return true
		}
	}
	return false
}

// ensureParentsProgrammedCondition ensures that provided route's parent statuses
// have Programmed condition set properly. It returns a boolean flag indicating
// whether an update to the provided route has been performed.
//
// Use the condition argument to specify the Reason, Status and Message.
// Type will be set to Programmed whereas ObservedGeneration and LastTransitionTime
// will be set accordingly based on the route's generation and current time.
func ensureParentsProgrammedCondition[
	routeT gatewayapi.RouteT,
](
	ctx context.Context,
	client client.SubResourceWriter,
	route routeT,
	routeParentStatuses []gatewayapi.RouteParentStatus,
	gateways []supportedGatewayWithCondition,
	condition metav1.Condition,
) (bool, error) {
	// map the existing parentStatues to avoid duplications
	parentStatuses := getParentStatuses(route, routeParentStatuses)

	condition.Type = ConditionTypeProgrammed
	condition.ObservedGeneration = route.GetGeneration()
	condition.LastTransitionTime = metav1.Now()

	statusChanged := false
	for _, g := range gateways {
		gateway := g.gateway

		parentRefKey := routeParentStatusKey(route, g)
		parentStatus, ok := parentStatuses[parentRefKey]
		if ok {
			// update existing parent in status.
			changed := setRouteParentStatusCondition(parentStatus, condition)
			if changed {
				parentStatuses[parentRefKey] = parentStatus
				setRouteParentInStatusForParent(route, *parentStatus, g)
			}
			statusChanged = statusChanged || changed
		} else {
			// add a new parent if the parent is not found in status.
			newParentStatus := gatewayapi.RouteParentStatus{
				ParentRef: gatewayapi.ParentReference{
					Namespace: lo.ToPtr(gatewayapi.Namespace(gateway.Namespace)),
					Name:      gatewayapi.ObjectName(gateway.Name),
					Kind:      lo.ToPtr(gatewayapi.Kind("Gateway")),
					Group:     lo.ToPtr(gatewayapi.Group(gatewayv1.GroupName)),
					SectionName: func() *gatewayapi.SectionName {
						// We don't need to check whether the listener matches route's spec
						// because that should already be done via getSupportedGatewayForRoute
						// at this point.
						if g.listenerName != "" {
							return lo.ToPtr(gatewayapi.SectionName(g.listenerName))
						}
						return nil
					}(),

					// TODO: set port after gateway port matching implemented:
					// https://github.com/Kong/kubernetes-ingress-controller/issues/3016
				},
				ControllerName: GetControllerName(),
				Conditions: []metav1.Condition{
					condition,
				},
			}
			setRouteParentInStatusForParent(route, newParentStatus, g)

			routeParentStatuses = append(routeParentStatuses, newParentStatus)
			parentStatuses[parentRefKey] = &newParentStatus
			statusChanged = true
		}
	}

	// update status if needed.
	if statusChanged {
		if err := client.Update(ctx, route); err != nil {
			return false, err
		}
		return true, nil
	}
	// no need to update if no status is changed.
	return false, nil
}

// setRouteParentInStatusForParent checks if the provided route Status, contains
// status for the provided parent and if it does it sets it to the provided
// RouteStatusParent. If it does not then it appends the provided RouteStatusParent
// to provided route's Status.Parents field.
//
// This might come in useful when the caller wants to set only one parent's
// status.
func setRouteParentInStatusForParent[
	routeT gatewayapi.RouteT,
	parentT namespacedNamer,
](
	route routeT,
	routeStatusParent gatewayapi.RouteParentStatus,
	parent parentT,
) {
	switch r := any(route).(type) {
	case *gatewayapi.HTTPRoute:
		r.Status.Parents = ensureRoutesParents(r.Status.Parents, routeStatusParent, parent)
	case *gatewayapi.TCPRoute:
		r.Status.Parents = ensureRoutesParents(r.Status.Parents, routeStatusParent, parent)
	case *gatewayapi.UDPRoute:
		r.Status.Parents = ensureRoutesParents(r.Status.Parents, routeStatusParent, parent)
	case *gatewayapi.TLSRoute:
		r.Status.Parents = ensureRoutesParents(r.Status.Parents, routeStatusParent, parent)
	case *gatewayapi.GRPCRoute:
		r.Status.Parents = ensureRoutesParents(r.Status.Parents, routeStatusParent, parent)
	}
}

// ensureRoutesParents ensures that the provided RouteStatusParents are updated.
// This function checks if the provided []RouteStatusParents contains a parentRef
// status for the provided status.
// If it doesn't then it adds it to the provided []RouteStatusParents and returns it.
// If it does then it overwrites the provious status for that parent and returns
// the updates []RouteStatusParents.
func ensureRoutesParents[
	parentT namespacedNamer,
](
	routeStatusParents []gatewayapi.RouteParentStatus,
	routeStatusParent gatewayapi.RouteParentStatus,
	parent parentT,
) []gatewayapi.RouteParentStatus {
	for i, p := range routeStatusParents {
		if ensureParentsStatusUpdated(p.ParentRef, parent, routeStatusParents, i, routeStatusParent) {
			return routeStatusParents
		}
	}
	routeStatusParents = append(routeStatusParents, routeStatusParent)
	return routeStatusParents
}

func ensureParentsStatusUpdated[
	parentT namespacedNamer,
](
	parentRef gatewayapi.ParentReference,
	parent parentT,
	routeStatusParents []gatewayapi.RouteParentStatus,
	i int,
	routeStatusParent gatewayapi.RouteParentStatus,
) bool {
	if !isParentRefEqualToParent(parentRef, parent) {
		return false
	}

	routeStatusParents[i] = routeStatusParent
	return true
}

func isParentRefEqualToParent[
	parentT namespacedNamer,
](
	parentRef gatewayapi.ParentReference,
	parent parentT,
) bool {
	if *parentRef.Group != gatewayv1.GroupName {
		return false
	}
	if *parentRef.Kind != "Gateway" {
		return false
	}
	if string(parentRef.Name) != parent.GetName() {
		return false
	}
	if parentRef.Namespace != nil && string(*parentRef.Namespace) != parent.GetNamespace() {
		return false
	}
	if parentRef.SectionName != nil && string(*parentRef.SectionName) != parent.GetSectionName().OrEmpty() {
		return false
	}

	return true
}

// isRouteAcceptedByListener checks the given route is accepted by the
// gateway's listener specified by a proper parentReference.
func isRouteAcceptedByListener[T gatewayapi.RouteT](ctx context.Context,
	mgrc client.Client,
	route T,
	gateway gatewayapi.Gateway,
	listenerIndex int,
	parentRef gatewayapi.ParentReference,
) (bool, error) {
	// Check if the route matches listener's AllowedRoutes.
	listener := gateway.Spec.Listeners[listenerIndex]
	if ok, err := routeMatchesListenerAllowedRoutes(ctx, mgrc, route, listener, gateway.Namespace, parentRef.Namespace); err != nil {
		return false, fmt.Errorf("failed matching listener %s to a route %s for gateway %s: %w",
			listener.Name, route.GetName(), gateway.Name, err,
		)
	} else if !ok {
		return false, nil
	}

	// Check the listeners statuses:
	// - Check if a listener status exists with a matching type (via SupportedKinds).
	// - Check if it matches the requested listener by name (if specified).
	if err := existsMatchingListenerInStatus(route, listener, gateway.Status.Listeners); err != nil {
		// return no error here, as we don't care of the reason why this check failed.
		return false, nil //nolint:nilerr
	}

	// Check if listener name matches.
	if parentRef.SectionName != nil && *parentRef.SectionName != "" && *parentRef.SectionName != listener.Name {
		return false, nil
	}

	// Perform the port matching as described in GEP-957.
	if parentRef.Port != nil && *parentRef.Port != listener.Port {
		// This ParentRef has a port specified and it's different
		// than current listener's port.
		return false, nil
	}

	if !routeTypeMatchesListenerType(route, listener) {
		return false, nil
	}

	if !routeHostnamesIntersectsWithListenerHostname(route, listener) {
		return false, nil
	}

	return true, nil
}

// ensureGatewayReferenceStatusRemoved uses the ControllerName provided by the Gateway
// implementation to prune status references to Gateways supported by this controller
// in the provided route.
func ensureGatewayReferenceStatusRemoved[routeT gatewayapi.RouteT](
	ctx context.Context, cl client.Client, log logr.Logger, route routeT,
) (bool, error) {
	debug(log, route, "Unsupported route found, processing to verify whether it was ever supported")
	kind := route.GetObjectKind().GroupVersionKind().Kind
	parents := getRouteStatusParents(route)

	// drop all status references to supported Gateway objects
	newStatuses := make([]gatewayapi.RouteParentStatus, 0)
	for _, status := range parents {
		if status.ControllerName != GetControllerName() {
			newStatuses = append(newStatuses, status)
		} else {
			parentRefNN := string(status.ParentRef.Name)
			if status.ParentRef.Namespace != nil {
				parentRefNN = fmt.Sprintf("%s/%s", *status.ParentRef.Namespace, parentRefNN)
			}
			debug(log, route, "Removing parentRef from route status", "parentRef", parentRefNN, "kind", kind)
		}
	}

	// if the new list of statuses is the same length as the old
	// nothing has changed and we're all done.
	if len(newStatuses) == len(parents) {
		return false, nil
	}

	// if the route doesn't have a supported Gateway+GatewayClass associated with
	// it it's possible it became orphaned after becoming queued. In either case
	// ensure that it's removed from the proxy cache to avoid orphaned data-plane
	// configurations.
	debug(log, route, "Ensuring that dataplane is updated to remove unsupported route (if applicable)")
	setRouteStatusParents(route, newStatuses)
	if err := cl.Status().Update(ctx, route); err != nil {
		return false, fmt.Errorf("failed to remove Gateway parentRef from %s status: %w", kind, err)
	}

	debug(log, route, "Unsupported route was previously supported, status was updated")
	// the status needed to be updated and it was updated successfully
	return true, nil
}

func getRouteParentRefs[T gatewayapi.RouteT](route T) []gatewayapi.ParentReference {
	switch r := any(route).(type) {
	case *gatewayapi.HTTPRoute:
		return r.Spec.ParentRefs
	case *gatewayapi.TCPRoute:
		return r.Spec.ParentRefs
	case *gatewayapi.UDPRoute:
		return r.Spec.ParentRefs
	case *gatewayapi.TLSRoute:
		return r.Spec.ParentRefs
	case *gatewayapi.GRPCRoute:
		return r.Spec.ParentRefs
	default:
		return nil
	}
}
