package envtest

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	admregv1 "k8s.io/api/admissionregistration/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

func setupValidatingWebhookConfiguration(
	ctx context.Context,
	t *testing.T,
	webhookServerListenPort int,
	cert []byte,
	ctrlClient client.Client,
) {
	webhookConfig := validatingWebhookConfigWithClientConfig(t, admregv1.WebhookClientConfig{
		URL:      lo.ToPtr(fmt.Sprintf("https://localhost:%d/", webhookServerListenPort)),
		CABundle: cert,
	})
	require.NoError(t, ctrlClient.Create(ctx, webhookConfig))
	_ = webhookConfig
}

func validatingWebhookConfigWithClientConfig(t *testing.T, clientConfig admregv1.WebhookClientConfig) *admregv1.ValidatingWebhookConfiguration {
	file, err := os.ReadFile("../../config/webhook/manifests.yaml")
	require.NoError(t, err)
	file = bytes.ReplaceAll(file, []byte("---"), []byte("")) // We're only expecting one document in the file.

	// Load the webhook configuration from the generated manifest.
	webhookConfig := &admregv1.ValidatingWebhookConfiguration{}
	require.NoError(t, yaml.Unmarshal(file, webhookConfig))

	// Set the client config.
	for i := range webhookConfig.Webhooks {
		webhookConfig.Webhooks[i].ClientConfig = clientConfig
	}

	return webhookConfig
}
