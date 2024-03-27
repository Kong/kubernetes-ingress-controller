//go:build integration_tests

package integration

import (
	"context"
	"crypto/tls"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/kong/go-kong/kong"
	ktfkong "github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	pb "github.com/moul/pb/grpcbin/go-grpc"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/builder"
	"github.com/kong/kubernetes-ingress-controller/v2/test/helpers/certificate"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
)

func grpcbinClient(ctx context.Context, url, hostname string) (pb.GRPCBinClient, func() error, error) {
	conn, err := grpc.DialContext(ctx, url,
		grpc.WithTransportCredentials(credentials.NewTLS(
			&tls.Config{
				ServerName:         hostname,
				InsecureSkipVerify: true,
			},
		)),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial GRPC server: %w", err)
	}

	client := pb.NewGRPCBinClient(conn)
	return client, conn.Close, nil
}

func grpcEchoResponds(ctx context.Context, url, hostname, input string) error {
	client, closeConn, err := grpcbinClient(ctx, url, hostname)
	if err != nil {
		return err
	}
	defer closeConn() //nolint:errcheck

	resp, err := client.DummyUnary(ctx, &pb.DummyMessage{
		FString: input,
	})
	if err != nil {
		return fmt.Errorf("failed to send GRPC request: %w", err)
	}
	if resp.FString != input {
		return fmt.Errorf("unexpected response from GRPC server: %s", resp.FString)
	}

	return nil
}

func TestGRPCRouteEssentials(t *testing.T) {
	ctx := context.Background()

	ns, cleaner := helpers.Setup(ctx, t, env)

	t.Log("getting a gateway client")
	gatewayClient, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("deploying a new gatewayClass")
	gatewayClassName := uuid.NewString()
	gwc, err := DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
	require.NoError(t, err)
	cleaner.Add(gwc)

	t.Log("configuring secret")
	tlsRouteExampleTLSCert, tlsRouteExampleTLSKey := certificate.MustGenerateSelfSignedCertPEMFormat(certificate.WithCommonName(tlsRouteHostname))

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			UID:       k8stypes.UID("7428fb98-180b-4702-a91f-61351a33c6e8"),
			Name:      tlsSecretName,
			Namespace: ns.Name,
		},
		Data: map[string][]byte{
			"tls.crt": tlsRouteExampleTLSCert,
			"tls.key": tlsRouteExampleTLSKey,
		},
	}

	t.Log("deploying secret")
	secret, err = env.Cluster().Client().CoreV1().Secrets(ns.Name).Create(ctx, secret, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(secret)

	t.Log("deploying a new gateway")
	testHostname := "cholpon.example"
	gateway, err := DeployGateway(ctx, gatewayClient, ns.Name, unmanagedGatewayClassName, func(gw *gatewayv1.Gateway) {
		gw.Spec.Listeners = builder.NewListener("https").
			HTTPS().
			WithPort(ktfkong.DefaultProxyTLSServicePort).
			WithHostname(testHostname).
			WithTLSConfig(&gatewayv1.GatewayTLSConfig{
				CertificateRefs: []gatewayv1.SecretObjectReference{
					{
						Name: gatewayv1.ObjectName(secret.Name),
					},
				},
			}).IntoSlice()
	})
	require.NoError(t, err)
	cleaner.Add(gateway)

	grpcPort := int32(9001)
	grpcPortNumber := gatewayv1.PortNumber(grpcPort)
	t.Log("deploying a minimal GRPC container deployment to test Ingress routes")
	container := generators.NewContainer("grpcbin", "moul/grpcbin", grpcPort)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err = env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Logf("creating an grpcroute to access deployment %s via kong", deployment.Name)

	grpcRoute := &gatewayv1alpha2.GRPCRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name: "cholpon-grpcroute",
		},
		Spec: gatewayv1alpha2.GRPCRouteSpec{
			CommonRouteSpec: gatewayv1.CommonRouteSpec{
				ParentRefs: []gatewayv1.ParentReference{{
					Name: gatewayv1.ObjectName(gateway.Name),
				}},
			},
			Hostnames: []gatewayv1alpha2.Hostname{
				gatewayv1alpha2.Hostname(testHostname),
			},
			Rules: []gatewayv1alpha2.GRPCRouteRule{{
				Matches: []gatewayv1alpha2.GRPCRouteMatch{
					{
						// this will match only the DummyUnary method without any headers
						Method: &gatewayv1alpha2.GRPCMethodMatch{
							Service: kong.String("grpcbin.GRPCBin"),
							Method:  kong.String("DummyUnary"),
						},
					},
					{
						// this will match all methods with the x-hello header
						Headers: []gatewayv1alpha2.GRPCHeaderMatch{
							{
								Name:  gatewayv1alpha2.GRPCHeaderName("x-hello"),
								Value: "bidi",
							},
						},
					},
				},
				BackendRefs: []gatewayv1alpha2.GRPCBackendRef{{
					BackendRef: gatewayv1alpha2.BackendRef{
						BackendObjectReference: gatewayv1.BackendObjectReference{
							Name: gatewayv1.ObjectName(service.Name),
							Port: &grpcPortNumber,
						},
					},
				}},
			}},
		},
	}

	grpcRoute, err = gatewayClient.GatewayV1alpha2().GRPCRoutes(ns.Name).Create(ctx, grpcRoute, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(grpcRoute)

	t.Log("verifying that the Gateway gets linked to the route via status")
	callback := GetGatewayIsLinkedCallback(ctx, t, gatewayClient, gatewayv1.HTTPProtocolType, ns.Name, grpcRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)
	t.Log("verifying that the grpcroute contains 'Programmed' condition")
	require.Eventually(t,
		GetVerifyProgrammedConditionCallback(t, gatewayClient, gatewayv1.HTTPProtocolType, ns.Name, grpcRoute.Name, metav1.ConditionTrue),
		ingressWait, waitTick,
	)

	grpcAddr := fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultProxyTLSServicePort)
	t.Log("waiting for routes from GRPCRoute to become operational")
	require.Eventually(t, func() bool {
		err := grpcEchoResponds(ctx, grpcAddr, testHostname, "kong")
		if err != nil {
			t.Log(err)
		}
		return err == nil
	}, ingressWait, waitTick)

	client, closeGrpcConn, err := grpcbinClient(ctx, grpcAddr, testHostname)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := closeGrpcConn()
		require.NoError(t, err)
	})

	t.Log("ensure that the method HeadersUnary is not matched when headers are not passed")
	require.Eventually(t, func() bool {
		_, err := client.HeadersUnary(ctx, &pb.EmptyMessage{})
		if err == nil {
			t.Log("expected error, got nil")
			return false
		}

		t.Log(err)
		return true
	}, ingressWait, waitTick)

	t.Log("ensure that the method HeadersUnary is matched when headers passed")
	require.Eventually(t, func() bool {
		// Set the headers in the context as that's how grpc-go propagates them.
		md := metadata.New(map[string]string{"x-hello": "bidi"})
		ctx := metadata.NewOutgoingContext(ctx, md)
		_, err := client.HeadersUnary(ctx, &pb.EmptyMessage{}, grpc.Header(&metadata.MD{"x-hello": []string{"bidi"}}))
		if err != nil {
			t.Logf("expected no error, got: %v", err)
			return false
		}

		return true
	}, ingressWait, waitTick)
}
