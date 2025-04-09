//go:build e2e_tests

package e2e

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"testing"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/google/uuid"
	gokong "github.com/kong/go-kong/kong"
	environment "github.com/kong/kubernetes-testing-framework/pkg/environments"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/konnect"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/konnect/nodes"
	rg "github.com/kong/kubernetes-ingress-controller/v2/internal/konnect/runtimegroups"
	rgc "github.com/kong/kubernetes-ingress-controller/v2/internal/konnect/runtimegroupsconfig"
	"github.com/kong/kubernetes-ingress-controller/v2/test/helpers/certificate"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
)

const (
	konnectRuntimeGroupsBaseURL          = "https://us.kic.api.konghq.tech/v2"
	konnectRuntimeGroupsConfigBaseURLFmt = "https://us.api.konghq.tech/konnect-api/api/runtime_groups/%s/v1"
	konnectRuntimeGroupAdminAPIBaseURL   = "https://us.kic.api.konghq.tech"
	konnectRolesBaseURL                  = "https://global.api.konghq.tech/v2"

	konnectNodeRegistrationTimeout = 5 * time.Minute
	konnectNodeRegistrationCheck   = 30 * time.Second
)

var konnectAccessToken = os.Getenv("TEST_KONG_KONNECT_ACCESS_TOKEN")

func TestKonnectConfigPush(t *testing.T) {
	t.Parallel()
	skipTestIfControllerVersionBelow(t, gatewayDiscoveryMinimalVersion)
	skipIfMissingRequiredKonnectEnvVariables(t)

	ctx, env := setupE2ETest(t)

	rgID := createTestRuntimeGroup(ctx, t)
	cert, key := createClientCertificate(ctx, t, rgID)
	createKonnectClientSecretAndConfigMap(ctx, t, env, cert, key, rgID)

	deployments := deployAllInOneKonnectManifest(ctx, t, env)

	t.Log("running ingress tests to verify all-in-one deployed ingress controller and proxy are functional")
	deployIngressWithEchoBackends(ctx, t, env, numberOfEchoBackends)
	verifyIngressWithEchoBackends(ctx, t, env, numberOfEchoBackends)

	t.Log("ensuring ingress resources are correctly populated in Konnect Runtime Group's Admin API")
	konnectAdminAPIClient := createKonnectAdminAPIClient(t, rgID, cert, key)
	verifyIngressWithEchoBackendsInAdminAPI(ctx, t, konnectAdminAPIClient.AdminAPIClient(), numberOfEchoBackends)

	t.Log("ensuring KIC nodes and controlled kong gateway nodes are present in konnect runtime group")
	requireKonnectNodesConsistentWithK8s(ctx, t, env, deployments, rgID, cert, key)
	requireAllProxyReplicasIDsConsistentWithKonnect(ctx, t, env, deployments.ProxyNN, rgID, cert, key)
}

func TestKonnectLicenseActivation(t *testing.T) {
	t.Parallel()
	skipTestIfControllerVersionBelow(t, gatewayDiscoveryMinimalVersion)
	skipIfMissingRequiredKonnectEnvVariables(t)

	ctx, env := setupE2ETest(t)

	rgID := createTestRuntimeGroup(ctx, t)
	cert, key := createClientCertificate(ctx, t, rgID)
	createKonnectClientSecretAndConfigMap(ctx, t, env, cert, key, rgID)

	manifestFile := "../../deploy/single/all-in-one-dbless-konnect-enterprise.yaml"
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
			return false
		}
		return license.License.Expiration != ""
	}, adminAPIWait, time.Second)
	t.Log("done")
}

func TestKonnectWhenMisconfiguredBasicIngressNotAffected(t *testing.T) {
	t.Parallel()
	skipIfMissingRequiredKonnectEnvVariables(t)
	skipTestIfControllerVersionBelow(t, gatewayDiscoveryMinimalVersion)
	ctx, env := setupE2ETest(t)

	rgID := createTestRuntimeGroup(ctx, t)
	cert, key := createClientCertificate(ctx, t, rgID)

	// create a Konnect client secret and config map with a non-existing runtime group ID to simulate misconfiguration
	notExistingRgID := "not-existing-rg-id"
	createKonnectClientSecretAndConfigMap(ctx, t, env, cert, key, notExistingRgID)

	deployAllInOneKonnectManifest(ctx, t, env)

	t.Log("running ingress tests to verify misconfiguration doesn't affect basic ingress functionality")
	deployIngressWithEchoBackends(ctx, t, env, numberOfEchoBackends)
	verifyIngressWithEchoBackends(ctx, t, env, numberOfEchoBackends)
}

func skipIfMissingRequiredKonnectEnvVariables(t *testing.T) {
	if konnectAccessToken == "" {
		t.Skip("missing TEST_KONG_KONNECT_ACCESS_TOKEN")
	}
}

// deployAllInOneKonnectManifest deploys all-in-one-dbless-konnect.yaml manifest, replacing the controller image
// if specified by environment variables.
func deployAllInOneKonnectManifest(ctx context.Context, t *testing.T, env environment.Environment) Deployments {
	const manifestFile = "../../deploy/single/all-in-one-dbless-konnect.yaml"
	t.Logf("deploying %s manifest file", manifestFile)

	return ManifestDeploy{Path: manifestFile}.Run(ctx, t, env)
}

func generateTestKonnectRuntimeGroupDescription(t *testing.T) string {
	t.Helper()

	desc := fmt.Sprintf("runtime group for test %s", t.Name())
	if githubServerURL != "" && githubRepo != "" && githubRunID != "" {
		githubRunURL := fmt.Sprintf("%s/%s/actions/runs/%s",
			githubServerURL, githubRepo, githubRunID)
		desc += ", github workflow run " + githubRunURL
	}

	return desc
}

