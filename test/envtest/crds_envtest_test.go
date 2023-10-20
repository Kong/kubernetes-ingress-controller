//go:build envtest

package envtest

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/go-logr/zapr"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/v2/test/consts"
)

// TestGatewayAPIControllersMayBeDynamicallyStarted ensures that in case of missing CRDs installation in the
// cluster, specific controllers are not started until the CRDs are installed.
func TestGatewayAPIControllersMayBeDynamicallyStarted(t *testing.T) {
	t.Parallel()

	scheme := Scheme(t, WithKong)
	envcfg := Setup(t, scheme, WithInstallGatewayCRDs(false))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	loggerHook := RunManager(ctx, t, envcfg,
		AdminAPIOptFns(),
		WithGatewayFeatureEnabled,
		WithPublishService("ns"),
	)

	controllers := []string{
		"Gateway",
		"HTTPRoute",
		"ReferenceGrant",
		"UDPRoute",
		"TCPRoute",
		"TLSRoute",
		"GRPCRoute",
	}

	requireLogForAllControllers := func(expectedLog string) {
		require.Eventually(t, func() bool {
			for _, controller := range controllers {
				if !lo.ContainsBy(loggerHook.All(), func(entry observer.LoggedEntry) bool {
					return strings.Contains(entry.LoggerName, controller) && strings.Contains(entry.Message, expectedLog)
				}) {
					t.Logf("expected log not found for %s controller", controller)
					return false
				}
			}
			return true
		}, time.Minute, time.Millisecond*500)
	}

	const (
		expectedLogOnStartup      = "Required CustomResourceDefinitions are not installed, setting up a watch for them in case they are installed afterward"
		expectedLogOnCRDInstalled = "All required CustomResourceDefinitions are installed, setting up the controller"
	)

	t.Log("waiting for all controllers to not start due to missing CRDs")
	requireLogForAllControllers(expectedLogOnStartup)

	t.Log("installing missing CRDs")
	installGatewayCRDs(t, scheme, envcfg)

	t.Log("waiting for all controllers to start after CRDs installation")
	requireLogForAllControllers(expectedLogOnCRDInstalled)
}

// TestNoKongCRDsInstalledIsFatal ensures that in case of missing Kong CRDs installation, the manager Run() eventually
// returns an error due to cache synchronisation timeout.
func TestNoKongCRDsInstalledIsFatal(t *testing.T) {
	t.Parallel()

	scheme := Scheme(t)
	envcfg := Setup(t, scheme, WithInstallKongCRDs(false))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg := ConfigForEnvConfig(t, envcfg)

	logger := zapr.NewLogger(zap.NewNop())
	ctrl.SetLogger(logger)

	// Reducing the cache sync timeout to speed up the test.
	cfg.CacheSyncTimeout = time.Millisecond * 500
	err := manager.Run(ctx, &cfg, util.ConfigDumpDiagnostic{}, logger)
	require.ErrorContains(t, err, "timed out waiting for cache to be synced")
}

func TestCRDValidations(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	scheme := Scheme(t, WithKong)
	envcfg := Setup(t, scheme)
	client := NewControllerClient(t, scheme, envcfg)

	testCases := []struct {
		name     string
		scenario func(ctx context.Context, t *testing.T, ns string)
	}{
		{
			name: "invalid TCPIngress service name",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createFaultyTCPIngress(ctx, t, envcfg, ns, func(ingress *kongv1beta1.TCPIngress) {
					ingress.Spec.Rules[0].Backend.ServiceName = ""
				})

				require.ErrorContains(t, err, "serviceName")
			},
		},
		{
			name: "invalid TCPIngress service port",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createFaultyTCPIngress(ctx, t, envcfg, ns, func(ingress *kongv1beta1.TCPIngress) {
					ingress.Spec.Rules[0].Backend.ServicePort = 0
				})

				require.ErrorContains(t, err, "servicePort")
			},
		},
		{
			name: "invalid TCPIngress rule port",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createFaultyTCPIngress(ctx, t, envcfg, ns, func(ingress *kongv1beta1.TCPIngress) {
					ingress.Spec.Rules[0].Port = 0
				})

				require.ErrorContains(t, err, "spec.rules[0].port")
			},
		},
		{
			name: "invalid UDPIngress service name",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createFaultyUDPIngress(ctx, t, envcfg, ns, func(ingress *kongv1beta1.UDPIngress) {
					ingress.Spec.Rules[0].Backend.ServiceName = ""
				})

				require.ErrorContains(t, err, "serviceName")
			},
		},
		{
			name: "invalid UDPIngress service port",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createFaultyUDPIngress(ctx, t, envcfg, ns, func(ingress *kongv1beta1.UDPIngress) {
					ingress.Spec.Rules[0].Backend.ServicePort = 0
				})

				require.ErrorContains(t, err, "servicePort")
			},
		},
		{
			name: "invalid UDPIngress rule port",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createFaultyUDPIngress(ctx, t, envcfg, ns, func(ingress *kongv1beta1.UDPIngress) {
					ingress.Spec.Rules[0].Port = 0
				})

				require.ErrorContains(t, err, "spec.rules[0].port")
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ns := CreateNamespace(ctx, t, client)
			tc.scenario(ctx, t, ns.Name)
		})
	}
}

func createFaultyTCPIngress(ctx context.Context, t *testing.T, envcfg *rest.Config, ns string, modifier func(*kongv1beta1.TCPIngress)) error {
	ingress := validTCPIngress()
	modifier(ingress)

	gatewayClient, err := clientset.NewForConfig(envcfg)
	require.NoError(t, err)

	c := gatewayClient.ConfigurationV1beta1().TCPIngresses(ns)
	ingress, err = c.Create(ctx, ingress, metav1.CreateOptions{})
	if !assert.Error(t, err) {
		t.Cleanup(func() { _ = c.Delete(ctx, ingress.Name, metav1.DeleteOptions{}) })
	}
	return err
}

func validTCPIngress() *kongv1beta1.TCPIngress {
	return &kongv1beta1.TCPIngress{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
			Annotations: map[string]string{
				annotations.IngressClassKey: consts.IngressClass,
			},
		},
		Spec: kongv1beta1.TCPIngressSpec{
			Rules: []kongv1beta1.IngressRule{
				{
					Port: 80,
					Backend: kongv1beta1.IngressBackend{
						ServiceName: "service-name",
						ServicePort: 80,
					},
				},
			},
		},
	}
}

func createFaultyUDPIngress(ctx context.Context, t *testing.T, envcfg *rest.Config, ns string, modifier func(ingress *kongv1beta1.UDPIngress)) error {
	ingress := validUDPIngress()
	modifier(ingress)

	gatewayClient, err := clientset.NewForConfig(envcfg)
	require.NoError(t, err)

	c := gatewayClient.ConfigurationV1beta1().UDPIngresses(ns)
	ingress, err = c.Create(ctx, ingress, metav1.CreateOptions{})
	if !assert.Error(t, err) {
		t.Cleanup(func() { _ = c.Delete(ctx, ingress.Name, metav1.DeleteOptions{}) })
	}
	return err
}

func validUDPIngress() *kongv1beta1.UDPIngress {
	return &kongv1beta1.UDPIngress{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
			Annotations: map[string]string{
				annotations.IngressClassKey: consts.IngressClass,
			},
		},
		Spec: kongv1beta1.UDPIngressSpec{
			Rules: []kongv1beta1.UDPIngressRule{
				{
					Port: 80,
					Backend: kongv1beta1.IngressBackend{
						ServiceName: "service-name",
						ServicePort: 80,
					},
				},
			},
		},
	}
}
