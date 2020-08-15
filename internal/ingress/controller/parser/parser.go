package parser

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/annotations"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/store"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/utils"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1"
	configurationv1beta1 "github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1beta1"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/sets"
	"knative.dev/pkg/network"
	knative "knative.dev/serving/pkg/apis/networking/v1alpha1"
)

// Route represents a Kong Route and holds a reference to the Ingress
// rule.
type Route struct {
	kong.Route

	// Ingress object associated with this route
	Ingress networking.Ingress
	// TCPIngress object associated with this route
	TCPIngress configurationv1beta1.TCPIngress
	// Is this route coming from TCPIngress or networking.Ingress?
	IsTCP   bool
	Plugins []kong.Plugin
}

type backend struct {
	Name string
	Port intstr.IntOrString
}

// Service represents a service in Kong and holds routes associated with the
// service and other k8s metadata.
type Service struct {
	kong.Service
	Backend    backend
	Namespace  string
	Routes     []Route
	Plugins    []kong.Plugin
	K8sService corev1.Service
}

// Upstream is a wrapper around Upstream object in Kong.
type Upstream struct {
	kong.Upstream
	Targets []Target
	// Service this upstream is asosciated with.
	Service Service
}

// Target is a wrapper around Target object in Kong.
type Target struct {
	kong.Target
}

// Consumer holds a Kong consumer and its plugins and credentials.
type Consumer struct {
	kong.Consumer
	Plugins    []kong.Plugin
	KeyAuths   []*kong.KeyAuth
	HMACAuths  []*kong.HMACAuth
	JWTAuths   []*kong.JWTAuth
	BasicAuths []*kong.BasicAuth
	ACLGroups  []*kong.ACLGroup

	Oauth2Creds []*kong.Oauth2Credential

	k8sKongConsumer configurationv1.KongConsumer
}

// KongState holds the configuration that should be applied to Kong.
type KongState struct {
	Services       []Service
	Upstreams      []Upstream
	Certificates   []Certificate
	CACertificates []kong.CACertificate
	Plugins        []Plugin
	Consumers      []Consumer
}

// Certificate represents the certificate object in Kong.
type Certificate struct {
	kong.Certificate
}

// Plugin represetns a plugin Object in Kong.
type Plugin struct {
	kong.Plugin
}

// Parser parses Kubernetes CRDs and Ingress rules and generates a
// Kong configuration.
type Parser struct {
	store  store.Storer
	Logger logrus.FieldLogger
}

type parsedIngressRules struct {
	SecretNameToSNIs      map[string][]string
	ServiceNameToServices map[string]Service
}

var supportedCreds = sets.NewString(
	"acl",
	"basic-auth",
	"hmac-auth",
	"jwt",
	"key-auth",
	"oauth2",
)

var validProtocols = regexp.MustCompile(`\Ahttps$|\Ahttp$|\Agrpc$|\Agrpcs|\Atcp|\Atls$`)
var validMethods = regexp.MustCompile(`\A[A-Z]+$`)

// New returns a new parser backed with store.
func New(store store.Storer, logger logrus.FieldLogger) Parser {
	return Parser{store: store, Logger: logger}
}

// Build creates a Kong configuration from Ingress and Custom resources
// defined in Kuberentes.
// It throws an error if there is an error returned from client-go.
func (p *Parser) Build() (*KongState, error) {
	var state KongState
	ings := p.store.ListIngresses()
	tcpIngresses, err := p.store.ListTCPIngresses()
	if err != nil {
		p.Logger.Errorf("failed to list TCPIngresses: %v", err)
	}
	// parse ingress rules
	parsedInfo := p.parseIngressRules(ings, tcpIngresses)

	knativeIngresses, err := p.store.ListKnativeIngresses()
	if err != nil {
		p.Logger.Errorf("failed to list Knative Ingresses: %v", err)
	}
	servicesFromKnative, secretToSNIsFromKnative := p.parseKnativeIngressRules(knativeIngresses)

	for name, service := range servicesFromKnative {
		parsedInfo.ServiceNameToServices[name] = service
	}

	for secret, snis := range secretToSNIsFromKnative {
		var combinedSNIs []string
		if snisFromIngress, ok := parsedInfo.SecretNameToSNIs[secret]; ok {
			combinedSNIs = append(combinedSNIs, snisFromIngress...)
		}
		combinedSNIs = append(combinedSNIs, snis...)
		parsedInfo.SecretNameToSNIs[secret] = combinedSNIs
	}

	// populate Kubernetes Service
	for key, service := range parsedInfo.ServiceNameToServices {
		k8sSvc, err := p.store.GetService(service.Namespace, service.Backend.Name)
		if err != nil {
			p.Logger.WithFields(logrus.Fields{
				"service_name":      service.Backend.Name,
				"service_namespace": service.Namespace,
			}).Errorf("failed to fetch service: %v", err)
		}
		if k8sSvc != nil {
			service.K8sService = *k8sSvc
		}
		secretName := annotations.ExtractClientCertificate(
			service.K8sService.GetAnnotations())
		if secretName != "" {
			secret, err := p.store.GetSecret(service.K8sService.Namespace,
				secretName)
			secretKey := service.K8sService.Namespace + "/" + secretName
			// ensure that the cert is loaded into Kong
			if _, ok := parsedInfo.SecretNameToSNIs[secretKey]; !ok {
				parsedInfo.SecretNameToSNIs[secretKey] = []string{}
			}
			if err == nil {
				service.ClientCertificate = &kong.Certificate{
					ID: kong.String(string(secret.UID)),
				}
			} else {
				p.Logger.WithFields(logrus.Fields{
					"secret_name":      secretName,
					"secret_namespace": service.K8sService.Namespace,
				}).Errorf("failed to fetch secret: %v", err)
			}
		}
		parsedInfo.ServiceNameToServices[key] = service
	}

	// add the routes and services to the state
	for _, service := range parsedInfo.ServiceNameToServices {
		state.Services = append(state.Services, service)
	}

	// generate Upstreams and Targets from service defs
	state.Upstreams = p.getUpstreams(parsedInfo.ServiceNameToServices)

	// merge KongIngress with Routes, Services and Upstream
	p.fillOverrides(state)

	// generate consumers and credentials
	p.fillConsumersAndCredentials(&state)

	// process annotation plugins
	state.Plugins = p.fillPlugins(state)

	// generate Certificates and SNIs
	state.Certificates = p.getCerts(parsedInfo.SecretNameToSNIs)

	// populate CA certificates in Kong
	state.CACertificates, err = p.getCACerts()
	if err != nil {
		return nil, err
	}

	return &state, nil
}

func (p *Parser) getCACerts() ([]kong.CACertificate, error) {
	caCertSecrets, err := p.store.ListCACerts()
	if err != nil {
		return nil, err
	}

	var caCerts []kong.CACertificate
	for _, certSecret := range caCertSecrets {
		secretName := certSecret.Namespace + "/" + certSecret.Name

		idbytes, idExists := certSecret.Data["id"]
		logger := p.Logger.WithFields(logrus.Fields{
			"secret_name":      secretName,
			"secret_namespace": certSecret.Namespace,
		})
		if !idExists {
			logger.Errorf("invalid CA certificate: missing 'id' field in data")
			continue
		}

		caCertbytes, certExists := certSecret.Data["cert"]
		if !certExists {
			logger.Errorf("invalid CA certificate: missing 'cert' field in data")
			continue
		}

		pemBlock, _ := pem.Decode(caCertbytes)
		if pemBlock == nil {
			logger.Errorf("invalid CA certificate: invalid PEM block")
			continue
		}
		x509Cert, err := x509.ParseCertificate(pemBlock.Bytes)
		if err != nil {
			logger.Errorf("invalid CA certificate: failed to parse certificate: %v", err)
			continue
		}
		if !x509Cert.IsCA {
			logger.Errorf("invalid CA certificate: certificate is missing the 'CA' basic constraint: %v", err)
			continue
		}

		caCerts = append(caCerts, kong.CACertificate{
			ID:   kong.String(string(idbytes)),
			Cert: kong.String(string(caCertbytes)),
		})
	}

	return caCerts, nil
}

