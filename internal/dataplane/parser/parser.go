package parser

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net"
	"sort"
	"strings"

	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/samber/mo"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager/featuregates"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1alpha1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
)

// -----------------------------------------------------------------------------
// Parser - Public Constants and Package Variables
// -----------------------------------------------------------------------------

const (
	KindGateway = gatewayapi.Kind("Gateway")

	// kongRouterFlavorExpressions is the value used in router_flavor of kong configuration
	// to enable expression based router of kong.
	kongRouterFlavorExpressions = "expressions"
)

// -----------------------------------------------------------------------------
// Parser - Public Types
// -----------------------------------------------------------------------------

// FeatureFlags are used to control the behavior of the parser.
type FeatureFlags struct {
	// ReportConfiguredKubernetesObjects turns on object reporting for this parser:
	// each subsequent call to BuildKongConfig() will track the Kubernetes objects which
	// were successfully parsed.
	ReportConfiguredKubernetesObjects bool

	// ExpressionRoutes indicates whether to translate Kubernetes objects to expression based Kong Routes.
	ExpressionRoutes bool

	// FillIDs enables the parser to fill in the IDs fields of Kong entities - Services, Routes, and Consumers - based
	// on their names. It ensures that IDs remain stable across restarts of the controller.
	FillIDs bool

	// RewriteURIs enables the parser to translate the konghq.com/rewrite annotation to the proper set of Kong plugins.
	RewriteURIs bool
}

func NewFeatureFlags(
	logger logr.Logger,
	featureGates featuregates.FeatureGates,
	routerFlavor string,
	updateStatusFlag bool,
) FeatureFlags {
	return FeatureFlags{
		ReportConfiguredKubernetesObjects: updateStatusFlag,
		ExpressionRoutes:                  shouldEnableParserExpressionRoutes(logger, featureGates, routerFlavor),
		FillIDs:                           featureGates.Enabled(featuregates.FillIDsFeature),
		RewriteURIs:                       featureGates.Enabled(featuregates.RewriteURIsFeature),
	}
}

func shouldEnableParserExpressionRoutes(
	logger logr.Logger,
	featureGates featuregates.FeatureGates,
	routerFlavor string,
) bool {
	if !featureGates.Enabled(featuregates.ExpressionRoutesFeature) {
		return false
	}
	if routerFlavor != kongRouterFlavorExpressions {
		logger.V(util.InfoLevel).Info("ExpressionRoutes feature gate enabled but Gateway is running with incompatible router flavor, using that instead", "flavor", routerFlavor)
		return false
	}
	logger.V(util.InfoLevel).Info("expression routes mode enabled")
	return true
}

// LicenseGetter is an interface for getting the Kong Enterprise license.
type LicenseGetter interface {
	// GetLicense returns an optional license.
	GetLicense() mo.Option[kong.License]
}

// Parser parses Kubernetes objects and configurations into their
// equivalent Kong objects and configurations, producing a complete
// state configuration for the Kong Admin API.
type Parser struct {
	logger        logr.Logger
	storer        store.Storer
	licenseGetter LicenseGetter
	featureFlags  FeatureFlags

	failuresCollector      *failures.ResourceFailuresCollector
	parsedObjectsCollector *ObjectsCollector
}

// NewParser produces a new Parser object provided a logging mechanism
// and a Kubernetes object store.
func NewParser(
	logger logr.Logger,
	storer store.Storer,
	featureFlags FeatureFlags,
) (*Parser, error) {
	failuresCollector := failures.NewResourceFailuresCollector(logger)

	// If the feature flag is enabled, create a new collector for parsed objects.
	var parsedObjectsCollector *ObjectsCollector
	if featureFlags.ReportConfiguredKubernetesObjects {
		parsedObjectsCollector = NewObjectsCollector()
	}

	return &Parser{
		logger:                 logger,
		storer:                 storer,
		featureFlags:           featureFlags,
		failuresCollector:      failuresCollector,
		parsedObjectsCollector: parsedObjectsCollector,
	}, nil
}

// -----------------------------------------------------------------------------
// Parser - Public Methods
// -----------------------------------------------------------------------------

