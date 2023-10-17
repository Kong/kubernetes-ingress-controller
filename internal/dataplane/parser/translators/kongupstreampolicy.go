package translators

import (
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"

	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
)

// TranslateKongUpstreamPolicy translates KongUpstreamPolicySpec to kong.Upstream. It makes assumption that
// KongUpstreamPolicySpec has been validated on the API level.
func TranslateKongUpstreamPolicy(policy kongv1beta1.KongUpstreamPolicySpec) *kong.Upstream {
	return &kong.Upstream{
		Algorithm:    policy.Algorithm,
		Slots:        policy.Slots,
		Healthchecks: translateHealthchecks(policy.Healthchecks),
		HostHeader:   policy.HostHeader,

		HashOn:           translateHashOn(policy.HashOn),
		HashOnHeader:     translateHashOnHeader(policy.HashOn),
		HashOnURICapture: translateHashOnURICapture(policy.HashOn),
		HashOnCookie:     translateHashOnCookie(policy.HashOn),
		HashOnCookiePath: translateHashOnCookiePath(policy.HashOn),
		HashOnQueryArg:   translateHashOnQueryArg(policy.HashOn),

		HashFallback:           translateHashOn(policy.HashOnFallback),
		HashFallbackHeader:     translateHashOnHeader(policy.HashOnFallback),
		HashFallbackURICapture: translateHashOnURICapture(policy.HashOnFallback),
		HashFallbackQueryArg:   translateHashOnQueryArg(policy.HashOnFallback),
	}
}

func translateHashOn(hashOn *kongv1beta1.KongUpstreamHash) *string {
	if hashOn == nil {
		return nil
	}
	// Only one of hashOn fields can be set.
	switch {
	case hashOn.Header != nil:
		return lo.ToPtr("header")
	case hashOn.Cookie != nil:
		return lo.ToPtr("cookie")
	case hashOn.QueryArg != nil:
		return lo.ToPtr("query_arg")
	case hashOn.URICapture != nil:
		return lo.ToPtr("uri_capture")
	default:
		return nil
	}
}

func translateHashOnHeader(hasOn *kongv1beta1.KongUpstreamHash) *string {
	if hasOn == nil {
		return nil
	}
	return hasOn.Header
}

func translateHashOnCookie(hashOn *kongv1beta1.KongUpstreamHash) *string {
	if hashOn == nil {
		return nil
	}
	return hashOn.Cookie
}

func translateHashOnQueryArg(hashOn *kongv1beta1.KongUpstreamHash) *string {
	if hashOn == nil {
		return nil
	}
	return hashOn.QueryArg
}

func translateHashOnURICapture(hashOn *kongv1beta1.KongUpstreamHash) *string {
	if hashOn == nil {
		return nil
	}
	return hashOn.URICapture
}

func translateHashOnCookiePath(hashOn *kongv1beta1.KongUpstreamHash) *string {
	if hashOn == nil {
		return nil
	}
	return hashOn.CookiePath
}

func translateHealthchecks(healthchecks *kongv1beta1.KongUpstreamHealthcheck) *kong.Healthcheck {
	if healthchecks == nil {
		return nil
	}
	return &kong.Healthcheck{
		Active:  translateActiveHealthcheck(healthchecks.Active),
		Passive: translatePassiveHealthcheck(healthchecks.Passive),
	}
}

func translateActiveHealthcheck(healthcheck *kongv1beta1.KongUpstreamActiveHealthcheck) *kong.ActiveHealthcheck {
	if healthcheck == nil {
		return nil
	}
	return &kong.ActiveHealthcheck{
		Concurrency:            healthcheck.Concurrency,
		HTTPPath:               healthcheck.HTTPPath,
		HTTPSSni:               healthcheck.HTTPSSNI,
		HTTPSVerifyCertificate: healthcheck.HTTPSVerifyCertificate,
		Type:                   healthcheck.Type,
		Timeout:                healthcheck.Timeout,
		Headers:                healthcheck.Headers,
		Healthy:                translateHealthy(healthcheck.Healthy),
		Unhealthy:              translateUnhealthy(healthcheck.Unhealthy),
	}
}

func translatePassiveHealthcheck(healthcheck *kongv1beta1.KongUpstreamPassiveHealthcheck) *kong.PassiveHealthcheck {
	if healthcheck == nil {
		return nil
	}
	return &kong.PassiveHealthcheck{
		Type:      healthcheck.Type,
		Healthy:   translateHealthy(healthcheck.Healthy),
		Unhealthy: translateUnhealthy(healthcheck.Unhealthy),
	}
}

func translateHealthy(healthy *kongv1beta1.KongUpstreamHealthcheckHealthy) *kong.Healthy {
	if healthy == nil {
		return nil
	}
	return &kong.Healthy{
		HTTPStatuses: healthy.HTTPStatuses,
		Interval:     healthy.Interval,
		Successes:    healthy.Successes,
	}
}

func translateUnhealthy(unhealthy *kongv1beta1.KongUpstreamHealthcheckUnhealthy) *kong.Unhealthy {
	if unhealthy == nil {
		return nil
	}
	return &kong.Unhealthy{
		HTTPFailures: unhealthy.HTTPFailures,
		HTTPStatuses: unhealthy.HTTPStatuses,
		TCPFailures:  unhealthy.TCPFailures,
		Timeouts:     unhealthy.Timeouts,
		Interval:     unhealthy.Interval,
	}
}
