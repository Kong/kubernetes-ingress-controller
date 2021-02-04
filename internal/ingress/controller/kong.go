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
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"sort"
	"strings"

	"github.com/kong/deck/diff"
	"github.com/kong/deck/dump"
	"github.com/kong/deck/file"
	"github.com/kong/deck/solver"
	"github.com/kong/deck/state"
	deckutils "github.com/kong/deck/utils"
	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/pkg/deckgen"
	"github.com/kong/kubernetes-ingress-controller/pkg/kongstate"
	"github.com/kong/kubernetes-ingress-controller/pkg/util"
	"github.com/sirupsen/logrus"
)

// OnUpdate is called periodically by syncQueue to keep the configuration in sync.
// returning nil implies the synchronization finished correctly.
// Returning an error means requeue the update.
func (n *KongController) OnUpdate(ctx context.Context, state *kongstate.KongState) error {
	var customEntities []byte
	var err error
	// process any custom entities
	if n.cfg.InMemory && n.cfg.KongCustomEntitiesSecret != "" {
		customEntities, err = n.fetchCustomEntities()
		if err != nil {
			// failure to fetch custom entities shouldn't block updates
			n.Logger.Errorf("failed to fetch custom entities: %v", err)
		}
	}

	newSHA, err := performUpdate(ctx,
		n.Logger,
		&n.cfg.Kong,
		n.cfg.InMemory,
		n.cfg.EnableReverseSync,
		state,
		&n.PluginSchemaStore,
		n.getIngressControllerTags(),
		customEntities,
		n.runningConfigHash,
	)

	n.runningConfigHash = newSHA
	return err
}

func equalSHA(a, b []byte) bool {
	return reflect.DeepEqual(a, b)
}

func performUpdate(ctx context.Context,
	log logrus.FieldLogger,
	kongConfig *Kong,
	inMemory bool,
	reverseSync bool,
	state *kongstate.KongState,
	schemas *util.PluginSchemaStore,
	selectorTags []string,
	customEntities []byte,
	oldSHA []byte,
) ([]byte, error) {
	targetContent := toDeckContent(ctx, log, state, schemas, selectorTags)

	newSHA, err := deckgen.GenerateSHA(targetContent, customEntities)
	if err != nil {
		return oldSHA, err
	}
	// disable optimization if reverse sync is enabled
	if !reverseSync {
		if equalSHA(oldSHA, newSHA) {
			log.Info("no configuration change, skipping sync to kong")
			return oldSHA, nil
		}
	}

	if inMemory {
		err = onUpdateInMemoryMode(ctx, log, targetContent, customEntities, kongConfig)
	} else {
		err = onUpdateDBMode(targetContent, kongConfig, selectorTags)
	}
	if err != nil {
		return nil, err
	}
	log.Info("successfully synced configuration to kong")
	return newSHA, nil
}

func renderConfigWithCustomEntities(log logrus.FieldLogger, state *file.Content,
	customEntitiesJSONBytes []byte) ([]byte, error) {

	var kongCoreConfig []byte
	var err error

	kongCoreConfig, err = json.Marshal(state)
	if err != nil {
		return nil, fmt.Errorf("marshaling kong config into json: %w", err)
	}

	// fast path
	if len(customEntitiesJSONBytes) == 0 {
		return kongCoreConfig, nil
	}

	// slow path
	mergeMap := map[string]interface{}{}
	var result []byte
	var customEntities map[string]interface{}

	// unmarshal core config into the merge map
	err = json.Unmarshal(kongCoreConfig, &mergeMap)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling kong config into map[string]interface{}: %w", err)
	}

	// unmarshal custom entities config into the merge map
	err = json.Unmarshal(customEntitiesJSONBytes, &customEntities)
	if err != nil {
		// do not error out when custom entities are messed up
		log.Errorf("failed to unmarshal custom entities from secret data: %v", err)
	} else {
		for k, v := range customEntities {
			if _, exists := mergeMap[k]; !exists {
				mergeMap[k] = v
			}
		}
	}

	// construct the final configuration
	result, err = json.Marshal(mergeMap)
	if err != nil {
		err = fmt.Errorf("marshaling final config into JSON: %w", err)
		return nil, err
	}

	return result, nil
}

