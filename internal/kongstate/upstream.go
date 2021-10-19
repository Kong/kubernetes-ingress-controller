package kongstate

import (
	"github.com/kong/go-kong/kong"

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

	// The upstream within the KongIngress has no name.
	// As this overwrites the entire upstream object, we must restore the
	// original name after.
	name := *u.Upstream.Name
	u.Upstream = *kongIngress.Upstream.DeepCopy()
	u.Name = &name
}

// override sets Upstream fields by KongIngress first, then by annotation
func (u *Upstream) override(kongIngress *configurationv1.KongIngress, anns map[string]string) {
	if u == nil {
		return
	}

	u.overrideByKongIngress(kongIngress)
	u.overrideByAnnotation(anns)
}
