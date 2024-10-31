//go:build integration_tests

package isolated

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/google/uuid"
	kongv1beta1 "github.com/kong/kubernetes-configuration/api/configuration/v1beta1"
	"github.com/kong/kubernetes-configuration/pkg/clientset"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	ktfkong "github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/integration/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testlabels"
)

func TestTCPIngressEssentials(t *testing.T) {
	const tcpIngressName = "tcp-ingress-essentials"
	const servicePort = ktfkong.DefaultTCPServicePort
	const serviceName = "tcpecho"
	testUUID := uuid.NewString()

	f := features.
		New("essentials").
		WithLabel(testlabels.NetworkingFamily, testlabels.NetworkingFamilyIngress).
		WithLabel(testlabels.Kind, testlabels.KindKongTCPIngress).
		Setup(SkipIfRouterNotExpressions).
		WithSetup("deploy kong addon into cluster", featureSetup()).
		WithSetup("configure a tcpecho Deployment and Service",
			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				t.Log("setting up the TCPIngress tests")

				kongClient, err := clientset.NewForConfig(cfg.Client().RESTConfig())
				assert.NoError(t, err)
				ctx = SetInCtxForT(ctx, t, kongClient)

				cleaner := GetFromCtxForT[*clusters.Cleaner](ctx, t)
				namespace := GetNamespaceForT(ctx, t)
				cluster := GetClusterFromCtx(ctx)

				t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
				container := generators.NewContainer(serviceName, test.HTTPBinImage, test.HTTPBinPort)
				// App go-echo sends a "Running on Pod <UUID>." immediately on connecting.
				container.Env = []corev1.EnvVar{
					{
						Name:  "POD_NAME",
						Value: testUUID,
					},
				}
				deployment := generators.NewDeploymentForContainer(container)
				deployment, err = cluster.Client().AppsV1().Deployments(namespace).Create(ctx, deployment, metav1.CreateOptions{})
				assert.NoError(t, err)
				cleaner.Add(deployment)

				t.Logf("exposing deployment %s/%s via service", deployment.Namespace, deployment.Name)
				service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeClusterIP)
				service.Name = serviceName
				// Use the same port as the default TCP port from the Kong Gateway deployment
				// to the tcpecho port, as this is what will be used to route the traffic at the Gateway.
				service.Spec.Ports = []corev1.ServicePort{{
					Name:       "tcp",
					Protocol:   corev1.ProtocolTCP,
					Port:       servicePort,
					TargetPort: intstr.FromInt(test.HTTPBinPort),
				}}
				service, err = cluster.Client().CoreV1().Services(namespace).Create(ctx, service, metav1.CreateOptions{})
				assert.NoError(t, err)
				cleaner.Add(service)

				t.Logf("creating a TCPIngress to access deployment %s via kong", deployment.Name)
				tcpIngress := &kongv1beta1.TCPIngress{
					ObjectMeta: metav1.ObjectMeta{
						Name: tcpIngressName,
						Annotations: map[string]string{
							annotations.IngressClassKey: GetIngressClassFromCtx(ctx),
						},
					},
					Spec: kongv1beta1.TCPIngressSpec{Rules: []kongv1beta1.IngressRule{
						{
							Port: ktfkong.DefaultTCPServicePort,
							Backend: kongv1beta1.IngressBackend{
								ServiceName: service.Name,
								ServicePort: servicePort,
							},
						},
					}},
				}
				_, err = kongClient.ConfigurationV1beta1().TCPIngresses(namespace).Create(ctx, tcpIngress, metav1.CreateOptions{})
				assert.NoError(t, err)

				return ctx
			}).
		Assess("basic test - status and connectivity", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			tcpGatewayURL := GetTCPURLFromCtx(ctx)
			ingressClient := GetFromCtxForT[*clientset.Clientset](ctx, t).ConfigurationV1beta1().TCPIngresses(GetNamespaceForT(ctx, t))

			t.Logf("verifying that TCPIngress becomes routable at %s", tcpGatewayURL)
			assert.EventuallyWithT(t, func(c *assert.CollectT) {
				configuredIngress, err := ingressClient.Get(ctx, tcpIngressName, metav1.GetOptions{})
				assert.NoError(c, err)
				assert.NotNil(c, configuredIngress)
				for _, ingress := range configuredIngress.Status.LoadBalancer.Ingress {
					ipReportedByIngress := ingress.IP
					assert.NotEmpty(c, ipReportedByIngress)
					ipOfKong, _, err := net.SplitHostPort(tcpGatewayURL)
					assert.NoError(c, err)
					assert.Equal(c, ipOfKong, ipReportedByIngress, "TCPIngress is not ready to redirect traffic")
				}
			}, consts.StatusWait, consts.WaitTick)

			tcpProxyURL := fmt.Sprintf("http://%s", tcpGatewayURL)
			t.Logf("verifying that the TCPIngress %s is responding ready", tcpProxyURL)
			assert.EventuallyWithT(t, func(c *assert.CollectT) {
				resp, err := helpers.DefaultHTTPClient().Get(tcpProxyURL)
				if assert.NoError(c, err) && assert.NotNil(c, resp) {
					// the assert.EventuallyWithT will collect all errors,
					// so the Close() can only be called if the response is not nil
					defer resp.Body.Close()
					assert.Equal(c, http.StatusOK, resp.StatusCode)
					// now that the ingress backend is routable, make sure the contents we're getting back are what we expect
					// Expected: "<title>httpbin.org</title>"
					b := new(bytes.Buffer)
					n, err := b.ReadFrom(resp.Body)
					assert.NoError(c, err)
					assert.Greater(c, n, int64(0))
					assert.Contains(c, b.String(), "<title>httpbin.org</title>")
				}
			}, consts.StatusWait, consts.WaitTick)

			return ctx
		}).
		Assess("test teardown - TCPIngress deletion", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			ingressClient := GetFromCtxForT[*clientset.Clientset](ctx, t).ConfigurationV1beta1().TCPIngresses(GetNamespaceForT(ctx, t))

			t.Log("deleting TCPIngress")
			assert.NoError(t, ingressClient.Delete(ctx, tcpIngressName, metav1.DeleteOptions{}))
			t.Logf("verifying that traffic is no longer routed")
			assertEventuallyNoResponseTCP(t, GetTCPURLFromCtx(ctx))

			return ctx
		}).
		Teardown(featureTeardown())

	tenv.Test(t, f.Feature())
}