// KongConfigBuildingResult is a result of Parser.BuildKongConfig method.
type KongConfigBuildingResult struct {
	// KongState is the Kong configuration used to configure the Gateway(s).
	KongState *kongstate.KongState

	// TranslationFailures is a list of resource failures that occurred during parsing.
	// They should be used to provide users with feedback on Kubernetes objects validity.
	TranslationFailures []failures.ResourceFailure

	// ConfiguredKubernetesObjects is a list of Kubernetes objects that were successfully parsed.
	ConfiguredKubernetesObjects []client.Object
}

// BuildKongConfig creates a Kong configuration from Ingress and Custom resources
// defined in Kubernetes.
func (p *Parser) BuildKongConfig() KongConfigBuildingResult {
	// parse and merge all rules together from all Kubernetes API sources
	ingressRules := mergeIngressRules(
		p.ingressRulesFromIngressV1(),
		p.ingressRulesFromTCPIngressV1beta1(),
		p.ingressRulesFromUDPIngressV1beta1(),
		p.ingressRulesFromHTTPRoutes(),
		p.ingressRulesFromUDPRoutes(),
		p.ingressRulesFromTCPRoutes(),
		p.ingressRulesFromTLSRoutes(),
		p.ingressRulesFromGRPCRoutes(),
	)

	// populate any Kubernetes Service objects relevant objects and get the
	// services to be skipped because of annotations inconsistency
	servicesToBeSkipped := ingressRules.populateServices(p.logger, p.storer, p.failuresCollector)

	// add the routes and services to the state
	var result kongstate.KongState

	// generate Upstreams and Targets from service defs
	// update ServiceNameToServices with resolved ports (translating any name references to their number, as Kong
	// services require a number)
	result.Upstreams, ingressRules.ServiceNameToServices = p.getUpstreams(ingressRules.ServiceNameToServices)

	for key, service := range ingressRules.ServiceNameToServices {
		// if the service doesn't need to be skipped, then add it to the
		// list of services.
		if _, ok := servicesToBeSkipped[key]; !ok {
			result.Services = append(result.Services, service)
		}
	}

	// merge KongIngress with Routes, Services and Upstream
	result.FillOverrides(p.logger, p.storer)

	// generate consumers and credentials
	result.FillConsumersAndCredentials(p.logger, p.storer, p.failuresCollector)
	for i := range result.Consumers {
		p.registerSuccessfullyParsedObject(&result.Consumers[i].K8sKongConsumer)
	}

	// process consumer groups
	result.FillConsumerGroups(p.logger, p.storer)
	for i := range result.ConsumerGroups {
		p.registerSuccessfullyParsedObject(&result.ConsumerGroups[i].K8sKongConsumerGroup)
	}

	// process annotation plugins
	result.FillPlugins(p.logger, p.storer, p.failuresCollector)
	for i := range result.Plugins {
		p.registerSuccessfullyParsedObject(result.Plugins[i].K8sParent)
	}

	// generate Certificates and SNIs
	ingressCerts := p.getCerts(ingressRules.SecretNameToSNIs)
	gatewayCerts := p.getGatewayCerts()
	// note that ingress-derived certificates will take precedence over gateway-derived certificates for SNI assignment
	result.Certificates = mergeCerts(p.logger, ingressCerts, gatewayCerts)

	// populate CA certificates in Kong
	result.CACertificates = p.getCACerts()

	if p.licenseGetter != nil {
		optionalLicense := p.licenseGetter.GetLicense()
		if l, ok := optionalLicense.Get(); ok {
			result.Licenses = append(result.Licenses, kongstate.License{License: l})
		}
	}

	if p.featureFlags.FillIDs {
		// generate IDs for Kong entities
		result.FillIDs(p.logger)
	}

	return KongConfigBuildingResult{
		KongState:                   &result,
		TranslationFailures:         p.popTranslationFailures(),
		ConfiguredKubernetesObjects: p.popConfiguredKubernetesObjects(),
	}
}

// -----------------------------------------------------------------------------
// Parser - Public Methods - Other Optional Features
// -----------------------------------------------------------------------------

// InjectLicenseGetter sets a license getter to be used by the parser.
func (p *Parser) InjectLicenseGetter(licenseGetter LicenseGetter) {
	p.licenseGetter = licenseGetter
}

// -----------------------------------------------------------------------------
// Parser - Private Methods
// -----------------------------------------------------------------------------