func (p *Parser) processCredential(credType string, consumer *Consumer,
	credConfig interface{}) error {
	switch credType {
	case "key-auth", "keyauth_credential":
		var cred kong.KeyAuth
		err := decodeCredential(credConfig, &cred)
		if err != nil {
			return fmt.Errorf("failed to decode key-auth credential: %w", err)

		}
		consumer.KeyAuths = append(consumer.KeyAuths, &cred)
	case "basic-auth", "basicauth_credential":
		var cred kong.BasicAuth
		err := decodeCredential(credConfig, &cred)
		if err != nil {
			return fmt.Errorf("failed to decode basic-auth credential: %w", err)
		}
		consumer.BasicAuths = append(consumer.BasicAuths, &cred)
	case "hmac-auth", "hmacauth_credential":
		var cred kong.HMACAuth
		err := decodeCredential(credConfig, &cred)
		if err != nil {
			return fmt.Errorf("failed to decode hmac-auth credential: %w", err)
		}
		consumer.HMACAuths = append(consumer.HMACAuths, &cred)
	case "oauth2":
		var cred kong.Oauth2Credential
		err := decodeCredential(credConfig, &cred)
		if err != nil {
			return fmt.Errorf("failed to decode oauth2 credential: %w", err)
		}
		consumer.Oauth2Creds = append(consumer.Oauth2Creds, &cred)
	case "jwt", "jwt_secret":
		var cred kong.JWTAuth
		err := decodeCredential(credConfig, &cred)
		if err != nil {
			p.Logger.Errorf("failed to process JWT credential: %v", err)
		}
		// This is treated specially because only this
		// field might be omitted by user under the expectation
		// that Kong will insert the default.
		// If we don't set it, decK will detect a diff and PUT this
		// credential everytime it performs a sync operation, which
		// leads to unnecessary cache invalidations in Kong.
		if cred.Algorithm == nil || *cred.Algorithm == "" {
			cred.Algorithm = kong.String("HS256")
		}
		consumer.JWTAuths = append(consumer.JWTAuths, &cred)
	case "acl":
		var cred kong.ACLGroup
		err := decodeCredential(credConfig, &cred)
		if err != nil {
			p.Logger.Errorf("failed to process ACL group: %v", err)
		}
		consumer.ACLGroups = append(consumer.ACLGroups, &cred)
	default:
		return fmt.Errorf("invalid credential type: '%v'", credType)
	}
	return nil
}

func decodeCredential(credConfig interface{},
	credStructPointer interface{}) error {
	decoder, err := mapstructure.NewDecoder(
		&mapstructure.DecoderConfig{TagName: "json",
			Result: credStructPointer,
		})
	if err != nil {
		return fmt.Errorf("failed to create a decoder: %w", err)
	}
	err = decoder.Decode(credConfig)
	if err != nil {
		return fmt.Errorf("failed to decode credential: %w", err)
	}
	return nil
}

func (p *Parser) fillConsumersAndCredentials(state *KongState) {
	consumerIndex := make(map[string]Consumer)

	// build consumer index
	for _, consumer := range p.store.ListKongConsumers() {
		var c Consumer
		if consumer.Username == "" && consumer.CustomID == "" {
			continue
		}
		if consumer.Username != "" {
			c.Username = kong.String(consumer.Username)
		}
		if consumer.CustomID != "" {
			c.CustomID = kong.String(consumer.CustomID)
		}
		c.k8sKongConsumer = *consumer

		logger := p.Logger.WithFields(logrus.Fields{
			"kongconsumer_name":      consumer.Name,
			"kongconsumer_namespace": consumer.Namespace,
		})
		for _, cred := range consumer.Credentials {
			logger = logger.WithFields(logrus.Fields{
				"secret_name":      cred,
				"secret_namespace": consumer.Namespace,
			})
			secret, err := p.store.GetSecret(consumer.Namespace, cred)
			if err != nil {
				logger.Errorf("failed to fetch secret: %v", err)
				continue
			}
			credConfig := map[string]interface{}{}
			for k, v := range secret.Data {
				// TODO populate these based on schema from Kong
				// and remove this workaround
				if k == "redirect_uris" {
					credConfig[k] = strings.Split(string(v), ",")
					continue
				}
				credConfig[k] = string(v)
			}
			credType, ok := credConfig["kongCredType"].(string)
			if !ok {
				logger.Errorf("failed to provision credential: invalid credType: %v", credType)
			}
			if !supportedCreds.Has(credType) {
				logger.Errorf("failed to provision credential: invalid credType: %v", credType)
				continue
			}
			if len(credConfig) <= 1 { // 1 key of credType itself
				logger.Errorf("failed to provision credential: empty secret")
				continue
			}
			err = p.processCredential(credType, &c, credConfig)
			if err != nil {
				logger.Errorf("failed to provision credential: %v", err)
				continue
			}
		}

		consumerIndex[consumer.Namespace+"/"+consumer.Name] = c
	}

	// legacy attach credentials
	credentials := p.store.ListKongCredentials()
	if len(credentials) > 0 {
		p.Logger.Warnf("deprecated KongCredential resource in use; " +
			"please use secret-based credentials, " +
			"KongCredential resource will be removed in future")
	}
	for _, credential := range credentials {
		logger := p.Logger.WithFields(logrus.Fields{
			"kongcredential_name":      credential.Name,
			"kongcredential_namespace": credential.Namespace,
			"consumerRef":              credential.ConsumerRef,
		})
		consumer, ok := consumerIndex[credential.Namespace+"/"+
			credential.ConsumerRef]
		if !ok {
			continue
		}
		if credential.Type == "" {
			logger.Errorf("invalid KongCredential: no Type provided")
			continue
		}
		if !supportedCreds.Has(credential.Type) {
			logger.Errorf("invalid KongCredential: invalid Type provided")
			continue
		}
		if credential.Config == nil {
			logger.Errorf("invalid KongCredential: empty config")
			continue
		}
		err := p.processCredential(credential.Type, &consumer, credential.Config)
		if err != nil {
			logger.Errorf("failed to provision credential: %v", err)
			continue
		}
		consumerIndex[credential.Namespace+"/"+credential.ConsumerRef] = consumer
	}

	// populate the consumer in the state
	for _, c := range consumerIndex {
		state.Consumers = append(state.Consumers, c)
	}
}

