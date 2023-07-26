package translators

import (
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

func TestGenerateKongExpressionRoutesFromGRPCRouteRule(t *testing.T) {
	testCases := []struct {
		name           string
		objectName     string
		annotations    map[string]string
		hostnames      []string
		rule           gatewayv1alpha2.GRPCRouteRule
		expectedRoutes []kongstate.Route
	}{
		{
			name:        "single match without hostname",
			objectName:  "single-match",
			annotations: map[string]string{},
			rule: gatewayv1alpha2.GRPCRouteRule{
				Matches: []gatewayv1alpha2.GRPCRouteMatch{
					{
						Method: &gatewayv1alpha2.GRPCMethodMatch{
							Service: lo.ToPtr("service0"),
							Method:  nil,
						},
						Headers: []gatewayv1alpha2.GRPCHeaderMatch{
							{
								Name:  gatewayv1alpha2.GRPCHeaderName("X-Foo"),
								Value: "Bar",
							},
						},
					},
				},
			},
			expectedRoutes: []kongstate.Route{
				{
					Ingress: util.K8sObjectInfo{
						Name:             "single-match",
						Namespace:        "default",
						Annotations:      map[string]string{},
						GroupVersionKind: grpcRouteGVK,
					},
					ExpressionRoutes: true,
					Route: kong.Route{
						Name:       kong.String("grpcroute.default.single-match.0.0"),
						Expression: kong.String(`(http.path ^= "/service0/") && (http.headers.x_foo == "Bar")`),
						Priority:   kong.Int(1),
					},
				},
			},
		},
		{
			name:        "single match with hostname",
			objectName:  "single-match-with-hostname",
			annotations: map[string]string{},
			hostnames:   []string{"foo.com", "*.foo.com"},
			rule: gatewayv1alpha2.GRPCRouteRule{
				Matches: []gatewayv1alpha2.GRPCRouteMatch{
					{
						Method: &gatewayv1alpha2.GRPCMethodMatch{
							Service: lo.ToPtr("service0"),
							Method:  lo.ToPtr("method0"),
						},
					},
				},
			},
			expectedRoutes: []kongstate.Route{
				{
					Ingress: util.K8sObjectInfo{
						Name:             "single-match-with-hostname",
						Namespace:        "default",
						Annotations:      map[string]string{},
						GroupVersionKind: grpcRouteGVK,
					},
					ExpressionRoutes: true,
					Route: kong.Route{
						Name:       kong.String("grpcroute.default.single-match-with-hostname.0.0"),
						Expression: kong.String(`(http.path == "/service0/method0") && ((http.host == "foo.com") || (http.host =^ ".foo.com"))`),
						Priority:   kong.Int(1),
					},
				},
			},
		},
		{
			name:       "multuple matches without hostname",
			objectName: "multiple-matches",
			rule: gatewayv1alpha2.GRPCRouteRule{
				Matches: []gatewayv1alpha2.GRPCRouteMatch{
					{
						Method: &gatewayv1alpha2.GRPCMethodMatch{
							Service: nil,
							Method:  lo.ToPtr("method0"),
						},
						Headers: []gatewayv1alpha2.GRPCHeaderMatch{
							{
								Name:  "Version",
								Value: "2",
							},
							{
								Name:  "Client",
								Value: "kong-test",
							},
						},
					},
					{
						Method: &gatewayv1alpha2.GRPCMethodMatch{
							Type:    lo.ToPtr(gatewayv1alpha2.GRPCMethodMatchRegularExpression),
							Service: lo.ToPtr("v[012]"),
						},
					},
				},
			},
			expectedRoutes: []kongstate.Route{
				{
					Ingress: util.K8sObjectInfo{
						Name:             "multiple-matches",
						Namespace:        "default",
						Annotations:      map[string]string{},
						GroupVersionKind: grpcRouteGVK,
					},
					ExpressionRoutes: true,
					Route: kong.Route{
						Name:       kong.String("grpcroute.default.multiple-matches.0.0"),
						Expression: kong.String(`(http.path =^ "/method0") && ((http.headers.client == "kong-test") && (http.headers.version == "2"))`),
						Priority:   kong.Int(1),
					},
				},
				{
					Ingress: util.K8sObjectInfo{
						Name:             "multiple-matches",
						Namespace:        "default",
						Annotations:      map[string]string{},
						GroupVersionKind: grpcRouteGVK,
					},
					ExpressionRoutes: true,
					Route: kong.Route{
						Name:       kong.String("grpcroute.default.multiple-matches.0.1"),
						Expression: kong.String(`http.path ~ "^/v[012]/.+"`),
						Priority:   kong.Int(1),
					},
				},
			},
		},
		{
			name:       "single match with annotations",
			objectName: "single-match-with-annotations",
			annotations: map[string]string{
				"konghq.com/methods":   "POST,GET",
				"konghq.com/protocols": "https",
				"konghq.com/snis":      "kong.foo.com",
			},
			rule: gatewayv1alpha2.GRPCRouteRule{
				Matches: []gatewayv1alpha2.GRPCRouteMatch{
					{
						Method: &gatewayv1alpha2.GRPCMethodMatch{
							Service: lo.ToPtr("service0"),
							Method:  lo.ToPtr("method0"),
						},
					},
				},
			},
			expectedRoutes: []kongstate.Route{
				{
					Ingress: util.K8sObjectInfo{
						Name:      "single-match-with-annotations",
						Namespace: "default",
						Annotations: map[string]string{
							"konghq.com/methods":   "POST,GET",
							"konghq.com/protocols": "https",
							"konghq.com/snis":      "kong.foo.com",
						},
						GroupVersionKind: grpcRouteGVK,
					},
					ExpressionRoutes: true,
					Route: kong.Route{
						Name:       kong.String("grpcroute.default.single-match-with-annotations.0.0"),
						Expression: kong.String(`(http.path == "/service0/method0") && (net.protocol == "https") && (tls.sni == "kong.foo.com")`),
						Priority:   kong.Int(1),
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			grpcroute := makeTestGRPCRoute(tc.objectName, "default", tc.annotations, tc.hostnames, []gatewayv1alpha2.GRPCRouteRule{tc.rule})
			routes := GenerateKongExpressionRoutesFromGRPCRouteRule(grpcroute, 0)
			require.Equal(t, tc.expectedRoutes, routes)
		})
	}
}

func TestMethodMatcherFromGRPCMethodMatch(t *testing.T) {
	testCases := []struct {
		name       string
		match      gatewayv1alpha2.GRPCMethodMatch
		expression string
	}{
		{
			name: "exact method match",
			match: gatewayv1alpha2.GRPCMethodMatch{
				Type:    lo.ToPtr(gatewayv1alpha2.GRPCMethodMatchExact),
				Service: lo.ToPtr("service0"),
				Method:  lo.ToPtr("method0"),
			},
			expression: `http.path == "/service0/method0"`,
		},
		{
			name: "exact match with service unspecified",
			match: gatewayv1alpha2.GRPCMethodMatch{
				Type:   lo.ToPtr(gatewayv1alpha2.GRPCMethodMatchExact),
				Method: lo.ToPtr("method0"),
			},
			expression: `http.path =^ "/method0"`,
		},
		{
			name: "regex method match",
			match: gatewayv1alpha2.GRPCMethodMatch{
				Type:    lo.ToPtr(gatewayv1alpha2.GRPCMethodMatchRegularExpression),
				Service: lo.ToPtr("auth-v[012]"),
				Method:  lo.ToPtr("[a-z]+"),
			},
			expression: `http.path ~ "^/auth-v[012]/[a-z]+"`,
		},
		{
			name: "empty regex match",
			match: gatewayv1alpha2.GRPCMethodMatch{
				Type: lo.ToPtr(gatewayv1alpha2.GRPCMethodMatchRegularExpression),
			},
			expression: `http.path ^= "/"`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expression, methodMatcherFromGRPCMethodMatch(&tc.match).Expression())
		})
	}
}

func TestSplitGRPCRoute(t *testing.T) {
	namesToBackendRefs := func(names []string) []gatewayv1alpha2.GRPCBackendRef {
		backendRefs := []gatewayv1alpha2.GRPCBackendRef{}
		for _, name := range names {
			backendRefs = append(backendRefs,
				gatewayv1alpha2.GRPCBackendRef{
					BackendRef: gatewayv1alpha2.BackendRef{
						BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
							Name: gatewayv1alpha2.ObjectName(name),
						},
					},
				},
			)
		}
		return backendRefs
	}

	testGRPCRoutes := []*gatewayv1alpha2.GRPCRoute{
		{
			TypeMeta: grpcRouteTypeMeta,
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "grpcroute-no-hostname-one-match",
			},
			Spec: gatewayv1alpha2.GRPCRouteSpec{
				Rules: []gatewayv1alpha2.GRPCRouteRule{
					{
						Matches: []gatewayv1alpha2.GRPCRouteMatch{
							{
								Method: &gatewayv1alpha2.GRPCMethodMatch{
									Service: lo.ToPtr("pets"),
									Method:  lo.ToPtr("list"),
								},
							},
						},
						BackendRefs: namesToBackendRefs([]string{"listpets"}),
					},
				},
			},
		},
		{
			TypeMeta: grpcRouteTypeMeta,
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "grpcroute-one-hostname-multiple-matches",
			},
			Spec: gatewayv1alpha2.GRPCRouteSpec{
				Hostnames: []gatewayv1alpha2.Hostname{
					gatewayv1alpha2.Hostname("petstore.com"),
				},
				Rules: []gatewayv1alpha2.GRPCRouteRule{
					{
						Matches: []gatewayv1alpha2.GRPCRouteMatch{
							{
								Method: &gatewayv1alpha2.GRPCMethodMatch{
									Service: lo.ToPtr("cats"),
									Method:  lo.ToPtr("list"),
								},
							},
							{
								Method: &gatewayv1alpha2.GRPCMethodMatch{
									Service: lo.ToPtr("dogs"),
									Method:  lo.ToPtr("list"),
								},
							},
						},
						BackendRefs: namesToBackendRefs([]string{"listpets"}),
					},
					{
						Matches: []gatewayv1alpha2.GRPCRouteMatch{
							{
								Method: &gatewayv1alpha2.GRPCMethodMatch{
									Service: lo.ToPtr("cats"),
									Method:  lo.ToPtr("create"),
								},
							},
							{
								Method: &gatewayv1alpha2.GRPCMethodMatch{
									Service: lo.ToPtr("dogs"),
									Method:  lo.ToPtr("create"),
								},
							},
						},
						BackendRefs: namesToBackendRefs([]string{"createpets"}),
					},
				},
			},
		},
		{
			TypeMeta: grpcRouteTypeMeta,
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "grpcroute-multiple-hostnames",
			},
			Spec: gatewayv1alpha2.GRPCRouteSpec{
				Hostnames: []gatewayv1alpha2.Hostname{
					gatewayv1alpha2.Hostname("petstore.com"),
					gatewayv1alpha2.Hostname("petstore.net"),
				},
				Rules: []gatewayv1alpha2.GRPCRouteRule{
					{
						Matches: []gatewayv1alpha2.GRPCRouteMatch{
							{
								Method: &gatewayv1alpha2.GRPCMethodMatch{
									Service: lo.ToPtr("pets"),
									Method:  lo.ToPtr("list"),
								},
							},
						},
						BackendRefs: namesToBackendRefs([]string{"listpets"}),
					},
				},
			},
		},
	}

	testCases := []struct {
		name                 string
		grpcRoute            *gatewayv1alpha2.GRPCRoute
		expectedSplitMatches []SplitGRPCRouteMatch
	}{
		{
			name:      "no hostname and one match",
			grpcRoute: testGRPCRoutes[0],
			expectedSplitMatches: []SplitGRPCRouteMatch{
				{
					Source:     testGRPCRoutes[0],
					Hostname:   "",
					Match:      testGRPCRoutes[0].Spec.Rules[0].Matches[0],
					RuleIndex:  0,
					MatchIndex: 0,
				},
			},
		},
		{
			name:      "single hostname and multiple rules",
			grpcRoute: testGRPCRoutes[1],
			expectedSplitMatches: []SplitGRPCRouteMatch{
				{
					Source:     testGRPCRoutes[1],
					Hostname:   string(testGRPCRoutes[1].Spec.Hostnames[0]),
					Match:      testGRPCRoutes[1].Spec.Rules[0].Matches[0],
					RuleIndex:  0,
					MatchIndex: 0,
				},
				{
					Source:     testGRPCRoutes[1],
					Hostname:   string(testGRPCRoutes[1].Spec.Hostnames[0]),
					Match:      testGRPCRoutes[1].Spec.Rules[0].Matches[1],
					RuleIndex:  0,
					MatchIndex: 1,
				},
				{
					Source:     testGRPCRoutes[1],
					Hostname:   string(testGRPCRoutes[1].Spec.Hostnames[0]),
					Match:      testGRPCRoutes[1].Spec.Rules[1].Matches[0],
					RuleIndex:  1,
					MatchIndex: 0,
				},
				{
					Source:     testGRPCRoutes[1],
					Hostname:   string(testGRPCRoutes[1].Spec.Hostnames[0]),
					Match:      testGRPCRoutes[1].Spec.Rules[1].Matches[1],
					RuleIndex:  1,
					MatchIndex: 1,
				},
			},
		},
		{
			name:      "multiple hostnames",
			grpcRoute: testGRPCRoutes[2],
			expectedSplitMatches: []SplitGRPCRouteMatch{
				{
					Source:     testGRPCRoutes[2],
					Hostname:   string(testGRPCRoutes[2].Spec.Hostnames[0]),
					Match:      testGRPCRoutes[2].Spec.Rules[0].Matches[0],
					RuleIndex:  0,
					MatchIndex: 0,
				},
				{
					Source:     testGRPCRoutes[2],
					Hostname:   string(testGRPCRoutes[2].Spec.Hostnames[1]),
					Match:      testGRPCRoutes[2].Spec.Rules[0].Matches[0],
					RuleIndex:  0,
					MatchIndex: 0,
				},
			},
		},
	}

	for i, tc := range testCases {
		tc := tc
		indexStr := strconv.Itoa(i)
		t.Run(indexStr+"-"+tc.name, func(t *testing.T) {
			splitMatches := SplitGRPCRoute(tc.grpcRoute)
			require.Len(t, splitMatches, len(tc.expectedSplitMatches), "should have same number of split matches with expected")
			for i, splitGRPCRoute := range tc.expectedSplitMatches {
				require.Truef(t, reflect.DeepEqual(splitGRPCRoute, splitMatches[i]),
					"should have the same GRPCRoute match as expected on index %d", i)
			}
		})
	}
}

