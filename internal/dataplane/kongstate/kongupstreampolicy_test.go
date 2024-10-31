package kongstate_test

import (
	"testing"

	"github.com/kong/go-kong/kong"
	kongv1beta1 "github.com/kong/kubernetes-configuration/api/configuration/v1beta1"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
)

func TestGetKongUpstreamPolicyForServices(t *testing.T) {
	testCases := []struct {
		name          string
		servicesGroup []*corev1.Service
		policies      []*kongv1beta1.KongUpstreamPolicy
		expectPolicy  bool
		expectError   string
	}{
		{
			name:         "no services in group gives no policy",
			expectPolicy: false,
		},
		{
			name: "no KongUpstreamPolicy in store while services are configured with one gives error",
			servicesGroup: []*corev1.Service{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "svc",
					Namespace: "default",
					Annotations: map[string]string{
						kongv1beta1.KongUpstreamPolicyAnnotationKey: "upstream-policy",
					},
				},
			}},
			expectError: "failed fetching KongUpstreamPolicy",
		},
		{
			name: "all services configured with no KongUpstreamPolicy gives no policy and no error",
			servicesGroup: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "svc-1",
						Namespace: "default",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "svc-2",
						Namespace: "default",
					},
				},
			},
			expectPolicy: false,
		},
		{
			name: "services in group with different KongUpstreamPolicy configurations gives error",
			servicesGroup: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "svc-1",
						Namespace: "default",
						Annotations: map[string]string{
							kongv1beta1.KongUpstreamPolicyAnnotationKey: "upstream-policy",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "svc-2",
						Namespace: "default",
						Annotations: map[string]string{
							kongv1beta1.KongUpstreamPolicyAnnotationKey: "other-upstream-policy",
						},
					},
				},
			},
			expectError: "inconsistent KongUpstreamPolicy configuration for services",
		},
		{
			name: "one service with and one without KongUpstreamPolicy configuration gives error",
			servicesGroup: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "svc-1",
						Namespace: "default",
						Annotations: map[string]string{
							kongv1beta1.KongUpstreamPolicyAnnotationKey: "upstream-policy",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "svc-2",
						Namespace: "default",
					},
				},
			},
			expectError: "inconsistent KongUpstreamPolicy configuration for services",
		},
		{
			name: "all the services configured with the same KongUpstreamPolicy gives the policy",
			servicesGroup: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "svc-1",
						Namespace: "default",
						Annotations: map[string]string{
							kongv1beta1.KongUpstreamPolicyAnnotationKey: "upstream-policy",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "svc-2",
						Namespace: "default",
						Annotations: map[string]string{
							kongv1beta1.KongUpstreamPolicyAnnotationKey: "upstream-policy",
						},
					},
				},
			},
			policies: []*kongv1beta1.KongUpstreamPolicy{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "upstream-policy",
						Namespace: "default",
					},
				},
			},
			expectPolicy: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, err := store.NewFakeStore(store.FakeObjects{
				Services:             tc.servicesGroup,
				KongUpstreamPolicies: tc.policies,
			})
			require.NoError(t, err)

			policy, err := kongstate.GetKongUpstreamPolicyForServices(s, tc.servicesGroup)
			if tc.expectError != "" {
				require.ErrorContains(t, err, tc.expectError)
				return
			}
			if tc.expectPolicy {
				require.NotNil(t, policy)
			} else {
				require.Nil(t, policy)
			}
		})
	}
}

