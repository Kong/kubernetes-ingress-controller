package kongstate

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
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

	// Parent is the parent object of this Service.
	// It is expected to be a Kubernetes object which translation resulted in creating this Kong Service.
	// For example, if this Service was created as a result of translating a Kubernetes Ingress, then
	// Parent is expected to be the Ingress object itself.
	Parent client.Object
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

func (s *Service) overrideConnectTimeout(anns map[string]string) {
	if s == nil {
		return
	}
	timeout, exists := annotations.ExtractConnectTimeout(anns)
	if !exists {
		return
	}
	val, err := strconv.Atoi(timeout)
	if err != nil {
		return
	}
	s.ConnectTimeout = kong.Int(val)
}

func (s *Service) overrideWriteTimeout(anns map[string]string) {
	if s == nil {
		return
	}
	timeout, exists := annotations.ExtractWriteTimeout(anns)
	if !exists {
		return
	}
	val, err := strconv.Atoi(timeout)
	if err != nil {
		return
	}
	s.WriteTimeout = kong.Int(val)
}

func (s *Service) overrideReadTimeout(anns map[string]string) {
	if s == nil {
		return
	}
	timeout, exists := annotations.ExtractReadTimeout(anns)
	if !exists {
		return
	}
	val, err := strconv.Atoi(timeout)
	if err != nil {
		return
	}
	s.ReadTimeout = kong.Int(val)
}

func (s *Service) overrideRetries(anns map[string]string) {
	if s == nil {
		return
	}
	retries, exists := annotations.ExtractRetries(anns)
	if !exists {
		return
	}
	val, err := strconv.Atoi(retries)
	if err != nil {
		return
	}
	s.Retries = kong.Int(val)
}

func (s *Service) overrideTLSVerify(anns map[string]string) {
	if s == nil {
		return
	}
	tlsVerify, exists := annotations.ExtractTLSVerify(anns)
	if !exists {
		return
	}
	s.TLSVerify = kong.Bool(tlsVerify)
}

func (s *Service) overrideTLSVerifyDepth(anns map[string]string) {
	if s == nil {
		return
	}
	tlsVerifyDepth, exists := annotations.ExtractTLSVerifyDepth(anns)
	if !exists {
		return
	}
	s.TLSVerifyDepth = kong.Int(tlsVerifyDepth)
}

// overrideByAnnotation modifies the Kong service based on annotations
// on the Kubernetes service.
func (s *Service) overrideByAnnotation(anns map[string]string) {
	if s == nil {
		return
	}
	s.overrideProtocol(anns)
	s.overridePath(anns)
	s.overrideConnectTimeout(anns)
	s.overrideWriteTimeout(anns)
	s.overrideReadTimeout(anns)
	s.overrideRetries(anns)
	s.overrideTLSVerify(anns)
	s.overrideTLSVerifyDepth(anns)
}

// override sets Service fields using Kubernetes Service annotations.
func (s *Service) override() error {
	if s == nil {
		return nil
	}

	// Apply overrides from Kubernetes Service annotations. As we keep them in a map, let's first sort its keys to ensure
	// deterministic order of overrides.
	servicesNames := lo.Keys(s.K8sServices)
	sort.Strings(servicesNames)
	for _, serviceName := range servicesNames {
		svc := s.K8sServices[serviceName]
		s.overrideByAnnotation(svc.Annotations)
		protocol := annotations.ExtractProtocolName(svc.Annotations)
		if !util.ValidateProtocol(protocol) {
			return fmt.Errorf("%s annotation has invalid value: %s", annotations.AnnotationPrefix+annotations.ProtocolKey, protocol)
		}
	}

	if s.Protocol != nil && (*s.Protocol == "grpc" || *s.Protocol == "grpcs") {
		// grpc(s) doesn't accept a path
		s.Path = nil
	}
	return nil
}
