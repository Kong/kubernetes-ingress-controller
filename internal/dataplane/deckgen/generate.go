package deckgen

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/go-logr/logr"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/versions"
)

// StubUpstreamName is a name of a stub upstream that is created when the configuration is empty.
const StubUpstreamName = "kong"

type PluginSchemaStore interface {
	Schema(ctx context.Context, pluginName string) (map[string]interface{}, error)
}

// GenerateDeckContentParams is the parameters used to generate deck contents.
type GenerateDeckContentParams struct {
	SelectorTags     []string
	ExpressionRoutes bool
	PluginSchemas    PluginSchemaStore

	// AppendStubEntityWhenConfigEmpty indicates whether to append a stub entity to the configuration when
	// the configuration is empty. It is used to workaround behavior in Kong where sending an empty configuration
	// does not make its `GET /status/ready` endpoint return 200s.
	AppendStubEntityWhenConfigEmpty bool
}

// ToDeckContent generates a decK configuration from `k8sState` and auxiliary parameters.
func ToDeckContent(
	ctx context.Context,
	logger logr.Logger,
	k8sState *kongstate.KongState,
	params GenerateDeckContentParams,
) *file.Content {
	var content file.Content
	content.FormatVersion = versions.DeckFileFormatVersion

	for _, s := range k8sState.Services {
		service := file.FService{Service: s.Service}
		for _, p := range s.Plugins {
			plugin := file.FPlugin{
				Plugin: *p.DeepCopy(),
			}
			if err := fillPlugin(ctx, &plugin, params.PluginSchemas); err != nil {
				logger.Error(err, "Failed to fill in defaults for plugin", "plugin_name", *plugin.Name)
			}
			service.Plugins = append(service.Plugins, &plugin)
			sort.SliceStable(service.Plugins, func(i, j int) bool {
				return strings.Compare(*service.Plugins[i].Name, *service.Plugins[j].Name) > 0
			})
		}

		for _, r := range s.Routes {
			route := file.FRoute{Route: r.Route}
			fillRoute(&route.Route, params.ExpressionRoutes)

			for _, p := range r.Plugins {
				plugin := file.FPlugin{
					Plugin: *p.DeepCopy(),
				}
				if err := fillPlugin(ctx, &plugin, params.PluginSchemas); err != nil {
					logger.Error(err, "Failed to fill in defaults for plugin", "plugin_name", *plugin.Name)
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
		if err := fillPlugin(ctx, &plugin, params.PluginSchemas); err != nil {
			logger.Error(err, "Failed to fill in defaults for plugin", "plugin_name", *plugin.Name)
		}
		content.Plugins = append(content.Plugins, plugin)
	}
	sort.SliceStable(content.Plugins, func(i, j int) bool {
		return strings.Compare(PluginString(content.Plugins[i]),
			PluginString(content.Plugins[j])) > 0
	})

	for _, cg := range k8sState.ConsumerGroups {
		consumerGroup := file.FConsumerGroupObject{ConsumerGroup: cg.ConsumerGroup}
		content.ConsumerGroups = append(content.ConsumerGroups, consumerGroup)
	}
	sort.SliceStable(content.ConsumerGroups, func(i, j int) bool {
		return strings.Compare(*content.ConsumerGroups[i].Name, *content.ConsumerGroups[j].Name) > 0
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
		cert := GetFCertificateFromKongCert(c.Certificate)
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

	for _, c := range k8sState.Licenses {
		content.Licenses = append(content.Licenses,
			file.FLicense{License: c.License})
	}

	sort.SliceStable(content.Licenses, func(i, j int) bool {
		return strings.Compare(*content.Licenses[i].Payload, *content.Licenses[j].Payload) > 0
	})

	for _, c := range k8sState.Consumers {
		consumer := file.FConsumer{Consumer: c.Consumer}

		// If a consumer with no username and no custom_id is provided deck wont be able to process it,
		// but we shouldn't fail the rest of the deckgen either or this will result in one bad consumer
		// being capable of stopping all updates to the Kong Admin API.
		// This shouldn't happen as we enforce either of those field being present in CRD CEL validation rules.
		if consumer.Username == nil && consumer.CustomID == nil {
			logger.Error(nil, "Invalid consumer received (username and custom_id were empty)")
			continue
		}

		for _, cg := range c.ConsumerGroups {
			consumer.Groups = append(consumer.Groups, &cg)
		}

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
		for _, v := range c.MTLSAuths {
			consumer.MTLSAuths = append(consumer.MTLSAuths, &v.MTLSAuth)
		}
		content.Consumers = append(content.Consumers, consumer)
	}
	sort.Stable(fConsumerByUsernameAndCustomID(content.Consumers))

	// convert vaults.
	for _, v := range k8sState.Vaults {
		vault := file.FVault{
			Vault: v.Vault,
		}
		content.Vaults = append(content.Vaults, vault)
	}
	sort.SliceStable(content.Vaults, func(i, j int) bool {
		return (*content.Vaults[i].Prefix) > (*content.Vaults[j].Prefix)
	})

	if len(params.SelectorTags) > 0 {
		content.Info = &file.Info{
			SelectorTags: params.SelectorTags,
		}
	}

	if params.AppendStubEntityWhenConfigEmpty && IsContentEmpty(&content) {
		content.Upstreams = append(content.Upstreams, file.FUpstream{
			Upstream: kong.Upstream{
				Name: lo.ToPtr(StubUpstreamName),
			},
		})
	}

	return &content
}

func fillRoute(route *kong.Route, expressionRoutes bool) {
	if route.HTTPSRedirectStatusCode == nil {
		route.HTTPSRedirectStatusCode = kong.Int(426)
	}
	if route.PathHandling == nil && !expressionRoutes {
		route.PathHandling = kong.String("v0")
	}
}

func fillUpstream(upstream *kong.Upstream) {
	if upstream.Algorithm == nil {
		upstream.Algorithm = kong.String("round-robin")
	}
}

func fillPlugin(ctx context.Context, plugin *file.FPlugin, schemas PluginSchemaStore) error {
	if plugin == nil {
		return fmt.Errorf("plugin is nil")
	}
	if plugin.Name == nil || *plugin.Name == "" {
		return fmt.Errorf("plugin doesn't have a name")
	}
	schema, err := schemas.Schema(ctx, *plugin.Name)
	if err != nil {
		return fmt.Errorf("error retrieving schema for plugin %s: %w", *plugin.Name, err)
	}
	if plugin.Config == nil {
		plugin.Config = make(kong.Configuration)
	}
	err = kong.FillPluginsDefaults(&plugin.Plugin, schema)
	if err != nil {
		return fmt.Errorf("error filling in default for plugin %s: %w", *plugin.Name, err)
	}
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
