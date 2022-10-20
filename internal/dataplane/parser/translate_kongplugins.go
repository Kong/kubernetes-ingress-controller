package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	"github.com/sirupsen/logrus"
)

func (p *Parser) getPlugins(log logrus.FieldLogger, s store.Storer, pluginRels map[string]util.ForeignRelations) []kongstate.Plugin {
	var plugins []kongstate.Plugin

	for pluginIdentifier, relations := range pluginRels {
		identifier := strings.Split(pluginIdentifier, ":")
		namespace, kongPluginName := identifier[0], identifier[1]
		plugin, err := p.getPlugin(s, namespace, kongPluginName)
		if err != nil {
			log.WithFields(logrus.Fields{
				"kongplugin_name":      kongPluginName,
				"kongplugin_namespace": namespace,
			}).WithError(err).Errorf("failed to fetch KongPlugin")
			continue
		}

		for _, rel := range relations.GetCombinations() {
			plugin := *plugin.DeepCopy()
			// ID is populated because that is read by decK and in_memory
			// translator too
			if rel.Service != "" {
				plugin.Service = &kong.Service{ID: kong.String(rel.Service)}
			}
			if rel.Route != "" {
				plugin.Route = &kong.Route{ID: kong.String(rel.Route)}
			}
			if rel.Consumer != "" {
				plugin.Consumer = &kong.Consumer{ID: kong.String(rel.Consumer)}
			}
			plugins = append(plugins, kongstate.Plugin{plugin})
		}
	}

	globalPlugins, err := globalPlugins(log, s)
	if err != nil {
		log.WithError(err).Error("failed to fetch global plugins")
	}
	plugins = append(plugins, globalPlugins...)

	return plugins
}

func globalPlugins(log logrus.FieldLogger, s store.Storer) ([]kongstate.Plugin, error) {
	// removed as of 0.10.0
	// only retrieved now to warn users
	globalPlugins, err := s.ListGlobalKongPlugins()
	if err != nil {
		return nil, fmt.Errorf("error listing global KongPlugins: %w", err)
	}
	if len(globalPlugins) > 0 {
		log.Warning("global KongPlugins found. These are no longer applied and",
			" must be replaced with KongClusterPlugins.",
			" Please run \"kubectl get kongplugin -l global=true --all-namespaces\" to list existing plugins")
	}
	res := make(map[string]kongstate.Plugin)
	var duplicates []string // keep track of duplicate
	// TODO respect the oldest CRD
	// Current behavior is to skip creating the plugin but in case
	// of duplicate plugin definitions, we should respect the oldest one
	// This is important since if a user comes in to k8s and creates a new
	// CRD, the user now deleted an older plugin

	globalClusterPlugins, err := s.ListGlobalKongClusterPlugins()
	if err != nil {
		return nil, fmt.Errorf("error listing global KongClusterPlugins: %w", err)
	}
	for i := 0; i < len(globalClusterPlugins); i++ {
		k8sPlugin := *globalClusterPlugins[i]
		pluginName := k8sPlugin.PluginName
		// empty pluginName skip it
		if pluginName == "" {
			log.WithFields(logrus.Fields{
				"kongclusterplugin_name": k8sPlugin.Name,
			}).Errorf("invalid KongClusterPlugin: empty plugin property")
			continue
		}
		if _, ok := res[pluginName]; ok {
			log.Error("multiple KongPlugin definitions found with"+
				" 'global' label for '", pluginName,
				"', the plugin will not be applied")
			duplicates = append(duplicates, pluginName)
			continue
		}
		if plugin, err := kongPluginFromK8SClusterPlugin(s, k8sPlugin); err == nil {
			res[pluginName] = kongstate.Plugin{
				Plugin: plugin,
			}
		} else {
			log.WithFields(logrus.Fields{
				"kongclusterplugin_name": k8sPlugin.Name,
			}).WithError(err).Error("failed to generate configuration from KongClusterPlugin")
		}
	}
	for _, plugin := range duplicates {
		delete(res, plugin)
	}
	var plugins []kongstate.Plugin
	for _, p := range res {
		plugins = append(plugins, p)
	}
	return plugins, nil
}

// getPlugin constructs a plugins from a KongPlugin resource.
func (p *Parser) getPlugin(s store.Storer, namespace, name string) (kong.Plugin, error) {
	var plugin kong.Plugin
	k8sPlugin, err := s.GetKongPlugin(namespace, name)
	if err != nil {
		// if no namespaced plugin definition, then
		// search for cluster level-plugin definition
		if errors.As(err, &store.ErrNotFound{}) {
			clusterPlugin, err := s.GetKongClusterPlugin(name)
			// not found
			if errors.As(err, &store.ErrNotFound{}) {
				return plugin, errors.New(
					"no KongPlugin or KongClusterPlugin was found")
			}
			if err != nil {
				return plugin, err
			}
			if clusterPlugin.PluginName == "" {
				return plugin, fmt.Errorf("invalid empty 'plugin' property")
			}
			plugin, err = kongPluginFromK8SClusterPlugin(s, *clusterPlugin)
			return plugin, err
		}
	}
	// ignore plugins with no name
	if k8sPlugin.PluginName == "" {
		return plugin, fmt.Errorf("invalid empty 'plugin' property")
	}

	plugin, err = kongPluginFromK8SPlugin(s, *k8sPlugin)
	if err != nil {
		return kong.Plugin{}, err
	}

	p.ReportKubernetesObjectUpdate(k8sPlugin)
	return plugin, err
}

