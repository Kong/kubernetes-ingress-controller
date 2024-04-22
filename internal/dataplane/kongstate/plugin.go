package kongstate

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
)

// Plugin represents a plugin Object in Kong.
type Plugin struct {
	kong.Plugin
	K8sParent           client.Object
	SensitiveFieldsMeta PluginSensitiveFieldsMetadata
}

func (p Plugin) DeepCopy() Plugin {
	return Plugin{
		Plugin:              *p.Plugin.DeepCopy(),
		K8sParent:           p.K8sParent,
		SensitiveFieldsMeta: p.SensitiveFieldsMeta,
	}
}

func (p Plugin) SanitizedCopy() Plugin {
	// We do not want to return an error if any of below fails - the best we can do
	// is to return a plugin with wholly redacted config.
	// Let's have a closure returning a plugin with wholly redacted config prepared.
	whollySanitized := func() Plugin {
		p := p.DeepCopy()
		p.Config = sanitizeWholePluginConfig(p.Config)
		return p
	}

	// If the whole config is sensitive, we need to redact the entire config.
	if p.SensitiveFieldsMeta.WholeConfigIsSensitive {
		return whollySanitized()
	}

	// If there are JSON paths, we need to redact them.
	if len(p.SensitiveFieldsMeta.JSONPaths) > 0 {
		var patchOperations []string
		for _, path := range p.SensitiveFieldsMeta.JSONPaths {
			// If the path is empty, we need to sanitize the whole config.
			// An empty path means that the patch is on the root of the config.
			if path == "" {
				return whollySanitized()
			}

			patchOperations = append(patchOperations, fmt.Sprintf(
				`{"op":"replace","path":"%s","value":"%s"}`,
				path,
				*redactedString,
			))
		}

		// Decode the patch and apply it to the config.
		// We need to marshal the config to JSON and then unmarshal it back to Configuration
		// because the patch library works with bytes.
		patch, err := jsonpatch.DecodePatch([]byte(fmt.Sprintf("[%s]", strings.Join(patchOperations, ","))))
		if err != nil {
			return whollySanitized()
		}
		configB, err := json.Marshal(p.Config)
		if err != nil {
			return whollySanitized()
		}
		sanitizedConfigB, err := patch.Apply(configB)
		if err != nil {
			return whollySanitized()
		}
		sanitizedConfig := kong.Configuration{}
		if err := json.Unmarshal(sanitizedConfigB, &sanitizedConfig); err != nil {
			return whollySanitized()
		}

		sanitized := p.DeepCopy()
		sanitized.Config = sanitizedConfig
		return sanitized
	}

	// Nothing to sanitize.
	return p
}

// sanitizeWholePluginConfig redacts the entire config of a plugin by replacing all of its
// values with a redacted string.
func sanitizeWholePluginConfig(config kong.Configuration) kong.Configuration {
	sanitized := config.DeepCopy()
	for k := range config {
		sanitized[k] = *redactedString
	}
	return sanitized
}

// PluginSensitiveFieldsMetadata holds metadata about sensitive fields in a plugin's configuration.
// It can be used to sanitize them before exposing the configuration to the user (e.g. in debug dumps
// or in Konnect Admin API).
type PluginSensitiveFieldsMetadata struct {
	// WholeConfigIsSensitive indicates that the entire configuration of the plugin is sensitive.
	// If this is true, the configuration should be redacted entirely (each of its fields' values
	// should be replaced with a redacted string).
	WholeConfigIsSensitive bool

	// JSONPaths holds a list of JSON paths to sensitive fields in the plugin's configuration.
	// If this is not empty, the configuration should be redacted by replacing the values of the
	// fields at these paths with a redacted string.
	JSONPaths []string
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
	config, err := RawConfigurationWithNamespacedPatchesToConfiguration(
		s,
		k8sPlugin.Config,
		k8sPlugin.ConfigPatches,
	)
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
		config, err = NamespacedSecretToConfiguration(
			s,
			k8sPlugin.ConfigFrom.SecretValue)
		if err != nil {
			return Plugin{},
				fmt.Errorf("error parsing config for KongClusterPlugin %s: %w",
					k8sPlugin.Name, err)
		}
	}

	// Prepare sensitive fields metadata for the plugin.
	sensitiveFieldsMeta := PluginSensitiveFieldsMetadata{
		JSONPaths: lo.Map(k8sPlugin.ConfigPatches, func(patch kongv1.NamespacedConfigPatch, _ int) string {
			return patch.Path
		}),
		WholeConfigIsSensitive: k8sPlugin.ConfigFrom != nil,
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
		K8sParent:           &k8sPlugin,
		SensitiveFieldsMeta: sensitiveFieldsMeta,
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
	config, err := RawConfigurationWithPatchesToConfiguration(
		s,
		k8sPlugin.Namespace,
		k8sPlugin.Config,
		k8sPlugin.ConfigPatches,
	)
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
			k8sPlugin.ConfigFrom.SecretValue, k8sPlugin.Namespace)
		if err != nil {
			return Plugin{},
				fmt.Errorf("error parsing config for KongPlugin '%s/%s': %w",
					k8sPlugin.Name, k8sPlugin.Namespace, err)
		}
	}

	// Prepare sensitive fields metadata for the plugin.
	sensitiveFieldsMeta := PluginSensitiveFieldsMetadata{
		JSONPaths: lo.Map(k8sPlugin.ConfigPatches, func(patch kongv1.ConfigPatch, _ int) string {
			return patch.Path
		}),
		WholeConfigIsSensitive: k8sPlugin.ConfigFrom != nil,
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
		K8sParent:           &k8sPlugin,
		SensitiveFieldsMeta: sensitiveFieldsMeta,
	}, nil
}

