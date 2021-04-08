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
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/kong/kubernetes-ingress-controller/pkg/deckgen"
	"github.com/kong/kubernetes-ingress-controller/pkg/parser"
	"github.com/kong/kubernetes-ingress-controller/pkg/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/pkg/store"
)

type SecretReconcilerParams struct {
	WatchName      string
	WatchNamespace string

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
	return object.GetName() == r.Params.WatchName &&
		object.GetNamespace() == r.Params.WatchNamespace
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
	logruslogger := logrus.New()

	// pull the configuration secret from the API
	configSecret := new(corev1.Secret)
	if err := r.Get(ctx, req.NamespacedName, configSecret); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// if the configuration secret is empty, something is wrong and we need to wait.
	if len(configSecret.Data) < 1 {
		return ctrl.Result{}, fmt.Errorf("no data is present in the configuration secret")
	}

	// build a new cache store from the objects present in the configuration secret
	yamls := make([][]byte, 0, len(configSecret.Data))
	for _, yaml := range configSecret.Data {
		yamls = append(yamls, yaml)
	}
	cache, err := store.NewCacheStoresFromObjYAML(yamls...)
	if err != nil {
		return ctrl.Result{}, err
	}

	// build the storer from the cached objects
	// TODO; verify these arguments are right
	storer := store.New(cache, "kong", false, false, false, logruslogger)

	// build the kongstate object from the Kubernetes objects in the storer
	kongstate, err := parser.Build(logruslogger, storer)
	if err != nil {
		return ctrl.Result{}, err
	}

	// generate the deck configuration to be applied to the admin API
	targetConfig := deckgen.ToDeckContent(ctx, logruslogger, kongstate, nil, nil)

	// apply the configuration update in Kong
	timedCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err = sendconfig.PerformUpdate(timedCtx, logruslogger, &r.Params.KongConfig, true, false, targetConfig, nil, nil, nil)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}
