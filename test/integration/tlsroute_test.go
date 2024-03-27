//go:build integration_tests

package integration

import (
	"bytes"
	"context"
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
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/builder"
	"github.com/kong/kubernetes-ingress-controller/v2/test"
	"github.com/kong/kubernetes-ingress-controller/v2/test/helpers/certificate"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
)

const (
	tlsRouteHostname      = "tlsroute.kong.example"
	tlsRouteExtraHostname = "extratlsroute.kong.example"
	tlsSecretName         = "secret-test"
)

const (
	tlsEchoPort = 1030
)

// TestTLSRouteReferenceGrant tests cross-namespace certificate references. These are technically implemented within
// Gateway Listeners, but require an attached Route to see the associated certificate behavior on the proxy.
func TestTLSRoutePassthroughReferenceGrant(t *testing.T) {
	skipTestForExpressionRouter(t)
	t.Log("locking Gateway TLS ports")
	tlsMutex.Lock()
	t.Cleanup(func() {
		t.Log("unlocking TLS port")
		tlsMutex.Unlock()
	})

	ctx := context.Background()
	ns, cleaner := helpers.Setup(ctx, t, env)

	otherNs, err := clusters.GenerateNamespace(ctx, env.Cluster(), t.Name())
	require.NoError(t, err)
	cleaner.AddNamespace(otherNs)

	t.Log("getting the gateway client")
	gatewayClient, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("configuring secrets")
	tlsRouteExampleTLSCert, tlsRouteExampleTLSKey := certificate.MustGenerateSelfSignedCertPEMFormat(certificate.WithCommonName(tlsRouteHostname))
	extraTLSRouteTLSCert, extraTLSRouteTLSKey := certificate.MustGenerateSelfSignedCertPEMFormat(certificate.WithCommonName(tlsRouteExtraHostname))

	secrets := []*corev1.Secret{
		{
			ObjectMeta: metav1.ObjectMeta{
				UID:       k8stypes.UID("7428fb98-180b-4702-a91f-61351a33c6e8"),
				Name:      tlsSecretName,
				Namespace: ns.Name,
			},
			Data: map[string][]byte{
				"tls.crt": tlsRouteExampleTLSCert,
				"tls.key": tlsRouteExampleTLSKey,
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				UID:  k8stypes.UID("7428fb98-180b-4702-a91f-61351a33c6e9"),
				Name: "secret2",
			},
			Data: map[string][]byte{
				"tls.crt": extraTLSRouteTLSCert,
				"tls.key": extraTLSRouteTLSKey,
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
	// we need to create the secret 2 in the namespace 1 as well because we need to mount in the deployment. The Gateway will be
	// using the secret in namespace 1 to test the referenceGrant.
	secret3, err := env.Cluster().Client().CoreV1().Secrets(ns.Name).Create(ctx, secrets[1], metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(secret3)

	modePassthrough := gatewayv1.TLSModePassthrough
	t.Log("deploying a gateway to the test cluster using unmanaged gateway mode")
	gateway, err := DeployGateway(ctx, gatewayClient, ns.Name, unmanagedGatewayClassName, func(gw *gatewayv1.Gateway) {
		otherNamespace := gatewayv1.Namespace(otherNs.Name)
		gw.Spec.Listeners = []gatewayv1.Listener{
			builder.NewListener("tls").
				TLS().
				WithPort(ktfkong.DefaultTLSServicePort).
				WithHostname(tlsRouteHostname).
				WithTLSConfig(&gatewayv1.GatewayTLSConfig{
					Mode: &modePassthrough,
					CertificateRefs: []gatewayv1.SecretObjectReference{
						{
							Name: gatewayv1.ObjectName(secrets[0].Name),
						},
					},
				}).Build(),
			builder.NewListener("tlsother").
				TLS().
				WithPort(ktfkong.DefaultTLSServicePort).
				WithHostname(tlsRouteExtraHostname).
				WithTLSConfig(&gatewayv1.GatewayTLSConfig{
					Mode: &modePassthrough,
					CertificateRefs: []gatewayv1.SecretObjectReference{
						{
							Name:      gatewayv1.ObjectName(secrets[1].Name),
							Namespace: &otherNamespace,
						},
					},
				}).Build(),
		}
	})

	require.NoError(t, err)
	cleaner.Add(gateway)

	secret2Name := gatewayv1alpha2.ObjectName(secrets[1].Name)
	t.Logf("creating a ReferenceGrant that permits gateway access from %s to secrets in %s", ns.Name, otherNs.Name)
	grant := &gatewayv1beta1.ReferenceGrant{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
		},
		Spec: gatewayv1beta1.ReferenceGrantSpec{
			From: []gatewayv1beta1.ReferenceGrantFrom{
				{
					Group:     gatewayv1alpha2.Group("gateway.networking.k8s.io"),
					Kind:      gatewayv1alpha2.Kind("Gateway"),
					Namespace: gatewayv1alpha2.Namespace(gateway.Namespace),
				},
			},
			To: []gatewayv1beta1.ReferenceGrantTo{
				{
					Group: gatewayv1alpha2.Group(""),
					Kind:  gatewayv1alpha2.Kind("Secret"),
					Name:  &secret2Name,
				},
			},
		},
	}

	grant, err = gatewayClient.GatewayV1beta1().ReferenceGrants(otherNs.Name).Create(ctx, grant, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(grant)

	t.Log("creating a tcpecho pod to test TLSRoute traffic routing")
	testUUID := uuid.NewString()
	deployment := generators.NewDeploymentForContainer(createTLSEchoContainer(tlsEchoPort, testUUID))
	deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, corev1.Volume{
		Name: tlsSecretName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: tlsSecretName,
			},
		},
	})
	deployment, err = env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(deployment)

	t.Log("creating another tcpecho pod to test TLSRoute traffic routing with referenceGrant")
	testUUID2 := uuid.NewString()
	deployment2 := generators.NewDeploymentForContainer(createTLSEchoContainer(tlsEchoPort, testUUID2))
	deployment2.Spec.Template.Spec.Volumes = append(deployment2.Spec.Template.Spec.Volumes, corev1.Volume{
		Name: tlsSecretName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: string(secret2Name),
			},
		},
	})
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

	backendTLSPort := gatewayv1alpha2.PortNumber(tlsEchoPort)
	t.Logf("creating a tlsroute to access deployment %s via kong", deployment.Name)
	tlsroute := &gatewayv1alpha2.TLSRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
		},
		Spec: gatewayv1alpha2.TLSRouteSpec{
			CommonRouteSpec: gatewayv1alpha2.CommonRouteSpec{
				ParentRefs: []gatewayv1alpha2.ParentReference{{
					Name: gatewayv1alpha2.ObjectName(gateway.Name),
				}},
			},
			Hostnames: []gatewayv1alpha2.Hostname{tlsRouteHostname, tlsRouteExtraHostname},
			Rules: []gatewayv1alpha2.TLSRouteRule{{
				BackendRefs: []gatewayv1alpha2.BackendRef{
					{
						BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
							Name: gatewayv1alpha2.ObjectName(service.Name),
							Port: &backendTLSPort,
						},
					},
					{
						BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
							Name: gatewayv1alpha2.ObjectName(service2.Name),
							Port: &backendTLSPort,
						},
					},
				},
			}},
		},
	}
	tlsroute, err = gatewayClient.GatewayV1alpha2().TLSRoutes(ns.Name).Create(ctx, tlsroute, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(tlsroute)

	proxyAddress := fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTLSServicePort)
	t.Log("verifying that the tcpecho is responding properly over TLS")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(proxyAddress, testUUID, tlsRouteHostname, tlsRouteHostname, true)
		if err != nil {
			t.Logf("failed accessing tcpecho at %s, err: %v", proxyAddress, err)
			return false
		}
		return err == nil && responded == true
	}, ingressWait, waitTick)

	t.Log("verifying that the tcpecho route can also serve certificates permitted by a ReferenceGrant with a named To")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(proxyAddress, testUUID2, tlsRouteExtraHostname, tlsRouteExtraHostname, true)
		if err != nil {
			t.Logf("failed accessing tcpecho at %s, err: %v", proxyAddress, err)
			return false
		}
		return err == nil && responded == true
	}, ingressWait, waitTick)

	t.Log("verifying that using the wrong name in the ReferenceGrant removes the related certificate")
	badName := gatewayv1alpha2.ObjectName("garbage")
	grant.Spec.To[0].Name = &badName
	grant, err = gatewayClient.GatewayV1beta1().ReferenceGrants(otherNs.Name).Update(ctx, grant, metav1.UpdateOptions{})
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(proxyAddress, testUUID2, tlsRouteExtraHostname, tlsRouteExtraHostname, true)
		return err != nil && responded == false
	}, ingressWait, waitTick)

	t.Log("verifying that a Listener has the invalid ref status condition")
	gateway, err = gatewayClient.GatewayV1beta1().Gateways(ns.Name).Get(ctx, gateway.Name, metav1.GetOptions{})
	require.NoError(t, err)
	invalid := false
	for _, status := range gateway.Status.Listeners {
		if ok := util.CheckCondition(
			status.Conditions,
			util.ConditionType(gatewayv1.ListenerConditionResolvedRefs),
			util.ConditionReason(gatewayv1.ListenerReasonRefNotPermitted),
			metav1.ConditionFalse,
			gateway.Generation,
		); ok {
			invalid = true
		}
	}
	require.True(t, invalid)

	t.Log("verifying the certificate returns when using a ReferenceGrant with no name restrictions")
	grant.Spec.To[0].Name = nil
	_, err = gatewayClient.GatewayV1beta1().ReferenceGrants(otherNs.Name).Update(ctx, grant, metav1.UpdateOptions{})
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(proxyAddress, testUUID2, tlsRouteExtraHostname, tlsRouteExtraHostname, true)
		if err != nil {
			t.Logf("failed accessing tcpecho at %s, err: %v", proxyAddress, err)
			return false
		}
		return err == nil && responded == true
	}, ingressWait, waitTick)
}

