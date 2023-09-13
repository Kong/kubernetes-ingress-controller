package v1

import (
	"encoding/json"

	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

// the need to support both array and single value fields makes this much more awkward than I'd hoped.
// I don't think there's a good way to design the type around supporting that and having ValueNested support either.
// we need to be able to distinguish between arrays of simple values and arrays of objects:
//
// {
//     "foo": [
//         {
//             "fooSubA": "valueA1",
//             "fooSubB": "valueB1",
//         },
//         {
//             "fooSubA": "valueA2",
//             "fooSubB": "valueB2",
//         }
//     ],
//     "bar": [
//         "stringA",
//         "stringB"
//     ]
// }
//
// we could maybe use a single []*ArbitraryObj field with the implicit rule that objects with no names convert to
// arrays of simple values, but that's probably more confusing than it's worth. references don't play nice with
// arrays either, since you can't selectively make only some items in an array secret. however, since they're all
// the same type, I wouldn't expect that to be an actual use case--if one item in an array of same-type values is
// worth protecting, the rest should be as well.

// This variant uses apiextensions.JSON for the entire CRD field
// type ArbitraryObj struct {
// 	Name             string             `json:"name"`
// 	Value            interface{}        `json:"value,omitempty"`
// 	ValueArray       []interface{}      `json:"valueArray,omitempty"`
// 	ValueFrom        *ArbitraryObjSource   `json:"valueFrom,omitempty"`
// 	ValueFromArray   []*ArbitraryObjSource `json:"valueFromArray,omitempty"`
// 	ValueNested      *ArbitraryObj         `json:"valueNested,omitempty"`
// 	ValueNestedArray []*ArbitraryObj       `json:"valueNestedArray,omitempty"`
// }

// This variant can be included in a CRD field

// TODO we want interface{} for Value and ValueArray. These can be any JSON primitive type, but schema validation
// doesn't like that: make manifests will return "invalid field type: interface{}".
// unsure if we can handle it with generics or some other approach. If not we may be stuck using apiextensionsv1.JSON
// in the CRD and not validating simple fields. For now, it's just a string for the PoC.

// ArbitraryObj represents a JSON blob whose fields are represented with literal values, references to external store,
// nested ArbitraryObj, or arrays of any of the above.
type ArbitraryObj struct {
	Name             string                  `json:"name"`
	Value            string                  `json:"value,omitempty"`
	ValueArray       []string                `json:"valueArray,omitempty"`
	ValueFrom        *ArbitraryObjSource     `json:"valueFrom,omitempty"`
	ValueFromArray   []*ArbitraryObjSource   `json:"valueFromArray,omitempty"`
	ValueNested      *apiextensionsv1.JSON   `json:"valueNested,omitempty"`
	ValueNestedArray []*apiextensionsv1.JSON `json:"valueNestedArray,omitempty"`
}

// ArbitraryObjSource is a Secret key reference with an optional namespace.
type ArbitraryObjSource struct {
	Namespace    string                    `json:"namespace,omitempty"`
	SecretKeyRef *corev1.SecretKeySelector `json:"secretKeyRef,omitempty"`
}

func ArbitraryObjsToJSONMap(sources []ArbitraryObj) (map[string]interface{}, error) {
	jsonMap := make(map[string]interface{})
	for _, src := range sources {
		subMap, err := ArbitraryObjToJSONMap(src)
		if err != nil {
			return map[string]interface{}{}, err
		}
		jsonMap[src.Name] = subMap
	}
	return jsonMap, nil
}

func ArbitraryObjToJSONMap(source ArbitraryObj) (map[string]interface{}, error) {
	jsonMap := make(map[string]interface{})
	// TODO properly ensure that there's only one source and error if not, rather than choose the first available
	// TODO make this not a string
	if source.Value != "" {
		jsonMap[source.Name] = source.Value
	} else if len(source.ValueArray) != 0 {
		jsonMap[source.Name] = source.ValueArray
	} else if source.ValueFrom != nil {
		// TODO actually retrieve Secret contents. this will require a client
		jsonMap[source.Name] = "TODO"
	} else if len(source.ValueFromArray) != 0 {
		array := make([]interface{}, len(source.ValueFromArray))
		// TODO ditto source.ValueFrom
		for i := range source.ValueFromArray {
			array[i] = "TODO"
		}
		jsonMap[source.Name] = array
	} else if source.ValueNested != nil {
		a := ArbitraryObj{}
		err := json.Unmarshal(source.ValueNested.Raw, &a)
		if err != nil {
			return map[string]interface{}{}, err
		}
		subMap, err := ArbitraryObjToJSONMap(a)
		if err != nil {
			return map[string]interface{}{}, err
		}
		jsonMap[source.Name] = subMap
	} else if len(source.ValueNestedArray) != 0 {
		array := make([]ArbitraryObj, len(source.ValueNestedArray))
		for i := range source.ValueNestedArray {
			a := ArbitraryObj{}
			err := json.Unmarshal(source.ValueNested.Raw, &a)
			if err != nil {
				return map[string]interface{}{}, err
			}
			array[i] = a
		}
		subMap, err := ArbitraryObjsToJSONMap(array)
		if err != nil {
			return map[string]interface{}{}, err
		}
		jsonMap[source.Name] = subMap
	}
	return jsonMap, nil
}