func filterHosts(secretNameToSNIs map[string][]string, hosts []string) []string {
	hostsToAdd := []string{}
	seenHosts := map[string]bool{}
	for _, hosts := range secretNameToSNIs {
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

func processTLSSections(tlsSections []networking.IngressTLS,
	namespace string, secretNameToSNIs map[string][]string) {
	// TODO: optmize: collect all TLS sections and process at the same
	// time to avoid regenerating the seen map; or use a seen map in the
	// parser struct itself.
	for _, tls := range tlsSections {
		if len(tls.Hosts) == 0 {
			continue
		}
		if tls.SecretName == "" {
			continue
		}
		hosts := tls.Hosts
		secretName := namespace + "/" + tls.SecretName
		hosts = filterHosts(secretNameToSNIs, hosts)
		if secretNameToSNIs[secretName] != nil {
			hosts = append(hosts, secretNameToSNIs[secretName]...)
		}
		secretNameToSNIs[secretName] = hosts
	}
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

func (p *Parser) parseKnativeIngressRules(
	ingressList []*knative.Ingress) (map[string]Service, map[string][]string) {

	sort.SliceStable(ingressList, func(i, j int) bool {
		return ingressList[i].CreationTimestamp.Before(
			&ingressList[j].CreationTimestamp)
	})

	services := map[string]Service{}
	secretToSNIs := map[string][]string{}

	for i := 0; i < len(ingressList); i++ {
		ingress := *ingressList[i]
		ingressSpec := ingress.Spec

		processTLSSections(knativeIngressToNetworkingTLS(ingress.Spec.TLS),
			ingress.Namespace, secretToSNIs)
		for i, rule := range ingressSpec.Rules {
			if rule.HTTP == nil {
				continue
			}
			hosts := explandClusterLocal(rule.Hosts)
			for j, rule := range rule.HTTP.Paths {
				path := rule.Path

				if path == "" {
					path = "/"
				}
				r := Route{
					Route: kong.Route{
						// TODO Figure out a way to name the routes
						// This is not a stable scheme
						// 1. If a user adds a route in the middle,
						// due to a shift, all the following routes will
						// be PATCHED
						// 2. Is it guaranteed that the order is stable?
						// Meaning, the routes will always appear in the same
						// order?
						Name:          kong.String(ingress.Namespace + "." + ingress.Name + "." + strconv.Itoa(i) + strconv.Itoa(j)),
						Paths:         kong.StringSlice(path),
						StripPath:     kong.Bool(false),
						PreserveHost:  kong.Bool(true),
						Protocols:     kong.StringSlice("http", "https"),
						RegexPriority: kong.Int(0),
					},
				}
				r.Hosts = kong.StringSlice(hosts...)

				knativeBackend := knativeSelectSplit(rule.Splits)
				serviceName := knativeBackend.ServiceNamespace + "." +
					knativeBackend.ServiceName + "." +
					knativeBackend.ServicePort.String()
				serviceHost := knativeBackend.ServiceName + "." +
					knativeBackend.ServiceNamespace + "." +
					knativeBackend.ServicePort.String() + ".svc"
				service, ok := services[serviceName]
				if !ok {

					var headers []string
					for key, value := range knativeBackend.AppendHeaders {
						headers = append(headers, key+":"+value)
					}
					for key, value := range rule.AppendHeaders {
						headers = append(headers, key+":"+value)
					}

					service = Service{
						Service: kong.Service{
							Name:           kong.String(serviceName),
							Host:           kong.String(serviceHost),
							Port:           kong.Int(80),
							Protocol:       kong.String("http"),
							Path:           kong.String("/"),
							ConnectTimeout: kong.Int(60000),
							ReadTimeout:    kong.Int(60000),
							WriteTimeout:   kong.Int(60000),
							Retries:        kong.Int(5),
						},
						Namespace: ingress.Namespace,
						Backend: backend{
							Name: knativeBackend.ServiceName,
							Port: knativeBackend.ServicePort,
						},
					}
					if len(headers) > 0 {
						service.Plugins = append(service.Plugins, kong.Plugin{
							Name: kong.String("request-transformer"),
							Config: kong.Configuration{
								"add": map[string]interface{}{
									"headers": headers,
								},
							},
						})
					}
				}
				service.Routes = append(service.Routes, r)
				services[serviceName] = service
			}
		}
	}

	return services, secretToSNIs
}

// explandClusterLocal expands cluster local domains to support short domains.
// e.g. hello.default.svc.cluster.local expands hello.default.svc.cluster.local, hello.default.svc and hello.default
func explandClusterLocal(hosts []string) []string {
	var expand []string
	localDomainSuffix := ".svc." + network.GetClusterDomainName()
	for _, host := range hosts {
		expand = append(expand, host)
		if strings.HasSuffix(host, localDomainSuffix) {
			expand = append(expand, strings.TrimSuffix(host, localDomainSuffix))
			expand = append(expand, strings.TrimSuffix(host, "."+network.GetClusterDomainName()))
		}
	}
	return expand
}

func knativeSelectSplit(splits []knative.IngressBackendSplit) knative.IngressBackendSplit {
	if len(splits) == 0 {
		return knative.IngressBackendSplit{}
	}
	res := splits[0]
	maxPercentage := splits[0].Percent
	if len(splits) == 1 {
		return res
	}
	for i := 1; i < len(splits); i++ {
		if splits[i].Percent > maxPercentage {
			res = splits[i]
			maxPercentage = res.Percent
		}
	}
	return res
}

func (p *Parser) parseIngressRules(
	ingressList []*networking.Ingress,
	tcpIngressList []*configurationv1beta1.TCPIngress) *parsedIngressRules {

	sort.SliceStable(ingressList, func(i, j int) bool {
		return ingressList[i].CreationTimestamp.Before(
			&ingressList[j].CreationTimestamp)
	})

	sort.SliceStable(tcpIngressList, func(i, j int) bool {
		return tcpIngressList[i].CreationTimestamp.Before(
			&tcpIngressList[j].CreationTimestamp)
	})

	// generate the following:
	// Services and Routes
	var allDefaultBackends []networking.Ingress
	secretNameToSNIs := make(map[string][]string)
	serviceNameToServices := make(map[string]Service)

	for i := 0; i < len(ingressList); i++ {
		ingress := *ingressList[i]
		ingressSpec := ingress.Spec
		logger := p.Logger.WithFields(logrus.Fields{
			"ingress_namespace": ingress.Namespace,
			"ingress_name":      ingress.Name,
		})

		if ingressSpec.Backend != nil {
			allDefaultBackends = append(allDefaultBackends, ingress)

		}

		processTLSSections(ingressSpec.TLS, ingress.Namespace, secretNameToSNIs)

		for i, rule := range ingressSpec.Rules {
			host := rule.Host
			if rule.HTTP == nil {
				continue
			}
			for j, rule := range rule.HTTP.Paths {
				path := rule.Path

				if strings.Contains(path, "//") {
					logger.Errorf("rule skipped: invalid path: '%v'", path)
					continue
				}
				if path == "" {
					path = "/"
				}
				r := Route{
					Ingress: ingress,
					Route: kong.Route{
						// TODO Figure out a way to name the routes
						// This is not a stable scheme
						// 1. If a user adds a route in the middle,
						// due to a shift, all the following routes will
						// be PATCHED
						// 2. Is it guaranteed that the order is stable?
						// Meaning, the routes will always appear in the same
						// order?
						Name:          kong.String(ingress.Namespace + "." + ingress.Name + "." + strconv.Itoa(i) + strconv.Itoa(j)),
						Paths:         kong.StringSlice(path),
						StripPath:     kong.Bool(false),
						PreserveHost:  kong.Bool(true),
						Protocols:     kong.StringSlice("http", "https"),
						RegexPriority: kong.Int(0),
					},
				}
				if host != "" {
					r.Hosts = kong.StringSlice(host)
				}

				serviceName := ingress.Namespace + "." +
					rule.Backend.ServiceName + "." +
					rule.Backend.ServicePort.String()
				service, ok := serviceNameToServices[serviceName]
				if !ok {
					service = Service{
						Service: kong.Service{
							Name: kong.String(serviceName),
							Host: kong.String(rule.Backend.ServiceName +
								"." + ingress.Namespace + "." +
								rule.Backend.ServicePort.String() + ".svc"),
							Port:           kong.Int(80),
							Protocol:       kong.String("http"),
							Path:           kong.String("/"),
							ConnectTimeout: kong.Int(60000),
							ReadTimeout:    kong.Int(60000),
							WriteTimeout:   kong.Int(60000),
							Retries:        kong.Int(5),
						},
						Namespace: ingress.Namespace,
						Backend: backend{
							Name: rule.Backend.ServiceName,
							Port: rule.Backend.ServicePort,
						},
					}
				}
				service.Routes = append(service.Routes, r)
				serviceNameToServices[serviceName] = service
			}
		}
	}

	for i := 0; i < len(tcpIngressList); i++ {
		ingress := *tcpIngressList[i]
		ingressSpec := ingress.Spec

		logger := p.Logger.WithFields(logrus.Fields{
			"tcpingress_namespace": ingress.Namespace,
			"tcpingress_name":      ingress.Name,
		})

		processTLSSections(tcpIngressToNetworkingTLS(ingressSpec.TLS),
			ingress.Namespace, secretNameToSNIs)

		for i, rule := range ingressSpec.Rules {

			if rule.Port <= 0 {
				logger.Errorf("invalid TCPIngress: invalid port: %v", rule.Port)
				continue
			}
			r := Route{
				IsTCP:      true,
				TCPIngress: ingress,
				Route: kong.Route{
					// TODO Figure out a way to name the routes
					// This is not a stable scheme
					// 1. If a user adds a route in the middle,
					// due to a shift, all the following routes will
					// be PATCHED
					// 2. Is it guaranteed that the order is stable?
					// Meaning, the routes will always appear in the same
					// order?
					Name:      kong.String(ingress.Namespace + "." + ingress.Name + "." + strconv.Itoa(i)),
					Protocols: kong.StringSlice("tcp", "tls"),
					Destinations: []*kong.CIDRPort{
						{
							Port: kong.Int(rule.Port),
						},
					},
				},
			}
			host := rule.Host
			if host != "" {
				r.SNIs = kong.StringSlice(host)
			}
			if rule.Backend.ServiceName == "" {
				logger.Errorf("invalid TCPIngress: empty serviceName")
				continue
			}
			if rule.Backend.ServicePort <= 0 {
				logger.Errorf("invalid TCPIngress: invalid servicePort: %v", rule.Backend.ServicePort)
				continue
			}

			serviceName := ingress.Namespace + "." +
				rule.Backend.ServiceName + "." +
				strconv.Itoa(rule.Backend.ServicePort)
			service, ok := serviceNameToServices[serviceName]
			if !ok {
				service = Service{
					Service: kong.Service{
						Name: kong.String(serviceName),
						Host: kong.String(rule.Backend.ServiceName +
							"." + ingress.Namespace + "." +
							strconv.Itoa(rule.Backend.ServicePort) + ".svc"),
						Port:           kong.Int(80),
						Protocol:       kong.String("tcp"),
						ConnectTimeout: kong.Int(60000),
						ReadTimeout:    kong.Int(60000),
						WriteTimeout:   kong.Int(60000),
						Retries:        kong.Int(5),
					},
					Namespace: ingress.Namespace,
					Backend: backend{
						Name: rule.Backend.ServiceName,
						Port: intstr.FromInt(rule.Backend.ServicePort),
					},
				}
			}
			service.Routes = append(service.Routes, r)
			serviceNameToServices[serviceName] = service
		}
	}

	sort.SliceStable(allDefaultBackends, func(i, j int) bool {
		return allDefaultBackends[i].CreationTimestamp.Before(&allDefaultBackends[j].CreationTimestamp)
	})

	// Process the default backend
	if len(allDefaultBackends) > 0 {
		ingress := allDefaultBackends[0]
		defaultBackend := allDefaultBackends[0].Spec.Backend
		serviceName := allDefaultBackends[0].Namespace + "." +
			defaultBackend.ServiceName + "." +
			defaultBackend.ServicePort.String()
		service, ok := serviceNameToServices[serviceName]
		if !ok {
			service = Service{
				Service: kong.Service{
					Name: kong.String(serviceName),
					Host: kong.String(defaultBackend.ServiceName + "." +
						ingress.Namespace + "." +
						defaultBackend.ServicePort.String() + ".svc"),
					Port:           kong.Int(80),
					Protocol:       kong.String("http"),
					ConnectTimeout: kong.Int(60000),
					ReadTimeout:    kong.Int(60000),
					WriteTimeout:   kong.Int(60000),
					Retries:        kong.Int(5),
				},
				Namespace: ingress.Namespace,
				Backend: backend{
					Name: defaultBackend.ServiceName,
					Port: defaultBackend.ServicePort,
				},
			}
		}
		r := Route{
			Ingress: ingress,
			Route: kong.Route{
				Name:          kong.String(ingress.Namespace + "." + ingress.Name),
				Paths:         kong.StringSlice("/"),
				StripPath:     kong.Bool(false),
				PreserveHost:  kong.Bool(true),
				Protocols:     kong.StringSlice("http", "https"),
				RegexPriority: kong.Int(0),
			},
		}
		service.Routes = append(service.Routes, r)
		serviceNameToServices[serviceName] = service
	}

	return &parsedIngressRules{
		SecretNameToSNIs:      secretNameToSNIs,
		ServiceNameToServices: serviceNameToServices,
	}
}

func (p *Parser) fillOverrides(state KongState) {
	for i := 0; i < len(state.Services); i++ {
		// Services
		anns := state.Services[i].K8sService.Annotations
		kongIngress, err := p.getKongIngressForService(
			state.Services[i].K8sService)
		if err != nil {
			p.Logger.WithFields(logrus.Fields{
				"service_name":      state.Services[i].K8sService.Name,
				"service_namespace": state.Services[i].K8sService.Namespace,
			}).Errorf("failed to fetch KongIngress resource for Service: %v", err)
		}
		overrideService(&state.Services[i], kongIngress, anns)

		// Routes
		for j := 0; j < len(state.Services[i].Routes); j++ {
			var kongIngress *configurationv1.KongIngress
			var err error
			if state.Services[i].Routes[j].IsTCP {
				kongIngress, err = p.getKongIngressFromTCPIngress(
					&state.Services[i].Routes[j].TCPIngress)
				if err != nil {
					p.Logger.WithFields(logrus.Fields{
						"tcpingress_name":      state.Services[i].Routes[j].TCPIngress.Name,
						"tcpingress_namespace": state.Services[i].Routes[j].TCPIngress.Namespace,
					}).Errorf("failed to fetch KongIngress resource for Ingress: %v", err)
				}
			} else {
				kongIngress, err = p.getKongIngressFromIngress(
					&state.Services[i].Routes[j].Ingress)
				if err != nil {
					p.Logger.WithFields(logrus.Fields{
						"ingress_name":      state.Services[i].Routes[j].Ingress.Name,
						"ingress_namespace": state.Services[i].Routes[j].Ingress.Namespace,
					}).Errorf("failed to fetch KongIngress resource for Ingress: %v", err)
				}
			}

			p.overrideRoute(&state.Services[i].Routes[j], kongIngress)
		}
	}

	// Upstreams
	for i := 0; i < len(state.Upstreams); i++ {
		kongIngress, err := p.getKongIngressForService(
			state.Upstreams[i].Service.K8sService)
		anns := state.Upstreams[i].Service.K8sService.Annotations
		if err != nil {
			p.Logger.WithFields(logrus.Fields{
				"service_name":      state.Upstreams[i].Service.K8sService.Name,
				"service_namespace": state.Upstreams[i].Service.K8sService.Namespace,
			}).Errorf("failed to fetch KongIngress resource for Service: %v", err)
			continue
		}
		overrideUpstream(&state.Upstreams[i], kongIngress, anns)
	}
}

// overrideServiceByKongIngress sets Service fields by KongIngress
func overrideServiceByKongIngress(service *Service,
	kongIngress *configurationv1.KongIngress) {
	if kongIngress == nil || kongIngress.Proxy == nil {
		return
	}
	s := kongIngress.Proxy
	if s.Protocol != nil {
		service.Protocol = kong.String(*s.Protocol)
	}
	if s.Path != nil {
		service.Path = kong.String(*s.Path)
	}
	if s.Retries != nil {
		service.Retries = kong.Int(*s.Retries)
	}
	if s.ConnectTimeout != nil {
		service.ConnectTimeout = kong.Int(*s.ConnectTimeout)
	}
	if s.ReadTimeout != nil {
		service.ReadTimeout = kong.Int(*s.ReadTimeout)
	}
	if s.WriteTimeout != nil {
		service.WriteTimeout = kong.Int(*s.WriteTimeout)
	}
}

func overrideServicePath(service *kong.Service, anns map[string]string) {
	if service == nil {
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
	service.Path = kong.String(path)
}

func overrideServiceProtocol(service *kong.Service, anns map[string]string) {
	if service == nil {
		return
	}
	protocol := annotations.ExtractProtocolName(anns)
	if protocol == "" || !validateProtocol(protocol) {
		return
	}
	service.Protocol = kong.String(protocol)
}

// overrideServiceByAnnotation modifies the Kong service based on annotations
// on the Kubernetes service.
func overrideServiceByAnnotation(service *kong.Service,
	anns map[string]string) {
	if service == nil {
		return
	}
	overrideServiceProtocol(service, anns)
	overrideServicePath(service, anns)
}

// overrideService sets Service fields by KongIngress first, then by annotation
func overrideService(service *Service,
	kongIngress *configurationv1.KongIngress,
	anns map[string]string) {
	if service == nil {
		return
	}
	overrideServiceByKongIngress(service, kongIngress)
	overrideServiceByAnnotation(&service.Service, anns)

	if *service.Protocol == "grpc" || *service.Protocol == "grpcs" {
		// grpc(s) doesn't accept a path
		service.Path = nil
	}
}

// overrideRouteByKongIngress sets Route fields by KongIngress
func (p *Parser) overrideRouteByKongIngress(route *Route,
	kongIngress *configurationv1.KongIngress) {
	if kongIngress == nil || kongIngress.Route == nil {
		return
	}

	r := kongIngress.Route
	if len(r.Methods) != 0 {
		invalid := false
		var methods []*string
		for _, method := range r.Methods {
			sanitizedMethod := strings.TrimSpace(strings.ToUpper(*method))
			if validMethods.MatchString(sanitizedMethod) {
				methods = append(methods, kong.String(sanitizedMethod))
			} else {
				// if any method is invalid (not an uppercase alpha string),
				// discard everything
				p.Logger.WithFields(logrus.Fields{
					"ingress_namespace": route.Ingress.Namespace,
					"ingress_name":      route.Ingress.Name,
				}).Errorf("ingress contains invalid method: '%v'", *method)
				invalid = true
			}
		}
		if !invalid {
			route.Methods = methods
		}
	}
	if len(r.Headers) != 0 {
		route.Headers = r.Headers
	}
	if len(r.Protocols) != 0 {
		route.Protocols = cloneStringPointerSlice(r.Protocols...)
	}
	if r.RegexPriority != nil {
		route.RegexPriority = kong.Int(*r.RegexPriority)
	}
	if r.StripPath != nil {
		route.StripPath = kong.Bool(*r.StripPath)
	}
	if r.PreserveHost != nil {
		route.PreserveHost = kong.Bool(*r.PreserveHost)
	}
	if r.HTTPSRedirectStatusCode != nil {
		route.HTTPSRedirectStatusCode = kong.Int(*r.HTTPSRedirectStatusCode)
	}
	if r.PathHandling != nil {
		route.PathHandling = kong.String(*r.PathHandling)
	}
}

// normalizeProtocols prevents users from mismatching grpc/http
func normalizeProtocols(route *Route) {
	protocols := route.Protocols
	var http, grpc bool

	for _, protocol := range protocols {
		if strings.Contains(*protocol, "grpc") {
			grpc = true
		}
		if strings.Contains(*protocol, "http") {
			http = true
		}
		if !validateProtocol(*protocol) {
			http = true
		}
	}

	if grpc && http {
		route.Protocols = kong.StringSlice("http", "https")
	}
}

// validateProtocol returns a bool of whether string is a valid protocol
func validateProtocol(protocol string) bool {
	match := validProtocols.MatchString(protocol)
	return match
}

// useSSLProtocol updates the protocol of the route to either https or grpcs, or https and grpcs
func useSSLProtocol(route *kong.Route) {
	var http, grpc bool
	var prots []*string

	for _, val := range route.Protocols {

		if strings.Contains(*val, "grpc") {
			grpc = true
		}

		if strings.Contains(*val, "http") {
			http = true
		}
	}

	if grpc {
		prots = append(prots, kong.String("grpcs"))
	}
	if http {
		prots = append(prots, kong.String("https"))
	}

	if !grpc && !http {
		prots = append(prots, kong.String("https"))
	}

	route.Protocols = prots
}
func overrideRouteStripPath(route *kong.Route, anns map[string]string) {
	if route == nil {
		return
	}

	stripPathValue := annotations.ExtractStripPath(anns)
	if stripPathValue == "" {
		return
	}
	stripPathValue = strings.ToLower(stripPathValue)
	switch stripPathValue {
	case "true":
		route.StripPath = kong.Bool(true)
	case "false":
		route.StripPath = kong.Bool(false)
	default:
		return
	}
}

func overrideRouteProtocols(route *kong.Route, anns map[string]string) {
	protocols := annotations.ExtractProtocolNames(anns)
	var prots []*string
	for _, prot := range protocols {
		if !validateProtocol(prot) {
			return
		}
		prots = append(prots, kong.String(prot))
	}

	route.Protocols = prots
}

func overrideRouteHTTPSRedirectCode(route *kong.Route, anns map[string]string) {

	if annotations.HasForceSSLRedirectAnnotation(anns) {
		route.HTTPSRedirectStatusCode = kong.Int(302)
		useSSLProtocol(route)
	}

	code := annotations.ExtractHTTPSRedirectStatusCode(anns)
	if code == "" {
		return
	}
	statusCode, err := strconv.Atoi(code)
	if err != nil {
		return
	}
	if statusCode != 426 &&
		statusCode != 301 &&
		statusCode != 302 &&
		statusCode != 307 &&
		statusCode != 308 {
		return
	}

	route.HTTPSRedirectStatusCode = kong.Int(statusCode)
}

func overrideRoutePreserveHost(route *kong.Route, anns map[string]string) {
	preserveHostValue := annotations.ExtractPreserveHost(anns)
	if preserveHostValue == "" {
		return
	}
	preserveHostValue = strings.ToLower(preserveHostValue)
	switch preserveHostValue {
	case "true":
		route.PreserveHost = kong.Bool(true)
	case "false":
		route.PreserveHost = kong.Bool(false)
	default:
		return
	}
}

func overrideRouteRegexPriority(route *kong.Route, anns map[string]string) {
	priority := annotations.ExtractRegexPriority(anns)
	if priority == "" {
		return
	}
	regexPriority, err := strconv.Atoi(priority)
	if err != nil {
		return
	}

	route.RegexPriority = kong.Int(regexPriority)
}

func (p *Parser) overrideRouteMethods(route *kong.Route, anns map[string]string) {
	annMethods := annotations.ExtractMethods(anns)
	if len(annMethods) == 0 {
		return
	}
	var methods []*string
	for _, method := range annMethods {
		sanitizedMethod := strings.TrimSpace(strings.ToUpper(method))
		if validMethods.MatchString(sanitizedMethod) {
			methods = append(methods, kong.String(sanitizedMethod))
		} else {
			// if any method is invalid (not an uppercase alpha string),
			// discard everything
			p.Logger.WithField("kongroute", route.Name).Errorf("invalid method: %v", method)
			return
		}
	}

	route.Methods = methods
}

// overrideRouteByAnnotation sets Route protocols via annotation
func (p *Parser) overrideRouteByAnnotation(route *Route) {
	anns := route.Ingress.Annotations
	if route.IsTCP {
		anns = route.TCPIngress.Annotations
	}
	overrideRouteProtocols(&route.Route, anns)
	overrideRouteStripPath(&route.Route, anns)
	overrideRouteHTTPSRedirectCode(&route.Route, anns)
	overrideRoutePreserveHost(&route.Route, anns)
	overrideRouteRegexPriority(&route.Route, anns)
	p.overrideRouteMethods(&route.Route, anns)
}

// overrideRoute sets Route fields by KongIngress first, then by annotation
func (p *Parser) overrideRoute(route *Route,
	kongIngress *configurationv1.KongIngress) {
	if route == nil {
		return
	}
	p.overrideRouteByKongIngress(route, kongIngress)
	p.overrideRouteByAnnotation(route)
	normalizeProtocols(route)
	for _, val := range route.Protocols {
		if *val == "grpc" || *val == "grpcs" {
			// grpc(s) doesn't accept strip_path
			route.StripPath = nil
			break
		}
	}
}

func cloneStringPointerSlice(array ...*string) (res []*string) {
	res = append(res, array...)
	return
}

func overrideUpstreamHostHeader(upstream *kong.Upstream, anns map[string]string) {
	if upstream == nil {
		return
	}
	host := annotations.ExtractHostHeader(anns)
	if host == "" {
		return
	}
	upstream.HostHeader = kong.String(host)
}

// overrideUpstreamByAnnotation modifies the Kong upstream based on annotations
// on the Kubernetes service.
func overrideUpstreamByAnnotation(upstream *kong.Upstream,
	anns map[string]string) {
	if upstream == nil {
		return
	}
	overrideUpstreamHostHeader(upstream, anns)
}

// overrideUpstreamByKongIngress modifies the Kong upstream based on KongIngresses
// associated with the Kubernetes service.
func overrideUpstreamByKongIngress(upstream *Upstream,
	kongIngress *configurationv1.KongIngress) {
	if upstream == nil {
		return
	}

	if kongIngress == nil || kongIngress.Upstream == nil {
		return
	}

	// The upstream within the KongIngress has no name.
	// As this overwrites the entire upstream object, we must restore the
	// original name after.
	name := *upstream.Upstream.Name
	upstream.Upstream = *kongIngress.Upstream.DeepCopy()
	upstream.Name = &name
}

// overrideUpstream sets Upstream fields by KongIngress first, then by annotation
func overrideUpstream(upstream *Upstream,
	kongIngress *configurationv1.KongIngress,
	anns map[string]string) {
	if upstream == nil {
		return
	}

	overrideUpstreamByKongIngress(upstream, kongIngress)
	overrideUpstreamByAnnotation(&upstream.Upstream, anns)
}

func (p *Parser) getUpstreams(serviceMap map[string]Service) []Upstream {
	var upstreams []Upstream
	for _, service := range serviceMap {
		upstreamName := service.Backend.Name + "." + service.Namespace + "." + service.Backend.Port.String() + ".svc"
		upstream := Upstream{
			Upstream: kong.Upstream{
				Name: kong.String(upstreamName),
			},
			Service: service,
		}
		targets := p.getServiceEndpoints(service.K8sService,
			service.Backend.Port.String())
		upstream.Targets = targets
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

func (p *Parser) getCerts(secretsToSNIs map[string][]string) []Certificate {
	snisAdded := make(map[string]bool)
	// map of cert public key + private key to certificate
	type certWrapper struct {
		cert              kong.Certificate
		CreationTimestamp metav1.Time
	}
	certs := make(map[string]certWrapper)

	for secretKey, SNIs := range secretsToSNIs {
		namespaceName := strings.Split(secretKey, "/")
		secret, err := p.store.GetSecret(namespaceName[0], namespaceName[1])
		if err != nil {
			p.Logger.WithFields(logrus.Fields{
				"secret_name":      namespaceName[1],
				"secret_namespace": namespaceName[0],
			}).Logger.Errorf("failed to fetch secret: %v", err)
			continue
		}
		cert, key, err := getCertFromSecret(secret)
		if err != nil {
			p.Logger.WithFields(logrus.Fields{
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
	var res []Certificate
	for _, cert := range certs {
		res = append(res, Certificate{cert.cert})
	}
	return res
}

type foreignRelations struct {
	Consumer, Route, Service []string
}

func getPluginRelations(state KongState) map[string]foreignRelations {
	// KongPlugin key (KongPlugin's name:namespace) to corresponding associations
	pluginRels := map[string]foreignRelations{}
	addConsumerRelation := func(namespace, pluginName, identifier string) {
		pluginKey := namespace + ":" + pluginName
		relations, ok := pluginRels[pluginKey]
		if !ok {
			relations = foreignRelations{}
		}
		relations.Consumer = append(relations.Consumer, identifier)
		pluginRels[pluginKey] = relations
	}
	addRouteRelation := func(namespace, pluginName, identifier string) {
		pluginKey := namespace + ":" + pluginName
		relations, ok := pluginRels[pluginKey]
		if !ok {
			relations = foreignRelations{}
		}
		relations.Route = append(relations.Route, identifier)
		pluginRels[pluginKey] = relations
	}
	addServiceRelation := func(namespace, pluginName, identifier string) {
		pluginKey := namespace + ":" + pluginName
		relations, ok := pluginRels[pluginKey]
		if !ok {
			relations = foreignRelations{}
		}
		relations.Service = append(relations.Service, identifier)
		pluginRels[pluginKey] = relations
	}

	for i := range state.Services {
		// service
		svc := state.Services[i].K8sService
		pluginList := annotations.ExtractKongPluginsFromAnnotations(
			svc.GetAnnotations())
		for _, pluginName := range pluginList {
			addServiceRelation(svc.Namespace, pluginName,
				*state.Services[i].Name)
		}
		// route
		for j := range state.Services[i].Routes {
			ingress := state.Services[i].Routes[j].Ingress
			pluginList := annotations.ExtractKongPluginsFromAnnotations(ingress.GetAnnotations())
			for _, pluginName := range pluginList {
				addRouteRelation(ingress.Namespace, pluginName, *state.Services[i].Routes[j].Name)
			}
		}
	}
	// consumer
	for _, c := range state.Consumers {
		pluginList := annotations.ExtractKongPluginsFromAnnotations(c.k8sKongConsumer.GetAnnotations())
		for _, pluginName := range pluginList {
			addConsumerRelation(c.k8sKongConsumer.Namespace, pluginName, *c.Username)
		}
	}
	return pluginRels
}

type rel struct {
	Consumer, Route, Service string
}

func getCombinations(relations foreignRelations) []rel {

	var cartesianProduct []rel

	if len(relations.Consumer) > 0 {
		consumers := relations.Consumer
		if len(relations.Route)+len(relations.Service) > 0 {
			for _, service := range relations.Service {
				for _, consumer := range consumers {
					cartesianProduct = append(cartesianProduct, rel{
						Service:  service,
						Consumer: consumer,
					})
				}
			}
			for _, route := range relations.Route {
				for _, consumer := range consumers {
					cartesianProduct = append(cartesianProduct, rel{
						Route:    route,
						Consumer: consumer,
					})
				}
			}
		} else {
			for _, consumer := range relations.Consumer {
				cartesianProduct = append(cartesianProduct, rel{Consumer: consumer})
			}
		}
	} else {
		for _, service := range relations.Service {
			cartesianProduct = append(cartesianProduct, rel{Service: service})
		}
		for _, route := range relations.Route {
			cartesianProduct = append(cartesianProduct, rel{Route: route})
		}
	}

	return cartesianProduct
}

func (p *Parser) fillPlugins(state KongState) []Plugin {
	var plugins []Plugin
	pluginRels := getPluginRelations(state)

	for pluginIdentifier, relations := range pluginRels {
		identifier := strings.Split(pluginIdentifier, ":")
		namespace, kongPluginName := identifier[0], identifier[1]
		plugin, err := p.getPlugin(namespace, kongPluginName)
		if err != nil {
			p.Logger.WithFields(logrus.Fields{
				"kongplugin_name":      kongPluginName,
				"kongplugin_namespace": namespace,
			}).Logger.Errorf("failed to fetch KongPlugin: %v", err)
			continue
		}

		for _, rel := range getCombinations(relations) {
			plugin := *plugin.DeepCopy()
			// ID is populated because that is read by decK and in_memory
			// translator too
			if rel.Service != "" {
				plugin.Service = &kong.Service{ID: kong.String(rel.Service)}
			}
			if rel.Route != "" {
				plugin.Route = &kong.Route{ID: kong.String(rel.Route)}
			}
			if rel.Consumer != "" {
				plugin.Consumer = &kong.Consumer{ID: kong.String(rel.Consumer)}
			}
			plugins = append(plugins, Plugin{plugin})
		}
	}

	globalPlugins, err := p.globalPlugins()
	if err != nil {
		p.Logger.Errorf("failed to fetch global plugins: %v", err)
	}
	plugins = append(plugins, globalPlugins...)

	return plugins
}

func (p *Parser) globalPlugins() ([]Plugin, error) {
	// removed as of 0.10.0
	// only retrieved now to warn users
	globalPlugins, err := p.store.ListGlobalKongPlugins()
	if err != nil {
		return nil, fmt.Errorf("error listing global KongPlugins: %w", err)
	}
	if len(globalPlugins) > 0 {
		p.Logger.Warning("global KongPlugins found. These are no longer applied and",
			" must be replaced with KongClusterPlugins.",
			" Please run \"kubectl get kongplugin -l global=true --all-namespaces\" to list existing plugins")
	}
	res := make(map[string]Plugin)
	var duplicates []string // keep track of duplicate
	// TODO respect the oldest CRD
	// Current behavior is to skip creating the plugin but in case
	// of duplicate plugin definitions, we should respect the oldest one
	// This is important since if a user comes in to k8s and creates a new
	// CRD, the user now deleted an older plugin

	globalClusterPlugins, err := p.store.ListGlobalKongClusterPlugins()
	if err != nil {
		return nil, fmt.Errorf("error listing global KongClusterPlugins: %w", err)
	}
	for i := 0; i < len(globalClusterPlugins); i++ {
		k8sPlugin := *globalClusterPlugins[i]
		pluginName := k8sPlugin.PluginName
		// empty pluginName skip it
		if pluginName == "" {
			p.Logger.WithFields(logrus.Fields{
				"kongclusterplugin_name": k8sPlugin.Name,
			}).Errorf("invalid KongClusterPlugin: empty plugin property")
			continue
		}
		if _, ok := res[pluginName]; ok {
			p.Logger.Error("multiple KongPlugin definitions found with"+
				" 'global' label for '", pluginName,
				"', the plugin will not be applied")
			duplicates = append(duplicates, pluginName)
			continue
		}
		if plugin, err := p.kongPluginFromK8SClusterPlugin(k8sPlugin); err == nil {
			res[pluginName] = Plugin{
				Plugin: plugin,
			}
		} else {
			p.Logger.WithFields(logrus.Fields{
				"kongclusterplugin_name": k8sPlugin.Name,
			}).Errorf("failed to generate configuration from KongClusterPlugin: %v ", err)
		}
	}
	for _, plugin := range duplicates {
		delete(res, plugin)
	}
	var plugins []Plugin
	for _, p := range res {
		plugins = append(plugins, p)
	}
	return plugins, nil
}

func (p *Parser) getServiceEndpoints(svc corev1.Service,
	backendPort string) []Target {
	var targets []Target
	var endpoints []utils.Endpoint
	var servicePort corev1.ServicePort

	logger := p.Logger.WithFields(logrus.Fields{
		"service_name":      svc.Name,
		"service_namespace": svc.Namespace,
	})

	for _, port := range svc.Spec.Ports {
		// targetPort could be a string, use the name or the port (int)
		if strconv.Itoa(int(port.Port)) == backendPort ||
			port.TargetPort.String() == backendPort ||
			port.Name == backendPort {
			servicePort = port
			break
		}
	}

	// Ingress with an ExternalName service and no port defined in the service.
	if len(svc.Spec.Ports) == 0 &&
		svc.Spec.Type == corev1.ServiceTypeExternalName {
		// nolint: gosec
		externalPort, err := strconv.Atoi(backendPort)
		if err != nil {
			logger.Warningf("invalid ExternalName Service (only numeric ports allowed): %v", backendPort)
			return targets
		}

		servicePort = corev1.ServicePort{
			Protocol:   "TCP",
			Port:       int32(externalPort),
			TargetPort: intstr.FromString(backendPort),
		}
	}

	endpoints = p.getEndpoints(&svc, &servicePort,
		corev1.ProtocolTCP, p.store.GetEndpointsForService)
	if len(endpoints) == 0 {
		logger.Warningf("no active endpionts")
	}
	for _, endpoint := range endpoints {
		target := Target{
			Target: kong.Target{
				Target: kong.String(endpoint.Address + ":" + endpoint.Port),
			},
		}
		targets = append(targets, target)
	}
	return targets
}

func (p *Parser) getKongIngressForService(service corev1.Service) (
	*configurationv1.KongIngress, error) {
	confName := annotations.ExtractConfigurationName(service.Annotations)
	if confName == "" {
		return nil, nil
	}
	return p.store.GetKongIngress(service.Namespace, confName)
}

func (p *Parser) getKongIngressFromIngressAnnotations(namespace, name string,
	anns map[string]string) (
	*configurationv1.KongIngress, error) {
	confName := annotations.ExtractConfigurationName(anns)
	if confName != "" {
		ki, err := p.store.GetKongIngress(namespace, confName)
		if err == nil {
			return ki, nil
		}
	}

	ki, err := p.store.GetKongIngress(namespace, name)
	if err == nil {
		return ki, nil
	}
	return nil, nil
}

// getKongIngressFromIngress checks if the Ingress
// contains an annotation for configuration
// or if exists a KongIngress object with the same name than the Ingress
func (p *Parser) getKongIngressFromIngress(ing *networking.Ingress) (
	*configurationv1.KongIngress, error) {
	return p.getKongIngressFromIngressAnnotations(ing.Namespace,
		ing.Name, ing.Annotations)
}

// getKongIngressFromTCPIngress checks if the TCPIngress contains an
// annotation for configuration
// or if exists a KongIngress object with the same name than the Ingress
func (p *Parser) getKongIngressFromTCPIngress(ing *configurationv1beta1.TCPIngress) (
	*configurationv1.KongIngress, error) {
	return p.getKongIngressFromIngressAnnotations(ing.Namespace,
		ing.Name, ing.Annotations)
}

// getPlugin constructs a plugins from a KongPlugin resource.
func (p *Parser) getPlugin(namespace, name string) (kong.Plugin, error) {
	var plugin kong.Plugin
	k8sPlugin, err := p.store.GetKongPlugin(namespace, name)
	if err != nil {
		// if no namespaced plugin definition, then
		// search for cluster level-plugin definition
		if errors.As(err, &store.ErrNotFound{}) {
			clusterPlugin, err := p.store.GetKongClusterPlugin(name)
			// not found
			if errors.As(err, &store.ErrNotFound{}) {
				return plugin, errors.New(
					"no KongPlugin or KongClusterPlugin was found")
			}
			if err != nil {
				return plugin, err
			}
			if clusterPlugin.PluginName == "" {
				return plugin, fmt.Errorf("invalid empty 'plugin' property")
			}
			plugin, err = p.kongPluginFromK8SClusterPlugin(*clusterPlugin)
			return plugin, err
		}
	}
	// ignore plugins with no name
	if k8sPlugin.PluginName == "" {
		return plugin, fmt.Errorf("invalid empty 'plugin' property")
	}

	plugin, err = p.kongPluginFromK8SPlugin(*k8sPlugin)
	return plugin, err
}

func (p *Parser) secretToConfiguration(
	reference configurationv1.SecretValueFromSource, namespace string) (
	configurationv1.Configuration, error) {
	secret, err := p.store.GetSecret(namespace, reference.Secret)
	if err != nil {
		return configurationv1.Configuration{}, fmt.Errorf(
			"error fetching plugin configuration secret '%v/%v': %v",
			namespace, reference.Secret, err)
	}
	secretVal, ok := secret.Data[reference.Key]
	if !ok {
		return configurationv1.Configuration{},
			fmt.Errorf("no key '%v' in secret '%v/%v'",
				reference.Key, namespace, reference.Secret)
	}
	var config configurationv1.Configuration
	if err := json.Unmarshal(secretVal, &config); err != nil {
		if err := yaml.Unmarshal(secretVal, &config); err != nil {
			return configurationv1.Configuration{},
				fmt.Errorf("key '%v' in secret '%v/%v' contains neither "+
					"valid JSON nor valid YAML)",
					reference.Key, namespace, reference.Secret)
		}
	}
	return config, nil
}

func (p *Parser) namespacedSecretToConfiguration(
	reference configurationv1.NamespacedSecretValueFromSource) (
	configurationv1.Configuration, error) {
	bareReference := configurationv1.SecretValueFromSource{
		Secret: reference.Secret,
		Key:    reference.Key}
	return p.secretToConfiguration(bareReference, reference.Namespace)
}

// plugin is a intermediate type to hold plugin related configuration
type plugin struct {
	Name   string
	Config configurationv1.Configuration

	RunOn     string
	Disabled  bool
	Protocols []string
}

func toKongPlugin(plugin plugin) kong.Plugin {
	result := kong.Plugin{
		Name:   kong.String(plugin.Name),
		Config: kong.Configuration(plugin.Config).DeepCopy(),
	}
	if plugin.RunOn != "" {
		result.RunOn = kong.String(plugin.RunOn)
	}
	if plugin.Disabled {
		result.Enabled = kong.Bool(false)
	}
	if len(plugin.Protocols) > 0 {
		result.Protocols = kong.StringSlice(plugin.Protocols...)
	}
	return result
}

func (p *Parser) kongPluginFromK8SClusterPlugin(
	k8sPlugin configurationv1.KongClusterPlugin) (kong.Plugin, error) {
	config := k8sPlugin.Config
	if k8sPlugin.ConfigFrom.SecretValue !=
		(configurationv1.NamespacedSecretValueFromSource{}) &&
		len(k8sPlugin.Config) > 0 {
		return kong.Plugin{},
			fmt.Errorf("KongClusterPlugin '/%v' has both "+
				"Config and ConfigFrom set", k8sPlugin.Name)
	}
	if k8sPlugin.ConfigFrom.SecretValue != (configurationv1.
		NamespacedSecretValueFromSource{}) {
		var err error
		config, err = p.namespacedSecretToConfiguration(
			k8sPlugin.ConfigFrom.SecretValue)
		if err != nil {
			return kong.Plugin{},
				fmt.Errorf("error parsing config for KongClusterPlugin %v: %w",
					k8sPlugin.Name, err)
		}
	}
	kongPlugin := toKongPlugin(plugin{
		Name:   k8sPlugin.PluginName,
		Config: config,

		RunOn:     k8sPlugin.RunOn,
		Disabled:  k8sPlugin.Disabled,
		Protocols: k8sPlugin.Protocols,
	})
	return kongPlugin, nil
}

func (p *Parser) kongPluginFromK8SPlugin(
	k8sPlugin configurationv1.KongPlugin) (kong.Plugin, error) {
	config := k8sPlugin.Config
	if k8sPlugin.ConfigFrom.SecretValue !=
		(configurationv1.SecretValueFromSource{}) &&
		len(k8sPlugin.Config) > 0 {
		return kong.Plugin{},
			fmt.Errorf("KongPlugin '%v/%v' has both "+
				"Config and ConfigFrom set",
				k8sPlugin.Namespace, k8sPlugin.Name)
	}
	if k8sPlugin.ConfigFrom.SecretValue !=
		(configurationv1.SecretValueFromSource{}) {
		var err error
		config, err = p.secretToConfiguration(
			k8sPlugin.ConfigFrom.SecretValue, k8sPlugin.Namespace)
		if err != nil {
			return kong.Plugin{},
				fmt.Errorf("error parsing config for KongPlugin '%v/%v': %w",
					k8sPlugin.Name, k8sPlugin.Namespace, err)
		}
	}
	kongPlugin := toKongPlugin(plugin{
		Name:   k8sPlugin.PluginName,
		Config: config,

		RunOn:     k8sPlugin.RunOn,
		Disabled:  k8sPlugin.Disabled,
		Protocols: k8sPlugin.Protocols,
	})
	return kongPlugin, nil
}

// getEndpoints returns a list of <endpoint ip>:<port> for a given service/target port combination.
func (p *Parser) getEndpoints(
	s *corev1.Service,
	port *corev1.ServicePort,
	proto corev1.Protocol,
	getEndpoints func(string, string) (*corev1.Endpoints, error),
) []utils.Endpoint {

	upsServers := []utils.Endpoint{}

	if s == nil || port == nil {
		return upsServers
	}

	logger := p.Logger.WithFields(logrus.Fields{
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
		logger.Debug("found service of type=ExternalName")

		targetPort := port.TargetPort.IntValue()
		// check for invalid port value
		if targetPort <= 0 {
			logger.Errorf("invalid service: invalid port: %v", targetPort)
			return upsServers
		}

		return append(upsServers, utils.Endpoint{
			Address: s.Spec.ExternalName,
			Port:    fmt.Sprintf("%v", targetPort),
		})
	}
	if annotations.HasServiceUpstreamAnnotation(s.Annotations) {
		return append(upsServers, utils.Endpoint{
			Address: s.Name + "." + s.Namespace + ".svc",
			Port:    fmt.Sprintf("%v", port.Port),
		})

	}

	logger.Debugf("fetching endpoints")
	ep, err := getEndpoints(s.Namespace, s.Name)
	if err != nil {
		logger.Errorf("failed to fetch endpoints: %v", err)
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
				ups := utils.Endpoint{
					Address: epAddress.IP,
					Port:    fmt.Sprintf("%v", targetPort),
				}
				upsServers = append(upsServers, ups)
				adus[ep] = true
			}
		}
	}

	logger.Debugf("found endpoints: %v", upsServers)
	return upsServers
}