func TestCalculateSplitGRCPRoutePriorityTraits(t *testing.T) {
	testCases := []struct {
		name           string
		match          SplitGRPCRouteMatch
		expectedTraits GRPCRoutePriorityTraits
	}{
		{
			name: "precise hostname with exact method match",
			match: SplitGRPCRouteMatch{
				Hostname: "petstore.com",
				Match: gatewayv1alpha2.GRPCRouteMatch{
					Method: &gatewayv1alpha2.GRPCMethodMatch{
						Type:    lo.ToPtr(gatewayv1alpha2.GRPCMethodMatchExact),
						Service: lo.ToPtr("pets"),
						Method:  lo.ToPtr("list"),
					},
				},
			},
			expectedTraits: GRPCRoutePriorityTraits{
				PreciseHostname: true,
				HostnameLength:  len("petstore.com"),
				MethodMatchType: gatewayv1alpha2.GRPCMethodMatchExact,
				ServiceLength:   len("pets"),
				MethodLength:    len("list"),
			},
		},
		{
			name: "wildcard hostname and partial method match",
			match: SplitGRPCRouteMatch{
				Hostname: "*.petstore.com",
				Match: gatewayv1alpha2.GRPCRouteMatch{
					Method: &gatewayv1alpha2.GRPCMethodMatch{
						Type:    lo.ToPtr(gatewayv1alpha2.GRPCMethodMatchExact),
						Service: lo.ToPtr("pets"),
					},
				},
			},
			expectedTraits: GRPCRoutePriorityTraits{
				PreciseHostname: false,
				HostnameLength:  len("*.petstore.com"),
				MethodMatchType: gatewayv1alpha2.GRPCMethodMatchExact,
				ServiceLength:   len("pets"),
				MethodLength:    0,
			},
		},
		{
			name: "no hostname with only header matches",
			match: SplitGRPCRouteMatch{
				Match: gatewayv1alpha2.GRPCRouteMatch{
					Headers: []gatewayv1alpha2.GRPCHeaderMatch{
						{
							Type:  lo.ToPtr(gatewayv1beta1.HeaderMatchExact),
							Name:  gatewayv1alpha2.GRPCHeaderName("key1"),
							Value: "value1",
						},
						{
							Type:  lo.ToPtr(gatewayv1beta1.HeaderMatchExact),
							Name:  gatewayv1alpha2.GRPCHeaderName("key2"),
							Value: "value2",
						},
					},
				},
			},
			expectedTraits: GRPCRoutePriorityTraits{
				HostnameLength: 0,
				ServiceLength:  0,
				MethodLength:   0,
				HeaderCount:    2,
			},
		},
	}

	for i, tc := range testCases {
		tc := tc
		indexStr := strconv.Itoa(i)
		t.Run(indexStr+"-"+tc.name, func(t *testing.T) {
			traits := CalculateGRCPRouteMatchPriorityTraits(tc.match)
			require.Equal(t, tc.expectedTraits, traits)
		})
	}
}

