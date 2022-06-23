//go:build integration_tests
// +build integration_tests

package integration

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	ktfkong "github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v2/test"
)

const (
	tlsRouteHostname      = "tlsroute.kong.example"
	tlsRouteExtraHostname = "extratlsroute.kong.example"
	tlsSecretName         = "secret-test"
)

var (
	tlsRouteTLSPairs = []TLSPair{
		{
			Cert: `-----BEGIN CERTIFICATE-----
MIIC/jCCAoSgAwIBAgIUVL6UYVDdH6peVNSOnOkCuYyhmrswCgYIKoZIzj0EAwIw
gbQxCzAJBgNVBAYTAlVTMRMwEQYDVQQIDApDYWxpZm9ybmlhMRYwFAYDVQQHDA1T
YW4gRnJhbmNpc2NvMRMwEQYDVQQKDApLb25nLCBJbmMuMRgwFgYDVQQLDA9UZWFt
IEt1YmVybmV0ZXMxHjAcBgNVBAMMFXRsc3JvdXRlLmtvbmcuZXhhbXBsZTEpMCcG
CSqGSIb3DQEJARYadGVzdEB0bHNyb3V0ZS5rb25nLmV4YW1wbGUwIBcNMjIwNjE2
MjExMjI4WhgPMjEyMjA1MjMyMTEyMjhaMIG0MQswCQYDVQQGEwJVUzETMBEGA1UE
CAwKQ2FsaWZvcm5pYTEWMBQGA1UEBwwNU2FuIEZyYW5jaXNjbzETMBEGA1UECgwK
S29uZywgSW5jLjEYMBYGA1UECwwPVGVhbSBLdWJlcm5ldGVzMR4wHAYDVQQDDBV0
bHNyb3V0ZS5rb25nLmV4YW1wbGUxKTAnBgkqhkiG9w0BCQEWGnRlc3RAdGxzcm91
dGUua29uZy5leGFtcGxlMHYwEAYHKoZIzj0CAQYFK4EEACIDYgAEQecfzsxmPwC0
6uNs3kyiLDb6brngM4ZtGXgwcGD393cbYmaunfBPRtxqh76RKdS9wzq4q+oB8dPs
QKgBNhlJTr+iFH9Di7bBZFcYqx+SnNUXZ0dDNBbW4rPVTJHQvdono1MwUTAdBgNV
HQ4EFgQU+OOVbqMcu+yXomZfnZ54LgIRNo4wHwYDVR0jBBgwFoAU+OOVbqMcu+yX
omZfnZ54LgIRNo4wDwYDVR0TAQH/BAUwAwEB/zAKBggqhkjOPQQDAgNoADBlAjBu
PMq+T+iTJ0yNvldYpB3BfdIhrv0EJQ9ALbB16nJwF91YV6YE7mdNP5rNVnoZ0nAC
MQDmnIpipMawjJWpfSPSZS1/iArz8YuBroWrGFXP62lwhCUp8RZweNnrLmmb/Aek
y3o=
-----END CERTIFICATE-----`,
			Key: `-----BEGIN PRIVATE KEY-----
MIG2AgEAMBAGByqGSM49AgEGBSuBBAAiBIGeMIGbAgEBBDDDRndgPYZaonVuqHiu
5uuYWI+A16BYLoUBnY0/9BL9U0s47G7LC/b05wE/7UPJEBKhZANiAARB5x/OzGY/
ALTq42zeTKIsNvpuueAzhm0ZeDBwYPf3dxtiZq6d8E9G3GqHvpEp1L3DOrir6gHx
0+xAqAE2GUlOv6IUf0OLtsFkVxirH5Kc1RdnR0M0Ftbis9VMkdC92ic=
-----END PRIVATE KEY-----`,
		},
		{
			Cert: `-----BEGIN CERTIFICATE-----
MIIDCDCCAo6gAwIBAgIUJB+Fq4hrxgiwhWLtqeAKp+NXigwwCgYIKoZIzj0EAwIw
gbkxCzAJBgNVBAYTAlVTMRMwEQYDVQQIDApDYWxpZm9ybmlhMRYwFAYDVQQHDA1T
YW4gRnJhbmNpc2NvMRMwEQYDVQQKDApLb25nLCBJbmMuMRgwFgYDVQQLDA9UZWFt
IEt1YmVybmV0ZXMxIzAhBgNVBAMMGmV4dHJhdGxzcm91dGUua29uZy5leGFtcGxl
MSkwJwYJKoZIhvcNAQkBFhp0ZXN0QHRsc3JvdXRlLmtvbmcuZXhhbXBsZTAgFw0y
MjA2MjIyMDIwNDlaGA8yMTIyMDUyOTIwMjA0OVowgbkxCzAJBgNVBAYTAlVTMRMw
EQYDVQQIDApDYWxpZm9ybmlhMRYwFAYDVQQHDA1TYW4gRnJhbmNpc2NvMRMwEQYD
VQQKDApLb25nLCBJbmMuMRgwFgYDVQQLDA9UZWFtIEt1YmVybmV0ZXMxIzAhBgNV
BAMMGmV4dHJhdGxzcm91dGUua29uZy5leGFtcGxlMSkwJwYJKoZIhvcNAQkBFhp0
ZXN0QHRsc3JvdXRlLmtvbmcuZXhhbXBsZTB2MBAGByqGSM49AgEGBSuBBAAiA2IA
BACgptITKMoxBz67FTxi9eP0CcnIabUu4AlkP7IOSkgprzsPGUfgn6sSv88IxHbn
0qSIxMi1OjoK+m12a5eayYYnr1kiy9Qvm0jCubCDog03534rrMqjKFTimMSk/4U4
A6NTMFEwHQYDVR0OBBYEFN3kitZnxny13r7TajZ74IkwCq4uMB8GA1UdIwQYMBaA
FN3kitZnxny13r7TajZ74IkwCq4uMA8GA1UdEwEB/wQFMAMBAf8wCgYIKoZIzj0E
AwIDaAAwZQIwGYtGE0xOKdiObmVUIxlc5Iif9cVwzfvaMF0wiuuth9Hxd3n40XPv
aof6F4WdQihFAjEA/heIDoQActLLXrhFvxS6JP/XPT3C086lK1mq3inRGIvYX/1r
/gHROKq7BRLjo6FS
-----END CERTIFICATE-----`,
			Key: `-----BEGIN PRIVATE KEY-----
MIG2AgEAMBAGByqGSM49AgEGBSuBBAAiBIGeMIGbAgEBBDBHOv6UxSf7MbyPOllv
0Sb/hnXf+UfTblLA8TeoKa4Hr9RjoB0QYLFHLDFPMg5eplGhZANiAAQAoKbSEyjK
MQc+uxU8YvXj9AnJyGm1LuAJZD+yDkpIKa87DxlH4J+rEr/PCMR259KkiMTItTo6
CvptdmuXmsmGJ69ZIsvUL5tIwrmwg6INN+d+K6zKoyhU4pjEpP+FOAM=
-----END PRIVATE KEY-----`,
		},
	}
)

