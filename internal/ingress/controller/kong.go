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
	"net/http"
	"reflect"
	"strconv"

	"github.com/blang/semver"
	"github.com/golang/glog"
	"github.com/hbagdi/deck/counter"
	"github.com/hbagdi/deck/diff"
	"github.com/hbagdi/deck/dump"
	"github.com/hbagdi/deck/solver"
	"github.com/hbagdi/deck/state"
	"github.com/hbagdi/deck/utils"
	"github.com/hbagdi/go-kong/kong"
	"github.com/hbagdi/go-kong/kong/custom"
	"github.com/imdario/mergo"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/controller/kong/dbless"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/controller/parser"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/sets"
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

	// XXX
	err := n.fillConsumersAndCredentials(state)
	if err != nil {
		return err
	}

	config := dbless.KongNativeState(state)
	jsonConfig, err := json.Marshal(&config)
	if err != nil {
		return errors.Wrap(err,
			"marshaling Kong declarative configuration to JSON")
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
	_, err = client.Do(nil, req, nil)
	if err != nil {
		return errors.Wrap(err, "posting new config to /config")
	}

	return err
}

func (n *KongController) onUpdateDBMode(state *parser.KongState) error {

	client := n.cfg.Kong.Client

	err := n.syncConsumers()
	if err != nil {
		return err
	}

	err = n.syncCredentials()
	if err != nil {
		return err
	}

	currentState, err := dump.GetState(client, dump.Config{})
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
	return nil
}

func (n *KongController) toDeckKongState(k8sState *parser.KongState) (*state.KongState, error) {
	targetState, err := state.NewKongState()
	if err != nil {
		return nil, errors.Wrap(err, "creating new parser.KongState")
	}

	for _, s := range k8sState.Services {
		service := state.Service{Service: s.Service}
		service.ID = kong.String("placeholder-" +
			strconv.FormatUint(count.Inc(), 10))
		err := targetState.Services.Add(service)
		if err != nil {
			return nil, errors.Wrap(err, "inserting service into state")
		}
		for _, p := range s.Plugins {
			plugin := state.Plugin{
				Plugin: *p.DeepCopy(),
			}
			err = n.fillPlugin(&plugin)
			if err != nil {
				glog.Errorf("error filling in defaults for plugin: %s", *plugin.Name)
			}
			plugin.ID = kong.String("placeholder-" +
				strconv.FormatUint(count.Inc(), 10))
			plugin.Service = s.DeepCopy()
			err = targetState.Plugins.Add(plugin)
			if err != nil {
				return nil, errors.Wrap(err, "inserting plugin into state")
			}
		}

		for _, r := range s.Routes {
			route := state.Route{Route: r.Route}
			route.ID = kong.String("placeholder-" +
				strconv.FormatUint(count.Inc(), 10))
			route.Service = service.DeepCopy()
			err := targetState.Routes.Add(route)
			if err != nil {
				return nil, errors.Wrap(err, "inserting route into state")
			}
			for _, p := range r.Plugins {
				plugin := state.Plugin{
					Plugin: *p.DeepCopy(),
				}
				err = n.fillPlugin(&plugin)
				if err != nil {
					glog.Errorf("error filling in defaults for plugin: %s", *plugin.Name)
				}
				plugin.ID = kong.String("placeholder-" +
					strconv.FormatUint(count.Inc(), 10))
				plugin.Route = r.DeepCopy()
				err := targetState.Plugins.Add(plugin)
				if err != nil {
					return nil, errors.Wrap(err, "inserting plugin into state")
				}
			}
		}
	}

	for _, plugin := range k8sState.GlobalPlugins {
		plugin := state.Plugin{
			Plugin: plugin.Plugin,
		}
		err = n.fillPlugin(&plugin)
		if err != nil {
			glog.Errorf("error filling in defaults for plugin: %s", *plugin.Name)
		}

		plugin.ID = kong.String("placeholder-" +
			strconv.FormatUint(count.Inc(), 10))
		err := targetState.Plugins.Add(plugin)
		if err != nil {
			return nil, errors.Wrap(err, "inserting plugin into state")
		}
	}

	for _, u := range k8sState.Upstreams {
		upstream := state.Upstream{Upstream: u.Upstream}
		upstream.ID = kong.String("placeholder-" +
			strconv.FormatUint(count.Inc(), 10))
		err = setDefaultsInUpstream(&upstream.Upstream)
		if err != nil {
			return nil, err
		}
		err := targetState.Upstreams.Add(upstream)
		if err != nil {
			return nil, errors.Wrap(err, "inserting upstream into state")
		}
		for _, t := range u.Targets {
			target := state.Target{Target: t.Target}
			target.ID = kong.String("placeholder-" +
				strconv.FormatUint(count.Inc(), 10))
			target.Upstream = upstream.DeepCopy()
			target.Weight = kong.Int(100)
			err := targetState.Targets.Add(target)
			if err != nil {
				return nil, errors.Wrap(err, "inserting target into state")
			}
		}
	}

	for _, c := range k8sState.Certificates {
		cert := state.Certificate{Certificate: c.Certificate}
		cert.ID = kong.String("placeholder-" +
			strconv.FormatUint(count.Inc(), 10))
		err := targetState.Certificates.Add(cert)
		if err != nil {
			return nil, errors.Wrap(err, "inserting certificate into state")
		}
	}
	return targetState, nil
}