// registerTranslationFailure should be called when any Kubernetes object translation failure is encountered.
func (p *Parser) registerTranslationFailure(reason string, causingObjects ...client.Object) {
	p.failuresCollector.PushResourceFailure(reason, causingObjects...)
}

func (p *Parser) popTranslationFailures() []failures.ResourceFailure {
	return p.failuresCollector.PopResourceFailures()
}

// registerSuccessfullyParsedObject should be called when any Kubernetes object is successfully parsed.
// It collects the object for reporting purposes.
func (p *Parser) registerSuccessfullyParsedObject(obj client.Object) {
	p.parsedObjectsCollector.Add(obj)
}

// popConfiguredKubernetesObjects provides a list of all the Kubernetes objects
// that have been successfully parsed as part of BuildKongConfig() call so far.
func (p *Parser) popConfiguredKubernetesObjects() []client.Object {
	return p.parsedObjectsCollector.Pop()
}

func tcpIngressToNetworkingTLS(tls []kongv1beta1.IngressTLS) []netv1.IngressTLS {
	var result []netv1.IngressTLS

	for _, t := range tls {
		result = append(result, netv1.IngressTLS{
			Hosts:      t.Hosts,
			SecretName: t.SecretName,
		})
	}
	return result
}

// findPort finds a port matching the specified definition in a Kubernetes Service.
func findPort(svc *corev1.Service, wantPort kongstate.PortDef) (*corev1.ServicePort, error) {
	switch wantPort.Mode {
	case kongstate.PortModeByNumber:
		// ExternalName Services have no port declaration of their own
		// We must assume that the user-requested port is valid and construct a ServicePort from it
		if svc.Spec.Type == corev1.ServiceTypeExternalName {
			return &corev1.ServicePort{
				Port:       wantPort.Number,
				TargetPort: intstr.FromInt(int(wantPort.Number)),
			}, nil
		}
		for _, port := range svc.Spec.Ports {
			port := port
			if port.Port == wantPort.Number {
				return &port, nil
			}
		}

	case kongstate.PortModeByName:
		if svc.Spec.Type == corev1.ServiceTypeExternalName {
			return nil, fmt.Errorf("rules with an ExternalName service must specify numeric ports")
		}
		for _, port := range svc.Spec.Ports {
			port := port
			if port.Name == wantPort.Name {
				return &port, nil
			}
			if port.TargetPort.Type == intstr.String && port.TargetPort.String() == wantPort.Name {
				return &port, nil
			}
		}

	case kongstate.PortModeImplicit:
		if svc.Spec.Type == corev1.ServiceTypeExternalName {
			return nil, fmt.Errorf("rules with an ExternalName service must specify numeric ports")
		}
		if len(svc.Spec.Ports) != 1 {
			return nil, fmt.Errorf("in implicit mode, service must have exactly 1 port, has %d", len(svc.Spec.Ports))
		}
		return &svc.Spec.Ports[0], nil

	default:
		return nil, fmt.Errorf("unknown mode %v", wantPort.Mode)
	}

	return nil, fmt.Errorf("no suitable port found")
}

