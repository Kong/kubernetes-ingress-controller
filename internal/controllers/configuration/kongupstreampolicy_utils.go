package configuration

import (
	"context"
	"fmt"
	"reflect"

	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	gatewaycontroller "github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/gateway"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1beta1"
)

type serviceStatus struct {
	service           corev1.Service
	acceptedCondition metav1.Condition
}

// -----------------------------------------------------------------------------
// KongUpstreamPolicy Controller - Reconciler Helpers
// -----------------------------------------------------------------------------

// enforceKongUpstreamPolicyStatus gets a list of services (ancestors) along with their desired status and enforce them
// in the KongUpstreamPolicy status.
func (r *KongUpstreamPolicyReconciler) enforceKongUpstreamPolicyStatus(ctx context.Context, oldPolicy *kongv1beta1.KongUpstreamPolicy, servicesStatus []serviceStatus) (bool, error) {
	newPolicyStatus := gatewayapi.Policystatus{}
	if len(servicesStatus) > 0 {
		newPolicyStatus.Ancestors = make([]gatewayapi.PolicyAncestorStatus, 0, len(servicesStatus))
	}
	for _, ss := range servicesStatus {
		newPolicyStatus.Ancestors = append(newPolicyStatus.Ancestors,
			gatewayapi.PolicyAncestorStatus{
				AncestorRef: gatewayapi.ParentReference{
					Group:     lo.ToPtr(gatewayapi.Group("core")),
					Kind:      lo.ToPtr(gatewayapi.Kind("Service")),
					Namespace: lo.ToPtr(gatewayapi.Namespace(ss.service.Namespace)),
					Name:      gatewayapi.ObjectName(ss.service.Name),
				},
				ControllerName: gatewaycontroller.GetControllerName(),
				Conditions: []metav1.Condition{
					ss.acceptedCondition,
				},
			},
		)
	}
	if isStatusUpdated := isPolicyStatusUpdated(oldPolicy.Status, newPolicyStatus); !isStatusUpdated {
		newPolicy := oldPolicy.DeepCopy()
		newPolicy.Status = newPolicyStatus
		return true, r.Client.Status().Patch(ctx, newPolicy, client.MergeFrom(oldPolicy))
	}
	return false, nil
}

// buildServicesStatus creates a list of services with their conditions associated.
func (r *KongUpstreamPolicyReconciler) buildServicesStatus(ctx context.Context, services []corev1.Service) ([]serviceStatus, error) {
	// prepare a service mapping to be used in subsequent operations
	mappedServices := make(map[string]serviceStatus)
	for _, service := range services {
		acceptedCondition := metav1.Condition{
			Type:               string(gatewayapi.PolicyConditionAccepted),
			Status:             metav1.ConditionTrue,
			Reason:             string(gatewayapi.PolicyReasonAccepted),
			LastTransitionTime: metav1.Now(),
		}
		mappedServices[buildServiceReference(service.Namespace, service.Name)] = serviceStatus{
			service:           service,
			acceptedCondition: acceptedCondition,
		}
	}

	for _, service := range services {
		httpRoutes := &gatewayapi.HTTPRouteList{}
		err := r.List(ctx, httpRoutes,
			client.MatchingFields{
				routeBackendRefServiceNameIndexKey: buildServiceReference(service.Namespace, service.Name),
			},
		)
		if err != nil {
			return nil, err
		}

		for _, httpRoute := range httpRoutes.Items {
			for _, rule := range httpRoute.Spec.Rules {
				var commonPolicy string
				if len(rule.BackendRefs) == 0 {
					continue
				}
				for i, br := range rule.BackendRefs {
					serviceRef := backendRefToServiceRef(httpRoute.Namespace, br.BackendRef)
					if serviceRef == "" {
						continue
					}
					policy := getPolicyByService(mappedServices[serviceRef].service)
					if _, ok := mappedServices[serviceRef]; !ok {
						continue
					}
					if i == 0 {
						commonPolicy = policy
					}
					if policy != commonPolicy {
						serviceStatus := mappedServices[serviceRef]
						serviceStatus.acceptedCondition.Status = metav1.ConditionFalse
						serviceStatus.acceptedCondition.Reason = string(gatewayapi.PolicyReasonConflicted)
						mappedServices[serviceRef] = serviceStatus
					}
				}
			}
		}
	}

	servicesStatus := make([]serviceStatus, len(mappedServices))
	var i int
	for _, ss := range mappedServices {
		servicesStatus[i] = ss
		i++
	}
	return servicesStatus, nil
}

// -----------------------------------------------------------------------------
// KongUpstreamPolicy Controller - Helpers
// -----------------------------------------------------------------------------

func isPolicyStatusUpdated(oldStatus, newStatus gatewayapi.Policystatus) bool {
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

func backendRefToServiceRef(routeNamespace string, br gatewayapi.BackendRef) string {
	if br.Group != nil && *br.Group != "" && *br.Group != "core" {
		return ""
	}
	if br.Kind != nil && *br.Kind != "" && *br.Kind != "Service" {
		return ""
	}
	namespace := routeNamespace
	if br.Namespace != nil {
		namespace = string(*br.Namespace)
	}
	return buildServiceReference(namespace, string(br.Name))
}

func getPolicyByService(service corev1.Service) string {
	return service.Annotations[kongv1beta1.KongUpstreamPolicyAnnotationKey]
}

func buildServiceReference(namespace, name string) string {
	return fmt.Sprintf("%s/%s", namespace, name)
}
