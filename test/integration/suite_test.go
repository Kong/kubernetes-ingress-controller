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
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/knative"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/metallb"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/gke"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/kind"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	testutils "github.com/kong/kubernetes-ingress-controller/v2/internal/util/test"
)

// For some reason we've ended up using kind v0.15.0 (which by default deploys k8s v1.25)
// even though that
//   - our CI runners use ubuntu-latest (which at the time of writing this comment was ubuntu20.04)
//     which uses kind v0.14 https://github.com/actions/runner-images/blob/main/images/linux/Ubuntu2004-Readme.md
//   - when used with ktf (which at the time of writing this comment was set to v0.19.0) we
//     should use kind v0.14 since that what ktf has set as dependency
//
// With all that said, we still managed to get kind v0.15 on our CI
// https://github.com/Kong/kubernetes-ingress-controller/runs/8211490522?check_suite_focus=true#step:5:6
// which causes issues down the line (metallb manifests using PSP which is not available
// in k8s v1.25+).
var defaultKindClusterVesion = semver.Version{
	Major: 1,
	Minor: 24,
	Patch: 4,
}

var k8sClient *kubernetes.Clientset

// -----------------------------------------------------------------------------
// Testing Main
// -----------------------------------------------------------------------------

// generateKongBuilder returns a Kong KTF addon builder and a string slice of controller arguments needed to interact
// with the addon.
func generateKongBuilder() (*kong.Builder, []string) {
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

	if kongImage != "" {
		if kongTag == "" {
			exitOnErrWithCode(fmt.Errorf("TEST_KONG_IMAGE requires TEST_KONG_TAG"), ExitCodeEnvSetupFailed)
		}
		kongbuilder = kongbuilder.WithProxyImage(kongImage, kongTag)
	}

	if kongPullUsername != "" || kongPullPassword != "" {
		if kongPullPassword == "" || kongPullUsername == "" {
			exitOnErrWithCode(fmt.Errorf("TEST_KONG_PULL_USERNAME requires TEST_KONG_PULL_PASSWORD"), ExitCodeEnvSetupFailed)
		}
		kongbuilder = kongbuilder.WithProxyImagePullSecret("", kongPullUsername, kongPullPassword, "")
	}

	if dbmode == "postgres" {
		kongbuilder = kongbuilder.WithPostgreSQL()
	}

	kongbuilder = kongbuilder.WithProxyEnvVar("router_flavor", "traditional")

	kongbuilder.WithControllerDisabled()

	return kongbuilder, extraControllerArgs
}

func TestMain(m *testing.M) {
	var code int
	defer func() {
		os.Exit(code)
	}()
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("INFO: setting up test environment")
	kongbuilder, extraControllerArgs := generateKongBuilder()
	kongAddon := kongbuilder.Build()
	builder := environments.NewBuilder().WithAddons(kongAddon)

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

		clusterVersion = defaultKindClusterVesion
		if clusterVersionStr != "" {
			var err error
			clusterVersion, err = semver.Parse(strings.TrimPrefix(clusterVersionStr, "v"))
			exitOnErr(err)

		}
		fmt.Printf("INFO: build a new KIND cluster with version %s\n", clusterVersion.String())
		builder.WithKubernetesVersion(clusterVersion)
	}

	fmt.Println("INFO: building test environment")
	var err error
	env, err = builder.Build(ctx)
	exitOnErr(err)
	k8sClient = env.Cluster().Client()

	cleaner := clusters.NewCleaner(env.Cluster())
	defer cleaner.Cleanup(ctx) //nolint:errcheck

	fmt.Printf("INFO: reconfiguring the kong admin service as LoadBalancer type\n")
	svc, err := env.Cluster().Client().CoreV1().Services(kongAddon.Namespace()).Get(ctx, kong.DefaultAdminServiceName, metav1.GetOptions{})
	exitOnErr(err)
	svc.Spec.Type = corev1.ServiceTypeLoadBalancer
	_, err = env.Cluster().Client().CoreV1().Services(kongAddon.Namespace()).Update(ctx, svc, metav1.UpdateOptions{})
	exitOnErr(err)

	clusterVersion, err = env.Cluster().Version()
	exitOnErr(err)
	if clusterVersion.GE(knativeMinKubernetesVersion) {
		fmt.Println("INFO: deploying knative addon")
		knativeBuilder := knative.NewBuilder()
		knativeAddon := knativeBuilder.Build()
		exitOnErr(env.Cluster().DeployAddon(ctx, knativeAddon))
	}

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
		fmt.Println("INFO: configuring feature gates")
		if controllerFeatureGates == "" {
			controllerFeatureGates = defaultFeatureGates
		}
		fmt.Printf("INFO: feature gates enabled: %s\n", controllerFeatureGates)
		fmt.Println("INFO: starting the controller manager")
		standardControllerArgs := []string{
			fmt.Sprintf("--ingress-class=%s", ingressClass),
			fmt.Sprintf("--admission-webhook-cert=%s", testutils.KongSystemServiceCert),
			fmt.Sprintf("--admission-webhook-key=%s", testutils.KongSystemServiceKey),
			fmt.Sprintf("--admission-webhook-listen=0.0.0.0:%d", testutils.AdmissionWebhookListenPort),
			"--profiling",
			"--dump-config",
			"--log-level=trace",
			"--debug-log-reduce-redundancy",
			fmt.Sprintf("--feature-gates=%s", controllerFeatureGates),
			fmt.Sprintf("--election-namespace=%s", kongAddon.Namespace()),
		}
		allControllerArgs := append(standardControllerArgs, extraControllerArgs...)
		exitOnErr(testutils.DeployControllerManagerForCluster(ctx, env.Cluster(), allControllerArgs...))
	}

	gatewayClient, err := gatewayclient.NewForConfig(env.Cluster().Config())
	exitOnErr(err)

	fmt.Println("INFO: Deploying the default GatewayClass")
	gwc, err := DeployGatewayClass(ctx, gatewayClient, unmanagedGatewayClassName)
	exitOnErr(err)
	cleaner.Add(gwc)

	fmt.Println("INFO: Deploying the controller's IngressClass")
	iclass := &netv1.IngressClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: ingressClass,
		},
		Spec: netv1.IngressClassSpec{
			Controller: store.IngressClassKongController,
		},
	}
	iclass, err = env.Cluster().Client().NetworkingV1().IngressClasses().Create(ctx, iclass, metav1.CreateOptions{})
	exitOnErr(err)
	defer func() {
		// deleting this directly instead of adding it to the cleaner because the cleaner always gets a 404 on it for
		// unknown reasons
		_ = env.Cluster().Client().NetworkingV1().IngressClasses().Delete(ctx, iclass.Name, metav1.DeleteOptions{})
	}()

	fmt.Printf("INFO: testing environment is ready KUBERNETES_VERSION=(%v): running tests\n", clusterVersion)
	code = m.Run()

	if keepTestCluster == "" && existingCluster == "" {
		ctx, cancel := context.WithTimeout(context.Background(), environmentCleanupTimeout)
		defer cancel()
		fmt.Printf("INFO: cluster %s is being deleted\n", env.Cluster().Name())
		exitOnErr(env.Cleanup(ctx))
	}
}
