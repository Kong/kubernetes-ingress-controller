package multiinstance_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-logr/logr/testr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"github.com/kong/kubernetes-ingress-controller/v3/pkg/manager"
	"github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/multiinstance"
)

const (
	waitTime = time.Second
	tickTime = time.Millisecond * 10
)

type MockInstance struct {
	id                 manager.ID
	returnErrOnRun     error
	wasStarted         atomic.Bool
	wasContextCanceled atomic.Bool
}

func newMockInstance(id manager.ID) *MockInstance {
	return &MockInstance{
		id: id,
	}
}

func (m *MockInstance) ID() manager.ID {
	return m.id
}

func (m *MockInstance) Run(ctx context.Context) error {
	m.wasStarted.Store(true)

	go func() {
		<-ctx.Done()
		m.wasContextCanceled.Store(true)
	}()

	return m.returnErrOnRun
}

func (m *MockInstance) IsReady() error {
	return nil
}

func TestManager(t *testing.T) {
	// At the end of the test, assert that there are no goroutines leaked.
	t.Cleanup(func() { goleak.VerifyNone(t) })

	// Such context will be canceled just before the test ends (Cleanups are run)
	// so we can ensure all goroutines are cleaned up.
	ctx := t.Context()

	multiManager := multiinstance.NewManager(testr.New(t))

	mockInstance1 := newMockInstance(manager.NewRandomID())
	mockInstance2 := newMockInstance(manager.NewRandomID())

	t.Run("can schedule instances before starting the manager", func(t *testing.T) {
		err := multiManager.ScheduleInstance(mockInstance1)
		require.NoError(t, err)

		err = multiManager.ScheduleInstance(mockInstance2)
		require.NoError(t, err)

		require.False(t, mockInstance1.wasStarted.Load(), "instance should not have been started yet as the manager is not running")
	})

	t.Run("schedling an instance with the same ID should fail", func(t *testing.T) {
		err := multiManager.ScheduleInstance(mockInstance1)
		require.Error(t, err)
		require.IsType(t, multiinstance.InstanceWithIDAlreadyScheduledError{}, err)
	})

	managerRunning := make(chan struct{})
	t.Run("can run the manager", func(t *testing.T) {
		go func() {
			close(managerRunning)
			require.NoError(t, multiManager.Run(ctx))
		}()
	})

	t.Run("can schedule instances after starting the manager", func(t *testing.T) {
		<-managerRunning // Wait for the manager to start.

		mockInstance3 := newMockInstance(manager.NewRandomID())
		err := multiManager.ScheduleInstance(mockInstance3)
		require.NoError(t, err)

		require.EventuallyWithT(t, func(t *assert.CollectT) {
			assert.True(t, mockInstance1.wasStarted.Load())
		}, waitTime, tickTime)
	})

	t.Run("can stop an instance", func(t *testing.T) {
		err := multiManager.StopInstance(mockInstance1.ID())
		require.NoError(t, err)

		require.EventuallyWithT(t, func(t *assert.CollectT) {
			assert.True(t, mockInstance1.wasContextCanceled.Load())
		}, waitTime, tickTime)
	})

	t.Run("can inspect instance readiness", func(t *testing.T) {
		err := multiManager.IsInstanceReady(mockInstance2.ID())
		require.NoError(t, err)

		err = multiManager.IsInstanceReady(mockInstance1.ID())
		require.Error(t, err)
		require.IsType(t, multiinstance.InstanceNotFoundError{}, err)
	})
}