var rawPatchPattern = `[{"op":"%s","path":"%s","value":%s}]`

type JSONPatchOp string

var (
	JSONPatchOpAdd     JSONPatchOp = "add"
	JSONPatchOpReplace JSONPatchOp = "replace"
)

func applyJSONPatchFromNamespacedSecretRef(s SecretGetter, raw []byte, path string, namespace string, secretName string, key string) ([]byte, error) {
	secret, err := s.GetSecret(namespace, secretName)
	if err != nil {
		return nil, err
	}
	secretVal, ok := secret.Data[key]
	if !ok {
		return nil,
			fmt.Errorf("no key '%v' in secret '%v/%v'",
				key, namespace, secretName)
	}

	// JSON patch (RFC6902) specifies the behavior of applying "add" on root,
	// but because the jsonpatch package could not do "add" on root path (path=""),
	// we have to use "replace" op on root to set the entire content of document if patch is on root path.
	// https://github.com/evanphx/json-patch/issues/188
	op := JSONPatchOpAdd
	if path == "" {
		op = JSONPatchOpReplace
	}

	rawPatch := fmt.Sprintf(rawPatchPattern, op, path, string(secretVal))
	p, err := jsonpatch.DecodePatch([]byte(rawPatch))
	if err != nil {
		return nil, err
	}

	// Set EnsurePathExistsOnAdd to true for adding to subpaths to a non-existing path, e.g:
	// Apply {"op":"add","path":"/add/headers","value":[{"h1":"v1"},{"h2":"v2"}]} on `{}`.
	opts := jsonpatch.NewApplyOptions()
	opts.EnsurePathExistsOnAdd = true
	raw, err = p.ApplyWithOptions(raw, opts)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

// RawConfigurationWithPatchesToConfiguration converts config and add patches from configPatches of KongPlugin.
func RawConfigurationWithPatchesToConfiguration(
	s SecretGetter, namespace string,
	rawConfig apiextensionsv1.JSON,
	patches []kongv1.ConfigPatch,
) (kong.Configuration, error) {
	raw := rawConfig.Raw
	if raw == nil {
		// In case the config is empty, we need to initialize it to an empty
		// JSON object so that the patches can be applied.
		raw = []byte("{}")
	}

	// apply patches
	for _, patch := range patches {
		var err error
		raw, err = applyJSONPatchFromNamespacedSecretRef(
			s,
			raw,
			patch.Path,
			namespace,
			patch.ValueFrom.SecretValue.Secret,
			patch.ValueFrom.SecretValue.Key,
		)
		if err != nil {
			return kong.Configuration{}, err
		}
	}
	return RawConfigToConfiguration(raw)
}

// RawConfigurationWithNamespacedPatchesToConfiguration converts config and add patches from configPatches of KongClusterPlugin.
func RawConfigurationWithNamespacedPatchesToConfiguration(
	s SecretGetter,
	rawConfig apiextensionsv1.JSON,
	patches []kongv1.NamespacedConfigPatch,
) (kong.Configuration, error) {
	raw := rawConfig.Raw

	if raw == nil {
		// In case the config is empty, we need to initialize it to an empty
		// JSON object so that the patches can be applied.
		raw = []byte("{}")
	}
	for _, patch := range patches {
		var err error
		raw, err = applyJSONPatchFromNamespacedSecretRef(
			s,
			raw,
			patch.Path,
			patch.ValueFrom.SecretValue.Namespace,
			patch.ValueFrom.SecretValue.Secret,
			patch.ValueFrom.SecretValue.Key,
		)
		if err != nil {
			return kong.Configuration{}, err
		}
	}
	return RawConfigToConfiguration(raw)
}

// NamespacedSecretToConfiguration fetches specified value from given namespace, secret and key,
// then parse the value to Kong plugin configurations.
// Exported primarily to be used in admission validators.
func NamespacedSecretToConfiguration(
	s SecretGetter,
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

// SecretToConfiguration fetches specified value from secret and key in the namespace,
// then parse the value to Kong plugin configurations.
// Exported primarily to be used in admission validators.
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