func TestTLSRoutePassthrough(t *testing.T) {
	skipTestForExpressionRouter(t)
	t.Log("locking Gateway TLS ports")
	tlsMutex.Lock()
	t.Cleanup(func() {
		t.Log("unlocking TLS port")
		tlsMutex.Unlock()
	})

	ctx := context.Background()
	ns, cleaner := helpers.Setup(ctx, t, env)

	t.Log("getting gateway client")
	gatewayClient, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("configuring secrets")
	tlsRouteExampleTLSCert, tlsRouteExampleTLSKey := certificate.MustGenerateSelfSignedCertPEMFormat(certificate.WithCommonName(tlsRouteHostname))
	secrets := []*corev1.Secret{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      tlsSecretName,
				Namespace: ns.Name,
			},
			Data: map[string][]byte{
				"tls.crt": tlsRouteExampleTLSCert,
				"tls.key": tlsRouteExampleTLSKey,
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

	t.Log("deploying a gateway to the test cluster using unmanaged gateway mode")
	modePassthrough := gatewayv1.TLSModePassthrough
	gatewayName := uuid.NewString()
	gateway, err := DeployGateway(ctx, gatewayClient, ns.Name, gatewayClassName, func(gw *gatewayv1.Gateway) {
		hostname := gatewayv1.Hostname(tlsRouteHostname)
		gw.Name = gatewayName
		gw.Spec.Listeners = []gatewayv1.Listener{
			{
				Name:     "tls-passthrough",
				Protocol: gatewayv1.TLSProtocolType,
				Port:     gatewayv1.PortNumber(ktfkong.DefaultTLSServicePort),
				Hostname: &hostname,
				TLS: &gatewayv1.GatewayTLSConfig{
					CertificateRefs: []gatewayv1.SecretObjectReference{
						{
							Name: gatewayv1.ObjectName(tlsSecretName),
						},
					},
					Mode: &modePassthrough,
				},
			},
		}
	})
	require.NoError(t, err)
	cleaner.Add(gateway)

	t.Log("creating a tcpecho deployment to test TLSRoute traffic routing")
	testUUID := uuid.NewString() // go-echo sends a "Running on Pod <UUID>." immediately on connecting
	deployment := generators.NewDeploymentForContainer(createTLSEchoContainer(tlsEchoPort, testUUID))
	deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, corev1.Volume{
		Name: tlsSecretName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: tlsSecretName,
			},
		},
	})
	deployment, err = env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(deployment)

	t.Log("creating an additional tcpecho pod to test TLSRoute multiple backendRef loadbalancing")
	testUUID2 := uuid.NewString() // go-echo sends a "Running on Pod <UUID>." immediately on connecting
	deployment2 := generators.NewDeploymentForContainer(createTLSEchoContainer(tlsEchoPort, testUUID2))
	deployment2.Spec.Template.Spec.Volumes = append(deployment2.Spec.Template.Spec.Volumes, corev1.Volume{
		Name: tlsSecretName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: tlsSecretName,
			},
		},
	})
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

	backendTLSPort := gatewayv1alpha2.PortNumber(tlsEchoPort)
	t.Logf("create a TLSRoute using passthrough listener")
	sectionName := gatewayv1alpha2.SectionName("tls-passthrough")
	tlsRoute := &gatewayv1alpha2.TLSRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
		},
		Spec: gatewayv1alpha2.TLSRouteSpec{
			CommonRouteSpec: gatewayv1alpha2.CommonRouteSpec{
				ParentRefs: []gatewayv1alpha2.ParentReference{{
					Name:        gatewayv1alpha2.ObjectName(gateway.Name),
					SectionName: &sectionName,
				}},
			},
			Hostnames: []gatewayv1alpha2.Hostname{tlsRouteHostname},
			Rules: []gatewayv1alpha2.TLSRouteRule{{
				BackendRefs: []gatewayv1alpha2.BackendRef{{
					BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
						Name: gatewayv1alpha2.ObjectName(service.Name),
						Port: &backendTLSPort,
					},
				}},
			}},
		},
	}
	tlsRoute, err = gatewayClient.GatewayV1alpha2().TLSRoutes(ns.Name).Create(ctx, tlsRoute, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(tlsRoute)

	proxyAddress := fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTLSServicePort)
	t.Log("verifying that the tcpecho is responding properly over TLS")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(proxyAddress, testUUID, tlsRouteHostname, tlsRouteHostname, true)
		if err != nil {
			t.Logf("failed accessing tcpecho at %s, err: %v", proxyAddress, err)
			return false
		}
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
	callback := GetGatewayIsUnlinkedCallback(ctx, t, gatewayClient, gatewayv1.TLSProtocolType, ns.Name, tlsRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that the tcpecho is no longer responding")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTLSServicePort),
			testUUID, tlsRouteHostname, tlsRouteHostname, false)
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
	callback = GetGatewayIsLinkedCallback(ctx, t, gatewayClient, gatewayv1.TLSProtocolType, ns.Name, tlsRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that putting the parentRefs back results in the routes becoming available again")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(proxyAddress, testUUID, tlsRouteHostname, tlsRouteHostname, true)
		if err != nil {
			t.Logf("failed accessing tcpecho at %s, err: %v", proxyAddress, err)
			return false
		}
		return err == nil && responded == true
	}, ingressWait, waitTick)

	t.Log("deleting the GatewayClass")
	require.NoError(t, gatewayClient.GatewayV1beta1().GatewayClasses().Delete(ctx, gwc.Name, metav1.DeleteOptions{}))

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	callback = GetGatewayIsUnlinkedCallback(ctx, t, gatewayClient, gatewayv1.TLSProtocolType, ns.Name, tlsRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that the data-plane configuration from the TLSRoute gets dropped with the GatewayClass now removed")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(proxyAddress, testUUID, tlsRouteHostname, tlsRouteHostname, true)
		return responded == false && errors.Is(err, io.EOF)
	}, ingressWait, waitTick)

	t.Log("putting the GatewayClass back")
	gwc, err = DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
	require.NoError(t, err)

	t.Log("verifying that the Gateway gets linked to the route via status")
	callback = GetGatewayIsLinkedCallback(ctx, t, gatewayClient, gatewayv1.TLSProtocolType, ns.Name, tlsRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that creating the GatewayClass again triggers reconciliation of TLSRoutes and the route becomes available again")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(proxyAddress, testUUID, tlsRouteHostname, tlsRouteHostname, true)
		if err != nil {
			t.Logf("failed accessing tcpecho at %s, err: %v", proxyAddress, err)
			return false
		}
		return err == nil && responded == true
	}, ingressWait, waitTick)

	t.Log("deleting the Gateway")
	require.NoError(t, gatewayClient.GatewayV1beta1().Gateways(ns.Name).Delete(ctx, gatewayName, metav1.DeleteOptions{}))

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	callback = GetGatewayIsUnlinkedCallback(ctx, t, gatewayClient, gatewayv1.TLSProtocolType, ns.Name, tlsRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that the data-plane configuration from the TLSRoute gets dropped with the Gateway now removed")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(proxyAddress, testUUID, tlsRouteHostname, tlsRouteHostname, true)
		return responded == false && errors.Is(err, io.EOF)
	}, ingressWait, waitTick)

	t.Log("putting the Gateway back")
	gateway, err = DeployGateway(ctx, gatewayClient, ns.Name, gatewayClassName, func(gw *gatewayv1.Gateway) {
		hostname := gatewayv1.Hostname(tlsRouteHostname)
		gw.Name = gatewayName
		gw.Spec.Listeners = []gatewayv1.Listener{
			{
				Name:     "tls-passthrough",
				Protocol: gatewayv1.TLSProtocolType,
				Port:     gatewayv1.PortNumber(ktfkong.DefaultTLSServicePort),
				Hostname: &hostname,
				TLS: &gatewayv1.GatewayTLSConfig{
					CertificateRefs: []gatewayv1.SecretObjectReference{
						{
							Name: gatewayv1.ObjectName(tlsSecretName),
						},
					},
					Mode: &modePassthrough,
				},
			},
		}
	})
	require.NoError(t, err)

	t.Log("verifying that the Gateway gets linked to the route via status")
	callback = GetGatewayIsLinkedCallback(ctx, t, gatewayClient, gatewayv1.TLSProtocolType, ns.Name, tlsRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that creating the Gateway again triggers reconciliation of TLSRoutes and the route becomes available again")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(proxyAddress, testUUID, tlsRouteHostname, tlsRouteHostname, true)
		if err != nil {
			t.Logf("failed accessing tcpecho at %s, err: %v", proxyAddress, err)
			return false
		}
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
					Port: &backendTLSPort,
				},
			},
			{
				BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
					Name: gatewayv1alpha2.ObjectName(service2.Name),
					Port: &backendTLSPort,
				},
			},
		}

		tlsRoute, err = gatewayClient.GatewayV1alpha2().TLSRoutes(ns.Name).Update(ctx, tlsRoute, metav1.UpdateOptions{})
		return err == nil
	}, ingressWait, waitTick)

	t.Log("verifying that the TLSRoute is now load-balanced between two services")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(proxyAddress, testUUID, tlsRouteHostname, tlsRouteHostname, true)
		if err != nil {
			t.Logf("failed accessing tcpecho at %s, err: %v", proxyAddress, err)
			return false
		}
		return err == nil && responded == true
	}, ingressWait, waitTick)
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(proxyAddress, testUUID2, tlsRouteHostname, tlsRouteHostname, true)
		if err != nil {
			t.Logf("failed accessing tcpecho at %s, err: %v", proxyAddress, err)
			return false
		}
		return err == nil && responded == true
	}, ingressWait, waitTick)

	t.Log("deleting both GatewayClass and Gateway rapidly")
	require.NoError(t, gatewayClient.GatewayV1beta1().GatewayClasses().Delete(ctx, gwc.Name, metav1.DeleteOptions{}))
	require.NoError(t, gatewayClient.GatewayV1beta1().Gateways(ns.Name).Delete(ctx, gateway.Name, metav1.DeleteOptions{}))

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	callback = GetGatewayIsUnlinkedCallback(ctx, t, gatewayClient, gatewayv1.TLSProtocolType, ns.Name, tlsRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that the data-plane configuration from the TLSRoute does not get orphaned with the GatewayClass and Gateway gone")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTLSServicePort),
			testUUID, tlsRouteHostname, tlsRouteHostname, true)
		return responded == false && errors.Is(err, io.EOF)
	}, ingressWait, waitTick)

	t.Log("testing port matching")
	t.Log("putting the Gateway back")
	_, err = DeployGateway(ctx, gatewayClient, ns.Name, gatewayClassName, func(gw *gatewayv1.Gateway) {
		hostname := gatewayv1.Hostname(tlsRouteHostname)
		gw.Name = gatewayName
		gw.Spec.Listeners = []gatewayv1.Listener{
			{
				Name:     "tls-passthrough",
				Protocol: gatewayv1.TLSProtocolType,
				Port:     gatewayv1.PortNumber(ktfkong.DefaultTLSServicePort),
				Hostname: &hostname,
				TLS: &gatewayv1.GatewayTLSConfig{
					CertificateRefs: []gatewayv1.SecretObjectReference{
						{
							Name: gatewayv1.ObjectName(tlsSecretName),
						},
					},
					Mode: &modePassthrough,
				},
			},
		}
	})
	require.NoError(t, err)

	t.Log("putting the GatewayClass back")
	_, err = DeployGatewayClass(ctx, gatewayClient, gatewayClassName)

	t.Log("ensuring tls echo responds after recreating gateway and gateway class")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(proxyAddress, testUUID, tlsRouteHostname, tlsRouteHostname, true)
		if err != nil {
			t.Logf("failed accessing tcpecho at %s, err: %v", proxyAddress, err)
			return false
		}
		return err == nil && responded == true
	}, ingressWait, waitTick)

	t.Log("setting the port in ParentRef which does not have a matching listener in Gateway")
	require.Eventually(t, func() bool {
		tlsRoute, err = gatewayClient.GatewayV1alpha2().TLSRoutes(ns.Name).Get(ctx, tlsRoute.Name, metav1.GetOptions{})
		if err != nil {
			return false
		}
		notExistingPort := gatewayv1alpha2.PortNumber(81)
		tlsRoute.Spec.ParentRefs[0].Port = &notExistingPort
		tlsRoute, err = gatewayClient.GatewayV1alpha2().TLSRoutes(ns.Name).Update(ctx, tlsRoute, metav1.UpdateOptions{})
		return err == nil
	}, time.Minute, time.Second)

	t.Log("ensuring tls echo does not respond after using not existing port")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTLSServicePort),
			testUUID, tlsRouteHostname, tlsRouteHostname, true)
		return responded == false && errors.Is(err, io.EOF)
	}, ingressWait, waitTick)
}

