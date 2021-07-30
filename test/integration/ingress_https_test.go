//+build integration_tests

package integration

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/internal/annotations"
)

type TLSPair struct {
	Key, Cert string
}

var (
	tlsPairs = []TLSPair{
		{
			Cert: `-----BEGIN CERTIFICATE-----
MIIC2DCCAcACCQC32eFOsWpKojANBgkqhkiG9w0BAQsFADAuMRcwFQYDVQQDDA5z
ZWN1cmUtZm9vLWJhcjETMBEGA1UECgwKa29uZ2hxLm9yZzAeFw0xODEyMTgyMTI4
MDBaFw0xOTEyMTgyMTI4MDBaMC4xFzAVBgNVBAMMDnNlY3VyZS1mb28tYmFyMRMw
EQYDVQQKDAprb25naHEub3JnMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKC
AQEAqhl/HSwV6PbMv+cMFU9X+HuM7QbNNPh39GKa4pkxzFgiAnuuJ4jw9V/bzsEy
S+ZIyjzo+QKB1LzmgdcX4vkdI22BjxUd9HPHdZxtv3XilbNmSk9UOl2Hh1fORJoS
7YH+VbvVwiz5lo7qKRepbg/jcKkbs6AUE0YWFygtDLTvhP2qkphQkxZ0m8qroW91
CWgI73Ar6U2W/YQBRI3+LwtsKo0p2ASDijvqxElQBgBIiyGIr0RZc5pkCJ1eQdDB
2F6XaMfpeEyBj0MxypNL4S9HHfchOt55J1KOzYnUPkQnSoxp6oEjef4Q/ZCj5BRL
EGZnTb3tbwzHZCxGtgl9KqO9pQIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQAKQ5BX
kkBL+alERL31hsOgWgRiUMw+sPDtRS96ozUlPtVvAg9XFdpY4ldtWkxFcmBnhKzp
UewjrHkf9rR16NISwUTjlGIwaJu/ACQrY15v+r301Crq2DV+GjiUJFVuT495dp/l
0LZbt2Sh/uD+r3UNTcJpJ7jb1V0UP7FWXFj8oafsoFSgmxAPjpKQySTC54JK4AYb
QSnWu1nQLyohnrB9qLZhe2+jOQZnkKuCcWJQ5njvU6SxT3SOKE5XaOZCezEQ6IVL
U47YCCXsq+7wKWXBhKl4H2Ztk6x3HOC56l0noXWezsMfrou/kjwGuuViGnrjqelS
WQ7uVeNCUBY+l+qY
-----END CERTIFICATE-----`,
			Key: `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCqGX8dLBXo9sy/
5wwVT1f4e4ztBs00+Hf0YprimTHMWCICe64niPD1X9vOwTJL5kjKPOj5AoHUvOaB
1xfi+R0jbYGPFR30c8d1nG2/deKVs2ZKT1Q6XYeHV85EmhLtgf5Vu9XCLPmWjuop
F6luD+NwqRuzoBQTRhYXKC0MtO+E/aqSmFCTFnSbyquhb3UJaAjvcCvpTZb9hAFE
jf4vC2wqjSnYBIOKO+rESVAGAEiLIYivRFlzmmQInV5B0MHYXpdox+l4TIGPQzHK
k0vhL0cd9yE63nknUo7NidQ+RCdKjGnqgSN5/hD9kKPkFEsQZmdNve1vDMdkLEa2
CX0qo72lAgMBAAECggEADxMTYNJ3Xp4Ap0EioQDXGv5YDul7ZiZe+xmCAHLzJtjo
qq+rT3WjZRuJr1kPzAosiT+8pdTDDMdw5jDZvRO2sV0TDksgzHk2RAYI897OpdWw
SwWcwU9oo2X0sb+1zbang5GR8BNsSxt/RQUDzu05itJx0gltvgeIDaVR2L5wO6ja
USa8OVuj/92XtIIve9OtyK9jAzgR6LQOTFrCCEv89/vmy5Bykv4Uz8s8swZmTs3v
XJmAmruHGuSLMfXk8lBRp/gVyNTi3uMsdph5AJbVKnra5TZLguEozZKbLdNUYk0p
+aAc7rxDcH2sPqa/7DwRvei9dvd5oB3VJlxGVgC8AQKBgQDfznRSSKAD15hoSDzt
cKNyhLgWAL+MD0jhHKUy3x+Z9OCvf0DVnmru5HfQKq5UfT0t8VTRPGKmOtAMD4cf
LYjIurvMvpVzQGSJfhtHQuULZTh3dfsM7xivMqSV+9txklMAakM7vGQlOQxhrScM
21Mp5LWDU6+e2pFCrQPop0IPkQKBgQDCkVE+dou2yFuJx3uytCH1yKPSy9tkdhQH
dGF12B5dq8MZZozAz5P9YN/COa9WjsNKDqWbEgLEksEQUq4t8SBjHnSV/D3x7rEF
qgwii0GETYxax6gms8nueIqWZQf+0NbX7Gc5mTqeVb7v3TrhsKr0VNMFRXXQwP2E
M/pxJq8q1QKBgQC3rH7oXLP+Ez0AMHDYSL3LKULOw/RvpMeh/9lQA6+ysTaIsP3r
kuSdhCEUVULXEiVYhBug0FcBp3jAvSmem8cLPb0Mjkim2mzoLfeDJ1JEZODPoaLU
fZEbj4tlj9oLvhOiXpMo/jaOGeCgdPN8aK86zXlt+wtBao0WVFnF4SalEQKBgQC1
uLfi2SGgs/0a8B/ORoO5ZY3s4c2lRMtsMvyb7iBeaIAuByPLKZUVABe89deXxnsL
fiaacPX41wBO2IoqCp2vNdC6DP9mKQNZQPtYgCvPAAbo+rVIgH9HpXn7AZ24FyGy
RfAbUcv3+in9KelGxZTF4zu8HqXtNXMSuOFeMT1FiQKBgF0R+IFDGHhD4nudAQvo
hncXsgyzK6QUzak6HmFji/CMZ6EU9q6A67JkiEWrYoKqIAKZ2Og8+Eucr/rDdGWc
kqlmLPBJAJeUsP/9KidBjTE5mIbn/2n089VPMBvnlt2xIcuB6+zrf2NjvlcZEyKS
Gn+T2uCyOP4a1DTUoPyoNJXo
-----END PRIVATE KEY-----`,
		},
		{
			Cert: `-----BEGIN CERTIFICATE-----
MIIDEzCCAfugAwIBAgIUOwYJvXJ+s0qX9uAKFjW0zExV51IwDQYJKoZIhvcNAQEF
BQAwJTEKMAgGA1UEAwwBIzEKMAgGA1UECgwBIzELMAkGA1UEBhMCQ04wHhcNMjEw
NDMwMDAwMDAwWhcNMjEwNDMwMDAwMDAwWjBDMQswCQYDVQQGEwJDTjEKMAgGA1UE
CAwBIzEKMAgGA1UEBwwBIzEKMAgGA1UECgwBIzEQMA4GA1UEAwwHZm9vLmNvbTCC
ASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAPRFV285ScP1ntF3zlj60GC0
DJyCEX6Ji38gBf+6An6Zk7D+3Aif9C/3e7V0811x0VoO4o9ZdQUSdxmE9fj/ADOU
OM/AYf62L51d/zdqXVaF89vpqPk8em4179wo6jg2IiCewGVLTtuMa/5Mud8XZOly
tcOXS7ZnCbfm/XklwGL1rAmWhOTSDHlIH5bbC46tmi3E9Cjp+VTiwzVCgVtrLkzY
0cjs72m2wb5uZ9TlT7n1TKYjdX74FvYp4X70YEcFYEUmFMxMV7otkJ7wTWWVNah/
ZsojaiJ48ueJFQR1S9utYA/h6LcA4T6UQJxw7+6SjJElLCHGht5UHFvQkjQvxZkC
AwEAAaMdMBswCwYDVR0RBAQwAoIAMAwGA1UdEwEB/wQCMAAwDQYJKoZIhvcNAQEF
BQADggEBAHE4U9SlCIVNjpfOyfH0NPhxLTAqH83GQKJc7TgQFmhby1dfQE7MOTaN
ayA1RJ0qKcNGlHP70M/Xc8TIF+E7pOASqa+zNztiv14zHIgJC9oGJcwt1sh8GADz
4EJSQ1mIRxbgs39BA9FDY91HBa3RfLxkmyTbQK1rhKdh8aBYr0/6R1oAdKEQF/vQ
HxD4NCpJruxp7g+RSet1PB12GOao1Ntfb7kOLAHzYW3yvTsCaQ7EdeueOs8dv+G/
Ncy+4n/l3audbi+WQFfEvb1bwyADPpp90C9OczHzpR4+dtuR4oUXB6ZXimB3MljC
BhoOkUOMjrKl/QDkB5pxa/IxURffFDs=
-----END CERTIFICATE-----`,
			Key: `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQD0RVdvOUnD9Z7R
d85Y+tBgtAycghF+iYt/IAX/ugJ+mZOw/twIn/Qv93u1dPNdcdFaDuKPWXUFEncZ
hPX4/wAzlDjPwGH+ti+dXf83al1WhfPb6aj5PHpuNe/cKOo4NiIgnsBlS07bjGv+
TLnfF2TpcrXDl0u2Zwm35v15JcBi9awJloTk0gx5SB+W2wuOrZotxPQo6flU4sM1
QoFbay5M2NHI7O9ptsG+bmfU5U+59UymI3V++Bb2KeF+9GBHBWBFJhTMTFe6LZCe
8E1llTWof2bKI2oiePLniRUEdUvbrWAP4ei3AOE+lECccO/ukoyRJSwhxobeVBxb
0JI0L8WZAgMBAAECggEAS3gBA4QXnkuMvfrd7e/P4ZC/NLua3BVy29vw/olhq+uX
DeAv6xpAP3Ge7vkrF3vKyqA+rztcRCzoTyIdPMjMLyNkSguOOzveiig4ac6a99h6
9z7Bzf85dEOFz+d0NdnDwYBKwQD7ZCmGVBMwehSoQNgIAF4GLly3S/I57ewT/H6A
GjknY/jCmk+L9388hjcL0jrEJR/br2O/o6f1zdRYWqqb9A1wDW4zkiu4Wrq5in9s
cxQ/7667eGkciD3HkJvwcbi7xg9ZJHxCWScVYGRBX9ek6fMKxML5hUsITZjt5zxF
p+HmOKJcii6hlR1RWaUsbrpQOHVui3US7CAJNR/gAQKBgQD7zE0z2XGv81XjIvMk
sS+IvtsSGpvoUI2QRbdnelC8ahCdKj5PmVQyfhPxrgNCTu6k7VzQUDL+wZDGKoRL
NEaRkoHz7tVzBE7DY7Y2SD0yfjT477w98iaF/nwortmhpXms0KyzPhZOF5d5166q
PDR31NIFvmy2H+Hh9bVM5BYaGQKBgQD4WOIgocc+pXn+3fehNT3qedrvGuXGYX7I
PAO/4zM/oP/0TtxKTz5wGAFR9heBKfogW2jYUBBOofraLMJq3X+T9jEOXuQ9+UQq
HaybHdQycxpTIWhtiAs9khvSbuEBs2SXyKussPGW8Do5uVfi4/KWWu/wcTzMlfEv
w207iaN3gQKBgQDAh7u0XJx4PCi871lZAf5logGiOyRhI07LNPOCtN0M5FDly4ov
lP7zSMH5NuQZDH+fLjucsOX9M4Z+b74OPt+CqbKiEUm2k2GiNxj5Mo1QkX3xpmWa
PBDGvgqzlNalqgB6amjS+TNW7OUO7iMI2dYIlnsslylKrOArxZOmQnS/6QKBgQCs
ZVcj++nKDSjwybk6yTDf8hMO5IcY/Vj7Ot4HeHp88xB60buOQhA/1AomkUSjvzYI
/Ct97aZET6FJjsSvVm9XkRFgvnKGquCss8i8LSq+krR1fL13O3dCGIkDvUCo45Uy
4HR7/qDWfJCOvaDKuh4OTbY+HP1tr7CrzWeoatV1AQKBgFFmlMWrIThfjtjVWDTg
+QPLQTofB1A3lrCmB52iBdUMi0qGnExLn8aiy54wPz/I7rEplsLzg2hmDNuBPM7q
QLAtVaZd9SSi4Z/RX6B4L3Rj0Mwfn+tbrtYO5Pyhi40hiXf4aMgbVDFYMR0MMmH0
4uiYeQPmK6USKjntOFQ0eNOe
-----END PRIVATE KEY-----`,
		},
	}
)

