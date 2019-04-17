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
	"fmt"
	"sync"
	"time"

	"github.com/golang/glog"

	"github.com/eapache/channels"
	extensions "k8s.io/api/extensions/v1beta1"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/flowcontrol"

	"github.com/kong/kubernetes-ingress-controller/internal/ingress/annotations/class"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/controller/parser"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/controller/store"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/election"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/status"
	"github.com/kong/kubernetes-ingress-controller/internal/task"
)

// NewKongController creates a new NGINX Ingress controller.
// If the environment variable NGINX_BINARY exists it will be used
// as source for nginx commands
func NewKongController(config *Configuration) (*KongController, error) {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{
		Interface: config.KubeClient.CoreV1().Events(config.Namespace),
	})

	n := &KongController{
		cfg:             config,
		syncRateLimiter: flowcontrol.NewTokenBucketRateLimiter(config.SyncRateLimit, 1),

		stopCh:   make(chan struct{}),
		updateCh: channels.NewRingChannel(1024),

		stopLock: &sync.Mutex{},
		PluginSchemaStore: *NewPluginSchemaStore(config.Kong.Client,
			config.Kong.URL),
	}

	n.store = store.New(
		config.Namespace,
		config.ResyncPeriod,
		config.KubeClient,
		config.KubeConf,
		n.updateCh)

	n.parser = parser.New(n.store)
	n.syncQueue = task.NewTaskQueue(n.syncIngress)

	ingressClass := class.DefaultClass
	if class.IngressClass != "" {
		ingressClass = class.IngressClass
	}
	electionID := config.ElectionID + "-" + ingressClass

	electionConfig := election.Config{
		Client:     config.KubeClient,
		ElectionID: electionID,
	}

	if config.UpdateStatus {
		n.syncStatus = status.NewStatusSyncer(status.Config{
			Client:                 config.KubeClient,
			PublishService:         config.PublishService,
			PublishStatusAddress:   config.PublishStatusAddress,
			IngressLister:          n.store,
			ElectionID:             electionID,
			IngressClass:           class.IngressClass,
			DefaultIngressClass:    class.DefaultClass,
			UpdateStatusOnShutdown: config.UpdateStatusOnShutdown,
			OnStartedLeading: func() {
				// force a sync
				n.syncQueue.Enqueue(&extensions.Ingress{})
			},
		})
		electionConfig.Callbacks = n.syncStatus.Callbacks()
	} else {
		glog.Warning("Update of ingress status is disabled (flag --update-status=false was specified)")
	}

	n.elector = election.NewElector(electionConfig)

	return n, nil
}

// KongController ...
type KongController struct {
	cfg *Configuration

	syncQueue *task.Queue

	syncStatus status.Sync

	elector election.Elector

	syncRateLimiter flowcontrol.RateLimiter

	// stopLock is used to enforce only a single call to Stop is active.
	// Needed because we allow stopping through an http endpoint and
	// allowing concurrent stoppers leads to stack traces.
	stopLock *sync.Mutex

	stopCh   chan struct{}
	updateCh *channels.RingChannel

	// runningConfig contains the running configuration in the Backend
	runningConfig *parser.KongState

	isShuttingDown bool

	store store.Storer

	parser parser.Parser

	PluginSchemaStore PluginSchemaStore
}

// Start start a new NGINX master process running in foreground.
func (n *KongController) Start() {
	glog.Infof("starting Ingress controller")

	n.store.Run(n.stopCh)

	go n.elector.Run()

	if n.syncStatus != nil {
		go n.syncStatus.Run()
	}

	go n.syncQueue.Run(time.Second, n.stopCh)
	// force initial sync
	n.syncQueue.Enqueue(&extensions.Ingress{})

	for {
		select {
		case event := <-n.updateCh.Out():
			if n.isShuttingDown {
				break
			}
			if evt, ok := event.(store.Event); ok {
				glog.V(3).Infof("Event %v received - object %v", evt.Type, evt.Obj)
				n.syncQueue.Enqueue(evt.Obj)
			} else {
				glog.Warningf("unexpected event type received %T", event)
			}
		case <-n.stopCh:
			break
		}
	}
}

// Stop gracefully stops the NGINX master process.
func (n *KongController) Stop() error {
	n.isShuttingDown = true

	n.stopLock.Lock()
	defer n.stopLock.Unlock()

	// Only try draining the workqueue if we haven't already.
	if n.syncQueue.IsShuttingDown() {
		return fmt.Errorf("shutdown already in progress")
	}

	glog.Infof("shutting down controller queues")
	close(n.stopCh)
	go n.syncQueue.Shutdown()
	if n.syncStatus != nil {
		n.syncStatus.Shutdown(n.elector.IsLeader())
	}

	return nil
}
