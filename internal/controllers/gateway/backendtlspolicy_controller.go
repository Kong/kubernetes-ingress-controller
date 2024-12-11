package gateway

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/go-logr/logr"
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers"
	ctrlref "github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/reference"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

// -----------------------------------------------------------------------------
// BackendTLSPolicy Controller - BackendTLSPolicyReconciler
// -----------------------------------------------------------------------------

// BackendTLSPolicyReconciler reconciles a BackendTLSPolicy object.
type BackendTLSPolicyReconciler struct {
	client.Client

	Log               logr.Logger
	DataplaneClient   controllers.DataPlane
	CacheSyncTimeout  time.Duration
	ReferenceIndexers ctrlref.CacheIndexers
	// If GatewayNN is set,
	// only resources managed by the specified Gateway are reconciled.
	GatewayNN controllers.OptionalNamespacedName
}

// SetupWithManager sets up the controller with the Manager.
func (r *BackendTLSPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := setupBackendTLSPolicyIndices(mgr); err != nil {
		return fmt.Errorf("failed to setup indexers: %w", err)
	}
	return ctrl.NewControllerManagedBy(mgr).
		Named("backendtlspolicy-controller").
		WithOptions(controller.Options{
			LogConstructor: func(_ *reconcile.Request) logr.Logger {
				return r.Log
			},
			CacheSyncTimeout: r.CacheSyncTimeout,
		}).
		For(&gatewayapi.BackendTLSPolicy{}).
		Watches(&corev1.Service{}, handler.EnqueueRequestsFromMapFunc(r.listBackendTLSPoliciesForServices)).
		Watches(&gatewayapi.HTTPRoute{}, handler.EnqueueRequestsFromMapFunc(r.listBackendTLSPoliciesForHTTPRoutes)).
		Watches(&gatewayapi.Gateway{}, handler.EnqueueRequestsFromMapFunc(r.listBackendTLSPoliciesForGateways)).
		Complete(r)
}

// -----------------------------------------------------------------------------
// BackendTLSPolicy Controller - Indexers
// -----------------------------------------------------------------------------

const (
	// backendTLSPolicyTargetRefIndexKey is the index key for BackendTLSPolicy objects by their target service reference.
	// The value is the name of the service.
	backendTLSPolicyTargetRefIndexKey = "backendtlspolicy-targetref"
	// httpRouteParentRefIndexKey is the index key for HTTPRoute objects by their parent Gateway reference.
	// The value is the namespace and name of the Gateway.
	httpRouteParentRefIndexKey = "httproute-parentref"
	// httpRouteBackendRefIndexKey is the index key for HTTPRoute objects by their backend service reference.
	// The value is the namespace and name of the Service.
	httpRouteBackendRefIndexKey = "httproute-backendref"
)

// indexBackendTLSPolicyOnTargetRef indexes BackendTLSPolicy objects by their target service reference.
func indexBackendTLSPolicyOnTargetRef(obj client.Object) []string {
	policy, ok := obj.(*gatewayapi.BackendTLSPolicy)
	if !ok {
		return []string{}
	}

	services := []string{}
	for _, targetRef := range policy.Spec.TargetRefs {
		if (targetRef.Group == "" || targetRef.Group == "core") && targetRef.Kind == "Service" {
			services = append(services, string(targetRef.Name))
		}
	}
	return services
}

// indexHTTPRouteOnParentRef indexes HTTPRoute objects by their parent Gateway references.
func indexHTTPRouteOnParentRef(obj client.Object) []string {
	httpRoute, ok := obj.(*gatewayapi.HTTPRoute)
	if !ok {
		return []string{}
	}

	gateways := []string{}
	for _, parentRef := range httpRoute.Spec.ParentRefs {
		// no need to check group and kind nilness, as they have a default value in case not specified
		if *parentRef.Group == gatewayapi.V1Group && *parentRef.Kind == "Gateway" {
			namespace := httpRoute.Namespace
			if parentRef.Namespace != nil {
				namespace = string(*parentRef.Namespace)
			}
			gateways = append(gateways, namespace+"/"+string(parentRef.Name))
		}
	}
	return gateways
}

