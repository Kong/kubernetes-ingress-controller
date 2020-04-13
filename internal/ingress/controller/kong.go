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
	"crypto/sha256"
	"encoding/json"
	"net/http"
	"reflect"
	"sort"
	"strings"

	"github.com/blang/semver"
	"github.com/golang/glog"
	"github.com/hbagdi/deck/diff"
	"github.com/hbagdi/deck/dump"
	"github.com/hbagdi/deck/file"
	"github.com/hbagdi/deck/solver"
	"github.com/hbagdi/deck/state"
	"github.com/hbagdi/deck/utils"
	"github.com/hbagdi/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/controller/parser"
	"github.com/pkg/errors"
)

// OnUpdate is called periodically by syncQueue to keep the configuration in sync.
// returning nil implies the synchronization finished correctly.
// Returning an error means requeue the update.
func (n *KongController) OnUpdate(state *parser.KongState) error {
	targetContent, err := n.toDeckContent(state)
	if err != nil {
		return err
	}

	jsonConfig, err := json.Marshal(targetContent)
	if err != nil {
		return errors.Wrap(err,
			"marshaling Kong declarative configuration to JSON")
	}
	shaSum := sha256.Sum256(jsonConfig)
	if reflect.DeepEqual(n.runningConfigHash, shaSum) {
		glog.Info("no configuration change, skipping sync to Kong")
		return nil
	}
	if n.cfg.InMemory {
		err = n.onUpdateInMemoryMode(targetContent)
	} else {
		err = n.onUpdateDBMode(targetContent)
	}
	if err == nil {
		glog.Info("successfully synced configuration to Kong")
		n.runningConfigHash = shaSum
	}
	return err
}

func cleanUpNullsInPluginConfigs(state *file.Content) {

	for _, s := range state.Services {
		for _, p := range s.Plugins {
			for k, v := range p.Config {
				if v == nil {
					delete(p.Config, k)
				}
			}
		}
		for _, r := range state.Routes {
			for _, p := range r.Plugins {
				for k, v := range p.Config {
					if v == nil {
						delete(p.Config, k)
					}
				}
			}
		}
	}

	for _, c := range state.Consumers {
		for _, p := range c.Plugins {
			for k, v := range p.Config {
				if v == nil {
					delete(p.Config, k)
				}
			}
		}
	}

	for _, p := range state.Plugins {
		for k, v := range p.Config {
			if v == nil {
				delete(p.Config, k)
			}
		}
	}
}

func (n *KongController) onUpdateInMemoryMode(state *file.Content) error {
	client := n.cfg.Kong.Client

	// Kong will error out if this is set
	state.Info = nil
	// Kong errors out if `null`s are present in `config` of plugins
	cleanUpNullsInPluginConfigs(state)

	config, err := json.Marshal(state)
	if err != nil {
		return errors.Wrap(err,
			"marshaling Kong declarative configuration to JSON")
	}
	req, err := http.NewRequest("POST", n.cfg.Kong.URL+"/config",
		bytes.NewReader(config))
	if err != nil {
		return errors.Wrap(err, "creating new HTTP request for /config")
	}
	req.Header.Add("content-type", "application/json")

	queryString := req.URL.Query()
	queryString.Add("check_hash", "1")

	req.URL.RawQuery = queryString.Encode()

	_, err = client.Do(nil, req, nil)
	if err != nil {
		return errors.Wrap(err, "posting new config to /config")
	}

	return err
}

func (n *KongController) onUpdateDBMode(targetContent *file.Content) error {
	client := n.cfg.Kong.Client

	// read the current state
	rawState, err := dump.Get(client, dump.Config{
		SelectorTags: n.getIngressControllerTags(),
	})
	if err != nil {
		return errors.Wrap(err, "loading configuration from kong")
	}
	currentState, err := state.Get(rawState)
	if err != nil {
		return err
	}

	// read the target state
	rawState, err = file.Get(targetContent, file.RenderConfig{
		CurrentState: currentState,
		KongVersion:  n.cfg.Kong.Version,
	})
	if err != nil {
		return err
	}
	targetState, err := state.Get(rawState)
	if err != nil {
		return err
	}

	syncer, err := diff.NewSyncer(currentState, targetState)
	if err != nil {
		return errors.Wrap(err, "creating a new syncer")
	}
	syncer.SilenceWarnings = true
	//client.SetDebugMode(true)
	_, errs := solver.Solve(nil, syncer, client, n.cfg.Kong.Concurrency, false)
	if errs != nil {
		return utils.ErrArray{Errors: errs}
	}
	return nil
}

