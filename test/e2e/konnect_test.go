//go:build e2e_tests

package e2e

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	environment "github.com/kong/kubernetes-testing-framework/pkg/environments"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/konnect"
	rg "github.com/kong/kubernetes-ingress-controller/v2/internal/konnect/runtimegroups"
	rgc "github.com/kong/kubernetes-ingress-controller/v2/internal/konnect/runtimegroupsconfig"
)

const (
	konnectRuntimeGroupsBaseURL          = "https://us.kic.api.konghq.tech/v2"
	konnectRuntimeGroupsConfigBaseURLFmt = "https://us.api.konghq.tech/konnect-api/api/runtime_groups/%s/v1"
	konnectRuntimeGroupAdminAPIBaseURL   = "https://us.kic.api.konghq.tech"

	konnectNodeRegistrationTimeout = 5 * time.Minute
	konnectNodeRegistrationCheck   = 30 * time.Second
)

var konnectAccessToken = os.Getenv("TEST_KONG_KONNECT_ACCESS_TOKEN")

func TestKonnectConfigPush(t *testing.T) {
	t.Parallel()
	skipIfMissingRequiredKonnectEnvVariables(t)

	ctx, env := setupE2ETest(t)

	rgID := createTestRuntimeGroup(ctx, t)
	cert, key := createClientCertificate(ctx, t, rgID)
	createKonnectClientSecretAndConfigMap(ctx, t, env, cert, key, rgID)

	deployAllInOneKonnectManifest(ctx, t, env)

	t.Log("running ingress tests to verify all-in-one deployed ingress controller and proxy are functional")
	deployIngress(ctx, t, env)
	verifyIngress(ctx, t, env)

	t.Log("ensuring ingress resources are correctly populated in Konnect Runtime Group's Admin API")
	konnectAdminAPIClient := createKonnectAdminAPIClient(t, rgID, cert, key)
	requireIngressConfiguredInAdminAPIEventually(ctx, t, konnectAdminAPIClient.AdminAPIClient())

	t.Log("ensuring KIC nodes and controlled kong gateway nodes are present in konnect runtime group")
	requireKonnectNodesConsistentWithK8s(ctx, t, env, rgID, cert, key)
}

func TestKonnectWhenMisconfiguredBasicIngressNotAffected(t *testing.T) {
	t.Parallel()
	skipIfMissingRequiredKonnectEnvVariables(t)

	ctx, env := setupE2ETest(t)

	rgID := createTestRuntimeGroup(ctx, t)
	cert, key := createClientCertificate(ctx, t, rgID)

	// create a Konnect client secret and config map with a non-existing runtime group ID to simulate misconfiguration
	notExistingRgID := "not-existing-rg-id"
	createKonnectClientSecretAndConfigMap(ctx, t, env, cert, key, notExistingRgID)

	deployAllInOneKonnectManifest(ctx, t, env)

	t.Log("running ingress tests to verify misconfiguration doesn't affect basic ingress functionality")
	deployIngress(ctx, t, env)
	verifyIngress(ctx, t, env)
}

func skipIfMissingRequiredKonnectEnvVariables(t *testing.T) {
	if konnectAccessToken == "" {
		t.Skip("missing TEST_KONG_KONNECT_ACCESS_TOKEN")
	}
}

// deployAllInOneKonnectManifest deploys all-in-one-dbless-konnect.yaml manifest, replacing the controller image
// if specified by environment variables.
func deployAllInOneKonnectManifest(ctx context.Context, t *testing.T, env environment.Environment) {
	const manifestFile = "../../deploy/single/all-in-one-dbless-konnect.yaml"
	t.Logf("deploying %s manifest file", manifestFile)

	manifest := getTestManifest(t, manifestFile)
	deployKong(ctx, t, env, manifest)
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

	createRgResp, err := rgClient.CreateRuntimeGroupWithResponse(ctx, rg.CreateRuntimeGroupRequest{
		Description: lo.ToPtr("This is a description"),
		Labels:      &rg.Labels{"created_in_tests": "true"},
		Name:        uuid.NewString(),
	})
	require.NoError(t, err, "failed to create runtime group")
	require.Equal(t, http.StatusCreated, createRgResp.StatusCode())
	require.NotNil(t, createRgResp.JSON201)
	require.NotNil(t, createRgResp.JSON201.Id)
	id := *createRgResp.JSON201.Id
	t.Cleanup(func() {
		_, err := rgClient.DeleteRuntimeGroupWithResponse(ctx, id)
		assert.NoErrorf(t, err, "failed to cleanup a runtime group: %q", id)
	})

	t.Logf("created test Konnect Runtime Group: %q", id.String())
	return id.String()
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
	)
	require.NoError(t, err)

	priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	require.NoError(t, err)
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Kong Inc."},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 180),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	require.NoError(t, err)

	out := &bytes.Buffer{}
	err = pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	require.NoError(t, err)
	cert := out.String()

	out.Reset()
	err = pem.Encode(out, pemBlockForKey(t, priv))
	require.NoError(t, err)
	key := out.String()

	t.Log("creating client certificate in Konnect")
	resp, err := rgConfigClient.PostDpClientCertificatesWithResponse(ctx, rgc.PostDpClientCertificatesJSONRequestBody{
		Cert: cert,
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode())

	return cert, key
}

func pemBlockForKey(t *testing.T, k *ecdsa.PrivateKey) *pem.Block {
	b, err := x509.MarshalECPrivateKey(k)
	require.NoError(t, err)
	return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
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
func createKonnectAdminAPIClient(t *testing.T, rgID, cert, key string) *adminapi.Client {
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

// createKonnectNodeClient creates a konnect.NodeAPIClient to get nodes in konnect runtime group.
func createKonnectNodeClient(t *testing.T, rgID, cert, key string) *konnect.NodeAPIClient {
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
	c, err := konnect.NewNodeAPIClient(cfg)
	require.NoError(t, err)
	return c
}

func requireKonnectNodesConsistentWithK8s(ctx context.Context, t *testing.T, env environment.Environment, rgID string, cert, key string) {
	konnectNodeClient := createKonnectNodeClient(t, rgID, cert, key)
	require.Eventually(t, func() bool {
		nodes, err := konnectNodeClient.ListAllNodes(ctx)
		if err != nil {
			t.Logf("list all nodes failed: %v", err)
			return false
		}

		kicPods, err := listPodsByLabels(ctx, env, "kong", map[string]string{"app": "ingress-kong"})
		if err != nil || len(kicPods) != 1 {
			return false
		}

		kongPods, err := listPodsByLabels(ctx, env, "kong", map[string]string{"app": "proxy-kong"})
		if err != nil || len(kongPods) != 2 {
			return false
		}

		kicNodes := []*konnect.NodeItem{}
		kongNodes := []*konnect.NodeItem{}

		for _, node := range nodes {
			if node.Type == konnect.NodeTypeIngressController {
				kicNodes = append(kicNodes, node)
			}
			if node.Type == konnect.NodeTypeKongProxy {
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
			if !lo.ContainsBy(kongNodes, func(n *konnect.NodeItem) bool {
				return n.Hostname == nsName
			}) {
				return false
			}
		}

		return true
	}, konnectNodeRegistrationTimeout, konnectNodeRegistrationCheck)
}
