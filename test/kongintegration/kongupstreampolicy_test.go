package kongintegration

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/v3/test/kongintegration/containers"
)

// TestKongUpstreamPolicyTranslation ensures that the Upstream Policy CRD is translated to the Kong Upstream
// object in a way that when it's sent to Kong, all the fields are correctly propagated.
func TestKongUpstreamPolicyTranslation(t *testing.T) {
	t.Parallel()

	const (
		timeout = time.Second * 1
		period  = time.Millisecond * 100
	)

	ctx := context.Background()

	kongC := containers.NewKong(ctx, t)
	kongClient, err := adminapi.NewKongAPIClient(kongC.AdminURL(ctx, t), &http.Client{})
	require.NoError(t, err)
	updateStrategy := sendconfig.NewUpdateStrategyInMemory(
		kongClient,
		sendconfig.DefaultContentToDBLessConfigConverter{},
		logr.Discard(),
	)

	testCases := []struct {
		name             string
		policySpec       kongv1beta1.KongUpstreamPolicySpec
		expectedUpstream *kong.Upstream
	}{
		{
			name: "KongUpstreamPolicySpec with no hash-on or hash-fallback",
			policySpec: kongv1beta1.KongUpstreamPolicySpec{
				Algorithm: lo.ToPtr("least-connections"),
				Slots:     lo.ToPtr(20),
			},
			expectedUpstream: &kong.Upstream{
				Algorithm: lo.ToPtr("least-connections"),
				Slots:     lo.ToPtr(20),
			},
		},
		{
			name: "KongUpstreamPolicySpec with hash-on header",
			policySpec: kongv1beta1.KongUpstreamPolicySpec{
				Algorithm: lo.ToPtr("consistent-hashing"),
				HashOn: &kongv1beta1.KongUpstreamHash{
					Header: lo.ToPtr("foo"),
				},
				HashOnFallback: &kongv1beta1.KongUpstreamHash{
					Header: lo.ToPtr("bar"),
				},
			},
			expectedUpstream: &kong.Upstream{
				Algorithm:          lo.ToPtr("consistent-hashing"),
				HashOn:             lo.ToPtr("header"),
				HashOnHeader:       lo.ToPtr("foo"),
				HashFallback:       lo.ToPtr("header"),
				HashFallbackHeader: lo.ToPtr("bar"),
			},
		},
		{
			name: "KongUpstreamPolicySpec with hash-on cookie",
			policySpec: kongv1beta1.KongUpstreamPolicySpec{
				Algorithm: lo.ToPtr("consistent-hashing"),
				HashOn: &kongv1beta1.KongUpstreamHash{
					Cookie:     lo.ToPtr("foo"),
					CookiePath: lo.ToPtr("/"),
				},
			},
			expectedUpstream: &kong.Upstream{
				Algorithm:        lo.ToPtr("consistent-hashing"),
				HashOn:           lo.ToPtr("cookie"),
				HashOnCookie:     lo.ToPtr("foo"),
				HashOnCookiePath: lo.ToPtr("/"),
			},
		},
		{
			name: "KongUpstreamPolicySpec with hash-on query-arg",
			policySpec: kongv1beta1.KongUpstreamPolicySpec{
				Algorithm: lo.ToPtr("consistent-hashing"),
				HashOn: &kongv1beta1.KongUpstreamHash{
					QueryArg: lo.ToPtr("foo"),
				},
			},
			expectedUpstream: &kong.Upstream{
				Algorithm:      lo.ToPtr("consistent-hashing"),
				HashOn:         lo.ToPtr("query_arg"),
				HashOnQueryArg: lo.ToPtr("foo"),
			},
		},
		{
			name: "KongUpstreamPolicySpec with predefined hash input",
			policySpec: kongv1beta1.KongUpstreamPolicySpec{
				Algorithm: lo.ToPtr("consistent-hashing"),
				HashOn: &kongv1beta1.KongUpstreamHash{
					Input: lo.ToPtr(kongv1beta1.HashInput("consumer")),
				},
				HashOnFallback: &kongv1beta1.KongUpstreamHash{
					Input: lo.ToPtr(kongv1beta1.HashInput("ip")),
				},
			},
			expectedUpstream: &kong.Upstream{
				Algorithm:    lo.ToPtr("consistent-hashing"),
				HashOn:       lo.ToPtr("consumer"),
				HashFallback: lo.ToPtr("ip"),
			},
		},
		{
			name: "KongUpstreamPolicySpec with hash-on uri-capture",
			policySpec: kongv1beta1.KongUpstreamPolicySpec{
				Algorithm: lo.ToPtr("consistent-hashing"),
				HashOn: &kongv1beta1.KongUpstreamHash{
					URICapture: lo.ToPtr("foo"),
				},
			},
			expectedUpstream: &kong.Upstream{
				Algorithm:        lo.ToPtr("consistent-hashing"),
				HashOn:           lo.ToPtr("uri_capture"),
				HashOnURICapture: lo.ToPtr("foo"),
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
							TCPFailures:  lo.ToPtr(5),
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
							HTTPStatuses: []kongv1beta1.HTTPStatus{200},
							Successes:    lo.ToPtr(100),
						},
						Unhealthy: &kongv1beta1.KongUpstreamHealthcheckUnhealthy{
							HTTPStatuses: []kongv1beta1.HTTPStatus{500},
							TCPFailures:  lo.ToPtr(110),
							Timeouts:     lo.ToPtr(120),
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
							TCPFailures:  lo.ToPtr(5),
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
							HTTPStatuses: []int{200},
							Successes:    lo.ToPtr(100),
						},
						Unhealthy: &kong.Unhealthy{
							HTTPFailures: lo.ToPtr(0),
							HTTPStatuses: []int{500},
							TCPFailures:  lo.ToPtr(110),
							Timeouts:     lo.ToPtr(120),
						},
					},
					Threshold: lo.ToPtr(0.),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			translatedUpstream := kongstate.TranslateKongUpstreamPolicy(tc.policySpec)
			const upstreamName = "test-upstream"
			translatedUpstream.Name = lo.ToPtr(upstreamName)
			tc.expectedUpstream.Name = lo.ToPtr(upstreamName)

			content := sendconfig.ContentWithHash{
				Content: &file.Content{
					FormatVersion: "3.0",
					Upstreams: []file.FUpstream{
						{
							Upstream: *translatedUpstream,
						},
					},
				},
			}

			// Update Kong with the Upstream.
			require.Eventually(t, func() bool {
				err = updateStrategy.Update(ctx, content)
				if err != nil {
					t.Logf("error updating Kong configuration: %v", err)
					return false
				}
				return true
			}, timeout, period)

			// Wait for the Upstream to be created in Kong and assert it matches the expected Upstream.
			require.Eventually(t, func() bool {
				upstreamInKong, err := kongClient.Upstreams.Get(ctx, lo.ToPtr(upstreamName))
				if err != nil {
					t.Logf("error getting upstream from Kong: %v", err)
					return false
				}

				// We ignore the values that are generated by Kong and are not known in advance.
				ignoreKongGeneratedValues := cmp.Transformer(
					"IgnoreKongGeneratedValues",
					func(upstream *kong.Upstream) *kong.Upstream {
						return dropIDsAndTimestamps(dropKongDefaults(upstream))
					},
				)
				diff := cmp.Diff(upstreamInKong, tc.expectedUpstream, ignoreKongGeneratedValues)
				if diff != "" {
					t.Logf("upstream in Kong does not match expected upstream, diff: %s", diff)
					return false
				}
				return true
			}, timeout, period)
		})
	}
}

