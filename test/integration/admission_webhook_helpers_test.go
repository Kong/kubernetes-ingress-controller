//go:build integration_tests

package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/networking"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	admregv1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/builder"
	testutils "github.com/kong/kubernetes-ingress-controller/v3/internal/util/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers/certificate"
)

// NOTE: Functions in this file use a shared environment.
// This can be refactored so that the client is provided as an argument to the functions.

func ensureWebhookService(ctx context.Context, t *testing.T, name string) {
	t.Logf("creating webhook service: %q in namespace: %q", name, consts.ControllerNamespace)
	validationsService, err := env.Cluster().Client().CoreV1().Services(consts.ControllerNamespace).Create(ctx, &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:       "default",
					Port:       443,
					TargetPort: intstr.FromInt(testutils.AdmissionWebhookListenPort),
				},
			},
		},
	}, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("creating webhook endpoints")
	endpoints, err := env.Cluster().Client().DiscoveryV1().EndpointSlices(consts.ControllerNamespace).Create(ctx, &discoveryv1.EndpointSlice{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-1", name),
			Labels: map[string]string{
				discoveryv1.LabelServiceName: name,
			},
		},
		AddressType: discoveryv1.AddressTypeIPv4,
		Endpoints: []discoveryv1.Endpoint{
			{
				Addresses: []string{testutils.GetAdmissionWebhookListenHost()},
			},
		},
		Ports: builder.NewEndpointPort(testutils.AdmissionWebhookListenPort).WithName("default").WithProtocol(corev1.ProtocolTCP).IntoSlice(),
	}, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Cleanup(func() {
		if err := env.Cluster().Client().CoreV1().Services(consts.ControllerNamespace).Delete(ctx, validationsService.Name, metav1.DeleteOptions{}); err != nil && !apierrors.IsNotFound(err) {
			require.NoError(t, err)
		}
		if err := env.Cluster().Client().DiscoveryV1().EndpointSlices(consts.ControllerNamespace).Delete(ctx, endpoints.Name, metav1.DeleteOptions{}); err != nil && !apierrors.IsNotFound(err) {
			require.NoError(t, err)
		}
	})
}

func ensureWebhookServiceIsConnective(ctx context.Context, t *testing.T, configResourceName string) {
	svcName := fmt.Sprintf("webhook-%s", configResourceName)
	const svcPort = 443
	waitCtx, cancel := context.WithTimeout(ctx, ingressWait)
	defer cancel()
	require.NoError(
		t,
		networking.WaitForConnectionOnServicePort(waitCtx, env.Cluster().Client(), consts.ControllerNamespace, svcName, svcPort, test.RequestTimeout),
	)
}

// ensureAdmissionRegistration registers a validating webhook for the given configuration, it validates objects only when applied to the given namespace.
func ensureAdmissionRegistration(ctx context.Context, t *testing.T, namespace, configResourceName string, rules []admregv1.RuleWithOperations) {
	svcName := fmt.Sprintf("webhook-%s", configResourceName)
	ensureWebhookService(ctx, t, svcName)

	cert, _ := certificate.GetKongSystemSelfSignedCerts()
	webhook, err := env.Cluster().Client().AdmissionregistrationV1().ValidatingWebhookConfigurations().Create(ctx,
		&admregv1.ValidatingWebhookConfiguration{
			ObjectMeta: metav1.ObjectMeta{Name: configResourceName},
			Webhooks: []admregv1.ValidatingWebhook{
				{
					Name:                    "validations.kong.konghq.com",
					FailurePolicy:           lo.ToPtr(admregv1.Ignore),
					SideEffects:             lo.ToPtr(admregv1.SideEffectClassNone),
					AdmissionReviewVersions: []string{"v1beta1", "v1"},
					Rules:                   rules,
					ClientConfig: admregv1.WebhookClientConfig{
						Service:  &admregv1.ServiceReference{Namespace: consts.ControllerNamespace, Name: svcName},
						CABundle: cert,
					},
					NamespaceSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"kubernetes.io/metadata.name": namespace,
						},
					},
				},
			},
		}, metav1.CreateOptions{})
	require.NoError(t, err)
	for _, r := range rules {
		t.Logf(
			"configured admission webhook for: %q that validates in namespace: %q",
			fmt.Sprintf("%s %s %s", r.Rule.APIGroups, r.Rule.APIVersions, r.Rule.Resources), namespace,
		)
	}

	t.Cleanup(func() {
		if err := env.Cluster().Client().AdmissionregistrationV1().ValidatingWebhookConfigurations().Delete(ctx, webhook.Name, metav1.DeleteOptions{}); err != nil && !apierrors.IsNotFound(err) {
			require.NoError(t, err)
		}
	})
}
