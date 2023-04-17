//go:build conformance_tests
// +build conformance_tests

package conformance

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/go-logr/logr"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/metallb"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/kind"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"

	testutils "github.com/kong/kubernetes-ingress-controller/v2/internal/util/test"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/testenv"
)

var (
	existingCluster = os.Getenv("KONG_TEST_CLUSTER")
	ingressClass    = "kong-conformance-tests"

	env                    environments.Environment
	ctx                    context.Context
	globalDeprecatedLogger logrus.FieldLogger
	globalLogger           logr.Logger
)

func TestMain(m *testing.M) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	// Logger needs to be configured before anything else happens.
	// This is because the controller manager has a timeout for
	// logger initialization, and if the logger isn't configured
	// after 30s from the start of controller manager package init function,
	// the controller manager will set up a no op logger and continue.
	// The logger cannot be configured after that point.
	deprecatedLogger, logger, logOutput, err := testutils.SetupLoggers("trace", "text", false)
	if err != nil {
		exitOnErr(fmt.Errorf("failed to setup loggers: %w", err))
	}
	if logOutput != "" {
		fmt.Printf("INFO: writing manager logs to %s\n", logOutput)
	}
	globalDeprecatedLogger = deprecatedLogger
	globalLogger = logger

	kongAddon := kong.NewBuilder().WithControllerDisabled().WithProxyAdminServiceTypeLoadBalancer().Build()
	builder := environments.NewBuilder().WithAddons(metallb.New(), kongAddon)
	useExistingClusterIfPresent(builder)

	env, err = builder.Build(ctx)
	exitOnErr(err)

	defer func() {
		if existingCluster == "" {
			exitOnErr(env.Cleanup(ctx))
		}
	}()

	fmt.Println("INFO: waiting for cluster and addons to be ready")
	envReadyCtx, envReadyCancel := context.WithTimeout(ctx, testenv.EnvironmentReadyTimeout())
	defer envReadyCancel()
	exitOnErr(<-env.WaitForReady(envReadyCtx))

	// To allow running conformance tests in a loop to e.g. detect flaky tests
	// let's ensure that conformance related namespaced are deleted from the cluster.
	exitOnErr(ensureConformanceTestsNamespacesAreNotPresent(ctx, env.Cluster().Client()))

	code := m.Run()

	os.Exit(code)
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
	if existingCluster != "" {
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
