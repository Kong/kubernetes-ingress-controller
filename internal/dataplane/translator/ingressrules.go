package translator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/logging"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

const (
	// defaultServiceProtocol is the default protocol used by Kong for services.
	defaultServiceProtocol = "http"
)

type ingressRules struct {
	SecretNameToSNIs      SecretNameToSNIs
	ServiceNameToServices map[string]kongstate.Service
	ServiceNameToParent   map[string]client.Object
}

func newIngressRules() ingressRules {
	return ingressRules{
		SecretNameToSNIs:      newSecretNameToSNIs(),
		ServiceNameToServices: make(map[string]kongstate.Service),
		ServiceNameToParent:   make(map[string]client.Object),
	}
}

func mergeIngressRules(objs ...ingressRules) ingressRules {
	result := newIngressRules()

	for _, obj := range objs {
		result.SecretNameToSNIs.merge(obj.SecretNameToSNIs)
		for k, v := range obj.ServiceNameToServices {
			result.ServiceNameToServices[k] = v
		}
		for k, v := range obj.ServiceNameToParent {
			result.ServiceNameToParent[k] = v
		}
	}
	return result
}

// populateServices populates the ServiceNameToServices map with additional information
// and returns a map of services to be skipped.
func (ir *ingressRules) populateServices(
	logger logr.Logger,
	s store.Storer,
	failuresCollector *failures.ResourceFailuresCollector,
	translatedObjectsCollector *ObjectsCollector,
) map[string]interface{} {
	serviceNamesToSkip := make(map[string]interface{})

	// populate Kubernetes Service
	for key, service := range ir.ServiceNameToServices {
		if service.K8sServices == nil {
			service.K8sServices = make(map[string]*corev1.Service)
		}

		// collect all the Kubernetes services configured for the service backends,
		// and all annotations with our prefix in use across all services (when applicable).
		serviceParent := ir.ServiceNameToParent[key]
		k8sServices, seenKongAnnotations := getK8sServicesForBackends(s, service.Backends, translatedObjectsCollector, failuresCollector, serviceParent)

		// if the Kubernetes services have been deemed invalid, log an error message
		// and skip the current service.
		if !collectInconsistentAnnotations(k8sServices, seenKongAnnotations, failuresCollector, key) {
			// The Kong services not having all the k8s services correctly annotated must be marked
			// as to be skipped.
			serviceNamesToSkip[key] = nil
			continue
		}

		for _, k8sService := range k8sServices {
			// We need to create a copy of the k8s service as we need to modify it. The original is read
			// by another routine and we may incur in a data race.
			k8sServiceCopy := k8sService.DeepCopy()

			// Convert the backendTLSPolicy targeting the service to the proper set of annotations.
			ir.handleBackendTLSPolices(s, k8sServiceCopy, failuresCollector)

			// Extract client certificates intended for use by the service.
			ir.handleServiceClientCertificates(s, k8sServiceCopy, &service, failuresCollector)

			// Extract CA certificates intended for use by the service.
			ir.handleServiceCACertificates(s, k8sServiceCopy, &service, failuresCollector)

			// at this point we know the Kubernetes service itself is valid and can be
			// used for traffic, so cache it amongst the kong Services k8s services.
			service.K8sServices[fmt.Sprintf("%s/%s", k8sServiceCopy.Namespace, k8sServiceCopy.Name)] = k8sServiceCopy
		}
		service.Tags = ir.generateKongServiceTags(k8sServices, service, logger)

		// Kubernetes Services have been populated for this Kong Service, so it can
		// now be cached.
		ir.ServiceNameToServices[key] = service
	}
	return serviceNamesToSkip
}

