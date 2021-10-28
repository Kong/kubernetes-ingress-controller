//go:build integration_tests
// +build integration_tests

package integration

import (
	"net"
	"strings"
	"testing"
	"time"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/kind"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	admregv1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestValidationWebhook(t *testing.T) {
	t.Parallel()
	ns, cleanup := namespace(t)
	defer cleanup()

	if env.Cluster().Type() != kind.KindClusterType {
		t.Skip("TODO: webhook tests are only supported on KIND based environments right now")
	}

	const webhookSvcName = "validations"
	_, err := env.Cluster().Client().CoreV1().Services(controllerNamespace).Create(ctx, &corev1.Service{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Service"},
		ObjectMeta: metav1.ObjectMeta{Name: webhookSvcName},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:       "default",
					Port:       443,
					TargetPort: intstr.FromInt(49023),
				},
			},
		},
	}, metav1.CreateOptions{})
	require.NoError(t, err, "creating webhook service")

	nodeName := "aaaa"
	_, err = env.Cluster().Client().CoreV1().Endpoints(controllerNamespace).Create(ctx, &corev1.Endpoints{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Endpoints"},
		ObjectMeta: metav1.ObjectMeta{Name: webhookSvcName},
		Subsets: []corev1.EndpointSubset{
			{
				Addresses: []corev1.EndpointAddress{
					{
						IP:       "172.17.0.1",
						NodeName: &nodeName,
					},
				},
				Ports: []corev1.EndpointPort{
					{
						Name:     "default",
						Port:     49023,
						Protocol: corev1.ProtocolTCP,
					},
				},
			},
		},
	}, metav1.CreateOptions{})
	require.NoError(t, err, "creating webhook endpoints")

	fail := admregv1.Fail
	none := admregv1.SideEffectClassNone
	_, err = env.Cluster().Client().AdmissionregistrationV1().ValidatingWebhookConfigurations().Create(ctx,
		&admregv1.ValidatingWebhookConfiguration{
			TypeMeta:   metav1.TypeMeta{APIVersion: "admissionregistration.k8s.io/v1", Kind: "ValidatingWebhookConfiguration"},
			ObjectMeta: metav1.ObjectMeta{Name: "kong-validations"},
			Webhooks: []admregv1.ValidatingWebhook{
				{
					Name:                    "validations.kong.konghq.com",
					FailurePolicy:           &fail,
					SideEffects:             &none,
					AdmissionReviewVersions: []string{"v1beta1", "v1"},
					Rules: []admregv1.RuleWithOperations{
						{
							Rule: admregv1.Rule{
								APIGroups:   []string{""},
								APIVersions: []string{"v1"},
								Resources:   []string{"secrets"},
							},
							Operations: []admregv1.OperationType{admregv1.Create, admregv1.Update},
						},
					},
					ClientConfig: admregv1.WebhookClientConfig{
						Service:  &admregv1.ServiceReference{Namespace: controllerNamespace, Name: webhookSvcName},
						CABundle: []byte(admissionWebhookCert),
					},
				},
			},
		}, metav1.CreateOptions{})
	require.NoError(t, err, "creating webhook config")
	require.Eventually(t, func() bool {
		_, err := net.DialTimeout("tcp", "172.17.0.1:49023", 1*time.Second)
		return err == nil
	}, ingressWait, waitTick, "waiting for the admission service to be up")

	// TODO: flakes were occurring in this test because proxy readiness isn't a consistent gate mechanism
	//       by which to determine readiness for the webhook validation tests. We will follow up on this by
	//       improving these tests, but for now (for speed at the time of writing) we just sleep.
	//
	//       See: https://github.com/Kong/kubernetes-ingress-controller/issues/1442
	time.Sleep(time.Second * 5)

	for _, tt := range []struct {
		name           string
		obj            corev1.Secret
		wantErr        bool
		wantPartialErr string
	}{
		{
			name: "validation passed: non-kong secret",
			obj: corev1.Secret{
				TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Secret"},
				ObjectMeta: metav1.ObjectMeta{Name: "unknown-kongcredtype"},
				StringData: map[string]string{"something": "something"},
			},
		},
		{
			name: "secret of unknown kongCredType",
			obj: corev1.Secret{
				TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Secret"},
				ObjectMeta: metav1.ObjectMeta{Name: "unknown-kongcredtype"},
				StringData: map[string]string{"kongCredType": "nonexistent"},
			},
			wantErr:        true,
			wantPartialErr: "invalid credential type: nonexistent",
		},
		{
			name: "secret with missing fields",
			obj: corev1.Secret{
				TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Secret"},
				ObjectMeta: metav1.ObjectMeta{Name: "basic-auth"},
				StringData: map[string]string{"kongCredType": "basic-auth", "username": "foo"},
			},
			wantErr:        true,
			wantPartialErr: "missing required field(s): password",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := env.Cluster().Client().CoreV1().Secrets(ns.Name).Create(ctx, &tt.obj, metav1.CreateOptions{})
			defer func() {
				if err := env.Cluster().Client().CoreV1().Secrets(ns.Name).Delete(ctx, tt.obj.ObjectMeta.Name, metav1.DeleteOptions{}); err != nil {
					if !errors.IsNotFound(err) {
						assert.NoError(t, err)
					}
				}
			}()
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, strings.Contains(err.Error(), tt.wantPartialErr),
					"got error string %q, want a superstring of %q", err.Error(), tt.wantPartialErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