func (n *KongController) fetchCustomEntities() ([]byte, error) {
	ns, name, err := util.ParseNameNS(n.cfg.KongCustomEntitiesSecret)
	if err != nil {
		return nil, fmt.Errorf("parsing kong custom entities secret: %w", err)
	}
	secret, err := n.store.GetSecret(ns, name)
	if err != nil {
		return nil, fmt.Errorf("fetching secret: %w", err)
	}
	config, ok := secret.Data["config"]
	if !ok {
		return nil, fmt.Errorf("'config' key not found in "+
			"custom entities secret '%v'", n.cfg.KongCustomEntitiesSecret)
	}
	return config, nil
}

func onUpdateInMemoryMode(ctx context.Context,
	log logrus.FieldLogger,
	state *file.Content,
	customEntities []byte,
	kongConfig *Kong,
) error {
	// Kong will error out if this is set
	state.Info = nil
	// Kong errors out if `null`s are present in `config` of plugins
	deckgen.CleanUpNullsInPluginConfigs(state)

	config, err := renderConfigWithCustomEntities(log, state, customEntities)
	if err != nil {
		return fmt.Errorf("constructing kong configuration: %w", err)
	}

	req, err := http.NewRequest("POST", kongConfig.URL+"/config",
		bytes.NewReader(config))
	if err != nil {
		return fmt.Errorf("creating new HTTP request for /config: %w", err)
	}
	req.Header.Add("content-type", "application/json")

	queryString := req.URL.Query()
	queryString.Add("check_hash", "1")

	req.URL.RawQuery = queryString.Encode()

	_, err = kongConfig.Client.Do(ctx, req, nil)
	if err != nil {
		return fmt.Errorf("posting new config to /config: %w", err)
	}

	return err
}

