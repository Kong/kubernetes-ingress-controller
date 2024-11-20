//go:build integration_tests

package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/networking"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/builder"
	testutils "github.com/kong/kubernetes-ingress-controller/v3/internal/util/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers/certificate"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers/webhook"
)

// ensureAdmissionRegistration registers a validating webhook for the given configuration, it validates objects
// only when applied to the given namespace.
func ensureAdmissionRegistration(
	ctx context.Context, t *testing.T, client *kubernetes.Clientset, webhookName, namespaceToCheck string,
) {
	t.Helper()

	const svcPort = 443
	webhookService := k8stypes.NamespacedName{
		Namespace: consts.ControllerNamespace,
		Name:      fmt.Sprintf("webhook-%s", webhookName),
	}
	ensureWebhookService(ctx, t, client, webhookService)

	webhookConfig := webhook.GetWebhookConfigWithKustomize(t)
	cert, _ := certificate.GetKongSystemSelfSignedCerts()
	for i := range webhookConfig.Webhooks {
		webhookConfig.Webhooks[i].ClientConfig.Service.Name = webhookService.Name
		webhookConfig.Webhooks[i].ClientConfig.Service.Namespace = webhookService.Namespace
		webhookConfig.Webhooks[i].ClientConfig.CABundle = cert
		webhookConfig.Webhooks[i].NamespaceSelector = &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"kubernetes.io/metadata.name": namespaceToCheck,
			},
		}
	}

	t.Log("creating webhook configuration")
	validationWebhookClient := client.AdmissionregistrationV1().ValidatingWebhookConfigurations()
	webhookConfig, err := validationWebhookClient.Create(ctx, webhookConfig, metav1.CreateOptions{})
	require.NoError(t, err)
	t.Cleanup(func() {
		if err := validationWebhookClient.Delete(ctx, webhookConfig.Name, metav1.DeleteOptions{}); err != nil && !apierrors.IsNotFound(err) {
			require.NoError(t, err)
		}
	})

	t.Log("waiting for webhook service to be connective")
	waitCtx, cancel := context.WithTimeout(ctx, ingressWait)
	defer cancel()
	require.NoError(
		t,
		networking.WaitForConnectionOnServicePort(
			waitCtx, client, webhookService.Namespace, webhookService.Name, svcPort, test.RequestTimeout,
		),
	)
}

func ensureWebhookService(ctx context.Context, t *testing.T, client *kubernetes.Clientset, nn k8stypes.NamespacedName) {
	t.Logf("creating webhook service: %s", nn)
	validationsService, err := client.CoreV1().Services(nn.Namespace).Create(ctx, &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: nn.Name,
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
	endpoints, err := client.DiscoveryV1().EndpointSlices(nn.Namespace).Create(ctx, &discoveryv1.EndpointSlice{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-1", nn.Name),
			Labels: map[string]string{
				discoveryv1.LabelServiceName: nn.Name,
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
		if err := client.CoreV1().Services(nn.Namespace).Delete(ctx, validationsService.Name, metav1.DeleteOptions{}); err != nil && !apierrors.IsNotFound(err) {
			require.NoError(t, err)
		}
		if err := client.DiscoveryV1().EndpointSlices(nn.Namespace).Delete(ctx, endpoints.Name, metav1.DeleteOptions{}); err != nil && !apierrors.IsNotFound(err) {
			require.NoError(t, err)
		}
	})
}
