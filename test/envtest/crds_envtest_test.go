//go:build envtest

package envtest

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/bombsimon/logrusr/v4"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/v2/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
)

func TestManagerDoesntStartUntilKubernetesAPIReachable(t *testing.T) {
	scheme := Scheme(t, WithKong)
	envcfg := Setup(t, scheme)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Log("Setting up a proxy for Kubernetes API server so that we can interrupt it")
	u, err := url.Parse(envcfg.Host)
	require.NoError(t, err)
	apiServerProxy, err := helpers.NewTCPProxy(u.Host)
	require.NoError(t, err)
	go func() {
		err := apiServerProxy.Run(ctx)
		assert.NoError(t, err)
	}()
	apiServerProxy.StopHandlingConnections()

	t.Log("Replacing Kubernetes API server address with the proxy address")
	envcfg.Host = fmt.Sprintf("https://%s", apiServerProxy.Address())

	loggerHook := RunManager(ctx, t, envcfg)
	hasLog := func(expectedLog string) bool {
		return lo.ContainsBy(loggerHook.AllEntries(), func(entry *logrus.Entry) bool {
			return strings.Contains(entry.Message, expectedLog)
		})
	}

	t.Log("Ensuring manager is waiting for Kubernetes API to be ready")
	const expectedKubernetesAPICheckErrorLog = "Retrying Kubernetes API readiness check after error"
	require.Eventually(t, func() bool { return hasLog(expectedKubernetesAPICheckErrorLog) }, time.Minute, time.Millisecond)

	t.Log("Ensure manager hasn't been started yet and no config sync has happened")
	const configurationSyncedToKongLog = "successfully synced configuration to Kong"
	const startingManagerLog = "Starting manager"
	require.False(t, hasLog(configurationSyncedToKongLog))
	require.False(t, hasLog(startingManagerLog))

	t.Log("Starting accepting connections in Kubernetes API proxy so that manager can start")
	apiServerProxy.StartHandlingConnections()

	t.Log("Ensuring manager has been started and config sync has happened")
	require.Eventually(t, func() bool {
		return hasLog(startingManagerLog) &&
			hasLog(configurationSyncedToKongLog)
	}, time.Minute, time.Millisecond)
}

// TestDynamicCRDController_StartsControllersWhenCRDsInstalled ensures that in case of missing CRDs installation in the
// cluster, specific controllers are not started until the CRDs are installed.
func TestDynamicCRDController_StartsControllersWhenCRDsInstalled(t *testing.T) {
	scheme := Scheme(t, WithKong)
	envcfg := Setup(t, scheme, WithInstallGatewayCRDs(false))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	loggerHook := RunManager(ctx, t, envcfg, WithGatewayFeatureEnabled)

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
				if !lo.ContainsBy(loggerHook.AllEntries(), func(entry *logrus.Entry) bool {
					loggerName, ok := entry.Data["logger"].(string)
					if !ok {
						return false
					}
					return strings.Contains(loggerName, controller) && strings.Contains(entry.Message, expectedLog)
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

func TestNoKongCRDsIsFatal(t *testing.T) {
	scheme := Scheme(t)
	envcfg := Setup(t, scheme, WithInstallKongCRDs(false))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg := ConfigForEnvConfig(t, envcfg)

	logrusLogger, _ := test.NewNullLogger()
	logger := logrusr.New(logrusLogger)
	ctrl.SetLogger(logger)

	err := manager.Run(ctx, &cfg, util.ConfigDumpDiagnostic{}, logrusLogger)
	require.ErrorContains(t, err, "requirements not satisfied")
}

func TestCRDValidations(t *testing.T) {
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