func TestGRPCRouteTraitsEncodeToPriority(t *testing.T) {
	testCases := []struct {
		name              string
		traits            GRPCRoutePriorityTraits
		exprectedPriority int
	}{
		{
			name: "precise hostname",
			traits: GRPCRoutePriorityTraits{
				PreciseHostname: true,
				HostnameLength:  15,
				ServiceLength:   7,
			},
			exprectedPriority: (1 << 50) | (1 << 49) | (15 << 41) | (7 << 30),
		},
		{
			name: "non precise hostname",
			traits: GRPCRoutePriorityTraits{
				PreciseHostname: false,
				HostnameLength:  15,
				ServiceLength:   7,
				MethodLength:    7,
				HeaderCount:     3,
			},
			exprectedPriority: (1 << 50) | (15 << 41) | (7 << 30) | (7 << 19) | (3 << 14),
		},
	}

	for i, tc := range testCases {
		tc := tc
		indexStr := strconv.Itoa(i)
		t.Run(indexStr+"-"+tc.name, func(t *testing.T) {
			priority := tc.traits.EncodeToPriority()
			require.Equal(t, tc.exprectedPriority, priority)
		})
	}
}

func TestAssignRoutePriorityToSplitGRPCRouteMatches(t *testing.T) {
	type splitGRPCRouteMatchIndex struct {
		namespace  string
		name       string
		hostname   string
		ruleIndex  int
		matchIndex int
	}
	now := time.Now()
	const maxRelativeOrderPriorityBits = (1 << 14) - 1

	testCases := []struct {
		name                  string
		splitGRPCRouteMatches []SplitGRPCRouteMatch
		// GRPCRoute index -> priority
		priorities map[splitGRPCRouteMatchIndex]int
	}{
		{
			name: "no dupelicated fixed priority",
			splitGRPCRouteMatches: []SplitGRPCRouteMatch{
				{
					Source: &gatewayv1alpha2.GRPCRoute{
						ObjectMeta: metav1.ObjectMeta{
							Namespace:         "default",
							Name:              "grpcroute-1",
							CreationTimestamp: metav1.NewTime(now.Add(-5 * time.Second)),
						},
						Spec: gatewayv1alpha2.GRPCRouteSpec{
							Hostnames: []gatewayv1beta1.Hostname{"foo.com"},
							Rules: []gatewayv1alpha2.GRPCRouteRule{
								{
									Matches: []gatewayv1alpha2.GRPCRouteMatch{
										{
											Method: &gatewayv1alpha2.GRPCMethodMatch{
												Service: lo.ToPtr("pets"),
												Method:  lo.ToPtr("list"),
											},
										},
									},
								},
							},
						},
					},
					Hostname: "foo.com",
					Match: gatewayv1alpha2.GRPCRouteMatch{
						Method: &gatewayv1alpha2.GRPCMethodMatch{
							Service: lo.ToPtr("pets"),
							Method:  lo.ToPtr("list"),
						},
					},
					RuleIndex:  0,
					MatchIndex: 0,
				},
				{
					Source: &gatewayv1alpha2.GRPCRoute{
						ObjectMeta: metav1.ObjectMeta{
							Namespace:         "default",
							Name:              "grpcroute-2",
							CreationTimestamp: metav1.NewTime(now.Add(-5 * time.Second)),
						},
						Spec: gatewayv1alpha2.GRPCRouteSpec{
							Hostnames: []gatewayv1beta1.Hostname{"*.bar.com"},
							Rules: []gatewayv1alpha2.GRPCRouteRule{
								{
									Matches: []gatewayv1alpha2.GRPCRouteMatch{
										{
											Method: &gatewayv1alpha2.GRPCMethodMatch{
												Service: lo.ToPtr("pets"),
												Method:  lo.ToPtr("list"),
											},
										},
									},
								},
							},
						},
					},
					Hostname: "*.bar.com",
					Match: gatewayv1alpha2.GRPCRouteMatch{
						Method: &gatewayv1alpha2.GRPCMethodMatch{
							Service: lo.ToPtr("pets"),
							Method:  lo.ToPtr("list"),
						},
					},
				},
			},
			priorities: map[splitGRPCRouteMatchIndex]int{
				{
					namespace:  "default",
					name:       "grpcroute-1",
					hostname:   "foo.com",
					ruleIndex:  0,
					matchIndex: 0,
				}: GRPCRoutePriorityTraits{
					PreciseHostname: true,
					HostnameLength:  len("foo.com"),
					ServiceLength:   len("pets"),
					MethodLength:    len("list"),
				}.EncodeToPriority() + maxRelativeOrderPriorityBits,
				{
					namespace:  "default",
					name:       "grpcroute-2",
					hostname:   "*.bar.com",
					ruleIndex:  0,
					matchIndex: 0,
				}: GRPCRoutePriorityTraits{
					PreciseHostname: false,
					HostnameLength:  len("*.bar.com"),
					ServiceLength:   len("pets"),
					MethodLength:    len("list"),
				}.EncodeToPriority() + maxRelativeOrderPriorityBits,
			},
		},
		{
			name: "break tie by creation timestamp",
			splitGRPCRouteMatches: []SplitGRPCRouteMatch{
				{
					Source: &gatewayv1alpha2.GRPCRoute{
						ObjectMeta: metav1.ObjectMeta{
							Namespace:         "default",
							Name:              "grpcroute-1",
							CreationTimestamp: metav1.NewTime(now.Add(-5 * time.Second)),
						},
						Spec: gatewayv1alpha2.GRPCRouteSpec{
							Hostnames: []gatewayv1beta1.Hostname{"foo.com"},
							Rules: []gatewayv1alpha2.GRPCRouteRule{
								{
									Matches: []gatewayv1alpha2.GRPCRouteMatch{
										{
											Method: &gatewayv1alpha2.GRPCMethodMatch{
												Service: lo.ToPtr("pets"),
												Method:  lo.ToPtr("list"),
											},
										},
									},
								},
							},
						},
					},
					Hostname: "foo.com",
					Match: gatewayv1alpha2.GRPCRouteMatch{
						Method: &gatewayv1alpha2.GRPCMethodMatch{
							Service: lo.ToPtr("pets"),
							Method:  lo.ToPtr("list"),
						},
					},
					RuleIndex:  0,
					MatchIndex: 0,
				},
				{
					Source: &gatewayv1alpha2.GRPCRoute{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "grpcroute-2",
							// created earlier
							CreationTimestamp: metav1.NewTime(now.Add(-10 * time.Second)),
						},
						Spec: gatewayv1alpha2.GRPCRouteSpec{
							Hostnames: []gatewayv1beta1.Hostname{"bar.com"},
							Rules: []gatewayv1alpha2.GRPCRouteRule{
								{
									Matches: []gatewayv1alpha2.GRPCRouteMatch{
										{
											Method: &gatewayv1alpha2.GRPCMethodMatch{
												Service: lo.ToPtr("pets"),
												Method:  lo.ToPtr("list"),
											},
										},
									},
								},
							},
						},
					},
					Hostname: "bar.com",
					Match: gatewayv1alpha2.GRPCRouteMatch{
						Method: &gatewayv1alpha2.GRPCMethodMatch{
							Service: lo.ToPtr("pets"),
							Method:  lo.ToPtr("list"),
						},
					},
					RuleIndex:  0,
					MatchIndex: 0,
				},
			},
			priorities: map[splitGRPCRouteMatchIndex]int{
				{
					namespace:  "default",
					name:       "grpcroute-1",
					hostname:   "foo.com",
					ruleIndex:  0,
					matchIndex: 0,
				}: GRPCRoutePriorityTraits{
					PreciseHostname: true,
					HostnameLength:  len("foo.com"),
					ServiceLength:   len("pets"),
					MethodLength:    len("list"),
				}.EncodeToPriority() + maxRelativeOrderPriorityBits - 1,
				{
					namespace:  "default",
					name:       "grpcroute-2",
					hostname:   "bar.com",
					ruleIndex:  0,
					matchIndex: 0,
				}: GRPCRoutePriorityTraits{
					PreciseHostname: true,
					HostnameLength:  len("bar.com"),
					ServiceLength:   len("pets"),
					MethodLength:    len("list"),
				}.EncodeToPriority() + maxRelativeOrderPriorityBits,
			},
		},
		{
			name: "break tie by name",
			splitGRPCRouteMatches: []SplitGRPCRouteMatch{
				{
					Source: &gatewayv1alpha2.GRPCRoute{
						ObjectMeta: metav1.ObjectMeta{
							Namespace:         "default",
							Name:              "grpcroute-1",
							CreationTimestamp: metav1.NewTime(now.Add(-5 * time.Second)),
						},
						Spec: gatewayv1alpha2.GRPCRouteSpec{
							Hostnames: []gatewayv1beta1.Hostname{"foo.com"},
							Rules: []gatewayv1alpha2.GRPCRouteRule{
								{
									Matches: []gatewayv1alpha2.GRPCRouteMatch{
										{
											Method: &gatewayv1alpha2.GRPCMethodMatch{
												Service: lo.ToPtr("pets"),
												Method:  lo.ToPtr("list"),
											},
										},
									},
								},
							},
						},
					},
					Hostname: "foo.com",
					Match: gatewayv1alpha2.GRPCRouteMatch{
						Method: &gatewayv1alpha2.GRPCMethodMatch{
							Service: lo.ToPtr("pets"),
							Method:  lo.ToPtr("list"),
						},
					},
					RuleIndex:  0,
					MatchIndex: 0,
				},
				{
					Source: &gatewayv1alpha2.GRPCRoute{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "grpcroute-2",
							Annotations: map[string]string{
								InternalRuleIndexAnnotationKey:  "0",
								InternalMatchIndexAnnotationKey: "0",
							},
							CreationTimestamp: metav1.NewTime(now.Add(-5 * time.Second)),
						},
						Spec: gatewayv1alpha2.GRPCRouteSpec{
							Hostnames: []gatewayv1beta1.Hostname{"bar.com"},
							Rules: []gatewayv1alpha2.GRPCRouteRule{
								{
									Matches: []gatewayv1alpha2.GRPCRouteMatch{
										{
											Method: &gatewayv1alpha2.GRPCMethodMatch{
												Service: lo.ToPtr("pets"),
												Method:  lo.ToPtr("list"),
											},
										},
									},
								},
							},
						},
					},
					Hostname: "bar.com",
					Match: gatewayv1alpha2.GRPCRouteMatch{
						Method: &gatewayv1alpha2.GRPCMethodMatch{
							Service: lo.ToPtr("pets"),
							Method:  lo.ToPtr("list"),
						},
					},
					RuleIndex:  0,
					MatchIndex: 0,
				},
			},
			priorities: map[splitGRPCRouteMatchIndex]int{
				{
					namespace:  "default",
					name:       "grpcroute-1",
					hostname:   "foo.com",
					ruleIndex:  0,
					matchIndex: 0,
				}: GRPCRoutePriorityTraits{
					PreciseHostname: true,
					HostnameLength:  len("foo.com"),
					ServiceLength:   len("pets"),
					MethodLength:    len("list"),
				}.EncodeToPriority() + maxRelativeOrderPriorityBits,
				{
					namespace:  "default",
					name:       "grpcroute-2",
					hostname:   "bar.com",
					ruleIndex:  0,
					matchIndex: 0,
				}: GRPCRoutePriorityTraits{
					PreciseHostname: true,
					HostnameLength:  len("bar.com"),
					ServiceLength:   len("pets"),
					MethodLength:    len("list"),
				}.EncodeToPriority() + maxRelativeOrderPriorityBits - 1,
			},
		},
		{
			name: "break tie by internal match order",
			splitGRPCRouteMatches: []SplitGRPCRouteMatch{
				{
					Source: &gatewayv1alpha2.GRPCRoute{
						ObjectMeta: metav1.ObjectMeta{
							Namespace:         "default",
							Name:              "grpcroute-1",
							CreationTimestamp: metav1.NewTime(now.Add(-5 * time.Second)),
						},
						Spec: gatewayv1alpha2.GRPCRouteSpec{
							Hostnames: []gatewayv1beta1.Hostname{"foo.com"},
							Rules: []gatewayv1alpha2.GRPCRouteRule{
								{
									Matches: []gatewayv1alpha2.GRPCRouteMatch{
										{
											Method: &gatewayv1alpha2.GRPCMethodMatch{
												Service: lo.ToPtr("cats"),
												Method:  lo.ToPtr("list"),
											},
										},
										{
											Method: &gatewayv1alpha2.GRPCMethodMatch{
												Service: lo.ToPtr("dogs"),
												Method:  lo.ToPtr("list"),
											},
										},
									},
								},
							},
						},
					},
					Hostname: "foo.com",
					Match: gatewayv1alpha2.GRPCRouteMatch{
						Method: &gatewayv1alpha2.GRPCMethodMatch{
							Service: lo.ToPtr("cats"),
							Method:  lo.ToPtr("list"),
						},
					},
					RuleIndex:  0,
					MatchIndex: 0,
				},
				{
					// The same as the one above
					Source: &gatewayv1alpha2.GRPCRoute{
						ObjectMeta: metav1.ObjectMeta{
							Namespace:         "default",
							Name:              "grpcroute-1",
							CreationTimestamp: metav1.NewTime(now.Add(-5 * time.Second)),
						},
						Spec: gatewayv1alpha2.GRPCRouteSpec{
							Hostnames: []gatewayv1beta1.Hostname{"foo.com"},
							Rules: []gatewayv1alpha2.GRPCRouteRule{
								{
									Matches: []gatewayv1alpha2.GRPCRouteMatch{
										{
											Method: &gatewayv1alpha2.GRPCMethodMatch{
												Service: lo.ToPtr("cats"),
												Method:  lo.ToPtr("list"),
											},
										},
										{
											Method: &gatewayv1alpha2.GRPCMethodMatch{
												Service: lo.ToPtr("dogs"),
												Method:  lo.ToPtr("list"),
											},
										},
									},
								},
							},
						},
					},
					Hostname: "foo.com",
					Match: gatewayv1alpha2.GRPCRouteMatch{
						Method: &gatewayv1alpha2.GRPCMethodMatch{
							Service: lo.ToPtr("dogs"),
							Method:  lo.ToPtr("list"),
						},
					},
					RuleIndex:  0,
					MatchIndex: 1,
				},
			},
			priorities: map[splitGRPCRouteMatchIndex]int{
				{
					namespace:  "default",
					name:       "grpcroute-1",
					hostname:   "foo.com",
					ruleIndex:  0,
					matchIndex: 0,
				}: GRPCRoutePriorityTraits{
					PreciseHostname: true,
					HostnameLength:  len("foo.com"),
					ServiceLength:   len("cats"),
					MethodLength:    len("list"),
				}.EncodeToPriority() + maxRelativeOrderPriorityBits,
				{
					namespace:  "default",
					name:       "grpcroute-1",
					hostname:   "foo.com",
					ruleIndex:  0,
					matchIndex: 1,
				}: GRPCRoutePriorityTraits{
					PreciseHostname: true,
					HostnameLength:  len("foo.com"),
					ServiceLength:   len("dogs"),
					MethodLength:    len("list"),
				}.EncodeToPriority() + maxRelativeOrderPriorityBits - 1,
			},
		},
	}

	for i, tc := range testCases {
		indexStr := strconv.Itoa(i)
		tc := tc
		t.Run(indexStr+"-"+tc.name, func(t *testing.T) {
			splitMatchesWithPriorities := AssignRoutePriorityToSplitGRPCRouteMatches(logr.Discard(), tc.splitGRPCRouteMatches)
			require.Lenf(t, splitMatchesWithPriorities, len(tc.priorities),
				"should have expeceted number of results")
			for _, m := range splitMatchesWithPriorities {
				grpcRoute := m.Match.Source

				require.Equalf(t, tc.priorities[splitGRPCRouteMatchIndex{
					namespace:  grpcRoute.Namespace,
					name:       grpcRoute.Name,
					hostname:   string(grpcRoute.Spec.Hostnames[0]),
					ruleIndex:  m.Match.RuleIndex,
					matchIndex: m.Match.MatchIndex,
				}], m.Priority,
					"grpcroute %s/%s: hostname %s, rule %d match %d does not have expected priority",
					grpcRoute.Namespace, grpcRoute.Name, grpcRoute.Spec.Hostnames[0], m.Match.RuleIndex, m.Match.MatchIndex)
			}
		})
	}
}