// getIngressControllerTags returns a tag to use if the current
// Kong entity supports tagging.
func (n *KongController) getIngressControllerTags() []string {
	var res []string
	if n.cfg.Kong.HasTagSupport {
		res = append(res, n.cfg.Kong.FilterTags...)
	}
	return res
}

func (n *KongController) toDeckContent(
	k8sState *parser.KongState) (*file.Content, error) {
	var content file.Content
	content.FormatVersion = "1.1"
	var err error

	for _, s := range k8sState.Services {
		service := file.FService{Service: s.Service}
		for _, p := range s.Plugins {
			plugin := file.FPlugin{
				Plugin: *p.DeepCopy(),
			}
			err = n.fillPlugin(&plugin)
			if err != nil {
				glog.Errorf("error filling in defaults for plugin: %s",
					*plugin.Name)
			}
			service.Plugins = append(service.Plugins, &plugin)
			sort.SliceStable(service.Plugins, func(i, j int) bool {
				return strings.Compare(*service.Plugins[i].Name, *service.Plugins[j].Name) > 0
			})
		}

		for _, r := range s.Routes {
			route := file.FRoute{Route: r.Route}
			n.fillRoute(&route.Route)

			for _, p := range r.Plugins {
				plugin := file.FPlugin{
					Plugin: *p.DeepCopy(),
				}
				err = n.fillPlugin(&plugin)
				if err != nil {
					glog.Errorf("error filling in defaults for plugin: %s",
						*plugin.Name)
				}
				route.Plugins = append(route.Plugins, &plugin)
				sort.SliceStable(route.Plugins, func(i, j int) bool {
					return strings.Compare(*route.Plugins[i].Name, *route.Plugins[j].Name) > 0
				})
			}
			service.Routes = append(service.Routes, &route)
		}
		sort.SliceStable(service.Routes, func(i, j int) bool {
			return strings.Compare(*service.Routes[i].Name, *service.Routes[j].Name) > 0
		})
		content.Services = append(content.Services, service)
	}
	sort.SliceStable(content.Services, func(i, j int) bool {
		return strings.Compare(*content.Services[i].Name, *content.Services[j].Name) > 0
	})

	for _, plugin := range k8sState.Plugins {
		plugin := file.FPlugin{
			Plugin: plugin.Plugin,
		}
		err = n.fillPlugin(&plugin)
		if err != nil {
			glog.Errorf("error filling in defaults for plugin: %s",
				*plugin.Name)
		}
		content.Plugins = append(content.Plugins, plugin)
	}
	sort.SliceStable(content.Plugins, func(i, j int) bool {
		return strings.Compare(pluginString(content.Plugins[i]),
			pluginString(content.Plugins[j])) > 0
	})

	for _, u := range k8sState.Upstreams {
		n.fillUpstream(&u.Upstream)
		upstream := file.FUpstream{Upstream: u.Upstream}
		for _, t := range u.Targets {
			target := file.FTarget{Target: t.Target}
			upstream.Targets = append(upstream.Targets, &target)
		}
		sort.SliceStable(upstream.Targets, func(i, j int) bool {
			return strings.Compare(*upstream.Targets[i].Target.Target, *upstream.Targets[j].Target.Target) > 0
		})
		content.Upstreams = append(content.Upstreams, upstream)
	}
	sort.SliceStable(content.Upstreams, func(i, j int) bool {
		return strings.Compare(*content.Upstreams[i].Name, *content.Upstreams[j].Name) > 0
	})

	for _, c := range k8sState.Certificates {
		cert := getFCertificateFromKongCert(c.Certificate)
		content.Certificates = append(content.Certificates, cert)
	}
	sort.SliceStable(content.Certificates, func(i, j int) bool {
		return strings.Compare(*content.Certificates[i].Cert, *content.Certificates[j].Cert) > 0
	})

	for _, c := range k8sState.Consumers {
		consumer := file.FConsumer{Consumer: c.Consumer}
		for _, p := range c.Plugins {
			consumer.Plugins = append(consumer.Plugins, &file.FPlugin{Plugin: p})
		}
		consumer.KeyAuths = c.KeyAuths
		consumer.HMACAuths = c.HMACAuths
		consumer.BasicAuths = c.BasicAuths
		consumer.JWTAuths = c.JWTAuths
		consumer.ACLGroups = c.ACLGroups
		consumer.Oauth2Creds = c.Oauth2Creds
		content.Consumers = append(content.Consumers, consumer)
	}
	sort.SliceStable(content.Consumers, func(i, j int) bool {
		return strings.Compare(*content.Consumers[i].Username, *content.Consumers[j].Username) > 0
	})
	selectorTags := n.getIngressControllerTags()
	if len(selectorTags) > 0 {
		content.Info = &file.Info{
			SelectorTags: selectorTags,
		}
	}

	return &content, nil
}
func getFCertificateFromKongCert(kongCert kong.Certificate) file.FCertificate {
	var res file.FCertificate
	if kongCert.ID != nil {
		res.ID = kong.String(*kongCert.ID)
	}
	if kongCert.Key != nil {
		res.Key = kong.String(*kongCert.Key)
	}
	if kongCert.Cert != nil {
		res.Cert = kong.String(*kongCert.Cert)
	}
	res.SNIs = getSNIs(kongCert.SNIs)
	return res
}