// tlsEchoResponds takes a TLS address URL and a Pod name and checks if a
// go-echo instance is running on that Pod at that address using hostname for SNI.
// It compares an expected message and its length against an expected message, returning true
// if it is and false and an error explanation if it is not.
func tlsEchoResponds(
	url string, podName string, hostname, certHostname string, passthrough bool,
) (bool, error) {
	dialer := net.Dialer{Timeout: time.Second * 10}
	conn, err := tls.DialWithDialer(&dialer,
		"tcp",
		url,
		&tls.Config{
			ServerName:         hostname,
			InsecureSkipVerify: true,
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
	// if we are testing with passthrough, the go-echo service should return a message
	// noting that it is listening in TLS mode.
	if passthrough {
		header = append(header, []byte("\nThrough TLS connection.")...)
	}
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

func createTLSEchoContainer(tlsEchoPort int32, sendMsg string) corev1.Container { //nolint:unparam
	container := generators.NewContainer("tcpecho-"+sendMsg, test.EchoImage, tlsEchoPort)
	const tlsCertDir = "/var/run/certs"
	container.Env = append(container.Env,
		corev1.EnvVar{
			Name:  "POD_NAME",
			Value: sendMsg,
		},
		corev1.EnvVar{
			Name:  "TLS_PORT",
			Value: fmt.Sprint(tlsEchoPort),
		},
		corev1.EnvVar{
			Name:  "TLS_CERT_FILE",
			Value: tlsCertDir + "/tls.crt",
		},
		corev1.EnvVar{
			Name:  "TLS_KEY_FILE",
			Value: tlsCertDir + "/tls.key",
		},
	)
	container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
		Name:      tlsSecretName,
		ReadOnly:  true,
		MountPath: tlsCertDir,
	})
	return container
}
