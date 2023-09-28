package parser

import (
	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/versions"
)

// This file contains unit test functions to test translation failures genreated by parser.

func newResourceFailure(t *testing.T, reason string, objects ...client.Object) failures.ResourceFailure {
	failure, err := failures.NewResourceFailure(reason, objects...)
	require.NoError(t, err)
	return failure
}

func TestTranslationFailureUnsupportedObjectsExpressionRoutes(t *testing.T) {
	testCases := []struct {
		name           string
		objects        store.FakeObjects
		causingObjects []client.Object
	}{
		{
			name: "TLSRoutes in gateway APIs are not supported",
			objects: store.FakeObjects{
				TLSRoutes: []*gatewayapi.TLSRoute{
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
				&gatewayapi.TLSRoute{
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
			parser.kongVersion = versions.ExpressionRouterL4Cutoff

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
