package envtest

import (
	"fmt"
	"testing"

	"github.com/phayes/freeport"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/rest"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager/featuregates"
)

// ConfigForEnvConfig prepares a manager.Config for use in tests
// It will start a mock Admin API server which will be set in KIC's config
// and which will be automatically stopped during test cleanup.
func ConfigForEnvConfig(t *testing.T, envcfg *rest.Config) manager.Config {
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

	cfg.KongAdminURLs = []string{StartAdminAPIServerMock(t).URL}
	cfg.UpdateStatus = false
	// Shorten the wait in tests.
	cfg.ProxySyncSeconds = 0.1
	cfg.InitCacheSyncDuration = 0

	p, err := freeport.GetFreePort()
	require.NoError(t, err)
	cfg.MetricsAddr = fmt.Sprintf("localhost:%d", p)

	// And other settings which are irrelevant here.
	cfg.Konnect.ConfigSynchronizationEnabled = false
	cfg.Konnect.LicenseSynchronizationEnabled = false
	cfg.AnonymousReports = false
	cfg.FeatureGates = featuregates.GetFeatureGatesDefaults()
	cfg.FeatureGates[featuregates.GatewayFeature] = false

	return cfg
}