func TestTLSRouteEssentials(t *testing.T) {
	backendPort := gatewayv1alpha2.PortNumber(tcpEchoPort)
	t.Log("locking TLS port")
	tlsMutex.Lock()
	defer func() {
		t.Log("unlocking TLS port")
		tlsMutex.Unlock()
	}()

	ns, cleaner := setup(t)
	defer func() { assert.NoError(t, cleaner.Cleanup(ctx)) }()

	t.Log("getting gateway client")
	gatewayClient, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("configuring secrets")
	secrets := []*corev1.Secret{
		{
			ObjectMeta: metav1.ObjectMeta{
				UID:       types.UID("7428fb98-180b-4702-a91f-61351a33c6e8"),
				Name:      tlsSecretName,
				Namespace: ns.Name,
			},
			Data: map[string][]byte{
				"tls.crt": []byte(tlsRouteTLSPairs[0].Cert),
				"tls.key": []byte(tlsRouteTLSPairs[0].Key),
			},
		},
	}

	t.Log("deploying secrets")
	secret1, err := env.Cluster().Client().CoreV1().Secrets(ns.Name).Create(ctx, secrets[0], metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(secret1)

	t.Log("deploying a supported gatewayclass to the test cluster")
	gatewayClassName := uuid.NewString()
	gwc, err := DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
	require.NoError(t, err)
	cleaner.Add(gwc)

	t.Log("deploying a gateway to the test cluster using unmanaged gateway mode and port 8899")
	gatewayName := uuid.NewString()
	hostname := gatewayv1alpha2.Hostname(tlsRouteHostname)
	gateway, err := DeployGateway(ctx, gatewayClient, ns.Name, gatewayClassName, func(gw *gatewayv1alpha2.Gateway) {
		gw.Name = gatewayName
		gw.Spec.Listeners = []gatewayv1alpha2.Listener{{
			Name:     "tls",
			Protocol: gatewayv1alpha2.TLSProtocolType,
			Port:     gatewayv1alpha2.PortNumber(ktfkong.DefaultTLSServicePort),
			Hostname: &hostname,
			TLS: &gatewayv1alpha2.GatewayTLSConfig{
				CertificateRefs: []*gatewayv1alpha2.SecretObjectReference{
					{
						Name: gatewayv1alpha2.ObjectName(tlsSecretName),
					},
				},
			},
		}}
	})
	require.NoError(t, err)
	cleaner.Add(gateway)

	t.Log("creating a tcpecho pod to test TLSRoute traffic routing")

	container := generators.NewContainer("tcpecho-1", test.TCPEchoImage, tcpEchoPort)
	// go-echo sends a "Running on Pod <UUID>." immediately on connecting
	testUUID := uuid.NewString()
	container.Env = []corev1.EnvVar{
		{
			Name:  "POD_NAME",
			Value: testUUID,
		},
	}
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err = env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(deployment)

	t.Log("creating an additional tcpecho pod to test TLSRoute multiple backendRef loadbalancing")
	container2 := generators.NewContainer("tcpecho-2", test.TCPEchoImage, tcpEchoPort)
	// go-echo sends a "Running on Pod <UUID>." immediately on connecting
	testUUID2 := uuid.NewString()
	container2.Env = []corev1.EnvVar{
		{
			Name:  "POD_NAME",
			Value: testUUID2,
		},
	}
	deployment2 := generators.NewDeploymentForContainer(container2)
	deployment2, err = env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment2, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(deployment2)

	t.Logf("exposing deployment %s/%s via service", deployment.Namespace, deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	service, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(service)

	t.Logf("exposing deployment %s/%s via service", deployment2.Namespace, deployment2.Name)
	service2 := generators.NewServiceForDeployment(deployment2, corev1.ServiceTypeLoadBalancer)
	service2, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service2, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(service2)

	t.Logf("creating a tlsroute to access deployment %s via kong", deployment.Name)
	tlsRoute := &gatewayv1alpha2.TLSRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:        uuid.NewString(),
			Annotations: map[string]string{},
		},
		Spec: gatewayv1alpha2.TLSRouteSpec{
			CommonRouteSpec: gatewayv1alpha2.CommonRouteSpec{
				ParentRefs: []gatewayv1alpha2.ParentReference{{
					Name: gatewayv1alpha2.ObjectName(gatewayName),
				}},
			},
			Hostnames: []gatewayv1alpha2.Hostname{tlsRouteHostname},
			Rules: []gatewayv1alpha2.TLSRouteRule{{
				BackendRefs: []gatewayv1alpha2.BackendRef{{
					BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
						Name: gatewayv1alpha2.ObjectName(service.Name),
						Port: &backendPort,
					},
				}},
			}},
		},
	}
	tlsRoute, err = gatewayClient.GatewayV1alpha2().TLSRoutes(ns.Name).Create(ctx, tlsRoute, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(tlsRoute)

	t.Log("verifying that the Gateway gets linked to the route via status")
	callback := GetGatewayIsLinkedCallback(t, gatewayClient, gatewayv1alpha2.TLSProtocolType, ns.Name, tlsRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that the tcpecho is responding properly over TLS")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTLSServicePort),
			testUUID, tlsRouteHostname, tlsRouteHostname)
		return err == nil && responded == true
	}, ingressWait, waitTick)

	t.Log("removing the parentrefs from the TLSRoute")
	oldParentRefs := tlsRoute.Spec.ParentRefs
	require.Eventually(t, func() bool {
		tlsRoute, err = gatewayClient.GatewayV1alpha2().TLSRoutes(ns.Name).Get(ctx, tlsRoute.Name, metav1.GetOptions{})
		require.NoError(t, err)
		tlsRoute.Spec.ParentRefs = nil
		tlsRoute, err = gatewayClient.GatewayV1alpha2().TLSRoutes(ns.Name).Update(ctx, tlsRoute, metav1.UpdateOptions{})
		return err == nil
	}, time.Minute, time.Second)

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	callback = GetGatewayIsUnlinkedCallback(t, gatewayClient, gatewayv1alpha2.TLSProtocolType, ns.Name, tlsRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that the tcpecho is no longer responding")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTLSServicePort),
			testUUID, tlsRouteHostname, tlsRouteHostname)
		return responded == false && errors.Is(err, io.EOF)
	}, ingressWait, waitTick)

	t.Log("putting the parentRefs back")
	require.Eventually(t, func() bool {
		tlsRoute, err = gatewayClient.GatewayV1alpha2().TLSRoutes(ns.Name).Get(ctx, tlsRoute.Name, metav1.GetOptions{})
		require.NoError(t, err)
		tlsRoute.Spec.ParentRefs = oldParentRefs
		tlsRoute, err = gatewayClient.GatewayV1alpha2().TLSRoutes(ns.Name).Update(ctx, tlsRoute, metav1.UpdateOptions{})
		return err == nil
	}, time.Minute, time.Second)

	t.Log("verifying that the Gateway gets linked to the route via status")
	callback = GetGatewayIsLinkedCallback(t, gatewayClient, gatewayv1alpha2.TLSProtocolType, ns.Name, tlsRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that putting the parentRefs back results in the routes becoming available again")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTLSServicePort),
			testUUID, tlsRouteHostname, tlsRouteHostname)
		return err == nil && responded == true
	}, ingressWait, waitTick)

	t.Log("deleting the GatewayClass")
	require.NoError(t, gatewayClient.GatewayV1alpha2().GatewayClasses().Delete(ctx, gwc.Name, metav1.DeleteOptions{}))

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	callback = GetGatewayIsUnlinkedCallback(t, gatewayClient, gatewayv1alpha2.TLSProtocolType, ns.Name, tlsRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that the data-plane configuration from the TLSRoute gets dropped with the GatewayClass now removed")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTLSServicePort),
			testUUID, tlsRouteHostname, tlsRouteHostname)
		return responded == false && errors.Is(err, io.EOF)
	}, ingressWait, waitTick)

	t.Log("putting the GatewayClass back")
	gwc, err = DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
	require.NoError(t, err)

	t.Log("verifying that the Gateway gets linked to the route via status")
	callback = GetGatewayIsLinkedCallback(t, gatewayClient, gatewayv1alpha2.TLSProtocolType, ns.Name, tlsRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that creating the GatewayClass again triggers reconciliation of TLSRoutes and the route becomes available again")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTLSServicePort),
			testUUID, tlsRouteHostname, tlsRouteHostname)
		return err == nil && responded == true
	}, ingressWait, waitTick)

	t.Log("deleting the Gateway")
	require.NoError(t, gatewayClient.GatewayV1alpha2().Gateways(ns.Name).Delete(ctx, gatewayName, metav1.DeleteOptions{}))

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	callback = GetGatewayIsUnlinkedCallback(t, gatewayClient, gatewayv1alpha2.TLSProtocolType, ns.Name, tlsRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that the data-plane configuration from the TLSRoute gets dropped with the Gateway now removed")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTLSServicePort),
			testUUID, tlsRouteHostname, tlsRouteHostname)
		return responded == false && errors.Is(err, io.EOF)
	}, ingressWait, waitTick)

	t.Log("putting the Gateway back")
	gateway, err = DeployGateway(ctx, gatewayClient, ns.Name, gatewayClassName, func(gw *gatewayv1alpha2.Gateway) {
		gw.Name = gatewayName
		gw.Spec.Listeners = []gatewayv1alpha2.Listener{{
			Name:     "tls",
			Protocol: gatewayv1alpha2.TLSProtocolType,
			Port:     gatewayv1alpha2.PortNumber(ktfkong.DefaultTLSServicePort),
			Hostname: &hostname,
			TLS: &gatewayv1alpha2.GatewayTLSConfig{
				CertificateRefs: []*gatewayv1alpha2.SecretObjectReference{
					{
						Name: gatewayv1alpha2.ObjectName(tlsSecretName),
					},
				},
			},
		}}
	})
	require.NoError(t, err)

	t.Log("verifying that the Gateway gets linked to the route via status")
	callback = GetGatewayIsLinkedCallback(t, gatewayClient, gatewayv1alpha2.TLSProtocolType, ns.Name, tlsRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that creating the Gateway again triggers reconciliation of TLSRoutes and the route becomes available again")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTLSServicePort),
			testUUID, tlsRouteHostname, tlsRouteHostname)
		return err == nil && responded == true
	}, ingressWait, waitTick)

	t.Log("adding an additional backendRef to the TLSRoute")
	require.Eventually(t, func() bool {
		tlsRoute, err = gatewayClient.GatewayV1alpha2().TLSRoutes(ns.Name).Get(ctx, tlsRoute.Name, metav1.GetOptions{})
		require.NoError(t, err)

		tlsRoute.Spec.Rules[0].BackendRefs = []gatewayv1alpha2.BackendRef{
			{
				BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
					Name: gatewayv1alpha2.ObjectName(service.Name),
					Port: &backendPort,
				},
			},
			{
				BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
					Name: gatewayv1alpha2.ObjectName(service2.Name),
					Port: &backendPort,
				},
			},
		}

		tlsRoute, err = gatewayClient.GatewayV1alpha2().TLSRoutes(ns.Name).Update(ctx, tlsRoute, metav1.UpdateOptions{})
		return err == nil
	}, ingressWait, waitTick)

	t.Log("verifying that the TLSRoute is now load-balanced between two services")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTLSServicePort),
			testUUID, tlsRouteHostname, tlsRouteHostname)
		return err == nil && responded == true
	}, ingressWait, waitTick)
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTLSServicePort),
			testUUID2, tlsRouteHostname, tlsRouteHostname)
		return err == nil && responded == true
	}, ingressWait, waitTick)

	t.Log("deleting both GatewayClass and Gateway rapidly")
	require.NoError(t, gatewayClient.GatewayV1alpha2().GatewayClasses().Delete(ctx, gwc.Name, metav1.DeleteOptions{}))
	require.NoError(t, gatewayClient.GatewayV1alpha2().Gateways(ns.Name).Delete(ctx, gateway.Name, metav1.DeleteOptions{}))

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	callback = GetGatewayIsUnlinkedCallback(t, gatewayClient, gatewayv1alpha2.TLSProtocolType, ns.Name, tlsRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that the data-plane configuration from the TLSRoute does not get orphaned with the GatewayClass and Gateway gone")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTLSServicePort),
			testUUID, tlsRouteHostname, tlsRouteHostname)
		return responded == false && errors.Is(err, io.EOF)
	}, ingressWait, waitTick)
}

