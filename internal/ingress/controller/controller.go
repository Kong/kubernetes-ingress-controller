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
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/blang/semver"
	"github.com/eapache/channels"
	"github.com/golang/glog"
	"github.com/hbagdi/go-kong/kong"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/internal/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/controller/parser"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/election"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/status"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/store"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/task"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	networking "k8s.io/api/networking/v1beta1"
	clientset "k8s.io/client-go/kubernetes"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/flowcontrol"
)

// Kong Represents a Kong client and connection information
type Kong struct {
	URL        string
	FilterTags []string
	// Headers are injected into every request to Kong's Admin API
	// to help with authorization/authentication.
	Client *kong.Client

	InMemory      bool
	HasTagSupport bool
	Enterprise    bool

	Version semver.Version

	Concurrency int
}

// Configuration contains all the settings required by an Ingress controller
type Configuration struct {
	Kong

	KubeClient clientset.Interface

	ResyncPeriod  time.Duration
	SyncRateLimit float32

	Namespace string

	IngressClass string

	// optional
	PublishService       string
	PublishStatusAddress string

	UpdateStatus           bool
	UpdateStatusOnShutdown bool
	ElectionID             string

	UseNetworkingV1beta1 bool
}

// sync collects all the pieces required to assemble the configuration file and
// then sends the content to the backend (OnUpdate) receiving the populated
// template as response reloading the backend if is required.
func (n *KongController) syncIngress(interface{}) error {
	n.syncRateLimiter.Accept()

	if n.syncQueue.IsShuttingDown() {
		return nil
	}

	// If in-memory mode, each Kong instance runs with it's own controller
	if !n.cfg.Kong.InMemory &&
		!n.elector.IsLeader() {
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

	return nil
}

// NewKongController creates a new NGINX Ingress controller.
// If the environment variable NGINX_BINARY exists it will be used
// as source for nginx commands
func NewKongController(config *Configuration,
	updateCh *channels.RingChannel,
	store store.Storer) (*KongController, error) {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{
		Interface: config.KubeClient.CoreV1().Events(config.Namespace),
	})

	n := &KongController{
		cfg:             config,
		syncRateLimiter: flowcontrol.NewTokenBucketRateLimiter(config.SyncRateLimit, 1),

		stopCh:   make(chan struct{}),
		updateCh: updateCh,

		stopLock:          &sync.Mutex{},
		PluginSchemaStore: *NewPluginSchemaStore(config.Kong.Client),
	}

	n.store = store
	n.parser = parser.New(n.store)
	n.syncQueue = task.NewTaskQueue(n.syncIngress)

	electionID := config.ElectionID + "-" + config.IngressClass

	electionConfig := election.Config{
		Client:     config.KubeClient,
		ElectionID: electionID,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
			},
			OnStoppedLeading: func() {
			},
			OnNewLeader: func(identity string) {
			},
		},
	}

	if config.UpdateStatus {
		n.syncStatus = status.NewStatusSyncer(status.Config{
			Client:                 config.KubeClient,
			PublishService:         config.PublishService,
			PublishStatusAddress:   config.PublishStatusAddress,
			IngressLister:          n.store,
			ElectionID:             electionID,
			UpdateStatusOnShutdown: config.UpdateStatusOnShutdown,
			UseNetworkingV1beta1:   config.UseNetworkingV1beta1,
			OnStartedLeading: func() {
				// force a sync
				n.syncQueue.Enqueue(&networking.Ingress{})
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

	runningConfigHash [32]byte

	isShuttingDown bool

	store store.Storer

	parser parser.Parser

	PluginSchemaStore PluginSchemaStore
}

// Start start a new NGINX master process running in foreground.
func (n *KongController) Start() {
	glog.Infof("starting Ingress controller")

	go n.elector.Run()

	if n.syncStatus != nil {
		go n.syncStatus.Run()
	}

	go n.syncQueue.Run(time.Second, n.stopCh)
	// force initial sync
	n.syncQueue.Enqueue(&networking.Ingress{})

	for {
		select {
		case event := <-n.updateCh.Out():
			if n.isShuttingDown {
				break
			}
			if evt, ok := event.(Event); ok {
				glog.V(3).Infof("Event %v received - object %v", evt.Type, evt.Obj)
				n.syncQueue.Enqueue(evt.Obj)
				// TODO retry for ephermal error conditions
				// This function is called outside the task queue because event
				// information is currently shielded from the sync function.
				// Sync function syncs everything, no matter what the event is
				err := n.handleBasicAuthUpdates(evt)
				if err != nil {
					glog.Errorf("error handling basic-auth update: %v", err)
				}
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

// handleBasicAuthUpdates updates basic-auth password field
// in Kong whenever it is changed.
// Kong hashes basic-auth passwords in DB and API responses once created.
// Due to this reason, one can't perform a 'diff' with them.
// This function filters for basic-auth password changes and applies them
// to Kong as they happen.
func (n *KongController) handleBasicAuthUpdates(event Event) error {
	if !n.elector.IsLeader() {
		return nil
	}
	if n.cfg.Kong.InMemory {
		return nil
	}
	if event.Type != UpdateEvent {
		return nil
	}
	newCred, ok := event.Obj.(*configurationv1.KongCredential)
	if !ok {
		return nil
	}
	if newCred.ConsumerRef == "" {
		return nil
	}
	oldCred, ok := event.Old.(*configurationv1.KongCredential)
	if !ok {
		return nil
	}
	// if the credential type was changed, then the basic-auth
	// credential will be either created or deleted by decK
	if oldCred.Type != "basic-auth" && newCred.Type != "basic-auth" {
		return nil
	}
	oldPassword := oldCred.Config["password"]
	newPassword := newCred.Config["password"]
	if oldPassword == newPassword {
		return nil
	}
	// there was an update to a basic-auth credential and the password
	// has changed, sync it
	var cred kong.BasicAuth
	decoder, err := mapstructure.NewDecoder(
		&mapstructure.DecoderConfig{TagName: "json",
			Result: &cred,
		})
	if err != nil {
		return errors.Wrap(err, "failed to create a decoder")
	}
	err = decoder.Decode(newCred.Config)
	if err != nil {
		return errors.Wrapf(err, "error decoding credential '%v/%v'",
			newCred.Namespace, newCred.Name)
	}

	kongConsumer, err := n.store.GetKongConsumer(newCred.Namespace,
		newCred.ConsumerRef)
	if err != nil {
		return errors.Wrapf(err, "error searching for consumer '%v/%v'",
			newCred.Namespace, newCred.ConsumerRef)
	}
	username := kongConsumer.Username
	client := n.cfg.Kong.Client

	// find the ID of the cred from Kong
	outdatedCred, err := client.BasicAuths.Get(nil, &username, cred.Username)
	if err != nil {
		return errors.Wrap(err, "fetching basic-auth credential")
	}
	cred.ID = outdatedCred.ID
	// update it
	_, err = client.BasicAuths.Create(nil, &username, &cred)
	if err != nil {
		return errors.Wrap(err, "updating basic-auth credential")
	}
	return nil
}
