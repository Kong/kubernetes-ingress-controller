package envtest

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/labels"
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
				} else if !assert.NoError(c, err) {
					t.Logf("Error: %v", err)
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
				} else if !assert.NoError(c, err) {
					t.Logf("Error: %v", err)
				}
			}, 30*time.Second, 100*time.Millisecond)
		})
	}
}

func TestAdmissionWebhook_KongConsumers(t *testing.T) {
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

	t.Logf("creating some static credentials in %s namespace which will be used to test global validation", ns.Name)
	for _, secret := range []*corev1.Secret{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "tuxcreds1",
				Labels: map[string]string{
					labels.CredentialTypeLabel: "basic-auth",
				},
			},
			StringData: map[string]string{
				"username": "tux1",
				"password": "testpass",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "tuxcreds2",
				Labels: map[string]string{
					labels.CredentialTypeLabel: "basic-auth",
				},
			},
			StringData: map[string]string{
				"username": "tux2",
				"password": "testpass",
			},
		},
	} {
		secret := secret.DeepCopy()
		require.NoError(t, ctrlClient.Create(ctx, secret))
		t.Cleanup(func() {
			if err := ctrlClient.Delete(ctx, secret); err != nil && !apierrors.IsNotFound(err) && !errors.Is(err, context.Canceled) {
				assert.NoError(t, err)
			}
		})
	}

	t.Logf("creating a static consumer in %s namespace which will be used to test global validation", ns.Name)
	consumer := &kongv1.KongConsumer{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "statis-consumer-",
			Annotations: map[string]string{
				annotations.IngressClassKey: annotations.DefaultIngressClass,
			},
		},
		Username: "tux",
		CustomID: uuid.NewString(),
		Credentials: []string{
			"tuxcreds1",
			"tuxcreds2",
		},
	}
	require.NoError(t, ctrlClient.Create(ctx, consumer))
	t.Cleanup(func() {
		if err := ctrlClient.Delete(ctx, consumer); err != nil && !apierrors.IsNotFound(err) && !errors.Is(err, context.Canceled) {
			assert.NoError(t, err)
		}
	})

	testCases := []struct {
		name           string
		consumer       *kongv1.KongConsumer
		credentials    []*corev1.Secret
		wantErr        bool
		wantPartialErr string
	}{
		{
			name: "a consumer with no credentials should pass validation",
			consumer: &kongv1.KongConsumer{
				ObjectMeta: metav1.ObjectMeta{
					Name: "testconsumer",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Username: uuid.NewString(),
				CustomID: uuid.NewString(),
			},
			credentials: nil,
			wantErr:     false,
		},
		{
			name: "a consumer with valid credentials should pass validation",
			consumer: &kongv1.KongConsumer{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Username:    "electron",
				CustomID:    uuid.NewString(),
				Credentials: []string{"electronscreds"},
			},
			credentials: []*corev1.Secret{{
				ObjectMeta: metav1.ObjectMeta{
					Name: "electronscreds",
					Labels: map[string]string{
						labels.CredentialTypeLabel: "basic-auth",
					},
				},
				StringData: map[string]string{

					"username": "electron",
					"password": "testpass",
				},
			}},
			wantErr: false,
		},
		{
			name: "a consumer with duplicate credentials which are NOT constrained should pass validation",
			consumer: &kongv1.KongConsumer{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Username: "proton",
				CustomID: uuid.NewString(),
				Credentials: []string{
					"protonscreds1",
					"protonscreds2",
				},
			},
			credentials: []*corev1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "protonscreds1",
						Labels: map[string]string{
							labels.CredentialTypeLabel: "basic-auth",
						},
					},
					StringData: map[string]string{
						"username": "proton",
						"password": "testpass",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "protonscreds2",
						Labels: map[string]string{
							labels.CredentialTypeLabel: "basic-auth",
						},
					},
					StringData: map[string]string{
						"username": "electron", // username is unique constrained
						"password": "testpass", // password is not unique constrained
					},
				},
			},
			wantErr: false,
		},
		{
			name: "a consumer referencing credentials secrets which do not yet exist should fail validation",
			consumer: &kongv1.KongConsumer{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Username: "repairedlawnmower",
				CustomID: uuid.NewString(),
				Credentials: []string{
					"nonexistentcreds",
				},
			},
			wantErr:        true,
			wantPartialErr: "not found",
		},
		{
			name: "a consumer with duplicate credentials which ARE constrained should fail validation",
			consumer: &kongv1.KongConsumer{
				ObjectMeta: metav1.ObjectMeta{
					Name: "brokenshovel",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Username: "neutron",
				CustomID: uuid.NewString(),
				Credentials: []string{
					"neutronscreds1",
					"neutronscreds2",
				},
			},
			credentials: []*corev1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "neutronscreds1",
						Labels: map[string]string{
							labels.CredentialTypeLabel: "basic-auth",
						},
					},
					StringData: map[string]string{
						"username": "neutron",
						"password": "testpass",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "neutronscreds2",
						Labels: map[string]string{
							labels.CredentialTypeLabel: "basic-auth",
						},
					},
					StringData: map[string]string{
						"username": "neutron", // username is unique constrained
						"password": "testpass",
					},
				},
			},
			wantErr:        true,
			wantPartialErr: "unique key constraint violated for username",
		},
		{
			name: "a consumer that provides duplicate credentials which are NOT in violation of unique key constraints should pass validation",
			consumer: &kongv1.KongConsumer{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Username: "reasonablehammer",
				CustomID: uuid.NewString(),
				Credentials: []string{
					"reasonablehammer",
				},
			},
			credentials: []*corev1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "reasonablehammer",
						Labels: map[string]string{
							labels.CredentialTypeLabel: "basic-auth",
						},
					},
					StringData: map[string]string{
						"username": "reasonablehammer",
						"password": "testpass", // not unique constrained, so even though someone else is using this password this should pass
					},
				},
			},
			wantErr: false,
		},
		{
			name: "a consumer that provides credentials that are in violation of unique constraints globally against other existing consumers should fail validation",
			consumer: &kongv1.KongConsumer{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "violating-uniqueness-",
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Username: "unreasonablehammer",
				CustomID: uuid.NewString(),
				Credentials: []string{
					"unreasonablehammer",
				},
			},
			credentials: []*corev1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "unreasonablehammer",
						Labels: map[string]string{
							labels.CredentialTypeLabel: "basic-auth",
						},
					},
					StringData: map[string]string{
						"username": "tux1", // unique constrained with previous created static consumer credentials
						"password": "testpass",
					},
				},
			},
			wantErr:        true,
			wantPartialErr: "unique key constraint violated for username",
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			for _, credential := range tc.credentials {
				require.NoError(t, ctrlClient.Create(ctx, credential))
				t.Cleanup(func() {
					if err := ctrlClient.Delete(ctx, credential); err != nil && !apierrors.IsNotFound(err) {
						assert.NoError(t, err)
					}
				})
			}

			err := ctrlClient.Create(ctx, tc.consumer)
			if tc.wantErr {
				require.Error(t, err, fmt.Sprintf("consumer %s should fail to create", tc.consumer.Name))
				assert.Contains(t, err.Error(), tc.wantPartialErr,
					"got error string %q, want a superstring of %q", err.Error(), tc.wantPartialErr,
				)
			} else {
				t.Cleanup(func() {
					if err := ctrlClient.Delete(ctx, tc.consumer); err != nil && !apierrors.IsNotFound(err) {
						assert.NoError(t, err)
					}
				})
				require.NoError(t, err, fmt.Sprintf("consumer %s should create successfully", tc.consumer.Name))
			}
		})
	}
}

