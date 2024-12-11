//go:build integration_tests

package isolated

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/configuration"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	constsmgr "github.com/kong/kubernetes-ingress-controller/v3/internal/manager/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers/certificate"
	"github.com/kong/kubernetes-ingress-controller/v3/test/integration/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testlabels"
)

// TestBackendTLSPolicy tests that BackendTLSPolicies are properly configured.
// It exposes an HTTPS goecho server listening with a certificate signed by an intermediate CA. The Service is associated
// with a root CA certificate which should be used to verify the upstream server's certificate. The test first sets the
// TLS verification depth to 0, which should fail the TLS handshake (as that would only pass with a self-signed cert).
// Then, it sets the TLS verification depth to 1, which should pass with one intermediate CA.
func TestBackendTLSPolicy(t *testing.T) {
	const (
		caConfigMapName        = "ca"
		anotherCAConfigMapName = "another-ca"
		certsVolumeName        = "certs"

		echoRoute            = "/echo"
		goEchoServerHostname = "goecho"
		goEchoTLSSecretName  = "goecho-tls"
		expectedGoEchoBody   = "Through HTTPS connection."
	)

	f := features.
		New("verify upstream TLS").
		WithLabel(testlabels.NetworkingFamily, testlabels.NetworkingFamilyIngress).
		WithLabel(testlabels.Kind, testlabels.KindBackendTLSPolicy).
		WithSetup("deploy kong addon into cluster", featureSetup()).
		WithSetup("prepare gateway API client", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			cluster := GetClusterFromCtx(ctx)
			gwapiClient, err := gatewayclient.NewForConfig(cluster.Config())
			require.NoError(t, err)
			return SetInCtxForT(ctx, t, gwapiClient)
		}).
		WithSetup("generate certificate", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			cleaner := GetFromCtxForT[*clusters.Cleaner](ctx, t)
			cluster := GetClusterFromCtx(ctx)
			namespace := GetNamespaceForT(ctx, t)

			caCert := certificate.MustGenerateCert(
				certificate.WithCommonName("ca"),
				certificate.WithCATrue(),
				certificate.WithMaxPathLen(1),
			)
			caCertPEM, _ := certificate.CertToPEMFormat(caCert)

			intermediateCert := certificate.MustGenerateCert(
				certificate.WithCommonName("intermediate"),
				certificate.WithCATrue(),
				certificate.WithMaxPathLen(0), // Can only sign end-entity certificates.
				certificate.WithParent(caCert),
			)
			intermediateCertPEM, _ := certificate.CertToPEMFormat(intermediateCert)

			cert, key := certificate.MustGenerateCertPEMFormat(
				certificate.WithDNSNames(goEchoServerHostname),
				certificate.WithParent(intermediateCert),
			)

			// Generate a chain with server and intermediate certificates to be used by the server.
			// It has to present the server certificate first, followed by the intermediate certificate.
			bundle := bytes.Join([][]byte{cert, intermediateCertPEM}, nil)

			t.Log("Deploying CA configmap")
			caConfigmap := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name: caConfigMapName,
					Labels: map[string]string{
						configuration.CACertLabelKey:       "true",
						constsmgr.DefaultConfigMapSelector: "true",
					},
				},
				Data: map[string]string{
					"id":     uuid.NewString(),
					"ca.crt": string(caCertPEM),
				},
			}
			caConfigmap, err := cluster.Client().CoreV1().ConfigMaps(namespace).Create(ctx, caConfigmap, metav1.CreateOptions{})
			assert.NoError(t, err)
			cleaner.Add(caConfigmap)

			t.Log("Deploying goecho certificate")
			goechoSecret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: goEchoTLSSecretName,
				},
				StringData: map[string]string{
					"tls.crt": string(bundle),
					"tls.key": string(key),
				},
			}
			_, err = cluster.Client().CoreV1().Secrets(namespace).Create(ctx, goechoSecret, metav1.CreateOptions{})
			assert.NoError(t, err)

			t.Log("Deploying another CA configmap to verify multiple CA certificates can be bound to a service")
			anotherCA, _ := certificate.MustGenerateCertPEMFormat(
				certificate.WithCommonName("another-ca"),
				certificate.WithCATrue(),
			)
			anotherCaConfigMap := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name: anotherCAConfigMapName,
					Labels: map[string]string{
						configuration.CACertLabelKey:       "true",
						constsmgr.DefaultConfigMapSelector: "true",
					},
				},
				Data: map[string]string{
					"id":     uuid.NewString(),
					"ca.crt": string(anotherCA),
				},
			}
			anotherCaConfigMap, err = cluster.Client().CoreV1().ConfigMaps(namespace).Create(ctx, anotherCaConfigMap, metav1.CreateOptions{})
			assert.NoError(t, err)
			cleaner.Add(anotherCaConfigMap)

			return ctx
		}).
		WithSetup("deploy goecho service and expose it via HTTPRoute", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			cluster := GetClusterFromCtx(ctx)
			gwapiClient := GetFromCtxForT[*gatewayclient.Clientset](ctx, t)

			ns := GetNamespaceForT(ctx, t)
			cleaner := GetFromCtxForT[*clusters.Cleaner](ctx, t)

			t.Log("Create GatewayClass")
			gatewayClassName := uuid.NewString()
			gatewayCLass, err := helpers.DeployGatewayClass(ctx, gwapiClient, gatewayClassName)
			assert.NoError(t, err)
			cleaner.Add(gatewayCLass)

			t.Log("Create Gateway")
			gatewayName := uuid.NewString()
			gateway, err := helpers.DeployGateway(ctx, gwapiClient, ns, gatewayClassName, func(gw *gatewayapi.Gateway) {
				gw.Name = gatewayName
			})
			assert.NoError(t, err)
			cleaner.Add(gateway)

			container := generators.NewContainer("goecho", test.EchoImage, test.EchoHTTPSPort)
			container.VolumeMounts = []corev1.VolumeMount{
				{
					Name:      certsVolumeName,
					MountPath: "/etc/certs",
				},
			}
			container.Env = []corev1.EnvVar{
				{
					Name:  "HTTPS_PORT",
					Value: strconv.Itoa(test.EchoHTTPSPort),
				},
				{
					Name:  "TLS_CERT_FILE",
					Value: "/etc/certs/tls.crt",
				},
				{
					Name:  "TLS_KEY_FILE",
					Value: "/etc/certs/tls.key",
				},
			}
			deployment := generators.NewDeploymentForContainer(container)
			deployment.Spec.Template.Spec.Volumes = []corev1.Volume{
				{
					Name: certsVolumeName,
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: goEchoTLSSecretName,
						},
					},
				},
			}
			deployment, err = cluster.Client().AppsV1().Deployments(ns).Create(ctx, deployment, metav1.CreateOptions{})
			assert.NoError(t, err)
			cleaner.Add(deployment)

			t.Logf("Exposing deployment %s via service", deployment.Name)
			service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeClusterIP)

			service, err = cluster.Client().CoreV1().Services(ns).Create(ctx, service, metav1.CreateOptions{})
			assert.NoError(t, err)
			cleaner.Add(service)
			ctx = SetInCtxForT(ctx, t, service)

			t.Logf("Expose service %s via HTTPRoute", service.Name)
			httpRoute := &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{
							{
								Name: gatewayapi.ObjectName(gateway.Name),
							},
						},
					},
					Rules: []gatewayapi.HTTPRouteRule{
						{
							BackendRefs: []gatewayapi.HTTPBackendRef{
								{
									BackendRef: gatewayapi.BackendRef{
										BackendObjectReference: gatewayapi.BackendObjectReference{
											Name: gatewayapi.ObjectName(service.Name),
											Port: lo.ToPtr(gatewayapi.PortNumber(1028)),
										},
									},
								},
							},
						},
					},
				},
			}
			_, err = gwapiClient.GatewayV1().HTTPRoutes(ns).Create(ctx, httpRoute, metav1.CreateOptions{})
			assert.NoError(t, err)

			t.Logf("Targeting service %s with backendTLSPolicy", service.Name)
			backendTLSPolicy := &gatewayapi.BackendTLSPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name: "goecho-tls-policy",
				},
				Spec: gatewayapi.BackendTLSPolicySpec{
					TargetRefs: []gatewayapi.LocalPolicyTargetReferenceWithSectionName{
						{
							LocalPolicyTargetReference: gatewayapi.LocalPolicyTargetReference{
								Group: "core",
								Kind:  "Service",
								Name:  gatewayapi.ObjectName(service.Name),
							},
						},
					},
					Validation: gatewayapi.BackendTLSPolicyValidation{
						CACertificateRefs: []gatewayapi.LocalObjectReference{
							{
								Group: "core",
								Kind:  "ConfigMap",
								Name:  caConfigMapName,
							},
							{
								Group: "core",
								Kind:  "ConfigMap",
								Name:  anotherCAConfigMapName,
							},
						},
						Hostname: goEchoServerHostname,
					},
					Options: map[gatewayapi.AnnotationKey]gatewayapi.AnnotationValue{
						gatewayapi.TLSVerifyDepthKey: gatewayapi.AnnotationValue("0"),
					},
				},
			}
			backendTLSPolicy, err = gwapiClient.GatewayV1alpha3().BackendTLSPolicies(ns).Create(ctx, backendTLSPolicy, metav1.CreateOptions{})
			assert.NoError(t, err)
			ctx = SetInCtxForT(ctx, t, backendTLSPolicy)

			return ctx
		}).
		Assess("verify that if the verify-depth is lower than chain length, TLS handshake fails", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			proxyURL := GetHTTPURLFromCtx(ctx)
			require.EventuallyWithT(t, func(t *assert.CollectT) {
				req, err := http.NewRequest("GET", proxyURL.String()+echoRoute, nil)
				if !assert.NoError(t, err) {
					return
				}
				req.Host = goEchoServerHostname

				resp, err := http.DefaultClient.Do(req)
				if !assert.NoError(t, err) {
					return
				}
				defer resp.Body.Close()
				assert.Equal(t, http.StatusBadGateway, resp.StatusCode)
			}, consts.StatusWait, consts.WaitTick)
			return ctx
		}).
		Assess("verify that Kong can access the upstream service when using correct host", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			backendTLSPolicy := GetFromCtxForT[*gatewayapi.BackendTLSPolicy](ctx, t)
			gwapiClient := GetFromCtxForT[*gatewayclient.Clientset](ctx, t)

			t.Log("Fixing TLS verification depth to a sufficient value = 1")
			require.Eventually(t, func() bool {
				backendTLSPolicy, err := gwapiClient.GatewayV1alpha3().BackendTLSPolicies(GetNamespaceForT(ctx, t)).Get(ctx, backendTLSPolicy.Name, metav1.GetOptions{})
				if err != nil {
					t.Logf("Failed to get BackendTLSPolicy: %v", err)
					return false
				}
				backendTLSPolicy.Spec.Options[gatewayapi.TLSVerifyDepthKey] = "1"
				_, err = gwapiClient.GatewayV1alpha3().BackendTLSPolicies(GetNamespaceForT(ctx, t)).Update(ctx, backendTLSPolicy, metav1.UpdateOptions{})
				if err != nil {
					t.Logf("Failed to update BackendTLSPolicy: %v", err)
					return false
				}
				return true
			}, consts.StatusWait, consts.WaitTick)

			proxyURL := GetHTTPURLFromCtx(ctx)
			require.EventuallyWithT(t, func(t *assert.CollectT) {
				req, err := http.NewRequest("GET", proxyURL.String()+echoRoute, nil)
				if !assert.NoError(t, err) {
					return
				}
				req.Host = goEchoServerHostname

				resp, err := http.DefaultClient.Do(req)
				if !assert.NoError(t, err) {
					return
				}
				defer resp.Body.Close()
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				b, _ := io.ReadAll(resp.Body)
				assert.Contains(t, string(b), expectedGoEchoBody)
			}, consts.StatusWait, consts.WaitTick)
			return ctx
		}).
		Teardown(featureTeardown())

	tenv.Test(t, f.Feature())
}
