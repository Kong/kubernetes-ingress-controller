package parser

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	netv1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	knative "knative.dev/networking/pkg/apis/networking/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	configurationv1alpha1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1alpha1"
	configurationv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
)

// -----------------------------------------------------------------------------
// Parser - Public Types
// -----------------------------------------------------------------------------

// Parser parses Kubernetes objects and configurations into their
// equivalent Kong objects and configurations, producing a complete
// state configuration for the Kong Admin API.
type Parser struct {
	logger                      logrus.FieldLogger
	storer                      store.Storer
	configuredKubernetesObjects []client.Object

	featureEnabledReportConfiguredKubernetesObjects bool
	featureEnabledCombinedServiceRoutes             bool

	flagEnabledRegexPathPrefix bool
}

// NewParser produces a new Parser object provided a logging mechanism
// and a Kubernetes object store.
func NewParser(
	logger logrus.FieldLogger,
	storer store.Storer,
) *Parser {
	return &Parser{
		logger: logger,
		storer: storer,
	}
}

// -----------------------------------------------------------------------------
// Parser - Public Methods
// -----------------------------------------------------------------------------

// Build creates a Kong configuration from Ingress and Custom resources
// defined in Kubernetes.
// It throws an error if there is an error returned from client-go.
func (p *Parser) Build() (*kongstate.KongState, error) {
	// parse and merge all rules together from all Kubernetes API sources
	ingressRules := mergeIngressRules(
		p.ingressRulesFromIngressV1beta1(),
		p.ingressRulesFromIngressV1(),
		p.ingressRulesFromTCPIngressV1beta1(),
		p.ingressRulesFromUDPIngressV1beta1(),
		p.ingressRulesFromKnativeIngress(),
		p.ingressRulesFromHTTPRoutes(),
		p.ingressRulesFromUDPRoutes(),
		p.ingressRulesFromTCPRoutes(),
		p.ingressRulesFromTLSRoutes(),
	)

	// populate any Kubernetes Service objects relevant objects and get the
	// services to be skipped because of annotations inconsistency
	servicesToBeSkipped := ingressRules.populateServices(p.logger, p.storer)

	// add the routes and services to the state
	var result kongstate.KongState
	for key, service := range ingressRules.ServiceNameToServices {
		// if the service doesn't need to be skipped, then add it to the
		// list of services.
		if _, ok := servicesToBeSkipped[key]; !ok {
			result.Services = append(result.Services, service)
		}
	}

	// generate Upstreams and Targets from service defs
	result.Upstreams = getUpstreams(p.logger, p.storer, ingressRules.ServiceNameToServices)

	// merge KongIngress with Routes, Services and Upstream
	result.FillOverrides(p.logger, p.storer)

	// generate consumers and credentials
	result.FillConsumersAndCredentials(p.logger, p.storer)

	// process annotation plugins
	result.FillPlugins(p.logger, p.storer)

	// generate Certificates and SNIs
	ingressCerts := getCerts(p.logger, p.storer, ingressRules.SecretNameToSNIs)
	gatewayCerts := getGatewayCerts(p.logger, p.storer)
	// note that ingress-derived certificates will take precedence over gateway-derived certificates for SNI assignment
	result.Certificates = mergeCerts(p.logger, ingressCerts, gatewayCerts)

	// populate CA certificates in Kong
	var err error
	caCertSecrets, err := p.storer.ListCACerts()
	if err != nil {
		return nil, err
	}
	result.CACertificates = toCACerts(p.logger, caCertSecrets)

	return &result, nil
}

// -----------------------------------------------------------------------------
// Parser - Public Methods - Kubernetes Object Reporting
// -----------------------------------------------------------------------------

// EnableKubernetesObjectReports turns on object reporting for this parser:
// each subsequent call to Build() will track the Kubernetes objects which
// were successfully parsed. Objects tracked this way can be retrieved by
// calling GenerateKubernetesObjectReport().
func (p *Parser) EnableKubernetesObjectReports() {
	p.featureEnabledReportConfiguredKubernetesObjects = true
}

// ReportKubernetesObjectUpdate reports an update to a Kubernetes object if
// updates have been requested. If the parser has not been configured to
// report Kubernetes object updates this is a no-op.
func (p *Parser) ReportKubernetesObjectUpdate(obj client.Object) {
	if p.featureEnabledReportConfiguredKubernetesObjects {
		p.configuredKubernetesObjects = append(p.configuredKubernetesObjects, obj)
	}
}

