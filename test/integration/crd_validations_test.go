//go:build integration_tests
// +build integration_tests

package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/v2/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
)

// TestCRDValidations ensures that CRD validations we expect to exist are properly disallowing faulty objects to be created.
func TestCRDValidations(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name          string
		scenario      func(t *testing.T, cleaner *clusters.Cleaner, ns string)
		expectedError string
	}{
		{
			name: "invalid TCPIngress service name",
			scenario: func(t *testing.T, cleaner *clusters.Cleaner, ns string) {
				err := createFaultyTCPIngress(t, cleaner, ns, func(ingress *kongv1beta1.TCPIngress) {
					ingress.Spec.Rules[0].Backend.ServiceName = ""
				})

				require.ErrorContains(t, err, "serviceName")
			},
		},
		{
			name: "invalid TCPIngress service port",
			scenario: func(t *testing.T, cleaner *clusters.Cleaner, ns string) {
				err := createFaultyTCPIngress(t, cleaner, ns, func(ingress *kongv1beta1.TCPIngress) {
					ingress.Spec.Rules[0].Backend.ServicePort = 0
				})

				require.ErrorContains(t, err, "servicePort")
			},
		},
		{
			name: "invalid TCPIngress rule port",
			scenario: func(t *testing.T, cleaner *clusters.Cleaner, ns string) {
				err := createFaultyTCPIngress(t, cleaner, ns, func(ingress *kongv1beta1.TCPIngress) {
					ingress.Spec.Rules[0].Port = 0
				})

				require.ErrorContains(t, err, "spec.rules[0].port")
			},
		},
		{
			name: "invalid UDPIngress service name",
			scenario: func(t *testing.T, cleaner *clusters.Cleaner, ns string) {
				err := createFaultyUDPIngress(t, cleaner, ns, func(ingress *kongv1beta1.UDPIngress) {
					ingress.Spec.Rules[0].Backend.ServiceName = ""
				})

				require.ErrorContains(t, err, "serviceName")
			},
		},
		{
			name: "invalid UDPIngress service port",
			scenario: func(t *testing.T, cleaner *clusters.Cleaner, ns string) {
				err := createFaultyUDPIngress(t, cleaner, ns, func(ingress *kongv1beta1.UDPIngress) {
					ingress.Spec.Rules[0].Backend.ServicePort = 0
				})

				require.ErrorContains(t, err, "servicePort")
			},
		},
		{
			name: "invalid UDPIngress rule port",
			scenario: func(t *testing.T, cleaner *clusters.Cleaner, ns string) {
				err := createFaultyUDPIngress(t, cleaner, ns, func(ingress *kongv1beta1.UDPIngress) {
					ingress.Spec.Rules[0].Port = 0
				})

				require.ErrorContains(t, err, "spec.rules[0].port")
			},
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ns, cleaner := helpers.Setup(ctx, t, env)

			tt.scenario(t, cleaner, ns.GetName())
		})
	}
}

func createFaultyTCPIngress(t *testing.T, cleaner *clusters.Cleaner, ns string, modifier func(*kongv1beta1.TCPIngress)) error {
	ingress := validTCPIngress()
	modifier(ingress)

	gatewayClient, err := clientset.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	ingress, err = gatewayClient.ConfigurationV1beta1().TCPIngresses(ns).Create(context.Background(), ingress, metav1.CreateOptions{})
	if !assert.Error(t, err) {
		cleaner.Add(ingress)
	}
	return err
}

func validTCPIngress() *kongv1beta1.TCPIngress {
	return &kongv1beta1.TCPIngress{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
			Annotations: map[string]string{
				annotations.IngressClassKey: consts.IngressClass,
			},
		},
		Spec: kongv1beta1.TCPIngressSpec{
			Rules: []kongv1beta1.IngressRule{
				{
					Port: 80,
					Backend: kongv1beta1.IngressBackend{
						ServiceName: "service-name",
						ServicePort: 80,
					},
				},
			},
		},
	}
}

func createFaultyUDPIngress(t *testing.T, cleaner *clusters.Cleaner, ns string, modifier func(ingress *kongv1beta1.UDPIngress)) error {
	ingress := validUDPIngress()
	modifier(ingress)

	gatewayClient, err := clientset.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	ingress, err = gatewayClient.ConfigurationV1beta1().UDPIngresses(ns).Create(context.Background(), ingress, metav1.CreateOptions{})
	if !assert.Error(t, err) {
		cleaner.Add(ingress)
	}
	return err
}

func validUDPIngress() *kongv1beta1.UDPIngress {
	return &kongv1beta1.UDPIngress{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
			Annotations: map[string]string{
				annotations.IngressClassKey: consts.IngressClass,
			},
		},
		Spec: kongv1beta1.UDPIngressSpec{
			Rules: []kongv1beta1.UDPIngressRule{
				{
					Port: 80,
					Backend: kongv1beta1.IngressBackend{
						ServiceName: "service-name",
						ServicePort: 80,
					},
				},
			},
		},
	}
}
