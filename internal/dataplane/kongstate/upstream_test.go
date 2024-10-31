package kongstate

import (
	"testing"

	"github.com/kong/go-kong/kong"
	kongv1 "github.com/kong/kubernetes-configuration/api/configuration/v1"
	kongv1beta1 "github.com/kong/kubernetes-configuration/api/configuration/v1beta1"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestOverrideUpstream(t *testing.T) {
	testTable := []struct {
		inUpstream     Upstream
		inKongIngresss *kongv1.KongIngress
		outUpstream    Upstream
		svc            *corev1.Service
	}{
		{
			inUpstream: Upstream{
				Upstream: kong.Upstream{
					Name: kong.String("foo.com"),
				},
			},
			inKongIngresss: nil,
			outUpstream: Upstream{
				Upstream: kong.Upstream{
					Name: kong.String("foo.com"),
				},
			},
		},
		{
			inUpstream: Upstream{
				Upstream: kong.Upstream{
					Name: kong.String("foo.com"),
				},
			},
			inKongIngresss: &kongv1.KongIngress{
				Upstream: &kongv1.KongIngressUpstream{
					HashOn:                 kong.String("HashOn"),
					HashOnCookie:           kong.String("HashOnCookie"),
					HashOnCookiePath:       kong.String("HashOnCookiePath"),
					HashOnHeader:           kong.String("HashOnHeader"),
					HashFallback:           kong.String("HashFallback"),
					HashFallbackHeader:     kong.String("HashFallbackHeader"),
					HashOnQueryArg:         kong.String("HashOnQueryArg"),
					HashFallbackQueryArg:   kong.String("HashFallbackQueryArg"),
					HashOnURICapture:       kong.String("HashOnURICapture"),
					HashFallbackURICapture: kong.String("HashFallbackURICapture"),
					Slots:                  kong.Int(42),
				},
			},
			outUpstream: Upstream{
				Upstream: kong.Upstream{
					Name:                   kong.String("foo.com"),
					HashOn:                 kong.String("HashOn"),
					HashOnCookie:           kong.String("HashOnCookie"),
					HashOnCookiePath:       kong.String("HashOnCookiePath"),
					HashOnHeader:           kong.String("HashOnHeader"),
					HashFallback:           kong.String("HashFallback"),
					HashFallbackHeader:     kong.String("HashFallbackHeader"),
					HostHeader:             kong.String("foo.com"),
					HashOnQueryArg:         kong.String("HashOnQueryArg"),
					HashFallbackQueryArg:   kong.String("HashFallbackQueryArg"),
					HashOnURICapture:       kong.String("HashOnURICapture"),
					HashFallbackURICapture: kong.String("HashFallbackURICapture"),
					Slots:                  kong.Int(42),
				},
			},
			svc: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						"konghq.com/host-header": "foo.com",
					},
				},
			},
		},
	}

	for _, testcase := range testTable {
		testcase.inUpstream.override(testcase.inKongIngresss, testcase.svc)
		require.Equal(t, testcase.inUpstream, testcase.outUpstream)
	}

	require.NotPanics(t, func() {
		var nilUpstream *Upstream
		nilUpstream.override(nil, nil)
	})
}