// GenerateKubernetesObjectReport provides a list of all the Kubernetes objects
// that have been successfully parsed as part of Build() calls so far. The
// objects are consumed: the parser's internal list will be emptied once this
// method is called, until more builds are run.
func (p *Parser) GenerateKubernetesObjectReport() []client.Object {
	report := p.configuredKubernetesObjects
	p.configuredKubernetesObjects = nil
	return report
}

// -----------------------------------------------------------------------------
// Parser - Public Methods - Other Optional Features
// -----------------------------------------------------------------------------

// EnableCombinedServiceRoutes changes the translation logic from the legacy
// mode which would create a kong.Route object per each individual path on
// an Ingress object to a mode that can combine routes for paths where the
// service name, host and port match for those paths.
func (p *Parser) EnableCombinedServiceRoutes() {
	p.featureEnabledCombinedServiceRoutes = true
}

// EnableRegexPathPrefix enables adding the Kong 3.x+ regex path prefix on regex paths generated by the controller
// (to satisfy the Ingress Prefix and Exact path types) or indicated by a resource (e.g. when an HTTPRoute uses a
// RegularExpression Match). It does _not_ enable heuristic regex path detection for Ingress ImplementationSpecific
// paths, which require an IngressClass setting.
func (p *Parser) EnableRegexPathPrefix() {
	p.flagEnabledRegexPathPrefix = true
}

// -----------------------------------------------------------------------------
// Parser - Private Methods
// -----------------------------------------------------------------------------

func toCACerts(log logrus.FieldLogger, caCertSecrets []*corev1.Secret) []kong.CACertificate {
	var caCerts []kong.CACertificate
	for _, certSecret := range caCertSecrets {
		secretName := certSecret.Namespace + "/" + certSecret.Name

		idbytes, idExists := certSecret.Data["id"]
		log = log.WithFields(logrus.Fields{
			"secret_name":      secretName,
			"secret_namespace": certSecret.Namespace,
		})
		if !idExists {
			log.Errorf("invalid CA certificate: missing 'id' field in data")
			continue
		}

		caCertbytes, certExists := certSecret.Data["cert"]
		if !certExists {
			log.Errorf("invalid CA certificate: missing 'cert' field in data")
			continue
		}

		pemBlock, _ := pem.Decode(caCertbytes)
		if pemBlock == nil {
			log.Errorf("invalid CA certificate: invalid PEM block")
			continue
		}
		x509Cert, err := x509.ParseCertificate(pemBlock.Bytes)
		if err != nil {
			log.WithError(err).Errorf("invalid CA certificate: failed to parse certificate")
			continue
		}
		if !x509Cert.IsCA {
			log.WithError(err).Errorf("invalid CA certificate: certificate is missing the 'CA' basic constraint")
			continue
		}

		caCerts = append(caCerts, kong.CACertificate{
			ID:   kong.String(string(idbytes)),
			Cert: kong.String(string(caCertbytes)),
		})
	}

	return caCerts
}

func knativeIngressToNetworkingTLS(tls []knative.IngressTLS) []netv1beta1.IngressTLS {
	var result []netv1beta1.IngressTLS

	for _, t := range tls {
		result = append(result, netv1beta1.IngressTLS{
			Hosts:      t.Hosts,
			SecretName: t.SecretName,
		})
	}
	return result
}

