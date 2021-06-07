//+build integration_tests

package integration

import (
	"context"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/pkg/annotations"
	"github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1beta1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/railgun/pkg/clientset"
	k8sgen "github.com/kong/kubernetes-testing-framework/pkg/generators/k8s"
)

func TestMinimalUDPIngress(t *testing.T) {
	// TODO: once KIC 2.0 lands and pre v2 is gone, we can remove this check
	if useLegacyKIC() {
		t.Skip("legacy KIC does not support UDPIngress, skipping")
	}
	if dbmode != "" && dbmode != "off" {
		t.Skip("v1beta1.UDPIngress is only supported on DBLESS backend proxies at this time")
	}

	testName := "minudp"
	namespace := corev1.NamespaceDefault
	ctx, cancel := context.WithTimeout(context.Background(), ingressWait)
	defer cancel()

	t.Log("configuring coredns corefile")
	cfgmap := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "coredns"}, Data: map[string]string{"Corefile": corefile}}
	_, err := cluster.Client().CoreV1().ConfigMaps(namespace).Create(ctx, cfgmap, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the coredns corefile %s", cfgmap.Name)
		assert.NoError(t, cluster.Client().CoreV1().ConfigMaps(namespace).Delete(ctx, cfgmap.Name, metav1.DeleteOptions{}))
	}()

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
	deployment, err = cluster.Client().AppsV1().Deployments(namespace).Create(ctx, deployment, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up deployment %s", deployment.Name)
		assert.NoError(t, cluster.Client().AppsV1().Deployments(namespace).Delete(ctx, deployment.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := k8sgen.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	service, err = cluster.Client().CoreV1().Services(namespace).Create(ctx, service, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the service %s", service.Name)
		assert.NoError(t, cluster.Client().CoreV1().Services(namespace).Delete(ctx, service.Name, metav1.DeleteOptions{}))
	}()

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
	udp, err = c.ConfigurationV1beta1().UDPIngresses(namespace).Create(ctx, udp, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("ensuring UDPIngress %s is cleaned up", udp.Name)
		if err := c.ConfigurationV1beta1().UDPIngresses(namespace).Delete(ctx, udp.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}
		}
	}()

	t.Log("configurating a net.Resolver to resolve DNS via the proxy")
	p := proxyReady()
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Second * 5,
			}
			return d.DialContext(ctx, network, fmt.Sprintf("%s:9999", p.ProxyUDPUrl.Hostname()))
		},
	}

	t.Logf("waiting for DNS to resolve via UDPIngress %s", udp.Name)
	assert.Eventually(t, func() bool {
		_, err := resolver.LookupHost(ctx, "kernel.org")
		if err != nil {
			t.Logf("WARNING: failure to lookup kernel.org: %v", err)
			return false
		}
		return true
	}, ingressWait, waitTick)

	t.Logf("tearing down UDPIngress %s and ensuring backends are torn down", udp.Name)
	assert.NoError(t, c.ConfigurationV1beta1().UDPIngresses(namespace).Delete(ctx, udp.Name, metav1.DeleteOptions{}))
	assert.Eventually(t, func() bool {
		_, err := resolver.LookupHost(ctx, "kernel.org")
		if err != nil {
			if strings.Contains(err.Error(), "i/o timeout") {
				return true
			}
		}
		return false
	}, ingressWait, waitTick)
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