var (
	testIngressHTTPSNamespace         = "ingress-https"
	testIngressHTTPSRedirectNamespace = "ingress-redirect"
)

func TestHTTPSRedirect(t *testing.T) {
	ctx := context.Background()
	opts := metav1.CreateOptions{}

	t.Logf("creating namespace %s for testing", testIngressHTTPSRedirectNamespace)
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: testIngressHTTPSRedirectNamespace}}
	ns, err := env.Cluster().Client().CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up namespace %s", testIngressHTTPSRedirectNamespace)
		require.NoError(t, env.Cluster().Client().CoreV1().Namespaces().Delete(ctx, ns.Name, metav1.DeleteOptions{}))
		require.Eventually(t, func() bool {
			_, err := env.Cluster().Client().CoreV1().Namespaces().Get(ctx, ns.Name, metav1.GetOptions{})
			if err != nil {
				if errors.IsNotFound(err) {
					return true
				}
			}
			return false
		}, ingressWait, waitTick)
	}()

	t.Log("creating an HTTP container via deployment to test redirect functionality")
	container := generators.NewContainer("alsohttpbin", httpBinImage, 80)
	deployment := generators.NewDeploymentForContainer(container)
	_, err = env.Cluster().Client().AppsV1().Deployments(corev1.NamespaceDefault).Create(ctx, deployment, opts)
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up deployment %s", deployment.Name)
		assert.NoError(t, env.Cluster().Client().AppsV1().Deployments(corev1.NamespaceDefault).Delete(ctx, deployment.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing deployment %s via Service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeClusterIP)
	service, err = env.Cluster().Client().CoreV1().Services(corev1.NamespaceDefault).Create(ctx, service, opts)
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up Service %s", service.Name)
		assert.NoError(t, env.Cluster().Client().CoreV1().Services(corev1.NamespaceDefault).Delete(ctx, service.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing Service %s via Ingress", service.Name)
	ingress := generators.NewIngressForService("/httpbin", map[string]string{
		annotations.IngressClassKey:             ingressClass,
		"konghq.com/protocols":                  "https",
		"konghq.com/https-redirect-status-code": "301",
	}, service)
	ingress, err = env.Cluster().Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Create(ctx, ingress, opts)
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up Ingress %s", ingress.Name)
		assert.NoError(t, env.Cluster().Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Delete(ctx, ingress.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("waiting for Ingress %s to be operational and properly redirect", ingress.Name)
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: time.Second * 3,
	}
	assert.Eventually(t, func() bool {
		resp, err := client.Get(fmt.Sprintf("%s/httpbin", proxyURL))
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusMovedPermanently
	}, ingressWait, waitTick)
}

func TestHTTPSIngress(t *testing.T) {
	ctx := context.Background()
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	testTransport := http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			if addr == "foo.example:443" {
				addr = fmt.Sprintf("%s:443", proxyURL.Hostname())
			}
			if addr == "bar.example:443" {
				addr = fmt.Sprintf("%s:443", proxyURL.Hostname())
			}
			if addr == "baz.example:443" {
				addr = fmt.Sprintf("%s:443", proxyURL.Hostname())
			}
			return dialer.DialContext(ctx, network, addr)
		},
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec
	}
	httpcStatic := http.Client{
		Timeout:   httpcTimeout,
		Transport: &testTransport,
	}

	t.Logf("creating namespace %s for testing", testIngressHTTPSNamespace)
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: testIngressHTTPSNamespace}}
	ns, err := env.Cluster().Client().CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up namespace %s", testIngressHTTPSNamespace)
		require.NoError(t, env.Cluster().Client().CoreV1().Namespaces().Delete(ctx, ns.Name, metav1.DeleteOptions{}))
		require.Eventually(t, func() bool {
			_, err := env.Cluster().Client().CoreV1().Namespaces().Get(ctx, ns.Name, metav1.GetOptions{})
			if err != nil {
				if errors.IsNotFound(err) {
					return true
				}
			}
			return false
		}, ingressWait, waitTick)
	}()

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", httpBinImage, 80)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err = env.Cluster().Client().AppsV1().Deployments(corev1.NamespaceDefault).Create(ctx, deployment, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the deployment %s", deployment.Name)
		assert.NoError(t, env.Cluster().Client().AppsV1().Deployments(corev1.NamespaceDefault).Delete(ctx, deployment.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(corev1.NamespaceDefault).Create(ctx, service, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the service %s", service.Name)
		assert.NoError(t, env.Cluster().Client().CoreV1().Services(corev1.NamespaceDefault).Delete(ctx, service.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("creating an ingress for service %s with ingress.class %s", service.Name, ingressClass)
	ingress1 := generators.NewIngressForService("/foo", map[string]string{
		annotations.IngressClassKey: ingressClass,
		"konghq.com/strip-path":     "true",
	}, service)
	ingress1.Spec.TLS = []networkingv1.IngressTLS{
		{
			SecretName: "secret1",
			Hosts:      []string{"foo.example"},
		},
	}
	ingress1.ObjectMeta.Name = "ingress1"
	ingress2 := generators.NewIngressForService("/bar", map[string]string{
		annotations.IngressClassKey: ingressClass,
		"konghq.com/strip-path":     "true",
	}, service)
	ingress2.Spec.TLS = []networkingv1.IngressTLS{
		{
			SecretName: "secret2",
			Hosts:      []string{"bar.example"},
		},
	}
	ingress2.ObjectMeta.Name = "ingress2"
	secrets := []*corev1.Secret{
		{
			ObjectMeta: metav1.ObjectMeta{
				UID:       types.UID("7428fb98-180b-4702-a91f-61351a33c6e4"),
				Name:      "secret1",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"tls.crt": []byte(tlsPairs[0].Cert),
				"tls.key": []byte(tlsPairs[0].Key),
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				UID:       types.UID("7428fb98-180b-4702-a91f-61351a33c6e5"),
				Name:      "secret2",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"tls.crt": []byte(tlsPairs[1].Cert),
				"tls.key": []byte(tlsPairs[1].Key),
			},
		},
	}
	secret1, err := env.Cluster().Client().CoreV1().Secrets(corev1.NamespaceDefault).Create(ctx, secrets[0], metav1.CreateOptions{})
	assert.NoError(t, err)
	secret2, err := env.Cluster().Client().CoreV1().Secrets(corev1.NamespaceDefault).Create(ctx, secrets[1], metav1.CreateOptions{})
	assert.NoError(t, err)
	ingress1, err = env.Cluster().Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Create(ctx, ingress1, metav1.CreateOptions{})
	assert.NoError(t, err)
	ingress2, err = env.Cluster().Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Create(ctx, ingress2, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("ensuring that Ingress %s is cleaned up", ingress1.Name)
		if err := env.Cluster().Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Delete(ctx, ingress1.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}
		}
		t.Logf("ensuring that Ingress %s is cleaned up", ingress2.Name)
		if err := env.Cluster().Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Delete(ctx, ingress2.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}
		}
		t.Logf("ensuring that Secret %s is cleaned up", secret1.Name)
		if err := env.Cluster().Client().CoreV1().Secrets(corev1.NamespaceDefault).Delete(ctx, secret1.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}
		}
		t.Logf("ensuring that Secret %s is cleaned up", secret2.Name)
		if err := env.Cluster().Client().CoreV1().Secrets(corev1.NamespaceDefault).Delete(ctx, secret2.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}
		}
	}()

	t.Logf("checking ingress %s status readiness.", ingress1.Name)
	require.Eventually(t, func() bool {
		curIng, err := env.Cluster().Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Get(ctx, ingress1.Name, metav1.GetOptions{})
		if err != nil || curIng == nil {
			return false
		}
		ingresses := curIng.Status.LoadBalancer.Ingress
		for _, ingress := range ingresses {
			if len(ingress.Hostname) > 0 || len(ingress.IP) > 0 {
				t.Logf("networkingv1 ingress1 hostname %s or ip %s is ready to redirect traffic.", ingress.Hostname, ingress.IP)
				return true
			}
		}
		return false
	}, 120*time.Second, 1*time.Second, true)

	t.Logf("checking ingress %s status readiness.", ingress2.Name)
	assert.Eventually(t, func() bool {
		curIng, err := env.Cluster().Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Get(ctx, ingress2.Name, metav1.GetOptions{})
		if err != nil || curIng == nil {
			return false
		}
		ingresses := curIng.Status.LoadBalancer.Ingress
		for _, ingress := range ingresses {
			if len(ingress.Hostname) > 0 || len(ingress.IP) > 0 {
				t.Logf("networkingv1 ingress2 hostname %s or ip %s is ready to redirect traffic.", ingress.Hostname, ingress.IP)
				return true
			}
		}
		return false
	}, 120*time.Second, 1*time.Second, true)

	t.Logf("waiting for routes from Ingress %s to be operational with expected certificate", ingress1.Name)
	assert.Eventually(t, func() bool {
		resp, err := httpcStatic.Get("https://foo.example:443/foo")
		if err != nil {
			t.Logf("WARNING: error while waiting for https://foo.example:443/foo: %v", err)
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			b := new(bytes.Buffer)
			n, err := b.ReadFrom(resp.Body)
			require.NoError(t, err)
			require.True(t, n > 0)
			return strings.Contains(b.String(), "<title>httpbin.org</title>") && resp.TLS.PeerCertificates[0].Subject.CommonName == "secure-foo-bar"
		}
		return false
	}, ingressWait, waitTick, true)

	t.Logf("waiting for routes from Ingress %s to be operational with expected certificate", ingress2.Name)
	assert.Eventually(t, func() bool {
		resp, err := httpcStatic.Get("https://bar.example:443/bar")
		if err != nil {
			t.Logf("WARNING: error while waiting for https://bar.example:443/bar: %v", err)
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			b := new(bytes.Buffer)
			n, err := b.ReadFrom(resp.Body)
			require.NoError(t, err)
			require.True(t, n > 0)
			return strings.Contains(b.String(), "<title>httpbin.org</title>") && resp.TLS.PeerCertificates[0].Subject.CommonName == "foo.com"
		}
		return false
	}, ingressWait, waitTick, true)

	// This should work currently. generators.NewIngressForService() only creates path rules by default, so while we don't
	// do anything for baz.example other than add fake DNS for it, the /bar still routes it through ingress2's route.
	// We're going to break it later, but need to confirm it does work first.
	t.Logf("confirm Ingress %s path routes available on other hostnames", ingress2.Name)
	assert.Eventually(t, func() bool {
		resp, err := httpcStatic.Get("https://baz.example:443/bar")
		if err != nil {
			t.Logf("WARNING: error while waiting for https://bar.example:443/baz: %v", err)
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			b := new(bytes.Buffer)
			n, err := b.ReadFrom(resp.Body)
			require.NoError(t, err)
			require.True(t, n > 0)
			return strings.Contains(b.String(), "<title>httpbin.org</title>")
		}
		return false
	}, ingressWait, waitTick)

	ingress2, err = env.Cluster().Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Get(ctx, ingress2.Name, metav1.GetOptions{})
	assert.NoError(t, err)
	ingress2.ObjectMeta.Annotations["konghq.com/snis"] = "bar.example"
	ingress2, err = env.Cluster().Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Update(ctx, ingress2, metav1.UpdateOptions{})
	assert.NoError(t, err)

	t.Logf("confirm Ingress %s no longer routes without matching SNI", ingress2.Name)
	assert.Eventually(t, func() bool {
		resp, err := httpcStatic.Get("https://baz.example:443/bar")
		if err != nil {
			t.Logf("WARNING: error while waiting for https://baz.example:443/bar: %v", err)
			return false
		}

		defer resp.Body.Close()
		return resp.StatusCode == http.StatusNotFound
	}, ingressWait, waitTick)

	t.Logf("confirm Ingress %s still routes with matching SNI", ingress2.Name)
	assert.Eventually(t, func() bool {
		resp, err := httpcStatic.Get("https://bar.example:443/bar")
		if err != nil {
			t.Logf("WARNING: error while waiting for https://bar.example:443/bar: %v", err)
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			b := new(bytes.Buffer)
			n, err := b.ReadFrom(resp.Body)
			require.NoError(t, err)
			require.True(t, n > 0)
			return strings.Contains(b.String(), "<title>httpbin.org</title>")
		}
		return false
	}, ingressWait, waitTick)
}
