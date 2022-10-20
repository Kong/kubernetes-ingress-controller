package kongstate

import (
	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	corev1 "k8s.io/api/core/v1"
)

func getKongIngressForServices(
	s store.Storer,
	services map[string]*corev1.Service,
) (*configurationv1.KongIngress, error) {
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
			return nil, err
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
	obj util.K8sObjectInfo,
) (
	*configurationv1.KongIngress, error,
) {
	return getKongIngressFromObjAnnotations(s, obj)
}

func getKongIngressFromObjAnnotations(
	s store.Storer,
	obj util.K8sObjectInfo,
) (
	*configurationv1.KongIngress, error,
) {
	confName := annotations.ExtractConfigurationName(obj.Annotations)
	if confName != "" {
		ki, err := s.GetKongIngress(obj.Namespace, confName)
		if err == nil {
			return ki, nil
		}
	}

	ki, err := s.GetKongIngress(obj.Namespace, obj.Name)
	if err == nil {
		return ki, nil
	}
	return nil, nil
}

// PrettyPrintServiceList makes a clean printable list of a map of Kubernetes
// services for the purpose of logging (errors, info, e.t.c.).
func PrettyPrintServiceList(services map[string]*corev1.Service) string {
	var serviceList string
	first := true
	for _, svc := range services {
		if !first {
			serviceList = serviceList + ", "
		}
		serviceList = serviceList + svc.Namespace + "/" + svc.Name
		if first {
			first = false
		}
	}
	return serviceList
}
