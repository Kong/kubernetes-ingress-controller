/*
Copyright 2017 The Kubernetes Authors.

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

package annotations

import (
	"fmt"
	"strings"
)

// pluginAnnotationSuffix sufix of kong annotations to configure plugins
const pluginAnnotationSuffix = "plugin.konghq.com"

const pluginListAnnotation = "plugins.konghq.com"

const configurationAnnotation = "configuration.konghq.com"

// ExtractKongPluginAnnotations extracts information about kong plugins
// configured using plugin.konghq.com annotation.
// DEPRECATED, please use ExtractKongPluginsFromAnnotations instead.
func ExtractKongPluginAnnotations(anns map[string]string) map[string][]string {
	ka := make(map[string][]string, 0)
	for k, v := range anns {
		if strings.HasSuffix(k, pluginAnnotationSuffix) {
			name := strings.TrimSuffix(k, fmt.Sprintf(".%v", pluginAnnotationSuffix))
			var values []string
			for _, line := range strings.Split(v, "\n") {
				s := strings.TrimSpace(strings.TrimPrefix(line, "-"))
				if s != "" {
					values = append(values, s)
				}
			}
			ka[name] = values
		}
	}

	return ka
}

// ExtractKongPluginsFromAnnotations extracts information about Kong
// Plugins configured using plugins.konghq.com annotation.
// This returns a list of KongPlugin resource names that should be applied.
func ExtractKongPluginsFromAnnotations(anns map[string]string) []string {
	var kongPluginCRs []string
	v, ok := anns[pluginListAnnotation]
	if !ok {
		return kongPluginCRs
	}
	for _, kongPlugin := range strings.Split(v, ",") {
		s := strings.TrimSpace(kongPlugin)
		if s != "" {
			kongPluginCRs = append(kongPluginCRs, s)
		}
	}
	return kongPluginCRs
}

// ExtractConfigurationName extracts the name of the KongIngress object that holds
// information about the configuration to use in Routes, Services and Upstreams
func ExtractConfigurationName(anns map[string]string) string {
	return anns[configurationAnnotation]
}