func (p *Parser) getUpstreams(serviceMap map[string]kongstate.Service) ([]kongstate.Upstream, map[string]kongstate.Service) {
	upstreamDedup := make(map[string]struct{}, len(serviceMap))
	var empty struct{}
	upstreams := make([]kongstate.Upstream, 0, len(serviceMap))
	for serviceName, service := range serviceMap {
		// the name of the Upstream for a service must match the service.Host
		// as the Gateway's internal DNS resolve mechanisms will fail to properly
		// resolve the host otherwise.
		name := *service.Host

		if _, exists := upstreamDedup[name]; !exists {
			// populate all the kong targets for the upstream given all the backends
			var targets []kongstate.Target
			for _, backend := range service.Backends {
				// gather the Kubernetes service for the backend
				backendNamespace := backend.Namespace
				if backendNamespace == "" {
					// if the backend namespace isn't specified, it's in the same namespace as the referee route (which is,
					// somewhat confusingly, the _service_ namespace in serviceMap services, as historically there was no option
					// to reference services outside the route namespace, and we could always stuff the route namespace into the
					// placeholder service.
					backendNamespace = service.Namespace
				}
				k8sService, ok := service.K8sServices[fmt.Sprintf("%s/%s", backendNamespace, backend.Name)]
				if !ok {
					p.registerTranslationFailure(
						fmt.Sprintf("can't add target for backend %s: no kubernetes service found", backend.Name),
						service.Parent,
					)
					continue
				}

				// determine the port for the backend
				port, err := findPort(k8sService, backend.PortDef)
				if err != nil {
					p.registerTranslationFailure(
						fmt.Sprintf("can't find port for backend kubernetes service: %v", err),
						k8sService, service.Parent,
					)
					continue
				}
				service.Port = lo.ToPtr(int(port.Port))
				serviceMap[serviceName] = service

				// get the new targets for this backend service
				newTargets := getServiceEndpoints(p.logger, p.storer, k8sService, port)

				if len(newTargets) == 0 {
					p.logger.V(util.InfoLevel).Info("no targets could be found for kubernetes service",
						"namespace", k8sService.Namespace, "name", k8sService.Name, "kong_service", *service.Name)
				}

				// if weights were set for the backend then that weight needs to be
				// distributed equally among all the targets.
				if backend.Weight != nil && len(newTargets) != 0 {
					// initialize the weight of the target based on the weight of the backend
					// which governs that target (and potentially more). If the weight of the
					// backend is 0 then this indicates an intention to drop all targets from
					// this backend from the load-balancer and is a special situation where
					// all derived targets will receive a weight of 0.
					targetWeight := int(*backend.Weight)

					// if the backend governing this target is not set to a weight of 0,
					// all targets derived from the backend split the weight, therefore
					// equally splitting the traffic load.
					if *backend.Weight != 0 {
						targetWeight = int(*backend.Weight) / len(newTargets)
						// minimum weight of 1 if weight zero was not specifically set.
						if targetWeight == 0 {
							targetWeight = 1
						}
					}

					for i := range newTargets {
						newTargets[i].Weight = &targetWeight
					}
				}

				// add the new targets to the existing pool of targets for the Upstream.
				targets = append(targets, newTargets...)
			}

			// warn if an upstream was created with 0 targets
			if len(targets) == 0 {
				p.logger.V(util.InfoLevel).Info("no targets found to create upstream", "service_name", *service.Name)
			}

			// define the upstream including all the newly populated targets
			// to load-balance traffic to.
			upstream := kongstate.Upstream{
				Upstream: kong.Upstream{
					Name: kong.String(name),
					Tags: service.Tags, // populated by populateServices already
				},
				Service: service,
				Targets: targets,
			}
			upstreams = append(upstreams, upstream)
			upstreamDedup[name] = empty
		}
	}
	return upstreams, serviceMap
}

func getCertFromSecret(secret *corev1.Secret) (string, string, error) {
	certData, okcert := secret.Data[corev1.TLSCertKey]
	keyData, okkey := secret.Data[corev1.TLSPrivateKeyKey]

	if !okcert || !okkey {
		return "", "", fmt.Errorf("no keypair could be found in"+
			" secret '%v/%v'", secret.Namespace, secret.Name)
	}

	cert := bytes.TrimSpace(certData)
	key := bytes.TrimSpace(keyData)

	if _, err := tls.X509KeyPair(cert, key); err != nil {
		return "", "", fmt.Errorf("parsing TLS key-pair in secret '%v/%v': %w",
			secret.Namespace, secret.Name, err)
	}

	return string(cert), string(key), nil
}

type certWrapper struct {
	identifier        string
	cert              kong.Certificate
	snis              []string
	CreationTimestamp metav1.Time
}

