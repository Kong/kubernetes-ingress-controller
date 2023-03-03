//go:build integration_tests
// +build integration_tests

package integration

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/kong/go-kong/kong"
	ktfkong "github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
)

func grpcRequest(input string) string {
	return fmt.Sprintf(`{"greeting": "%s"}`, input)
}

func grpcResponse(input string) string {
	return fmt.Sprintf("{\n  \"reply\": \"hello %s\"\n}\n", input)
}

func grpcEchoResponds(ctx context.Context, url, hostname, input string) (bool, error) {
	args := []string{
		"-d",
		grpcRequest(input),
		"-insecure",
		"-servername",
		hostname,
		url,
		"hello.HelloService.SayHello",
	}
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)

	cmd := exec.CommandContext(ctx, "grpcurl", args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		return false, fmt.Errorf("failed to echo GRPC server STDOUT=(%s) STDERR=(%s): %w", strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()), err)
	}

	return stdout.String() == grpcResponse(input), nil
}

func grpcCurl(ctx context.Context, url, hostname, input, service, method string, headers map[string]string) (bool, error) {
	args := []string{
		"-d",
		grpcRequest(input),
		"-insecure",
		"-servername",
		hostname,
	}
	for name, value := range headers {
		args = append(args, "-rpc-header", fmt.Sprintf("%s:%s", name, value))
	}
	args = append(args, url, fmt.Sprintf("%s.%s", service, method))
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)

	cmd := exec.CommandContext(ctx, "grpcurl", args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		return false, fmt.Errorf("failed to echo GRPC server STDOUT=(%s) STDERR=(%s): %w", strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()), err)
	}

	return stdout.String() == grpcResponse(input), nil
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

	t.Log("deploying a new gateway")
	gatewayName := uuid.NewString()
	gateway, err := DeployGateway(ctx, gatewayClient, ns.Name, gatewayClassName, func(gw *gatewayv1beta1.Gateway) {
		gw.Name = gatewayName
	})
	require.NoError(t, err)
	cleaner.Add(gateway)

	grpcPort := int32(9001)
	grpcPortNumber := gatewayv1beta1.PortNumber(grpcPort)
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
			CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
				ParentRefs: []gatewayv1beta1.ParentReference{{
					Name: gatewayv1beta1.ObjectName(gateway.Name),
				}},
			},
			Hostnames: []gatewayv1alpha2.Hostname{
				"cholpon.example",
			},
			Rules: []gatewayv1alpha2.GRPCRouteRule{{
				Matches: []gatewayv1alpha2.GRPCRouteMatch{
					{
						Method: &gatewayv1alpha2.GRPCMethodMatch{
							Service: kong.String("hello.HelloService"),
							Method:  kong.String("SayHello"),
						},
					},
					{
						Method: &gatewayv1alpha2.GRPCMethodMatch{
							Service: kong.String("hello.HelloService"),
							Method:  kong.String("BidiHello"),
						},
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
						BackendObjectReference: gatewayv1beta1.BackendObjectReference{
							Name: gatewayv1beta1.ObjectName(service.Name),
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
	callback := GetGatewayIsLinkedCallback(t, gatewayClient, gatewayv1beta1.HTTPProtocolType, ns.Name, grpcRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)
	t.Log("verifying that the grpcroute contains 'Programmed' condition")
	require.Eventually(t,
		GetVerifyProgrammedConditionCallback(t, gatewayClient, gatewayv1beta1.HTTPProtocolType, ns.Name, grpcRoute.Name, metav1.ConditionTrue),
		ingressWait, waitTick,
	)

	t.Log("waiting for routes from GRPCRoute to become operational")
	require.Eventually(t, func() bool {
		responded, err := grpcEchoResponds(ctx, fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultProxyTLSServicePort), "cholpon.example", "kong")
		if err != nil {
			t.Log(err)
		}
		return err == nil && responded
	}, ingressWait, waitTick)

	require.Eventually(t, func() bool {
		responded, err := grpcCurl(ctx, fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultProxyTLSServicePort),
			"cholpon.example", "kong",
			"hello.HelloService", "SayHello", map[string]string{})
		if err != nil {
			t.Log(err)
		}
		return err == nil && responded
	}, ingressWait, waitTick)

	require.Eventually(t, func() bool {
		responded, err := grpcCurl(ctx, fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultProxyTLSServicePort),
			"cholpon.example", "kong",
			"hello.HelloService", "BidiHello", map[string]string{"x-hello": "bidi"})
		if err != nil {
			t.Log(err)
		}
		return err == nil && responded
	}, ingressWait, waitTick)
}
