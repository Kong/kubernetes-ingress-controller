package kongstate

import (
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"

	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1beta1"
)

const (
	KongHashOnTypeHeader     string = "header"
	KongHashOnTypeCookie     string = "cookie"
	KongHashOnTypeQueryArg   string = "query_arg"
	KongHashOnTypeURICapture string = "uri_capture"
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
	// CRD validations will ensure only one of hashOn fields can be set, therefore the order doesn't matter.
	switch {
	case hashOn.Input != nil:
		return lo.ToPtr(string(*hashOn.Input))
	case hashOn.Header != nil:
		return lo.ToPtr(KongHashOnTypeHeader)
	case hashOn.Cookie != nil:
		return lo.ToPtr(KongHashOnTypeCookie)
	case hashOn.QueryArg != nil:
		return lo.ToPtr(KongHashOnTypeQueryArg)
	case hashOn.URICapture != nil:
		return lo.ToPtr(KongHashOnTypeURICapture)
	default:
		return nil
	}
}

func translateHashOnHeader(hashOn *kongv1beta1.KongUpstreamHash) *string {
	if hashOn == nil {
		return nil
	}
	return hashOn.Header
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
		HTTPStatuses: translateHTTPStatuses(healthy.HTTPStatuses),
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
		HTTPStatuses: translateHTTPStatuses(unhealthy.HTTPStatuses),
		TCPFailures:  unhealthy.TCPFailures,
		Timeouts:     unhealthy.Timeouts,
		Interval:     unhealthy.Interval,
	}
}

func translateHTTPStatuses(statuses []kongv1beta1.HTTPStatus) []int {
	if statuses == nil {
		return nil
	}
	// Using lo.Map only in case healthy.HTTPStatuses is not nil, because lo.Map creates a non-nil slice even
	// if the input slice is nil.
	return lo.Map(statuses, func(s kongv1beta1.HTTPStatus, _ int) int { return int(s) })
}
