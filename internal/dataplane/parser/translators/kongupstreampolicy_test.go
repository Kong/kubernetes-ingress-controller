package translators_test

import (
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser/translators"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
)

func TestTranslateKongUpstreamPolicy(t *testing.T) {
	testCases := []struct {
		name             string
		policySpec       kongv1beta1.KongUpstreamPolicySpec
		expectedUpstream *kong.Upstream
	}{
		{
			name: "KongUpstreamPolicySpec with no hash-on or hash-fallback",
			policySpec: kongv1beta1.KongUpstreamPolicySpec{
				HostHeader: lo.ToPtr("foo"),
				Algorithm:  lo.ToPtr("least-connections"),
				Slots:      lo.ToPtr(10),
			},
			expectedUpstream: &kong.Upstream{
				HostHeader: lo.ToPtr("foo"),
				Algorithm:  lo.ToPtr("least-connections"),
				Slots:      lo.ToPtr(10),
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
			actualUpstream := translators.TranslateKongUpstreamPolicy(tc.policySpec)
			require.Equal(t, tc.expectedUpstream, actualUpstream)
		})
	}
}