// indexHTTPRouteOnBackendRef indexes HTTPRoute objects by their backend service references.
func indexHTTPRouteOnBackendRef(obj client.Object) []string {
	httpRoute, ok := obj.(*gatewayapi.HTTPRoute)
	if !ok {
		return []string{}
	}

	services := []string{}
	for _, rule := range httpRoute.Spec.Rules {
		for _, backendRef := range rule.BackendRefs {
			// no need to check group and kind nilness, as they have a default value in case not specified
			if (*backendRef.Group != "core" && *backendRef.Group != "") || *backendRef.Kind != "Service" {
				continue
			}
			namespace := httpRoute.Namespace
			if backendRef.Namespace != nil {
				namespace = string(*backendRef.Namespace)
			}
			services = append(services, namespace+"/"+string(backendRef.Name))
		}
	}
	return services
}

// setupIndexers sets up the indexers for the BackendTLSPolicy controller.
func setupBackendTLSPolicyIndices(mgr ctrl.Manager) error {
	if err := mgr.GetCache().IndexField(
		context.Background(),
		&gatewayapi.BackendTLSPolicy{},
		backendTLSPolicyTargetRefIndexKey,
		indexBackendTLSPolicyOnTargetRef,
	); err != nil {
		return fmt.Errorf("failed to index backendTLSPolicies on service reference: %w", err)
	}

	if err := mgr.GetCache().IndexField(
		context.Background(),
		&gatewayapi.HTTPRoute{},
		httpRouteParentRefIndexKey,
		indexHTTPRouteOnParentRef,
	); err != nil {
		return fmt.Errorf("failed to index httpRoute on ParentRef reference: %w", err)
	}

	if err := mgr.GetCache().IndexField(
		context.Background(),
		&gatewayapi.HTTPRoute{},
		httpRouteBackendRefIndexKey,
		indexHTTPRouteOnBackendRef,
	); err != nil {
		return fmt.Errorf("failed to index httpRoute on Backend reference: %w", err)
	}

	return nil
}

// -----------------------------------------------------------------------------
// BackendTLSPolicy Controller - Event Handlers
// -----------------------------------------------------------------------------

// listBackendTLSPoliciesForServices returns the list of BackendTLSPolicies that targets the given Service.
func (r *BackendTLSPolicyReconciler) listBackendTLSPoliciesForServices(ctx context.Context, obj client.Object) []reconcile.Request {
	service, ok := obj.(*corev1.Service)
	if !ok {
		r.Log.Error(fmt.Errorf("invalid type"), "Found invalid type in event handlers", "expected", "Service", "found", reflect.TypeOf(obj))
		return nil
	}
	policies := &gatewayapi.BackendTLSPolicyList{}
	if err := r.List(ctx, policies,
		client.InNamespace(service.Namespace),
		client.MatchingFields{backendTLSPolicyTargetRefIndexKey: service.Name},
	); err != nil {
		r.Log.Error(err, "Failed to list BackendTLSPolicies for Service", "service", service)
		return nil
	}
	requests := make([]reconcile.Request, 0, len(policies.Items))
	for _, policy := range policies.Items {
		requests = append(requests, reconcile.Request{
			NamespacedName: client.ObjectKeyFromObject(&policy),
		})
	}
	return requests
}

// listBackendTLSPoliciesForHTTPRoutes returns the list of BackendTLSPolicies that targets a service which is used as a backend by
// the given HTTPRoute.
func (r *BackendTLSPolicyReconciler) listBackendTLSPoliciesForHTTPRoutes(ctx context.Context, obj client.Object) []reconcile.Request {
	httpRoute, ok := obj.(*gatewayapi.HTTPRoute)
	if !ok {
		r.Log.Error(fmt.Errorf("invalid type"), "Found invalid type in event handlers", "expected", "HTTPRoute", "found", reflect.TypeOf(obj))
		return nil
	}
	policiesNN, err := r.getBackendTLSPoliciesByHTTPRoute(ctx, *httpRoute)
	if err != nil {
		r.Log.Error(err, "Failed to list BackendTLSPolicies for HTTPRoute", "httpRoute", httpRoute)
		return nil
	}
	return lo.Map(policiesNN, func(policy gatewayapi.BackendTLSPolicy, _ int) reconcile.Request {
		return reconcile.Request{
			NamespacedName: k8stypes.NamespacedName{
				Namespace: policy.Namespace,
				Name:      policy.Name,
			},
		}
	})
}

