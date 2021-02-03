package kongstate

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/pkg/annotations"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/pkg/store"
	"github.com/kong/kubernetes-ingress-controller/pkg/util"
	corev1 "k8s.io/api/core/v1"
)

func getKongIngressForService(s store.Storer, service corev1.Service) (
	*configurationv1.KongIngress, error) {
	confName := annotations.ExtractConfigurationName(service.Annotations)
	if confName == "" {
		return nil, nil
	}
	return s.GetKongIngress(service.Namespace, confName)
}

func getKongIngressFromObjectMeta(s store.Storer, obj *util.K8sObjectInfo) (
	*configurationv1.KongIngress, error) {
	return getKongIngressFromIngressAnnotations(s, obj.Namespace, obj.Name, obj.Annotations)
}

func getKongIngressFromIngressAnnotations(s store.Storer, namespace, name string,
	anns map[string]string) (
	*configurationv1.KongIngress, error) {
	confName := annotations.ExtractConfigurationName(anns)
	if confName != "" {
		ki, err := s.GetKongIngress(namespace, confName)
		if err == nil {
			return ki, nil
		}
	}

	ki, err := s.GetKongIngress(namespace, name)
	if err == nil {
		return ki, nil
	}
	return nil, nil
}

// getPlugin constructs a plugins from a KongPlugin resource.
func getPlugin(s store.Storer, namespace, name string) (kong.Plugin, error) {
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
	return plugin, err
}

func kongPluginFromK8SClusterPlugin(
	s store.Storer,
	k8sPlugin configurationv1.KongClusterPlugin) (kong.Plugin, error) {
	config := k8sPlugin.Config
	if k8sPlugin.ConfigFrom.SecretValue !=
		(configurationv1.NamespacedSecretValueFromSource{}) &&
		len(k8sPlugin.Config) > 0 {
		return kong.Plugin{},
			fmt.Errorf("KongClusterPlugin '/%v' has both "+
				"Config and ConfigFrom set", k8sPlugin.Name)
	}
	if k8sPlugin.ConfigFrom.SecretValue != (configurationv1.
		NamespacedSecretValueFromSource{}) {
		var err error
		config, err = namespacedSecretToConfiguration(
			s,
			k8sPlugin.ConfigFrom.SecretValue)
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
		Disabled:  k8sPlugin.Disabled,
		Protocols: k8sPlugin.Protocols,
	}.toKongPlugin()
	return kongPlugin, nil
}

func cloneStringPointerSlice(array ...*string) (res []*string) {
	res = append(res, array...)
	return
}

func kongPluginFromK8SPlugin(
	s store.Storer,
	k8sPlugin configurationv1.KongPlugin) (kong.Plugin, error) {
	config := k8sPlugin.Config
	if k8sPlugin.ConfigFrom.SecretValue !=
		(configurationv1.SecretValueFromSource{}) &&
		len(k8sPlugin.Config) > 0 {
		return kong.Plugin{},
			fmt.Errorf("KongPlugin '%v/%v' has both "+
				"Config and ConfigFrom set",
				k8sPlugin.Namespace, k8sPlugin.Name)
	}
	if k8sPlugin.ConfigFrom.SecretValue !=
		(configurationv1.SecretValueFromSource{}) {
		var err error
		config, err = SecretToConfiguration(s,
			k8sPlugin.ConfigFrom.SecretValue, k8sPlugin.Namespace)
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
		Disabled:  k8sPlugin.Disabled,
		Protocols: k8sPlugin.Protocols,
	}.toKongPlugin()
	return kongPlugin, nil
}

func namespacedSecretToConfiguration(
	s store.Storer,
	reference configurationv1.NamespacedSecretValueFromSource) (
	configurationv1.Configuration, error) {
	bareReference := configurationv1.SecretValueFromSource{
		Secret: reference.Secret,
		Key:    reference.Key}
	return SecretToConfiguration(s, bareReference, reference.Namespace)
}

func SecretToConfiguration(
	s store.Storer,
	reference configurationv1.SecretValueFromSource, namespace string) (
	configurationv1.Configuration, error) {
	secret, err := s.GetSecret(namespace, reference.Secret)
	if err != nil {
		return configurationv1.Configuration{}, fmt.Errorf(
			"error fetching plugin configuration secret '%v/%v': %v",
			namespace, reference.Secret, err)
	}
	secretVal, ok := secret.Data[reference.Key]
	if !ok {
		return configurationv1.Configuration{},
			fmt.Errorf("no key '%v' in secret '%v/%v'",
				reference.Key, namespace, reference.Secret)
	}
	var config configurationv1.Configuration
	if err := json.Unmarshal(secretVal, &config); err != nil {
		if err := yaml.Unmarshal(secretVal, &config); err != nil {
			return configurationv1.Configuration{},
				fmt.Errorf("key '%v' in secret '%v/%v' contains neither "+
					"valid JSON nor valid YAML)",
					reference.Key, namespace, reference.Secret)
		}
	}
	return config, nil
}

// plugin is a intermediate type to hold plugin related configuration
type plugin struct {
	Name   string
	Config configurationv1.Configuration

	RunOn     string
	Disabled  bool
	Protocols []string
}

func (p plugin) toKongPlugin() kong.Plugin {
	result := kong.Plugin{
		Name:   kong.String(p.Name),
		Config: kong.Configuration(p.Config).DeepCopy(),
	}
	if p.RunOn != "" {
		result.RunOn = kong.String(p.RunOn)
	}
	if p.Disabled {
		result.Enabled = kong.Bool(false)
	}
	if len(p.Protocols) > 0 {
		result.Protocols = kong.StringSlice(p.Protocols...)
	}
	return result
}
