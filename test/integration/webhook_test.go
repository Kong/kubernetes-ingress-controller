//go:build integration_tests
// +build integration_tests

package integration

import (
	"context"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/kind"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	admregv1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	testutils "github.com/kong/kubernetes-ingress-controller/v2/internal/util/test"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/v2/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
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

	closer, err := ensureAdmissionRegistration(ctx,
		"kong-validations-consumer",
		[]admregv1.RuleWithOperations{
			{
				Rule: admregv1.Rule{
					APIGroups:   []string{""},
					APIVersions: []string{"v1"},
					Resources:   []string{"secrets"},
				},
				Operations: []admregv1.OperationType{admregv1.Update},
			},
			{
				Rule: admregv1.Rule{
					APIGroups:   []string{"configuration.konghq.com"},
					APIVersions: []string{"v1"},
					Resources:   []string{"kongconsumers"},
				},
				Operations: []admregv1.OperationType{admregv1.Create, admregv1.Update},
			},
		},
	)
	assert.NoError(t, err, "creating webhook config")
	defer func() {
		assert.NoError(t, closer())
	}()

	err = waitForWebhookServiceConnective(ctx, "kong-validations-consumer")
	require.NoError(t, err)

	t.Log("creating a large number of consumers on the cluster to verify the performance of the cached client during validation")
	kongClient, err := clientset.NewForConfig(env.Cluster().Config())
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
				t.Logf("failed to create consumer, will retry: %s", err)
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
			name: "a consumer with an invalid credential type should fail validation",
			consumer: &kongv1.KongConsumer{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
					Annotations: map[string]string{
						annotations.IngressClassKey: consts.IngressClass,
					},
				},
				Username: "junklawnmower",
				CustomID: uuid.NewString(),
				Credentials: []string{
					"junklawnmowercreds",
				},
			},
			credentials: []*corev1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "junklawnmowercreds",
					},
					StringData: map[string]string{
						"kongCredType": "invalid-auth",
						"username":     "junklawnmower",
						"password":     "testpass",
					},
				},
			},
			wantErr:        true,
			wantPartialErr: "invalid credential type",
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
		{
			name: "secret with missing fields",
			consumer: &kongv1.KongConsumer{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
					Annotations: map[string]string{
						annotations.IngressClassKey: consts.IngressClass,
					},
				},
				Username: "missingpassword",
				CustomID: uuid.NewString(),
				Credentials: []string{
					"basic-auth-with-missing-fields",
				},
			},
			credentials: []*corev1.Secret{
				{
					TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Secret"},
					ObjectMeta: metav1.ObjectMeta{Name: "basic-auth-with-missing-fields"},
					StringData: map[string]string{"kongCredType": "basic-auth", "username": "foo"},
				},
			},
			wantErr:        true,
			wantPartialErr: "missing required field(s): password",
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

	t.Log("verifying that an invalid credential secret not yet referenced by a KongConsumer is not validated")
	invalidCredential := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "brokenfence",
		},
		StringData: map[string]string{
			"kongCredType": "invalid-auth", // not a valid credential type, but wont be validated until referenced by consumer
			"username":     "brokenfence",
			"password":     "testpass",
		},
	}
	invalidCredential, err = env.Cluster().Client().CoreV1().Secrets(ns.Name).Create(ctx, invalidCredential, metav1.CreateOptions{})
	require.NoError(t, err)
	defer func() {
		if err := env.Cluster().Client().CoreV1().Secrets(ns.Name).Delete(ctx, invalidCredential.Name, metav1.DeleteOptions{}); err != nil {
			if !apierrors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
	}()

	t.Log("an existing invalid credential that becomes referenced by a consumer fails consumer validation")
	validConsumerLinkedToInvalidCredentials := &kongv1.KongConsumer{
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
	_, err = kongClient.ConfigurationV1().KongConsumers(ns.Name).Create(ctx, validConsumerLinkedToInvalidCredentials, metav1.CreateOptions{})
	require.Error(t, err, "a consumer that references an invalid credential can not be created")
	require.Contains(t, err.Error(), "invalid credential type")
	defer func() {
		if err := kongClient.ConfigurationV1().KongConsumers(ns.Name).Delete(ctx, validConsumerLinkedToInvalidCredentials.Name, metav1.DeleteOptions{}); err != nil {
			if !apierrors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
	}()

	t.Log("fixing the invalid credentials")
	invalidCredential.Data["kongCredType"] = []byte("basic-auth")
	validCredential, err := env.Cluster().Client().CoreV1().Secrets(ns.Name).Update(ctx, invalidCredential, metav1.UpdateOptions{})
	require.NoError(t, err)

	t.Log("verifying that now that the credentials are fixed the consumer passes validation")
	_, err = kongClient.ConfigurationV1().KongConsumers(ns.Name).Create(ctx, validConsumerLinkedToInvalidCredentials, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("verifying that the valid credentials which include a unique-constrained key can be updated in place")
	validCredential.Data["value"] = []byte("newpassword")
	validCredential, err = env.Cluster().Client().CoreV1().Secrets(ns.Name).Update(ctx, validCredential, metav1.UpdateOptions{})
	require.NoError(t, err)

	t.Log("verifying that validation fails if the now referenced and valid credential gets updated to become invalid")
	validCredential.Data["kongCredType"] = []byte("invalid-auth")
	_, err = env.Cluster().Client().CoreV1().Secrets(ns.Name).Update(ctx, validCredential, metav1.UpdateOptions{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid credential type")

	t.Log("verifying that if the referent consumer goes away the validation passes for updates that would make the credential invalid")
	require.NoError(t, kongClient.ConfigurationV1().KongConsumers(ns.Name).Delete(ctx, validConsumerLinkedToInvalidCredentials.Name, metav1.DeleteOptions{}))
	_, err = env.Cluster().Client().CoreV1().Secrets(ns.Name).Update(ctx, validCredential, metav1.UpdateOptions{})
	require.NoError(t, err)

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
	require.NoError(t, err)
	jwtConsumer := &kongv1.KongConsumer{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
			Annotations: map[string]string{
				annotations.IngressClassKey: consts.IngressClass,
			},
		},
		Username: "bad-jwt-consumer",
		CustomID: uuid.NewString(),
		Credentials: []string{
			invalidJWTName,
		},
	}
	_, err = kongClient.ConfigurationV1().KongConsumers(ns.Name).Create(ctx, jwtConsumer, metav1.CreateOptions{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "some fields were invalid due to missing data: rsa_public_key, key, secret")
}

func ensureWebhookService(ctx context.Context, name string) (func() error, error) {
	validationsService, err := env.Cluster().Client().CoreV1().Services(consts.ControllerNamespace).Create(ctx, &corev1.Service{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Service"},
		ObjectMeta: metav1.ObjectMeta{Name: name},
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
	if err != nil {
		return nil, fmt.Errorf("creating webhook service: %w", err)
	}

	nodeName := "aaaa"
	endpoints, err := env.Cluster().Client().CoreV1().Endpoints(consts.ControllerNamespace).Create(ctx, &corev1.Endpoints{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Endpoints"},
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Subsets: []corev1.EndpointSubset{
			{
				Addresses: []corev1.EndpointAddress{
					{
						IP:       testutils.AdmissionWebhookListenHost,
						NodeName: &nodeName,
					},
				},
				Ports: []corev1.EndpointPort{
					{
						Name:     "default",
						Port:     testutils.AdmissionWebhookListenPort,
						Protocol: corev1.ProtocolTCP,
					},
				},
			},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("creating webhook endpoints: %w", err)
	}

	closer := func() error {
		if err := env.Cluster().Client().CoreV1().Services(consts.ControllerNamespace).Delete(ctx, validationsService.Name, metav1.DeleteOptions{}); err != nil && !apierrors.IsNotFound(err) {
			return err
		}

		if err := env.Cluster().Client().CoreV1().Endpoints(consts.ControllerNamespace).Delete(ctx, endpoints.Name, metav1.DeleteOptions{}); err != nil && !apierrors.IsNotFound(err) {
			return err
		}
		return nil
	}

	return closer, nil
}

func waitForWebhookServiceConnective(ctx context.Context, configResourceName string) error {
	svcName := fmt.Sprintf("webhook-%s", configResourceName)
	svcPort := 443
	waitCtx, cancel := context.WithTimeout(ctx, ingressWait)
	defer cancel()

	return waitForConnectionOnServicePort(waitCtx, env.Cluster().Client(), consts.ControllerNamespace, svcName, svcPort, 10*time.Second)
}

func ensureAdmissionRegistration(ctx context.Context, configResourceName string, rules []admregv1.RuleWithOperations) (func() error, error) {
	svcName := fmt.Sprintf("webhook-%s", configResourceName)
	svcCloser, err := ensureWebhookService(ctx, svcName)
	if err != nil {
		return nil, err
	}

	fail := admregv1.Fail
	none := admregv1.SideEffectClassNone
	webhook, err := env.Cluster().Client().AdmissionregistrationV1().ValidatingWebhookConfigurations().Create(ctx,
		&admregv1.ValidatingWebhookConfiguration{
			TypeMeta:   metav1.TypeMeta{APIVersion: "admissionregistration.k8s.io/v1", Kind: "ValidatingWebhookConfiguration"},
			ObjectMeta: metav1.ObjectMeta{Name: configResourceName},
			Webhooks: []admregv1.ValidatingWebhook{
				{
					Name:                    "validations.kong.konghq.com",
					FailurePolicy:           &fail,
					SideEffects:             &none,
					AdmissionReviewVersions: []string{"v1beta1", "v1"},
					Rules:                   rules,
					ClientConfig: admregv1.WebhookClientConfig{
						Service:  &admregv1.ServiceReference{Namespace: consts.ControllerNamespace, Name: svcName},
						CABundle: []byte(testutils.KongSystemServiceCert),
					},
				},
			},
		}, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	closer := func() error {
		if err := env.Cluster().Client().AdmissionregistrationV1().ValidatingWebhookConfigurations().Delete(ctx, webhook.Name, metav1.DeleteOptions{}); err != nil && !apierrors.IsNotFound(err) {
			return err
		}
		return svcCloser()
	}

	return closer, nil
}

// TODO: Below helper functions were copied from https://github.com/Kong/kubernetes-testing-framework/pull/643
// and should be removed once a new version of KTF (with that PR) will be released and it will be possible
// to upgrade it in KIC. Read more about required steps in https://github.com/Kong/kubernetes-ingress-controller/issues/3981.

// waitForServiceLoadBalancerAddress waits for a service provided by
// namespace/name to have an ingress IP or Host provisioned and returns that
// address. This function will throw an error if the service gets provisioned
// more than a single address, that is not supported. The context provided
// should have a timeout associated with it or you're going to have a bad time.
func waitForServiceLoadBalancerAddress(ctx context.Context, c kubernetes.Interface, namespace, name string) (string, bool, error) {
	for {
		select {
		case <-ctx.Done():
			return "", false, fmt.Errorf("context completed while waiting for loadbalancer service to provision: %w", ctx.Err())
		default:
			// retrieve a fresh copy of the service
			service, err := c.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
			if err != nil {
				return "", false, fmt.Errorf("error while trying to retrieve registry service: %w", err)
			}
			lbing := service.Status.LoadBalancer.Ingress

			// don't support services which have multiple addresses
			if len(lbing) > 1 {
				return "", false, fmt.Errorf("services with more than one load balancer address are not supported (found %d)", len(lbing))
			}

			// verify whether the loadbalancer details are provisioned
			if len(lbing) == 1 {
				for _, ing := range lbing {
					if ing.Hostname != "" {
						return ing.Hostname, false, nil
					}
					if ing.IP != "" {
						return ing.IP, true, nil
					}
				}
			}
		}
	}
}

// waitForConnectionOnServicePort waits until it can make successful TCP connections
// to a service (provided by namespace/name). This will temporarily create a LoadBalancer
// type Service to allow connections to the Service and port from outside the cluster while
// the connection attempts are made using the LoadBalancer public address.
func waitForConnectionOnServicePort(ctx context.Context, c kubernetes.Interface, namespace, name string, port int, dialTimeout time.Duration) error {
	svcClient := c.CoreV1().Services(namespace)
	service, err := svcClient.Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	const correspondingSvcNameLabel = "corresponding-service"
	lbServiceName := "templb-" + name
	tempLoadBalancer := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      lbServiceName,
			Labels: map[string]string{
				correspondingSvcNameLabel: name,
			},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeLoadBalancer,
			// Copy the selector and ports of the service to check.
			Selector: service.Spec.Selector,
			Ports:    service.Spec.Ports,
		},
	}

	// Empty selector, we should create the endpoints separately.
	// If the target service does not have a selector, it usually means that
	// the endpoints of the target server is manually created, but not chosen from pods by labels in selector.
	// so we need to manually create the same endpoints as the target service has here.
	if len(service.Spec.Selector) == 0 {
		epsClient := c.DiscoveryV1().EndpointSlices(namespace)
		endpointSlices, err := epsClient.List(
			ctx, metav1.ListOptions{LabelSelector: discoveryv1.LabelServiceName + "=" + name},
		)
		if err != nil {
			return err
		}

		// Recreate EndpointSlices for the lb service with proper metadata.
		tempEndpointSlices := endpointSlices.DeepCopy().Items
		for i := range tempEndpointSlices {
			epsName := fmt.Sprintf("%s-%d", lbServiceName, i)
			tempEndpointSlices[i].ObjectMeta = metav1.ObjectMeta{
				Namespace: namespace,
				Name:      epsName,
				Labels: map[string]string{
					discoveryv1.LabelServiceName: lbServiceName, // Maps EndpointSlice to Service.
					correspondingSvcNameLabel:    name,
				},
			}
			if _, err = epsClient.Create(ctx, &tempEndpointSlices[i], metav1.CreateOptions{}); err != nil {
				return err
			}
			// For each successfully created temporary EndpointSlice ensure deletion on return from the function.
			defer func(epsName string) {
				err := epsClient.Delete(ctx, epsName, metav1.DeleteOptions{})
				if err != nil && !apierrors.IsNotFound(err) {
					fmt.Printf("failed to delete EndpointSlice %s/%s after testing, error %v\n",
						namespace, epsName, err,
					)
				}
			}(epsName)
		}
	}

	if _, err = svcClient.Create(ctx, tempLoadBalancer, metav1.CreateOptions{}); err != nil {
		return err
	}
	defer func() {
		err := svcClient.Delete(ctx, lbServiceName, metav1.DeleteOptions{})
		if err != nil && !apierrors.IsNotFound(err) {
			fmt.Printf("failed to delete service %s/%s after testing, error %v\n",
				namespace, lbServiceName, err)
		}
	}()

	ip, _, err := waitForServiceLoadBalancerAddress(ctx, c, namespace, lbServiceName)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(time.Second)
	address := fmt.Sprintf("%s:%d", ip, port)
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context completed or dialTimeout reached while waiting for %s:%d to be connected", ip, port)
		case <-ticker.C:
			dialer := &net.Dialer{Timeout: dialTimeout}
			if _, err := dialer.DialContext(ctx, "tcp", address); err == nil {
				return nil
			}
		}
	}
}
