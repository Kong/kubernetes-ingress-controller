package configuration

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers"
	ctrlref "github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/reference"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/labels"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/logging"
)

// -----------------------------------------------------------------------------
// CoreV1 ConfigMap - Reconciler
// -----------------------------------------------------------------------------

// CoreV1ConfigMapReconciler reconciles ConfigMap resources.
type CoreV1ConfigMapReconciler struct {
	client.Client

	Log              logr.Logger
	Scheme           *runtime.Scheme
	DataplaneClient  controllers.DataPlane
	CacheSyncTimeout time.Duration

	ReferenceIndexers ctrlref.CacheIndexers
	LabelSelector     string
}

var _ controllers.Reconciler = &CoreV1ConfigMapReconciler{}

// SetupWithManager sets up the controller with the Manager.
func (r *CoreV1ConfigMapReconciler) SetupWithManager(mgr ctrl.Manager) error {
	predicateFuncs := predicate.NewPredicateFuncs(r.shouldReconcileConfigMap)
	// we should always try to delete configmaps in caches when they are deleted in cluster.
	predicateFuncs.DeleteFunc = func(_ event.DeleteEvent) bool { return true }

	var (
		labelPredicate predicate.Predicate
		labelSelector  metav1.LabelSelector
		err            error
	)
	if r.LabelSelector != "" {
		labelSelector = metav1.LabelSelector{
			MatchLabels: map[string]string{r.LabelSelector: "true"},
		}
	}

	labelPredicate, err = predicate.LabelSelectorPredicate(labelSelector)
	if err != nil {
		return fmt.Errorf("failed to create secret label selector predicate: %w", err)
	}

	return ctrl.NewControllerManagedBy(mgr).
		Named("CoreV1ConfigMap").
		WithOptions(controller.Options{
			LogConstructor: func(_ *reconcile.Request) logr.Logger {
				return r.Log
			},
			CacheSyncTimeout: r.CacheSyncTimeout,
		}).
		For(&corev1.ConfigMap{},
			builder.WithPredicates(
				predicate.And(
					predicateFuncs,
					labelPredicate,
				)),
		).
		Complete(r)
}

// SetLogger sets the logger.
func (r *CoreV1ConfigMapReconciler) SetLogger(l logr.Logger) {
	r.Log = l
}

// shouldReconcileConfigMap is the filter function to judge whether the ConfigMap should be reconciled
// and stored in cache of the controller. It returns true for the ConfigMap should be reconciled when
// the ConfigMap is referred by objects we care (BackendTLSPolicy).
func (r *CoreV1ConfigMapReconciler) shouldReconcileConfigMap(obj client.Object) bool {
	configMap, ok := obj.(*corev1.ConfigMap)
	if !ok {
		return false
	}

	l := configMap.Labels
	if l != nil {
		if l[CACertLabelKey] == "true" {
			return true
		}

		if _, ok := l[labels.CredentialTypeLabel]; ok {
			return true
		}
	}

	referred, err := r.ReferenceIndexers.ObjectReferred(configMap)
	if err != nil {
		r.Log.Error(err, "Failed to check whether configmap referred",
			"namespace", configMap.Namespace, "name", configMap.Name)
		return false
	}

	return referred
}

// +kubebuilder:rbac:groups="",resources=configmaps,verbs=list;watch

// Reconcile processes the watched objects.
func (r *CoreV1ConfigMapReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("CoreV1ConfigMap", req.NamespacedName)

	// get the relevant object
	configMap := new(corev1.ConfigMap)
	if err := r.Get(ctx, req.NamespacedName, configMap); err != nil {
		if apierrors.IsNotFound(err) {
			configMap.Namespace = req.Namespace
			configMap.Name = req.Name
			return ctrl.Result{}, r.DataplaneClient.DeleteObject(configMap)
		}
		return ctrl.Result{}, err
	}

	log.V(logging.DebugLevel).Info("Reconciling ConfigMap resource", "namespace", req.Namespace, "name", req.Name)

	// clean the object up if it's being deleted
	if !configMap.DeletionTimestamp.IsZero() && time.Now().After(configMap.DeletionTimestamp.Time) {
		log.V(logging.DebugLevel).Info("Resource is being deleted, its configuration will be removed", "type", "ConfigMap", "namespace", req.Namespace, "name", req.Name)
		objectExistsInCache, err := r.DataplaneClient.ObjectExists(configMap)
		if err != nil {
			return ctrl.Result{}, err
		}
		if objectExistsInCache {
			if err := r.DataplaneClient.DeleteObject(configMap); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil // wait until the object is no longer present in the cache
		}
		return ctrl.Result{}, nil
	}

	// update the kong Admin API with the changes
	if err := r.DataplaneClient.UpdateObject(configMap); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}
