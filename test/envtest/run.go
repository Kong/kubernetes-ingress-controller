package envtest

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest/observer"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/cmd/rootcmd"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/featuregates"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/mocks"
)

const (
	// PublishServiceName is the name of the publish service used in Gateway API tests.
	PublishServiceName = "publish-svc"

	// ManagerStartupWaitTime is the time to wait for the manager to start.
	ManagerStartupWaitTime = 5 * time.Second

	// ManagerStartupWaitInterval is the interval to wait for the manager to start.
	ManagerStartupWaitInterval = time.Millisecond
)

// ConfigForEnvConfig prepares a manager.Config for use in tests
// It will start a mock Admin API server which will be set in KIC's config
// and which will be automatically stopped during test cleanup.
func ConfigForEnvConfig(t *testing.T, envcfg *rest.Config, opts ...mocks.AdminAPIHandlerOpt) manager.Config {
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

	cfg.KongAdminURLs = []string{StartAdminAPIServerMock(t, opts...).URL}
	cfg.UpdateStatus = false
	// Shorten the wait in tests.
	cfg.ProxySyncSeconds = 0.1
	cfg.InitCacheSyncDuration = 0

	p := helpers.GetFreePort(t)
	cfg.MetricsAddr = fmt.Sprintf("localhost:%d", p)

	// And other settings which are irrelevant here.
	cfg.Konnect.ConfigSynchronizationEnabled = false
	cfg.Konnect.LicenseSynchronizationEnabled = false
	cfg.AnonymousReports = false
	cfg.FeatureGates = featuregates.GetFeatureGatesDefaults()

	// Set the GracefulShutdownTimeout to 0 to prevent errors:
	// failed waiting for all runnables to end within grace period of 30s: context deadline exceeded
	// Ref: https://github.com/kubernetes-sigs/controller-runtime/blob/e59161ee/pkg/manager/internal.go#L543-L548
	cfg.GracefulShutdownTimeout = lo.ToPtr(time.Duration(0))

	// Disable Gateway API controllers, enable those only in tests that use them.
	cfg.GatewayAPIGatewayController = false
	cfg.GatewayAPIHTTPRouteController = false
	cfg.GatewayAPIReferenceGrantController = false

	return cfg
}

type ModifyManagerConfigFn func(cfg *manager.Config)

func WithGatewayFeatureEnabled(cfg *manager.Config) {
	cfg.FeatureGates[featuregates.GatewayAlphaFeature] = true
}

func WithGatewayAPIControllers() func(cfg *manager.Config) {
	return func(cfg *manager.Config) {
		cfg.GatewayAPIGatewayController = true
		cfg.GatewayAPIHTTPRouteController = true
		cfg.GatewayAPIReferenceGrantController = true
	}
}

func WithGatewayToReconcile(gatewayNN string) func(cfg *manager.Config) {
	parts := strings.SplitN(gatewayNN, "/", 3)
	if len(parts) != 2 {
		panic("the expected format if namespace/name")
	}
	return func(cfg *manager.Config) {
		cfg.GatewayToReconcile = mo.Some(k8stypes.NamespacedName{
			Namespace: parts[0],
			Name:      parts[1],
		})
	}
}

func WithPublishService(namespace string) func(cfg *manager.Config) {
	return func(cfg *manager.Config) {
		cfg.PublishService = mo.Some(k8stypes.NamespacedName{
			Name:      PublishServiceName,
			Namespace: namespace,
		})
	}
}

func WithPublishStatusAddress(addresses []string, udps []string) func(cfg *manager.Config) {
	return func(cfg *manager.Config) {
		cfg.PublishStatusAddress = addresses
		cfg.PublishStatusAddressUDP = udps
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

func WithUpdateStatus() func(cfg *manager.Config) {
	return func(cfg *manager.Config) {
		cfg.UpdateStatus = true
	}
}

func WithKongServiceFacadeFeatureEnabled() func(cfg *manager.Config) {
	return func(cfg *manager.Config) {
		cfg.FeatureGates[featuregates.KongServiceFacade] = true
	}
}

func WithKongAdminURLs(urls ...string) func(cfg *manager.Config) {
	return func(cfg *manager.Config) {
		cfg.KongAdminURLs = urls
	}
}

func WithAdmissionWebhookEnabled(key, cert []byte, addr string) func(cfg *manager.Config) {
	return func(cfg *manager.Config) {
		cfg.AdmissionServer.ListenAddr = addr
		cfg.AdmissionServer.Key = string(key)
		cfg.AdmissionServer.Cert = string(cert)
	}
}

// AdminAPIOptFns wraps a variadic list of mocks.AdminAPIHandlerOpt and returns
// a slice containing all of them.
// The purpose of this is func is to make the call sites a bit less verbose.
//
// NOTE: Ideally we'd refactor the RunManager() so that it'd not need to accept
// an empty slice of mocks.AdminAPIHandlerOpt or a call to AdminAPIOptFns() with
// no arguments but we can't accept 2 variadic list parameters.
// A slight refactor might be beneficial here.
func AdminAPIOptFns(fns ...mocks.AdminAPIHandlerOpt) []mocks.AdminAPIHandlerOpt {
	return fns
}

// RunManager runs the manager in a goroutine. It's possible to modify the manager's configuration
// by passing in modifyCfgFns. The manager is stopped when the context is canceled.
func RunManager(
	ctx context.Context,
	t *testing.T,
	envcfg *rest.Config,
	adminAPIOpts []mocks.AdminAPIHandlerOpt,
	modifyCfgFns ...func(cfg *manager.Config),
) (manager.Config, LogsObserver) {
	cfg := ConfigForEnvConfig(t, envcfg, adminAPIOpts...)

	for _, modifyCfgFn := range modifyCfgFns {
		modifyCfgFn(&cfg)
	}

	ctx, logger, logs := CreateTestLogger(ctx)

	// This wait group makes it so that we wait for manager to exit.
	// This way we get clean test logs not mixing between tests.
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		diagServer, err := rootcmd.StartDiagnosticsServer(ctx, cfg.DiagnosticServerPort, &cfg, logger)
		require.NoError(t, err)

		require.NoError(t, manager.Run(ctx, &cfg, diagServer.ConfigDumps(), logger))
	}()
	t.Cleanup(func() {
		wg.Wait()
		DumpLogsIfTestFailed(t, logs)
	})

	return cfg, logs
}

// WaitForManagerStart waits for the manager to start. The indication of the manager starting is
// the "Starting manager" log entry that is emitted just before the manager starts.
// Note: We cannot rely here on the manager's readiness probe because it returns 200 OK as soon as it
// starts listening which happens before the manager actually starts.
func WaitForManagerStart(t *testing.T, logsObserver LogsObserver) {
	t.Helper()
	t.Log("Waiting for manager to start...")
	require.Eventually(t, func() bool {
		const expectedLog = "Starting manager"
		return lo.ContainsBy(logsObserver.All(), func(item observer.LoggedEntry) bool {
			return strings.Contains(item.Message, expectedLog)
		})
	}, ManagerStartupWaitTime, ManagerStartupWaitInterval)
}
