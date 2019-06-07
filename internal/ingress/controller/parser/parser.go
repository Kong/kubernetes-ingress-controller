package parser

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"strconv"

	"github.com/golang/glog"
	"github.com/hbagdi/go-kong/kong"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/internal/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/annotations"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/store"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/utils"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Route represents a Kong Route and holds a reference to the Ingress
// rule.
type Route struct {
	kong.Route
	// Ingress object associated with this route
	Ingress extensions.Ingress
	Plugins []kong.Plugin
}

// Service represents a service in Kong and holds routes associated with the
// service and other k8s metadata.
type Service struct {
	kong.Service
	Backend   extensions.IngressBackend
	Namespace string
	Routes    []Route
	Plugins   []kong.Plugin
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
	Plugins []kong.Plugin
	// Credentials type(key-auth, basic-auth) to credential mapping
	Credentials map[string][]map[string]interface{}
}

// KongState holds the configuration that should be applied to Kong.
type KongState struct {
	Services      []Service
	Upstreams     []Upstream
	Certificates  []Certificate
	GlobalPlugins []Plugin
	Consumers     []Consumer
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
	// parse ingress rules
	parsedInfo, err := p.parseIngressRules(ings)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing ingress rules")
	}

	// add the routes and services to the state
	for _, service := range parsedInfo.ServiceNameToServices {
		state.Services = append(state.Services, service)
	}

	// generate Upstreams and Targets from service defs
	// TODO add support for DNS based service records
	state.Upstreams, err = p.getUpstreams(parsedInfo.ServiceNameToServices)
	if err != nil {
		return nil, errors.Wrap(err, "building upstreams and targets")
	}

	// merge KongIngress with Routes, Services and Upstream
	err = p.fillOverrides(state)
	if err != nil {
		return nil, errors.Wrap(err, "overriding KongIngress values")
	}

	// TODO add processors for annotations on Ingress object

	// process annotation plugins
	err = p.fillPlugins(state)

	// generate Certificates and SNIs
	state.Certificates, err = p.getCerts(parsedInfo.SecretNameToSNIs)
	if err != nil {
		return nil, err
	}

	// Global plugins
	state.GlobalPlugins, err = p.globalPlugins()
	if err != nil {
		return nil, err
	}

	// TODO add support for consumers and credentials

	return &state, nil
}

func (p *Parser) parseIngressRules(
	ingressList []*extensions.Ingress) (*parsedIngressRules, error) {

	// generate the following:
	// Services and Routes
	var allDefaultBackends []extensions.Ingress
	secretNameToSNIs := make(map[string][]string)
	serviceNameToServices := make(map[string]Service)

	for i := 0; i < len(ingressList); i++ {
		ingress := *ingressList[i]
		ingressSpec := ingress.Spec

		if ingressSpec.Backend != nil {
			allDefaultBackends = append(allDefaultBackends, ingress)

		}

		for _, tls := range ingressSpec.TLS {
			if len(tls.Hosts) == 0 {
				continue
			}
			if tls.SecretName == "" {
				continue
			}
			hosts := tls.Hosts
			secretName := ingress.Namespace + "/" + tls.SecretName
			if secretNameToSNIs[secretName] != nil {
				hosts = append(hosts, secretNameToSNIs[secretName]...)
			}
			secretNameToSNIs[secretName] = hosts
		}

		for i, rule := range ingressSpec.Rules {
			host := rule.Host
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
							Name:           kong.String(serviceName),
							Host:           kong.String(rule.Backend.ServiceName + "." + ingress.Namespace + ".svc"),
							Port:           kong.Int(80),
							Protocol:       kong.String("http"),
							Path:           kong.String("/"),
							ConnectTimeout: kong.Int(60000),
							ReadTimeout:    kong.Int(60000),
							WriteTimeout:   kong.Int(60000),
							Retries:        kong.Int(5),
						},
						Namespace: ingress.Namespace,
						Backend:   rule.Backend,
					}
				}
				service.Routes = append(service.Routes, r)
				serviceNameToServices[serviceName] = service
			}
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
					Host: kong.String(defaultBackend.ServiceName + "." + ingress.Namespace + ".svc"),
					Port: kong.Int(80),
				},
				Namespace: ingress.Namespace,
				Backend:   *defaultBackend,
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
		kongIngress, err := p.getKongIngressForService(
			state.Services[i].Namespace, state.Services[i].Backend.ServiceName)
		if err != nil {
			return errors.Wrapf(err, "fetching KongIngress for service '%v' in namespace '%v'",
				state.Services[i].Backend.ServiceName, state.Services[i].Namespace)
		}
		overrideService(&state.Services[i], kongIngress)

		// Routes
		for j := 0; j < len(state.Services[i].Routes); j++ {
			kongIngress, err := p.getKongIngressFromIngress(
				&state.Services[i].Routes[j].Ingress)
			if err != nil {
				return errors.Wrapf(err, "fetching KongIngress for Ingress '%v' in namespace '%v'",
					state.Services[i].Routes[j].Ingress.Name, state.Services[i].Routes[j].Ingress.Namespace)
			}
			overrideRoute(&state.Services[i].Routes[j], kongIngress)
		}
	}

	// Upstreams
	for i := 0; i < len(state.Upstreams); i++ {
		kongIngress, err := p.getKongIngressForService(
			state.Upstreams[i].Service.Namespace, state.Upstreams[i].Service.Backend.ServiceName)
		if err != nil {
			return errors.Wrapf(err, "fetching KongIngress for service '%v' in namespace '%v'",
				state.Upstreams[i].Service.Backend.ServiceName, state.Upstreams[i].Service.Namespace)
		}
		overrideUpstream(&state.Upstreams[i], kongIngress)
	}
	return nil
}

