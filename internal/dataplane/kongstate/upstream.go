package kongstate

import (
	"github.com/kong/go-kong/kong"
	corev1 "k8s.io/api/core/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
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

// override sets Upstream using k8s Service's annotations.
func (u *Upstream) override(
	svc *corev1.Service,
) {
	if u == nil || svc == nil {
		return
	}

	u.overrideByAnnotation(svc.Annotations)
}
