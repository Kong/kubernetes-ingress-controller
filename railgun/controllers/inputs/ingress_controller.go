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

package inputs

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	netv1beta1 "k8s.io/api/extensions/v1beta1"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// -----------------------------------------------------------------------------
// V1 Ingress
// -----------------------------------------------------------------------------

// IngressReconciler reconciles a Ingress object
type IngressReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// SetupWithManager sets up the controller with the Manager.
func (r *IngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&netv1.Ingress{}).Complete(r)
}

//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses/finalizers,verbs=update

// Reconcile adds any v1.Ingress configured for use by Kong to the combined configuration secret used to configure
// the Kong Admin API to configure and add new Services and Routes for the Ingress object.
func (r *IngressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("ingress", req.NamespacedName)

	ing := new(netv1.Ingress)
	if err := r.Get(ctx, req.NamespacedName, ing); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !ing.DeletionTimestamp.IsZero() && time.Now().After(ing.DeletionTimestamp.Time) {
		// TODO: finalizer
		log.Info("ingress resource being deleted, its configuration will be removed", "namespace", req.Namespace, "name", req.Name)
		return cleanupIngress(ctx, r.Client, log, req.NamespacedName, ing)
	}

	return storeIngressUpdates(ctx, r.Client, log, req.NamespacedName, ing)
}

// -----------------------------------------------------------------------------
// V1Beta1 Ingress
// -----------------------------------------------------------------------------

// V1Beta1IngressReconciler reconciles a Ingress object
type V1Beta1IngressReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// SetupWithManager sets up the controller with the Manager.
func (r *V1Beta1IngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&netv1beta1.Ingress{}).Complete(r)
}

//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses/finalizers,verbs=update

// Reconcile adds any v1beta1.Ingress configured for use by Kong to the combined configuration secret used to configure
// the Kong Admin API to configure and add new Services and Routes for the Ingress object.
func (r *V1Beta1IngressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("v1beta1ingress", req.NamespacedName)

	ing := new(netv1beta1.Ingress)
	if err := r.Get(ctx, req.NamespacedName, ing); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !ing.DeletionTimestamp.IsZero() && time.Now().After(ing.DeletionTimestamp.Time) {
		// TODO: finalizer
		log.Info("ingress resource being deleted, its configuration will be removed", "namespace", req.Namespace, "name", req.Name)
		return cleanupIngress(ctx, r.Client, log, req.NamespacedName, ing)
	}

	return storeIngressUpdates(ctx, r.Client, log, req.NamespacedName, ing)
}
