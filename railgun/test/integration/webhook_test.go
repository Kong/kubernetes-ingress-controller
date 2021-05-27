//+build integration_tests

package integration

import (
	"context"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	admregv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const defaultNs = "default"

func TestValidationWebhook(t *testing.T) {
	if useLegacyKIC() {
		t.Skip("not testing validation webhook for KIC 1.x")
	}
	ctx := context.Background()

	const webhookSvcName = "validations"
	_, err := cluster.Client().CoreV1().Services(controllerNamespace).Create(ctx, &corev1.Service{
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
	assert.NoError(t, err, "creating webhook service")

	nodeName := "aaaa"
	_, err = cluster.Client().CoreV1().Endpoints(controllerNamespace).Create(ctx, &corev1.Endpoints{
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
	assert.NoError(t, err, "creating webhook endpoints")

	fail := admregv1beta1.Fail
	none := admregv1beta1.SideEffectClassNone
	_, err = cluster.Client().AdmissionregistrationV1beta1().ValidatingWebhookConfigurations().Create(ctx,
		&admregv1beta1.ValidatingWebhookConfiguration{
			TypeMeta:   metav1.TypeMeta{APIVersion: "admissionregistration.k8s.io/v1", Kind: "ValidatingWebhookConfiguration"},
			ObjectMeta: metav1.ObjectMeta{Name: "kong-validations"},
			Webhooks: []admregv1beta1.ValidatingWebhook{
				{
					Name:                    "validations.kong.konghq.com",
					FailurePolicy:           &fail,
					SideEffects:             &none,
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
						Service:  &admregv1beta1.ServiceReference{Namespace: controllerNamespace, Name: webhookSvcName},
						CABundle: []byte(admissionWebhookCert),
					},
				},
			},
		}, metav1.CreateOptions{})
	assert.NoError(t, err, "creating webhook config")

	t.Log("waiting for proxy ready")
	_ = proxyReady()
	t.Log("waiting for proxy ready done")
	assert.Eventually(t, func() bool {
		_, err := net.DialTimeout("tcp", "172.17.0.1:49023", 1*time.Second)
		return err == nil
	}, ingressWait, waitTick, "waiting for the admission service to be up")

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
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := cluster.Client().CoreV1().Secrets(defaultNs).Create(ctx, &tt.obj, metav1.CreateOptions{})
			defer cluster.Client().CoreV1().Secrets(defaultNs).Delete(ctx, tt.obj.ObjectMeta.Name, metav1.DeleteOptions{})
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