func onUpdateDBMode(
	targetContent *file.Content,
	kongConfig *Kong,
	selectorTags []string,
) error {
	// read the current state
	rawState, err := dump.Get(kongConfig.Client, dump.Config{
		SelectorTags: selectorTags,
	})
	if err != nil {
		return fmt.Errorf("loading configuration from kong: %w", err)
	}
	currentState, err := state.Get(rawState)
	if err != nil {
		return err
	}

	// read the target state
	rawState, err = file.Get(targetContent, file.RenderConfig{
		CurrentState: currentState,
		KongVersion:  kongConfig.Version,
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
		return fmt.Errorf("creating a new syncer: %w", err)
	}
	syncer.SilenceWarnings = true
	//client.SetDebugMode(true)
	_, errs := solver.Solve(nil, syncer, kongConfig.Client, kongConfig.Concurrency, false)
	if errs != nil {
		return deckutils.ErrArray{Errors: errs}
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

func toDeckContent(
	ctx context.Context,
	log logrus.FieldLogger,
	k8sState *kongstate.KongState,
	schemas *util.PluginSchemaStore,
	selectorTags []string,
) *file.Content {
	var content file.Content
	content.FormatVersion = "1.1"
	var err error

	for _, s := range k8sState.Services {
		service := file.FService{Service: s.Service}
		for _, p := range s.Plugins {
			plugin := file.FPlugin{
				Plugin: *p.DeepCopy(),
			}
			err = fillPlugin(ctx, &plugin, schemas)
			if err != nil {
				log.Errorf("failed to fill-in defaults for plugin: %s", *plugin.Name)
			}
			service.Plugins = append(service.Plugins, &plugin)
			sort.SliceStable(service.Plugins, func(i, j int) bool {
				return strings.Compare(*service.Plugins[i].Name, *service.Plugins[j].Name) > 0
			})
		}

		for _, r := range s.Routes {
			route := file.FRoute{Route: r.Route}
			fillRoute(&route.Route)

			for _, p := range r.Plugins {
				plugin := file.FPlugin{
					Plugin: *p.DeepCopy(),
				}
				err = fillPlugin(ctx, &plugin, schemas)
				if err != nil {
					log.Errorf("failed to fill-in defaults for plugin: %s", *plugin.Name)
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
		err = fillPlugin(ctx, &plugin, schemas)
		if err != nil {
			log.Errorf("failed to fill-in defaults for plugin: %s", *plugin.Name)
		}
		content.Plugins = append(content.Plugins, plugin)
	}
	sort.SliceStable(content.Plugins, func(i, j int) bool {
		return strings.Compare(deckgen.PluginString(content.Plugins[i]),
			deckgen.PluginString(content.Plugins[j])) > 0
	})

	for _, u := range k8sState.Upstreams {
		fillUpstream(&u.Upstream)
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
		cert := deckgen.GetFCertificateFromKongCert(c.Certificate)
		content.Certificates = append(content.Certificates, cert)
	}
	sort.SliceStable(content.Certificates, func(i, j int) bool {
		return strings.Compare(*content.Certificates[i].Cert, *content.Certificates[j].Cert) > 0
	})

	for _, c := range k8sState.CACertificates {
		content.CACertificates = append(content.CACertificates,
			file.FCACertificate{CACertificate: c})
	}
	sort.SliceStable(content.CACertificates, func(i, j int) bool {
		return strings.Compare(*content.CACertificates[i].Cert, *content.CACertificates[j].Cert) > 0
	})

	for _, c := range k8sState.Consumers {
		consumer := file.FConsumer{Consumer: c.Consumer}
		for _, p := range c.Plugins {
			consumer.Plugins = append(consumer.Plugins, &file.FPlugin{Plugin: p})
		}

		for _, v := range c.KeyAuths {
			consumer.KeyAuths = append(consumer.KeyAuths, &v.KeyAuth)
		}
		for _, v := range c.HMACAuths {
			consumer.HMACAuths = append(consumer.HMACAuths, &v.HMACAuth)
		}
		for _, v := range c.BasicAuths {
			consumer.BasicAuths = append(consumer.BasicAuths, &v.BasicAuth)
		}
		for _, v := range c.JWTAuths {
			consumer.JWTAuths = append(consumer.JWTAuths, &v.JWTAuth)
		}
		for _, v := range c.ACLGroups {
			consumer.ACLGroups = append(consumer.ACLGroups, &v.ACLGroup)
		}
		for _, v := range c.Oauth2Creds {
			consumer.Oauth2Creds = append(consumer.Oauth2Creds, &v.Oauth2Credential)
		}
		content.Consumers = append(content.Consumers, consumer)
	}
	sort.SliceStable(content.Consumers, func(i, j int) bool {
		return strings.Compare(*content.Consumers[i].Username, *content.Consumers[j].Username) > 0
	})
	if len(selectorTags) > 0 {
		content.Info = &file.Info{
			SelectorTags: selectorTags,
		}
	}

	return &content
}

func fillRoute(route *kong.Route) {
	if route.HTTPSRedirectStatusCode == nil {
		route.HTTPSRedirectStatusCode = kong.Int(426)
	}
	if route.PathHandling == nil {
		route.PathHandling = kong.String("v0")
	}
}

func fillUpstream(upstream *kong.Upstream) {
	if upstream.Algorithm == nil {
		upstream.Algorithm = kong.String("round-robin")
	}
}

func fillPlugin(ctx context.Context, plugin *file.FPlugin, schemas *util.PluginSchemaStore) error {
	if plugin == nil {
		return fmt.Errorf("plugin is nil")
	}
	if plugin.Name == nil || *plugin.Name == "" {
		return fmt.Errorf("plugin doesn't have a name")
	}
	schema, err := schemas.Schema(ctx, *plugin.Name)
	if err != nil {
		return fmt.Errorf("error retrieveing schema for plugin %s: %w", *plugin.Name, err)
	}
	if plugin.Config == nil {
		plugin.Config = make(kong.Configuration)
	}
	newConfig, err := deckgen.FillPluginConfig(schema, plugin.Config)
	if err != nil {
		return fmt.Errorf("error filling in default for plugin %s: %w", *plugin.Name, err)
	}
	plugin.Config = newConfig
	if plugin.RunOn == nil {
		plugin.RunOn = kong.String("first")
	}
	if plugin.Enabled == nil {
		plugin.Enabled = kong.Bool(true)
	}
	if len(plugin.Protocols) == 0 {
		// TODO read this from the schema endpoint
		plugin.Protocols = kong.StringSlice("http", "https")
	}
	plugin.RunOn = nil
	return nil
}
