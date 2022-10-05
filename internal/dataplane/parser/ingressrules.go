package parser

import (
	"strings"

	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	netv1beta1 "k8s.io/api/networking/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
)

type ingressRules struct {
	SecretNameToSNIs      SecretNameToSNIs
	ServiceNameToServices map[string]kongstate.Service
}

func newIngressRules() ingressRules {
	return ingressRules{
		SecretNameToSNIs:      newSecretNameToSNIs(),
		ServiceNameToServices: make(map[string]kongstate.Service),
	}
}

func mergeIngressRules(objs ...ingressRules) ingressRules {
	result := newIngressRules()

	for _, obj := range objs {
		for k, v := range obj.SecretNameToSNIs {
			result.SecretNameToSNIs[k] = append(result.SecretNameToSNIs[k], v...)
		}
		for k, v := range obj.ServiceNameToServices {
			result.ServiceNameToServices[k] = v
		}
	}
	return result
}

// populateServices populates the ServiceNameToServices map with additional information
// and returns a map of services to be skipped.
func (ir *ingressRules) populateServices(log logrus.FieldLogger, s store.Storer) map[string]interface{} {
	serviceNamesToSkip := make(map[string]interface{})

	// populate Kubernetes Service
	for key, service := range ir.ServiceNameToServices {
		if service.K8sServices == nil {
			service.K8sServices = make(map[string]*corev1.Service)
		}

		// collect all the Kubernetes services configured for the service backends,
		// and all the annotations in use across all services (when applicable).
		k8sServices, seenAnnotations := getK8sServicesForBackends(log, s, service.Namespace, service.Backends)

		// if the Kubernetes services have been deemed invalid, log an error message
		// and skip the current service.
		if !servicesAllUseTheSameKongAnnotations(log, k8sServices, seenAnnotations) {
			log.Errorf("the Kubernetes Services %v cannot have different sets of konghq.com annotations. "+
				"These Services are used in the same Gateway Route BackendRef together to create the Kong Service %s"+
				"and must use the same Kong annotations", k8sServices, *service.Name)
			// The Kong services not having all the k8s services correctly annotated must be marked
			// as to be skipped.
			serviceNamesToSkip[key] = nil
			continue
		}

		for _, k8sService := range k8sServices {
			// at this point we know the Kubernetes service itself is valid and can be
			// used for traffic, so cache it amongst the kong Services k8s services.
			service.K8sServices[k8sService.Name] = k8sService

			// extract client certificates intended for use by the service
			secretName := annotations.ExtractClientCertificate(k8sService.Annotations)
			if secretName != "" {
				secret, err := s.GetSecret(k8sService.Namespace, secretName)
				secretKey := k8sService.Namespace + "/" + secretName
				// ensure that the cert is loaded into Kong
				if _, ok := ir.SecretNameToSNIs[secretKey]; !ok {
					ir.SecretNameToSNIs[secretKey] = []string{}
				}
				if err == nil {
					service.ClientCertificate = &kong.Certificate{
						ID: kong.String(string(secret.UID)),
					}
				} else {
					log.WithFields(logrus.Fields{
						"secret_name":      secretName,
						"secret_namespace": k8sService.Namespace,
					}).Errorf("failed to fetch secret: %v", err)
				}
			}
		}

		// Kubernetes Services have been populated for this Kong Service, so it can
		// now be cached.
		ir.ServiceNameToServices[key] = service
	}
	return serviceNamesToSkip
}

type SecretNameToSNIs map[string][]string

func newSecretNameToSNIs() SecretNameToSNIs {
	return SecretNameToSNIs(map[string][]string{})
}

func (m SecretNameToSNIs) addFromIngressV1beta1TLS(tlsSections []netv1beta1.IngressTLS, namespace string) {
	// Assume that v1beta1 and v1 tlsSections have identical semantics and field-wise content.
	var v1 []netv1.IngressTLS
	for _, item := range tlsSections {
		v1 = append(v1, netv1.IngressTLS{Hosts: item.Hosts, SecretName: item.SecretName})
	}
	m.addFromIngressV1TLS(v1, namespace)
}

