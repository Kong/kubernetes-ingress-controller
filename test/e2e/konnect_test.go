//go:build e2e_tests

package e2e

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os/exec"
	"sync"
	"testing"
	"time"

	environment "github.com/kong/kubernetes-testing-framework/pkg/environments"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/konnect"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/konnect/nodes"
	testkonnect "github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers/konnect"
)

const (
	konnectControlPlaneAdminAPIBaseURL = "https://us.kic.api.konghq.tech"

	konnectNodeRegistrationTimeout = 5 * time.Minute
	konnectNodeRegistrationCheck   = 30 * time.Second
)

func TestKonnectConfigPush(t *testing.T) {
	t.Parallel()
	testkonnect.SkipIfMissingRequiredKonnectEnvVariables(t)

	ctx, env := setupE2ETest(t)

	cpID := testkonnect.CreateTestControlPlane(ctx, t)
	cert, key := testkonnect.CreateClientCertificate(ctx, t, cpID)
	createKonnectClientSecretAndConfigMap(ctx, t, env, cert, key, cpID)

	deployments := deployAllInOneKonnectManifest(ctx, t, env)

	t.Log("running ingress tests to verify all-in-one deployed ingress controller and proxy are functional")
	deployIngressWithEchoBackends(ctx, t, env, numberOfEchoBackends)
	verifyIngressWithEchoBackends(ctx, t, env, numberOfEchoBackends)

	t.Log("ensuring ingress resources are correctly populated in Konnect Control Plane's Admin API")
	konnectAdminAPIClient := createKonnectAdminAPIClient(t, cpID, cert, key)
	verifyIngressWithEchoBackendsInAdminAPI(ctx, t, konnectAdminAPIClient.AdminAPIClient(), numberOfEchoBackends)

	t.Log("ensuring KIC nodes and controlled kong gateway nodes are present in konnect control plane")
	requireKonnectNodesConsistentWithK8s(ctx, t, env, deployments, cpID, cert, key)
	requireAllProxyReplicasIDsConsistentWithKonnect(ctx, t, env, deployments.ProxyNN, cpID, cert, key)
}

func TestKonnectLicenseActivation(t *testing.T) {
	t.Parallel()
	testkonnect.SkipIfMissingRequiredKonnectEnvVariables(t)

	ctx, env := setupE2ETest(t)

	rgID := testkonnect.CreateTestControlPlane(ctx, t)
	cert, key := testkonnect.CreateClientCertificate(ctx, t, rgID)
	createKonnectClientSecretAndConfigMap(ctx, t, env, cert, key, rgID)

	const manifestFile = "manifests/all-in-one-dbless-konnect-enterprise.yaml"
	ManifestDeploy{Path: manifestFile}.Run(ctx, t, env)

	exposeAdminAPI(ctx, t, env, k8stypes.NamespacedName{Namespace: "kong", Name: "proxy-kong"})

	t.Log("disabling license management")
	kubeconfig := getTemporaryKubeconfig(t, env)
	require.NoError(t, setEnv(setEnvParams{
		kubeCfgPath:   kubeconfig,
		namespace:     namespace,
		target:        fmt.Sprintf("deployment/%s", controllerDeploymentName),
		containerName: controllerContainerName,
		variableName:  "CONTROLLER_KONNECT_LICENSING_ENABLED",
		value:         "",
	}))

	t.Log("restarting proxy")
	cmd := exec.Command("kubectl", "--kubeconfig", kubeconfig, "rollout", "-n", "kong", "restart", "deployment", "proxy-kong")
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	require.NoErrorf(t, err, "restarting proxy failed: STDOUT(%s) STDERR(%s)", stdout.String(), stderr.String())

	t.Log("confirming that the license is empty")
	require.Eventually(t, func() bool {
		license, err := getLicenseFromAdminAPI(ctx, env, "")
		if err != nil {
			t.Logf("failed to get license: %v", err)
			return false
		}
		return license.License.Expiration == ""
	}, adminAPIWait, time.Second)

	t.Log("re-enabling license management")
	require.NoError(t, setEnv(setEnvParams{
		kubeCfgPath:   kubeconfig,
		namespace:     namespace,
		target:        fmt.Sprintf("deployment/%s", controllerDeploymentName),
		containerName: controllerContainerName,
		variableName:  "CONTROLLER_KONNECT_LICENSING_ENABLED",
		value:         "true",
	}))

	t.Log("confirming that the license is set")
	assert.Eventually(t, func() bool {
		license, err := getLicenseFromAdminAPI(ctx, env, "")
		if err != nil {
			t.Logf("failed to get license: %v", err)
			return false
		}
		return license.License.Expiration != ""
	}, adminAPIWait, time.Second)
	t.Log("done")
}

