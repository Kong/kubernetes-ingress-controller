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
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/health"
	"github.com/kong/kubernetes-ingress-controller/v3/pkg/manager"
	managercfg "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
	"github.com/kong/kubernetes-ingress-controller/v3/pkg/telemetry"
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

// WithDefaultEnvTestsConfig modifies a managercfg.Config for use in envtests.
func WithDefaultEnvTestsConfig(envcfg *rest.Config) func(cfg *managercfg.Config) {
	return func(cfg *managercfg.Config) {
		// Override the APIServer.
		cfg.APIServerHost = envcfg.Host
		cfg.APIServerCertData = envcfg.CertData
		cfg.APIServerKeyData = envcfg.KeyData
		cfg.APIServerCAData = envcfg.CAData

		// Shorten the wait in tests.
		cfg.UpdateStatus = false
		cfg.ProxySyncSeconds = 0.1
		cfg.InitCacheSyncDuration = 0

		cfg.MetricsAddr = "0"

		// And other settings which are irrelevant here.
		cfg.AnonymousReports = false

		// Set the GracefulShutdownTimeout to 0 to prevent errors:
		// failed waiting for all runnables to end within grace period of 30s: context deadline exceeded
		// Ref: https://github.com/kubernetes-sigs/controller-runtime/blob/e59161ee/pkg/manager/internal.go#L543-L548
		cfg.GracefulShutdownTimeout = lo.ToPtr(time.Duration(0))

		// Disable Gateway API controllers, enable those only in tests that use them.
		cfg.GatewayAPIGatewayController = false
		cfg.GatewayAPIHTTPRouteController = false
		cfg.GatewayAPIReferenceGrantController = false

		// Disable leader election, which doesn't work outside the cluster and is irrelevant for single-instance tests.
		cfg.LeaderElectionForce = managercfg.LeaderElectionDisabled
	}
}

func WithGatewayFeatureEnabled(cfg *managercfg.Config) {
	cfg.FeatureGates[managercfg.GatewayAlphaFeature] = true
}

func WithGatewayAPIControllers() func(cfg *managercfg.Config) {
	return func(cfg *managercfg.Config) {
		cfg.GatewayAPIGatewayController = true
		cfg.GatewayAPIHTTPRouteController = true
		cfg.GatewayAPIReferenceGrantController = true
	}
}

func WithGatewayToReconcile(gatewayNN string) func(cfg *managercfg.Config) {
	parts := strings.SplitN(gatewayNN, "/", 3)
	if len(parts) != 2 {
		panic("the expected format if namespace/name")
	}
	return func(cfg *managercfg.Config) {
		cfg.GatewayToReconcile = mo.Some(k8stypes.NamespacedName{
			Namespace: parts[0],
			Name:      parts[1],
		})
	}
}

func WithPublishService(namespace string) func(cfg *managercfg.Config) {
	return func(cfg *managercfg.Config) {
		cfg.PublishService = mo.Some(k8stypes.NamespacedName{
			Name:      PublishServiceName,
			Namespace: namespace,
		})
	}
}

func WithPublishStatusAddress(addresses []string, udps []string) func(cfg *managercfg.Config) {
	return func(cfg *managercfg.Config) {
		cfg.PublishStatusAddress = addresses
		cfg.PublishStatusAddressUDP = udps
	}
}

func WithIngressClass(name string) func(cfg *managercfg.Config) {
	return func(cfg *managercfg.Config) {
		cfg.IngressClassName = name
	}
}

func WithProxySyncSeconds(period float32) func(cfg *managercfg.Config) {
	return func(cfg *managercfg.Config) {
		cfg.ProxySyncSeconds = period
	}
}

func WithDiagnosticsServer(port int) func(cfg *managercfg.Config) {
	return func(cfg *managercfg.Config) {
		cfg.DiagnosticServerPort = port
		cfg.EnableConfigDumps = true
	}
}

func WithDiagnosticsWithoutServer() func(cfg *managercfg.Config) {
	return func(cfg *managercfg.Config) {
		cfg.EnableConfigDumps = true
		cfg.DisableRunningDiagnosticsServer = true
	}
}

func WithHealthProbePort(port int) func(cfg *managercfg.Config) {
	return func(cfg *managercfg.Config) {
		cfg.ProbeAddr = fmt.Sprintf("localhost:%d", port)
	}
}

func WithProfiling() func(cfg *managercfg.Config) {
	return func(cfg *managercfg.Config) {
		cfg.EnableProfiling = true
	}
}

func WithUpdateStatus() func(cfg *managercfg.Config) {
	return func(cfg *managercfg.Config) {
		cfg.UpdateStatus = true
	}
}

func WithKongServiceFacadeFeatureEnabled() func(cfg *managercfg.Config) {
	return func(cfg *managercfg.Config) {
		cfg.FeatureGates[managercfg.KongServiceFacadeFeature] = true
	}
}

