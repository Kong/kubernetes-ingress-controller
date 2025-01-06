package manager

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/controller-runtime/pkg/healthz"

	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers"
)

func TestHealthCheckServer(t *testing.T) {
	passChecker := func(_ *http.Request) error {
		return nil
	}
	failChecker := func(_ *http.Request) error {
		return errors.New("you shall not pass")
	}

	testCases := []struct {
		name               string
		healthzChecker     healthz.Checker
		readyzChecker      healthz.Checker
		livenessProbeCode  int
		readinessProbeCode int
	}{
		{
			name:               "healthz and readyz both pass should return both 200",
			healthzChecker:     passChecker,
			readyzChecker:      passChecker,
			livenessProbeCode:  http.StatusOK,
			readinessProbeCode: http.StatusOK,
		},
		{
			name:               "healthz fail, readyz not set should return 500 & 404",
			healthzChecker:     failChecker,
			livenessProbeCode:  http.StatusInternalServerError,
			readinessProbeCode: http.StatusNotFound,
		},
		{
			name:               "healthz pass, readyz fail should return 200 & 500",
			healthzChecker:     passChecker,
			readyzChecker:      failChecker,
			livenessProbeCode:  http.StatusOK,
			readinessProbeCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := &healthCheckServer{}
			h.setHealthzCheck(tc.healthzChecker)
			h.setReadyzCheck(tc.readyzChecker)
			s := httptest.NewServer(h)
			defer s.Close()

			livenessResp, err := http.Get(s.URL + "/healthz")
			require.NoError(t, err)
			defer livenessResp.Body.Close()
			require.Equal(t, tc.livenessProbeCode, livenessResp.StatusCode)

			readinessResp, err := http.Get(s.URL + "/readyz")
			require.NoError(t, err)
			defer readinessResp.Body.Close()
			require.Equal(t, tc.readinessProbeCode, readinessResp.StatusCode)
		})
	}
}

func TestHealthCheckServer_Start(t *testing.T) {
	h := &healthCheckServer{}
	h.setHealthzCheck(healthz.Ping)

	// Get free local port.
	port := helpers.GetFreePort(t)

	ctx, cancel := context.WithCancel(context.Background())
	addr := fmt.Sprintf("localhost:%d", port)
	// Use discard logger to prevent:
	// panic: Log in goroutine after TestHealthCheckServer_Start has completed: "level"=0 "msg"="healthz server closed"
	h.Start(ctx, addr, logr.Discard())

	healtzEndpoint := fmt.Sprintf("http://%s/healthz", addr)
	// Allow some failures just after the server gets started.
	require.Eventually(t, func() bool {
		resp, err := http.Get(healtzEndpoint)
		if err != nil {
			t.Logf("got error: %v but none expected", err)
			return false
		}

		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Logf("got %d status code but expected 200", resp.StatusCode)
			return false
		}

		return true
	}, time.Second, time.Millisecond)

	// Cancel the context to stop the server and check it is no longer listening.
	cancel()

	require.Eventually(t, func() bool {
		resp, err := http.Get(healtzEndpoint)
		if err == nil {
			defer resp.Body.Close()
		}
		if err != nil && !strings.Contains(err.Error(), "connection refused") {
			t.Log("expected error to contain 'connection refused', got:", err)
			return false
		}
		return true
	}, time.Second, time.Millisecond)
}
