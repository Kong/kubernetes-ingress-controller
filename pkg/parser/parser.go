package parser

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"reflect"
	"strings"

	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/pkg/annotations"
	configurationv1beta1 "github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/pkg/kongstate"
	"github.com/kong/kubernetes-ingress-controller/pkg/store"
	"github.com/kong/kubernetes-ingress-controller/pkg/util"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	knative "knative.dev/networking/pkg/apis/networking/v1alpha1"
)

func parseAll(log logrus.FieldLogger, s store.Storer) ingressRules {
	parsedIngressV1beta1 := fromIngressV1beta1(log, s.ListIngressesV1beta1())
	parsedIngressV1 := fromIngressV1(log, s.ListIngressesV1())

	tcpIngresses, err := s.ListTCPIngresses()
	if err != nil {
		log.Errorf("failed to list TCPIngresses: %v", err)
	}
	parsedTCPIngress := fromTCPIngressV1beta1(log, tcpIngresses)

	knativeIngresses, err := s.ListKnativeIngresses()
	if err != nil {
		log.Errorf("failed to list Knative Ingresses: %v", err)
	}
	parsedKnative := fromKnativeIngress(log, knativeIngresses)

	return mergeIngressRules(parsedIngressV1beta1, parsedIngressV1, parsedTCPIngress, parsedKnative)
}

// Build creates a Kong configuration from Ingress and Custom resources
// defined in Kuberentes.
// It throws an error if there is an error returned from client-go.
func Build(log logrus.FieldLogger, s store.Storer) (*kongstate.KongState, error) {
	parsedAll := parseAll(log, s)
	parsedAll.populateServices(log, s)

	var result kongstate.KongState
	// add the routes and services to the state
	for _, service := range parsedAll.ServiceNameToServices {
		result.Services = append(result.Services, service)
	}

	// generate Upstreams and Targets from service defs
	result.Upstreams = getUpstreams(log, s, parsedAll.ServiceNameToServices)

	// merge KongIngress with Routes, Services and Upstream
	result.FillOverrides(log, s)

	// generate consumers and credentials
	result.FillConsumersAndCredentials(log, s)

	// process annotation plugins
	result.FillPlugins(log, s)

	// generate Certificates and SNIs
	result.Certificates = getCerts(log, s, parsedAll.SecretNameToSNIs)

	// populate CA certificates in Kong
	var err error
	caCertSecrets, err := s.ListCACerts()
	if err != nil {
		return nil, err
	}
	result.CACertificates = toCACerts(log, caCertSecrets)

	return &result, nil
}

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
			log.Errorf("invalid CA certificate: failed to parse certificate: %v", err)
			continue
		}
		if !x509Cert.IsCA {
			log.Errorf("invalid CA certificate: certificate is missing the 'CA' basic constraint: %v", err)
			continue
		}

		caCerts = append(caCerts, kong.CACertificate{
			ID:   kong.String(string(idbytes)),
			Cert: kong.String(string(caCertbytes)),
		})
	}

	return caCerts
}

func knativeIngressToNetworkingTLS(tls []knative.IngressTLS) []networking.IngressTLS {
	var result []networking.IngressTLS

	for _, t := range tls {
		result = append(result, networking.IngressTLS{
			Hosts:      t.Hosts,
			SecretName: t.SecretName,
		})
	}
	return result
}

