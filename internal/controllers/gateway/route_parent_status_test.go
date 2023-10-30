package gateway

import (
	"testing"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

func TestGetParentStatuses(t *testing.T) {
	t.Run("HTTPRoute", func(t *testing.T) {
		tests := []struct {
			name  string
			route *gatewayapi.HTTPRoute
			want  map[string]*gatewayapi.RouteParentStatus
		}{
			{
				name: "basic",
				route: &gatewayapi.HTTPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      uuid.NewString(),
						Namespace: uuid.NewString(),
					},
					Status: gatewayapi.HTTPRouteStatus{
						RouteStatus: gatewayapi.RouteStatus{
							Parents: []gatewayapi.RouteParentStatus{
								{
									ParentRef: gatewayapi.ParentReference{
										Group:       lo.ToPtr(gatewayapi.Group("group")),
										Kind:        lo.ToPtr(gatewayapi.Kind("kind")),
										Namespace:   lo.ToPtr(gatewayapi.Namespace("namespace")),
										Name:        gatewayapi.ObjectName("name"),
										SectionName: lo.ToPtr(gatewayapi.SectionName("section1")),
									},
								},
							},
						},
					},
				},
				want: map[string]*gatewayapi.RouteParentStatus{
					"namespace/name/section1": {
						ParentRef: gatewayapi.ParentReference{
							Group:       lo.ToPtr(gatewayapi.Group("group")),
							Kind:        lo.ToPtr(gatewayapi.Kind("kind")),
							Namespace:   lo.ToPtr(gatewayapi.Namespace("namespace")),
							Name:        gatewayapi.ObjectName("name"),
							SectionName: lo.ToPtr(gatewayapi.SectionName("section1")),
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
			route *gatewayapi.UDPRoute
			want  map[string]*gatewayapi.RouteParentStatus
		}{
			{
				name: "basic",
				route: &gatewayapi.UDPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      uuid.NewString(),
						Namespace: uuid.NewString(),
					},
					Status: gatewayapi.UDPRouteStatus{
						RouteStatus: gatewayapi.RouteStatus{
							Parents: []gatewayapi.RouteParentStatus{
								{
									ParentRef: gatewayapi.ParentReference{
										Group:     lo.ToPtr(gatewayapi.Group("group")),
										Kind:      lo.ToPtr(gatewayapi.Kind("kind")),
										Namespace: lo.ToPtr(gatewayapi.Namespace("namespace")),
										Name:      gatewayapi.ObjectName("name"),
									},
								},
							},
						},
					},
				},
				want: map[string]*gatewayapi.RouteParentStatus{
					"namespace/name": {
						ParentRef: gatewayapi.ParentReference{
							Group:     lo.ToPtr(gatewayapi.Group("group")),
							Kind:      lo.ToPtr(gatewayapi.Kind("kind")),
							Namespace: lo.ToPtr(gatewayapi.Namespace("namespace")),
							Name:      gatewayapi.ObjectName("name"),
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
			route *gatewayapi.TCPRoute
			want  map[string]*gatewayapi.RouteParentStatus
		}{
			{
				name: "basic",
				route: &gatewayapi.TCPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      uuid.NewString(),
						Namespace: uuid.NewString(),
					},
					Status: gatewayapi.TCPRouteStatus{
						RouteStatus: gatewayapi.RouteStatus{
							Parents: []gatewayapi.RouteParentStatus{
								{
									ParentRef: gatewayapi.ParentReference{
										Group:     lo.ToPtr(gatewayapi.Group("group")),
										Kind:      lo.ToPtr(gatewayapi.Kind("kind")),
										Namespace: lo.ToPtr(gatewayapi.Namespace("namespace")),
										Name:      gatewayapi.ObjectName("name"),
									},
								},
							},
						},
					},
				},
				want: map[string]*gatewayapi.RouteParentStatus{
					"namespace/name": {
						ParentRef: gatewayapi.ParentReference{
							Group:     lo.ToPtr(gatewayapi.Group("group")),
							Kind:      lo.ToPtr(gatewayapi.Kind("kind")),
							Namespace: lo.ToPtr(gatewayapi.Namespace("namespace")),
							Name:      gatewayapi.ObjectName("name"),
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
			route *gatewayapi.TLSRoute
			want  map[string]*gatewayapi.RouteParentStatus
		}{
			{
				name: "basic",
				route: &gatewayapi.TLSRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      uuid.NewString(),
						Namespace: uuid.NewString(),
					},
					Status: gatewayapi.TLSRouteStatus{
						RouteStatus: gatewayapi.RouteStatus{
							Parents: []gatewayapi.RouteParentStatus{
								{
									ParentRef: gatewayapi.ParentReference{
										Group:     lo.ToPtr(gatewayapi.Group("group")),
										Kind:      lo.ToPtr(gatewayapi.Kind("kind")),
										Namespace: lo.ToPtr(gatewayapi.Namespace("namespace")),
										Name:      gatewayapi.ObjectName("name"),
									},
								},
							},
						},
					},
				},
				want: map[string]*gatewayapi.RouteParentStatus{
					"namespace/name": {
						ParentRef: gatewayapi.ParentReference{
							Group:     lo.ToPtr(gatewayapi.Group("group")),
							Kind:      lo.ToPtr(gatewayapi.Kind("kind")),
							Namespace: lo.ToPtr(gatewayapi.Namespace("namespace")),
							Name:      gatewayapi.ObjectName("name"),
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
