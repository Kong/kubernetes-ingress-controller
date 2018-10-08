/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/golang/glog"

	apiv1 "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/sets"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	kong "github.com/kong/kubernetes-ingress-controller/internal/apis/admin"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress"
)

const (
	defUpstreamName = "upstream-default-backend"
	defServerName   = "_"
	rootLocation    = "/"
)

type Kong struct {
	URL    string
	Client *kong.RestClient
}

// Configuration contains all the settings required by an Ingress controller
type Configuration struct {
	Kong

	APIServerHost  string
	KubeConfigFile string
	KubeClient     clientset.Interface
	KubeConf       *rest.Config

	ResyncPeriod time.Duration

	DefaultService string

	Namespace string

	ForceNamespaceIsolation bool

	DefaultHealthzURL     string
	DefaultSSLCertificate string

	// optional
	PublishService       string
	PublishStatusAddress string

	UpdateStatus           bool
	ElectionID             string
	UpdateStatusOnShutdown bool

	EnableProfiling bool

	EnableSSLChainCompletion bool

	FakeCertificatePath string
	FakeCertificateSHA  string

	SyncRateLimit float32
}

// GetPublishService returns the configured service used to set ingress status
func (n NGINXController) GetPublishService() *apiv1.Service {
	s, err := n.store.GetService(n.cfg.PublishService)
	if err != nil {
		return nil
	}

	return s
}

// sync collects all the pieces required to assemble the configuration file and
// then sends the content to the backend (OnUpdate) receiving the populated
// template as response reloading the backend if is required.
func (n *NGINXController) syncIngress(interface{}) error {
	n.syncRateLimiter.Accept()

	if n.syncQueue.IsShuttingDown() {
		return nil
	}

	if !n.syncStatus.IsLeader() {
		glog.V(2).Infof("skipping synchronization of configuration because I am not the leader.")
		return nil
	}

	// Sort ingress rules using the ResourceVersion field
	ings := n.store.ListIngresses()
	sort.SliceStable(ings, func(i, j int) bool {
		ir := ings[i].ResourceVersion
		jr := ings[j].ResourceVersion
		return ir < jr
	})

	upstreams, servers := n.getBackendServers(ings)

	pcfg := ingress.Configuration{
		Backends: upstreams,
		Servers:  servers,
	}

	glog.Infof("syncing Ingress configuration...")
	err := n.OnUpdate(&pcfg)
	if err != nil {
		glog.Errorf("unexpected failure updating Kong configuration: \n%v", err)
		return err
	}

	n.runningConfig = &pcfg

	return nil
}

// getDefaultUpstream returns an upstream associated with the
// default backend service. In case of error retrieving information
// configure the upstream to return http code 503.
func (n *NGINXController) getDefaultUpstream() *ingress.Backend {
	upstream := &ingress.Backend{
		Name: defUpstreamName,
	}
	svcKey := n.cfg.DefaultService
	svc, err := n.store.GetService(svcKey)
	if err != nil {
		return upstream
	}

	endps := getEndpoints(svc, &svc.Spec.Ports[0], apiv1.ProtocolTCP, n.store.GetServiceEndpoints)
	if len(endps) == 0 {
		glog.Warningf("service %v does not have any active endpoints", svcKey)
	}

	upstream.Service = svc
	upstream.Endpoints = append(upstream.Endpoints, endps...)
	return upstream
}

