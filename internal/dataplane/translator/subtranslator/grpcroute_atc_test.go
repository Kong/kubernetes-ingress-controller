package subtranslator

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

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

func TestGenerateKongExpressionRoutesFromGRPCRouteRule(t *testing.T) {
	testCases := []struct {
		name           string
		objectName     string
		annotations    map[string]string
		hostnames      []string
		rule           gatewayapi.GRPCRouteRule
		expectedRoutes []kongstate.Route
	}{
		{
			name:        "single match without hostname",
			objectName:  "single-match",
			annotations: map[string]string{},
			rule: gatewayapi.GRPCRouteRule{
				Matches: []gatewayapi.GRPCRouteMatch{
					{
						Method: &gatewayapi.GRPCMethodMatch{
							Service: lo.ToPtr("service0"),
							Method:  nil,
						},
						Headers: []gatewayapi.GRPCHeaderMatch{
							{
								Name:  gatewayapi.GRPCHeaderName("X-Foo"),
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
						Priority:   kong.Uint64(1),
					},
				},
			},
		},
		{
			name:        "single match with hostname",
			objectName:  "single-match-with-hostname",
			annotations: map[string]string{},
			hostnames:   []string{"foo.com", "*.foo.com"},
			rule: gatewayapi.GRPCRouteRule{
				Matches: []gatewayapi.GRPCRouteMatch{
					{
						Method: &gatewayapi.GRPCMethodMatch{
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
						Priority:   kong.Uint64(1),
					},
				},
			},
		},
		{
			name:       "multuple matches without hostname",
			objectName: "multiple-matches",
			rule: gatewayapi.GRPCRouteRule{
				Matches: []gatewayapi.GRPCRouteMatch{
					{
						Method: &gatewayapi.GRPCMethodMatch{
							Service: nil,
							Method:  lo.ToPtr("method0"),
						},
						Headers: []gatewayapi.GRPCHeaderMatch{
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
						Method: &gatewayapi.GRPCMethodMatch{
							Type:    lo.ToPtr(gatewayapi.GRPCMethodMatchRegularExpression),
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
						GroupVersionKind: grpcRouteGVK,
					},
					ExpressionRoutes: true,
					Route: kong.Route{
						Name:       kong.String("grpcroute.default.multiple-matches.0.0"),
						Expression: kong.String(`(http.path =^ "/method0") && ((http.headers.client == "kong-test") && (http.headers.version == "2"))`),
						Priority:   kong.Uint64(1),
					},
				},
				{
					Ingress: util.K8sObjectInfo{
						Name:             "multiple-matches",
						Namespace:        "default",
						GroupVersionKind: grpcRouteGVK,
					},
					ExpressionRoutes: true,
					Route: kong.Route{
						Name:       kong.String("grpcroute.default.multiple-matches.0.1"),
						Expression: kong.String(`http.path ~ "^/v[012]/.+"`),
						Priority:   kong.Uint64(1),
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
			rule: gatewayapi.GRPCRouteRule{
				Matches: []gatewayapi.GRPCRouteMatch{
					{
						Method: &gatewayapi.GRPCMethodMatch{
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
						Expression: kong.String(`(http.path == "/service0/method0") && (tls.sni == "kong.foo.com")`),
						Priority:   kong.Uint64(1),
					},
				},
			},
		},
		{
			name:       "single hostname with no match",
			objectName: "hostname-only",
			hostnames:  []string{"foo.com"},
			rule: gatewayapi.GRPCRouteRule{
				Matches: []gatewayapi.GRPCRouteMatch{},
			},
			expectedRoutes: []kongstate.Route{
				{
					Ingress: util.K8sObjectInfo{
						Name:             "hostname-only",
						Namespace:        "default",
						GroupVersionKind: grpcRouteGVK,
					},
					ExpressionRoutes: true,
					Route: kong.Route{
						Name:       kong.String("grpcroute.default.hostname-only.0.0"),
						Expression: kong.String(`http.host == "foo.com"`),
						Priority:   kong.Uint64(1),
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			grpcroute := makeTestGRPCRoute(tc.objectName, "default", tc.annotations, tc.hostnames, []gatewayapi.GRPCRouteRule{tc.rule}, nil)
			routes := GenerateKongExpressionRoutesFromGRPCRouteRule(grpcroute, 0)
			require.Equal(t, tc.expectedRoutes, routes)
		})
	}
}

func TestMethodMatcherFromGRPCMethodMatch(t *testing.T) {
	testCases := []struct {
		name       string
		match      gatewayapi.GRPCMethodMatch
		expression string
	}{
		{
			name: "exact method match",
			match: gatewayapi.GRPCMethodMatch{
				Type:    lo.ToPtr(gatewayapi.GRPCMethodMatchExact),
				Service: lo.ToPtr("service0"),
				Method:  lo.ToPtr("method0"),
			},
			expression: `http.path == "/service0/method0"`,
		},
		{
			name: "exact match with service unspecified",
			match: gatewayapi.GRPCMethodMatch{
				Type:   lo.ToPtr(gatewayapi.GRPCMethodMatchExact),
				Method: lo.ToPtr("method0"),
			},
			expression: `http.path =^ "/method0"`,
		},
		{
			name: "regex method match",
			match: gatewayapi.GRPCMethodMatch{
				Type:    lo.ToPtr(gatewayapi.GRPCMethodMatchRegularExpression),
				Service: lo.ToPtr("auth-v[012]"),
				Method:  lo.ToPtr("[a-z]+"),
			},
			expression: `http.path ~ "^/auth-v[012]/[a-z]+"`,
		},
		{
			name: "empty regex match",
			match: gatewayapi.GRPCMethodMatch{
				Type: lo.ToPtr(gatewayapi.GRPCMethodMatchRegularExpression),
			},
			expression: `http.path ^= "/"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expression, methodMatcherFromGRPCMethodMatch(&tc.match).Expression())
		})
	}
}

func TestSplitGRPCRoute(t *testing.T) {
	namesToBackendRefs := func(names []string) []gatewayapi.GRPCBackendRef {
		backendRefs := []gatewayapi.GRPCBackendRef{}
		for _, name := range names {
			backendRefs = append(backendRefs,
				gatewayapi.GRPCBackendRef{
					BackendRef: gatewayapi.BackendRef{
						BackendObjectReference: gatewayapi.BackendObjectReference{
							Name: gatewayapi.ObjectName(name),
						},
					},
				},
			)
		}
		return backendRefs
	}

	testGRPCRoutes := []*gatewayapi.GRPCRoute{
		{
			TypeMeta: grpcRouteTypeMeta,
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "grpcroute-no-hostname-one-match",
			},
			Spec: gatewayapi.GRPCRouteSpec{
				Rules: []gatewayapi.GRPCRouteRule{
					{
						Matches: []gatewayapi.GRPCRouteMatch{
							{
								Method: &gatewayapi.GRPCMethodMatch{
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
			Spec: gatewayapi.GRPCRouteSpec{
				Hostnames: []gatewayapi.Hostname{
					gatewayapi.Hostname("petstore.com"),
				},
				Rules: []gatewayapi.GRPCRouteRule{
					{
						Matches: []gatewayapi.GRPCRouteMatch{
							{
								Method: &gatewayapi.GRPCMethodMatch{
									Service: lo.ToPtr("cats"),
									Method:  lo.ToPtr("list"),
								},
							},
							{
								Method: &gatewayapi.GRPCMethodMatch{
									Service: lo.ToPtr("dogs"),
									Method:  lo.ToPtr("list"),
								},
							},
						},
						BackendRefs: namesToBackendRefs([]string{"listpets"}),
					},
					{
						Matches: []gatewayapi.GRPCRouteMatch{
							{
								Method: &gatewayapi.GRPCMethodMatch{
									Service: lo.ToPtr("cats"),
									Method:  lo.ToPtr("create"),
								},
							},
							{
								Method: &gatewayapi.GRPCMethodMatch{
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
			Spec: gatewayapi.GRPCRouteSpec{
				Hostnames: []gatewayapi.Hostname{
					gatewayapi.Hostname("petstore.com"),
					gatewayapi.Hostname("petstore.net"),
				},
				Rules: []gatewayapi.GRPCRouteRule{
					{
						Matches: []gatewayapi.GRPCRouteMatch{
							{
								Method: &gatewayapi.GRPCMethodMatch{
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
				Name:      "grpcroute-multiple-hostnames-no-match",
			},
			Spec: gatewayapi.GRPCRouteSpec{
				Hostnames: []gatewayapi.Hostname{
					gatewayapi.Hostname("pets.com"),
					gatewayapi.Hostname("pets.net"),
				},
				Rules: []gatewayapi.GRPCRouteRule{
					{
						BackendRefs: namesToBackendRefs([]string{"listpets"}),
					},
				},
			},
		},
	}

	testCases := []struct {
		name                 string
		grpcRoute            *gatewayapi.GRPCRoute
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
		{
			name:      "multiple hostnames and no match",
			grpcRoute: testGRPCRoutes[3],
			expectedSplitMatches: []SplitGRPCRouteMatch{
				{
					Source:     testGRPCRoutes[3],
					Hostname:   string(testGRPCRoutes[3].Spec.Hostnames[0]),
					Match:      gatewayapi.GRPCRouteMatch{},
					RuleIndex:  0,
					MatchIndex: 0,
				},
				{
					Source:     testGRPCRoutes[3],
					Hostname:   string(testGRPCRoutes[3].Spec.Hostnames[1]),
					Match:      gatewayapi.GRPCRouteMatch{},
					RuleIndex:  0,
					MatchIndex: 0,
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
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
				Match: gatewayapi.GRPCRouteMatch{
					Method: &gatewayapi.GRPCMethodMatch{
						Type:    lo.ToPtr(gatewayapi.GRPCMethodMatchExact),
						Service: lo.ToPtr("pets"),
						Method:  lo.ToPtr("list"),
					},
				},
			},
			expectedTraits: GRPCRoutePriorityTraits{
				PreciseHostname: true,
				HostnameLength:  len("petstore.com"),
				MethodMatchType: gatewayapi.GRPCMethodMatchExact,
				ServiceLength:   len("pets"),
				MethodLength:    len("list"),
			},
		},
		{
			name: "wildcard hostname and partial method match",
			match: SplitGRPCRouteMatch{
				Hostname: "*.petstore.com",
				Match: gatewayapi.GRPCRouteMatch{
					Method: &gatewayapi.GRPCMethodMatch{
						Type:    lo.ToPtr(gatewayapi.GRPCMethodMatchExact),
						Service: lo.ToPtr("pets"),
					},
				},
			},
			expectedTraits: GRPCRoutePriorityTraits{
				PreciseHostname: false,
				HostnameLength:  len("*.petstore.com"),
				MethodMatchType: gatewayapi.GRPCMethodMatchExact,
				ServiceLength:   len("pets"),
				MethodLength:    0,
			},
		},
		{
			name: "no hostname with only header matches",
			match: SplitGRPCRouteMatch{
				Match: gatewayapi.GRPCRouteMatch{
					Headers: []gatewayapi.GRPCHeaderMatch{
						{
							Type:  lo.ToPtr(gatewayapi.GRPCHeaderMatchExact),
							Name:  gatewayapi.GRPCHeaderName("key1"),
							Value: "value1",
						},
						{
							Type:  lo.ToPtr(gatewayapi.GRPCHeaderMatchExact),
							Name:  gatewayapi.GRPCHeaderName("key2"),
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
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			traits := CalculateGRCPRouteMatchPriorityTraits(tc.match)
			require.Equal(t, tc.expectedTraits, traits)
		})
	}
}

func TestGRPCRouteTraitsEncodeToPriority(t *testing.T) {
	testCases := []struct {
		name              string
		traits            GRPCRoutePriorityTraits
		exprectedPriority RoutePriorityType
	}{
		{
			name: "precise hostname",
			traits: GRPCRoutePriorityTraits{
				PreciseHostname: true,
				HostnameLength:  15,
				ServiceLength:   7,
			},
			exprectedPriority: (1 << 44) | (1 << 43) | (15 << 35) | (7 << 24),
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
			exprectedPriority: (1 << 44) | (15 << 35) | (7 << 24) | (7 << 13) | (3 << 8),
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
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
	const maxRelativeOrderPriorityBits = (1 << 8) - 1

	testCases := []struct {
		name                  string
		splitGRPCRouteMatches []SplitGRPCRouteMatch
		// GRPCRoute index -> priority
		priorities map[splitGRPCRouteMatchIndex]RoutePriorityType
	}{
		{
			name: "no dupelicated fixed priority",
			splitGRPCRouteMatches: []SplitGRPCRouteMatch{
				{
					Source: &gatewayapi.GRPCRoute{
						ObjectMeta: metav1.ObjectMeta{
							Namespace:         "default",
							Name:              "grpcroute-1",
							CreationTimestamp: metav1.NewTime(now.Add(-5 * time.Second)),
						},
						Spec: gatewayapi.GRPCRouteSpec{
							Hostnames: []gatewayapi.Hostname{"foo.com"},
							Rules: []gatewayapi.GRPCRouteRule{
								{
									Matches: []gatewayapi.GRPCRouteMatch{
										{
											Method: &gatewayapi.GRPCMethodMatch{
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
					Match: gatewayapi.GRPCRouteMatch{
						Method: &gatewayapi.GRPCMethodMatch{
							Service: lo.ToPtr("pets"),
							Method:  lo.ToPtr("list"),
						},
					},
					RuleIndex:  0,
					MatchIndex: 0,
				},
				{
					Source: &gatewayapi.GRPCRoute{
						ObjectMeta: metav1.ObjectMeta{
							Namespace:         "default",
							Name:              "grpcroute-2",
							CreationTimestamp: metav1.NewTime(now.Add(-5 * time.Second)),
						},
						Spec: gatewayapi.GRPCRouteSpec{
							Hostnames: []gatewayapi.Hostname{"*.bar.com"},
							Rules: []gatewayapi.GRPCRouteRule{
								{
									Matches: []gatewayapi.GRPCRouteMatch{
										{
											Method: &gatewayapi.GRPCMethodMatch{
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
					Match: gatewayapi.GRPCRouteMatch{
						Method: &gatewayapi.GRPCMethodMatch{
							Service: lo.ToPtr("pets"),
							Method:  lo.ToPtr("list"),
						},
					},
				},
			},
			priorities: map[splitGRPCRouteMatchIndex]RoutePriorityType{
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
					Source: &gatewayapi.GRPCRoute{
						ObjectMeta: metav1.ObjectMeta{
							Namespace:         "default",
							Name:              "grpcroute-1",
							CreationTimestamp: metav1.NewTime(now.Add(-5 * time.Second)),
						},
						Spec: gatewayapi.GRPCRouteSpec{
							Hostnames: []gatewayapi.Hostname{"foo.com"},
							Rules: []gatewayapi.GRPCRouteRule{
								{
									Matches: []gatewayapi.GRPCRouteMatch{
										{
											Method: &gatewayapi.GRPCMethodMatch{
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
					Match: gatewayapi.GRPCRouteMatch{
						Method: &gatewayapi.GRPCMethodMatch{
							Service: lo.ToPtr("pets"),
							Method:  lo.ToPtr("list"),
						},
					},
					RuleIndex:  0,
					MatchIndex: 0,
				},
				{
					Source: &gatewayapi.GRPCRoute{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "grpcroute-2",
							// created earlier
							CreationTimestamp: metav1.NewTime(now.Add(-10 * time.Second)),
						},
						Spec: gatewayapi.GRPCRouteSpec{
							Hostnames: []gatewayapi.Hostname{"bar.com"},
							Rules: []gatewayapi.GRPCRouteRule{
								{
									Matches: []gatewayapi.GRPCRouteMatch{
										{
											Method: &gatewayapi.GRPCMethodMatch{
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
					Match: gatewayapi.GRPCRouteMatch{
						Method: &gatewayapi.GRPCMethodMatch{
							Service: lo.ToPtr("pets"),
							Method:  lo.ToPtr("list"),
						},
					},
					RuleIndex:  0,
					MatchIndex: 0,
				},
			},
			priorities: map[splitGRPCRouteMatchIndex]RoutePriorityType{
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
					Source: &gatewayapi.GRPCRoute{
						ObjectMeta: metav1.ObjectMeta{
							Namespace:         "default",
							Name:              "grpcroute-1",
							CreationTimestamp: metav1.NewTime(now.Add(-5 * time.Second)),
						},
						Spec: gatewayapi.GRPCRouteSpec{
							Hostnames: []gatewayapi.Hostname{"foo.com"},
							Rules: []gatewayapi.GRPCRouteRule{
								{
									Matches: []gatewayapi.GRPCRouteMatch{
										{
											Method: &gatewayapi.GRPCMethodMatch{
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
					Match: gatewayapi.GRPCRouteMatch{
						Method: &gatewayapi.GRPCMethodMatch{
							Service: lo.ToPtr("pets"),
							Method:  lo.ToPtr("list"),
						},
					},
					RuleIndex:  0,
					MatchIndex: 0,
				},
				{
					Source: &gatewayapi.GRPCRoute{
						ObjectMeta: metav1.ObjectMeta{
							Namespace:         "default",
							Name:              "grpcroute-2",
							CreationTimestamp: metav1.NewTime(now.Add(-5 * time.Second)),
						},
						Spec: gatewayapi.GRPCRouteSpec{
							Hostnames: []gatewayapi.Hostname{"bar.com"},
							Rules: []gatewayapi.GRPCRouteRule{
								{
									Matches: []gatewayapi.GRPCRouteMatch{
										{
											Method: &gatewayapi.GRPCMethodMatch{
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
					Match: gatewayapi.GRPCRouteMatch{
						Method: &gatewayapi.GRPCMethodMatch{
							Service: lo.ToPtr("pets"),
							Method:  lo.ToPtr("list"),
						},
					},
					RuleIndex:  0,
					MatchIndex: 0,
				},
			},
			priorities: map[splitGRPCRouteMatchIndex]RoutePriorityType{
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
					Source: &gatewayapi.GRPCRoute{
						ObjectMeta: metav1.ObjectMeta{
							Namespace:         "default",
							Name:              "grpcroute-1",
							CreationTimestamp: metav1.NewTime(now.Add(-5 * time.Second)),
						},
						Spec: gatewayapi.GRPCRouteSpec{
							Hostnames: []gatewayapi.Hostname{"foo.com"},
							Rules: []gatewayapi.GRPCRouteRule{
								{
									Matches: []gatewayapi.GRPCRouteMatch{
										{
											Method: &gatewayapi.GRPCMethodMatch{
												Service: lo.ToPtr("cats"),
												Method:  lo.ToPtr("list"),
											},
										},
										{
											Method: &gatewayapi.GRPCMethodMatch{
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
					Match: gatewayapi.GRPCRouteMatch{
						Method: &gatewayapi.GRPCMethodMatch{
							Service: lo.ToPtr("cats"),
							Method:  lo.ToPtr("list"),
						},
					},
					RuleIndex:  0,
					MatchIndex: 0,
				},
				{
					// The same as the one above
					Source: &gatewayapi.GRPCRoute{
						ObjectMeta: metav1.ObjectMeta{
							Namespace:         "default",
							Name:              "grpcroute-1",
							CreationTimestamp: metav1.NewTime(now.Add(-5 * time.Second)),
						},
						Spec: gatewayapi.GRPCRouteSpec{
							Hostnames: []gatewayapi.Hostname{"foo.com"},
							Rules: []gatewayapi.GRPCRouteRule{
								{
									Matches: []gatewayapi.GRPCRouteMatch{
										{
											Method: &gatewayapi.GRPCMethodMatch{
												Service: lo.ToPtr("cats"),
												Method:  lo.ToPtr("list"),
											},
										},
										{
											Method: &gatewayapi.GRPCMethodMatch{
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
					Match: gatewayapi.GRPCRouteMatch{
						Method: &gatewayapi.GRPCMethodMatch{
							Service: lo.ToPtr("dogs"),
							Method:  lo.ToPtr("list"),
						},
					},
					RuleIndex:  0,
					MatchIndex: 1,
				},
			},
			priorities: map[splitGRPCRouteMatchIndex]RoutePriorityType{
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
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
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
					Source: &gatewayapi.GRPCRoute{
						TypeMeta: grpcRouteTypeMeta,
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "no-hostname-exact-method",
						},
						Spec: gatewayapi.GRPCRouteSpec{
							Rules: []gatewayapi.GRPCRouteRule{
								{
									Matches: []gatewayapi.GRPCRouteMatch{
										{
											Method: &gatewayapi.GRPCMethodMatch{
												Type:    lo.ToPtr(gatewayapi.GRPCMethodMatchExact),
												Service: lo.ToPtr("pets"),
												Method:  lo.ToPtr("list"),
											},
										},
									},
								},
							},
						},
					},
					Match: gatewayapi.GRPCRouteMatch{
						Method: &gatewayapi.GRPCMethodMatch{
							Type:    lo.ToPtr(gatewayapi.GRPCMethodMatchExact),
							Service: lo.ToPtr("pets"),
							Method:  lo.ToPtr("list"),
						},
					},
					RuleIndex:  0,
					MatchIndex: 0,
				},
				Priority: 1 << 14,
			},
			expectedRoute: kongstate.Route{
				Route: kong.Route{
					Name:         kong.String("grpcroute.default.no-hostname-exact-method._.0.0"),
					PreserveHost: kong.Bool(true),
					Expression:   kong.String(`http.path == "/pets/list"`),
					Priority:     kong.Uint64(1 << 14),
				},
				ExpressionRoutes: true,
			},
		},
		{
			name: "precise hostname and regex method match",
			splitGRPCMatchWithPriority: SplitGRPCRouteMatchToPriority{
				Match: SplitGRPCRouteMatch{
					Source: &gatewayapi.GRPCRoute{
						TypeMeta: grpcRouteTypeMeta,
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "precise-hostname-regex-method",
						},
						Spec: gatewayapi.GRPCRouteSpec{
							Hostnames: []gatewayapi.Hostname{
								"foo.com",
							},
							Rules: []gatewayapi.GRPCRouteRule{
								{
									Matches: []gatewayapi.GRPCRouteMatch{
										{
											Method: &gatewayapi.GRPCMethodMatch{
												Type:    lo.ToPtr(gatewayapi.GRPCMethodMatchRegularExpression),
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
					Match: gatewayapi.GRPCRouteMatch{
						Method: &gatewayapi.GRPCMethodMatch{
							Type:    lo.ToPtr(gatewayapi.GRPCMethodMatchRegularExpression),
							Service: lo.ToPtr("name"),
							Method:  lo.ToPtr("[a-z0-9]+"),
						},
					},
					RuleIndex:  0,
					MatchIndex: 0,
				},
				Priority: (1 << 31) + 1,
			},
			expectedRoute: kongstate.Route{
				Route: kong.Route{
					Name:         kong.String("grpcroute.default.precise-hostname-regex-method.foo.com.0.0"),
					Expression:   kong.String(`(http.path ~ "^/name/[a-z0-9]+") && (http.host == "foo.com")`),
					PreserveHost: kong.Bool(true),
					Priority:     kong.Uint64((1 << 31) + 1),
				},
				ExpressionRoutes: true,
			},
		},
		{
			name: "wildcard hostname and header match",
			splitGRPCMatchWithPriority: SplitGRPCRouteMatchToPriority{
				Match: SplitGRPCRouteMatch{
					Source: &gatewayapi.GRPCRoute{
						TypeMeta: grpcRouteTypeMeta,
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "wildcard-hostname-header-match",
						},
						Spec: gatewayapi.GRPCRouteSpec{
							Hostnames: []gatewayapi.Hostname{
								"*.foo.com",
							},
							Rules: []gatewayapi.GRPCRouteRule{
								{
									Matches: []gatewayapi.GRPCRouteMatch{
										{
											Method: &gatewayapi.GRPCMethodMatch{
												Type:    lo.ToPtr(gatewayapi.GRPCMethodMatchExact),
												Service: lo.ToPtr("names"),
												Method:  lo.ToPtr("create"),
											},
										},
										{
											Method: &gatewayapi.GRPCMethodMatch{
												Type:    lo.ToPtr(gatewayapi.GRPCMethodMatchExact),
												Service: lo.ToPtr("name"),
											},
											Headers: []gatewayapi.GRPCHeaderMatch{
												{
													Name:  gatewayapi.GRPCHeaderName("foo"),
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
					Match: gatewayapi.GRPCRouteMatch{
						Method: &gatewayapi.GRPCMethodMatch{
							Type:    lo.ToPtr(gatewayapi.GRPCMethodMatchExact),
							Service: lo.ToPtr("name"),
						},
						Headers: []gatewayapi.GRPCHeaderMatch{
							{
								Name:  gatewayapi.GRPCHeaderName("foo"),
								Value: "bar",
							},
						},
					},
					RuleIndex:  0,
					MatchIndex: 1,
				},
				Priority: (1 << 42) + 1,
			},
			expectedRoute: kongstate.Route{
				Route: kong.Route{
					Name:         kong.String("grpcroute.default.wildcard-hostname-header-match._.foo.com.0.1"),
					Expression:   kong.String(`(http.path ^= "/name/") && (http.headers.foo == "bar") && (http.host =^ ".foo.com")`),
					PreserveHost: kong.Bool(true),
					Priority:     kong.Uint64((1 << 42) + 1),
				},
				ExpressionRoutes: true,
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
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
