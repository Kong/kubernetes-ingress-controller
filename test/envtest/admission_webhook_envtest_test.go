package envtest

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
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

func TestAdmissionWebhook_KongPlugins(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		scheme           = Scheme(t, WithKong)
		envcfg           = Setup(t, scheme)
		ctrlClientGlobal = NewControllerClient(t, scheme, envcfg)
		ns               = CreateNamespace(ctx, t, ctrlClientGlobal)
		ctrlClient       = client.NewNamespacedClient(ctrlClientGlobal, ns.Name)

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
	)
	WaitForManagerStart(t, logs)
	setupValidatingWebhookConfiguration(ctx, t, admissionWebhookPort, webhookCert, ctrlClient)

	testCases := []struct {
		name                string
		kongPlugin          *kongv1.KongPlugin
		expectErrorContains string
		secretBefore        *corev1.Secret
		secretAfter         *corev1.Secret
		errorOnUpdate       bool
		errorContains       string
	}{
		{
			name: "should fail the validation if secret used in ConfigFrom of KongPlugin generates invalid plugin configuration",
			kongPlugin: &kongv1.KongPlugin{
				ObjectMeta: metav1.ObjectMeta{
					Name: "rate-limiting-invalid-config-from",
				},
				PluginName: "rate-limiting",
				ConfigFrom: &kongv1.ConfigSource{
					SecretValue: kongv1.SecretValueFromSource{
						Secret: "conf-secret-invalid-config",
						Key:    "rate-limiting-config",
					},
				},
			},
			secretBefore: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "conf-secret-invalid-config",
				},
				Data: map[string][]byte{
					"rate-limiting-config": []byte(`{"limit_by":"consumer","policy":"local","minute":5}`),
				},
			},
			secretAfter: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "conf-secret-invalid-config",
				},
				Data: map[string][]byte{
					"rate-limiting-config": []byte(`{"limit_by":"consumer","policy":"local","minute":"5"}`),
				},
			},
			errorOnUpdate: true,
			errorContains: "Change on secret will generate invalid configuration for KongPlugin",
		},
		{
			name: "should fail the validation if the secret is used in ConfigPatches of KongPlugin and generates invalid config",
			kongPlugin: &kongv1.KongPlugin{
				ObjectMeta: metav1.ObjectMeta{
					Name: "rate-limiting-invalid-config-patches",
				},
				PluginName: "rate-limiting",
				Config: apiextensionsv1.JSON{
					Raw: []byte(`{"limit_by":"consumer","policy":"local"}`),
				},
				ConfigPatches: []kongv1.ConfigPatch{
					{
						Path: "/minute",
						ValueFrom: kongv1.ConfigSource{
							SecretValue: kongv1.SecretValueFromSource{
								Secret: "conf-secret-invalid-field",
								Key:    "rate-limiting-config-minutes",
							},
						},
					},
				},
			},
			secretBefore: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "conf-secret-invalid-field",
				},
				Data: map[string][]byte{
					"rate-limiting-config-minutes": []byte("10"),
				},
			},
			secretAfter: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "conf-secret-invalid-field",
				},
				Data: map[string][]byte{
					"rate-limiting-config-minutes": []byte(`"10"`),
				},
			},
			errorOnUpdate: true,
			errorContains: "Change on secret will generate invalid configuration for KongPlugin",
		},
		{
			name: "should pass the validation if the secret used in ConfigPatches of KongPlugin and generates valid config",
			kongPlugin: &kongv1.KongPlugin{
				ObjectMeta: metav1.ObjectMeta{
					Name: "rate-limiting-valid-config",
				},
				PluginName: "rate-limiting",
				Config: apiextensionsv1.JSON{
					Raw: []byte(`{"limit_by":"consumer","policy":"local"}`),
				},
				ConfigPatches: []kongv1.ConfigPatch{
					{
						Path: "/minute",
						ValueFrom: kongv1.ConfigSource{
							SecretValue: kongv1.SecretValueFromSource{
								Secret: "conf-secret-valid-field",
								Key:    "rate-limiting-config-minutes",
							},
						},
					},
				},
			},
			secretBefore: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "conf-secret-valid-field",
				},
				Data: map[string][]byte{
					"rate-limiting-config-minutes": []byte(`10`),
				},
			},
			secretAfter: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "conf-secret-valid-field",
				},
				Data: map[string][]byte{
					"rate-limiting-config-minutes": []byte(`15`),
				},
			},
			errorOnUpdate: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, ctrlClient.Create(ctx, tc.secretBefore))
			t.Cleanup(func() {
				require.NoError(t, ctrlClient.Delete(ctx, tc.secretBefore))
			})

			require.NoError(t, ctrlClient.Create(ctx, tc.kongPlugin))
			t.Cleanup(func() {
				require.NoError(t, ctrlClient.Delete(ctx, tc.kongPlugin))
			})

			require.EventuallyWithT(t, func(c *assert.CollectT) {
				err := ctrlClient.Update(ctx, tc.secretAfter)
				if tc.errorOnUpdate {
					if !assert.Error(c, err) {
						return
					}
					assert.Contains(c, err.Error(), tc.expectErrorContains)
				} else {
					if !assert.NoError(c, err) {
						t.Logf("Error: %v", err)
					}
				}
			}, 10*time.Second, 100*time.Millisecond)
		})
	}
}

