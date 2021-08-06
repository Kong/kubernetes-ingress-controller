//+build integration_tests

package integration

import (
	"context"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1beta1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/pkg/clientset"
)

func TestUDPIngressEssentials(t *testing.T) {
	t.Parallel()
	ns, cleanup := namespace(t)
	defer cleanup()

	t.Log("configuring coredns corefile")
	cfgmap := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "coredns"}, Data: map[string]string{"Corefile": corefile}}
	cfgmap, err := env.Cluster().Client().CoreV1().ConfigMaps(ns.Name).Create(ctx, cfgmap, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the coredns corefile %s", cfgmap.Name)
		assert.NoError(t, env.Cluster().Client().CoreV1().ConfigMaps(ns.Name).Delete(ctx, cfgmap.Name, metav1.DeleteOptions{}))
	}()

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
	deployment, err = env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up deployment %s", deployment.Name)
		assert.NoError(t, env.Cluster().Client().AppsV1().Deployments(ns.Name).Delete(ctx, deployment.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	service, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the service %s", service.Name)
		assert.NoError(t, env.Cluster().Client().CoreV1().Services(ns.Name).Delete(ctx, service.Name, metav1.DeleteOptions{}))
	}()

	t.Log("exposing DNS service via UDPIngress")
	udp := &kongv1beta1.UDPIngress{
		ObjectMeta: metav1.ObjectMeta{
			Name: "minudp",
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
	c, err := clientset.NewForConfig(env.Cluster().Config())
	assert.NoError(t, err)
	udp, err = c.ConfigurationV1beta1().UDPIngresses(ns.Name).Create(ctx, udp, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("ensuring UDPIngress %s is cleaned up", udp.Name)
		if err := c.ConfigurationV1beta1().UDPIngresses(ns.Name).Delete(ctx, udp.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}
		}
	}()

	t.Log("configurating a net.Resolver to resolve DNS via the proxy")
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Second * 5,
			}
			return d.DialContext(ctx, network, fmt.Sprintf("%s:9999", proxyUDPURL.Hostname()))
		},
	}

	t.Logf("checking udpingress %s status readiness.", udp.Name)
	ingCli := c.ConfigurationV1beta1().UDPIngresses(ns.Name)
	assert.Eventually(t, func() bool {
		curIng, err := ingCli.Get(ctx, udp.Name, metav1.GetOptions{})
		if err != nil || curIng == nil {
			return false
		}
		ingresses := curIng.Status.LoadBalancer.Ingress
		for _, ingress := range ingresses {
			if len(ingress.Hostname) > 0 || len(ingress.IP) > 0 {
				t.Logf("udpingress hostname %s or ip %s is ready to redirect traffic.", ingress.Hostname, ingress.IP)
				return true
			}
		}
		return false
	}, 120*time.Second, 1*time.Second, true)

	t.Logf("checking DNS to resolve via UDPIngress %s", udp.Name)
	assert.Eventually(t, func() bool {
		_, err := resolver.LookupHost(ctx, "kernel.org")
		return err == nil
	}, ingressWait, waitTick)

	t.Logf("tearing down UDPIngress %s and ensuring backends are torn down", udp.Name)
	assert.NoError(t, c.ConfigurationV1beta1().UDPIngresses(ns.Name).Delete(ctx, udp.Name, metav1.DeleteOptions{}))
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
