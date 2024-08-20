package configuration

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	k8stypes "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers"
	ctrlutils "github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/utils"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/logging"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/kubernetes/object/status"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1beta1"
	incubatorv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/incubator/v1alpha1"
)

// -----------------------------------------------------------------------------
// KongUpstreamPolicy Controller - Reconciler
// -----------------------------------------------------------------------------

// KongUpstreamPolicyReconciler reconciles KongUpstreamPolicy resources.
type KongUpstreamPolicyReconciler struct {
	client.Client

	Log              logr.Logger
	Scheme           *runtime.Scheme
	DataplaneClient  controllers.DataPlane
	CacheSyncTimeout time.Duration
	StatusQueue      *status.Queue

	IngressClassName           string
	DisableIngressClassLookups bool

	// KongServiceFacadeEnabled determines whether the controller should populate the KongUpstreamPolicy's ancestor
	// status for KongServiceFacades.
	KongServiceFacadeEnabled bool
	// HTTPRouteEnabled determines whether the controller should populate the KongUpstreamPolicy's
	// ancestor status for Services used in HTTPRoutes.
	HTTPRouteEnabled bool
}

// SetupWithManager sets up the controller with the Manager.
func (r *KongUpstreamPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := r.setupIndices(mgr); err != nil {
		return err
	}

	blder := ctrl.NewControllerManagedBy(mgr).
		Named("KongUpstreamPolicy").
		WithOptions(controller.Options{
			LogConstructor: func(_ *reconcile.Request) logr.Logger {
				return r.Log
			},
			CacheSyncTimeout: r.CacheSyncTimeout,
		}).
		Watches(&corev1.Service{},
			handler.EnqueueRequestsFromMapFunc(r.getUpstreamPolicyForObject),
			builder.WithPredicates(predicate.NewPredicateFuncs(doesObjectReferUpstreamPolicy)),
		).
		Watches(&netv1.Ingress{},
			// Watch for Ingress changes to trigger reconciliation for the KongUpstreamPolicies referenced by the Services
			// used as backend of the Ingress.
			// REVIEW: add predicate here to filter Ingresses not reconciled by current controller?
			handler.EnqueueRequestsFromMapFunc(r.getUpstreamPoliciesForIngressServices),
		)

	if r.HTTPRouteEnabled {
		// Watch for HTTPRoute changes to trigger reconciliation for the KongUpstreamPolicies referenced by the Services
		// of the HTTPRoute.
		blder.Watches(&gatewayapi.HTTPRoute{},
			handler.EnqueueRequestsFromMapFunc(r.getUpstreamPoliciesForHTTPRouteServices),
		)
	}

	if r.KongServiceFacadeEnabled {
		blder.Watches(&incubatorv1alpha1.KongServiceFacade{},
			handler.EnqueueRequestsFromMapFunc(r.getUpstreamPolicyForObject),
			builder.WithPredicates(predicate.NewPredicateFuncs(doesObjectReferUpstreamPolicy)),
		)
	}

	if r.StatusQueue != nil {
		// Watch for notifications on the status queue from Services and KongServiceFacades as their status change
		// needs to be propagated to the KongUpstreamPolicy's ancestor Programmed status.
		blder.WatchesRawSource(
			source.Channel(
				r.StatusQueue.Subscribe(schema.GroupVersionKind{
					Version: "v1",
					Kind:    "Service",
				}),
				handler.EnqueueRequestsFromMapFunc(r.getUpstreamPolicyForObject),
				source.WithPredicates[client.Object, reconcile.Request](predicate.NewPredicateFuncs(doesObjectReferUpstreamPolicy)),
			),
		).
			WatchesRawSource(
				source.Channel(r.StatusQueue.Subscribe(schema.GroupVersionKind{
					Version: incubatorv1alpha1.SchemeGroupVersion.Version,
					Group:   incubatorv1alpha1.SchemeGroupVersion.Group,
					Kind:    incubatorv1alpha1.KongServiceFacadeKind,
				}),
					handler.EnqueueRequestsFromMapFunc(r.getUpstreamPolicyForObject),
					source.WithPredicates[client.Object, reconcile.Request](predicate.NewPredicateFuncs(doesObjectReferUpstreamPolicy)),
				),
			)
	}

	return blder.For(&kongv1beta1.KongUpstreamPolicy{}).
		Complete(r)
}

