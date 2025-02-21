package multiinstance_test

import (
	"context"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/samber/lo"

	"github.com/kong/kubernetes-ingress-controller/v3/pkg/manager"
)

// mockInstance is a mock implementation of multiinstance.ManagerInstance.
type mockInstance struct {
	id                 manager.ID
	returnErrOnRun     error
	wasStarted         atomic.Bool
	wasContextCanceled atomic.Bool
}

func newMockInstance(id manager.ID) *mockInstance {
	return &mockInstance{
		id: id,
	}
}

func (m *mockInstance) ID() manager.ID {
	return m.id
}

func (m *mockInstance) Start(ctx context.Context) error {
	m.wasStarted.Store(true)

	go func() {
		<-ctx.Done()
		m.wasContextCanceled.Store(true)
	}()

	return m.returnErrOnRun
}

func (m *mockInstance) IsReady() error {
	return nil
}

func (m *mockInstance) DiagnosticsHandler() http.Handler {
	return nil
}

// mockDiagnosticsExposer is a mock implementation of multiinstance.DiagnosticsExposer.
type mockDiagnosticsExposer struct {
	registeredInstances map[manager.ID]struct{}
	lock                sync.Mutex
}

func newMockDiagnosticsExposer() *mockDiagnosticsExposer {
	return &mockDiagnosticsExposer{
		registeredInstances: make(map[manager.ID]struct{}),
	}
}

func (m *mockDiagnosticsExposer) RegisterInstance(id manager.ID, _ http.Handler) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.registeredInstances[id] = struct{}{}
}

func (m *mockDiagnosticsExposer) UnregisterInstance(id manager.ID) {
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.registeredInstances, id)
}

func (m *mockDiagnosticsExposer) RegisteredInstances() []manager.ID {
	m.lock.Lock()
	defer m.lock.Unlock()

	return lo.Keys(m.registeredInstances)
}
