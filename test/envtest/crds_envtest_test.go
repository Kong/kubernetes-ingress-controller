//go:build envtest

package envtest

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1alpha1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/v2/test/consts"
)

// TestMissingCRDsDontCrashTheManager ensures that in case of missing CRDs installation in the cluster, specific
// controllers are disabled, this fact is properly logged, and the manager does not crash.
func TestMissingCRDsDontCrashTheManager(t *testing.T) {
	emptyScheme := runtime.NewScheme()
	envcfg := Setup(t, emptyScheme)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	loggerHook := RunManager(ctx, t, envcfg, func(cfg *manager.Config) {
		// Reducing controllers' cache synchronisation timeout in order to trigger the possible sync timeout quicker.
		// It's a regression test for https://github.com/Kong/gateway-operator/issues/326.
		cfg.CacheSyncTimeout = time.Millisecond * 500
	})

	require.Eventually(t, func() bool {
		gvrs := []schema.GroupVersionResource{
			{
				Group:    kongv1beta1.GroupVersion.Group,
				Version:  kongv1beta1.GroupVersion.Version,
				Resource: "udpingresses",
			},
			{
				Group:    kongv1beta1.GroupVersion.Group,
				Version:  kongv1beta1.GroupVersion.Version,
				Resource: "tcpingresses",
			},
			{
				Group:    kongv1.GroupVersion.Group,
				Version:  kongv1.GroupVersion.Version,
				Resource: "kongingresses",
			},
			{
				Group:    kongv1alpha1.GroupVersion.Group,
				Version:  kongv1alpha1.GroupVersion.Version,
				Resource: "ingressclassparameterses",
			},
			{
				Group:    kongv1.GroupVersion.Group,
				Version:  kongv1.GroupVersion.Version,
				Resource: "kongplugins",
			},
			{
				Group:    kongv1.GroupVersion.Group,
				Version:  kongv1.GroupVersion.Version,
				Resource: "kongconsumers",
			},
			{
				Group:    kongv1beta1.GroupVersion.Group,
				Version:  kongv1beta1.GroupVersion.Version,
				Resource: "kongconsumergroups",
			},
			{
				Group:    kongv1.GroupVersion.Group,
				Version:  kongv1.GroupVersion.Version,
				Resource: "kongclusterplugins",
			},
		}

		for _, gvr := range gvrs {
			expectedLog := fmt.Sprintf("Disabling controller for Group=%s, Resource=%s due to missing CRD", gvr.GroupVersion(), gvr.Resource)
			if !lo.ContainsBy(loggerHook.AllEntries(), func(entry *logrus.Entry) bool {
				return strings.Contains(entry.Message, expectedLog)
			}) {
				t.Logf("expected log not found: %s", expectedLog)
				return false
			}
		}
		return true
	}, time.Minute, time.Millisecond*500)
}

func TestCRDValidations(t *testing.T) {
	ctx := context.Background()
	scheme := Scheme(t, WithKong)
	envcfg := Setup(t, scheme, WithInstallKongCRDs(true))
	client := NewControllerClient(t, scheme, envcfg)

	testCases := []struct {
		name     string
		scenario func(ctx context.Context, t *testing.T, ns string)
	}{
		{
			name: "invalid TCPIngress service name",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createFaultyTCPIngress(ctx, t, envcfg, ns, func(ingress *kongv1beta1.TCPIngress) {
					ingress.Spec.Rules[0].Backend.ServiceName = ""
				})

				require.ErrorContains(t, err, "serviceName")
			},
		},
		{
			name: "invalid TCPIngress service port",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createFaultyTCPIngress(ctx, t, envcfg, ns, func(ingress *kongv1beta1.TCPIngress) {
					ingress.Spec.Rules[0].Backend.ServicePort = 0
				})

				require.ErrorContains(t, err, "servicePort")
			},
		},
		{
			name: "invalid TCPIngress rule port",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createFaultyTCPIngress(ctx, t, envcfg, ns, func(ingress *kongv1beta1.TCPIngress) {
					ingress.Spec.Rules[0].Port = 0
				})

				require.ErrorContains(t, err, "spec.rules[0].port")
			},
		},
		{
			name: "invalid UDPIngress service name",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createFaultyUDPIngress(ctx, t, envcfg, ns, func(ingress *kongv1beta1.UDPIngress) {
					ingress.Spec.Rules[0].Backend.ServiceName = ""
				})

				require.ErrorContains(t, err, "serviceName")
			},
		},
		{
			name: "invalid UDPIngress service port",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createFaultyUDPIngress(ctx, t, envcfg, ns, func(ingress *kongv1beta1.UDPIngress) {
					ingress.Spec.Rules[0].Backend.ServicePort = 0
				})

				require.ErrorContains(t, err, "servicePort")
			},
		},
		{
			name: "invalid UDPIngress rule port",
			scenario: func(ctx context.Context, t *testing.T, ns string) {
				err := createFaultyUDPIngress(ctx, t, envcfg, ns, func(ingress *kongv1beta1.UDPIngress) {
					ingress.Spec.Rules[0].Port = 0
				})

				require.ErrorContains(t, err, "spec.rules[0].port")
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ns := CreateNamespace(ctx, t, client)
			tc.scenario(ctx, t, ns.Name)
		})
	}
}

func createFaultyTCPIngress(ctx context.Context, t *testing.T, envcfg *rest.Config, ns string, modifier func(*kongv1beta1.TCPIngress)) error {
	ingress := validTCPIngress()
	modifier(ingress)

	gatewayClient, err := clientset.NewForConfig(envcfg)
	require.NoError(t, err)

	c := gatewayClient.ConfigurationV1beta1().TCPIngresses(ns)
	ingress, err = c.Create(ctx, ingress, metav1.CreateOptions{})
	if !assert.Error(t, err) {
		t.Cleanup(func() { _ = c.Delete(ctx, ingress.Name, metav1.DeleteOptions{}) })
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

func createFaultyUDPIngress(ctx context.Context, t *testing.T, envcfg *rest.Config, ns string, modifier func(ingress *kongv1beta1.UDPIngress)) error {
	ingress := validUDPIngress()
	modifier(ingress)

	gatewayClient, err := clientset.NewForConfig(envcfg)
	require.NoError(t, err)

	c := gatewayClient.ConfigurationV1beta1().UDPIngresses(ns)
	ingress, err = c.Create(ctx, ingress, metav1.CreateOptions{})
	if !assert.Error(t, err) {
		t.Cleanup(func() { _ = c.Delete(ctx, ingress.Name, metav1.DeleteOptions{}) })
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
