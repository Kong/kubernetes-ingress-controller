package kongplugin

import (
	"fmt"
	"strings"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
)

type kongPluginGetter interface {
	GetKongPlugin(namespace, name string) (*kongv1.KongPlugin, error)
}

type k8sObjectWithAnnotations interface {
	GetAnnotations() map[string]string
	GetNamespace() string
}

// ValidatePluginUniquenessPerObject validates that only one plugin of a particular type is attached to an object.
// It can only check when instance of particular plugin has been already created (KongPlugin has been applied before).
// Otherwise it will not be able to check and will return nil.
func ValidatePluginUniquenessPerObject(pluginStore kongPluginGetter, obj k8sObjectWithAnnotations) error {
	pluginNames := annotations.ExtractKongPluginsFromAnnotations(obj.GetAnnotations())
	if len(pluginNames) == 0 {
		return nil
	}

	pluginsForType := make(map[string][]string)
	for _, pn := range pluginNames {
		plugin, err := pluginStore.GetKongPlugin(obj.GetNamespace(), pn)
		if err != nil {
			continue
		}
		// How many plugins of particular type e.g. "rate-limiting" are attached and what are those names.
		pluginsForType[plugin.PluginName] = append(pluginsForType[plugin.PluginName], plugin.Name)
	}
	// Check if the plugin is unique.
	for pluginType, pluginNames := range pluginsForType {
		if len(pluginNames) > 1 {
			return fmt.Errorf(
				"cannot attach multiple plugins: %s of the same type %s",
				strings.Join(pluginNames, ", "), pluginType,
			)
		}
	}
	return nil
}
