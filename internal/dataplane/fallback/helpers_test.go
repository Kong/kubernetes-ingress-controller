package fallback_test

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
	incubatorv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/incubator/v1alpha1"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers"
)

const (
	testNamespace = "test-namespace"
)

func testIngressClass(t *testing.T, name string) *netv1.IngressClass {
	return helpers.WithTypeMeta(t, &netv1.IngressClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	})
}

func testService(t *testing.T, name string) *corev1.Service {
	return helpers.WithTypeMeta(t, &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
		},
	})
}

func testKongServiceFacade(t *testing.T, name string) *incubatorv1alpha1.KongServiceFacade {
	return helpers.WithTypeMeta(t, &incubatorv1alpha1.KongServiceFacade{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
		},
	})
}

func testKongPlugin(t *testing.T, name string) *kongv1.KongPlugin {
	return helpers.WithTypeMeta(t, &kongv1.KongPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
		},
	})
}

func testKongClusterPlugin(t *testing.T, name string) *kongv1.KongClusterPlugin {
	return helpers.WithTypeMeta(t, &kongv1.KongClusterPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
		},
	})
}