func TestAdmissionWebhook_SecretCredentials(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// highEndConsumerUsageCount indicates a number of consumers with credentials
	// that we consider a large number and is used to generate background
	// consumers for testing validation (since validation relies on listing all
	// consumers from the controller runtime cached client).
	const highEndConsumerUsageCount = 50

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

	createKongConsumers(ctx, t, ctrlClient, highEndConsumerUsageCount)

	t.Run("attaching secret to consumer", func(t *testing.T) {
		t.Log("verifying that an invalid credential secret not yet referenced by a KongConsumer fails validation")
		require.Error(t,
			ctrlClient.Create(ctx,
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name: "brokenfence",
						Labels: map[string]string{
							labels.CredentialTypeLabel: "invalid-auth",
						},
					},
					StringData: map[string]string{
						"username": "brokenfence",
						"password": "testpass",
					},
				},
			),
			"invalid credential type",
		)

		t.Log("creating a valid credential secret to be referenced by a KongConsumer")
		validCredential := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: "brokenfence",
				Labels: map[string]string{
					labels.CredentialTypeLabel: "basic-auth",
				},
			},
			StringData: map[string]string{
				"username": "brokenfence",
				"password": "testpass",
			},
		}
		require.NoError(t, ctrlClient.Create(ctx, validCredential))
		t.Cleanup(func() {
			err := ctrlClient.Delete(ctx, validCredential)
			if err != nil && !apierrors.IsNotFound(err) && !errors.Is(err, context.Canceled) {
				assert.NoError(t, err)
			}
		})

		t.Log("verifying that valid credentials assigned to a consumer pass validation")
		validConsumerLinkedToValidCredentials := &kongv1.KongConsumer{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "valid-consumer-",
				Annotations: map[string]string{
					annotations.IngressClassKey: annotations.DefaultIngressClass,
				},
			},
			Username: "brokenfence",
			CustomID: uuid.NewString(),
			Credentials: []string{
				"brokenfence",
			},
		}
		require.NoError(t, ctrlClient.Create(ctx, validConsumerLinkedToValidCredentials))
		t.Cleanup(func() {
			err := ctrlClient.Delete(ctx, validConsumerLinkedToValidCredentials)
			if err != nil && !apierrors.IsNotFound(err) && !errors.Is(err, context.Canceled) {
				assert.NoError(t, err)
			}
		})

		t.Log("verifying that the valid credentials which include a unique-constrained key can be updated in place")
		validCredential.Data["value"] = []byte("newpassword")
		require.NoError(t, ctrlClient.Update(ctx, validCredential))

		t.Log("verifying that validation fails if the now referenced and valid credential gets updated to become invalid")
		validCredential.ObjectMeta.Labels[labels.CredentialTypeLabel] = "invalid-auth"
		err := ctrlClient.Update(ctx, validCredential)
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid credential type")

		t.Log("verifying that if the referent consumer goes away the validation fails for updates that make the credential invalid")
		require.NoError(t, ctrlClient.Delete(ctx, validConsumerLinkedToValidCredentials))
		require.ErrorContains(t, ctrlClient.Update(ctx, validCredential), "invalid credential type")
	})

	t.Run("JWT", func(t *testing.T) {
		t.Log("verifying that a JWT credential which has keys with missing values fails validation")
		require.ErrorContains(t,
			ctrlClient.Create(ctx, &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "invalid-jwt-",
					Labels: map[string]string{
						labels.CredentialTypeLabel: "jwt",
					},
				},
				StringData: map[string]string{
					"algorithm": "RS256",
				},
			}),
			"missing required field(s): rsa_public_key, key, secret",
		)

		hmacAlgos := []string{"HS256", "HS384", "HS512"}

		t.Log("verifying that a JWT credentials with hmac algorithms do not require rsa_public_key field")
		for _, algo := range hmacAlgos {
			t.Run(algo, func(t *testing.T) {
				require.NoError(t, ctrlClient.Create(ctx, &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						GenerateName: "valid-jwt-" + strings.ToLower(algo) + "-",
						Labels: map[string]string{
							labels.CredentialTypeLabel: "jwt",
						},
					},
					StringData: map[string]string{
						"algorithm": algo,
						"key":       "key-name",
						"secret":    "secret-name",
					},
				}), "failed to create JWT credential with algorithm %s", algo)
			})
		}

		nonHmacAlgos := []string{"RS256", "RS384", "RS512", "ES256", "ES384", "ES512"}
		t.Log("verifying that a JWT credentials with non hmac algorithms do require rsa_public_key field")
		for _, algo := range nonHmacAlgos {
			t.Run(algo, func(t *testing.T) {
				require.Error(t, ctrlClient.Create(ctx, &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						GenerateName: "invalid-jwt-" + strings.ToLower(algo) + "-",
						Labels: map[string]string{
							labels.CredentialTypeLabel: "jwt",
						},
					},
					StringData: map[string]string{
						"algorithm": algo,
						"key":       "key-name",
						"secret":    "secret-name",
					},
				}), "expected failure when creating JWT %s", algo)
			})
		}
	})
}

