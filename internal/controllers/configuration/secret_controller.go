package configuration

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers"
	ctrlref "github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/reference"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/labels"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

// -----------------------------------------------------------------------------
// CoreV1 Secret - Reconciler
// -----------------------------------------------------------------------------

const (
	CACertLabelKey = "konghq.com/ca-cert"
)

// CoreV1SecretReconciler reconciles Secret resources.
type CoreV1SecretReconciler struct {
	client.Client

	Log              logr.Logger
	Scheme           *runtime.Scheme
	DataplaneClient  controllers.DataPlane
	CacheSyncTimeout time.Duration

	ReferenceIndexers ctrlref.CacheIndexers
}

var _ controllers.Reconciler = &CoreV1SecretReconciler{}

// SetupWithManager sets up the controller with the Manager.
func (r *CoreV1SecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	predicateFuncs := predicate.NewPredicateFuncs(r.shouldReconcileSecret)
	// we should always try to delete secrets in caches when they are deleted in cluster.
	predicateFuncs.DeleteFunc = func(_ event.DeleteEvent) bool { return true }

	return ctrl.NewControllerManagedBy(mgr).
		Named("CoreV1Secret").
		WithOptions(controller.Options{
			LogConstructor: func(_ *reconcile.Request) logr.Logger {
				return r.Log
			},
			CacheSyncTimeout: r.CacheSyncTimeout,
		}).
		Watches(&corev1.Secret{},
			&handler.EnqueueRequestForObject{},
			builder.WithPredicates(predicateFuncs),
		).
		Complete(r)
}

// SetLogger sets the logger.
func (r *CoreV1SecretReconciler) SetLogger(l logr.Logger) {
	r.Log = l
}

// shouldReconcileSecret is the filter function to judge whether the secret should be reconciled
// and stored in cache of the controller. It returns true for the secret should be reconciled when:
// - the secret has label: konghq.com/ca-cert:true
// - or the secret is referred by objects we care (service, ingress, gateway, ...)
func (r *CoreV1SecretReconciler) shouldReconcileSecret(obj client.Object) bool {
	secret, ok := obj.(*corev1.Secret)
	if !ok {
		return false
	}

	l := secret.Labels
	if l != nil {
		if l[CACertLabelKey] == "true" {
			return true
		}

		if _, ok := l[labels.CredentialTypeLabel]; ok {
			return true
		}
	}

	referred, err := r.ReferenceIndexers.ObjectReferred(secret)
	if err != nil {
		r.Log.Error(err, "Failed to check whether secret referred",
			"namespace", secret.Namespace, "name", secret.Name)
		return false
	}

	return referred
}

// +kubebuilder:rbac:groups="",resources=secrets,verbs=list;watch

// Reconcile processes the watched objects.
func (r *CoreV1SecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("CoreV1Secret", req.NamespacedName)

	// get the relevant object
	secret := new(corev1.Secret)
	if err := r.Get(ctx, req.NamespacedName, secret); err != nil {
		if apierrors.IsNotFound(err) {
			secret.Namespace = req.Namespace
			secret.Name = req.Name
			return ctrl.Result{}, r.DataplaneClient.DeleteObject(secret)
		}
		return ctrl.Result{}, err
	}

	log.V(util.DebugLevel).Info("Reconciling resource", "namespace", req.Namespace, "name", req.Name)

	// clean the object up if it's being deleted
	if !secret.DeletionTimestamp.IsZero() && time.Now().After(secret.DeletionTimestamp.Time) {
		log.V(util.DebugLevel).Info("Resource is being deleted, its configuration will be removed", "type", "Secret", "namespace", req.Namespace, "name", req.Name)
		objectExistsInCache, err := r.DataplaneClient.ObjectExists(secret)
		if err != nil {
			return ctrl.Result{}, err
		}
		if objectExistsInCache {
			if err := r.DataplaneClient.DeleteObject(secret); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil // wait until the object is no longer present in the cache
		}
		return ctrl.Result{}, nil
	}

	// update the kong Admin API with the changes
	if err := r.DataplaneClient.UpdateObject(secret); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}
