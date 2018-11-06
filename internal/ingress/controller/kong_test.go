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

	"github.com/hbagdi/go-kong/kong"
)

func TestPluginDeepEqual(t *testing.T) {
	var equal bool

	equal = pluginDeepEqual(map[string]interface{}{}, &kong.Plugin{Config: map[string]interface{}{}})
	if !equal {
		t.Errorf("Comparing empty maps failed")
	}

	equal = pluginDeepEqual(map[string]interface{}{
		"key1": "vaule1",
		"key2": "value2",
		"key3": "value3",
	}, &kong.Plugin{Config: map[string]interface{}{
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
	}, &kong.Plugin{Config: map[string]interface{}{
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
	}, &kong.Plugin{Config: map[string]interface{}{
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
		"key1": 1,
		"key2": 2,
		"key3": 8,
	}, &kong.Plugin{Config: map[string]interface{}{
		"key1": 1.0,
		"key2": 2.0,
		"key3": 8,
	}})
	if !equal {
		t.Errorf("Comparing maps with numeric values in different type failed")
	}

	equal = pluginDeepEqual(map[string]interface{}{
		"key1": map[string]string{},
		"key2": "value2",
		"key3": "value3",
	}, &kong.Plugin{Config: map[string]interface{}{
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
	}, &kong.Plugin{Config: map[string]interface{}{
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
	}, &kong.Plugin{Config: map[string]interface{}{
		"key2": "value2",
		"key3": "value3",
		"key1": [3]string{
			"arr2", "arr3", "arr1",
		},
	}})
	if !equal {
		t.Errorf("Comparing maps with nested string array with different order failed")
	}

	equal = pluginDeepEqual(map[string]interface{}{
		"key1": [3]int{
			1, 2, 3,
		},
		"key2": "value2",
		"key3": "value3",
	}, &kong.Plugin{Config: map[string]interface{}{
		"key2": "value2",
		"key3": "value3",
		"key1": [3]float64{
			1.0, 2.0, 3.0,
		},
	}})
	if !equal {
		t.Errorf("Comparing maps with nested numeric value array failed")
	}

	equal = pluginDeepEqual(map[string]interface{}{
		"key1": []string{},
		"key2": "value2",
		"key3": "value3",
	}, &kong.Plugin{Config: map[string]interface{}{
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
	}, &kong.Plugin{Config: map[string]interface{}{
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
	}, &kong.Plugin{Config: map[string]interface{}{
		"key3": "value3",
		"key1": map[string]string{
			"nestedkey1": "nestedvalue1",
		},
		"key4": "value4",
	}})
	if equal {
		t.Errorf("Comparing maps with unmatched keys failed")
	}

	equal = pluginDeepEqual(map[string]interface{}{
		"key2": "value2",
		"key1": map[string]string{
			"nestedkey1": "nestedvalue1",
		},
	}, &kong.Plugin{Config: map[string]interface{}{
		"key2": "value2",
		"key1": map[string]string{
			"nestedkey1": "nestedvalue1",
		},
		"default3": "value3",
	}})
	if !equal {
		t.Errorf("Comparing maps with default configs failed")
	}

	equal = pluginDeepEqual(map[string]interface{}{
		"key2": "value2",
		"key1": map[string]string{
			"nestedkey1": "nestedvalue1",
		},
	}, &kong.Plugin{Config: map[string]interface{}{
		"key2": "value2",
		"key1": map[string]string{
			"nestedkey1":     "nestedvalue1",
			"defaultnested2": "nestedvalue2",
		},
		"default3": "value3",
	}})
	if !equal {
		t.Errorf("Comparing maps with nested default configs failed")
	}
}
