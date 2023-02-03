package gateway

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/types"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// -----------------------------------------------------------------------------
// Route Utilities
// -----------------------------------------------------------------------------

const (
	unsupportedGW = "no supported Gateway found for route"
)

const (
	ConditionTypeProgrammed                                                = "Programmed"
	ConditionReasonProgrammedUnknown   gatewayv1beta1.RouteConditionReason = "Unknown"
	ConditionReasonConfiguredInGateway gatewayv1beta1.RouteConditionReason = "ConfiguredInGateway"
	ConditionReasonTranslationError    gatewayv1beta1.RouteConditionReason = "TranslationError"
)

var ErrNoMatchingListenerHostname = fmt.Errorf("no matching hostnames in listener")

// supportedGatewayWithCondition is a struct that wraps a gateway and some further info
// such as the condition Status condition Accepted of the gateway and the listenerName.
type supportedGatewayWithCondition struct {
	gateway      *Gateway
	condition    metav1.Condition
	listenerName string
}

// parentRefsForRoute provides a list of the parentRefs given a Gateway APIs route object
// (e.g. HTTPRoute, TCPRoute, e.t.c.) which refer to the Gateway resource(s) which manage it.
func parentRefsForRoute[T types.RouteT](route T) ([]ParentReference, error) {
	// Note: Ideally we wouldn't have to do this but it's hard to juggle around types
	// and support ParentReference and gatewayv1alpha2.ParentReference
	// at the same time so we just copy v1alpha2 refs to a new v1beta1 slice.
	convertV1Alpha2ToV1Beta1ParentReference := func(
		refsAlpha []gatewayv1alpha2.ParentReference,
	) []ParentReference {
		ret := make([]ParentReference, len(refsAlpha))
		for i, v := range refsAlpha {
			ret[i] = ParentReference{
				Group:       (*Group)(v.Group),
				Kind:        (*Kind)(v.Kind),
				Namespace:   (*Namespace)(v.Namespace),
				Name:        (ObjectName)(v.Name),
				SectionName: (*SectionName)(v.SectionName),
				Port:        (*PortNumber)(v.Port),
			}
		}
		return ret
	}

	switch r := (interface{})(route).(type) {
	case *gatewayv1beta1.HTTPRoute:
		return r.Spec.ParentRefs, nil
	case *gatewayv1alpha2.UDPRoute:
		return convertV1Alpha2ToV1Beta1ParentReference(r.Spec.ParentRefs), nil
	case *gatewayv1alpha2.TCPRoute:
		return convertV1Alpha2ToV1Beta1ParentReference(r.Spec.ParentRefs), nil
	case *gatewayv1alpha2.TLSRoute:
		return convertV1Alpha2ToV1Beta1ParentReference(r.Spec.ParentRefs), nil
	default:
		return nil, fmt.Errorf("cant determine parent gateway for unsupported route type %s", reflect.TypeOf(route))
	}
}

const (
	// This reason is used with the "Accepted" condition when there are
	// no matching Parents. In the case of Gateways, this can occur when
	// a Route ParentRef specifies a Port and/or SectionName that does not
	// match any Listeners in the Gateway.
	//
	// NOTE: This is already in uptsream, albeit unreleased:
	// https://github.com/kubernetes-sigs/gateway-api/pull/1516
	// TODO: swap this out with upstream const when released.
	RouteReasonNoMatchingParent gatewayv1beta1.RouteConditionReason = "NoMatchingParent"
)