func (ir *ingressRules) handleBackendTLSPolices(
	s store.Storer,
	k8sService *corev1.Service,
	failuresCollector *failures.ResourceFailuresCollector,
) {
	policies, err := s.ListBackendTLSPoliciesByTargetService(client.ObjectKeyFromObject(k8sService))
	if err != nil {
		failuresCollector.PushResourceFailure(
			fmt.Sprintf("Failed to list backendTLSPolicies: %v", err), k8sService,
		)
		return
	}
	if len(policies) == 0 {
		return
	}
	if len(policies) > 1 {
		failuresCollector.PushResourceFailure(
			"Multiple BackendTLSPolicies attached to service", k8sService,
		)
		return
	}
	policy := policies[0]

	if k8sService.Annotations == nil {
		k8sService.Annotations = make(map[string]string)
	}

	annotations.SetTLSVerify(k8sService.Annotations, true)
	annotations.SetHostHeader(k8sService.Annotations, string(policy.Spec.Validation.Hostname))
	annotations.SetProtocol(k8sService.Annotations, "https")
	annotations.SetCACertificates(k8sService.Annotations,
		lo.Map(policy.Spec.Validation.CACertificateRefs, func(ref gatewayapi.LocalObjectReference, _ int) string {
			return string(ref.Name)
		}),
	)
	if depth, ok := getTLSVerifyDepthOption(policy.Spec.Options); ok {
		annotations.SetTLSVerifyDepth(k8sService.Annotations, depth)
	}
}

func getTLSVerifyDepthOption(options map[gatewayapi.AnnotationKey]gatewayapi.AnnotationValue) (int, bool) {
	const (
		tlsVerifyDepthKey = "tls-verify-depth"
	)

	// If the annotation is not set, return no depth.
	depthStr, ok := options[tlsVerifyDepthKey]
	if !ok {
		return 0, false
	}

	// If the annotation is not an int, return no depth.
	depth, err := strconv.Atoi(string(depthStr))
	if err != nil {
		return 0, false
	}
	// If the annotation is < 0, return no depth.
	if depth < 0 {
		return 0, false
	}

	return depth, true
}

func (ir *ingressRules) handleServiceClientCertificates(
	s store.Storer,
	k8sService *corev1.Service,
	service *kongstate.Service,
	failuresCollector *failures.ResourceFailuresCollector,
) {
	secretName := annotations.ExtractClientCertificate(k8sService.Annotations)
	if secretName != "" {
		secretKey := k8sService.Namespace + "/" + secretName
		secret, err := s.GetSecret(k8sService.Namespace, secretName)
		if err != nil {
			failuresCollector.PushResourceFailure(
				fmt.Sprintf("Failed to fetch secret '%s': %v", secretKey, err), k8sService,
			)
			return
		}

		// override protocol isn't set yet, need to get it from the annotation
		protocol := getEffectiveServiceProtocol(k8sService)
		if isNonTLSProtocol(protocol) {
			failuresCollector.PushResourceFailure(
				fmt.Sprintf("Client certificate requested for incompatible service protocol '%s'", *service.Protocol),
				k8sService,
			)
			return
		}
		// ensure that the cert is loaded into Kong
		ir.SecretNameToSNIs.addUniqueParents(secretKey, k8sService)
		service.ClientCertificate = &kong.Certificate{
			ID: kong.String(string(secret.UID)),
		}
	}
}

