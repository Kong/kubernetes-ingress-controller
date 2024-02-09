package envtest

import (
	"context"
	"fmt"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	admregv1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func setupValidatingWebhookConfiguration(
	ctx context.Context,
	t *testing.T,
	webhookServerListenPort int,
	cert []byte,
	ctrlClient client.Client,
	rules ...admregv1.RuleWithOperations,
) {
	webhookConfig := admregv1.ValidatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{Name: "kong-vault-admission-webhook"},
		Webhooks: []admregv1.ValidatingWebhook{
			{
				Name:                    "kong-vault-admission-webhook.konghq.com",
				FailurePolicy:           lo.ToPtr(admregv1.Fail),
				SideEffects:             lo.ToPtr(admregv1.SideEffectClassNone),
				TimeoutSeconds:          lo.ToPtr[int32](30),
				AdmissionReviewVersions: []string{"v1"},
				ClientConfig: admregv1.WebhookClientConfig{
					URL:      lo.ToPtr(fmt.Sprintf("https://localhost:%d/", webhookServerListenPort)),
					CABundle: cert,
				},
				Rules: rules,
			},
		},
	}
	require.NoError(t, ctrlClient.Create(ctx, &webhookConfig))
}
