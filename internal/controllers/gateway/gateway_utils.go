package gateway

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
)

// -----------------------------------------------------------------------------
// Gateway Utilities
// -----------------------------------------------------------------------------

const (
	// maxConds is the maximum number of status conditions a Gateway can have at one time.
	maxConds = 8
)

// isGatewayScheduled returns boolean whether or not the gateway object was scheduled
// previously by the gateway controller.
func isGatewayScheduled(gateway *gatewayv1alpha2.Gateway) bool {
	for _, cond := range gateway.Status.Conditions {
		if cond.Type == string(gatewayv1alpha2.GatewayConditionScheduled) &&
			cond.Reason == string(gatewayv1alpha2.GatewayReasonScheduled) &&
			cond.Status == metav1.ConditionTrue {
			return true
		}
	}
	return false
}

// isGatewayReady returns boolean whether the ready condition exists
// for the given gateway object if it matches the currently known generation of that object.
func isGatewayReady(gateway *gatewayv1alpha2.Gateway) bool {
	for _, cond := range gateway.Status.Conditions {
		if cond.Type == string(gatewayv1alpha2.GatewayConditionReady) && cond.Reason == string(gatewayv1alpha2.GatewayReasonReady) && cond.ObservedGeneration == gateway.Generation {
			return true
		}
	}
	return false
}

// isGatewayInClassAndUnmanaged returns boolean if the provided combination of gateway and class
// is controlled by this controller and the gateway is configured for unmanaged mode.
func isGatewayInClassAndUnmanaged(gatewayClass *gatewayv1alpha2.GatewayClass, gateway gatewayv1alpha2.Gateway) bool {
	_, ok := annotations.ExtractUnmanagedGatewayMode(gateway.Annotations)
	return ok && gatewayClass.Spec.ControllerName == ControllerName
}

// getRefFromPublishService splits a publish service string in the format namespace/name into a types.NamespacedName
// and verifies the contents producing an error if they don't match namespace/name format.
func getRefFromPublishService(publishService string) (types.NamespacedName, error) {
	publishServiceSplit := strings.SplitN(publishService, "/", 3)
	if len(publishServiceSplit) != 2 {
		return types.NamespacedName{}, fmt.Errorf("--publish-service expected in format 'namespace/name' but got %s", publishService)
	}
	return types.NamespacedName{
		Namespace: publishServiceSplit[0],
		Name:      publishServiceSplit[1],
	}, nil
}

// pruneGatewayStatusConds cleans out old status conditions if the Gateway currently has more
// status conditions set than the 8 maximum allowed by the Kubernetes API.
func pruneGatewayStatusConds(gateway *gatewayv1alpha2.Gateway) *gatewayv1alpha2.Gateway {
	if len(gateway.Status.Conditions) > maxConds {
		gateway.Status.Conditions = gateway.Status.Conditions[len(gateway.Status.Conditions)-maxConds:]
	}
	return gateway
}

// reconcileGatewaysIfClassMatches is a filter function to convert a list of gateways into a list
// of reconciliation requests for those gateways based on which match the given class.
func reconcileGatewaysIfClassMatches(gatewayClass client.Object, gateways []gatewayv1alpha2.Gateway) (recs []reconcile.Request) {
	for _, gateway := range gateways {
		if string(gateway.Spec.GatewayClassName) == gatewayClass.GetName() {
			recs = append(recs, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: gateway.Namespace,
					Name:      gateway.Name,
				},
			})
		}
	}
	return
}

// ListenerTracker holds Gateway Listeners and their statuses, and provides methods to update statuses upon
// reconciliation
type ListenerTracker struct {
	// actual listeners
	Listeners map[gatewayv1alpha2.SectionName]gatewayv1alpha2.Listener

	// statuses
	Statuses map[gatewayv1alpha2.SectionName]gatewayv1alpha2.ListenerStatus
	// protocol to port to number map (var protocols)
	protocolToPort map[gatewayv1alpha2.ProtocolType]map[gatewayv1alpha2.PortNumber]bool
	// port to protocol map (portsToProtocol)
	portToProtocol map[gatewayv1alpha2.PortNumber]gatewayv1alpha2.ProtocolType
	// port to hostname to listener name map (portsToHostnames)
	portsToHostnames map[gatewayv1alpha2.PortNumber]map[gatewayv1alpha2.Hostname]gatewayv1alpha2.SectionName
}

// update from existing becomes moot if we're stateful, correct?
// we just keep the existing maps around
// need to detect changes, will still receive the full set

