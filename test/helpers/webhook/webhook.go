package webhook

import (
	"testing"

	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/kubectl"
	"github.com/stretchr/testify/require"
	admregv1 "k8s.io/api/admissionregistration/v1"
	"sigs.k8s.io/yaml"
)

func GetWebhookConfigWithKustomize(t *testing.T) *admregv1.ValidatingWebhookConfiguration {
	t.Helper()

	webhookKustomize, err := kubectl.RunKustomize("../../config/webhook")
	require.NoError(t, err)

	webhookConfig := &admregv1.ValidatingWebhookConfiguration{}
	require.NoError(t, yaml.Unmarshal(webhookKustomize, webhookConfig))
	return webhookConfig
}