func (r *KongUpstreamPolicyReconciler) setupIndices(mgr ctrl.Manager) error {
	if err := mgr.GetCache().IndexField(
		context.Background(),
		&corev1.Service{},
		upstreamPolicyIndexKey,
		indexServicesOnUpstreamPolicyAnnotation,
	); err != nil {
		return fmt.Errorf("failed to index services on annotation %s: %w", kongv1beta1.KongUpstreamPolicyAnnotationKey, err)
	}

	if err := mgr.GetCache().IndexField(
		context.Background(),
		&netv1.Ingress{},
		routeBackendRefServiceNameIndexKey,
		indexIngressesOnBackendServiceName,
	); err != nil {
		return fmt.Errorf("failed to index Ingresses on backendServiceName: %w", err)
	}

	if r.HTTPRouteEnabled {
		if err := mgr.GetCache().IndexField(
			context.Background(),
			&gatewayapi.HTTPRoute{},
			routeBackendRefServiceNameIndexKey,
			indexRoutesOnBackendRefServiceName,
		); err != nil {
			return fmt.Errorf("failed to index HTTPRoutes on Services in backendReferences: %w", err)
		}
	}

	if r.KongServiceFacadeEnabled {
		if err := mgr.GetCache().IndexField(
			context.Background(),
			&incubatorv1alpha1.KongServiceFacade{},
			upstreamPolicyIndexKey,
			indexServiceFacadesOnUpstreamPolicyAnnotation,
		); err != nil {
			return fmt.Errorf("failed to index KongServiceFacades on annotation %s: %w", kongv1beta1.KongUpstreamPolicyAnnotationKey, err)
		}
		if err := mgr.GetCache().IndexField(
			context.Background(),
			&gatewayapi.HTTPRoute{},
			routeBackendRefServiceFacadeIndexKey,
			indexRoutesOnBackendRefServiceFacadeName,
		); err != nil {
			return fmt.Errorf("failed to index HTTPRoutes on ServiceFacades in backendReferences: %w", err)
		}
	}

	return nil
}

// -----------------------------------------------------------------------------
// KongUpstreamPolicy Controller - Indexers
// -----------------------------------------------------------------------------

const (
	upstreamPolicyIndexKey               = "upstreamPolicy"
	routeBackendRefServiceNameIndexKey   = "serviceRef"
	routeBackendRefServiceFacadeIndexKey = "serviceFacadeRef"
)

// indexServicesOnUpstreamPolicyAnnotation indexes the services on the annotation konghq.com/upstream-policy.
func indexServicesOnUpstreamPolicyAnnotation(o client.Object) []string {
	service, ok := o.(*corev1.Service)
	if !ok {
		return []string{}
	}
	if service.Annotations != nil {
		if policy, ok := service.Annotations[kongv1beta1.KongUpstreamPolicyAnnotationKey]; ok {
			return []string{policy}
		}
	}
	return []string{}
}

// indexIngressesOnBackendServiceName indexes the Ingresses on the backends.
func indexIngressesOnBackendServiceName(o client.Object) []string {
	ingress, ok := o.(*netv1.Ingress)
	if !ok {
		return []string{}
	}
	var indexes []string
	for _, rule := range ingress.Spec.Rules {
		if rule.HTTP == nil {
			continue
		}
		for _, path := range rule.HTTP.Paths {
			if service := path.Backend.Service; service != nil {
				indexes = append(indexes, string(buildServiceReference(ingress.Namespace, service.Name)))
			}
		}
	}
	return indexes
}

// indexRoutesOnBackendRefServiceName indexes the HTTPRoutes on the backendReferences.
func indexRoutesOnBackendRefServiceName(o client.Object) []string {
	httpRoute, ok := o.(*gatewayapi.HTTPRoute)
	if !ok {
		return []string{}
	}

	var indexes []string
	for _, rule := range httpRoute.Spec.Rules {
		for _, br := range rule.BackendRefs {
			serviceRef := backendRefToServiceRef(httpRoute.Namespace, br.BackendRef)
			if serviceRef == "" {
				continue
			}
			indexes = append(indexes, string(serviceRef))
		}
	}
	return indexes
}

// indexRoutesOnBackendRefServiceFacadeName indexes the HTTPRoutes on the backendReferences using KongServiceFacades as backends.
func indexRoutesOnBackendRefServiceFacadeName(o client.Object) []string {
	httpRoute, ok := o.(*gatewayapi.HTTPRoute)
	if !ok {
		return []string{}
	}

	var indexes []string
	for _, rule := range httpRoute.Spec.Rules {
		for _, br := range rule.BackendRefs {
			serviceRef := backendRefToServiceFacadeRef(httpRoute.Namespace, br.BackendRef)
			if serviceRef == "" {
				continue
			}
			indexes = append(indexes, string(serviceRef))
		}
	}
	return indexes
}