func TestUpstreamOverrideByKongUpstreamPolicy(t *testing.T) {
	testCases := []struct {
		name               string
		upstream           kong.Upstream
		kongUpstreamPolicy *kongv1beta1.KongUpstreamPolicy
		expected           kong.Upstream
	}{
		{
			name: "algorithm, slots, healthchecks, hash_on, hash_fallback, hash_fallback_header",
			upstream: kong.Upstream{
				Algorithm: kong.String("Algorithm"),
				Slots:     kong.Int(42),
				Healthchecks: &kong.Healthcheck{
					Active: &kong.ActiveHealthcheck{
						Concurrency: kong.Int(1),
					},
				},
				HashOn:             kong.String("HashOn"),
				HashFallback:       kong.String("HashFallback"),
				HashFallbackHeader: kong.String("HashFallbackHeader"),
			},
			kongUpstreamPolicy: &kongv1beta1.KongUpstreamPolicy{
				Spec: kongv1beta1.KongUpstreamPolicySpec{
					Algorithm: lo.ToPtr("least-connections"),
					Slots:     lo.ToPtr(10),
					Healthchecks: &kongv1beta1.KongUpstreamHealthcheck{
						Active: &kongv1beta1.KongUpstreamActiveHealthcheck{
							Concurrency: lo.ToPtr(2),
						},
					},
					HashOn: &kongv1beta1.KongUpstreamHash{
						Input: lo.ToPtr(kongv1beta1.HashInput("consumer")),
					},
					HashOnFallback: &kongv1beta1.KongUpstreamHash{
						Header: lo.ToPtr("foo"),
					},
				},
			},
			expected: kong.Upstream{
				Algorithm: kong.String("least-connections"),
				Slots:     kong.Int(10),
				Healthchecks: &kong.Healthcheck{
					Active: &kong.ActiveHealthcheck{
						Concurrency: kong.Int(2),
					},
				},
				HashOn:             kong.String("consumer"),
				HashFallback:       kong.String("header"),
				HashFallbackHeader: kong.String("foo"),
			},
		},
		{
			name: "hash_on_header, hash_fallback_query_arg",
			upstream: kong.Upstream{
				HashOn:               kong.String("HashOn"),
				HashFallback:         kong.String("HashFallback"),
				HashOnHeader:         kong.String("HashOnHeader"),
				HashFallbackQueryArg: kong.String("HashOnQueryArg"),
			},
			kongUpstreamPolicy: &kongv1beta1.KongUpstreamPolicy{
				Spec: kongv1beta1.KongUpstreamPolicySpec{
					HashOn: &kongv1beta1.KongUpstreamHash{
						Header: lo.ToPtr("foo"),
					},
					HashOnFallback: &kongv1beta1.KongUpstreamHash{
						QueryArg: lo.ToPtr("foo"),
					},
				},
			},
			expected: kong.Upstream{
				HashOn:               kong.String("header"),
				HashFallback:         kong.String("query_arg"),
				HashOnHeader:         kong.String("foo"),
				HashFallbackQueryArg: kong.String("foo"),
			},
		},
		{
			name: "hash_on_cookie, hash_on_cookie_path, hash_fallback_uri_capture",
			upstream: kong.Upstream{
				HashOn:                 kong.String("HashOn"),
				HashFallback:           kong.String("HashFallback"),
				HashOnCookie:           kong.String("HashOnCookie"),
				HashOnCookiePath:       kong.String("HashOnCookiePath"),
				HashFallbackURICapture: kong.String("HashFallbackURICapture"),
			},
			kongUpstreamPolicy: &kongv1beta1.KongUpstreamPolicy{
				Spec: kongv1beta1.KongUpstreamPolicySpec{
					HashOn: &kongv1beta1.KongUpstreamHash{
						Cookie:     lo.ToPtr("foo"),
						CookiePath: lo.ToPtr("/"),
					},
					HashOnFallback: &kongv1beta1.KongUpstreamHash{
						URICapture: lo.ToPtr("foo"),
					},
				},
			},
			expected: kong.Upstream{
				HashOn:                 kong.String("cookie"),
				HashFallback:           kong.String("uri_capture"),
				HashOnCookie:           kong.String("foo"),
				HashOnCookiePath:       kong.String("/"),
				HashFallbackURICapture: kong.String("foo"),
			},
		},
		{
			name: "hash_on_uri_capture",
			upstream: kong.Upstream{
				HashOn:           kong.String("HashOn"),
				HashOnURICapture: kong.String("HashOnURICapture"),
			},
			kongUpstreamPolicy: &kongv1beta1.KongUpstreamPolicy{
				Spec: kongv1beta1.KongUpstreamPolicySpec{
					HashOn: &kongv1beta1.KongUpstreamHash{
						URICapture: lo.ToPtr("foo"),
					},
				},
			},
			expected: kong.Upstream{
				HashOn:           kong.String("uri_capture"),
				HashOnURICapture: kong.String("foo"),
			},
		},
		{
			name: "hash_on_query_arg",
			upstream: kong.Upstream{
				HashOn:         kong.String("HashOn"),
				HashOnQueryArg: kong.String("HashOnQueryArg"),
			},
			kongUpstreamPolicy: &kongv1beta1.KongUpstreamPolicy{
				Spec: kongv1beta1.KongUpstreamPolicySpec{
					HashOn: &kongv1beta1.KongUpstreamHash{
						QueryArg: lo.ToPtr("foo"),
					},
				},
			},
			expected: kong.Upstream{
				HashOn:         kong.String("query_arg"),
				HashOnQueryArg: kong.String("foo"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			upstream := Upstream{Upstream: tc.upstream}
			upstream.overrideByKongUpstreamPolicy(tc.kongUpstreamPolicy)
			require.Equal(t, tc.expected, upstream.Upstream)
		})
	}

	require.NotPanics(t, func() {
		var nilUpstream *Upstream
		nilUpstream.overrideByKongUpstreamPolicy(nil)
	})
}
