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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	konghqcomv1 "github.com/kong/railgun/api/v1"
)

// KongIngressReconciler reconciles a KongIngress object
type KongIngressReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// SetupWithManager sets up the controller with the Manager.
func (r *KongIngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&konghqcomv1.KongIngress{}).Complete(r)
}

//+kubebuilder:rbac:groups=konghq.com.my.domain,resources=kongingresses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=konghq.com.my.domain,resources=kongingresses/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=konghq.com.my.domain,resources=kongingresses/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the KongIngress object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *KongIngressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("kongingress", req.NamespacedName)

	log := r.Log.WithValues("ingress", req.NamespacedName)

	ing := new(konghqcomv1.KongIngress)
	if err := r.Get(ctx, req.NamespacedName, ing); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !ing.DeletionTimestamp.IsZero() && time.Now().After(ing.DeletionTimestamp.Time) {
		log.Info("resource being deleted, its configuration will be removed", "namespace", req.Namespace, "name", req.Name)
		return cleanupObj(ctx, r.Client, log, req.NamespacedName, ing)
	}

	return storeObjUpdates(ctx, r.Client, log, req.NamespacedName, ing)

	//return ctrl.Result{}, nil
}
