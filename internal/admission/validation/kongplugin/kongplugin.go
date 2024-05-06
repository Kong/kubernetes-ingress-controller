package kongplugin

import (
	"context"
	"fmt"
	"strings"

	"github.com/samber/lo"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
)

// clusterObjectsLister is an interface controller-runtime client.Client.
type clusterObjectsLister interface {
	List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error
}

type k8sObjectWithAnnotations interface {
	GetAnnotations() map[string]string
	GetNamespace() string
}

// ValidatePluginUniquenessPerObject validates that only one plugin of a particular type is attached to an object.
// It can only check when instance of particular plugin has been already created (KongPlugin has been applied before and saved by Kubernetes).
// Otherwise it will not be able to check and will return nil.
func ValidatePluginUniquenessPerObject(ctx context.Context, objectsLister clusterObjectsLister, obj k8sObjectWithAnnotations) error {
	attachedPlugins := annotations.ExtractKongPluginsFromAnnotations(obj.GetAnnotations())
	if len(attachedPlugins) == 0 {
		return nil
	}

	// Getting plugins directly from the cluster instead of operator's cache to make it more real-time.
	pluginList := &kongv1.KongPluginList{}
	if err := objectsLister.List(ctx, pluginList, &client.ListOptions{Namespace: obj.GetNamespace()}); err != nil {
		// If we cannot list plugins, we cannot check for uniqueness, hence assume it is valid.
		return nil //nolint:nilerr
	}

	pluginsForType := make(map[string][]string)
	for _, plugin := range pluginList.Items {
		if lo.Contains(attachedPlugins, plugin.Name) {
			// How many plugins of particular type e.g. "rate-limiting" are attached and what are those names.
			pluginsForType[plugin.PluginName] = append(pluginsForType[plugin.PluginName], plugin.Name)
		}
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
