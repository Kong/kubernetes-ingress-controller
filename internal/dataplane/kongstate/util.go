package kongstate

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/samber/mo"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/yaml"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
)

func getKongIngressForServices(
	s store.Storer,
	services []*corev1.Service,
) (*kongv1.KongIngress, error) {
	// loop through each service and retrieve the attached KongIngress resources.
	// there can only be one KongIngress for a group of services: either one of
	// them is configured with a KongIngress and this configures the Kong Service
	// or Upstream OR all of them can be configured but they must be configured
	// with the same KongIngress.
	for _, svc := range services {
		// check if the service is even configured with a KongIngress
		confName := annotations.ExtractConfigurationName(svc.Annotations)
		if confName == "" {
			continue // some other service in the group may yet have a KongIngress attachment
		}

		// retrieve the attached KongIngress for the service
		kongIngress, err := s.GetKongIngress(svc.Namespace, confName)
		if err != nil {
			return nil, err
		}

		// we found the KongIngress for these services. We don't have to check any
		// further services as validation is expected to ensure all these Services
		// already are annotated with the exact same overrides.
		return kongIngress, nil
	}

	// there are no KongIngress resources for these services.
	return nil, nil
}

// getKongUpstreamPolicyForServices scans all Services in the group to see if their KongUpstreamPolicy is consistent
// and returns the KongUpstreamPolicy if it is.
// We require either:
// - all of them to be configured with the same KongUpstreamPolicy.
// - none of them to be configured with a KongUpstreamPolicy.
func getKongUpstreamPolicyForServices(
	s store.Storer,
	servicesGroup []*corev1.Service,
) (*kongv1beta1.KongUpstreamPolicy, error) {
	if len(servicesGroup) == 0 {
		return nil, nil
	}

	servicesGroupedByUpstreamPolicy := lo.GroupBy(servicesGroup, func(svc *corev1.Service) mo.Option[string] {
		policyName, ok := annotations.ExtractUpstreamPolicy(svc.Annotations)
		if !ok {
			return mo.None[string]()
		}
		return mo.Some(policyName)
	})

	// If there's more than one group, then there are services with different KongUpstreamPolicy configurations.
	if len(servicesGroupedByUpstreamPolicy) > 1 {
		return nil, fmt.Errorf("inconsistent KongUpstreamPolicy configuration for services %s",
			PrettyPrintServiceList(servicesGroup))
	}

	// If there's one group (must be at least one, since we checked len(servicesGroup) == 0 above), then
	// there's either one KongUpstreamPolicy for all services, or none.
	upstreamPolicyName, ok := lo.Keys(servicesGroupedByUpstreamPolicy)[0].Get()
	if !ok {
		return nil, nil
	}

	policy, err := s.GetKongUpstreamPolicy(servicesGroup[0].Namespace, upstreamPolicyName)
	if err != nil {
		return nil, fmt.Errorf("failed fetching KongUpstreamPolicy: %w", err)
	}

	return policy, nil
}

func getKongIngressFromObjectMeta(
	s store.Storer,
	obj util.K8sObjectInfo,
) (
	*kongv1.KongIngress, error,
) {
	return getKongIngressFromObjAnnotations(s, obj)
}

func getKongIngressFromObjAnnotations(
	s store.Storer,
	obj util.K8sObjectInfo,
) (
	*kongv1.KongIngress, error,
) {
	confName := annotations.ExtractConfigurationName(obj.Annotations)
	if confName != "" {
		ki, err := s.GetKongIngress(obj.Namespace, confName)
		if err == nil {
			return ki, nil
		}
	}

	ki, err := s.GetKongIngress(obj.Namespace, obj.Name)
	if err == nil {
		return ki, nil
	}
	return nil, nil
}

// getKongPluginOrKongClusterPlugin fetches a KongPlugin or KongClusterPlugin (as fallback) from the store.
// If both are not found, an error is returned.
func getKongPluginOrKongClusterPlugin(s store.Storer, namespace, name string) (
	*kongv1.KongPlugin,
	*kongv1.KongClusterPlugin,
	error,
) {
	plugin, pluginErr := s.GetKongPlugin(namespace, name)
	if pluginErr != nil {
		if !errors.As(pluginErr, &store.NotFoundError{}) {
			return nil, nil, fmt.Errorf("failed fetching KongPlugin: %w", pluginErr)
		}

		// If KongPlugin is not found, try to fetch KongClusterPlugin.
		clusterPlugin, err := s.GetKongClusterPlugin(name)
		if err != nil {
			if !errors.As(err, &store.NotFoundError{}) {
				return nil, nil, fmt.Errorf("failed fetching KongClusterPlugin: %w", err)
			}

			// Both KongPlugin and KongClusterPlugin are not found.
			return nil, nil, fmt.Errorf("no KongPlugin or KongClusterPlugin was found for %s/%s", namespace, name)
		}

		return nil, clusterPlugin, nil
	}

	return plugin, nil, nil
}