func (ir *ingressRules) handleServiceCACertificates(
	s store.Storer,
	service *corev1.Service,
	k *kongstate.Service,
	collector *failures.ResourceFailuresCollector,
) {
	secretcertificates := annotations.ExtractCACertificateSecretNames(service.Annotations)
	configMapCertificates := annotations.ExtractCACertificateConfigMapNames(service.Annotations)
	if len(secretcertificates)+len(configMapCertificates) == 0 {
		// No CA certificates to process.
		return
	}

	// Validate that the service has TLS verification turned on.
	if v, ok := annotations.ExtractTLSVerify(service.Annotations); !ok || !v {
		collector.PushResourceFailure(
			"CA certificates requested for service without TLS verification enabled",
			service,
		)
		return
	}

	// Validate that the effective service protocol is compatible with the CA certificates.
	protocol := getEffectiveServiceProtocol(service)
	if isNonTLSProtocol(protocol) {
		collector.PushResourceFailure(
			fmt.Sprintf("CA certificates requested for incompatible service protocol '%s'", protocol),
			service,
		)
		return
	}

	// Process each CA certificate from secret and add it to the Kong Service.
	for _, certificate := range secretcertificates {
		secretKey := service.Namespace + "/" + certificate
		secret, err := s.GetSecret(service.Namespace, certificate)
		if err != nil {
			collector.PushResourceFailure(
				fmt.Sprintf("Failed to fetch secret for CA Certificate '%s': %v", secretKey, err), service,
			)
			continue
		}

		certID, ok := secret.Data["id"]
		if !ok {
			collector.PushResourceFailure(
				fmt.Sprintf("Invalid CA certificate '%s': missing 'id' field in data", secretKey), secret, service,
			)
			continue
		}
		k.CACertificates = append(k.CACertificates, lo.ToPtr(string(certID)))
	}

	// Process each CA certificate from ConfigMap and add it to the Kong Service.
	for _, certificate := range configMapCertificates {
		configmapKey := service.Namespace + "/" + certificate
		configMap, err := s.GetConfigMap(service.Namespace, certificate)
		if err != nil {
			collector.PushResourceFailure(
				fmt.Sprintf("Failed to fetch configmap for CA Certificate '%s': %v", configmapKey, err), service,
			)
			continue
		}

		certID, ok := configMap.Data["id"]
		if !ok {
			collector.PushResourceFailure(
				fmt.Sprintf("Invalid CA certificate '%s': missing 'id' field in data", configmapKey), configMap, service,
			)
			continue
		}
		k.CACertificates = append(k.CACertificates, lo.ToPtr(certID))
	}
}

func getEffectiveServiceProtocol(svc *corev1.Service) string {
	protocol := annotations.ExtractProtocolName(svc.Annotations)
	if protocol == "" {
		// Annotation value does not indicate the effective default on Kong side.
		protocol = defaultServiceProtocol
	}
	return protocol
}

func (ir *ingressRules) generateKongServiceTags(
	k8sServices []*corev1.Service,
	service kongstate.Service,
	logger logr.Logger,
) []*string {
	// For multi-backend Services we expect ServiceNameToParent to be populated.
	if len(k8sServices) > 1 {
		if parent, ok := ir.ServiceNameToParent[*service.Name]; ok {
			return util.GenerateTagsForObject(parent)
		}
		logger.Error(nil, "Multi-service backend lacks parent info, cannot generate tags",
			"service", *service.Name)
		return nil
	}

	// For single-backend Services we ...
	if len(k8sServices) == 1 {
		// ... either use the parent object of the Service when its backend is a KongServiceFacade ...
		if len(service.Backends) == 1 && service.Backends[0].IsServiceFacade() {
			return util.GenerateTagsForObject(service.Parent)
		}
		// ... or use the backing Kubernetes Service.
		return util.GenerateTagsForObject(k8sServices[0])
	}

	// This tag generation code runs _before_ we would discard routes that are invalid because their backend
	// Service doesn't actually exist. attempting to generate tags for that Service would trigger a panic.
	// The translator should discard this invalid route later, but this adds a placeholder value in case it doesn't.
	// If you encounter an actual config where a service has these tags, something strange has happened.
	logger.V(logging.DebugLevel).Info("Service has zero k8sServices backends, cannot generate tags for it properly",
		"service", *service.Name)
	return kong.StringSlice(
		util.K8sNameTagPrefix+"UNKNOWN",
		util.K8sNamespaceTagPrefix+"UNKNOWN",
		util.K8sKindTagPrefix+"Service",
		util.K8sUIDTagPrefix+"00000000-0000-0000-0000-000000000000",
		util.K8sGroupTagPrefix+"core",
		util.K8sVersionTagPrefix+"v1",
	)
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
	parents map[k8stypes.UID]client.Object
	hosts   []string
}