// createKongConsumers creates a provider number of consumers on the cluster.
// Resources will be created in client's default namespace. When using controller-runtime's
// client you can specify that by calling client.NewNamespacedClient(client, namespace).
func createKongConsumers(ctx context.Context, t *testing.T, cl client.Client, count int) {
	t.Helper()

	t.Logf("creating #%d of consumers on the cluster to verify the performance of the cached client during validation", count)

	errg := errgroup.Group{}
	for i := 0; i < count; i++ {
		i := i
		errg.Go(func() error {
			consumerName := fmt.Sprintf("background-noise-consumer-%d", i)

			// create 5 credentials for each consumer
			for j := 0; j < 5; j++ {
				credentialName := fmt.Sprintf("%s-credential-%d", consumerName, j)
				credential := &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name: credentialName,
						Labels: map[string]string{
							labels.CredentialTypeLabel: "basic-auth",
						},
					},
					StringData: map[string]string{
						"username": credentialName,
						"password": "testpass",
					},
				}
				t.Logf("creating %s Secret that contains credentials", credentialName)
				require.NoError(t, cl.Create(ctx, credential))
				t.Cleanup(func() {
					if err := cl.Delete(ctx, credential); err != nil && !apierrors.IsNotFound(err) && !errors.Is(err, context.Canceled) {
						assert.NoError(t, err)
					}
				})
			}

			// create the consumer referencing its credentials
			consumer := &kongv1.KongConsumer{
				ObjectMeta: metav1.ObjectMeta{
					Name: consumerName,
					Annotations: map[string]string{
						annotations.IngressClassKey: annotations.DefaultIngressClass,
					},
				},
				Username: consumerName,
				CustomID: uuid.NewString(),
			}
			for j := 0; j < 5; j++ {
				credentialName := fmt.Sprintf("%s-credential-%d", consumerName, j)
				consumer.Credentials = append(consumer.Credentials, credentialName)
			}
			t.Logf("creating %s KongConsumer", consumerName)
			require.EventuallyWithT(t, func(c *assert.CollectT) {
				assert.NoError(c, cl.Create(ctx, consumer))
			}, 10*time.Second, 100*time.Millisecond)
			t.Cleanup(func() {
				if err := cl.Delete(ctx, consumer); err != nil && !apierrors.IsNotFound(err) && !errors.Is(err, context.Canceled) {
					assert.NoError(t, err)
				}
			})
			return nil
		})
	}
	require.NoError(t, errg.Wait())
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