// NewListenerTracker returns a ListenerTracker with empty maps
func NewListenerTracker() ListenerTracker {
	return ListenerTracker{
		Statuses:         map[gatewayv1alpha2.SectionName]gatewayv1alpha2.ListenerStatus{},
		Listeners:        map[gatewayv1alpha2.SectionName]gatewayv1alpha2.Listener{},
		protocolToPort:   map[gatewayv1alpha2.ProtocolType]map[gatewayv1alpha2.PortNumber]bool{},
		portToProtocol:   map[gatewayv1alpha2.PortNumber]gatewayv1alpha2.ProtocolType{},
		portsToHostnames: map[gatewayv1alpha2.PortNumber]map[gatewayv1alpha2.Hostname]gatewayv1alpha2.SectionName{},
	}
}

type protocolPortMap map[gatewayv1alpha2.ProtocolType]map[gatewayv1alpha2.PortNumber]bool
type portProtocolMap map[gatewayv1alpha2.PortNumber]gatewayv1alpha2.ProtocolType
type portHostnameMap map[gatewayv1alpha2.PortNumber]map[gatewayv1alpha2.Hostname]bool
type listenerAttachedMap map[gatewayv1alpha2.SectionName]int32

func buildKongPortMap(listens []gatewayv1alpha2.Listener) protocolPortMap {
	p := make(map[gatewayv1alpha2.ProtocolType]map[gatewayv1alpha2.PortNumber]bool, len(listens))
	for _, listen := range listens {
		_, ok := p[listen.Protocol]
		if !ok {
			p[listen.Protocol] = map[gatewayv1alpha2.PortNumber]bool{}
		}
		p[listen.Protocol][listen.Port] = true
	}
	return p
}

// initializeListenerMaps takes a Gateway and builds indices used in status updates and conflict detection. It returns
// empty maps from port to protocol to listener name and from port to hostnames, and a populated map from listener name
// to attached route count from their status
func initializeListenerMaps(gateway *gatewayv1alpha2.Gateway) (
	portProtocolMap,
	portHostnameMap,
	listenerAttachedMap,
) {
	portToProtocol := make(portProtocolMap, len(gateway.Status.Listeners))
	portToHostname := make(portHostnameMap, len(gateway.Status.Listeners))
	listenerToAttached := make(listenerAttachedMap, len(gateway.Status.Listeners))

	existingStatuses := make(map[gatewayv1alpha2.SectionName]gatewayv1alpha2.ListenerStatus,
		len(gateway.Status.Listeners))
	for _, listenerStatus := range gateway.Status.Listeners {
		existingStatuses[listenerStatus.Name] = listenerStatus
	}

	for _, listener := range gateway.Spec.Listeners {
		portToHostname[listener.Port] = make(map[gatewayv1alpha2.Hostname]bool)
		if existingStatus, ok := existingStatuses[listener.Name]; ok {
			listenerToAttached[listener.Name] = existingStatus.AttachedRoutes
		} else {
			listenerToAttached[listener.Name] = 0
		}
	}
	return portToProtocol, portToHostname, listenerToAttached
}

func canSharePort(requested gatewayv1alpha2.ProtocolType, existing gatewayv1alpha2.ProtocolType) bool {
	switch requested {
	// TCP and UDP listeners must always use unique ports
	case gatewayv1alpha2.TCPProtocolType, gatewayv1alpha2.UDPProtocolType:
		return false
	// HTTPS and TLS Listeners can share ports with others of their type or the other TLS type
	// note that this is not actually possible in Kong: TLS is a stream listen and HTTPS is an http listen
	// however, this section implements the spec ignoring Kong's reality
	case gatewayv1alpha2.HTTPSProtocolType:
		if existing == gatewayv1alpha2.HTTPSProtocolType || existing == gatewayv1alpha2.TLSProtocolType {
			return true
		}
		return false
	case gatewayv1alpha2.TLSProtocolType:
		if existing == gatewayv1alpha2.HTTPSProtocolType || existing == gatewayv1alpha2.TLSProtocolType {
			return true
		}
		return false
	// HTTP Listeners can share ports with others of the same protocol only
	case gatewayv1alpha2.HTTPProtocolType:
		if existing == gatewayv1alpha2.HTTPProtocolType {
			return true
		}
		return false
	default:
		return false
	}
}