func newSNIs() *SNIs {
	return &SNIs{
		parents: map[k8stypes.UID]client.Object{},
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
	storer store.Storer,
	backends kongstate.ServiceBackends,
	translatedObjectsCollector *ObjectsCollector,
	failuresCollector *failures.ResourceFailuresCollector,
	parent client.Object,
) ([]*corev1.Service, map[string]string) {
	// we collect all annotations seen for this group of services so that these
	// can be later validated.
	seenAnnotationsForK8sServices := make(map[string]string)

	// for each backend (which is a reference to a Kubernetes Service object)
	// retreieve that backend and capture any Kong annotations its using.
	k8sServices := make([]*corev1.Service, 0, len(backends))
	for _, backend := range backends {
		k8sService, err := resolveKubernetesServiceForBackend(storer, backend, translatedObjectsCollector)
		if err != nil {
			failuresCollector.PushResourceFailure(fmt.Sprintf("failed to resolve Kubernetes Service for backend: %s", err), parent)
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

func resolveKubernetesServiceForBackend(
	storer store.Storer,
	backend kongstate.ServiceBackend,
	translatedObjectsCollector *ObjectsCollector,
) (*corev1.Service, error) {
	// In case of KongServiceFacade, we need to fetch it to determine the Kubernetes Service backing it.
	// We also want to use its annotations as they override the annotations of the Kubernetes Service.
	if backend.IsServiceFacade() {
		svcFacade, err := storer.GetKongServiceFacade(backend.Namespace(), backend.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to fetch KongServiceFacade %s/%s: %w", backend.Namespace(), backend.Name(), err)
		}
		k8sService, err := storer.GetService(backend.Namespace(), svcFacade.Spec.Backend.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch Service %s/%s: %w", backend.Namespace(), svcFacade.Spec.Backend.Name, err)
		}

		// Make a copy of the Kubernetes Service to avoid mutating the cache (the k8sService we
		// get from the storer is a pointer).
		k8sService = k8sService.DeepCopy()

		// Merge the annotations from the KongServiceFacade with the annotations from the Service.
		// KongServiceFacade overrides the Service annotations if they have the same key.
		for k, v := range svcFacade.GetAnnotations() {
			if k8sService.Annotations == nil {
				k8sService.Annotations = make(map[string]string)
			}
			k8sService.Annotations[k] = v
		}

		// After KongServiceFacade's backing Service is fetched successfully, we can consider it a translated object.
		translatedObjectsCollector.Add(svcFacade)

		return k8sService, nil
	}

	// In case of Kubernetes Service, we just need to fetch it.
	k8sService, err := storer.GetService(backend.Namespace(), backend.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Service %s/%s: %w", backend.Namespace(), backend.Name(), err)
	}

	// After Kubernetes Service is fetched successfully, we can consider it a translated object.
	translatedObjectsCollector.Add(k8sService)

	return k8sService, nil
}

// collectInconsistentAnnotations takes a list of services and annotation+value pairs and confirms that all services
// have those annotations with those values. If any service does not have one of the annotation+value pairs, push
// a resource failure to the provided collector for all services indicating the problem annotation.
func collectInconsistentAnnotations(
	services []*corev1.Service,
	annotations map[string]string,
	failuresCollector *failures.ResourceFailuresCollector,
	kongServiceName string,
) bool {
	match := true
	badAnnotations := sets.Set[string]{}
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
				badAnnotations.Insert(k)
				match = false
				// continue as it doesn't make sense to verify value of not existing annotation
				continue
			}

			if valueForThisObject != v {
				badAnnotations.Insert(k)
				match = false
			}
		}
	}

	for _, service := range services {
		for _, annotation := range badAnnotations.UnsortedList() {
			failuresCollector.PushResourceFailure(
				fmt.Sprintf("Service has inconsistent %s annotation and is used in multi-Service backend %s. "+
					"This annotation must have the same value across all Services in the backend.",
					annotation, kongServiceName),
				service.DeepCopy(),
			)
		}
	}

	return match
}

// isNonTLSProtocol returns true if the protocol is a non-TLS protocol.
func isNonTLSProtocol(proto string) bool {
	// Strings used here for comparison reflect Kong's protocol names.
	nonTLSProtocols := []string{"http", "grpc", "tcp", "tls_passthrough", "udp", "ws"}
	return lo.Contains(nonTLSProtocols, proto)
}
