package gateway

import (
	"testing"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

func TestGetParentStatuses(t *testing.T) {
	t.Run("HTTPRoute", func(t *testing.T) {
		tests := []struct {
			name  string
			route *gatewayv1beta1.HTTPRoute
			want  map[string]*gatewayv1beta1.RouteParentStatus
		}{
			{
				name: "basic",
				route: &gatewayv1beta1.HTTPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      uuid.NewString(),
						Namespace: uuid.NewString(),
					},
					Status: gatewayv1beta1.HTTPRouteStatus{
						RouteStatus: gatewayv1beta1.RouteStatus{
							Parents: []gatewayv1beta1.RouteParentStatus{
								{
									ParentRef: gatewayv1beta1.ParentReference{
										Group:       lo.ToPtr(gatewayv1beta1.Group("group")),
										Kind:        lo.ToPtr(gatewayv1beta1.Kind("kind")),
										Namespace:   lo.ToPtr(gatewayv1beta1.Namespace("namespace")),
										Name:        gatewayv1beta1.ObjectName("name"),
										SectionName: lo.ToPtr(gatewayv1beta1.SectionName("section1")),
									},
								},
							},
						},
					},
				},
				want: map[string]*gatewayv1beta1.RouteParentStatus{
					"namespace/name/section1": {
						ParentRef: gatewayv1beta1.ParentReference{
							Group:       lo.ToPtr(gatewayv1beta1.Group("group")),
							Kind:        lo.ToPtr(gatewayv1beta1.Kind("kind")),
							Namespace:   lo.ToPtr(gatewayv1beta1.Namespace("namespace")),
							Name:        gatewayv1beta1.ObjectName("name"),
							SectionName: lo.ToPtr(gatewayv1beta1.SectionName("section1")),
						},
					},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, tt.want, getParentStatuses(tt.route, tt.route.Status.Parents))
			})
		}
	})

	t.Run("UDPRoute", func(t *testing.T) {
		tests := []struct {
			name  string
			route *gatewayv1alpha2.UDPRoute
			want  map[string]*gatewayv1alpha2.RouteParentStatus
		}{
			{
				name: "basic",
				route: &gatewayv1alpha2.UDPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      uuid.NewString(),
						Namespace: uuid.NewString(),
					},
					Status: gatewayv1alpha2.UDPRouteStatus{
						RouteStatus: gatewayv1alpha2.RouteStatus{
							Parents: []gatewayv1alpha2.RouteParentStatus{
								{
									ParentRef: gatewayv1alpha2.ParentReference{
										Group:     lo.ToPtr(gatewayv1alpha2.Group("group")),
										Kind:      lo.ToPtr(gatewayv1alpha2.Kind("kind")),
										Namespace: lo.ToPtr(gatewayv1alpha2.Namespace("namespace")),
										Name:      gatewayv1alpha2.ObjectName("name"),
									},
								},
							},
						},
					},
				},
				want: map[string]*gatewayv1alpha2.RouteParentStatus{
					"namespace/name": {
						ParentRef: gatewayv1alpha2.ParentReference{
							Group:     lo.ToPtr(gatewayv1alpha2.Group("group")),
							Kind:      lo.ToPtr(gatewayv1alpha2.Kind("kind")),
							Namespace: lo.ToPtr(gatewayv1alpha2.Namespace("namespace")),
							Name:      gatewayv1alpha2.ObjectName("name"),
						},
					},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, tt.want, getParentStatuses(tt.route, tt.route.Status.Parents))
			})
		}
	})

	t.Run("TCPRoute", func(t *testing.T) {
		tests := []struct {
			name  string
			route *gatewayv1alpha2.TCPRoute
			want  map[string]*gatewayv1alpha2.RouteParentStatus
		}{
			{
				name: "basic",
				route: &gatewayv1alpha2.TCPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      uuid.NewString(),
						Namespace: uuid.NewString(),
					},
					Status: gatewayv1alpha2.TCPRouteStatus{
						RouteStatus: gatewayv1alpha2.RouteStatus{
							Parents: []gatewayv1alpha2.RouteParentStatus{
								{
									ParentRef: gatewayv1alpha2.ParentReference{
										Group:     lo.ToPtr(gatewayv1alpha2.Group("group")),
										Kind:      lo.ToPtr(gatewayv1alpha2.Kind("kind")),
										Namespace: lo.ToPtr(gatewayv1alpha2.Namespace("namespace")),
										Name:      gatewayv1alpha2.ObjectName("name"),
									},
								},
							},
						},
					},
				},
				want: map[string]*gatewayv1alpha2.RouteParentStatus{
					"namespace/name": {
						ParentRef: gatewayv1alpha2.ParentReference{
							Group:     lo.ToPtr(gatewayv1alpha2.Group("group")),
							Kind:      lo.ToPtr(gatewayv1alpha2.Kind("kind")),
							Namespace: lo.ToPtr(gatewayv1alpha2.Namespace("namespace")),
							Name:      gatewayv1alpha2.ObjectName("name"),
						},
					},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, tt.want, getParentStatuses(tt.route, tt.route.Status.Parents))
			})
		}
	})

	t.Run("TLSRoute", func(t *testing.T) {
		tests := []struct {
			name  string
			route *gatewayv1alpha2.TLSRoute
			want  map[string]*gatewayv1alpha2.RouteParentStatus
		}{
			{
				name: "basic",
				route: &gatewayv1alpha2.TLSRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      uuid.NewString(),
						Namespace: uuid.NewString(),
					},
					Status: gatewayv1alpha2.TLSRouteStatus{
						RouteStatus: gatewayv1alpha2.RouteStatus{
							Parents: []gatewayv1alpha2.RouteParentStatus{
								{
									ParentRef: gatewayv1alpha2.ParentReference{
										Group:     lo.ToPtr(gatewayv1alpha2.Group("group")),
										Kind:      lo.ToPtr(gatewayv1alpha2.Kind("kind")),
										Namespace: lo.ToPtr(gatewayv1alpha2.Namespace("namespace")),
										Name:      gatewayv1alpha2.ObjectName("name"),
									},
								},
							},
						},
					},
				},
				want: map[string]*gatewayv1alpha2.RouteParentStatus{
					"namespace/name": {
						ParentRef: gatewayv1alpha2.ParentReference{
							Group:     lo.ToPtr(gatewayv1alpha2.Group("group")),
							Kind:      lo.ToPtr(gatewayv1alpha2.Kind("kind")),
							Namespace: lo.ToPtr(gatewayv1alpha2.Namespace("namespace")),
							Name:      gatewayv1alpha2.ObjectName("name"),
						},
					},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, tt.want, getParentStatuses(tt.route, tt.route.Status.Parents))
			})
		}
	})
}
