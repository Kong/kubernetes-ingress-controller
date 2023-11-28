//go:build integration_tests

package isolated

import (
	"context"
	"fmt"
	"testing"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	ktfkong "github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers/certificate"
	"github.com/kong/kubernetes-ingress-controller/v3/test/integration/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testlabels"
)

func TestIngressGRPC(t *testing.T) {
	const testHostname = "grpcs-over-ingress.example"

	f := features.
		New("essentials").
		WithLabel(testlabels.NetworkingFamily, testlabels.NetworkingFamilyIngress).
		WithLabel(testlabels.Kind, testlabels.KindIngress).
		WithSetup("deploy kong addon into cluster", featureSetup(
			withKongProxyEnvVars(map[string]string{
				"PROXY_LISTEN": `0.0.0.0:8000 http2\, 0.0.0.0:8443 http2 ssl`,
			}),
		)).
		WithSetup("deploying gRPC service exposed via Ingress", func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
			cleaner := GetFromCtxForT[*clusters.Cleaner](ctx, t)
			cluster := GetClusterFromCtx(ctx)
			namespace := GetNamespaceForT(ctx, t)

			t.Log("configuring secret")
			tlsRouteExampleTLSCert, tlsRouteExampleTLSKey := certificate.MustGenerateSelfSignedCertPEMFormat(certificate.WithCommonName(testHostname))
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "secret-test",
				},
				Data: map[string][]byte{
					"tls.crt": tlsRouteExampleTLSCert,
					"tls.key": tlsRouteExampleTLSKey,
				},
			}

			t.Log("deploying secret")
			secret, err := cluster.Client().CoreV1().Secrets(namespace).Create(ctx, secret, metav1.CreateOptions{})
			assert.NoError(t, err)
			cleaner.Add(secret)

			t.Log("deploying a minimal GRPC service exposed with Ingress via HTTPS")
			ctx = deployGRPCServiceWithIngress(ctx, t, true)

			t.Log("deploying a minimal GRPC service exposed with Ingress via HTTP")
			return deployGRPCServiceWithIngress(ctx, t, false)
		}).
		Assess("checking whether Ingress status is updated and gRPC traffic over HTTPS is properly routed", func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
			t.Log("verifying service connectivity via HTTPS and Ingress status update")
			ctx = verifyGRPCServiceAndIngress(ctx, t, testHostname)

			t.Log("verifying service connectivity via HTTP and Ingress status update")
			return verifyGRPCServiceAndIngress(ctx, t, "")
		}).
		Teardown(featureTeardown())

	tenv.Test(t, f.Feature())
}

func deployGRPCServiceWithIngress(ctx context.Context, t *testing.T, gRPCS bool) context.Context {
	cleaner := GetFromCtxForT[*clusters.Cleaner](ctx, t)
	cluster := GetClusterFromCtx(ctx)
	namespace := GetNamespaceForT(ctx, t)
	ingressClass := GetIngressClassFromCtx(ctx)
	t.Log("deploying a minimal GRPC container deployment to test Ingress routes")

	// Kong distinguishes gRPC over HTTP - 'grpcs' and gRPC over HTTPS - 'grpc'.
	// Furthermore example service uses different ports for each protocol.
	kongProtocol := "grpc"
	grpcBinPort := int32(9000)
	if gRPCS {
		kongProtocol = "grpcs"
		grpcBinPort = test.GRPCBinPort
	}

	container := generators.NewContainer("grpcbin", test.GRPCBinImage, grpcBinPort)
	deployment := generators.NewDeploymentForContainer(container)
	deployment.Name += kongProtocol
	deployment, err := cluster.Client().AppsV1().Deployments(namespace).Create(ctx, deployment, metav1.CreateOptions{})
	assert.NoError(t, err)
	cleaner.Add(deployment)

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	service.Name += kongProtocol
	service.Annotations = map[string]string{
		annotations.AnnotationPrefix + annotations.ProtocolKey: kongProtocol,
	}
	_, err = cluster.Client().CoreV1().Services(namespace).Create(ctx, service, metav1.CreateOptions{})
	assert.NoError(t, err)
	cleaner.Add(deployment)

	t.Logf("creating an ingress for service %s with ingress.class %s", service.Name, ingressClass)
	ingress := generators.NewIngressForService("/", map[string]string{
		annotations.AnnotationPrefix + annotations.ProtocolsKey: kongProtocol,
	}, service)
	ingress.Spec.IngressClassName = &ingressClass
	assert.NoError(t, clusters.DeployIngress(ctx, cluster, namespace, ingress))
	cleaner.Add(ingress)
	ctx = SetInCtxForT(ctx, t, ingress)

	return ctx
}

func verifyGRPCServiceAndIngress(ctx context.Context, t *testing.T, hostname string) context.Context {
	t.Log("waiting for updated ingress status to include IP")
	assert.Eventually(t, func() bool {
		cluster := GetClusterFromCtx(ctx)
		namespace := GetNamespaceForT(ctx, t)
		ingress := GetFromCtxForT[*netv1.Ingress](ctx, t)

		lbstatus, err := clusters.GetIngressLoadbalancerStatus(ctx, cluster, namespace, ingress)
		if err != nil {
			return false
		}
		return len(lbstatus.Ingress) > 0
	}, consts.StatusWait, consts.WaitTick)

	proxyURL := GetProxyURLFromCtx(ctx)
	t.Log("verifying that gRPC service can be accessed via Ingress")
	var tlsEnabled bool
	proxyPort := ktfkong.DefaultProxyHTTPPort
	if hostname != "" {
		tlsEnabled = true
		proxyPort = ktfkong.DefaultProxyTLSServicePort
	}
	assert.Eventually(t, func() bool {
		if err := grpcEchoResponds(
			ctx, fmt.Sprintf("%s:%d", proxyURL.Hostname(), proxyPort), hostname, fmt.Sprintf("echo %q kong", hostname), tlsEnabled,
		); err != nil {
			t.Log(err)
			return false
		}
		return true
	}, consts.IngressWait, consts.WaitTick)

	return ctx
}