func kongPluginFromK8SClusterPlugin(
	s store.Storer,
	k8sPlugin configurationv1.KongClusterPlugin,
) (kong.Plugin, error) {
	var config kong.Configuration
	config, err := RawConfigToConfiguration(k8sPlugin.Config)
	if err != nil {
		return kong.Plugin{}, fmt.Errorf("could not parse KongPlugin %v/%v config: %w",
			k8sPlugin.Namespace, k8sPlugin.Name, err)
	}
	if k8sPlugin.ConfigFrom != nil && len(config) > 0 {
		return kong.Plugin{},
			fmt.Errorf("KongClusterPlugin '/%v' has both "+
				"Config and ConfigFrom set", k8sPlugin.Name)
	}
	if k8sPlugin.ConfigFrom != nil {
		var err error
		config, err = namespacedSecretToConfiguration(
			s,
			(*k8sPlugin.ConfigFrom).SecretValue)
		if err != nil {
			return kong.Plugin{},
				fmt.Errorf("error parsing config for KongClusterPlugin %v: %w",
					k8sPlugin.Name, err)
		}
	}
	kongPlugin := plugin{
		Name:   k8sPlugin.PluginName,
		Config: config,

		RunOn:     k8sPlugin.RunOn,
		Ordering:  k8sPlugin.Ordering,
		Disabled:  k8sPlugin.Disabled,
		Protocols: protocolsToStrings(k8sPlugin.Protocols),
	}.toKongPlugin()
	return kongPlugin, nil
}

func protocolsToStrings(protocols []configurationv1.KongProtocol) (res []string) {
	for _, protocol := range protocols {
		res = append(res, string(protocol))
	}
	return
}

func kongPluginFromK8SPlugin(
	s store.Storer,
	k8sPlugin configurationv1.KongPlugin,
) (kong.Plugin, error) {
	var config kong.Configuration
	config, err := RawConfigToConfiguration(k8sPlugin.Config)
	if err != nil {
		return kong.Plugin{}, fmt.Errorf("could not parse KongPlugin %v/%v config: %w",
			k8sPlugin.Namespace, k8sPlugin.Name, err)
	}
	if k8sPlugin.ConfigFrom != nil && len(config) > 0 {
		return kong.Plugin{},
			fmt.Errorf("KongPlugin '%v/%v' has both "+
				"Config and ConfigFrom set",
				k8sPlugin.Namespace, k8sPlugin.Name)
	}
	if k8sPlugin.ConfigFrom != nil {
		var err error
		config, err = SecretToConfiguration(s,
			(*k8sPlugin.ConfigFrom).SecretValue, k8sPlugin.Namespace)
		if err != nil {
			return kong.Plugin{},
				fmt.Errorf("error parsing config for KongPlugin '%v/%v': %w",
					k8sPlugin.Name, k8sPlugin.Namespace, err)
		}
	}
	kongPlugin := plugin{
		Name:   k8sPlugin.PluginName,
		Config: config,

		RunOn:     k8sPlugin.RunOn,
		Ordering:  k8sPlugin.Ordering,
		Disabled:  k8sPlugin.Disabled,
		Protocols: protocolsToStrings(k8sPlugin.Protocols),
	}.toKongPlugin()
	return kongPlugin, nil
}

func RawConfigToConfiguration(config apiextensionsv1.JSON) (kong.Configuration, error) {
	if len(config.Raw) == 0 {
		return kong.Configuration{}, nil
	}
	var kongConfig kong.Configuration
	err := json.Unmarshal(config.Raw, &kongConfig)
	if err != nil {
		return kong.Configuration{}, err
	}
	return kongConfig, nil
}

func namespacedSecretToConfiguration(
	s store.Storer,
	reference configurationv1.NamespacedSecretValueFromSource) (
	kong.Configuration, error,
) {
	bareReference := configurationv1.SecretValueFromSource{
		Secret: reference.Secret,
		Key:    reference.Key,
	}
	return SecretToConfiguration(s, bareReference, reference.Namespace)
}

type SecretGetter interface {
	GetSecret(namespace, name string) (*corev1.Secret, error)
}

func SecretToConfiguration(
	s SecretGetter,
	reference configurationv1.SecretValueFromSource, namespace string) (
	kong.Configuration, error,
) {
	secret, err := s.GetSecret(namespace, reference.Secret)
	if err != nil {
		return kong.Configuration{}, fmt.Errorf(
			"error fetching plugin configuration secret '%v/%v': %w",
			namespace, reference.Secret, err)
	}
	secretVal, ok := secret.Data[reference.Key]
	if !ok {
		return kong.Configuration{},
			fmt.Errorf("no key '%v' in secret '%v/%v'",
				reference.Key, namespace, reference.Secret)
	}
	var config kong.Configuration
	if err := json.Unmarshal(secretVal, &config); err != nil {
		if err := yaml.Unmarshal(secretVal, &config); err != nil {
			return kong.Configuration{},
				fmt.Errorf("key '%v' in secret '%v/%v' contains neither "+
					"valid JSON nor valid YAML)",
					reference.Key, namespace, reference.Secret)
		}
	}
	return config, nil
}

// plugin is a intermediate type to hold plugin related configuration.
type plugin struct {
	Name   string
	Config kong.Configuration

	RunOn     string
	Ordering  *kong.PluginOrdering
	Disabled  bool
	Protocols []string
}

func (p plugin) toKongPlugin() kong.Plugin {
	result := kong.Plugin{
		Name:   kong.String(p.Name),
		Config: p.Config.DeepCopy(),
	}
	if p.RunOn != "" {
		result.RunOn = kong.String(p.RunOn)
	}
	if p.Disabled {
		result.Enabled = kong.Bool(false)
	}
	result.Ordering = p.Ordering
	if len(p.Protocols) > 0 {
		result.Protocols = kong.StringSlice(p.Protocols...)
	}
	return result
}
