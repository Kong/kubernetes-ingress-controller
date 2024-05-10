package crds

import (
	"context"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"github.com/samber/lo"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/utils"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

// +kubebuilder:rbac:groups="apiextensions.k8s.io",resources=customresourcedefinitions,verbs=list;watch

type Controller interface {
	SetupWithManager(mgr ctrl.Manager) error
}

// DynamicCRDController ensures that RequiredCRDs are installed in the cluster and only then sets up its Controller
// that depends on them.
// In case the CRDs are not installed at start-up time, DynamicCRDController will set up a watch for CustomResourceDefinition
// and will dynamically set up its Controller once it detects that all RequiredCRDs are already in place.
type DynamicCRDController struct {
	Log              logr.Logger
	Manager          ctrl.Manager
	CacheSyncTimeout time.Duration
	Controller       Controller
	RequiredCRDs     []schema.GroupVersionResource

	// startControllerOnce ensures that the controller is started only once.
	startControllerOnce sync.Once
}

func (r *DynamicCRDController) SetupWithManager(mgr ctrl.Manager) error {
	if r.allRequiredCRDsInstalled() {
		r.Log.V(util.DebugLevel).Info("All required CustomResourceDefinitions are installed, skipping DynamicCRDController set up")
		return r.setupController(mgr)
	}

	r.Log.Info("Required CustomResourceDefinitions are not installed, setting up a watch for them in case they are installed afterward")

	return ctrl.NewControllerManagedBy(mgr).
		Named("DynamicCRDController").
		WithOptions(controller.Options{
			LogConstructor: func(_ *reconcile.Request) logr.Logger {
				return r.Log
			},
			CacheSyncTimeout: r.CacheSyncTimeout,
		}).
		Watches(&apiextensionsv1.CustomResourceDefinition{},
			&handler.EnqueueRequestForObject{},
			builder.WithPredicates(predicate.NewPredicateFuncs(r.isOneOfRequiredCRDs)),
		).
		Complete(r)
}

func (r *DynamicCRDController) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("CustomResourceDefinition", req.NamespacedName)

	crd := new(apiextensionsv1.CustomResourceDefinition)
	if err := r.Manager.GetClient().Get(ctx, req.NamespacedName, crd); err != nil {
		if apierrors.IsNotFound(err) {
			log.V(util.DebugLevel).Info("Object enqueued no longer exists, skipping", "name", req.Name)
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	log.V(util.DebugLevel).Info("Processing CustomResourceDefinition", "name", req.Name)

	if !r.allRequiredCRDsInstalled() {
		log.V(util.DebugLevel).Info("Still not all required CustomResourceDefinitions are installed, waiting")
		return ctrl.Result{}, nil
	}

	var startControllerErr error
	r.startControllerOnce.Do(func() {
		log.V(util.InfoLevel).Info("All required CustomResourceDefinitions are installed, setting up the controller")
		startControllerErr = r.setupController(r.Manager)
	})
	if startControllerErr != nil {
		return ctrl.Result{}, startControllerErr
	}

	return ctrl.Result{}, nil
}

func (r *DynamicCRDController) SetLogger(logger logr.Logger) {
	r.Log = logger
}

func (r *DynamicCRDController) allRequiredCRDsInstalled() bool {
	return lo.EveryBy(r.RequiredCRDs, func(gvr schema.GroupVersionResource) bool {
		return utils.CRDExists(r.Manager.GetClient().RESTMapper(), gvr)
	})
}

func (r *DynamicCRDController) isOneOfRequiredCRDs(obj client.Object) bool {
	crd, ok := obj.(*apiextensionsv1.CustomResourceDefinition)
	if !ok {
		return false
	}

	return lo.ContainsBy(r.RequiredCRDs, func(gvr schema.GroupVersionResource) bool {
		versionMatches := lo.ContainsBy(crd.Spec.Versions, func(crdv apiextensionsv1.CustomResourceDefinitionVersion) bool {
			return crdv.Name == gvr.Version
		})

		return crd.Spec.Group == gvr.Group &&
			crd.Status.AcceptedNames.Plural == gvr.Resource &&
			versionMatches
	})
}

func (r *DynamicCRDController) setupController(mgr ctrl.Manager) error {
	return r.Controller.SetupWithManager(mgr)
}