// getSupportedGatewayForRoute will retrieve the Gateway and GatewayClass object for any
// Gateway APIs route object (e.g. HTTPRoute, TCPRoute, e.t.c.) from the provided cached
// client if they match this controller. If there are no gateways present for this route
// OR the present gateways are references to missing objects, this will return a unsupportedGW error.
func getSupportedGatewayForRoute[T types.RouteT](ctx context.Context, mgrc client.Client, route T) ([]supportedGatewayWithCondition, error) {
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

		// pull the Gateway object from the cached client
		gateway := gatewayv1beta1.Gateway{}
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

		// pull the GatewayClass for the Gateway object from the cached client
		gatewayClass := gatewayv1beta1.GatewayClass{}
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
		)

		for _, listener := range gateway.Spec.Listeners {
			// Check if the route matches listener's AllowedRoutes.
			if ok, err := routeMatchesListenerAllowedRoutes(ctx, mgrc, route, listener, gateway.Namespace, parentRef.Namespace); err != nil {
				return nil, fmt.Errorf("failed matching listener %s to a route %s for gateway %s: %w",
					listener.Name, route.GetName(), gateway.Name, err,
				)
			} else if !ok {
				continue
			} else {
				allowedByAllowedRoutes = true
			}

			// Check the listeners statuses:
			// - Check if a listener status exists with a matching type (via SupportedKinds).
			// - Check if it matches the requested listener by name (if specified).
			// - And finally check if that listeners is marked as Ready.
			if err := existsMatchingReadyListenerInStatus(route, listener, gateway.Status.Listeners); err != nil {
				continue
			} else {
				allowedBySupportedKinds = true
			}

			// Check if listener name matches.
			if parentRef.SectionName != nil {
				if *parentRef.SectionName != "" && *parentRef.SectionName != listener.Name {
					continue
				}
				allowedByListenerName = true
			}

			// Perform the port matching as described in GEP-957.
			if parentRef.Port != nil {
				if *parentRef.Port != listener.Port {
					// This ParentRef has a port specified and it's different
					// than current listener's port.
					continue
				}
				portMatched = true
			}

			if !routeTypeMatchesListenerType(route, listener) {
				continue
			}

			if routeHostnamesIntersectsWithListenerHostname(route, listener) {
				condTrue := metav1.ConditionTrue
				matchingHostname = &condTrue
			} else {
				condFalse := metav1.ConditionFalse
				matchingHostname = &condFalse
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
					Type:   string(gatewayv1beta1.RouteConditionAccepted),
					Status: metav1.ConditionTrue,
					Reason: string(gatewayv1beta1.RouteReasonAccepted),
				},
			})
		} else {
			// We failed to match a listener with this route

			// This will also catch a case of not matching listener/section name.
			reason := RouteReasonNoMatchingParent

			if matchingHostname != nil && *matchingHostname == metav1.ConditionFalse {
				// If there is no matchingHostname, the gateway Status Condition Accepted
				// must be set to False with reason NoMatchingListenerHostname
				reason = gatewayv1beta1.RouteReasonNoMatchingListenerHostname
			} else if (parentRef.SectionName) != nil && !allowedByListenerName {
				// If ParentRef specified listener names but none of the listeners matches the name,
				// the gateway Status Condition Accepted must be set to False with reason RouteReasonNoMatchingParent.
				reason = RouteReasonNoMatchingParent
			} else if (parentRef.Port != nil) && !portMatched {
				// If ParentRef specified a Port but none of the listeners matched, the gateway Status
				// Condition Accepted must be set to False with reason NoMatchingListenerPort
				reason = RouteReasonNoMatchingParent
			} else if !allowedByAllowedRoutes || !allowedBySupportedKinds {
				reason = gatewayv1beta1.RouteReasonNotAllowedByListeners
			}

			var listenerName string
			if parentRef.SectionName != nil && *parentRef.SectionName != "" {
				listenerName = string(*parentRef.SectionName)
			}

			gateways = append(gateways, supportedGatewayWithCondition{
				gateway:      &gateway,
				listenerName: listenerName,
				condition: metav1.Condition{
					Type:   string(gatewayv1beta1.RouteConditionAccepted),
					Status: metav1.ConditionFalse,
					Reason: string(reason),
				},
			})
		}
	}

	if len(gateways) == 0 {
		// TODO https://github.com/Kong/kubernetes-ingress-controller/issues/2417 separate out various rejected reasons
		// and apply specific statuses for those failures in the Route controllers
		return nil, fmt.Errorf(unsupportedGW)
	}

	return gateways, nil
}

func routeHostnamesIntersectsWithListenerHostname[T types.RouteT](route T, listener Listener) bool {
	switch r := (any)(route).(type) {
	case *gatewayv1beta1.HTTPRoute:
		return listenerHostnameIntersectWithRouteHostnames(listener, r.Spec.Hostnames)
	case *gatewayv1alpha2.TCPRoute:
		return true
	case *gatewayv1alpha2.UDPRoute:
		return true
	case *gatewayv1alpha2.TLSRoute:
		return listenerHostnameIntersectWithRouteHostnames(listener, r.Spec.Hostnames)
	default:
		return false
	}
}