func WithKongAdminURLs(urls ...string) func(cfg *managercfg.Config) {
	return func(cfg *managercfg.Config) {
		cfg.KongAdminURLs = urls
	}
}

func WithAdmissionWebhookEnabled(key, cert []byte, port int) func(cfg *managercfg.Config) {
	return func(cfg *managercfg.Config) {
		cfg.AdmissionServer.ListenAddr = fmt.Sprintf(":%d", port)
		cfg.AdmissionServer.Key = string(key)
		cfg.AdmissionServer.Cert = string(cert)
	}
}

func WithMetricsAddr(addr string) func(cfg *managercfg.Config) {
	return func(cfg *managercfg.Config) {
		cfg.MetricsAddr = addr
	}
}

func WithTelemetry(splunkEndpoint string, telemetryPeriod time.Duration) managercfg.Opt {
	return func(cfg *managercfg.Config) {
		cfg.AnonymousReports = true
		cfg.SplunkEndpoint = splunkEndpoint
		cfg.SplunkEndpointInsecureSkipVerify = true
		cfg.TelemetryPeriod = telemetryPeriod
		cfg.EnableProfiling = false
		cfg.EnableConfigDumps = false
	}
}

func WithCacheSyncTimeout(d time.Duration) func(cfg *managercfg.Config) {
	return func(cfg *managercfg.Config) {
		cfg.CacheSyncTimeout = d
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

func SetupManager(
	ctx context.Context,
	t *testing.T,
	mgrID manager.ID,
	envcfg *rest.Config,
	adminAPIOpts []mocks.AdminAPIHandlerOpt,
	modifyCfgFns ...managercfg.Opt,
) *manager.Manager {
	adminAPIServerURL := StartAdminAPIServerMock(t, adminAPIOpts...).URL

	modifyCfgFns = append([]managercfg.Opt{
		WithDefaultEnvTestsConfig(envcfg),
		WithKongAdminURLs(adminAPIServerURL),
	},
		modifyCfgFns..., // Add the user-provided modifyCfgFns last so they can override the defaults.
	)

	cfg, err := manager.NewConfig(modifyCfgFns...)
	require.NoError(t, err)

	logger := ctrl.LoggerFrom(ctx)
	mgr, err := manager.NewManager(ctx, mgrID, logger, cfg)
	require.NoError(t, err)

	if cfg.ProbeAddr != "" {
		t.Log("Starting standalone health check server")
		health.NewHealthCheckServer(
			healthz.Ping, health.NewHealthCheckerFromFunc(mgr.IsReady),
		).Start(ctx, cfg.ProbeAddr, logger.WithName("health-check"))
	}
	if cfg.AnonymousReports {
		stopAnonymousReports, err := telemetry.SetupAnonymousReports(
			ctx,
			mgr.GetKubeconfig(),
			mgr.GetClientsManager(),
			telemetry.ReportConfig{
				SplunkEndpoint:                   cfg.SplunkEndpoint,
				SplunkEndpointInsecureSkipVerify: cfg.SplunkEndpointInsecureSkipVerify,
				TelemetryPeriod:                  cfg.TelemetryPeriod,
				ReportValues: telemetry.ReportValues{
					PublishServiceNN:               cfg.PublishService.OrEmpty(),
					FeatureGates:                   cfg.FeatureGates,
					MeshDetection:                  len(cfg.WatchNamespaces) == 0,
					KonnectSyncEnabled:             cfg.Konnect.ConfigSynchronizationEnabled,
					GatewayServiceDiscoveryEnabled: cfg.KongAdminSvc.IsPresent(),
				},
			},
			mgrID,
		)
		// TODO: it's not closed properly (should be called after Run(...) returns), but for test it's fine for now.
		// It will leak only in case of running it with telemetry enabled.
		_ = stopAnonymousReports
		if err != nil {
			logger.Error(err, "Failed setting up anonymous reports, continuing without telemetry")
		} else {
			logger.Info("Anonymous reports enabled")
		}
	} else {
		logger.Info("Anonymous reports disabled, skipping")
	}

	return mgr
}

// RunManager runs the manager in a goroutine. It's possible to modify the manager's configuration
// by passing in modifyCfgFns. The manager is stopped when the context is canceled.
func RunManager(
	ctx context.Context,
	t *testing.T,
	envcfg *rest.Config,
	adminAPIOpts []mocks.AdminAPIHandlerOpt,
	modifyCfgFns ...managercfg.Opt,
) LogsObserver {
	mgrID, err := manager.NewID(t.Name())
	require.NoError(t, err)

	ctx, _, logs := CreateTestLogger(ctx)

	// This wait group makes it so that we wait for manager to exit.
	// This way we get clean test logs not mixing between tests.
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		mgr := SetupManager(ctx, t, mgrID, envcfg, adminAPIOpts, modifyCfgFns...)
		require.NoError(t, mgr.Start(ctx))
	}()
	t.Cleanup(func() {
		wg.Wait()
		DumpLogsIfTestFailed(t, logs)
	})

	return logs
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
