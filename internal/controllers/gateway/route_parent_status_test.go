package gateway

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/builder"
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
									ParentRef: builder.NewParentReference().
										Group("group").
										Kind("kind").
										Namespace("namespace").
										Name("name").
										SectionName("section1").
										Build(),
								},
							},
						},
					},
				},
				want: map[string]*gatewayapi.RouteParentStatus{
					"namespace/name/section1": {
						ParentRef: builder.NewParentReference().
							Group("group").
							Kind("kind").
							Namespace("namespace").
							Name("name").
							SectionName("section1").
							Build(),
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
									ParentRef: builder.NewParentReference().
										Group("group").
										Kind("kind").
										Namespace("namespace").
										Name("name").
										Build(),
								},
							},
						},
					},
				},
				want: map[string]*gatewayapi.RouteParentStatus{
					"namespace/name": {
						ParentRef: builder.NewParentReference().
							Group("group").
							Kind("kind").
							Namespace("namespace").
							Name("name").
							Build(),
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
									ParentRef: builder.NewParentReference().
										Group("group").
										Kind("kind").
										Namespace("namespace").
										Name("name").
										Build(),
								},
							},
						},
					},
				},
				want: map[string]*gatewayapi.RouteParentStatus{
					"namespace/name": {
						ParentRef: builder.NewParentReference().
							Group("group").
							Kind("kind").
							Namespace("namespace").
							Name("name").
							Build(),
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
									ParentRef: builder.NewParentReference().
										Group("group").
										Kind("kind").
										Namespace("namespace").
										Name("name").
										Build(),
								},
							},
						},
					},
				},
				want: map[string]*gatewayapi.RouteParentStatus{
					"namespace/name": {
						ParentRef: builder.NewParentReference().
							Group("group").
							Kind("kind").
							Namespace("namespace").
							Name("name").
							Build(),
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

func TestParentStatusesForRoute(t *testing.T) {
	t.Run("HTTPRoute", func(t *testing.T) {
		tests := []struct {
			name     string
			route    *gatewayapi.HTTPRoute
			statuses []gatewayapi.RouteParentStatus
			gateways []supportedGatewayWithCondition
			want     map[string]*gatewayapi.RouteParentStatus
			changed  bool
		}{
			{
				name: "no Programmed condition on a route, update and return true",
				route: &gatewayapi.HTTPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:       uuid.NewString(),
						Namespace:  uuid.NewString(),
						Generation: 7,
					},
					Status: gatewayapi.HTTPRouteStatus{
						RouteStatus: gatewayapi.RouteStatus{
							Parents: []gatewayapi.RouteParentStatus{
								{
									ParentRef: builder.NewParentReference().
										Group("gateway.networking.k8s.io").
										Kind("Gateway").
										Namespace("namespace").
										Name("gateway1").
										SectionName("section1").
										Build(),
									ControllerName: GetControllerName(),
								},
							},
						},
					},
				},
				gateways: []supportedGatewayWithCondition{
					{
						gateway: &gatewayapi.Gateway{
							ObjectMeta: metav1.ObjectMeta{
								Name:      "gateway1",
								Namespace: "namespace",
							},
						},
						condition: metav1.Condition{
							Type:   "Ready",
							Status: metav1.ConditionTrue,
							Reason: "Ready",
						},
						listenerName: "section1",
					},
				},
				want: map[string]*gatewayapi.RouteParentStatus{
					"namespace/gateway1/section1": {
						ParentRef: builder.NewParentReference().
							Group("gateway.networking.k8s.io").
							Kind("Gateway").
							Namespace("namespace").
							Name("gateway1").
							SectionName("section1").
							Build(),
						ControllerName: GetControllerName(),
						Conditions: []metav1.Condition{
							newCondition("Ready", metav1.ConditionTrue, "Ready", 7),
						},
					},
				},
				changed: true,
			},
			{
				name: "Programmed condition is as exptected on a route, don't update and return false",
				route: &gatewayapi.HTTPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:       uuid.NewString(),
						Namespace:  uuid.NewString(),
						Generation: 7,
					},
					Status: gatewayapi.HTTPRouteStatus{
						RouteStatus: gatewayapi.RouteStatus{
							Parents: []gatewayapi.RouteParentStatus{
								{
									ParentRef: builder.NewParentReference().
										Group("gateway.networking.k8s.io").
										Kind("Gateway").
										Namespace("namespace").
										Name("gateway1").
										SectionName("section1").
										Build(),
									ControllerName: GetControllerName(),
									Conditions: []metav1.Condition{
										newCondition("Ready", metav1.ConditionTrue, "Ready", 7),
									},
								},
							},
						},
					},
				},
				gateways: []supportedGatewayWithCondition{
					{
						gateway: &gatewayapi.Gateway{
							ObjectMeta: metav1.ObjectMeta{
								Name:      "gateway1",
								Namespace: "namespace",
							},
						},
						condition: metav1.Condition{
							Type:   "Ready",
							Status: metav1.ConditionTrue,
							Reason: "Ready",
						},
						listenerName: "section1",
					},
				},
				want: map[string]*gatewayapi.RouteParentStatus{
					"namespace/gateway1/section1": {
						ParentRef: builder.NewParentReference().
							Group("gateway.networking.k8s.io").
							Kind("Gateway").
							Namespace("namespace").
							Name("gateway1").
							SectionName("section1").
							Build(),
						ControllerName: GetControllerName(),
						Conditions: []metav1.Condition{
							newCondition("Ready", metav1.ConditionTrue, "Ready", 7),
						},
					},
				},
				changed: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, changed := parentStatusesForRoute(tt.route, tt.route.Status.Parents, tt.gateways...)
				ignoreLastTransitionTime := cmpopts.IgnoreFields(metav1.Condition{}, "LastTransitionTime")
				if !cmp.Equal(tt.want, got, ignoreLastTransitionTime) {
					assert.Equal(t, tt.want, got)
				}
				assert.Equal(t, tt.changed, changed)
			})
		}
	})

	t.Run("UDPRoute", func(t *testing.T) {
		tests := []struct {
			name     string
			route    *gatewayapi.UDPRoute
			statuses []gatewayapi.RouteParentStatus
			gateways []supportedGatewayWithCondition
			want     map[string]*gatewayapi.RouteParentStatus
			changed  bool
		}{
			{
				name: "no Programmed condition on a route, update and return true",
				route: &gatewayapi.UDPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:       uuid.NewString(),
						Namespace:  uuid.NewString(),
						Generation: 7,
					},
					Status: gatewayapi.UDPRouteStatus{
						RouteStatus: gatewayapi.RouteStatus{
							Parents: []gatewayapi.RouteParentStatus{
								{
									ParentRef: builder.NewParentReference().
										Group("gateway.networking.k8s.io").
										Kind("Gateway").
										Namespace("namespace").
										Name("gateway1").
										Build(),
									ControllerName: GetControllerName(),
								},
							},
						},
					},
				},
				gateways: []supportedGatewayWithCondition{
					{
						gateway: &gatewayapi.Gateway{
							ObjectMeta: metav1.ObjectMeta{
								Name:      "gateway1",
								Namespace: "namespace",
							},
						},
						condition: metav1.Condition{
							Type:   "Ready",
							Status: metav1.ConditionTrue,
							Reason: "Ready",
						},
					},
				},
				want: map[string]*gatewayapi.RouteParentStatus{
					"namespace/gateway1": {
						ParentRef: builder.NewParentReference().
							Group("gateway.networking.k8s.io").
							Kind("Gateway").
							Namespace("namespace").
							Name("gateway1").
							Build(),
						ControllerName: GetControllerName(),
						Conditions: []metav1.Condition{
							newCondition("Ready", metav1.ConditionTrue, "Ready", 7),
						},
					},
				},
				changed: true,
			},
			{
				name: "Programmed condition is as exptected on a route, don't update and return false",
				route: &gatewayapi.UDPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:       uuid.NewString(),
						Namespace:  uuid.NewString(),
						Generation: 7,
					},
					Status: gatewayapi.UDPRouteStatus{
						RouteStatus: gatewayapi.RouteStatus{
							Parents: []gatewayapi.RouteParentStatus{
								{
									ParentRef: builder.NewParentReference().
										Group("gateway.networking.k8s.io").
										Kind("Gateway").
										Namespace("namespace").
										Name("gateway1").
										Build(),
									ControllerName: GetControllerName(),
									Conditions: []metav1.Condition{
										newCondition("Ready", metav1.ConditionTrue, "Ready", 7),
									},
								},
							},
						},
					},
				},
				gateways: []supportedGatewayWithCondition{
					{
						gateway: &gatewayapi.Gateway{
							ObjectMeta: metav1.ObjectMeta{
								Name:      "gateway1",
								Namespace: "namespace",
							},
						},
						condition: metav1.Condition{
							Type:   "Ready",
							Status: metav1.ConditionTrue,
							Reason: "Ready",
						},
					},
				},
				want: map[string]*gatewayapi.RouteParentStatus{
					"namespace/gateway1": {
						ParentRef: builder.NewParentReference().
							Group("gateway.networking.k8s.io").
							Kind("Gateway").
							Namespace("namespace").
							Name("gateway1").
							Build(),
						ControllerName: GetControllerName(),
						Conditions: []metav1.Condition{
							newCondition("Ready", metav1.ConditionTrue, "Ready", 7),
						},
					},
				},
				changed: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, changed := parentStatusesForRoute(tt.route, tt.route.Status.Parents, tt.gateways...)
				ignoreLastTransitionTime := cmpopts.IgnoreFields(metav1.Condition{}, "LastTransitionTime")
				if !cmp.Equal(tt.want, got, ignoreLastTransitionTime) {
					assert.Equal(t, tt.want, got)
				}
				assert.Equal(t, tt.changed, changed)
			})
		}
	})

	t.Run("TCPRoute", func(t *testing.T) {
		tests := []struct {
			name     string
			route    *gatewayapi.TCPRoute
			statuses []gatewayapi.RouteParentStatus
			gateways []supportedGatewayWithCondition
			want     map[string]*gatewayapi.RouteParentStatus
			changed  bool
		}{
			{
				name: "no Programmed condition on a route, update and return true",
				route: &gatewayapi.TCPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:       uuid.NewString(),
						Namespace:  uuid.NewString(),
						Generation: 7,
					},
					Status: gatewayapi.TCPRouteStatus{
						RouteStatus: gatewayapi.RouteStatus{
							Parents: []gatewayapi.RouteParentStatus{
								{
									ParentRef: builder.NewParentReference().
										Group("gateway.networking.k8s.io").
										Kind("Gateway").
										Namespace("namespace").
										Name("gateway1").
										Build(),
									ControllerName: GetControllerName(),
								},
							},
						},
					},
				},
				gateways: []supportedGatewayWithCondition{
					{
						gateway: &gatewayapi.Gateway{
							ObjectMeta: metav1.ObjectMeta{
								Name:      "gateway1",
								Namespace: "namespace",
							},
						},
						condition: metav1.Condition{
							Type:   "Ready",
							Status: metav1.ConditionTrue,
							Reason: "Ready",
						},
					},
				},
				want: map[string]*gatewayapi.RouteParentStatus{
					"namespace/gateway1": {
						ParentRef: builder.NewParentReference().
							Group("gateway.networking.k8s.io").
							Kind("Gateway").
							Namespace("namespace").
							Name("gateway1").
							Build(),
						ControllerName: GetControllerName(),
						Conditions: []metav1.Condition{
							newCondition("Ready", metav1.ConditionTrue, "Ready", 7),
						},
					},
				},
				changed: true,
			},
			{
				name: "Programmed condition is as exptected on a route, don't update and return false",
				route: &gatewayapi.TCPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:       uuid.NewString(),
						Namespace:  uuid.NewString(),
						Generation: 7,
					},
					Status: gatewayapi.TCPRouteStatus{
						RouteStatus: gatewayapi.RouteStatus{
							Parents: []gatewayapi.RouteParentStatus{
								{
									ParentRef: builder.NewParentReference().
										Group("gateway.networking.k8s.io").
										Kind("Gateway").
										Namespace("namespace").
										Name("gateway1").
										Build(),
									ControllerName: GetControllerName(),
									Conditions: []metav1.Condition{
										newCondition("Ready", metav1.ConditionTrue, "Ready", 7),
									},
								},
							},
						},
					},
				},
				gateways: []supportedGatewayWithCondition{
					{
						gateway: &gatewayapi.Gateway{
							ObjectMeta: metav1.ObjectMeta{
								Name:      "gateway1",
								Namespace: "namespace",
							},
						},
						condition: metav1.Condition{
							Type:   "Ready",
							Status: metav1.ConditionTrue,
							Reason: "Ready",
						},
					},
				},
				want: map[string]*gatewayapi.RouteParentStatus{
					"namespace/gateway1": {
						ParentRef: builder.NewParentReference().
							Group("gateway.networking.k8s.io").
							Kind("Gateway").
							Namespace("namespace").
							Name("gateway1").
							Build(),
						ControllerName: GetControllerName(),
						Conditions: []metav1.Condition{
							newCondition("Ready", metav1.ConditionTrue, "Ready", 7),
						},
					},
				},
				changed: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, changed := parentStatusesForRoute(tt.route, tt.route.Status.Parents, tt.gateways...)
				ignoreLastTransitionTime := cmpopts.IgnoreFields(metav1.Condition{}, "LastTransitionTime")
				if !cmp.Equal(tt.want, got, ignoreLastTransitionTime) {
					assert.Equal(t, tt.want, got)
				}
				assert.Equal(t, tt.changed, changed)
			})
		}
	})

	t.Run("TLSRoute", func(t *testing.T) {
		tests := []struct {
			name     string
			route    *gatewayapi.TLSRoute
			statuses []gatewayapi.RouteParentStatus
			gateways []supportedGatewayWithCondition
			want     map[string]*gatewayapi.RouteParentStatus
			changed  bool
		}{
			{
				name: "no Programmed condition on a route, update and return true",
				route: &gatewayapi.TLSRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:       uuid.NewString(),
						Namespace:  uuid.NewString(),
						Generation: 7,
					},
					Status: gatewayapi.TLSRouteStatus{
						RouteStatus: gatewayapi.RouteStatus{
							Parents: []gatewayapi.RouteParentStatus{
								{
									ParentRef: builder.NewParentReference().
										Group("gateway.networking.k8s.io").
										Kind("Gateway").
										Namespace("namespace").
										Name("gateway1").
										Build(),
									ControllerName: GetControllerName(),
								},
							},
						},
					},
				},
				gateways: []supportedGatewayWithCondition{
					{
						gateway: &gatewayapi.Gateway{
							ObjectMeta: metav1.ObjectMeta{
								Name:      "gateway1",
								Namespace: "namespace",
							},
						},
						condition: metav1.Condition{
							Type:   "Ready",
							Status: metav1.ConditionTrue,
							Reason: "Ready",
						},
					},
				},
				want: map[string]*gatewayapi.RouteParentStatus{
					"namespace/gateway1": {
						ParentRef: builder.NewParentReference().
							Group("gateway.networking.k8s.io").
							Kind("Gateway").
							Namespace("namespace").
							Name("gateway1").
							Build(),
						ControllerName: GetControllerName(),
						Conditions: []metav1.Condition{
							newCondition("Ready", metav1.ConditionTrue, "Ready", 7),
						},
					},
				},
				changed: true,
			},
			{
				name: "Programmed condition is as exptected on a route, don't update and return false",
				route: &gatewayapi.TLSRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:       uuid.NewString(),
						Namespace:  uuid.NewString(),
						Generation: 7,
					},
					Status: gatewayapi.TLSRouteStatus{
						RouteStatus: gatewayapi.RouteStatus{
							Parents: []gatewayapi.RouteParentStatus{
								{
									ParentRef: builder.NewParentReference().
										Group("gateway.networking.k8s.io").
										Kind("Gateway").
										Namespace("namespace").
										Name("gateway1").
										Build(),
									ControllerName: GetControllerName(),
									Conditions: []metav1.Condition{
										newCondition("Ready", metav1.ConditionTrue, "Ready", 7),
									},
								},
							},
						},
					},
				},
				gateways: []supportedGatewayWithCondition{
					{
						gateway: &gatewayapi.Gateway{
							ObjectMeta: metav1.ObjectMeta{
								Name:      "gateway1",
								Namespace: "namespace",
							},
						},
						condition: metav1.Condition{
							Type:   "Ready",
							Status: metav1.ConditionTrue,
							Reason: "Ready",
						},
					},
				},
				want: map[string]*gatewayapi.RouteParentStatus{
					"namespace/gateway1": {
						ParentRef: builder.NewParentReference().
							Group("gateway.networking.k8s.io").
							Kind("Gateway").
							Namespace("namespace").
							Name("gateway1").
							Build(),
						ControllerName: GetControllerName(),
						Conditions: []metav1.Condition{
							newCondition("Ready", metav1.ConditionTrue, "Ready", 7),
						},
					},
				},
				changed: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, changed := parentStatusesForRoute(tt.route, tt.route.Status.Parents, tt.gateways...)
				ignoreLastTransitionTime := cmpopts.IgnoreFields(metav1.Condition{}, "LastTransitionTime")
				if !cmp.Equal(tt.want, got, ignoreLastTransitionTime) {
					assert.Equal(t, tt.want, got)
				}
				assert.Equal(t, tt.changed, changed)
			})
		}
	})

	t.Run("GRPCRoute", func(t *testing.T) {
		tests := []struct {
			name     string
			route    *gatewayapi.GRPCRoute
			statuses []gatewayapi.RouteParentStatus
			gateways []supportedGatewayWithCondition
			want     map[string]*gatewayapi.RouteParentStatus
			changed  bool
		}{
			{
				name: "no Programmed condition on a route, update and return true",
				route: &gatewayapi.GRPCRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:       uuid.NewString(),
						Namespace:  uuid.NewString(),
						Generation: 7,
					},
					Status: gatewayapi.GRPCRouteStatus{
						RouteStatus: gatewayapi.RouteStatus{
							Parents: []gatewayapi.RouteParentStatus{
								{
									ParentRef: builder.NewParentReference().
										Group("gateway.networking.k8s.io").
										Kind("Gateway").
										Namespace("namespace").
										Name("gateway1").
										SectionName("section1").
										Build(),
									ControllerName: GetControllerName(),
								},
							},
						},
					},
				},
				gateways: []supportedGatewayWithCondition{
					{
						gateway: &gatewayapi.Gateway{
							ObjectMeta: metav1.ObjectMeta{
								Name:      "gateway1",
								Namespace: "namespace",
							},
						},
						condition: metav1.Condition{
							Type:   "Ready",
							Status: metav1.ConditionTrue,
							Reason: "Ready",
						},
						listenerName: "section1",
					},
				},
				want: map[string]*gatewayapi.RouteParentStatus{
					"namespace/gateway1/section1": {
						ParentRef: builder.NewParentReference().
							Group("gateway.networking.k8s.io").
							Kind("Gateway").
							Namespace("namespace").
							Name("gateway1").
							SectionName("section1").
							Build(),
						ControllerName: GetControllerName(),
						Conditions: []metav1.Condition{
							newCondition("Ready", metav1.ConditionTrue, "Ready", 7),
						},
					},
				},
				changed: true,
			},
			{
				name: "Programmed condition is as exptected on a route, don't update and return false",
				route: &gatewayapi.GRPCRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:       uuid.NewString(),
						Namespace:  uuid.NewString(),
						Generation: 7,
					},
					Status: gatewayapi.GRPCRouteStatus{
						RouteStatus: gatewayapi.RouteStatus{
							Parents: []gatewayapi.RouteParentStatus{
								{
									ParentRef: builder.NewParentReference().
										Group("gateway.networking.k8s.io").
										Kind("Gateway").
										Namespace("namespace").
										Name("gateway1").
										SectionName("section1").
										Build(),
									ControllerName: GetControllerName(),
									Conditions: []metav1.Condition{
										newCondition("Ready", metav1.ConditionTrue, "Ready", 7),
									},
								},
							},
						},
					},
				},
				gateways: []supportedGatewayWithCondition{
					{
						gateway: &gatewayapi.Gateway{
							ObjectMeta: metav1.ObjectMeta{
								Name:      "gateway1",
								Namespace: "namespace",
							},
						},
						condition: metav1.Condition{
							Type:   "Ready",
							Status: metav1.ConditionTrue,
							Reason: "Ready",
						},
						listenerName: "section1",
					},
				},
				want: map[string]*gatewayapi.RouteParentStatus{
					"namespace/gateway1/section1": {
						ParentRef: builder.NewParentReference().
							Group("gateway.networking.k8s.io").
							Kind("Gateway").
							Namespace("namespace").
							Name("gateway1").
							SectionName("section1").
							Build(),
						ControllerName: GetControllerName(),
						Conditions: []metav1.Condition{
							newCondition("Ready", metav1.ConditionTrue, "Ready", 7),
						},
					},
				},
				changed: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, changed := parentStatusesForRoute(tt.route, tt.route.Status.Parents, tt.gateways...)
				ignoreLastTransitionTime := cmpopts.IgnoreFields(metav1.Condition{}, "LastTransitionTime")
				if !cmp.Equal(tt.want, got, ignoreLastTransitionTime) {
					assert.Equal(t, tt.want, got)
				}
				assert.Equal(t, tt.changed, changed)
			})
		}
	})
}