// indexServiceFacadesOnUpstreamPolicyAnnotation indexes the KongServiceFacades on the annotation konghq.com/upstream-policy.
func indexServiceFacadesOnUpstreamPolicyAnnotation(o client.Object) []string {
	service, ok := o.(*incubatorv1alpha1.KongServiceFacade)
	if !ok {
		return []string{}
	}
	if service.Annotations != nil {
		if policy, ok := service.Annotations[kongv1beta1.KongUpstreamPolicyAnnotationKey]; ok {
			return []string{policy}
		}
	}
	return []string{}
}

// -----------------------------------------------------------------------------
// KongUpstreamPolicy Controller - Watch Predicates
// -----------------------------------------------------------------------------

// getUpstreamPolicyForObject enqueues a new reconcile request for the KongUpstreamPolicy referenced by an object.
func (r *KongUpstreamPolicyReconciler) getUpstreamPolicyForObject(ctx context.Context, obj client.Object) []reconcile.Request {
	annotations := obj.GetAnnotations()
	if annotations == nil {
		return nil
	}
	policyName, ok := annotations[kongv1beta1.KongUpstreamPolicyAnnotationKey]
	if !ok {
		return nil
	}

	kongUpstreamPolicy := &kongv1beta1.KongUpstreamPolicy{}
	if err := r.Get(ctx, k8stypes.NamespacedName{
		Namespace: obj.GetNamespace(),
		Name:      policyName,
	}, kongUpstreamPolicy); err != nil {
		if !apierrors.IsNotFound(err) {
			r.Log.Error(err, "Failed to retrieve KongUpstreamPolicy in watch predicates",
				"KongUpstreamPolicy", fmt.Sprintf("%s/%s", obj.GetNamespace(), policyName),
			)
		}
		return []reconcile.Request{}
	}

	return []reconcile.Request{
		{
			NamespacedName: k8stypes.NamespacedName{
				Namespace: obj.GetNamespace(),
				Name:      policyName,
			},
		},
	}
}

// getUpstreamPoliciesForIngressServices enqueues a new reconcile request for the KongUpstreamPolicies referenced by
// the Services of an Ingress.
func (r *KongUpstreamPolicyReconciler) getUpstreamPoliciesForIngressServices(ctx context.Context, obj client.Object) []reconcile.Request {
	ingress, ok := obj.(*netv1.Ingress)
	if !ok {
		return nil
	}
	var requests []reconcile.Request
	for _, rule := range ingress.Spec.Rules {
		for _, path := range rule.HTTP.Paths {
			if path.Backend.Service == nil {
				continue
			}
			service := &corev1.Service{}
			if err := r.Client.Get(ctx, k8stypes.NamespacedName{
				Namespace: ingress.Namespace,
				Name:      path.Backend.Service.Name,
			}, service); err != nil {
				if !apierrors.IsNotFound(err) {
					r.Log.Error(err, "Failed to retrieve Service in watch predicates",
						"Service", fmt.Sprintf("%s/%s", ingress.Namespace, path.Backend.Service.Name),
					)
				}
				continue
			}

			if service.Annotations == nil {
				continue
			}
			upstreamPolicy, ok := service.Annotations[kongv1beta1.KongUpstreamPolicyAnnotationKey]
			if !ok {
				continue
			}
			requests = append(requests, reconcile.Request{
				NamespacedName: k8stypes.NamespacedName{
					Namespace: ingress.Namespace,
					Name:      upstreamPolicy,
				},
			})
		}
	}
	return requests
}

// getUpstreamPoliciesForHTTPRouteServices enqueues a new reconcile request for the KongUpstreamPolicies referenced by
// the Services of an HTTPRoute.
func (r *KongUpstreamPolicyReconciler) getUpstreamPoliciesForHTTPRouteServices(ctx context.Context, obj client.Object) []reconcile.Request {
	httpRoute, ok := obj.(*gatewayapi.HTTPRoute)
	if !ok {
		return nil
	}
	var requests []reconcile.Request
	for _, rule := range httpRoute.Spec.Rules {
		for _, br := range rule.BackendRefs {
			if !isSupportedHTTPRouteBackendRef(br.BackendRef) {
				continue
			}

			namespace := httpRoute.Namespace
			if br.BackendRef.Namespace != nil {
				namespace = string(*br.BackendRef.Namespace)
			}
			service := &corev1.Service{}
			if err := r.Client.Get(ctx, k8stypes.NamespacedName{
				Namespace: namespace,
				Name:      string(br.BackendRef.Name),
			}, service); err != nil {
				if !apierrors.IsNotFound(err) {
					r.Log.Error(err, "Failed to retrieve Service in watch predicates",
						"Service", fmt.Sprintf("%s/%s", namespace, string(br.BackendRef.Name)),
					)
				}
				continue
			}

			if service.Annotations == nil {
				continue
			}
			upstreamPolicy, ok := service.Annotations[kongv1beta1.KongUpstreamPolicyAnnotationKey]
			if !ok {
				continue
			}
			requests = append(requests, reconcile.Request{
				NamespacedName: k8stypes.NamespacedName{
					Namespace: httpRoute.Namespace,
					Name:      upstreamPolicy,
				},
			})
		}
	}
	return requests
}