func overrideService(service *Service,
	kongIngress *configurationv1.KongIngress) {
	if service == nil || kongIngress == nil || kongIngress.Proxy == nil {
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

func overrideRoute(route *Route,
	kongIngress *configurationv1.KongIngress) {
	if route == nil || kongIngress == nil || kongIngress.Route == nil {
		return
	}
	r := kongIngress.Route

	if len(r.Methods) != 0 {
		route.Methods = cloneStringPointerSlice(r.Methods...)
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
		upstreamName := service.Backend.ServiceName + "." + service.Namespace + ".svc"
		upstream := Upstream{
			Upstream: kong.Upstream{
				Name: kong.String(upstreamName),
			},
			Service: service,
		}
		svcKey := service.Namespace + "/" + service.Backend.ServiceName
		targets, err := p.getServiceEndpoints(service.Namespace,
			service.Backend.ServiceName,
			service.Backend.ServicePort.String())
		if err != nil {
			return nil, errors.Wrapf(err,
				"error getting endpoints for '%v' service", svcKey)
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

	return cert, key, nil
}

func (p *Parser) getCerts(secretsToSNIs map[string][]string) ([]Certificate,
	error) {
	var res []Certificate
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
		kongCert := Certificate{
			Certificate: kong.Certificate{
				Cert: kong.String(cert),
				Key:  kong.String(key),
				SNIs: kong.StringSlice(SNIs...),
			},
		}
		res = append(res, kongCert)
	}
	return res, nil
}

func (p *Parser) fillPlugins(state KongState) error {
	for i := range state.Services {
		// service
		svcKey := state.Services[i].Namespace + "/" +
			state.Services[i].Backend.ServiceName
		svc, err := p.store.GetService(state.Services[i].Namespace,
			state.Services[i].Backend.ServiceName)
		if err != nil {
			return errors.Wrapf(err, "fetching service '%s'", svcKey)
		}
		plugins, err := p.getPluginsFromAnnotations(state.Services[i].Namespace,
			svc.GetAnnotations())
		if err != nil {
			return errors.Wrapf(err, "fetching KongPlugins for service '%s'", svcKey)
		}
		state.Services[i].Plugins = plugins

		// route
		for j := range state.Services[i].Routes {
			plugins, err := p.getPluginsFromAnnotations(state.Services[i].Routes[j].Ingress.Namespace, state.Services[i].Routes[j].Ingress.GetAnnotations())
			if err != nil {
				return errors.Wrapf(err, "fetching KongPlugins for a route in Ingress '%s'", svcKey)
			}
			state.Services[i].Routes[j].Plugins = plugins
		}
	}
	return nil
}

func (p *Parser) globalPlugins() ([]Plugin, error) {
	globalPlugins, err := p.store.ListGlobalKongPlugins()
	if err != nil {
		return nil, errors.Wrap(err, "error listing global plugins:")
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
			Plugin: kong.Plugin{
				Name:   kong.String(pluginName),
				Config: kong.Configuration(k8sPlugin.Config),
			},
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

func (p *Parser) getServiceEndpoints(namespace, svcName string,
	backendPort string) ([]Target, error) {
	var targets []Target
	var endpoints []utils.Endpoint
	var servicePort corev1.ServicePort
	svc, err := p.store.GetService(namespace, svcName)
	svcKey := namespace + "/" + svcName
	if err != nil {
		return nil, errors.Wrapf(err,
			"error getting service '%v' from the cache", svcKey)
	}
	glog.V(3).Infof("obtaining port information for service %v", svcKey)

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

	endpoints = getEndpoints(svc, &servicePort,
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

func (p *Parser) getKongIngressForService(namespace, serviceName string) (
	*configurationv1.KongIngress, error) {
	svc, err := p.store.GetService(namespace, serviceName)
	if err != nil {
		return nil, errors.Wrapf(err, "fetching service '%s' from cache",
			namespace+"/"+serviceName)
	}
	confName := annotations.ExtractConfigurationName(svc.Annotations)
	if confName != "" {
		ki, err := p.store.GetKongIngress(svc.Namespace, confName)
		if err == nil {
			return ki, nil
		}
	}
	return nil, nil
}

// getKongIngress checks if the Ingress contains an annotation for configuration
// or if exists a KongIngress object with the same name than the Ingress
func (p *Parser) getKongIngressFromIngress(ing *extensions.Ingress) (
	*configurationv1.KongIngress, error) {
	confName := annotations.ExtractConfigurationName(ing.Annotations)
	if confName != "" {
		ki, err := p.store.GetKongIngress(ing.Namespace, confName)
		if err == nil {
			return ki, nil
		}
	}

	ki, err := p.store.GetKongIngress(ing.Namespace, ing.Name)
	if err == nil {
		return ki, err
	}
	return nil, nil
}

// getPluginsFromAnnotations extracts plugins to be applied on an ingress/service from annotations
func (p *Parser) getPluginsFromAnnotations(namespace string, anns map[string]string) ([]kong.Plugin, error) {
	pluginsInk8s := make(map[string]*configurationv1.KongPlugin)
	pluginList := annotations.ExtractKongPluginsFromAnnotations(anns)
	// override plugins configured by new annotation
	for _, plugin := range pluginList {
		k8sPlugin, err := p.store.GetKongPlugin(namespace, plugin)
		if err != nil {
			return nil, errors.Wrapf(err, "fetching KongPlugin %v", plugin)
		}
		// ignore plugins with no name
		if k8sPlugin.PluginName == "" {
			glog.Errorf("KongPlugin Custom resource '%v' has no `plugin` property, the plugin will not be configured", k8sPlugin.Name)
			continue
		}
		pluginsInk8s[k8sPlugin.PluginName] = k8sPlugin
	}

	var plugins []kong.Plugin
	for _, p := range pluginsInk8s {
		plugin := kong.Plugin{
			Name:   kong.String(p.PluginName),
			Config: kong.Configuration(p.Config).DeepCopy(),
		}
		if p.RunOn != "" {
			plugin.RunOn = kong.String(p.RunOn)
		}
		if p.Disabled {
			plugin.Enabled = kong.Bool(false)
		}
		if len(p.Protocols) > 0 {
			plugin.Protocols = kong.StringSlice(p.Protocols...)
		}
		plugins = append(plugins, plugin)
	}
	return plugins, nil
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