// getBackendServers returns a list of Upstream and Server to be used by the backend
// An upstream can be used in multiple servers if the namespace, service name and port are the same
func (n *NGINXController) getBackendServers(ingresses []*extensions.Ingress) ([]*ingress.Backend, []*ingress.Server) {
	du := n.getDefaultUpstream()
	upstreams := n.createUpstreams(ingresses, du)
	servers := n.createServers(ingresses, upstreams, du)

	for _, ing := range ingresses {
		for _, rule := range ing.Spec.Rules {
			host := rule.Host
			if host == "" {
				host = defServerName
			}
			server := servers[host]
			if server == nil {
				server = servers[defServerName]
			}

			if rule.HTTP == nil &&
				host != defServerName {
				glog.V(3).Infof("ingress rule %v/%v does not contain HTTP rules, using default backend", ing.Namespace, ing.Name)
				continue
			}

			for _, path := range rule.HTTP.Paths {
				upsName := fmt.Sprintf("%v.%v.%v",
					ing.GetNamespace(),
					path.Backend.ServiceName,
					path.Backend.ServicePort.String())

				ups := upstreams[upsName]

				// if there's no path defined we assume /
				nginxPath := rootLocation
				if path.Path != "" {
					nginxPath = path.Path
				}

				addLoc := true
				for _, loc := range server.Locations {
					if loc.Path == nginxPath {
						addLoc = false

						if !loc.IsDefBackend {
							glog.V(3).Infof("avoiding replacement of ingress rule %v/%v location %v upstream %v (%v)", ing.Namespace, ing.Name, loc.Path, ups.Name, loc.Backend)
							break
						}

						glog.V(3).Infof("replacing ingress rule %v/%v location %v upstream %v (%v)", ing.Namespace, ing.Name, loc.Path, ups.Name, loc.Backend)
						loc.Backend = ups.Name
						loc.IsDefBackend = false
						loc.Backend = ups.Name
						loc.Port = ups.Port
						loc.Service = ups.Service
						loc.Ingress = ing

						break
					}
				}
				// is a new location
				if addLoc {
					glog.V(3).Infof("adding location %v in ingress rule %v/%v upstream %v", nginxPath, ing.Namespace, ing.Name, ups.Name)
					loc := &ingress.Location{
						Path:         nginxPath,
						Backend:      ups.Name,
						IsDefBackend: false,
						Service:      ups.Service,
						Port:         ups.Port,
						Ingress:      ing,
					}

					server.Locations = append(server.Locations, loc)
				}
			}
		}
	}

	aUpstreams := make([]*ingress.Backend, 0, len(upstreams))

	for _, upstream := range upstreams {
		//isHTTPSfrom := []*ingress.Server{}
		for _, server := range servers {
			for _, location := range server.Locations {
				if upstream.Name == location.Backend {
					if len(upstream.Endpoints) == 0 {
						glog.V(3).Infof("upstream %v does not have any active endpoints.", upstream.Name)
						location.Backend = ""

						// check if the location contains endpoints and a custom default backend
						if location.DefaultBackend != nil {
							sp := location.DefaultBackend.Spec.Ports[0]
							endps := getEndpoints(location.DefaultBackend, &sp, apiv1.ProtocolTCP, n.store.GetServiceEndpoints)
							if len(endps) > 0 {
								glog.V(3).Infof("using custom default backend in server %v location %v (service %v/%v)",
									server.Hostname, location.Path, location.DefaultBackend.Namespace, location.DefaultBackend.Name)
								nb := upstream.DeepCopy()
								name := fmt.Sprintf("custom-default-backend.%v", upstream.Name)
								nb.Name = name
								nb.Endpoints = endps
								aUpstreams = append(aUpstreams, nb)
								location.Backend = name
							}
						}
					}
				}
			}
		}
	}

	// create the list of upstreams and skip those without endpoints
	for _, upstream := range upstreams {
		if len(upstream.Endpoints) == 0 {
			continue
		}
		aUpstreams = append(aUpstreams, upstream)
	}

	aServers := make([]*ingress.Server, 0, len(servers))
	for _, value := range servers {
		sort.SliceStable(value.Locations, func(i, j int) bool {
			return value.Locations[i].Path > value.Locations[j].Path
		})
		aServers = append(aServers, value)
	}

	sort.SliceStable(aServers, func(i, j int) bool {
		return aServers[i].Hostname < aServers[j].Hostname
	})

	return aUpstreams, aServers
}

