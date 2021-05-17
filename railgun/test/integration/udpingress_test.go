//+build integration_tests

package integration

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/pkg/annotations"
	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1alpha1"
	"github.com/kong/kubernetes-ingress-controller/railgun/pkg/clientset"
)

func TestMinimalUDPIngress(t *testing.T) {
	// TODO: once KIC 2.0 lands and pre v2 is gone, we can remove this check
	if useLegacyKIC() {
		t.Skip("legacy KIC does not support UDPIngress, skipping")
	}
	if dbmode := os.Getenv("TEST_DATABASE_MODE"); dbmode != "" && dbmode != "off" {
		t.Skip("v1alpha1.UDPIngress is only supported on DBLESS backend proxies at this time")
	}

	namespace := "default"
	testName := "minudp"
	ctx, cancel := context.WithTimeout(context.Background(), ingressWait)
	defer cancel()

	t.Log("configurating a net.Resolver to resolve DNS via the proxy")
	c, err := clientset.NewForConfig(cluster.Config())
	assert.NoError(t, err)
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

	t.Log("exposing DNS service via UDPIngress")
	udp := &kongv1alpha1.UDPIngress{
		ObjectMeta: metav1.ObjectMeta{
			Name: testName,
			Annotations: map[string]string{
				annotations.IngressClassKey: ingressClass,
			},
		},
		Spec: kongv1alpha1.UDPIngressSpec{
			Host:       "9.9.9.9",
			ListenPort: 9999,
			TargetPort: 53,
		},
	}
	udp, err = c.ConfigurationV1alpha1().UDPIngresses(namespace).Create(ctx, udp, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("ensuring UDPIngress %s is cleaned up", udp.Name)
		if err := c.ConfigurationV1alpha1().UDPIngresses(namespace).Delete(ctx, udp.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}
		}
	}()

	t.Logf("waiting for DNS to resolve via UDPIngress %s", udp.Name)
	assert.Eventually(t, func() bool {
		_, err := resolver.LookupHost(ctx, "kernel.org")
		if err != nil {
			return false
		}
		return true
	}, ingressWait, waitTick)

	t.Logf("tearing down UDPIngress %s and ensuring backends are torn down", udp.Name)
	assert.NoError(t, c.ConfigurationV1alpha1().UDPIngresses(namespace).Delete(ctx, udp.Name, metav1.DeleteOptions{}))
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
