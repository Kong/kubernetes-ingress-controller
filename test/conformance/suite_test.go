//go:build conformance_tests

package conformance

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/metallb"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/kind"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/gateway"
	dpconf "github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/config"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	testutils "github.com/kong/kubernetes-ingress-controller/v3/internal/util/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers/certificate"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testenv"
)

var (
	conformanceTestsBaseManifests = fmt.Sprintf("%s/conformance/base/manifests.yaml", consts.GatewayRawRepoURL)

	env          environments.Environment
	ctx          context.Context
	globalLogger logr.Logger
)

func TestMain(m *testing.M) {
	var code int
	defer func() {
		os.Exit(code)
	}()
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	// Logger needs to be configured before anything else happens.
	// This is because the controller manager has a timeout for
	// logger initialization, and if the logger isn't configured
	// after 30s from the start of controller manager package init function,
	// the controller manager will set up a no op logger and continue.
	// The logger cannot be configured after that point.
	logger, logOutput, err := testutils.SetupLoggers("trace", "text")
	if err != nil {
		exitOnErr(fmt.Errorf("failed to setup loggers: %w", err))
	}
	if logOutput != "" {
		fmt.Printf("INFO: writing manager logs to %s\n", logOutput)
	}
	globalLogger = logger

	// In order to pass conformance tests, the expression router is required.
	kongBuilder := kong.NewBuilder().WithControllerDisabled().WithProxyAdminServiceTypeLoadBalancer().
		WithNamespace(consts.ControllerNamespace)
	if testenv.KongRouterFlavor() == dpconf.RouterFlavorExpressions {
		fmt.Println("INFO: expression routes enabled")
		kongBuilder = kongBuilder.WithProxyEnvVar("router_flavor", string(dpconf.RouterFlavorExpressions))
	}

	// The test cases for GRPCRoute in the current GatewayAPI all use the h2c protocol.
	// In order to pass conformance tests, the proxy must listen http2 and http on the same port.
	kongBuilder.WithProxyEnvVar("PROXY_LISTEN", `0.0.0.0:8000 http2\, 0.0.0.0:8443 http2 ssl`)

	// Pin the Helm chart version.
	kongBuilder.WithHelmChartVersion(testenv.KongHelmChartVersion())

	kongAddon := kongBuilder.Build()
	builder := environments.NewBuilder().WithAddons(metallb.New(), kongAddon)
	useExistingClusterIfPresent(builder)

	env, err = builder.Build(ctx)
	exitOnErr(err)

	fmt.Println("INFO: waiting for cluster and addons to be ready")
	envReadyCtx, envReadyCancel := context.WithTimeout(ctx, testenv.EnvironmentReadyTimeout())
	defer envReadyCancel()
	exitOnErr(<-env.WaitForReady(envReadyCtx))

	// To allow running conformance tests in a loop to e.g. detect flaky tests
	// let's ensure that conformance related namespaced are deleted from the cluster.
	exitOnErr(ensureConformanceTestsNamespacesAreNotPresent(ctx, env.Cluster().Client()))

	code = m.Run()
	if testenv.IsCI() {
		fmt.Printf("INFO: running in ephemeral CI environment, skipping cluster %s teardown\n", env.Cluster().Name())
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), test.EnvironmentCleanupTimeout)
		defer cancel()
		exitOnErr(helpers.RemoveCluster(ctx, env.Cluster()))
	}
}

// prepareEnvForGatewayConformanceTests prepares the environment for running the gateway conformance test suite
// from the gateway-api project. Used as a helper function for both the stable and experimental conformance tests.
func prepareEnvForGatewayConformanceTests(t *testing.T) (c client.Client, gatewayClassName string) {
	t.Log("configuring environment for gateway conformance tests")
	client, err := client.New(env.Cluster().Config(), client.Options{})
	require.NoError(t, err)
	require.NoError(t, gatewayv1alpha2.Install(client.Scheme()))
	require.NoError(t, gatewayv1beta1.Install(client.Scheme()))
	require.NoError(t, gatewayv1.Install(client.Scheme()))
	require.NoError(t, apiextensionsv1.AddToScheme(client.Scheme()))

	featureGateFlag := fmt.Sprintf("--feature-gates=%s", consts.DefaultFeatureGates)

	t.Log("Preparing the environment to run the controller manager")
	require.NoError(t, testutils.PrepareClusterForRunningControllerManager(ctx, env.Cluster()))

	t.Log("starting the controller manager")
	cert, key := certificate.GetKongSystemSelfSignedCerts()
	args := []string{
		"--ingress-class=kong-conformance-tests",
		fmt.Sprintf("--admission-webhook-cert=%s", cert),
		fmt.Sprintf("--admission-webhook-key=%s", key),
		fmt.Sprintf("--admission-webhook-listen=%s:%d", testutils.GetAdmissionWebhookListenHost(), testutils.AdmissionWebhookListenPort),
		"--profiling",
		"--dump-config",
		"--log-level=trace",
		featureGateFlag,
		"--anonymous-reports=false",
	}
	cancel, err := testutils.DeployControllerManagerForCluster(ctx, globalLogger, env.Cluster(), nil, args)
	require.NoError(t, err)
	t.Cleanup(func() { cancel() })

	t.Log("creating GatewayClass for gateway conformance tests")
	gatewayClass := &gatewayapi.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
			Annotations: map[string]string{
				annotations.AnnotationPrefix + annotations.GatewayClassUnmanagedKey: annotations.GatewayClassUnmanagedAnnotationValuePlaceholder,
			},
		},
		Spec: gatewayapi.GatewayClassSpec{
			ControllerName: gateway.GetControllerName(),
		},
	}
	require.NoError(t, client.Create(ctx, gatewayClass))
	t.Cleanup(func() { require.NoError(t, client.Delete(ctx, gatewayClass)) })
	t.Logf("created GatewayClass %q", gatewayClass.Name)

	return client, gatewayClass.Name
}