func getListenerStatus(
	gateway *gatewayv1alpha2.Gateway,
	kongListens []gatewayv1alpha2.Listener,
) []gatewayv1alpha2.ListenerStatus {
	statuses := make(map[gatewayv1alpha2.SectionName]gatewayv1alpha2.ListenerStatus, len(gateway.Spec.Listeners))
	portToProtocol, portToHostname, listenerToAttached := initializeListenerMaps(gateway)
	kongProtocolsToPort := buildKongPortMap(kongListens)
	conflictedPorts := make(map[gatewayv1alpha2.PortNumber]bool, len(gateway.Spec.Listeners))
	conflictedHostnames := make(map[gatewayv1alpha2.PortNumber]map[gatewayv1alpha2.Hostname]bool, len(gateway.Spec.Listeners))

	// TODO we should check transition time rather than always nowing, which we do throughout the below
	// https://github.com/Kong/kubernetes-ingress-controller/issues/2556
	for _, listener := range gateway.Spec.Listeners {
		var hostname gatewayv1alpha2.Hostname
		if listener.Hostname != nil {
			hostname = *listener.Hostname
		}
		status := gatewayv1alpha2.ListenerStatus{
			Name:           listener.Name,
			Conditions:     []metav1.Condition{},
			SupportedKinds: supportedRouteGroupKinds,
			// this has been populated by initializeListenerMaps()
			AttachedRoutes: listenerToAttached[listener.Name],
		}
		// TODO this only handles some Listener conditions and reasons as needed to check cross-listener compatibility
		// and unattachability due to missing Kong configuration. There are others available and it may be appropriate
		// for us to add them https://github.com/Kong/kubernetes-ingress-controller/issues/2558
		if _, ok := portToProtocol[listener.Port]; !ok {
			// unoccupied ports are free game
			portToProtocol[listener.Port] = listener.Protocol
			portToHostname[listener.Port][hostname] = true
		} else {
			if !canSharePort(listener.Protocol, portToProtocol[listener.Port]) {
				status.Conditions = append(status.Conditions, metav1.Condition{
					Type:               string(gatewayv1alpha2.ListenerConditionConflicted),
					Status:             metav1.ConditionTrue,
					ObservedGeneration: gateway.Generation,
					LastTransitionTime: metav1.Now(),
					Reason:             string(gatewayv1alpha2.ListenerReasonProtocolConflict),
				})
				conflictedPorts[listener.Port] = true
			} else {
				// shareable ports determine conflicts by hostname
				// Each Listener within the group specifies a Hostname that is unique within the group.
				// As a special case, one Listener within a group may omit Hostname, in which case this Listener
				// matches when no other Listener matches.
				if conflictedHostnames[listener.Port] == nil {
					conflictedHostnames[listener.Port] = map[gatewayv1alpha2.Hostname]bool{}
				}
				if _, exists := portToHostname[listener.Port][hostname]; !exists {
					portToHostname[listener.Port][hostname] = true
				} else {
					status.Conditions = append(status.Conditions, metav1.Condition{
						Type:               string(gatewayv1alpha2.ListenerConditionConflicted),
						Status:             metav1.ConditionTrue,
						ObservedGeneration: gateway.Generation,
						LastTransitionTime: metav1.Now(),
						Reason:             string(gatewayv1alpha2.ListenerReasonHostnameConflict),
					})
					conflictedHostnames[listener.Port][hostname] = true
				}
			}
		}

		// independent of conflict detection. for example, two TCP Listeners both requesting the same port that Kong
		// does not provide should be both Conflicted and Detached
		if len(kongProtocolsToPort[listener.Protocol]) == 0 {
			status.Conditions = append(status.Conditions, metav1.Condition{
				Type:               string(gatewayv1alpha2.ListenerConditionDetached),
				Status:             metav1.ConditionTrue,
				ObservedGeneration: gateway.Generation,
				LastTransitionTime: metav1.Now(),
				Reason:             string(gatewayv1alpha2.ListenerReasonUnsupportedProtocol),
				Message:            "no Kong listen with the requested protocol is configured",
			})
		}
		if _, ok := kongProtocolsToPort[listener.Protocol][listener.Port]; !ok {
			status.Conditions = append(status.Conditions, metav1.Condition{
				Type:               string(gatewayv1alpha2.ListenerConditionDetached),
				Status:             metav1.ConditionTrue,
				ObservedGeneration: gateway.Generation,
				LastTransitionTime: metav1.Now(),
				Reason:             string(gatewayv1alpha2.ListenerReasonPortUnavailable),
				Message:            "no Kong listen with the requested protocol is configured for the requested port",
			})
		}

		// finalize adding any general conditions
		// TODO these (and really the others too) do not account for the conditions maybe having already been present
		// we simply generate them from scratch each time and mark the current generation the observed generation,
		// whereas we should preserve the original observed generation
		// https://github.com/Kong/kubernetes-ingress-controller/issues/2556
		if len(status.Conditions) == 0 {
			// if we've gotten this far with no conditions, the listener is good to go
			status.Conditions = append(status.Conditions,
				metav1.Condition{
					Type:               string(gatewayv1alpha2.ListenerConditionConflicted),
					Status:             metav1.ConditionFalse,
					ObservedGeneration: gateway.Generation,
					LastTransitionTime: metav1.Now(),
					Reason:             string(gatewayv1alpha2.ListenerReasonNoConflicts),
				},
				metav1.Condition{
					Type:               string(gatewayv1alpha2.ListenerConditionReady),
					Status:             metav1.ConditionTrue,
					ObservedGeneration: gateway.Generation,
					LastTransitionTime: metav1.Now(),
					Reason:             string(gatewayv1alpha2.ListenerReasonReady),
					Message:            "the listener is ready and available for routing",
				},
			)
		} else {
			// any conditions we added above will prevent the Listener from becoming ready
			// unsure if we want to add the ready=false condition on a per-failure basis or use this else to just mark
			// it generic unready if we hit anything bad. do any failure conditions block readiness? do we care about
			// having distinct ready false messages, assuming we have more descriptive messages in the other conditions?
			status.Conditions = append(status.Conditions, metav1.Condition{
				Type:               string(gatewayv1alpha2.ListenerConditionReady),
				Status:             metav1.ConditionFalse,
				ObservedGeneration: gateway.Generation,
				LastTransitionTime: metav1.Now(),
				Reason:             string(gatewayv1alpha2.ListenerReasonPending),
				Message:            "the listener is not ready and cannot route requests",
			})
		}

		statuses[listener.Name] = status
	}

	// any conflict applies to all listeners sharing the conflicted resource (see
	// https://github.com/Kong/kubernetes-ingress-controller/pull/2555#issuecomment-1154579046 for discussion)
	// if we encountered conflicts, we must strip the ready status we originally set
	for _, listener := range gateway.Spec.Listeners {
		var reason string

		var hostname gatewayv1alpha2.Hostname
		if listener.Hostname != nil {
			hostname = *listener.Hostname
		}
		// there's no filter for protocols that don't use Hostname, but this won't be populated from earlier for those
		if _, ok := conflictedHostnames[listener.Port][hostname]; ok {
			reason = string(gatewayv1alpha2.ListenerReasonHostnameConflict)
		}

		if _, ok := conflictedPorts[listener.Port]; ok {
			reason = string(gatewayv1alpha2.ListenerReasonProtocolConflict)
		}

		if len(reason) > 0 {
			newConditions := []metav1.Condition{}
			for _, cond := range statuses[listener.Name].Conditions {
				switch gatewayv1alpha2.ListenerConditionType(cond.Type) {
				case gatewayv1alpha2.ListenerConditionReady, gatewayv1alpha2.ListenerConditionConflicted:
					continue
				default:
					newConditions = append(newConditions, cond)
				}
			}
			newConditions = append(newConditions,
				metav1.Condition{
					Type:               string(gatewayv1alpha2.ListenerConditionConflicted),
					Status:             metav1.ConditionTrue,
					ObservedGeneration: gateway.Generation,
					LastTransitionTime: metav1.Now(),
					Reason:             reason,
				},
				metav1.Condition{
					Type:               string(gatewayv1alpha2.ListenerConditionReady),
					Status:             metav1.ConditionFalse,
					ObservedGeneration: gateway.Generation,
					LastTransitionTime: metav1.Now(),
					Reason:             string(gatewayv1alpha2.ListenerReasonReady),
					Message:            "the listener is not ready and cannot route requests",
				},
			)
			status := statuses[listener.Name]
			status.Conditions = newConditions
			statuses[listener.Name] = status
		}
	}
	statusArray := []gatewayv1alpha2.ListenerStatus{}
	for _, status := range statuses {
		statusArray = append(statusArray, status)
	}

	return statusArray
}

