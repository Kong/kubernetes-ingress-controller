//go:build integration_tests

package integration

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/kind"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/networking"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	admregv1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/builder"
	testutils "github.com/kong/kubernetes-ingress-controller/v3/internal/util/test"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/v3/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers/certificate"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
)

// extraWebhookNamespace is an additional namespace used by tests when needing
// to run tests that need multiple namespaces.
const extraWebhookNamespace = "webhookextra"

// highEndConsumerUsageCount indicates a number of consumers with credentials
// that we consider a large number and is used to generate background
// consumers for testing validation (since validation relies on listing all
// consumers from the controller runtime cached client).
const highEndConsumerUsageCount = 50

func TestValidationWebhook(t *testing.T) {
	ctx := context.Background()

	t.Parallel()
	ns := helpers.Namespace(ctx, t, env)

	if env.Cluster().Type() != kind.KindClusterType {
		t.Skip("webhook tests are only available on KIND clusters currently")
	}

	t.Log("creating an extra namespace for testing global consumer credentials validation")
	require.NoError(t, clusters.CreateNamespace(ctx, env.Cluster(), extraWebhookNamespace))
	defer func() {
		if err := env.Cluster().Client().CoreV1().Namespaces().Delete(ctx, extraWebhookNamespace, metav1.DeleteOptions{}); err != nil {
			if !apierrors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
	}()

	ensureAdmissionRegistration(
		ctx,
		t,
		ns.Name,
		"kong-validations",
		[]admregv1.RuleWithOperations{
			{
				Rule: admregv1.Rule{
					APIGroups:   []string{""},
					APIVersions: []string{"v1"},
					Resources:   []string{"secrets"},
				},
				Operations: []admregv1.OperationType{admregv1.Create, admregv1.Update},
			},
			{
				Rule: admregv1.Rule{
					APIGroups:   []string{"configuration.konghq.com"},
					APIVersions: []string{"v1"},
					Resources:   []string{"kongconsumers"},
				},
				Operations: []admregv1.OperationType{admregv1.Create, admregv1.Update},
			},
			{
				Rule: admregv1.Rule{
					APIGroups:   []string{"configuration.konghq.com"},
					APIVersions: []string{"v1alpha1"},
					Resources:   []string{"kongvaults"},
				},
				Operations: []admregv1.OperationType{admregv1.Create, admregv1.Update},
			},
		},
	)
	t.Log("waiting for webhook service to be connective")
	ensureWebhookServiceIsConnective(ctx, t, "kong-validations")

	kongClient, err := clientset.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("creating a large number of consumers on the cluster to verify the performance of the cached client during validation")
	for i := 0; i < highEndConsumerUsageCount; i++ {
		consumerName := fmt.Sprintf("background-noise-consumer-%d", i)
		// create 5 credentials for each consumer
		for j := 0; j < 5; j++ {
			credentialName := fmt.Sprintf("%s-credential-%d", consumerName, j)
			credential := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: credentialName,
				},
				StringData: map[string]string{
					"kongCredType": "basic-auth",
					"username":     credentialName,
					"password":     "testpass",
				},
			}
			_, err := env.Cluster().Client().CoreV1().Secrets(extraWebhookNamespace).Create(ctx, credential, metav1.CreateOptions{})
			require.NoError(t, err)
			defer func() {
				if err := env.Cluster().Client().CoreV1().Secrets(extraWebhookNamespace).Delete(ctx, credentialName, metav1.DeleteOptions{}); err != nil {
					if !apierrors.IsNotFound(err) {
						assert.NoError(t, err)
					}
				}
			}()
		}
	}

	for i := 0; i < highEndConsumerUsageCount; i++ {
		consumerName := fmt.Sprintf("background-noise-consumer-%d", i)

		// create the consumer referencing its credentials
		consumer := &kongv1.KongConsumer{
			ObjectMeta: metav1.ObjectMeta{
				Name: consumerName,
				Annotations: map[string]string{
					annotations.IngressClassKey: consts.IngressClass,
				},
			},
			Username: consumerName,
			CustomID: uuid.NewString(),
		}
		for j := 0; j < 5; j++ {
			credentialName := fmt.Sprintf("%s-credential-%d", consumerName, j)
			consumer.Credentials = append(consumer.Credentials, credentialName)
		}
		assert.Eventually(t, func() bool {
			_, err = kongClient.ConfigurationV1().KongConsumers(extraWebhookNamespace).Create(ctx, consumer, metav1.CreateOptions{})
			if err != nil {
				t.Logf("Failed to create consumer, will retry: %s", err)
			}
			return err == nil
		}, time.Second*10, time.Second*1)
		require.NoError(t, err)
		defer func() {
			if err := kongClient.ConfigurationV1().KongConsumers(extraWebhookNamespace).Delete(ctx, consumerName, metav1.DeleteOptions{}); err != nil {
				if !apierrors.IsNotFound(err) {
					assert.NoError(t, err)
				}
			}
		}()
	}

	t.Log("creating some static credentials in an extra namespace which will be used to test global validation")
	for _, secret := range []*corev1.Secret{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "tuxcreds1",
			},
			StringData: map[string]string{
				"kongCredType": "basic-auth",
				"username":     "tux1",
				"password":     "testpass",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "tuxcreds2",
			},
			StringData: map[string]string{
				"kongCredType": "basic-auth",
				"username":     "tux2",
				"password":     "testpass",
			},
		},
	} {
		secret, err = env.Cluster().Client().CoreV1().Secrets(extraWebhookNamespace).Create(ctx, secret, metav1.CreateOptions{})
		require.NoError(t, err)
		secretName := secret.Name
		defer func() {
			if err := env.Cluster().Client().CoreV1().Secrets(extraWebhookNamespace).Delete(ctx, secretName, metav1.DeleteOptions{}); err != nil {
				if !apierrors.IsNotFound(err) {
					assert.NoError(t, err)
				}
			}
		}()
	}

	t.Log("creating a static consumer in an extra namespace which will be used to test global validation")
	require.NoError(t, err)
	consumer := &kongv1.KongConsumer{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
			Annotations: map[string]string{
				annotations.IngressClassKey: consts.IngressClass,
			},
		},
		Username: "tux",
		CustomID: uuid.NewString(),
		Credentials: []string{
			"tuxcreds1",
			"tuxcreds2",
		},
	}
	consumer, err = kongClient.ConfigurationV1().KongConsumers(extraWebhookNamespace).Create(ctx, consumer, metav1.CreateOptions{})
	require.NoError(t, err)
	defer func() {
		if err := kongClient.ConfigurationV1().KongConsumers(extraWebhookNamespace).Delete(ctx, consumer.Name, metav1.DeleteOptions{}); err != nil {
			if !apierrors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
	}()

	t.Log("testing consumer credentials validation")
	for _, tt := range []struct {
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
						annotations.IngressClassKey: consts.IngressClass,
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
						annotations.IngressClassKey: consts.IngressClass,
					},
				},
				Username:    "electron",
				CustomID:    uuid.NewString(),
				Credentials: []string{"electronscreds"},
			},
			credentials: []*corev1.Secret{{
				ObjectMeta: metav1.ObjectMeta{
					Name: "electronscreds",
				},
				StringData: map[string]string{
					"kongCredType": "basic-auth",
					"username":     "electron",
					"password":     "testpass",
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
						annotations.IngressClassKey: consts.IngressClass,
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
					},
					StringData: map[string]string{
						"kongCredType": "basic-auth",
						"username":     "proton",
						"password":     "testpass",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "protonscreds2",
					},
					StringData: map[string]string{
						"kongCredType": "basic-auth",
						"username":     "electron", // username is unique constrained
						"password":     "testpass", // password is not unique constrained
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
						annotations.IngressClassKey: consts.IngressClass,
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
						annotations.IngressClassKey: consts.IngressClass,
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
					},
					StringData: map[string]string{
						"kongCredType": "basic-auth",
						"username":     "neutron",
						"password":     "testpass",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "neutronscreds2",
					},
					StringData: map[string]string{
						"kongCredType": "basic-auth",
						"username":     "neutron", // username is unique constrained
						"password":     "testpass",
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
						annotations.IngressClassKey: consts.IngressClass,
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
					},
					StringData: map[string]string{
						"kongCredType": "basic-auth",
						"username":     "reasonablehammer",
						"password":     "testpass", // not unique constrained, so even though someone else is using this password this should pass
					},
				},
			},
			wantErr: false,
		},
		{
			name: "a consumer that provides credentials that are in violation of unique constraints globally against other existing consumers should fail validation",
			consumer: &kongv1.KongConsumer{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
					Annotations: map[string]string{
						annotations.IngressClassKey: consts.IngressClass,
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
					},
					StringData: map[string]string{
						"kongCredType": "basic-auth",
						"username":     "tux1", // unique constrained with previous created static consumer credentials
						"password":     "testpass",
					},
				},
			},
			wantErr:        true,
			wantPartialErr: "unique key constraint violated for username",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, credential := range tt.credentials {
				credential, err = env.Cluster().Client().CoreV1().Secrets(ns.Name).Create(ctx, credential, metav1.CreateOptions{})
				require.NoError(t, err)
				credentialName := credential.Name
				defer func() {
					if err := env.Cluster().Client().CoreV1().Secrets(ns.Name).Delete(ctx, credentialName, metav1.DeleteOptions{}); err != nil {
						if !apierrors.IsNotFound(err) {
							assert.NoError(t, err)
						}
					}
				}()
			}

			defer func() {
				if err := kongClient.ConfigurationV1().KongConsumers(ns.Name).Delete(ctx, tt.consumer.Name, metav1.DeleteOptions{}); err != nil {
					if !apierrors.IsNotFound(err) {
						assert.NoError(t, err)
					}
				}
			}()

			consumer, err := kongClient.ConfigurationV1().KongConsumers(ns.Name).Create(ctx, tt.consumer, metav1.CreateOptions{})
			if tt.wantErr {
				require.Error(t, err, fmt.Sprintf("consumer %s should fail to create", consumer.Name))
				assert.True(t, strings.Contains(err.Error(), tt.wantPartialErr),
					"got error string %q, want a superstring of %q", err.Error(), tt.wantPartialErr)
			} else {
				require.NoError(t, err, fmt.Sprintf("consumer %s should create successfully", consumer.Name))
			}
		})
	}

	t.Log("verifying that an invalid credential secret not yet referenced by a KongConsumer fails validation")
	invalidCredential := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "brokenfence",
		},
		StringData: map[string]string{
			"kongCredType": "invalid-auth", // not a valid credential type
			"username":     "brokenfence",
			"password":     "testpass",
		},
	}
	_, err = env.Cluster().Client().CoreV1().Secrets(ns.Name).Create(ctx, invalidCredential, metav1.CreateOptions{})
	require.ErrorContains(t, err, "invalid credential type")

	t.Log("creating a valid credential secret to be referenced by a KongConsumer")
	validCredential, err := env.Cluster().Client().CoreV1().Secrets(ns.Name).Create(ctx, &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "brokenfence",
		},
		StringData: map[string]string{
			"kongCredType": "basic-auth",
			"username":     "brokenfence",
			"password":     "testpass",
		},
	}, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("verifying that valid credentials assigned to a consumer pass validation")
	validConsumerLinkedToValidCredentials := &kongv1.KongConsumer{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
			Annotations: map[string]string{
				annotations.IngressClassKey: consts.IngressClass,
			},
		},
		Username: "brokenfence",
		CustomID: uuid.NewString(),
		Credentials: []string{
			"brokenfence",
		},
	}
	validConsumerLinkedToValidCredentials, err = kongClient.ConfigurationV1().KongConsumers(ns.Name).Create(ctx, validConsumerLinkedToValidCredentials, metav1.CreateOptions{})
	require.NoError(t, err)
	defer func() {
		if err := kongClient.ConfigurationV1().KongConsumers(ns.Name).Delete(ctx, validConsumerLinkedToValidCredentials.Name, metav1.DeleteOptions{}); err != nil {
			if !apierrors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
	}()

	t.Log("verifying that the valid credentials which include a unique-constrained key can be updated in place")
	validCredential.Data["value"] = []byte("newpassword")
	validCredential, err = env.Cluster().Client().CoreV1().Secrets(ns.Name).Update(ctx, validCredential, metav1.UpdateOptions{})
	require.NoError(t, err)

	t.Log("verifying that validation fails if the now referenced and valid credential gets updated to become invalid")
	validCredential.Data["kongCredType"] = []byte("invalid-auth")
	_, err = env.Cluster().Client().CoreV1().Secrets(ns.Name).Update(ctx, validCredential, metav1.UpdateOptions{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid credential type")

	t.Log("verifying that if the referent consumer goes away the validation fails for updates that make the credential invalid")
	require.NoError(t, kongClient.ConfigurationV1().KongConsumers(ns.Name).Delete(ctx, validConsumerLinkedToValidCredentials.Name, metav1.DeleteOptions{}))
	_, err = env.Cluster().Client().CoreV1().Secrets(ns.Name).Update(ctx, validCredential, metav1.UpdateOptions{})
	require.ErrorContains(t, err, "invalid credential type")

	t.Log("verifying that a JWT credential which has keys with missing values fails validation")
	invalidJWTName := uuid.NewString()
	invalidJWT := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: invalidJWTName,
		},
		StringData: map[string]string{
			"kongCredType":   "jwt",
			"algorithm":      "RS256",
			"key":            "",
			"rsa_public_key": "",
			"secret":         "",
		},
	}
	_, err = env.Cluster().Client().CoreV1().Secrets(ns.Name).Create(ctx, invalidJWT, metav1.CreateOptions{})
	require.ErrorContains(t, err, "some fields were invalid due to missing data: rsa_public_key, key, secret")

	t.Log("verifying that the validation fails when secret generates invalid plugin configuration for KongPlugin")
	for _, tt := range []struct {
		name          string
		KongPlugin    *kongv1.KongPlugin
		secretBefore  *corev1.Secret
		secretAfter   *corev1.Secret
		errorOnUpdate bool
		errorContains string
	}{
		{
			name: "should fail the validation if secret used in ConfigFrom of KongPlugin generates invalid plugin configuration",
			KongPlugin: &kongv1.KongPlugin{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: ns.Name,
					Name:      "rate-limiting-invalid-config-from",
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
					Namespace: ns.Name,
					Name:      "conf-secret-invalid-config",
				},
				Data: map[string][]byte{
					"rate-limiting-config": []byte(`{"limit_by":"consumer","policy":"local","minute":5}`),
				},
			},
			secretAfter: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: ns.Name,
					Name:      "conf-secret-invalid-config",
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
			KongPlugin: &kongv1.KongPlugin{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: ns.Name,
					Name:      "rate-limiting-invalid-config-patches",
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
					Namespace: ns.Name,
					Name:      "conf-secret-invalid-field",
				},
				Data: map[string][]byte{
					"rate-limiting-config-minutes": []byte("10"),
				},
			},
			secretAfter: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: ns.Name,
					Name:      "conf-secret-invalid-field",
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
			KongPlugin: &kongv1.KongPlugin{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: ns.Name,
					Name:      "rate-limiting-valid-config",
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
					Namespace: ns.Name,
					Name:      "conf-secret-valid-field",
				},
				Data: map[string][]byte{
					"rate-limiting-config-minutes": []byte(`10`),
				},
			},
			secretAfter: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: ns.Name,
					Name:      "conf-secret-valid-field",
				},
				Data: map[string][]byte{
					"rate-limiting-config-minutes": []byte(`15`),
				},
			},
			errorOnUpdate: false,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, err := env.Cluster().Client().CoreV1().Secrets(ns.Name).Create(ctx, tt.secretBefore, metav1.CreateOptions{})
			require.NoError(t, err)
			_, err = kongClient.ConfigurationV1().KongPlugins(ns.Name).Create(ctx, tt.KongPlugin, metav1.CreateOptions{})
			require.NoError(t, err)
			defer func() {
				err := kongClient.ConfigurationV1().KongPlugins(ns.Name).Delete(ctx, tt.KongPlugin.Name, metav1.DeleteOptions{})
				require.NoError(t, err)
			}()

			_, err = env.Cluster().Client().CoreV1().Secrets(ns.Name).Update(ctx, tt.secretAfter, metav1.UpdateOptions{})
			if tt.errorOnUpdate {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorContains)
			} else {
				require.NoError(t, err)
			}
		})
	}

	t.Log("verifying that the validation fails when secret generates invalid plugin configuration for KongClusterPlugin")
	for _, tt := range []struct {
		name              string
		kongClusterPlugin *kongv1.KongClusterPlugin
		secretBefore      *corev1.Secret
		secretAfter       *corev1.Secret
		errorOnUpdate     bool
		errorContains     string
	}{
		{
			name: "should pass the validation if the secret used in ConfigFrom of KongClusterPlugin generates valid configuration",
			kongClusterPlugin: &kongv1.KongClusterPlugin{
				ObjectMeta: metav1.ObjectMeta{
					Name: "cluster-rate-limiting-valid",
					Annotations: map[string]string{
						annotations.IngressClassKey: consts.IngressClass,
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
					Namespace: ns.Name,
					Name:      "cluster-conf-secret-valid",
				},
				Data: map[string][]byte{
					"rate-limiting-config": []byte(`{"limit_by":"consumer","policy":"local","minute":5}`),
				},
			},
			secretAfter: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: ns.Name,
					Name:      "cluster-conf-secret-valid",
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
						annotations.IngressClassKey: consts.IngressClass,
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
					Namespace: ns.Name,
					Name:      "cluster-conf-secret-invalid",
				},
				Data: map[string][]byte{
					"rate-limiting-config": []byte(`{"limit_by":"consumer","policy":"local","minute":5}`),
				},
			},
			secretAfter: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: ns.Name,
					Name:      "cluster-conf-secret-invalid",
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
						annotations.IngressClassKey: consts.IngressClass,
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
								Namespace: ns.Namespace,
								Secret:    "cluster-conf-secret-valid-patch",
								Key:       "rate-limiting-minute",
							},
						},
					},
				},
			},
			secretBefore: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: ns.Name,
					Name:      "cluster-conf-secret-valid-patch",
				},
				Data: map[string][]byte{
					"rate-limiting-minute": []byte(`5`),
				},
			},
			secretAfter: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: ns.Name,
					Name:      "cluster-conf-secret-valid-patch",
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
						annotations.IngressClassKey: consts.IngressClass,
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
					Namespace: ns.Name,
					Name:      "cluster-conf-secret-invalid-patch",
				},
				Data: map[string][]byte{
					"rate-limiting-minute": []byte(`5`),
				},
			},
			secretAfter: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: ns.Name,
					Name:      "cluster-conf-secret-invalid-patch",
				},
				Data: map[string][]byte{
					"rate-limiting-minute": []byte(`"10"`),
				},
			},
			errorOnUpdate: true,
			errorContains: "Change on secret will generate invalid configuration for KongClusterPlugin",
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, err := env.Cluster().Client().CoreV1().Secrets(ns.Name).Create(ctx, tt.secretBefore, metav1.CreateOptions{})
			require.NoError(t, err)
			_, err = kongClient.ConfigurationV1().KongClusterPlugins().Create(ctx, tt.kongClusterPlugin, metav1.CreateOptions{})
			require.NoError(t, err)
			defer func() {
				err := kongClient.ConfigurationV1().KongClusterPlugins().Delete(ctx, tt.kongClusterPlugin.Name, metav1.DeleteOptions{})
				require.NoError(t, err)
			}()

			_, err = env.Cluster().Client().CoreV1().Secrets(ns.Name).Update(ctx, tt.secretAfter, metav1.UpdateOptions{})
			if tt.errorOnUpdate {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

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