// createUpstreams creates the NGINX upstreams for each service referenced in
// Ingress rules. The servers inside the upstream are endpoints.
func (n *NGINXController) createUpstreams(data []*extensions.Ingress, du *ingress.Backend) map[string]*ingress.Backend {
	upstreams := make(map[string]*ingress.Backend)
	upstreams[defUpstreamName] = du

	for _, ing := range data {
		var defBackend string
		if ing.Spec.Backend != nil {
			defBackend = fmt.Sprintf("%v.%v.%v",
				ing.GetNamespace(),
				ing.Spec.Backend.ServiceName,
				ing.Spec.Backend.ServicePort.String())

			glog.V(3).Infof("creating upstream %v", defBackend)
			upstreams[defBackend] = newUpstream(defBackend)

			svcKey := fmt.Sprintf("%v/%v", ing.GetNamespace(), ing.Spec.Backend.ServiceName)

			if len(upstreams[defBackend].Endpoints) == 0 {
				endps, err := n.serviceEndpoints(svcKey, ing.Spec.Backend.ServicePort.String())
				upstreams[defBackend].Endpoints = append(upstreams[defBackend].Endpoints, endps...)
				if err != nil {
					glog.Warningf("error creating upstream %v: %v", defBackend, err)
				}
			}

		}

		for _, rule := range ing.Spec.Rules {
			if rule.HTTP == nil {
				continue
			}

			for _, path := range rule.HTTP.Paths {
				name := fmt.Sprintf("%v.%v.%v",
					ing.GetNamespace(),
					path.Backend.ServiceName,
					path.Backend.ServicePort.String())

				if _, ok := upstreams[name]; ok {
					continue
				}

				glog.V(3).Infof("creating upstream %v", name)
				upstreams[name] = newUpstream(name)
				upstreams[name].Port = path.Backend.ServicePort

				svcKey := fmt.Sprintf("%v/%v", ing.GetNamespace(), path.Backend.ServiceName)

				if len(upstreams[name].Endpoints) == 0 {
					endp, err := n.serviceEndpoints(svcKey, path.Backend.ServicePort.String())
					if err != nil {
						glog.Warningf("error obtaining service endpoints: %v", err)
						continue
					}
					upstreams[name].Endpoints = endp
				}

				s, err := n.store.GetService(svcKey)
				if err != nil {
					glog.Warningf("error obtaining service: %v", err)
					continue
				}

				upstreams[name].Service = s
			}
		}
	}

	return upstreams
}

// serviceEndpoints returns the upstream servers (endpoints) associated
// to a service.
func (n *NGINXController) serviceEndpoints(svcKey, backendPort string) ([]ingress.Endpoint, error) {
	svc, err := n.store.GetService(svcKey)

	var upstreams []ingress.Endpoint
	if err != nil {
		return upstreams, fmt.Errorf("error getting service %v from the cache: %v", svcKey, err)
	}

	glog.V(3).Infof("obtaining port information for service %v", svcKey)
	for _, servicePort := range svc.Spec.Ports {
		// targetPort could be a string, use the name or the port (int)
		if strconv.Itoa(int(servicePort.Port)) == backendPort ||
			servicePort.TargetPort.String() == backendPort ||
			servicePort.Name == backendPort {

			endps := getEndpoints(svc, &servicePort, apiv1.ProtocolTCP, n.store.GetServiceEndpoints)
			if len(endps) == 0 {
				glog.Warningf("service %v does not have any active endpoints", svcKey)
			}

			upstreams = append(upstreams, endps...)
			break
		}
	}

	// Ingress with an ExternalName service and no port defined in the service.
	if len(svc.Spec.Ports) == 0 && svc.Spec.Type == apiv1.ServiceTypeExternalName {
		externalPort, err := strconv.Atoi(backendPort)
		if err != nil {
			glog.Warningf("only numeric ports are allowed in ExternalName services: %v is not valid as a TCP/UDP port", backendPort)
			return upstreams, nil
		}

		servicePort := apiv1.ServicePort{
			Protocol:   "TCP",
			Port:       int32(externalPort),
			TargetPort: intstr.FromString(backendPort),
		}
		endps := getEndpoints(svc, &servicePort, apiv1.ProtocolTCP, n.store.GetServiceEndpoints)
		if len(endps) == 0 {
			glog.Warningf("service %v does not have any active endpoints", svcKey)
			return upstreams, nil
		}

		upstreams = append(upstreams, endps...)
		return upstreams, nil
	}

	return upstreams, nil
}

