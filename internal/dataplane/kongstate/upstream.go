package kongstate

import (
	"github.com/kong/go-kong/kong"
	corev1 "k8s.io/api/core/v1"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1beta1"
)

// Upstream is a wrapper around Upstream object in Kong.
type Upstream struct {
	kong.Upstream
	Targets []Target
	// Service this upstream is associated with.
	Service Service
}

func (u *Upstream) overrideHostHeader(anns map[string]string) {
	if u == nil {
		return
	}
	host := annotations.ExtractHostHeader(anns)
	if host == "" {
		return
	}
	u.HostHeader = kong.String(host)
}

// overrideByAnnotation modifies the Kong upstream based on annotations
// on the Kubernetes service.
func (u *Upstream) overrideByAnnotation(anns map[string]string) {
	if u == nil {
		return
	}
	u.overrideHostHeader(anns)
}

// overrideByKongIngress modifies the Kong upstream based on KongIngresses
// associated with the Kubernetes service.
func (u *Upstream) overrideByKongIngress(kongIngress *kongv1.KongIngress) {
	if u == nil {
		return
	}

	if kongIngress == nil || kongIngress.Upstream == nil {
		return
	}
	k := kongIngress.Upstream
	if k.HostHeader != nil {
		u.HostHeader = kong.String(*k.HostHeader)
	}
	if k.Algorithm != nil {
		u.Algorithm = kong.String(*k.Algorithm)
	}
	if k.Slots != nil {
		u.Slots = kong.Int(*k.Slots)
	}
	if k.Healthchecks != nil {
		u.Healthchecks = k.Healthchecks
	}
	if k.HashOn != nil {
		u.HashOn = kong.String(*k.HashOn)
	}
	if k.HashFallback != nil {
		u.HashFallback = kong.String(*k.HashFallback)
	}
	if k.HashOnHeader != nil {
		u.HashOnHeader = kong.String(*k.HashOnHeader)
	}
	if k.HashFallbackHeader != nil {
		u.HashFallbackHeader = kong.String(*k.HashFallbackHeader)
	}
	if k.HashOnCookie != nil {
		u.HashOnCookie = kong.String(*k.HashOnCookie)
	}
	if k.HashOnCookiePath != nil {
		u.HashOnCookiePath = kong.String(*k.HashOnCookiePath)
	}
	if k.HashOnQueryArg != nil {
		u.HashOnQueryArg = kong.String(*k.HashOnQueryArg)
	}
	if k.HashFallbackQueryArg != nil {
		u.HashFallbackQueryArg = kong.String(*k.HashFallbackQueryArg)
	}
	if k.HashOnURICapture != nil {
		u.HashOnURICapture = kong.String(*k.HashOnURICapture)
	}
	if k.HashFallbackURICapture != nil {
		u.HashFallbackURICapture = kong.String(*k.HashFallbackURICapture)
	}
}

func (u *Upstream) overrideByKongUpstreamPolicy(policy *kongv1beta1.KongUpstreamPolicy) {
	if u == nil {
		return
	}

	kongUpstreamOverrides := TranslateKongUpstreamPolicy(policy.Spec)
	if kongUpstreamOverrides.Algorithm != nil {
		u.Algorithm = kongUpstreamOverrides.Algorithm
	}
	if kongUpstreamOverrides.Slots != nil {
		u.Slots = kongUpstreamOverrides.Slots
	}
	if kongUpstreamOverrides.Healthchecks != nil {
		u.Healthchecks = kongUpstreamOverrides.Healthchecks
	}
	if kongUpstreamOverrides.HashOn != nil {
		u.HashOn = kongUpstreamOverrides.HashOn
	}
	if kongUpstreamOverrides.HashFallback != nil {
		u.HashFallback = kongUpstreamOverrides.HashFallback
	}
	if kongUpstreamOverrides.HashOnHeader != nil {
		u.HashOnHeader = kongUpstreamOverrides.HashOnHeader
	}
	if kongUpstreamOverrides.HashFallbackHeader != nil {
		u.HashFallbackHeader = kongUpstreamOverrides.HashFallbackHeader
	}
	if kongUpstreamOverrides.HashOnCookie != nil {
		u.HashOnCookie = kongUpstreamOverrides.HashOnCookie
	}
	if kongUpstreamOverrides.HashOnCookiePath != nil {
		u.HashOnCookiePath = kongUpstreamOverrides.HashOnCookiePath
	}
	if kongUpstreamOverrides.HashOnQueryArg != nil {
		u.HashOnQueryArg = kongUpstreamOverrides.HashOnQueryArg
	}
	if kongUpstreamOverrides.HashFallbackQueryArg != nil {
		u.HashFallbackQueryArg = kongUpstreamOverrides.HashFallbackQueryArg
	}
	if kongUpstreamOverrides.HashOnURICapture != nil {
		u.HashOnURICapture = kongUpstreamOverrides.HashOnURICapture
	}
	if kongUpstreamOverrides.HashFallbackURICapture != nil {
		u.HashFallbackURICapture = kongUpstreamOverrides.HashFallbackURICapture
	}
}

// override sets Upstream fields by KongIngress first, then by k8s Service's annotations.
func (u *Upstream) override(
	kongIngress *kongv1.KongIngress,
	svc *corev1.Service,
) {
	if u == nil {
		return
	}

	if u.Service.Parent != nil && kongIngress != nil {
		// If the parent object behind Kong Upstream's is a Gateway API object
		// (probably *Route) then check if we're trying to override said Service
		// configuration with a KongIngress object and if that's the case then
		// skip it since those should not be affected.
		gvk := u.Service.Parent.GetObjectKind().GroupVersionKind()
		if gvk.Group == gatewayv1.GroupName {
			// No log needed here as there will be one issued already from Kong's
			// Service override. The reason for this is that there is no other
			// object in Kubernetes that creates a Kong's Upstream and Kubernetes
			// Service will already trigger Kong's Service creation and log issuance.
			return
		}
	}

	u.overrideByKongIngress(kongIngress)
	if svc != nil {
		u.overrideByAnnotation(svc.Annotations)
	}
}
