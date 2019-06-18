/*
Copyright 2017 The Kubernetes Authors.

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

package main

import (
	"errors"
	"flag"
	"os"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/spf13/pflag"

	apiv1 "k8s.io/api/core/v1"

	"github.com/kong/kubernetes-ingress-controller/internal/ingress/annotations"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/controller"
)

type headers []string

func (h *headers) String() string {
	return "my string representation"
}

func (h *headers) Set(value string) error {
	if len(strings.Split(value, ":")) < 2 {
		return errors.New("header should be of form key:value")
	}
	*h = append(*h, value)
	return nil
}

func parseFlags() (bool, *controller.Configuration, error) {
	var (
		flags = pflag.NewFlagSet("", pflag.ExitOnError)

		apiserverHost = flags.String("apiserver-host", "",
			`The address of the Kubernetes Apiserver to connect to in the format of 
protocol://address:port, e.g., "http://localhost:8080.
If not specified, the assumption is that the binary runs inside a 
Kubernetes cluster and local discovery is attempted.`)
		kubeConfigFile = flags.String("kubeconfig", "", "Path to kubeconfig file with authorization and master location information.")

		ingressClass = flags.String("ingress-class", annotations.DefaultIngressClass,
			`Name of the ingress class to route through this controller.`)

		publishSvc = flags.String("publish-service", "",
			`Service fronting the ingress controllers. Takes the form namespace/name.
		The controller will set the endpoint records on the ingress objects to reflect those on the service.`)

		resyncPeriod = flags.Duration("sync-period", 600*time.Second,
			`Relist and confirm cloud resources this often. Default is 10 minutes`)

		watchNamespace = flags.String("watch-namespace", apiv1.NamespaceAll,
			`Namespace to watch for Ingress. Default is to watch all namespaces`)

		profiling = flags.Bool("profiling", true, `Enable profiling via web interface host:port/debug/pprof/`)

		updateStatus = flags.Bool("update-status", true, `Indicates if the
		ingress controller should update the Ingress status IP/hostname. Default is true`)

		electionID = flags.String("election-id", "ingress-controller-leader", `Election id to use for status update.`)

		updateStatusOnShutdown = flags.Bool("update-status-on-shutdown", true,
			`Indicates if the ingress controller should update the Ingress status 
IP/hostname when the controller is being stopped. Default is true`)

		showVersion = flags.Bool("version", false,
			`Shows release information about the Kong Ingress controller`)

		syncRateLimit = flags.Float32("sync-rate-limit", 0.3,
			`Define the sync frequency upper limit`)

		publishStatusAddress = flags.String("publish-status-address", "",
			`User customized address to be set in the status of ingress resources.
The controller will set the endpoint records on the ingress using this address.`)

		kongURL = flags.String("kong-url", "http://localhost:8001",
			"The address of the Kong Admin URL to connect to in the format of protocol://address:port")

		kongTLSSkipVerify = flag.Bool("admin-tls-skip-verify", false,
			"Disable verification of TLS certificate of Kong's Admin endpoint.")
		kongTLSServerName = flag.String("admin-tls-server-name", "",
			"SNI name to use to verify the certificate presented by Kong in TLS.")
		kongCACert = flag.String("admin-ca-cert-file", "",
			"Path to PEM-encoded CA certificate file to verify the Kong's Admin SSL certificate.")
		workspace = flag.String("kong-workspace", "",
			"Workspace in Kong Enterprise to be configured")

		kongHeaders headers
	)

	flag.Var(&kongHeaders, "admin-header",
		"add a header (key:value) to every Admin API call, flag can be used multiple times")

	flag.Set("logtostderr", "true")

	flags.AddGoFlagSet(flag.CommandLine)
	flags.Parse(os.Args)

	// Workaround for this issue:
	// https://github.com/kubernetes/kubernetes/issues/17162
	flag.CommandLine.Parse([]string{})

	pflag.VisitAll(func(flag *pflag.Flag) {
		glog.V(2).Infof("FLAG: --%s=%q", flag.Name, flag.Value)
	})

	if *showVersion {
		return true, nil, nil
	}

	config := &controller.Configuration{
		Kong: controller.Kong{
			URL:       *kongURL,
			Headers:   kongHeaders,
			Workspace: *workspace,

			TLSServerName: *kongTLSServerName,
			TLSSkipVerify: *kongTLSSkipVerify,
			CACert:        *kongCACert,
		},
		APIServerHost:          *apiserverHost,
		KubeConfigFile:         *kubeConfigFile,
		UpdateStatus:           *updateStatus,
		ElectionID:             *electionID,
		EnableProfiling:        *profiling,
		ResyncPeriod:           *resyncPeriod,
		IngressClass:           *ingressClass,
		Namespace:              *watchNamespace,
		PublishService:         *publishSvc,
		PublishStatusAddress:   *publishStatusAddress,
		UpdateStatusOnShutdown: *updateStatusOnShutdown,
		SyncRateLimit:          *syncRateLimit,
	}

	return false, config, nil
}
