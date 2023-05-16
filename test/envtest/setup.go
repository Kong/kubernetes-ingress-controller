package envtest

import (
	"go/build"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/test/consts"
)

// Setup sets up the envtest environment which will be stopped on test cleanup
// using t.Cleanup().
//
// Note: If you want apiserver output on stdout set
// KUBEBUILDER_ATTACH_CONTROL_PLANE_OUTPUT to true when running tests.
func Setup(t *testing.T, scheme *runtime.Scheme) *rest.Config {
	t.Helper()

	gatewayCRDPath := filepath.Join(build.Default.GOPATH, "pkg", "mod", "sigs.k8s.io", "gateway-api@"+consts.GatewayAPIVersion, "config", "crd", "experimental")
	testEnv := &envtest.Environment{
		ControlPlaneStopTimeout: time.Second * 60,
		CRDDirectoryPaths: []string{
			gatewayCRDPath,
		},
		CRDInstallOptions: envtest.CRDInstallOptions{
			CleanUpAfterUse: false,
			Scheme:          scheme,
		},
		Scheme: scheme,
	}

	t.Logf("starting envtest environment...")
	cfg, err := testEnv.Start()
	require.NoError(t, err)

	t.Logf("waiting for Gateway API CRDs to be available...")
	require.NoError(t, envtest.WaitForCRDs(cfg, []*apiextensionsv1.CustomResourceDefinition{
		{
			Spec: apiextensionsv1.CustomResourceDefinitionSpec{
				Group: gatewayv1beta1.GroupVersion.Group,
				Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
					{
						Name:   gatewayv1beta1.GroupVersion.Version,
						Served: true,
					},
				},
				Names: apiextensionsv1.CustomResourceDefinitionNames{
					Plural: "gateways",
				},
			},
		},
		{
			Spec: apiextensionsv1.CustomResourceDefinitionSpec{
				Group: gatewayv1beta1.GroupVersion.Group,
				Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
					{
						Name:   gatewayv1beta1.GroupVersion.Version,
						Served: true,
					},
				},
				Names: apiextensionsv1.CustomResourceDefinitionNames{
					Plural: "httproutes",
				},
			},
		},
		{
			Spec: apiextensionsv1.CustomResourceDefinitionSpec{
				Group: gatewayv1beta1.GroupVersion.Group,
				Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
					{
						Name:   gatewayv1beta1.GroupVersion.Version,
						Served: true,
					},
				},
				Names: apiextensionsv1.CustomResourceDefinitionNames{
					Plural: "referencegrants",
				},
			},
		},
	}, envtest.CRDInstallOptions{}))

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

	t.Cleanup(func() {
		t.Logf("stopping envtest environment for test %s", t.Name())
		close(done)
		wg.Wait()
	})

	return cfg
}