func TestKonnectWhenMisconfiguredBasicIngressNotAffected(t *testing.T) {
	t.Parallel()
	testkonnect.SkipIfMissingRequiredKonnectEnvVariables(t)
	ctx, env := setupE2ETest(t)

	rgID := testkonnect.CreateTestControlPlane(ctx, t)
	cert, key := testkonnect.CreateClientCertificate(ctx, t, rgID)

	// create a Konnect client secret and config map with a non-existing control plane ID to simulate misconfiguration
	notExistingRgID := "not-existing-cp-id"
	createKonnectClientSecretAndConfigMap(ctx, t, env, cert, key, notExistingRgID)

	deployAllInOneKonnectManifest(ctx, t, env)

	t.Log("running ingress tests to verify misconfiguration doesn't affect basic ingress functionality")
	deployIngressWithEchoBackends(ctx, t, env, numberOfEchoBackends)
	verifyIngressWithEchoBackends(ctx, t, env, numberOfEchoBackends)
}

// deployAllInOneKonnectManifest deploys all-in-one-dbless-konnect.yaml manifest, replacing the controller image
// if specified by environment variables.
func deployAllInOneKonnectManifest(ctx context.Context, t *testing.T, env environment.Environment) Deployments {
	const manifestFile = "manifests/all-in-one-dbless-konnect.yaml"
	t.Logf("deploying %s manifest file", manifestFile)

	return ManifestDeploy{Path: manifestFile}.Run(ctx, t, env)
}

// createKonnectClientSecretAndConfigMap creates a Secret with client TLS certificate that is used by KIC to communicate
// with Konnect Admin API. It also creates a ConfigMap that specifies a Control Plane ID and Konnect Admin API URL.
// Both Secret and ConfigMap are used by all-in-one-dbless-konnect.yaml manifest and need to be populated before
// deploying it.
func createKonnectClientSecretAndConfigMap(ctx context.Context, t *testing.T, env environment.Environment, tlsCert, tlsKey, rgID string) {
	t.Helper()

	// create a namespace in case it doesn't exist yet
	t.Log("creating kong namespace")
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}
	_, err := env.Cluster().Client().CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	if !apierrors.IsAlreadyExists(err) {
		require.NoError(t, err)
	}

	t.Log("creating konnect client tls secret")
	_, err = env.Cluster().Client().CoreV1().Secrets(namespace).Create(ctx, &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "konnect-client-tls",
		},
		Data: map[string][]byte{
			"tls.crt": []byte(tlsCert),
			"tls.key": []byte(tlsKey),
		},
		Type: corev1.SecretTypeTLS,
	}, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("creating konnect config map")
	_, err = env.Cluster().Client().CoreV1().ConfigMaps(namespace).Create(ctx, &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "konnect-config",
		},
		Data: map[string]string{
			"CONTROLLER_KONNECT_CONTROL_PLANE_ID": rgID,
			"CONTROLLER_KONNECT_ADDRESS":          konnectControlPlaneAdminAPIBaseURL,
		},
	}, metav1.CreateOptions{})
	require.NoError(t, err)
}

// createKonnectAdminAPIClient creates an *kong.Client that will communicate with Konnect Control Plane's Admin API.
func createKonnectAdminAPIClient(t *testing.T, rgID, cert, key string) *adminapi.KonnectClient {
	t.Helper()

	c, err := adminapi.NewKongClientForKonnectControlPlane(adminapi.KonnectConfig{
		ControlPlaneID: rgID,
		Address:        konnectControlPlaneAdminAPIBaseURL,
		TLSClient: adminapi.TLSClientConfig{
			Cert: cert,
			Key:  key,
		},
	})
	require.NoError(t, err)
	return c
}

// createKonnectNodeClient creates a konnect.NodeClient to get nodes in konnect control plane.
func createKonnectNodeClient(t *testing.T, rgID, cert, key string) *nodes.Client {
	cfg := adminapi.KonnectConfig{
		ConfigSynchronizationEnabled: true,
		ControlPlaneID:               rgID,
		Address:                      konnectControlPlaneAdminAPIBaseURL,
		RefreshNodePeriod:            konnect.MinRefreshNodePeriod,
		TLSClient: adminapi.TLSClientConfig{
			Cert: cert,
			Key:  key,
		},
	}
	c, err := nodes.NewClient(cfg)
	require.NoError(t, err)
	return c
}