func routeTypeMatchesListenerType[T types.RouteT](route T, listener Listener) bool {
	switch (any)(route).(type) {
	case *gatewayv1beta1.HTTPRoute:
		// HTTPRoutes support Terminate only
		// Note: this is a guess we are doing as the upstream documentation is unclear at the moment.
		// see https://github.com/kubernetes-sigs/gateway-api/issues/1474
		if !(listener.Protocol == HTTPProtocolType || listener.Protocol == HTTPSProtocolType) {
			return false
		}
		if listener.TLS != nil && *listener.TLS.Mode != gatewayv1beta1.TLSModeTerminate {
			return false
		}
	case *gatewayv1alpha2.TCPRoute:
		if listener.Protocol != TCPProtocolType {
			return false
		}
		// TCPRoutes support Terminate only
		// Note: this is a guess we are doing as the upstream documentation is unclear at the moment.
		// see https://github.com/kubernetes-sigs/gateway-api/issues/1474
		if listener.TLS != nil && *listener.TLS.Mode != gatewayv1beta1.TLSModeTerminate {
			return false
		}
	case *gatewayv1alpha2.UDPRoute:
		if listener.Protocol != UDPProtocolType {
			return false
		}
		// TLS should not be set in UDP listeners
		if listener.TLS != nil {
			return false
		}
	case *gatewayv1alpha2.TLSRoute:
		if listener.Protocol != TLSProtocolType {
			return false
		}
		// TLSRoutes currently support Passthrough only
		// Note: this is a guess we are doing as the upstream documentation is unclear at the moment.
		// see https://github.com/kubernetes-sigs/gateway-api/issues/1474
		if listener.TLS != nil && *listener.TLS.Mode != gatewayv1beta1.TLSModePassthrough {
			return false
		}
	default:
		return false
	}
	return true
}

