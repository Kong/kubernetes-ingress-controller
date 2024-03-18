package envtest

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1alpha1"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers/certificate"
)

func TestAdmissionWebhook_KongVault(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		scheme     = Scheme(t, WithKong)
		envcfg     = Setup(t, scheme)
		ctrlClient = NewControllerClient(t, scheme, envcfg)
		ns         = CreateNamespace(ctx, t, ctrlClient)

		webhookCert, webhookKey = certificate.MustGenerateSelfSignedCertPEMFormat(
			certificate.WithDNSNames("localhost"),
		)
		admissionWebhookPort = helpers.GetFreePort(t)

		kongContainer = runKongEnterprise(ctx, t)
	)

	_, logs := RunManager(ctx, t, envcfg,
		AdminAPIOptFns(),
		WithPublishService(ns.Name),
		WithAdmissionWebhookEnabled(webhookKey, webhookCert, fmt.Sprintf(":%d", admissionWebhookPort)),
		WithKongAdminURLs(kongContainer.AdminURL(ctx, t)),
		WithUpdateStatus(),
	)
	WaitForManagerStart(t, logs)
	setupValidatingWebhookConfiguration(ctx, t, admissionWebhookPort, webhookCert, ctrlClient)

	const prefixForDuplicationTest = "duplicate-prefix"
	prepareKongVaultAlreadyProgrammedInGateway(ctx, t, ctrlClient, prefixForDuplicationTest)

	testCases := []struct {
		name                string
		kongVault           *kongv1alpha1.KongVault
		expectErrorContains string
	}{
		{
			name: "should pass the validation if the configuration is correct",
			kongVault: &kongv1alpha1.KongVault{
				ObjectMeta: metav1.ObjectMeta{
					Name: "vault-valid",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: kongv1alpha1.KongVaultSpec{
					Backend:     "env",
					Prefix:      "env-test",
					Description: "test env vault",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"prefix":"kong_vault_test_"}`),
					},
				},
			},
		},
		{
			name: "should also pass the validation if the description is empty",
			kongVault: &kongv1alpha1.KongVault{
				ObjectMeta: metav1.ObjectMeta{
					Name: "vault-empty-description",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: kongv1alpha1.KongVaultSpec{
					Backend: "env",
					Prefix:  "env-empty-desc",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"prefix":"kong_vault_test_"}`),
					},
				},
			},
		},
		{
			name: "should fail the validation if the backend is not supported by Kong gateway",
			kongVault: &kongv1alpha1.KongVault{
				ObjectMeta: metav1.ObjectMeta{
					Name: "vault-unsupported-backend",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: kongv1alpha1.KongVaultSpec{
					Backend:     "env1",
					Prefix:      "unsupported-backend",
					Description: "test env vault",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"prefix":"kong-env-test"}`),
					},
				},
			},
			expectErrorContains: `vault configuration in invalid: schema violation (name: vault 'env1' is not installed)`,
		},
		{
			name: "should fail the validation if the spec.config does not pass the schema check of Kong gateway",
			kongVault: &kongv1alpha1.KongVault{
				ObjectMeta: metav1.ObjectMeta{
					Name: "vault-invalid-config",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: kongv1alpha1.KongVaultSpec{
					Backend:     "env",
					Prefix:      "invalid-config",
					Description: "test env vault",
					Config: apiextensionsv1.JSON{
						Raw: []byte(`{"prefix":"kong-env-test","foo":"bar"}`),
					},
				},
			},
			expectErrorContains: `vault configuration in invalid: schema violation (config.foo: unknown field)`,
		},
		{
			name: "should fail the validation if spec.prefix is duplicate",
			kongVault: &kongv1alpha1.KongVault{
				ObjectMeta: metav1.ObjectMeta{
					Name: "vault-dupe",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Spec: kongv1alpha1.KongVaultSpec{
					Backend:     "env",
					Prefix:      prefixForDuplicationTest, // This is the same prefix as the one created in setup.
					Description: "test env vault",
				},
			},
			expectErrorContains: fmt.Sprintf(`spec.prefix "%s" is duplicate with existing KongVault`, prefixForDuplicationTest),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ctrlClient.Create(ctx, tc.kongVault)
			if tc.expectErrorContains != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectErrorContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// prepareKongVaultAlreadyProgrammedInGateway creates a KongVault and waits until it gets programmed in Gateway.
func prepareKongVaultAlreadyProgrammedInGateway(
	ctx context.Context,
	t *testing.T,
	ctrlClient client.Client,
	vaultPrefix string,
) {
	t.Helper()

	const (
		programmedWaitTimeout  = 30 * time.Second
		programmedWaitInterval = 20 * time.Millisecond
	)

	name := uuid.NewString()
	vault := &kongv1alpha1.KongVault{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Annotations: map[string]string{
				annotations.IngressClassKey: annotations.DefaultIngressClass,
			},
		},
		Spec: kongv1alpha1.KongVaultSpec{
			Backend:     "env",
			Prefix:      vaultPrefix,
			Description: "vault description",
		},
	}
	err := ctrlClient.Create(ctx, vault)
	require.NoError(t, err)

	t.Logf("Waiting for KongVault %s to be programmed...", name)
	require.Eventuallyf(t, func() bool {
		kv := &kongv1alpha1.KongVault{}
		err := ctrlClient.Get(ctx, k8stypes.NamespacedName{Name: name}, kv)
		if err != nil {
			return false
		}
		programmed, ok := lo.Find(kv.Status.Conditions, func(c metav1.Condition) bool {
			return c.Type == "Programmed"
		})
		return ok && programmed.Status == metav1.ConditionTrue
	}, programmedWaitTimeout, programmedWaitInterval, "KongVault %s was expected to be programmed", name)
}