func getSNIs(names []*string) []kong.SNI {
	var snis []kong.SNI
	for _, name := range names {
		snis = append(snis, kong.SNI{
			Name: kong.String(*name),
		})
	}
	return snis
}

func pluginString(plugin file.FPlugin) string {
	result := ""
	if plugin.Name != nil {
		result = *plugin.Name
	}
	if plugin.Consumer != nil && plugin.Consumer.ID != nil {
		result += *plugin.Consumer.ID
	}
	if plugin.Route != nil && plugin.Route.ID != nil {
		result += *plugin.Route.ID
	}
	if plugin.Service != nil && plugin.Service.ID != nil {
		result += *plugin.Service.ID
	}
	return result
}

var (
	kong110version = semver.MustParse("1.1.0")
	kong120version = semver.MustParse("1.2.0")
	kong130version = semver.MustParse("1.3.0")
	kong200version = semver.MustParse("2.0.0")

	kongEnterprise036version = semver.MustParse("0.36.0")
)

func (n *KongController) fillRoute(route *kong.Route) {
	if n.cfg.Kong.Version.GTE(kong120version) ||
		(n.cfg.Kong.Enterprise &&
			n.cfg.Kong.Version.GTE(kongEnterprise036version)) {
		if route.HTTPSRedirectStatusCode == nil {
			route.HTTPSRedirectStatusCode = kong.Int(426)
		}
	}
	if n.cfg.Kong.Version.GTE(kong200version) {
		if route.PathHandling == nil {
			route.PathHandling = kong.String("v0")
		}
	}
}

func (n *KongController) fillUpstream(upstream *kong.Upstream) {
	if n.cfg.Kong.Version.GTE(kong130version) {
		if upstream.Algorithm == nil {
			upstream.Algorithm = kong.String("round-robin")
		}
	}
}

func (n *KongController) fillPlugin(plugin *file.FPlugin) error {
	if plugin == nil {
		return errors.New("plugin is nil")
	}
	if plugin.Name == nil || *plugin.Name == "" {
		return errors.New("plugin doesn't have a name")
	}
	schema, err := n.PluginSchemaStore.Schema(*plugin.Name)
	if err != nil {
		return errors.Wrapf(err, "error retrieveing schema for plugin %s",
			*plugin.Name)
	}
	if plugin.Config == nil {
		plugin.Config = make(kong.Configuration)
	}
	newConfig, err := fill(schema, plugin.Config)
	if err != nil {
		return errors.Wrapf(err, "error filling in default for plugin %s",
			*plugin.Name)
	}
	plugin.Config = newConfig
	if plugin.RunOn == nil {
		plugin.RunOn = kong.String("first")
	}
	if plugin.Enabled == nil {
		plugin.Enabled = kong.Bool(true)
	}
	if n.cfg.Kong.Version.GTE(kong110version) {
		if len(plugin.Protocols) == 0 {
			// TODO read this from the schema endpoint
			plugin.Protocols = kong.StringSlice("http", "https")
		}
	}
	if n.cfg.Kong.Version.GTE(kong200version) {
		plugin.RunOn = nil
	}
	return nil
}