// TestTLSRouteReferencePolicy tests cross-namespace certificate references. These are technically implemented within
// Gateway Listeners, but require an attached Route to see the associated certificate behavior on the proxy
func TestTLSRouteReferencePolicy(t *testing.T) {
	backendPort := gatewayv1alpha2.PortNumber(tcpEchoPort)
	t.Log("locking TLS port")
	tlsMutex.Lock()
	defer func() {
		t.Log("unlocking TLS port")
		tlsMutex.Unlock()
	}()

	ns, cleaner := setup(t)
	defer func() { assert.NoError(t, cleaner.Cleanup(ctx)) }()

	otherNs, err := clusters.GenerateNamespace(ctx, env.Cluster(), t.Name())
	require.NoError(t, err)
	cleaner.AddNamespace(otherNs)

	t.Log("getting the gateway client")
	gatewayClient, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("configuring secrets")
	secrets := []*corev1.Secret{
		{
			ObjectMeta: metav1.ObjectMeta{
				UID:       types.UID("7428fb98-180b-4702-a91f-61351a33c6e8"),
				Name:      tlsSecretName,
				Namespace: ns.Name,
			},
			Data: map[string][]byte{
				"tls.crt": []byte(tlsRouteTLSPairs[0].Cert),
				"tls.key": []byte(tlsRouteTLSPairs[0].Key),
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				UID:  types.UID("7428fb98-180b-4702-a91f-61351a33c6e9"),
				Name: "secret2",
			},
			Data: map[string][]byte{
				"tls.crt": []byte(tlsRouteTLSPairs[1].Cert),
				"tls.key": []byte(tlsRouteTLSPairs[1].Key),
			},
		},
	}

	t.Log("deploying secrets")
	secret1, err := env.Cluster().Client().CoreV1().Secrets(ns.Name).Create(ctx, secrets[0], metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(secret1)
	secret2, err := env.Cluster().Client().CoreV1().Secrets(otherNs.Name).Create(ctx, secrets[1], metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(secret2)

	t.Log("deploying a gateway to the test cluster using unmanaged gateway mode")
	gateway, err := DeployGateway(ctx, gatewayClient, ns.Name, managedGatewayClassName, func(gw *gatewayv1alpha2.Gateway) {
		hostname := gatewayv1alpha2.Hostname(tlsRouteHostname)
		otherHostname := gatewayv1alpha2.Hostname(tlsRouteExtraHostname)
		otherNamespace := gatewayv1alpha2.Namespace(otherNs.Name)
		gw.Spec.Listeners = []gatewayv1alpha2.Listener{
			{
				Name:     "tls",
				Protocol: gatewayv1alpha2.TLSProtocolType,
				Port:     gatewayv1alpha2.PortNumber(ktfkong.DefaultTLSServicePort),
				Hostname: &hostname,
				TLS: &gatewayv1alpha2.GatewayTLSConfig{
					CertificateRefs: []*gatewayv1alpha2.SecretObjectReference{
						{
							Name: gatewayv1alpha2.ObjectName(secrets[0].Name),
						},
					},
				},
			},
			{
				Name:     "tlsother",
				Protocol: gatewayv1alpha2.TLSProtocolType,
				Port:     gatewayv1alpha2.PortNumber(ktfkong.DefaultTLSServicePort),
				Hostname: &otherHostname,
				TLS: &gatewayv1alpha2.GatewayTLSConfig{
					CertificateRefs: []*gatewayv1alpha2.SecretObjectReference{
						{
							Name:      gatewayv1alpha2.ObjectName(secrets[1].Name),
							Namespace: &otherNamespace,
						},
					},
				},
			},
		}
	})

	require.NoError(t, err)
	cleaner.Add(gateway)

	secret2Name := gatewayv1alpha2.ObjectName(secrets[1].Name)
	t.Logf("creating a reference policy that permits tcproute access from %s to services in %s", ns.Name, otherNs.Name)
	policy := &gatewayv1alpha2.ReferencePolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:        uuid.NewString(),
			Annotations: map[string]string{},
		},
		Spec: gatewayv1alpha2.ReferencePolicySpec{
			From: []gatewayv1alpha2.ReferencePolicyFrom{
				{
					Group:     gatewayv1alpha2.Group("gateway.networking.k8s.io"),
					Kind:      gatewayv1alpha2.Kind("Gateway"),
					Namespace: gatewayv1alpha2.Namespace(gateway.Namespace),
				},
			},
			To: []gatewayv1alpha2.ReferencePolicyTo{
				{
					Group: gatewayv1alpha2.Group(""),
					Kind:  gatewayv1alpha2.Kind("Secret"),
					Name:  &secret2Name,
				},
			},
		},
	}

	policy, err = gatewayClient.GatewayV1alpha2().ReferencePolicies(otherNs.Name).Create(ctx, policy, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(policy)

	t.Log("creating a tcpecho pod to test TLSRoute traffic routing")
	container := generators.NewContainer("tcpecho", test.TCPEchoImage, tcpEchoPort)
	// go-echo sends a "Running on Pod <UUID>." immediately on connecting
	testUUID := uuid.NewString()
	container.Env = []corev1.EnvVar{
		{
			Name:  "POD_NAME",
			Value: testUUID,
		},
	}
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err = env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(deployment)

	t.Logf("exposing deployment %s/%s via service", deployment.Namespace, deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	service, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(service)

	t.Logf("creating a tlsroute to access deployment %s via kong", deployment.Name)
	tlsroute := &gatewayv1alpha2.TLSRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:        uuid.NewString(),
			Annotations: map[string]string{},
		},
		Spec: gatewayv1alpha2.TLSRouteSpec{
			CommonRouteSpec: gatewayv1alpha2.CommonRouteSpec{
				ParentRefs: []gatewayv1alpha2.ParentReference{{
					Name: gatewayv1alpha2.ObjectName(gateway.Name),
				}},
			},
			Hostnames: []gatewayv1alpha2.Hostname{tlsRouteHostname, tlsRouteExtraHostname},
			Rules: []gatewayv1alpha2.TLSRouteRule{{
				BackendRefs: []gatewayv1alpha2.BackendRef{{
					BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
						Name: gatewayv1alpha2.ObjectName(service.Name),
						Port: &backendPort,
					},
				}},
			}},
		},
	}
	tlsroute, err = gatewayClient.GatewayV1alpha2().TLSRoutes(ns.Name).Create(ctx, tlsroute, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(tlsroute)

	t.Log("verifying that the tcpecho is responding properly over TLS")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTLSServicePort),
			testUUID, tlsRouteHostname, tlsRouteHostname)
		return err == nil && responded == true
	}, ingressWait, waitTick)

	t.Log("verifying that the tcpecho route can also serve certificates permitted by a ReferencePolicy with a named To")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTLSServicePort),
			testUUID, tlsRouteExtraHostname, tlsRouteExtraHostname)
		return err == nil && responded == true
	}, ingressWait, waitTick)

	t.Log("verifying that using the wrong name in the ReferencePolicy removes the related certificate")
	badName := gatewayv1alpha2.ObjectName("garbage")
	policy.Spec.To[0].Name = &badName
	policy, err = gatewayClient.GatewayV1alpha2().ReferencePolicies(otherNs.Name).Update(ctx, policy, metav1.UpdateOptions{})
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTLSServicePort),
			testUUID, tlsRouteExtraHostname, tlsRouteExtraHostname)
		return err != nil && responded == false
	}, ingressWait, waitTick)

	t.Log("verifying the certificate returns when using a ReferencePolicy with no name restrictions")
	policy.Spec.To[0].Name = nil
	policy, err = gatewayClient.GatewayV1alpha2().ReferencePolicies(otherNs.Name).Update(ctx, policy, metav1.UpdateOptions{})
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTLSServicePort),
			testUUID, tlsRouteExtraHostname, tlsRouteExtraHostname)
		return err == nil && responded == true
	}, ingressWait, waitTick)
}

