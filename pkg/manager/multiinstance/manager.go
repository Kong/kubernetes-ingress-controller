package multiinstance

import (
	"context"
	"net/http"
	"runtime/pprof"
	"sync"

	"github.com/go-logr/logr"

	"github.com/kong/kubernetes-ingress-controller/v3/pkg/manager"
)

const (
	// SchedulingQueueSize is the size of the scheduling queue for manager.Manager instances. It should be large enough
	// to handle all reasonable cases of manager.Manager instances being scheduled at the same time.
	SchedulingQueueSize = 100
)

// ManagerInstance is an interface that represents a single instance of a manager.Manager, exposing only the methods
// needed by the multi-instance manager.
type ManagerInstance interface {
	ID() manager.ID
	Start(context.Context) error
	IsReady() error
	DiagnosticsHandler() http.Handler
}

// DiagnosticsExposer is an interface that represents an object that can expose diagnostics data of manager.Manager
// instances.
type DiagnosticsExposer interface {
	// RegisterInstance registers a new manager.Manager instance with the diagnostics exposer.
	RegisterInstance(manager.ID, http.Handler)

	// UnregisterInstance unregisters a manager.Manager instance with the diagnostics exposer.
	UnregisterInstance(manager.ID)
}

// Manager is able to dynamically run multiple instances of manager.Manager and manage their lifecycle.
// It is responsible for things like:
// - Making sure there's only one instance of a manager.Manager with a given ID.
// - Starting and stopping manager.Manager instances as needed.
// - Registering instances' diagnostic handlers in a DiagnosticsExposer when configured.
type Manager struct {
	logger logr.Logger

	instances          map[manager.ID]*instance
	instancesLock      sync.RWMutex
	schedulingQueue    chan manager.ID
	diagnosticsExposer DiagnosticsExposer
}

// ManagerOption is a functional option that can be used to configure a new multi-instance manager.
type ManagerOption func(*Manager)

func WithDiagnosticsExposer(exposer DiagnosticsExposer) ManagerOption {
	return func(m *Manager) {
		m.diagnosticsExposer = exposer
	}
}

// NewManager creates a new multi-instance manager.
func NewManager(logger logr.Logger, opts ...ManagerOption) *Manager {
	m := &Manager{
		logger:          logger,
		instances:       make(map[manager.ID]*instance),
		schedulingQueue: make(chan manager.ID, SchedulingQueueSize),
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

// Start starts the multi-instance manager and blocks until the context is canceled. It should only be called once.
func (m *Manager) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case instanceID := <-m.schedulingQueue:
			go m.runInstance(ctx, instanceID)
		}
	}
}

// ScheduleInstance adds a new manager.Manager instance to the multi-instance manager and starts it immediately in a
// separate goroutine. If an instance with the same ID already exists, it returns a InstanceWithIDAlreadyScheduledError error.
func (m *Manager) ScheduleInstance(in ManagerInstance) error {
	m.logger.Info("Scheduling instance", "instanceID", in.ID())

	m.instancesLock.Lock()
	defer m.instancesLock.Unlock()

	if _, exists := m.instances[in.ID()]; exists {
		return NewInstanceWithIDAlreadyScheduledError(in.ID())
	}

	// Keep track of the instance, but do not start it from here.
	m.instances[in.ID()] = newInstance(in, m.logger)

	// Send a signal to the scheduling channel to start the instance.
	m.schedulingQueue <- in.ID()

	return nil
}

// StopInstance stops a manager.Manager instance with the given ID. If no instance with the given ID exists, it returns
// a InstanceNotFoundError error.
func (m *Manager) StopInstance(instanceID manager.ID) error {
	m.logger.Info("Stopping instance", "instanceID", instanceID)

	m.instancesLock.Lock()
	defer m.instancesLock.Unlock()

	in, exists := m.instances[instanceID]
	if !exists {
		return NewInstanceNotFoundError(instanceID)
	}

	// If diagnostics are enabled, unregister the instance from the diagnostics exposer.
	if m.diagnosticsExposer != nil {
		m.diagnosticsExposer.UnregisterInstance(instanceID)
	}

	// Send a signal to the instance to stop and let the running goroutine handle the cleanup.
	in.Stop()

	return nil
}

// IsInstanceReady checks if a manager.Manager instance with the given ID is ready. If no instance with the given ID
// exists, it returns a InstanceNotFoundError error.
func (m *Manager) IsInstanceReady(id manager.ID) error {
	m.instancesLock.RLock()
	defer m.instancesLock.RUnlock()
	in, ok := m.instances[id]
	if !ok {
		return NewInstanceNotFoundError(id)
	}
	return in.IsReady()
}

func (m *Manager) runInstance(ctx context.Context, instanceID manager.ID) {
	m.instancesLock.RLock()
	in, exists := m.instances[instanceID]
	m.instancesLock.RUnlock()

	if !exists {
		// Instance was removed while waiting for the lock.
		m.logger.WithValues("instanceID", instanceID).Info("Instance was removed while waiting for the lock")
		return
	}

	m.logger.Info("Starting instance", "instanceID", instanceID)

	// Wrap with pprof.Do to add instanceID to the pprof labels. That will make it easier to identify which instance
	// is responsible for the CPU consumption.
	pprof.Do(ctx, pprof.Labels("instanceID", instanceID.String()), func(ctx context.Context) {
		go in.Start(ctx)
	})

	// If diagnostics are enabled, register the instance with the diagnostics exposer.
	if m.diagnosticsExposer != nil {
		m.diagnosticsExposer.RegisterInstance(instanceID, in.DiagnosticsHandler())
	}

	// Wait for the instance to stop or the parent context be done.
	select {
	case <-in.StopChannel():
		m.logger.Info("Instance stopped, removing it from managed instances", "instanceID", instanceID)
		m.instancesLock.Lock()
		delete(m.instances, instanceID)
		m.instancesLock.Unlock()
	case <-ctx.Done():
	}
}
