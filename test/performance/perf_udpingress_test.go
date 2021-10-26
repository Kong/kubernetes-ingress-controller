//go:build performance_tests
// +build performance_tests

package performance

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset"
)

func TestUDPIngressPerformance(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), ingressWait)
	cluster := env.Cluster()
	defer cancel()

	cnt := 1
	cost := 0
	for cnt <= max_ingress {
		testName := "minudp"
		testUDPIngressNamespace := fmt.Sprintf("udpingress-%d", cnt)
		t.Logf("creating namespace %s for testing", testUDPIngressNamespace)
		err := CreateNamespace(ctx, testUDPIngressNamespace, t)
		assert.NoError(t, err)

		t.Log("configuring coredns corefile")
		cfgmap := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "coredns"}, Data: map[string]string{"Corefile": corefile}}
		cfgmap, err = cluster.Client().CoreV1().ConfigMaps(testUDPIngressNamespace).Create(ctx, cfgmap, metav1.CreateOptions{})
		assert.NoError(t, err)

		t.Log("configuring a coredns deployent to deploy for UDP testing")
		container := generators.NewContainer("coredns", "coredns/coredns", 53)
		container.Ports[0].Protocol = corev1.ProtocolUDP
		container.VolumeMounts = []corev1.VolumeMount{{Name: "config-volume", MountPath: "/etc/coredns"}}
		container.Args = []string{"-conf", "/etc/coredns/Corefile"}
		deployment := generators.NewDeploymentForContainer(container)

		t.Log("configuring the coredns pod with a custom corefile")
		configVolume := corev1.Volume{
			Name: "config-volume",
			VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: cfgmap.Name},
				Items:                []corev1.KeyToPath{{Key: "Corefile", Path: "Corefile"}}}}}
		deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, configVolume)

		t.Log("deploying coredns")
		deployment, err = cluster.Client().AppsV1().Deployments(testUDPIngressNamespace).Create(ctx, deployment, metav1.CreateOptions{})
		assert.NoError(t, err)

		t.Logf("exposing deployment %s via service", deployment.Name)
		service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
		service, err = cluster.Client().CoreV1().Services(testUDPIngressNamespace).Create(ctx, service, metav1.CreateOptions{})
		assert.NoError(t, err)

		t.Log("exposing DNS service via UDPIngress")
		udp := &kongv1beta1.UDPIngress{
			ObjectMeta: metav1.ObjectMeta{
				Name: testName,
				Annotations: map[string]string{
					annotations.IngressClassKey: ingressClass,
				},
			},
			Spec: kongv1beta1.UDPIngressSpec{Rules: []kongv1beta1.UDPIngressRule{
				{
					Port: 9999,
					Backend: v1beta1.IngressBackend{
						ServiceName: service.Name,
						ServicePort: int(service.Spec.Ports[0].Port),
					},
				},
			}},
		}
		c, err := clientset.NewForConfig(cluster.Config())
		assert.NoError(t, err)
		udp, err = c.ConfigurationV1beta1().UDPIngresses(testUDPIngressNamespace).Create(ctx, udp, metav1.CreateOptions{})
		assert.NoError(t, err)
		start_time := time.Now().Nanosecond()

		t.Logf("checking udpingress %s status readiness.", udp.Name)
		ingCli := c.ConfigurationV1beta1().UDPIngresses(testUDPIngressNamespace)
		assert.Eventually(t, func() bool {
			curIng, err := ingCli.Get(ctx, udp.Name, metav1.GetOptions{})
			if err != nil || curIng == nil {
				return false
			}
			ingresses := curIng.Status.LoadBalancer.Ingress
			for _, ingress := range ingresses {
				if len(ingress.Hostname) > 0 || len(ingress.IP) > 0 {
					end_time := time.Now().Nanosecond()
					loop := end_time - start_time
					t.Logf("udpingress hostname %s or ip %s is ready to redirect traffic after %d nanoseconds.", ingress.Hostname, ingress.IP, loop)
					cost += loop
					return true
				}
			}
			return false
		}, 120*time.Second, 1*time.Second, true)
		cnt += 1
	}
	t.Logf("udp ingress average cost %d millisecond", cost/cnt/1000)
}

const corefile = `
.:53 {
    errors
    health {
       lameduck 5s
    }
    ready
    kubernetes cluster.local in-addr.arpa ip6.arpa {
       pods insecure
       fallthrough in-addr.arpa ip6.arpa
       ttl 30
    }
    forward . /etc/resolv.conf {
       max_concurrent 1000
    }
    cache 30
    loop
    reload
    loadbalance
}
`
