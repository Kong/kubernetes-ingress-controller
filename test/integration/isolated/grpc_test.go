//go:build integration_tests

package isolated

import (
	"context"
	"crypto/tls"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	ktfkong "github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	pb "github.com/moul/pb/grpcbin/go-grpc"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/builder"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers/certificate"
	"github.com/kong/kubernetes-ingress-controller/v3/test/integration/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testlabels"
)

func TestGRPCRouteEssentials(t *testing.T) {
	const testHostname = "cholpon.example"

	f := features.
		New("essentials").
		WithLabel(testlabels.NetworkingFamily, testlabels.NetworkingFamilyGatewayAPI).
		WithLabel(testlabels.Kind, testlabels.KindGRPCRoute).
		WithSetup("deploy kong addon into cluster", featureSetup()).
		Assess("deploying Gateway and example GRPC service (without konghq.com/protocol annotation) exposed via GRPCRoute over HTTPS", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			// On purpose omit protocol annotation to test defaulting to "grpcs" that is preserved to not break users' configs.
			cleaner := GetFromCtxForT[*clusters.Cleaner](ctx, t)
			cluster := GetClusterFromCtx(ctx)
			namespace := GetNamespaceForT(ctx, t)

			t.Log("getting a gateway client")
			gatewayClient, err := gatewayclient.NewForConfig(cluster.Config())
			assert.NoError(t, err)
			ctx = SetInCtxForT(ctx, t, gatewayClient)

			t.Log("deploying a new gatewayClass")
			gatewayClassName := uuid.NewString()
			gwc, err := helpers.DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
			assert.NoError(t, err)
			cleaner.Add(gwc)

			t.Log("configuring secret")
			const tlsRouteHostname = "tls-route.example"
			tlsRouteExampleTLSCert, tlsRouteExampleTLSKey := certificate.MustGenerateSelfSignedCertPEMFormat(certificate.WithCommonName(tlsRouteHostname))
			const tlsSecretName = "secret-test"
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					UID:       k8stypes.UID("7428fb98-180b-4702-a91f-61351a33c6e8"),
					Name:      tlsSecretName,
					Namespace: namespace,
				},
				Data: map[string][]byte{
					"tls.crt": tlsRouteExampleTLSCert,
					"tls.key": tlsRouteExampleTLSKey,
				},
			}

			t.Log("deploying secret")
			secret, err = cluster.Client().CoreV1().Secrets(namespace).Create(ctx, secret, metav1.CreateOptions{})
			assert.NoError(t, err)
			cleaner.Add(secret)

			t.Log("deploying a new gateway")
			gateway, err := helpers.DeployGateway(ctx, gatewayClient, namespace, gatewayClassName, func(gw *gatewayapi.Gateway) {
				// Besides default HTTP listener, add a HTTPS listener.
				gw.Spec.Listeners = append(
					gw.Spec.Listeners,
					builder.NewListener("https").
						HTTPS().
						WithPort(ktfkong.DefaultProxyTLSServicePort).
						WithHostname(testHostname).
						WithTLSConfig(&gatewayapi.GatewayTLSConfig{
							CertificateRefs: []gatewayapi.SecretObjectReference{
								{
									Name: gatewayapi.ObjectName(secret.Name),
								},
							},
						}).
						Build(),
				)
			})
			assert.NoError(t, err)
			cleaner.Add(gateway)

			t.Log("deploying a minimal GRPC container deployment to test Ingress routes")
			container := generators.NewContainer("grpcbin", test.GRPCBinImage, test.GRPCBinPort)
			deployment := generators.NewDeploymentForContainer(container)
			deployment, err = cluster.Client().AppsV1().Deployments(namespace).Create(ctx, deployment, metav1.CreateOptions{})
			assert.NoError(t, err)

			t.Logf("exposing deployment %s via service", deployment.Name)
			service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
			_, err = cluster.Client().CoreV1().Services(namespace).Create(ctx, service, metav1.CreateOptions{})
			assert.NoError(t, err)

			t.Logf("creating an GRPCRoute to access deployment %s via Kong", deployment.Name)
			grpcRoute := &gatewayapi.GRPCRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name: "cholpon-grpcroute",
				},
				Spec: gatewayapi.GRPCRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{{
							Name: gatewayapi.ObjectName(gateway.Name),
						}},
					},
					Hostnames: []gatewayapi.Hostname{
						gatewayapi.Hostname(testHostname),
					},
					Rules: []gatewayapi.GRPCRouteRule{{
						Matches: []gatewayapi.GRPCRouteMatch{
							{
								// this will match only the DummyUnary method without any headers
								Method: &gatewayapi.GRPCMethodMatch{
									Service: kong.String("grpcbin.GRPCBin"),
									Method:  kong.String("DummyUnary"),
								},
							},
							{
								// this will match all methods with the x-hello header
								Headers: []gatewayapi.GRPCHeaderMatch{
									{
										Name:  gatewayapi.GRPCHeaderName("x-hello"),
										Value: "bidi",
									},
								},
							},
						},
						BackendRefs: []gatewayapi.GRPCBackendRef{{
							BackendRef: gatewayapi.BackendRef{
								BackendObjectReference: gatewayapi.BackendObjectReference{
									Name: gatewayapi.ObjectName(service.Name),
									Port: lo.ToPtr(gatewayapi.PortNumber(test.GRPCBinPort)),
								},
							},
						}},
					}},
				},
			}

			grpcRoute, err = gatewayClient.GatewayV1alpha2().GRPCRoutes(namespace).Create(ctx, grpcRoute, metav1.CreateOptions{})
			assert.NoError(t, err)
			cleaner.Add(grpcRoute)
			ctx = SetInCtxForT(ctx, t, grpcRoute)

			return ctx
		}).
		Assess("checking if GRPCRoute is linked correctly and client can connect properly to the exposed service", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			proxyURL := GetProxyURLFromCtx(ctx)
			namespace := GetNamespaceForT(ctx, t)
			gatewayClient := GetFromCtxForT[*gatewayclient.Clientset](ctx, t)
			grpcRoute := GetFromCtxForT[*gatewayapi.GRPCRoute](ctx, t)

			t.Log("verifying that the Gateway gets linked to the route via status")
			callback := helpers.GetGatewayIsLinkedCallback(ctx, t, gatewayClient, gatewayapi.HTTPProtocolType, namespace, grpcRoute.Name)
			assert.Eventually(t, callback, consts.IngressWait, consts.WaitTick)
			t.Log("verifying that the GRPCRoute contains 'Programmed' condition")
			assert.Eventually(t,
				helpers.GetVerifyProgrammedConditionCallback(t, gatewayClient, gatewayapi.HTTPProtocolType, namespace, grpcRoute.Name, metav1.ConditionTrue),
				consts.IngressWait, consts.WaitTick,
			)

			grpcAddr := fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultProxyTLSServicePort)
			t.Log("waiting for routes from GRPCRoute to become operational")
			assert.Eventually(t, func() bool {
				err := grpcEchoResponds(ctx, grpcAddr, testHostname, "kong", true)
				if err != nil {
					t.Log(err)
				}
				return err == nil
			}, consts.IngressWait, consts.WaitTick)

			client, closeGrpcConn, err := grpcBinClient(ctx, grpcAddr, testHostname, true)
			assert.NoError(t, err)
			t.Cleanup(func() {
				err := closeGrpcConn()
				assert.NoError(t, err)
			})

			t.Log("ensure that the method HeadersUnary is not matched when headers are not passed")
			assert.Eventually(t, func() bool {
				_, err := client.HeadersUnary(ctx, &pb.EmptyMessage{})
				if err == nil {
					t.Log("expected error, got nil")
					return false
				}

				t.Log(err)
				return true
			}, consts.IngressWait, consts.WaitTick)

			t.Log("ensure that the method HeadersUnary is matched when headers passed")
			assert.Eventually(t, func() bool {
				// Set the headers in the context as that's how grpc-go propagates them.
				md := metadata.New(map[string]string{"x-hello": "bidi"})
				ctx := metadata.NewOutgoingContext(ctx, md)
				_, err := client.HeadersUnary(ctx, &pb.EmptyMessage{}, grpc.Header(&metadata.MD{"x-hello": []string{"bidi"}}))
				if err != nil {
					t.Logf("expected no error, got: %v", err)
					return false
				}

				return true
			}, consts.IngressWait, consts.WaitTick)

			return ctx
		}).
		Teardown(featureTeardown())

	tenv.Test(t, f.Feature())
}

func grpcEchoResponds(ctx context.Context, url, hostname, input string, enableTLS bool) error {
	client, closeConn, err := grpcBinClient(ctx, url, hostname, enableTLS)
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

func grpcBinClient(ctx context.Context, url, hostname string, enableTLS bool) (pb.GRPCBinClient, func() error, error) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithAuthority(hostname)}
	if enableTLS {
		opts = []grpc.DialOption{grpc.WithTransportCredentials(credentials.NewTLS(
			&tls.Config{
				ServerName:         hostname,
				InsecureSkipVerify: true,
			},
		))}
	}
	conn, err := grpc.DialContext(ctx, url, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial GRPC server: %w", err)
	}

	client := pb.NewGRPCBinClient(conn)
	return client, conn.Close, nil
}
