package envtest

import (
	"context"
	"go/build"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/gateway"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/builder"
	"github.com/kong/kubernetes-ingress-controller/v2/test/consts"
)

type Options struct {
	InstallGatewayCRDs bool
	InstallKongCRDs    bool
}

var DefaultEnvTestOpts = Options{
	InstallGatewayCRDs: true,
	InstallKongCRDs:    true,
}

type OptionModifier func(Options) Options

func WithInstallKongCRDs(install bool) OptionModifier {
	return func(opts Options) Options {
		opts.InstallKongCRDs = install
		return opts
	}
}

func WithInstallGatewayCRDs(install bool) OptionModifier {
	return func(opts Options) Options {
		opts.InstallGatewayCRDs = install
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
	if opts.InstallKongCRDs {
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

func deployIngressClass(ctx context.Context, t *testing.T, name string, client ctrlclient.Client) {
	ingress := &netv1.IngressClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: netv1.IngressClassSpec{
			Controller: store.IngressClassKongController,
		},
	}
	require.NoError(t, client.Create(ctx, ingress))
}

// deployGateway deploys a Gateway, GatewayClass, and ingress service for use in tests.
func deployGateway(ctx context.Context, t *testing.T, client ctrlclient.Client) gatewayapi.Gateway {
	ns := CreateNamespace(ctx, t, client)

	publishSvc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns.Name,
			Name:      IngressServiceName,
		},
		Spec: corev1.ServiceSpec{
			Ports: builder.NewServicePort().
				WithName("http").
				WithProtocol(corev1.ProtocolTCP).
				WithPort(8000).
				IntoSlice(),
		},
	}
	require.NoError(t, client.Create(ctx, &publishSvc))
	t.Cleanup(func() { _ = client.Delete(ctx, &publishSvc) })

	gwc := gatewayapi.GatewayClass{
		Spec: gatewayapi.GatewayClassSpec{
			ControllerName: gateway.GetControllerName(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
			Annotations: map[string]string{
				"konghq.com/gatewayclass-unmanaged": "placeholder",
			},
		},
	}

	require.NoError(t, client.Create(ctx, &gwc))
	t.Cleanup(func() { _ = client.Delete(ctx, &gwc) })

	gw := gatewayapi.Gateway{
		Spec: gatewayapi.GatewaySpec{
			GatewayClassName: gatewayapi.ObjectName(gwc.Name),
			Listeners: []gatewayapi.Listener{
				{
					Name:     "http",
					Protocol: gatewayapi.HTTPProtocolType,
					Port:     gatewayapi.PortNumber(8000),
				},
			},
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns.Name,
			Name:      uuid.NewString(),
		},
	}
	require.NoError(t, client.Create(ctx, &gw))
	t.Cleanup(func() { _ = client.Delete(ctx, &gw) })

	return gw
}
