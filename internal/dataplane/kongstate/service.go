package kongstate

import (
	"strings"

	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
)

// Services is a list of kongstate.Service objects with sorting enabled based
// on a lexographical comparison of the underlying kong.Service names which are
// always expected to be unique.
type Services []*Service

func (s Services) Len() int {
	return len(s)
}

func (s Services) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Services) Less(i, j int) bool {
	a := ""
	if s[i].Service.Name != nil {
		a = *s[i].Service.Name
	}

	b := ""
	if s[j].Service.Name != nil {
		b = *s[j].Service.Name
	}

	return strings.Compare(a, b) == -1
}

// Service represents a service in Kong and holds routes associated with the
// service and other k8s metadata.
type Service struct {
	kong.Service
	Namespace string
	Routes    []Route
	Plugins   []kong.Plugin

	Backends    []ServiceBackend
	K8sServices map[string]*corev1.Service
	Parent      client.Object
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

// override sets Service fields by KongIngress first, then by k8s Service's annotations
func (s *Service) override(
	log logrus.FieldLogger,
	kongIngress *configurationv1.KongIngress,
	svc *corev1.Service,
) {
	if s == nil {
		return
	}

	if s.Parent != nil && kongIngress != nil {
		kongIngressFromSvcAnnotation := annotations.ExtractConfigurationName(svc.Annotations)
		if kongIngressFromSvcAnnotation != "" {
			// If the parent object behind Kong Service is a Gateway API object
			// (probably *Route but log a warning for all other objects as well)
			// then check if we're trying to override said Service configuration with
			// a KongIngress object and if that's the case then skip it since those
			// should not be affected.
			gvk := s.Parent.GetObjectKind().GroupVersionKind()
			if gvk.Group == gatewayv1alpha2.GroupName {
				obj := s.Parent
				fields := logrus.Fields{
					"resource_name":      obj.GetName(),
					"resource_namespace": obj.GetNamespace(),
					"resource_kind":      gvk.Kind,
				}
				if svc != nil {
					fields["service_name"] = svc.Name
					fields["service_namespace"] = svc.Namespace
				}
				log.WithFields(fields).
					Warn("KongIngress annotation is not allowed on Services " +
						"referenced by Gateway API *Route objects.")
				return
			}
		}
	}

	s.overrideByKongIngress(kongIngress)
	if svc != nil {
		s.overrideByAnnotation(svc.Annotations)
	}

	if *s.Protocol == "grpc" || *s.Protocol == "grpcs" {
		// grpc(s) doesn't accept a path
		s.Path = nil
	}
}