func ensureConformanceTestsNamespacesAreNotPresent(ctx context.Context, client *kubernetes.Clientset) error {
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(3) //nolint:gomnd

	for _, namespace := range []string{"gateway-conformance-infra", "gateway-conformance-web-backend", "gateway-conformance-app-backend"} {
		namespace := namespace
		g.Go(func() error {
			return ensureNamespaceDeleted(ctx, namespace, client)
		})
	}
	return g.Wait()
}

func ensureNamespaceDeleted(ctx context.Context, ns string, client *kubernetes.Clientset) error {
	namespaceClient := client.CoreV1().Namespaces()
	namespace, err := namespaceClient.Get(ctx, ns, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			// If the namespace cannot be found then we're good to go.
			return nil
		}
		return err
	}

	if namespace.Status.Phase == corev1.NamespaceActive {
		fmt.Printf("INFO: removing %s namespace for clean test run\n", ns)
		err := namespaceClient.Delete(ctx, ns, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}

	w, err := namespaceClient.Watch(ctx, metav1.ListOptions{
		LabelSelector: "kubernetes.io/metadata.name=" + ns,
	})
	if err != nil {
		return err
	}

	defer w.Stop()
	for {
		select {
		case event := <-w.ResultChan():
			if event.Type == watch.Deleted {
				return nil
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func useExistingClusterIfPresent(builder *environments.Builder) {
	if existingCluster := testenv.ExistingClusterName(); existingCluster != "" {
		parts := strings.Split(existingCluster, ":")
		if len(parts) != 2 {
			exitOnErr(fmt.Errorf("%s is not a valid value for KONG_TEST_CLUSTER", existingCluster))
		}
		if parts[0] != "kind" {
			exitOnErr(fmt.Errorf("%s is not a supported cluster type for this test suite yet", parts[0]))
		}
		cluster, err := kind.NewFromExisting(parts[1])
		exitOnErr(err)
		fmt.Printf("INFO: using existing kind cluster for test (name: %s)\n", parts[1])
		builder.WithExistingCluster(cluster)
	} else {
		fmt.Println("INFO: creating new kind cluster for conformance tests")
	}
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

// patchGatewayClassToPassTestGatewayClassObservedGenerationBump - wait for the GatewayClass
// (call is blocking run in a goroutine)  created by the test GatewayClassObservedGenerationBump
// and patch it with the unmanaged annotation to make it reconciled by the GatewayClass controller.
// The timeout and the tick are pretty loose because of the nondeterministic test order execution.
func patchGatewayClassToPassTestGatewayClassObservedGenerationBump(ctx context.Context, t *testing.T, k8sClient client.Client) {
	ensureTestGatewayClassIsUnmanaged := func(ctx context.Context, k8sClient client.Client) bool {
		gwcNamespacedName := k8stypes.NamespacedName{Name: "gatewayclass-observed-generation-bump"}
		gwc := &gatewayapi.GatewayClass{}
		if err := k8sClient.Get(ctx, gwcNamespacedName, gwc); err != nil {
			return false
		}
		if gwc.Annotations == nil {
			gwc.Annotations = map[string]string{}
		}
		gwc.Annotations[annotations.AnnotationPrefix+annotations.GatewayClassUnmanagedKey] = annotations.GatewayClassUnmanagedAnnotationValuePlaceholder
		if err := k8sClient.Update(ctx, gwc); err != nil {
			return false
		}
		return true
	}

	require.Eventually(t, func() bool {
		return ensureTestGatewayClassIsUnmanaged(ctx, k8sClient)
	}, 10*time.Minute, test.RequestTimeout)
}