func TestKongExpressionRouteFromSplitGRPCRouteWithPriority(t *testing.T) {
	testCases := []struct {
		name                       string
		splitGRPCMatchWithPriority SplitGRPCRouteMatchToPriority
		expectedRoute              kongstate.Route
	}{
		{
			name: "no host and exact method match",
			splitGRPCMatchWithPriority: SplitGRPCRouteMatchToPriority{
				Match: SplitGRPCRouteMatch{
					Source: &gatewayv1alpha2.GRPCRoute{
						TypeMeta: grpcRouteTypeMeta,
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "no-hostname-exact-method",
						},
						Spec: gatewayv1alpha2.GRPCRouteSpec{
							Rules: []gatewayv1alpha2.GRPCRouteRule{
								{
									Matches: []gatewayv1alpha2.GRPCRouteMatch{
										{
											Method: &gatewayv1alpha2.GRPCMethodMatch{
												Type:    lo.ToPtr(gatewayv1alpha2.GRPCMethodMatchExact),
												Service: lo.ToPtr("pets"),
												Method:  lo.ToPtr("list"),
											},
										},
									},
								},
							},
						},
					},
					Match: gatewayv1alpha2.GRPCRouteMatch{
						Method: &gatewayv1alpha2.GRPCMethodMatch{
							Type:    lo.ToPtr(gatewayv1alpha2.GRPCMethodMatchExact),
							Service: lo.ToPtr("pets"),
							Method:  lo.ToPtr("list"),
						},
					},
					RuleIndex:  0,
					MatchIndex: 0,
				},
				Priority: 1024,
			},
			expectedRoute: kongstate.Route{
				Route: kong.Route{
					Name:         kong.String("grpcroute.default.no-hostname-exact-method._.0.0"),
					PreserveHost: kong.Bool(true),
					Expression:   kong.String(`http.path == "/pets/list"`),
					Priority:     kong.Int(1024),
				},
				ExpressionRoutes: true,
			},
		},
		{
			name: "precise hostname and regex method match",
			splitGRPCMatchWithPriority: SplitGRPCRouteMatchToPriority{
				Match: SplitGRPCRouteMatch{
					Source: &gatewayv1alpha2.GRPCRoute{
						TypeMeta: grpcRouteTypeMeta,
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "precise-hostname-regex-method",
							Annotations: map[string]string{
								InternalRuleIndexAnnotationKey:  "0",
								InternalMatchIndexAnnotationKey: "0",
							},
						},
						Spec: gatewayv1alpha2.GRPCRouteSpec{
							Hostnames: []gatewayv1alpha2.Hostname{
								"foo.com",
							},
							Rules: []gatewayv1alpha2.GRPCRouteRule{
								{
									Matches: []gatewayv1alpha2.GRPCRouteMatch{
										{
											Method: &gatewayv1alpha2.GRPCMethodMatch{
												Type:    lo.ToPtr(gatewayv1alpha2.GRPCMethodMatchRegularExpression),
												Service: lo.ToPtr("name"),
												Method:  lo.ToPtr("[a-z0-9]+"),
											},
										},
									},
								},
							},
						},
					},
					Hostname: "foo.com",
					Match: gatewayv1alpha2.GRPCRouteMatch{
						Method: &gatewayv1alpha2.GRPCMethodMatch{
							Type:    lo.ToPtr(gatewayv1alpha2.GRPCMethodMatchRegularExpression),
							Service: lo.ToPtr("name"),
							Method:  lo.ToPtr("[a-z0-9]+"),
						},
					},
					RuleIndex:  0,
					MatchIndex: 0,
				},
				Priority: 1024,
			},
			expectedRoute: kongstate.Route{
				Route: kong.Route{
					Name:         kong.String("grpcroute.default.precise-hostname-regex-method.foo.com.0.0"),
					Expression:   kong.String(`(http.path ~ "^/name/[a-z0-9]+") && (http.host == "foo.com")`),
					PreserveHost: kong.Bool(true),
					Priority:     kong.Int(1024),
				},
				ExpressionRoutes: true,
			},
		},
		{
			name: "wildcard hostname and header match",
			splitGRPCMatchWithPriority: SplitGRPCRouteMatchToPriority{
				Match: SplitGRPCRouteMatch{
					Source: &gatewayv1alpha2.GRPCRoute{
						TypeMeta: grpcRouteTypeMeta,
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "wildcard-hostname-header-match",
						},
						Spec: gatewayv1alpha2.GRPCRouteSpec{
							Hostnames: []gatewayv1alpha2.Hostname{
								"*.foo.com",
							},
							Rules: []gatewayv1alpha2.GRPCRouteRule{
								{
									Matches: []gatewayv1alpha2.GRPCRouteMatch{
										{
											Method: &gatewayv1alpha2.GRPCMethodMatch{
												Type:    lo.ToPtr(gatewayv1alpha2.GRPCMethodMatchExact),
												Service: lo.ToPtr("names"),
												Method:  lo.ToPtr("create"),
											},
										},
										{
											Method: &gatewayv1alpha2.GRPCMethodMatch{
												Type:    lo.ToPtr(gatewayv1alpha2.GRPCMethodMatchExact),
												Service: lo.ToPtr("name"),
											},
											Headers: []gatewayv1alpha2.GRPCHeaderMatch{
												{
													Name:  gatewayv1alpha2.GRPCHeaderName("foo"),
													Value: "bar",
												},
											},
										},
									},
								},
							},
						},
					},
					Hostname: "*.foo.com",
					Match: gatewayv1alpha2.GRPCRouteMatch{
						Method: &gatewayv1alpha2.GRPCMethodMatch{
							Type:    lo.ToPtr(gatewayv1alpha2.GRPCMethodMatchExact),
							Service: lo.ToPtr("name"),
						},
						Headers: []gatewayv1alpha2.GRPCHeaderMatch{
							{
								Name:  gatewayv1alpha2.GRPCHeaderName("foo"),
								Value: "bar",
							},
						},
					},
					RuleIndex:  0,
					MatchIndex: 1,
				},
				Priority: 1024,
			},
			expectedRoute: kongstate.Route{
				Route: kong.Route{
					Name:         kong.String("grpcroute.default.wildcard-hostname-header-match._.foo.com.0.1"),
					Expression:   kong.String(`(http.path ^= "/name/") && (http.headers.foo == "bar") && (http.host =^ ".foo.com")`),
					PreserveHost: kong.Bool(true),
					Priority:     kong.Int(1024),
				},
				ExpressionRoutes: true,
			},
		},
	}

	for i, tc := range testCases {
		indexStr := strconv.Itoa(i)
		tc := tc
		t.Run(indexStr+"-"+tc.name, func(t *testing.T) {
			r := KongExpressionRouteFromSplitGRPCRouteMatchWithPriority(tc.splitGRPCMatchWithPriority)
			grpcRoute := tc.splitGRPCMatchWithPriority.Match.Source
			tc.expectedRoute.Route.Tags = util.GenerateTagsForObject(grpcRoute)
			require.Equal(t, tc.expectedRoute.Route, r.Route)
			require.True(t, r.ExpressionRoutes)
			require.Equal(t, grpcRoute.Namespace, r.Ingress.Namespace)
			require.Equal(t, grpcRoute.Name, r.Ingress.Name)
		})
	}
}
