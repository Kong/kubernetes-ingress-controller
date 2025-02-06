package multiinstance_test

import (
	"context"
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

func TestManager_Scheduling(t *testing.T) {
	onCleanupVerifyThereAreNoLeakedGoroutines(t)

	// Create a context that will be canceled when the test ends so we can ensure all goroutines are cleaned up.
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

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

	t.Run("scheduling an instance with the same ID should fail", func(t *testing.T) {
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

func TestManager_WithDiagnosticsExposer(t *testing.T) {
	onCleanupVerifyThereAreNoLeakedGoroutines(t)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Log("Configuring a manager with a diagnostics exposer")
	diagnosticsExposer := newMockDiagnosticsExposer()
	multiManager := multiinstance.NewManager(testr.New(t), multiinstance.WithDiagnosticsExposer(diagnosticsExposer))

	managerRunning := make(chan struct{})
	go func() {
		close(managerRunning)
		require.NoError(t, multiManager.Run(ctx))
	}()
	<-managerRunning // Wait for the manager to start.

	instanceID1 := manager.NewRandomID()
	instanceID2 := manager.NewRandomID()

	t.Log("Scheduling first instance")
	err := multiManager.ScheduleInstance(newMockInstance(instanceID1))
	require.NoError(t, err)

	t.Log("Expecting the diagnostics exposer to have the first instance registered")
	require.EventuallyWithT(t, func(t *assert.CollectT) {
		require.Contains(t, diagnosticsExposer.RegisteredInstances(), instanceID1)
		require.NotContains(t, diagnosticsExposer.RegisteredInstances(), instanceID2)
	}, waitTime, tickTime)

	t.Log("Scheduling second instance")
	err = multiManager.ScheduleInstance(newMockInstance(instanceID2))
	require.NoError(t, err)

	t.Log("Expecting the diagnostics exposer to have both instances registered")
	require.EventuallyWithT(t, func(t *assert.CollectT) {
		require.ElementsMatch(t, diagnosticsExposer.RegisteredInstances(), []manager.ID{instanceID1, instanceID2})
	}, waitTime, tickTime)

	t.Log("Stopping first instance")
	err = multiManager.StopInstance(instanceID1)
	require.NoError(t, err)

	t.Log("Expecting the diagnostics exposer to have only the second instance registered")
	require.EventuallyWithT(t, func(t *assert.CollectT) {
		require.Contains(t, diagnosticsExposer.RegisteredInstances(), instanceID2)
		require.NotContains(t, diagnosticsExposer.RegisteredInstances(), instanceID1)
	}, waitTime, tickTime)

	t.Log("Stopping second instance")
	err = multiManager.StopInstance(instanceID2)
	require.NoError(t, err)

	t.Log("Expecting the diagnostics exposer to have no instances registered")
	require.EventuallyWithT(t, func(t *assert.CollectT) {
		require.Empty(t, diagnosticsExposer.RegisteredInstances())
	}, waitTime, tickTime)
}

// onCleanupVerifyThereAreNoLeakedGoroutines is a helper function that sets up a cleanup function to verify there are no
// leaked goroutines at the end of the test.
func onCleanupVerifyThereAreNoLeakedGoroutines(t *testing.T) {
	t.Helper()
	t.Cleanup(func() { goleak.VerifyNone(t) })
}
