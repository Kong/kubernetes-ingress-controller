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
	"testing"
)

func TestExtractKongPluginAnnotations(t *testing.T) {
	data := map[string]string{
		"rate-limiting.plugin.konghq.com":      "v1",
		"key-authentication.plugin.konghq.com": "v2",
	}

	ka := ExtractKongPluginAnnotations(data)
	if len(ka) != 2 {
		t.Errorf("expected two keys but %v returned", len(ka))
	}

	if _, ok := ka["rate-limiting"]; !ok {
		t.Errorf("expected a rate limiting plugin but none returned")
	}
}

func TestExtractKongPluginsFromAnnotations(t *testing.T) {
	data := map[string]string{
		"plugins.konghq.com": "kp-rl, kp-cors",
	}

	ka := ExtractKongPluginsFromAnnotations(data)
	if len(ka) != 2 {
		t.Errorf("expected two keys but %v returned", len(ka))
	}
	if ka[0] != "kp-rl" {
		t.Errorf("expected first element to be 'kp-rl'")
	}
}

func TestExtractConfigurationName(t *testing.T) {
	data := map[string]string{
		"configuration.konghq.com": "demo",
	}

	cn := ExtractConfigurationName(data)
	if cn != "demo" {
		t.Errorf("expected demo as configuration name but got %v", cn)
	}
}
