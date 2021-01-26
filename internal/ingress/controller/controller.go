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
	"io/ioutil"
	"sync"
	"sync/atomic"
	"time"

	"github.com/blang/semver"
	"github.com/eapache/channels"
	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/controller/parser"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/election"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/status"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/store"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/task"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/utils"
	configClientSet "github.com/kong/kubernetes-ingress-controller/pkg/client/configuration/clientset/versioned"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	networking "k8s.io/api/networking/v1beta1"
	clientset "k8s.io/client-go/kubernetes"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/flowcontrol"
	knativeClientSet "knative.dev/networking/pkg/client/clientset/versioned"
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
	KongCustomEntitiesSecret string

	KubeClient       clientset.Interface
	KongConfigClient configClientSet.Interface
	KnativeClient    knativeClientSet.Interface

	ResyncPeriod      time.Duration
	SyncRateLimit     float32
	EnableReverseSync bool

	Namespace string

	IngressClass string

	// optional
	PublishService       string
	PublishStatusAddress string

	UpdateStatus           bool
	UpdateStatusOnShutdown bool
	ElectionID             string

	IngressAPI utils.IngressAPI

	EnableKnativeIngressSupport bool

	Logger     logrus.FieldLogger
	LogLevel   string
	DumpConfig string
}

// sync collects all the pieces required to assemble the configuration file and
// then sends the content to the backend (OnUpdate) receiving the populated
// template as response reloading the backend if is required.
func (n *KongController) syncIngress(interface{}) error {
	ctx := context.Background()
	n.syncRateLimiter.Accept()

	if n.syncQueue.IsShuttingDown() {
		return nil
	}

	// If in-memory mode, each Kong instance runs with its own controller
	if !n.cfg.Kong.InMemory &&
		!n.elector.IsLeader() {
		n.Logger.Debugf("node is a follower, skipping update")
		return nil
	}

	n.Logger.Infof("syncing configuration")
	state, err := parser.Build(n.Logger.WithField("component", "store"), n.store)
	if err != nil {
		return fmt.Errorf("error building kong state: %w", err)
	}
	err = n.OnUpdate(ctx, state)
	if err != nil {
		n.Logger.Errorf("failed to update kong configuration: %v", err)
		return err
	}

	return nil
}

// NewKongController creates a new Ingress controller.
func NewKongController(ctx context.Context,
	config *Configuration,
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

		Logger: config.Logger,
	}

	n.store = store
	n.syncQueue = task.NewTaskQueue(n.syncIngress,
		config.Logger.WithField("component", "sync-queue"))

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
		Logger: config.Logger.WithField("component", "elector"),
	}

	if config.UpdateStatus {
		var err error
		n.syncStatus, err = status.NewStatusSyncer(ctx, status.Config{
			CoreClient:             config.KubeClient,
			KongConfigClient:       config.KongConfigClient,
			KnativeClient:          config.KnativeClient,
			PublishService:         config.PublishService,
			PublishStatusAddress:   config.PublishStatusAddress,
			IngressLister:          n.store,
			UpdateStatusOnShutdown: config.UpdateStatusOnShutdown,
			IngressAPI:             config.IngressAPI,
			OnStartedLeading: func() {
				// force a sync
				n.syncQueue.Enqueue(&networking.Ingress{})
			},
			Logger: n.Logger.WithField("component", "status-syncer"),
		})
		if err != nil {
			return nil, fmt.Errorf("initializing status syncer: %v", err)
		}
		electionConfig.Callbacks = n.syncStatus.Callbacks()
	} else {

		n.Logger.Warnf("ingress status updates is disabled, flag --update-status=false was specified")
	}

	n.elector = election.NewElector(ctx, electionConfig)

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

	// backgroundGroup tracks running background goroutines, on which we'll wait to stop in Stop.
	backgroundGroup errgroup.Group

	runningConfigHash []byte
	lastConfig        []byte

	isShuttingDown uint32

	store store.Storer

	PluginSchemaStore PluginSchemaStore

	Logger logrus.FieldLogger

	tmpDir string
}

// Start starts a new master process running in foreground, blocking until the next call to
// Stop.
func (n *KongController) Start() {
	n.Logger.Debugf("startin up controller")
	var err error
	n.tmpDir, err = ioutil.TempDir("", "controller")
	if err != nil {
		panic(fmt.Errorf("failed to create a temporary working directory: %v", err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		<-n.stopCh
		cancel()
	}()

	var group errgroup.Group

	group.Go(func() error {
		n.elector.Run(ctx)
		return nil
	})

	if n.syncStatus != nil {
		group.Go(func() error {
			n.syncStatus.Run()
			return nil
		})
	}

	group.Go(func() error {
		n.syncQueue.Run(time.Second, n.stopCh)
		return nil
	})
	// Force initial sync.
	n.syncQueue.Enqueue(&networking.Ingress{})

	for {
		select {
		case event := <-n.updateCh.Out():
			if v := atomic.LoadUint32(&n.isShuttingDown); v != 0 {
				return
			}
			if evt, ok := event.(Event); ok {
				n.Logger.WithField("event_type", evt.Type).Debugf("event received")
				n.syncQueue.Enqueue(evt.Obj)
			} else {
				n.Logger.WithField("event_type", evt.Type).Errorf("invalid event received")
			}
		case <-n.stopCh:
			return
		}
	}
}

// Stop stops the master process gracefully.
func (n *KongController) Stop() error {
	atomic.StoreUint32(&n.isShuttingDown, 1)

	n.stopLock.Lock()
	defer n.stopLock.Unlock()

	// Only try draining the workqueue if we haven't already.
	if n.syncQueue.IsShuttingDown() {
		return fmt.Errorf("shutdown already in progress")
	}

	if n.syncStatus != nil {
		// Finish writing to any Ingress objects before giving up leadership.
		n.syncStatus.Shutdown(n.elector.IsLeader())
	}
	n.Logger.Infof("shutting down controller queues")
	n.syncQueue.Shutdown()
	// Closing the stop channel will cause us to give up leadership.
	close(n.stopCh)
	n.Logger.Infof("awaiting completion of shutdown procedures")
	if err := n.backgroundGroup.Wait(); err != nil {
		n.Logger.Errorf("background controller task failed: %v", err)
	}

	return nil
}