func requireKonnectNodesConsistentWithK8s(ctx context.Context, t *testing.T, env environment.Environment, deployments Deployments, rgID string, cert, key string) {
	konnectNodeClient := createKonnectNodeClient(t, rgID, cert, key)
	require.Eventually(t, func() bool {
		ns, err := konnectNodeClient.ListAllNodes(ctx)
		if err != nil {
			t.Logf("list all nodes failed: %v", err)
			return false
		}

		kicPods, err := listPodsByLabels(ctx, env, "kong", map[string]string{"app": deployments.ControllerNN.Name})
		if err != nil || len(kicPods) != 1 {
			return false
		}

		kongPods, err := listPodsByLabels(ctx, env, "kong", map[string]string{"app": deployments.ProxyNN.Name})
		if err != nil || len(kongPods) != 2 {
			return false
		}

		kicNodes := []*nodes.NodeItem{}
		kongNodes := []*nodes.NodeItem{}

		for _, node := range ns {
			if node.Type == nodes.NodeTypeIngressController {
				kicNodes = append(kicNodes, node)
			}
			if node.Type == nodes.NodeTypeKongProxy {
				kongNodes = append(kongNodes, node)
			}
		}

		// check for number of nodes in Konnect.
		if len(kicNodes) != 1 || len(kongNodes) != 2 {
			return false
		}

		if kicNodes[0].Hostname != fmt.Sprintf("%s/%s", kicPods[0].Namespace, kicPods[0].Name) {
			return false
		}

		for _, pod := range kongPods {
			nsName := fmt.Sprintf("%s/%s", pod.Namespace, pod.Name)
			if !lo.ContainsBy(kongNodes, func(n *nodes.NodeItem) bool {
				return n.Hostname == nsName
			}) {
				return false
			}
		}

		return true
	}, konnectNodeRegistrationTimeout, konnectNodeRegistrationCheck)
}

// requireAllProxyReplicasIDsConsistentWithKonnect ensures that all proxy replicas are registered in Konnect's Node API
// with their respective Admin API Node IDs.
// It's required because when a proxy replica connects with Konnect (e.g. to report Analytics data), it uses its locally
// generated Node ID (KIC knows it via calling gateway's Admin API) to identify itself. If the Node is not registered
// in Konnect using the same ID, it won't be possible to associate requests with the correct node.
func requireAllProxyReplicasIDsConsistentWithKonnect(
	ctx context.Context,
	t *testing.T,
	env environment.Environment,
	proxyDeploymentNN k8stypes.NamespacedName,
	rg, cert, key string,
) {
	pods, err := listPodsByLabels(ctx, env, proxyDeploymentNN.Namespace, map[string]string{"app": proxyDeploymentNN.Name})
	require.NoError(t, err)

	nodeAPIClient := createKonnectNodeClient(t, rg, cert, key)

	getNodeIDFromAdminAPI := func(proxyPod corev1.Pod) string {
		client := &http.Client{
			Timeout: time.Second * 30,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}

		forwardCtx, cancel := context.WithCancel(ctx)
		defer cancel()
		localPort := startPortForwarder(forwardCtx, t, env, proxyDeploymentNN.Namespace, proxyPod.Name, "8444")
		address := fmt.Sprintf("https://localhost:%d", localPort)

		kongClient, err := adminapi.NewKongAPIClient(address, client)
		require.NoError(t, err)

		nodeID, err := adminapi.NewClient(kongClient).NodeID(ctx)
		require.NoError(t, err)
		return nodeID
	}

	t.Logf("ensuring all %d proxy replicas have consistent IDs assigned in Node API", len(pods))
	wg := sync.WaitGroup{}
	for _, pod := range pods {
		pod := pod
		wg.Add(1)
		go func() {
			defer wg.Done()
			nodeIDInAdminAPI := getNodeIDFromAdminAPI(pod)

			require.Eventually(t, func() bool {
				_, err := nodeAPIClient.GetNode(ctx, nodeIDInAdminAPI)
				if err != nil {
					t.Logf("Failed to get node %s from Node API: %v", nodeIDInAdminAPI, err)
					return false
				}

				return true
			}, konnectNodeRegistrationTimeout, konnectNodeRegistrationCheck)

			t.Logf("proxy pod %s/%s has consistent ID %s in Node API", pod.Namespace, pod.Name, nodeIDInAdminAPI)
		}()
	}

	wg.Wait()
}
