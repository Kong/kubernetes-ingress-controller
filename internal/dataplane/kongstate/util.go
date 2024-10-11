package kongstate

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kong/go-kong/kong"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
)

func getKongIngressForServices(
	s store.Storer,
	services []*corev1.Service,
) (*kongv1.KongIngress, error) {
	// loop through each service and retrieve the attached KongIngress resources.
	// there can only be one KongIngress for a group of services: either one of
	// them is configured with a KongIngress and this configures the Kong Service
	// or Upstream OR all of them can be configured but they must be configured
	// with the same KongIngress.
	for _, svc := range services {
		// check if the service is even configured with a KongIngress
		confName := annotations.ExtractConfigurationName(svc.Annotations)
		if confName == "" {
			continue // some other service in the group may yet have a KongIngress attachment
		}

		// retrieve the attached KongIngress for the service
		kongIngress, err := s.GetKongIngress(svc.Namespace, confName)
		if err != nil {
			return nil, fmt.Errorf("failed to get KongIngress: %w", err)
		}

		// we found the KongIngress for these services. We don't have to check any
		// further services as validation is expected to ensure all these Services
		// already are annotated with the exact same overrides.
		return kongIngress, nil
	}

	// there are no KongIngress resources for these services.
	return nil, nil
}

func getKongIngressFromObjectMeta(
	s store.Storer,
	obj client.Object,
) (
	*kongv1.KongIngress, error,
) {
	return getKongIngressFromObjAnnotations(s, obj)
}

func getKongIngressFromObjAnnotations(
	s store.Storer,
	obj client.Object,
) (
	*kongv1.KongIngress, error,
) {
	confName := annotations.ExtractConfigurationName(obj.GetAnnotations())
	if confName != "" {
		ki, err := s.GetKongIngress(obj.GetNamespace(), confName)
		if err == nil {
			return ki, nil
		}
	}

	ki, err := s.GetKongIngress(obj.GetNamespace(), obj.GetName())
	if err == nil {
		return ki, nil
	}
	return nil, nil
}

// prettyPrintServiceList makes a clean printable list of Kubernetes
// services for the purpose of logging (errors, info, etc.).
func prettyPrintServiceList(services []*corev1.Service) string {
	serviceList := make([]string, 0, len(services))
	for _, svc := range services {
		serviceList = append(serviceList, svc.Namespace+"/"+svc.Name)
	}
	return strings.Join(serviceList, ", ")
}

// RawConfigToConfiguration decodes raw JSON to the format of Kong configuration.
// it is run after all patches applied to the initial config.
func RawConfigToConfiguration(raw []byte) (kong.Configuration, error) {
	if len(raw) == 0 {
		return kong.Configuration{}, nil
	}
	var kongConfig kong.Configuration
	err := json.Unmarshal(raw, &kongConfig)
	if err != nil {
		return kong.Configuration{}, err
	}
	return kongConfig, nil
}