func TestTranslateKongUpstreamPolicy(t *testing.T) {
	testCases := []struct {
		name             string
		policySpec       kongv1beta1.KongUpstreamPolicySpec
		expectedUpstream *kong.Upstream
	}{
		{
			name: "KongUpstreamPolicySpec with no hash-on or hash-fallback",
			policySpec: kongv1beta1.KongUpstreamPolicySpec{
				Algorithm: lo.ToPtr("least-connections"),
				Slots:     lo.ToPtr(10),
			},
			expectedUpstream: &kong.Upstream{
				Algorithm: lo.ToPtr("least-connections"),
				Slots:     lo.ToPtr(10),
			},
		},
		{
			name: "KongUpstreamPolicySpec with hash-on header",
			policySpec: kongv1beta1.KongUpstreamPolicySpec{
				HashOn: &kongv1beta1.KongUpstreamHash{
					Header: lo.ToPtr("foo"),
				},
				HashOnFallback: &kongv1beta1.KongUpstreamHash{
					Header: lo.ToPtr("bar"),
				},
			},
			expectedUpstream: &kong.Upstream{
				HashOn:             lo.ToPtr("header"),
				HashOnHeader:       lo.ToPtr("foo"),
				HashFallback:       lo.ToPtr("header"),
				HashFallbackHeader: lo.ToPtr("bar"),
			},
		},
		{
			name: "KongUpstreamPolicySpec with hash-on cookie",
			policySpec: kongv1beta1.KongUpstreamPolicySpec{
				HashOn: &kongv1beta1.KongUpstreamHash{
					Cookie:     lo.ToPtr("foo"),
					CookiePath: lo.ToPtr("/"),
				},
			},
			expectedUpstream: &kong.Upstream{
				HashOn:           lo.ToPtr("cookie"),
				HashOnCookie:     lo.ToPtr("foo"),
				HashOnCookiePath: lo.ToPtr("/"),
			},
		},
		{
			name: "KongUpstreamPolicySpec with hash-on query-arg",
			policySpec: kongv1beta1.KongUpstreamPolicySpec{
				HashOn: &kongv1beta1.KongUpstreamHash{
					QueryArg: lo.ToPtr("foo"),
				},
			},
			expectedUpstream: &kong.Upstream{
				HashOn:         lo.ToPtr("query_arg"),
				HashOnQueryArg: lo.ToPtr("foo"),
			},
		},
		{
			name: "KongUpstreamPolicySpec with hash-on uri-capture",
			policySpec: kongv1beta1.KongUpstreamPolicySpec{
				HashOn: &kongv1beta1.KongUpstreamHash{
					URICapture: lo.ToPtr("foo"),
				},
			},
			expectedUpstream: &kong.Upstream{
				HashOn:           lo.ToPtr("uri_capture"),
				HashOnURICapture: lo.ToPtr("foo"),
			},
		},
		{
			name: "KongUpstreamPolicySpec with predefined hash input",
			policySpec: kongv1beta1.KongUpstreamPolicySpec{
				HashOn: &kongv1beta1.KongUpstreamHash{
					Input: lo.ToPtr(kongv1beta1.HashInput("consumer")),
				},
				HashOnFallback: &kongv1beta1.KongUpstreamHash{
					Input: lo.ToPtr(kongv1beta1.HashInput("ip")),
				},
			},
			expectedUpstream: &kong.Upstream{
				HashOn:       lo.ToPtr("consumer"),
				HashFallback: lo.ToPtr("ip"),
			},
		},
		{
			name: "KongUpstreamPolicySpec with healthchecks",
			policySpec: kongv1beta1.KongUpstreamPolicySpec{
				Healthchecks: &kongv1beta1.KongUpstreamHealthcheck{
					Active: &kongv1beta1.KongUpstreamActiveHealthcheck{
						Type:        lo.ToPtr("http"),
						Concurrency: lo.ToPtr(10),
						Healthy: &kongv1beta1.KongUpstreamHealthcheckHealthy{
							HTTPStatuses: []kongv1beta1.HTTPStatus{200},
							Interval:     lo.ToPtr(20),
							Successes:    lo.ToPtr(30),
						},
						Unhealthy: &kongv1beta1.KongUpstreamHealthcheckUnhealthy{
							HTTPFailures: lo.ToPtr(40),
							HTTPStatuses: []kongv1beta1.HTTPStatus{500},
							Timeouts:     lo.ToPtr(60),
							Interval:     lo.ToPtr(70),
						},
						HTTPPath:               lo.ToPtr("/foo"),
						HTTPSSNI:               lo.ToPtr("foo.com"),
						HTTPSVerifyCertificate: lo.ToPtr(true),
						Timeout:                lo.ToPtr(80),
						Headers:                map[string][]string{"foo": {"bar"}},
					},
					Passive: &kongv1beta1.KongUpstreamPassiveHealthcheck{
						Type: lo.ToPtr("tcp"),
						Healthy: &kongv1beta1.KongUpstreamHealthcheckHealthy{
							Successes: lo.ToPtr(100),
						},
						Unhealthy: &kongv1beta1.KongUpstreamHealthcheckUnhealthy{
							TCPFailures: lo.ToPtr(110),
							Timeouts:    lo.ToPtr(120),
						},
					},
					Threshold: lo.ToPtr(140),
				},
			},
			expectedUpstream: &kong.Upstream{
				Healthchecks: &kong.Healthcheck{
					Active: &kong.ActiveHealthcheck{
						Type:        lo.ToPtr("http"),
						Concurrency: lo.ToPtr(10),
						Healthy: &kong.Healthy{
							HTTPStatuses: []int{200},
							Interval:     lo.ToPtr(20),
							Successes:    lo.ToPtr(30),
						},
						Unhealthy: &kong.Unhealthy{
							HTTPFailures: lo.ToPtr(40),
							HTTPStatuses: []int{500},
							Timeouts:     lo.ToPtr(60),
							Interval:     lo.ToPtr(70),
						},
						HTTPPath:               lo.ToPtr("/foo"),
						HTTPSSni:               lo.ToPtr("foo.com"),
						HTTPSVerifyCertificate: lo.ToPtr(true),
						Headers:                map[string][]string{"foo": {"bar"}},
						Timeout:                lo.ToPtr(80),
					},
					Passive: &kong.PassiveHealthcheck{
						Type: lo.ToPtr("tcp"),
						Healthy: &kong.Healthy{
							Successes: lo.ToPtr(100),
						},
						Unhealthy: &kong.Unhealthy{
							TCPFailures: lo.ToPtr(110),
							Timeouts:    lo.ToPtr(120),
						},
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualUpstream := kongstate.TranslateKongUpstreamPolicy(tc.policySpec)
			require.Equal(t, tc.expectedUpstream, actualUpstream)
		})
	}
}
