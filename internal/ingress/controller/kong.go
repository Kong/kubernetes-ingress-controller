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
	"fmt"
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
	"github.com/hbagdi/go-kong/kong/custom"
	"github.com/imdario/mergo"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/controller/kong/dbless"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/controller/parser"
	"github.com/pkg/errors"
)

var count counter.Counter

const ingressControllerTag = "managed-by-ingress-controller"

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

func (n *KongController) onUpdateDBMode(state *parser.KongState) error {

	client := n.cfg.Kong.Client

	currentState, err := dump.GetState(client,
		dump.Config{
			SelectorTags: n.getIngressControllerTags(),
		})
	if err != nil {
		return err
	}
	targetState, err := n.toDeckKongState(state)
	if err != nil {
		return err
	}
	syncer, err := diff.NewSyncer(currentState, targetState)
	if err != nil {
		return errors.Wrap(err, "creating a new syncer")
	}
	errs := solver.Solve(nil, syncer, client, false)
	if errs != nil {
		return utils.ErrArray{Errors: errs}
	}

	// credentials are not synced using decK
	err = n.syncCredentials(state)
	return err
}

// getIngressControllerTags returns a tag to use if the current
// Kong entity supports tagging.
func (n *KongController) getIngressControllerTags() []string {
	var res []string
	if n.cfg.Kong.HasTagSupport {
		res = append(res, ingressControllerTag)
	}
	return res
}

func (n *KongController) toDeckKongState(
	k8sState *parser.KongState) (*state.KongState, error) {
	var content file.Content
	var err error

	for _, s := range k8sState.Services {
		service := file.Service{Service: s.Service}
		for _, p := range s.Plugins {
			plugin := file.Plugin{
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
			route := file.Route{Route: r.Route}
			n.fillRoute(&route.Route)

			for _, p := range r.Plugins {
				plugin := file.Plugin{
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
		plugin := file.Plugin{
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
		upstream := file.Upstream{Upstream: u.Upstream}
		for _, t := range u.Targets {
			target := file.Target{Target: t.Target}
			upstream.Targets = append(upstream.Targets, &target)
		}
		content.Upstreams = append(content.Upstreams, upstream)
	}

	for _, c := range k8sState.Certificates {
		cert := file.Certificate{Certificate: c.Certificate}
		content.Certificates = append(content.Certificates, cert)
	}

	for _, c := range k8sState.Consumers {
		consumer := file.Consumer{Consumer: c.Consumer}
		for _, p := range c.Plugins {
			consumer.Plugins = append(consumer.Plugins, &file.Plugin{Plugin: p})
		}
		content.Consumers = append(content.Consumers, consumer)
	}

	content.Info.SelectorTags = n.getIngressControllerTags()
	targetState, _, err := file.GetStateFromContent(&content)
	if err != nil {
		return nil, errors.Wrap(err, "error creating a valid state for Kong")
	}
	return targetState, nil
}

var kong110version = semver.MustParse("1.1.0")

var kong120version = semver.MustParse("1.2.0")

func (n *KongController) fillRoute(route *kong.Route) {
	if n.cfg.Kong.Version.GTE(kong120version) {
		if route.HTTPSRedirectStatusCode == nil {
			route.HTTPSRedirectStatusCode = kong.Int(426)
		}
	}
}

func (n *KongController) fillPlugin(plugin *file.Plugin) error {
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

// syncCredentials synchronizes the state between KongCredential
// (Kubernetes CRD) and Kong credentials.
func (n *KongController) syncCredentials(state *parser.KongState) error {
	client := n.cfg.Kong.Client

	// List all consumers in Kong to obtain an ID
	usernameToConsumer := make(map[string]*kong.Consumer)
	customIDToConsumer := make(map[string]*kong.Consumer)
	kongConsumers, err := client.Consumers.ListAll(nil)
	if err != nil {
		return err
	}

	// create simple indexes
	for i := range kongConsumers {
		username := kongConsumers[i].Username
		customID := kongConsumers[i].CustomID

		if !utils.Empty(username) {
			usernameToConsumer[*username] = kongConsumers[i]
		}
		if !utils.Empty(customID) {
			customIDToConsumer[*customID] = kongConsumers[i]
		}
	}

	for _, consumer := range state.Consumers {
		for credType, creds := range consumer.Credentials {
			for _, credData := range creds {
				// lookup consumerID
				var consumerID string
				if !utils.Empty(consumer.Username) {
					if c, ok := usernameToConsumer[*consumer.Username]; ok {
						consumerID = *c.ID
					}
				}
				if !utils.Empty(consumer.CustomID) {
					if c, ok := customIDToConsumer[*consumer.CustomID]; ok {
						consumerID = *c.ID
					}
				}
				if consumerID == "" {
					continue
				}

				// lookup credential in Kong
				credInKong := custom.NewEntityObject(custom.Type(credType))
				credentialID := fmt.Sprintf("%v", credData["id"])
				credInKong.AddRelation("consumer_id", consumerID)
				credInKong.SetObject(map[string]interface{}{"id": credentialID})

				_, err = client.CustomEntities.Get(nil, credInKong)
				if !kong.IsNotFoundErr(err) || err == nil {
					return err
				}

				// if not found, then create it
				credInKong.SetObject(map[string]interface{}(credData))
				_, err := client.CustomEntities.Create(nil, credInKong)
				if err != nil {
					glog.Errorf("Unexpected error creating credential: %v", err)
					return err
				}
				// TODO: allow changes in credentials?
				// TODO sync this object using PUT (whenever the ResourceVersion is out of date)
				// An alternate solution would be to handle UPDATE events separately
			}
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
