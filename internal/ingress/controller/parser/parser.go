package parser

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"strconv"

	"github.com/golang/glog"
	"github.com/hbagdi/go-kong/kong"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/internal/apis/configuration/v1"
	configurationv1beta1 "github.com/kong/kubernetes-ingress-controller/internal/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/annotations"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/store"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/sets"
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

// Consumer holds a Kong consumer and it's plugins and credentials.
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
	Services     []Service
	Upstreams    []Upstream
	Certificates []Certificate
	Plugins      []Plugin
	Consumers    []Consumer
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
	store store.Storer
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

// New returns a new parser backed with store.
func New(store store.Storer) Parser {
	return Parser{store: store}
}

// Build creates a Kong configuration from Ingress and Custom resources
// defined in Kuberentes.
// It throws an error if there is an error returned from client-go.
func (p *Parser) Build() (*KongState, error) {
	var state KongState
	ings := p.store.ListIngresses()
	tcpIngresses, err := p.store.ListTCPIngresses()
	if err != nil {
		glog.Errorf("error listing TCPIngresses: %v", err)
	}
	// parse ingress rules
	parsedInfo, err := p.parseIngressRules(ings, tcpIngresses)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing ingress rules")
	}

	// populate Kubernetes Service
	for key, service := range parsedInfo.ServiceNameToServices {
		k8sSvc, err := p.store.GetService(service.Namespace, service.Backend.Name)
		if err != nil {
			glog.Errorf("getting service: %v", err)
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
				glog.Errorf("getting secret: %v: %v", secretKey, err)
			}
		}
		parsedInfo.ServiceNameToServices[key] = service
	}

	// add the routes and services to the state
	for _, service := range parsedInfo.ServiceNameToServices {
		state.Services = append(state.Services, service)
	}

	// generate Upstreams and Targets from service defs
	state.Upstreams, err = p.getUpstreams(parsedInfo.ServiceNameToServices)
	if err != nil {
		return nil, errors.Wrap(err, "building upstreams and targets")
	}

	// merge KongIngress with Routes, Services and Upstream
	err = p.fillOverrides(state)
	if err != nil {
		return nil, errors.Wrap(err, "overriding KongIngress values")
	}

	// generate consumers and credentials
	err = p.fillConsumersAndCredentials(&state)
	if err != nil {
		return nil, errors.Wrap(err, "building consumers and credentials")
	}

	// process annotation plugins
	state.Plugins = p.fillPlugins(state)

	// generate Certificates and SNIs
	state.Certificates, err = p.getCerts(parsedInfo.SecretNameToSNIs)
	if err != nil {
		return nil, err
	}

	return &state, nil
}

func processCredential(credType string, consumer *Consumer,
	credConfig interface{}) error {
	switch credType {
	case "key-auth", "keyauth_credential":
		var cred kong.KeyAuth
		err := decodeCredential(credConfig, &cred)
		if err != nil {
			return errors.Wrap(err, "failed to decode key-auth credential")

		}
		consumer.KeyAuths = append(consumer.KeyAuths, &cred)
	case "basic-auth", "basicauth_credential":
		var cred kong.BasicAuth
		err := decodeCredential(credConfig, &cred)
		if err != nil {
			return errors.Wrap(err, "failed to decode basic-auth credential")
		}
		consumer.BasicAuths = append(consumer.BasicAuths, &cred)
	case "hmac-auth", "hmacauth_credential":
		var cred kong.HMACAuth
		err := decodeCredential(credConfig, &cred)
		if err != nil {
			return errors.Wrap(err, "failed to decode hmac-auth credential")
		}
		consumer.HMACAuths = append(consumer.HMACAuths, &cred)
	case "oauth2":
		var cred kong.Oauth2Credential
		err := decodeCredential(credConfig, &cred)
		if err != nil {
			return errors.Wrap(err, "failed to decode oauth2 credential")
		}
		consumer.Oauth2Creds = append(consumer.Oauth2Creds, &cred)
	case "jwt", "jwt_secret":
		var cred kong.JWTAuth
		err := decodeCredential(credConfig, &cred)
		if err != nil {
			glog.Error("failed to process JWT credential", err)
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
			glog.Error("failed to process ACL group", err)
		}
		consumer.ACLGroups = append(consumer.ACLGroups, &cred)
	default:
		return errors.Errorf("invalid credential type: '%v'", credType)
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
		return errors.Wrap(err, "failed to create a decoder")
	}
	err = decoder.Decode(credConfig)
	if err != nil {
		return errors.Wrapf(err, "failed to decode credential")
	}
	return nil
}

