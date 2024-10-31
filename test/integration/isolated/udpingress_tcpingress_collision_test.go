//go:build integration_tests

package isolated

import (
	"context"
	"strconv"
	"testing"

	"github.com/google/uuid"
	kongv1beta1 "github.com/kong/kubernetes-configuration/api/configuration/v1beta1"
	"github.com/kong/kubernetes-configuration/pkg/clientset"
	"github.com/kong/kubernetes-configuration/pkg/clientset/typed/configuration/v1beta1"
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
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testlabels"
)

func TestUDPIngressTCPIngressCollision(t *testing.T) {
	// Constants shared in many steps of this test that doesn't change.
	const (
		serviceName         = "udp-and-tcp-echo"
		commonTCPandUDPPort = 7777
		udpIngressName      = "udp-echo-ingress"
		tcpIngressName      = "tcp-echo-ingress"
	)
	testUUID := uuid.NewString()

	f := features.
		New("essentials").
		WithLabel(testlabels.NetworkingFamily, testlabels.NetworkingFamilyIngress).
		WithLabel(testlabels.Kind, testlabels.KindKongUDPIngress).
		WithLabel(testlabels.Kind, testlabels.KindKongTCPIngress).
		Setup(SkipIfRouterNotExpressions).
		WithSetup("deploy kong addon into cluster", featureSetup()).
		WithSetup("configure tcpecho and udpecho Deployment and Service on the same port and expose with ingresses", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			c, err := clientset.NewForConfig(cfg.Client().RESTConfig())
			assert.NoError(t, err)
			confV1Beta1Client := c.ConfigurationV1beta1()
			ctx = SetInCtxForT(ctx, t, confV1Beta1Client)

			cleaner := GetFromCtxForT[*clusters.Cleaner](ctx, t)
			namespace := GetNamespaceForT(ctx, t)
			client := GetClusterFromCtx(ctx).Client()

			t.Log("creating a udpecho Deployment and Service")
			container := generators.NewContainer(serviceName, test.EchoImage, test.EchoUDPPort)
			// App go-echo sends a "Running on Pod <UUID>." immediately on connecting.
			container.Env = []corev1.EnvVar{
				{
					Name:  "POD_NAME",
					Value: testUUID,
				},
				{
					Name:  "UDP_PORT",
					Value: strconv.Itoa(commonTCPandUDPPort),
				},
				{
					Name:  "TCP_PORT",
					Value: strconv.Itoa(commonTCPandUDPPort),
				},
			}
			deployment := generators.NewDeploymentForContainer(container)
			deployment, err = client.AppsV1().Deployments(namespace).Create(ctx, deployment, metav1.CreateOptions{})
			assert.NoError(t, err)
			cleaner.Add(deployment)

			t.Logf("exposing deployment %s/%s via service", deployment.Namespace, deployment.Name)
			service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeClusterIP)
			service.Name = serviceName
			service.Spec.Ports = []corev1.ServicePort{
				{
					Name:       "udp",
					Protocol:   corev1.ProtocolUDP,
					Port:       commonTCPandUDPPort,
					TargetPort: intstr.FromInt(commonTCPandUDPPort),
				},
				{
					Name:       "tcp",
					Protocol:   corev1.ProtocolTCP,
					Port:       commonTCPandUDPPort,
					TargetPort: intstr.FromInt(commonTCPandUDPPort),
				},
			}
			service, err = client.CoreV1().Services(namespace).Create(ctx, service, metav1.CreateOptions{})
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
							ServicePort: commonTCPandUDPPort,
						},
					},
				}},
			}
			_, err = confV1Beta1Client.UDPIngresses(namespace).Create(ctx, udpIngress, metav1.CreateOptions{})
			assert.NoError(t, err)

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
							ServicePort: commonTCPandUDPPort,
						},
					},
				}},
			}
			_, err = confV1Beta1Client.TCPIngresses(namespace).Create(ctx, tcpIngress, metav1.CreateOptions{})
			assert.NoError(t, err)

			return ctx
		}).
		Assess("basic test - status and connectivity", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			t.Log("verifying that the udpecho is responding properly")
			assertEventuallyResponseUDP(t, GetUDPURLFromCtx(ctx), testUUID)
			t.Log("verifying that the tcpecho is responding properly")
			assertEventuallyResponseTCP(t, GetTCPURLFromCtx(ctx), testUUID)

			return ctx
		}).
		Assess("test teardown - UDPIngress deletion", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			confV1Beta1Client := GetFromCtxForT[v1beta1.ConfigurationV1beta1Interface](ctx, t)
			namespace := GetNamespaceForT(ctx, t)

			t.Log("deleting UDPIngress")
			assert.NoError(t, confV1Beta1Client.UDPIngresses(namespace).Delete(ctx, udpIngressName, metav1.DeleteOptions{}))
			t.Log("deleting TCPIngress")
			assert.NoError(t, confV1Beta1Client.TCPIngresses(namespace).Delete(ctx, tcpIngressName, metav1.DeleteOptions{}))
			t.Log("verifying that traffic is no longer routed to udpecho")
			assertEventuallyNoResponseUDP(t, GetUDPURLFromCtx(ctx))
			t.Log("verifying that traffic is no longer routed to tcpecho")
			assertEventuallyNoResponseTCP(t, GetTCPURLFromCtx(ctx))

			return ctx
		}).
		Teardown(featureTeardown())

	tenv.Test(t, f.Feature())
}
