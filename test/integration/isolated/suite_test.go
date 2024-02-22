//go:build integration_tests

package isolated

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/blang/semver/v4"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/metallb"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/kind"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/e2e-framework/klient/conf"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	testutils "github.com/kong/kubernetes-ingress-controller/v3/internal/util/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/consts"
	testhelpers "github.com/kong/kubernetes-ingress-controller/v3/test/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers/certificate"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testenv"
)

var tenv env.Environment

// -----------------------------------------------------------------------------
// Testing Main
// -----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := envconf.NewFromFlags()
	helpers.ExitOnErrWithCode(ctx, err, consts.ExitCodeEnvSetupFailed)
	cfg.WithKubeconfigFile(conf.ResolveKubeConfigFile())
	tenv = env.NewWithConfig(cfg)

	var (
		// Specifying a run ID so that multiple runs wouldn't collide.
		// It is used when creating tests namespaces and their labels.
		runID = envconf.RandomName("", 3)
		// The env is shared and built only once.
		env environments.Environment
	)

	fmt.Printf("INFO: runID %s\n", runID)

	builder := environments.NewBuilder()
	fmt.Println("INFO: configuring cluster for testing environment")
	if existingCluster := testenv.ExistingClusterName(); existingCluster != "" {
		if cv := testenv.ClusterVersion(); cv != "" {
			helpers.ExitOnErrWithCode(ctx,
				fmt.Errorf("can't flag cluster version (%s) & provide an existing cluster at the same time", cv),
				consts.ExitCodeIncompatibleOptions)
		}
		clusterParts := strings.Split(existingCluster, ":")
		if len(clusterParts) != 2 {
			helpers.ExitOnErrWithCode(ctx, fmt.Errorf("existing cluster in wrong format (%s): format is <TYPE>:<NAME> (e.g. kind:test-cluster)", existingCluster), consts.ExitCodeCantUseExistingCluster)
		}
		clusterType, clusterName := clusterParts[0], clusterParts[1]

		fmt.Printf("INFO: using existing %s cluster %s\n", clusterType, clusterName)
		switch clusterType {
		case string(kind.KindClusterType):
			cluster, err := kind.NewFromExisting(clusterName)
			helpers.ExitOnErr(ctx, err)
			builder.WithExistingCluster(cluster)
			builder.WithAddons(metallb.New())
		default:
			helpers.ExitOnErrWithCode(ctx, fmt.Errorf("unknown cluster type: %s", clusterType), consts.ExitCodeCantUseExistingCluster)
		}

	} else {
		fmt.Println("INFO: no existing cluster found, deploying using Kubernetes In Docker (KIND)")

		builder.WithAddons(metallb.New())

		if testenv.ClusterVersion() != "" {
			clusterVersion, err := semver.Parse(strings.TrimPrefix(testenv.ClusterVersion(), "v"))
			helpers.ExitOnErr(ctx, err)

			fmt.Printf("INFO: build a new KIND cluster with version %s\n", clusterVersion.String())
			builder.WithKubernetesVersion(clusterVersion)
		}
	}

	fmt.Println("INFO: building test environment")
	env, err = builder.Build(ctx)
	helpers.ExitOnErr(ctx, err)

	fmt.Printf("INFO: waiting for cluster %s and all addons to become ready\n", env.Cluster().Name())
	envReadyCtx, envReadyCancel := context.WithTimeout(ctx, testenv.EnvironmentReadyTimeout())
	defer envReadyCancel()
	helpers.ExitOnErr(ctx, <-env.WaitForReady(envReadyCtx))

	if err := testutils.PrepareClusterForRunningControllerManager(ctx, env.Cluster()); err != nil {
		helpers.ExitOnErr(ctx, fmt.Errorf("failed to prepare cluster for running the controller manager: %w", err))
	}

	ctx = SetClusterInCtx(ctx, env.Cluster())
	ctx = SetRunIDInCtx(ctx, runID)
	tenv = tenv.WithContext(ctx)

	// TODO: can't use any of AfterEachFeature,BeforeEachFeature,AfterEachTest,BeforeEachTest
	// to get conditional setup and teardown.
	// Related: https://github.com/Kong/kubernetes-ingress-controller/issues/4847

	var l sync.RWMutex
	tenv.BeforeEachFeature(
		// TODO: Prevent a data race by using a mutex explicitly when first creating the client.
		// Related: https://github.com/Kong/kubernetes-ingress-controller/issues/4848
		func(ctx context.Context, c *envconf.Config, _ *testing.T, _ features.Feature) (context.Context, error) {
			l.Lock()
			defer l.Unlock()
			_, err = c.NewClient()
			return ctx, err
		},
	)

	code := tenv.Run(m)
	defer func() {
		os.Exit(code)
	}()

	if testenv.IsCI() {
		fmt.Printf("INFO: running in ephemeral CI environment, skipping cluster %s teardown\n", env.Cluster().Name())
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), test.EnvironmentCleanupTimeout)
		defer cancel()
		helpers.ExitOnErr(ctx, helpers.RemoveCluster(ctx, env.Cluster()))
	}
}

