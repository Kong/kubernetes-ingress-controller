package envtest

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/go-logr/logr/testr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-configuration/pkg/clientset/scheme"

	"github.com/kong/kubernetes-ingress-controller/v3/pkg/manager"
	"github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/multiinstance"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers"
)

func TestMultiInstanceManagerDiagnostics(t *testing.T) {
	t.Parallel()

	const (
		waitTime = 10 * time.Second
		tickTime = 10 * time.Millisecond
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	envcfg := Setup(t, scheme.Scheme)
	diagPort := helpers.GetFreePort(t)

	t.Log("Starting the diagnostics server and the multi-instance manager")
	diagServer := multiinstance.NewDiagnosticsServer(diagPort)
	go func() {
		require.ErrorIs(t, diagServer.Start(ctx), http.ErrServerClosed)
	}()
	multimgr := multiinstance.NewManager(testr.New(t), multiinstance.WithDiagnosticsExposer(diagServer))
	go func() {
		require.NoError(t, multimgr.Run(ctx))
	}()

	t.Log("Setting up two instances of the manager and scheduling them in the multi-instance manager")
	mgrInstance1 := SetupManager(ctx, t, manager.NewRandomID(), envcfg, AdminAPIOptFns(), WithDiagnosticsWithoutServer())
	mgrInstance2 := SetupManager(ctx, t, manager.NewRandomID(), envcfg, AdminAPIOptFns(), WithDiagnosticsWithoutServer())
	require.NoError(t, multimgr.ScheduleInstance(mgrInstance1))
	require.NoError(t, multimgr.ScheduleInstance(mgrInstance2))

	t.Log("Waiting for the diagnostics server to expose instances' diagnostics endpoints")
	require.EventuallyWithT(t, func(t *assert.CollectT) {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/%s/debug/config/successful", diagPort, mgrInstance1.ID()))
		if assert.NoError(t, err) {
			resp.Body.Close()
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}

		resp, err = http.Get(fmt.Sprintf("http://localhost:%d/%s/debug/config/successful", diagPort, mgrInstance2.ID()))
		if assert.NoError(t, err) {
			resp.Body.Close()
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}
	}, waitTime, tickTime, "diagnostics should be exposed under /{instanceID}/debug/config prefix for both instances")

	t.Log("Stopping the first instance and waiting for its diagnostics endpoints to be removed from the server")
	require.NoError(t, multimgr.StopInstance(mgrInstance1.ID()))
	require.EventuallyWithT(t, func(t *assert.CollectT) {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/%s/debug/config/successful", diagPort, mgrInstance1.ID()))
		require.NoError(t, err)
		resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	}, waitTime, tickTime, "diagnostics should no longer be available after stopping the instance")
}