// createTestRuntimeGroup creates a runtime group to be used in tests. It returns the created runtime group's ID.
// It also sets up a cleanup function for it to be deleted.
func createTestRuntimeGroup(ctx context.Context, t *testing.T) string {
	t.Helper()
	rgClient, err := rg.NewClientWithResponses(konnectRuntimeGroupsBaseURL, rg.WithRequestEditorFn(
		func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+konnectAccessToken)
			return nil
		}),
	)
	require.NoError(t, err)

	var rgID uuid.UUID
	createRgErr := retry.Do(func() error {
		rgName := uuid.NewString()
		createRgResp, err := rgClient.CreateRuntimeGroupWithResponse(ctx, rg.CreateRuntimeGroupRequest{
			Description: lo.ToPtr(generateTestKonnectRuntimeGroupDescription(t)),
			Labels: &rg.Labels{
				"created_in_tests": "true",
			},
			Name:        rgName,
			ClusterType: rg.ClusterTypeKubernetesIngressController,
		})
		if err != nil {
			return fmt.Errorf("failed to create runtime group: %w", err)
		}
		if createRgResp.StatusCode() != http.StatusCreated {
			return fmt.Errorf("failed to create RG: code %d, message %s", createRgResp.StatusCode(), string(createRgResp.Body))
		}
		if createRgResp.JSON201 == nil || createRgResp.JSON201.Id == nil {
			return errors.New("No runtime group ID in response")
		}

		rgID = *createRgResp.JSON201.Id
		return nil
	}, retry.Attempts(5), retry.Delay(time.Second))
	require.NoError(t, createRgErr)

	t.Cleanup(func() {
		t.Logf("deleting test Konnect Runtime Group: %q", rgID)
		err := retry.Do(
			func() error {
				_, err := rgClient.DeleteRuntimeGroupWithResponse(ctx, rgID)
				return err
			},
			retry.Attempts(5), retry.Delay(time.Second),
		)
		assert.NoErrorf(t, err, "failed to cleanup a runtime group: %q", rgID)
	})

	t.Logf("created test Konnect Runtime Group: %q", rgID.String())
	return rgID.String()
}

// createClientCertificate creates a TLS client certificate and POSTs it to Konnect Runtime Group configuration API
// so that KIC can use the certificates to authenticate against Konnect Admin API.
func createClientCertificate(ctx context.Context, t *testing.T, rgID string) (certPEM string, keyPEM string) {
	t.Helper()

	rgConfigClient, err := rgc.NewClientWithResponses(fmt.Sprintf(konnectRuntimeGroupsConfigBaseURLFmt, rgID), rgc.WithRequestEditorFn(
		func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+konnectAccessToken)
			return nil
		}),
		rgc.WithHTTPClient(helpers.RetryableHTTPClient(helpers.DefaultHTTPClient())),
	)
	require.NoError(t, err)

	cert, key := certificate.MustGenerateSelfSignedCertPEMFormat()

	t.Log("creating client certificate in Konnect")
	resp, err := rgConfigClient.PostDpClientCertificatesWithResponse(ctx, rgc.PostDpClientCertificatesJSONRequestBody{
		Cert: string(cert),
	})
	require.NoError(t, err)
	require.Equalf(t, http.StatusCreated, resp.StatusCode(), "failed creating client certificate: %s", string(resp.Body))

	return string(cert), string(key)
}

// createKonnectClientSecretAndConfigMap creates a Secret with client TLS certificate that is used by KIC to communicate
// with Konnect Admin API. It also creates a ConfigMap that specifies a Runtime Group ID and Konnect Admin API URL.
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
			"CONTROLLER_KONNECT_RUNTIME_GROUP_ID": rgID,
			"CONTROLLER_KONNECT_ADDRESS":          konnectRuntimeGroupAdminAPIBaseURL,
		},
	}, metav1.CreateOptions{})
	require.NoError(t, err)
}

// createKonnectAdminAPIClient creates an *kong.Client that will communicate with Konnect Runtime Group's Admin API.
func createKonnectAdminAPIClient(t *testing.T, rgID, cert, key string) *adminapi.KonnectClient {
	t.Helper()

	c, err := adminapi.NewKongClientForKonnectRuntimeGroup(adminapi.KonnectConfig{
		RuntimeGroupID: rgID,
		Address:        konnectRuntimeGroupAdminAPIBaseURL,
		TLSClient: adminapi.TLSClientConfig{
			Cert: cert,
			Key:  key,
		},
	})
	require.NoError(t, err)
	return c
}

// createKonnectNodeClient creates a konnect.NodeClient to get nodes in konnect runtime group.
func createKonnectNodeClient(t *testing.T, rgID, cert, key string) *nodes.Client {
	cfg := adminapi.KonnectConfig{
		ConfigSynchronizationEnabled: true,
		RuntimeGroupID:               rgID,
		Address:                      konnectRuntimeGroupAdminAPIBaseURL,
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

		kongClient, err := gokong.NewClient(lo.ToPtr(address), client)
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
					t.Logf("failed to get node %s from Node API: %v", nodeIDInAdminAPI, err)
					return false
				}

				return true
			}, konnectNodeRegistrationTimeout, konnectNodeRegistrationCheck)

			t.Logf("proxy pod %s/%s has consistent ID %s in Node API", pod.Namespace, pod.Name, nodeIDInAdminAPI)
		}()
	}

	wg.Wait()
}
