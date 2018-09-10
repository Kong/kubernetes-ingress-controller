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
	"bytes"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/imdario/mergo"

	"github.com/fatih/structs"
	"github.com/golang/glog"
	kong "github.com/kong/kubernetes-ingress-controller/internal/apis/admin"
	kongadminv1 "github.com/kong/kubernetes-ingress-controller/internal/apis/admin/v1"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/internal/apis/configuration/v1"
	pluginv1 "github.com/kong/kubernetes-ingress-controller/internal/apis/plugin/v1"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/annotations"
	"github.com/kong/kubernetes-ingress-controller/internal/net/ssl"
	"github.com/pkg/errors"
	extensions "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/util/sets"
)

// OnUpdate is called periodically by syncQueue to keep the configuration in sync.
// returning nil implies the synchronization finished correctly.
// Returning an error means requeue the update.
func (n *NGINXController) OnUpdate(ingressCfg *ingress.Configuration) error {
	// Synchronizde the state between Kubernetes and Kong with this order:
	//	- SSL Certificates
	//	- SNIs
	// 	- Upstreams
	//	- Upstream targets
	//  - Consumers
	//  - Credentials
	// 	- Services (and Plugins)
	//  - Routes (and Plugins)

	// TODO: All the resources created by the ingress controller add the annotations
	// kong-ingress-controller and kubernetes
	// This allows the identification of reources that can be removed if they
	// are not present in Kubernetes when the sync process is executed.
	// For instance an Ingress, Service or Secret is removed.

	for _, server := range ingressCfg.Servers {
		if server.Hostname == "_" {
			// there is no catch all server in kong
			continue
		}

		// check the certificate is present in kong
		if server.SSLCert != nil {
			err := n.syncCertificate(server)
			if err != nil {
				return err
			}
		}

		err := n.syncUpstreams(server.Locations, ingressCfg.Backends)
		if err != nil {
			return err
		}
	}

	err := n.syncConsumers()
	if err != nil {
		return err
	}

	err = n.syncCredentials()
	if err != nil {
		return err
	}

	err = n.syncGlobalPlugins()
	if err != nil {
		return err
	}

	checkServices, err := n.syncServices(ingressCfg)
	if err != nil {
		return err
	}

	checkRoutes, err := n.syncRoutes(ingressCfg)
	if err != nil {
		return err
	}

	// trigger a new sync event to ensure routes and services are up to date
	// this is required because the plugins configuration could be incorrect
	// if some delete occurred.
	if checkServices || checkRoutes {
		defer func() {
			n.syncQueue.Enqueue(&extensions.Ingress{})
		}()
	}

	return nil
}

