/*
Copyright 2018 Kong Inc

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
	"testing"

	kongadminv1 "github.com/kong/kubernetes-ingress-controller/internal/apis/admin/v1"
)

func TestPluginDeepEqual(t *testing.T) {
	var equal bool

	equal = pluginDeepEqual(map[string]interface{}{}, &kongadminv1.Plugin{Config: map[string]interface{}{}})
	if !equal {
		t.Errorf("Comparing empty maps failed")
	}

	equal = pluginDeepEqual(map[string]interface{}{
		"key1": "vaule1",
		"key2": "value2",
		"key3": "value3",
	}, &kongadminv1.Plugin{Config: map[string]interface{}{
		"key1": "vaule1",
		"key2": "value2",
		"key3": "value3",
	}})
	if !equal {
		t.Errorf("Comparing maps with same keys and same order failed")
	}

	equal = pluginDeepEqual(map[string]interface{}{
		"key1": "vaule1",
		"key2": "value2",
		"key3": "value3",
	}, &kongadminv1.Plugin{Config: map[string]interface{}{
		"key2": "value2",
		"key3": "value3",
		"key1": "vaule1",
	}})
	if !equal {
		t.Errorf("Comparing maps with same keys and different order failed")
	}

	equal = pluginDeepEqual(map[string]interface{}{
		"key1": map[string]string{
			"nestedkey1": "nestedvalue1",
		},
		"key2": "value2",
		"key3": "value3",
	}, &kongadminv1.Plugin{Config: map[string]interface{}{
		"key2": "value2",
		"key3": "value3",
		"key1": map[string]string{
			"nestedkey1": "nestedvalue1",
		},
	}})
	if !equal {
		t.Errorf("Comparing maps with nested map failed")
	}

	equal = pluginDeepEqual(map[string]interface{}{
		"key1": map[string]string{},
		"key2": "value2",
		"key3": "value3",
	}, &kongadminv1.Plugin{Config: map[string]interface{}{
		"key2": "value2",
		"key3": "value3",
		"key1": map[string]string{},
	}})
	if !equal {
		t.Errorf("Comparing maps with empty map failed")
	}

	equal = pluginDeepEqual(map[string]interface{}{
		"key1": [3]string{
			"arr1", "arr2", "arr3",
		},
		"key2": "value2",
		"key3": "value3",
	}, &kongadminv1.Plugin{Config: map[string]interface{}{
		"key2": "value2",
		"key3": "value3",
		"key1": [3]string{
			"arr1", "arr2", "arr3",
		},
	}})
	if !equal {
		t.Errorf("Comparing maps with nested array failed")
	}

	equal = pluginDeepEqual(map[string]interface{}{
		"key1": [3]string{
			"arr1", "arr2", "arr3",
		},
		"key2": "value2",
		"key3": "value3",
	}, &kongadminv1.Plugin{Config: map[string]interface{}{
		"key2": "value2",
		"key3": "value3",
		"key1": [3]string{
			"arr2", "arr3", "arr1",
		},
	}})
	if equal {
		t.Errorf("Comparing maps with nested array with different order failed")
	}

	equal = pluginDeepEqual(map[string]interface{}{
		"key1": []string{},
		"key2": "value2",
		"key3": "value3",
	}, &kongadminv1.Plugin{Config: map[string]interface{}{
		"key2": "value2",
		"key3": "value3",
		"key1": []string{},
	}})
	if !equal {
		t.Errorf("Comparing maps with empty array failed")
	}

	equal = pluginDeepEqual(map[string]interface{}{
		"key1": map[string]string{
			"nestedkey1": "nestedvalue1",
		},
		"key2": "value2",
		"key3": "value3",
	}, &kongadminv1.Plugin{Config: map[string]interface{}{
		"key2": "value2",
		"key1": map[string]string{
			"nestedkey1": "nestedvalue1",
		},
	}})
	if equal {
		t.Errorf("Comparing maps with missing keys failed")
	}

	equal = pluginDeepEqual(map[string]interface{}{
		"key1": map[string]string{
			"nestedkey1": "nestedvalue1",
		},
		"key2": "value2",
		"key3": "value3",
	}, &kongadminv1.Plugin{Config: map[string]interface{}{
		"key3": "value3",
		"key1": map[string]string{
			"nestedkey1": "nestedvalue1",
		},
		"key4": "value4",
	}})
	if equal {
		t.Errorf("Comparing maps with unmatched keys failed")
	}
}
