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

package gateway

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

// ReferenceGrantReconciler reconciles a ReferenceGrant object.
type ReferenceGrantReconciler struct {
	client.Client
	Log             logr.Logger
	Scheme          *runtime.Scheme
	DataplaneClient controllers.DataPlane

	CacheSyncTimeout time.Duration
}

// SetupWithManager sets up the controller with the Manager.
func (r *ReferenceGrantReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// set the controller name
		Named("referencegrant-controller").
		WithOptions(controller.Options{
			LogConstructor: func(_ *reconcile.Request) logr.Logger {
				return r.Log
			},
			CacheSyncTimeout: r.CacheSyncTimeout,
		}).
		// watch Referencegrant objects
		For(&gatewayapi.ReferenceGrant{}).
		Complete(r)
}

// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=referencegrants,verbs=get;list;watch
// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=referencegrants/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ReferenceGrantReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("GatewayV1Alpha2ReferenceGrant", req.NamespacedName)
	grant := new(gatewayapi.ReferenceGrant)
	if err := r.Get(ctx, req.NamespacedName, grant); err != nil {
		// if the queued object is no longer present in the proxy cache we need
		// to ensure that if it was ever added to the cache, it gets removed.
		if apierrors.IsNotFound(err) {
			debug(log, grant, "Object does not exist, ensuring it is not present in the proxy cache")
			grant.Namespace = req.Namespace
			grant.Name = req.Name
			return ctrl.Result{}, r.DataplaneClient.DeleteObject(grant)
		}

		// for any error other than 404, requeue
		return ctrl.Result{}, err
	}

	debug(log, grant, "Processing referencegrant")

	debug(log, grant, "Checking deletion timestamp")
	if grant.DeletionTimestamp != nil {
		debug(log, grant, "Referencegrant is being deleted, re-configuring data-plane")
		if err := r.DataplaneClient.DeleteObject(grant); err != nil {
			debug(log, grant, "failed to delete object from data-plane, requeuing")
			return ctrl.Result{}, err
		}
		debug(log, grant, "Ensured object was removed from the data-plane (if ever present)")
		return ctrl.Result{}, r.DataplaneClient.DeleteObject(grant)
	}

	if err := r.DataplaneClient.UpdateObject(grant); err != nil {
		debug(log, grant, "Failed to update object in data-plane, requeueing")
		return ctrl.Result{}, err
	}
	info(log, grant, "Referencegrant has been configured on the data-plane")
	return ctrl.Result{}, nil
}
