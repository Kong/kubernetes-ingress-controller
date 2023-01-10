package configuration

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// KongAdminAPIServiceReconciler reconciles Kong Admin API Secrets
type KongAdminAPIServiceReconciler struct {
	client.Client

	ServiceNN types.NamespacedName
	Service   *corev1.Service

	Log              logr.Logger
	CacheSyncTimeout time.Duration

	EndpointsNotifier EndpointsNotifier
}

type EndpointsNotifier interface {
	Notify(addresses []string)
}

// SetupWithManager sets up the controller with the Manager.
func (r *KongAdminAPIServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// var service corev1.Service
	// if err := r.Get(context.Background(), r.ServiceNN, &service); err != nil {
	// 	return fmt.Errorf("failed to get kong Admin API service %s: %w", r.ServiceNN, err)
	// }
	// r.Service = &service

	c, err := controller.New("KongAdminAPIEndpoints", mgr, controller.Options{
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
		&source.Kind{Type: &discoveryv1.EndpointSlice{}},
		&handler.EnqueueRequestForObject{},
		predicate.NewPredicateFuncs(r.shouldReconcileEndpointSlice),
	)
}

func (r *KongAdminAPIServiceReconciler) shouldReconcileEndpointSlice(obj client.Object) bool {
	endpoints, ok := obj.(*discoveryv1.EndpointSlice)
	if !ok {
		return false
	}

	if !lo.ContainsBy(endpoints.OwnerReferences, func(ref metav1.OwnerReference) bool {
		r.Log.Info("checking ref", "ref", ref)
		// if ref.UID != r.Service.UID {
		// 	return false
		// }

		if ref.Kind != "Service" || ref.Name != r.ServiceNN.Name {
			return false
		}
		r.Log.Info("ref ok", "ref", ref)
		return true
	}) {
		return false
	}

	return true // TODO xxx
}

//+kubebuilder:rbac:groups="discovery.k8s.io",resources=endpointslices,verbs=get;list;watch

// Reconcile processes the watched objects
func (r *KongAdminAPIServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// get the relevant object
	// var endpoints corev1.Endpoints
	var endpoints discoveryv1.EndpointSlice
	if err := r.Get(ctx, req.NamespacedName, &endpoints); err != nil {
		// if err := r.Get(ctx, req.NamespacedName, &endpoints); err != nil {
		return ctrl.Result{}, err
	}
	r.Log.Info("reconciling resource", "namespace", req.Namespace, "name", req.Name)

	if !endpoints.DeletionTimestamp.IsZero() {
		r.Log.Info("resource is being deleted", "type", "Service", "namespace", req.Namespace, "name", req.Name)
		return ctrl.Result{}, nil
	}

	addresses := AddressesFromEndpointSlice(endpoints)
	r.EndpointsNotifier.Notify(addresses)

	return ctrl.Result{}, nil
}

func AddressesFromEndpointSlice(endpoints discoveryv1.EndpointSlice) []string {
	var addresses []string
	for _, p := range endpoints.Ports {
		if p.Name == nil {
			continue
		}
		// TODO
		if *p.Name != "admin" && *p.Name != "kong-admin" {
			continue
		}

		for _, e := range endpoints.Endpoints {
			if e.Conditions.Ready == nil || !*e.Conditions.Ready {
				continue
			}

			for _, addr := range e.Addresses {
				addresses = append(addresses, fmt.Sprintf("https://%s:%d", addr, *p.Port))
			}
		}
	}
	return addresses
}
