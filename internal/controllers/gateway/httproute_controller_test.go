package gateway

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

func TestEnsureNoStaleParentStatus(t *testing.T) {
	testCases := []struct {
		name                     string
		httproute                *gatewayapi.HTTPRoute
		expectedAnyStatusRemoved bool
		expectedStatusesForRefs  []gatewayapi.ParentReference
	}{
		{
			name: "no stale status",
			httproute: &gatewayapi.HTTPRoute{
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{
							{Name: "defined-in-spec"},
						},
					},
				},
			},
			expectedAnyStatusRemoved: false,
			expectedStatusesForRefs:  nil, // There was no status for `defined-in-spec` created yet.
		},
		{
			name: "no stale status with existing status for spec parent ref",
			httproute: &gatewayapi.HTTPRoute{
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{
							{Name: "defined-in-spec"},
						},
					},
				},
				Status: gatewayapi.HTTPRouteStatus{
					RouteStatus: gatewayapi.RouteStatus{
						Parents: []gatewayapi.RouteParentStatus{
							{
								ControllerName: GetControllerName(),
								ParentRef:      gatewayapi.ParentReference{Name: "defined-in-spec"},
							},
						},
					},
				},
			},
			expectedStatusesForRefs: []gatewayapi.ParentReference{
				{Name: "defined-in-spec"},
			},
			expectedAnyStatusRemoved: false,
		},
		{
			name: "stale status",
			httproute: &gatewayapi.HTTPRoute{
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{
							{Name: "defined-in-spec"},
						},
					},
				},
				Status: gatewayapi.HTTPRouteStatus{
					RouteStatus: gatewayapi.RouteStatus{
						Parents: []gatewayapi.RouteParentStatus{
							{
								ControllerName: GetControllerName(),
								ParentRef:      gatewayapi.ParentReference{Name: "not-defined-in-spec"},
							},
						},
					},
				},
			},
			expectedStatusesForRefs:  nil, // There was no status for `defined-in-spec` created yet.
			expectedAnyStatusRemoved: true,
		},
		{
			name: "stale status with existing status for spec parent ref",
			httproute: &gatewayapi.HTTPRoute{
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{
							{Name: "defined-in-spec"},
						},
					},
				},
				Status: gatewayapi.HTTPRouteStatus{
					RouteStatus: gatewayapi.RouteStatus{
						Parents: []gatewayapi.RouteParentStatus{
							{
								ControllerName: GetControllerName(),
								ParentRef:      gatewayapi.ParentReference{Name: "not-defined-in-spec"},
							},
							{
								ControllerName: GetControllerName(),
								ParentRef:      gatewayapi.ParentReference{Name: "defined-in-spec"},
							},
						},
					},
				},
			},
			expectedStatusesForRefs: []gatewayapi.ParentReference{
				{Name: "defined-in-spec"},
			},
			expectedAnyStatusRemoved: true,
		},
		{
			name: "do not remove status for other controllers",
			httproute: &gatewayapi.HTTPRoute{
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{
							{Name: "defined-in-spec"},
						},
					},
				},
				Status: gatewayapi.HTTPRouteStatus{
					RouteStatus: gatewayapi.RouteStatus{
						Parents: []gatewayapi.RouteParentStatus{
							{
								ControllerName: gatewayapi.GatewayController("another-controller"),
								ParentRef:      gatewayapi.ParentReference{Name: "not-defined-in-spec"},
							},
							{
								ControllerName: GetControllerName(),
								ParentRef:      gatewayapi.ParentReference{Name: "defined-in-spec"},
							},
						},
					},
				},
			},
			expectedStatusesForRefs: []gatewayapi.ParentReference{
				{Name: "not-defined-in-spec"},
				{Name: "defined-in-spec"},
			},
			expectedAnyStatusRemoved: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			wasAnyStatusRemoved := ensureNoStaleParentStatus(tc.httproute)
			assert.Equal(t, tc.expectedAnyStatusRemoved, wasAnyStatusRemoved)

			actualStatuses := lo.Map(tc.httproute.Status.Parents, func(status gatewayapi.RouteParentStatus, _ int) string {
				return parentReferenceKey(tc.httproute.Namespace, status.ParentRef)
			})
			expectedStatuses := lo.Map(tc.expectedStatusesForRefs, func(ref gatewayapi.ParentReference, _ int) string {
				return parentReferenceKey(tc.httproute.Namespace, ref)
			})
			assert.ElementsMatch(t, expectedStatuses, actualStatuses)
		})
	}
}

func TestParentReferenceKey(t *testing.T) {
	const routeNamespace = "route-ns"
	testCases := []struct {
		name     string
		ref      gatewayapi.ParentReference
		expected string
	}{
		{
			name: "only name",
			ref: gatewayapi.ParentReference{
				Name: "foo",
			},
			expected: "route-ns/foo//",
		},
		{
			name: "name with namespace",
			ref: gatewayapi.ParentReference{
				Name:      "foo",
				Namespace: lo.ToPtr(gatewayapi.Namespace("bar")),
			},
			expected: "bar/foo//",
		},
		{
			name: "name with port number",
			ref: gatewayapi.ParentReference{
				Name: "foo",
				Port: lo.ToPtr(gatewayapi.PortNumber(1234)),
			},
			expected: "route-ns/foo//1234",
		},
		{
			name: "name with section name",
			ref: gatewayapi.ParentReference{
				Name:        "foo",
				SectionName: lo.ToPtr(gatewayapi.SectionName("section")),
			},
			expected: "route-ns/foo/section/",
		},
		{
			name: "all fields",
			ref: gatewayapi.ParentReference{
				Name:        "foo",
				Namespace:   lo.ToPtr(gatewayapi.Namespace("bar")),
				Port:        lo.ToPtr(gatewayapi.PortNumber(1234)),
				SectionName: lo.ToPtr(gatewayapi.SectionName("section")),
			},
			expected: "bar/foo/section/1234",
		},
		{
			name: "group and kind are ignored",
			ref: gatewayapi.ParentReference{
				Name:  "foo",
				Group: lo.ToPtr(gatewayapi.Group("group")),
				Kind:  lo.ToPtr(gatewayapi.Kind("kind")),
			},
			expected: "route-ns/foo//",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := parentReferenceKey(routeNamespace, tc.ref)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