func (p *Parser) getGatewayCerts() []certWrapper {
	logger := p.logger
	s := p.storer
	certs := []certWrapper{}
	gateways, err := s.ListGateways()
	if err != nil {
		logger.Error(err, "failed to list Gateways")
		return certs
	}
	for _, gateway := range gateways {
		statuses := make(map[gatewayapi.SectionName]gatewayapi.ListenerStatus, len(gateway.Status.Listeners))
		for _, status := range gateway.Status.Listeners {
			statuses[status.Name] = status
		}

		for _, listener := range gateway.Spec.Listeners {
			status, ok := statuses[listener.Name]
			if !ok {
				logger.V(util.DebugLevel).Info("listener missing status information",
					"gateway", gateway.Name,
					"listener", listener.Name,
					"listener_protocol", listener.Protocol,
					"listener_port", listener.Port,
				)
				continue
			}

			// Check if listener is marked as programmed
			if !util.CheckCondition(
				status.Conditions,
				util.ConditionType(gatewayapi.ListenerConditionProgrammed),
				util.ConditionReason(gatewayapi.ListenerReasonProgrammed),
				metav1.ConditionTrue,
				gateway.Generation,
			) {
				continue
			}

			if listener.TLS != nil {
				if len(listener.TLS.CertificateRefs) > 0 {
					if len(listener.TLS.CertificateRefs) > 1 {
						// TODO support cert_alt and key_alt if there are 2 SecretObjectReferences
						// https://github.com/Kong/kubernetes-ingress-controller/issues/2604
						p.registerTranslationFailure("listener '%s' has more than one certificateRef, it's not supported", gateway)
						continue
					}

					// determine the Secret Namespace
					ref := listener.TLS.CertificateRefs[0]
					namespace := gateway.Namespace
					if ref.Namespace != nil {
						namespace = string(*ref.Namespace)
					}

					// retrieve the Secret and extract the PEM strings
					secret, err := s.GetSecret(namespace, string(ref.Name))
					if err != nil {
						logger.Error(err, "failed to fetch secret",
							"gateway", gateway.Name,
							"listener", listener.Name,
							"secret_name", string(ref.Name),
							"secret_namespace", namespace,
						)
						continue
					}
					cert, key, err := getCertFromSecret(secret)
					if err != nil {
						p.registerTranslationFailure("failed to construct certificate from secret", secret, gateway)
						continue
					}

					// determine the SNI
					hostname := "*"
					if listener.Hostname != nil {
						hostname = string(*listener.Hostname)
					}

					// create a Kong certificate, wrap it in metadata, and add it to the certs slice
					certs = append(certs, certWrapper{
						identifier: cert + key,
						cert: kong.Certificate{
							ID:   kong.String(string(secret.UID)),
							Cert: kong.String(cert),
							Key:  kong.String(key),
							Tags: util.GenerateTagsForObject(secret),
						},
						CreationTimestamp: secret.CreationTimestamp,
						snis:              []string{hostname},
					})
				}
			}
		}
	}
	return certs
}

func (p *Parser) getCerts(secretsToSNIs SecretNameToSNIs) []certWrapper {
	certs := []certWrapper{}

	for secretKey, SNIs := range secretsToSNIs.secretToSNIs {
		namespaceName := strings.Split(secretKey, "/")
		secret, err := p.storer.GetSecret(namespaceName[0], namespaceName[1])
		if err != nil {
			p.registerTranslationFailure(fmt.Sprintf("failed to fetch the secret (%s)", secretKey), SNIs.Parents()...)
			continue
		}
		cert, key, err := getCertFromSecret(secret)
		if err != nil {
			causingObjects := append(SNIs.Parents(), secret)
			p.registerTranslationFailure("failed to construct certificate from secret", causingObjects...)
			continue
		}
		certs = append(certs, certWrapper{
			identifier: cert + key,
			cert: kong.Certificate{
				ID:   kong.String(string(secret.UID)),
				Cert: kong.String(cert),
				Key:  kong.String(key),
				Tags: util.GenerateTagsForObject(secret),
			},
			CreationTimestamp: secret.CreationTimestamp,
			snis:              SNIs.Hosts(),
		})
	}

	return certs
}