// createServers initializes a map that contains information about the list of
// FDQN referenced by ingress rules and the common name field in the referenced
// SSL certificates. Each server is configured with location / using a default
// backend specified by the user or the one inside the ingress spec.
func (n *NGINXController) createServers(data []*extensions.Ingress,
	upstreams map[string]*ingress.Backend,
	du *ingress.Backend) map[string]*ingress.Server {

	servers := make(map[string]*ingress.Server, len(data))
	// If a server has a hostname equivalent to a pre-existing alias, then we
	// remove the alias to avoid conflicts.
	aliases := make(map[string]string, len(data))

	// initialize the default server
	servers[defServerName] = &ingress.Server{
		Hostname: defServerName,
		Locations: []*ingress.Location{
			{
				Path:         rootLocation,
				IsDefBackend: true,
				Backend:      du.Name,
				Service:      du.Service,
			},
		}}

	// initialize all the servers
	for _, ing := range data {
		// default upstream server
		un := du.Name

		if ing.Spec.Backend != nil {
			// replace default backend
			defUpstream := fmt.Sprintf("%v.%v.%v", ing.GetNamespace(), ing.Spec.Backend.ServiceName, ing.Spec.Backend.ServicePort.String())
			if backendUpstream, ok := upstreams[defUpstream]; ok {
				un = backendUpstream.Name

				// Special case:
				// ingress only with a backend and no rules
				// this case defines a "catch all" server
				defLoc := servers[defServerName].Locations[0]
				if defLoc.IsDefBackend && len(ing.Spec.Rules) == 0 {
					defLoc.IsDefBackend = false
					defLoc.Backend = backendUpstream.Name
					defLoc.Service = backendUpstream.Service
					defLoc.Ingress = ing
				}
			}
		}

		for _, rule := range ing.Spec.Rules {
			host := rule.Host
			if host == "" {
				host = defServerName
			}
			if _, ok := servers[host]; ok {
				// server already configured
				continue
			}

			servers[host] = &ingress.Server{
				Hostname: host,
				Locations: []*ingress.Location{
					{
						Path:         rootLocation,
						IsDefBackend: true,
						Backend:      un,
						Service:      &apiv1.Service{},
					},
				},
			}
		}
	}

	// configure default location, alias, and SSL
	for _, ing := range data {
		for _, rule := range ing.Spec.Rules {
			host := rule.Host
			if host == "" {
				host = defServerName
			}

			// only add a certificate if the server does not have one previously configured
			if servers[host].SSLCert != nil {
				continue
			}

			if len(ing.Spec.TLS) == 0 {
				glog.V(3).Infof("ingress %v/%v for host %v does not contains a TLS section", ing.Namespace, ing.Name, host)
				continue
			}

			tlsSecretName := ""
			found := false
			for _, tls := range ing.Spec.TLS {
				if sets.NewString(tls.Hosts...).Has(host) {
					tlsSecretName = tls.SecretName
					found = true
					break
				}
			}

			if !found {
				// does contain a TLS section but none of the
				// host match or there is no hosts in the TLS section
				// as last resort we validate the host against the certificate
				// if is valid we use it.
				for _, tls := range ing.Spec.TLS {
					key := fmt.Sprintf("%v/%v", ing.Namespace, tls.SecretName)
					cert, err := n.store.GetCertFromSecret(key)
					if err != nil {
						glog.Warningf("ssl certificate \"%v\" does not exist in k8s secret store", key)
						continue
					}

					if sets.NewString(cert.CN...).Has(host) {
						tlsSecretName = tls.SecretName
						found = true
						break
					}
				}

				if !found {
					continue
				}
			}

			key := fmt.Sprintf("%v/%v", ing.Namespace, tlsSecretName)
			cert, err := n.store.GetCertFromSecret(key)
			if err != nil {
				glog.Warningf("ssl certificate \"%v\" does not exist in store", key)
				continue
			}

			servers[host].SSLCert = cert

			if cert.ExpireTime.Before(time.Now().Add(240 * time.Hour)) {
				glog.Warningf("ssl certificate for host %v is about to expire in 10 days", host)
			}
		}
	}

	for alias, host := range aliases {
		if _, ok := servers[alias]; ok {
			glog.Warningf("There is a conflict with server hostname '%v' and alias '%v'. Removing alias to avoid conflicts.", host, alias)
			servers[host].Alias = ""
		}
	}

	return servers
}

// newUpstream creates an upstream without servers.
func newUpstream(name string) *ingress.Backend {
	return &ingress.Backend{
		Name:      name,
		Endpoints: []ingress.Endpoint{},
		Service:   &apiv1.Service{},
	}
}
