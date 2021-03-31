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

	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1alpha1"
	"github.com/kong/kubernetes-ingress-controller/railgun/pkg/clientset"
	"github.com/stretchr/testify/assert"
)

func TestMinimalUDPIngress(t *testing.T) {
	if useLegacyKIC() {
		t.Skip("legacy KIC does not support UDPIngress, skipping")
	}

	ctx := context.Background()
	namespace := "default"
	testName := "minudp"

	// UDPIngress requires an update to the proxy to open up a new listen port
	proxyLB, cleanup, err := updateProxyListeners(ctx, testName, "0.0.0.0:9999 udp reuseport", corev1.ContainerPort{
		Name:          testName,
		ContainerPort: 9999,
		Protocol:      corev1.ProtocolUDP,
	})
	assert.NoError(t, err)
	defer cleanup()

	// build a kong kubernetes clientset
	c, err := clientset.NewForConfig(cluster.Config())
	assert.NoError(t, err)

	// create the UDPIngress record
	udp := &kongv1alpha1.UDPIngress{
		ObjectMeta: metav1.ObjectMeta{
			Name: testName,
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
	_, err = c.ConfigurationV1alpha1().UDPIngresses(namespace).Create(ctx, udp, metav1.CreateOptions{})
	assert.NoError(t, err)

	// configure a net.Resolver that will go through our proxy
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Second * 5,
			}
			return d.DialContext(ctx, network, fmt.Sprintf("%s:9999", proxyLB.Status.LoadBalancer.Ingress[0].IP))
		},
	}

	// ensure that we can eventually make a successful DNS request through the proxy
	assert.Eventually(t, func() bool {
		_, err := resolver.LookupHost(ctx, "kernel.org")
		if err != nil {
			return false
		}
		return true
	}, udpWait, waitTick)
}