// dropIDsAndTimestamps drops the ID and CreatedAt fields from the Upstream. These fields are generated by Kong and
// are not known in advance, so we want to ignore them when comparing the Upstream in Kong with the expected Upstream.
func dropIDsAndTimestamps(upstream *kong.Upstream) *kong.Upstream {
	upstream = upstream.DeepCopy()
	upstream.ID = nil
	upstream.CreatedAt = nil
	return upstream
}

// dropKongDefaults drops the default values that Kong sets for some fields. We offload the responsibility of setting
// these default values to Kong, so we want to ignore them when comparing the Upstream in Kong with the expected one.
func dropKongDefaults(upstream *kong.Upstream) *kong.Upstream {
	upstream = upstream.DeepCopy()

	defaultHealthcheck := &kong.Healthcheck{
		Active: &kong.ActiveHealthcheck{
			Concurrency: lo.ToPtr(10),
			Healthy: &kong.Healthy{
				HTTPStatuses: []int{200, 302},
				Successes:    lo.ToPtr(0),
				Interval:     lo.ToPtr(0),
			},
			HTTPPath:               lo.ToPtr("/"),
			HTTPSVerifyCertificate: lo.ToPtr(true),
			Type:                   lo.ToPtr("http"),
			Timeout:                lo.ToPtr(1),
			Unhealthy: &kong.Unhealthy{
				HTTPFailures: lo.ToPtr(0),
				HTTPStatuses: []int{429, 404, 500, 501, 502, 503, 504, 505},
				TCPFailures:  lo.ToPtr(0),
				Timeouts:     lo.ToPtr(0),
				Interval:     lo.ToPtr(0),
			},
		},
		Passive: &kong.PassiveHealthcheck{
			Healthy: &kong.Healthy{
				HTTPStatuses: []int{200, 201, 202, 203, 204, 205, 206, 207, 208, 226, 300, 301, 302, 303, 304, 305, 306, 307, 308},
				Successes:    lo.ToPtr(0),
			},
			Type: lo.ToPtr("http"),
			Unhealthy: &kong.Unhealthy{
				HTTPFailures: lo.ToPtr(0),
				HTTPStatuses: []int{429, 500, 503},
				TCPFailures:  lo.ToPtr(0),
				Timeouts:     lo.ToPtr(0),
			},
		},
		Threshold: lo.ToPtr(0.),
	}

	if diff := cmp.Diff(upstream.Healthchecks, defaultHealthcheck); diff == "" {
		upstream.Healthchecks = nil
	}
	if upstream.HashOn != nil && *upstream.HashOn == "none" {
		upstream.HashOn = nil
	}
	if upstream.HashFallback != nil && *upstream.HashFallback == "none" {
		upstream.HashFallback = nil
	}
	if upstream.HashOnCookiePath != nil && *upstream.HashOnCookiePath == "/" {
		upstream.HashOnCookiePath = nil
	}
	if upstream.UseSrvName != nil && *upstream.UseSrvName == false {
		upstream.UseSrvName = nil
	}
	if upstream.Slots != nil && *upstream.Slots == 10000 {
		upstream.Slots = nil
	}
	if upstream.Algorithm != nil && *upstream.Algorithm == "round-robin" {
		upstream.Algorithm = nil
	}

	return upstream
}