// listBackendTLSPoliciesForGateways returns the list of BackendTLSPolicies that targets a service which is used as a backend by
// HTTPRoutes connected to the given Gateway.
func (r *BackendTLSPolicyReconciler) listBackendTLSPoliciesForGateways(ctx context.Context, obj client.Object) []reconcile.Request {
	gateway, ok := obj.(*gatewayapi.Gateway)
	if !ok {
		r.Log.Error(fmt.Errorf("invalid type"), "Found invalid type in event handlers", "expected", "Gateway", "found", reflect.TypeOf(obj))
		return nil
	}

	if !r.GatewayNN.Matches(gateway) {
		return nil
	}

	httpRoutes := &gatewayapi.HTTPRouteList{}
	if err := r.List(ctx, httpRoutes,
		client.MatchingFields{httpRouteParentRefIndexKey: gateway.Namespace + "/" + gateway.Name},
	); err != nil {
		r.Log.Error(err, "Failed to list HTTPRoutes for Gateway", "gateway", gateway)
		return nil
	}
	policies := []reconcile.Request{}
	for _, httpRoute := range httpRoutes.Items {
		policiesUsedByHTTPRoute, err := r.getBackendTLSPoliciesByHTTPRoute(ctx, httpRoute)
		if err != nil {
			r.Log.Error(err, "Failed to list BackendTLSPolicies for HTTPRoute", "httpRoute", httpRoute)
			return nil
		}
		policies = append(policies, lo.Map(policiesUsedByHTTPRoute, func(policy gatewayapi.BackendTLSPolicy, _ int) reconcile.Request {
			return reconcile.Request{
				NamespacedName: k8stypes.NamespacedName{
					Namespace: policy.Namespace,
					Name:      policy.Name,
				},
			}
		})...)
	}
	return policies
}

// -----------------------------------------------------------------------------
// BackendTLSPolicy Controller - Reconciliation
// -----------------------------------------------------------------------------

// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=httproutes;gateways;gatewayclasses,verbs=get;list;watch
// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=backendtlspolicies,verbs=get;list;watch;patch;update
// +kubebuilder:rbac:groups="",resources=services;configmaps,verbs=get;list;watch

// Reconcile processes the watched objects.
func (r *BackendTLSPolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("GatewayV1alpha3BackendTLSPolicy", req.NamespacedName)

	backendTLSPolicy := new(gatewayapi.BackendTLSPolicy)
	if err := r.Get(ctx, req.NamespacedName, backendTLSPolicy); err != nil {
		if apierrors.IsNotFound(err) {
			backendTLSPolicy.Namespace = req.Namespace
			backendTLSPolicy.Name = req.Name
			return ctrl.Result{}, r.DataplaneClient.DeleteObject(backendTLSPolicy)
		}
		return ctrl.Result{}, err
	}

	debug(log, backendTLSPolicy, "Processing backendTLSPolicy")

	ancestors, err := r.getBackendTLSPolicyAncestors(ctx, *backendTLSPolicy)
	if err != nil {
		return ctrl.Result{}, err
	}

	if len(ancestors) > 0 {
		acceptedCondition, err := r.validateBackendTLSPolicy(ctx, *backendTLSPolicy)
		if err != nil {
			return ctrl.Result{}, err
		}

		// If the policy is accepted, update the policy in the dataplane.
		if acceptedCondition.Status == metav1.ConditionTrue {
			if err := r.DataplaneClient.UpdateObject(backendTLSPolicy); err != nil {
				return ctrl.Result{}, err
			}

			// Update references to ConfigMaps in the dataplane cache.
			referredConfigMapNames := listConfigMapNamesReferredByBackendTLSPolicy(backendTLSPolicy)
			if err := ctrlref.UpdateReferencesToSecretOrConfigMap(
				ctx,
				r.Client,
				r.ReferenceIndexers,
				r.DataplaneClient,
				backendTLSPolicy,
				referredConfigMapNames,
				&corev1.ConfigMap{}); err != nil {
				if apierrors.IsNotFound(err) {
					return ctrl.Result{Requeue: true}, nil
				}
			}
		} else {
			// In case the policy is not accepted, ensure it gets deleted from the dataplane cache
			if err := r.DataplaneClient.DeleteObject(backendTLSPolicy); err != nil {
				return ctrl.Result{}, err
			}
		}

		if err := r.setPolicyStatus(ctx, *backendTLSPolicy, ancestors, *acceptedCondition); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}
