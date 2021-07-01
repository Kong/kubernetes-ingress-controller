package configuration

import (
	"context"

	"github.com/go-logr/logr"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/kong/kubernetes-ingress-controller/railgun/internal/proxy"
)

// -----------------------------------------------------------------------------
// APIExtensions CustomResourceDefinition
// -----------------------------------------------------------------------------

// CustomResourceDefinition reconciles a Ingress object
type CustomResourceDefinitionReconciler struct {
	client.Client
	Mgr              manager.Manager
	IngressClassName string

	Log    logr.Logger
	Scheme *runtime.Scheme
	Proxy  proxy.Proxy
}

// SetupWithManager sets up the controller with the Manager.
func (r *CustomResourceDefinitionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&apiextv1beta1.CustomResourceDefinition{}).Complete(r)
}

//+kubebuilder:rbac:groups="",resources=customresourcedefinitions,verbs=get;list;watch

// Reconcile processes the watched objects
func (r *CustomResourceDefinitionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("CustomResourceDefinition", req.NamespacedName)

	// get the relevant object
	obj := new(apiextv1beta1.CustomResourceDefinition)
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		log.Error(err, "object was queued for reconcilation but could not be retrieved", "namespace", req.Namespace, "name", req.Name)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// don't act on any CRDs until they have been established (i.e. they are usable)
	established := false
	for _, condition := range obj.Status.Conditions {
		established = condition.Type == apiextv1beta1.Established && condition.Status == apiextv1beta1.ConditionTrue
	}
	if !established {
		return ctrl.Result{Requeue: true}, nil
	}

	// for some 3rd party CRDs we support (e.g. knative.Ingress) we will dynamically load a controller
	// for those types once the CRD becomes present in the API.
	loaded, err := r.dynamicallyLoadControllerForCRD(obj)
	if err != nil {
		log.Error(err, "supported CRD found but could not load the controller", "name", obj.Name)
		return ctrl.Result{}, err
	}

	if loaded {
		log.Info("found supported CRD, controller for this supported API has now beed started")
	}

	return ctrl.Result{}, nil
}

// -----------------------------------------------------------------------------
// Dynamically Loaded Controllers
// -----------------------------------------------------------------------------

const (
	supportedKnativeCRD            = "ingresses.networking.internal.knative.dev"
	supportedKnativeIngressVersion = "v1alpha1"
)

// dynamicallyLoadControllerForCRD accepts a provided CRD and controller manager and if
// that CRD is a supported 3rd party type which we provide a controller for, it will ensure
// that controller is loaded into the manager now that the API has become available in Kubernetes.
// if this returns false, nil then it was a no-op and the CRD is not supported.
func (r *CustomResourceDefinitionReconciler) dynamicallyLoadControllerForCRD(crd *apiextv1beta1.CustomResourceDefinition) (loaded bool, err error) {
	if crd.Name == supportedKnativeCRD {
		for _, version := range crd.Spec.Versions {
			if version.Name == supportedKnativeIngressVersion {
				knative := Knativev1alpha1IngressReconciler{
					Client:           r.Mgr.GetClient(),
					Log:              ctrl.Log.WithName("controllers").WithName("Ingress").WithName("KnativeV1Alpha1"),
					Scheme:           r.Mgr.GetScheme(),
					IngressClassName: r.IngressClassName,
					Proxy:            r.Proxy,
				}
				err = knative.SetupWithManager(r.Mgr)
				loaded = true
			}
		}
	}
	return
}
