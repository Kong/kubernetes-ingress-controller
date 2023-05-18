package parser

import (
	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	knative "knative.dev/networking/pkg/apis/networking/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
)

// This file contains unit test functions to test translation failures genreated by parser.

func TestTranslationFailureUnsupportedObjectsExpressionRoutes(t *testing.T) {
	testCases := []struct {
		name           string
		objects        store.FakeObjects
		causingObjects []client.Object
	}{
		{
			name: "knative.Ingresses are not supported",
			objects: store.FakeObjects{
				KnativeIngresses: []*knative.Ingress{
					{
						TypeMeta: metav1.TypeMeta{
							Kind:       "Ingress",
							APIVersion: knative.SchemeGroupVersion.String(),
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:      "knative-ing-1",
							Namespace: "default",
							Annotations: map[string]string{
								annotations.KnativeIngressClassKey: annotations.DefaultIngressClass,
							},
						},
						Spec: knative.IngressSpec{},
					},
				},
			},
			causingObjects: []client.Object{
				&knative.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "knative-ing-1",
						Namespace: "default",
					},
				},
			},
		},
		{
			name: "TCPIngresses and UDPIngresses are not supported",
			objects: store.FakeObjects{
				TCPIngresses: []*kongv1beta1.TCPIngress{
					{
						TypeMeta: metav1.TypeMeta{
							Kind:       "TCPIngress",
							APIVersion: kongv1beta1.GroupVersion.String(),
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:      "tcpingress-1",
							Namespace: "default",
							Annotations: map[string]string{
								annotations.IngressClassKey: annotations.DefaultIngressClass,
							},
						},
					},
				},
				UDPIngresses: []*kongv1beta1.UDPIngress{
					{
						TypeMeta: metav1.TypeMeta{
							Kind:       "UDPIngress",
							APIVersion: kongv1beta1.GroupVersion.String(),
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:      "udpingress-1",
							Namespace: "default",
							Annotations: map[string]string{
								annotations.IngressClassKey: annotations.DefaultIngressClass,
							},
						},
					},
				},
			},
			causingObjects: []client.Object{
				&kongv1beta1.TCPIngress{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "tcpingress-1",
						Namespace: "default",
					},
				},
				&kongv1beta1.UDPIngress{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "udpingress-1",
						Namespace: "default",
					},
				},
			},
		},
		{
			name: "TCPRoutes, UDPRoutes and TLSRoutes in gateway APIs are not supported",
			objects: store.FakeObjects{
				TCPRoutes: []*gatewayv1alpha2.TCPRoute{
					{
						TypeMeta: metav1.TypeMeta{
							Kind:       "TCPRoute",
							APIVersion: gatewayv1alpha2.GroupVersion.String(),
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:      "tcproute-1",
							Namespace: "default",
						},
					},
				},
				UDPRoutes: []*gatewayv1alpha2.UDPRoute{
					{
						TypeMeta: metav1.TypeMeta{
							Kind:       "UDPRoute",
							APIVersion: gatewayv1alpha2.GroupVersion.String(),
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:      "udproute-1",
							Namespace: "default",
						},
					},
				},
				TLSRoutes: []*gatewayv1alpha2.TLSRoute{
					{
						TypeMeta: metav1.TypeMeta{
							Kind:       "TLSRoute",
							APIVersion: gatewayv1alpha2.GroupVersion.String(),
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:      "tlsroute-1",
							Namespace: "default",
						},
					},
				},
			},
			causingObjects: []client.Object{
				&gatewayv1alpha2.TCPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "tcproute-1",
						Namespace: "default",
					},
				},
				&gatewayv1alpha2.UDPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "udproute-1",
						Namespace: "default",
					},
				},
				&gatewayv1alpha2.TLSRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "tlsroute-1",
						Namespace: "default",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			storer, err := store.NewFakeStore(tc.objects)
			require.NoError(t, err)

			parser := mustNewParser(t, storer)
			parser.featureFlags.ExpressionRoutes = true

			result := parser.BuildKongConfig()
			t.Log(result.TranslationFailures)
			for _, object := range tc.causingObjects {
				require.True(t, lo.ContainsBy(result.TranslationFailures, func(f failures.ResourceFailure) bool {
					msg := f.Message()
					if !strings.Contains(msg, "not supported when expression routes enabled") {
						return false
					}

					causingObjects := f.CausingObjects()
					if len(causingObjects) != 1 {
						return false
					}
					causingObject := causingObjects[0]
					return object.GetNamespace() == causingObject.GetNamespace() &&
						object.GetName() == causingObject.GetName()
				}))
			}
		})

	}
}