var kong110version = semver.MustParse("1.1.0")

func (n *KongController) fillPlugin(plugin *state.Plugin) error {
	if plugin == nil {
		return errors.New("plugin is nil")
	}
	if plugin.Name == nil || *plugin.Name == "" {
		return errors.New("plugin doesn't have a name")
	}
	schema, err := n.PluginSchemaStore.Schema(*plugin.Name)
	if err != nil {
		return errors.Wrapf(err, "error retrieveing schema for plugin %s", *plugin.Name)
	}
	if plugin.Config == nil {
		plugin.Config = make(kong.Configuration)
	}
	newConfig, err := fill(schema, plugin.Config)
	if err != nil {
		return errors.Wrapf(err, "error filling in default for plugin %s", *plugin.Name)
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

func (n *KongController) fillConsumersAndCredentials(state *parser.KongState) error {
	consumers := make(map[string]parser.Consumer)
	for _, consumer := range n.store.ListKongConsumers() {
		if consumer.Username == "" && consumer.CustomID == "" {
			continue
		}

		c := parser.Consumer{}

		if consumer.Username != "" {
			c.Username = kong.String(consumer.Username)
		}

		if consumer.CustomID != "" {
			c.CustomID = kong.String(consumer.CustomID)
		}

		consumers[consumer.Namespace+"/"+consumer.Name] = c
	}

	for _, credential := range n.store.ListKongCredentials() {
		consumer, ok := consumers[credential.Namespace+"/"+credential.ConsumerRef]
		if !ok {
			continue
		}
		if consumer.Credentials == nil {
			consumer.Credentials = make(map[string][]map[string]interface{})
		}
		if credential.Type == "" {
			continue
		}
		consumer.Credentials[credential.Type] = append(consumer.Credentials[credential.Type], credential.Config)

		consumers[credential.Namespace+"/"+credential.ConsumerRef] = consumer
	}

	for _, c := range consumers {
		state.Consumers = append(state.Consumers, c)
	}
	return nil
}

// syncConsumers synchronizes the state between KongConsumer (Kubernetes CRD) type and Kong consumers.
// This loop only creates new consumers in Kong.
func (n *KongController) syncConsumers() error {

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
		glog.V(2).Infof("checking if Kong consumer %v exists", consumer.Name)
		consumerID := fmt.Sprintf("%v", consumer.GetUID())

		kc, ok := consumersInKong[consumerID]

		if !ok {
			glog.V(2).Infof("Creating Kong consumer %v", consumerID)
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
func (n *KongController) syncCredentials() error {
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

type intTransformer struct {
}

func (t intTransformer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	var a *int
	var ar []int
	if typ == reflect.TypeOf(ar) {
		return func(dst, src reflect.Value) error {
			if dst.CanSet() {
				if reflect.DeepEqual(reflect.Zero(dst.Type()).Interface(), dst.Interface()) {
					return nil
				}
			}
			return nil
		}
	}
	if typ == reflect.TypeOf(a) {
		return func(dst, src reflect.Value) error {
			if dst.CanSet() {
				if reflect.DeepEqual(reflect.Zero(dst.Type()).Interface(), dst.Interface()) {
					return nil
				}
			}
			return nil
		}
	}
	return nil
}

func setDefaultsInUpstream(upstream *kong.Upstream) error {
	err := mergo.Merge(upstream, upstreamDefaults, mergo.WithTransformers(intTransformer{}))
	if err != nil {
		return errors.Wrap(err, "error overriding upstream")
	}
	return err
}