func (n *NGINXController) syncGlobalPlugins() error {
	glog.Infof("syncing global plugins")

	targetGlobalPlugins, err := n.store.ListGlobalKongPlugins()
	if err != nil {
		return err
	}
	targetPluginMap := make(map[string]*pluginv1.KongPlugin)
	var duplicates []string // keep track of duplicate

	for i := 0; i < len(targetGlobalPlugins); i++ {
		name := targetGlobalPlugins[i].Name
		// empty name skip it
		if name == "" {
			continue
		}
		if _, ok := targetPluginMap[name]; ok {
			glog.Error("Multiple KongPlugin definitions found with 'global' annotation for :", name,
				", the plugin will not be applied")
			duplicates = append(duplicates, name)
			continue
		}
		targetPluginMap[name] = targetGlobalPlugins[i]
	}

	// remove duplicates
	for _, plugin := range duplicates {
		delete(targetPluginMap, plugin)
	}

	client := n.cfg.Kong.Client
	plugins, err := client.Plugins().List(nil)
	if err != nil {
		return err
	}

	// plugins in Kong
	currentGlobalPlugins := make(map[string]kongadminv1.Plugin)
	for _, plugin := range plugins.Items {
		if plugin.Route == "" && plugin.Service == "" && plugin.Consumer == "" {
			currentGlobalPlugins[plugin.Name] = plugin
		}
	}

	// sync plugins to Kong
	for pluginName, kongPlugin := range targetPluginMap {
		// plugin exists?
		if pluginInKong, ok := currentGlobalPlugins[pluginName]; !ok {
			// no, create it
			p := &kongadminv1.Plugin{
				Name:   pluginName,
				Config: kongadminv1.Configuration(kongPlugin.Config),
			}
			_, res := client.Plugins().CreateGlobal(p)
			if res.StatusCode != http.StatusCreated {
				return errors.Wrap(res.Error(), fmt.Sprintf("creating a global Kong plugin %v", p))
			}
		} else {
			// plugin exists, is the configuration up to date
			if !pluginDeepEqual(kongPlugin.Config, &pluginInKong) {
				// no, update it
				p := &kongadminv1.Plugin{
					Name:   pluginName,
					Config: kongadminv1.Configuration(kongPlugin.Config),
				}
				_, res := client.Plugins().Patch(pluginInKong.ID, p)
				if res.StatusCode != http.StatusOK {
					return errors.Wrap(res.Error(), fmt.Sprintf("updating a global Kong plugin %v", p))
				}
			}
		}
		// remove from the current list, all that remain in the current list will be deleted
		delete(currentGlobalPlugins, pluginName)
	}

	// delete the ones not configured in k8s
	for _, plugin := range currentGlobalPlugins {
		err := client.Plugins().Delete(plugin.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// syncTargets reconciles the state between the ingress controller and
// kong comparing the endpoints in Kubernetes and the targets in a
// particular kong upstream. To avoid downtimes we create the new targets
// first and then remove the killed ones.
func syncTargets(upstream string, ingressEndpopint *ingress.Backend, client *kong.RestClient) error {
	glog.V(3).Infof("syncing Kong targets")
	b, res := client.Upstreams().Get(upstream)
	if res.StatusCode == http.StatusNotFound {
		glog.Errorf("there is no upstream with name %v in Kong", upstream)
		return nil
	}

	kongTargets, err := client.Targets().List(nil, upstream)
	if err != nil {
		return err
	}

	oldTargets := sets.NewString()
	for _, kongTarget := range kongTargets.Items {
		if !oldTargets.Has(kongTarget.Target) {
			oldTargets.Insert(kongTarget.Target)
		}
	}

	newTargets := sets.NewString()
	for _, endpoint := range ingressEndpopint.Endpoints {
		nt := fmt.Sprintf("%v:%v", endpoint.Address, endpoint.Port)
		if !newTargets.Has(nt) {
			newTargets.Insert(nt)
		}
	}

	add := newTargets.Difference(oldTargets).List()
	remove := oldTargets.Difference(newTargets).List()

	for _, endpoint := range add {
		target := &kongadminv1.Target{
			Target:   endpoint,
			Upstream: b.ID,
		}
		glog.Infof("creating Kong Target %v for upstream %v", endpoint, b.ID)
		_, res := client.Targets().Create(target, upstream)
		if res.StatusCode != http.StatusCreated {
			glog.Errorf("Unexpected error creating Kong Upstream: %v", res)
			return res.Error()
		}
	}

	// wait to avoid hitting the kong API server too fast
	time.Sleep(100 * time.Millisecond)

	for _, endpoint := range remove {
		for _, kongTarget := range kongTargets.Items {
			if kongTarget.Target != endpoint {
				continue
			}
			glog.Infof("deleting Kong Target %v from upstream %v", kongTarget.ID, kongTarget.Upstream)
			err := client.Targets().Delete(kongTarget.ID, b.ID)
			if err != nil {
				glog.Errorf("Unexpected error deleting Kong Upstream: %v", err)
				return err
			}
		}
	}

	return nil
}

// syncServices reconciles the state between the ingress controller and
// kong services. After the synchronization of services we also check
// if there was any changes to the kong plugins applied to the service
func (n *NGINXController) syncServices(ingressCfg *ingress.Configuration) (bool, error) {
	client := n.cfg.Kong.Client

	servicesToKeep := sets.NewString()

	// triggerReload indicates if the sync process altered configuration with services
	// and require an additional run
	var triggerReload bool

	// Check if the endpoints exists as a service in kong
	for _, server := range ingressCfg.Servers {
		if server.Hostname == "_" {
			// there is no catch all server in kong
			continue
		}

		for _, location := range server.Locations {
			backend := location.Backend
			if backend == "default-backend" {
				// there is no default backend in Kong
				continue
			}

			ingress := location.Ingress
			if ingress == nil {
				// location is the default backend (not mapped against Kong)
				continue
			}

			if backend == "" {
				glog.Warningf("the service defined in the ingress %v/%v does not exists", ingress.Namespace, ingress.Name)
				continue
			}

			kongIngress, err := n.getKongIngress(ingress)
			if err != nil {
				glog.Warningf("there is no custom Ingress configuration for rule %v/%v", ingress.Namespace, ingress.Name)
			}

			name := buildName(backend, location)
			for _, upstream := range ingressCfg.Backends {
				if upstream.Name != backend {
					continue
				}

				// defaults
				proto := "http"
				port := 80

				s, res := client.Services().Get(name)
				if res.StatusCode == http.StatusNotFound {
					s = &kongadminv1.Service{
						Name:     name,
						Path:     "/",
						Protocol: proto,
						Host:     name,
						Port:     port,
						Retries:  5,
					}
				}

				if s != nil {
					if kongIngress != nil && kongIngress.Proxy != nil {
						if kongIngress.Proxy.Path != "" {
							s.Path = kongIngress.Proxy.Path
						}

						if kongIngress.Proxy.Protocol != "" &&
							(kongIngress.Proxy.Protocol == "http" || kongIngress.Proxy.Protocol == "https") {
							s.Protocol = kongIngress.Proxy.Protocol
						}

						if kongIngress.Proxy.ConnectTimeout > 0 {
							s.ConnectTimeout = kongIngress.Proxy.ConnectTimeout
						}

						if kongIngress.Proxy.ReadTimeout > 0 {
							s.ReadTimeout = kongIngress.Proxy.ReadTimeout
						}

						if kongIngress.Proxy.WriteTimeout > 0 {
							s.WriteTimeout = kongIngress.Proxy.WriteTimeout
						}

						if kongIngress.Proxy.Retries > 0 {
							s.Retries = kongIngress.Proxy.Retries
						}
					}

					if res.StatusCode == http.StatusNotFound {
						glog.Infof("Creating Kong Service name %v", name)
						_, res := client.Services().Create(s)
						if res.StatusCode != http.StatusCreated {
							glog.Errorf("Unexpected error creating Kong Service: %v", res)
							return false, res.Error()
						}
					} else {
						glog.Infof("Patching Kong Service name %v", name)
						_, res := client.Services().Patch(s.ID, s)
						if res.StatusCode != http.StatusOK {
							glog.Errorf("Unexpected error patching Kong Service: %v", res)
							return false, res.Error()
						}
					}
				}

				// do not remove the service
				if !servicesToKeep.Has(name) {
					servicesToKeep.Insert(name)
				}

				break
			}

			svc, res := client.Services().Get(name)
			if res.StatusCode == http.StatusNotFound || svc.ID == "" {
				glog.Warningf("service %v does not exists in kong", name)
				continue
			}

			plugins := annotations.ExtractKongPluginAnnotations(location.Service.GetAnnotations())

			if len(plugins) == 0 {
				glog.Infof("service %v/%v does not contain any plugins. Checking if it is required to remove plugins...",
					location.Service.Namespace, location.Service.Name)
				// remove all the plugins from the service.
				plugins, err := client.Plugins().GetAllByService(svc.ID)
				if err != nil {
					glog.Errorf("Unexpected error obtaining Kong plugins for service %v: %v", svc.ID, err)
					continue
				}

				for _, plugin := range plugins {
					glog.Infof("Removing plugin %v from service %v", plugin.ID, svc.ID)
					err := client.Plugins().Delete(plugin.ID)
					if err != nil {
						return false, errors.Wrap(err, "deleting Kong plugin")
					}

					glog.Infof("Plugin %v successfully removed from service %v", plugin.ID, svc.ID)
				}
			} else {
				glog.Infof("configuring plugins '%v' for service %v...", plugins, svc.ID)
			}

			// configure plugins poresent in the service
			for plugin, crdNames := range plugins {
				for _, crdName := range crdNames {
					// search configured plugin CRD in k8s
					k8sPlugin, err := n.store.GetKongPlugin(location.Ingress.Namespace, crdName)
					if err != nil {
						return false, errors.Wrap(err, fmt.Sprintf("searching plugin KongPlugin %v", crdName))
					}

					pluginID := fmt.Sprintf("%v", k8sPlugin.GetUID())
					// The plugin is not defined in the service.
					// check if the route has the plugin or is required to
					// create a new one
					configuredPlugin, err := client.Plugins().GetByID(pluginID)
					if err != nil {
						if !kong.IsPluginNotConfiguredError(err) {
							return false, errors.Wrap(err, fmt.Sprintf("getting Kong plugin %v", pluginID))
						}

						// there is no plugin, create a new one
						p := &kongadminv1.Plugin{
							Name:    plugin,
							Service: svc.ID,
							Config:  kongadminv1.Configuration(k8sPlugin.Config),
							Required: kongadminv1.Required{
								ID: pluginID,
							},
						}

						if k8sPlugin.ConsumerRef != "" {
							consumer, err := n.store.GetKongConsumer(k8sPlugin.Namespace, k8sPlugin.ConsumerRef)
							if err != nil {
								glog.Errorf("Unexpected error searching plugin configuration %v for service %v: %v", plugin, svc.ID, err)
							} else {
								p.Consumer = fmt.Sprintf("%v", consumer.GetUID())
							}
						}

						_, res := client.Plugins().CreateInService(svc.ID, p)
						if res.StatusCode != http.StatusCreated {
							return false, errors.Wrap(res.Error(), fmt.Sprintf("creating a Kong plugin %v in service %v", plugin, svc.ID))
						}

						continue
					}

					// check the kong plugin configuration is up to date
					if !pluginDeepEqual(k8sPlugin.Config, configuredPlugin) ||
						(!k8sPlugin.Disabled != configuredPlugin.Enabled) {
						glog.Infof("plugin %v configuration in kong is outdated. Updating...", k8sPlugin.Name)
						p := &kongadminv1.Plugin{
							Name:    configuredPlugin.Name,
							Config:  kongadminv1.Configuration(k8sPlugin.Config),
							Enabled: k8sPlugin.Disabled,
							Service: svc.ID,
							Route:   "",
						}

						if k8sPlugin.ConsumerRef != "" && configuredPlugin.Consumer == "" {
							consumer, err := n.store.GetKongConsumer(k8sPlugin.Namespace, k8sPlugin.ConsumerRef)
							if err != nil {
								glog.Errorf("Unexpected error searching plugin configuration %v for service %v: %v",
									plugin, svc.ID, err)
								continue
							}
							p.Consumer = fmt.Sprintf("%v", consumer.GetUID())
						}
						_, res := client.Plugins().Patch(configuredPlugin.ID, p)
						if res.StatusCode != http.StatusOK {
							glog.Errorf("Unexpected error updating plugin configuration %v for service %v: %v",
								plugin, svc.ID, res)
						}

						continue
					}

					glog.Infof("plugin %v configuration in kong is up to date.", k8sPlugin.Name)
				}
			}
		}
	}

	kongServices, err := client.Services().List(nil)
	if err != nil {
		return false, err
	}

	serviceNames := sets.NewString()

	for _, svc := range kongServices.Items {
		if !serviceNames.Has(svc.Name) {
			serviceNames.Insert(svc.Name)
		}
	}

	servicesToRemove := serviceNames.Difference(servicesToKeep)

	// remove all those services that are present in Kong but not in the current Kubernetes state
	for _, svcName := range servicesToRemove.List() {
		svc, res := client.Services().Get(svcName)
		if res.StatusCode == http.StatusNotFound || svc == nil {
			glog.Warningf("service %v does not exists in kong", svcName)
			continue
		}

		glog.Infof("deleting Kong Service %v", svcName)
		// before deleting the service we need to remove the upstream and th targets that reference the service
		err := deleteServiceUpstream(svc.Name, client)
		if err != nil {
			glog.Errorf("Unexpected error deleting Kong upstreams and targets that depend on service %v: %v", svc.Name, err)
			continue
		}
		err = client.Services().Delete(svc.ID)
		if err != nil {
			// this means the service is being referenced by a route
			// during the next sync it will be removed
			glog.V(3).Infof("Unexpected error deleting Kong Service: %v", err)
		}
	}

	if len(servicesToRemove) > 0 {
		triggerReload = true
	}

	return triggerReload, nil
}

// syncConsumers synchronizes the state between KongConsumer (Kubernetes CRD) type and Kong consumers.
// This loop only creates new consumers in Kong.
func (n *NGINXController) syncConsumers() error {

	consumersInKong := make(map[string]*kongadminv1.Consumer)
	client := n.cfg.Kong.Client

	// List all consumers in Kong
	kongConsumers, err := client.Consumers().List(nil)
	if err != nil {
		return err
	}

	for i := range kongConsumers.Items {
		consumersInKong[kongConsumers.Items[i].ID] = &kongConsumers.Items[i]
	}

	// List existing Consumers in Kubernetes
	for _, consumer := range n.store.ListKongConsumers() {
		glog.Infof("checking if Kong consumer %v exists", consumer.Name)
		consumerID := fmt.Sprintf("%v", consumer.GetUID())

		kc, ok := consumersInKong[consumerID]

		if !ok {
			glog.Infof("Creating Kong consumer %v", consumerID)
			c := &kongadminv1.Consumer{
				Username: consumer.Username,
				CustomID: consumer.CustomID,
				Required: kongadminv1.Required{
					ID: consumerID,
				},
			}

			c, res := n.cfg.Kong.Client.Consumers().Create(c)
			if res.StatusCode != http.StatusCreated {
				return errors.Wrap(res.Error(), "creating a Kong consumer")
			}
		} else {
			// check the consumers are equals
			if consumer.Username != kc.Username || consumer.CustomID != kc.CustomID {
				kc.Username = consumer.Username
				kc.CustomID = consumer.CustomID
				_, res := n.cfg.Kong.Client.Consumers().Patch(consumerID, kc)
				if res.StatusCode != http.StatusOK {
					return errors.Wrap(res.Error(), "patching a Kong consumer")
				}
			}
		}
		delete(consumersInKong, consumerID)
	}
	// remaining entries in the map should be deleted

	for _, consumer := range consumersInKong {
		err := client.Consumers().Delete(consumer.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// set of valid authentication plugins for consumers
var validCredentialTypes = sets.NewString(
	// Kong CE and EE
	"oauth2",
	"jwt",
	"basic-auth",
	"acl",
	"key-auth",
	"hmac-auth",
	"ldap-authentication",
	// Kong EE only
	"openid-connect",
	"oauth2-introspection",
)

// syncCredentials synchronizes the state between KongCredential (Kubernetes CRD) and Kong credentials.
func (n *NGINXController) syncCredentials() error {
	// List existing credentials in Kubernetes
	for _, credential := range n.store.ListKongCredentials() {
		if !validCredentialTypes.Has(credential.Type) {
			glog.Errorf("the credential does contains a valid type: %v", credential.Type)
			continue
		}
		credentialID := fmt.Sprintf("%v", credential.GetUID())

		consumer, err := n.store.GetKongConsumer(credential.Namespace, credential.ConsumerRef)
		if err != nil {
			glog.Errorf("Unexpected error searching KongConsumer: %v", err)
			continue
		}
		consumerID := fmt.Sprintf("%v", consumer.GetUID())

		// TODO: allow changes in credentials?
		_, res := n.cfg.Kong.Client.Credentials().GetByType(consumerID, credentialID, credential.Type)
		if res.StatusCode == http.StatusNotFound {
			// use the configuration
			data := credential.Config
			// create a credential with the same id of k8s
			data["id"] = credentialID
			data["consumer_id"] = consumerID
			res := n.cfg.Kong.Client.Credentials().CreateByType(data, consumerID, credential.Type)
			if res.StatusCode != http.StatusCreated {
				glog.Errorf("Unexpected error updating Kong Route: %v", res)
				return res.Error()
			}
		}
	}

	return nil
}

// syncRoutes synchronizes the state between the Ingress configuration model and Kong Routes.
func (n *NGINXController) syncRoutes(ingressCfg *ingress.Configuration) (bool, error) {
	client := n.cfg.Kong.Client

	kongRoutes, err := client.Routes().List(nil)
	if err != nil {
		return false, err
	}

	// triggerReload indicates if the sync process altered
	// configuration with services and require an additional run
	var triggerReload bool

	// create a copy of the existing routes to be able to run a comparison
	routesToRemove := sets.NewString()
	for _, old := range kongRoutes.Items {
		if !routesToRemove.Has(old.ID) {
			routesToRemove.Insert(old.ID)
		}
	}

	// Routes
	for _, server := range ingressCfg.Servers {
		if server.Hostname == "_" {
			// there is no catch all server in kong
			continue
		}

		protos := []string{"http"}
		if server.SSLCert != nil {
			protos = append(protos, "https")
		}

		for _, location := range server.Locations {
			backend := location.Backend
			if backend == "default-backend" {
				// there is no default backend in Kong
				continue
			}

			ingress := location.Ingress
			if ingress == nil {
				// location is the default backend (not mapped against Kong)
				continue
			}

			if backend == "" {
				glog.Warningf("the service defined in the ingress %v/%v does not exists", ingress.Namespace, ingress.Name)
				continue
			}

			kongIngress, err := n.getKongIngress(ingress)
			if err != nil {
				glog.Warningf("there is no custom Ingress configuration for rule %v/%v", ingress.Namespace, ingress.Name)
			}

			name := buildName(backend, location)
			svc, res := client.Services().Get(name)
			if res.StatusCode == http.StatusNotFound || svc.ID == "" {
				glog.Warningf("service %v does not exists in kong", name)
				continue
			}

			r := &kongadminv1.Route{
				Paths:     []string{location.Path},
				Protocols: []string{"http", "https"}, // default
				Hosts:     []string{server.Hostname},
				Service:   kongadminv1.InlineService{ID: svc.ID},
			}

			if kongIngress != nil && kongIngress.Route != nil {
				if len(kongIngress.Route.Methods) > 0 {
					r.Methods = kongIngress.Route.Methods
				}

				if kongIngress.Route.PreserveHost != r.PreserveHost {
					r.PreserveHost = kongIngress.Route.PreserveHost
				}

				if kongIngress.Route.RegexPriority != r.RegexPriority {
					r.RegexPriority = kongIngress.Route.RegexPriority
				}

				if kongIngress.Route.StripPath != r.StripPath {
					r.StripPath = kongIngress.Route.StripPath
				}
				if len(kongIngress.Route.Protocols) != 0 {
					r.Protocols = kongIngress.Route.Protocols
				}
			}

			if !isRouteInKong(r, kongRoutes.Items) {
				glog.Infof("creating Kong Route for host %v, path %v and service %v", server.Hostname, location.Path, svc.ID)
				_, res := client.Routes().Create(r)
				if res.StatusCode != http.StatusCreated {
					glog.Errorf("Unexpected error creating Kong Route: %v", res)
					return false, res.Error()
				}
			} else {
				// the route exists but the properties could have changed
				route := getKongRoute(server.Hostname, location.Path, kongRoutes.Items)

				if routesToRemove.Has(route.ID) {
					routesToRemove.Delete(route.ID)
				}
				outOfSync := false
				sort.Strings(protos)
				sort.Strings(route.Protocols)

				if !reflect.DeepEqual(protos, route.Protocols) {
					outOfSync = true
					route.Protocols = protos
				}

				if kongIngress != nil && kongIngress.Route != nil {
					if kongIngress.Route.Methods != nil {
						sort.Strings(kongIngress.Route.Methods)
						sort.Strings(route.Methods)
						if !reflect.DeepEqual(route.Methods, kongIngress.Route.Methods) {
							outOfSync = true
							route.Methods = kongIngress.Route.Methods
						}
					}

					if kongIngress.Route.PreserveHost != route.PreserveHost {
						outOfSync = true
						route.PreserveHost = kongIngress.Route.PreserveHost
					}

					if kongIngress.Route.RegexPriority != route.RegexPriority {
						outOfSync = true
						route.RegexPriority = kongIngress.Route.RegexPriority
					}

					if kongIngress.Route.StripPath != route.StripPath {
						outOfSync = true
						route.StripPath = kongIngress.Route.StripPath
					}
					if len(kongIngress.Route.Protocols) > 0 {
						sort.Strings(kongIngress.Route.Protocols)
						if !reflect.DeepEqual(route.Protocols, kongIngress.Route.Protocols) {
							outOfSync = true
							glog.Infof("protocols changed form ", route.Protocols, " to ", kongIngress.Route.Protocols)
							route.Protocols = kongIngress.Route.Protocols
						}
					}
				}

				if outOfSync {
					glog.Infof("updating Kong Route for host %v, path %v and service %v", server.Hostname, location.Path, svc.ID)
					_, res := client.Routes().Patch(route.ID, route)
					if res.StatusCode != http.StatusOK {
						glog.Errorf("Unexpected error updating Kong Route: %v", res)
						return false, res.Error()
					}
				}
			}

			kongRoutes, err = client.Routes().List(nil)
			if err != nil {
				return false, err
			}

			route := getKongRoute(server.Hostname, location.Path, kongRoutes.Items)

			plugins := annotations.ExtractKongPluginAnnotations(location.Ingress.GetAnnotations())

			if len(plugins) == 0 {
				glog.Errorf("Route %v does not contain any plugins. Checking if it is required to remove plugins...", route.ID)
				plugins, err := client.Plugins().GetAllByRoute(route.ID)
				if err != nil {
					glog.Errorf("Unexpected error obtaining Kong plugins for route %v: %v", route.ID, err)
					continue
				}

				for _, plugin := range plugins {
					glog.Errorf("Removing plugin %v from route %v", plugin.ID, route.ID)
					err := client.Plugins().Delete(plugin.ID)
					if err != nil {
						return false, errors.Wrap(err, "deleting Kong plugin")
					}

					glog.Infof("Plugin %v successfully removed from route %v", plugin.ID, svc.ID)
				}
			} else {
				glog.Infof("configuring plugins '%v' for route %v...", plugins, route.ID)
			}

			// pluginsInService contains the list of plugins configured in the
			// service and the source of truth to remove plugins from routes
			var pluginsInService []string

			// The Ingress contains at least one plugin.
			// Before starting to process the configuration we need to check
			// if the service does not have the same plugin configured.
			// In that case we need to skip the configuration in the route
			if len(location.Service.GetAnnotations()) > 0 {
				sa := annotations.ExtractKongPluginAnnotations(location.Service.GetAnnotations())
				for p := range sa {
					_, ok := plugins[p]
					if ok {
						glog.Warningf("Plugin %v is already configured in service %v/%v. Omitting plugin creation in Kong Route",
							p, location.Service.Namespace, location.Service.Name)
						delete(plugins, p)
						pluginsInService = append(pluginsInService, p)
					}
				}
			}

			for plugin, crdNames := range plugins {
				for _, crdName := range crdNames {
					glog.Infof("configuring plugin stored in KongPlugin CRD %v", crdName)
					// search configured plugin CRD in k8s
					k8sPlugin, err := n.store.GetKongPlugin(location.Ingress.Namespace, crdName)
					if err != nil {
						glog.Errorf("Unexpected error searching plugin %v for route %v: %v", plugin, route.ID, err)
						continue
					}

					pluginID := fmt.Sprintf("%v", k8sPlugin.GetUID())
					// The plugin is not defined in the service.
					// check if the route has the plugin or is required to
					// create a new one
					configuredPlugin, err := client.Plugins().GetByID(pluginID)
					if err != nil {
						if !kong.IsPluginNotConfiguredError(err) {
							glog.Errorf("%v", err)
							continue
						}

						// there is no plugin, create a new one
						p := &kongadminv1.Plugin{
							Name:   plugin,
							Route:  route.ID,
							Config: kongadminv1.Configuration(k8sPlugin.Config),
							Required: kongadminv1.Required{
								ID: pluginID,
							},
						}

						_, res := client.Plugins().CreateInRoute(route.ID, p)
						if res.StatusCode != http.StatusCreated {
							glog.Errorf("Unexpected error creating plugin %v for route %v: %v", plugin, route.ID, res)
						}

						continue
					}

					// check the kong plugin configuration is up to date
					if !pluginDeepEqual(k8sPlugin.Config, configuredPlugin) ||
						(!k8sPlugin.Disabled != configuredPlugin.Enabled) {
						glog.Infof("plugin %v configuration in kong is outdated. Updating...", k8sPlugin.Name)
						p := &kongadminv1.Plugin{
							Name:    configuredPlugin.Name,
							Config:  kongadminv1.Configuration(k8sPlugin.Config),
							Enabled: !k8sPlugin.Disabled,
							Service: "",
							Route:   route.ID,
						}

						_, res := client.Plugins().Patch(configuredPlugin.ID, p)
						if res.StatusCode != http.StatusOK {
							glog.Errorf("Unexpected error updating plugin %v for route %v: %v", plugin, route.ID, res)
						}

						continue
					}

					glog.Infof("plugin %v configuration in kong is up to date.", k8sPlugin.Name)
				}
			}

			if len(pluginsInService) == 0 {
				continue
			}

			// delete plugins from route configured in a service
			glog.Infof("deleting Kong Route plugins in route %v already defined in a service: '%v'", route.ID, pluginsInService)
			for _, plugin := range pluginsInService {
				configuredPlugin, err := client.Plugins().GetByRoute(plugin, route.ID)
				if err != nil {
					if !kong.IsPluginNotConfiguredError(err) {
						glog.Errorf("%v", err)
					}
					continue
				}

				err = client.Plugins().Delete(configuredPlugin.ID)
				if err != nil {
					glog.Errorf("Unexpected error deleting Kong plugin %v: %v", configuredPlugin.ID, err)
					continue
				}

				triggerReload = true
			}
		}
	}

	// remove all those routes that are present in Kong but not in the current Kubernetes state
	for _, route := range routesToRemove.List() {
		glog.Infof("deleting Kong Route %v", route)
		err := client.Routes().Delete(route)
		if err != nil {
			glog.Errorf("Unexpected error deleting Kong Route: %v", err)
		}
		// TODO: remove plugins from kong?
	}

	if len(routesToRemove) > 0 {
		triggerReload = true
	}

	return triggerReload, nil
}

// syncUpstreams synchronizes the state of Ingress backends against Kong upstreams
// This process only creates new Kong upstreams and synchronizes targets (Kubernetes endpoints)
func (n *NGINXController) syncUpstreams(locations []*ingress.Location, backends []*ingress.Backend) error {
	client := n.cfg.Kong.Client

	glog.V(3).Infof("syncing Kong upstreams")

	for _, location := range locations {
		backend := location.Backend
		if backend == "default-backend" {
			// there is no default backend in Kong
			continue
		}

		ingress := location.Ingress
		if ingress == nil {
			// location is the default backend (not mapped against Kong)
			continue
		}

		if backend == "" {
			glog.Warningf("the service defined in the ingress %v/%v does not exists", ingress.Namespace, ingress.Name)
			continue
		}

		kongIngress, err := n.getKongIngress(ingress)
		if err != nil {
			glog.V(5).Infof("there is no custom Ingress configuration for rule %v/%v", ingress.Namespace, ingress.Name)
		}

		for _, upstream := range backends {
			if upstream.Name != backend {
				continue
			}

			upstreamName := buildName(backend, location)

			_, res := client.Upstreams().Get(upstreamName)
			if res.StatusCode == http.StatusNotFound {
				upstream := kongadminv1.NewUpstream(upstreamName)

				if kongIngress != nil && kongIngress.Upstream != nil {
					if kongIngress.Upstream.HashOn != "" {
						upstream.HashOn = kongIngress.Upstream.HashOn
					}

					if kongIngress.Upstream.HashOnHeader != "" {
						upstream.HashOnHeader = kongIngress.Upstream.HashOnHeader
					}

					if kongIngress.Upstream.HashFallback != "" {
						upstream.HashFallback = kongIngress.Upstream.HashFallback
					}

					if kongIngress.Upstream.Slots != 0 {
						upstream.Slots = kongIngress.Upstream.Slots
					}

					if kongIngress.Upstream.Healthchecks != nil {
						if kongIngress.Upstream.Healthchecks.Active != nil {
							m := structs.Map(kongIngress.Upstream.Healthchecks.Active)
							if upstream.Healthchecks == nil {
								upstream.Healthchecks = &kongadminv1.Healthchecks{}
							}
							if upstream.Healthchecks.Active == nil {
								upstream.Healthchecks.Active = &kongadminv1.ActiveHealthCheck{}
							}

							mergo.MapWithOverwrite(upstream.Healthchecks.Active, m)
						}

						if kongIngress.Upstream.Healthchecks.Passive != nil {
							m := structs.Map(kongIngress.Upstream.Healthchecks.Passive)

							if upstream.Healthchecks == nil {
								upstream.Healthchecks = &kongadminv1.Healthchecks{}
							}
							if upstream.Healthchecks.Passive == nil {
								upstream.Healthchecks.Passive = &kongadminv1.Passive{}
							}

							mergo.MapWithOverwrite(upstream.Healthchecks.Passive, m)
						}
					}
				}

				glog.Infof("creating Kong Upstream with name %v", upstreamName)

				_, res := client.Upstreams().Create(upstream)
				if res.StatusCode != http.StatusCreated {
					glog.Errorf("Unexpected error creating Kong Upstream: %v", res)
					return res.Error()
				}
			}

			err := syncTargets(upstreamName, upstream, client)
			if err != nil {
				return errors.Wrap(err, "syncing targets")
			}

			//TODO: check if an update is required (change in kongIngress)
		}
	}

	return nil
}

// syncCertificate synchronizes the state of referenced secrets by Ingress
// rules with Kong certificates.
// This process only create or update certificates in Kong
func (n *NGINXController) syncCertificate(server *ingress.Server) error {
	if server.SSLCert == nil {
		return nil
	}

	if server.SSLCert.ID == "" {
		glog.Warningf("certificate %v/%v is invalid", server.SSLCert.Namespace, server.SSLCert.Name)
		return nil
	}

	client := n.cfg.Client

	sc := bytes.NewBuffer(server.SSLCert.Raw.Cert).String()
	sk := bytes.NewBuffer(server.SSLCert.Raw.Key).String()

	//sync cert
	cert, res := client.Certificates().Get(server.SSLCert.ID)

	switch res.StatusCode {
	case http.StatusOK:
		// check if an update is required
		name := fmt.Sprintf("temporal-cert-%v", time.Now().UnixNano())
		pem, err := ssl.AddOrUpdateCertAndKey(name,
			[]byte(strings.TrimSpace(cert.Cert)),
			[]byte(strings.TrimSpace(cert.Key)),
			[]byte{},
			n.fileSystem)
		if err != nil {
			return err
		}

		defer func() {
			os.Remove(pem.PemFileName)
			os.Remove(pem.FullChainPemFileName)
		}()

		if server.SSLCert.PemSHA != pem.PemSHA {
			glog.Infof("updating Kong SSL Certificate for host %v located in Secret %v/%v",
				server.Hostname, server.SSLCert.Namespace, server.SSLCert.Name)

			cert = &kongadminv1.Certificate{
				Cert: sc,
				Key:  sk,
			}
			cert, res = client.Certificates().Patch(server.SSLCert.ID, cert)
			if res.StatusCode != http.StatusOK {
				return errors.Wrap(res.Error(), "patching a Kong consumer")
			}
		}
	case http.StatusNotFound:
		cert = &kongadminv1.Certificate{
			Required: kongadminv1.Required{
				ID: server.SSLCert.ID,
			},
			Cert:  sc,
			Key:   sk,
			Hosts: []string{server.Hostname},
		}

		glog.Infof("creating Kong SSL Certificate for host %v located in Secret %v/%v",
			server.Hostname, server.SSLCert.Namespace, server.SSLCert.Name)

		cert, res = client.Certificates().Create(cert)
		if res.StatusCode != http.StatusCreated {
			glog.Errorf("Unexpected error creating Kong Certificate: %v", res)
			return res.Error()
		}
	default:
		glog.Errorf("Unexpected response searching a Kong Certificate: %v", res)
		return res.Error()
	}

	//sync SNI
	sni, res := client.SNIs().Get(server.Hostname)

	switch res.StatusCode {
	case http.StatusOK:
		// check if it is using the right certificate
		if sni.Certificate != cert.ID {
			glog.Infof("updating certificate for host %v to certificate id %v", server.Hostname, cert.ID)

			sni.Certificate = cert.ID
			sni, res = client.SNIs().Patch(sni.ID, sni)
			if res.StatusCode != http.StatusOK {
				return errors.Wrap(res.Error(), "patching a Kong consumer")
			}
		}
	case http.StatusNotFound:
		sni = &kongadminv1.SNI{
			Name:        server.Hostname,
			Certificate: cert.ID,
		}
		glog.Infof("creating Kong SNI for host %v and certificate id %v", server.Hostname, cert.ID)

		_, res = client.SNIs().Create(sni)
		if res.StatusCode != http.StatusCreated {
			glog.Errorf("Unexpected error creating Kong SNI (%v): %v", server.Hostname, res)
		}
	default:
		glog.Errorf("Unexpected error looking up Kong SNI: %v", res)
		return res.Error()
	}

	return nil
}

// buildName returns a string valid as a hostnames taking a backend and
// location as input. The format of backend is <namespace>-<service name>-<port>
// For the location the field Path is used. If the path is / only the backend is used
// This process ensures the returned name is unique.
func buildName(backend string, location *ingress.Location) string {
	return backend
}

// getKongService returns a Route from a list using the path and hosts as filters
func getKongRoute(hostname, path string, routes []kongadminv1.Route) *kongadminv1.Route {
	for _, r := range routes {
		if sets.NewString(r.Paths...).Has(path) &&
			sets.NewString(r.Hosts...).Has(hostname) {
			return &r
		}
	}

	return nil
}

// isRouteInKong checks if a route exists or not in Kong
func isRouteInKong(route *kongadminv1.Route, routes []kongadminv1.Route) bool {
	for _, eRoute := range routes {
		if route.Equal(&eRoute) {
			return true
		}
	}

	return false
}

// deleteServiceUpstream deletes multiple Kong upstreams for a particular
// Kong service. This process requires the removal of all the targets that
// reference the upstream to be removed
func deleteServiceUpstream(host string, client *kong.RestClient) error {
	kongUpstreams, err := client.Upstreams().List(nil)
	if err != nil {
		return err
	}

	upstreamsToRemove := sets.NewString()
	for _, upstream := range kongUpstreams.Items {
		if upstream.Name == host {
			if !upstreamsToRemove.Has(upstream.ID) {
				upstreamsToRemove.Insert(upstream.ID)
			}
		}
	}

	for _, upstream := range upstreamsToRemove.List() {
		kongTargets, err := client.Targets().List(nil, upstream)
		if err != nil {
			return err
		}

		for _, target := range kongTargets.Items {
			if target.Upstream == upstream {
				err := client.Targets().Delete(target.ID, upstream)
				if err != nil {
					return errors.Wrap(err, "removing a Kong target")
				}
			}
		}

		err = client.Upstreams().Delete(upstream)
		if err != nil {
			return errors.Wrap(err, "removing a Kong upstream")
		}
	}

	return nil
}

// pluginDeepEqual compares the configuration of a Plugin (CRD) against
// the persisted state in the Kong database
// This is required because a plugin has defaults that could not exists in the CRD.
func pluginDeepEqual(config map[string]interface{}, kong *kongadminv1.Plugin) bool {
	for k, v := range config {
		kv, ok := kong.Config[k]
		if !ok {
			return false
		}

		if !reflect.DeepEqual(v, kv) {
			return false
		}
	}

	return true
}

// getKongIngress checks if the Ingress contains an annotation for configuration
// or if exists a KongIngress object with the same name than the Ingress
func (n *NGINXController) getKongIngress(ing *extensions.Ingress) (*configurationv1.KongIngress, error) {
	confName := annotations.ExtractConfigurationName(ing.Annotations)
	if confName != "" {
		return n.store.GetKongIngress(ing.Namespace, confName)
	}

	return n.store.GetKongIngress(ing.Namespace, ing.Name)
}
