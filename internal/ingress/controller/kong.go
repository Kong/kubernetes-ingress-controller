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

	"github.com/blang/semver"
	"github.com/golang/glog"
	"github.com/hbagdi/deck/counter"
	"github.com/hbagdi/deck/diff"
	"github.com/hbagdi/deck/dump"
	"github.com/hbagdi/deck/file"
	"github.com/hbagdi/deck/solver"
	"github.com/hbagdi/deck/state"
	"github.com/hbagdi/deck/utils"
	"github.com/hbagdi/go-kong/kong"
	"github.com/imdario/mergo"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/controller/kong/dbless"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/controller/parser"
	"github.com/pkg/errors"
)

var count counter.Counter

var upstreamDefaults = kong.Upstream{
	Slots: kong.Int(10000),
	Healthchecks: &kong.Healthcheck{
		Active: &kong.ActiveHealthcheck{
			Concurrency: kong.Int(10),
			Healthy: &kong.Healthy{
				HTTPStatuses: []int{200, 302},
				Interval:     kong.Int(0),
				Successes:    kong.Int(0),
			},
			HTTPPath:               kong.String("/"),
			Timeout:                kong.Int(1),
			HTTPSVerifyCertificate: kong.Bool(true),
			Type:                   kong.String("http"),
			Unhealthy: &kong.Unhealthy{
				HTTPFailures: kong.Int(0),
				TCPFailures:  kong.Int(0),
				Timeouts:     kong.Int(0),
				HTTPStatuses: []int{429, 404, 500, 501, 502, 503, 504, 505},
			},
		},
		Passive: &kong.PassiveHealthcheck{
			Healthy: &kong.Healthy{
				HTTPStatuses: []int{200, 201, 202, 203, 204, 205,
					206, 207, 208, 226, 300, 301, 302, 303, 304, 305,
					306, 307, 308},
				Successes: kong.Int(0),
			},
			Unhealthy: &kong.Unhealthy{
				HTTPFailures: kong.Int(0),
				TCPFailures:  kong.Int(0),
				Timeouts:     kong.Int(0),
				HTTPStatuses: []int{429, 500, 503},
			},
		},
	},
	HashOn:           kong.String("none"),
	HashFallback:     kong.String("none"),
	HashOnCookiePath: kong.String("/"),
}

// OnUpdate is called periodically by syncQueue to keep the configuration in sync.
// returning nil implies the synchronization finished correctly.
// Returning an error means requeue the update.
func (n *KongController) OnUpdate(state *parser.KongState) error {
	if n.cfg.InMemory {
		return n.onUpdateInMemoryMode(state)
	}
	return n.onUpdateDBMode(state)
}

func (n *KongController) onUpdateInMemoryMode(state *parser.KongState) error {
	client := n.cfg.Kong.Client

	config := dbless.KongNativeState(state)
	jsonConfig, err := json.Marshal(&config)
	if err != nil {
		return errors.Wrap(err,
			"marshaling Kong declarative configuration to JSON")
	}

	if reflect.DeepEqual(n.runningConfigHash, sha256.Sum256(jsonConfig)) {
		glog.V(2).Info("no configuration change, skipping call sync config")
		return nil
	}
	type reqBody struct {
		Config string `json:"config"`
	}
	json, err := json.Marshal(&reqBody{Config: string(jsonConfig)})
	if err != nil {
		//TODO annotate
		return errors.Wrap(err,
			"marshaling /config request body")
	}
	req, err := http.NewRequest("POST", n.cfg.Kong.URL+"/config",
		bytes.NewReader(json))
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
	n.runningConfigHash = sha256.Sum256(jsonConfig)

	return err
}

func (n *KongController) onUpdateDBMode(parserState *parser.KongState) error {
	client := n.cfg.Kong.Client

	// parse to target file
	targetContent, err := n.toDeckContent(parserState)
	if err != nil {
		return err
	}

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
	rawState, err = file.Get(targetContent, currentState)
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
	// TODO configurable parallelism
	//client.SetDebugMode(true)
	errs := solver.Solve(nil, syncer, client, 10, false)
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
			}
			service.Routes = append(service.Routes, &route)
		}
		content.Services = append(content.Services, service)
	}

	for _, plugin := range k8sState.GlobalPlugins {
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

	for _, u := range k8sState.Upstreams {
		n.fillUpstream(&u.Upstream)
		upstream := file.FUpstream{Upstream: u.Upstream}
		for _, t := range u.Targets {
			target := file.FTarget{Target: t.Target}
			upstream.Targets = append(upstream.Targets, &target)
		}
		content.Upstreams = append(content.Upstreams, upstream)
	}

	for _, c := range k8sState.Certificates {
		cert := file.FCertificate{Certificate: c.Certificate}
		content.Certificates = append(content.Certificates, cert)
	}

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
	selectorTags := n.getIngressControllerTags()
	if len(selectorTags) > 0 {
		content.Info = &file.Info{
			SelectorTags: selectorTags,
		}
	}

	return &content, nil
}

var (
	kong110version = semver.MustParse("1.1.0")
	kong120version = semver.MustParse("1.2.0")
	kong130version = semver.MustParse("1.3.0")

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
	return nil
}

type intTransformer struct {
}

func (t intTransformer) Transformer(typ reflect.Type) func(dst,
	src reflect.Value) error {
	var a *int
	var ar []int
	if typ == reflect.TypeOf(ar) {
		return func(dst, src reflect.Value) error {
			if dst.CanSet() {
				if reflect.DeepEqual(reflect.Zero(dst.Type()).Interface(),
					dst.Interface()) {
					return nil
				}
			}
			return nil
		}
	}
	if typ == reflect.TypeOf(a) {
		return func(dst, src reflect.Value) error {
			if dst.CanSet() {
				if reflect.DeepEqual(reflect.Zero(dst.Type()).Interface(),
					dst.Interface()) {
					return nil
				}
			}
			return nil
		}
	}
	return nil
}

func setDefaultsInUpstream(upstream *kong.Upstream) error {
	err := mergo.Merge(upstream, upstreamDefaults,
		mergo.WithTransformers(intTransformer{}))
	if err != nil {
		return errors.Wrap(err, "error overriding upstream")
	}
	return err
}
