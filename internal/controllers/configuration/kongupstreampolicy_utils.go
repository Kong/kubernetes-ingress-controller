package configuration

import (
	"context"
	"fmt"
	"reflect"
	"sort"

	"github.com/samber/lo"
	"github.com/samber/mo"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	gatewaycontroller "github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/gateway"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1beta1"
	incubatorv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/incubator/v1alpha1"
)

// maxNAncestors is the maximum number of ancestors that can be stored in the KongUpstreamPolicy status.
// This is a limitation of the Gateway API.
const maxNAncestors = 16

// upstreamPolicyAncestorKind represents kind of KongUpstreamPolicy ancestor (Service or KongServiceFacade).
type upstreamPolicyAncestorKind string

const (
	upstreamPolicyAncestorKindService           upstreamPolicyAncestorKind = "Service"
	upstreamPolicyAncestorKindKongServiceFacade upstreamPolicyAncestorKind = "KongServiceFacade"
)

// ancestorStatus represents the status of an ancestor (Service or KongServiceFacade).
// A collection of all ancestors' statuses is used to build the KongUpstreamPolicy status.
type ancestorStatus struct {
	namespacedName      k8stypes.NamespacedName
	ancestorKind        upstreamPolicyAncestorKind
	acceptedCondition   metav1.Condition
	programmedCondition metav1.Condition
	creationTimestamp   metav1.Time
}

// serviceKey is used as a key for indexing Services by "namespace/name".
type serviceKey string

// servicesSet is a set of serviceKeys.
type servicesSet map[serviceKey]struct{}

// enforceKongUpstreamPolicyStatus gets a list of services (ancestors) along with their desired status and enforce them
// in the KongUpstreamPolicy status.
func (r *KongUpstreamPolicyReconciler) enforceKongUpstreamPolicyStatus(
	ctx context.Context,
	oldPolicy *kongv1beta1.KongUpstreamPolicy,
) (bool, error) {
	policyNN := k8stypes.NamespacedName{
		Namespace: oldPolicy.Namespace,
		Name:      oldPolicy.Name,
	}

	// Get all objects (Services and KongServiceFacades) that reference this KongUpstreamPolicy.
	services, err := r.getServicesReferencingUpstreamPolicy(ctx, policyNN)
	if err != nil {
		return false, err
	}
	serviceFacades, err := r.maybeGetServiceFacadesReferencingUpstreamPolicy(ctx, policyNN)
	if err != nil {
		return false, err
	}

	// Build the status for each ancestor.
	ancestorsStatus, err := r.buildAncestorsStatus(ctx, services, serviceFacades)
	if err != nil {
		return false, err
	}

	// Build the desired KongUpstreamPolicy status.
	newPolicyStatus, err := r.buildPolicyStatus(policyNN, ancestorsStatus)
	if err != nil {
		return false, err
	}

	// If the status is not updated, we don't need to patch the KongUpstreamPolicy.
	if isStatusUpdated := isPolicyStatusUpdated(oldPolicy.Status, newPolicyStatus); !isStatusUpdated {
		newPolicy := oldPolicy.DeepCopy()
		newPolicy.Status = newPolicyStatus
		return true, r.Client.Status().Patch(ctx, newPolicy, client.MergeFrom(oldPolicy))
	}
	return false, nil
}