// tlsEchoResponds takes a TLS address URL and a Pod name and checks if a
// go-echo instance is running on that Pod at that address using hostname for SNI.
// It compares an expected message and its length against an expected message, returning true
// if it is and false and an error explanation if it is not
func tlsEchoResponds(url string, podName string, hostname, certHostname string) (bool, error) {
	dialer := net.Dialer{Timeout: time.Second * 10}
	conn, err := tls.DialWithDialer(&dialer,
		"tcp",
		url,
		&tls.Config{
			InsecureSkipVerify: true, //nolint:gosec
			ServerName:         hostname,
		})
	if err != nil {
		return false, err
	}
	defer conn.Close()

	cert := conn.ConnectionState().PeerCertificates[0]
	if cert.Subject.CommonName != certHostname {
		return false, fmt.Errorf("expected certificate with cn=%s, got cn=%s", certHostname, cert.Subject.CommonName)
	}

	header := []byte(fmt.Sprintf("Running on Pod %s.", podName))
	message := []byte("testing tlsroute")

	wrote, err := conn.Write(message)
	if err != nil {
		return false, err
	}

	if wrote != len(message) {
		return false, fmt.Errorf("wrote message of size %d, expected %d", wrote, len(message))
	}

	if err := conn.SetDeadline(time.Now().Add(time.Second * 5)); err != nil {
		return false, err
	}

	headerResponse := make([]byte, len(header)+1)
	read, err := conn.Read(headerResponse)
	if err != nil {
		return false, err
	}

	if read != len(header)+1 { // add 1 for newline
		return false, fmt.Errorf("read %d bytes but expected %d", read, len(header)+1)
	}

	if !bytes.Contains(headerResponse, header) {
		return false, fmt.Errorf(`expected header response "%s", received: "%s"`, string(header), string(headerResponse))
	}

	messageResponse := make([]byte, wrote+1)
	read, err = conn.Read(messageResponse)
	if err != nil {
		return false, err
	}

	if read != len(message) {
		return false, fmt.Errorf("read %d bytes but expected %d", read, len(message))
	}

	if !bytes.Contains(messageResponse, message) {
		return false, fmt.Errorf(`expected message response "%s", received: "%s"`, string(message), string(messageResponse))
	}

	return true, nil
}