func tcpIngressToNetworkingTLS(tls []configurationv1beta1.IngressTLS) []networking.IngressTLS {
	var result []networking.IngressTLS

	for _, t := range tls {
		result = append(result, networking.IngressTLS{
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
	log logrus.FieldLogger, s store.Storer, serviceMap map[string]kongstate.Service) []kongstate.Upstream {
	var upstreams []kongstate.Upstream
	for _, service := range serviceMap {

		var targets []kongstate.Target
		port, err := findPort(&service.K8sService, service.Backend.Port)
		if err == nil {
			targets = getServiceEndpoints(log, s, service.K8sService, port)
		} else {
			log.WithField("service_name", *service.Name).Warnf("skipping service - getServiceEndpoints failed: %v", err)
		}

		upstream := kongstate.Upstream{
			Upstream: kong.Upstream{
				Name: kong.String(
					fmt.Sprintf("%s.%s.%s.svc", service.Backend.Name, service.Namespace, service.Backend.Port.CanonicalString())),
			},
			Service: service,
			Targets: targets,
		}
		upstreams = append(upstreams, upstream)
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
		return "", "", fmt.Errorf("parsing TLS key-pair in secret '%v/%v': %v",
			secret.Namespace, secret.Name, err)
	}

	return cert, key, nil
}

func getCerts(log logrus.FieldLogger, s store.Storer, secretsToSNIs map[string][]string) []kongstate.Certificate {
	snisAdded := make(map[string]bool)
	// map of cert public key + private key to certificate
	type certWrapper struct {
		cert              kong.Certificate
		CreationTimestamp metav1.Time
	}
	certs := make(map[string]certWrapper)

	for secretKey, SNIs := range secretsToSNIs {
		namespaceName := strings.Split(secretKey, "/")
		secret, err := s.GetSecret(namespaceName[0], namespaceName[1])
		if err != nil {
			log.WithFields(logrus.Fields{
				"secret_name":      namespaceName[1],
				"secret_namespace": namespaceName[0],
			}).Logger.Errorf("failed to fetch secret: %v", err)
			continue
		}
		cert, key, err := getCertFromSecret(secret)
		if err != nil {
			log.WithFields(logrus.Fields{
				"secret_name":      namespaceName[1],
				"secret_namespace": namespaceName[0],
			}).Logger.Errorf("failed to construct certificate from secret: %v", err)
			continue
		}
		kongCert, ok := certs[cert+key]
		if !ok {
			kongCert = certWrapper{
				cert: kong.Certificate{
					ID:   kong.String(string(secret.UID)),
					Cert: kong.String(cert),
					Key:  kong.String(key),
				},
				CreationTimestamp: secret.CreationTimestamp,
			}
		} else {
			if kongCert.CreationTimestamp.After(secret.CreationTimestamp.Time) {
				kongCert.cert.ID = kong.String(string(secret.UID))
				kongCert.CreationTimestamp = secret.CreationTimestamp
			}
		}

		for _, sni := range SNIs {
			if !snisAdded[sni] {
				snisAdded[sni] = true
				kongCert.cert.SNIs = append(kongCert.cert.SNIs, kong.String(sni))
			}
		}
		certs[cert+key] = kongCert
	}
	var res []kongstate.Certificate
	for _, cert := range certs {
		res = append(res, kongstate.Certificate{Certificate: cert.cert})
	}
	return res
}

func getServiceEndpoints(log logrus.FieldLogger, s store.Storer, svc corev1.Service,
	servicePort *corev1.ServicePort) []kongstate.Target {
	var targets []kongstate.Target
	var endpoints []util.Endpoint

	log = log.WithFields(logrus.Fields{
		"service_name":      svc.Name,
		"service_namespace": svc.Namespace,
		"service_port":      servicePort,
	})

	endpoints = getEndpoints(log, &svc, servicePort, corev1.ProtocolTCP, s.GetEndpointsForService)
	if len(endpoints) == 0 {
		log.Warningf("no active endpionts")
	}
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

// getEndpoints returns a list of <endpoint ip>:<port> for a given service/target port combination.
func getEndpoints(
	log logrus.FieldLogger,
	s *corev1.Service,
	port *corev1.ServicePort,
	proto corev1.Protocol,
	getEndpoints func(string, string) (*corev1.Endpoints, error),
) []util.Endpoint {

	upsServers := []util.Endpoint{}

	if s == nil || port == nil {
		return upsServers
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
			log.Errorf("invalid service: invalid port: %v", targetPort)
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
		log.Errorf("failed to fetch endpoints: %v", err)
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
