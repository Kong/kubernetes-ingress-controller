//go:build envtest
// +build envtest

package rootcmd_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/rest"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/cmd/rootcmd"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager/featuregates"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset/scheme"
	"github.com/kong/kubernetes-ingress-controller/v2/test/envtest"
)

func TestDebugEndpoints(t *testing.T) {
	t.Parallel()

	const (
		waitTime = time.Minute
		tickTime = 10 * time.Millisecond
	)

	envcfg := envtest.Setup(t, scheme.Scheme)
	cfg := configForEnvConfig(t, envcfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func(ctx context.Context) {
		rootcmd.Run(ctx, &cfg, io.Discard) //nolint:errcheck
	}(ctx)

	urls := []struct {
		name string
		port int
	}{
		{
			port: manager.DiagnosticsPort,
			name: "debug/pprof/",
		},
		{
			port: manager.DiagnosticsPort,
			name: "debug/config/successful",
		},
		{
			port: manager.DiagnosticsPort,
			name: "debug/config/failed",
		},
		{
			port: manager.HealthzPort,
			name: "healthz",
		},
		{
			port: manager.HealthzPort,
			name: "readyz",
		},
		{
			port: manager.MetricsPort,
			name: "metrics",
		},
	}

	for _, u := range urls {
		u := u
		t.Run(u.name, func(t *testing.T) {
			url := fmt.Sprintf("http://localhost:%d/%s", u.port, u.name)
			eventuallHTTPGet(t, http.DefaultClient, url, waitTime, tickTime)
		})
	}
}

func eventuallHTTPGet(t *testing.T, h *http.Client, url string, waitTime, tickTime time.Duration) {
	t.Helper()

	t.Logf("HTTP GET %s", url)
	assert.EventuallyWithT(t, func(t *assert.CollectT) {
		resp, err := h.Get(url)
		if !assert.NoError(t, err) {
			return
		}
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}, waitTime, tickTime)
}

func startAdminAPIServerMock(t *testing.T) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type Root struct {
			Configuration map[string]any `json:"configuration"`
			Version       string         `json:"version"`
		}
		body, err := json.Marshal(Root{
			Version: "3.3.0",
			Configuration: map[string]any{
				"database":      "off",
				"router_flavor": "traditional",
			},
		})
		if err != nil {
			w.WriteHeader(500)
			return
		}

		if _, err := w.Write(body); err != nil {
			w.WriteHeader(500)
			return
		}
	})
	s := httptest.NewServer(handler)
	t.Cleanup(func() {
		s.Close()
	})
	return s
}

func configForEnvConfig(t *testing.T, envcfg *rest.Config) manager.Config {
	t.Helper()

	cfg := manager.Config{}
	cfg.FlagSet() // Just set the defaults.

	// Enable debugging endpoints.
	cfg.EnableProfiling = true
	cfg.EnableConfigDumps = true

	// Override the APIServer.
	cfg.APIServerHost = envcfg.Host
	cfg.APIServerCertData = envcfg.CertData
	cfg.APIServerKeyData = envcfg.KeyData
	cfg.APIServerCAData = envcfg.CAData

	cfg.KongAdminURLs = []string{startAdminAPIServerMock(t).URL}
	cfg.UpdateStatus = false
	cfg.ProxySyncSeconds = 0.1

	// And other settings which are irrelevant.
	cfg.Konnect.ConfigSynchronizationEnabled = false
	cfg.Konnect.LicenseSynchronizationEnabled = false
	cfg.AnonymousReports = false
	cfg.FeatureGates = featuregates.GetFeatureGatesDefaults()
	cfg.FeatureGates[featuregates.GatewayFeature] = false

	return cfg
}