type featureSetupCfg struct {
	controllerManagerOpts []helpers.ControllerManagerOpt
	kongProxyEnvVars      map[string]string
}

type featureSetupOpt func(*featureSetupCfg)

func withControllerManagerOpts(opts ...helpers.ControllerManagerOpt) featureSetupOpt {
	return func(o *featureSetupCfg) {
		o.controllerManagerOpts = opts
	}
}

func withKongProxyEnvVars(envVars map[string]string) featureSetupOpt {
	return func(o *featureSetupCfg) {
		o.kongProxyEnvVars = envVars
	}
}

func featureSetup(opts ...featureSetupOpt) func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
	var setupCfg featureSetupCfg
	for _, opt := range opts {
		opt(&setupCfg)
	}
	return func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		// TODO: this is temporary to allow things like:
		// clusters.ApplyManifestByYAML(ctx, cluster, s)
		// in tests.
		// Ideally this wouldn't be possible because it allows tests to break
		// a shared cluster but we don't have means to use kubectl against a cluster
		// without including a 3rd party package.
		// We should remove the cluster from the context as the last setup step.
		cluster := GetClusterFromCtx(ctx)

		runID := GetRunIDFromCtx(ctx)

		ctx, err := CreateNSForTest(ctx, cfg, t, runID)
		if !assert.NoError(t, err) {
			return ctx
		}

		namespace := GetNamespaceForT(ctx, t)

		// Logger needs to be configured before anything else happens.
		// This is because the controller manager has a timeout for
		// logger initialization, and if the logger isn't configured
		// after 30s from the start of controller manager package init function,
		// the controller manager will set up a no op logger and continue.
		// The logger cannot be configured after that point.
		logger, logOutput, err := testutils.SetupLoggers("trace", "text")
		if !assert.NoError(t, err, "failed to setup loggers") {
			return ctx
		}
		if logOutput != "" {
			t.Logf("writing manager logs to %s", logOutput)
		}

		t.Logf("setting up test environment")
		kongBuilder, extraControllerArgs, err := helpers.GenerateKongBuilder(ctx)
		if !assert.NoError(t, err) {
			return ctx
		}
		if testenv.KongImage() != "" && testenv.KongTag() != "" {
			fmt.Printf("INFO: custom kong image specified via env: %s:%s\n", testenv.KongImage(), testenv.KongTag())
		}

		// Pin the Helm chart version.
		kongBuilder.WithHelmChartVersion(testenv.KongHelmChartVersion())
		kongBuilder.WithNamespace(namespace)
		kongBuilder.WithName(NameFromT(t))
		kongBuilder.WithAdditionalValue("readinessProbe.initialDelaySeconds", "1")
		for name, value := range setupCfg.kongProxyEnvVars {
			kongBuilder.WithProxyEnvVar(name, value)
		}

		kongAddon := kongBuilder.Build()
		t.Logf("deploying kong addon to cluster %s in namespace %s", cluster.Name(), namespace)
		if !assert.NoError(t, cluster.DeployAddon(ctx, kongAddon)) {
			return ctx
		}
		ctx = SetInCtxForT(ctx, t, kongAddon)

		cleaner := clusters.NewCleaner(cluster)
		t.Cleanup(func() {
			if err := cleaner.Cleanup(ctx); err != nil {
				fmt.Printf("ERROR: failed cleaning up the cluster: %v\n", err)
			}
		})
		ctx = SetInCtxForT(ctx, t, cleaner)

		if !assert.Eventually(t, func() bool {
			_, ok, err := kongAddon.Ready(ctx, cluster)
			if err != nil {
				t.Logf("error checking if kong addon is ready: %v", err)
				return false
			}

			return ok
		}, time.Minute*3, 100*time.Millisecond, "failed waiting for kong addon to become ready") {
			return ctx
		}

		t.Logf("collecting urls from the kong proxy deployment in namespace: %s", namespace)
		if !assert.NoError(t, err) {
			return ctx
		}
		proxyAdminURL, err := kongAddon.ProxyAdminURL(ctx, cluster)
		if !assert.NoError(t, err) {
			return ctx
		}
		ctx = SetAdminURLInCtx(ctx, proxyAdminURL)

		proxyUDPURL, err := kongAddon.ProxyUDPURL(ctx, cluster)
		if !assert.NoError(t, err) {
			return ctx
		}
		ctx = SetUDPURLInCtx(ctx, proxyUDPURL)

		proxyTCPURL, err := kongAddon.ProxyTCPURL(ctx, cluster)
		if !assert.NoError(t, err) {
			return ctx
		}
		ctx = SetTCPURLInCtx(ctx, proxyTCPURL)

		proxyTLSURL, err := kongAddon.ProxyTLSURL(ctx, cluster)
		if !assert.NoError(t, err) {
			return ctx
		}
		ctx = SetTLSURLInCtx(ctx, proxyTLSURL)

		proxyHTTPURL, err := kongAddon.ProxyHTTPURL(ctx, cluster)
		if !assert.NoError(t, err) {
			return ctx
		}
		ctx = SetHTTPURLInCtx(ctx, proxyHTTPURL)

		proxyHTTPSURL, err := kongAddon.ProxyHTTPSURL(ctx, cluster)
		if !assert.NoError(t, err) {
			return ctx
		}
		ctx = SetHTTPSURLInCtx(ctx, proxyHTTPSURL)

		if !assert.NoError(t, retry.Do(
			func() error {
				version, err := helpers.GetKongVersion(ctx, proxyAdminURL, consts.KongTestPassword)
				if err != nil {
					return err
				}
				t.Logf("using Kong instance (version: %s) reachable at %s\n", version, proxyAdminURL)
				return nil
			},
			retry.OnRetry(
				func(n uint, err error) {
					t.Logf("WARNING: try to get Kong Gateway version attempt %d/10 - error: %s\n", n+1, err)
				},
			),
			retry.LastErrorOnly(true),
			retry.Attempts(10),
		), "failed getting Kong's version") {
			return ctx
		}

		t.Logf("configuring feature gates")
		// TODO: https://github.com/Kong/kubernetes-ingress-controller/issues/4849
		featureGates := consts.DefaultFeatureGates
		t.Logf("feature gates enabled: %s", featureGates)

		t.Logf("starting the controller manager")
		cert, key := certificate.GetKongSystemSelfSignedCerts()
		metricsPort := testhelpers.GetFreePort(t)
		healthProbePort := testhelpers.GetFreePort(t)
		ingressClass := envconf.RandomName("ingressclass", 16)
		standardControllerArgs := []string{
			fmt.Sprintf("--health-probe-bind-address=localhost:%d", healthProbePort),
			fmt.Sprintf("--metrics-bind-address=localhost:%d", metricsPort),
			fmt.Sprintf("--ingress-class=%s", ingressClass),
			fmt.Sprintf("--admission-webhook-cert=%s", cert),
			fmt.Sprintf("--admission-webhook-key=%s", key),
			fmt.Sprintf("--admission-webhook-listen=0.0.0.0:%d", testutils.AdmissionWebhookListenPort),
			"--anonymous-reports=false",
			"--log-level=trace",
			fmt.Sprintf("--feature-gates=%s", featureGates),
			// Use fixed election namespace `kong` because RBAC roles for leader election are in the namespace,
			// so we create resources for leader election in the namespace to make sure that KIC can operate these resources.
			fmt.Sprintf("--election-namespace=%s", consts.ControllerNamespace),
			fmt.Sprintf("--watch-namespace=%s", kongAddon.Namespace()),
		}
		allControllerArgs := append(standardControllerArgs, extraControllerArgs...)
		for _, opt := range setupCfg.controllerManagerOpts {
			allControllerArgs = opt(allControllerArgs)
		}

		cancel, err := testutils.DeployControllerManagerForCluster(ctx, logger, cluster, kongAddon, allControllerArgs...)
		t.Cleanup(func() { cancel() })
		if !assert.NoError(t, err, "failed deploying controller manager") {
			return ctx
		}

		t.Logf("deploying the controller's IngressClass %q", ingressClass)
		if !assert.NoError(t, helpers.CreateIngressClass(ctx, ingressClass, cluster.Client()), "failed creating IngressClass") {
			return ctx
		}
		defer func() {
			// deleting this directly instead of adding it to the cleaner because
			// the cleaner always gets a 404 on it for unknown reasons
			_ = cluster.Client().NetworkingV1().IngressClasses().Delete(ctx, ingressClass, metav1.DeleteOptions{})
		}()
		ctx = setInCtx(ctx, _ingressClass{}, ingressClass)

		clusterVersion, err := cluster.Version()
		if !assert.NoError(t, err, "failed getting cluster version") {
			return ctx
		}
		t.Logf("testing environment is ready KUBERNETES_VERSION=(%v): running tests", clusterVersion)

		// TODO refactor. Perhaps there's a better way than just storing the cancel func in context.
		ctx = SetInCtxForT(ctx, t, cancel)

		return ctx
	}
}

func featureTeardown() func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
	return func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
		// Call cancel to stop the manager - this prevents Feature tests from running until the whole suite ends.
		cancel := GetFromCtxForT[func()](ctx, t)
		cancel()

		cluster := GetClusterFromCtx(ctx)
		runID := GetRunIDFromCtx(ctx)

		kongAddon := GetFromCtxForT[*kong.Addon](ctx, t)
		assert.NoError(t, cluster.DeleteAddon(ctx, kongAddon))

		ctx, err := deleteNSForTest(ctx, c, t, runID)
		assert.NoError(t, err)
		return ctx
	}
}
