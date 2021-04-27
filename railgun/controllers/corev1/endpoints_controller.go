package corev1

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/pkg/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/railgun/internal/ctrlutils"
	"github.com/kong/kubernetes-ingress-controller/railgun/internal/mgrutils"
)

// -----------------------------------------------------------------------------
// CoreV1 Endpoints
// -----------------------------------------------------------------------------

// CoreV1Endpoints reconciles a Endpoint object
type CoreV1EndpointsReconciler struct {
	client.Client

	Log        logr.Logger
	Scheme     *runtime.Scheme
	KongConfig sendconfig.Kong
}

// SetupWithManager sets up the controller with the Manager.
func (r *CoreV1EndpointsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// TODO: this is too broad, we need to sweep for Endpoints referred to by Services we support.
	//       See: https://github.com/Kong/kubernetes-ingress-controller/issues/1259
	return ctrl.NewControllerManagedBy(mgr).For(&v1.Endpoints{}).Complete(r)
}

//+kubebuilder:rbac:groups=v1,resources=endpoints,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=v1,resources=endpoints/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=v1,resources=endpoints/finalizers,verbs=update

// Reconcile processes the watched objects
func (r *CoreV1EndpointsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("CoreV1Endpoint", req.NamespacedName)

	// TODO: this just reduces log noise for now, we can clean this up when we clean up other TODO items
	//       See: https://github.com/Kong/kubernetes-ingress-controller/issues/1259
	if req.Namespace == "local-path-storage" || req.Namespace == "kong-system" || req.Namespace == "kube-system" || req.Name == "kubernetes" {
		return ctrl.Result{}, nil
	}

	// get the relevant object
	obj := new(v1.Endpoints)
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log.Info("reconciling resource", "namespace", req.Namespace, "name", req.Name)

	// clean the object up if it's being deleted
	if !obj.DeletionTimestamp.IsZero() && time.Now().After(obj.DeletionTimestamp.Time) {
		log.Info("resource is being deleted, its configuration will be removed", "type", "Endpoint", "namespace", req.Namespace, "name", req.Name)
		if err := mgrutils.CacheStores.Endpoint.Delete(obj); err != nil {
			return ctrl.Result{}, err
		}
		return ctrlutils.CleanupFinalizer(ctx, r.Client, log, req.NamespacedName, obj)
	}

	// before we store cache data for this object, ensure that it has our finalizer set
	if !ctrlutils.HasFinalizer(obj, ctrlutils.KongIngressFinalizer) {
		log.Info("finalizer is not set for ingress object, setting it", req.Namespace, req.Name)
		finalizers := obj.GetFinalizers()
		obj.SetFinalizers(append(finalizers, ctrlutils.KongIngressFinalizer))
		if err := r.Client.Update(ctx, obj); err != nil { // TODO: patch here instead of update
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// cache the new object
	if err := mgrutils.CacheStores.Endpoint.Add(obj); err != nil {
		return ctrl.Result{}, err
	}

	// update the kong Admin API with the changes
	return ctrl.Result{}, ctrlutils.UpdateKongAdmin(ctx, &r.KongConfig)
}
