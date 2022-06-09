package kongstate

import (
	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
)

// Upstream is a wrapper around Upstream object in Kong.
type Upstream struct {
	kong.Upstream
	Targets []Target
	// Service this upstream is asosciated with.
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
func (u *Upstream) overrideByKongIngress(kongIngress *configurationv1.KongIngress) {
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
	// TODO https://github.com/Kong/kubernetes-ingress-controller/issues/2075
}

// override sets Upstream fields by KongIngress first, then by k8s Service's annotations
func (u *Upstream) override(
	log logrus.FieldLogger,
	kongIngress *configurationv1.KongIngress,
	svc *corev1.Service,
) {
	if u == nil {
		return
	}

	if u.Service.Parent != nil && kongIngress != nil {
		// If the parent object behind Kong Upstream's is a Gateway API object
		// (probably *Route but log a warning for all other objects as well)
		// then check if we're trying to override said Service configuration with
		// a KongIngress object and if that's the case then skip it since those
		// should not be affected.
		gvk := u.Service.Parent.GetObjectKind().GroupVersionKind()
		if gvk.Group == gatewayv1alpha2.GroupName {
			obj := u.Service.Parent
			fields := logrus.Fields{
				"resource_name":      obj.GetName(),
				"resource_namespace": obj.GetNamespace(),
				"resource_kind":      gvk.Kind,
			}
			if svc != nil {
				fields["service_name"] = svc.Name
				fields["service_namespace"] = svc.Namespace
			}
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