func TestAdmissionWebhook_KongClusterPlugins(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		scheme           = Scheme(t, WithKong)
		envcfg           = Setup(t, scheme)
		ctrlClientGlobal = NewControllerClient(t, scheme, envcfg)
		ns               = CreateNamespace(ctx, t, ctrlClientGlobal)
		ctrlClient       = client.NewNamespacedClient(ctrlClientGlobal, ns.Name)

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
	)
	WaitForManagerStart(t, logs)
	setupValidatingWebhookConfiguration(ctx, t, admissionWebhookPort, webhookCert, ctrlClient)

	testCases := []struct {
		name                string
		kongClusterPlugin   *kongv1.KongClusterPlugin
		expectErrorContains string
		secretBefore        *corev1.Secret
		secretAfter         *corev1.Secret
		errorOnUpdate       bool
		errorContains       string
	}{
		{
			name: "should pass the validation if the secret used in ConfigFrom of KongClusterPlugin generates valid configuration",
			kongClusterPlugin: &kongv1.KongClusterPlugin{
				ObjectMeta: metav1.ObjectMeta{
					Name: "cluster-rate-limiting-valid",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				PluginName: "rate-limiting",
				ConfigFrom: &kongv1.NamespacedConfigSource{
					SecretValue: kongv1.NamespacedSecretValueFromSource{
						Namespace: ns.Name,
						Secret:    "cluster-conf-secret-valid",
						Key:       "rate-limiting-config",
					},
				},
			},
			secretBefore: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "cluster-conf-secret-valid",
				},
				Data: map[string][]byte{
					"rate-limiting-config": []byte(`{"limit_by":"consumer","policy":"local","minute":5}`),
				},
			},
			secretAfter: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "cluster-conf-secret-valid",
				},
				Data: map[string][]byte{
					"rate-limiting-config": []byte(`{"limit_by":"consumer","policy":"local","minute":10}`),
				},
			},
			errorOnUpdate: false,
		},
		{
			name: "should fail the validation if the secret in ConfigFrom of KongClusterPlugin generates invalid configuration",
			kongClusterPlugin: &kongv1.KongClusterPlugin{
				ObjectMeta: metav1.ObjectMeta{
					Name: "cluster-rate-limiting-invalid",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				PluginName: "rate-limiting",
				ConfigFrom: &kongv1.NamespacedConfigSource{
					SecretValue: kongv1.NamespacedSecretValueFromSource{
						Namespace: ns.Name,
						Secret:    "cluster-conf-secret-invalid",
						Key:       "rate-limiting-config",
					},
				},
			},
			secretBefore: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "cluster-conf-secret-invalid",
				},
				Data: map[string][]byte{
					"rate-limiting-config": []byte(`{"limit_by":"consumer","policy":"local","minute":5}`),
				},
			},
			secretAfter: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "cluster-conf-secret-invalid",
				},
				Data: map[string][]byte{
					"rate-limiting-config": []byte(`{"limit_by":"consumer","policy":"local","minute":"5"}`),
				},
			},
			errorOnUpdate: true,
			errorContains: "Change on secret will generate invalid configuration for KongClusterPlugin",
		},
		{
			name: "should pass the validation if the secret in ConfigPatches of KongClusterPlugin generates valid configuration",
			kongClusterPlugin: &kongv1.KongClusterPlugin{
				ObjectMeta: metav1.ObjectMeta{
					Name: "cluster-rate-limiting-valid-config-patches",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				PluginName: "rate-limiting",
				Config: apiextensionsv1.JSON{
					Raw: []byte(`{"limit_by":"consumer","policy":"local"}`),
				},
				ConfigPatches: []kongv1.NamespacedConfigPatch{
					{
						Path: "/minute",
						ValueFrom: kongv1.NamespacedConfigSource{
							SecretValue: kongv1.NamespacedSecretValueFromSource{
								Namespace: ns.Name,
								Secret:    "cluster-conf-secret-valid-patch",
								Key:       "rate-limiting-minute",
							},
						},
					},
				},
			},
			secretBefore: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "cluster-conf-secret-valid-patch",
				},
				Data: map[string][]byte{
					"rate-limiting-minute": []byte(`5`),
				},
			},
			secretAfter: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "cluster-conf-secret-valid-patch",
				},
				Data: map[string][]byte{
					"rate-limiting-minute": []byte(`10`),
				},
			},
			errorOnUpdate: false,
		},
		{
			name: "should fail the validation if the secret in ConfigPatches of KongClusterPlugin generates invalid configuration",
			kongClusterPlugin: &kongv1.KongClusterPlugin{
				ObjectMeta: metav1.ObjectMeta{
					Name: "cluster-rate-limiting-invalid-config-patches",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				PluginName: "rate-limiting",
				Config: apiextensionsv1.JSON{
					Raw: []byte(`{"limit_by":"consumer","policy":"local"}`),
				},
				ConfigPatches: []kongv1.NamespacedConfigPatch{
					{
						Path: "/minute",
						ValueFrom: kongv1.NamespacedConfigSource{
							SecretValue: kongv1.NamespacedSecretValueFromSource{
								Namespace: ns.Name,
								Secret:    "cluster-conf-secret-invalid-patch",
								Key:       "rate-limiting-minute",
							},
						},
					},
				},
			},
			secretBefore: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "cluster-conf-secret-invalid-patch",
				},
				Data: map[string][]byte{
					"rate-limiting-minute": []byte(`5`),
				},
			},
			secretAfter: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "cluster-conf-secret-invalid-patch",
				},
				Data: map[string][]byte{
					"rate-limiting-minute": []byte(`"10"`),
				},
			},
			errorOnUpdate: true,
			errorContains: "Change on secret will generate invalid configuration for KongClusterPlugin",
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, ctrlClient.Create(ctx, tc.secretBefore))
			t.Cleanup(func() {
				require.NoError(t, ctrlClient.Delete(ctx, tc.secretBefore))
			})

			require.NoError(t, ctrlClientGlobal.Create(ctx, tc.kongClusterPlugin))
			t.Cleanup(func() {
				require.NoError(t, ctrlClientGlobal.Delete(ctx, tc.kongClusterPlugin))
			})

			require.EventuallyWithT(t, func(c *assert.CollectT) {
				err := ctrlClient.Update(ctx, tc.secretAfter)
				if tc.errorOnUpdate {
					if !assert.Error(c, err) {
						return
					}
					assert.Contains(c, err.Error(), tc.expectErrorContains)
				} else {
					if !assert.NoError(c, err) {
						t.Logf("Error: %v", err)
					}
				}
			}, 30*time.Second, 100*time.Millisecond)
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
