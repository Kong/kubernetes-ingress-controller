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
	"github.com/stretchr/testify/require"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/cmd/rootcmd"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager/featuregates"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

const (
	// IngressServiceName is the name of the ingress service used in Gateway API tests.
	IngressServiceName = "ingress-svc"
)

// ConfigForEnvConfig prepares a manager.Config for use in tests
// It will start a mock Admin API server which will be set in KIC's config
// and which will be automatically stopped during test cleanup.
func ConfigForEnvConfig(t *testing.T, envcfg *rest.Config) manager.Config {
	t.Helper()

	cfg := manager.Config{}
	cfg.FlagSet() // Just set the defaults.

	// Disable debugging endpoints.
	// If need be those can be enabled by manipulating the returned config.
	cfg.EnableProfiling = false
	cfg.EnableConfigDumps = false

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

func WithIngressService(namespace string) func(cfg *manager.Config) {
	return func(cfg *manager.Config) {
		cfg.IngressAddresses = []string{"127.0.0.1"}
		cfg.IngressService = mo.Some(k8stypes.NamespacedName{
			Name:      IngressServiceName,
			Namespace: namespace,
		})
	}
}

func WithIngressClass(name string) func(cfg *manager.Config) {
	return func(cfg *manager.Config) {
		cfg.IngressClassName = name
	}
}

func WithProxySyncSeconds(period float32) func(cfg *manager.Config) {
	return func(cfg *manager.Config) {
		cfg.ProxySyncSeconds = period
	}
}

func WithDiagnosticsServer(port int) func(cfg *manager.Config) {
	return func(cfg *manager.Config) {
		cfg.DiagnosticServerPort = port
		cfg.EnableConfigDumps = true
	}
}

func WithHealthProbePort(port int) func(cfg *manager.Config) {
	return func(cfg *manager.Config) {
		cfg.ProbeAddr = fmt.Sprintf("localhost:%d", port)
	}
}

func WithProfiling() func(cfg *manager.Config) {
	return func(cfg *manager.Config) {
		cfg.EnableProfiling = true
	}
}

// buffer is a goroutine safe bytes.Buffer.
type buffer struct {
	buffer bytes.Buffer
	mutex  sync.RWMutex
}

// Write appends the contents of p to the buffer, growing the buffer as needed.
// It returns the number of bytes written.
func (s *buffer) Write(p []byte) (n int, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.buffer.Write(p)
}

// String returns the contents of the unread portion of the buffer
// as a string. If the Buffer is a nil pointer, it returns "<nil>".
func (s *buffer) String() string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.buffer.String()
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
	b := &buffer{}
	logrusLogger.Out = b
	logger := logrusr.New(logrusLogger)
	ctx = ctrl.LoggerInto(ctx, logger)
	ctrl.SetLogger(logger)

	// This wait group makes it so that we wait for manager to exit.
	// This way we get clean test logs not mixing between tests.
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		var configDumps util.ConfigDumpDiagnostic
		if cfg.EnableConfigDumps {
			diag, err := rootcmd.StartDiagnosticsServer(ctx, cfg.DiagnosticServerPort, &cfg, logger)
			require.NoError(t, err)
			configDumps = diag.ConfigDumps
		}

		require.NoError(t, manager.Run(ctx, &cfg, configDumps, logrusLogger))
	}()
	t.Cleanup(func() {
		wg.Wait()
		if t.Failed() {
			t.Logf("manager logs:\n%s", b.String())
		}
	})

	return loggerHook
}