// doesObjectReferUpstreamPolicy filters out all the objects not referencing KongUpstreamPolicies.
func doesObjectReferUpstreamPolicy(obj client.Object) bool {
	annotations := obj.GetAnnotations()
	if annotations == nil {
		return false
	}
	_, ok := annotations[kongv1beta1.KongUpstreamPolicyAnnotationKey]
	return ok
}

// -----------------------------------------------------------------------------
// KongUpstreamPolicy Controller - Reconciliation
// -----------------------------------------------------------------------------

// +kubebuilder:rbac:groups=configuration.konghq.com,resources=kongupstreampolicies,verbs=get;list;watch
// +kubebuilder:rbac:groups=configuration.konghq.com,resources=kongupstreampolicies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=httproutes,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch
// +kubebuilder:rbac:groups=incubator.ingress-controller.konghq.com,resources=kongservicefacades,verbs=get;list;watch

// Reconcile processes the watched objects.
func (r *KongUpstreamPolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("KongV1beta1KongUpstreamPolicy", req.NamespacedName)

	// get the relevant object
	kongUpstreamPolicy := new(kongv1beta1.KongUpstreamPolicy)

	if err := r.Get(ctx, req.NamespacedName, kongUpstreamPolicy); err != nil {
		if apierrors.IsNotFound(err) {
			kongUpstreamPolicy.Namespace = req.Namespace
			kongUpstreamPolicy.Name = req.Name

			return ctrl.Result{}, r.DataplaneClient.DeleteObject(kongUpstreamPolicy)
		}
		return ctrl.Result{}, err
	}
	log.V(logging.DebugLevel).Info("Reconciling resource", "namespace", req.Namespace, "name", req.Name)

	// clean the object up if it's being deleted
	if !kongUpstreamPolicy.DeletionTimestamp.IsZero() && time.Now().After(kongUpstreamPolicy.DeletionTimestamp.Time) {
		log.V(logging.DebugLevel).Info("Resource is being deleted, its configuration will be removed", "type", "KongUpstreamPolicy", "namespace", req.Namespace, "name", req.Name)

		objectExistsInCache, err := r.DataplaneClient.ObjectExists(kongUpstreamPolicy)
		if err != nil {
			return ctrl.Result{}, err
		}
		if objectExistsInCache {
			if err := r.DataplaneClient.DeleteObject(kongUpstreamPolicy); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil // wait until the object is no longer present in the cache
		}
		return ctrl.Result{}, nil
	}

	// Do not store the KongUpstreamPolicy into the cache if it is not used by some Service being the backend of Ingress/HTTPRoute
	// that matches the ingress class/gateway class of the controller.
	class := new(netv1.IngressClass)
	if !r.DisableIngressClassLookups {
		if err := r.Get(ctx, k8stypes.NamespacedName{Name: r.IngressClassName}, class); err != nil {
			// we log this without taking action to support legacy configurations that only set ingressClassName or
			// used the class annotation and did not create a corresponding IngressClass. We only need this to determine
			// if the IngressClass is default or to configure default settings, and can assume no/no additional defaults
			// if none exists.
			log.V(logging.DebugLevel).Info("Could not retrieve IngressClass", "ingressclass", r.IngressClassName)
		}
	}
	shouldStore, err := r.upstreamPolicyUsedByBackendsOfMatchingClass(ctx, req.NamespacedName, ctrlutils.IsDefaultIngressClass(class))
	if err != nil {
		return ctrl.Result{}, err
	}
	if !shouldStore {
		log.V(logging.DebugLevel).Info("KongUpstreamPolicy is not referenced by Services used as backend of Ingress or HTTPRoute with matching class, skipping reconciliation")
		return ctrl.Result{}, nil
	}

	// enforce the desired KongUpstreamPolicy status
	updated, err := r.enforceKongUpstreamPolicyStatus(ctx, kongUpstreamPolicy)
	if err != nil {
		return ctrl.Result{}, err
	}
	if updated {
		// status update will re-trigger reconciliation
		return ctrl.Result{}, nil
	}

	// update the kong Admin API with the changes
	if err := r.DataplaneClient.UpdateObject(kongUpstreamPolicy); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetLogger sets the logger.
func (r *KongUpstreamPolicyReconciler) SetLogger(l logr.Logger) {
	r.Log = l
}
