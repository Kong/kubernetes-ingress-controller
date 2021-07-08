//+build performance_tests

package performance

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/pkg/annotations"
	"github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1beta1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/railgun/pkg/clientset"
	k8sgen "github.com/kong/kubernetes-testing-framework/pkg/generators/k8s"
)

func TestUDPIngressPerformance(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), ingressWait)
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
		container := k8sgen.NewContainer("coredns", "coredns/coredns", 53)
		container.Ports[0].Protocol = corev1.ProtocolUDP
		container.VolumeMounts = []corev1.VolumeMount{{Name: "config-volume", MountPath: "/etc/coredns"}}
		container.Args = []string{"-conf", "/etc/coredns/Corefile"}
		deployment := k8sgen.NewDeploymentForContainer(container)

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
		service := k8sgen.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
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
		s := time.Now().Nanosecond()
		t.Log("configurating a net.Resolver to resolve DNS via the proxy")
		resolver := &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: time.Second * 5,
				}
				return d.DialContext(ctx, network, fmt.Sprintf("%s:9999", KongInfo.ProxyUDPUrl.Hostname()))
			},
		}

		t.Logf("waiting for DNS to resolve via UDPIngress %s", udp.Name)
		assert.Eventually(t, func() bool {
			_, err := resolver.LookupHost(ctx, "kernel.org")
			if err != nil {
				return false
			}
			e := time.Now().Nanosecond()
			loop := e - s
			t.Logf("udp ingress loop cost %d nanosecond", loop)
			cost += loop
			return true
		}, ingressWait, waitTick)
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
