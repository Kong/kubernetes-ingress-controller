/*
Copyright 2021.

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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"gopkg.in/yaml.v2"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

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

//+kubebuilder:rbac:groups=networking.k8s.io.my.domain,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=networking.k8s.io.my.domain,resources=ingresses/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=networking.k8s.io.my.domain,resources=ingresses/finalizers,verbs=update

// Reconcile adds any v1.Ingress configured for use by Kong to the combined configuration secret used to configure
// the Kong Admin API to configure and add new Services and Routes for the Ingress object.
func (r *IngressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("ingress", req.NamespacedName)

	ing := new(netv1.Ingress)
	if err := r.Get(ctx, req.NamespacedName, ing); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// ensure this Ingress is managed by KONG
	if !isManaged(ing.ObjectMeta) {
		return ctrl.Result{}, nil
	}

	// marshal to YAML for later storage
	cfg, err := yaml.Marshal(ing)
	if err != nil {
		return ctrl.Result{}, err
	}

	// get the configuration secret
	secret, created, err := getOrCreateConfigSecret(ctx, r.Client, req.Namespace)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, err
	}
	if created {
		return ctrl.Result{Requeue: true}, nil
	}

	// get the storage key for this ingress object
	key := keyFor(ing, req.NamespacedName)
	if _, ok := secret.Data[key]; ok {
		// TODO: for debugging, but need to remove later
		r.Log.Info("ingress entry already exists and will be overwritten", "key", key)
	}

	// TODO: patch instead of update for perf
	secret.Data[key] = cfg
	if err := r.Update(ctx, secret); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}
