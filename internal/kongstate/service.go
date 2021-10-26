package kongstate

import (
	"strings"

	"github.com/kong/go-kong/kong"
	corev1 "k8s.io/api/core/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
)

// Service represents a service in Kong and holds routes associated with the
// service and other k8s metadata.
type Service struct {
	kong.Service
	Backend    ServiceBackend
	Namespace  string
	Routes     []Route
	Plugins    []kong.Plugin
	K8sService corev1.Service
}

// overrideByKongIngress sets Service fields by KongIngress
func (s *Service) overrideByKongIngress(kongIngress *configurationv1.KongIngress) {
	if kongIngress == nil || kongIngress.Proxy == nil {
		return
	}
	p := kongIngress.Proxy
	if p.Protocol != nil {
		s.Protocol = kong.String(*p.Protocol)
	}
	if p.Path != nil {
		s.Path = kong.String(*p.Path)
	}
	if p.Retries != nil {
		s.Retries = kong.Int(*p.Retries)
	}
	if p.ConnectTimeout != nil {
		s.ConnectTimeout = kong.Int(*p.ConnectTimeout)
	}
	if p.ReadTimeout != nil {
		s.ReadTimeout = kong.Int(*p.ReadTimeout)
	}
	if p.WriteTimeout != nil {
		s.WriteTimeout = kong.Int(*p.WriteTimeout)
	}
}

func (s *Service) overridePath(anns map[string]string) {
	if s == nil {
		return
	}
	path := annotations.ExtractPath(anns)
	if path == "" {
		return
	}
	// kong errors if path doesn't start with `/`
	if !strings.HasPrefix(path, "/") {
		return
	}
	s.Path = kong.String(path)
}

func (s *Service) overrideProtocol(anns map[string]string) {
	if s == nil {
		return
	}
	protocol := annotations.ExtractProtocolName(anns)
	if protocol == "" || !util.ValidateProtocol(protocol) {
		return
	}
	s.Protocol = kong.String(protocol)
}

// overrideByAnnotation modifies the Kong service based on annotations
// on the Kubernetes service.
func (s *Service) overrideByAnnotation(anns map[string]string) {
	if s == nil {
		return
	}
	s.overrideProtocol(anns)
	s.overridePath(anns)
}

// override sets Service fields by KongIngress first, then by annotation
func (s *Service) override(kongIngress *configurationv1.KongIngress,
	anns map[string]string) {
	if s == nil {
		return
	}

	s.overrideByKongIngress(kongIngress)
	s.overrideByAnnotation(anns)

	if *s.Protocol == "grpc" || *s.Protocol == "grpcs" {
		// grpc(s) doesn't accept a path
		s.Path = nil
	}
}
