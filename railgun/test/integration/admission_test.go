//+build integration_tests

package integration

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	admregv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const defaultNs = "default"

func TestValidateCredential(t *testing.T) {
	ctx := context.Background()

	const webhookSvcName = "webhook-svc"
	_, err := cluster.Client().CoreV1().Services(controllerNamespace).Create(ctx, &corev1.Service{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Service"},
		ObjectMeta: metav1.ObjectMeta{Name: webhookSvcName},
		Spec:       corev1.ServiceSpec{
			// XXX finish this
		},
	}, metav1.CreateOptions{})
	_, err := cluster.Client().AdmissionregistrationV1beta1().ValidatingWebhookConfigurations().Create(ctx,
		admregv1beta1.ValidatingWebhookConfiguration{
			TypeMeta:   metav1.TypeMeta{APIVersion: "admissionregistration.k8s.io/v1", Kind: "ValidatingWebhookConfiguration"},
			ObjectMeta: metav1.ObjectMeta{Name: "kong-validations"},
			Webhooks: admregv1beta1.ValidatingWebhook{
				Name:                    "validations.kong.konghq.com",
				FailurePolicy:           &admregv1beta1.Fail,
				SideEffects:             &admregv1beta1.SideEffectClassNone,
				AdmissionReviewVersions: []string{"v1beta1", "v1"},
				Rules: []admregv1beta1.RuleWithOperations{
					{
						Rule: admregv1beta1.Rule{
							APIGroups:   []string{""},
							APIVersions: []string{"v1"},
							Resources:   []string{"secrets"},
						},
						Operations: []admregv1beta1.OperationType{admregv1beta1.Create, admregv1beta1.Update},
					},
				},
				ClientConfig: admregv1beta1.WebhookClientConfig{
					Service: &admregv1beta1.ServiceReference{Namespace: controllerNamespace, Name: webhookSvcName},
				},
			},
		}, metav1.CreateOptions{})

	for _, tt := range []struct {
		name           string
		obj            corev1.Secret
		wantErr        bool
		wantPartialErr string
	}{
		{
			name: "secret of unknown kongCredType",
			obj: corev1.Secret{
				TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Secret"},
				ObjectMeta: metav1.ObjectMeta{Name: "unknown-kongcredtype"},
				StringData: map[string]string{"kongCredType": "nonexistent"},
			},
			wantErr:        true,
			wantPartialErr: "ukulele",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := cluster.Client().CoreV1().Secrets(defaultNs).Create(ctx, &tt.obj, metav1.CreateOptions{})
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
