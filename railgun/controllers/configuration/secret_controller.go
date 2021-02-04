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
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"time"

	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	"github.com/kong/railgun/controllers"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SecretReconciler reconciles a Secret object
type SecretReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// SetupWithManager sets up the controller with the Manager.
func (r *SecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// TODO: something to keep in mind: long term we're still considering use a custom API instead of a secret for the Configuration.
	return ctrl.NewControllerManagedBy(mgr).For(&corev1.Secret{}).Complete(r)
}

//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=secrets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core,resources=secrets/finalizers,verbs=update

// Reconcile manages the configuration secret for ingresses and parses that into a Kong configuration
// which is posted to all available Proxy APIs.
func (r *SecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("secret", req.NamespacedName)

	configSecret := new(corev1.Secret)
	if err := r.Get(ctx, req.NamespacedName, configSecret); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// get the configuration secret namespace and filter out others
	// TODO: we can do this with Watcher filters instead later, but if we're going to switch to a CRD
	// implementation anyway we may not have to bother with this.
	secretNamespace := os.Getenv(controllers.CtrlNamespaceEnv)
	if secretNamespace == "" {
		return ctrl.Result{}, fmt.Errorf("kong can not be configured because the required %s env var is not present", controllers.CtrlNamespaceEnv)
	}
	if req.Namespace != secretNamespace {
		return ctrl.Result{}, nil
	}
	if req.Name != controllers.ConfigSecretName {
		return ctrl.Result{}, nil
	}

	// collect all proxy instances to be configured
	proxyInstances := new(corev1.PodList)

	matchingLabels := client.MatchingLabels{controllers.ProxyInstanceLabel: "true"}
	if err := r.List(ctx, proxyInstances, matchingLabels); err != nil {
		return ctrl.Result{}, err
	}
	if len(proxyInstances.Items) < 1 {
		return ctrl.Result{}, fmt.Errorf("no proxy instances found (pods with label %s) is Kong running?", controllers.ProxyInstanceLabel)
	}
	log.Info("found proxy instances which need to be configured", "count", len(proxyInstances.Items))

	for _, pod := range proxyInstances.Items {
		ip := pod.Status.PodIP
		if ip == "" {
			log.Info("proxy instance pod found but waiting until it has an IP address", "pod", pod.Name)
			return ctrl.Result{Requeue: true}, nil
		}

		// FIXME - super gross hack to for testing, will remove this soon
		commandCtx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cmd := exec.CommandContext(commandCtx, "kubectl", "-n", pod.Namespace, "port-forward", fmt.Sprintf("pod/%s", pod.Name), "8444:8444")
		if v := os.Getenv(controllers.ExternalCtrlEnv); v == "true" {
			ip = "127.0.0.1"
			go func() {
				cmd.Run()
			}()
			time.Sleep(time.Second * 1)
		}
		// FIXME

		url, err := url.Parse(fmt.Sprintf("https://%s:8444", ip))
		if err != nil {
			return ctrl.Result{}, err
		}

		// TODO: add mTLS
		transport := http.DefaultTransport
		transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		httpc := http.Client{Transport: transport}
		kongc, err := kong.NewClient(kong.String(url.String()), &httpc)
		if err != nil {
			return ctrl.Result{}, err
		}

		status, err := kongc.Status(ctx)
		if err != nil {
			return ctrl.Result{}, err
		}
		log.Info("proxy instance in pod returned status successfully", "pod", pod.Name, "connections_accepted", status.Server.ConnectionsAccepted)

		// FIXME - super gross testing hack cleanup, will remove this soon
		cancel()
		time.Sleep(time.Second * 1)
		// FIXME
	}

	// TODO: parse configSecret into Kong configuration (using KIC libs)

	// TODO: post configuration updates to the Kong Admin API of each proxy

	return ctrl.Result{}, nil
}