func (p *Parser) fillConsumersAndCredentials(state *KongState) error {
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

		for _, cred := range consumer.Credentials {
			secret, err := p.store.GetSecret(consumer.Namespace, cred)
			if err != nil {
				glog.Errorf("error fetching credential secret '%v/%v': %v",
					consumer.Namespace, cred, err)
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
				glog.Errorf("invalid credType in secret '%v/%v'",
					consumer.Namespace, cred)
			}
			if !supportedCreds.Has(credType) {
				glog.Errorf("invalid credType '%v' in secret '%v/%v'",
					credType, consumer.Namespace, cred)
				continue
			}
			if len(credConfig) <= 1 { // 1 key of credType itself
				glog.Errorf("empty config in '%v' in secret '%v/%v'",
					credType, consumer.Namespace, cred)
				continue
			}
			err = processCredential(credType, &c, credConfig)
			if err != nil {
				glog.Errorf("failed to process credential in "+
					"secret '%v/%v': %v", consumer.Namespace, cred, err)
				continue
			}
		}

		consumerIndex[consumer.Namespace+"/"+consumer.Name] = c
	}

	// legacy attach credentials
	credentials := p.store.ListKongCredentials()
	if len(credentials) > 0 {
		glog.Warningf("Deprecated KongCredential in use, " +
			"please use secret-based credentials. " +
			"KongCredential resource will be removed in future.")
	}
	for _, credential := range credentials {
		consumer, ok := consumerIndex[credential.Namespace+"/"+
			credential.ConsumerRef]
		if !ok {
			continue
		}
		if credential.Type == "" {
			glog.Errorf("empty credential type in KongCredential '%v/%v'",
				credential.Namespace, credential.Name)
			continue
		}
		if !supportedCreds.Has(credential.Type) {
			glog.Errorf("invalid credType '%v' in KongCredential '%v/%v'",
				credential.Type, credential.Namespace, credential.Name)
			continue
		}
		if credential.Config == nil {
			glog.Errorf("empty config in '%v' in KongCredential '%v/%v'",
				credential.Type, credential.Namespace, credential.Name)
			continue
		}
		err := processCredential(credential.Type, &consumer, credential.Config)
		if err != nil {
			glog.Errorf("failed to process credential in "+
				"KongCredential '%v/%v': %v", credential.Namespace,
				credential.Name, err)
			continue
		}
		consumerIndex[credential.Namespace+"/"+credential.ConsumerRef] = consumer
	}

	// populate the consumer in the state
	for _, c := range consumerIndex {
		state.Consumers = append(state.Consumers, c)
	}
	return nil
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

func toNetworkingTLS(tls []configurationv1beta1.IngressTLS) []networking.IngressTLS {
	var result []networking.IngressTLS

	for _, t := range tls {
		result = append(result, networking.IngressTLS{
			Hosts:      t.Hosts,
			SecretName: t.SecretName,
		})
	}
	return result
}

