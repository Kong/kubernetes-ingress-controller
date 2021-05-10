/*
Copyright 2021 Kong, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by Kong; DO NOT EDIT.

package configuration

import (
	"context"
	"time"

	"github.com/go-logr/logr"

	corev1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	netv1 "k8s.io/api/networking/v1"
	netv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/kong/kubernetes-ingress-controller/pkg/annotations"
	kongv1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1"
	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1alpha1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/railgun/internal/ctrlutils"
	"github.com/kong/kubernetes-ingress-controller/railgun/internal/proxy"
)

// -----------------------------------------------------------------------------
// CoreV1 Service
// -----------------------------------------------------------------------------

// CoreV1Service reconciles a Ingress object
type CoreV1ServiceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
	Proxy  proxy.Proxy
}

// SetupWithManager sets up the controller with the Manager.
func (r *CoreV1ServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&corev1.Service{}).Complete(r)
}

//+kubebuilder:rbac:groups=v1,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=v1,resources=services/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=v1,resources=services/finalizers,verbs=update

// Reconcile processes the watched objects
func (r *CoreV1ServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("CoreV1Service", req.NamespacedName)

	// get the relevant object
	obj := new(corev1.Service)
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log.Info("reconciling resource", "namespace", req.Namespace, "name", req.Name)

	// clean the object up if it's being deleted
	if !obj.DeletionTimestamp.IsZero() && time.Now().After(obj.DeletionTimestamp.Time) {
		log.Info("resource is being deleted, its configuration will be removed", "type", "Service", "namespace", req.Namespace, "name", req.Name)
		if err := r.Proxy.DeleteObject(obj); err != nil {
			return ctrl.Result{}, err
		}
		return ctrlutils.CleanupFinalizer(ctx, r.Client, log, req.NamespacedName, obj)
	}

	// before we store cache data for this object, ensure that it has our finalizer set
	if !ctrlutils.HasFinalizer(obj, ctrlutils.KongIngressFinalizer) {
		log.Info("finalizer is not set for ingress object, setting it", req.Namespace, req.Name)
		finalizers := obj.GetFinalizers()
		obj.SetFinalizers(append(finalizers, ctrlutils.KongIngressFinalizer))
		if err := r.Client.Update(ctx, obj); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// update the kong Admin API with the changes
	log.Info("updating the proxy with new Service", "namespace", obj.Namespace, "name", obj.Name)
	if err := r.Proxy.UpdateObject(obj); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// -----------------------------------------------------------------------------
// CoreV1 Endpoints
// -----------------------------------------------------------------------------

// CoreV1Endpoints reconciles a Ingress object
type CoreV1EndpointsReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
	Proxy  proxy.Proxy
}

// SetupWithManager sets up the controller with the Manager.
func (r *CoreV1EndpointsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&corev1.Endpoints{}).Complete(r)
}

//+kubebuilder:rbac:groups=v1,resources=endpoints,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=v1,resources=endpoints/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=v1,resources=endpoints/finalizers,verbs=update

// Reconcile processes the watched objects
func (r *CoreV1EndpointsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("CoreV1Endpoints", req.NamespacedName)

	// get the relevant object
	obj := new(corev1.Endpoints)
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log.Info("reconciling resource", "namespace", req.Namespace, "name", req.Name)

	// clean the object up if it's being deleted
	if !obj.DeletionTimestamp.IsZero() && time.Now().After(obj.DeletionTimestamp.Time) {
		log.Info("resource is being deleted, its configuration will be removed", "type", "Endpoints", "namespace", req.Namespace, "name", req.Name)
		if err := r.Proxy.DeleteObject(obj); err != nil {
			return ctrl.Result{}, err
		}
		return ctrlutils.CleanupFinalizer(ctx, r.Client, log, req.NamespacedName, obj)
	}

	// before we store cache data for this object, ensure that it has our finalizer set
	if !ctrlutils.HasFinalizer(obj, ctrlutils.KongIngressFinalizer) {
		log.Info("finalizer is not set for ingress object, setting it", req.Namespace, req.Name)
		finalizers := obj.GetFinalizers()
		obj.SetFinalizers(append(finalizers, ctrlutils.KongIngressFinalizer))
		if err := r.Client.Update(ctx, obj); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// update the kong Admin API with the changes
	log.Info("updating the proxy with new Endpoints", "namespace", obj.Namespace, "name", obj.Name)
	if err := r.Proxy.UpdateObject(obj); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// -----------------------------------------------------------------------------
// NetV1 Ingress
// -----------------------------------------------------------------------------

// NetV1Ingress reconciles a Ingress object
type NetV1IngressReconciler struct {
	client.Client
	Log              logr.Logger
	Scheme           *runtime.Scheme
	Proxy            proxy.Proxy
	IngressClassName string
}

// SetupWithManager sets up the controller with the Manager.
func (r *NetV1IngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	preds := predicate.NewPredicateFuncs(func(obj client.Object) bool {
		if v, ok := obj.GetAnnotations()[annotations.IngressClassKey]; ok {
			if v == r.IngressClassName {
				return true
			}
		}
		if ing, ok := obj.(*netv1.Ingress); ok {
			if ing.Spec.IngressClassName != nil && *ing.Spec.IngressClassName == r.IngressClassName {
				return true
			}
		}
		return false
	})
	preds.UpdateFunc = func(e event.UpdateEvent) bool {
		// at least one of the objects (old or new) needs to be configured with the relevant ingress.class to be supported.
		if v, ok := e.ObjectOld.GetAnnotations()[annotations.IngressClassKey]; ok {
			if v == r.IngressClassName {
				return true
			}
		}
		if v, ok := e.ObjectNew.GetAnnotations()[annotations.IngressClassKey]; ok {
			if v == r.IngressClassName {
				return true
			}
		}
		if ing, ok := e.ObjectOld.(*netv1.Ingress); ok {
			if ing.Spec.IngressClassName != nil && *ing.Spec.IngressClassName == r.IngressClassName {
				return true
			}
		}
		if ing, ok := e.ObjectNew.(*netv1.Ingress); ok {
			if ing.Spec.IngressClassName != nil && *ing.Spec.IngressClassName == r.IngressClassName {
				return true
			}
		}
		return false
	}
	return ctrl.NewControllerManagedBy(mgr).For(&netv1.Ingress{}, builder.WithPredicates(preds)).Complete(r)
}

//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses/finalizers,verbs=update

// Reconcile processes the watched objects
func (r *NetV1IngressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("NetV1Ingress", req.NamespacedName)

	// get the relevant object
	obj := new(netv1.Ingress)
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log.Info("reconciling resource", "namespace", req.Namespace, "name", req.Name)

	// clean the object up if it's being deleted
	if !obj.DeletionTimestamp.IsZero() && time.Now().After(obj.DeletionTimestamp.Time) {
		log.Info("resource is being deleted, its configuration will be removed", "type", "Ingress", "namespace", req.Namespace, "name", req.Name)
		if err := r.Proxy.DeleteObject(obj); err != nil {
			return ctrl.Result{}, err
		}
		return ctrlutils.CleanupFinalizer(ctx, r.Client, log, req.NamespacedName, obj)
	}

	// if the object is not configured with our ingress.class, then we need to ensure it's removed from the cache
	if !ctrlutils.MatchesIngressClassName(obj, r.IngressClassName) {
		log.Info("object missing ingress class, ensuring it's removed from configuration", req.Namespace, req.Name)
		if err := r.Proxy.DeleteObject(obj); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// before we store cache data for this object, ensure that it has our finalizer set
	if !ctrlutils.HasFinalizer(obj, ctrlutils.KongIngressFinalizer) {
		log.Info("finalizer is not set for ingress object, setting it", req.Namespace, req.Name)
		finalizers := obj.GetFinalizers()
		obj.SetFinalizers(append(finalizers, ctrlutils.KongIngressFinalizer))
		if err := r.Client.Update(ctx, obj); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// update the kong Admin API with the changes
	log.Info("updating the proxy with new Ingress", "namespace", obj.Namespace, "name", obj.Name)
	if err := r.Proxy.UpdateObject(obj); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// -----------------------------------------------------------------------------
// NetV1Beta1 Ingress
// -----------------------------------------------------------------------------

// NetV1Beta1Ingress reconciles a Ingress object
type NetV1Beta1IngressReconciler struct {
	client.Client
	Log              logr.Logger
	Scheme           *runtime.Scheme
	Proxy            proxy.Proxy
	IngressClassName string
}

// SetupWithManager sets up the controller with the Manager.
func (r *NetV1Beta1IngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	preds := predicate.NewPredicateFuncs(func(obj client.Object) bool {
		if v, ok := obj.GetAnnotations()[annotations.IngressClassKey]; ok {
			if v == r.IngressClassName {
				return true
			}
		}
		return false
	})
	preds.UpdateFunc = func(e event.UpdateEvent) bool {
		// at least one of the objects (old or new) needs to be configured with the relevant ingress.class to be supported.
		if v, ok := e.ObjectOld.GetAnnotations()[annotations.IngressClassKey]; ok {
			if v == r.IngressClassName {
				return true
			}
		}
		if v, ok := e.ObjectNew.GetAnnotations()[annotations.IngressClassKey]; ok {
			if v == r.IngressClassName {
				return true
			}
		}
		return false
	}
	return ctrl.NewControllerManagedBy(mgr).For(&netv1beta1.Ingress{}, builder.WithPredicates(preds)).Complete(r)
}

//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses/finalizers,verbs=update

// Reconcile processes the watched objects
func (r *NetV1Beta1IngressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("NetV1Beta1Ingress", req.NamespacedName)

	// get the relevant object
	obj := new(netv1beta1.Ingress)
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log.Info("reconciling resource", "namespace", req.Namespace, "name", req.Name)

	// clean the object up if it's being deleted
	if !obj.DeletionTimestamp.IsZero() && time.Now().After(obj.DeletionTimestamp.Time) {
		log.Info("resource is being deleted, its configuration will be removed", "type", "Ingress", "namespace", req.Namespace, "name", req.Name)
		if err := r.Proxy.DeleteObject(obj); err != nil {
			return ctrl.Result{}, err
		}
		return ctrlutils.CleanupFinalizer(ctx, r.Client, log, req.NamespacedName, obj)
	}

	// if the object is not configured with our ingress.class, then we need to ensure it's removed from the cache
	if !ctrlutils.MatchesIngressClassName(obj, r.IngressClassName) {
		log.Info("object missing ingress class, ensuring it's removed from configuration", req.Namespace, req.Name)
		if err := r.Proxy.DeleteObject(obj); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// before we store cache data for this object, ensure that it has our finalizer set
	if !ctrlutils.HasFinalizer(obj, ctrlutils.KongIngressFinalizer) {
		log.Info("finalizer is not set for ingress object, setting it", req.Namespace, req.Name)
		finalizers := obj.GetFinalizers()
		obj.SetFinalizers(append(finalizers, ctrlutils.KongIngressFinalizer))
		if err := r.Client.Update(ctx, obj); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// update the kong Admin API with the changes
	log.Info("updating the proxy with new Ingress", "namespace", obj.Namespace, "name", obj.Name)
	if err := r.Proxy.UpdateObject(obj); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// -----------------------------------------------------------------------------
// ExtV1Beta1 Ingress
// -----------------------------------------------------------------------------

// ExtV1Beta1Ingress reconciles a Ingress object
type ExtV1Beta1IngressReconciler struct {
	client.Client
	Log              logr.Logger
	Scheme           *runtime.Scheme
	Proxy            proxy.Proxy
	IngressClassName string
}

// SetupWithManager sets up the controller with the Manager.
func (r *ExtV1Beta1IngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	preds := predicate.NewPredicateFuncs(func(obj client.Object) bool {
		if v, ok := obj.GetAnnotations()[annotations.IngressClassKey]; ok {
			if v == r.IngressClassName {
				return true
			}
		}
		return false
	})
	preds.UpdateFunc = func(e event.UpdateEvent) bool {
		// at least one of the objects (old or new) needs to be configured with the relevant ingress.class to be supported.
		if v, ok := e.ObjectOld.GetAnnotations()[annotations.IngressClassKey]; ok {
			if v == r.IngressClassName {
				return true
			}
		}
		if v, ok := e.ObjectNew.GetAnnotations()[annotations.IngressClassKey]; ok {
			if v == r.IngressClassName {
				return true
			}
		}
		return false
	}
	return ctrl.NewControllerManagedBy(mgr).For(&extv1beta1.Ingress{}, builder.WithPredicates(preds)).Complete(r)
}

//+kubebuilder:rbac:groups=apiextensions.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apiextensions.k8s.io,resources=ingresses/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apiextensions.k8s.io,resources=ingresses/finalizers,verbs=update

// Reconcile processes the watched objects
func (r *ExtV1Beta1IngressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("ExtV1Beta1Ingress", req.NamespacedName)

	// get the relevant object
	obj := new(extv1beta1.Ingress)
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log.Info("reconciling resource", "namespace", req.Namespace, "name", req.Name)

	// clean the object up if it's being deleted
	if !obj.DeletionTimestamp.IsZero() && time.Now().After(obj.DeletionTimestamp.Time) {
		log.Info("resource is being deleted, its configuration will be removed", "type", "Ingress", "namespace", req.Namespace, "name", req.Name)
		if err := r.Proxy.DeleteObject(obj); err != nil {
			return ctrl.Result{}, err
		}
		return ctrlutils.CleanupFinalizer(ctx, r.Client, log, req.NamespacedName, obj)
	}

	// if the object is not configured with our ingress.class, then we need to ensure it's removed from the cache
	if !ctrlutils.MatchesIngressClassName(obj, r.IngressClassName) {
		log.Info("object missing ingress class, ensuring it's removed from configuration", req.Namespace, req.Name)
		if err := r.Proxy.DeleteObject(obj); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// before we store cache data for this object, ensure that it has our finalizer set
	if !ctrlutils.HasFinalizer(obj, ctrlutils.KongIngressFinalizer) {
		log.Info("finalizer is not set for ingress object, setting it", req.Namespace, req.Name)
		finalizers := obj.GetFinalizers()
		obj.SetFinalizers(append(finalizers, ctrlutils.KongIngressFinalizer))
		if err := r.Client.Update(ctx, obj); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// update the kong Admin API with the changes
	log.Info("updating the proxy with new Ingress", "namespace", obj.Namespace, "name", obj.Name)
	if err := r.Proxy.UpdateObject(obj); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// -----------------------------------------------------------------------------
// KongV1 KongIngress
// -----------------------------------------------------------------------------

// KongV1KongIngress reconciles a Ingress object
type KongV1KongIngressReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
	Proxy  proxy.Proxy
}

// SetupWithManager sets up the controller with the Manager.
func (r *KongV1KongIngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&kongv1.KongIngress{}).Complete(r)
}

//+kubebuilder:rbac:groups=configuration.konghq.com,resources=kongingresses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=configuration.konghq.com,resources=kongingresses/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=configuration.konghq.com,resources=kongingresses/finalizers,verbs=update

// Reconcile processes the watched objects
func (r *KongV1KongIngressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("KongV1KongIngress", req.NamespacedName)

	// get the relevant object
	obj := new(kongv1.KongIngress)
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log.Info("reconciling resource", "namespace", req.Namespace, "name", req.Name)

	// clean the object up if it's being deleted
	if !obj.DeletionTimestamp.IsZero() && time.Now().After(obj.DeletionTimestamp.Time) {
		log.Info("resource is being deleted, its configuration will be removed", "type", "KongIngress", "namespace", req.Namespace, "name", req.Name)
		if err := r.Proxy.DeleteObject(obj); err != nil {
			return ctrl.Result{}, err
		}
		return ctrlutils.CleanupFinalizer(ctx, r.Client, log, req.NamespacedName, obj)
	}

	// before we store cache data for this object, ensure that it has our finalizer set
	if !ctrlutils.HasFinalizer(obj, ctrlutils.KongIngressFinalizer) {
		log.Info("finalizer is not set for ingress object, setting it", req.Namespace, req.Name)
		finalizers := obj.GetFinalizers()
		obj.SetFinalizers(append(finalizers, ctrlutils.KongIngressFinalizer))
		if err := r.Client.Update(ctx, obj); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// update the kong Admin API with the changes
	log.Info("updating the proxy with new KongIngress", "namespace", obj.Namespace, "name", obj.Name)
	if err := r.Proxy.UpdateObject(obj); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// -----------------------------------------------------------------------------
// KongV1 KongPlugin
// -----------------------------------------------------------------------------

// KongV1KongPlugin reconciles a Ingress object
type KongV1KongPluginReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
	Proxy  proxy.Proxy
}

// SetupWithManager sets up the controller with the Manager.
func (r *KongV1KongPluginReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&kongv1.KongPlugin{}).Complete(r)
}

//+kubebuilder:rbac:groups=configuration.konghq.com,resources=kongplugins,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=configuration.konghq.com,resources=kongplugins/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=configuration.konghq.com,resources=kongplugins/finalizers,verbs=update

// Reconcile processes the watched objects
func (r *KongV1KongPluginReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("KongV1KongPlugin", req.NamespacedName)

	// get the relevant object
	obj := new(kongv1.KongPlugin)
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log.Info("reconciling resource", "namespace", req.Namespace, "name", req.Name)

	// clean the object up if it's being deleted
	if !obj.DeletionTimestamp.IsZero() && time.Now().After(obj.DeletionTimestamp.Time) {
		log.Info("resource is being deleted, its configuration will be removed", "type", "KongPlugin", "namespace", req.Namespace, "name", req.Name)
		if err := r.Proxy.DeleteObject(obj); err != nil {
			return ctrl.Result{}, err
		}
		return ctrlutils.CleanupFinalizer(ctx, r.Client, log, req.NamespacedName, obj)
	}

	// before we store cache data for this object, ensure that it has our finalizer set
	if !ctrlutils.HasFinalizer(obj, ctrlutils.KongIngressFinalizer) {
		log.Info("finalizer is not set for ingress object, setting it", req.Namespace, req.Name)
		finalizers := obj.GetFinalizers()
		obj.SetFinalizers(append(finalizers, ctrlutils.KongIngressFinalizer))
		if err := r.Client.Update(ctx, obj); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// update the kong Admin API with the changes
	log.Info("updating the proxy with new KongPlugin", "namespace", obj.Namespace, "name", obj.Name)
	if err := r.Proxy.UpdateObject(obj); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// -----------------------------------------------------------------------------
// KongV1 KongClusterPlugin
// -----------------------------------------------------------------------------

// KongV1KongClusterPlugin reconciles a Ingress object
type KongV1KongClusterPluginReconciler struct {
	client.Client
	Log              logr.Logger
	Scheme           *runtime.Scheme
	Proxy            proxy.Proxy
	IngressClassName string
}

// SetupWithManager sets up the controller with the Manager.
func (r *KongV1KongClusterPluginReconciler) SetupWithManager(mgr ctrl.Manager) error {
	preds := predicate.NewPredicateFuncs(func(obj client.Object) bool {
		if v, ok := obj.GetAnnotations()[annotations.IngressClassKey]; ok {
			if v == r.IngressClassName {
				return true
			}
		}
		return false
	})
	preds.UpdateFunc = func(e event.UpdateEvent) bool {
		// at least one of the objects (old or new) needs to be configured with the relevant ingress.class to be supported.
		if v, ok := e.ObjectOld.GetAnnotations()[annotations.IngressClassKey]; ok {
			if v == r.IngressClassName {
				return true
			}
		}
		if v, ok := e.ObjectNew.GetAnnotations()[annotations.IngressClassKey]; ok {
			if v == r.IngressClassName {
				return true
			}
		}
		return false
	}
	return ctrl.NewControllerManagedBy(mgr).For(&kongv1.KongClusterPlugin{}, builder.WithPredicates(preds)).Complete(r)
}

//+kubebuilder:rbac:groups=configuration.konghq.com,resources=kongclusterplugins,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=configuration.konghq.com,resources=kongclusterplugins/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=configuration.konghq.com,resources=kongclusterplugins/finalizers,verbs=update

// Reconcile processes the watched objects
func (r *KongV1KongClusterPluginReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("KongV1KongClusterPlugin", req.NamespacedName)

	// get the relevant object
	obj := new(kongv1.KongClusterPlugin)
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log.Info("reconciling resource", "namespace", req.Namespace, "name", req.Name)

	// clean the object up if it's being deleted
	if !obj.DeletionTimestamp.IsZero() && time.Now().After(obj.DeletionTimestamp.Time) {
		log.Info("resource is being deleted, its configuration will be removed", "type", "KongClusterPlugin", "namespace", req.Namespace, "name", req.Name)
		if err := r.Proxy.DeleteObject(obj); err != nil {
			return ctrl.Result{}, err
		}
		return ctrlutils.CleanupFinalizer(ctx, r.Client, log, req.NamespacedName, obj)
	}

	// if the object is not configured with our ingress.class, then we need to ensure it's removed from the cache
	if !ctrlutils.MatchesIngressClassName(obj, r.IngressClassName) {
		log.Info("object missing ingress class, ensuring it's removed from configuration", req.Namespace, req.Name)
		if err := r.Proxy.DeleteObject(obj); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// before we store cache data for this object, ensure that it has our finalizer set
	if !ctrlutils.HasFinalizer(obj, ctrlutils.KongIngressFinalizer) {
		log.Info("finalizer is not set for ingress object, setting it", req.Namespace, req.Name)
		finalizers := obj.GetFinalizers()
		obj.SetFinalizers(append(finalizers, ctrlutils.KongIngressFinalizer))
		if err := r.Client.Update(ctx, obj); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// update the kong Admin API with the changes
	log.Info("updating the proxy with new KongClusterPlugin", "namespace", obj.Namespace, "name", obj.Name)
	if err := r.Proxy.UpdateObject(obj); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// -----------------------------------------------------------------------------
// KongV1 KongConsumer
// -----------------------------------------------------------------------------

// KongV1KongConsumer reconciles a Ingress object
type KongV1KongConsumerReconciler struct {
	client.Client
	Log              logr.Logger
	Scheme           *runtime.Scheme
	Proxy            proxy.Proxy
	IngressClassName string
}

// SetupWithManager sets up the controller with the Manager.
func (r *KongV1KongConsumerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	preds := predicate.NewPredicateFuncs(func(obj client.Object) bool {
		if v, ok := obj.GetAnnotations()[annotations.IngressClassKey]; ok {
			if v == r.IngressClassName {
				return true
			}
		}
		return false
	})
	preds.UpdateFunc = func(e event.UpdateEvent) bool {
		// at least one of the objects (old or new) needs to be configured with the relevant ingress.class to be supported.
		if v, ok := e.ObjectOld.GetAnnotations()[annotations.IngressClassKey]; ok {
			if v == r.IngressClassName {
				return true
			}
		}
		if v, ok := e.ObjectNew.GetAnnotations()[annotations.IngressClassKey]; ok {
			if v == r.IngressClassName {
				return true
			}
		}
		return false
	}
	return ctrl.NewControllerManagedBy(mgr).For(&kongv1.KongConsumer{}, builder.WithPredicates(preds)).Complete(r)
}

//+kubebuilder:rbac:groups=configuration.konghq.com,resources=kongconsumers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=configuration.konghq.com,resources=kongconsumers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=configuration.konghq.com,resources=kongconsumers/finalizers,verbs=update

// Reconcile processes the watched objects
func (r *KongV1KongConsumerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("KongV1KongConsumer", req.NamespacedName)

	// get the relevant object
	obj := new(kongv1.KongConsumer)
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log.Info("reconciling resource", "namespace", req.Namespace, "name", req.Name)

	// clean the object up if it's being deleted
	if !obj.DeletionTimestamp.IsZero() && time.Now().After(obj.DeletionTimestamp.Time) {
		log.Info("resource is being deleted, its configuration will be removed", "type", "KongConsumer", "namespace", req.Namespace, "name", req.Name)
		if err := r.Proxy.DeleteObject(obj); err != nil {
			return ctrl.Result{}, err
		}
		return ctrlutils.CleanupFinalizer(ctx, r.Client, log, req.NamespacedName, obj)
	}

	// if the object is not configured with our ingress.class, then we need to ensure it's removed from the cache
	if !ctrlutils.MatchesIngressClassName(obj, r.IngressClassName) {
		log.Info("object missing ingress class, ensuring it's removed from configuration", req.Namespace, req.Name)
		if err := r.Proxy.DeleteObject(obj); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// before we store cache data for this object, ensure that it has our finalizer set
	if !ctrlutils.HasFinalizer(obj, ctrlutils.KongIngressFinalizer) {
		log.Info("finalizer is not set for ingress object, setting it", req.Namespace, req.Name)
		finalizers := obj.GetFinalizers()
		obj.SetFinalizers(append(finalizers, ctrlutils.KongIngressFinalizer))
		if err := r.Client.Update(ctx, obj); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// update the kong Admin API with the changes
	log.Info("updating the proxy with new KongConsumer", "namespace", obj.Namespace, "name", obj.Name)
	if err := r.Proxy.UpdateObject(obj); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// -----------------------------------------------------------------------------
// KongV1Alpha1 UDPIngress
// -----------------------------------------------------------------------------

// KongV1Alpha1UDPIngress reconciles a Ingress object
type KongV1Alpha1UDPIngressReconciler struct {
	client.Client
	Log              logr.Logger
	Scheme           *runtime.Scheme
	Proxy            proxy.Proxy
	IngressClassName string
}

// SetupWithManager sets up the controller with the Manager.
func (r *KongV1Alpha1UDPIngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	preds := predicate.NewPredicateFuncs(func(obj client.Object) bool {
		if v, ok := obj.GetAnnotations()[annotations.IngressClassKey]; ok {
			if v == r.IngressClassName {
				return true
			}
		}
		return false
	})
	preds.UpdateFunc = func(e event.UpdateEvent) bool {
		// at least one of the objects (old or new) needs to be configured with the relevant ingress.class to be supported.
		if v, ok := e.ObjectOld.GetAnnotations()[annotations.IngressClassKey]; ok {
			if v == r.IngressClassName {
				return true
			}
		}
		if v, ok := e.ObjectNew.GetAnnotations()[annotations.IngressClassKey]; ok {
			if v == r.IngressClassName {
				return true
			}
		}
		return false
	}
	return ctrl.NewControllerManagedBy(mgr).For(&kongv1alpha1.UDPIngress{}, builder.WithPredicates(preds)).Complete(r)
}

//+kubebuilder:rbac:groups=configuration.konghq.com,resources=udpingresses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=configuration.konghq.com,resources=udpingresses/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=configuration.konghq.com,resources=udpingresses/finalizers,verbs=update

// Reconcile processes the watched objects
func (r *KongV1Alpha1UDPIngressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("KongV1Alpha1UDPIngress", req.NamespacedName)

	// get the relevant object
	obj := new(kongv1alpha1.UDPIngress)
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log.Info("reconciling resource", "namespace", req.Namespace, "name", req.Name)

	// clean the object up if it's being deleted
	if !obj.DeletionTimestamp.IsZero() && time.Now().After(obj.DeletionTimestamp.Time) {
		log.Info("resource is being deleted, its configuration will be removed", "type", "UDPIngress", "namespace", req.Namespace, "name", req.Name)
		if err := r.Proxy.DeleteObject(obj); err != nil {
			return ctrl.Result{}, err
		}
		return ctrlutils.CleanupFinalizer(ctx, r.Client, log, req.NamespacedName, obj)
	}

	// if the object is not configured with our ingress.class, then we need to ensure it's removed from the cache
	if !ctrlutils.MatchesIngressClassName(obj, r.IngressClassName) {
		log.Info("object missing ingress class, ensuring it's removed from configuration", req.Namespace, req.Name)
		if err := r.Proxy.DeleteObject(obj); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// before we store cache data for this object, ensure that it has our finalizer set
	if !ctrlutils.HasFinalizer(obj, ctrlutils.KongIngressFinalizer) {
		log.Info("finalizer is not set for ingress object, setting it", req.Namespace, req.Name)
		finalizers := obj.GetFinalizers()
		obj.SetFinalizers(append(finalizers, ctrlutils.KongIngressFinalizer))
		if err := r.Client.Update(ctx, obj); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// update the kong Admin API with the changes
	log.Info("updating the proxy with new UDPIngress", "namespace", obj.Namespace, "name", obj.Name)
	if err := r.Proxy.UpdateObject(obj); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// -----------------------------------------------------------------------------
// KongV1Beta1 TCPIngress
// -----------------------------------------------------------------------------

// KongV1Beta1TCPIngress reconciles a Ingress object
type KongV1Beta1TCPIngressReconciler struct {
	client.Client
	Log              logr.Logger
	Scheme           *runtime.Scheme
	Proxy            proxy.Proxy
	IngressClassName string
}

// SetupWithManager sets up the controller with the Manager.
func (r *KongV1Beta1TCPIngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	preds := predicate.NewPredicateFuncs(func(obj client.Object) bool {
		if v, ok := obj.GetAnnotations()[annotations.IngressClassKey]; ok {
			if v == r.IngressClassName {
				return true
			}
		}
		return false
	})
	preds.UpdateFunc = func(e event.UpdateEvent) bool {
		// at least one of the objects (old or new) needs to be configured with the relevant ingress.class to be supported.
		if v, ok := e.ObjectOld.GetAnnotations()[annotations.IngressClassKey]; ok {
			if v == r.IngressClassName {
				return true
			}
		}
		if v, ok := e.ObjectNew.GetAnnotations()[annotations.IngressClassKey]; ok {
			if v == r.IngressClassName {
				return true
			}
		}
		return false
	}
	return ctrl.NewControllerManagedBy(mgr).For(&kongv1beta1.TCPIngress{}, builder.WithPredicates(preds)).Complete(r)
}

//+kubebuilder:rbac:groups=configuration.konghq.com,resources=tcpingresses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=configuration.konghq.com,resources=tcpingresses/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=configuration.konghq.com,resources=tcpingresses/finalizers,verbs=update

// Reconcile processes the watched objects
func (r *KongV1Beta1TCPIngressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("KongV1Beta1TCPIngress", req.NamespacedName)

	// get the relevant object
	obj := new(kongv1beta1.TCPIngress)
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log.Info("reconciling resource", "namespace", req.Namespace, "name", req.Name)

	// clean the object up if it's being deleted
	if !obj.DeletionTimestamp.IsZero() && time.Now().After(obj.DeletionTimestamp.Time) {
		log.Info("resource is being deleted, its configuration will be removed", "type", "TCPIngress", "namespace", req.Namespace, "name", req.Name)
		if err := r.Proxy.DeleteObject(obj); err != nil {
			return ctrl.Result{}, err
		}
		return ctrlutils.CleanupFinalizer(ctx, r.Client, log, req.NamespacedName, obj)
	}

	// if the object is not configured with our ingress.class, then we need to ensure it's removed from the cache
	if !ctrlutils.MatchesIngressClassName(obj, r.IngressClassName) {
		log.Info("object missing ingress class, ensuring it's removed from configuration", req.Namespace, req.Name)
		if err := r.Proxy.DeleteObject(obj); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// before we store cache data for this object, ensure that it has our finalizer set
	if !ctrlutils.HasFinalizer(obj, ctrlutils.KongIngressFinalizer) {
		log.Info("finalizer is not set for ingress object, setting it", req.Namespace, req.Name)
		finalizers := obj.GetFinalizers()
		obj.SetFinalizers(append(finalizers, ctrlutils.KongIngressFinalizer))
		if err := r.Client.Update(ctx, obj); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// update the kong Admin API with the changes
	log.Info("updating the proxy with new TCPIngress", "namespace", obj.Namespace, "name", obj.Name)
	if err := r.Proxy.UpdateObject(obj); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}
