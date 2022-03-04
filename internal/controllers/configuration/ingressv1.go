package configuration

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/go-logr/logr"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	ctrlutils "github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/utils"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/kubernetes/object/status"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// -----------------------------------------------------------------------------
// NetV1 Ingress - Reconciler
// -----------------------------------------------------------------------------

// NetV1IngressReconciler reconciles Ingress resources
type NetV1IngressReconciler struct {
	client.Client

	Log             logr.Logger
	Scheme          *runtime.Scheme
	DataplaneClient *dataplane.KongClient

	DataplaneAddressFinder *dataplane.AddressFinder
	StatusQueue            *status.Queue

	IngressClassName string
}

// SetupWithManager sets up the controller with the Manager.
func (r *NetV1IngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("NetV1Ingress", mgr, controller.Options{
		Reconciler: r,
		Log:        r.Log,
	})
	if err != nil {
		return err
	}
	// if configured, start the status updater controller
	if r.StatusQueue != nil {
		if err := c.Watch(
			&source.Channel{Source: r.StatusQueue.Subscribe(schema.GroupVersionKind{
				Group:   "networking.k8s.io",
				Version: "v1",
				Kind:    "Ingress",
			})},
			&handler.EnqueueRequestForObject{},
		); err != nil {
			return err
		}
	}
	err = c.Watch(
		&source.Kind{Type: &netv1.IngressClass{}},
		handler.EnqueueRequestsFromMapFunc(r.reconcileIngressesWithoutIngressClass),
		predicate.NewPredicateFuncs(ctrlutils.IsDefaultIngressClass),
	)
	if err != nil {
		return err
	}
	preds := ctrlutils.GeneratePredicateFuncsForIngressClassFilter(r.IngressClassName, true, true)
	return c.Watch(
		&source.Kind{Type: &netv1.Ingress{}},
		&handler.EnqueueRequestForObject{},
		preds,
	)
}

// reconcileIngressesWithoutIngressClass finds and reconciles all Ingress objects without spec.ingressClassName set
func (r *NetV1IngressReconciler) reconcileIngressesWithoutIngressClass(obj client.Object) []reconcile.Request {
	ingresses := &netv1.IngressList{}
	if err := r.Client.List(context.Background(), ingresses); err != nil {
		r.Log.Error(err, "failed to list classless ingresses for default class")
		return nil
	}
	var recs []reconcile.Request
	for _, ingress := range ingresses.Items {
		if ingress.Spec.IngressClassName == nil {
			recs = append(recs, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: ingress.Namespace,
					Name:      ingress.Name,
				},
			})
		}
	}
	return recs
}

//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses/status,verbs=get;update;patch

// Reconcile processes the watched objects
func (r *NetV1IngressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("NetV1Ingress", req.NamespacedName)

	// get the relevant object
	obj := new(netv1.Ingress)
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		if errors.IsNotFound(err) {
			obj.Namespace = req.Namespace
			obj.Name = req.Name
			return ctrl.Result{}, r.DataplaneClient.DeleteObject(obj)
		}
		return ctrl.Result{}, err
	}
	log.V(util.DebugLevel).Info("reconciling resource", "namespace", req.Namespace, "name", req.Name)

	// clean the object up if it's being deleted
	if !obj.DeletionTimestamp.IsZero() && time.Now().After(obj.DeletionTimestamp.Time) {
		log.V(util.DebugLevel).Info("resource is being deleted, its configuration will be removed", "type", "Ingress", "namespace", req.Namespace, "name", req.Name)
		objectExistsInCache, err := r.DataplaneClient.ObjectExists(obj)
		if err != nil {
			return ctrl.Result{}, err
		}
		if objectExistsInCache {
			if err := r.DataplaneClient.DeleteObject(obj); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil // wait until the object is no longer present in the cache
		}
		return ctrl.Result{}, nil
	}

	// retrieve the configured IngressClass, to check if it's the default IngressClass
	class := new(netv1.IngressClass)
	if err := r.Get(ctx, types.NamespacedName{Name: r.IngressClassName}, class); err != nil {
		// we log this without taking action to support legacy configurations that only set ingressClassName or
		// used the class annotation and did not create a corresponding IngressClass. We only need this to determine
		// if the IngressClass is default or to configure default settings, and can assume no/no additional defaults
		// if none exists.
		log.V(util.DebugLevel).Info("could not retrieve IngressClass", "ingressclass", r.IngressClassName)
	}

	// if the object is not configured with our ingress.class, then we need to ensure it's removed from the cache
	if !ctrlutils.MatchesIngressClassName(obj, r.IngressClassName, ctrlutils.IsDefaultIngressClass(class)) {
		log.V(util.DebugLevel).Info("object missing ingress class, ensuring it's removed from configuration", "namespace", req.Namespace, "name", req.Name)
		return ctrl.Result{}, r.DataplaneClient.DeleteObject(obj)
	}

	// update the kong Admin API with the changes
	if err := r.DataplaneClient.UpdateObject(obj); err != nil {
		return ctrl.Result{}, err
	}
	// if status updates are enabled report the status for the object
	if r.DataplaneClient.AreKubernetesObjectReportsEnabled() {
		log.V(util.DebugLevel).Info("determining whether data-plane configuration has succeeded", "namespace", req.Namespace, "name", req.Name)
		if !r.DataplaneClient.KubernetesObjectIsConfigured(obj) {
			log.V(util.DebugLevel).Error(fmt.Errorf("resource not yet configured in the data-plane"), "namespace", req.Namespace, "name", req.Name)
			return ctrl.Result{Requeue: true}, nil // requeue until the object has been properly configured
		}

		log.V(util.DebugLevel).Info("determining gateway addresses for object status updates", "namespace", req.Namespace, "name", req.Name)
		addrs, err := r.DataplaneAddressFinder.GetLoadBalancerAddresses()
		if err != nil {
			return ctrl.Result{}, err
		}

		log.V(util.DebugLevel).Info("found addresses for data-plane updating object status", "namespace", req.Namespace, "name", req.Name)
		if len(obj.Status.LoadBalancer.Ingress) != len(addrs) || !reflect.DeepEqual(obj.Status.LoadBalancer.Ingress, addrs) {
			obj.Status.LoadBalancer.Ingress = addrs
			return ctrl.Result{}, r.Status().Update(ctx, obj)
		}
		log.V(util.DebugLevel).Info("status update not needed", "namespace", req.Namespace, "name", req.Name)
	}

	return ctrl.Result{}, nil
}