func tcpIngressToNetworkingTLS(tls []configurationv1beta1.IngressTLS) []netv1beta1.IngressTLS {
	var result []netv1beta1.IngressTLS

	for _, t := range tls {
		result = append(result, netv1beta1.IngressTLS{
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
			if port.Port == wantPort.Number {
				return &port, nil
			}
		}

	case kongstate.PortModeByName:
		if svc.Spec.Type == corev1.ServiceTypeExternalName {
			return nil, fmt.Errorf("rules with an ExternalName service must specify numeric ports")
		}
		for _, port := range svc.Spec.Ports {
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

func getUpstreams(
	log logrus.FieldLogger,
	s store.Storer,
	serviceMap map[string]kongstate.Service,
) []kongstate.Upstream {
	upstreamDedup := make(map[string]struct{}, len(serviceMap))
	var empty struct{}
	upstreams := make([]kongstate.Upstream, 0, len(serviceMap))
	for _, service := range serviceMap {
		// the name of the Upstream for a service must match the service.Host
		// as the Gateway's internal DNS resolve mechanisms will fail to properly
		// resolve the host otherwise.
		name := *service.Host

		if _, exists := upstreamDedup[name]; !exists {
			// populate all the kong targets for the upstream given all the backends
			var targets []kongstate.Target
			for _, backend := range service.Backends {
				// gather the Kubernetes service for the backend
				k8sService, ok := service.K8sServices[backend.Name]
				if !ok {
					log.WithField("service_name", *service.Name).Errorf("can't add target for backend %s: no kubernetes service found", backend.Name)
					continue
				}

				// determine the port for the backend
				port, err := findPort(k8sService, backend.PortDef)
				if err != nil {
					log.WithField("service_name", *service.Name).Errorf("can't find port for backend kubernetes service %s/%s: %v", k8sService.Namespace, k8sService.Name, err)
					continue
				}

				// get the new targets for this backend service
				newTargets := getServiceEndpoints(log, s, k8sService, port)

				if len(newTargets) == 0 {
					log.WithField("service_name", *service.Name).Infof("no targets could be found for kubernetes service %s/%s", k8sService.Namespace, k8sService.Name)
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
				log.WithField("service_name", *service.Name).Infof("no targets found to create upstream")
			}

			// define the upstream including all the newly populated targets
			// to load-balance traffic to.
			upstream := kongstate.Upstream{
				Upstream: kong.Upstream{
					Name: kong.String(name),
				},
				Service: service,
				Targets: targets,
			}
			upstreams = append(upstreams, upstream)
			upstreamDedup[name] = empty
		}
	}
	return upstreams
}

func getCertFromSecret(secret *corev1.Secret) (string, string, error) {
	certData, okcert := secret.Data[corev1.TLSCertKey]
	keyData, okkey := secret.Data[corev1.TLSPrivateKeyKey]

	if !okcert || !okkey {
		return "", "", fmt.Errorf("no keypair could be found in"+
			" secret '%v/%v'", secret.Namespace, secret.Name)
	}

	cert := strings.TrimSpace(bytes.NewBuffer(certData).String())
	key := strings.TrimSpace(bytes.NewBuffer(keyData).String())

	_, err := tls.X509KeyPair([]byte(cert), []byte(key))
	if err != nil {
		return "", "", fmt.Errorf("parsing TLS key-pair in secret '%v/%v': %w",
			secret.Namespace, secret.Name, err)
	}

	return cert, key, nil
}

type certWrapper struct {
	identifier        string
	cert              kong.Certificate
	snis              []string
	CreationTimestamp metav1.Time
}

func getGatewayCerts(log logrus.FieldLogger, s store.Storer) []certWrapper {
	certs := []certWrapper{}
	gateways, err := s.ListGateways()
	if err != nil {
		log.WithError(err).Error("failed to list Gateways")
		return certs
	}
	grants, err := s.ListReferenceGrants()
	if err != nil {
		log.WithError(err).Error("failed to list ReferenceGrants")
		return certs
	}
	for _, gateway := range gateways {
		statuses := make(map[gatewayv1beta1.SectionName]gatewayv1beta1.ListenerStatus, len(gateway.Status.Listeners))
		for _, status := range gateway.Status.Listeners {
			statuses[status.Name] = status
		}
		for _, listener := range gateway.Spec.Listeners {
			ready := false
			if status, ok := statuses[listener.Name]; ok {
				log.WithFields(logrus.Fields{
					"gateway":  gateway.Name,
					"listener": listener.Name,
				}).Debug("listener missing status information")
				if ok := util.CheckCondition(
					status.Conditions,
					util.ConditionType(gatewayv1alpha2.ListenerConditionReady),
					util.ConditionReason(gatewayv1alpha2.ListenerReasonReady),
					metav1.ConditionTrue,
					gateway.Generation,
				); ok {
					ready = true
				}
			}
			if !ready {
				continue
			}
			if listener.TLS != nil {
				if len(listener.TLS.CertificateRefs) > 0 {
					if len(listener.TLS.CertificateRefs) > 1 {
						// TODO support cert_alt and key_alt if there are 2 SecretObjectReferences
						// https://github.com/Kong/kubernetes-ingress-controller/issues/2604
						log.WithFields(logrus.Fields{
							"gateway":  gateway.Name,
							"listener": listener.Name,
						}).Error("Gateway Listeners with more than one certificateRef are not supported")
						continue
					}

					ref := listener.TLS.CertificateRefs[0]

					// SecretObjectReference is a misnomer; it can reference any Group+Kind, but defaults to Secret
					if ref.Kind != nil {
						if string(*ref.Kind) != "Secret" {
							log.WithFields(logrus.Fields{
								"gateway":  gateway.Name,
								"listener": listener.Name,
							}).Error("CertificateRefs to kinds other than Secret are not supported")
						}
					}

					// determine the Secret Namespace and validate against ReferenceGrant if needed
					namespace := gateway.Namespace
					if ref.Namespace != nil {
						namespace = string(*ref.Namespace)
					}
					if namespace != gateway.Namespace {
						allowed := getPermittedForReferenceGrantFrom(gatewayv1alpha2.ReferenceGrantFrom{
							Group:     gatewayv1alpha2.Group(gateway.GetObjectKind().GroupVersionKind().Group),
							Kind:      gatewayv1alpha2.Kind(gateway.GetObjectKind().GroupVersionKind().Kind),
							Namespace: gatewayv1alpha2.Namespace(gateway.GetNamespace()),
						}, grants)

						if !newRefChecker(ref).IsRefAllowedByGrant(allowed) {
							log.WithFields(logrus.Fields{
								"gateway":           gateway.Name,
								"gateway_namespace": gateway.Namespace,
								"listener":          listener.Name,
								"secret_name":       string(ref.Name),
								"secret_namespace":  namespace,
							}).WithError(err).Error("secret reference not allowed by ReferenceGrant")
							continue
						}
					}

					// retrieve the Secret and extract the PEM strings
					secret, err := s.GetSecret(namespace, string(ref.Name))
					if err != nil {
						log.WithFields(logrus.Fields{
							"gateway":          gateway.Name,
							"listener":         listener.Name,
							"secret_name":      string(ref.Name),
							"secret_namespace": namespace,
						}).WithError(err).Error("failed to fetch secret")
						continue
					}
					cert, key, err := getCertFromSecret(secret)
					if err != nil {
						log.WithFields(logrus.Fields{
							"gateway":          gateway.Name,
							"listener":         listener.Name,
							"secret_name":      string(ref.Name),
							"secret_namespace": namespace,
						}).WithError(err).Error("failed to construct certificate from secret")
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

func getCerts(log logrus.FieldLogger, s store.Storer, secretsToSNIs map[string][]string) []certWrapper {
	certs := []certWrapper{}

	for secretKey, SNIs := range secretsToSNIs {
		namespaceName := strings.Split(secretKey, "/")
		secret, err := s.GetSecret(namespaceName[0], namespaceName[1])
		if err != nil {
			log.WithFields(logrus.Fields{
				"secret_name":      namespaceName[1],
				"secret_namespace": namespaceName[0],
			}).WithError(err).Error("failed to fetch secret")
			continue
		}
		cert, key, err := getCertFromSecret(secret)
		if err != nil {
			log.WithFields(logrus.Fields{
				"secret_name":      namespaceName[1],
				"secret_namespace": namespaceName[0],
			}).WithError(err).Error("failed to construct certificate from secret")
			continue
		}
		certs = append(certs, certWrapper{
			identifier: cert + key,
			cert: kong.Certificate{
				ID:   kong.String(string(secret.UID)),
				Cert: kong.String(cert),
				Key:  kong.String(key),
			},
			CreationTimestamp: secret.CreationTimestamp,
			snis:              SNIs,
		})
	}

	return certs
}

func mergeCerts(log logrus.FieldLogger, certLists ...[]certWrapper) []kongstate.Certificate {
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
					log.WithFields(logrus.Fields{
						"served_secret_cert":    seen,
						"requested_secret_cert": *current.cert.ID,
						"sni":                   sni,
					}).Error("same SNI requested for multiple certs, can only serve one cert")
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
	log logrus.FieldLogger,
	s store.Storer,
	svc *corev1.Service,
	servicePort *corev1.ServicePort,
) []kongstate.Target {
	log = log.WithFields(logrus.Fields{
		"service_name":      svc.Name,
		"service_namespace": svc.Namespace,
		"service_port":      servicePort,
	})

	// in theory a Service could have multiple port protocols, we need to ensure we gather
	// endpoints based on all the protocols the service is configured for. We always check
	// for TCP as this is the default protocol for service ports.
	protocols := listProtocols(svc)

	// Check if the service is an upstream service either by annotation or controller configuration.
	var isSvcUpstream bool
	ingressClassParameters, err := getIngressClassParametersOrDefault(s)
	if err != nil {
		log.Debugf("error getting an IngressClassParameters: %v", err)
	} else {
		isSvcUpstream = ingressClassParameters.ServiceUpstream
	}

	// check all protocols for associated endpoints
	endpoints := []util.Endpoint{}
	for protocol := range protocols {
		newEndpoints := getEndpoints(log, svc, servicePort, protocol, s.GetEndpointsForService, isSvcUpstream)
		if len(newEndpoints) > 0 {
			endpoints = append(endpoints, newEndpoints...)
		}
	}
	if len(endpoints) == 0 {
		log.Warningf("no active endpoints")
	}

	return targetsForEndpoints(endpoints)
}

// getIngressClassParametersOrDefault returns the parameters for the current ingress class.
// If the cluster operators have specified a set of parameters explicitly, it returns those.
// Otherwise, it returns a default set of parameters.
func getIngressClassParametersOrDefault(s store.Storer) (configurationv1alpha1.IngressClassParametersSpec, error) {
	params, err := s.GetIngressClassParametersV1Alpha1()
	if err != nil {
		return configurationv1alpha1.IngressClassParametersSpec{}, err
	}

	return params.Spec, nil
}

// getEndpoints returns a list of <endpoint ip>:<port> for a given service/target port combination.
// It also checks if the service is an upstream service either by its annotations
// of by IngressClassParameters configuration provided as a flag.
func getEndpoints(
	log logrus.FieldLogger,
	s *corev1.Service,
	port *corev1.ServicePort,
	proto corev1.Protocol,
	getEndpoints func(string, string) (*corev1.Endpoints, error),
	isSvcUpstream bool,
) []util.Endpoint {
	upsServers := []util.Endpoint{}

	if s == nil || port == nil {
		return upsServers
	}

	// If service is an upstream service...
	if isSvcUpstream || annotations.HasServiceUpstreamAnnotation(s.Annotations) {
		// ... return its address as the only endpoint.
		return append(upsServers, util.Endpoint{
			Address: s.Name + "." + s.Namespace + ".svc",
			Port:    fmt.Sprintf("%v", port.Port),
		})
	}

	log = log.WithFields(logrus.Fields{
		"service_name":      s.Name,
		"service_namespace": s.Namespace,
		"service_port":      port.String(),
	})

	// avoid duplicated upstream servers when the service
	// contains multiple port definitions sharing the same
	// targetport.
	adus := make(map[string]bool)

	// ExternalName services
	if s.Spec.Type == corev1.ServiceTypeExternalName {
		log.Debug("found service of type=ExternalName")

		targetPort := port.TargetPort.IntValue()
		// check for invalid port value
		if targetPort <= 0 {
			err := fmt.Errorf("invalid port: %v", targetPort)
			log.WithError(err).Error("invalid service")
			return upsServers
		}

		return append(upsServers, util.Endpoint{
			Address: s.Spec.ExternalName,
			Port:    fmt.Sprintf("%v", targetPort),
		})
	}
	if annotations.HasServiceUpstreamAnnotation(s.Annotations) {
		return append(upsServers, util.Endpoint{
			Address: s.Name + "." + s.Namespace + ".svc",
			Port:    fmt.Sprintf("%v", port.Port),
		})
	}

	log.Debugf("fetching endpoints")
	ep, err := getEndpoints(s.Namespace, s.Name)
	if err != nil {
		log.WithError(err).Error("failed to fetch endpoints")
		return upsServers
	}

	for _, ss := range ep.Subsets {
		for _, epPort := range ss.Ports {

			if !reflect.DeepEqual(epPort.Protocol, proto) {
				continue
			}

			var targetPort int32

			if port.Name == "" {
				// port.Name is optional if there is only one port
				targetPort = epPort.Port
			} else if port.Name == epPort.Name {
				targetPort = epPort.Port
			}

			// check for invalid port value
			if targetPort <= 0 {
				continue
			}

			for _, epAddress := range ss.Addresses {
				ep := fmt.Sprintf("%v:%v", epAddress.IP, targetPort)
				if _, exists := adus[ep]; exists {
					continue
				}
				ups := util.Endpoint{
					Address: epAddress.IP,
					Port:    fmt.Sprintf("%v", targetPort),
				}
				upsServers = append(upsServers, ups)
				adus[ep] = true
			}
		}
	}

	log.Debugf("found endpoints: %v", upsServers)
	return upsServers
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
		target := kongstate.Target{
			Target: kong.Target{
				Target: kong.String(endpoint.Address + ":" + endpoint.Port),
			},
		}
		targets = append(targets, target)
	}
	return targets
}
