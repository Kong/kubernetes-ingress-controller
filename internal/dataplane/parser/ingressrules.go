package parser

import (
	"fmt"
	"strings"

	netv1beta1 "k8s.io/api/networking/v1beta1"

	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
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
			result.SecretNameToSNIs.addHosts(k, v.hosts)
			result.SecretNameToSNIs.addParents(k, v.parents)
		}
		for k, v := range obj.ServiceNameToServices {
			result.ServiceNameToServices[k] = v
		}
	}
	return result
}

// populateServices populates the ServiceNameToServices map with additional information
// and returns a map of services to be skipped.
func (ir *ingressRules) populateServices(log logrus.FieldLogger, s store.Storer, failuresCollector *TranslationFailuresCollector) map[string]interface{} {
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
		if !servicesAllUseTheSameKongAnnotations(k8sServices, seenAnnotations, failuresCollector) {
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
				secretKey := k8sService.Namespace + "/" + secretName
				secret, err := s.GetSecret(k8sService.Namespace, secretName)
				if err != nil {
					failuresCollector.PushTranslationFailure(
						fmt.Sprintf("failed to fetch secret '%s': %v", secretKey, err), k8sService,
					)
					continue
				}

				// ensure that the cert is loaded into Kong
				if _, ok := ir.SecretNameToSNIs[secretKey]; !ok {
					ir.SecretNameToSNIs[secretKey] = &SNIs{
						parents: []client.Object{k8sService},
						hosts:   []string{},
					}
				}
				service.ClientCertificate = &kong.Certificate{
					ID: kong.String(string(secret.UID)),
				}
			}
		}

		// Kubernetes Services have been populated for this Kong Service, so it can
		// now be cached.
		ir.ServiceNameToServices[key] = service
	}
	return serviceNamesToSkip
}

type SecretNameToSNIs map[string]*SNIs

type SNIs struct {
	parents []client.Object
	hosts   []string
}

func newSecretNameToSNIs() SecretNameToSNIs {
	return map[string]*SNIs{}
}

func (m SecretNameToSNIs) addFromIngressV1TLS(tlsSections []netv1.IngressTLS, parent client.Object) {
	for _, tls := range tlsSections {
		if len(tls.Hosts) == 0 {
			continue
		}
		if tls.SecretName == "" {
			continue
		}

		secretKey := parent.GetNamespace() + "/" + tls.SecretName
		m.addHosts(secretKey, tls.Hosts)
		m.addParents(secretKey, []client.Object{parent})
	}
}

func (m SecretNameToSNIs) addHosts(secretKey string, hosts []string) {
	if _, ok := m[secretKey]; !ok {
		m[secretKey] = &SNIs{}
	}

	seenHosts := map[string]bool{}
	for _, snis := range m {
		for _, host := range snis.hosts {
			seenHosts[host] = true
		}
	}

	var hostsToAdd []string
	for _, host := range hosts {
		if !seenHosts[host] {
			hostsToAdd = append(hostsToAdd, host)
		}
	}

	m[secretKey].hosts = append(m[secretKey].hosts, hostsToAdd...)
}

func (m SecretNameToSNIs) addParents(secretKey string, parents []client.Object) {
	if _, ok := m[secretKey]; !ok {
		m[secretKey] = &SNIs{}
	}

	seenHosts := map[types.UID]bool{}
	for _, parent := range m[secretKey].parents {
		seenHosts[parent.GetUID()] = true
	}

	var parentsToAdd []client.Object
	for _, parent := range parents {
		if !seenHosts[parent.GetUID()] {
			parentsToAdd = append(parentsToAdd, parent)
		}
	}

	m[secretKey].parents = append(m[secretKey].parents, parentsToAdd...)
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
	services []*corev1.Service,
	annotations map[string]string,
	failuresCollector *TranslationFailuresCollector,
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
				failuresCollector.PushTranslationFailure(
					fmt.Sprintf("in the backend group of %d kubernetes services some have the %s annotation while others don't. "+
						"this is not supported: when multiple services comprise a backend all kong annotations "+
						"between them must be set to the same value", len(services), k),
					service.DeepCopy(),
				)
				match = false
				// continue as it doesn't make sense to verify value of not existing annotation
				continue
			}

			if valueForThisObject != v {
				failuresCollector.PushTranslationFailure(
					fmt.Sprintf("the value of annotation %s is different between the %d services which comprise this backend. "+
						"this is not supported: when multiple services comprise a backend all kong annotations "+
						"between them must be set to the same value", k, len(services)),
					service.DeepCopy(),
				)
				match = false
			}
		}
	}

	return match
}

func v1beta1toV1TLS(tlsSections []netv1beta1.IngressTLS) []netv1.IngressTLS {
	var v1 []netv1.IngressTLS
	for _, item := range tlsSections {
		v1 = append(v1, netv1.IngressTLS{Hosts: item.Hosts, SecretName: item.SecretName})
	}
	return v1
}