// -----------------------------------------------------------------------------
// Gateway Utils - Watch Predicate Helpers
// -----------------------------------------------------------------------------

// isGatewayClassEventInClass produces a boolean whether or not a given event which contains
// one or more GatewayClass objects is supported by this controller according to those
// objects ControllerName.
func isGatewayClassEventInClass(log logr.Logger, watchEvent interface{}) bool {
	objs := make([]client.Object, 0, 2)
	switch e := watchEvent.(type) {
	case event.CreateEvent:
		objs = append(objs, e.Object)
	case event.DeleteEvent:
		objs = append(objs, e.Object)
	case event.GenericEvent:
		objs = append(objs, e.Object)
	case event.UpdateEvent:
		objs = append(objs, e.ObjectOld)
		objs = append(objs, e.ObjectNew)
	default:
		log.Error(fmt.Errorf("invalid type"), "received invalid event type in event handlers", "found", reflect.TypeOf(watchEvent))
		return false
	}

	for _, obj := range objs {
		gwc, ok := obj.(*gatewayv1alpha2.GatewayClass)
		if !ok {
			log.Error(fmt.Errorf("invalid type"), "received invalid object type in event handlers", "expected", "GatewayClass", "found", reflect.TypeOf(obj))
			continue
		}
		if gwc.Spec.ControllerName == ControllerName {
			return true
		}
	}

	return false
}