func (p *Parser) parseIngressRules(
	ingressList []*networking.Ingress,
	tcpIngressList []*configurationv1beta1.TCPIngress) (*parsedIngressRules, error) {

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

				isACMEChallenge := strings.HasPrefix(path, "/.well-known/acme-challenge/")

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
						StripPath:     kong.Bool(!isACMEChallenge),
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

		processTLSSections(toNetworkingTLS(ingressSpec.TLS),
			ingress.Namespace, secretNameToSNIs)

		for i, rule := range ingressSpec.Rules {

			if rule.Port <= 0 {
				glog.Errorf("invalid port value (%v) in TCPIngress %v/%v",
					rule.Port, ingress.Namespace, ingress.Name)
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
				glog.Errorf("invalid empty serviceName in"+
					"TCPIngress %v/%v", ingress.Namespace, ingress.Name)
				continue
			}
			if rule.Backend.ServicePort <= 0 {
				glog.Errorf("invalid servicePort (%v) in"+
					"TCPIngress %v/%v", rule.Backend.ServicePort,
					ingress.Namespace, ingress.Name)
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
				StripPath:     kong.Bool(true),
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
	}, nil
}

func (p *Parser) fillOverrides(state KongState) error {
	for i := 0; i < len(state.Services); i++ {
		// Services
		anns := state.Services[i].K8sService.Annotations
		kongIngress, err := p.getKongIngressForService(
			state.Services[i].K8sService)
		if err != nil {
			glog.Errorf("error getting kongIngress %v", err)
		}
		overrideService(&state.Services[i], kongIngress, anns)

		// Routes
		for j := 0; j < len(state.Services[i].Routes); j++ {
			var kongIngress *configurationv1.KongIngress
			var err error
			if state.Services[i].Routes[j].IsTCP {
				kongIngress, err = p.getKongIngressFromTCPIngress(
					&state.Services[i].Routes[j].TCPIngress)
			} else {
				kongIngress, err = p.getKongIngressFromIngress(
					&state.Services[i].Routes[j].Ingress)
			}

			if err != nil {
				glog.Errorf("error getting kongIngress %v", err)
			}
			overrideRoute(&state.Services[i].Routes[j], kongIngress)
		}
	}

	// Upstreams
	for i := 0; i < len(state.Upstreams); i++ {
		kongIngress, err := p.getKongIngressForService(
			state.Upstreams[i].Service.K8sService)
		if err == nil {
			overrideUpstream(&state.Upstreams[i], kongIngress)
		} else {
			glog.Error(errors.Wrapf(err, "fetching KongIngress for service '%v' in namespace '%v'",
				state.Upstreams[i].Service.Backend.Name, state.Upstreams[i].Service.Namespace))
		}
	}
	return nil
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

// overrideServiceByAnnotation sets the Service protocol via annotation
func overrideServiceByAnnotation(service *Service,
	anns map[string]string) {
	protocol := annotations.ExtractProtocolName(anns)
	if protocol == "" || validateProtocol(protocol) != true {
		return
	}
	service.Protocol = kong.String(protocol)
}

// overrideService sets Service fields by KongIngress first, then by annotation
func overrideService(service *Service,
	kongIngress *configurationv1.KongIngress,
	anns map[string]string) {
	if service == nil {
		return
	}
	overrideServiceByKongIngress(service, kongIngress)
	overrideServiceByAnnotation(service, anns)

	if *service.Protocol == "grpc" || *service.Protocol == "grpcs" {
		// grpc(s) doesn't accept a path
		service.Path = nil
	}
}

// overrideRouteByKongIngress sets Route fields by KongIngress
func overrideRouteByKongIngress(route *Route,
	kongIngress *configurationv1.KongIngress) {
	if kongIngress == nil || kongIngress.Route == nil {
		return
	}

	r := kongIngress.Route
	if len(r.Methods) != 0 {
		route.Methods = cloneStringPointerSlice(r.Methods...)
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
		if validateProtocol(*protocol) != true {
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

// overrideRouteByAnnotation sets Route protocols via annotation
func overrideRouteByAnnotation(route *Route, anns map[string]string) {
	if anns == nil {
		return
	}
	protocols := annotations.ExtractProtocolNames(anns)
	var prots []*string
	for _, prot := range protocols {
		if validateProtocol(prot) != true {
			return
		}
		prots = append(prots, kong.String(prot))
	}

	route.Protocols = prots
}

// overrideRoute sets Route fields by KongIngress first, then by annotation
func overrideRoute(route *Route,
	kongIngress *configurationv1.KongIngress) {
	if route == nil {
		return
	}
	overrideRouteByKongIngress(route, kongIngress)
	anns := route.Ingress.Annotations
	if route.IsTCP {
		anns = route.TCPIngress.Annotations
	}
	overrideRouteByAnnotation(route, anns)
	normalizeProtocols(route)
	for _, val := range route.Protocols {
		if *val == "grpc" || *val == "grpcs" {
			// grpc(s) doesn't accept strip_path
			route.StripPath = nil
			break
		}
	}
}

func cloneStringPointerSlice(array ...*string) []*string {
	var res []*string
	for _, s := range array {
		res = append(res, &(*s))
	}
	return res
}

func overrideUpstream(upstream *Upstream,
	kongIngress *configurationv1.KongIngress) {
	if kongIngress == nil || kongIngress.Upstream == nil || upstream == nil {
		return
	}
	// name is the only field that is set
	name := *upstream.Upstream.Name
	upstream.Upstream = *kongIngress.Upstream.DeepCopy()
	upstream.Name = &name
}

func (p *Parser) getUpstreams(serviceMap map[string]Service) ([]Upstream, error) {
	var upstreams []Upstream
	for _, service := range serviceMap {
		upstreamName := service.Backend.Name + "." + service.Namespace + "." + service.Backend.Port.String() + ".svc"
		upstream := Upstream{
			Upstream: kong.Upstream{
				Name: kong.String(upstreamName),
			},
			Service: service,
		}
		svcKey := service.Namespace + "/" + service.Backend.Name
		targets, err := p.getServiceEndpoints(service.K8sService,
			service.Backend.Port.String())
		if err != nil {
			glog.Errorf("error getting endpoints for '%v' service: %v",
				svcKey, err)
		}
		upstream.Targets = targets
		upstreams = append(upstreams, upstream)
	}
	return upstreams, nil
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

func (p *Parser) getCerts(secretsToSNIs map[string][]string) ([]Certificate,
	error) {
	snisAdded := make(map[string]bool)
	// map of cert public key + private key to certificate
	certs := make(map[string]Certificate)

	for secretKey, SNIs := range secretsToSNIs {
		namespaceName := strings.Split(secretKey, "/")
		secret, err := p.store.GetSecret(namespaceName[0], namespaceName[1])
		if err != nil {
			glog.Errorf("error fetching certificate '%v': %v", secretKey, err)
			continue
		}
		cert, key, err := getCertFromSecret(secret)
		if err != nil {
			glog.Errorf("error finding a certificate in '%v': %v",
				secretKey, err)
			continue
		}
		kongCert, ok := certs[cert+key]
		if !ok {
			kongCert = Certificate{
				Certificate: kong.Certificate{
					ID:   kong.String(string(secret.UID)),
					Cert: kong.String(cert),
					Key:  kong.String(key),
				},
			}
		}

		for _, sni := range SNIs {
			if !snisAdded[sni] {
				snisAdded[sni] = true
				kongCert.SNIs = append(kongCert.SNIs, kong.String(sni))
			}
		}
		certs[cert+key] = kongCert
	}
	var res []Certificate
	for _, cert := range certs {
		res = append(res, cert)
	}
	return res, nil
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
			glog.Errorf("fetching KongPlugin '%v/%v': %v", namespace,
				kongPluginName, err)
			continue
		}

		for _, rel := range getCombinations(relations) {
			plugin := *plugin.DeepCopy()
			// ID is populated because that is read by decK and in_memory
			// translater too
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
		glog.Errorf("fetching global plugins: %v", err)
	}
	plugins = append(plugins, globalPlugins...)

	return plugins
}

func (p *Parser) globalPlugins() ([]Plugin, error) {
	globalPlugins, err := p.store.ListGlobalKongPlugins()
	if err != nil {
		return nil, errors.Wrap(err, "error listing global KongPlugins:")
	}
	res := make(map[string]Plugin)
	var duplicates []string // keep track of duplicate
	// TODO respect the oldest CRD
	// Current behavior is to skip creating the plugin but in case
	// of duplicate plugin definitions, we should respect the oldest one
	// This is important since if a user comes in to k8s and creates a new
	// CRD, the user now deleted an older plugin

	for i := 0; i < len(globalPlugins); i++ {
		k8sPlugin := *globalPlugins[i]
		pluginName := k8sPlugin.PluginName
		// empty pluginName skip it
		if pluginName == "" {
			glog.Errorf("KongPlugin '%v' does not specify a plugin name",
				k8sPlugin.Name)
			continue
		}
		if _, ok := res[pluginName]; ok {
			glog.Error("Multiple KongPlugin definitions found with"+
				" 'global' annotation for '", pluginName,
				"', the plugin will not be applied")
			duplicates = append(duplicates, pluginName)
			continue
		}
		res[pluginName] = Plugin{
			Plugin: kongPluginFromK8SPlugin(k8sPlugin),
		}
	}

	globalClusterPlugins, err := p.store.ListGlobalKongClusterPlugins()
	if err != nil {
		return nil, errors.Wrap(err, "error listing global KongClusterPlugins")
	}
	for i := 0; i < len(globalClusterPlugins); i++ {
		k8sPlugin := *globalClusterPlugins[i]
		pluginName := k8sPlugin.PluginName
		// empty pluginName skip it
		if pluginName == "" {
			glog.Errorf("KongPlugin '%v' does not specify a plugin name",
				k8sPlugin.Name)
			continue
		}
		if _, ok := res[pluginName]; ok {
			glog.Error("Multiple KongPlugin definitions found with"+
				" 'global' annotation for '", pluginName,
				"', the plugin will not be applied")
			duplicates = append(duplicates, pluginName)
			continue
		}
		res[pluginName] = Plugin{
			Plugin: kongPluginFromK8SClusterPlugin(k8sPlugin),
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
	backendPort string) ([]Target, error) {
	var targets []Target
	var endpoints []utils.Endpoint
	var servicePort corev1.ServicePort
	svcKey := svc.Namespace + "/" + svc.Name

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
		externalPort, err := strconv.Atoi(backendPort)
		if err != nil {
			glog.Warningf("only numeric ports are allowed in"+
				" ExternalName services: %v is not valid as a TCP/UDP port",
				backendPort)
			return targets, nil
		}

		servicePort = corev1.ServicePort{
			Protocol:   "TCP",
			Port:       int32(externalPort),
			TargetPort: intstr.FromString(backendPort),
		}
	}

	endpoints = getEndpoints(&svc, &servicePort,
		corev1.ProtocolTCP, p.store.GetEndpointsForService)
	if len(endpoints) == 0 {
		glog.Warningf("service %v does not have any active endpoints",
			svcKey)
	}
	for _, endpoint := range endpoints {
		target := Target{
			Target: kong.Target{
				Target: kong.String(endpoint.Address + ":" + endpoint.Port),
			},
		}
		targets = append(targets, target)
	}
	return targets, nil
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
		return ki, err
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
				return plugin, errors.Errorf("invalid empty 'plugin' property")
			}
			plugin = kongPluginFromK8SClusterPlugin(*clusterPlugin)
			return plugin, err
		}
		// handle other errors
		if err != nil {
			return plugin, err
		}
	}
	// ignore plugins with no name
	if k8sPlugin.PluginName == "" {
		return plugin, errors.Errorf("invalid empty 'plugin' property")
	}
	plugin = kongPluginFromK8SPlugin(*k8sPlugin)
	return plugin, nil
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
		Config: kong.Configuration(*plugin.Config.DeepCopy()),
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

func kongPluginFromK8SClusterPlugin(k8sPlugin configurationv1.KongClusterPlugin) kong.Plugin {
	return toKongPlugin(plugin{
		Name:   k8sPlugin.PluginName,
		Config: k8sPlugin.Config,

		RunOn:     k8sPlugin.RunOn,
		Disabled:  k8sPlugin.Disabled,
		Protocols: k8sPlugin.Protocols,
	})
}

func kongPluginFromK8SPlugin(k8sPlugin configurationv1.KongPlugin) kong.Plugin {
	return toKongPlugin(plugin{
		Name:   k8sPlugin.PluginName,
		Config: k8sPlugin.Config,

		RunOn:     k8sPlugin.RunOn,
		Disabled:  k8sPlugin.Disabled,
		Protocols: k8sPlugin.Protocols,
	})
}

// getEndpoints returns a list of <endpoint ip>:<port> for a given service/target port combination.
func getEndpoints(
	s *corev1.Service,
	port *corev1.ServicePort,
	proto corev1.Protocol,
	getEndpoints func(string, string) (*corev1.Endpoints, error),
) []utils.Endpoint {

	upsServers := []utils.Endpoint{}

	if s == nil || port == nil {
		return upsServers
	}

	// avoid duplicated upstream servers when the service
	// contains multiple port definitions sharing the same
	// targetport.
	adus := make(map[string]bool)

	// ExternalName services
	if s.Spec.Type == corev1.ServiceTypeExternalName {
		glog.V(3).Infof("Ingress using a service %v of type=ExternalName", s.Name)

		targetPort := port.TargetPort.IntValue()
		// check for invalid port value
		if targetPort <= 0 {
			glog.Errorf("ExternalName service with an invalid port: %v", targetPort)
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

	glog.V(3).Infof("getting endpoints for service %v/%v and port %v", s.Namespace, s.Name, port.String())
	ep, err := getEndpoints(s.Namespace, s.Name)
	if err != nil {
		glog.Warningf("unexpected error obtaining service endpoints: %v", err)
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

	glog.V(3).Infof("endpoints found: %v", upsServers)
	return upsServers
}
