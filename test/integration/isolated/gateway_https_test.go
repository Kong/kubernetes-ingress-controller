//go:build integration_tests

package isolated

import (
	"context"
	"crypto/x509"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers/certificate"
	"github.com/kong/kubernetes-ingress-controller/v3/test/integration/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testlabels"
)

func TestGatewayHTTPSMultipleCertificates(t *testing.T) {
	const (
		wildcardExample     = "*.example.com"
		nameWildcardExample = "example"
		testExampleURL      = "https://test.example.com"

		wildCardOneInternalExample     = "*.one.internal.example.com"
		nameWildcardOneInternalExample = "one-internal-example"
		testOneInternalExampleURL      = "https://test.one.internal.example.com"
	)

	f := features.
		New("essentials").
		WithLabel(testlabels.NetworkingFamily, testlabels.NetworkingFamilyGatewayAPI).
		WithLabel(testlabels.Kind, testlabels.KindHTTPRoute).
		Setup(SkipIfRouterNotExpressions).
		WithSetup("deploy kong addon into cluster", featureSetup(
			withControllerManagerOpts(helpers.ControllerManagerOptAdditionalWatchNamespace("default")),
			withKongProxyEnvVars(map[string]string{
				"PROXY_LISTEN": `0.0.0.0:8443 http2 ssl`, // Ensure that only HTTPS is available.
			}),
		)).
		Assess(
			"deploying gateway with configuration (certificates and routing) to cluster",
			func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
				cleaner := GetFromCtxForT[*clusters.Cleaner](ctx, t)
				cluster := GetClusterFromCtx(ctx)
				namespace := GetNamespaceForT(ctx, t)

				certPool := x509.NewCertPool()
				ctx = SetCertPoolInCtx(ctx, certPool)
				t.Log("deploying secrets with certificates")
				createK8sCertSecret := func(
					ctx context.Context, t *testing.T, cluster clusters.Cluster, nn k8stypes.NamespacedName, certPool *x509.CertPool, domainName string,
				) {
					t.Helper()
					certificateCrt, firstCertificateKey := certificate.MustGenerateCertPEMFormat(
						certificate.WithCommonName(domainName),
						certificate.WithDNSNames(domainName),
					)
					certificateSecret := &corev1.Secret{
						ObjectMeta: metav1.ObjectMeta{
							Name: nn.Name,
						},
						Type: corev1.SecretTypeTLS,
						Data: map[string][]byte{
							"tls.crt": certificateCrt,
							"tls.key": firstCertificateKey,
						},
					}
					_, err := cluster.Client().CoreV1().Secrets(nn.Namespace).Create(ctx, certificateSecret, metav1.CreateOptions{})
					require.NoError(t, err)
					require.True(t, certPool.AppendCertsFromPEM(certificateCrt))
				}
				createK8sCertSecret(
					ctx, t, cluster, k8stypes.NamespacedName{Namespace: namespace, Name: nameWildcardExample}, certPool, wildcardExample,
				)
				createK8sCertSecret(
					ctx, t, cluster, k8stypes.NamespacedName{Namespace: namespace, Name: nameWildcardOneInternalExample}, certPool, wildCardOneInternalExample,
				)

				t.Log("getting a gateway client")
				gatewayClient, err := gatewayclient.NewForConfig(cluster.Config())
				require.NoError(t, err)

				t.Log("deploying a new gatewayClass")
				gatewayClassName := uuid.NewString()
				gwc, err := helpers.DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
				require.NoError(t, err)
				cleaner.Add(gwc)

				t.Log("deploying a new gateway")
				createHTTPSListener := func(t *testing.T, domainName string, certSecretName string) gatewayapi.Listener {
					t.Helper()
					return gatewayapi.Listener{
						Name:     gatewayapi.SectionName(fmt.Sprintf("https-%s", certSecretName)),
						Protocol: gatewayapi.HTTPSProtocolType,
						Port:     gatewayapi.PortNumber(443),
						Hostname: lo.ToPtr(gatewayapi.Hostname(domainName)),
						TLS: &gatewayapi.GatewayTLSConfig{
							Mode: lo.ToPtr(gatewayapi.TLSModeTerminate),
							CertificateRefs: []gatewayapi.SecretObjectReference{
								{
									Name: gatewayapi.ObjectName(certSecretName),
									Kind: lo.ToPtr(gatewayapi.Kind("Secret")),
								},
							},
						},
					}
				}
				gatewayName := uuid.NewString()
				gateway, err := helpers.DeployGateway(ctx, gatewayClient, namespace, gatewayClassName, func(gw *gatewayapi.Gateway) {
					gw.Name = gatewayName
					gw.Spec.Listeners = []gatewayapi.Listener{
						createHTTPSListener(t, wildcardExample, nameWildcardExample),
						createHTTPSListener(t, wildCardOneInternalExample, nameWildcardOneInternalExample),
					}
				})
				require.NoError(t, err)
				cleaner.Add(gateway)

				t.Log("deploying a minimal HTTP container")
				deployment := generators.NewDeploymentForContainer(
					generators.NewContainer("echo", test.EchoImage, test.EchoHTTPPort),
				)
				deployment.Spec.Template.Spec.Containers[0].Env = append(deployment.Spec.Template.Spec.Containers[0].Env, corev1.EnvVar{
					Name:  "POD_NAME",
					Value: "echo",
				})
				deployment, err = cluster.Client().AppsV1().Deployments(namespace).Create(ctx, deployment, metav1.CreateOptions{})
				require.NoError(t, err)

				t.Logf("exposing deployment %q via service", deployment.Name)
				service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeClusterIP)
				_, err = cluster.Client().CoreV1().Services(namespace).Create(ctx, service, metav1.CreateOptions{})
				require.NoError(t, err)

				t.Logf("creating a HTTPRoute to access deployment %q via Kong", deployment.Name)
				httpRoute := &gatewayapi.HTTPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name: "echo-httproute",
					},
					Spec: gatewayapi.HTTPRouteSpec{
						CommonRouteSpec: gatewayapi.CommonRouteSpec{
							ParentRefs: []gatewayapi.ParentReference{{
								Name: gatewayapi.ObjectName(gateway.Name),
							}},
						},
						Hostnames: []gatewayapi.Hostname{
							gatewayapi.Hostname(wildcardExample),
							gatewayapi.Hostname(wildCardOneInternalExample),
						},
						Rules: []gatewayapi.HTTPRouteRule{
							{
								BackendRefs: []gatewayapi.HTTPBackendRef{
									{
										BackendRef: gatewayapi.BackendRef{
											BackendObjectReference: gatewayapi.BackendObjectReference{
												Name: gatewayapi.ObjectName(service.Name),
												Port: lo.ToPtr(gatewayapi.PortNumber(test.EchoHTTPPort)),
												Kind: lo.ToPtr(gatewayapi.Kind("Service")),
											},
										},
									},
								},
							},
						},
					},
				}
				_, err = gatewayClient.GatewayV1().HTTPRoutes(namespace).Create(ctx, httpRoute, metav1.CreateOptions{})
				require.NoError(t, err)

				return ctx
			},
		).
		Assess(
			"verifying that certs match and HTTPS traffic is routed to the service",
			func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
				verifyEventuallyGet := func(url string) {
					helpers.EventuallyGETPath(
						t,
						GetHTTPSURLFromCtx(ctx),
						url,
						"",
						GetCertPoolFromCtx(ctx),
						http.StatusOK,
						"echo",
						nil,
						consts.IngressWait,
						consts.WaitTick,
					)
				}

				t.Logf("verifying that %s is routed by %s", testExampleURL, wildcardExample)
				verifyEventuallyGet(testExampleURL)
				t.Logf("verifying that %s is routed by %s", testOneInternalExampleURL, wildCardOneInternalExample)
				verifyEventuallyGet(testOneInternalExampleURL)

				return ctx
			},
		).
		Teardown(featureTeardown())

	tenv.Test(t, f.Feature())
}
