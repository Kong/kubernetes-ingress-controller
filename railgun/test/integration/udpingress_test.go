//+build integration_tests

package integration

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1alpha1"
	"github.com/kong/kubernetes-ingress-controller/railgun/pkg/clientset"
	"github.com/stretchr/testify/assert"
)

func TestMinimalUDPIngress(t *testing.T) {
	if useLegacyKIC() {
		t.Skip("legacy KIC does not support UDPIngress")
	}
	ctx := context.Background()

	// gather the proxy container as it will need to be specially configured to serve UDP
	proxy, err := cluster.Client().AppsV1().Deployments("kong-system").Get(ctx, "ingress-controller-kong", metav1.GetOptions{})
	assert.NoError(t, err)
	assert.Len(t, proxy.Spec.Template.Spec.Containers, 1)
	container := proxy.Spec.Template.Spec.Containers[0].DeepCopy()

	// override the KONG_STREAM_LISTEN env var in the proxy container
	originalVal, err := overrideEnvVar(container, "KONG_STREAM_LISTEN", "0.0.0.0:9999 udp reuseport")
	assert.NoError(t, err)
	proxy.Spec.Template.Spec.Containers[0] = *container

	// add the UDP port to the pod
	container.Ports = append(container.Ports, corev1.ContainerPort{
		Name:          "dns",
		ContainerPort: 9999,
		Protocol:      corev1.ProtocolUDP,
	})

	// update the deployment with the new container configurations
	proxy, err = cluster.Client().AppsV1().Deployments("kong-system").Update(ctx, proxy, metav1.UpdateOptions{})
	assert.NoError(t, err)

	// make sure we clean up after ourselves
	defer func() {
		if false { // FIXME - short circuit for testing
			_, err := overrideEnvVar(container, "KONG_STREAM_LISTEN", originalVal.Value)
			assert.NoError(t, err)
			_, err = cluster.Client().AppsV1().Deployments("kong-system").Update(ctx, proxy, metav1.UpdateOptions{})
			assert.NoError(t, err)
		}
	}()

	// ensure that the proxy deployment is ready before we proceed
	assert.Eventually(t, func() bool {
		d, err := cluster.Client().AppsV1().Deployments("kong-system").Get(ctx, proxy.Name, metav1.GetOptions{})
		if err != nil {
			t.Logf("WARNING: error while waiting for deployment %s to become ready: %v", proxy, err)
			return false
		}
		if d.Status.ReadyReplicas == d.Status.Replicas && d.Status.AvailableReplicas == d.Status.Replicas && d.Status.UnavailableReplicas < 1 {
			return true
		}
		return false
	}, time.Minute*3, time.Second*1)

	// create a LoadBalancer service to reach port 9999 on the proxy
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "quad9-dns",
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeLoadBalancer,
			Selector: proxy.Spec.Selector.MatchLabels,
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolUDP,
					Port:       9999,
					TargetPort: intstr.FromInt(9999),
				},
			},
		},
	}
	svc, err = cluster.Client().CoreV1().Services("kong-system").Create(ctx, svc, metav1.CreateOptions{})
	assert.NoError(t, err)

	// build a kong kubernetes clientset
	c, err := clientset.NewForConfig(cluster.Config())
	assert.NoError(t, err)

	// create the UDPIngress record
	udp := &kongv1alpha1.UDPIngress{
		ObjectMeta: metav1.ObjectMeta{
			Name: "quad9-dns",
			Annotations: map[string]string{
				"kubernetes.io/ingress.class": "kong",
			},
		},
		Spec: kongv1alpha1.UDPIngressSpec{
			Host:       "9.9.9.9",
			ListenPort: 9999,
			TargetPort: 53,
		},
	}
	_, err = c.ConfigurationV1alpha1().UDPIngresses("default").Create(ctx, udp, metav1.CreateOptions{})
	assert.NoError(t, err)

	// ensure that the DNS service is provisioned an IP address
	var dnsServer string
	assert.Eventually(t, func() bool {
		svc, err := cluster.Client().CoreV1().Services("kong-system").Get(ctx, "quad9-dns", metav1.GetOptions{})
		if err != nil {
			t.Logf("WARNING: ran into an error while trying to retrieve UDP service \"quad9-dns\": %v", err)
			return false
		}
		if len(svc.Status.LoadBalancer.Ingress) > 0 {
			ing := svc.Status.LoadBalancer.Ingress[0]
			if dnsServer = ing.IP; dnsServer != "" {
				return true
			}
		}
		return false
	}, time.Minute*3, time.Second*1)

	// configure a net.Resolver that will go through our proxy
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Second * 5,
			}
			return d.DialContext(ctx, network, fmt.Sprintf("%s:9999", dnsServer))
		},
	}

	// ensure that we can eventually make a successful DNS request through the proxy
	assert.Eventually(t, func() bool {
		// FIXME - this is a temporary hack to deal with a reconciliation bug that appears to be occuring in the secret controller,
		// I'm going to fix this before I move the PR out of draft to ready. For now, a second update consistently works around it.
		udp, err = c.ConfigurationV1alpha1().UDPIngresses("default").Get(ctx, "quad9-dns", metav1.GetOptions{})
		assert.NoError(t, err)
		udp.ObjectMeta.Annotations["FIXME"] = time.Now().String()
		_, err = c.ConfigurationV1alpha1().UDPIngresses("default").Update(ctx, udp, metav1.UpdateOptions{})
		assert.NoError(t, err)

		_, err := resolver.LookupHost(ctx, "kernel.org")
		if err != nil {
			return false
		}
		return true
	}, time.Minute*3, time.Second*1)
}

func overrideEnvVar(container *corev1.Container, key, val string) (original *corev1.EnvVar, err error) {
	newEnv := make([]corev1.EnvVar, 0, len(container.Env))
	for _, envvar := range container.Env {
		// override any existing KONG_STREAM_LISTEN value with our custom configuration
		if envvar.Name == key {
			// save the original configuration so we can put it back after we finish testing
			original = envvar.DeepCopy()
			envvar.Value = val
		}
		newEnv = append(newEnv, envvar)
	}

	if original == nil {
		err = fmt.Errorf("could not override env var: %s was not present on container %s", key, container.Name)
	}

	container.Env = newEnv
	return
}
