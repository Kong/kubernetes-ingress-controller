//go:build integration_tests
// +build integration_tests

package integration

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/blang/semver/v4"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/knative"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/metallb"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/gke"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/kind"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"

	testutils "github.com/kong/kubernetes-ingress-controller/v2/internal/test/util"
)

// -----------------------------------------------------------------------------
// Testing Main
// -----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("INFO: setting up test environment")
	kongbuilder := kong.NewBuilder()
	extraControllerArgs := []string{}
	if kongEnterpriseEnabled == "true" {
		licenseJSON, err := kong.GetLicenseJSONFromEnv()
		exitOnErr(err)
		extraControllerArgs = append(extraControllerArgs, fmt.Sprintf("--kong-admin-token=%s", kongTestPassword))
		extraControllerArgs = append(extraControllerArgs, "--kong-workspace=notdefault")
		kongbuilder = kongbuilder.WithProxyEnterpriseEnabled(licenseJSON).
			WithProxyEnterpriseSuperAdminPassword(kongTestPassword).
			WithProxyAdminServiceTypeLoadBalancer()
	}

	if dbmode == "postgres" {
		kongbuilder = kongbuilder.WithPostgreSQL()
	}

	kongbuilder.WithControllerDisabled()
	kongAddon := kongbuilder.Build()
	builder := environments.NewBuilder().WithAddons(kongAddon, knative.New())

	fmt.Println("INFO: configuring cluster for testing environment")
	if existingCluster != "" {
		if clusterVersionStr != "" {
			exitOnErrWithCode(fmt.Errorf("can't flag cluster version & provide an existing cluster at the same time"), ExitCodeIncompatibleOptions)
		}
		clusterParts := strings.Split(existingCluster, ":")
		if len(clusterParts) != 2 {
			exitOnErrWithCode(fmt.Errorf("existing cluster in wrong format (%s): format is <TYPE>:<NAME> (e.g. kind:test-cluster)", existingCluster), ExitCodeCantUseExistingCluster)
		}
		clusterType, clusterName := clusterParts[0], clusterParts[1]

		fmt.Printf("INFO: using existing %s cluster %s\n", clusterType, clusterName)
		switch clusterType {
		case string(kind.KindClusterType):
			cluster, err := kind.NewFromExisting(clusterName)
			exitOnErr(err)
			builder.WithExistingCluster(cluster)
			builder.WithAddons(metallb.New())
		case string(gke.GKEClusterType):
			cluster, err := gke.NewFromExistingWithEnv(ctx, clusterName)
			exitOnErr(err)
			builder.WithExistingCluster(cluster)
		default:
			exitOnErrWithCode(fmt.Errorf("unknown cluster type: %s", clusterType), ExitCodeCantUseExistingCluster)
		}
	} else {
		fmt.Println("INFO: no existing cluster found, deploying using Kubernetes In Docker (KIND)")
		builder.WithAddons(metallb.New())
	}
	if clusterVersionStr != "" {
		clusterVersion, err := semver.Parse(strings.TrimPrefix(clusterVersionStr, "v"))
		exitOnErr(err)
		cluster, err := kind.NewBuilder().WithClusterVersion(clusterVersion).Build(ctx)
		exitOnErr(err)
		builder.WithExistingCluster(cluster)
	}

	fmt.Println("INFO: building test environment")
	var err error
	env, err = builder.Build(ctx)
	exitOnErr(err)

	fmt.Printf("INFO: reconfiguring the kong admin service as LoadBalancer type\n")
	svc, err := env.Cluster().Client().CoreV1().Services(kongAddon.Namespace()).Get(ctx, kong.DefaultAdminServiceName, metav1.GetOptions{})
	exitOnErr(err)
	svc.Spec.Type = corev1.ServiceTypeLoadBalancer
	_, err = env.Cluster().Client().CoreV1().Services(kongAddon.Namespace()).Update(ctx, svc, metav1.UpdateOptions{})
	exitOnErr(err)

	fmt.Printf("INFO: waiting for cluster %s and all addons to become ready\n", env.Cluster().Name())
	exitOnErr(<-env.WaitForReady(ctx))

	fmt.Println("INFO: collecting urls from the kong proxy deployment")
	proxyURL, err = kongAddon.ProxyURL(ctx, env.Cluster())
	exitOnErr(err)
	proxyAdminURL, err = kongAddon.ProxyAdminURL(ctx, env.Cluster())
	exitOnErr(err)
	proxyUDPURL, err = kongAddon.ProxyUDPURL(ctx, env.Cluster())
	exitOnErr(err)

	fmt.Println("INFO: generating unique namespaces for each test case")
	testCases, err := identifyTestCasesForDir("./")
	exitOnErr(err)
	for _, testCase := range testCases {
		namespaceForTestCase, err := clusters.GenerateNamespace(ctx, env.Cluster(), testCase)
		exitOnErr(err)
		namespaces[testCase] = namespaceForTestCase
		watchNamespaces = fmt.Sprintf("%s,%s", watchNamespaces, namespaceForTestCase.Name)
	}

	if v := os.Getenv("KONG_BRING_MY_OWN_KIC"); v == "true" {
		fmt.Println("WARNING: caller indicated that they will manage their own controller")
	} else {
		fmt.Println("INFO: creating additional controller namespaces")
		ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: controllerNamespace}}
		if _, err := env.Cluster().Client().CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{}); err != nil {
			if !errors.IsAlreadyExists(err) {
				exitOnErr(err)
			}
		}
		fmt.Println("INFO: starting the controller manager")
		standardControllerArgs := []string{
			fmt.Sprintf("--ingress-class=%s", ingressClass),
			fmt.Sprintf("--admission-webhook-cert=%s", admissionWebhookCert),
			fmt.Sprintf("--admission-webhook-key=%s", admissionWebhookKey),
			fmt.Sprintf("--watch-namespace=%s", watchNamespaces),
			"--admission-webhook-listen=172.17.0.1:49023",
			"--profiling",
			"--dump-config",
			"--log-level=trace",
			"--debug-log-reduce-redundancy",
			"--feature-gates=Gateway=true",
		}
		allControllerArgs := append(standardControllerArgs, extraControllerArgs...)
		exitOnErr(testutils.DeployControllerManagerForCluster(ctx, env.Cluster(), allControllerArgs...))
	}

	fmt.Println("INFO: running final testing environment checks")
	clusterVersion, err = env.Cluster().Version()
	exitOnErr(err)

	fmt.Printf("INFO: testing environment is ready KUBERNETES_VERSION=(%v): running tests\n", clusterVersion)
	code := m.Run()

	if keepTestCluster == "" && existingCluster == "" {
		ctx, cancel := context.WithTimeout(context.Background(), environmentCleanupTimeout)
		defer cancel()
		fmt.Printf("INFO: cluster %s is being deleted\n", env.Cluster().Name())
		exitOnErr(env.Cleanup(ctx))
	}

	os.Exit(code)
}