func (r *KongUpstreamPolicyReconciler) getServicesReferencingUpstreamPolicy(
	ctx context.Context,
	upstreamPolicyNN k8stypes.NamespacedName,
) ([]corev1.Service, error) {
	services := &corev1.ServiceList{}
	err := r.List(ctx, services,
		client.InNamespace(upstreamPolicyNN.Namespace),
		client.MatchingFields{
			upstreamPolicyIndexKey: upstreamPolicyNN.Name,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed listing Services: %w", err)
	}
	return services.Items, nil
}

// maybeGetServiceFacadesReferencingUpstreamPolicy returns a list of KongServiceFacades that reference the given KongUpstreamPolicy.
// Skips the lookup if KongServiceFacade is not enabled.
func (r *KongUpstreamPolicyReconciler) maybeGetServiceFacadesReferencingUpstreamPolicy(
	ctx context.Context,
	upstreamPolicyNN k8stypes.NamespacedName,
) ([]incubatorv1alpha1.KongServiceFacade, error) {
	if !r.KongServiceFacadeEnabled {
		// KongServiceFacade is not enabled, so we don't need to check for it.
		return nil, nil
	}
	serviceFacades := &incubatorv1alpha1.KongServiceFacadeList{}
	err := r.List(ctx, serviceFacades,
		client.InNamespace(upstreamPolicyNN.Namespace),
		client.MatchingFields{
			upstreamPolicyIndexKey: upstreamPolicyNN.Name,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed listing KongServiceFacades: %w", err)
	}
	return serviceFacades.Items, nil
}

// buildAncestorsStatus creates a list of services with their conditions associated.
func (r *KongUpstreamPolicyReconciler) buildAncestorsStatus(
	ctx context.Context,
	services []corev1.Service,
	serviceFacades []incubatorv1alpha1.KongServiceFacade,
) ([]ancestorStatus, error) {
	// Check if any Services have conflicts. We do not verify conflicts for KongServiceFacades as there's
	// no scenario in which they would have one.
	conflictedServices, err := r.getConflictedServices(ctx, services)
	if err != nil {
		return nil, err
	}

	// Prepare conditions.
	acceptedCondition := metav1.Condition{
		Type:               string(gatewayapi.PolicyConditionAccepted),
		Status:             metav1.ConditionTrue,
		Reason:             string(gatewayapi.PolicyReasonAccepted),
		LastTransitionTime: metav1.Now(),
	}
	programmedCondition := metav1.Condition{
		Type:               string(gatewayapi.GatewayConditionProgrammed),
		Status:             metav1.ConditionTrue,
		Reason:             string(gatewayapi.GatewayReasonProgrammed),
		LastTransitionTime: metav1.Now(),
	}

	// Build the status for each ancestor (Services and KongServiceFacades).
	ancestorsStatus := make([]ancestorStatus, 0, len(services)+len(serviceFacades))
	for _, service := range services {
		service := service
		acceptedCondition := acceptedCondition
		programmedCondition := programmedCondition

		if _, isConflicted := conflictedServices[buildServiceReference(service.Namespace, service.Name)]; isConflicted {
			// If the Service is conflicted, we change both conditions to False.
			acceptedCondition.Status = metav1.ConditionFalse
			acceptedCondition.Reason = string(gatewayapi.PolicyReasonConflicted)
			programmedCondition.Status = metav1.ConditionFalse
			programmedCondition.Reason = string(gatewayapi.GatewayReasonPending)
		}

		if !r.DataplaneClient.KubernetesObjectIsConfigured(&service) {
			// If the Service is not configured, we change it to False.
			programmedCondition.Status = metav1.ConditionFalse
			programmedCondition.Reason = string(gatewayapi.GatewayReasonPending)
		}

		ancestorsStatus = append(ancestorsStatus, ancestorStatus{
			namespacedName: k8stypes.NamespacedName{
				Namespace: service.Namespace,
				Name:      service.Name,
			},
			ancestorKind:        upstreamPolicyAncestorKindService,
			acceptedCondition:   acceptedCondition,
			programmedCondition: programmedCondition,
		})
	}
	for _, serviceFacade := range serviceFacades {
		serviceFacade := serviceFacade
		programmedCondition := programmedCondition
		if !r.DataplaneClient.KubernetesObjectIsConfigured(&serviceFacade) {
			// If the KongServiceFacade is not configured, we change it to False.
			programmedCondition.Status = metav1.ConditionFalse
			programmedCondition.Reason = string(gatewayapi.GatewayReasonPending)
		}

		ancestorsStatus = append(ancestorsStatus, ancestorStatus{
			namespacedName: k8stypes.NamespacedName{
				Namespace: serviceFacade.Namespace,
				Name:      serviceFacade.Name,
			},
			ancestorKind:        upstreamPolicyAncestorKindKongServiceFacade,
			acceptedCondition:   acceptedCondition,
			programmedCondition: programmedCondition,
		})
	}

	return ancestorsStatus, nil
}

// getConflictedServices returns a set of services that have conflicts.
func (r *KongUpstreamPolicyReconciler) getConflictedServices(ctx context.Context, services []corev1.Service) (servicesSet, error) {
	// return directly when HTTPRoute is not enabled, as it only check conflicted services in HTTPRoute backends only.
	if !r.HTTPRouteEnabled {
		return make(servicesSet), nil
	}
	// Prepare a mapping for efficient lookups if a Service uses this KongUpstreamPolicy.
	upstreamPolicyServices := make(servicesSet)
	for _, service := range services {
		upstreamPolicyServices[buildServiceReference(service.Namespace, service.Name)] = struct{}{}
	}

	conflictedServices := make(servicesSet)
	for serviceKey := range upstreamPolicyServices {
		// We fetch all the HTTPRoutes that reference this service.
		httpRoutes := &gatewayapi.HTTPRouteList{}
		err := r.List(ctx, httpRoutes,
			client.MatchingFields{
				routeBackendRefServiceNameIndexKey: string(serviceKey),
			},
		)
		if err != nil {
			return nil, err
		}
		hasConflict := lo.ContainsBy(httpRoutes.Items, func(httpRoute gatewayapi.HTTPRoute) bool {
			return httpRouteHasUpstreamPolicyConflictedBackendRefsWithService(httpRoute, upstreamPolicyServices, serviceKey)
		})
		if hasConflict {
			conflictedServices[serviceKey] = struct{}{}
		}
	}
	return conflictedServices, nil
}

// httpRouteHasUpstreamPolicyConflictedBackendRefsWithService checks if there's any HTTPRoute's rule that uses multiple backendRefs
// AND they're not all using the same KongUpstreamPolicy.
// If so, that means that we have a conflict because we cannot apply multiple KongUpstreamPolicy to the same Kong Service.
func httpRouteHasUpstreamPolicyConflictedBackendRefsWithService(
	httpRoute gatewayapi.HTTPRoute,
	upstreamPolicyServices servicesSet,
	serviceKey serviceKey,
) bool {
	backendRefsUsedWithThisService := getAllBackendRefsUsedWithService(httpRoute, serviceKey)
	hasAnyBackendRefNotUsingSameUpstreamPolicy := lo.ContainsBy(backendRefsUsedWithThisService, func(br gatewayapi.HTTPBackendRef) bool {
		serviceRef := backendRefToServiceRef(httpRoute.Namespace, br.BackendRef)
		if serviceRef == "" {
			return false
		}
		// If the serviceRef is not in the upstreamPolicyServices, it means it doesn't use this KongUpstreamPolicy.
		_, ok := upstreamPolicyServices[serviceRef]
		return !ok
	})
	return hasAnyBackendRefNotUsingSameUpstreamPolicy
}

// getAllBackendRefsUsedWithService returns HTTPRoute's backendRefs that use the given service (excluding the given service).
func getAllBackendRefsUsedWithService(httpRoute gatewayapi.HTTPRoute, serviceKey serviceKey) []gatewayapi.HTTPBackendRef {
	var backendRefs []gatewayapi.HTTPBackendRef
	for _, rule := range httpRoute.Spec.Rules {
		// We will look for a backendRef that matches the given service and keep its index if found.
		backendRefMatchingServiceIdx := mo.None[int]()
		for i, br := range rule.BackendRefs {
			serviceRef := backendRefToServiceRef(httpRoute.Namespace, br.BackendRef)
			if serviceRef == serviceKey {
				// We found a backendRef that matches the given service, no need to look further.
				backendRefMatchingServiceIdx = mo.Some(i)
				break
			}
		}
		if matchingIdx, ok := backendRefMatchingServiceIdx.Get(); ok {
			// We found a backendRef that matches the given service. We will keep all the backendRefs that are together
			// with this backendRef in the rule.

			backendRefs = append(backendRefs, rule.BackendRefs[:matchingIdx]...)   // We do not keep the reference to rule.BackendRefs, but copy it.
			backendRefs = append(backendRefs, rule.BackendRefs[matchingIdx+1:]...) // We do not keep the reference to rule.BackendRefs, but copy it.
		}
	}
	return backendRefs
}

// buildPolicyStatus builds the KongUpstreamPolicy status from the ancestors' statuses.
// It ensures that the number of ancestors is not greater than the maximum allowed by the Gateway API
// and that the oldest ancestors are kept.
func (r *KongUpstreamPolicyReconciler) buildPolicyStatus(
	upstreamPolicyNN k8stypes.NamespacedName,
	ancestorsStatus []ancestorStatus,
) (gatewayapi.PolicyStatus, error) {
	// Sort the ancestors by creation timestamp and keep only the oldest ones.
	sort.Slice(ancestorsStatus, func(i, j int) bool {
		return ancestorsStatus[i].creationTimestamp.Before(&ancestorsStatus[j].creationTimestamp)
	})
	if len(ancestorsStatus) > maxNAncestors {
		r.Log.Info("status has more ancestors than the Gateway API permits, the newest ones will be ignored",
			"KongUpstreamPolicy", upstreamPolicyNN.String(),
			"ancestorsCount", len(ancestorsStatus),
			"maxAllowedAncestors", maxNAncestors,
		)
		ancestorsStatus = ancestorsStatus[:maxNAncestors]
	}

	// Populate the KongUpstreamPolicy status with the ancestors' statuses.
	policyStatus := gatewayapi.PolicyStatus{}
	if len(ancestorsStatus) > 0 {
		policyStatus.Ancestors = make([]gatewayapi.PolicyAncestorStatus, 0, len(ancestorsStatus))
	}
	for _, ss := range ancestorsStatus {
		ancestorRef, err := ancestorRef(ss.namespacedName, ss.ancestorKind)
		if err != nil {
			return gatewayapi.PolicyStatus{}, fmt.Errorf("failed to build ancestor reference: %w", err)
		}
		policyStatus.Ancestors = append(policyStatus.Ancestors,
			gatewayapi.PolicyAncestorStatus{
				AncestorRef:    ancestorRef,
				ControllerName: gatewaycontroller.GetControllerName(),
				Conditions: []metav1.Condition{
					ss.acceptedCondition,
					ss.programmedCondition,
				},
			},
		)
	}

	return policyStatus, nil
}

func ancestorRef(nn k8stypes.NamespacedName, kind upstreamPolicyAncestorKind) (gatewayapi.ParentReference, error) {
	switch kind {
	case upstreamPolicyAncestorKindService:
		return gatewayapi.ParentReference{
			Group:     lo.ToPtr(gatewayapi.Group("core")),
			Kind:      lo.ToPtr(gatewayapi.Kind("Service")),
			Namespace: lo.ToPtr(gatewayapi.Namespace(nn.Namespace)),
			Name:      gatewayapi.ObjectName(nn.Name),
		}, nil
	case upstreamPolicyAncestorKindKongServiceFacade:
		return gatewayapi.ParentReference{
			Group:     lo.ToPtr(gatewayapi.Group(incubatorv1alpha1.GroupVersion.Group)),
			Kind:      lo.ToPtr(gatewayapi.Kind(incubatorv1alpha1.KongServiceFacadeKind)),
			Namespace: lo.ToPtr(gatewayapi.Namespace(nn.Namespace)),
			Name:      gatewayapi.ObjectName(nn.Name),
		}, nil
	}

	return gatewayapi.ParentReference{}, fmt.Errorf("unknown ancestor kind %q", kind)
}

func isPolicyStatusUpdated(oldStatus, newStatus gatewayapi.PolicyStatus) bool {
	if len(oldStatus.Ancestors) != len(newStatus.Ancestors) {
		return false
	}
	for i, oldAncestor := range oldStatus.Ancestors {
		newAncestor := newStatus.Ancestors[i]
		if newAncestor.ControllerName != oldAncestor.ControllerName {
			return false
		}
		if !reflect.DeepEqual(newAncestor.AncestorRef, oldAncestor.AncestorRef) {
			return false
		}

		if len(oldAncestor.Conditions) != len(newAncestor.Conditions) {
			return false
		}
		for j, oldCondition := range oldAncestor.Conditions {
			newCondition := newAncestor.Conditions[j]
			if newCondition.Type != oldCondition.Type ||
				newCondition.Status != oldCondition.Status ||
				newCondition.Reason != oldCondition.Reason ||
				newCondition.Message != oldCondition.Message {
				return false
			}
		}
	}

	return true
}

func backendRefToServiceRef(routeNamespace string, br gatewayapi.BackendRef) serviceKey {
	if !isSupportedHTTPRouteBackendRef(br) {
		return ""
	}
	namespace := routeNamespace
	if br.Namespace != nil {
		namespace = string(*br.Namespace)
	}
	return buildServiceReference(namespace, string(br.Name))
}

func buildServiceReference(namespace, name string) serviceKey {
	return serviceKey(fmt.Sprintf("%s/%s", namespace, name))
}

func isSupportedHTTPRouteBackendRef(br gatewayapi.BackendRef) bool {
	groupIsCoreOrNilOrEmpty := br.Group == nil || *br.Group == "core" || *br.Group == ""
	kindIsServiceOrNil := br.Kind == nil || *br.Kind == "Service"

	// We only support core Services.
	// For Group the specification says when it's unspecified (nil or empty string), core API group should be inferred.
	// For Kind nil case should never happen as it defaults on the API level to 'Service'. We can safely consider
	// nil to be treated as 'Service' if it would happen for any reason.
	return groupIsCoreOrNilOrEmpty && kindIsServiceOrNil
}
