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
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"sort"
	"time"

	"github.com/fatih/structs"
	"github.com/golang/glog"
	"github.com/hbagdi/go-kong/kong"
	"github.com/hbagdi/go-kong/kong/custom"
	"github.com/imdario/mergo"
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

	err := n.syncCertificates(ingressCfg.Servers)
	if err != nil {
		return err
	}

	for _, server := range ingressCfg.Servers {
		if server.Hostname == "_" {
			// there is no catch all server in kong
			continue
		}

		err := n.syncUpstreams(server.Locations, ingressCfg.Backends)
		if err != nil {
			return err
		}
	}

	err = n.syncConsumers()
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
	plugins, err := client.Plugins.ListAll(nil)
	if err != nil {
		return err
	}

	// plugins in Kong
	currentGlobalPlugins := make(map[string]kong.Plugin)
	for _, plugin := range plugins {
		if isEmpty(plugin.RouteID) && isEmpty(plugin.ServiceID) && isEmpty(plugin.ConsumerID) {
			currentGlobalPlugins[*plugin.Name] = *plugin
		}
	}

	// sync plugins to Kong
	for pluginName, kongPlugin := range targetPluginMap {
		// plugin exists?
		if pluginInKong, ok := currentGlobalPlugins[pluginName]; !ok {
			// no, create it
			p := &kong.Plugin{
				Name:   kong.String(pluginName),
				Config: kong.Configuration(kongPlugin.Config),
			}
			_, err := client.Plugins.Create(nil, p)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("creating a global Kong plugin %v", p))
			}
		} else {
			// plugin exists, is the configuration up to date
			if !pluginDeepEqual(kongPlugin.Config, &pluginInKong) {
				// no, update it
				p := &kong.Plugin{
					ID:     kong.String(*pluginInKong.ID),
					Name:   kong.String(pluginName),
					Config: kong.Configuration(kongPlugin.Config),
				}
				_, err := client.Plugins.Update(nil, p)
				if err != nil {
					return errors.Wrap(err, fmt.Sprintf("updating a global Kong plugin %v", p))
				}
			}
		}
		// remove from the current list, all that remain in the current list will be deleted
		delete(currentGlobalPlugins, pluginName)
	}

	// delete the ones not configured in k8s
	for _, plugin := range currentGlobalPlugins {
		err := client.Plugins.Delete(nil, plugin.ID)
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
func syncTargets(upstream string, ingressEndpopint *ingress.Backend, client *kong.Client) error {
	glog.V(3).Infof("syncing Kong targets")
	b, err := client.Upstreams.Get(nil, &upstream)
	if kong.IsNotFoundErr(err) {
		glog.Errorf("there is no upstream with name %v in Kong", upstream)
		return nil
	} else if err != nil {
		glog.Errorf("err fetching upstream %v from Kong: %v", upstream, err)
		return err
	}

	kongTargets, err := client.Targets.ListAll(nil, &upstream)
	if err != nil {
		return err
	}

	oldTargets := sets.NewString()
	for _, kongTarget := range kongTargets {
		if !oldTargets.Has(*kongTarget.Target) {
			oldTargets.Insert(*kongTarget.Target)
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
		target := &kong.Target{
			Target:     kong.String(endpoint),
			UpstreamID: b.ID,
		}
		glog.Infof("creating Kong Target %v for upstream %v", endpoint, b.ID)
		_, err := client.Targets.Create(nil, &upstream, target)
		if err != nil {
			glog.Errorf("Unexpected error creating Kong Upstream: %v", err)
			return err
		}
	}

	// wait to avoid hitting the kong API server too fast
	time.Sleep(100 * time.Millisecond)

	for _, endpoint := range remove {
		for _, kongTarget := range kongTargets {
			if *kongTarget.Target != endpoint {
				continue
			}
			glog.Infof("deleting Kong Target %v from upstream %v", kongTarget.ID, kongTarget.UpstreamID)
			err := client.Targets.Delete(nil, b.ID, kongTarget.ID)
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

				s, err := client.Services.Get(nil, &name)
				if kong.IsNotFoundErr(err) {
					s = &kong.Service{
						Name:     kong.String(name),
						Path:     kong.String("/"),
						Protocol: kong.String(proto),
						Host:     kong.String(name),
						Port:     kong.Int(port),
						Retries:  kong.Int(5),
					}
				}

				if s != nil {
					// the flag that if we need to patch KongService later
					outOfSync := false
					if kongIngress != nil && kongIngress.Proxy != nil {
						if kongIngress.Proxy.Path != "" && (s.Path == nil || *s.Path != kongIngress.Proxy.Path) {
							s.Path = kong.String(kongIngress.Proxy.Path)
							outOfSync = true
						}

						if kongIngress.Proxy.Protocol != "" &&
							(kongIngress.Proxy.Protocol == "http" || kongIngress.Proxy.Protocol == "https") &&
							(s.Path == nil || *s.Protocol != kongIngress.Proxy.Protocol) {
							s.Protocol = kong.String(kongIngress.Proxy.Protocol)
							outOfSync = true
						}

						if kongIngress.Proxy.ConnectTimeout > 0 &&
							(s.ConnectTimeout == nil || *s.ConnectTimeout != kongIngress.Proxy.ConnectTimeout) {
							s.ConnectTimeout = kong.Int(kongIngress.Proxy.ConnectTimeout)
							outOfSync = true
						}

						if kongIngress.Proxy.ReadTimeout > 0 &&
							(s.ReadTimeout == nil || *s.ReadTimeout != kongIngress.Proxy.ReadTimeout) {
							s.ReadTimeout = kong.Int(kongIngress.Proxy.ReadTimeout)
							outOfSync = true
						}

						if kongIngress.Proxy.WriteTimeout > 0 &&
							(s.WriteTimeout == nil || *s.WriteTimeout != kongIngress.Proxy.WriteTimeout) {
							s.WriteTimeout = kong.Int(kongIngress.Proxy.WriteTimeout)
							outOfSync = true
						}

						if kongIngress.Proxy.Retries > 0 &&
							(s.Retries == nil || *s.Retries != kongIngress.Proxy.Retries) {
							s.Retries = kong.Int(kongIngress.Proxy.Retries)
							outOfSync = true
						}
					}
					if kong.IsNotFoundErr(err) {
						glog.Infof("Creating Kong Service name %v", name)
						_, err := client.Services.Create(nil, s)
						if err != nil {
							glog.Errorf("Unexpected error creating Kong Service: %v", err)
							return false, err
						}
					} else if outOfSync {
						glog.Infof("Patching Kong Service name %v", name)
						_, err := client.Services.Update(nil, s)
						if err != nil {
							glog.Errorf("Unexpected error patching Kong Service: %v", err)
							return false, err
						}
					}
				}

				// do not remove the service
				if !servicesToKeep.Has(name) {
					servicesToKeep.Insert(name)
				}

				break
			}

			svc, err := client.Services.Get(nil, &name)
			if err != nil {
				glog.Warningf("service %v does not exists in kong", name)
				continue
			}

			// Get plugin annotations from k8s, these plugins should be configured for this service
			anns := location.Service.GetAnnotations()
			pluginsInk8s, err := n.getPluginsFromAnnotations(location.Ingress.Namespace, anns)

			// Get plugins configured in Kong currently
			plugins, err := client.Plugins.ListAllForService(nil, svc.ID)
			if err != nil {
				glog.Errorf("Unexpected error obtaining Kong plugins for service %v: %v", svc.ID, err)
				continue
			}
			pluginsInKong := make(map[string]kong.Plugin)
			for _, plugin := range plugins {
				pluginsInKong[*plugin.Name] = *plugin
			}
			var pluginsToDelete []kong.Plugin
			var pluginsToCreate []kong.Plugin
			var pluginsToUpdate []kong.Plugin

			// Plugins present in Kong but not in k8s should be deleted
			for _, plugin := range pluginsInKong {
				if _, ok := pluginsInk8s[*plugin.Name]; !ok {
					pluginsToDelete = append(pluginsToDelete, plugin)
				}
			}

			// Plugins present in k8s but not in Kong should be created
			for pluginName, pluginInk8s := range pluginsInk8s {
				if pluginInKong, ok := pluginsInKong[pluginName]; !ok {
					p := kong.Plugin{
						Name:      kong.String(pluginName),
						ServiceID: kong.String(*svc.ID),
						Config:    kong.Configuration(pluginInk8s.Config),
					}
					pluginsToCreate = append(pluginsToCreate, p)
				} else {
					enabled := false
					if pluginInKong.Enabled != nil {
						enabled = *pluginInKong.Enabled
					}
					// Plugins present in Kong and k8s need should have same configuration
					if !pluginDeepEqual(pluginInk8s.Config, &pluginInKong) || // plugin conf
						(pluginInk8s.ConsumerRef != "" && isEmpty(pluginInKong.ConsumerID)) || // consumerID
						// plugin disabled?
						(!pluginInk8s.Disabled != enabled) {
						glog.Infof("plugin %v configuration in kong is outdated. Updating...", pluginInk8s.Name)
						p := kong.Plugin{
							Name:   pluginInKong.Name,
							Config: kong.Configuration(pluginInk8s.Config),

							Enabled:   kong.Bool(!pluginInk8s.Disabled),
							ServiceID: kong.String(*svc.ID),
							ID:        kong.String(*pluginInKong.ID),
						}
						if pluginInk8s.ConsumerRef != "" {
							consumer, err := n.store.GetKongConsumer(pluginInk8s.Namespace, pluginInk8s.ConsumerRef)
							if err != nil {
								glog.Errorf("Unexpected error searching for consumer %v: %v",
									pluginInk8s.ConsumerRef, err)
								continue
							}
							p.ConsumerID = kong.String(fmt.Sprintf("%v", consumer.GetUID()))
						}
						pluginsToUpdate = append(pluginsToUpdate, p)
					}
				}
			}

			for _, plugin := range pluginsToDelete {
				err := client.Plugins.Delete(nil, plugin.ID)
				if err != nil {
					return false, errors.Wrap(err, "deleting Kong plugin")
				}
			}

			for _, plugin := range pluginsToCreate {
				_, err := client.Plugins.Create(nil, &plugin)
				if err != nil {
					return false, errors.Wrap(err, fmt.Sprintf("creating a Kong plugin %v in service %v", plugin, svc.ID))
				}
			}

			for _, plugin := range pluginsToUpdate {
				_, err := client.Plugins.Update(nil, &plugin)
				if err != nil {
					return false, errors.Wrap(err, fmt.Sprintf("updating a Kong plugin %v in service %v", plugin, svc.ID))
				}
			}
		}
	}

	kongServices, err := client.Services.ListAll(nil)
	if err != nil {
		return false, err
	}

	serviceNames := sets.NewString()

	for _, svc := range kongServices {
		if !serviceNames.Has(*svc.Name) {
			serviceNames.Insert(*svc.Name)
		}
	}

	servicesToRemove := serviceNames.Difference(servicesToKeep)

	// remove all those services that are present in Kong but not in the current Kubernetes state
	for _, svcName := range servicesToRemove.List() {
		svc, err := client.Services.Get(nil, &svcName)
		if kong.IsNotFoundErr(err) {
			glog.Warningf("service %v does not exists in kong", svcName)
			continue
		} else if err != nil {
			glog.Errorf("Unexpected error looking up service in Kong: %v", svc.Name)
			continue
		}

		glog.Infof("deleting Kong Service %v", svcName)
		// before deleting the service we need to remove the upstream and th targets that reference the service
		err = deleteServiceUpstream(*svc.Name, client)
		if err != nil {
			glog.Errorf("Unexpected error deleting Kong upstreams and targets that depend on service %v: %v", svc.Name, err)
			continue
		}
		err = client.Services.Delete(nil, svc.ID)
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

	consumersInKong := make(map[string]*kong.Consumer)
	client := n.cfg.Kong.Client

	// List all consumers in Kong
	kongConsumers, err := client.Consumers.ListAll(nil)
	if err != nil {
		return err
	}

	for i := range kongConsumers {
		consumersInKong[*kongConsumers[i].ID] = kongConsumers[i]
	}

	// List existing Consumers in Kubernetes
	for _, consumer := range n.store.ListKongConsumers() {
		glog.Infof("checking if Kong consumer %v exists", consumer.Name)
		consumerID := fmt.Sprintf("%v", consumer.GetUID())

		kc, ok := consumersInKong[consumerID]

		if !ok {
			glog.Infof("Creating Kong consumer %v", consumerID)
			c := &kong.Consumer{
				ID: kong.String(consumerID),
			}
			if consumer.Username != "" {
				c.Username = kong.String(consumer.Username)
			}
			if consumer.CustomID != "" {
				c.CustomID = kong.String(consumer.CustomID)
			}
			c, err := n.cfg.Kong.Client.Consumers.Create(nil, c)
			if err != nil {
				return errors.Wrap(err, "creating a Kong consumer")
			}
		} else {
			// check the consumers are equals
			outOfSync := false
			if consumer.Username != "" && (kc.Username == nil || *kc.Username != consumer.Username) {
				outOfSync = true
				kc.Username = kong.String(consumer.Username)
			}
			if consumer.CustomID != "" && (kc.CustomID == nil || *kc.CustomID == consumer.CustomID) {
				outOfSync = true
				kc.CustomID = kong.String(consumer.CustomID)
			}
			if outOfSync {
				_, err := n.cfg.Kong.Client.Consumers.Update(nil, kc)
				if err != nil {
					return errors.Wrap(err, "patching a Kong consumer")
				}
			}
		}
		delete(consumersInKong, consumerID)
	}
	// remaining entries in the map should be deleted

	for _, consumer := range consumersInKong {
		err := client.Consumers.Delete(nil, consumer.ID)
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
	client := n.cfg.Kong.Client

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
		credInKong := custom.NewEntityObject(custom.Type(credential.Type))
		credInKong.AddRelation("consumer_id", consumerID)
		credInKong.SetObject(map[string]interface{}{"id": credentialID})
		_, err = client.CustomEntities.Get(nil, credInKong)
		if kong.IsNotFoundErr(err) {
			// use the configuration
			data := credential.Config
			if data == nil {
				data = make(map[string]interface{})
			}
			// create a credential with the same id of k8s
			data["id"] = credentialID
			credInKong.SetObject(map[string]interface{}(data))
			_, err := client.CustomEntities.Create(nil, credInKong)
			if err != nil {
				glog.Errorf("Unexpected error creating credential: %v", err)
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

// syncRoutes synchronizes the state between the Ingress configuration model and Kong Routes.
func (n *NGINXController) syncRoutes(ingressCfg *ingress.Configuration) (bool, error) {
	client := n.cfg.Kong.Client

	kongRoutes, err := client.Routes.ListAll(nil)
	if err != nil {
		return false, err
	}

	// triggerReload indicates if the sync process altered
	// configuration with services and require an additional run
	var triggerReload bool

	// create a copy of the existing routes to be able to run a comparison
	routesToRemove := sets.NewString()
	for _, old := range kongRoutes {
		if !routesToRemove.Has(*old.ID) {
			routesToRemove.Insert(*old.ID)
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
			svc, err := client.Services.Get(nil, kong.String(name))
			if kong.IsNotFoundErr(err) || svc == nil || isEmpty(svc.ID) {
				glog.Warningf("service %v does not exists in kong", name)
				continue
			}

			r := &kong.Route{
				Paths:     []*string{&location.Path},
				Protocols: []*string{kong.String("http"), kong.String("https")}, // default
				Service:   &kong.Service{ID: svc.ID},
			}
			if server.Hostname != "_" {
				r.Hosts = []*string{kong.String(server.Hostname)}
			}

			if kongIngress != nil && kongIngress.Route != nil {
				// ignore updated on create
				mergeRouteAndKongIngress(r, kongIngress)
			}

			if !isRouteInKong(r, kongRoutes) {
				glog.Infof("creating Kong Route for host %v, path %v and service %v", server.Hostname, location.Path, *svc.ID)
				_, err := client.Routes.Create(nil, r)
				if err != nil {
					glog.Errorf("Unexpected error creating Kong Route: %v", err)
					return false, err
				}
			} else {
				// the route exists but the properties could have changed
				route := getKongRoute(server.Hostname, location.Path, kongRoutes)

				if routesToRemove.Has(*route.ID) {
					routesToRemove.Delete(*route.ID)
				}
				routesOutOfSync, ingressOutOfSync := false, false
				sort.Strings(protos)
				p := toStringArray(route.Protocols)
				sort.Strings(p)

				if !reflect.DeepEqual(protos, p) {
					routesOutOfSync = true
					route.Protocols = toStringPtrArray(protos)
				}

				if kongIngress != nil && kongIngress.Route != nil {
					ingressOutOfSync = mergeRouteAndKongIngress(route, kongIngress)
				}

				if routesOutOfSync || ingressOutOfSync {
					glog.Infof("updating Kong Route for host %v, path %v and service %v",
						server.Hostname, location.Path, svc.ID)
					_, err := client.Routes.Update(nil, route)
					if err != nil {
						glog.Errorf("Unexpected error updating Kong Route: %v", err)
						return false, err
					}
				}
			}

			kongRoutes, err = client.Routes.ListAll(nil)
			if err != nil {
				return false, err
			}

			route := getKongRoute(server.Hostname, location.Path, kongRoutes)

			if route == nil {
				continue
			}

			anns := location.Ingress.GetAnnotations()
			pluginsInk8s, err := n.getPluginsFromAnnotations(location.Ingress.Namespace, anns)

			// Get plugins configured in Kong currently
			plugins, err := client.Plugins.ListAllForRoute(nil, route.ID)
			if err != nil {
				glog.Errorf("Unexpected error obtaining Kong plugins for route %v: %v", route.ID, err)
				continue
			}
			pluginsInKong := make(map[string]kong.Plugin)
			for _, plugin := range plugins {
				pluginsInKong[*plugin.Name] = *plugin
			}
			var pluginsToDelete []kong.Plugin
			var pluginsToCreate []kong.Plugin
			var pluginsToUpdate []kong.Plugin

			// Plugins present in Kong but not in k8s should be deleted
			for _, plugin := range pluginsInKong {
				if _, ok := pluginsInk8s[*plugin.Name]; !ok {
					pluginsToDelete = append(pluginsToDelete, plugin)
				}
			}
			// Plugins present in k8s but not in Kong should be created
			for pluginName, pluginInk8s := range pluginsInk8s {
				if pluginInKong, ok := pluginsInKong[pluginName]; !ok {
					p := kong.Plugin{
						Name:    kong.String(pluginName),
						RouteID: kong.String(*route.ID),
						Config:  kong.Configuration(pluginInk8s.Config),
					}
					pluginsToCreate = append(pluginsToCreate, p)
				} else {
					if pluginInKong.Enabled == nil {
						pluginInKong.Enabled = kong.Bool(false)
					}
					// Plugins present in Kong and k8s need should have same configuration
					if !pluginDeepEqual(pluginInk8s.Config, &pluginInKong) ||
						(!pluginInk8s.Disabled != *pluginInKong.Enabled) ||
						(pluginInk8s.ConsumerRef != "" &&
							(pluginInKong.ConsumerID == nil || *pluginInKong.ConsumerID != pluginInk8s.ConsumerRef)) {
						glog.Infof("plugin %v configuration in kong is outdated. Updating...", pluginInk8s.Name)
						p := kong.Plugin{
							Name:    pluginInKong.Name,
							Config:  kong.Configuration(pluginInk8s.Config),
							Enabled: kong.Bool(!pluginInk8s.Disabled),
							RouteID: route.ID,
							ID:      pluginInKong.ID,
						}

						if pluginInk8s.ConsumerRef != "" {
							consumer, err := n.store.GetKongConsumer(pluginInk8s.Namespace, pluginInk8s.ConsumerRef)
							if err != nil {
								glog.Errorf("Unexpected error searching for consumer %v: %v",
									pluginInk8s.ConsumerRef, err)
								continue
							}
							p.ConsumerID = kong.String(fmt.Sprintf("%v", consumer.GetUID()))
						}
						pluginsToUpdate = append(pluginsToUpdate, p)
					}
					glog.Infof("plugin %v configuration in kong is up to date.", pluginInk8s.PluginName)

				}
			}
			for _, plugin := range pluginsToDelete {
				err := client.Plugins.Delete(nil, plugin.ID)
				if err != nil {
					return false, errors.Wrap(err, "deleting Kong plugin")
				}
			}

			for _, plugin := range pluginsToCreate {
				_, err := client.Plugins.Create(nil, &plugin)
				if err != nil {
					return false, errors.Wrap(err, fmt.Sprintf("creating a Kong plugin %v in service %v", plugin, svc.ID))
				}
			}

			for _, plugin := range pluginsToUpdate {
				_, err := client.Plugins.Update(nil, &plugin)
				if err != nil {
					return false, errors.Wrap(err, fmt.Sprintf("updating a Kong plugin %v in service %v", plugin, svc.ID))
				}
			}
		}
	}

	// remove all those routes that are present in Kong but not in the current Kubernetes state
	for _, route := range routesToRemove.List() {
		glog.Infof("deleting Kong Route %v", route)
		err := client.Routes.Delete(nil, &route)
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

			_, err := client.Upstreams.Get(nil, &upstreamName)
			if kong.IsNotFoundErr(err) {
				upstream := &kong.Upstream{Name: &upstreamName}

				if kongIngress != nil && kongIngress.Upstream != nil {
					if kongIngress.Upstream.HashOn != "" {
						upstream.HashOn = kong.String(kongIngress.Upstream.HashOn)
					}

					if kongIngress.Upstream.HashOnCookie != "" {
						upstream.HashOnCookie = kong.String(kongIngress.Upstream.HashOnCookie)
					}

					if kongIngress.Upstream.HashOnCookiePath != "" {
						upstream.HashOnCookiePath = kong.String(kongIngress.Upstream.HashOnCookiePath)
					}

					if kongIngress.Upstream.HashOnHeader != "" {
						upstream.HashOnHeader = kong.String(kongIngress.Upstream.HashOnHeader)
					}

					if kongIngress.Upstream.HashFallback != "" {
						upstream.HashFallback = kong.String(kongIngress.Upstream.HashFallback)
					}

					if kongIngress.Upstream.HashFallbackHeader != "" {
						upstream.HashFallbackHeader = kong.String(kongIngress.Upstream.HashFallbackHeader)
					}

					if kongIngress.Upstream.Slots != 0 {
						upstream.Slots = kong.Int(kongIngress.Upstream.Slots)
					}

					if kongIngress.Upstream.Healthchecks != nil {
						if kongIngress.Upstream.Healthchecks.Active != nil {
							m := structs.Map(kongIngress.Upstream.Healthchecks.Active)
							if upstream.Healthchecks == nil {
								upstream.Healthchecks = &kong.Healthcheck{}
							}
							if upstream.Healthchecks.Active == nil {
								upstream.Healthchecks.Active = &kong.ActiveHealthcheck{}
							}

							mergo.MapWithOverwrite(upstream.Healthchecks.Active, m)
						}

						if kongIngress.Upstream.Healthchecks.Passive != nil {
							m := structs.Map(kongIngress.Upstream.Healthchecks.Passive)

							if upstream.Healthchecks == nil {
								upstream.Healthchecks = &kong.Healthcheck{}
							}
							if upstream.Healthchecks.Passive == nil {
								upstream.Healthchecks.Passive = &kong.PassiveHealthcheck{}
							}

							mergo.MapWithOverwrite(upstream.Healthchecks.Passive, m)
						}
					}
				}

				glog.Infof("creating Kong Upstream with name %v", upstreamName)

				_, err := client.Upstreams.Create(nil, upstream)
				if err != nil {
					glog.Errorf("Unexpected error creating Kong Upstream: %v", err)
					return err
				}
			}

			err = syncTargets(upstreamName, upstream, client)
			if err != nil {
				return errors.Wrap(err, "syncing targets")
			}

			//TODO: check if an update is required (change in kongIngress)
		}
	}

	return nil
}

func (n *NGINXController) syncCertificates(servers []*ingress.Server) error {

	certsToKeep := make(map[string]bool)
	for _, server := range servers {
		if server.Hostname == "_" {
			// there is no catch all server in kong
			continue
		}

		// check the certificate is present in kong
		if server.SSLCert != nil {
			certID, err := n.syncCertificate(server)
			if err != nil {
				return err
			}
			if certID != "" {
				certsToKeep[certID] = true
			}
		}
	}
	client := n.cfg.Kong.Client
	certsInKong, err := client.Certificates.ListAll(nil)
	if err != nil {
		return err
	}

	for _, c := range certsInKong {
		glog.Infof("cert: %v", c.ID)
		if _, ok := certsToKeep[*c.ID]; !ok {
			// delete the cert
			err := client.Certificates.Delete(nil, c.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// syncCertificate synchronizes the state of referenced secrets by Ingress
// rules with Kong certificates.
// This process only create or update certificates in Kong
func (n *NGINXController) syncCertificate(server *ingress.Server) (string, error) {
	if server.SSLCert == nil {
		return "", nil
	}

	if server.SSLCert.ID == "" {
		glog.Warningf("certificate %v/%v is invalid", server.SSLCert.Namespace, server.SSLCert.Name)
		return "", nil
	}

	client := n.cfg.Client

	sc := bytes.NewBuffer(server.SSLCert.Raw.Cert).String()
	sk := bytes.NewBuffer(server.SSLCert.Raw.Key).String()

	//sync cert
	cert, err := client.Certificates.Get(nil, &server.SSLCert.ID)
	if err == nil {
		// check if an update is required

		if *cert.Cert != sc || *cert.Key != sk {
			glog.Infof("updating Kong SSL Certificate for host %v located in Secret %v/%v",
				server.Hostname, server.SSLCert.Namespace, server.SSLCert.Name)

			cert = &kong.Certificate{
				Cert: kong.String(sc),
				Key:  kong.String(sk),
				ID:   kong.String(server.SSLCert.ID),
			}
			cert, err = client.Certificates.Update(nil, cert)
			if err != nil {
				return "", errors.Wrap(err, "patching a Kong certificate")
			}
		}
	} else if kong.IsNotFoundErr(err) {

		cert = &kong.Certificate{
			ID:   kong.String(server.SSLCert.ID),
			Cert: kong.String(sc),
			Key:  kong.String(sk),
		}

		glog.Infof("creating Kong SSL Certificate for host %v located in Secret %v/%v",
			server.Hostname, server.SSLCert.Namespace, server.SSLCert.Name)

		cert, err = client.Certificates.Create(nil, cert)
		if err != nil {
			glog.Errorf("Unexpected error creating Kong Certificate: %v", err)
			return "", err
		}
	} else {
		glog.Errorf("Unexpected response searching a Kong Certificate: %v", err)
		return "", err
	}

	//sync SNI
	sni, err := client.SNIs.Get(nil, &server.Hostname)

	if err == nil {
		// check if it is using the right certificate
		if *sni.Certificate.ID != *cert.ID {
			glog.Infof("updating certificate for host %v to certificate id %v", server.Hostname, cert.ID)

			sni.Certificate.ID = kong.String(*cert.ID)
			sni, err = client.SNIs.Update(nil, sni)
			if err != nil {
				return "", errors.Wrap(err, "patching a Kong SNI")
			}
		}
	} else if kong.IsNotFoundErr(err) {
		sni = &kong.SNI{
			Name:        kong.String(server.Hostname),
			Certificate: &kong.Certificate{ID: kong.String(*cert.ID)},
		}
		glog.Infof("creating Kong SNI for host %v and certificate id %v", server.Hostname, cert.ID)

		_, err = client.SNIs.Create(nil, sni)
		if err != nil {
			glog.Errorf("Unexpected error creating Kong SNI (%v): %v", server.Hostname, err)
		}
	} else {
		glog.Errorf("Unexpected error looking up Kong SNI: %v", err)
		return "", err
	}

	return *cert.ID, nil
}

// buildName returns a string valid as a hostnames taking a backend and
// location as input. The format of backend is <namespace>-<service name>-<port>
// For the location the field Path is used. If the path is / only the backend is used
// This process ensures the returned name is unique.
func buildName(backend string, location *ingress.Location) string {
	return backend
}

// getKongService returns a Route from a list using the path and hosts as filters
func getKongRoute(hostname, path string, routes []*kong.Route) *kong.Route {
	for _, r := range routes {
		if hostname != "_" {
			if sets.NewString(toStringArray(r.Paths)...).Has(path) &&
				sets.NewString(toStringArray(r.Hosts)...).Has(hostname) {
				return r
			}
		} else {
			if sets.NewString(toStringArray(r.Paths)...).Has(path) {
				return r
			}
		}
	}

	return nil
}

// isRouteInKong checks if a route exists or not in Kong
func isRouteInKong(route *kong.Route, routes []*kong.Route) bool {
	for _, eRoute := range routes {
		if compareRoute(route, eRoute) {
			return true
		}
	}
	return false
}

// deleteServiceUpstream deletes multiple Kong upstreams for a particular
// Kong service. This process requires the removal of all the targets that
// reference the upstream to be removed
func deleteServiceUpstream(host string, client *kong.Client) error {
	kongUpstreams, err := client.Upstreams.ListAll(nil)
	if err != nil {
		return err
	}

	upstreamsToRemove := sets.NewString()
	for _, upstream := range kongUpstreams {
		if *upstream.Name == host {
			if !upstreamsToRemove.Has(*upstream.ID) {
				upstreamsToRemove.Insert(*upstream.ID)
			}
		}
	}

	for _, upstream := range upstreamsToRemove.List() {
		kongTargets, err := client.Targets.ListAll(nil, &upstream)
		if err != nil {
			return err
		}

		for _, target := range kongTargets {
			if *target.UpstreamID == upstream {
				err := client.Targets.Delete(nil, target.UpstreamID, target.ID)
				if err != nil {
					return errors.Wrap(err, "removing a Kong target")
				}
			}
		}

		err = client.Upstreams.Delete(nil, &upstream)
		if err != nil {
			return errors.Wrap(err, "removing a Kong upstream")
		}
	}

	return nil
}

// pluginDeepEqual compares the configuration of a Plugin (CRD) against
// the persisted state in the Kong database
// This is required because a plugin has defaults that could not exists in the CRD.
func pluginDeepEqual(config map[string]interface{}, kong *kong.Plugin) bool {
	return pluginDeepEqualWrapper(config, kong.Config)
}

func pluginDeepEqualWrapper(config1 map[string]interface{}, config2 map[string]interface{}) bool {
	return interfaceDeepEqual(config1, config2)
}

func interfaceDeepEqual(i1 interface{}, i2 interface{}) bool {
	v1 := reflect.ValueOf(i1)
	v2 := reflect.ValueOf(i2)

	k1 := v1.Type().Kind()
	k2 := v2.Type().Kind()
	if k1 == k2 {
		if k1 == reflect.Map {
			return mapDeepEqual(v1, v2)
		} else if k1 == reflect.Slice || k1 == reflect.Array {
			return listUnorderedDeepEqual(v1, v2)
		}
	}
	j1, e1 := json.Marshal(v1.Interface())
	j2, e2 := json.Marshal(v2.Interface())
	return e1 == nil && e2 == nil && string(j1) == string(j2)
}

func mapDeepEqual(m1 reflect.Value, m2 reflect.Value) bool {
	keys1 := m1.MapKeys()
	for _, k := range keys1 {
		v2 := m2.MapIndex(k)
		if !v2.IsValid() { // k not found in m2
			return false
		}
		v1 := m1.MapIndex(k)

		if v1.IsValid() && !interfaceDeepEqual(v1.Interface(), v2.Interface()) {
			return false
		}
	}
	return true
}

func listUnorderedDeepEqual(l1 reflect.Value, l2 reflect.Value) bool {
	length := l1.Len()
	if length != l2.Len() {
		return false
	}
	for i := 0; i < length; i++ {
		v1 := l1.Index(i)
		if !v1.IsValid() {
			return false // this shouldn't happen
		}
		found := false
		for j := 0; j < length; j++ {
			v2 := l2.Index(j)
			if v2.IsValid() && interfaceDeepEqual(v1.Interface(), v2.Interface()) {
				found = true
				break
			}
		}
		if !found {
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
