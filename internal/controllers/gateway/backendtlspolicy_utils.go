package gateway

import (
	"context"
	"sort"
	"strings"

	"github.com/samber/lo"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

// getBackendTLSPoliciesByHTTPRoute returns the list of BackendTLSPolicies that targets a service
// used as backend by the given HTTPRoute.
func (r *BackendTLSPolicyReconciler) getBackendTLSPoliciesByHTTPRoute(ctx context.Context, httpRoute gatewayapi.HTTPRoute) ([]gatewayapi.BackendTLSPolicy, error) {
	objects := []gatewayapi.BackendTLSPolicy{}
	for _, rule := range httpRoute.Spec.Rules {
		for _, backend := range rule.BackendRefs {
			// no need to check group and kind nilness, as they have a default value in case not specified
			if !((*backend.Group == "" || *backend.Group == "core") && *backend.Kind == "Service") {
				continue
			}

			namespace := httpRoute.Namespace
			if backend.Namespace != nil {
				namespace = string(*backend.Namespace)
			}
			policies := &gatewayapi.BackendTLSPolicyList{}
			if err := r.List(ctx, policies,
				client.InNamespace(namespace),
				client.MatchingFields{backendTLSPolicyTargetRefIndexKey: string(backend.Name)},
			); err != nil {
				return nil, err
			}
			objects = append(objects, policies.Items...)
		}
	}
	return objects, nil
}

// getBackendTLSPolicyAncestors returns the list of Gateways associated to the given BackendTLSPolicy.
// To retrieve such a list, the following steps are performed:
// 1. Get all the the HTTPRoutes that reference the backends targeted by the policy;
// 2. Find all the parents in the HTTPRoute status that have already properly resolved by KIC and have the resolvedRefs condition set to true.
// 3. Return all the successfully resolved Gateways.
func (r *BackendTLSPolicyReconciler) getBackendTLSPolicyAncestors(ctx context.Context, policy gatewayapi.BackendTLSPolicy) ([]gatewayapi.Gateway, error) {
	gateways := []gatewayapi.Gateway{}
	for _, targetRef := range policy.Spec.TargetRefs {
		if (targetRef.Group != "core" && targetRef.Group != "") && targetRef.Kind != "Service" {
			continue
		}

		httpRoutes := gatewayapi.HTTPRouteList{}
		if err := r.Client.List(ctx, &httpRoutes,
			client.MatchingFields{httpRouteBackendRefIndexKey: policy.Namespace + "/" + string(targetRef.Name)},
		); err != nil {
			return nil, err
		}

		for _, httpRoute := range httpRoutes.Items {
			for _, parentRef := range httpRoute.Spec.ParentRefs {
				if *parentRef.Group != gatewayapi.V1Group || *parentRef.Kind != "Gateway" {
					continue
				}

				namespace := httpRoute.Namespace
				if parentRef.Namespace != nil {
					namespace = string(*parentRef.Namespace)
				}

				var resolvedRefsStatus bool
				// Check the resolvedRefs condition is set to true to ensure that all the references are properly resolved
				// and granted by ReferenceGrants.
				for _, parentStatus := range httpRoute.Status.Parents {
					if parentStatus.ControllerName == GetControllerName() &&
						*parentStatus.ParentRef.Group == *parentRef.Group &&
						*parentStatus.ParentRef.Kind == *parentRef.Kind &&
						parentStatus.ParentRef.Name == parentRef.Name &&
						*parentStatus.ParentRef.Namespace == gatewayapi.Namespace(namespace) &&
						r.GatewayNN.MatchesNN(k8stypes.NamespacedName{Namespace: namespace, Name: string(parentRef.Name)}) {
						if _, found := lo.Find(parentStatus.Conditions, func(c metav1.Condition) bool {
							return c.Type == string(gatewayapi.RouteConditionResolvedRefs) && c.Status == metav1.ConditionTrue
						}); found {
							resolvedRefsStatus = true
							break
						}
					}
				}
				if resolvedRefsStatus {
					gateway := &gatewayapi.Gateway{}
					if err := r.Get(ctx, client.ObjectKey{Namespace: namespace, Name: string(parentRef.Name)}, gateway); err != nil {
						// In case the object is not found, we don't want to return an error, but we want to continue.
						if apierrors.IsNotFound(err) {
							continue
						}
						return nil, err
					}
					gateways = append(gateways, *gateway)
				}
			}
		}
	}

	return gateways, nil
}

// setPolicyStatus enforces an ancestorStatus for each Gateway associated to the given policy.
func (r *BackendTLSPolicyReconciler) setPolicyStatus(ctx context.Context, policy gatewayapi.BackendTLSPolicy, gateways []gatewayapi.Gateway, acceptedCondition metav1.Condition) error {
	ancestors := []gatewayapi.PolicyAncestorStatus{}

	var completeAcceptedCondition *metav1.Condition
	// First copy all the ancestorstatuses managed by other controllers.
	kicAncestors := []gatewayapi.PolicyAncestorStatus{}
	for _, ancestor := range policy.Status.Ancestors {
		if ancestor.ControllerName == GetControllerName() {
			kicAncestors = append(kicAncestors, ancestor)
			if completeAcceptedCondition == nil {
				completeAcceptedCondition = getCompleteAcceptedCondition(ancestor, acceptedCondition)
			}
			continue
		}
		ancestors = append(ancestors, ancestor)
	}
	if completeAcceptedCondition == nil {
		completeAcceptedCondition = &acceptedCondition
		completeAcceptedCondition.LastTransitionTime = metav1.Now()
	}

	// Sort the Gateways to be consistent across subsequent reconciliation loops.
	sortGateways(gateways, kicAncestors, policy.Namespace)

	// Then enforces all the ancestorsStatuses for the Gateways managed by this controller.
	for _, gateway := range gateways {
		// The ancestors are limited to 16, as per the Gateway API specification. in case more Gateways are found, we stop.
		if len(ancestors) >= 16 {
			break
		}
		ancestor := gatewayapi.PolicyAncestorStatus{
			AncestorRef: gatewayapi.ParentReference{
				Group:     lo.ToPtr(gatewayapi.V1Group),
				Kind:      lo.ToPtr(gatewayapi.Kind("Gateway")),
				Name:      gatewayapi.ObjectName(gateway.Name),
				Namespace: lo.ToPtr(gatewayapi.Namespace(gateway.Namespace)),
			},
			ControllerName: GetControllerName(),
			Conditions:     []metav1.Condition{*completeAcceptedCondition},
		}

		ancestors = append(ancestors, ancestor)
	}

	newPolicy := policy.DeepCopy()
	newPolicy.Status.Ancestors = ancestors

	return r.Status().Patch(ctx, newPolicy, client.MergeFrom(&policy))
}

func getCompleteAcceptedCondition(ancestors gatewayapi.PolicyAncestorStatus, acceptedCondition metav1.Condition) *metav1.Condition {
	for _, condition := range ancestors.Conditions {
		if condition.Type == acceptedCondition.Type &&
			condition.Status == acceptedCondition.Status &&
			condition.Reason == acceptedCondition.Reason &&
			condition.Message == acceptedCondition.Message {
			acceptedCondition.LastTransitionTime = condition.LastTransitionTime
			return &acceptedCondition
		}
	}
	acceptedCondition.LastTransitionTime = metav1.Now()
	return &acceptedCondition
}

// sortGateways sorts the given slice of Gateway objects by namespace and name.
func sortGateways(gateways []gatewayapi.Gateway, kicAncestors []gatewayapi.PolicyAncestorStatus, policyNamespace string) {
	kicAncestorsMap := lo.SliceToMap(kicAncestors, func(ancestor gatewayapi.PolicyAncestorStatus) (string, gatewayapi.PolicyAncestorStatus) {
		namespace := policyNamespace
		if ancestor.AncestorRef.Namespace != nil {
			namespace = string(*ancestor.AncestorRef.Namespace)
		}
		return namespace + "/" + string(ancestor.AncestorRef.Name), ancestor
	})
	sort.Slice(gateways, func(i, j int) bool {
		_, foundi := kicAncestorsMap[gateways[i].Namespace+"/"+gateways[i].Name]
		_, foundj := kicAncestorsMap[gateways[j].Namespace+"/"+gateways[j].Name]
		switch {
		// the precedence is on Gateways already set in the policy status.
		case foundi && !foundj:
			return true
		case !foundi && foundj:
			return false
		// then we sort by namespace/name.
		case gateways[i].Namespace < gateways[j].Namespace:
			return true
		case gateways[i].Namespace > gateways[j].Namespace:
			return false
		default:
			return gateways[i].Name < gateways[j].Name
		}
	})
}

// validateBackendTLSPolicy validates the given BackendTLSPolicy and returns the accepted Condition related to the policy.
func (r *BackendTLSPolicyReconciler) validateBackendTLSPolicy(ctx context.Context, policy gatewayapi.BackendTLSPolicy) (acceptedCondition *metav1.Condition, err error) {
	acceptedCondition = &metav1.Condition{
		Type:               string(gatewayapi.PolicyConditionAccepted),
		Status:             metav1.ConditionTrue,
		Reason:             string(gatewayapi.PolicyConditionAccepted),
		ObservedGeneration: policy.Generation,
	}

	for _, targetRef := range policy.Spec.TargetRefs {
		if (targetRef.Group != "core" && targetRef.Group != "") || targetRef.Kind != "Service" {
			continue
		}
		policies := &gatewayapi.BackendTLSPolicyList{}
		if err := r.List(ctx, policies,
			client.InNamespace(policy.Namespace),
			client.MatchingFields{backendTLSPolicyTargetRefIndexKey: string(targetRef.Name)},
		); err != nil {
			return nil, err
		}

		if len(policies.Items) > 1 {
			acceptedCondition = &metav1.Condition{
				Type:    string(gatewayapi.PolicyConditionAccepted),
				Status:  metav1.ConditionFalse,
				Reason:  string(gatewayapi.PolicyReasonConflicted),
				Message: "Multiple BackendTLSPolicies target the same service",
			}
			return acceptedCondition, nil
		}
	}

	var invalidMessages []string
	for _, caCert := range policy.Spec.Validation.CACertificateRefs {
		if (caCert.Group != "core" && caCert.Group != "") || caCert.Kind != "ConfigMap" {
			invalidMessages = append(invalidMessages, "CACertificateRefs must reference ConfigMaps in the core group")
			break
		}
	}
	if len(policy.Spec.Validation.SubjectAltNames) > 0 {
		invalidMessages = append(invalidMessages, "SubjectAltNames feature is not currently supported")
	}
	if policy.Spec.Validation.WellKnownCACertificates != nil {
		invalidMessages = append(invalidMessages, "WellKnownCACertificates feature is not currently supported")
	}
	if len(invalidMessages) > 0 {
		acceptedCondition.Status = metav1.ConditionFalse
		acceptedCondition.Reason = string(gatewayapi.PolicyReasonInvalid)
		acceptedCondition.Message = strings.Join(invalidMessages, " - ")
	}

	return acceptedCondition, nil
}

// list namespaced names of configmaps referred by the gateway.
func listConfigMapNamesReferredByBackendTLSPolicy(policy *gatewayapi.BackendTLSPolicy) map[k8stypes.NamespacedName]struct{} {
	// no need to check group and kind, as if they were different from core/Configmap, the policy would have been marked as invalid.
	nsNames := make(map[k8stypes.NamespacedName]struct{})
	for _, certRef := range policy.Spec.Validation.CACertificateRefs {
		nsName := k8stypes.NamespacedName{
			Namespace: policy.Namespace,
			Name:      string(certRef.Name),
		}
		nsNames[nsName] = struct{}{}
	}
	return nsNames
}
