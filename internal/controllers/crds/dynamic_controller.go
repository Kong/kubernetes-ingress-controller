package crds

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"github.com/samber/lo"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/utils"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// +kubebuilder:rbac:groups="apiextensions.k8s.io",resources=customresourcedefinitions,verbs=list;watch

type Controller interface {
	SetupWithManager(mgr ctrl.Manager) error
}

// DynamicCRDController ensures that RequiredCRDs are installed in the cluster and only then sets up all of its Controllers
// that depends on them.
// In case the CRDs are not installed at start-up time, DynamicCRDController will set up a watch for CustomResourceDefinition
// and will dynamically set up its Controllers once it detects that all RequiredCRDs are already in place.
type DynamicCRDController struct {
	Log              logr.Logger
	Manager          ctrl.Manager
	CacheSyncTimeout time.Duration
	Controllers      []Controller
	RequiredCRDs     []schema.GroupVersionResource

	startControllersOnce sync.Once
}

func (r *DynamicCRDController) SetupWithManager(mgr ctrl.Manager) error {
	if r.allRequiredCRDsInstalled() {
		r.Log.V(util.DebugLevel).Info("All required CustomResourceDefinitions are installed, skipping DynamicCRDController set up")
		return r.setupControllers(mgr)
	}

	r.Log.Info("Required CustomResourceDefinitions are not installed, setting up a watch for them in case they are installed afterward")

	c, err := controller.New("DynamicCRDController", mgr, controller.Options{
		Reconciler: r,
		LogConstructor: func(_ *reconcile.Request) logr.Logger {
			return r.Log
		},
		CacheSyncTimeout: r.CacheSyncTimeout,
	})
	if err != nil {
		return err
	}

	return c.Watch(
		&source.Kind{Type: &apiextensionsv1.CustomResourceDefinition{}},
		&handler.EnqueueRequestForObject{},
		predicate.NewPredicateFuncs(r.isOneOfRequiredCRDs),
	)
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

	var startControllersErr error
	r.startControllersOnce.Do(func() {
		log.V(util.InfoLevel).Info("All required CustomResourceDefinitions are installed, setting up the controllers")
		startControllersErr = r.setupControllers(r.Manager)
	})
	if startControllersErr != nil {
		return ctrl.Result{}, startControllersErr
	}

	return ctrl.Result{}, nil
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

func (r *DynamicCRDController) setupControllers(mgr ctrl.Manager) error {
	errs := lo.FilterMap(r.Controllers, func(c Controller, _ int) (error, bool) {
		if err := c.SetupWithManager(mgr); err != nil {
			return err, true
		}
		return nil, false
	})

	return errors.Join(errs...)
}
