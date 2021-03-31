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

package configuration

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"github.com/kong/kubernetes-ingress-controller/pkg/deckgen"
	"github.com/kong/kubernetes-ingress-controller/pkg/parser"
	"github.com/kong/kubernetes-ingress-controller/pkg/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/railgun/pkg/store"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type SecretReconcilerParams struct {
	WatchNamespacedName *types.NamespacedName

	KongConfig sendconfig.Kong
}

// SecretReconciler reconciles a Secret object
type SecretReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme

	Params SecretReconcilerParams
}

func (r *SecretReconciler) matchNsName(object client.Object) bool {
	return object.GetName() == r.Params.WatchNamespacedName.Name &&
		object.GetNamespace() == r.Params.WatchNamespacedName.Namespace
}

// SetupWithManager sets up the controller with the Manager.
func (r *SecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// TODO: something to keep in mind: long term we're still considering use a custom API instead of a secret for the Configuration.
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Secret{}).
		WithEventFilter(predicate.NewPredicateFuncs(r.matchNsName)).
		Complete(r)
}

//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=secrets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core,resources=secrets/finalizers,verbs=update

// Reconcile manages the configuration secret for ingresses and parses that into a Kong configuration
// which is posted to all available Proxy APIs.
func (r *SecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	//	log := r.Log.WithValues("secret", req.NamespacedName)
	logruslogger := logrus.New()

	configSecret := new(corev1.Secret)
	if err := r.Get(ctx, req.NamespacedName, configSecret); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	storer := store.New(r.Client)
	kongstate, err := parser.Build(logruslogger, storer)
	if err != nil {
		return ctrl.Result{}, err
	}

	targetConfig := deckgen.ToDeckContent(ctx, logruslogger, kongstate, nil, nil)

	timedCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err = sendconfig.PerformUpdate(timedCtx, logruslogger, &r.Params.KongConfig, true, false, targetConfig, nil, nil, nil)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}
