package configuration

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/kong/go-kong/kong"
	"github.com/samber/mo"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/cache"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/kubernetes/object/status"
	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1alpha1"
)

// -----------------------------------------------------------------------------
// KongV1Alpha1 KongLicense - Reconciler
// -----------------------------------------------------------------------------

// KongV1Alpha1KongLicenseReconciler reconciles KongLicense resources.
type KongV1Alpha1KongLicenseReconciler struct {
	client.Client

	Log              logr.Logger
	Scheme           *runtime.Scheme
	LicenseCache     cache.Store
	CacheSyncTimeout time.Duration
	StatusQueue      *status.Queue

	controllerName string
}

var _ controllers.Reconciler = &KongV1Alpha1KongLicenseReconciler{}

func NewLicenseCache() cache.Store {
	return cache.NewStore(kongLicenseKeyFunc)
}

func kongLicenseKeyFunc(obj interface{}) (string, error) {
	l, ok := obj.(*kongv1alpha1.KongLicense)
	if !ok {
		return "", fmt.Errorf("object is type %T, not KongLicense", obj)
	}
	return l.Name, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KongV1Alpha1KongLicenseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("KongV1Alpha1KongLicense", mgr, controller.Options{
		Reconciler: r,
		LogConstructor: func(_ *reconcile.Request) logr.Logger {
			return r.Log
		},
		CacheSyncTimeout: r.CacheSyncTimeout,
	})
	if err != nil {
		return err
	}
	// if configured, start the status updater controller
	if r.StatusQueue != nil {
		if err := c.Watch(
			&source.Channel{Source: r.StatusQueue.Subscribe(schema.GroupVersionKind{
				Group:   "configuration.konghq.com",
				Version: "v1alpha1",
				Kind:    "KongLicense",
			})},
			&handler.EnqueueRequestForObject{},
		); err != nil {
			return err
		}
	}
	return c.Watch(
		source.Kind(mgr.GetCache(), &kongv1alpha1.KongLicense{}),
		&handler.EnqueueRequestForObject{},
		predicate.NewPredicateFuncs(isKongLicenseEnabled),
	)
}

// SetLogger sets the logger.
func (r *KongV1Alpha1KongLicenseReconciler) SetLogger(l logr.Logger) {
	r.Log = l
}

//+kubebuilder:rbac:groups=configuration.konghq.com,resources=konglicenses,verbs=get;list;watch
//+kubebuilder:rbac:groups=configuration.konghq.com,resources=konglicenses/status,verbs=get;update;patch

// Reconcile processes the watched objects.
func (r *KongV1Alpha1KongLicenseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("KongV1Alpha1KongLicense", req.NamespacedName)

	// get the relevant object
	obj := new(kongv1alpha1.KongLicense)

	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		if apierrors.IsNotFound(err) {
			obj.Namespace = req.Namespace
			obj.Name = req.Name

			return ctrl.Result{}, r.LicenseCache.Delete(obj)
		}
		return ctrl.Result{}, err
	}
	log.V(util.DebugLevel).Info("Reconciling resource", "namespace", req.Namespace, "name", req.Name)

	// clean the object up if it's being deleted
	if !obj.DeletionTimestamp.IsZero() && time.Now().After(obj.DeletionTimestamp.Time) {
		log.V(util.DebugLevel).Info("Resource is being deleted, its configuration will be removed", "type", "KongLicense", "namespace", req.Namespace, "name", req.Name)

		_, objectExistsInCache, err := r.LicenseCache.Get(obj)
		if err != nil {
			return ctrl.Result{}, err
		}
		if objectExistsInCache {
			if err := r.LicenseCache.Delete(obj); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil // wait until the object is no longer present in the cache
		}
		return ctrl.Result{}, nil
	}

	// update the kong Admin API with the changes
	if err := r.LicenseCache.Add(obj); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func isKongLicenseEnabled(obj client.Object) bool {
	kongLicense, ok := obj.(*kongv1alpha1.KongLicense)
	if !ok {
		return false
	}
	return kongLicense.Enabled
}

func compareKongLicense(license1, license2 *kongv1alpha1.KongLicense) bool {
	if license1.CreationTimestamp.Before(&license2.CreationTimestamp) {
		return true
	}
	if license2.CreationTimestamp.Before(&license1.CreationTimestamp) {
		return false
	}
	return license1.Name < license2.Name
}

func (r *KongV1Alpha1KongLicenseReconciler) GetLicense() mo.Option[kong.License] {
	licenseList := r.LicenseCache.List()
	var chosenLicense *kongv1alpha1.KongLicense
	for _, obj := range licenseList {
		license, ok := obj.(*kongv1alpha1.KongLicense)
		if !ok {
			continue
		}
		if chosenLicense == nil || compareKongLicense(license, chosenLicense) {
			chosenLicense = license
		}
	}
	if chosenLicense == nil {
		r.Log.V(util.DebugLevel).Info("No available KongLicenses found in cluster")
		return mo.None[kong.License]()
	}
	// convert chosen KongLicense to kong.License
	// TODO: validate the license against Kong gateway.
	r.Log.V(util.DebugLevel).Info("Get license from KongLicense resource", "name", chosenLicense.Name)
	return mo.Some(kong.License{
		ID:      kong.String(uuid.NewSHA1(uuid.Nil, []byte("KongLicense:"+chosenLicense.Name)).String()),
		Payload: kong.String(chosenLicense.RawLicenseString),
	})
}
