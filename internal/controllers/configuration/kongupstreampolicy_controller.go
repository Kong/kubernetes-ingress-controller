package configuration

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	k8stypes "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
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

	// KongServiceFacadeEnabled determines whether the controller should populate the KongUpstreamPolicy's ancestor
	// status for KongServiceFacades.
	KongServiceFacadeEnabled bool
}

// SetupWithManager sets up the controller with the Manager.
func (r *KongUpstreamPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("KongUpstreamPolicy", mgr, controller.Options{
		Reconciler: r,
		LogConstructor: func(_ *reconcile.Request) logr.Logger {
			return r.Log
		},
		CacheSyncTimeout: r.CacheSyncTimeout,
	})
	if err != nil {
		return err
	}

	if err := r.setupIndices(mgr); err != nil {
		return err
	}

	if err := c.Watch(
		source.Kind(mgr.GetCache(), &corev1.Service{}),
		handler.EnqueueRequestsFromMapFunc(r.getUpstreamPolicyForObject),
		predicate.NewPredicateFuncs(doesObjectReferUpstreamPolicy),
	); err != nil {
		return err
	}

	if r.KongServiceFacadeEnabled {
		if err := c.Watch(
			source.Kind(mgr.GetCache(), &incubatorv1alpha1.KongServiceFacade{}),
			handler.EnqueueRequestsFromMapFunc(r.getUpstreamPolicyForObject),
			predicate.NewPredicateFuncs(doesObjectReferUpstreamPolicy),
		); err != nil {
			return err
		}
	}

	return c.Watch(
		source.Kind(mgr.GetCache(), &kongv1beta1.KongUpstreamPolicy{}),
		&handler.EnqueueRequestForObject{},
	)
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
		&gatewayapi.HTTPRoute{},
		routeBackendRefServiceNameIndexKey,
		indexRoutesOnBackendRefServiceName,
	); err != nil {
		return fmt.Errorf("failed to index HTTPRoutes on backendReferences: %w", err)
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
	}

	return nil
}

// -----------------------------------------------------------------------------
// KongUpstreamPolicy Controller - Indexers
// -----------------------------------------------------------------------------

const (
	upstreamPolicyIndexKey             = "upstreamPolicy"
	routeBackendRefServiceNameIndexKey = "serviceRef"
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
	log.V(util.DebugLevel).Info("Reconciling resource", "namespace", req.Namespace, "name", req.Name)

	// clean the object up if it's being deleted
	if !kongUpstreamPolicy.DeletionTimestamp.IsZero() && time.Now().After(kongUpstreamPolicy.DeletionTimestamp.Time) {
		log.V(util.DebugLevel).Info("Resource is being deleted, its configuration will be removed", "type", "KongUpstreamPolicy", "namespace", req.Namespace, "name", req.Name)

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
