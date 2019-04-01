/*
Copyright 2015 The Kubernetes Authors.

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

package controller

import (
	"sort"
	"time"

	"github.com/golang/glog"
	"github.com/pkg/errors"

	apiv1 "k8s.io/api/core/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/hbagdi/go-kong/kong"
)

const (
	defUpstreamName = "upstream-default-backend"
	defServerName   = "_"
	rootLocation    = "/"
)

// Kong Represents a Kong client and connection information
type Kong struct {
	URL string
	// Headers are injected into every request to Kong's Admin API
	// to help with authorization/authentication.
	Headers []string
	Client  *kong.Client

	TLSSkipVerify bool
	TLSServerName string
	CACert        string
}

// Configuration contains all the settings required by an Ingress controller
type Configuration struct {
	Kong

	APIServerHost  string
	KubeConfigFile string
	KubeClient     clientset.Interface
	KubeConf       *rest.Config

	ResyncPeriod time.Duration

	DefaultService string

	Namespace string

	DefaultHealthzURL     string
	DefaultSSLCertificate string

	// optional
	PublishService       string
	PublishStatusAddress string

	UpdateStatus           bool
	ElectionID             string
	UpdateStatusOnShutdown bool

	EnableProfiling bool

	SyncRateLimit float32
}

// GetPublishService returns the configured service used to set ingress status
func (n NGINXController) GetPublishService() *apiv1.Service {
	s, err := n.store.GetService(n.cfg.PublishService)
	if err != nil {
		return nil
	}

	return s
}

// sync collects all the pieces required to assemble the configuration file and
// then sends the content to the backend (OnUpdate) receiving the populated
// template as response reloading the backend if is required.
func (n *NGINXController) syncIngress(interface{}) error {
	n.syncRateLimiter.Accept()

	if n.syncQueue.IsShuttingDown() {
		return nil
	}

	if !n.syncStatus.IsLeader() {
		glog.V(2).Infof("skipping synchronization of configuration because I am not the leader.")
		return nil
	}

	// Sort ingress rules using the ResourceVersion field
	ings := n.store.ListIngresses()
	sort.SliceStable(ings, func(i, j int) bool {
		ir := ings[i].ResourceVersion
		jr := ings[j].ResourceVersion
		return ir < jr
	})

	glog.V(2).Infof("syncing Ingress configuration...")
	state, err := n.parser.Build()
	if err != nil {
		return errors.Wrap(err, "error building kong state")
	}
	err = n.OnUpdate(state)
	if err != nil {
		glog.Errorf("unexpected failure updating Kong configuration: \n%v", err)
		return err
	}
	glog.Info("successfully synced configuration to Kong")

	n.runningConfig = state

	return nil
}
