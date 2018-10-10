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
	"net"
	"net/http"
	"reflect"
	"sort"
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
		name := targetGlobalPlugins[i].PluginName
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
		nt := net.JoinHostPort(endpoint.Address, endpoint.Port)
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

// getPluginsFromAnnotations extracts plugins to be applied on an ingress/service from annotations
func (n *NGINXController) getPluginsFromAnnotations(namespace string, anns map[string]string) (map[string]*pluginv1.KongPlugin, error) {
	pluginAnnotations := annotations.ExtractKongPluginAnnotations(anns)
	pluginsInk8s := make(map[string]*pluginv1.KongPlugin)
	for plugin, crdNames := range pluginAnnotations {
		for _, crdName := range crdNames {
			// search configured plugin CRD in k8s
			k8sPlugin, err := n.store.GetKongPlugin(namespace, crdName)
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("searching plugin KongPlugin %v", crdName))
			}
			pluginsInk8s[plugin] = k8sPlugin
		}
	}
	pluginList := annotations.ExtractKongPluginsFromAnnotations(anns)
	// override plugins configured by new annotation
	for _, plugin := range pluginList {
		k8sPlugin, err := n.store.GetKongPlugin(namespace, plugin)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("searching plugin KongPlugin %v", plugin))
		}
		// ignore plugins with no name
		if k8sPlugin.PluginName == "" {
			glog.Errorf("KongPlugin Custom resource '%v' has no `plugin` property, the plugin will not be configured", k8sPlugin.Name)
			continue
		}
		pluginsInk8s[k8sPlugin.PluginName] = k8sPlugin
	}
	return pluginsInk8s, nil
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
					// the flag that if we need to patch KongService later
					outOfSync := false
					if kongIngress != nil && kongIngress.Proxy != nil {
						if kongIngress.Proxy.Path != "" && s.Path != kongIngress.Proxy.Path {
							s.Path = kongIngress.Proxy.Path
							outOfSync = true
						}

						if kongIngress.Proxy.Protocol != "" &&
							(kongIngress.Proxy.Protocol == "http" || kongIngress.Proxy.Protocol == "https") &&
							s.Protocol != kongIngress.Proxy.Protocol {
							s.Protocol = kongIngress.Proxy.Protocol
							outOfSync = true
						}

						if kongIngress.Proxy.ConnectTimeout > 0 &&
							s.ConnectTimeout != kongIngress.Proxy.ConnectTimeout {
							s.ConnectTimeout = kongIngress.Proxy.ConnectTimeout
							outOfSync = true
						}

						if kongIngress.Proxy.ReadTimeout > 0 &&
							s.ReadTimeout != kongIngress.Proxy.ReadTimeout {
							s.ReadTimeout = kongIngress.Proxy.ReadTimeout
							outOfSync = true
						}

						if kongIngress.Proxy.WriteTimeout > 0 &&
							s.WriteTimeout != kongIngress.Proxy.WriteTimeout {
							s.WriteTimeout = kongIngress.Proxy.WriteTimeout
							outOfSync = true
						}

						if kongIngress.Proxy.Retries > 0 && s.Retries != kongIngress.Proxy.Retries {
							s.Retries = kongIngress.Proxy.Retries
							outOfSync = true
						}
					}

					if res.StatusCode == http.StatusNotFound {
						glog.Infof("Creating Kong Service name %v", name)
						_, res := client.Services().Create(s)
						if res.StatusCode != http.StatusCreated {
							glog.Errorf("Unexpected error creating Kong Service: %v", res)
							return false, res.Error()
						}
					} else if outOfSync {
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

			// Get plugin annotations from k8s, these plugins should be configured for this service
			anns := location.Service.GetAnnotations()
			pluginsInk8s, err := n.getPluginsFromAnnotations(location.Ingress.Namespace, anns)

			// Get plugins configured in Kong currently
			plugins, err := client.Plugins().GetAllByService(svc.ID)
			if err != nil {
				glog.Errorf("Unexpected error obtaining Kong plugins for service %v: %v", svc.ID, err)
				continue
			}
			pluginsInKong := make(map[string]kongadminv1.Plugin)
			for _, plugin := range plugins {
				pluginsInKong[plugin.Name] = plugin
			}
			var pluginsToDelete []kongadminv1.Plugin
			var pluginsToCreate []kongadminv1.Plugin
			var pluginsToUpdate []kongadminv1.Plugin

			// Plugins present in Kong but not in k8s should be deleted
			for _, plugin := range pluginsInKong {
				if _, ok := pluginsInk8s[plugin.Name]; !ok {
					pluginsToDelete = append(pluginsToDelete, plugin)
				}
			}

			// Plugins present in k8s but not in Kong should be created
			for pluginName, pluginInk8s := range pluginsInk8s {
				if pluginInKong, ok := pluginsInKong[pluginName]; !ok {
					p := kongadminv1.Plugin{
						Name:    pluginName,
						Service: svc.ID,
						Config:  kongadminv1.Configuration(pluginInk8s.Config),
					}
					pluginsToCreate = append(pluginsToCreate, p)
				} else {
					// Plugins present in Kong and k8s need should have same configuration
					if !pluginDeepEqual(pluginInk8s.Config, &pluginInKong) ||
						(!pluginInk8s.Disabled != pluginInKong.Enabled) ||
						(pluginInk8s.ConsumerRef != "" && pluginInKong.Consumer == "") {
						glog.Infof("plugin %v configuration in kong is outdated. Updating...", pluginInk8s.Name)
						p := kongadminv1.Plugin{
							Name:    pluginInKong.Name,
							Config:  kongadminv1.Configuration(pluginInk8s.Config),
							Enabled: pluginInk8s.Disabled,
							Service: svc.ID,
							Required: kongadminv1.Required{
								ID: pluginInKong.ID,
							},
						}
						if pluginInk8s.ConsumerRef != "" {
							consumer, err := n.store.GetKongConsumer(pluginInk8s.Namespace, pluginInk8s.ConsumerRef)
							if err != nil {
								glog.Errorf("Unexpected error searching for consumer %v: %v",
									pluginInk8s.ConsumerRef, err)
								continue
							}
							p.Consumer = fmt.Sprintf("%v", consumer.GetUID())
						}
						pluginsToUpdate = append(pluginsToUpdate, p)
					}
				}
			}

			for _, plugin := range pluginsToDelete {
				err := client.Plugins().Delete(plugin.ID)
				if err != nil {
					return false, errors.Wrap(err, "deleting Kong plugin")
				}
			}

			for _, plugin := range pluginsToCreate {
				_, res := client.Plugins().CreateInService(svc.ID, &plugin)
				if res.StatusCode != http.StatusCreated {
					return false, errors.Wrap(res.Error(), fmt.Sprintf("creating a Kong plugin %v in service %v", plugin, svc.ID))
				}
			}

			for _, plugin := range pluginsToUpdate {
				_, res := client.Plugins().Patch(plugin.Required.ID, &plugin)
				if res.StatusCode != http.StatusOK {
					return false, errors.Wrap(res.Error(), fmt.Sprintf("updating a Kong plugin %v in service %v", plugin, svc.ID))
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
				Service:   kongadminv1.InlineService{ID: svc.ID},
			}
			if server.Hostname != "_" {
				r.Hosts = []string{server.Hostname}
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
							glog.Infof("protocols changed form %v to %v", route.Protocols, kongIngress.Route.Protocols)
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

			if route == nil {
				continue
			}

			anns := location.Ingress.GetAnnotations()
			pluginsInk8s, err := n.getPluginsFromAnnotations(location.Ingress.Namespace, anns)

			// Get plugins configured in Kong currently
			plugins, err := client.Plugins().GetAllByRoute(route.ID)
			if err != nil {
				glog.Errorf("Unexpected error obtaining Kong plugins for route %v: %v", route.ID, err)
				continue
			}
			pluginsInKong := make(map[string]kongadminv1.Plugin)
			for _, plugin := range plugins {
				pluginsInKong[plugin.Name] = plugin
			}
			var pluginsToDelete []kongadminv1.Plugin
			var pluginsToCreate []kongadminv1.Plugin
			var pluginsToUpdate []kongadminv1.Plugin

			// Plugins present in Kong but not in k8s should be deleted
			for _, plugin := range pluginsInKong {
				if _, ok := pluginsInk8s[plugin.Name]; !ok {
					pluginsToDelete = append(pluginsToDelete, plugin)
				}
			}
			// Plugins present in k8s but not in Kong should be created
			for pluginName, pluginInk8s := range pluginsInk8s {
				if pluginInKong, ok := pluginsInKong[pluginName]; !ok {
					p := kongadminv1.Plugin{
						Name:   pluginName,
						Route:  route.ID,
						Config: kongadminv1.Configuration(pluginInk8s.Config),
					}
					pluginsToCreate = append(pluginsToCreate, p)
				} else {
					// Plugins present in Kong and k8s need should have same configuration
					if !pluginDeepEqual(pluginInk8s.Config, &pluginInKong) ||
						(!pluginInk8s.Disabled != pluginInKong.Enabled) ||
						(pluginInk8s.ConsumerRef != "" && pluginInKong.Consumer == "") {
						glog.Infof("plugin %v configuration in kong is outdated. Updating...", pluginInk8s.Name)
						p := kongadminv1.Plugin{
							Name:    pluginInKong.Name,
							Config:  kongadminv1.Configuration(pluginInk8s.Config),
							Enabled: pluginInk8s.Disabled,
							Route:   route.ID,
							Required: kongadminv1.Required{
								ID: pluginInKong.ID,
							},
						}

						if pluginInk8s.ConsumerRef != "" {
							consumer, err := n.store.GetKongConsumer(pluginInk8s.Namespace, pluginInk8s.ConsumerRef)
							if err != nil {
								glog.Errorf("Unexpected error searching for consumer %v: %v",
									pluginInk8s.ConsumerRef, err)
								continue
							}
							p.Consumer = fmt.Sprintf("%v", consumer.GetUID())
						}
						pluginsToUpdate = append(pluginsToUpdate, p)
					}
					glog.Infof("plugin %v configuration in kong is up to date.", pluginInk8s.PluginName)

				}
			}
			for _, plugin := range pluginsToDelete {
				err := client.Plugins().Delete(plugin.ID)
				if err != nil {
					return false, errors.Wrap(err, "deleting Kong plugin")
				}
			}

			for _, plugin := range pluginsToCreate {
				_, res := client.Plugins().CreateInRoute(route.ID, &plugin)
				if res.StatusCode != http.StatusCreated {
					return false, errors.Wrap(res.Error(), fmt.Sprintf("creating a Kong plugin %v in service %v", plugin, svc.ID))
				}
			}

			for _, plugin := range pluginsToUpdate {
				_, res := client.Plugins().Patch(plugin.Required.ID, &plugin)
				if res.StatusCode != http.StatusOK {
					return false, errors.Wrap(res.Error(), fmt.Sprintf("updating a Kong plugin %v in service %v", plugin, svc.ID))
				}
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

					if kongIngress.Upstream.HashOnCookie != "" {
						upstream.HashOnCookie = kongIngress.Upstream.HashOnCookie
					}

					if kongIngress.Upstream.HashOnCookiePath != "" {
						upstream.HashOnCookiePath = kongIngress.Upstream.HashOnCookiePath
					}

					if kongIngress.Upstream.HashOnHeader != "" {
						upstream.HashOnHeader = kongIngress.Upstream.HashOnHeader
					}

					if kongIngress.Upstream.HashFallback != "" {
						upstream.HashFallback = kongIngress.Upstream.HashFallback
					}

					if kongIngress.Upstream.HashFallbackHeader != "" {
						upstream.HashFallbackHeader = kongIngress.Upstream.HashFallbackHeader
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

		if cert.Cert != sc || cert.Key != sk {
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
		if sni.Certificate.ID != cert.ID {
			glog.Infof("updating certificate for host %v to certificate id %v", server.Hostname, cert.ID)

			sni.Certificate.ID = cert.ID
			sni, res = client.SNIs().Patch(sni.ID, sni)
			if res.StatusCode != http.StatusOK {
				return errors.Wrap(res.Error(), "patching a Kong consumer")
			}
		}
	case http.StatusNotFound:
		sni = &kongadminv1.SNI{
			Name:        server.Hostname,
			Certificate: kongadminv1.InlineCertificate{ID: cert.ID},
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
		if hostname != "_" {
			if sets.NewString(r.Paths...).Has(path) &&
				sets.NewString(r.Hosts...).Has(hostname) {
				return &r
			}
		} else {
			if sets.NewString(r.Paths...).Has(path) {
				return &r
			}
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
