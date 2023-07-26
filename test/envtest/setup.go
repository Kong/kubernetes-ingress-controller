package envtest

import (
	"go/build"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	"github.com/kong/kubernetes-ingress-controller/v2/test/consts"
)

type Options struct {
	InstallGatewayCRDs bool
	InstanllKongCRDs   bool
}

var DefaultEnvTestOpts = Options{
	InstallGatewayCRDs: true,
	InstanllKongCRDs:   false,
}

type OptionModifier func(Options) Options

func WithInstallKongCRDs(install bool) OptionModifier {
	return func(opts Options) Options {
		opts.InstanllKongCRDs = install
		return opts
	}
}

// Setup sets up the envtest environment which will be stopped on test cleanup
// using t.Cleanup().
//
// Note: If you want apiserver output on stdout set
// KUBEBUILDER_ATTACH_CONTROL_PLANE_OUTPUT to true when running tests.
func Setup(t *testing.T, scheme *k8sruntime.Scheme, optModifiers ...OptionModifier) *rest.Config {
	t.Helper()

	testEnv := &envtest.Environment{
		ControlPlaneStopTimeout: time.Second * 60,
	}

	t.Logf("starting envtest environment for test %s...", t.Name())
	cfg, err := testEnv.Start()
	require.NoError(t, err)

	opts := DefaultEnvTestOpts
	for _, mod := range optModifiers {
		opts = mod(opts)
	}

	if opts.InstallGatewayCRDs {
		installGatewayCRDs(t, scheme, cfg)
	}
	if opts.InstanllKongCRDs {
		installKongCRDs(t, scheme, cfg)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	done := make(chan struct{})
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		defer wg.Done()
		select {
		case <-ch:
			_ = testEnv.Stop()
		case <-done:
			_ = testEnv.Stop()
		}
	}()

	config, err := clientcmd.BuildConfigFromFlags(cfg.Host, "")
	require.NoError(t, err)
	config.CertData = cfg.CertData
	config.CAData = cfg.CAData
	config.KeyData = cfg.KeyData

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	require.NoError(t, err)

	i, err := discoveryClient.ServerVersion()
	require.NoError(t, err)

	t.Logf("envtest environment (%s) started at %s", i, cfg.Host)

	t.Cleanup(func() {
		t.Helper()
		t.Logf("stopping envtest environment for test %s", t.Name())
		close(done)
		wg.Wait()
	})

	return cfg
}

func installGatewayCRDs(t *testing.T, scheme *k8sruntime.Scheme, cfg *rest.Config) {
	gatewayCRDPath := filepath.Join(build.Default.GOPATH, "pkg", "mod", "sigs.k8s.io", "gateway-api@"+consts.GatewayAPIVersion, "config", "crd", "experimental")
	_, err := envtest.InstallCRDs(cfg, envtest.CRDInstallOptions{
		Scheme:             scheme,
		Paths:              []string{gatewayCRDPath},
		ErrorIfPathMissing: true,
	})
	require.NoError(t, err, "failed installing Gateway API CRDs")
}

func installKongCRDs(t *testing.T, scheme *k8sruntime.Scheme, cfg *rest.Config) {
	// extract project root path.
	_, thisFilePath, _, _ := runtime.Caller(0) //nolint:dogsled
	projectRoot := filepath.Join(filepath.Dir(thisFilePath), "..", "..")
	// install Kong CRDs from config/crd/bases.
	kongCRDPath := filepath.Join(projectRoot, "config", "crd", "bases")
	t.Logf("install Kong CRDs from manifests in %s", kongCRDPath)
	_, err := envtest.InstallCRDs(cfg, envtest.CRDInstallOptions{
		Scheme:             scheme,
		Paths:              []string{kongCRDPath},
		ErrorIfPathMissing: true,
	})
	require.NoError(t, err)
}
