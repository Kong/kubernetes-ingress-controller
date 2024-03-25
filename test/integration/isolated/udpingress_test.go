//go:build integration_tests

package isolated

import (
	"context"
	"net"
	"os"
	"testing"

	"github.com/google/uuid"
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
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/v3/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/integration/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testlabels"
)

func TestUDPIngressEssentials(t *testing.T) {
	// Constants shared in many steps of this test that doesn't change.
	const udpIngressName = "upd-ingress-essentials"
	const servicePort = ktfkong.DefaultUDPServicePort
	const serviceName = "udpecho"
	testUUID := uuid.NewString()

	// Helpers used in this test.
	requireNoResponse := func(t *testing.T, udpGatewayURL string) {
		t.Helper()
		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			// For UDP lack of response (a timeout) means that we can't reach a service.
			err := test.EchoResponds(test.ProtocolUDP, udpGatewayURL, "irrelevant")
			assert.True(c, os.IsTimeout(err), "unexpected error: %v", err)
		}, consts.IngressWait, consts.WaitTick)
	}
	requireResponse := func(t *testing.T, udpGatewayURL, expectedMsg string) {
		t.Helper()
		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			assert.NoError(c, test.EchoResponds(test.ProtocolUDP, udpGatewayURL, expectedMsg))
		}, consts.IngressWait, consts.WaitTick)
	}

	f := features.
		New("essentials").
		WithLabel(testlabels.NetworkingFamily, testlabels.NetworkingFamilyGatewayAPI).
		WithLabel(testlabels.Kind, testlabels.KindKongUDPIngress).
		Setup(SkipIfRouterNotExpressions).
		WithSetup("deploy kong addon into cluster", featureSetup()).
		WithSetup("configure a udpecho Deployment and Service", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			kongClient, err := clientset.NewForConfig(cfg.Client().RESTConfig())
			assert.NoError(t, err)
			ctx = SetInCtxForT(ctx, t, kongClient)

			cleaner := GetFromCtxForT[*clusters.Cleaner](ctx, t)
			namespace := GetNamespaceForT(ctx, t)
			cluster := GetClusterFromCtx(ctx)

			t.Log("creating a udpecho Deployment and Service")
			container := generators.NewContainer(serviceName, test.EchoImage, test.EchoUDPPort)
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
			// Use the same port as the default UDP port from the Kong Gateway deployment
			// to the udpecho port, as this is what will be used to route the traffic at the Gateway.
			service.Spec.Ports = []corev1.ServicePort{{
				Name:       "udp",
				Protocol:   corev1.ProtocolUDP,
				Port:       servicePort,
				TargetPort: intstr.FromInt(test.EchoUDPPort),
			}}
			service, err = cluster.Client().CoreV1().Services(namespace).Create(ctx, service, metav1.CreateOptions{})
			assert.NoError(t, err)
			cleaner.Add(service)

			t.Logf("creating a UDPIngress to access deployment %s via kong", deployment.Name)
			udpIngress := &kongv1beta1.UDPIngress{
				ObjectMeta: metav1.ObjectMeta{
					Name: udpIngressName,
					Annotations: map[string]string{
						annotations.IngressClassKey: GetIngressClassFromCtx(ctx),
					},
				},
				Spec: kongv1beta1.UDPIngressSpec{Rules: []kongv1beta1.UDPIngressRule{
					{
						Port: ktfkong.DefaultUDPServicePort,
						Backend: kongv1beta1.IngressBackend{
							ServiceName: service.Name,
							ServicePort: servicePort,
						},
					},
				}},
			}
			_, err = kongClient.ConfigurationV1beta1().UDPIngresses(namespace).Create(ctx, udpIngress, metav1.CreateOptions{})
			assert.NoError(t, err)

			return ctx
		}).
		Assess("basic test - status and connectivity", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			udpGatewayURL := GetUDPURLFromCtx(ctx)
			ingressClient := GetFromCtxForT[*clientset.Clientset](ctx, t).ConfigurationV1beta1().UDPIngresses(GetNamespaceForT(ctx, t))

			t.Log("verifying UDPIngress status readiness")
			assert.EventuallyWithT(t, func(c *assert.CollectT) {
				configuredIngress, err := ingressClient.Get(ctx, udpIngressName, metav1.GetOptions{})
				assert.NoError(c, err)
				assert.NotNil(c, configuredIngress)
				for _, ingress := range configuredIngress.Status.LoadBalancer.Ingress {
					ipReportedByIngress := ingress.IP
					assert.NotEmpty(c, ipReportedByIngress)
					ipOfKong, _, err := net.SplitHostPort(udpGatewayURL)
					assert.NoError(c, err)
					assert.Equal(c, ipOfKong, ipReportedByIngress, "UDPIngress is not ready to redirect traffic")
				}
			}, consts.StatusWait, consts.WaitTick)

			t.Log("verifying that the udpecho is responding properly")
			requireResponse(t, udpGatewayURL, testUUID)

			return ctx
		}).
		Assess("test teardown - UDPIngress deletion", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			ingressClient := GetFromCtxForT[*clientset.Clientset](ctx, t).ConfigurationV1beta1().UDPIngresses(GetNamespaceForT(ctx, t))

			t.Log("deleting UDPIngress")
			assert.NoError(t, ingressClient.Delete(ctx, udpIngressName, metav1.DeleteOptions{}))
			t.Log("verifying that traffic is no longer routed")
			requireNoResponse(t, GetUDPURLFromCtx(ctx))
			return ctx
		}).
		Teardown(featureTeardown())

	tenv.Test(t, f.Feature())
}