// routeMatchesListenerAllowedRoutes checks if the provided route matches the
// criteria defined in listener's AllowedRoutes field.
func routeMatchesListenerAllowedRoutes[T types.RouteT](
	ctx context.Context,
	mgrc client.Client,
	route T,
	listener Listener,
	gatewayNamespace string,
	parentRefNamespace *Namespace,
) (bool, error) {
	if listener.AllowedRoutes == nil {
		return true, nil
	}

	if len(listener.AllowedRoutes.Kinds) > 0 {
		// Find if the route has a type that's within the listener's supported types.
		_, ok := lo.Find(listener.AllowedRoutes.Kinds, func(rgk gatewayv1beta1.RouteGroupKind) bool {
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
	case gatewayv1beta1.NamespacesFromAll:
		return true, nil

	case gatewayv1beta1.NamespacesFromSame:
		// If parentRef didn't specify the namespace then we check if
		// the gateway is from the same namespace as the route
		if parentRefNamespace == nil {
			return gatewayNamespace == route.GetNamespace(), nil
		}
		// Otherwise compare routes namespace with parentRef's one.
		return route.GetNamespace() == string(*parentRefNamespace), nil

	case gatewayv1beta1.NamespacesFromSelector:
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
	errUnsupportedRouteKind             = errors.New("unsupported route kind")
	errUnmatchedListenerName            = errors.New("unmatched listener name")
	errNoReadyConditionFoundForListener = errors.New("no Ready condition found for listener")
	errListenerNotReadyYet              = errors.New("listener not ready yet")
)

// existsMatchingReadyListenerInStatus checks if:
// - If a listener status exists with a matching type (via SupportedKinds).
// - If it matches the requested listener by name (if specified).
// - And finally check if the provided listener is marked as Ready.
func existsMatchingReadyListenerInStatus[T types.RouteT](route T, listener Listener, lss []ListenerStatus) error {
	listenerFound := false

	// Find listener's status...
	listenerStatus, ok := lo.Find(lss, func(ls gatewayv1beta1.ListenerStatus) bool {
		if ls.Name != listener.Name {
			return false
		}
		listenerFound = true

		// Find if the route has a type that's within the supported types, listed
		// in listener's status.
		_, ok := lo.Find(ls.SupportedKinds, func(rgk gatewayv1beta1.RouteGroupKind) bool {
			gvk := route.GetObjectKind().GroupVersionKind()
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

	// ... and verify if it's ready.
	lReadyCond, ok := lo.Find(listenerStatus.Conditions, func(c metav1.Condition) bool {
		return c.Type == string(gatewayv1beta1.ListenerConditionReady)
	})
	if !ok {
		return errNoReadyConditionFoundForListener
	}
	if lReadyCond.Status != "True" {
		return errListenerNotReadyYet // Listener is not ready yet.
	}

	return nil
}

func listenerHostnameIntersectWithRouteHostnames[H types.HostnameT, L types.ListenerT](listener L, hostnames []H) bool {
	if len(hostnames) == 0 {
		return true
	}

	// if the listener has no hostname, all hostnames automatically intersect
	switch l := (any)(listener).(type) {
	case gatewayv1alpha2.Listener:
		if l.Hostname == nil || *l.Hostname == "" {
			return true
		}

		// iterate over all the hostnames and check that at least one intersect with the listener hostname
		for _, hostname := range hostnames {
			if util.HostnamesIntersect(*l.Hostname, hostname) {
				return true
			}
		}
	case Listener:
		if l.Hostname == nil || *l.Hostname == "" {
			return true
		}

		// iterate over all the hostnames and check that at least one intersect with the listener hostname
		for _, hostname := range hostnames {
			if util.HostnamesIntersect(*l.Hostname, hostname) {
				return true
			}
		}
	}

	return false
}

// isListenerHostnameEffective returns true if the listener can specify an effective
// hostname to match hostnames in requests.
// It basically checks if the listener is using any these protocols: HTTP, HTTPS, or TLS.
func isListenerHostnameEffective(listener gatewayv1beta1.Listener) bool {
	return listener.Protocol == gatewayv1beta1.HTTPProtocolType ||
		listener.Protocol == gatewayv1beta1.HTTPSProtocolType ||
		listener.Protocol == gatewayv1beta1.TLSProtocolType
}

// filterHostnames accepts a HTTPRoute and returns a version of the same object with only a subset of the
// hostnames, the ones matching with the listeners' hostname.
// it returns an error if the intersection of hostname match in httproute and listeners is empty.
func filterHostnames(gateways []supportedGatewayWithCondition, httpRoute *gatewayv1beta1.HTTPRoute) (*gatewayv1beta1.HTTPRoute, error) {
	filteredHostnames := make([]gatewayv1beta1.Hostname, 0)
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
func getUnionOfGatewayHostnames(gateways []supportedGatewayWithCondition) ([]gatewayv1beta1.Hostname, bool) {
	hostnames := make([]gatewayv1beta1.Hostname, 0)
	for _, gateway := range gateways {
		if gateway.listenerName != "" {
			if listener := extractListenerSpecFromGateway(
				gateway.gateway,
				gatewayv1beta1.SectionName(gateway.listenerName),
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
func getMinimumHostnameIntersection(gateways []supportedGatewayWithCondition, hostname gatewayv1beta1.Hostname) gatewayv1beta1.Hostname {
	for _, gateway := range gateways {
		for _, listener := range gateway.gateway.Spec.Listeners {
			// if the listenerName is specified and matches the name of the gateway listener proceed
			if (SectionName)(gateway.listenerName) == "" ||
				(SectionName)(gateway.listenerName) == (listener.Name) {
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
		if gateway.condition.Type == string(gatewayv1alpha2.RouteConditionAccepted) && gateway.condition.Status == metav1.ConditionTrue {
			return true
		}
	}
	return false
}

// isHTTPReferenceGranted checks that the backendRef referenced by the HTTPRoute is granted by a ReferenceGrant.
func isHTTPReferenceGranted(grantSpec gatewayv1alpha2.ReferenceGrantSpec, backendRef gatewayv1beta1.HTTPBackendRef, fromNamespace string) bool {
	var backendRefGroup gatewayv1beta1.Group
	var backendRefKind Kind

	if backendRef.Group != nil {
		backendRefGroup = *backendRef.Group
	}
	if backendRef.Kind != nil {
		backendRefKind = *backendRef.Kind
	}
	for _, from := range grantSpec.From {
		if from.Group != gatewayv1beta1.GroupName || from.Kind != "HTTPRoute" || fromNamespace != string(from.Namespace) {
			continue
		}

		for _, to := range grantSpec.To {
			if backendRefGroup == (gatewayv1beta1.Group)(to.Group) &&
				backendRefKind == (Kind)(to.Kind) &&
				(to.Name == nil || (gatewayv1beta1.ObjectName)(*to.Name) == backendRef.Name) {
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
		a.Message == b.Message
}

func setRouteParentStatusCondition[T types.ParentStatusT](parentStatus T, newCondition metav1.Condition) bool {
	var conditionFound, changed bool
	switch p := (any)(parentStatus).(type) {
	case *gatewayv1beta1.RouteParentStatus:
		for i, condition := range p.Conditions {
			if condition.Type == newCondition.Type {
				conditionFound = true
				if !sameCondition(condition, newCondition) {
					p.Conditions[i] = newCondition
					changed = true
				}
			}
		}

		if !conditionFound {
			p.Conditions = append(p.Conditions, newCondition)
			changed = true
		}
	case *gatewayv1alpha2.RouteParentStatus:
		for i, condition := range p.Conditions {
			if condition.Type == newCondition.Type {
				conditionFound = true
				if !sameCondition(condition, newCondition) {
					p.Conditions[i] = newCondition
					changed = true
				}
			}
		}

		if !conditionFound {
			p.Conditions = append(p.Conditions, newCondition)
			changed = true
		}
	}
	return changed
}

func parentStatusHasProgrammedCondition[T types.ParentStatusT](parentStatus T) bool {
	switch p := (any)(parentStatus).(type) {
	case *gatewayv1beta1.RouteParentStatus:
		for _, condition := range p.Conditions {
			if condition.Type == ConditionTypeProgrammed {
				return true
			}
		}
		return false
	case *gatewayv1alpha2.RouteParentStatus:
		for _, condition := range p.Conditions {
			if condition.Type == ConditionTypeProgrammed {
				return true
			}
		}
		return false
	}
	return false
}