func kongPluginFromK8SClusterPlugin(
	s store.Storer,
	k8sPlugin kongv1.KongClusterPlugin,
) (Plugin, error) {
	var config kong.Configuration
	config, err := RawConfigToConfiguration(k8sPlugin.Config)
	if err != nil {
		return Plugin{}, fmt.Errorf("could not parse KongPlugin %s/%s config: %w",
			k8sPlugin.Namespace, k8sPlugin.Name, err)
	}
	if k8sPlugin.ConfigFrom != nil && len(config) > 0 {
		return Plugin{},
			fmt.Errorf("KongClusterPlugin '/%v' has both "+
				"Config and ConfigFrom set", k8sPlugin.Name)
	}
	if k8sPlugin.ConfigFrom != nil {
		var err error
		config, err = namespacedSecretToConfiguration(
			s,
			(*k8sPlugin.ConfigFrom).SecretValue)
		if err != nil {
			return Plugin{},
				fmt.Errorf("error parsing config for KongClusterPlugin %s: %w",
					k8sPlugin.Name, err)
		}
	}

	return Plugin{
		Plugin: plugin{
			Name:   k8sPlugin.PluginName,
			Config: config,

			RunOn:        k8sPlugin.RunOn,
			Ordering:     k8sPlugin.Ordering,
			InstanceName: k8sPlugin.InstanceName,
			Disabled:     k8sPlugin.Disabled,
			Protocols:    protocolsToStrings(k8sPlugin.Protocols),
			Tags:         util.GenerateTagsForObject(&k8sPlugin),
		}.toKongPlugin(),
		K8sParent: &k8sPlugin,
	}, nil
}

func protocolsToStrings(protocols []kongv1.KongProtocol) (res []string) {
	for _, protocol := range protocols {
		res = append(res, string(protocol))
	}
	return
}

func kongPluginFromK8SPlugin(
	s store.Storer,
	k8sPlugin kongv1.KongPlugin,
) (Plugin, error) {
	var config kong.Configuration
	config, err := RawConfigToConfiguration(k8sPlugin.Config)
	if err != nil {
		return Plugin{}, fmt.Errorf("could not parse KongPlugin %s/%s config: %w",
			k8sPlugin.Namespace, k8sPlugin.Name, err)
	}
	if k8sPlugin.ConfigFrom != nil && len(config) > 0 {
		return Plugin{},
			fmt.Errorf("KongPlugin '%s/%s' has both Config and ConfigFrom set",
				k8sPlugin.Namespace, k8sPlugin.Name)
	}
	if k8sPlugin.ConfigFrom != nil {
		var err error
		config, err = SecretToConfiguration(s,
			(*k8sPlugin.ConfigFrom).SecretValue, k8sPlugin.Namespace)
		if err != nil {
			return Plugin{},
				fmt.Errorf("error parsing config for KongPlugin '%s/%s': %w",
					k8sPlugin.Name, k8sPlugin.Namespace, err)
		}
	}

	return Plugin{
		Plugin: plugin{
			Name:   k8sPlugin.PluginName,
			Config: config,

			RunOn:        k8sPlugin.RunOn,
			Ordering:     k8sPlugin.Ordering,
			InstanceName: k8sPlugin.InstanceName,
			Disabled:     k8sPlugin.Disabled,
			Protocols:    protocolsToStrings(k8sPlugin.Protocols),
			Tags:         util.GenerateTagsForObject(&k8sPlugin),
		}.toKongPlugin(),
		K8sParent: &k8sPlugin,
	}, nil
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
	reference kongv1.NamespacedSecretValueFromSource) (
	kong.Configuration, error,
) {
	bareReference := kongv1.SecretValueFromSource{
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
	reference kongv1.SecretValueFromSource, namespace string) (
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

// PrettyPrintServiceList makes a clean printable list of Kubernetes
// services for the purpose of logging (errors, info, e.t.c.).
func PrettyPrintServiceList(services []*corev1.Service) string {
	var serviceList []string
	for _, svc := range services {
		serviceList = append(serviceList, svc.Namespace+"/"+svc.Name)
	}
	return strings.Join(serviceList, ", ")
}

// plugin is a intermediate type to hold plugin related configuration.
type plugin struct {
	Name   string
	Config kong.Configuration

	RunOn        string
	Ordering     *kong.PluginOrdering
	InstanceName string
	Disabled     bool
	Protocols    []string
	Tags         []*string
}

func (p plugin) toKongPlugin() kong.Plugin {
	result := kong.Plugin{
		Name:   kong.String(p.Name),
		Config: p.Config.DeepCopy(),
		Tags:   p.Tags,
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
	if p.InstanceName != "" {
		result.InstanceName = kong.String(p.InstanceName)
	}
	return result
}
