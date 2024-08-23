package configuration

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

func TestIndexIngressesOnBackendServiceName(t *testing.T) {
	testCases := []struct {
		name            string
		object          client.Object
		expectedIndexes []string
	}{
		{
			name:            "Object not Ingress should return empty index",
			object:          &corev1.Service{},
			expectedIndexes: []string{},
		},
		{
			name: "Ingress with single backend",
			object: &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "test-ns",
					Name:      "test-ingress",
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{Name: "svc", Port: netv1.ServiceBackendPort{Number: 80}},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedIndexes: []string{"test-ns/svc"},
		},
		{
			name: "Ingress with multiple backends and multiple rules",
			object: &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "test-ns",
					Name:      "test-ingress",
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							Host: "foo.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/foo",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{Name: "svc-1", Port: netv1.ServiceBackendPort{Number: 80}},
											},
										},
										{
											Path: "/bar",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{Name: "svc-2", Port: netv1.ServiceBackendPort{Number: 80}},
											},
										},
									},
								},
							},
						},
						{
							Host: "bar.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/foo",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{Name: "svc-4", Port: netv1.ServiceBackendPort{Number: 80}},
											},
										},
										{
											Path: "/bar",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{Name: "svc-3", Port: netv1.ServiceBackendPort{Number: 80}},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedIndexes: []string{"test-ns/svc-1", "test-ns/svc-2", "test-ns/svc-3", "test-ns/svc-4"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			indexes := indexIngressesOnBackendServiceName(tc.object)
			require.ElementsMatch(t, tc.expectedIndexes, indexes, "Should return expected indexes")
		})
	}
}

func TestIndexRoutesOnBackendRefServiceFacadeName(t *testing.T) {
	testCases := []struct {
		name            string
		object          client.Object
		expectedIndexes []string
	}{
		{
			name:            "Objects not HTTPRoute should return empty index",
			object:          &netv1.Ingress{},
			expectedIndexes: []string{},
		},
		{
			name: "HTTPRoute with ServiceFacade backendRef and non-ServiceFacade backendRef",
			object: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "test-ns",
					Name:      "test-httproute",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					Rules: []gatewayapi.HTTPRouteRule{
						{
							BackendRefs: []gatewayapi.HTTPBackendRef{
								{
									BackendRef: gatewayapi.BackendRef{
										BackendObjectReference: gatewayapi.BackendObjectReference{
											Name: gatewayapi.ObjectName("svc-1"),
										},
									},
								},
								{
									BackendRef: gatewayapi.BackendRef{
										BackendObjectReference: gatewayapi.BackendObjectReference{
											Group: lo.ToPtr(gatewayapi.Group("incubator.ingress-controller.konghq.com")),
											Kind:  lo.ToPtr(gatewayapi.Kind("KongServiceFacade")),
											Name:  gatewayapi.ObjectName("service-facade-1"),
										},
									},
								},
								{
									BackendRef: gatewayapi.BackendRef{
										BackendObjectReference: gatewayapi.BackendObjectReference{
											Group:     lo.ToPtr(gatewayapi.Group("incubator.ingress-controller.konghq.com")),
											Kind:      lo.ToPtr(gatewayapi.Kind("KongServiceFacade")),
											Namespace: lo.ToPtr(gatewayapi.Namespace("another-ns")),
											Name:      gatewayapi.ObjectName("service-facade-1"),
										},
									},
								},
							},
						},
					},
				},
			},
			expectedIndexes: []string{"test-ns/service-facade-1", "another-ns/service-facade-1"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			indexes := indexRoutesOnBackendRefServiceFacadeName(tc.object)
			require.ElementsMatch(t, tc.expectedIndexes, indexes, "Should return expected indexes")
		})
	}
}