func mergeCerts(logger logr.Logger, certLists ...[]certWrapper) []kongstate.Certificate {
	snisSeen := make(map[string]string)
	certsSeen := make(map[string]certWrapper)
	for _, cl := range certLists {
		for _, cw := range cl {
			current, ok := certsSeen[cw.identifier]
			if !ok {
				current = cw
			} else {
				// multiple Secrets that contain identical certificates are collapsed, because we only create one
				// Kong resource for a given cert+key pair. however, because we reuse the Secret ID and creation time
				// for the Kong resource equivalents, the selection of those needs to be deterministic to avoid
				// pointless configuration updates
				if current.CreationTimestamp.After(cw.CreationTimestamp.Time) {
					current.cert.ID = cw.cert.ID
					current.CreationTimestamp = cw.CreationTimestamp
				} else if current.CreationTimestamp.Time.Equal(cw.CreationTimestamp.Time) && (current.cert.ID == nil || *current.cert.ID > *cw.cert.ID) {
					current.cert.ID = cw.cert.ID
					current.CreationTimestamp = cw.CreationTimestamp
				}
				current.snis = append(current.snis, cw.snis...)
			}

			// although we use current in the end, we only warn/exclude on new ones here. SNIs already in the slice
			// have already been vetted by some previous iteration and /are/ in the seen list, but they're in the seen
			// list because the current we retrieved from certsSeen added them
			for _, sni := range cw.snis {
				if seen, ok := snisSeen[sni]; !ok {
					snisSeen[sni] = *current.cert.ID
					current.cert.SNIs = append(current.cert.SNIs, kong.String(sni))
				} else {
					// TODO this should really log information about the requesting Listener or Ingress-like, which is
					// what binds the SNI to a given Secret. Knowing the Secret ID isn't of great use beyond knowing
					// what cert will be served. however, the secretToSNIs input to getCerts does not provide this info
					// https://github.com/Kong/kubernetes-ingress-controller/issues/2605
					logger.Error(nil, "same SNI requested for multiple certs, can only serve one cert",
						"served_secret_cert", seen,
						"requested_secret_cert", *current.cert.ID,
						"sni", sni)
				}
			}
			certsSeen[current.identifier] = current
		}
	}
	var res []kongstate.Certificate
	for _, cw := range certsSeen {
		sort.SliceStable(cw.cert.SNIs, func(i, j int) bool {
			return strings.Compare(*cw.cert.SNIs[i], *cw.cert.SNIs[j]) < 0
		})
		res = append(res, kongstate.Certificate{Certificate: cw.cert})
	}
	return res
}

func getServiceEndpoints(
	logger logr.Logger,
	s store.Storer,
	svc *corev1.Service,
	servicePort *corev1.ServicePort,
) []kongstate.Target {
	logger = logger.WithValues(
		"service_name", svc.Name,
		"service_namespace", svc.Namespace,
		"service_port", servicePort,
	)

	// In theory a Service could have multiple port protocols, we need to ensure we gather
	// endpoints based on all the protocols the service is configured for. We always check
	// for TCP as this is the default protocol for service ports.
	protocols := listProtocols(svc)

	// Check if the service is an upstream service through Ingress Class parameters.
	var isSvcUpstream bool
	ingressClassParameters, err := getIngressClassParametersOrDefault(s)
	if err != nil {
		logger.V(util.DebugLevel).Info("unable to retrieve IngressClassParameters", "error", err)
	} else {
		isSvcUpstream = ingressClassParameters.ServiceUpstream
	}

	// Check all protocols for associated endpoints.
	endpoints := []util.Endpoint{}
	for protocol := range protocols {
		newEndpoints := getEndpoints(logger, svc, servicePort, protocol, s.GetEndpointSlicesForService, isSvcUpstream)
		endpoints = append(endpoints, newEndpoints...)
	}
	if len(endpoints) == 0 {
		logger.V(util.DebugLevel).Info("no active endpoints")
	}

	return targetsForEndpoints(endpoints)
}

// getIngressClassParametersOrDefault returns the parameters for the current ingress class.
// If the cluster operators have specified a set of parameters explicitly, it returns those.
// Otherwise, it returns a default set of parameters.
func getIngressClassParametersOrDefault(s store.Storer) (kongv1alpha1.IngressClassParametersSpec, error) {
	ingressClassName := s.GetIngressClassName()
	ingressClass, err := s.GetIngressClassV1(ingressClassName)
	if err != nil {
		return kongv1alpha1.IngressClassParametersSpec{}, err
	}

	params, err := s.GetIngressClassParametersV1Alpha1(ingressClass)
	if err != nil {
		return kongv1alpha1.IngressClassParametersSpec{}, err
	}

	return params.Spec, nil
}

