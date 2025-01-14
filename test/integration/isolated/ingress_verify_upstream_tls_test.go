//go:build integration_tests

package isolated

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strconv"
	"strings"
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

	"github.com/kong/kubernetes-configuration/pkg/clientset"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/configuration"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers/certificate"
	"github.com/kong/kubernetes-ingress-controller/v3/test/integration/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testlabels"
)

// TestIngressVerifyUpstreamTLS tests that user-provided annotations for TLS verification are respected accordingly.
// It exposes an HTTPS goecho server listening with a certificate signed by an intermediate CA. The Service is associated
// with a root CA certificate which should be used to verify the upstream server's certificate. The test first sets the
// TLS verification depth to 0, which should fail the TLS handshake (as that would only pass with a self-signed cert).
// Then, it sets the TLS verification depth to 1, which should pass with one intermediate CA.
func TestIngressVerifyUpstreamTLS(t *testing.T) {
	const (
		caSecretName        = "ca"
		anotherCASecretName = "another-ca"
		certsVolumeName     = "certs"

		echoRoute            = "/echo"
		goEchoServerHostname = "goecho"
		goEchoTLSSecretName  = "goecho-tls"
		expectedGoEchoBody   = "Through HTTPS connection."
	)

	f := features.
		New("verify upstream TLS").
		WithLabel(testlabels.NetworkingFamily, testlabels.NetworkingFamilyIngress).
		WithLabel(testlabels.Kind, testlabels.KindIngress).
		WithSetup("deploy kong addon into cluster", featureSetup()).
		WithSetup("prepare Kong clients", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			cluster := GetClusterFromCtx(ctx)
			kongClients, err := clientset.NewForConfig(cluster.Config())
			require.NoError(t, err)
			return SetInCtxForT(ctx, t, kongClients)
		}).
		WithSetup("generate certificate", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			cleaner := GetFromCtxForT[*clusters.Cleaner](ctx, t)
			cluster := GetClusterFromCtx(ctx)
			namespace := GetNamespaceForT(ctx, t)
			ingressClass := GetIngressClassFromCtx(ctx)

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

			t.Log("Deploying CA secret")
			caSecret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: caSecretName,
					Labels: map[string]string{
						configuration.CACertLabelKey: "true",
					},
					Annotations: map[string]string{
						annotations.IngressClassKey: ingressClass,
					},
				},
				StringData: map[string]string{
					"id":   uuid.NewString(),
					"cert": string(caCertPEM),
				},
			}
			caSecret, err := cluster.Client().CoreV1().Secrets(namespace).Create(ctx, caSecret, metav1.CreateOptions{})
			assert.NoError(t, err)
			cleaner.Add(caSecret)

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

			t.Log("Deploying another CA secret to verify multiple CA certificates can be bound to a service")
			anotherCA, _ := certificate.MustGenerateCertPEMFormat(
				certificate.WithCommonName("another-ca"),
				certificate.WithCATrue(),
			)
			anotherCaSecret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: anotherCASecretName,
					Labels: map[string]string{
						configuration.CACertLabelKey: "true",
					},
					Annotations: map[string]string{
						annotations.IngressClassKey: ingressClass,
					},
				},
				StringData: map[string]string{
					"id":   uuid.NewString(),
					"cert": string(anotherCA),
				},
			}
			anotherCaSecret, err = cluster.Client().CoreV1().Secrets(namespace).Create(ctx, anotherCaSecret, metav1.CreateOptions{})
			assert.NoError(t, err)
			cleaner.Add(anotherCaSecret)

			return ctx
		}).
		WithSetup("deploy goecho service and expose it via ingress", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			cluster := GetClusterFromCtx(ctx)
			ns := GetNamespaceForT(ctx, t)
			cleaner := GetFromCtxForT[*clusters.Cleaner](ctx, t)
			ingressClass := GetIngressClassFromCtx(ctx)
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
			deployment, err := cluster.Client().AppsV1().Deployments(ns).Create(ctx, deployment, metav1.CreateOptions{})
			assert.NoError(t, err)
			cleaner.Add(deployment)

			t.Logf("Exposing deployment %s via service", deployment.Name)
			service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeClusterIP)

			t.Logf("Setting up service annotations for TLS verification")
			service.Annotations = map[string]string{
				annotations.AnnotationPrefix + annotations.ProtocolKey:              "https", // Only https or tls are supported.
				annotations.AnnotationPrefix + annotations.TLSVerifyKey:             "true",
				annotations.AnnotationPrefix + annotations.CACertificatesSecretsKey: strings.Join([]string{caSecretName, anotherCASecretName}, ","),
				annotations.AnnotationPrefix + annotations.TLSVerifyDepthKey:        "0", // First, we'll set it to 0 to make sure it fails.
			}
			service, err = cluster.Client().CoreV1().Services(ns).Create(ctx, service, metav1.CreateOptions{})
			assert.NoError(t, err)
			cleaner.Add(service)
			ctx = SetInCtxForT(ctx, t, service)

			t.Logf("Exposing service %s via ingress", service.Name)
			ingress := generators.NewIngressForService(echoRoute, map[string]string{}, service)
			ingress.Spec.IngressClassName = lo.ToPtr(ingressClass)
			ingress.Spec.Rules[0].Host = goEchoServerHostname
			_, err = cluster.Client().NetworkingV1().Ingresses(ns).Create(ctx, ingress, metav1.CreateOptions{})
			assert.NoError(t, err)
			cleaner.Add(ingress)

			t.Log("Waiting for ingress status readiness")
			assert.EventuallyWithT(t, func(t *assert.CollectT) {
				lbstatus, err := clusters.GetIngressLoadbalancerStatus(ctx, cluster, ns, ingress)
				if !assert.NoError(t, err) {
					return
				}
				assert.NotEmpty(t, lbstatus.Ingress)
			}, consts.StatusWait, consts.WaitTick)

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
			service := GetFromCtxForT[*corev1.Service](ctx, t)
			cluster := GetClusterFromCtx(ctx)

			t.Log("Fixing TLS verification depth to a sufficient value = 1")
			service.Annotations[annotations.AnnotationPrefix+annotations.TLSVerifyDepthKey] = "1"
			_, err := cluster.Client().CoreV1().Services(service.Namespace).Update(ctx, service, metav1.UpdateOptions{})
			require.NoError(t, err)

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
