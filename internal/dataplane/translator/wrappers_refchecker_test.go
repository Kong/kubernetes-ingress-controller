package translator

import (
	"testing"

	"github.com/go-logr/logr"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

func TestRefChecker_IsRefAllowedByGrant(t *testing.T) {
	testcases := []struct {
		name        string
		route       *gatewayapi.HTTPRoute
		backendRef  gatewayapi.BackendRef
		allowedRefs map[gatewayapi.Namespace][]gatewayapi.ReferenceGrantTo
		expected    bool
	}{
		{
			name: "allowed by grant",
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "example-namespace",
					Name:      "example-name",
				},
			},
			backendRef: gatewayapi.BackendRef{
				BackendObjectReference: gatewayapi.BackendObjectReference{
					Kind:  lo.ToPtr(gatewayapi.Kind("example-kind")),
					Group: lo.ToPtr(gatewayapi.Group("example-group")),
					Name:  "example-name",
				},
			},
			allowedRefs: map[gatewayapi.Namespace][]gatewayapi.ReferenceGrantTo{
				"example-namespace": {
					{
						Kind:  gatewayapi.Kind("example-kind"),
						Group: gatewayapi.Group("example-group"),
						Name:  lo.ToPtr(gatewayapi.ObjectName("example-name")),
					},
				},
			},
			expected: true,
		},
		{
			name: "allowed by grant with namespace",
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "example-namespace",
					Name:      "example-name",
				},
			},
			backendRef: gatewayapi.BackendRef{
				BackendObjectReference: gatewayapi.BackendObjectReference{
					Kind:      lo.ToPtr(gatewayapi.Kind("example-kind")),
					Group:     lo.ToPtr(gatewayapi.Group("example-group")),
					Name:      "example-name",
					Namespace: lo.ToPtr(gatewayapi.Namespace("example-namespace")),
				},
			},
			allowedRefs: map[gatewayapi.Namespace][]gatewayapi.ReferenceGrantTo{
				"example-namespace": {
					{
						Kind:  gatewayapi.Kind("example-kind"),
						Group: gatewayapi.Group("example-group"),
						Name:  lo.ToPtr(gatewayapi.ObjectName("example-name")),
					},
				},
			},
			expected: true,
		},
		{
			name: "allowed because backendRef and route use the same namespace",
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "example-namespace",
					Name:      "example-name",
				},
			},
			backendRef: gatewayapi.BackendRef{
				BackendObjectReference: gatewayapi.BackendObjectReference{
					Kind:      lo.ToPtr(gatewayapi.Kind("example-kind")),
					Group:     lo.ToPtr(gatewayapi.Group("example-group")),
					Name:      "example-name",
					Namespace: lo.ToPtr(gatewayapi.Namespace("example-namespace")),
				},
			},
			expected: true,
		},
		{
			name: "not allowed when not in grant and using different namespace",
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "example-namespace-2",
					Name:      "example-name",
				},
			},
			backendRef: gatewayapi.BackendRef{
				BackendObjectReference: gatewayapi.BackendObjectReference{
					Kind:      lo.ToPtr(gatewayapi.Kind("example-kind")),
					Group:     lo.ToPtr(gatewayapi.Group("example-group")),
					Name:      "example-name",
					Namespace: lo.ToPtr(gatewayapi.Namespace("example-namespace")),
				},
			},
			expected: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rc := gatewayapi.NewRefCheckerForRoute(logr.Discard(), tc.route, tc.backendRef)
			result := rc.IsRefAllowedByGrant(tc.allowedRefs)
			require.Equal(t, tc.expected, result)
		})
	}
}
