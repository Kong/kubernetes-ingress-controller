package envtest

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/bombsimon/logrusr/v4"
	"github.com/phayes/freeport"
	"github.com/samber/mo"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager/featuregates"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

const (
	// PublishServiceName is the name of the publish service used in Gateway API tests.
	PublishServiceName = "publish-svc"
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

type ModifyManagerConfigFn func(cfg *manager.Config)

func WithGatewayFeatureEnabled(cfg *manager.Config) {
	cfg.FeatureGates[featuregates.GatewayFeature] = true
	cfg.FeatureGates[featuregates.GatewayAlphaFeature] = true
}

func WithPublishService(namespace string) func(cfg *manager.Config) {
	return func(cfg *manager.Config) {
		cfg.PublishStatusAddress = []string{"127.0.0.1"}
		cfg.PublishService = mo.Some(k8stypes.NamespacedName{
			Name:      PublishServiceName,
			Namespace: namespace,
		})
	}
}

// RunManager runs the manager in a goroutine. It's possible to modify the manager's configuration
// by passing in modifyCfgFns. The manager is stopped when the context is canceled.
func RunManager(
	ctx context.Context,
	t *testing.T,
	envcfg *rest.Config,
	modifyCfgFns ...func(cfg *manager.Config),
) (loggerHook *test.Hook) {
	cfg := ConfigForEnvConfig(t, envcfg)

	for _, modifyCfgFn := range modifyCfgFns {
		modifyCfgFn(&cfg)
	}

	logrusLogger, loggerHook := test.NewNullLogger()
	b := &bytes.Buffer{}
	logrusLogger.Out = b
	logger := logrusr.New(logrusLogger)
	ctx = ctrl.LoggerInto(ctx, logger)

	// This wait group makes it so that we wait for manager to exit.
	// This way we get clean test logs not mixing between tests.
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := manager.Run(ctx, &cfg, util.ConfigDumpDiagnostic{}, logrusLogger)
		assert.NoError(t, err)
	}()
	t.Cleanup(func() {
		wg.Wait()
		if t.Failed() {
			t.Logf("manager logs:\n%s", b.String())
		}
	})

	return loggerHook
}