func (m SecretNameToSNIs) addFromIngressV1TLS(tlsSections []netv1.IngressTLS, namespace string) {
	for _, tls := range tlsSections {
		if len(tls.Hosts) == 0 {
			continue
		}
		if tls.SecretName == "" {
			continue
		}
		hosts := tls.Hosts
		secretName := namespace + "/" + tls.SecretName
		hosts = m.filterHosts(hosts)
		if m[secretName] != nil {
			hosts = append(hosts, m[secretName]...)
		}
		m[secretName] = hosts
	}
}

func (m SecretNameToSNIs) filterHosts(hosts []string) []string {
	hostsToAdd := []string{}
	seenHosts := map[string]bool{}
	for _, hosts := range m {
		for _, host := range hosts {
			seenHosts[host] = true
		}
	}
	for _, host := range hosts {
		if !seenHosts[host] {
			hostsToAdd = append(hostsToAdd, host)
		}
	}
	return hostsToAdd
}

func getK8sServicesForBackends(
	log logrus.FieldLogger,
	storer store.Storer,
	namespace string,
	backends kongstate.ServiceBackends,
) ([]*corev1.Service, map[string]string) {
	// we collect all annotations seen for this group of services so that these
	// can be later validated.
	seenAnnotationsForK8sServices := make(map[string]string)

	// for each backend (which is a reference to a Kubernetes Service object)
	// retreieve that backend and capture any Kong annotations its using.
	k8sServices := make([]*corev1.Service, 0, len(backends))
	for _, backend := range backends {
		backendNamespace := namespace
		if backend.Namespace != "" {
			backendNamespace = backend.Namespace
		}
		k8sService, err := storer.GetService(backendNamespace, backend.Name)
		if err != nil {
			log.WithFields(logrus.Fields{
				"service_name":      backend.PortDef.Name,
				"service_namespace": backendNamespace,
			}).Errorf("failed to fetch service: %v", err)
			continue
		}
		if k8sService != nil {
			// record all Kong annotations in use by the service
			for k, v := range k8sService.GetAnnotations() {
				if strings.HasPrefix(k, annotations.AnnotationPrefix) {
					seenAnnotationsForK8sServices[k] = v
				}
			}

			// add the service to the list of backend services
			k8sServices = append(k8sServices, k8sService)
		}
	}

	return k8sServices, seenAnnotationsForK8sServices
}

func servicesAllUseTheSameKongAnnotations(
	log logrus.FieldLogger,
	services []*corev1.Service,
	annotations map[string]string,
) bool {
	match := true
	for _, service := range services {
		// all services grouped together via backends must have identical annotations
		// to avoid unexpected routing behaviors.
		//
		// TODO: ultimately we should be able to do this validation in our normal
		// validation layer, but we're limited at present on where and how that
		// validation can work. We should be able to move this validation there
		// once https://github.com/Kong/kubernetes-ingress-controller/issues/2195
		// is resolved.
		for k, v := range annotations {
			valueForThisObject, ok := service.Annotations[k]
			if !ok {
				log.WithFields(logrus.Fields{
					"service_name":      service.Name,
					"service_namespace": service.Namespace,
				}).Errorf("in the backend group of %d kubernetes services some have the %s annotation while others don't. "+
					"this is not supported: when multiple services comprise a backend all kong annotations "+
					"between them must be set to the same value", len(services), k)
				match = false
			}

			if valueForThisObject != v {
				log.WithFields(logrus.Fields{
					"service_name":      service.Name,
					"service_namespace": service.Namespace,
				}).Errorf("the value of annotation %s is different between the %d services which comprise this backend. "+
					"this is not supported: when multiple services comprise a backend all kong annotations "+
					"between them must be set to the same value", k, len(services))
				match = false
			}
		}
	}

	return match
}