// getEndpoints returns a list of <endpoint ip>:<port> for a given service/target port combination.
// It also checks if the service is an upstream service either by its annotations
// of by IngressClassParameters configuration provided as a flag.
func getEndpoints(
	logger logr.Logger,
	service *corev1.Service,
	port *corev1.ServicePort,
	proto corev1.Protocol,
	getEndpointSlices func(string, string) ([]*discoveryv1.EndpointSlice, error),
	isSvcUpstream bool,
) []util.Endpoint {
	if service == nil || port == nil {
		return []util.Endpoint{}
	}

	// If service is an upstream service...
	if isSvcUpstream || annotations.HasServiceUpstreamAnnotation(service.Annotations) {
		// ... return its address as the only endpoint.
		return []util.Endpoint{
			{
				Address: service.Name + "." + service.Namespace + ".svc",
				Port:    fmt.Sprint(port.Port),
			},
		}
	}

	logger = logger.WithValues(
		"service_name", service.Name,
		"service_namespace", service.Namespace,
		"service_port", port.String(),
	)

	// ExternalName services
	if service.Spec.Type == corev1.ServiceTypeExternalName {
		logger.V(util.DebugLevel).Info("found service of type=ExternalName")
		return []util.Endpoint{
			{
				Address: service.Spec.ExternalName,
				Port:    port.TargetPort.String(),
			},
		}
	}

	logger.V(util.DebugLevel).Info("fetching EndpointSlices")
	endpointSlices, err := getEndpointSlices(service.Namespace, service.Name)
	if err != nil {
		logger.Error(err, "error fetching EndpointSlices")
		return []util.Endpoint{}
	}
	logger.V(util.DebugLevel).Info("fetched EndpointSlices", "count", len(endpointSlices))

	// Avoid duplicated upstream servers when the service contains
	// multiple port definitions sharing the same target port.
	uniqueUpstream := make(map[util.Endpoint]struct{})
	upstreamServers := make([]util.Endpoint, 0)
	for _, endpointSlice := range endpointSlices {
		for _, p := range endpointSlice.Ports {
			if p.Port == nil || *p.Port < 0 || *p.Protocol != proto || *p.Name != port.Name {
				continue
			}
			upstreamPort := fmt.Sprint(*p.Port)
			for _, endpoint := range endpointSlice.Endpoints {
				// Ready indicates that this endpoint is prepared to receive traffic, according to whatever
				// system is managing the endpoint. A nil value indicates an unknown state.
				// In most cases consumers should interpret this unknown state as ready.
				// Field Ready has the same semantic as Endpoints from CoreV1 in Addresses.
				// https://kubernetes.io/docs/concepts/services-networking/endpoint-slices/#conditions
				if endpoint.Conditions.Ready != nil && !*endpoint.Conditions.Ready {
					continue
				}
				// One address per endpoint is rather expected (allowing multiple is due to historical reasons)
				// read more https://github.com/kubernetes/kubernetes/issues/106267#issuecomment-978770401.
				// These are all assumed to be fungible and clients may choose to only use the first element.
				upstreamServer := util.Endpoint{
					Address: endpoint.Addresses[0],
					Port:    upstreamPort,
				}
				if _, exists := uniqueUpstream[upstreamServer]; !exists {
					upstreamServers = append(upstreamServers, upstreamServer)
					uniqueUpstream[upstreamServer] = struct{}{}
				}
			}
		}
	}
	logger.V(util.DebugLevel).Info("found endpoints", "endpoints", upstreamServers)
	return upstreamServers
}

// listProtocols is a helper function to map out all the in-use corev1.Protocols
// for a service given a corev1.Service object.
//
// TODO: due to historical logic this function defaults to assuming TCP protocol
// is valid for the Service and its endpoints, however we need to follow up
// on this as this is not technically correct and causes waste.
// See: https://github.com/Kong/kubernetes-ingress-controller/issues/1429
func listProtocols(svc *corev1.Service) map[corev1.Protocol]bool {
	protocols := map[corev1.Protocol]bool{corev1.ProtocolTCP: true}
	for _, port := range svc.Spec.Ports {
		if port.Protocol != "" {
			protocols[port.Protocol] = true
		}
	}
	return protocols
}

// targetsForEndpoints generates kongstate.Target objects for each util.Endpoint provided.
func targetsForEndpoints(endpoints []util.Endpoint) []kongstate.Target {
	targets := []kongstate.Target{}
	for _, endpoint := range endpoints {
		addr := endpoint.Address
		parsed := net.ParseIP(endpoint.Address)
		if parsed != nil {
			if parsed.To4() == nil {
				// If we have an IPv6 endpoint, we need to surround it with brackets, else the port concat after this will
				// treat the port as part of the address.
				addr = fmt.Sprintf("[%s]", endpoint.Address)
			}
		}
		target := kongstate.Target{
			Target: kong.Target{
				Target: kong.String(addr + ":" + endpoint.Port),
			},
		}
		targets = append(targets, target)
	}
	return targets
}
