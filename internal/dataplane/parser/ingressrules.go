package parser

import (
	"fmt"
	"strings"

	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	netv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/failures"
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
		result.SecretNameToSNIs.merge(obj.SecretNameToSNIs)
		for k, v := range obj.ServiceNameToServices {
			result.ServiceNameToServices[k] = v
		}
	}
	return result
}

// populateServices populates the ServiceNameToServices map with additional information
// and returns a map of services to be skipped.
func (ir *ingressRules) populateServices(log logrus.FieldLogger, s store.Storer, failuresCollector *failures.ResourceFailuresCollector) map[string]interface{} {
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

			if konnectServiceName := annotations.ExtractKonnectService(k8sService.Annotations); konnectServiceName != "" {
				service.Tags = append(service.Tags, kong.String(fmt.Sprintf("_KonnectService:%s", konnectServiceName)))
			}

			// extract client certificates intended for use by the service
			secretName := annotations.ExtractClientCertificate(k8sService.Annotations)
			if secretName != "" {
				secretKey := k8sService.Namespace + "/" + secretName
				secret, err := s.GetSecret(k8sService.Namespace, secretName)
				if err != nil {
					failuresCollector.PushResourceFailure(
						fmt.Sprintf("failed to fetch secret '%s': %v", secretKey, err), k8sService,
					)
					continue
				}

				// ensure that the cert is loaded into Kong
				ir.SecretNameToSNIs.addUniqueParents(secretKey, k8sService)
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

type SecretNameToSNIs struct {
	// secretToSNIs maps secrets (by 'namespace/name' key) to SNIs they are related to.
	secretToSNIs map[string]*SNIs

	// seenHosts keeps global hosts registry to make sure only one secret can refer a host.
	seenHosts map[string]struct{}
}

func newSecretNameToSNIs() SecretNameToSNIs {
	return SecretNameToSNIs{
		secretToSNIs: map[string]*SNIs{},
		seenHosts:    map[string]struct{}{},
	}
}

func (m SecretNameToSNIs) Parents(secretKey string) []client.Object {
	if _, ok := m.secretToSNIs[secretKey]; !ok {
		return nil
	}
	return m.secretToSNIs[secretKey].Parents()
}

func (m SecretNameToSNIs) Hosts(secretKey string) []string {
	if _, ok := m.secretToSNIs[secretKey]; !ok {
		return nil
	}
	return m.secretToSNIs[secretKey].Hosts()
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
		m.addUniqueHosts(secretKey, tls.Hosts...)
		m.addUniqueParents(secretKey, parent)
	}
}

// addUniqueHosts adds hosts to SNIs stored under a secretKey.
// It ensures that a host is not assigned to any secret yet. If it's assigned already, it will get skipped.
func (m SecretNameToSNIs) addUniqueHosts(secretKey string, hosts ...string) {
	m.ensureSNIsEntry(secretKey)

	for _, host := range hosts {
		if _, ok := m.seenHosts[host]; ok {
			// Skip this host, it's already assigned, possibly to another secret.
			continue
		}

		m.secretToSNIs[secretKey].hosts = append(m.secretToSNIs[secretKey].hosts, host)
		m.seenHosts[host] = struct{}{}
	}
}

// addUniqueParents adds parents to SNIs stored under a secretKey, ensuring their uniqueness by the object UID.
func (m SecretNameToSNIs) addUniqueParents(secretKey string, parents ...client.Object) {
	m.ensureSNIsEntry(secretKey)

	for _, parent := range parents {
		m.secretToSNIs[secretKey].parents[parent.GetUID()] = parent
	}
}

func (m SecretNameToSNIs) ensureSNIsEntry(secretKey string) {
	if _, ok := m.secretToSNIs[secretKey]; !ok {
		m.secretToSNIs[secretKey] = newSNIs()
	}
}

// merge merges other SecretNameToSNIs into m in place.
func (m SecretNameToSNIs) merge(o SecretNameToSNIs) {
	for secretKey, snis := range o.secretToSNIs {
		for _, obj := range snis.parents {
			m.addUniqueParents(secretKey, obj)
		}
		for _, hostKey := range snis.hosts {
			m.addUniqueHosts(secretKey, hostKey)
		}
	}
}

type SNIs struct {
	// parents are objects that the SNIs are inherited from
	parents map[types.UID]client.Object
	hosts   []string
}

func newSNIs() *SNIs {
	return &SNIs{
		parents: map[types.UID]client.Object{},
	}
}

func (s SNIs) Parents() []client.Object {
	parents := make([]client.Object, 0, len(s.parents))
	for _, p := range s.parents {
		parents = append(parents, p)
	}
	return parents
}

func (s SNIs) Hosts() []string {
	return s.hosts
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
	failuresCollector *failures.ResourceFailuresCollector,
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
				failuresCollector.PushResourceFailure(
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
				failuresCollector.PushResourceFailure(
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
